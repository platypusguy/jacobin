/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

import Foundation

/// handles the lineNumberTable attribute of a method's code
// The layout of the table, per the JVM spec:
//    u2 attribute_name_index; (this has already been obtained when this class is created
//    u4 attribute_length;
//    u2 line_number_table_length;
//    {   u2 start_pc;
//        u2 line_number;
//    } line_number_table[line_number_table_length];

class LineNumberTable {
    struct Entry { var pc: Int, line: Int }
    var entryCount = 0
    var entries : [Entry] = []

    func load( klass: [UInt8], loc: Int ) -> Int {
        var currLoc = loc
        let attrLength = Utility.getIntfrom4Bytes(bytes: klass, index: currLoc + 1 )
        currLoc += 4
        print( "Line number table size: \(attrLength)" )

        entryCount = Utility.getIntFrom2Bytes(bytes: klass, index: currLoc + 1 )
        print( "Line number entries: \(entryCount)" )

        currLoc += 2
        for _ in 1...entryCount {
            let pc   = Utility.getIntFrom2Bytes(bytes: klass, index: currLoc+1 )
            let line = Utility.getIntFrom2Bytes(bytes: klass, index: currLoc+3 )
            let e = Entry( pc: pc, line: line )
            entries.append( e )
            print( "Line number entry: pc: \(pc) line# \(line)" )
            currLoc += 4
        }

        return currLoc
    }
}
