//go:build !windows

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package ghelpers

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"jacobin/src/globals"
)

func TestIANAZoneHelpersFixture(t *testing.T) {
	fixture := makeTestTZDB()
	db, err := parseIANATZDB(fixture)
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	if db.version != "test-version" {
		t.Errorf("version = %q, want test-version", db.version)
	}
	if len(db.ids) != 2 || db.ids[0] != "Test/Alias" || db.ids[1] != "Test/Zone" {
		t.Errorf("IDs = %v, want [Test/Alias Test/Zone]", db.ids)
	}

	tests := []struct {
		name     string
		instant  time.Time
		standard int64
		total    int64
	}{
		{"before standard transition", time.Date(2019, time.December, 31, 23, 59, 59, 0, time.UTC), 0, 0},
		{"after standard transition", time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), 3600000, 3600000},
		{"historic daylight time", time.Date(2020, time.July, 1, 0, 0, 0, 0, time.UTC), 3600000, 7200000},
		{"historic standard time", time.Date(2020, time.November, 1, 0, 0, 0, 0, time.UTC), 3600000, 3600000},
		{"future daylight time", time.Date(2021, time.July, 1, 0, 0, 0, 0, time.UTC), 3600000, 7200000},
		{"future standard time", time.Date(2021, time.December, 1, 0, 0, 0, 0, time.UTC), 3600000, 3600000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seconds := floorDiv(tt.instant.UnixMilli(), 1000)
			zone := db.zones["Test/Zone"]
			if got := int64(zone.standardOffsetAt(seconds)) * 1000; got != tt.standard {
				t.Errorf("standard offset = %dms, want %dms", got, tt.standard)
			}
			if got := int64(zone.totalOffsetAt(seconds)) * 1000; got != tt.total {
				t.Errorf("total offset = %dms, want %dms", got, tt.total)
			}
			alias := db.zones["Test/Alias"]
			if alias.totalOffsetAt(seconds) != zone.totalOffsetAt(seconds) {
				t.Errorf("alias total offset differs from canonical zone")
			}
		})
	}

	if _, ok := db.zones["Missing/Zone"]; ok {
		t.Fatal("unexpected missing zone mapping")
	}
}

func TestParseIANATZDBRejectsMalformedData(t *testing.T) {
	fixture := makeTestTZDB()
	for _, tt := range []struct {
		name string
		data []byte
	}{
		{"truncated header", fixture[:3]},
		{"truncated rule payload", fixture[:len(fixture)-1]},
		{"unsupported format", append([]byte{2}, fixture[1:]...)},
		{"trailing bytes", append(append([]byte(nil), fixture...), 0)},
		{"invalid mapping index", func() []byte {
			data := append([]byte(nil), fixture...)
			mappingOffset := testTZDBMappingOffset(t, data)
			binary.BigEndian.PutUint16(data[mappingOffset+12:mappingOffset+14], 0xffff)
			return data
		}()},
		{"invalid alias index", func() []byte {
			data := append([]byte(nil), fixture...)
			binary.BigEndian.PutUint16(data[len(data)-2:], 0xffff)
			return data
		}()},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := parseIANATZDB(tt.data); err == nil {
				t.Fatal("parseIANATZDB() unexpectedly succeeded")
			}
		})
	}
}

func TestParseIANATZDBAcceptsMissingAliasTable(t *testing.T) {
	fixture := makeTestTZDB()
	withoutAliases := fixture[:len(fixture)-6]
	db, err := parseIANATZDB(withoutAliases)
	if err != nil {
		t.Fatalf("parse fixture without optional alias table: %v", err)
	}
	if len(db.ids) != 2 {
		t.Errorf("IDs = %v, want both mapped regions", db.ids)
	}
}

