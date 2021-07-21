/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

import Foundation

//various utility functions

class Utility {

    // converts two successive bytes into a 16-bit integer
    static func getInt16from2Bytes( msb: UInt8, lsb: UInt8 ) -> Int16 {
        return Int16(( Int( msb ) * 256 ) + Int( lsb ))
    }

    // converts two successive bytes into an Int
    static func getIntFrom2Bytes( bytes: [UInt8], index: Int ) -> Int {
        Int(Int( bytes[index] )*256 ) + Int( bytes[index+1] )
    }

    // converts four successive bytes into an Int. Note index points to first byte.
    static func getIntfrom4Bytes( bytes: [UInt8], index: Int ) -> Int {
        var loc = index
        let i1 = getIntFrom2Bytes( bytes: bytes, index: loc )
        loc += 2
        let i2 = getIntFrom2Bytes( bytes: bytes, index: loc )
        return(( i1 * 65536 ) + i2 )
    }

    // converts four bytes into a 32-bit integer
    static func getInt32from4Bytes( byte1: UInt8, byte2: UInt8,
                                    byte3: UInt8, byte4: UInt8 ) -> Int32 {
        let firstBytes = getInt16from2Bytes( msb: byte1, lsb: byte2 )
        let lastBytes  = getInt16from2Bytes( msb: byte3, lsb: byte4 )
        return (( Int32( firstBytes ) * 65536 ) + ( Int32( lastBytes )))
    }

    // returns a UTF8 string pointed to by an index into the constant pool
    static func getUTF8stringFromConstantPoolIndex( klass: LoadedClass, index: Int ) -> String {
        let cpEntry = klass.cp[index]
        if cpEntry.type != ConstantPool.RecType.UTF8 {
            jacobin.log.log( msg: "Error: Class: \(klass.path) - invalid UTF8 index \(index)",
                             level: Logger.Level.WARNING )
            return ""
        }
        let UTF8entry = cpEntry as! CpEntryUTF8
        return UTF8entry.string
    }

    // safety mechanism to insert in unreachable code and to throw an error
    // if the code is unexpectedly reached
    static func unreachableCode() throws {
        let msg = "Unreachable code violation: \(#file) at line\(#line)"
        jacobin.log.log(msg: msg, level: Logger.Level.SEVERE )
        throw JVMerror.UnreachableError( msg: "" )
    }
}
