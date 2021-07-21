/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

import Foundation

/// An optional Class attribute that holds the name of the source file that
/// generated the class. Further details here:
/// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.7.10
class SourceFileAttribute: Attribute {
    var fileName = ""

    override init( name: String, length: Int ) {
        super.init( name: name, length: length )
    }

    /// get the two-byte index into the constant pool that points to the
    /// name of the source file and then put it in the fileName field
    /// - Parameters:
    ///   - bytes: the raw bytes of the class
    ///   - loc: the present location in the class bytes so far
    func load( klass: LoadedClass, loc: Int ) {
        let nameIndex =
            Utility.getIntFrom2Bytes(bytes: klass.rawBytes, index: loc+1 )
        fileName =
            Utility.getUTF8stringFromConstantPoolIndex(klass: klass, index: nameIndex )
    }

    func log( className: String ) {
        jacobin.log.log( msg: "Class \(className), source file: \( fileName )",
                         level: Logger.Level.FINEST )
    }
}