func TestParseZoneRulesUncompressedEncodings(t *testing.T) {
	var data bytes.Buffer
	data.WriteByte(zoneRulesType)
	writeInt32(&data, 1)
	writeEpochSecond(&data, -4575744001)
	writeTestOffset(&data, 64799)
	writeTestOffset(&data, 3600)
	writeInt32(&data, 0)
	writeTestOffset(&data, 64799)
	data.WriteByte(0)

	rules, err := parseZoneRules(data.Bytes())
	if err != nil {
		t.Fatalf("parse extended rule encodings: %v", err)
	}
	if rules.standardTransitions[0] != -4575744001 {
		t.Errorf("extended transition = %d, want -4575744001", rules.standardTransitions[0])
	}
	if rules.standardOffsets[0] != 64799 || rules.wallOffsets[0] != 64799 {
		t.Errorf("extended offset = (%d, %d), want (64799, 64799)", rules.standardOffsets[0], rules.wallOffsets[0])
	}
}

func TestRecurringTransitionDefinitionsAndYearBoundary(t *testing.T) {
	localStart := time.Date(2024, time.March, 1, 2, 0, 0, 0, time.UTC).Unix()
	for _, tt := range []struct {
		name       string
		definition uint8
		want       int64
	}{
		{"wall", 0, localStart - 7200},
		{"standard", 1, localStart - 3600},
		{"UTC", 2, localStart},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rule := ianaTransitionRule{
				month: time.March, dayOfMonth: 1, secondOfDay: 2 * 60 * 60,
				timeDefinition: tt.definition, standardOffset: 3600,
				offsetBefore: 7200, offsetAfter: 3600,
			}
			got, err := rule.transition(2024)
			if err != nil || got != tt.want {
				t.Fatalf("transition = (%d, %v), want (%d, nil)", got, err, tt.want)
			}
		})
	}

	crossYearRules := &ianaZoneRules{
		wallOffsets: []int32{3600},
		lastRules: []ianaTransitionRule{
			{month: time.January, dayOfMonth: 1, offsetBefore: 3600, offsetAfter: 0},
			{month: time.July, dayOfMonth: 1, offsetBefore: 0, offsetAfter: 3600},
		},
	}
	beforeNewYear := time.Date(2024, time.December, 31, 22, 30, 0, 0, time.UTC).Unix()
	afterNewYear := time.Date(2024, time.December, 31, 23, 30, 0, 0, time.UTC).Unix()
	if got := crossYearRules.totalOffsetAt(beforeNewYear); got != 3600 {
		t.Errorf("offset before cross-year transition = %d, want 3600", got)
	}
	if got := crossYearRules.totalOffsetAt(afterNewYear); got != 0 {
		t.Errorf("offset after cross-year transition = %d, want 0", got)
	}
}

func TestIANAZoneHelpersMissingTZDB(t *testing.T) {
	isolateIANATZDB(t)
	oldHome := globals.GetSystemProperty("java.home")
	globals.SetSystemProperty("java.home", t.TempDir())
	t.Cleanup(func() { globals.SetSystemProperty("java.home", oldHome) })
	if _, err := GetTZDBVersion(); err == nil {
		t.Fatal("GetTZDBVersion() unexpectedly succeeded without tzdb.dat")
	}
}

