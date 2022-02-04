## Welcome to Jacobin JVM

Jacobin is an implementation of the [JVM specification for Java 11](https://docs.oracle.com/javase/specs/jvms/se11/html/) as published by Oracle Corp. and filled out by numerous articles and technical reports on the workings of the machinery in the JVM. It is written entirely in Go. 

The goal is to provide a more-than-minimal implementation of the JVM that can run most class files and JARs and deliver the same results as the OpenJDK-based JVMs (that is, the majority of JVM implementations today). A paramount consideration in the design and implementation of Jacobin is the codebase: making it cohesive and containing clear code. The former aspect distinguishes it from the OpenJDK implementations: those codebases are spread out across a vast set of different code repositories, some VM-exclusive, some from the compiler, some from other Java tools and resources--making it very difficult to navigate the codebase unless you already know what you're looking for and how it's been implemented. Because Jacobin is strictly a JVM, its code is smaller and tightly focused on Java program execution. An additional factor in reducing the size of the codebase is that Jacobin relies on Go's built-in memory management to perform garbage collection and so it contains no GC code, which is a substantial part of the OpenJDK implementations.

Jacobin is heavily tested during development. As of February 2022, the testing code base is 2.1x the size of the production code. We're committed to increasing that ratio. When Jacobin advances some more, we intend to run the OpenJDK test suites against it. 

### Current Status

The current status is shown [here](https://github.com/platypusguy/jacobin). Updates are also posted in realtime on the [Jacobin Twitter account](https://twitter.com/jacobin_jvm).There are currently no packaged releases of Jacobin available (although you can always compile the code). We'll issue releases when Jacobin is mature enough to run classes as expected.

At present, all tasks and defects are logged in and instance of JetBrains' [YouTrack](https://www.jetbrains.com/youtrack/) (kindly provided at no cost). The task numbers appear at the start of the comment for every commit and push. The GitHub 'issues' facility is used strictly for issues posted by users. This design allows users to find solutions without needing to dig through numerous unrelated matters. 

### Contents

As we progress, we post short explanations of project decisions and explanations of how the JVM works. Current material can be found below:

#### Project Posts
[Why was Go chosen for this project?](http://binstock.blogspot.com/2021/08/a-whole-new-project-jvm.html)

[Detailed Jacobin status after 6 months](http://binstock.blogspot.com/2022/02/jacobin-jvm-project-after-six-months.html)

#### How the Jacobin JVM works
[Command-line processing](https://github.com/platypusguy/jacobin/wiki/Command-line-Processing)

[How Jacobin searches for methods](http://binstock.blogspot.com/2021/12/how-jacobin-jvm-accesses-methods.html)

### The Team (and Thanks)
Jacobin is presently being developed by Andrew Binstock ([platypusguy](https://github.com/platypusguy/)). Contributors are more than welcome. If you'd like to show your support the project but can't contribute code, we'd love a GitHub star or for you to follow the project. 

This project could not have been possible without Github (for the excellent platform), [JetBrains](https://www.jetbrains.com/go/) (for superb tools), Oracle's Java team (for the great technology and [best-in-class documentation](https://docs.oracle.com/javase/specs/index.html)), and these JVM experts: [Ben Evans](https://github.com/kittylyst), [Aleksey Shipilev](https://shipilev.net/), [Chris Newland](https://github.com/sponsors/chriswhocodes), and [Bill Venners](https://github.com/bvenners) who have written helpful, in-depth articles on the machinery of the JVM. A big thanks to all!
