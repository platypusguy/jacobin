/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

/// Outputs most messages to users
///
import Foundation

class UserMsgs {

    // shows the usage (switches and params) to the user, in response to error or command-line 'help' request
    // input Streams identifies whether message is shown to stdout or stderr
    static func showUsage( stream:  Streams ) {
        let usage =
                """
                Usage: jacobin [options] <mainclass> [args...]
                          (to execute a class)
                    or jacobin [options] -jar <jarfile> [args...]
                          (to execute a jar file)
                Arguments following the main class, source file, -jar <jarfile>,
                are passed as the arguments to main class.

                where options include:

                    -? -h -help
                                  print this help message to the error stream
                    --help        print this help message to the output stream
                    -version      print product version to the error stream and exit
                    --version     print product version to the output stream and exit
                    -showversion  print product version to the error stream and continue
                    --show-version
                                  print product version to the output stream and continue


                """
        threads.wait() // prevents logging info being partially overwritten by this
        fputs( usage + "\n", stream == Streams.sout ? stdout : stderr )
    }

    // shows the version number of this instance of Jacobin JVM
    // input Streams identifies whether message is shown to stdout or stderr
    static func showVersion( stream: Streams ) {
        let version =
                """
                jacobin JVM v. \(globals.version) 2021
                64-bit server JVM
                """
        threads.wait() // prevents logging info being partially overwritten by this
        fputs( version + "\n", stream == Streams.sout ? stdout : stderr )
    }
}
