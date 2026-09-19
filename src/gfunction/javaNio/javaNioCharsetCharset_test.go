package javaNio

import (
	"jacobin/src/gfunction/ghelpers"
	"reflect"
	"testing"
)

func TestLoad_Nio_Charset_Charset_RegistersAllMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Nio_Charset_Charset()

	checks := []struct {
		key   string
		slots int
		fn    func([]interface{}) interface{}
	}{
		{"java/nio/charset/Charset.<clinit>()V", 0, ghelpers.TrapClass},
		{"java/nio/charset/Charset.aliases()Ljava/util/Set;", 0, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.availableCharsets()Ljava/util/SortedMap;", 0, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.canEncode()Z", 0, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.compareTo(Ljava/nio/charset/Charset;)I", 1, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.contains(Ljava/nio/charset/Charset;)Z", 1, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.decode(Ljava/nio/ByteBuffer;)Ljava/nio/CharBuffer;", 1, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.defaultCharset()Ljava/nio/charset/Charset;", 0, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.displayName()Ljava/lang/String;", 0, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.displayName(Ljava/util/Locale;)Ljava/lang/String;", 1, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.encode(Ljava/lang/String;)Ljava/nio/ByteBuffer;", 1, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.encode(Ljava/nio/CharBuffer;)Ljava/nio/ByteBuffer;", 1, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.equals(Ljava/lang/Object;)Z", 1, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.forName(Ljava/lang/String;)Ljava/nio/charset/Charset;", 1, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.hashCode()I", 0, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.isRegistered()Z", 0, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.isSupported(Ljava/lang/String;)Z", 1, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.name()Ljava/lang/String;", 0, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.newDecoder()Ljava/nio/charset/CharsetDecoder;", 0, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.newEncoder()Ljava/nio/charset/CharsetEncoder;", 0, ghelpers.TrapFunction},
		{"java/nio/charset/Charset.toString()Ljava/lang/String;", 0, ghelpers.TrapFunction},
	}

	if len(ghelpers.MethodSignatures) != len(checks) {
		t.Fatalf("expected %d registered methods, got %d", len(checks), len(ghelpers.MethodSignatures))
	}

	for _, c := range checks {
		gm, ok := ghelpers.MethodSignatures[c.key]
		if !ok {
			t.Fatalf("missing MethodSignatures entry for %s", c.key)
		}
		if gm.ParamSlots != c.slots {
			t.Fatalf("%s ParamSlots expected %d, got %d", c.key, c.slots, gm.ParamSlots)
		}
		if gm.GFunction == nil {
			t.Fatalf("%s GFunction expected non-nil", c.key)
		}
		if reflect.ValueOf(gm.GFunction).Pointer() != reflect.ValueOf(c.fn).Pointer() {
			t.Fatalf("%s GFunction mismatch", c.key)
		}
	}
}
