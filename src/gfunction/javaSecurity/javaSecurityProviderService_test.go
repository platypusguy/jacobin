package javaSecurity

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"strings"
	"testing"
)

// TestLoadSecurityProviderService verifies that Load_Security_Provider_Service() properly initializes method signatures
func TestLoadSecurityProviderService(t *testing.T) {
	globals.InitGlobals("test")

	// Clear any existing signatures to ensure clean test
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	// Load the security provider service methods
	Load_Security_Provider_Service()

	// Verify constructor
	if _, exists := ghelpers.MethodSignatures["java/security/Provider$Service.<init>(Ljava/security/Provider;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;[Ljava/lang/String;)V"]; !exists {
		t.Errorf("Expected constructor to be registered")
	}

	// Verify member functions
	expectedMethods := []string{
		"java/security/Provider$Service.getAlgorithm()Ljava/lang/String;",
		"java/security/Provider$Service.getAliases()Ljava/util/List;",
		"java/security/Provider$Service.getAttribute(Ljava/lang/String;)Ljava/lang/String;",
		"java/security/Provider$Service.getClassName()Ljava/lang/String;",
		"java/security/Provider$Service.getProvider()Ljava/security/Provider;",
		"java/security/Provider$Service.getType()Ljava/lang/String;",
		"java/security/Provider$Service.newInstance(Ljava/lang/Object[])Ljava/lang/Object;",
		"java/security/Provider$Service.toString()Ljava/lang/String;",
	}

	for _, method := range expectedMethods {
		if _, exists := ghelpers.MethodSignatures[method]; !exists {
			t.Errorf("Expected method %s to be registered", method)
		}
	}
}

// TestLoadSecurityProviderServiceParamSlots verifies that param slots are correctly set
func TestLoadSecurityProviderServiceParamSlots(t *testing.T) {
	globals.InitGlobals("test")

	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)
	Load_Security_Provider_Service()

	tests := []struct {
		method string
		slots  int
	}{
		{"java/security/Provider$Service.<init>(Ljava/security/Provider;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;[Ljava/lang/String;)V", 5},
		{"java/security/Provider$Service.getAlgorithm()Ljava/lang/String;", 0},
		{"java/security/Provider$Service.getAliases()Ljava/util/List;", 0},
		{"java/security/Provider$Service.getAttribute(Ljava/lang/String;)Ljava/lang/String;", 1},
		{"java/security/Provider$Service.getClassName()Ljava/lang/String;", 0},
		{"java/security/Provider$Service.getProvider()Ljava/security/Provider;", 0},
		{"java/security/Provider$Service.getType()Ljava/lang/String;", 0},
		{"java/security/Provider$Service.newInstance(Ljava/lang/Object[])Ljava/lang/Object;", 1},
		{"java/security/Provider$Service.toString()Ljava/lang/String;", 0},
	}

	for _, tt := range tests {
		if gmeth, exists := ghelpers.MethodSignatures[tt.method]; exists {
			if gmeth.ParamSlots != tt.slots {
				t.Errorf("Method %s: expected %d param slots, got %d", tt.method, tt.slots, gmeth.ParamSlots)
			}
		}
	}
}

