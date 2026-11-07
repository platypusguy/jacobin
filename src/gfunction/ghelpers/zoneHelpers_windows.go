//go:build windows

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package ghelpers

// GetZoneIDs is currently a Windows placeholder and returns a list of a single zone ID.
func GetZoneIDs() ([]string, error) {
	return []string{"Europe/London"}, nil
}

// GetTZDBVersion is currently a Windows placeholder and returns a version = 1.
func GetTZDBVersion() (string, error) {
	return "1", nil
}

// GetZoneOffsets is currently a Windows placeholder and returns zero offsets.
func GetZoneOffsets(string, int64) (standardOffsetMillis, totalOffsetMillis int64, err error) {
	return 0, 0, nil
}
