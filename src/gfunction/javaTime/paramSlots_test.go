/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"strings"
	"testing"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/util"
)

func loadAllJavaTime() {
	globals.InitGlobals("test")
	Load_Time_DayOfWeek()
	Load_Time_Month()
	Load_Time_Year()
	Load_Time_LocalDate()
	Load_Time_LocalTime()
	Load_Time_LocalDateTime()
	Load_Time_Duration()
	Load_Time_Period()
	Load_Time_Chrono_Chronology()
	Load_Time_Chrono_ChronoLocalDate()
	Load_Time_Chrono_ChronoLocalDateTime()
	Load_Time_Chrono_ChronoPeriod()
	Load_Time_Chrono_ChronoZonedDateTime()
	Load_Time_Chrono_IsoChronology()
	Load_Time_DateTimeException()
	Load_Time_Format_DateTimeFormatter()
	Load_Time_Format_DateTimeFormatterBuilder()
	Load_Time_Format_DateTimeParseException()
	Load_Time_Format_DecimalStyle()
	Load_Time_Format_FormatStyle()
	Load_Time_Format_ResolverStyle()
	Load_Time_Format_SignStyle()
	Load_Time_Format_TextStyle()
	Load_Time_Temporal_TemporalAccessor()
	Load_Time_Temporal_Temporal()
	Load_Time_Temporal_TemporalAdjuster()
	Load_Time_Temporal_TemporalAdjusters()
	Load_Time_Temporal_TemporalAmount()
	Load_Time_Temporal_TemporalUnit()
	Load_Time_Temporal_TemporalField()
	Load_Time_Temporal_ChronoField()
	Load_Time_Temporal_ChronoUnit()
	Load_Time_Temporal_TemporalQuery()
	Load_Time_Temporal_TemporalQueries()
	Load_Time_Temporal_UnsupportedTemporalTypeException()
	Load_Time_Temporal_ValueRange()
	Load_Time_ZoneOffset()
	Load_Time_Zone_ZoneRules()
	Load_Time_Zone_ZoneRulesException()
	Load_Time_Zone_ZoneRulesProvider()
	Load_Time_Zone_ZoneOffsetTransition()
	Load_Time_ZoneId()
	Load_Time_ZoneRegion()
	Load_Time_Instant()
	Load_Time_OffsetDateTime()
	Load_Time_ZonedDateTime()
	Load_Time_Clock()
}

func TestJavaTimeParamSlotsConsistency(t *testing.T) {
	loadAllJavaTime()

	var mismatches []string
	for k, v := range ghelpers.MethodSignatures {
		if !strings.HasPrefix(k, "java/time/") {
			continue
		}
		_, _, methType := util.ParseFQN(k)
		params := util.ParseIncomingParamsFromMethTypeString(methType)
		expectedSlots := len(params)
		if v.ParamSlots != expectedSlots {
			t.Errorf("Signature mismatch for %s: ParamSlots is %d, expected %d (params: %v)", k, v.ParamSlots, expectedSlots, params)
			mismatches = append(mismatches, k)
		}
	}
	t.Logf("Total java/time signatures checked: %d, mismatches: %d", len(ghelpers.MethodSignatures), len(mismatches))
	if len(mismatches) > 0 {
		t.Fatalf("Found %d ParamSlots mismatches", len(mismatches))
	}
}
