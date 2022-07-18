package classloader

import (
	"os"
	"path/filepath"
	"testing"
)

var GOOD_JAR_NAME = "hello.jar"
var NO_MANIFEST_JAR_NAME = "nomanifest.jar"

func getJarFileName(name string) (string, error) {
	pwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return filepath.Join(pwd, "..", "..", "testdata", name), nil
}

func getJar(name string, t *testing.T) (*Archive, error) {
	fileName, err := getJarFileName(name)

	if err != nil {
		t.Error("Unable to get jar file", err)
		return nil, err
	}

	return NewJarFile(fileName)
}

func TestGoodJarFile(t *testing.T) {
	jar, err := getJar(GOOD_JAR_NAME, t)

	if err != nil {
		return
	}

	if err := jar.scanArchive(); err != nil {
		t.Error("Error scanning archive", err)
	}
}

func TestManifestParsing(t *testing.T) {
	jar, err := getJar(GOOD_JAR_NAME, t)

	if err != nil {
		return
	}

	if err := jar.scanArchive(); err != nil {
		t.Error("Error scanning archive", err)
		return
	}

	value, ok := jar.manifest["Main-Class"]

	if !ok {
		t.Error("Main-Class attribute should have been there, but wasn't")
	}

	if value != "jacobin.HelloWorld" {
		t.Error("Expected Main-Class to be 'jacobin.HelloWorld', but was " + value)
	}
}

func TestLoadClassSuccess(t *testing.T) {
	jar, err := getJar(GOOD_JAR_NAME, t)

	if err != nil {
		return
	}

	result, err := jar.loadClass("jacobin.HelloWorld")

	if err != nil {
		t.Error("Error loading class", err)
	}

	if !result.Success {
		t.Error("Loading class was not successful")
	}
}

func TestLoadClassDoesNotExist(t *testing.T) {
	jar, err := getJar(NO_MANIFEST_JAR_NAME, t)

	if err != nil {
		return
	}

	_, err = jar.loadClass("jacobin.HelloWorld")

	if err == nil {
		t.Error("Expected error loading class, but didn't get one.")
	}
}
