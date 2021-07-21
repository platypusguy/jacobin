/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */
import Foundation

class ConstantPool {

    public enum RecType: Int {
        case invalid        = -1 // used for initialization and for dummy entries (viz. for longs, doubles)
        case UTF8           =  1
        case intConst       =  3
        case floatConst     =  4
        case longConst      =  5
        case doubleConst    =  6
        case classRef       =  7
        case string         =  8
        case field          =  9
        case method         = 10
        case interface      = 11
        case nameAndType    = 12
        case methodHandle   = 15
        case methodType     = 16
        case dynamic        = 17
        case invokeDynamic  = 18
        case module         = 19
        case package        = 20
    }

    // the constant pool of a class is a collection of individual entries that point to classes, methods, strings, etc.
    // This method parses through them and creates an array of parsed entries in the class being loaded. The entries in
    // the array inherit from cpEntryTemplate. Note that the first entry in all constant pools is non-existent, which I
    // believe was done to avoid off-by-one errors in lookups, but not sure. This is why the loop through entries begins
    // at 1, rather than 0.
    //
    // returns the byte number of the end of constant pool
    //
    // Refer to: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4-140
    static func load( klass: LoadedClass ) -> Int {
        var byteCounter = 9 //the number of bytes we're into the class file (zero-based)

        let dummyEntry = CpDummyEntry()
        klass.cp.append( dummyEntry ) // entry[0] is a dummy entry that's never used
        var i = 1
        while( i <= klass.constantPoolCount-1 ) {
            byteCounter += 1
            let cpeType = Int(klass.rawBytes[byteCounter])
            switch( RecType( rawValue: cpeType ) ) {
            case .UTF8: // UTF-8 string
                let length =
                        Utility.getInt16from2Bytes( msb: klass.rawBytes[byteCounter+1],
                                                    lsb: klass.rawBytes[byteCounter+2] )
                byteCounter += 2
                var buffer = [UInt8]()
                for n in 0...Int(length-1)  {
                    buffer.append( klass.rawBytes[byteCounter+n+1] )
                }
                let UTF8string = String( bytes: buffer, encoding: String.Encoding.utf8 ) ?? ""
                let UTF8entry = CpEntryUTF8( contents: UTF8string )
                klass.cp.append( UTF8entry )
                byteCounter += Int(length)
                print( "UTF-8 string: \( UTF8string )" )

            case .intConst: // integer constant
                let value =
                        Utility.getIntfrom4Bytes(bytes: klass.rawBytes, index: byteCounter+1 )
                let integerConstantEntry = CpIntegerConstant( value: value )
                klass.cp.append( integerConstantEntry )
                byteCounter += 4
                print( "Integer constant: \( value )" )

            case .floatConst: // floating-point constant (32 bits) Convert 4 bytes into a Float
                let tPointer = UnsafeMutablePointer<UInt8>.allocate( capacity:4 )
                let pointer = UnsafeRawPointer( tPointer )

                tPointer[0]=klass.rawBytes[byteCounter+1]
                tPointer[1]=klass.rawBytes[byteCounter+2]
                tPointer[2]=klass.rawBytes[byteCounter+3]
                tPointer[3]=klass.rawBytes[byteCounter+4]

                let value: Float = pointer.load( fromByteOffset: 00, as: Float.self )

                let floatConstantEntry = CpFloatConstant( value: value )
                klass.cp.append( floatConstantEntry )
                byteCounter += 4
                print( "Float constant: \( value )" )

            case .longConst: // long constant (fills two slots in the constant pool)
                let highBytes =
                        Utility.getIntfrom4Bytes( bytes: klass.rawBytes, index: byteCounter+1 )
                let lowBytes =
                        Utility.getIntfrom4Bytes( bytes: klass.rawBytes, index: byteCounter+5 )
                let longValue : Int64 = Int64(( highBytes << 32) + lowBytes )
                let longConstantEntry = CpLongConstant( value: longValue )
                klass.cp.append( longConstantEntry )
                // longs take up two slots in the constant pool, of which the second slot is
                // never accessed. So set up a dummy entry for that slot.
                klass.cp.append( CpDummyEntry() )
                klass.constantPoolCount -= 1 // decrease the total number of entries to create due to dummy
                byteCounter += 8
                print( "Long constant: \( longValue )" )

            case .doubleConst: // double floating-point constant (64 bits). Fills two slots in the constant pool
                let tPointer = UnsafeMutablePointer<UInt8>.allocate( capacity:8 )
                let pointer = UnsafeRawPointer( tPointer )

                tPointer[0]=klass.rawBytes[byteCounter+1]
                tPointer[1]=klass.rawBytes[byteCounter+2]
                tPointer[2]=klass.rawBytes[byteCounter+3]
                tPointer[3]=klass.rawBytes[byteCounter+4]
                tPointer[4]=klass.rawBytes[byteCounter+5]
                tPointer[5]=klass.rawBytes[byteCounter+6]
                tPointer[6]=klass.rawBytes[byteCounter+7]
                tPointer[7]=klass.rawBytes[byteCounter+8]

                let value: Double = pointer.load( fromByteOffset: 00, as: Double.self )

                let doubleConstantEntry = CpDoubleConstant( value: value )
                klass.cp.append( doubleConstantEntry )
                // now add the dummy entry, which happen with longs and doubles
                klass.cp.append( CpDummyEntry() )
                klass.constantPoolCount -= 1 // decrease the total number of entries to create due to dummy
                byteCounter += 8
                print( "Double constant: \( value )" )

            case .classRef: // class reference
                let classNameIndex =
                        Utility.getInt16from2Bytes( msb: klass.rawBytes[byteCounter+1],
                                                    lsb: klass.rawBytes[byteCounter+2] )
                let classNameRef = CpEntryClassRef( index: classNameIndex )
                klass.cp.append( classNameRef )
                byteCounter += 2
                print( "Class name reference: index: \( classNameIndex )" )

            case .string: // string reference
                let stringIndex =
                        Utility.getInt16from2Bytes( msb: klass.rawBytes[byteCounter+1],
                                                    lsb: klass.rawBytes[byteCounter+2] )
                let stringRef = CpEntryStringRef( index: stringIndex )
                klass.cp.append( stringRef )
                byteCounter += 2
                print( "String reference: string index: \(stringIndex) ")

            case  .field: // field reference
                let classIndex =
                        Utility.getInt16from2Bytes( msb: klass.rawBytes[byteCounter+1],
                                                    lsb: klass.rawBytes[byteCounter+2] )
                let nameAndTypeIndex =
                        Utility.getInt16from2Bytes( msb: klass.rawBytes[byteCounter+3],
                                                    lsb: klass.rawBytes[byteCounter+4] )
                byteCounter += 4
                let fieldRef = CpEntryFieldRef( classIndex: classIndex,
                                                nameAndTypeIndex: nameAndTypeIndex );
                klass.cp.append( fieldRef )
                print( "Field reference: class index: \(classIndex) nameAndTypeIndex: \(nameAndTypeIndex)" )

            case .method: // method reference
                let classIndex =
                        Utility.getInt16from2Bytes( msb: klass.rawBytes[byteCounter+1],
                                                    lsb: klass.rawBytes[byteCounter+2] )
                let nameAndTypeIndex =
                        Utility.getInt16from2Bytes( msb: klass.rawBytes[byteCounter+3],
                                                    lsb: klass.rawBytes[byteCounter+4] )
                byteCounter += 4
                let methodRef = CpEntryMethodRef( classIndex: classIndex,
                                                  nameAndTypeIndex: nameAndTypeIndex );
                klass.cp.append( methodRef )
                print( "Method reference: class index: \(classIndex) nameAndTypeIndex: \(nameAndTypeIndex)" )

            case .interface: // interface method reference
                let classIndex =
                        Utility.getInt16from2Bytes( msb: klass.rawBytes[byteCounter+1],
                                                    lsb: klass.rawBytes[byteCounter+2] )
                let nameAndTypeIndex =
                        Utility.getInt16from2Bytes( msb: klass.rawBytes[byteCounter+3],
                                                    lsb: klass.rawBytes[byteCounter+4] )
                byteCounter += 4
                let interfaceMethodRef = CpEntryInterfaceMethodRef( classIndex: classIndex,
                                                                    nameAndTypeIndex: nameAndTypeIndex );
                klass.cp.append( interfaceMethodRef )
                print( "Interface reference: class index: \(classIndex) nameAndTypeIndex: \(nameAndTypeIndex)" )

            case .nameAndType: // name and type info
                let nameIndex =
                        Utility.getInt16from2Bytes( msb: klass.rawBytes[byteCounter+1], lsb: klass.rawBytes[byteCounter+2] )
                let descriptorIndex =
                        Utility.getInt16from2Bytes( msb: klass.rawBytes[byteCounter+3], lsb: klass.rawBytes[byteCounter+4] )
                byteCounter += 4
                let nameAndType : CpNameAndType =
                        CpNameAndType( nameIdx: Int(nameIndex), descriptorIdx: Int(descriptorIndex) )
                klass.cp.append( nameAndType )
                print( "Name and type info: name index: \(nameIndex) descriptorIndex: \(descriptorIndex)")

            case .methodHandle: // method handle
                let methodKind  = klass.rawBytes[byteCounter+1]
                let methodIndex =
                        Utility.getInt16from2Bytes( msb: klass.rawBytes[byteCounter+2],
                                                    lsb: klass.rawBytes[byteCounter+3] )
                byteCounter += 3
                let methodHandle = CpEntryMethodHandle( kind: methodKind, index: methodIndex )
                klass.cp.append( methodHandle )
                print( "Method handle kind: \(methodKind) index: \(methodIndex)" )

            case .methodType: // method type https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.9
                let methodIndex =
                        Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: byteCounter+1 )
                byteCounter += 2
                let methodType = CpMethodType( index: methodIndex )
                klass.cp.append( methodType )
                print( "Method type: \(methodIndex)" )

            case .invokeDynamic: // invokedynamic https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.10
                let bootstrapIndex =
                        Utility.getInt16from2Bytes( msb: klass.rawBytes[byteCounter+1],
                                                    lsb: klass.rawBytes[byteCounter+2] )
                let nameAndTypeIndex =
                        Utility.getInt16from2Bytes( msb: klass.rawBytes[byteCounter+3],
                                                    lsb: klass.rawBytes[byteCounter+4] )
                byteCounter += 4
                let invokedynamic = CpInvokedynamic( bootstrap: bootstrapIndex, nameAndType: nameAndTypeIndex)
                klass.cp.append( invokedynamic )
                print( "Invokedynamic boostrap idx: \(bootstrapIndex), name and type: \(nameAndTypeIndex)" )

            case .dynamic: // dynamic, which points to info re dynamically computed constants
                let bootstrapIndex  = Utility.getIntFrom2Bytes(bytes: klass.rawBytes, index: byteCounter+1 )
                let nameAndDefIndex = Utility.getIntFrom2Bytes(bytes: klass.rawBytes, index: byteCounter+3 )
                byteCounter += 4

                let nameAndType : CpNameAndType = klass.cp[nameAndDefIndex] as! CpNameAndType
                let dynamic = CpDynamic( bootstrap: bootstrapIndex,
                                         name: nameAndType.nameIndex,
                                         desc: nameAndType.descriptorIndex )
                klass.cp.append( dynamic )
                print( "Dynamic entry found" )

            case .module: // module (valid for Java 9+)
                let nameIndex = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: byteCounter+1 )
                byteCounter += 2
                guard nameIndex > 0 && nameIndex < klass.constantPoolCount else {
                    jacobin.log.log( msg: "Class \(klass.shortName), invalid module name",
                                     level: Logger.Level.WARNING )
                    klass.cp.append( CpModuleName( moduleName: "" ))
                    continue
                }
                let moduleName =
                    Utility.getUTF8stringFromConstantPoolIndex( klass: klass, index: nameIndex )
                klass.cp.append( CpModuleName( moduleName: moduleName ))

            case .package: // package name (valid for Java 9+)
                let nameIndex = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: byteCounter+1 )
                byteCounter += 2
                guard nameIndex > 0 && nameIndex < klass.constantPoolCount else {
                    jacobin.log.log( msg: "Class \(klass.shortName), invalid package name",
                                     level: Logger.Level.WARNING )
                    klass.cp.append( CpPackageName( packageName: "" ))
                    continue
                }
                let packageName =
                    Utility.getUTF8stringFromConstantPoolIndex( klass: klass, index: nameIndex )
                klass.cp.append( CpPackageName( packageName: packageName ))

