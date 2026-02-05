package javaSecurity

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

// helper: make Java byte[] object from Go bytes (MessageDigest tests only)
func mdMakeByteArrayObject(b []byte) *object.Object {
	jb := object.JavaByteArrayFromGoByteArray(b)
	return object.StringObjectFromJavaByteArray(jb)
}

// helper: extract Go bytes from Java byte[] object (MessageDigest tests only)
func mdBytesFromArrayObject(o *object.Object) []byte {
	jb := o.FieldTable["value"].Fvalue.([]types.JavaByte)
	return object.GoByteArrayFromJavaByteArray(jb)
}

func TestLoadSecurityMessageDigest_Registers(t *testing.T) {
	globals.InitGlobals("test")
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)
	Load_Security_MessageDigest()

	checks := []struct {
		sig   string
		slots int
	}{
		{"java/security/MessageDigest.getInstance(Ljava/lang/String;)Ljava/security/MessageDigest;", 1},
		{"java/security/MessageDigest.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljava/security/MessageDigest;", 2},
		{"java/security/MessageDigest.getAlgorithm()Ljava/lang/String;", 0},
		{"java/security/MessageDigest.getDigestLength()I", 0},
		{"java/security/MessageDigest.update(B)V", 1},
		{"java/security/MessageDigest.update([B)V", 1},
		{"java/security/MessageDigest.update([BII)V", 3},
		{"java/security/MessageDigest.digest()[B", 0},
		{"java/security/MessageDigest.digest([B)[B", 1},
		{"java/security/MessageDigest.digest([BII)I", 3},
		{"java/security/MessageDigest.reset()V", 0},
		{"java/security/MessageDigest.isEqual([B[B)Z", 2},
		{"java/security/MessageDigest.toString()Ljava/lang/String;", 0},
		{"java/security/MessageDigest.clone()Ljava/lang/Object;", 0},
	}
	for _, c := range checks {
		gm, ok := ghelpers.MethodSignatures[c.sig]
		if !ok {
			t.Fatalf("missing method signature: %s", c.sig)
		}
		if gm.ParamSlots != c.slots {
			t.Errorf("%s: expected %d slots, got %d", c.sig, c.slots, gm.ParamSlots)
		}
	}
}

func TestMsgDig_GetInstance_ValidAlgorithms(t *testing.T) {
	globals.InitGlobals("test")
	Load_Security_Provider()
	algos := []string{"MD5", "SHA-1", "SHA-224", "SHA-256", "SHA-384", "SHA-512", "SHA-512/224", "SHA-512/256"}
	for _, a := range algos {
		t.Run(a, func(t *testing.T) {
			obj := object.StringObjectFromGoString(a)
			ret := msgdigGetInstance([]any{obj})
			md, ok := ret.(*object.Object)
			if !ok {
				t.Fatalf("expected *object.Object, got %T", ret)
			}
			// algorithm stored
			gotAlg := object.GoStringFromStringObject(md.FieldTable["algorithm"].Fvalue.(*object.Object))
			if gotAlg != a {
				t.Errorf("algorithm stored mismatch: want %s got %s", a, gotAlg)
			}
			// provider present
			if md.FieldTable["provider"].Fvalue == nil {
				t.Errorf("expected provider set")
			}
			// getDigestLength() returns non-zero for supported algorithms
			if l := msgdigGetDigestLength([]any{md}).(int64); l == 0 {
				t.Errorf("expected non-zero digest length for %s", a)
			}
		})
	}
}

func TestMsgDig_GetInstance_Unsupported(t *testing.T) {
	globals.InitGlobals("test")
	Load_Security_Provider()
	obj := object.StringObjectFromGoString("FOO")
	ret := msgdigGetInstance([]any{obj})
	ge := ret.(*ghelpers.GErrBlk)
	if ge.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException, got %d (%s)", ge.ExceptionType, ge.ErrMsg)
	}
}

func TestMsgDig_GetInstance_WithProvider(t *testing.T) {
	globals.InitGlobals("test")
	Load_Security_Provider()
	alg := object.StringObjectFromGoString("SHA-256")
	// wrong provider
	badProv := object.StringObjectFromGoString("OtherProv")
	ret := msgdigGetInstanceProvider([]any{alg, badProv})
	ge := ret.(*ghelpers.GErrBlk)
	if ge.ExceptionType != excNames.ProviderNotFoundException {
		t.Fatalf("expected ProviderNotFoundException, got %d", ge.ExceptionType)
	}
	// correct provider name
	goodProv := object.StringObjectFromGoString(types.SecurityProviderName)
	ret2 := msgdigGetInstanceProvider([]any{alg, goodProv})
	if _, ok := ret2.(*object.Object); !ok {
		t.Fatalf("expected *object.Object for good provider, got %T", ret2)
	}
}