// TestSecurityProvSvcInitHappyPath tests securityProvSvcInit with valid parameters
func TestSecurityProvSvcInitHappyPath(t *testing.T) {
	globals.InitGlobals("test")

	// Create provider object
	providerClassName := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&providerClassName)

	// Create service object
	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	// Create parameters
	typeObj := object.StringObjectFromGoString("MessageDigest")
	algorithmObj := object.StringObjectFromGoString("SHA-256")
	classNameObj := object.StringObjectFromGoString("sun.security.provider.SHA256")
	aliases := []*object.Object{
		object.StringObjectFromGoString("SHA256"),
		object.StringObjectFromGoString("2.16.840.1.101.3.4.2.1"),
	}

	params := []any{service, provider, typeObj, algorithmObj, classNameObj, aliases}

	result := securityProvSvcInit(params)

	if result != nil {
		t.Errorf("securityProvSvcInit should return nil on success, got %v", result)
	}

	// Verify provider field
	providerField, exists := service.FieldTable["provider"]
	if !exists {
		t.Fatalf("provider field not set")
	}
	if providerField.Ftype != types.Ref {
		t.Errorf("Expected provider Ftype %s, got %s", types.Ref, providerField.Ftype)
	}
	if providerField.Fvalue.(*object.Object) != provider {
		t.Errorf("provider field value mismatch")
	}

	// Verify type field
	typeField, exists := service.FieldTable["type"]
	if !exists {
		t.Fatalf("type field not set")
	}
	if typeField.Ftype != types.StringClassName {
		t.Errorf("Expected type Ftype %s, got %s", types.StringClassName, typeField.Ftype)
	}
	typeStr := object.GoStringFromStringObject(typeField.Fvalue.(*object.Object))
	if typeStr != "MessageDigest" {
		t.Errorf("Expected type 'MessageDigest', got '%s'", typeStr)
	}

	// Verify algorithm field
	algorithmField, exists := service.FieldTable["algorithm"]
	if !exists {
		t.Fatalf("algorithm field not set")
	}
	algorithmStr := object.GoStringFromStringObject(algorithmField.Fvalue.(*object.Object))
	if algorithmStr != "SHA-256" {
		t.Errorf("Expected algorithm 'SHA-256', got '%s'", algorithmStr)
	}

	// Verify className field
	classNameField, exists := service.FieldTable["className"]
	if !exists {
		t.Fatalf("className field not set")
	}
	classNameStr := object.GoStringFromStringObject(classNameField.Fvalue.(*object.Object))
	if classNameStr != "sun.security.provider.SHA256" {
		t.Errorf("Expected className 'sun.security.provider.SHA256', got '%s'", classNameStr)
	}

	// Verify aliases field
	aliasesField, exists := service.FieldTable["aliases"]
	if !exists {
		t.Fatalf("aliases field not set")
	}
	if aliasesField.Ftype != types.StringArrayClassName {
		t.Errorf("Expected aliases Ftype %s, got %s", types.StringArrayClassName, aliasesField.Ftype)
	}

	// Verify attributes field is initialized
	attributesField, exists := service.FieldTable["attributes"]
	if !exists {
		t.Fatalf("attributes field not set")
	}
	if attributesField.Ftype != types.Map {
		t.Errorf("Expected attributes Ftype %s, got %s", types.Map, attributesField.Ftype)
	}
}

// TestSecurityProvSvcInitWithTrimming tests that securityProvSvcInit trims whitespace
func TestSecurityProvSvcInitWithTrimming(t *testing.T) {
	globals.InitGlobals("test")

	providerClassName := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&providerClassName)

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	typeObj := object.StringObjectFromGoString("  Cipher  ")
	algorithmObj := object.StringObjectFromGoString("  AES  ")
	classNameObj := object.StringObjectFromGoString("  javax.crypto.Cipher  ")
	aliases := []*object.Object{}

	params := []any{service, provider, typeObj, algorithmObj, classNameObj, aliases}

	result := securityProvSvcInit(params)

	if result != nil {
		t.Errorf("securityProvSvcInit should return nil on success, got %v", result)
	}

	typeStr := object.GoStringFromStringObject(service.FieldTable["type"].Fvalue.(*object.Object))
	if typeStr != "Cipher" {
		t.Errorf("Expected trimmed type 'Cipher', got '%s'", typeStr)
	}

	algorithmStr := object.GoStringFromStringObject(service.FieldTable["algorithm"].Fvalue.(*object.Object))
	if algorithmStr != "AES" {
		t.Errorf("Expected trimmed algorithm 'AES', got '%s'", algorithmStr)
	}

	classNameStr := object.GoStringFromStringObject(service.FieldTable["className"].Fvalue.(*object.Object))
	if classNameStr != "javax.crypto.Cipher" {
		t.Errorf("Expected trimmed className 'javax.crypto.Cipher', got '%s'", classNameStr)
	}
}

// TestSecurityProvSvcInitNullProvider tests securityProvSvcInit with null provider
func TestSecurityProvSvcInitNullProvider(t *testing.T) {
	globals.InitGlobals("test")

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	typeObj := object.StringObjectFromGoString("Type")
	algorithmObj := object.StringObjectFromGoString("Algorithm")
	classNameObj := object.StringObjectFromGoString("ClassName")

	params := []any{service, nil, typeObj, algorithmObj, classNameObj}

	result := securityProvSvcInit(params)

	gerr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("Expected *ghelpers.GErrBlk for null provider, got %T", result)
	}
	if gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException, got %v", gerr.ExceptionType)
	}
	if !strings.Contains(gerr.ErrMsg, "null") {
		t.Errorf("Expected error message to contain 'null', got '%s'", gerr.ErrMsg)
	}
}

// TestSecurityProvSvcInitNullType tests securityProvSvcInit with null type
func TestSecurityProvSvcInitNullType(t *testing.T) {
	globals.InitGlobals("test")

	providerClassName := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&providerClassName)

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	algorithmObj := object.StringObjectFromGoString("Algorithm")
	classNameObj := object.StringObjectFromGoString("ClassName")

	params := []any{service, provider, nil, algorithmObj, classNameObj}

	result := securityProvSvcInit(params)

	gerr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("Expected *ghelpers.GErrBlk for null type, got %T", result)
	}
	if gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException, got %v", gerr.ExceptionType)
	}
}

