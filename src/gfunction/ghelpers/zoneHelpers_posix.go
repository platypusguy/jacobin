//go:build !windows

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package ghelpers

import (
	"encoding/binary"
	"errors"
	"fmt"
	"jacobin/src/globals"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	tzdbFormatVersion = 1
	zoneRulesType     = 1
	maxZoneOffset     = 18 * 60 * 60
)

var zoneDb *ianaTZDB = nil
var zoneDbMu sync.Mutex

// GetZoneIDs returns the zone IDs in the latest database bundled with the JDK.
func GetZoneIDs() ([]string, error) {
	if err := ensureIANATZDB(); err != nil {
		return nil, err
	}
	return append([]string(nil), zoneDb.ids...), nil
}

// GetTZDBVersion returns the version of the JDK's bundled IANA time-zone database.
func GetTZDBVersion() (string, error) {
	if err := ensureIANATZDB(); err != nil {
		return "", err
	}
	return zoneDb.version, nil
}

// GetZoneOffsets returns standard and total UTC offsets in milliseconds at an epoch instant.
func GetZoneOffsets(zoneID string, epochMillis int64) (standardOffsetMillis, totalOffsetMillis int64, err error) {
	if err := ensureIANATZDB(); err != nil {
		return 0, 0, err
	}
	rules, ok := zoneDb.zones[zoneID]
	if !ok {
		return 0, 0, fmt.Errorf("unknown IANA time-zone ID %q", zoneID)
	}
	seconds := floorDiv(epochMillis, 1000)
	return int64(rules.standardOffsetAt(seconds)) * 1000,
		int64(rules.totalOffsetAt(seconds)) * 1000, nil
}

type ianaTZDB struct {
	version string
	ids     []string
	zones   map[string]*ianaZoneRules
}

type ianaZoneRules struct {
	standardTransitions []int64
	standardOffsets     []int32
	savingsTransitions  []int64
	wallOffsets         []int32
	lastRules           []ianaTransitionRule
}

type ianaTransitionRule struct {
	month          time.Month
	dayOfMonth     int
	dayOfWeek      time.Weekday // Java weekday, with Sunday equal to seven; zero means unspecified.
	secondOfDay    int
	timeDefinition uint8 // 0=wall, 1=standard, 2=UTC
	standardOffset int32
	offsetBefore   int32
	offsetAfter    int32
}

type tzdbReader struct {
	data []byte
	off  int
}

func ensureIANATZDB() error {
	zoneDbMu.Lock()
	defer zoneDbMu.Unlock()
	if zoneDb == nil {
		return loadIANATZDB()
	}
	return nil
}

func loadIANATZDB() error {
	javaHome := globals.GetSystemProperty("java.home")
	if javaHome == "" {
		return errors.New("java.home is not configured")
	}
	data, err := os.ReadFile(filepath.Join(javaHome, "lib", "tzdb.dat"))
	if err != nil {
		return fmt.Errorf("read JDK TZDB file: %w", err)
	}
	parsedDB, err := parseIANATZDB(data)
	if err != nil {
		return fmt.Errorf("parse JDK TZDB file: %w", err)
	}
	zoneDb = parsedDB
	return nil
}

