package gfunction

import (
    "jacobin/src/globals"
    "jacobin/src/object"
    "jacobin/src/types"
    "testing"
)

// Helpers for building/extracting byte array objects used by Base64 gfunctions.
func makeByteArrayObject(b []byte) *object.Object {
    jb := object.JavaByteArrayFromGoByteArray(b)
    return object.StringObjectFromJavaByteArray(jb)
}

func getJavaBytesFromArrayObject(arrObj *object.Object) []types.JavaByte {
    // Base64 encode/decode returns an object of class "[B" with FieldTable["value"] as []types.JavaByte
    // but to be robust, read FieldTable["value"] and coerce if needed.
    if arrObj == nil {
        return nil
    }
    if fld, ok := arrObj.FieldTable["value"]; ok {
        switch v := fld.Fvalue.(type) {
        case []types.JavaByte:
            return v
        case []byte:
            return object.JavaByteArrayFromGoByteArray(v)
        default:
            return nil
        }
    }
    return nil
}

func bytesEqual(a, b []types.JavaByte) bool {
    return object.JavaByteArrayEquals(a, b)
}

func TestBase64_Getters_And_WithoutPadding(t *testing.T) {
    globals.InitStringPool()

    // Standard encoder/decoder
    enc := base64GetStdEncoder([]interface{}{}).(*object.Object)
    dec := base64GetStdDecoder([]interface{}{}).(*object.Object)

    if cls := object.GoStringFromStringPoolIndex(enc.KlassName); cls != classNameBase64Encoder {
        t.Fatalf("expected encoder class %s, got %s", classNameBase64Encoder, cls)
    }
    if cls := object.GoStringFromStringPoolIndex(dec.KlassName); cls != classNameBase64Decoder {
        t.Fatalf("expected decoder class %s, got %s", classNameBase64Decoder, cls)
    }

    // withoutPadding on std -> stdRaw; on URL -> urlRaw; on MIME stays MIME
    encRaw := base64WithoutPadding([]interface{}{enc})
    if _, ok := encRaw.(*object.Object); !ok {
        t.Fatalf("withoutPadding should return an encoder object, got %T", encRaw)
    }

    urlEnc := base64GetUrlEncoder([]interface{}{}).(*object.Object)
    urlRaw := base64WithoutPadding([]interface{}{urlEnc}).(*object.Object)
    // Encode a value that normally has padding to verify removal
    src := []byte("hi") // base64 std "aGk="
    out := base64EncodeBsrcToString([]interface{}{urlRaw, makeByteArrayObject(src)}).(*object.Object)
    s := object.GoStringFromStringObject(out)
    if len(s) == 0 || s[len(s)-1] == '=' {
        t.Fatalf("expected no padding in URL raw encoding, got %q", s)
    }
}

func TestBase64_Std_Encode_Decode_RoundTrip(t *testing.T) {
    globals.InitStringPool()

    enc := base64GetStdEncoder([]interface{}{}).(*object.Object)
    dec := base64GetStdDecoder([]interface{}{}).(*object.Object)

    inputs := [][]byte{
        []byte(""),
        []byte("f"),
        []byte("fo"),
        []byte("foo"),
        []byte("hello world"),
        {0x00, 0xFF, 0x10, 0x20, 0x7F},
    }

    for _, in := range inputs {
        srcObj := makeByteArrayObject(in)

        // encode([B)[B -> primitive byte-array object
        encOut := base64EncodeBsrc([]interface{}{enc, srcObj}).(*object.Object)
        encJB := getJavaBytesFromArrayObject(encOut)
        if len(encJB) == 0 && len(in) != 0 {
            t.Fatalf("unexpected empty encoded output for input %v", in)
        }

        // encodeToString([B)Ljava/lang/String;
        encStrObj := base64EncodeBsrcToString([]interface{}{enc, srcObj}).(*object.Object)
        encStr := object.GoStringFromStringObject(encStrObj)
        if len(encStr) != len(encJB) {
            // String should carry same bytes length as encoded output
            t.Fatalf("encodeToString length mismatch: %d vs %d", len(encStr), len(encJB))
        }

        // decode(String)[B path: feed the encoded string as a String object
        encodedAsString := object.StringObjectFromGoString(encStr)
        decOut1 := base64Decode([]interface{}{dec, encodedAsString}).(*object.Object)
        decJB1 := getJavaBytesFromArrayObject(decOut1)
        if !bytesEqual(decJB1, object.JavaByteArrayFromGoByteArray(in)) {
            t.Fatalf("decode(String) round-trip mismatch: in=%v out=%v", in, decJB1)
        }

        // decode([B)[B path: feed the same bytes as a byte-array object
        encBytesObj := makeByteArrayObject(object.GoByteArrayFromJavaByteArray(encJB))
        decOut2 := base64Decode([]interface{}{dec, encBytesObj}).(*object.Object)
        decJB2 := getJavaBytesFromArrayObject(decOut2)
        if !bytesEqual(decJB2, object.JavaByteArrayFromGoByteArray(in)) {
            t.Fatalf("decode([B) round-trip mismatch: in=%v out=%v", in, decJB2)
        }
    }
}

