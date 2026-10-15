package javaIo

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaLang"
	"jacobin/src/gfunction/javaUtil"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBufferedReaderInit_Success(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := []byte("Hello\nWorld\n")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			ghelpers.FilePath: {Ftype: types.Array, Fvalue: object.JavaByteArrayFromGoString(filePath)},
		},
	}

	brObj := &object.Object{FieldTable: map[string]object.Field{}}

	params := []interface{}{brObj, fileObj}
	res := bufferedReaderInit(params)
	if res != nil {
		t.Fatalf("Expected success, got error: %v", res)
	}

	if _, ok := brObj.FieldTable[ghelpers.FilePath]; !ok {
		t.Errorf("ghelpers.FilePath field not set on BufferedReader object")
	}
	if fh, ok := brObj.FieldTable[ghelpers.FileHandle]; !ok {
		t.Errorf("ghelpers.FileHandle field not set on BufferedReader object")
	} else {
		if f, ok := fh.Fvalue.(*os.File); ok {
			_ = f.Close()
		}
	}
}

func TestBufferedReaderInit_FileNotFound(t *testing.T) {
	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			ghelpers.FilePath: {Ftype: types.Array, Fvalue: object.JavaByteArrayFromGoString("/no/such/file")},
		},
	}
	brObj := &object.Object{FieldTable: map[string]object.Field{}}

	params := []interface{}{brObj, fileObj}
	res := bufferedReaderInit(params)
	if res == nil {
		t.Errorf("Expected FileNotFoundException error, got nil")
		return
	}
	errObj, ok := res.(*ghelpers.GErrBlk)
	if !ok || errObj.ExceptionType != excNames.FileNotFoundException {
		t.Errorf("Expected FileNotFoundException, got %#v", res)
	}
}

func TestBufferedReaderInit_MissingFilePathField(t *testing.T) {
	fileObj := &object.Object{FieldTable: map[string]object.Field{}}
	brObj := &object.Object{FieldTable: map[string]object.Field{}}

	params := []interface{}{brObj, fileObj}
	res := bufferedReaderInit(params)
	if res == nil {
		t.Errorf("Expected InvalidTypeException error, got nil")
		return
	}
	errObj, ok := res.(*ghelpers.GErrBlk)
	if !ok || errObj.ExceptionType != excNames.InvalidTypeException {
		t.Errorf("Expected InvalidTypeException, got %#v", res)
	}
}

func TestBufferedReaderMarkSupported_False(t *testing.T) {
	params := []interface{}{}
	res := bufferedReaderMarkSupported(params)
	if res != types.JavaBoolFalse {
		t.Errorf("Expected JavaBoolFalse (0), got %v", res)
	}
}

func TestBufferedReaderReadLine_FirstLine(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := []byte("line1\nline2\n")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	f, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer f.Close()

	brObj := &object.Object{
		FieldTable: map[string]object.Field{
			ghelpers.FileHandle: {Ftype: types.Ref, Fvalue: f},
		},
	}

	params := []interface{}{brObj}
	res := bufferedReaderReadLine(params)
	if !object.IsStringObject(res) {
		t.Fatalf("Expected Java String object, got %#v", res)
	}
	goStr := object.GoStringFromStringObject(res.(*object.Object))
	if goStr != "line1" {
		t.Errorf("Expected 'line1', got %q", goStr)
	}
}

func TestBufferedReaderReadLine_AtEOF(t *testing.T) {
	brObj := &object.Object{
		FieldTable: make(map[string]object.Field),
	}
	ghelpers.EofSet(brObj, true)

	params := []interface{}{brObj}
	res := bufferedReaderReadLine(params)
	if res != object.Null {
		t.Errorf("Expected Null at EOF, got %#v", res)
	}
}

func TestBufferedReaderReadLine_MissingFileHandle(t *testing.T) {
	brObj := &object.Object{FieldTable: map[string]object.Field{}}
	params := []interface{}{brObj}
	res := bufferedReaderReadLine(params)
	if res == nil {
		t.Errorf("Expected IOException error, got nil")
		return
	}
	errObj, ok := res.(*ghelpers.GErrBlk)
	if !ok || errObj.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException, got %#v", res)
	}
}

func TestBufferedReaderReadLine_MultiLineSequential(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := []byte("first\nsecond\nthird")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	f, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer f.Close()

	brObj := &object.Object{
		FieldTable: map[string]object.Field{
			ghelpers.FileHandle: {Ftype: types.Ref, Fvalue: f},
		},
	}
	params := []interface{}{brObj}

	expected := []string{"first", "second", "third"}
	for i, want := range expected {
		res := bufferedReaderReadLine(params)
		if res == object.Null {
			t.Fatalf("Unexpected EOF on line %d", i+1)
		}
		if !object.IsStringObject(res) {
			t.Fatalf("Expected Java String object, got %#v", res)
		}
		got := object.GoStringFromStringObject(res.(*object.Object))
		if got != want {
			t.Errorf("Line %d: expected %q, got %q", i+1, want, got)
		}
	}
}

