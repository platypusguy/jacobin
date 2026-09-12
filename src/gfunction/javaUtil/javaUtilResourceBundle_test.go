package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"reflect"
	"testing"
)

func TestLoadUtilResourceBundle_RegistersAllMethods(t *testing.T) {
	Load_Util_ResourceBundle()

	expected := map[string]int{
		"java/util/ResourceBundle.clearCache()V":                                                                                                                      0,
		"java/util/ResourceBundle.clearCache(Ljava/lang/ClassLoader;)V":                                                                                               1,
		"java/util/ResourceBundle.containsKey(Ljava/lang/String;)Z":                                                                                                   1,
		"java/util/ResourceBundle.getBaseBundleName()Ljava/lang/String;":                                                                                              0,
		"java/util/ResourceBundle.getBundle(Ljava/lang/String;)Ljava/util/ResourceBundle;":                                                                            1,
		"java/util/ResourceBundle.getBundle(Ljava/lang/String;Ljava/util/Locale;)Ljava/util/ResourceBundle;":                                                          2,
		"java/util/ResourceBundle.getBundle(Ljava/lang/String;Ljava/util/Locale;Ljava/lang/ClassLoader;)Ljava/util/ResourceBundle;":                                   3,
		"java/util/ResourceBundle.getBundle(Ljava/lang/String;Ljava/util/Locale;Ljava/lang/Module;)Ljava/util/ResourceBundle;":                                        3,
		"java/util/ResourceBundle.getBundle(Ljava/lang/String;Ljava/util/Locale;Ljava/lang/ClassLoader;Ljava/util/ResourceBundle$Control;)Ljava/util/ResourceBundle;": 4,
		"java/util/ResourceBundle.getBundle(Ljava/lang/String;Ljava/util/Locale;Ljava/util/ResourceBundle$Control;)Ljava/util/ResourceBundle;":                        3,
		"java/util/ResourceBundle.getBundle(Ljava/lang/String;Ljava/util/ResourceBundle$Control;)Ljava/util/ResourceBundle;":                                          2,
		"java/util/ResourceBundle.getBundle(Ljava/lang/String;Ljava/lang/Module;)Ljava/util/ResourceBundle;":                                                          2,
		"java/util/ResourceBundle.getKeys()Ljava/util/Enumeration;":                                                                                                   0,
		"java/util/ResourceBundle.getLocale()Ljava/util/Locale;":                                                                                                      0,
		"java/util/ResourceBundle.getObject(Ljava/lang/String;)Ljava/lang/Object;":                                                                                    1,
		"java/util/ResourceBundle.getString(Ljava/lang/String;)Ljava/lang/String;":                                                                                    1,
		"java/util/ResourceBundle.getStringArray(Ljava/lang/String;)[Ljava/lang/String;":                                                                              1,
		"java/util/ResourceBundle.handleGetObject(Ljava/lang/String;)Ljava/lang/Object;":                                                                              1,
		"java/util/ResourceBundle.handleKeySet()Ljava/util/Set;":                                                                                                      0,
		"java/util/ResourceBundle.keySet()Ljava/util/Set;":                                                                                                            0,
		"java/util/ResourceBundle.setParent(Ljava/util/ResourceBundle;)V":                                                                                             1,
	}

	for key, paramSlots := range expected {
		gmeth, ok := ghelpers.MethodSignatures[key]
		if !ok {
			t.Fatalf("expected method signature %q to be registered", key)
		}
		if gmeth.ParamSlots != paramSlots {
			t.Fatalf("key %q: expected ParamSlots %d, got %d", key, paramSlots, gmeth.ParamSlots)
		}
		if reflect.ValueOf(gmeth.GFunction).Pointer() != reflect.ValueOf(ghelpers.TrapFunction).Pointer() {
			t.Fatalf("key %q: expected GFunction to be ghelpers.TrapFunction", key)
		}
	}
}
