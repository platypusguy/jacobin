/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

// Exceptions and error handlers

import Foundation

enum JVMerror : Error {
    case ClassFormatError( msg: String )
    case ClassVerificationError( msg: String )
    case InvalidParameterError( msg: String )
    case UnreachableError( msg: String )
}

// Private note: for many errors, use this message:
// JVMerror.ClassFormatError( msg: "in: \(#file), func: \(#function) line: \(#line)" )