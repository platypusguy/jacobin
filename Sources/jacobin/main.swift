/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

/// main line of jacobin

import Dispatch
import Foundation

var globals  = Globals( startTime: DispatchTime.now() )

let logQueue = DispatchQueue( label: "logQueue" )
let threads  = DispatchGroup()
let log = Logger()
try main()
shutdown( successFlag: true )


func main() throws {
    do {
        if ( CommandLine.arguments.contains( "-vverbose" ) ) {
            globals.logLevel = Logger.Level.FINEST;
        }
        globals.logLevel = Logger.Level.FINEST; //for the nonce -- remove eventually
        log.log( msg: "starting Jacobin VM", level: Logger.Level.FINE )
        processCommandLine( args: CommandLine.arguments )
        try loadClasses( startingClass: globals.startingClass )
    } catch {
        log.log( msg: "Unexpected error in Jacobin VM. Exiting", level: Logger.Level.SEVERE )
        shutdown( successFlag: false )
    }
}

// load classes starting with the main class of the app
func loadClasses( startingClass: String ) throws {
    globals.systemLoader.add( name: startingClass )
}

// parse the command line and capture all the various settings it specifies.
func processCommandLine( args: [String]) {
    let cp = CommandLineProcessor()
    cp.process( args: args )
    if  cp.dispatch( commandLine: globals.commandLine ) == execStop {
        shutdown( successFlag: true )
    }
}

/// shuts downt the JVM. If passed 'true' it's a normal shutdown, if 'false' this indicates an error was the cause
func shutdown( successFlag : Bool ) {
    log.log ( msg: "shutting down Jacobin VM", level: Logger.Level.FINE )
    threads.wait()
    exit( successFlag ? 0 : -1 )
}


