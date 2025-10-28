/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */
package classloader

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"jacobin/src/trace"
	"strings"
)

// This file contains the code for loading and managing JAR files.

type ResourceType int16
type ArchiveType int16

const (
	Resource ResourceType = iota
	ClassFile
	Manifest
)

type ResourceEntry struct {
	Location string
	Name     string
	Type     ResourceType
}

type Archive struct {
	Filename   string
	entryCache map[string]ResourceEntry
	manifest   map[string]string
}

type LoadResult struct {
	Success       bool
	Data          *[]byte
	ResourceEntry ResourceEntry
}

func NewJarFile(filename string) (*Archive, error) {
	jarFile := new(Archive)
	jarFile.Filename = filename
	jarFile.entryCache = make(map[string]ResourceEntry)
	jarFile.manifest = make(map[string]string)
	err := jarFile.scanArchive()

	if err != nil {
		return nil, err
	}

	return jarFile, err
}

func (archive *Archive) scanArchive() error {
	reader, err := zip.OpenReader(archive.Filename)
	if reader != nil {
		defer reader.Close()
	}

	if reader == nil || err != nil {
		trace.Error("Invalid, corrupt, or inaccessible jarfile " + archive.Filename)
		return err
	}

	for _, file := range reader.File {
		entry := archive.recordFile(file)
		if entry.Type == Manifest {
			if archive.parseManifest(file); err != nil {
				return err
			}
		}
	}

	return nil
}

func (archive *Archive) recordFile(file *zip.File) ResourceEntry {
	fileType := Resource
	resourceName := file.Name

	if strings.HasSuffix(file.Name, ".class") {
		fileType = ClassFile
		resourceName = strings.ReplaceAll(resourceName, "/", ".")
		resourceName = strings.TrimSuffix(resourceName, ".class")
	} else if file.Name == "META-INF/MANIFEST.MF" {
		fileType = Manifest
	}

	entry := ResourceEntry{
		Location: file.Name,
		Name:     resourceName,
		Type:     fileType,
	}

	archive.entryCache[resourceName] = entry

	return entry
}

func (archive *Archive) parseManifest(file *zip.File) error {

	// Get all the string lines from the archive file.
	rc, err := file.Open()
	defer rc.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("archives.parseManifest file.Open[%s], err: %s", archive.Filename, err.Error()))
	}
	data, err := io.ReadAll(rc)
	if err != nil {
		if err != nil {
			return errors.New(fmt.Sprintf("archives.parseManifest io.ReadAll[%s], err: %s", archive.Filename, err.Error()))
		}
	}
	contents := string(data)
	lines := strings.Split(contents, "\n") // supports all OSes

	// Process each string line
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) > 1 {
			archive.manifest[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return nil
}

func (archive *Archive) hasResource(name string, resourceType ResourceType) bool {
	item, ok := archive.entryCache[name]
	if !ok {
		return false
	}
	return item.Type == resourceType
}

func (archive *Archive) loadClass(className string) (*LoadResult, error) {
	item, ok := archive.entryCache[className]
	if !ok {
		err := errors.New(fmt.Sprintf("archives.loadClass archive.entryCache[%s] failed in archive %s", className, archive.Filename))
		return nil, err
	}

	if item.Type != ClassFile {
		return nil, errors.New(fmt.Sprintf("archives.loadClass file %s in archive %s is not a classfile", className, archive.Filename))
	}

	reader, err := zip.OpenReader(archive.Filename)
	defer reader.Close()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("archives.loadClass zip.OpenReader(%s) failed, err: %s", archive.Filename, err.Error()))
	}

	file, err := reader.Open(item.Location)
	defer file.Close()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("archives.loadClass reader.Open(%s, Location %s) failed, err: %s", archive.Filename, item.Location, err.Error()))
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("archives.loadClass io.ReadAll(%s, Location %s) failed, err: %s", archive.Filename, item.Location, err.Error()))
	}

	return &LoadResult{Data: &bytes, Success: true, ResourceEntry: item}, nil
}

func (archive *Archive) getMainClass() string {
	mainClass, exists := archive.manifest["Main-Class"]

	if exists {
		return mainClass
	} else {
		return ""
	}
}

func (archive *Archive) getClassPath() []string {

	preSplit, ok := archive.manifest["Class-Path"]
	classPath := []string{archive.Filename}
	if !ok {
		// Classpath is just the path of this jar.
		return classPath
	}

	// Append the manifest classpath to the jar path, giving the full classpath.
	postSplit := strings.Split(preSplit, " ")
	for ndx := 0; ndx < len(postSplit); ndx++ {
		classPath = append(classPath, postSplit[ndx])
	}
	return classPath
}
