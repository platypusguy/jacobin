# Jacobin VM Guidelines for AI Agents

### Project Overview
Jacobin is a Go-based Java Virtual Machine (JVM). Its primary objective is correctness, specification compliance, and high runtime efficiency without reliance on native JNI overhead where Go-implemented standard library functions (`gfunctions`) can be used.

---

### Core Runtime Architecture & Invariants

#### 1. Uniform 64-Bit Call-Stack Slot Layout
- Every stack entry and local variable slot in Jacobin is strictly 64 bits:
  - Reference pointer (`*object.Object` or `nil`)
  - 64-bit integer (`int64`)
  - 64-bit floating point (`float64`)
- **Single-Slot Convention**: Unlike HotSpot, Jacobin allocates **exactly one parameter slot** for `long` (`J`) and `double` (`D`) primitive arguments in method parameter lists. When registering `ParamSlots` in `ghelpers.MethodSignatures`, count `J` and `D` as 1 slot each.

#### 2. G-Function Return & Error Contract
Every G-function must strictly adhere to the two-way return convention:
- **Success**: Return the type matching the descriptor in the method signature:
  - Reference types (`L...;`, `[`): `*object.Object` (or `nil` for Java `null`).
  - Integer / Boolean primitives (`Z`, `B`, `C`, `S`, `I`, `J`): 64-bit integer (`int64`, `types.JavaBoolTrue`, `types.JavaBoolFalse`).
  - Floating-point primitives (`F`, `D`): `float64` / `float32`.
  - `void` (`V`): `nil` or `ghelpers.VoidReturn()`.
- **Error / Exception**: Return `*ghelpers.GErrBlk` created via `ghelpers.GetGErrBlk(exceptionType, message)` to allow the runtime dispatcher to throw the appropriate JVM exception.

#### 3. Java Arithmetic Bit-Width Semantics
While values are passed as 64-bit entries, all arithmetic operations must strictly obey JVM bit-width truncation and wrapping rules:
- `types.JavaByte` / `byte`: 8-bit signed (`int8`)
- `short` / `char`: 16-bit signed / unsigned (`int16` / `uint16`)
- `int`: 32-bit signed (`int32`)
- `long`: 64-bit signed (`int64`)
- `float`: 32-bit IEEE 754 precision (`float32`)
- `double`: 64-bit IEEE 754 precision (`float64`)

#### 4. Java Boolean Conventions
- Always return `types.JavaBoolTrue` (`int64(1)`) or `types.JavaBoolFalse` (`int64(0)`).
- Never return native Go `bool` from G-functions.
- Incoming boolean parameters must be extracted as 64-bit integers using `ghelpers.GetInt64Val(arg)`.

#### 5. Static Initialization
- Class static fields and formatter caches should be initialized via `<clinit>()V` method signatures to ensure deterministic loading and re-initialization across runtime states.

---

### Package Structure & Responsibilities

- `src/gfunction/`: Standard library method table loaders and execution dispatchers.
  - `src/gfunction/ghelpers/`: Universal G-function helpers, conversions, and the platform-neutral JDK TZDB zone helpers (`zoneHelpers.go`).
  - `src/gfunction/javaTime/`: Implementation of `java.time.*` and `java.time.zone.*` classes.
  - `src/gfunction/javaUtil/`: Collections, utilities, and currency/locale support.
  - `src/gfunction/javaLang/`: Core language primitives, reflection, and system classes.
  - `src/gfunction/javaMath/`: Arbitrary precision and mathematical operations (`BigInteger`, `BigDecimal`).
- `src/object/`: In-memory object representation, field tables, and string object helpers.
- `src/types/`: System-wide constants, primitive type indicators, and boolean representations.
- `src/statics/`: Static field storage and global tables.
- `src/excNames/`: Standard JVM exception FQNs.

---

### Verification and Testing Standards

- **Unit Tests**: Every new G-function or module must be accompanied by unit tests covering:
  - Normal execution paths and return values.
  - Edge cases (null inputs, boundary values, arithmetic overflow).
  - Exception and error block returns (`*ghelpers.GErrBlk`).
- **Running Tests**:
  - Run package-specific tests: `go test -v ./src/gfunction/javaTime/...`
  - Run full test suite: `go test ./...`
  - Cross-platform check (Windows compatibility): `GOOS=windows go test -c ./src/gfunction/...`
- **Zero Regressions**: All existing tests must pass before concluding any modification.
