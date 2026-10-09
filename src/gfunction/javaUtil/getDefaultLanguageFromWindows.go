//go:build windows

package javaUtil

import (
	"strings"
	"syscall"
	"unsafe"
)

// LOCALE_NAME_MAX_LENGTH from winnls.h, includes the terminating null.
const localeNameMaxLength = 85

// _getDefaultLanguageFromWindows returns the user's default locale in Java form
// (e.g. en_US), or "en_US" if the Win32 call fails.
func _getDefaultLanguageFromWindows() string {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("GetUserDefaultLocaleName")

	buf := make([]uint16, localeNameMaxLength)
	// Returns the number of UTF-16 units written, including the null, or 0 on failure.
	n, _, _ := proc.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if n == 0 {
		return "en_US"
	}

	name := syscall.UTF16ToString(buf) // BCP 47 tag, e.g. en-US
	if name == "" {
		return "en_US"
	}
	return strings.ReplaceAll(name, "-", "_")
}
