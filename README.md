![Go version](https://img.shields.io/github/go-mod/go-version/platypusguy/jacobin?filename=src%2Fgo.mod)
![Workflow](https://github.com/platypusguy/jacobin/actions/workflows/go.yml/badge.svg)
[![Go_report_card](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)](https://goreportcard.com/report/github.com/platypusguy/jacobin)
![GitHub](https://img.shields.io/github/license/platypusguy/jacobin)

# jacobin

A more-than-minimal JVM written in Go. 

<img src="https://github.com/platypusguy/jacobin/blob/0aedac33af431ca3befd67d96d0d95db84096b0c/assets/img/JacobinLogo.jpg" width=30% height=30%>


# Status
## Intended feature set:
* Java 17 functionality, but...
* No JNI (Oracle intends to replace it; see [JEP 389](https://openjdk.java.net/jeps/389))
* No security manager (Oracle intends to remove it; see [JEP 411](https://openjdk.java.net/jeps/411))
* No JIT
* Somewhat less stringent bytecode verification
* Does not enforce Java 17's sealed classes

## What we've done so far and what we need to do:
### Command-line parsing
* Gets options from the three environment variables. [Details here](https://github.com/platypusguy/jacobin/wiki/Command-line-Processing)
* Parses the command line; identify JVM options and application options
* Responds to most options listed in the `java -help` output

**To do**:
 * Handling @files (which contain command-line options)
 * Parsing complex classpaths

### Class loading
* Correctly reads and parses most classes
* Extracts bytecode and parameters needed for execution
* Automate loading of core Java classes (Object, etc.)
* Handles straightforward JAR files
  
**To do**:
* Handle more-complex classes
* Handle interfaces
* Handle inner classes

### Verification, Linking, Preparation, Initialization
* Performs [format check](https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.8) of class file.
* Linking, preparation, and initialization -- minimally and only as needed at execution time

**To do:**
* Verification
* Robust preparation and initialization

### Execution
* Execution of bytecode :pencil2: The primary focus of current coding work<br>
  180 bytecodes fully operational, including one- and multi-dimensional arrays
  
**To do:**
* invokespecial, invokedynamic
* Calls to superclasses
* Inner and nested classes
* Exception-tree walking
* Annotations

### Instrumentation
* Instruction-level tracing (use `-trace:inst` to enable this feature)
* Extensive logging data (use `-verbose:finest` to enable. Caveat: this produces *a lot* of data)

**To do:**
* Emit instrumented data to a port, for reading/display by a separate program.

## Garbage Collection
GC is handled by the golang runtime, which has its own GC

## Understanding the Code
A detailed roadmap to the code base can be found [in the wiki](https://github.com/platypusguy/jacobin/wiki/Roadmap-to-Jacobin-source-code).

# Thanks
The project's [home page](https://jacobin.org/) carries a lengthy note at the bottom that expresses our thanks to vendors and programmers who have made the Jacobin project possible. They are many and we are deeply grateful to them.
