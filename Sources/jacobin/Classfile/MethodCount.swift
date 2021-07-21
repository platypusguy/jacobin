/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

import Foundation

/// handles extracting the the number of methods implemented by this class from the constant pool.
/// This class is called from the classloader

class MethodCount {

    // read the number of methods (a 16-bit integer)
    static func readMethodCount( klass: LoadedClass, location: Int ) {
        let methodCount = Int(Utility.getInt16from2Bytes( msb: klass.rawBytes[location+1],
                lsb: klass.rawBytes[location+2] ))
        klass.methodCount = methodCount
    }

    // log the value (mostly used for diagnostic purposes)
    static func log( klass: LoadedClass ) {
        jacobin.log.log( msg: "Class: \( klass.path ) - # of methods: \( klass.methodCount )",
                level: Logger.Level.FINEST )
    }
}