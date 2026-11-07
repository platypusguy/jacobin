//go:build windows

package ghelpers

import "testing"

func TestWindowsZoneHelpersReturnDummyValues(t *testing.T) {
	ids, err := GetZoneIDs()
	if err != nil || len(ids) != 0 {
		t.Fatalf("GetZoneIDs() = (%v, %v), want empty IDs and nil error", ids, err)
	}
	version, err := GetTZDBVersion()
	if err != nil || version != "" {
		t.Fatalf("GetTZDBVersion() = (%q, %v), want empty version and nil error", version, err)
	}
	standard, total, err := GetZoneOffsets("Etc/UTC", 0)
	if err != nil || standard != 0 || total != 0 {
		t.Fatalf("GetZoneOffsets() = (%d, %d, %v), want zero offsets and nil error", standard, total, err)
	}
}
