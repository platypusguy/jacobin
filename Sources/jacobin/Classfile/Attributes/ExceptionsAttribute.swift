/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

import Foundation

/// handles data on the checked exceptions a method can throw. Details at:
/// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.7.5
class ExceptionsAttribute: Attribute {
    var exceptionCount = 0

    override init( name: String, length: Int ) {
        super.init( name: name, length: length )
    }

    /// load up this class with the attribute data
    /// for the moment, just record the number of exceptions
    /// - Parameters:
    ///   - klass: the class whose exceptions attribute we're processing
    ///   - loc: where in the class's bytecode we are currently
    func load( klass: LoadedClass, loc: Int ) {
        exceptionCount =
            Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: loc + 1 )
    }

    func log( klass: LoadedClass, method: Method ) {
        jacobin.log.log( msg: "Class \( klass.shortName ), Method \( method.name ), " +
                             "# of exceptions: \( exceptionCount )",
                         level: Logger.Level.FINEST )
    }
}
