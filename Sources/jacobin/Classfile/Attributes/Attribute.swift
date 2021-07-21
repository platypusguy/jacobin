/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

import Foundation


// Superclass for all the attributes found in class files
class Attribute {
    var attrName = ""
    var attrLength = 0

    init( name: String, length: Int ){
        attrName = name
        attrLength = length
    }
}

class AttributeInfo: Attribute {
    var attrInfo : [UInt8] = []
}
