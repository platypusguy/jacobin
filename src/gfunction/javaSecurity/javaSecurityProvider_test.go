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

// TestLoadSecurityProvider verifies that Load_Security_Provider() properly initializes method signatures
func TestLoadSecurityProvider(t *testing.T) {
	globals.InitGlobals("test")

	// Clear any existing signatures to ensure clean test
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	// Load the security provider methods
	Load_Security_Provider()

	// Verify clinit
	if _, exists := ghelpers.MethodSignatures["java/security/Provider.<clinit>()V"]; !exists {
		t.Errorf("Expected <clinit>()V method to be registered")
	}

	// Verify constructors
	if _, exists := ghelpers.MethodSignatures["java/security/Provider.<init>(Ljava/lang/String;DLjava/lang/String;)V"]; !exists {
		t.Errorf("Expected deprecated constructor to be registered")
	}

	if _, exists := ghelpers.MethodSignatures["java/security/Provider.<init>(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"]; !exists {
		t.Errorf("Expected main constructor to be registered")
	}

	// Verify member functions
	expectedMethods := []string{
		"java/security/Provider.clear()V",
		"java/security/Provider.getInfo()Ljava/lang/String;",
		"java/security/Provider.getName()Ljava/lang/String;",
		"java/security/Provider.getProperty(Ljava/lang/String;)Ljava/lang/String;",
		"java/security/Provider.getService(Ljava/lang/String;Ljava/lang/String;)Ljava/security/Provider$Service;",
		"java/security/Provider.getVersion()D",
		"java/security/Provider.put(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;",
		"java/security/Provider.putAll(Ljava/util/Map;)V",
		"java/security/Provider.putService(Ljava/security/Provider$Service;)V",
		"java/security/Provider.remove(Ljava/lang/Object;)Ljava/lang/Object;",
		"java/security/Provider.removeService(Ljava/security/Provider$Service;)V",
		"java/security/Provider.toString()Ljava/lang/String;",
	}

	for _, method := range expectedMethods {
		if _, exists := ghelpers.MethodSignatures[method]; !exists {
			t.Errorf("Expected method %s to be registered", method)
		}
	}
}

// TestLoadSecurityProviderParamSlots verifies that param slots are correctly set
func TestLoadSecurityProviderParamSlots(t *testing.T) {
	globals.InitGlobals("test")

	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)
	Load_Security_Provider()

	tests := []struct {
		method string
		slots  int
	}{
		{"java/security/Provider.<clinit>()V", 0},
		{"java/security/Provider.<init>(Ljava/lang/String;DLjava/lang/String;)V", 3},
		{"java/security/Provider.<init>(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V", 3},
		{"java/security/Provider.clear()V", 0},
		{"java/security/Provider.getInfo()Ljava/lang/String;", 0},
		{"java/security/Provider.getName()Ljava/lang/String;", 0},
		{"java/security/Provider.getProperty(Ljava/lang/String;)Ljava/lang/String;", 1},
		{"java/security/Provider.getService(Ljava/lang/String;Ljava/lang/String;)Ljava/security/Provider$Service;", 2},
		{"java/security/Provider.getVersion()D", 0},
		{"java/security/Provider.putService(Ljava/security/Provider$Service;)V", 1},
		{"java/security/Provider.toString()Ljava/lang/String;", 0},
	}

	for _, tt := range tests {
		if gmeth, exists := ghelpers.MethodSignatures[tt.method]; exists {
			if gmeth.ParamSlots != tt.slots {
				t.Errorf("Method %s: expected %d param slots, got %d", tt.method, tt.slots, gmeth.ParamSlots)
			}
		}
	}
}