func parseIANATZDB(data []byte) (*ianaTZDB, error) {
	rdr := tzdbReader{data: data}
	format, err := rdr.readByte()
	if err != nil {
		return nil, err
	}
	if format != tzdbFormatVersion {
		return nil, fmt.Errorf("unsupported TZDB format version %d", format)
	}
	group, err := rdr.readUTF()
	if err != nil {
		return nil, err
	}
	if group != "TZDB" {
		return nil, fmt.Errorf("unexpected TZDB group %q", group)
	}

	versionCount, err := rdr.readUint16()
	if err != nil {
		return nil, err
	}
	if versionCount == 0 {
		return nil, errors.New("TZDB has no versions")
	}
	versions := make([]string, versionCount)
	for i := range versions {
		versions[i], err = rdr.readUTF()
		if err != nil {
			return nil, fmt.Errorf("read version %d: %w", i, err)
		}
	}

	regionCount, err := rdr.readUint16()
	if err != nil {
		return nil, err
	}
	regions := make([]string, regionCount)
	for i := range regions {
		regions[i], err = rdr.readUTF()
		if err != nil {
			return nil, fmt.Errorf("read region %d: %w", i, err)
		}
	}

	ruleCount, err := rdr.readUint16()
	if err != nil {
		return nil, err
	}
	rules := make([]*ianaZoneRules, ruleCount)
	for i := range rules {
		length, readErr := rdr.readUint16()
		if readErr != nil {
			return nil, fmt.Errorf("read rule %d length: %w", i, readErr)
		}
		serialized, readErr := rdr.readN(int(length))
		if readErr != nil {
			return nil, fmt.Errorf("read rule %d: %w", i, readErr)
		}
		rules[i], readErr = parseZoneRules(serialized)
		if readErr != nil {
			return nil, fmt.Errorf("decode rule %d: %w", i, readErr)
		}
	}

	zones := make(map[string]*ianaZoneRules)
	for versionIndex := range versions {
		mappingCount, readErr := rdr.readUint16()
		if readErr != nil {
			return nil, fmt.Errorf("read version %q mapping count: %w", versions[versionIndex], readErr)
		}
		versionZones := make(map[string]*ianaZoneRules, mappingCount)
		for i := 0; i < int(mappingCount); i++ {
			regionIndex, readErr := rdr.readUint16()
			if readErr != nil {
				return nil, fmt.Errorf("read version %q mapping %d region: %w", versions[versionIndex], i, readErr)
			}
			ruleIndex, readErr := rdr.readUint16()
			if readErr != nil {
				return nil, fmt.Errorf("read version %q mapping %d rule: %w", versions[versionIndex], i, readErr)
			}
			if int(regionIndex) >= len(regions) || int(ruleIndex) >= len(rules) {
				return nil, fmt.Errorf("version %q mapping %d contains an out-of-range index", versions[versionIndex], i)
			}
			id := regions[regionIndex]
			if _, exists := versionZones[id]; exists {
				return nil, fmt.Errorf("duplicate zone mapping for %q in version %q", id, versions[versionIndex])
			}
			versionZones[id] = rules[ruleIndex]
		}
		if versionIndex == len(versions)-1 {
			zones = versionZones
		}
	}
	if rdr.off < len(data) {
		aliasCount, readErr := rdr.readUint16()
		if readErr != nil {
			return nil, fmt.Errorf("read alias count: %w", readErr)
		}
		for i := 0; i < int(aliasCount); i++ {
			aliasIndex, readErr := rdr.readUint16()
			if readErr != nil {
				return nil, fmt.Errorf("read alias %d ID: %w", i, readErr)
			}
			targetIndex, readErr := rdr.readUint16()
			if readErr != nil {
				return nil, fmt.Errorf("read alias %d target: %w", i, readErr)
			}
			if int(aliasIndex) >= len(regions) || int(targetIndex) >= len(regions) {
				return nil, fmt.Errorf("alias %d contains an out-of-range region index", i)
			}
		}
	}
	if rdr.off != len(data) {
		return nil, fmt.Errorf("TZDB has %d trailing bytes", len(data)-rdr.off)
	}
	if len(zones) == 0 {
		return nil, errors.New("latest TZDB version has no zone mappings")
	}
	ids := make([]string, 0, len(zones))
	for id := range zones {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return &ianaTZDB{version: versions[len(versions)-1], ids: ids, zones: zones}, nil
}

func parseZoneRules(data []byte) (*ianaZoneRules, error) {
	rdr := tzdbReader{data: data}
	typeByte, err := rdr.readByte()
	if err != nil {
		return nil, err
	}
	if typeByte != zoneRulesType {
		return nil, fmt.Errorf("unsupported serialized rule type %d", typeByte)
	}

	standardTransitions, err := rdr.readTransitions()
	if err != nil {
		return nil, fmt.Errorf("standard transitions: %w", err)
	}
	standardOffsets, err := rdr.readOffsets(len(standardTransitions) + 1)
	if err != nil {
		return nil, fmt.Errorf("standard offsets: %w", err)
	}
	savingsTransitions, err := rdr.readTransitions()
	if err != nil {
		return nil, fmt.Errorf("savings transitions: %w", err)
	}
	wallOffsets, err := rdr.readOffsets(len(savingsTransitions) + 1)
	if err != nil {
		return nil, fmt.Errorf("wall offsets: %w", err)
	}
	ruleCount, err := rdr.readByte()
	if err != nil {
		return nil, err
	}
	if ruleCount > 16 {
		return nil, fmt.Errorf("too many recurring transition rules: %d", ruleCount)
	}
	lastRules := make([]ianaTransitionRule, int(ruleCount))
	for i := range lastRules {
		lastRules[i], err = rdr.readTransitionRule()
		if err != nil {
			return nil, fmt.Errorf("recurring rule %d: %w", i, err)
		}
	}
	if rdr.off != len(data) {
		return nil, fmt.Errorf("serialized rule has %d trailing bytes", len(data)-rdr.off)
	}
	return &ianaZoneRules{
		standardTransitions: standardTransitions,
		standardOffsets:     standardOffsets,
		savingsTransitions:  savingsTransitions,
		wallOffsets:         wallOffsets,
		lastRules:           lastRules,
	}, nil
}

func (rdr *tzdbReader) readTransitions() ([]int64, error) {
	count, err := rdr.readInt32()
	if err != nil {
		return nil, err
	}
	if count < 0 || count > 1024 {
		return nil, fmt.Errorf("invalid transition count %d", count)
	}
	transitions := make([]int64, int(count))
	for i := range transitions {
		transitions[i], err = rdr.readEpochSecond()
		if err != nil {
			return nil, fmt.Errorf("transition %d: %w", i, err)
		}
		if i > 0 && transitions[i] <= transitions[i-1] {
			return nil, errors.New("transition instants are not strictly increasing")
		}
	}
	return transitions, nil
}

func (rdr *tzdbReader) readOffsets(count int) ([]int32, error) {
	offsets := make([]int32, count)
	for i := range offsets {
		encoded, err := rdr.readByte()
		if err != nil {
			return nil, err
		}
		var offset int32
		if int8(encoded) == 127 {
			offset, err = rdr.readInt32()
			if err != nil {
				return nil, err
			}
		} else {
			offset = int32(int8(encoded)) * 900
		}
		if !validZoneOffset(offset) {
			return nil, fmt.Errorf("invalid offset %d seconds", offset)
		}
		offsets[i] = offset
	}
	return offsets, nil
}

func (rdr *tzdbReader) readEpochSecond() (int64, error) {
	first, err := rdr.readByte()
	if err != nil {
		return 0, err
	}
	if first == 0xff {
		return rdr.readInt64()
	}
	second, err := rdr.readByte()
	if err != nil {
		return 0, err
	}
	third, err := rdr.readByte()
	if err != nil {
		return 0, err
	}
	compact := int64(first)<<16 | int64(second)<<8 | int64(third)
	return compact*900 - 4575744000, nil
}

func (rdr *tzdbReader) readTransitionRule() (ianaTransitionRule, error) {
	packed, err := rdr.readInt32()
	if err != nil {
		return ianaTransitionRule{}, err
	}
	month := time.Month(uint32(packed) >> 28)
	dayOfMonth := int((uint32(packed)>>22)&0x3f) - 32
	dayOfWeekCode := (uint32(packed) >> 19) & 7
	timeCode := int((uint32(packed) >> 14) & 0x1f)
	timeDefinition := uint8((uint32(packed) >> 12) & 3)
	standardCode := int((uint32(packed) >> 4) & 0xff)
	beforeCode := int((uint32(packed) >> 2) & 3)
	afterCode := int(uint32(packed) & 3)

	if month < time.January || month > time.December || dayOfMonth == 0 || dayOfMonth < -28 || dayOfMonth > 31 {
		return ianaTransitionRule{}, errors.New("invalid month or day-of-month indicator")
	}
	var dayOfWeek time.Weekday
	if dayOfWeekCode != 0 {
		dayOfWeek = time.Weekday(dayOfWeekCode)
	}
	secondOfDay := timeCode * 3600
	if timeCode == 31 {
		var encodedTime int32
		encodedTime, err = rdr.readInt32()
		if err != nil {
			return ianaTransitionRule{}, err
		}
		secondOfDay = int(encodedTime)
	} else if timeCode > 24 {
		return ianaTransitionRule{}, fmt.Errorf("invalid time-of-day code %d", timeCode)
	}
	if secondOfDay < 0 || secondOfDay > 86400 || timeDefinition > 2 {
		return ianaTransitionRule{}, errors.New("invalid transition time or time definition")
	}

	var standardOffset int32
	if standardCode == 255 {
		standardOffset, err = rdr.readInt32()
		if err != nil {
			return ianaTransitionRule{}, err
		}
	} else {
		standardOffset = int32(standardCode-128) * 900
	}
	beforeOffset, err := rdr.readRuleOffset(standardOffset, beforeCode)
	if err != nil {
		return ianaTransitionRule{}, fmt.Errorf("offset before: %w", err)
	}
	afterOffset, err := rdr.readRuleOffset(standardOffset, afterCode)
	if err != nil {
		return ianaTransitionRule{}, fmt.Errorf("offset after: %w", err)
	}
	if !validZoneOffset(standardOffset) || beforeOffset == afterOffset {
		return ianaTransitionRule{}, errors.New("invalid recurring-rule offsets")
	}
	return ianaTransitionRule{
		month:          month,
		dayOfMonth:     dayOfMonth,
		dayOfWeek:      dayOfWeek,
		secondOfDay:    secondOfDay,
		timeDefinition: timeDefinition,
		standardOffset: standardOffset,
		offsetBefore:   beforeOffset,
		offsetAfter:    afterOffset,
	}, nil
}

func (rdr *tzdbReader) readRuleOffset(standardOffset int32, code int) (int32, error) {
	var offset int32
	if code == 3 {
		var err error
		offset, err = rdr.readInt32()
		if err != nil {
			return 0, err
		}
	} else {
		offset = standardOffset + int32(code)*1800
	}
	if !validZoneOffset(offset) {
		return 0, fmt.Errorf("invalid offset %d seconds", offset)
	}
	return offset, nil
}

func (rdr *tzdbReader) readUTF() (string, error) {
	length, err := rdr.readUint16()
	if err != nil {
		return "", err
	}
	value, err := rdr.readN(int(length))
	if err != nil {
		return "", err
	}
	if !utf8.Valid(value) {
		return "", errors.New("invalid UTF-8 string in TZDB")
	}
	return string(value), nil
}

func (rdr *tzdbReader) readByte() (byte, error) {
	value, err := rdr.readN(1)
	if err != nil {
		return 0, err
	}
	return value[0], nil
}

func (rdr *tzdbReader) readUint16() (uint16, error) {
	value, err := rdr.readN(2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(value), nil
}

func (rdr *tzdbReader) readInt32() (int32, error) {
	value, err := rdr.readN(4)
	if err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(value)), nil
}

func (rdr *tzdbReader) readInt64() (int64, error) {
	value, err := rdr.readN(8)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(value)), nil
}

