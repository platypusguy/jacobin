## Welcome to Jacobin JVM

Jacobin is an implementation of the [JVM specification for Java 21](https://docs.oracle.com/javase/specs/jvms/se21/html/). It is written entirely in Go with no non-Go dependencies. 

The goal is to provide a more-than-minimal implementation of the JVM that can run most class files and JARs and deliver the same results as the OpenJDK-based JVMs (that is, the majority of JVM implementations today). A paramount consideration in the design and implementation of Jacobin is the codebase: making it cohesive and containing clear code. The cohesiveness, extensive commenting, and large test suite enable professionals who want to know more about how the JVM works to find the information quickly and in an easily accessible setting. Additional information on the [Jacobin wiki](https://github.com/platypusguy/jacobin/wiki/Jacobin-Documentation-Home) provides more background and insight. 

Eventually, we hope to enable developers to run their code and have a UI that shows the bytecodes executing, along with the contents of the various JVM stacks, heap, and counters. If you want to see how your Java code executes, this end product will show you in real time. (This is difficult to do with OpenJDK-based JVMs because of the size of the codebase, which is written in multiple languages.)


Due to our desire for a strong, reliable product, Jacobin is heavily tested during development. As of August 2025, the test code is 208% the size of the production code and consists of more than 1,710 tests. We're committed to increasing these numbers. When Jacobin advances some more, we intend to run the OpenJDK test suites against it. 

### Current Status

The current status is shown [here](https://github.com/platypusguy/jacobin). Updates are also posted in realtime on the [Jacobin Twitter account](https://twitter.com/jacobin_jvm).There are currently no packaged releases of Jacobin available (although you can always compile the code). We'll issue releases when Jacobin is mature enough for beta testers to be able to run classes as expected.

At present, all tasks and defects are logged in an instance of JetBrains' [YouTrack](https://www.jetbrains.com/youtrack/) (kindly provided at no cost). The task numbers appear at the start of the comment for every commit and push. The GitHub 'issues' facility is used strictly for issues posted by users. This design allows users to find solutions without needing to dig through numerous unrelated matters. 

### Contents

As we progress, we post short explanations of project decisions and explanations of how the JVM works. Current material can be found below:

#### Project Posts

[Jacobin after 4 years](https://binstock.blogspot.com/2025/09/jvm-at-4-year-mark-we-just-completed.html) | [3.5 years](https://binstock.blogspot.com/2025/03/jacobin-jvm-after-3-of-development-on.html) | [3 years](https://binstock.blogspot.com/2024/09/at-3-year-mark-brief-look-at-history.html) | 
[2.5 years ]( http://binstock.blogspot.com/2024/02/jacobin-jvm-at-30-months.html) | [2years](http://binstock.blogspot.com/2023/08/jacobin-at-2-year-mark.html) | [18 months](http://binstock.blogspot.com/2023/02/jacobin-jvm-at-18-months.html) | [after 1 year](http://binstock.blogspot.com/2022/08/jacobin-jvm-at-1-year-mark.html)

[Jacobin update log](https://github.com/platypusguy/jacobin/wiki/Update-Log)

[Why was Go chosen for this project?](http://binstock.blogspot.com/2021/08/a-whole-new-project-jvm.html)

#### How the Jacobin JVM works
[Command-line processing](https://github.com/platypusguy/jacobin/wiki/Command-line-Processing)

[Inside Java class files: the constant pool](https://blogs.oracle.com/javamagazine/post/java-class-file-constant-pool)

[Inside the JVM: Arrays](https://blogs.oracle.com/javamagazine/post/java-array-objects)

[Inside the JVM: How fields are handled](https://github.com/platypusguy/jacobin/wiki/How-Fields-are-Handled-in-the-JVM)

### The Team (and Thanks)
Jacobin is presently being developed by Andrew Binstock ([platypusguy](https://github.com/platypusguy/)) and [Richard Elkins](https://twitter.com/texadactyl), with significant early contributions by [Spencer Uresk](https://twitter.com/suresk). Contributors are more than welcome and on our [wiki](https://github.com/platypusguy/jacobin/wiki/Jacobin-Documentation-Home), we've provided copious information on the project structure to help contributors get aboard. If you can't contribute code but would like to show support, we'd love a GitHub star or for you to follow the project. 

This project could not have been possible without Github (for the excellent platform), [JetBrains](https://www.jetbrains.com/go/) (for superb tools), Oracle's Java team (for the great technology and [best-in-class documentation](https://docs.oracle.com/javase/specs/index.html)), and these JVM experts: [Ben Evans](https://github.com/kittylyst), [Aleksey Shipilev](https://shipilev.net/), [Chris Newland](https://github.com/sponsors/chriswhocodes), and [Bill Venners](https://github.com/bvenners) who have written helpful, in-depth articles on the machinery of the JVM. A big thanks to all!
