/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package util

import "strings"

// DecodeModifiedUTF8 decodes Modified UTF-8 null: 0xC0 0x80 → 0x00
func DecodeModifiedUTF8(baIn []byte) []byte {
	var baOut []byte
	for i := 0; i < len(baIn); {
		// Modified UTF-8 null: 0xC0 0x80 → 0x00
		if baIn[i] == 0xC0 && i+1 < len(baIn) && baIn[i+1] == 0x80 {
			baOut = append(baOut, 0x00)
			i += 2
		} else {
			baOut = append(baOut, baIn[i])
			i++
		}
	}
	return baOut
}

// IsClassPartOfJDK accepts a classname and returns true if the classname
// is part of the JDK distribution (includes packages in in early Java releases
// that were subsequently removed, such as CORBA, etc.)
var jdkPrefixes = []string{
	"java.",
	"jdk.",
	"sun.",

	// org.*
	"org.ietf.jgss.",
	"org.w3c.dom.",
	"org.omg.",
	"org.xml.sax.",
}

var jdkPrefixesJavax = []string{
	// javax.* (JDK-owned only; excludes EE/Jakarta-era packages)
	"javax.accessibility.",
	"javax.annotation.processing.",
	"javax.crypto.",
	"javax.imageio.",
	"javax.lang.model.",
	"javax.management.",
	"javax.naming.",
	"javax.net.",
	"javax.print.",
	"javax.rmi.",
	"javax.script.",
	"javax.security.auth.",
	"javax.security.cert.",
	"javax.security.sasl.",
	"javax.smartcardio.",
	"javax.sound.",
	"javax.sql.",
	"javax.swing.",
	"javax.tools.",
	"javax.transaction.xa.",
	"javax.xml.catalog.",
	"javax.xml.crypto.",
	"javax.xml.datatype.",
	"javax.xml.namespace.",
	"javax.xml.parsers.",
	"javax.xml.stream.",
	"javax.xml.transform.",
	"javax.xml.validation.",
	"javax.xml.xpath.",
}

var jdkPrefixesComSun = []string{
	// com.sun.* (JDK-owned only; excludes jna, mail, jersey, xml.bind, istack, ...)
	"com.sun.awt.",
	"com.sun.beans.",
	"com.sun.crypto.provider.",
	"com.sun.imageio.",
	"com.sun.java.accessibility.",
	"com.sun.java.swing.",
	"com.sun.java.util.jar.pack.",
	"com.sun.jarsigner.",
	"com.sun.jdi.",
	"com.sun.jmx.",
	"com.sun.jndi.",
	"com.sun.management.",
	"com.sun.media.sound.",
	"com.sun.naming.",
	"com.sun.net.httpserver.",
	"com.sun.net.ssl.",
	"com.sun.nio.file.",
	"com.sun.nio.sctp.",
	"com.sun.org.apache.",
	"com.sun.rowset.",
	"com.sun.security.auth.",
	"com.sun.security.jgss.",
	"com.sun.security.ntlm.",
	"com.sun.security.sasl.",
	"com.sun.source.",
	"com.sun.swing.",
	"com.sun.tools.",
	"com.sun.xml.internal.",
}

func IsClassPartOfJDK(arg string) bool {
	className := ConvertInternalClassNameToUserFormat(arg)

	for _, prefix := range jdkPrefixes {
		if strings.HasPrefix(className, prefix) {
			return true
		}
	}

	if strings.HasPrefix(className, "javax.") {
		for _, prefix := range jdkPrefixesJavax {
			if strings.HasPrefix(className, prefix) {
				return true
			}
		}
	}

	if strings.HasPrefix(className, "com.sun.") {
		for _, prefix := range jdkPrefixesComSun {
			if strings.HasPrefix(className, prefix) {
				return true
			}
		}
	}
	return false
}