// TestSecurityProviderInitWithFloatVersion tests securityProviderInit with float64 version
func TestSecurityProviderInitWithFloatVersion(t *testing.T) {
	globals.InitGlobals("test")

	// Create provider object
	className := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&className)

	// Create parameters
	nameObj := object.StringObjectFromGoString("TestProvider")
	version := 2.5
	infoObj := object.StringObjectFromGoString("Test Provider Info")

	params := []any{provider, nameObj, version, infoObj}

	result := securityProviderInit(params)

	if result != nil {
		t.Errorf("securityProviderInit should return nil on success, got %v", result)
	}

	// Verify fields
	nameField := provider.FieldTable["name"]
	if nameField.Ftype != types.StringClassName {
		t.Errorf("name field type should be %s, got %s", types.StringClassName, nameField.Ftype)
	}
	nameStr := object.GoStringFromStringObject(nameField.Fvalue.(*object.Object))
	if nameStr != "TestProvider" {
		t.Errorf("name should be 'TestProvider', got %s", nameStr)
	}

	versionField := provider.FieldTable["version"]
	if versionField.Ftype != types.Double {
		t.Errorf("version field type should be %s, got %s", types.Double, versionField.Ftype)
	}
	versionVal := versionField.Fvalue.(float64)
	if versionVal != 2.5 {
		t.Errorf("version should be 2.5, got %f", versionVal)
	}

	infoField := provider.FieldTable["info"]
	if infoField.Ftype != types.StringClassName {
		t.Errorf("info field type should be %s, got %s", types.StringClassName, infoField.Ftype)
	}
	infoStr := object.GoStringFromStringObject(infoField.Fvalue.(*object.Object))
	if infoStr != "Test Provider Info" {
		t.Errorf("info should be 'Test Provider Info', got %s", infoStr)
	}

	// Verify services map is initialized
	servicesField := provider.FieldTable["services"]
	if servicesField.Ftype != types.Map {
		t.Errorf("services field type should be %s, got %s", types.Map, servicesField.Ftype)
	}
	if _, ok := servicesField.Fvalue.(map[string]*object.Object); !ok {
		t.Errorf("services field value should be map[string]*object.Object")
	}
}

// TestSecurityProviderInitWithStringVersion tests securityProviderInit with string version
func TestSecurityProviderInitWithStringVersion(t *testing.T) {
	globals.InitGlobals("test")

	className := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&className)

	nameObj := object.StringObjectFromGoString("TestProvider")
	versionObj := object.StringObjectFromGoString("3.14")
	infoObj := object.StringObjectFromGoString("Test Info")

	params := []any{provider, nameObj, versionObj, infoObj}

	result := securityProviderInit(params)

	if result != nil {
		t.Errorf("securityProviderInit should return nil on success, got %v", result)
	}

	versionField := provider.FieldTable["version"]
	versionVal := versionField.Fvalue.(float64)
	if versionVal != 3.14 {
		t.Errorf("version should be 3.14, got %f", versionVal)
	}
}

// TestSecurityProviderInitWithInvalidStringVersion tests error handling for invalid version string
func TestSecurityProviderInitWithInvalidStringVersion(t *testing.T) {
	globals.InitGlobals("test")

	className := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&className)

	nameObj := object.StringObjectFromGoString("TestProvider")
	versionObj := object.StringObjectFromGoString("invalid")
	infoObj := object.StringObjectFromGoString("Test Info")

	params := []any{provider, nameObj, versionObj, infoObj}

	result := securityProviderInit(params)

	gerr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk for invalid version, got %T", result)
	}

	if gerr.ExceptionType != excNames.VirtualMachineError {
		t.Errorf("expected VirtualMachineError, got %v", gerr.ExceptionType)
	}

	if !strings.Contains(gerr.ErrMsg, "failed parsing version") {
		t.Errorf("error message should mention failed parsing, got: %s", gerr.ErrMsg)
	}
}

// TestSecurityProviderInitWithInvalidVersionType tests error handling for invalid version type
func TestSecurityProviderInitWithInvalidVersionType(t *testing.T) {
	globals.InitGlobals("test")

	className := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&className)

	nameObj := object.StringObjectFromGoString("TestProvider")
	infoObj := object.StringObjectFromGoString("Test Info")

	// Pass an invalid type (int) for version
	params := []any{provider, nameObj, 123, infoObj}

	result := securityProviderInit(params)

	gerr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk for invalid version type, got %T", result)
	}

	if gerr.ExceptionType != excNames.VirtualMachineError {
		t.Errorf("expected VirtualMachineError, got %v", gerr.ExceptionType)
	}

	if !strings.Contains(gerr.ErrMsg, "invalid version type") {
		t.Errorf("error message should mention invalid version type, got: %s", gerr.ErrMsg)
	}
}

// TestSecurityProviderInitWithWhitespace tests that whitespace is trimmed from strings
func TestSecurityProviderInitWithWhitespace(t *testing.T) {
	globals.InitGlobals("test")

	className := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&className)

	nameObj := object.StringObjectFromGoString("  TestProvider  ")
	versionObj := object.StringObjectFromGoString("  1.0  ")
	infoObj := object.StringObjectFromGoString("  Test Info  ")

	params := []any{provider, nameObj, versionObj, infoObj}

	result := securityProviderInit(params)

	if result != nil {
		t.Errorf("securityProviderInit should return nil on success, got %v", result)
	}

	nameStr := object.GoStringFromStringObject(provider.FieldTable["name"].Fvalue.(*object.Object))
	if nameStr != "TestProvider" {
		t.Errorf("name should be trimmed to 'TestProvider', got '%s'", nameStr)
	}

	infoStr := object.GoStringFromStringObject(provider.FieldTable["info"].Fvalue.(*object.Object))
	if infoStr != "Test Info" {
		t.Errorf("info should be trimmed to 'Test Info', got '%s'", infoStr)
	}
}

