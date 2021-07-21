/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

import Foundation

/// struct holding method parameters as stored in the method_info data structure
/// name can be "" in certain cases (such as compiler-generated parameters)
/// accessMask is a series of bits. Details on both fields here:
/// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.7.24

struct MethodParm {
        var name: String = ""
        var accessMask: UInt16 = 0x00
}
