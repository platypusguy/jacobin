package classloader

import (
	"os"
	"path"
	"strings"
	"testing"
)

func TestJmodFile(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Error("Unable to get cwd")
		return
	}

	jmodFileName := path.Join(pwd, "..", "..", "testdata", "jmod", "jacobin.jmod")

	jmodFile, err := os.Open(jmodFileName)
	if err != nil {
		t.Error("Unable to open jmod file", err)
		return
	}

	jmod := Jmod{*jmodFile}

	filesFound := make(map[string]any, 10)

	var empty struct{}

	jmod.Walk(func(bytes []byte, filename string) error {
		fname := strings.Split(filename, "+")[1]
		filesFound[fname] = empty
		return nil
	})

	if _, ok := filesFound["classes/org/jacobin/test/Hello.class"]; !ok {
		t.Error("Expected org.jacobin.test.Hello, but it wasn't there.")
	}

	if _, ok := filesFound["classes/module-info.class"]; ok {
		t.Error("Didn't expect module-info, but it was there.")
	}
}

func TestJmodFileNoClasslist(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Error("Unable to get cwd")
		return
	}

	jmodFileName := path.Join(pwd, "..", "..", "testdata", "jmod", "jacobinfull.jmod")

	jmodFile, err := os.Open(jmodFileName)
	if err != nil {
		t.Error("Unable to open jmod file", err)
		return
	}

	jmod := Jmod{*jmodFile}

	filesFound := make(map[string]any, 10)

	var empty struct{}

	jmod.Walk(func(bytes []byte, filename string) error {
		fname := strings.Split(filename, "+")[1]
		filesFound[fname] = empty
		return nil
	})

	if _, ok := filesFound["classes/org/jacobin/test/Hello.class"]; !ok {
		t.Error("Expected org.jacobin.test.Hello, but it wasn't there.")
	}

	if _, ok := filesFound["classes/module-info.class"]; !ok {
		t.Error("Expected module-info, but it wasn't there.")
	}
}