func TestProcessBuilderBufferedReaderPipeline(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()

	Load_Io_BufferedReader()
	Load_Io_InputStreamReader()
	javaLang.Load_Lang_Process()
	javaLang.Load_Lang_Process_Builder()
	javaUtil.Load_Util_Hash_Map()

	allocEnvVar := "ALLOC_ENV_VAR"
	tempDir := t.TempDir()
	tempFileName := filepath.Join(tempDir, "alloc_temp_file.txt")

	// Equivalent Java code:
	// ProcessBuilder pb = new ProcessBuilder(command);
	// pb.environment().put(ALLOC_ENV_VAR, tempFileName);
	// pb.redirectErrorStream(true);
	// Process process = pb.start();
	// BufferedReader reader = new BufferedReader(new InputStreamReader(process.getInputStream()));
	//
	// StringBuilder text = new StringBuilder();
	// String line;
	// while ((line = reader.readLine()) != null) {
	//     text.append(line).append('\n');
	// }

	command := []string{"sh", "-c", fmt.Sprintf("echo line_one\necho error_output >&2\necho $%s\necho line_three", allocEnvVar)}

	// 1. ProcessBuilder pb = new ProcessBuilder(command);
	pbClassName := "java/lang/ProcessBuilder"
	pbObj := object.MakeEmptyObjectWithClassName(&pbClassName)
	initPbRes := ghelpers.MethodSignatures["java/lang/ProcessBuilder.<init>([Ljava/lang/String;)V"].GFunction([]any{pbObj, command})
	if initPbRes != nil {
		t.Fatalf("ProcessBuilder.<init> failed: %v", initPbRes)
	}

	// 2. pb.environment().put(ALLOC_ENV_VAR, tempFileName);
	envRes := ghelpers.MethodSignatures["java/lang/ProcessBuilder.environment()Ljava/util/Map;"].GFunction([]any{pbObj})
	envMapObj, ok := envRes.(*object.Object)
	if !ok || envMapObj == nil {
		t.Fatalf("ProcessBuilder.environment() expected *object.Object, got %T (%v)", envRes, envRes)
	}
	putRes := ghelpers.MethodSignatures["java/util/HashMap.put(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"].GFunction([]any{
		envMapObj,
		object.StringObjectFromGoString(allocEnvVar),
		object.StringObjectFromGoString(tempFileName),
	})
	if errBlk, isErr := putRes.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("HashMap.put failed: %v", errBlk)
	}

	// 3. pb.redirectErrorStream(true);
	redirRes := ghelpers.MethodSignatures["java/lang/ProcessBuilder.redirectErrorStream(Z)Ljava/lang/ProcessBuilder;"].GFunction([]any{pbObj, true})
	if redirRes != pbObj {
		t.Fatalf("ProcessBuilder.redirectErrorStream(true) expected self, got %v", redirRes)
	}

	// 4. Process process = pb.start();
	procRes := ghelpers.MethodSignatures["java/lang/ProcessBuilder.start()Ljava/lang/Process;"].GFunction([]any{pbObj})
	procObj, ok := procRes.(*object.Object)
	if !ok || procObj == nil {
		t.Fatalf("ProcessBuilder.start() expected *object.Object, got %T (%v)", procRes, procRes)
	}

	// 5. process.getInputStream()
	inStreamRes := ghelpers.MethodSignatures["java/lang/Process.getInputStream()Ljava/io/InputStream;"].GFunction([]any{procObj})
	inStreamObj, ok := inStreamRes.(*object.Object)
	if !ok || inStreamObj == nil {
		t.Fatalf("Process.getInputStream() expected *object.Object, got %T (%v)", inStreamRes, inStreamRes)
	}

	// 6. new InputStreamReader(process.getInputStream())
	isrClassName := "java/io/InputStreamReader"
	isrObj := object.MakeEmptyObjectWithClassName(&isrClassName)
	initIsrRes := ghelpers.MethodSignatures["java/io/InputStreamReader.<init>(Ljava/io/InputStream;)V"].GFunction([]any{isrObj, inStreamObj})
	if initIsrRes != nil {
		t.Fatalf("InputStreamReader.<init> failed: %v", initIsrRes)
	}

	// 7. BufferedReader reader = new BufferedReader(new InputStreamReader(...))
	brClassName := "java/io/BufferedReader"
	brObj := object.MakeEmptyObjectWithClassName(&brClassName)
	initBrRes := ghelpers.MethodSignatures["java/io/BufferedReader.<init>(Ljava/io/Reader;)V"].GFunction([]any{brObj, isrObj})
	if initBrRes != nil {
		t.Fatalf("BufferedReader.<init> failed: %v", initBrRes)
	}

	// 8. Read lines in a loop:
	// StringBuilder text = new StringBuilder();
	// String line;
	// while ((line = reader.readLine()) != null) {
	//     text.append(line).append('\n');
	// }
	var text strings.Builder
	for {
		lineRes := ghelpers.MethodSignatures["java/io/BufferedReader.readLine()Ljava/lang/String;"].GFunction([]any{brObj})
		if lineRes == object.Null || lineRes == nil {
			break
		}
		if errBlk, isErr := lineRes.(*ghelpers.GErrBlk); isErr {
			t.Fatalf("BufferedReader.readLine returned error: %v", errBlk)
		}
		lineObj, ok := lineRes.(*object.Object)
		if !ok {
			t.Fatalf("BufferedReader.readLine expected *object.Object, got %T (%v)", lineRes, lineRes)
		}
		lineStr := object.GoStringFromStringObject(lineObj)
		text.WriteString(lineStr)
		text.WriteByte('\n')
	}

	output := text.String()
	if !strings.Contains(output, "line_one") {
		t.Errorf("expected output to contain 'line_one', got %q", output)
	}
	if !strings.Contains(output, "error_output") {
		t.Errorf("expected output to contain redirected stderr 'error_output', got %q", output)
	}
	if !strings.Contains(output, tempFileName) {
		t.Errorf("expected output to contain environment variable value %q, got %q", tempFileName, output)
	}
	if !strings.Contains(output, "line_three") {
		t.Errorf("expected output to contain 'line_three', got %q", output)
	}
}
