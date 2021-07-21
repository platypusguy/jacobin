[![Swift](https://github.com/platypusguy/jacobin/actions/workflows/swift.yml/badge.svg)](https://github.com/platypusguy/jacobin/actions/workflows/swift.yml)
![GitHub](https://img.shields.io/github/license/platypusguy/jacobin)

# jacobin

A more-than-minimal JVM written in Swift. 

[This overview](https://github.com/platypusguy/jacobin/wiki/Jacobin-Overview) gives the background on this project, including its aspirations and the features that it supports. The remaining pages discuss the basics of JVM operation and, where applicable, how Jacobin implements the various steps, noting any items that would be of particular interest to JVM cognoscenti. I've included references to the official JVM docs, where I can both as a reference for you, the reader, and for the Jacobin team's easy reference. 

# Status
## Intended feature set:
* Java 11 functionality, but...
* No JNI (Oracle intends to replace it; see [JEP 389](https://openjdk.java.net/jeps/389))
* No security manager (Oracle intends to remove it; see [JEP 411](https://openjdk.java.net/jeps/411))
* No JIT
* Somewhat less stringent bytecode verification

## What we've done so far and what we need to do:
### Command-line parsing
* Gets options from the three environment variables. [Details here](https://github.com/platypusguy/jacobin/wiki/Command-line-parameters)
* Parses the command line; identify JVM options and application options
* Responds to most options listed in the `java -help` output

**To do**:
  * Handling JAR files
  * Handling @files (which contain command-line options)
  * Parsing the classpath

### Class loading
* Correctly reads and parses basic classes
* Extracts bytecode and params needed for execution

**To do**:
* Handle more-complex classes
* Handle interfaces
* Handle arrays
* Handle inner classes
* Automate loading of core Java classes (Object, etc.)

### Verification, Linking, Preparation, Initialization
* Performs integrity check bytecode is correct. :pencil2: This is the focus of current coding work

**To do:**
* Linking and verification
* Preparation
* Initialization

### Execution
Not started yet

## Garbage Collection
GC is handled by the Swift runtime, which has its own GC

# Thanks
The overview page carries a right-side panel that expresses our thanks to vendors and programmers who have made the Jacobin project possible. They are many and we are deeply grateful to them.
