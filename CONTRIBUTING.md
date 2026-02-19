## The basic layout of the Jacobin JVM and the coding conventions

## Jacobin JVM Operation

The Jacobin JVM is a multi-threaded interpreter that executes Java 21 bytecode.

### Jacobin JVM Startup Sequence

Jacobin starts up in `jvm/jvmStart.go`

* It initializes the global variables defined in `globals/globals.go`

* It checks environment variables and command-line arguments. The latter are handled
in `jvm/cli.go` based on the flags defined in `jvm/option_table_loader.go`

* It then loads the base classes from the JDK. These are the files located in the java.base.jmod
file, which consists of approximately 1700 classes that are parsed and loaded into the method area,
which is the area in the JVM where all class metadata is stored. In Jacobin, it's a map
defined in `classloader/methArea.go`

* It loads required static values into the `Statics` map defined in `statics/statics.go`

* It creates the `classloader.MTable`, which is a cache of methods. Every executed method
is first loaded into the `MTable`, from which it is executed. 

* It loads references to all gfunctions into the MTable. gfunctions are go implementations
of Java library functions (primarily library functions that are implemented as native
code in the JDK). The gfunction directory contains all gfunctions, arranged in subdirectories
by Java package. For example, `java/lang/Float` is found in `gfunction/javaLang/javaLangFloat.go`.

* It then creates thread groups and the main thread.

* It then instantiates the main class and start the `main()` method in the main thread

### JVM Operation (Similarities and differences from HotSpot JVM)

Jacobin JVM hews closely to the way HotSpot JVM works. At present, it is an interpreter
only, although initial work on JIT'ing is being discussed. The project 
[status page](https://github.com/platypusguy/jacobin/blob/main/README.md) gives
details on the current progress.

When instantiating a class--including performing the various preparation and link steps--Jacobin
verifies the code (see `classloader/codecheck.go`) except for classes in the JDK. 

Jacobin JVM uses a 64-bit wide operand stack, as opposed to HotSpot, which uses 32-bit wide operand stack.
This affects certain bytecode instructions, which are fixed as needed in codecheck.go.

HotSpot stores a class's static fields in the `java/lang/Class` instance for that class, whereas
Jacobin stores them in the JVM-wide statics map.

Jacobin uses compact strings to hold the names of classes. These are stored in the
string pool (see `stringPool/StringPool.go`). A greater use of compact strings is planned.

Jacobin class metadata stores the entries for the constant pool in an unresolved form, meaning
they need to be resolved each time they're accessed. They will eventually be migrated to
resolved form, when our focus returns to performance optimization.

Jacobin supports multiple execution threads. We have tested it up to 64 threads, but we are quite
confident it can easily handle more.

## Code layout

main.go contains the main function. All other source code is located in the `src/` directory.
The main packages/directories are:
* `classloader` contains all the class parsing and loading logic
* `gfunction` and its subdirectories contain all the gfunctions, which implement Java library methods in go
* `jvm` is the heart of Jacobin: it contains the interpreter and numerous run utitlities
* `object` contains a definition of objects
* `types` contains definitions of various types as well as a variety of constants

## Coding Conventions

* All objects are passed as pointers to `object.Object`
* All integers used internally are `int64` with the exception of `uint16` values used by Java in the class constant pool
* All raw strings stored in objects are `object.javaByteArray` instances, rather than golang strings, or as
instances of Java string objects (as defined in `object/string.go`)
* All bytes used for data are not golang `byte`s (which are unsigned) but rather `types.JavaByte`, which are signed
* All floating point values are `float64`
* When exceptions or errors are thrown in Java methods, we first create a variable named `errMsg`, which contains error data for the user;
the exception is thrown using `exceptions.ThrowEx(excNames.NoSuchMethodException, errMsg, fr)`
where the first argument is the name of the exception class, the second is the error message, and the third is the frame
where the exception was thrown. If the frame data is not known, it can be passed as `nil`

### Coding conventions for gfunctions
Gfunctions are golang implementations of Java library functions. Not all library methods,
but mostly methods that are implemented as native code in the JDK. They have tight coding constraints:
* All gfunctions have the same signature: `funcName(params []any)any` All parameters are passed
in the `params` array. This includes the `this` pointer, which is a reference to the object whose
the method is called. For example, `s.length()`, s is the `this` object. If there is no `this`
and no argument, `params` is an empty array. 
* When errors occur in gfunctions, they are handled by creating an `errMsg` variable, which contains
error data for the user. This error message and the name of the exception/error to throw are returned
using `ghelpers.GetGErrBlk(excNames.ClassNotLoadedException, errMsg)`

### Coding conventions for unit tests
* Tests that involve any variable or reference outside the function under test should start with
a call to `globals.InitGlobals("test")`, which initializes many global variables required
for almost any action outside the function under test.
* If it's necessary to capture stderr, the code to be used is shown 
[here](https://github.com/platypusguy/jacobin/blob/main/notes/StderrInTests.txt). This
code enables the stderr messages to be viewed as a single golang string.
* As a rule, we do not use tables in unit tests. 
* All unit tests should be in alphabetical order in the file; helper functions should be
segregated at the top or bottom of the file with conspicuous comments.
* All unit tests should be preceded with a comment of 1 or more lines explaining the test.
If a specific condition is being tested, it should be mentioned, so that we can quickly
see which conditions have/have not been tested without reading the code.

### General coding conventions
* All comments should be in English.
* We follow the `go fmt` coding style (not because we think it's better, but because it's
the standard for golang).
* All functions should have a preceding comment that explains what they do. Each 
file should start with an introductory block comment explaining what functionality it
implements, so that we can quickly tell whether the file is relevant to our present
search.
right away if it contains the functionality we're looking for.
* We generally refer to the Go language as "golang," because this facilitates search