func (rdr *tzdbReader) readN(count int) ([]byte, error) {
	if count < 0 || rdr.off > len(rdr.data)-count {
		return nil, errors.New("unexpected end of TZDB data")
	}
	value := rdr.data[rdr.off : rdr.off+count]
	rdr.off += count
	return value, nil
}

func (r *ianaZoneRules) standardOffsetAt(epochSecond int64) int32 {
	index := sort.Search(len(r.standardTransitions), func(i int) bool {
		return r.standardTransitions[i] > epochSecond
	})
	return r.standardOffsets[index]
}

func (r *ianaZoneRules) totalOffsetAt(epochSecond int64) int32 {
	index := sort.Search(len(r.savingsTransitions), func(i int) bool {
		return r.savingsTransitions[i] > epochSecond
	})
	offset := r.wallOffsets[index]
	if len(r.lastRules) == 0 || (len(r.savingsTransitions) > 0 && epochSecond <= r.savingsTransitions[len(r.savingsTransitions)-1]) {
		return offset
	}

	year := time.Unix(epochSecond, 0).UTC().Year()
	var latest int64 = -1 << 63
	for candidateYear := year - 2; candidateYear <= year+2; candidateYear++ {
		for _, rule := range r.lastRules {
			transition, err := rule.transition(candidateYear)
			if err == nil && transition > latest && transition <= epochSecond {
				latest = transition
				offset = rule.offsetAfter
			}
		}
	}
	return offset
}

