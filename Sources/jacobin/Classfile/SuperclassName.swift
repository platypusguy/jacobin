/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

import Foundation

/// handles extracting the short name of the superclass of this class from the constant pool.
/// This class is called from the classloader

class SuperClassName {
    static var superClassRef = 0

    // reads the entry in the class file that points to the superclass for this class
    static func readName( klass: LoadedClass, location: Int ) {
        superClassRef = Int(Utility.getInt16from2Bytes( msb: klass.rawBytes[location+1],
                                                        lsb: klass.rawBytes[location+2] ))
    }

    // verifies that the entry points to the right type of record.
    static func verify( klass: LoadedClass ) throws {
        if( klass.cp[superClassRef].type != .classRef &&  // must point to valid class unless this class is Object.class
                         klass.shortName != "java/lang/Object" ) {
            jacobin.log.log( msg: "ClassFormatError in \( klass.path ): Invalid superClassReference",
                             level: Logger.Level.SEVERE )
            throw JVMerror.ClassFormatError( msg: "in: \(#file), func: \(#function) line: \(#line)" )
        }
    }

    // looks up the pointed-to name for the superclass and inserts it into klass.shortName; and logs it
    static func process( klass: LoadedClass ){
        if klass.shortName == "java/lang/Object" {
            klass.superClassName = ""
        }
        else {
            let cRef: CpEntryClassRef = klass.cp[superClassRef] as! CpEntryClassRef
            let pointerToName = cRef.classNameIndex
            let superNameEntry: CpEntryUTF8 = klass.cp[pointerToName] as! CpEntryUTF8
            klass.superClassName = superNameEntry.string
        }
    }

    // log the name of the superclass (Mostly used for diagnostic purposes)
    static func log( klass: LoadedClass ) {
        jacobin.log.log( msg: "Class: \( klass.path ) - superclass: \( klass.superClassName )",
                level: Logger.Level.FINEST )
    }
}