/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */
package classloader

import (
    "archive/zip"
    "errors"
    "fmt"
    "io"
    "jacobin/src/trace"
    "path/filepath"
    "strings"
)

// This file contains the code for loading and managing JAR files.

type ResourceType int16
type ArchiveType int16

const (
	TypeResource ResourceType = iota
	TypeClassFile
	TypeManifest
)

type ResourceEntry struct {
	Location string
	Name     string
	Type     ResourceType
}

type Archive struct {
	FilePath     string
	EntryCache   map[string]ResourceEntry
	Manifest     map[string]string
	ClasspathRaw string
	Classpath    []string
}

type LoadResult struct {
	Success       bool
	Data          *[]byte
	ResourceEntry ResourceEntry
}

/*
Open an archive file (jar or zip) and return a handle to it.
*/
func OpenArchive(filePath string) (*Archive, error) {
	archive := new(Archive)
	archive.FilePath = filePath
	archive.EntryCache = make(map[string]ResourceEntry)
	archive.Manifest = make(map[string]string)
	err := archive.scanArchive()

	if err != nil {
		return nil, err
	}

	return archive, err
}

/*
Scan the given archive file (jar or zip) and populate the EntryCache and TypeManifest maps.
*/
func (archive *Archive) scanArchive() error {
	reader, err := zip.OpenReader(archive.FilePath)
	if err != nil {
		trace.Error(fmt.Sprintf("scanArchive: zip.OpenReader(%s) failed, err: %s", archive.FilePath, err.Error()))
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		entry := archive.recordFile(file)
		if entry.Type == TypeManifest {
			if archive.parseManifest(file); err != nil {
				return err
			}
		}
	}

	return nil
}

func (archive *Archive) recordFile(file *zip.File) ResourceEntry {
	fileType := TypeResource
	resourceName := file.Name

	if strings.HasSuffix(file.Name, ".class") {
		fileType = TypeClassFile
		resourceName = strings.ReplaceAll(resourceName, "/", ".")
		resourceName = strings.TrimSuffix(resourceName, ".class")
	} else if file.Name == "META-INF/MANIFEST.MF" {
		fileType = TypeManifest
	}

	entry := ResourceEntry{
		Location: file.Name,
		Name:     resourceName,
		Type:     fileType,
	}

	archive.EntryCache[resourceName] = entry

	return entry
}

func (archive *Archive) parseManifest(file *zip.File) error {

	// Get all the string lines from the archive file.
	rc, err := file.Open()
	defer rc.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("archives.parseManifest file.Open[%s], err: %s", archive.FilePath, err.Error()))
	}
	data, err := io.ReadAll(rc)
	if err != nil {
		if err != nil {
			return errors.New(fmt.Sprintf("archives.parseManifest io.ReadAll[%s], err: %s", archive.FilePath, err.Error()))
		}
	}
	contents := string(data)
	lines := strings.Split(contents, "\n") // supports all OSes

	// Process each string line
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) > 1 {
			archive.Manifest[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return nil
}

func (archive *Archive) hasResource(name string, resourceType ResourceType) bool {
	item, ok := archive.EntryCache[name]
	if !ok {
		return false
	}
	return item.Type == resourceType
}

func (archive *Archive) loadClass(className string) (*LoadResult, error) {
	item, ok := archive.EntryCache[className]
	if !ok {
		err := errors.New(fmt.Sprintf("archives.loadClass archive.EntryCache[%s] failed in archive %s", className, archive.FilePath))
		return nil, err
	}

	if item.Type != TypeClassFile {
		return nil, errors.New(fmt.Sprintf("archives.loadClass file %s in archive %s is not a classfile", className, archive.FilePath))
	}

	reader, err := zip.OpenReader(archive.FilePath)
	defer reader.Close()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("archives.loadClass zip.OpenReader(%s) failed, err: %s", archive.FilePath, err.Error()))
	}

	file, err := reader.Open(item.Location)
	defer file.Close()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("archives.loadClass reader.Open(%s, Location %s) failed, err: %s", archive.FilePath, item.Location, err.Error()))
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("archives.loadClass io.ReadAll(%s, Location %s) failed, err: %s", archive.FilePath, item.Location, err.Error()))
	}

	return &LoadResult{Data: &bytes, Success: true, ResourceEntry: item}, nil
}

func (archive *Archive) getMainClass() string {
	mainClass, exists := archive.Manifest["Main-Class"]

	if exists {
		return mainClass
	} else {
		return ""
	}
}

func (archive *Archive) UpdateArchiveWithClassPath() {

    preSplit, ok := archive.Manifest["Class-Path"]
    classPathArray := []string{archive.FilePath}
    if !ok {
		// There is no Class-Path attribute in the manifest.
		// Therefore, the classpath is just the singleton string = path of this jar.
		archive.ClasspathRaw = preSplit
		archive.Classpath = classPathArray
		return
	}

    // There is a Class-Path attribute in the manifest.
    // Per the JAR File Specification, each entry is space-separated and
    // relative entries are resolved against the directory containing this JAR.
    postSplit := strings.Fields(preSplit)
    baseDir := filepath.Dir(archive.FilePath)
    for _, entry := range postSplit {
        e := strings.TrimSpace(entry)
        if e == "" {
            continue
        }
        // If the entry looks like an absolute path, keep as-is; otherwise join with baseDir.
        // We do not attempt to resolve URLs here; they will be left as-is for higher-level handling if needed.
        if filepath.IsAbs(e) {
            classPathArray = append(classPathArray, e)
        } else {
            classPathArray = append(classPathArray, filepath.Join(baseDir, e))
        }
    }
    archive.ClasspathRaw = preSplit
    archive.Classpath = classPathArray
    return
}
