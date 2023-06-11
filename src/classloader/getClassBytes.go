package classloader

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"jacobin/globals"
	"jacobin/log"
	"os"
)

const ExpectedMagicNumber = 0x4A4D

func GetClassBytes(jmodFileName string, className string) ([]byte, error) {

	global := globals.GetGlobalRef()
	jmodPath := global.JavaHome + string(os.PathSeparator) + "jmods" + string(os.PathSeparator) + jmodFileName
	classFileName := "classes/" + className + ".class"

	msg := fmt.Sprintf("GetClassBytes: jmodPath %s, className %s\n", jmodPath, className)
	log.Log(msg, log.TRACE_INST)

	// Read entire jmod file contents
	jmodBytes, err := os.ReadFile(jmodPath)
	if err != nil {
		msg = fmt.Sprintf("GetClassBytes: os.ReadFile(%s) failed", jmodPath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		return nil, err
	}

	// Validate the file's magic number
	fileMagicNumber := binary.BigEndian.Uint16(jmodBytes[:2])
	if fileMagicNumber != ExpectedMagicNumber {
		msg = fmt.Sprintf("GetClassBytes: fileMagicNumber != ExpectedMagicNumber in jmod file %s\n", jmodPath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		return nil, err
	}

	// Skip over the jmod header so that it is recognized as a ZIP file
	ioReader := bytes.NewReader(jmodBytes[4:])

	// Prepare the reader for the zip archive
	zipReader, err := zip.NewReader(ioReader, int64(len(jmodBytes)-4))
	if err != nil {
		msg = fmt.Sprintf("GetClassBytes: zip.NewReader(%s) failed", jmodPath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		return nil, err
	}

	// Open the file within the zip archive
	fileHandle, err := zipReader.Open(classFileName)
	if err != nil {
		msg = fmt.Sprintf("GetClassBytes: zipReader.Open(class file %s in jmod file %s) failed", classFileName, jmodPath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		return nil, err
	}

	// Read entire class file contents
	bytes, err := io.ReadAll(fileHandle)
	if err != nil {
		msg = fmt.Sprintf("GetClassBytes: os.ReadAll(class file %s in jmod file %s) failed", classFileName, jmodPath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		return nil, err
	}

	// Success!
	return bytes, nil

}