// TestSecurityProvSvcInitNullAlgorithm tests securityProvSvcInit with null algorithm
func TestSecurityProvSvcInitNullAlgorithm(t *testing.T) {
	globals.InitGlobals("test")

	providerClassName := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&providerClassName)

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	typeObj := object.StringObjectFromGoString("Type")
	classNameObj := object.StringObjectFromGoString("ClassName")

	params := []any{service, provider, typeObj, nil, classNameObj}

	result := securityProvSvcInit(params)

	gerr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("Expected *ghelpers.GErrBlk for null algorithm, got %T", result)
	}
	if gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException, got %v", gerr.ExceptionType)
	}
}

// TestSecurityProvSvcInitNullClassName tests securityProvSvcInit with null className
func TestSecurityProvSvcInitNullClassName(t *testing.T) {
	globals.InitGlobals("test")

	providerClassName := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&providerClassName)

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	typeObj := object.StringObjectFromGoString("Type")
	algorithmObj := object.StringObjectFromGoString("Algorithm")

	params := []any{service, provider, typeObj, algorithmObj, nil}

	result := securityProvSvcInit(params)

	gerr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("Expected *ghelpers.GErrBlk for null className, got %T", result)
	}
	if gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException, got %v", gerr.ExceptionType)
	}
}

// TestSecurityProvSvcGetProvider tests securityProvSvcGetProvider
func TestSecurityProvSvcGetProvider(t *testing.T) {
	globals.InitGlobals("test")

	providerClassName := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&providerClassName)

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	service.FieldTable["provider"] = object.Field{Ftype: types.Ref, Fvalue: provider}

	result := securityProvSvcGetProvider([]any{service})

	resultProvider, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}
	if resultProvider != provider {
		t.Errorf("Expected provider to match")
	}
}

// TestSecurityProvSvcGetType tests securityProvSvcGetType
func TestSecurityProvSvcGetType(t *testing.T) {
	globals.InitGlobals("test")

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	typeObj := object.StringObjectFromGoString("Signature")
	service.FieldTable["type"] = object.Field{Ftype: types.StringClassName, Fvalue: typeObj}

	result := securityProvSvcGetType([]any{service})

	resultObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}
	typeStr := object.GoStringFromStringObject(resultObj)
	if typeStr != "Signature" {
		t.Errorf("Expected type 'Signature', got '%s'", typeStr)
	}
}

// TestSecurityProvSvcGetAlgorithm tests securityProvSvcGetAlgorithm
func TestSecurityProvSvcGetAlgorithm(t *testing.T) {
	globals.InitGlobals("test")

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	algorithmObj := object.StringObjectFromGoString("RSA")
	service.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: algorithmObj}

	result := securityProvSvcGetAlgorithm([]any{service})

	resultObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}
	algorithmStr := object.GoStringFromStringObject(resultObj)
	if algorithmStr != "RSA" {
		t.Errorf("Expected algorithm 'RSA', got '%s'", algorithmStr)
	}
}

// TestSecurityProvSvcGetClassName tests securityProvSvcGetClassName
func TestSecurityProvSvcGetClassName(t *testing.T) {
	globals.InitGlobals("test")

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	classNameObj := object.StringObjectFromGoString("sun.security.rsa.RSASignature")
	service.FieldTable["className"] = object.Field{Ftype: types.StringClassName, Fvalue: classNameObj}

	result := securityProvSvcGetClassName([]any{service})

	resultObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}
	classNameStr := object.GoStringFromStringObject(resultObj)
	if classNameStr != "sun.security.rsa.RSASignature" {
		t.Errorf("Expected className 'sun.security.rsa.RSASignature', got '%s'", classNameStr)
	}
}

// TestSecurityProvSvcGetAliases tests securityProvSvcGetAliases
func TestSecurityProvSvcGetAliases(t *testing.T) {
	globals.InitGlobals("test")

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	aliases := object.StringObjectArrayFromGoStringArray([]string{"SHA256", "2.16.840.1.101.3.4.2.1"})
	service.FieldTable["aliases"] = object.Field{Ftype: types.StringArrayClassName, Fvalue: aliases}

	result := securityProvSvcGetAliases([]any{service})

	resultArray, ok := result.([]*object.Object)
	if !ok {
		t.Fatalf("Expected []*object.Object array, got %T", result)
	}
	// Verify it returns the aliases slice
	if len(resultArray) != 2 {
		t.Errorf("Expected 2 aliases, got %d", len(resultArray))
	}
	if resultArray != nil && len(resultArray) > 0 {
		alias1 := object.GoStringFromStringObject(resultArray[0])
		if alias1 != "SHA256" {
			t.Errorf("Expected first alias 'SHA256', got '%s'", alias1)
		}
	}
}