func TestIANAZoneOffsetsRejectUnknownZone(t *testing.T) {
	isolateIANATZDB(t)
	oldHome := globals.GetSystemProperty("java.home")
	javaHome := t.TempDir()
	if err := os.MkdirAll(filepath.Join(javaHome, "lib"), 0o755); err != nil {
		t.Fatalf("create fixture JDK lib directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(javaHome, "lib", "tzdb.dat"), makeTestTZDB(), 0o600); err != nil {
		t.Fatalf("write fixture TZDB: %v", err)
	}
	globals.SetSystemProperty("java.home", javaHome)
	t.Cleanup(func() { globals.SetSystemProperty("java.home", oldHome) })
	if _, _, err := GetZoneOffsets("Missing/Zone", 0); err == nil || !strings.Contains(err.Error(), "unknown IANA time-zone ID") {
		t.Fatalf("GetZoneOffsets() error = %v, want unknown-zone error", err)
	}
}

func TestIANAZoneHelpersUseJDKTZDB(t *testing.T) {
	isolateIANATZDB(t)
	javaHome := globals.GetSystemProperty("java.home")
	if javaHome == "" {
		t.Skip("java.home is not configured")
	}
	tzdbPath := filepath.Join(javaHome, "lib", "tzdb.dat")
	if _, err := os.Stat(tzdbPath); os.IsNotExist(err) {
		t.Skipf("JDK TZDB is absent: %s", tzdbPath)
	} else if err != nil {
		t.Fatalf("stat JDK TZDB %s: %v", tzdbPath, err)
	}
	version, err := GetTZDBVersion()
	if err != nil {
		t.Fatalf("JDK TZDB exists but could not be parsed: %v", err)
	}
	if version == "" {
		t.Fatal("GetTZDBVersion() returned an empty version")
	}
	ids, err := GetZoneIDs()
	if err != nil {
		t.Fatalf("GetZoneIDs(): %v", err)
	}
	for _, id := range []string{"America/New_York", "Australia/Sydney"} {
		found := false
		for _, available := range ids {
			if available == id {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("JDK TZDB does not list %q", id)
		}
	}

	checks := []struct {
		id       string
		date     time.Time
		standard int64
		total    int64
	}{
		{"America/New_York", time.Date(2050, time.January, 15, 12, 0, 0, 0, time.UTC), -18000000, -18000000},
		{"America/New_York", time.Date(2050, time.July, 15, 12, 0, 0, 0, time.UTC), -18000000, -14400000},
		{"Australia/Sydney", time.Date(2050, time.January, 15, 12, 0, 0, 0, time.UTC), 36000000, 39600000},
		{"Australia/Sydney", time.Date(2050, time.July, 15, 12, 0, 0, 0, time.UTC), 36000000, 36000000},
	}
	for _, tt := range checks {
		standard, total, err := GetZoneOffsets(tt.id, tt.date.UnixMilli())
		if err != nil {
			t.Fatalf("GetZoneOffsets(%q): %v", tt.id, err)
		}
		if standard != tt.standard || total != tt.total {
			t.Errorf("%s at %s: offsets = (%d, %d), want (%d, %d)",
				tt.id, tt.date.Format(time.RFC3339), standard, total, tt.standard, tt.total)
		}
	}
	if _, _, err = GetZoneOffsets("Not/A_Real_Zone", 0); err == nil {
		t.Error("unknown zone ID unexpectedly succeeded")
	}
}

func TestIANATZDBCacheReuse(t *testing.T) {
	isolateIANATZDB(t)
	javaHome := writeTestJDKTZDB(t)
	oldHome := globals.GetSystemProperty("java.home")
	globals.SetSystemProperty("java.home", javaHome)
	t.Cleanup(func() { globals.SetSystemProperty("java.home", oldHome) })

	version, err := GetTZDBVersion()
	if err != nil || version != "test-version" {
		t.Fatalf("initial GetTZDBVersion() = (%q, %v), want (test-version, nil)", version, err)
	}
	globals.SetSystemProperty("java.home", t.TempDir())
	ids, err := GetZoneIDs()
	if err != nil {
		t.Fatalf("GetZoneIDs() after changing java.home: %v", err)
	}
	if len(ids) != 2 || ids[0] != "Test/Alias" || ids[1] != "Test/Zone" {
		t.Fatalf("GetZoneIDs() = %v, want [Test/Alias Test/Zone]", ids)
	}
}

func TestIANATZDBRetriesAfterLoadFailure(t *testing.T) {
	isolateIANATZDB(t)
	javaHome := t.TempDir()
	if err := os.MkdirAll(filepath.Join(javaHome, "lib"), 0o755); err != nil {
		t.Fatalf("create fixture JDK lib directory: %v", err)
	}
	oldHome := globals.GetSystemProperty("java.home")
	globals.SetSystemProperty("java.home", javaHome)
	t.Cleanup(func() { globals.SetSystemProperty("java.home", oldHome) })

	if _, err := GetTZDBVersion(); err == nil {
		t.Fatal("GetTZDBVersion() unexpectedly succeeded before tzdb.dat was created")
	}
	if err := os.WriteFile(filepath.Join(javaHome, "lib", "tzdb.dat"), makeTestTZDB(), 0o600); err != nil {
		t.Fatalf("write fixture TZDB: %v", err)
	}
	version, err := GetTZDBVersion()
	if err != nil || version != "test-version" {
		t.Fatalf("GetTZDBVersion() after creating tzdb.dat = (%q, %v), want (test-version, nil)", version, err)
	}
}

func TestIANATZDBConcurrentInitialization(t *testing.T) {
	isolateIANATZDB(t)
	javaHome := writeTestJDKTZDB(t)
	oldHome := globals.GetSystemProperty("java.home")
	globals.SetSystemProperty("java.home", javaHome)
	t.Cleanup(func() { globals.SetSystemProperty("java.home", oldHome) })

	type result struct {
		version         string
		standard, total int64
		err             error
	}
	const workers = 32
	start := make(chan struct{})
	results := make(chan result, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			version, err := GetTZDBVersion()
			if err != nil {
				results <- result{err: err}
				return
			}
			standard, total, err := GetZoneOffsets("Test/Zone", time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC).UnixMilli())
			results <- result{version: version, standard: standard, total: total, err: err}
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	for got := range results {
		if got.err != nil {
			t.Errorf("concurrent helper call: %v", got.err)
			continue
		}
		if got.version != "test-version" || got.standard != 3600000 || got.total != 3600000 {
			t.Errorf("concurrent result = (%q, %d, %d), want (test-version, 3600000, 3600000)",
				got.version, got.standard, got.total)
		}
	}
}

func isolateIANATZDB(t *testing.T) {
	t.Helper()
	zoneDbMu.Lock()
	previous := zoneDb
	zoneDb = nil
	zoneDbMu.Unlock()
	t.Cleanup(func() {
		zoneDbMu.Lock()
		zoneDb = previous
		zoneDbMu.Unlock()
	})
}

func writeTestJDKTZDB(t *testing.T) string {
	t.Helper()
	javaHome := t.TempDir()
	libDir := filepath.Join(javaHome, "lib")
	if err := os.MkdirAll(libDir, 0o755); err != nil {
		t.Fatalf("create fixture JDK lib directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(libDir, "tzdb.dat"), makeTestTZDB(), 0o600); err != nil {
		t.Fatalf("write fixture TZDB: %v", err)
	}
	return javaHome
}

func makeTestTZDB() []byte {
	var out bytes.Buffer
	out.WriteByte(tzdbFormatVersion)
	writeUTF(&out, "TZDB")
	writeUint16(&out, 2)
	writeUTF(&out, "old-version")
	writeUTF(&out, "test-version")
	writeUint16(&out, 2)
	writeUTF(&out, "Test/Zone")
	writeUTF(&out, "Test/Alias")
	oldRule := makeFixedTestZoneRule(0)
	rule := makeTestZoneRule()
	writeUint16(&out, 2)
	writeUint16(&out, uint16(len(oldRule)))
	out.Write(oldRule)
	writeUint16(&out, uint16(len(rule)))
	out.Write(rule)
	for ruleIndex := uint16(0); ruleIndex < 2; ruleIndex++ {
		writeUint16(&out, 2)
		for region := uint16(0); region < 2; region++ {
			writeUint16(&out, region)
			writeUint16(&out, ruleIndex)
		}
	}
	writeUint16(&out, 1)
	writeUint16(&out, 1)
	writeUint16(&out, 0)
	return out.Bytes()
}

func makeFixedTestZoneRule(offset int) []byte {
	var out bytes.Buffer
	out.WriteByte(zoneRulesType)
	writeInt32(&out, 0)
	writeTestOffset(&out, offset)
	writeInt32(&out, 0)
	writeTestOffset(&out, offset)
	out.WriteByte(0)
	return out.Bytes()
}

func testTZDBMappingOffset(t *testing.T, data []byte) int {
	t.Helper()
	r := tzdbReader{data: data}
	if _, err := r.readByte(); err != nil {
		t.Fatal(err)
	}
	if _, err := r.readUTF(); err != nil {
		t.Fatal(err)
	}
	versionCount, err := r.readUint16()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < int(versionCount); i++ {
		if _, err := r.readUTF(); err != nil {
			t.Fatal(err)
		}
	}
	regionCount, err := r.readUint16()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < int(regionCount); i++ {
		if _, err := r.readUTF(); err != nil {
			t.Fatal(err)
		}
	}
	ruleCount, err := r.readUint16()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < int(ruleCount); i++ {
		length, err := r.readUint16()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := r.readN(int(length)); err != nil {
			t.Fatal(err)
		}
	}
	return r.off
}

func makeTestZoneRule() []byte {
	var out bytes.Buffer
	out.WriteByte(zoneRulesType)
	writeInt32(&out, 1)
	writeEpochSecond(&out, time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC).Unix())
	writeTestOffset(&out, 0)
	writeTestOffset(&out, 3600)
	writeInt32(&out, 3)
	writeEpochSecond(&out, time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC).Unix())
	writeEpochSecond(&out, time.Date(2020, time.March, 29, 1, 0, 0, 0, time.UTC).Unix())
	writeEpochSecond(&out, time.Date(2020, time.October, 25, 1, 0, 0, 0, time.UTC).Unix())
	writeTestOffset(&out, 0)
	writeTestOffset(&out, 3600)
	writeTestOffset(&out, 7200)
	writeTestOffset(&out, 3600)
	out.WriteByte(2)
	writeTestTransitionRule(&out, 3, -1, 7, 2, 3600, 3600, 7200)
	writeTestTransitionRule(&out, 10, -1, 7, 3, 3600, 7200, 3600)
	return out.Bytes()
}

func writeTestTransitionRule(out *bytes.Buffer, month time.Month, dayOfMonth, dayOfWeek, hour, standard, before, after int) {
	encodeRelativeOffset := func(offset int) uint32 {
		if (offset-standard)%1800 != 0 {
			panic("fixture offset is not a relative half-hour multiple")
		}
		code := (offset - standard) / 1800
		if code < 0 || code > 2 {
			panic("fixture offset is outside compact relative encoding")
		}
		return uint32(code)
	}
	packed := uint32(month)<<28 |
		uint32(dayOfMonth+32)<<22 |
		uint32(dayOfWeek)<<19 |
		uint32(hour)<<14 |
		uint32(standard/900+128)<<4 |
		encodeRelativeOffset(before)<<2 |
		encodeRelativeOffset(after)
	writeInt32(out, int32(packed))
}

func writeTestOffset(out *bytes.Buffer, seconds int) {
	if seconds%900 == 0 && seconds/900 >= -128 && seconds/900 <= 126 {
		out.WriteByte(byte(int8(seconds / 900)))
		return
	}
	out.WriteByte(127)
	writeInt32(out, int32(seconds))
}

func writeEpochSecond(out *bytes.Buffer, seconds int64) {
	const base = int64(4575744000)
	if seconds >= -base && seconds < base+int64(0xffffff)*900 && (seconds+base)%900 == 0 {
		compact := uint32((seconds + base) / 900)
		out.WriteByte(byte(compact >> 16))
		out.WriteByte(byte(compact >> 8))
		out.WriteByte(byte(compact))
		return
	}
	out.WriteByte(0xff)
	writeInt64(out, seconds)
}

func writeUTF(out *bytes.Buffer, value string) {
	writeUint16(out, uint16(len(value)))
	out.WriteString(value)
}

func writeUint16(out *bytes.Buffer, value uint16) {
	var encoded [2]byte
	binary.BigEndian.PutUint16(encoded[:], value)
	out.Write(encoded[:])
}

func writeInt32(out *bytes.Buffer, value int32) {
	var encoded [4]byte
	binary.BigEndian.PutUint32(encoded[:], uint32(value))
	out.Write(encoded[:])
}

func writeInt64(out *bytes.Buffer, value int64) {
	var encoded [8]byte
	binary.BigEndian.PutUint64(encoded[:], uint64(value))
	out.Write(encoded[:])
}
