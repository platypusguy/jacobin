# Windows Path Cheat Sheet (Java NIO / Jacobin)

This cheat sheet summarizes **Windows path semantics** for Java NIO (`Path`) — both **Jacobin** (HotSpot-emulated behavior) and HotSpot-compliant rules.

---

## Path Categories

| Category | Example | isAbsolute() | getRoot() | Notes |
|----------|---------|--------------|-----------|-------|
| Absolute drive path | `C:\foo\bar` | true | `C:\` | Fully-qualified absolute path with drive letter. |
| Drive-only path | `C:` | false | `C:` | Refers to current directory on drive `C:`. Not absolute. |
| Drive-relative path | `C:foo` | false | `C:` | Relative to current directory on drive `C:`. |
| Drive-rooted but drive-unspecified | `\foo` | false | `<null>` (Jacobin) / `\` (HotSpot) | Rooted on current drive. HotSpot returns `\` for getRoot(). |
| UNC path | `\\server\share\dir\file.txt` | true | `\\server\share\` | Network path. Root includes server and share. |
| Relative path | `foo\bar` | false | `<null>` | Relative to current working directory. |
| Mixed separators | `C:/a/b\c` | true | `C:\` | Java NIO normalizes mixed `/` and `\` to `\`. |

---

## Rules

1. **Absolute paths**
   * Must include a **drive letter** or **UNC prefix** to be absolute.  
   * Example: `C:\path` or `\\server\share\path`.

2. **Drive-relative paths**
   * Format: `C:` or `C:foo`.  
   * **Not absolute**, root is the drive letter.

3. **Drive-rooted but drive-unspecified**
   * Format: `\foo`.  
   * Refers to the current drive’s root.  
   * **Jacobin:** `getRoot() = null`  
   * **HotSpot:** `getRoot() = "\"`

4. **UNC paths**
   * Start with `\\server\share`.  
   * Always absolute.  
   * Root includes server and share: `\\server\share\`.

5. **Relative paths**
   * Do not start with a drive or `\`.  
   * Root is `<null>`.  
   * Resolved relative to current working directory.

6. **Mixed separators**
   * Windows accepts both `/` and `\`.  
   * Java NIO **normalizes** all to `\`.

---

## Examples

| Path | Normalized | isAbsolute() | getRoot() | NameCount | Name(0) |
|------|------------|--------------|-----------|-----------|---------|
| `C:\foo\bar` | `C:\foo\bar` | true | `C:\` | 2 | `foo` |
| `C:` | `C:` | false | `C:` | 0 | N/A |
| `C:foo` | `C:foo` | false | `C:` | 1 | `foo` |
| `\foo` | `\foo` | false | `<null>` / `\` | 1 | `foo` |
| `\\server\share\dir\file.txt` | `\\server\share\dir\file.txt` | true | `\\server\share\` | 2 | `dir` |
| `foo\bar` | `foo\bar` | false | `<null>` | 2 | `foo` |
| `C:/a/b\c` | `C:\a\b\c` | true | `C:\` | 3 | `a` |

---

## Notes

* `\foo` is a **critical edge case**: HotSpot and Jacobin differ slightly in `getRoot()`.  
* Use `Path.resolve()` and `Path.relativize()` carefully with drive-relative paths.  
* `Path.normalize()` removes `.` and resolves `..` segments.  
* UNC paths are always absolute; mixing UNC with relative segments follows standard normalization rules.  

---

> This cheat sheet reflects **Java NIO spec** semantics and is useful for building **HotSpot-faithful path emulators** like Jacobin.