// TestSecurityProvSvcGetAttributeFound tests securityProvSvcGetAttribute when attribute exists
func TestSecurityProvSvcGetAttributeFound(t *testing.T) {
	globals.InitGlobals("test")

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	attributes := map[string]*object.Object{
		"KeySize":     object.StringObjectFromGoString("2048"),
		"BlockSize":   object.StringObjectFromGoString("16"),
		"Implementor": object.StringObjectFromGoString("SunProvider"),
	}
	service.FieldTable["attributes"] = object.Field{Ftype: types.Map, Fvalue: attributes}

	keyObj := object.StringObjectFromGoString("KeySize")
	result := securityProvSvcGetAttribute([]any{service, keyObj})

	resultObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}
	valueStr := object.GoStringFromStringObject(resultObj)
	if valueStr != "2048" {
		t.Errorf("Expected attribute value '2048', got '%s'", valueStr)
	}
}

// TestSecurityProvSvcGetAttributeNotFound tests securityProvSvcGetAttribute when attribute doesn't exist
func TestSecurityProvSvcGetAttributeNotFound(t *testing.T) {
	globals.InitGlobals("test")

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	attributes := map[string]*object.Object{
		"KeySize": object.StringObjectFromGoString("2048"),
	}
	service.FieldTable["attributes"] = object.Field{Ftype: types.Map, Fvalue: attributes}

	keyObj := object.StringObjectFromGoString("NonExistent")
	result := securityProvSvcGetAttribute([]any{service, keyObj})

	if !object.IsNull(result) {
		t.Errorf("Expected object.IsNull(result) for non-existent attribute, got %v", result)
	}
}

// TestSecurityProvSvcGetAttributeNilKey tests securityProvSvcGetAttribute with nil key
func TestSecurityProvSvcGetAttributeNilKey(t *testing.T) {
	globals.InitGlobals("test")

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	attributes := map[string]*object.Object{
		"KeySize": object.StringObjectFromGoString("2048"),
	}
	service.FieldTable["attributes"] = object.Field{Ftype: types.Map, Fvalue: attributes}

	result := securityProvSvcGetAttribute([]any{service, object.StringObjectFromGoString("")})

	if !object.IsNull(result) {
		t.Errorf("Expected nil for nil key, got %v", result)
	}
}

// TestSecurityProvSvcGetAttributeInvalidKeyType tests securityProvSvcGetAttribute with invalid key type
func TestSecurityProvSvcGetAttributeInvalidKeyType(t *testing.T) {
	globals.InitGlobals("test")

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	attributes := map[string]*object.Object{
		"KeySize": object.StringObjectFromGoString("2048"),
	}
	service.FieldTable["attributes"] = object.Field{Ftype: types.Map, Fvalue: attributes}
	result := securityProvSvcGetAttribute([]any{service, 123})
	if !object.IsNull(result) {
		t.Errorf("Expected nil for invalid key type, got %v", result)
	}
}

// TestSecurityProvSvcToString tests securityProvSvcToString
func TestSecurityProvSvcToString(t *testing.T) {
	globals.InitGlobals("test")

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	typeObj := object.StringObjectFromGoString("MessageDigest")
	algorithmObj := object.StringObjectFromGoString("SHA-256")

	service.FieldTable["type"] = object.Field{Ftype: types.StringClassName, Fvalue: typeObj}
	service.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: algorithmObj}

	result := securityProvSvcToString([]any{service})

	resultObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}
	resultStr := object.GoStringFromStringObject(resultObj)
	expected := "MessageDigest.SHA-256"
	if resultStr != expected {
		t.Errorf("Expected toString '%s', got '%s'", expected, resultStr)
	}
}

// TestSecurityProvSvcToStringDifferentValues tests securityProvSvcToString with different type/algorithm
func TestSecurityProvSvcToStringDifferentValues(t *testing.T) {
	globals.InitGlobals("test")

	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)

	typeObj := object.StringObjectFromGoString("Cipher")
	algorithmObj := object.StringObjectFromGoString("AES")

	service.FieldTable["type"] = object.Field{Ftype: types.StringClassName, Fvalue: typeObj}
	service.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: algorithmObj}

	result := securityProvSvcToString([]any{service})

	resultObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}
	resultStr := object.GoStringFromStringObject(resultObj)
	expected := "Cipher.AES"
	if resultStr != expected {
		t.Errorf("Expected toString '%s', got '%s'", expected, resultStr)
	}
}
