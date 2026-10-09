//go:build !windows

package javaUtil

// _getDefaultLanguageFromWindows is never called off Windows; it exists so the
// package compiles on every platform.
func _getDefaultLanguageFromWindows() string {
	return "Nonsense!"
}
