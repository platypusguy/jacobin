package jdkInternalReflect

import (
	"container/list"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/frames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

// classNameJavaUtilMap is the class name used for the empty Map object
// returned by ReflectionRegisterFilter.
var classNameJavaUtilMap = "java/util/HashMap"

// fieldNameMap is the field name used by java/util/HashMap objects to hold
// their underlying Go map (types.DefHashMap).
var fieldNameMap = "map"

// methodFilterMap corresponds to the Java static field:
//
//	private static Map<Class<?>, Set<String>> methodFilterMap = Map.of();
var methodFilterMap map[*object.Object]map[string]bool = nil

// fieldFilterMap corresponds to the Java static field:
//
//	private static Map<Class<?>, Set<String>> fieldFilterMap = Map.of();
var fieldFilterMap map[*object.Object]map[string]bool = nil

func Load_Internal_Jdk_Reflect() {

	ghelpers.MethodSignatures["jdk/internal/reflect/Reflection.<clinit>()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["jdk/internal/reflect/Reflection.getCallerClass()Ljava/lang/Class;"] =
		ghelpers.GMeth{
			ParamSlots:   0,
			GFunction:    ReflectionGetCallerClass,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["jdk/internal/reflect/Reflection.registerFilter(Ljava/util/Map;Ljava/lang/Class;Ljava/util/Set;)Ljava/util/Map;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ReflectionRegisterFilter,
		}

}

// ReflectionGetCallerClass is the same function as "java/lang/Object.getClass()Ljava/lang/Class;"
// except that is gets the class object of the current frame (parent frame).
func ReflectionGetCallerClass(params []interface{}) interface{} {
	fs, ok := params[0].(*list.List)
	if !ok || object.IsNull(fs) {
		errMsg := fmt.Sprintf("ReflectionGetCallerClass: Invalid or null frame stack reference: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	fr := frames.PeekFrame(fs, 1)

	klassName := fr.ClName
	klass := classloader.MethAreaFetch(klassName)
	if klass == nil || klass.Data == nil {
		errMsg := fmt.Sprintf("ReflectionGetCallerClass: Class %s from parent frame not loaded",
			klassName)
		return ghelpers.GetGErrBlk(excNames.ClassNotLoadedException, errMsg)
	}
	return klass.Data.ClassObject
}

// ReflectionRegisterFilter corresponds to the Java method:
//
//	private static Map<Class<?>, Set<String>> ReflectionRegisterFilter(Map<Class<?>, Set<String>> map,
//	                                                         Class<?> containingClass,
//	                                                         Set<String> names)
//
// In this implementation it always returns a newly-created, empty
// java/util/HashMap object (equivalent to Map.of()).
func ReflectionRegisterFilter(params []interface{}) interface{} {
	var classNameHashMap = "java/util/HashMap"
	var fieldNameMap = "map"
	nilMap := make(types.DefHashMap)
	obj := object.MakeEmptyObjectWithClassName(&classNameHashMap)
	fld := object.Field{Ftype: types.HashMap, Fvalue: nilMap}
	obj.FieldTable[fieldNameMap] = fld
	return obj
}
