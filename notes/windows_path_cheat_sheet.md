# Windows Path Cheat Sheet (Java NIO / HotSpot & Jacobin)

This cheat sheet summarizes **Windows path semantics** for Java NIO (`Path`) on Windows.  
Jacobin emulates HotSpot faithfully, which internally behaves in a **Unix-like manner**, even on Windows.  

---

## Path Hierarchy vs Root

* **Hierarchy**: Paths are treated as sequences of name elements, separated by `/` internally.  
* **Root**: Windows-specific concept — drive letter, UNC prefix, or `\` for drive-rooted paths.  

| Path | Hierarchy (name elements) | Root | Notes |
|------|--------------------------|------|-------|
| `C:\foo\bar` | `[foo, bar]` | `C:\` | Absolute path on drive C. |
| `C:` | `[]` | `C:` | Refers to current directory on C. Not absolute. |
| `C:foo` | `[foo]` | `C:` | Drive-relative path, not absolute. |
| `\foo\bar` | `[foo, bar]` | `\` | Drive-rooted, current drive. `isAbsolute() = false`. |
| `\\server\share\dir\file.txt` | `[dir, file.txt]` | `\\server\share\` | UNC path, absolute. |
| `foo\bar` | `[foo, bar]` | `<null>` | Pure relative path. |
| `C:/a/b\c` | `[a, b, c]` | `C:\` | Mixed separators are normalized to `\`. |

---

## Key Rules

1. **Absolute paths**
   * Must include **drive letter** or **UNC prefix**.  
   * Example: `C:\path` or `\\server\share\path`.

2. **Drive-relative paths**
   * Format: `C:` or `C:foo`.  
   * **Not absolute**, root is drive letter.  

3. **Drive-rooted but drive-unspecified**
   * Format: `\foo`.  
   * Rooted on the current drive.  
   * `isAbsolute() = false`, `getRoot() = "\"`.  

4. **UNC paths**
   * Always absolute; root includes server and share: `\\server\share\`.  

5. **Relative paths**
   * Do not start with a drive or `\`.  
   * Root is `<null>`.  

6. **Mixed separators**
   * Both `/` and `\` are accepted on Windows.  
   * `Path.normalize()` converts all separators to `\`.

---

## Quick Examples

| Path | Normalized | isAbsolute() | getRoot() | NameCount | Name(0) |
|------|------------|--------------|-----------|-----------|---------|
| `C:\foo\bar` | `C:\foo\bar` | true | `C:\` | 2 | `foo` |
| `C:` | `C:` | false | `C:` | 0 | N/A |
| `C:foo` | `C:foo` | false | `C:` | 1 | `foo` |
| `\foo` | `\foo` | false | `\` | 1 | `foo` |
| `\\server\share\dir\file.txt` | `\\server\share\dir\file.txt` | true | `\\server\share\` | 2 | `dir` |
| `foo\bar` | `foo\bar` | false | `<null>` | 2 | `foo` |
| `C:/a/b\c` | `C:\a\b\c` | true | `C:\` | 3 | `a` |

---

## Notes

* **HotSpot and Jacobin**: Both treat the **path hierarchy as Unix-like**, but root handling is Windows-aware.  
* `Path.resolve()` and `Path.relativize()` work using the **internal hierarchy**, independent of OS separators.  
* Edge case `\foo` is **drive-rooted, not absolute**, but `getRoot()` is `\`.  
* UNC paths always include the server and share as root.  
* Mixed separators normalize to `\` on Windows.

> This cheat sheet is HotSpot-faithful and suitable for developers implementing or testing Windows path handling in Java or Jacobin.

