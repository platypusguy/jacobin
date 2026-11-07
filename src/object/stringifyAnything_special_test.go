/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package object_test

import (
	"jacobin/src/gfunction/javaTime"
	"math/big"
	"testing"

	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaLang"
	"jacobin/src/gfunction/javaMath"
	"jacobin/src/gfunction/javaUtil"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

func TestStringifyAnythingGo_SpecialObjects(t *testing.T) {
	globals.InitGlobals("test")
	global := globals.GetGlobalRef()
	originalInvoker := global.FuncInvokeGFunction
	originalMethodSignatures := make(map[string]ghelpers.GMeth, len(ghelpers.MethodSignatures))
	for signature, method := range ghelpers.MethodSignatures {
		originalMethodSignatures[signature] = method
	}
	t.Cleanup(func() {
		global.FuncInvokeGFunction = originalInvoker
		ghelpers.MethodSignatures = originalMethodSignatures
	})

	javaMath.Load_Math_Big_Integer()
	javaMath.Load_Math_Big_Decimal()
	javaTime.Load_Time_LocalDate()
	javaTime.Load_Time_LocalTime()
	javaTime.Load_Time_LocalDateTime()
	javaLang.Load_Lang_Thread()
	javaLang.Load_Lang_Thread_State()
	javaUtil.Load_Util_Locale()
	global.FuncInvokeGFunction = ghelpers.Invoke

	bigInteger := object.MakePrimitiveObject(types.ClassNameBigInteger, types.Ref,
		new(big.Int).SetInt64(-9876543210))

	bigDecimal := object.MakeEmptyObjectWithClassName(&types.ClassNameBigDecimal)
	if result := ghelpers.Invoke(types.ClassNameBigDecimal+".<init>(Ljava/lang/String;)V",
		[]any{bigDecimal, object.StringObjectFromGoString("123.4500")}); result != nil {
		t.Fatalf("BigDecimal constructor returned %v", result)
	}

	localDateTimeResult := ghelpers.Invoke("java/time/LocalDateTime.of(IIIIII)Ljava/time/LocalDateTime;",
		[]any{int64(2026), int64(6), int64(24), int64(14), int64(30), int64(15)})
	localDateTime, ok := localDateTimeResult.(*object.Object)
	if !ok {
		t.Fatalf("LocalDateTime factory returned %T, want *object.Object", localDateTimeResult)
	}
	localDateResult := ghelpers.Invoke("java/time/LocalDate.of(III)Ljava/time/LocalDate;",
		[]any{int64(2026), int64(6), int64(24)})
	localDate, ok := localDateResult.(*object.Object)
	if !ok {
		t.Fatalf("LocalDate factory returned %T, want *object.Object", localDateResult)
	}
	localTimeResult := ghelpers.Invoke("java/time/LocalTime.of(III)Ljava/time/LocalTime;",
		[]any{int64(14), int64(30), int64(15)})
	localTime, ok := localTimeResult.(*object.Object)
	if !ok {
		t.Fatalf("LocalTime factory returned %T, want *object.Object", localTimeResult)
	}

	locale := ghelpers.LocToObject(ghelpers.Loc{Language: "en", Country: "US", Script: "Latn"})
	localeClassName := "java/util/Locale"
	locale.KlassName = object.StringPoolIndexFromGoString(localeClassName)

	thread := object.MakeEmptyObjectWithClassName(&types.ClassNameThread)
	thread.FieldTable["ID"] = object.Field{Ftype: types.Long, Fvalue: int64(17)}
	thread.FieldTable["name"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString("worker")}
	thread.FieldTable["priority"] = object.Field{Ftype: types.Int, Fvalue: int64(7)}
	thread.FieldTable["state"] = object.Field{Ftype: types.Int, Fvalue: int64(3)}

	tests := []struct {
		name  string
		value any
		want  string
	}{
		{name: "String", value: object.StringObjectFromGoString("string value"), want: "string value"},
		{name: "BigInteger", value: bigInteger, want: "-9876543210"},
		{name: "BigDecimal with scale", value: bigDecimal, want: "123.4500"},
		{name: "LocalDate", value: localDate, want: "2026-06-24"},
		{name: "LocalTime", value: localTime, want: "14:30:15"},
		{name: "LocalDateTime", value: localDateTime, want: "2026-06-24T14:30:15"},
		{name: "Locale", value: locale, want: "en_US_#Latn"},
		{name: "Thread", value: thread, want: "Thread[ID=17, Name=worker, Priority=7, State=WAITING]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := object.StringifyAnythingGo(tt.value); got != tt.want {
				t.Errorf("StringifyAnythingGo() = %q, want %q", got, tt.want)
			}
			javaString := object.StringifyAnythingJava(tt.value)
			if got := object.GoStringFromStringObject(javaString); got != tt.want {
				t.Errorf("StringifyAnythingJava() content = %q, want %q", got, tt.want)
			}
		})
	}

	states := []struct {
		value int64
		want  string
	}{
		{value: 0, want: "NEW"},
		{value: 1, want: "RUNNABLE"},
		{value: 2, want: "BLOCKED"},
		{value: 3, want: "WAITING"},
		{value: 4, want: "TIMED_WAITING"},
		{value: 5, want: "TERMINATED"},
	}
	for _, state := range states {
		t.Run("Thread.State/"+state.want, func(t *testing.T) {
			stateObject := object.MakePrimitiveObject(types.ClassNameThreadState, types.Int, state.value)
			if got := object.StringifyAnythingGo(stateObject); got != state.want {
				t.Errorf("StringifyAnythingGo(Thread.State) = %q, want %q", got, state.want)
			}
			javaString := object.StringifyAnythingJava(stateObject)
			if got := object.GoStringFromStringObject(javaString); got != state.want {
				t.Errorf("StringifyAnythingJava(Thread.State) content = %q, want %q", got, state.want)
			}
		})
	}
}

func TestStringifyAnythingGo_MissingRegisteredToStringFallsBack(t *testing.T) {
	globals.InitGlobals("test")
	global := globals.GetGlobalRef()
	global.FuncInvokeGFunction = ghelpers.Invoke
	className := "example/Widget"
	if got, want := object.StringifyAnythingGo(object.MakeEmptyObjectWithClassName(&className)), "Widget{}"; got != want {
		t.Errorf("StringifyAnythingGo() = %q, want %q", got, want)
	}
}
