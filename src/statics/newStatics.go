/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package statics

// === Handling of statics.
// In the JVM, static fields are stored in the mirroring java/lang/Class object,
// which is stored in globals.JlcMap. That is a table whose key is the class name.
//
// For each class, we create a java/lang/Class object, Jlc, and store its static fields
// in the Jlc.Statics table. That Staics table is a map whose key is the name of the
// field concatenated with the field type. For example: timeI or outputLfile.
//
// The corresponding value is a struct, StaticField, which contains field metadata
// and the field's value.
//

// instances of java/lang/Class stored in global.JlcNap
// type Jlc struct {
// 	Lock        sync.RWMutex
// 	Statics     []string          // all static fields
// 	KlassPtr    *classdata.ClData // points back to the class's data in the method area
// 	Initialized bool              // has the class been initialized?
// }
//
// // StaticField is final represenation of a staic; it contains field metadata and the field's
// // value, of course
// type StaticField struct {
// 	FieldMetaData classdata.Field
// 	Value         any
// }
