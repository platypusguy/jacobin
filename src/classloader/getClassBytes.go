package classloader

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"io"
	"jacobin/globals"
	"log"
	"os"
)

const ExpectedMagicNumber = 0x4A4D

func getClassBytes(jmodFileName string, classFileName string) []byte {

	global := globals.GetGlobalRef()
	jmodPath := global.JavaHome + string(os.PathSeparator) + "jmods" + string(os.PathSeparator) + jmodFileName

	// Read entire jmod file contents
	jmodBytes, err := os.ReadFile(jmodPath)
	if err != nil {
		log.Fatalf("getClassBytes: os.ReadFile(%s) failed:\n%s\n", jmodPath, err.Error())
	}

	// Validate the file's magic number
	fileMagicNumber := binary.BigEndian.Uint16(jmodBytes[:2])
	if fileMagicNumber != ExpectedMagicNumber {
		log.Fatalf("getClassBytes: fileMagicNumber != ExpectedMagicNumber in jmod file %s\n", jmodPath)
	}

	// Skip over the jmod header so that it is recognized as a ZIP file
	ioReader := bytes.NewReader(jmodBytes[4:])

	// Prepare the reader for the zip archive
	zipReader, err := zip.NewReader(ioReader, int64(len(jmodBytes)-4))
	if err != nil {
		log.Fatalf("getClassBytes: zip.NewReader(%s) failed:\n%s\n", jmodPath, err.Error())
	}

	// Open the file within the zip archive
	fileHandle, err := zipReader.Open(classFileName)
	if err != nil {
		log.Fatalf("getClassBytes: zipReader.Open(%s in %s) failed:\n%s\n", classFileName, jmodPath, err.Error())
	}

	// Read entire class file contents
	bytes, err := io.ReadAll(fileHandle)
	if err != nil {
		log.Fatalf("getClassBytes: os.ReadAll(%s in %s) failed:\n%s\n", classFileName, jmodPath, err.Error())
	}

	return bytes

}
