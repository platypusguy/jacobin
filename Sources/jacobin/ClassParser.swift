/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

import Foundation

class ClassParser {

    /// reads in a class file, parses it, and puts the values into the fields of the
    /// class that will be loaded into the classloader. Some verification performed
    /// - Parameters:
    ///   - name: name of the class to read in and parse
    ///   - klass: the Swift object that will hold the data extracted from the parsing
    /// - Throws:
    ///   - JVMerror.ClassFormatError - if the parser finds anything unexpected
    static func parseClassfile( name: String, klass: LoadedClass ) throws {

        klass.path = name
        let fileURL = URL( string: "file:" + name )!
        do {
            let data = try Data( contentsOf: fileURL, options: [.uncached] )
            klass.rawBytes = [UInt8]( data )
        } catch {
            log.log( msg: "Error reading file: \(name) Exiting", level: Logger.Level.SEVERE )
            throw JVMerror.ClassFormatError( msg: "" )
        }

        do {
            //check that the class file begins with the magic number 0xCAFEBABE
            if klass.rawBytes[0] != 0xCA || klass.rawBytes[1] != 0xFE ||
               klass.rawBytes[2] != 0xBA || klass.rawBytes[3] != 0xBE {
                throw JVMerror.ClassFormatError( msg: name )
            }

            //check that the file version is not above JDK 11 (that is, 55)
            let version = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: 6 )
            if version > 55 {
                log.log(
                        msg: "Error: this version of Jacobin supports only Java classes at or below Java 11. Exiting.",
                        level: Logger.Level.SEVERE )
                throw JVMerror.ClassFormatError( msg: "" )
            } else {
                klass.version = version;
                klass.status = classStatus.PRELIM_VERIFIED
            }

            // get the constant pool count
            let cpCount = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: 8 )
            if cpCount < 2 {
                throw JVMerror.ClassFormatError( msg: name + " constant pool count." )
            } else {
                klass.constantPoolCount = cpCount
                log.log( msg: "Class \(name) constant pool should have \(cpCount) entries",
                         level: Logger.Level.FINEST )
            }

            // load and verify the constant pool
            var location: Int = ConstantPool.load( klass: klass ) //location = index of last byte examined
            try ConstantPool.verify( klass: klass )
            ConstantPool.log( klass: klass )

            // load and verify the class access masks
            AccessFlags.readAccessFlags( klass: klass, location: location )
            AccessFlags.processClassAccessMask( klass: klass )
            AccessFlags.verify( klass: klass )
            AccessFlags.log( klass: klass )
            location += 2

            // get the pointer to this class name and extract the name
            ThisClassName.readName( klass: klass, location: location )
            ThisClassName.verify( klass: klass )
            ThisClassName.process( klass: klass )
            ThisClassName.log( klass: klass )
            location += 2

            // get the pointer to the superclass for this class and extract the name
            SuperClassName.readName( klass: klass, location: location )
            try SuperClassName.verify( klass: klass )
            SuperClassName.process( klass: klass )
            SuperClassName.log( klass: klass )
            location += 2

            // get the count of interfaces implemented by this class
            InterfaceCount.readInterfaceCount( klass: klass, location: location )
            InterfaceCount.log( klass: klass )
            location += 2

            // get the names of all the interfaces implemented by this class
            if klass.interfaceCount > 0 {
                for _ in 1...klass.interfaceCount {
                    let index = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: location+1 )
                    location += 2
                    Interface.process( klass: klass, index: index )
                }
                Interface.logAll( klass: klass )
            }

            // get the count of fields in this class, put it into klass.fieldCount
            FieldCount.readFieldCount( klass: klass, location: location )
            FieldCount.log( klass: klass )
            location += 2

            // ...and process the fields
            if klass.fieldCount > 0  {
                for _ in 1...klass.fieldCount {
                    let field = Field()
                    location = field.load( klass: klass, location: location )
                    klass.fields.append( field )
                }
            }

            // get the count of methods in this class
            MethodCount.readMethodCount( klass: klass, location: location )
            MethodCount.log( klass: klass )
            location += 2

            if klass.methodCount > 0 {
                for _ in 1...( klass.methodCount ) {
                    let mi = MethodInfo()
                    location = mi.read( klass: klass, location: location )
                    mi.log( klass: klass )
                    klass.methodInfo.append( mi.methodData )
                }
            }

            let attrCount = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: location+1 )
            location += 2

            // get any remaining attributes
            for _ in 1...attrCount  {
                let attributeID = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: location+1 )
                let attrName =
                    Utility.getUTF8stringFromConstantPoolIndex( klass: klass, index: attributeID )
                location += 2
                let attrSize = Utility.getIntfrom4Bytes( bytes: klass.rawBytes, index: location+1 )
                location += 4

                switch attrName {
                case "BootstrapMethods":
                    let bsm = BootstrapMethodsAttribute( name: attrName, length: attrSize )
                    try bsm.load( klass: klass, loc: location )
                    bsm.log( klass: klass )
                    location += attrSize

                case "SourceFile":
                    let sfa = SourceFileAttribute( name: attrName, length: attrSize )
                    sfa.load( klass: klass, loc: location )
                    klass.attributes.append( sfa )
                    sfa.log( className: klass.shortName )
                    location += attrSize

                default:
                    print( "Attribute: \(attrName) not processed in class \(klass.shortName)" )
                    location += attrSize
                }
            }
        }
        catch JVMerror.ClassFormatError( let msg ) {
            log.log( msg: "ClassFormatError: \( msg )", level: Logger.Level.SEVERE )
            throw JVMerror.ClassFormatError( msg: "" )
        }

        //TODO: and validate with the logic here: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.1
        //TODO: from 4.10:
        //Ensuring that final classes are not subclassed.
        //
        //Ensuring that final methods are not overridden (§5.4.5).
        //
        //Checking that every class (except Object) has a direct superclass.

    }
}