func TestBase64_Std_BsrcBdst_Variants(t *testing.T) {
    globals.InitStringPool()

    enc := base64GetStdEncoder([]interface{}{}).(*object.Object)
    dec := base64GetStdDecoder([]interface{}{}).(*object.Object)

    input := []byte("hello") // aGVsbG8=
    srcObj := makeByteArrayObject(input)

    // Destination objects must have a preexisting value field; initialize with empty
    dstEnc := makeByteArrayObject([]byte{})
    nEnc := base64EncodeBsrcBdst([]interface{}{enc, srcObj, dstEnc}).(int64)
    encJB := dstEnc.FieldTable["value"].Fvalue.([]types.JavaByte)
    if int64(len(encJB)) != nEnc {
        t.Fatalf("encode([B[B)I length mismatch: field=%d ret=%d", len(encJB), nEnc)
    }

    // Now decode into destination
    encBytesObj := makeByteArrayObject(object.GoByteArrayFromJavaByteArray(encJB))
    dstDec := makeByteArrayObject([]byte{})
    nDec := base64DecodeBsrcBdst([]interface{}{dec, encBytesObj, dstDec}).(int64)
    decJB := dstDec.FieldTable["value"].Fvalue.([]types.JavaByte)
    if int64(len(decJB)) != nDec {
        t.Fatalf("decode([B[B)I length mismatch: field=%d ret=%d", len(decJB), nDec)
    }
    if !bytesEqual(decJB, object.JavaByteArrayFromGoByteArray(input)) {
        t.Fatalf("decode into dst mismatch: expected %v got %v", input, decJB)
    }
}

func TestBase64_Url_Encoding_NoPlusSlash(t *testing.T) {
    globals.InitStringPool()

    urlEnc := base64GetUrlEncoder([]interface{}{}).(*object.Object)
    urlDec := base64GetUrlDecoder([]interface{}{}).(*object.Object)

    // Choose bytes that will produce '+' and '/' in standard encoding: 0xFB 0xEF 0xFF ("+++//") patterns
    input := []byte{0xFB, 0xEF, 0xFF}
    srcObj := makeByteArrayObject(input)

    // URL encoding
    encOut := base64EncodeBsrcToString([]interface{}{urlEnc, srcObj}).(*object.Object)
    s := object.GoStringFromStringObject(encOut)
    for _, ch := range s {
        if ch == '+' || ch == '/' {
            t.Fatalf("URL encoding must not contain '+' or '/', got %q", s)
        }
    }

    // Round-trip via URL decoder
    encodedAsString := object.StringObjectFromGoString(s)
    decOut := base64Decode([]interface{}{urlDec, encodedAsString}).(*object.Object)
    decJB := getJavaBytesFromArrayObject(decOut)
    if !bytesEqual(decJB, object.JavaByteArrayFromGoByteArray(input)) {
        t.Fatalf("URL decode round-trip mismatch: expected %v got %v", input, decJB)
    }
}

func TestBase64_Mime_RoundTrip(t *testing.T) {
    globals.InitStringPool()

    mimeEnc := base64GetMimeEncoder([]interface{}{}).(*object.Object)
    mimeDec := base64GetMimeDecoder([]interface{}{}).(*object.Object)

    // Use a longer input to exercise streaming encoder
    input := []byte("The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog.")
    srcObj := makeByteArrayObject(input)

    encStrObj := base64EncodeBsrcToString([]interface{}{mimeEnc, srcObj}).(*object.Object)
    encStr := object.GoStringFromStringObject(encStrObj)

    // Decode back using MIME decoder
    decOut := base64Decode([]interface{}{mimeDec, object.StringObjectFromGoString(encStr)}).(*object.Object)
    decJB := getJavaBytesFromArrayObject(decOut)

    if !bytesEqual(decJB, object.JavaByteArrayFromGoByteArray(input)) {
        t.Fatalf("MIME round-trip mismatch")
    }
}

func TestBase64_Decode_InvalidInput(t *testing.T) {
    globals.InitStringPool()

    dec := base64GetStdDecoder([]interface{}{}).(*object.Object)

    // Invalid base64 string
    bad := object.StringObjectFromGoString("***not_base64***")
    res := base64Decode([]interface{}{dec, bad})
    if _, ok := res.(*GErrBlk); !ok {
        t.Fatalf("expected error block for invalid base64 input, got %T", res)
    }
}
