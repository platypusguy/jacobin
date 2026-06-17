## INVOKEDYNAMIC implementation status

jvm.interpret() -> jvm.doInvokeDynamic() -> 
mhResoution.go::jvm.resolveCallSite() does:
1. Fetch the InvokeDynamic entry -- OK
2. Get the Bootstrap Method info from the class attributes -- OK
3. Resolve the Bootstrap Method Handle -> classloader.ResolveMethodHandle() 



## How INOVKEDYNAMIC works

INVOKEDYNAMIC (opcode 186 or 0xBA) is a 5-byte instruction for dynamic method invocation, introduced in Java 7 
(JSR 292) to support dynamically typed languages and features like lambdas.

After the class file is parsed (during class loading and linking), the JVM handles it via the run-time constant pool 
and lazy linking/resolution on first execution. Here's the detailed process:

### 1. Class File Representation (Parsed at Load Time)
- The instruction in bytecode: `invokedynamic indexbyte1 indexbyte2 0 0`.
- The 2-byte index points to a `CONSTANT_InvokeDynamic_info` entry in the constant pool.
- This structure contains:
    - `bootstrap_method_attr_index`: Index into the class's `BootstrapMethods` attribute (which lists bootstrap method specifiers).
    - `name_and_type_index`: Points to a `CONSTANT_NameAndType_info` with the *dynamic invocation name* (e.g., 
  a method name like "lambda$...") and *method descriptor* (e.g., `(I)Ljava/lang/String;`).
- The `BootstrapMethods` attribute provides:
    - A `CONSTANT_MethodHandle_info` for the bootstrap method (BSM).
    - Optional static bootstrap arguments (loadable constants: strings, classes, primitives, method handles, method types, etc.).

These are turned into entries in the **run-time constant pool** (one per class/interface). Each `invokedynamic` site is a unique *dynamic call site*.

### 2. Initial State: Unlinked Call Site
- Each lexical `invokedynamic` occurrence represents a *dynamic call site*.
- Initially, it is **unlinked** — no target method is bound. The JVM does not know what code to execute.

### 3. Execution and Linking (First Invocation)
When the JVM executes `invokedynamic` for the first time (or if the call site becomes invalid):

- *Resolve the call site specifier* (dynamically-computed call site) from the run-time constant pool. This triggers §5.4.3.6 in the JVMS.
  - *Resolve necessary constants* (bootstrap method handle, name, method type, static arguments).
      - Resolve the Bootstrap Method Handle (BSM)
      - The JVM looks at the bootstrap_method_attr_index in the CONSTANT_InvokeDynamic_info entry. This index points to the `BootstrapMethods` attribute table for the class.
      - From that table, it gets a `CONSTANT_MethodHandle_info` entry.
      - The JVM then fully resolves this method handle. This is a recursive process that involves finding the actual method the handle points to (which must be a static method) and ensuring it is accessible. If this fails, a `BootstrapMethodError` is thrown.
  - Resolve the Static Arguments
      -  The `BootstrapMethods` attribute also contains a list of constant pool indices for any static arguments that need to be passed to the bootstrap method.
        -  The JVM resolves each of these arguments into live Java objects. This is a critical step where the LDC instruction's logic is reused internally:
          - A `CONSTANT_String_info` becomes a `java.lang.String` object.
          - A `CONSTANT_Class_info` becomes a `java.lang.Class` object.
          - A `CONSTANT_Integer_info` becomes an `int` (or `Integer`).
          - A `CONSTANT_MethodType_info` is resolved into a `java.lang.invoke.MethodType` object.
          - A `CONSTANT_MethodHandle_info` is resolved into a `java.lang.invoke.MethodHandle` object.
        - These resolved objects are prepared to be passed as arguments to the BSM.
          - 
  - Prepare the Dynamic Arguments for the BSM
    - The JVM prepares three special arguments that are always passed to the bootstrap method, before any of the static arguments.
    a. MethodHandles.Lookup: An object that represents the access rights and context of the calling class. This is crucial for security and encapsulation, as it allows the BSM to look up methods only with the permissions of the code that contains the invokedynamic instruction.
    b. String: The name of the method being invoked (e.g., "apply" for a functional interface, or "myDynamicMethod"). This is retrieved from the `CONSTANT_NameAndType_info` entry associated with the invokedynamic instruction.
    c. MethodType: An object representing the signature of the call (the parameter types and return type). This is also resolved from the `CONSTANT_NameAndType_info` entry.
    - 
