/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */
import Foundation

/// the principal data structure for holding method info. Attaches to a LoadedClass.
/// Most of the data items are extracted from the class file by MethodInfo.swift,
/// and described in section 4.x of the JVM spec.

class Method {
    var name = ""
    var descriptor = ""
    var maxStack = 0
    var maxLocals = 0
    var attributeCount = 0
    var codeLength = 0
    var code : [UInt8] = []
    var exceptionTableLength = 0
    struct ExceptionEntry {
        var startPc = 0
        var endPc = 0
        var handlerPc = 0
        var catchType = 0
    }
    var exceptionTable : [ExceptionEntry] = []

    var accessFlags: Int16 = 0
    func isPublic()       -> Bool {( accessFlags & 0x0001 ) > 0 }
    func isPrivate()      -> Bool {( accessFlags & 0x0002 ) > 0 }
    func isProtected()    -> Bool {( accessFlags & 0x0004 ) > 0 }
    func isStatic()       -> Bool {( accessFlags & 0x0008 ) > 0 }
    func isFinal()        -> Bool {( accessFlags & 0x0010 ) > 0 }
    func isSynchronized() -> Bool {( accessFlags & 0x0020 ) > 0 }
    func isBridge()       -> Bool {( accessFlags & 0x0040 ) > 0 }
    func isVarargs()      -> Bool {( accessFlags & 0x0080 ) > 0 }
    func isNative()       -> Bool {( accessFlags & 0x0100 ) > 0 }
    func isAbstract()     -> Bool {( accessFlags & 0x0400 ) > 0 }
    func isStrictFP()     -> Bool {( accessFlags & 0x0800 ) > 0 }
    func isSynthetic()    -> Bool {( accessFlags & 0x1000 ) > 0 }

    typealias LineNumber = [Int]
    var lineNumTable : [LineNumber] = []

    var parameters: [MethodParm] = []
    var deprecated = false
    var synthetic = false
}
