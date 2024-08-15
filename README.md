![Go version](https://img.shields.io/github/go-mod/go-version/platypusguy/jacobin?filename=src%2Fgo.mod)
![Workflow](https://github.com/platypusguy/jacobin/actions/workflows/go.yml/badge.svg)
<img alt="GitHub commit activity" src="https://img.shields.io/github/commit-activity/m/platypusguy/jacobin">
<!--
[![Go_report_card](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)](https://goreportcard.com/report/github.com/platypusguy/jacobin) -->
<!-- ![GitHub](https://img.shields.io/github/license/platypusguy/jacobin) -->

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
* Automated pre-loading of core Java classes (`Object`, etc.)
* `java.*`, `javax.*`, `jdk.*`, `sun.*` classes are loaded from the `JAVA_HOME` directory (i.e., from JDK binaries)
* Handles JAR files
* Handles interfaces
  
**To do**:
* Handle more-complex classes (called via method handles, etc.)
* Handle inner classes

### Verification, Linking, Preparation, Initialization
* Performs [format check](https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.8) of class file.
* Linking, preparation, and initialization -- minimally and only as needed at execution time

**To do:**
* Verification
* Robust preparation and initialization

### Execution
* Execution of bytecodes :pencil2: The primary focus of current coding work<br>
  All bytecodes except invokedynamic implemented, including one- and multi-dimensional arrays
* Static initialization blocks
* Throwing and catching exceptions
* Running native functions (written in go). [Details here.](https://github.com/platypusguy/jacobin/wiki/Native-golang-functions-methods )
  
**To do:**
* invokedynamic
* Calls to superclasses
* Inner and nested classes
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

# If you want to test Jacobin
At present, we're not looking for testers because we know the missing features and we're working on them. Testing, at this point, will likely result in frustration. However, if for your own enjoyment, you still want to try it out, see directions and cautions on our [Release Page](https://github.com/platypusguy/jacobin/releases). (If you want some fun, run your program on Jacobin with the `-trace:inst` option and watch the executing Java bytecodes whiz by along with the changing contents of the operand stack.) 

We expect/hope/trust that by the end of this year, we'll be ready to ask interested users to test Jacobin on real programs and share their feedback with us. 

# Thanks
The project's [home page](https://jacobin.org/) carries a lengthy note at the bottom that expresses our thanks to vendors and programmers who have made the Jacobin project possible. They are many and we are deeply grateful to them.
