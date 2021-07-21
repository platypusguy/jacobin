/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

import Foundation

/// handles extracting the the name of interfaces implemented by this class
/// from the constant pool. This class is called from the classloader

class Interface {

    /// read the interface data, which is at the constant pool
    /// entry pointed to by index, which should be a CONSTANT_Class
    /// - Parameters:
    ///   - klass: the class containing the constant pool and the interface
    ///   - index: pointer to the constant pool entry for this interface
    static func process( klass: LoadedClass, index: Int ) {
        let interface = klass.cp[index] as! CpEntryClassRef
        let intNameIdx = interface.classNameIndex
        if intNameIdx < 1 || intNameIdx >= klass.constantPoolCount {
            jacobin.log.log( msg: "Error: In \(klass.shortName) - invalid UTF8 index \(index) in Interface entry",
                             level: Logger.Level.INFO )
            return
        }

        let interfaceName =
            Utility.getUTF8stringFromConstantPoolIndex( klass: klass, index: intNameIdx)
        klass.interfaces.append( interfaceName )
    }

    // log the value (mostly used for diagnostic purposes)
    static func logAll( klass: LoadedClass ) {
        if ( klass.interfaceCount > 0 ) {
            for i in 1...klass.interfaceCount {
                jacobin.log.log( msg: "Class: \(klass.path) implements: \(klass.interfaces[i-1])",
                                 level: Logger.Level.FINEST )
            }
        }
    }
}