            default:
                print( "** Unhandled constant pool entry found: \(cpeType) at byte \(byteCounter)" )
                break
            }
        i += 1
        }
        return byteCounter
    }

    // make sure all the pointers point to the correct items and that values are within the right range
    static func verify( klass: LoadedClass ) throws {
        let className = klass.path
        for n in 1..<klass.constantPoolCount {
            switch ( klass.cp[n].type ) {
            case .invalid: // dummy entries created by the JVM ( for 0th element, longs, doubles, etc.)
                continue
            case .UTF8: //UTF8 string
                let currTemp: CpEntryTemplate = klass.cp[n]
                let currEntry = currTemp as! CpEntryUTF8
                let UTF8string = currEntry.string
                if UTF8string.contains( Character( UnicodeScalar( 0x00 ) ) ) || //Ox00 and OxF0 through 0xFF are disallowed
                           UTF8string.contains( Character( UnicodeScalar( 0xF0 ) ) ) ||
                           UTF8string.contains( Character( UnicodeScalar( 0xF1 ) ) ) ||
                           UTF8string.contains( Character( UnicodeScalar( 0xF2 ) ) ) ||
                           UTF8string.contains( Character( UnicodeScalar( 0xF3 ) ) ) ||
                           UTF8string.contains( Character( UnicodeScalar( 0xF4 ) ) ) ||
                           UTF8string.contains( Character( UnicodeScalar( 0xF5 ) ) ) ||
                           UTF8string.contains( Character( UnicodeScalar( 0xF6 ) ) ) ||
                           UTF8string.contains( Character( UnicodeScalar( 0xF7 ) ) ) ||
                           UTF8string.contains( Character( UnicodeScalar( 0xF8 ) ) ) ||
                           UTF8string.contains( Character( UnicodeScalar( 0xF9 ) ) ) ||
                           UTF8string.contains( Character( UnicodeScalar( 0xFA ) ) ) ||
                           UTF8string.contains( Character( UnicodeScalar( 0xFB ) ) ) ||
                           UTF8string.contains( Character( UnicodeScalar( 0xFC ) ) ) ||
                           UTF8string.contains( Character( UnicodeScalar( 0xFD ) ) ) ||
                           UTF8string.contains( Character( UnicodeScalar( 0xFE ) ) ) ||
                           UTF8string.contains( Character( UnicodeScalar( 0xFF ) ) ) {
                    jacobin.log.log( msg: "Error validating constant pool in class \(className) Exiting.",
                                     level: Logger.Level.SEVERE )
                    throw JVMerror.ClassFormatError( msg: "" )
                }

            case .classRef: // class reference must point to UTF8 string
                let currTemp: CpEntryTemplate = klass.cp[n]
                let currEntry: CpEntryClassRef = currTemp as! CpEntryClassRef
                let index = currEntry.classNameIndex
                let pointedToEntry = klass.cp[index]
                if pointedToEntry.type != .UTF8 {
                    jacobin.log.log( msg: "Error validating constant pool in class \(className) Exiting.",
                                     level: Logger.Level.SEVERE )
                }

            case .string: // constant string must point to a UTF8 string
                let currTemp: CpEntryTemplate = klass.cp[n]
                let currEntry: CpEntryStringRef = currTemp as! CpEntryStringRef
                let index = currEntry.stringIndex
                let pointedToEntry = klass.cp[index]
                if pointedToEntry.type != .UTF8 {
                    jacobin.log.log( msg: "Error validating constant pool in class \(className) Exiting.",
                                     level: Logger.Level.SEVERE )
                }

            case .field: // field ref must point to a class and to a nameAndType
                let currTemp: CpEntryTemplate = klass.cp[n]
                let currEntry: CpEntryFieldRef = currTemp as! CpEntryFieldRef
                let classIndex = currEntry.classIndex
                let nameAndTypeIndex = currEntry.nameAndTypeIndex
                let pointedToEntry = klass.cp[classIndex]
                let pointedToField = klass.cp[nameAndTypeIndex]
                if pointedToEntry.type != .classRef || pointedToField.type != .nameAndType {
                    jacobin.log.log( msg: "Error validating constant pool in class \(className) Exiting.",
                                     level: Logger.Level.SEVERE )
                }

            case .method: // method reference
                let currTemp: CpEntryTemplate = klass.cp[n]
                let currEntry: CpEntryMethodRef = currTemp as! CpEntryMethodRef
                let classIndex = currEntry.classIndex
                var pointedToEntry = klass.cp[classIndex]
                if pointedToEntry.type != .classRef { //method ref must point to a class reference
                    jacobin.log.log( msg: "Error validating constant pool in class \(className) Exiting.",
                                     level: Logger.Level.SEVERE )
                }
                let nameIndex = currEntry.nameAndTypeIndex
                pointedToEntry = klass.cp[nameIndex]
                if pointedToEntry.type != .nameAndType { //method ref name index must point to a name and type entry
                    jacobin.log.log( msg: "Error validating constant pool in class \(className) Exiting.",
                                     level: Logger.Level.SEVERE )
                } else { // make sure the name and type entry's name is pointing to a correctly named method
                    let nameAndTypEntry: CpNameAndType = pointedToEntry as! CpNameAndType
                    let namePointer = nameAndTypEntry.nameIndex
                    pointedToEntry = klass.cp[namePointer]
                    if pointedToEntry.type != .UTF8 { //the name must be a UTF8 string
                        jacobin.log.log( msg: "Error validating constant pool in class \(className) Exiting.",
                                         level: Logger.Level.SEVERE )
                    } else { // if the name begins with a < it must only be <init>
                        let utf8Entry = pointedToEntry as! CpEntryUTF8
                        let methodName = utf8Entry.string
                        if methodName.starts( with: "<" ) && !( methodName.starts( with: "<init>" ) ) {
                            jacobin.log.log( msg: "Error validating constant pool in class \(className) Exiting.",
                                             level: Logger.Level.SEVERE )
                        }
                    }
                }

            case .nameAndType: // name and type info
                let currTemp: CpEntryTemplate = klass.cp[n]
                let nameAndTypEntry: CpNameAndType = currTemp as! CpNameAndType
                let namePointer = nameAndTypEntry.nameIndex
                var cpEntry = klass.cp[namePointer]
                if cpEntry.type != .UTF8 { //the name pointer must point to a UTF8 string
                    jacobin.log.log( msg: "Error validating constant pool in class \(className) Exiting.",
                                     level: Logger.Level.SEVERE )
                }
                let typePointer = nameAndTypEntry.descriptorIndex
                cpEntry = klass.cp[typePointer]
                if cpEntry.type != .UTF8 { //the name pointer must point to a UTF8 string
                    jacobin.log.log( msg: "Error validating constant pool in class \(className) Exiting.",
                                     level: Logger.Level.SEVERE )
                }

            case .methodHandle: // method handle
                // see https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.8
                if klass.version < 51 { // methodHandle requires Java 7 at minimum
                    jacobin.log.log( msg: "Class\(klass.shortName) has invalid instruction version",
                        level: Logger.Level.SEVERE )
                    throw JVMerror.ClassFormatError( msg: "\(#file), func: \(#function) line: \(#line)" )
                }
                let cpMethHandle =  klass.cp[n] as! CpEntryMethodHandle
                switch( cpMethHandle.referenceKind ) {
                    case 1, 2, 3, 4:
                        if klass.cp[ cpMethHandle.referenceIndex ].type != .field {
                            jacobin.log.log( msg: "Class\(klass.shortName) has a method handle w/ invalid reference index " +
                                                  "for reference kind \( cpMethHandle.referenceKind)",
                                level: Logger.Level.SEVERE )
                            throw JVMerror.ClassFormatError( msg: "\(#file), func: \(#function) line: \(#line)" )
                        }
                    case 5, 8:
                        if klass.cp[ cpMethHandle.referenceIndex ].type != .method {
                            jacobin.log.log( msg: "Class\(klass.shortName) has a method handle w/ invalid reference index " +
                                                  "for reference kind \(cpMethHandle.referenceKind)",
                                level: Logger.Level.SEVERE )
                            throw JVMerror.ClassFormatError( msg: "\(#file), func: \(#function) line: \(#line)" )
                        }
                    case 6, 7:
                        if klass.version < 52 {
                            if klass.cp[ cpMethHandle.referenceIndex ].type != .method {
                                jacobin.log.log( msg: "Class\(klass.shortName) has a method handle w/ invalid reference index " +
                                                      "for reference kind \(cpMethHandle.referenceKind)",
                                    level: Logger.Level.SEVERE )
                                throw JVMerror.ClassFormatError( msg: "\(#file), func: \(#function) line: \(#line)" )
                            }
                        }
                        else {
                            if klass.cp[ cpMethHandle.referenceIndex ].type != .method ||
                               klass.cp[ cpMethHandle.referenceIndex ].type != .interface {
                                jacobin.log.log( msg: "Class\(klass.shortName) has a method handle w/ invalid reference index " +
                                                      "for reference kind \(cpMethHandle.referenceKind)",
                                    level: Logger.Level.SEVERE )
                                throw JVMerror.ClassFormatError( msg: "\(#file), func: \(#function) line: \(#line)" )
                            }
                        }
                    case 9:
                        if klass.cp[ cpMethHandle.referenceIndex ].type != .interface {
                            jacobin.log.log( msg: "Class\(klass.shortName) has a method handle w/ invalid reference index " +
                                                  "for reference kind \(cpMethHandle.referenceKind)",
                                level: Logger.Level.SEVERE )
                            throw JVMerror.ClassFormatError( msg: "\(#file), func: \(#function) line: \(#line)" )
                        }
                    default:
                        jacobin.log.log( msg: "Class\(klass.shortName) has a method handle w/ invalid reference kind: " +
                                              "\(cpMethHandle.referenceKind)",
                            level: Logger.Level.SEVERE )
                        throw JVMerror.ClassFormatError( msg: "\(#file), func: \(#function) line: \(#line)" )
                }

            case .methodType: // method type
                // see https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.9
                if klass.version < 51 {
                    jacobin.log.log( msg: "Error invalid methodType entry in \(className) Exiting.",
                                     level: Logger.Level.SEVERE )
                    throw JVMerror.ClassFormatError( msg: "" )
                }
                let cpMethType = klass.cp[n] as! CpMethodType
                if klass.cp[cpMethType.constantMethodIndex].type != .UTF8 {
                    jacobin.log.log( msg: "Class\(klass.shortName) has method type that does not point to UTF8 record",
                        level: Logger.Level.SEVERE )
                    throw JVMerror.ClassFormatError( msg: "\(#file), func: \(#function) line: \(#line)" )
                }

            case .dynamic: // constant dynamic
                if klass.version < 55 { // requires Java 11 at minimum
                    jacobin.log.log( msg: "Class\(klass.shortName) has invalid dynamic constant pool entry (version #)",
                                     level: Logger.Level.SEVERE )
                    throw JVMerror.ClassFormatError( msg: "\(#file), func: \(#function) line: \(#line)" )
                }

                case .invokeDynamic: // invokedynamic bytecode
                if klass.version < 51 {
                    jacobin.log.log( msg: "Class\(klass.shortName) has invalid invokedynamic constant pool entry (version #)",
                                     level: Logger.Level.SEVERE )
                    throw JVMerror.ClassFormatError( msg: "\(#file), func: \(#function) line: \(#line)" )
                }

            case .module: // JPMS module
                if klass.version < 53 { // requires Java 9 at minimum
                    jacobin.log.log( msg: "Class\(klass.shortName) has invalid module constant pool entry (version #)",
                        level: Logger.Level.SEVERE )
                    throw JVMerror.ClassFormatError( msg: "\(#file), func: \(#function) line: \(#line)" )
                }

            case .package: // JPMS package
                if klass.version < 53 { // requires Java 9 at minimum
                    jacobin.log.log( msg: "Class\(klass.shortName) has invalid package constant pool entry (version #)",
                        level: Logger.Level.SEVERE )
                    throw JVMerror.ClassFormatError( msg: "\(#file), func: \(#function) line: \(#line)" )
                }
            default: continue // for the nonce, eventually should be an error.
            }
        }
    }

    // a quick statistical point if we're at the highest level of verbosity
    static func log( klass: LoadedClass ) {
        jacobin.log.log(msg: "Class: \( klass.path ) - constant pool has: \( klass.cp.count ) entries",
                        level: Logger.Level.FINEST )
    }
}
