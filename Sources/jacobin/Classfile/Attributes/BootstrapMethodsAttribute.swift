/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

import Foundation

/// An optional Class attribute that holds the data for bootstrap methods
/// used by invokedynamic. Consult:
/// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.7.23
/// BootstrapMethods_attribute {
//    u2 attribute_name_index;
//    u4 attribute_length;
//    u2 num_bootstrap_methods;
//    {   u2 bootstrap_method_ref;
//        u2 num_bootstrap_arguments;
//        u2 bootstrap_arguments[num_bootstrap_arguments];
//    } bootstrap_methods[num_bootstrap_methods];
//}
class BootstrapMethodsAttribute: Attribute {
    var fileName = ""

    override init( name: String, length: Int ) {
        super.init( name: name, length: length )
    }

    /// Fill in the table of boostrap methods
    ///   - bytes: the raw bytes of the class
    ///   - loc: the present location in the class bytes so far
    func load( klass: LoadedClass, loc: Int ) throws {
        var currLoc = loc
        let entriesCount = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: currLoc + 1 )
        currLoc += 2
        guard entriesCount > 0 else {
            jacobin.log.log( msg: "Invalid BoostrapMethods attribute in: \(klass.shortName)",
                level: Logger.Level.SEVERE )
            throw JVMerror.ClassFormatError( msg: "\(#file), func: \(#function) line: \(#line)" )
        }

        for _ in 1...entriesCount {
            var bsm = BootstrapMethod()
            bsm.methodRef = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: currLoc + 1 )
            bsm.argCount  = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: currLoc + 3 )
            currLoc += 4

            if bsm.argCount > 0 {
                for _ in 1..<bsm.argCount {
                    let bsmArg = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: currLoc + 1 )
                    currLoc += 2
                    bsm.args.append( bsmArg )
                }
            }
            klass.bsms.append( bsm )
        }
    }

    func log( klass: LoadedClass ) {
        jacobin.log.log( msg: "Class \( klass.shortName ) has \(klass.bsms.count) bootstrap methods",
                         level: Logger.Level.FINEST )
    }
}


/// the bootstrap method that is added to the loaded class
struct BootstrapMethod {
    var methodRef = 0
    var argCount = 0
    var args : [Int] = []
}