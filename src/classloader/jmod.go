package classloader

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"jacobin/log"
	"os"
	"strings"
)

type WalkEntryFunc func(bytes []byte, filename string) error

// Jmod Holds the file referring to a Java Module (JMOD)
type Jmod struct {
	File os.File
}

// Walk Walks a JMOD file and invokes `walk` for all classes found in the classlist
func (j *Jmod) Walk(walk WalkEntryFunc) error {
	b, err := ioutil.ReadFile(j.File.Name())
	if err != nil {
		return err
	}

	// Skip over the JMOD header so that it is recognized as a ZIP file
	offsetReader := bytes.NewReader(b[4:])

	r, err := zip.NewReader(offsetReader, int64(len(b)-4))
	if err != nil {
		log.Log(err.Error(), log.WARNING)
		return err
	}

	classSet := getClasslist(*r)

	useClassSet := len(classSet) > 0

	for _, f := range r.File {
		if !strings.HasPrefix(f.Name, "classes") {
			continue
		}

		classFileName := strings.Replace(f.Name, "classes/", "", 1)

		if useClassSet {
			_, ok := classSet[classFileName]
			if !ok {
				continue
			}
		} else {
			if !strings.HasSuffix(f.Name, ".class") {
				continue
			}
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		b, err := ioutil.ReadAll(rc)
		if err != nil {
			return err
		}

		walk(b, j.File.Name()+"+"+f.Name)

		_ = rc.Close()
	}

	return nil
}

// Returns lib/classlist from the JMOD file, returning an empty map if the classlist cannot be found or read
func getClasslist(reader zip.Reader) map[string]struct{} {
	classSet := make(map[string]struct{})

	classlist, err := reader.Open("lib/classlist")
	if err != nil {
		log.Log(err.Error(), log.CLASS)
		log.Log("Unable to read lib/classlist from jmod file. Loading all classes in jmod file.", log.CLASS)
		return classSet
	}

	classlistContent, err := ioutil.ReadAll(classlist)
	if err != nil {
		log.Log(err.Error(), log.CLASS)
		log.Log("Unable to read lib/classlist from jmod file. Loading all classes in jmod file.", log.CLASS)
		return classSet
	}

	classes := strings.Split(string(classlistContent), "\n")

	var empty struct{}

	for _, c := range classes {
		classSet[c+".class"] = empty
	}

	return classSet
}
