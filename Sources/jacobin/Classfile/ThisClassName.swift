/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

import Foundation

/// handles extracting the short name of this class from the constant pool. Ex: Hello.class returns: Hello
/// This class is called from the classloader

class ThisClassName {
    static var thisClassRef = 0;

    // reads the entry in the class file that points to the short name for this class
    static func readName( klass: LoadedClass, location: Int ) {
        thisClassRef = Int(Utility.getInt16from2Bytes( msb: klass.rawBytes[location+1],
                                                       lsb: klass.rawBytes[location+2] ))
    }

    // verifies that the entry points to the right type of record.
    static func verify( klass: LoadedClass ) {
        if( klass.cp[thisClassRef].type != .classRef ) { // must point to a class reference
            jacobin.log.log( msg: "ClassFormatError in \(klass.path): Invalid thisClassReference",
                             level: Logger.Level.SEVERE )
            shutdown( successFlag: false )
        }
    }

    // looks up the pointed-to name for this class and inserts it into klass.shortName; and logs it
    static func process( klass: LoadedClass ){
        let cRef : CpEntryClassRef = klass.cp[thisClassRef] as! CpEntryClassRef
        let pointerToName = cRef.classNameIndex
        let shortNameEntry : CpEntryUTF8 = klass.cp[pointerToName] as! CpEntryUTF8
        klass.shortName = shortNameEntry.string
    }

    // log the class name (mostly used for diagnostic purposes)
    static func log( klass: LoadedClass ) {
        jacobin.log.log( msg: "Class: \(klass.path) - short name: \(klass.shortName)",
                         level: Logger.Level.FINEST )
    }
}
