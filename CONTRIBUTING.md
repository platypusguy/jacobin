## The basic layout of the Jacobin JVM and the coding conventions

## Jacobin JVM Operation

The Jacobin JVM is a multi-threaded interpreter that executes Java 21 classes.

### Jacobin JVM Startup Sequence

Jacobin starts up in `jvm/jvmStart.go` which performs the following steps:

* Initialization of the global variables and tables defined in `globals/globals.go`.

* Processing O/S environment variables and command-line arguments in `jvm/cli.go` based on the flags defined in `jvm/option_table_loader.go`.

* Parsing and loading the base classes from the JDK. There are approximately 1700 class files stored in the `java.base.jmod`
file. The target repository is the Jacobin method area where all class metadata is stored, a map
defined in `classloader/methArea.go`.

* Loading required static values into the `Statics` map defined in `statics/statics.go`

* Creation of the `classloader.MTable`, which is a cache of methods. Every executed method
is first loaded into the `MTable`, from which it is executed. 

* Addition of all gfunctions and their method signatures into the MTable. These are go replacements
for some of the Java library functions (implemented as Java classes or native
code in the JDK). The gfunction directory tree contains all gfunctions, arranged in subdirectories
by Java package. For example, `java/lang/Float` is found in `gfunction/javaLang/javaLangFloat.go`.

* Creates thread groups and the main thread.

* Finally, the main class is instantiated and the `main()` method of the main thread begins.

### JVM Operation (Similarities and differences from HotSpot JVM)

Jacobin JVM hews closely to the way HotSpot JVM works. At present, it is an interpreter
only, although the possibility of incorporating JIT capability is being discussed. The project 
[status page](https://github.com/platypusguy/jacobin/blob/main/README.md) gives
details on the current progress.

When instantiating a class, including performing the various preparation and link steps, Jacobin
verifies the code (see `classloader/codecheck.go`) except for classes in the JDK. 

Jacobin JVM uses a 64-bit wide operand stack, as opposed to HotSpot, which uses 32-bit wide operand stack.
This affects certain bytecode instructions, which are adjusted as needed in codecheck.go.

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

main.go in directory `src` contains the main function. All other source code is located in one of the directory trees subordinate to `src`.
The main packages/directories are:
* `classloader` contains all the class parsing and loading logic
* `gfunction` and its subdirectories contain all the gfunctions, which implement Java library methods in go
* `jvm` is the heart of Jacobin: it contains the interpreter and numerous run utitlities
* `object` contains a definition of objects
* `types` contains definitions of various types as well as a variety of constants

## Coding Conventions

* All objects are passed as pointers to `object.Object`.
* The most common integer format is `int64`. This is the format of all user arguments and return values (int, long, short). Internally, a variey of values are employed to be efficient or convenient for interfacing with the Go run-time.
* Booleans are handled as `int64` quantities having a value of `types.JavaBoolTrue` (int64 1) or `types.JavaBoolFalse` (int64 0).
However, boolean arrays store booleans as golang bytes. They are converted to/from `int64`s when operated on in the JVM.
* All raw strings stored in objects as developer-accessible values are `object.javaByteArray`s, rather than golang strings, or as instances of Java string objects (as defined in `object/string.go`).
* Byte-format used for data are not golang `byte`s (which are unsigned) but rather `types.JavaByte`, which are signed.
* Floating point values are `float64` whether they are Java float or double.
* When exceptions or errors are thrown in Java methods, we first create a variable named `errMsg`, which contains error data for the user;
the exception is thrown using a call such as `exceptions.ThrowEx(excNames.NoSuchMethodException, errMsg, fr)`
where the first argument is the name of the exception class, the second is the error message, and the third is the frame
where the exception was thrown. If the frame data is not known, it can be passed as `nil`.

### Coding conventions for gfunctions
Gfunctions have the following coding constraints:
* All gfunctions have the same signature: `funcName(params []any)any` All parameters are passed
in the `params` array. This includes the `this` pointer, which is a reference to the object whose
the method is called. For example, `s.length()`, s is the `this` object. If there is no `this`
and no argument, `params` is an empty array. 
* When errors occur in gfunctions, they are handled by creating an `errMsg` variable, which contains
error data for the user. This error message and the name of the exception/error to throw are returned
using `ghelpers.GetGErrBlk(excNames.ClassNotLoadedException, errMsg)`
* gfunctions should appear in alphabetical order in the file, with helper functions placed at the end
after a row of `=============`

### Coding conventions for gfunction loaders of MethodSignatures

* gfunctions are loaded into the `MTable` by the `gfunction. MTableLoadGFunctions()`.
* the loader calls the Loadxxx function in every gfunction file.
* the `Loadxxx` function should be named `Load_Package_Class`; for example, `Load_Lang_UTF16()` for java.lang.UTF16
* the `Loadxxx` function loads the gfunctions for methods specified in the list of `MethodSignatures`
* You should create a `MethodSignatures` entry for every method of the class including constructors
* The entry should be a `MethodSignature` struct, which contains the following fields:
    * `name` is the fully qualifed name (FQN) of the function
    * `GMeth` is a struct containing the following fields:
        * `ParamSlots` is an integer indicating the number of parameters passed into the Java function
        * `GFunction` is the name of the gfunction to be called 

* When creating the gfunction for a new library class, it is preferred to create `MethodSignatures` for all
non-private methods. Methods that are not implemented should call `ghelpers.TrapFunction` The complete list
of possible trap functions includes: 
* `ghelpers.TrapClass` - The entire class is not yet supported.
* `ghelpers.TrapDeprecated` - Deprecated classes, interfaces, and functions are not supported.
* `ghelpers.TrapFunction` - This individual function is not yet supported.
* `ghelpers.TrapProtected` - Protected functions are not supported.
* `ghelpers.TrapUndocumented` - Undocumented (hidden) functions are not supported.
* `ghelpers.TrapUnicode` - References to unicode are not yet supported.

When in doubt, use `ghelpers.TrapFunction`. All entries in `MethodSignatures` should be in alphabetical order by 
the FQN of the Java method.



### Coding conventions for unit tests
* Tests that involve any variable or reference outside the function under test should start with
a call to `globals.InitGlobals("test")`, which initializes many global variables required
for almost any action outside the function under test.
* If it's necessary to capture stderr, the code to be used is shown 
[here](https://github.com/platypusguy/jacobin/blob/main/notes/StderrInTests.txt). This
code enables the stderr messages to be viewed as a single golang string.
* As a rule, we do not use tables in unit tests. An exception could be in saving on unit code by driving a function multiple times with parameters.
* All unit tests should be in alphabetical order in the file; helper functions should be
segregated at the top or bottom of the file with conspicuous comments.
* All unit tests should be preceded with a comment of 1 or more lines explaining the test.
If a specific condition is being tested, it should be mentioned, so that we can quickly
see which conditions have/have not been tested without reading the code.

### General coding conventions
* All comments should be in English.
* We follow the `go fmt` coding style.
* All functions should have a preceding comment that explains what they do. Each 
file should start with an introductory block comment explaining what functionality it
implements, so that we can quickly tell whether the file is relevant to our present
search.
* We generally refer to the Go language as "golang," because this facilitates search

## AI Assistant Guidelines
*   **No Conversational Fillers:** Do not include summaries, "I have finished", or "Here is the code" banners in the chat response. Just perform the action.
*   **No Header Comments:** Do not add block comments summarizing the file or function at the top of files unless specifically requested.
*   **Concise Output:** Keep chat responses to the absolute minimum required to confirm the action.
