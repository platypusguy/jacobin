/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */
package classloader

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"jacobin/log"
	"strings"
)

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
		_ = log.Log("Error: Invalid or corrupt jarfile "+archive.Filename, log.SEVERE)
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
	rc, err := file.Open()

	defer rc.Close()

	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(rc)

	contents := string(data)

	lines := strings.Split(contents, "\r\n")

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
		err := errors.New(fmt.Sprintf("Unable to load class %s in archive %s", className, archive.Filename))
		return nil, err
	}

	if item.Type != ClassFile {
		return nil, errors.New(fmt.Sprintf("Class %s in archive %s is not a classfile", className, archive.Filename))
	}

	reader, err := zip.OpenReader(archive.Filename)

	defer reader.Close()

	if err != nil {
		return nil, err
	}

	file, err := reader.Open(item.Location)

	defer file.Close()

	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
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