func TestMsgDig_UpdateAndDigest_CorrectnessAndReset(t *testing.T) {
	globals.InitGlobals("test")
	Load_Security_Provider()
	alg := "SHA-256"
	mdObj := msgdigGetInstance([]any{object.StringObjectFromGoString(alg)}).(*object.Object)

	// update with "abc"
	data := []byte("abc")
	arr := mdMakeByteArrayObject(data)
	if r := msgdigUpdateBytes([]any{mdObj, arr}); r != nil {
		t.Fatalf("unexpected error updating bytes: %#v", r)
	}
	// digest
	d := msgdigDigest([]any{mdObj}).(*object.Object)
	got := mdBytesFromArrayObject(d)
	exp := sha256.Sum256(data)
	if hex.EncodeToString(got) != hex.EncodeToString(exp[:]) {
		t.Errorf("digest mismatch: got %s want %s", hex.EncodeToString(got), hex.EncodeToString(exp[:]))
	}
	// buffer should be reset; next digest should be of empty string
	d2 := msgdigDigest([]any{mdObj}).(*object.Object)
	got2 := mdBytesFromArrayObject(d2)
	exp2 := sha256.Sum256(nil)
	if hex.EncodeToString(got2) != hex.EncodeToString(exp2[:]) {
		t.Errorf("post-reset digest mismatch: got %s want %s", hex.EncodeToString(got2), hex.EncodeToString(exp2[:]))
	}
}

func TestMsgDig_UpdateVariantsAndBounds(t *testing.T) {
	globals.InitGlobals("test")
	Load_Security_Provider()
	mdObj := msgdigGetInstance([]any{object.StringObjectFromGoString("MD5")}).(*object.Object)

	// update single bytes 'a','b'
	msgdigUpdateByte([]any{mdObj, int64('a')})
	msgdigUpdateByte([]any{mdObj, int64('b')})

	// update with [c d e] using [BII] to take d,e
	arr := mdMakeByteArrayObject([]byte{'c', 'd', 'e'})
	if r := msgdigUpdateBytesII([]any{mdObj, arr, int64(1), int64(2)}); r != nil {
		t.Fatalf("unexpected error from update [BII]: %#v", r)
	}

	// digest should equal MD5("abde")
	d := msgdigDigest([]any{mdObj}).(*object.Object)
	got := mdBytesFromArrayObject(d)
	h := md5.Sum([]byte("ab" + "de"))
	if hex.EncodeToString(got) != hex.EncodeToString(h[:]) {
		t.Errorf("md5 mismatch got %s want %s", hex.EncodeToString(got), hex.EncodeToString(h[:]))
	}

	// bounds check: offset/length invalid
	mdObj2 := msgdigGetInstance([]any{object.StringObjectFromGoString("MD5")}).(*object.Object)
	bad := msgdigUpdateBytesII([]any{mdObj2, arr, int64(5), int64(1)})
	ge := bad.(*ghelpers.GErrBlk)
	if ge.ExceptionType != excNames.IndexOutOfBoundsException {
		t.Fatalf("expected IndexOutOfBoundsException, got %d", ge.ExceptionType)
	}
}

func TestMsgDig_DigestWithInput(t *testing.T) {
	globals.InitGlobals("test")
	Load_Security_Provider()
	mdObj := msgdigGetInstance([]any{object.StringObjectFromGoString("SHA-1")}).(*object.Object)
	data := []byte("hello world")
	out := msgdigDigestBytes([]any{mdObj, mdMakeByteArrayObject(data)})
	if _, ok := out.(*object.Object); !ok {
		t.Fatalf("expected byte[] object, got %T", out)
	}
	got := mdBytesFromArrayObject(out.(*object.Object))
	exp := sha1.Sum(data)
	if hex.EncodeToString(got) != hex.EncodeToString(exp[:]) {
		t.Errorf("sha1 mismatch got %s want %s", hex.EncodeToString(got), hex.EncodeToString(exp[:]))
	}
}