func (rule ianaTransitionRule) transition(year int) (int64, error) {
	firstDay := time.Date(year, rule.month, 1, 0, 0, 0, 0, time.UTC)
	var date time.Time
	if rule.dayOfMonth < 0 {
		date = time.Date(year, rule.month+1, 0, 0, 0, 0, 0, time.UTC).AddDate(0, 0, rule.dayOfMonth+1)
		if rule.dayOfWeek != 0 {
			date = date.AddDate(0, 0, -int((date.Weekday()-rule.dayOfWeek+7)%7))
		}
	} else {
		date = firstDay.AddDate(0, 0, rule.dayOfMonth-1)
		if rule.dayOfWeek != 0 {
			date = date.AddDate(0, 0, int((rule.dayOfWeek-date.Weekday()+7)%7))
		}
	}
	localSeconds := date.Unix() + int64(rule.secondOfDay)
	switch rule.timeDefinition {
	case 0:
		return localSeconds - int64(rule.offsetBefore), nil
	case 1:
		return localSeconds - int64(rule.standardOffset), nil
	case 2:
		return localSeconds, nil
	default:
		return 0, errors.New("invalid time definition")
	}
}

func validZoneOffset(offset int32) bool {
	return offset >= -maxZoneOffset && offset <= maxZoneOffset
}

func floorDiv(value, divisor int64) int64 {
	quotient := value / divisor
	if value%divisor < 0 {
		quotient--
	}
	return quotient
}
