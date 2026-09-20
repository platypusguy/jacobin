/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"testing"
)

func TestLoad_Util_Zip_ZipInputStream_RegistersAllMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Util_Zip_ZipInputStream()

	if len(ghelpers.MethodSignatures) == 0 {
		t.Fatalf("expected method signatures to be registered")
	}

	clinit, ok := ghelpers.MethodSignatures["java/util/zip/ZipInputStream.<clinit>()V"]
	if !ok {
		t.Fatalf("expected <clinit> to be registered")
	}
	if clinit.ParamSlots != 0 {
		t.Errorf("expected <clinit> ParamSlots 0, got %d", clinit.ParamSlots)
	}
	if clinit.GFunction == nil {
		t.Errorf("expected <clinit> GFunction to be non-nil")
	}

	for key, gm := range ghelpers.MethodSignatures {
		if key == "java/util/zip/ZipInputStream.<clinit>()V" {
			continue
		}
		if gm.GFunction == nil {
			t.Errorf("entry %q has nil GFunction", key)
		}
	}

	expectedSlots := map[string]int{
		"java/util/zip/ZipInputStream.<init>(Ljava/io/InputStream;)V":                             1,
		"java/util/zip/ZipInputStream.<init>(Ljava/io/InputStream;Ljava/nio/charset/Charset;)V":   2,
		"java/util/zip/ZipInputStream.available()I":                                               0,
		"java/util/zip/ZipInputStream.close()V":                                                   0,
		"java/util/zip/ZipInputStream.closeEntry()V":                                              0,
		"java/util/zip/ZipInputStream.createZipEntry(Ljava/lang/String;)Ljava/util/zip/ZipEntry;": 1,
		"java/util/zip/ZipInputStream.getNextEntry()Ljava/util/zip/ZipEntry;":                     0,
		"java/util/zip/ZipInputStream.mark(I)V":                                                   1,
		"java/util/zip/ZipInputStream.markSupported()Z":                                           0,
		"java/util/zip/ZipInputStream.read()I":                                                    0,
		"java/util/zip/ZipInputStream.read([BII)I":                                                3,
		"java/util/zip/ZipInputStream.reset()V":                                                   0,
		"java/util/zip/ZipInputStream.skip(J)J":                                                   1,
	}

	for key, expectedParamSlots := range expectedSlots {
		gm, ok := ghelpers.MethodSignatures[key]
		if !ok {
			t.Errorf("expected entry %q to be registered", key)
			continue
		}
		if gm.ParamSlots != expectedParamSlots {
			t.Errorf("entry %q expected ParamSlots %d, got %d", key, expectedParamSlots, gm.ParamSlots)
		}
	}

	if len(ghelpers.MethodSignatures) != len(expectedSlots)+1 {
		t.Errorf("expected %d entries, got %d", len(expectedSlots)+1, len(ghelpers.MethodSignatures))
	}
}
