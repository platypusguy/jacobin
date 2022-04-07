What is this?

This directory contains .java source files, as well as their binary output, necessary for testing the module reader.

* classes/ contains the java source files and their .class counterparts. You can regenerate them with: `javac org/jacobin/test/Hello.java module-info.java`.
* lib/ contains the classlist, which limits the classes that are loaded in the module.
* jacobin.jmod is the output module. You can regenerate it with `jmod create --class-path classes/ --libs lib jacobin.jmod`.
* jacobinfull.jmod is the output module that has no classlist, for testing purposes. You can regenerate it with `jmod create --class-path classes/ jacobinfull.jmod`.
