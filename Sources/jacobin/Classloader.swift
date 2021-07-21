/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

/// handles the loading of classes

import Foundation

class Classloader {
    var name = ""
    var parent = ""
    var cl : [String: LoadedClass] = [:]

    /// create Classloader with name and pointing to parent.
    init( name: String, parent: String ) {
        self.name = name
        self.parent = parent
    }

    /// inserts a parsed class into the classloader, if it's not already there
    /// - Parameters:
    ///   - name: the name of the class
    ///   - klass: the Swift object containing all the needed parsed data
    private func insert( name: String, klass: LoadedClass ) {
        if( cl[name] == nil ) {
            cl[name] = klass;
        }
    }

    /// add a class for which we have only the name, provided that it's not already
    /// in this classloader
    /// - Parameter name: the name of the class
    func add( name: String ) {
        if cl[name] != nil { return } // do nothing if the class is already loaded
                                      // TODO: go up the chain of classloaders
        let klass = LoadedClass()
        do {
            try ClassParser.parseClassfile(name: name, klass: klass )

            // format check of the parsed class defaults to true for non-
            // bootstrap classes. However, this can be changed on the command
            // line with the -Xverify option, which we consult here via the
            // globals.verifyBytecode setting

            if globals.verifyBytecode != .none {
                if( name == "bootstrap " && globals.verifyBytecode == .all ) ||
                    name != "bootstrap" {
                    try FormatCheck.check( klass: klass )
                    klass.status = .CHECKED
                    log.log( msg: "Class \(klass.shortName) format checked", level: Logger.Level.FINEST )
                }
            }
            insert( name: name, klass: klass )
        }
        catch JVMerror.ClassFormatError {
            shutdown( successFlag: false ) // error msg has already been shown to user
        }
        catch JVMerror.ClassVerificationError( let msg ) {
            log.log( msg: "Error info: \(msg)", level: Logger.Level.SEVERE )
        }
        catch { // any other errors are unexpected, we should tell the user
            log.log( msg: "Unexpected error loading class \(name)",
                     level: Logger.Level.SEVERE )
        }
    }
}

