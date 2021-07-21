/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */
import Foundation

/// handles data on method parameters. Details at:
/// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.7.24
class MethodParmsAttribute: Attribute {
    /*
     MethodParameters_attribute {
        u2 attribute_name_index;
        u4 attribute_length;
        u1 parameters_count;
        {
            u2 name_index;
            u2 access_flags;
        } parameters[parameters_count];
    }
    */

    var parmsCount = 0
    var parms: [MethodParm] = []


    override init( name: String, length: Int ) {
        super.init( name: name, length: length )
    }

    /// load up this class with the method parameter data, which consists of an array of
    /// structs containing the name of the parameter and the access flags for that parameter
    /// - Parameters:
    ///   - klass: the class whose exceptions attribute we're processing
    ///   - loc: where in the class's bytecode we are currently
    func load( klass: LoadedClass, loc: Int ) {
        var currLoc = loc
        parmsCount = Int(klass.rawBytes[currLoc+1])
        currLoc += 1

        for _ in 1...parmsCount {
            var parm = MethodParm()
            let nameIndex = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: currLoc+1 )
            currLoc += 2
            if nameIndex == 0 {
                parm.name = ""
            }
            else {
                parm.name =
                     Utility.getUTF8stringFromConstantPoolIndex( klass: klass, index: nameIndex )
            }
            parm.accessMask =
                    UInt16( Utility.getIntFrom2Bytes(bytes: klass.rawBytes, index: currLoc ))
            parms.append( parm )
        }
    }


    func log( klass: LoadedClass, method: Method ) {
        jacobin.log.log( msg: "Class \( klass.shortName ), Method \( method.name ), " +
                "# of parameters: \( parmsCount )",
                level: Logger.Level.FINEST )
    }
}
