/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

import Foundation

class FieldCount {

    // extract the number of fields in this class
    static func readFieldCount( klass: LoadedClass, location: Int ) {
        let fieldCount = Utility.getIntFrom2Bytes( bytes: klass.rawBytes, index: location+1 )
        klass.fieldCount = fieldCount
    }

    // log the number of fields. Mostly used for diagnostic purposes.
    static func log( klass: LoadedClass ) {
        jacobin.log.log( msg: "Class: \( klass.path ) - # of fields: \( klass.fieldCount )",
                level: Logger.Level.FINEST )
    }
}
