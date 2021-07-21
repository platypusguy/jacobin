/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

import Foundation

/// the data structure for holding a field's field_info from the class file.
/// the fields are explained in the JVM spec at:
/// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.5
///
// Layout:
// field_info {
//    u2             access_flags;
//    u2             name_index;
//    u2             descriptor_index;
//    u2             attributes_count;
//    attribute_info attributes[attributes_count];
//  }
class Field {
    var accessFlags : Int16 = 0x00
    var name = ""
    var description = ""
    var attributes : [FieldInitAttribute] = []

    /// loads the above data fields with the data from the classfile. Details here:
    /// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.5
    /// - Parameters:
    ///   - klass: the klass whose fields we're loading
    ///   - location: where we are in the class file, as an index to array of bytes
    /// - Returns: update location
    func load( klass: LoadedClass, location: Int ) -> Int {
        var loc = location
        accessFlags = Utility.getInt16from2Bytes( msb: klass.rawBytes[loc+1],
                                                  lsb: klass.rawBytes[loc+2] )
        loc += 2

        // get the name string
        let nameIndex = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: loc+1 )
        loc += 2
        name = Utility.getUTF8stringFromConstantPoolIndex( klass: klass, index: nameIndex )

        // get the description string
        let descIndex = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: loc+1 )
        loc += 2
        description = Utility.getUTF8stringFromConstantPoolIndex( klass: klass, index: descIndex )

        // get the attribute count
        let attrCount = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: loc+1 )
        loc += 2

        print( "Class \(klass.shortName), field: \(name), description: \(description), attributes: \(attrCount)" )

        if attrCount > 0 {
            for _ in 1...attrCount {
                // var containing the constant value, if any is specified
                var value: Any = ""

                // get attr name
                let attrNameIndex = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: loc+1 )
                guard klass.cp[attrNameIndex].type == .UTF8 else { // verify we're pointing at a UTF8 rec
                    jacobin.log.log( msg: "Class \(klass.shortName), field: \(name), description: invalid attribute pointer. skipped",
                                     level: Logger.Level.WARNING )
                    loc += 8
                    continue
                }
                loc += 2
                let attrLen = Utility.getIntfrom4Bytes(bytes: klass.rawBytes, index: loc+1)
                loc += 4
                let attrName = Utility.getUTF8stringFromConstantPoolIndex( klass: klass, index: attrNameIndex )

                // the ConstantValue entry points to a record in the constant pool that
                // contains the value to initialize the field to. This logic gets that
                // record and displays the number.
                if attrName == "ConstantValue" { //TODO: Check that ACC_STATIC is set per section 4.7
                    let cpPointer = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: loc+1 )
                    guard cpPointer > 0 && cpPointer < klass.constantPoolCount else {
                        jacobin.log.log( msg: "Class \(klass.shortName), field: \(name), invalid constant value ptr",
                                         level: Logger.Level.WARNING )
                        loc += attrLen
                        continue
                    }
                    let cpRec = klass.cp[cpPointer]

                    switch cpRec.type {
                        case .intConst:
                            let cpIntConst =  cpRec as! CpIntegerConstant
                            value = cpIntConst.int
                            print( "field \(name) is integer initialized to: \(value)" )
                        default:
                            print( "field \(name) is initialized" )
                            value = "?"
                    } //CURR: add the other constant types (long, double, etc.)
                }
                let fInit = FieldInitAttribute( type: description, value: value )
                attributes.append( fInit )
                loc += attrLen
            }
        }
        return loc
    }
}
