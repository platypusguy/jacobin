/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import (
	"jacobin/statics"
	"jacobin/types"
	"testing"
)

func TestNewString(t *testing.T) {
	str := *NewString()

	if *str.Klass != "java/lang/String" {
		t.Errorf("Klass should be java/lang/String, got: %s", *str.Klass)
	}

	value := str.FieldTable["value"].Fvalue.([]byte)
	valueStr := string(value)
	if len(valueStr) != 0 {
		t.Errorf("value field should be empty string, got: %s", string(value))
	}

	coder := str.FieldTable["coder"].Fvalue.(int64)
	if coder != 0 && coder != 1 {
		t.Errorf("coder field should be 0 or 1, got: %d", coder)
	}

	hash := str.FieldTable["hash"].Fvalue.(int64)
	if hash != 0 {
		t.Errorf("hash field should be 0, got: %d", hash)
	}

	hashIsZero := str.FieldTable["hashIsZero"].Fvalue.(int64)
	if hash != types.JavaBoolFalse {
		t.Errorf("hashIsZero field should be false, got: %d", hashIsZero)
	}
}

func TestNewStringFromGoString(t *testing.T) {
	statics.LoadStaticsString()

	s := NewStringFromGoString("hello")
	newString := string(s.FieldTable["value"].Fvalue.([]byte))
	if newString != "hello" {
		t.Errorf("expected strint to be 'hello', got: %s", newString)
	}
}

func TestCreateCompactStringFromGoString(t *testing.T) {
	goString := "You say hello!"
	s := CreateCompactStringFromGoString(&goString)
	compactString := string(s.FieldTable["value"].Fvalue.([]byte))

	if compactString != "You say hello!" {
		t.Errorf("expected string to be 'You say hello!', got: %s",
			compactString)
	}
}
