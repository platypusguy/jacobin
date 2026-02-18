## The basic layout of the Jacobin JVM and the coding conventions


### Jacobin JVM Startup Sequence

Jacobin starts up in jvm/jvmStart.go

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

### Code layout

main.go contains the main function. All other source code is located in the `src/` directory.
The main packages/directories are:
* `classloader` contains all the class parsing and loading logic
* `gfunction` and its subdirectories contain all the gfunctions, which implement Java library methods
* `jvm` is the heart of Jacobin: it contains the interpreter and numerous run utitlities
* `object` contains a definition 
* `types` contains definitions of various types as well as a variety of constants

### Coding Conventions