// TestSecurityProviderToString tests securityProviderToString
func TestSecurityProviderToString(t *testing.T) {
	globals.InitGlobals("test")

	className := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&className)

	nameObj := object.StringObjectFromGoString("MyProvider")
	infoObj := object.StringObjectFromGoString("My Provider Info")
	provider.FieldTable["name"] = object.Field{Ftype: types.StringClassName, Fvalue: nameObj}
	provider.FieldTable["version"] = object.Field{Ftype: types.Double, Fvalue: 1.5}
	provider.FieldTable["info"] = object.Field{Ftype: types.StringClassName, Fvalue: infoObj}

	result := securityProviderToString([]any{provider})

	resultObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", result)
	}

	str := object.GoStringFromStringObject(resultObj)
	expected := "MyProvider 1.5\nMy Provider Info"
	if str != expected {
		t.Errorf("expected '%s', got '%s'", expected, str)
	}
}

// TestSecurityProviderGetServiceWithNilParams tests getService with nil parameters
func TestSecurityProviderGetServiceWithNilParams(t *testing.T) {
	globals.InitGlobals("test")

	className := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&className)
	provider.FieldTable["services"] = object.Field{Ftype: types.Map, Fvalue: map[string]*object.Object{}}

	// Test with nil type parameter
	result := securityProviderGetService([]any{provider, nil, object.StringObjectFromGoString("algo")})
	if result != nil {
		t.Errorf("expected nil when type is nil, got %v", result)
	}

	// Test with nil algorithm parameter
	result = securityProviderGetService([]any{provider, object.StringObjectFromGoString("type"), nil})
	if result != nil {
		t.Errorf("expected nil when algorithm is nil, got %v", result)
	}
}

// TestSecurityProviderGetServiceNotFound tests getService when service doesn't exist
func TestSecurityProviderGetServiceNotFound(t *testing.T) {
	globals.InitGlobals("test")

	className := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&className)
	provider.FieldTable["services"] = object.Field{Ftype: types.Map, Fvalue: map[string]*object.Object{}}

	typeObj := object.StringObjectFromGoString("Cipher")
	algoObj := object.StringObjectFromGoString("AES")

	result := securityProviderGetService([]any{provider, typeObj, algoObj})
	if errBlk, ok := result.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NoSuchAlgorithmException {
		t.Errorf("expected NoSuchAlgorithmException for non-existent service, got %v", result)
	}
}

// TestSecurityProviderGetServiceFound tests getService when service exists
func TestSecurityProviderGetServiceFound(t *testing.T) {
	globals.InitGlobals("test")

	// Create provider
	providerClassName := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&providerClassName)
	services := map[string]*object.Object{}
	provider.FieldTable["services"] = object.Field{Ftype: types.Map, Fvalue: services}

	// Create service
	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)
	service.FieldTable["type"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString("Runtime")}
	service.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString("Security")}

	// Add service to provider
	services["Runtime/Security"] = service

	// Test retrieval
	typeObj := object.StringObjectFromGoString("Runtime")
	algoObj := object.StringObjectFromGoString("Security")

	result := securityProviderGetService([]any{provider, typeObj, algoObj})
	if result != service {
		t.Errorf("expected service to be returned, got %v", result)
	}
}

// TestSecurityProviderPutService tests securityProviderPutService
func TestSecurityProviderPutService(t *testing.T) {
	globals.InitGlobals("test")

	// Create provider
	providerClassName := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&providerClassName)
	services := map[string]*object.Object{}
	provider.FieldTable["services"] = object.Field{Ftype: types.Map, Fvalue: services}

	// Create service
	serviceClassName := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&serviceClassName)
	service.FieldTable["type"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString("MessageDigest")}
	service.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString("SHA-256")}

	// Put service
	result := securityProviderPutService([]any{provider, service})
	if result != nil {
		t.Errorf("securityProviderPutService should return nil, got %v", result)
	}

	// Verify service was added
	if len(services) != 1 {
		t.Errorf("expected 1 service in map, got %d", len(services))
	}

	if services["MessageDigest/SHA-256"] != service {
		t.Errorf("service not properly stored in map")
	}
}