func TestMsgDig_DigestIntoBuffer_TooSmall(t *testing.T) {
	globals.InitGlobals("test")
	Load_Security_Provider()
	mdObj := msgdigGetInstance([]any{object.StringObjectFromGoString("SHA-512/256")}).(*object.Object)
	// small buffer
	buf := mdMakeByteArrayObject(make([]byte, 16))
	ret := msgdigDigestBytesII([]any{mdObj, buf, int64(0), int64(16)})
	ge := ret.(*ghelpers.GErrBlk)
	if ge.ExceptionType != excNames.IllegalStateException {
		t.Fatalf("expected IllegalStateException, got %d", ge.ExceptionType)
	}
}

func TestMsgDig_DigestIntoBuffer_WritesAndReturnsLen(t *testing.T) {
	globals.InitGlobals("test")
	Load_Security_Provider()
	mdObj := msgdigGetInstance([]any{object.StringObjectFromGoString("SHA-384")}).(*object.Object)
	// supply some data so digest isn't of empty string
	_ = msgdigUpdateBytes([]any{mdObj, mdMakeByteArrayObject([]byte("xyz"))})

	buf := make([]byte, 100)
	arr := mdMakeByteArrayObject(buf)
	wrote := msgdigDigestBytesII([]any{mdObj, arr, int64(2), int64(len(buf) - 2)})
	// expect number of bytes written equals SHA-384 len 48
	if wrote.(int64) != 48 {
		t.Fatalf("expected 48 bytes written, got %d", wrote.(int64))
	}
	out := mdBytesFromArrayObject(arr)
	// compute expected digest for "xyz"
	h := sha512.New384()
	h.Write([]byte("xyz"))
	exp := h.Sum(nil)
	if hex.EncodeToString(out[2:2+48]) != hex.EncodeToString(exp) {
		t.Errorf("written digest mismatch")
	}
}

func TestMsgDig_Reset(t *testing.T) {
	globals.InitGlobals("test")
	Load_Security_Provider()
	mdObj := msgdigGetInstance([]any{object.StringObjectFromGoString("SHA-512")}).(*object.Object)
	_ = msgdigUpdateBytes([]any{mdObj, mdMakeByteArrayObject([]byte("data"))})
	// reset then digest empty
	_ = msgdigReset([]any{mdObj})
	d := msgdigDigest([]any{mdObj}).(*object.Object)
	got := mdBytesFromArrayObject(d)
	exp := sha512.Sum512(nil)
	if hex.EncodeToString(got) != hex.EncodeToString(exp[:]) {
		t.Errorf("reset not effective; got %s want %s", hex.EncodeToString(got), hex.EncodeToString(exp[:]))
	}
}

func TestMsgDig_IsEqual(t *testing.T) {
	globals.InitGlobals("test")
	a := mdMakeByteArrayObject([]byte{1, 2, 3})
	b := mdMakeByteArrayObject([]byte{1, 2, 3})
	c := mdMakeByteArrayObject([]byte{1, 2, 4})
	if msgdigIsEqual([]any{a, b}).(int64) != types.JavaBoolTrue {
		t.Errorf("expected true for equal arrays")
	}
	if msgdigIsEqual([]any{a, c}).(int64) != types.JavaBoolFalse {
		t.Errorf("expected false for different arrays")
	}
}

func TestMsgDig_ToString_And_Clone(t *testing.T) {
	globals.InitGlobals("test")
	Load_Security_Provider()
	mdObj := msgdigGetInstance([]any{object.StringObjectFromGoString("SHA-256")}).(*object.Object)
	// toString contains algorithm
	sObj := msgdigToString([]any{mdObj}).(*object.Object)
	s := object.GoStringFromStringObject(sObj)
	if s != "MessageDigest[SHA-256]" {
		t.Errorf("unexpected toString: %s", s)
	}
	// add data and clone
	_ = msgdigUpdateBytes([]any{mdObj, mdMakeByteArrayObject([]byte("abc"))})
	clone := msgdigClone([]any{mdObj}).(*object.Object)
	// mutate original, ensure clone buffer unchanged
	_ = msgdigUpdateByte([]any{mdObj, int64('Z')})
	origBuf := mdObj.FieldTable["buffer"].Fvalue.([]types.JavaByte)
	cloneBuf := clone.FieldTable["buffer"].Fvalue.([]types.JavaByte)
	if len(cloneBuf) == len(origBuf) {
		t.Errorf("expected clone buffer to be independent copy")
	}
}
