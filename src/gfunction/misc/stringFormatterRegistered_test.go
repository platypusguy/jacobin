package misc_test

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaLang"
	"jacobin/src/gfunction/javaMath"
	"jacobin/src/gfunction/javaTime"
	"jacobin/src/gfunction/misc"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestStringFormatter_StringConversionsUseRegisteredToStringMethods(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	global := globals.GetGlobalRef()
	originalInvoker := global.FuncInvokeGFunction
	originalSignatures := make(map[string]ghelpers.GMeth, len(ghelpers.MethodSignatures))
	for signature, method := range ghelpers.MethodSignatures {
		originalSignatures[signature] = method
	}
	t.Cleanup(func() {
		global.FuncInvokeGFunction = originalInvoker
		ghelpers.MethodSignatures = originalSignatures
	})

	javaLang.Load_Lang_Object()
	javaLang.Load_Lang_String()
	javaLang.Load_Lang_Integer()
	javaMath.Load_Math_Big_Integer()
	javaMath.Load_Math_Big_Decimal()
	javaTime.Load_Time_LocalDate()
	javaTime.Load_Time_LocalTime()
	javaTime.Load_Time_LocalDateTime()
	global.FuncInvokeGFunction = ghelpers.Invoke

	integer := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(42))
	objectClass := "java/lang/Object"
	baseObject := object.MakeEmptyObjectWithClassName(&objectClass)
	bigInteger := object.MakeEmptyObjectWithClassName(&types.ClassNameBigInteger)
	ghelpers.InitBigIntegerField(bigInteger, 9876543210)
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

	fmtObj := object.StringObjectFromGoString("%s|%s|%s|%s|%s|%s")
	args := object.MakeArrayFromRawArray([]*object.Object{
		object.StringObjectFromGoString("plain string"), baseObject, integer, bigInteger, bigDecimal, localDateTime,
	})
	out := misc.StringFormatter([]interface{}{fmtObj, args})
	formatted, ok := out.(*object.Object)
	if !ok {
		t.Fatalf("StringFormatter returned %T, want *object.Object", out)
	}
	got := object.GoStringFromStringObject(formatted)
	if want := "plain string|Object[]|42|9876543210|123.4500|2026-06-24T14:30:15"; got != want {
		t.Errorf("StringFormatter() = %q, want %q", got, want)
	}
}