// TestSecurityProviderPutServiceOverwrite tests that putService overwrites existing service
func TestSecurityProviderPutServiceOverwrite(t *testing.T) {
	globals.InitGlobals("test")

	// Create provider
	providerClassName := "java/security/Provider"
	provider := object.MakeEmptyObjectWithClassName(&providerClassName)
	services := map[string]*object.Object{}
	provider.FieldTable["services"] = object.Field{Ftype: types.Map, Fvalue: services}

	// Create first service
	serviceClassName := "java/security/Provider$Service"
	service1 := object.MakeEmptyObjectWithClassName(&serviceClassName)
	service1.FieldTable["type"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString("Cipher")}
	service1.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString("AES")}

	// Put first service
	securityProviderPutService([]any{provider, service1})

	// Create second service with same type/algorithm
	service2 := object.MakeEmptyObjectWithClassName(&serviceClassName)
	service2.FieldTable["type"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString("Cipher")}
	service2.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString("AES")}

	// Put second service (should overwrite)
	securityProviderPutService([]any{provider, service2})

	// Verify only one service exists and it's service2
	if len(services) != 1 {
		t.Errorf("expected 1 service in map, got %d", len(services))
	}

	if services["Cipher/AES"] != service2 {
		t.Errorf("service2 should have overwritten service1")
	}
}

// TestNewGoRuntimeProvider tests NewGoRuntimeProvider helper function
func TestNewGoRuntimeProvider(t *testing.T) {
	globals.InitGlobals("test")

	provider := NewGoRuntimeProvider()

	if provider == nil {
		t.Fatal("NewGoRuntimeProvider should return non-nil provider")
	}

	// Verify provider class name
	className := object.GoStringFromStringPoolIndex(provider.KlassName)
	if className != "java/security/Provider" {
		t.Errorf("expected class name 'java/security/Provider', got '%s'", className)
	}

	// Verify name field
	nameField := provider.FieldTable["name"]
	nameStr := object.GoStringFromStringObject(nameField.Fvalue.(*object.Object))
	if nameStr != types.SecurityProviderName {
		t.Errorf("expected name '%s', got '%s'", types.SecurityProviderName, nameStr)
	}

	// Verify version field
	versionField := provider.FieldTable["version"]
	versionVal := versionField.Fvalue.(float64)
	if versionVal != 1.0 {
		t.Errorf("expected version 1.0, got %f", versionVal)
	}

	// Verify info field
	infoField := provider.FieldTable["info"]
	infoStr := object.GoStringFromStringObject(infoField.Fvalue.(*object.Object))
	if infoStr != types.SecurityProviderInfo {
		t.Errorf("expected info '%s', got '%s'", types.SecurityProviderInfo, infoStr)
	}

	// Verify services map is initialized and contains the default service
	servicesField := provider.FieldTable["services"]
	services, ok := servicesField.Fvalue.(map[string]*object.Object)
	if !ok {
		t.Fatal("services field should be map[string]*object.Object")
	}

	if len(services) < 1 {
		t.Errorf("expected at least 1 default service, got %d", len(services))
	}

	// Verify the default service
	defaultService := services["Runtime/Security"]
	if defaultService == nil {
		t.Fatal("expected default 'Runtime/Security' service to be present")
	}

	// Verify service class name
	serviceClassName := object.GoStringFromStringPoolIndex(defaultService.KlassName)
	if serviceClassName != "java/security/Provider$Service" {
		t.Errorf("expected service class name 'java/security/Provider$Service', got '%s'", serviceClassName)
	}

	// Verify service type
	serviceType := object.GoStringFromStringObject(defaultService.FieldTable["type"].Fvalue.(*object.Object))
	if serviceType != "Runtime" {
		t.Errorf("expected service type 'Runtime', got '%s'", serviceType)
	}

	// Verify service algorithm
	serviceAlgo := object.GoStringFromStringObject(defaultService.FieldTable["algorithm"].Fvalue.(*object.Object))
	if serviceAlgo != "Security" {
		t.Errorf("expected service algorithm 'Security', got '%s'", serviceAlgo)
	}

	// Verify service className
	serviceClassNameObj := defaultService.FieldTable["className"].Fvalue.(*object.Object)
	serviceClassNameStr := object.GoStringFromStringObject(serviceClassNameObj)
	if serviceClassNameStr != types.SecurityProviderName {
		t.Errorf("expected service className '%s', got '%s'", types.SecurityProviderName, serviceClassNameStr)
	}

	// Verify service provider reference
	serviceProvider := defaultService.FieldTable["provider"].Fvalue.(*object.Object)
	if serviceProvider != provider {
		t.Errorf("service should reference its parent provider")
	}
}