- Invoke the bootstrap method (BSM):
    - The BSM is called as if via `MethodHandle.invoke` (or `invokeWithArguments` in some descriptions).
    - Arguments passed to the BSM (in order):
        1. `MethodHandles.Lookup` — a lookup object with access rights of the caller (for reflective access control).
        2. `String` — the name from the NameAndType.
        3. `MethodType` — the method descriptor from the NameAndType.
        4. Any static bootstrap arguments from the `BootstrapMethods` entry.
      5. 
    - The BSM is user-provided except for lambdas (For lambdas, this is where `LambdaMetafactory.metafactory` runs, dynamically generating a new class that implements the functional interface.) It runs the Java code to decide the target.
    - This is a standard Java method call. The BSM can now execute any Java code it needs to determine what method should ultimately be called. For lambdas, this is where `LambdaMetafactory.metafactory` runs, dynamically generating a new class that implements the functional interface.
- BSM returns a `java.lang.inovke.CallSite` object (or throws a `BootstrapMethodError` exception, causing linking failure).
    - Common implementations: `ConstantCallSite` (immutable target), `MutableCallSite` (changeable target), `VolatileCallSite`.
    - The `CallSite` object contains a crucial piece of information: its target MethodHandle. This handle points to the actual code that should be executed for this `invokedynamic` instruction from now on.

- *Bind ("link") the call site*: The JVM associates the returned `CallSite` (and its target `MethodHandle`) with this specific `invokedynamic` instruction. This binding is typically stored in a per-call-site data structure in the JVM's internal state (e.g., in HotSpot, involving call site objects and possibly JIT optimizations).


### 4. Subsequent Invocations (Linked State)
- The JVM extracts the target `MethodHandle` from the bound `CallSite`.
- It invokes the `MethodHandle` (as if via `invokeExact` or equivalent), passing the arguments that were on the operand stack for the `invokedynamic`.
- Operand stack effect: Pops the arguments `[arg1, arg2, ...]`, invokes the target, and pushes the result (if any). The call site's method type must match the arguments and return type.
- No re-invocation of the BSM unless the call site is invalidated (e.g., via `MutableCallSite` or class redefinition).

### 5. Key Implementation Details and Optimizations
- Method Handles: These are the core of the dynamic behavior. They support direct invocation, argument transformation (e.g., `asType`, `insertArguments`, `filterArguments`), and are highly optimizable by the JIT (often inlined to near-native performance after warmup).
- CallSite stability: Once linked, performance can match `invokevirtual`/`invokestatic` thanks to JIT (e.g., monomorphic call site optimization in HotSpot).
- Lambda example: For Java lambdas, the compiler emits `invokedynamic` with a BSM that points to `LambdaMetafactory.metafactory` (or `altMetafactory`). This factory creates a `CallSite` whose target is a synthetic method or proxy implementing the functional interface.
- Thread safety and concurrency: Linking is atomic per call site; multiple threads may trigger the BSM, but the JVM ensures it runs at most once successfully per site (with synchronization).
- Dynamic updates: Using `MutableCallSite` or `VolatileCallSite` allows the target to change at runtime (e.g., for language runtimes with evolving semantics). Changing the target can cause deoptimization in the JIT.

### 6. Errors and Edge Cases
- Bootstrap method throws → wrapped in `BootstrapMethodError`.
- Wrong return type from BSM (must be `CallSite` subtype) → linking failure.
- Type mismatches on invocation → `WrongMethodTypeException` (from `MethodHandle`).
- The two zero bytes in the instruction are reserved (for future extensions).

In summary, after parsing, `INVOKEDYNAMIC` turns static bytecode into a **pluggable, lazily-resolved invocation point** where a user-defined bootstrap method decides (and can later change) the exact behavior. This is what enables flexible dynamic languages and efficient higher-order features on the JVM. For the authoritative details, consult the latest JVMS sections on `invokedynamic` (6.5) and dynamically-computed call site resolution (5.4.3.6).
