# Windows Path Cheat Sheet (Jacobin / Java NIO Spec)

| Path string                  | isAbsolute() | getRoot()           | Notes / Name decomposition                     |
|-------------------------------|--------------|-------------------|-----------------------------------------------|
| `C:\foo\bar`                  | true         | `C:\`             | Drive-absolute. Name count: 2 (`foo`, `bar`) |
| `C:`                          | false        | `C:`              | Drive-relative, empty path. Name count: 0     |
| `C:foo`                       | false        | `C:`              | Drive-relative. Name count: 1 (`foo`)        |
| `\foo\bar`                     | true         | `\`               | Rooted at current drive. Name count: 2 (`foo`, `bar`) |
| `\\server\share\dir\file.txt` | true         | `\\server\share\` | UNC path. Name count: 2 (`dir`, `file.txt`) |
| `foo\bar`                      | false        | null              | Relative path. Name count: 2 (`foo`, `bar`) |
| `.`                            | false        | null              | Current directory reference. Name count: 1 (`.`) |
| `..`                           | false        | null              | Parent directory reference. Name count: 1 (`..`) |
| `C:\foo\..\bar`                | true         | `C:\`             | Normalize → `C:\bar`                          |
| `C:/foo/bar`                   | true         | `C:\`             | Mixed separators normalize → `C:\foo\bar`    |
| `C:\`                          | true         | `C:\`             | Root of drive. Name count: 0                  |
| `\`                            | true         | `\`               | Root of current drive. Name count: 0          |
| `\\server\share\`              | true         | `\\server\share\` | UNC root only. Name count: 0                 |
| `\\server\share`               | true         | `\\server\share\` | UNC root without trailing `\`. NIO adds it.  |
| `C:\.`                         | true         | `C:\`             | Normalizes to `C:\`                           |
| `C:\..\foo`                     | true         | `C:\`             | Normalizes to `C:\foo`                        |

## Key Takeaways

1. **Drive-relative vs absolute**  
   * `C:foo` → relative, root = `"C:"`  
   * `C:\foo` → absolute, root = `"C:\"`  

2. **Root-only paths**  
   * `\foo` → absolute, root = `\`  
   * `\` → absolute, root = `\`  

3. **UNC paths**  
   * Root includes `\\server\share\`  
   * Name elements start after the share  

4. **Relative paths**  
   * No drive, no leading `\` → `isAbsolute() = false`  
   * `getRoot() = null`  

5. **Normalization**  
   * Both mixed separators `/` and `\` are normalized to `\`  
   * `.` and `..` handled according to Java spec

