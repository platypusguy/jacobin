/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"strings"
	"time"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
)

var zipEntryClassName = "java/util/zip/ZipEntry"
var fileTimeClassName = "java/nio/file/attribute/FileTime"

// Compression method constants
const (
	ZipEntryStored   int64 = 0
	ZipEntryDeflated int64 = 8
)

// Load_Util_Zip_ZipEntry registers all method signatures for java.util.zip.ZipEntry.
func Load_Util_Zip_ZipEntry() {
	// <clinit>()V
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryClinit,
		}

	// <init>(Ljava/lang/String;)V
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipEntryInitString,
		}

	// <init>(Ljava/util/zip/ZipEntry;)V
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.<init>(Ljava/util/zip/ZipEntry;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipEntryInitCopy,
		}

	// getName()Ljava/lang/String;
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.getName()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryGetName,
		}

	// setTime(J)V
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.setTime(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipEntrySetTime,
		}

	// getTime()J
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.getTime()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryGetTime,
		}

	// setTimeLocal(Ljava/time/LocalDateTime;)V
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.setTimeLocal(Ljava/time/LocalDateTime;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipEntrySetTimeLocal,
		}

	// getTimeLocal()Ljava/time/LocalDateTime;
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.getTimeLocal()Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryGetTimeLocal,
		}

	// setLastModifiedTime(Ljava/nio/file/attribute/FileTime;)Ljava/util/zip/ZipEntry;
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.setLastModifiedTime(Ljava/nio/file/attribute/FileTime;)Ljava/util/zip/ZipEntry;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipEntrySetLastModifiedTime,
		}

	// getLastModifiedTime()Ljava/nio/file/attribute/FileTime;
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.getLastModifiedTime()Ljava/nio/file/attribute/FileTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryGetLastModifiedTime,
		}

	// setLastAccessTime(Ljava/nio/file/attribute/FileTime;)Ljava/util/zip/ZipEntry;
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.setLastAccessTime(Ljava/nio/file/attribute/FileTime;)Ljava/util/zip/ZipEntry;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipEntrySetLastAccessTime,
		}

	// getLastAccessTime()Ljava/nio/file/attribute/FileTime;
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.getLastAccessTime()Ljava/nio/file/attribute/FileTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryGetLastAccessTime,
		}

	// setCreationTime(Ljava/nio/file/attribute/FileTime;)Ljava/util/zip/ZipEntry;
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.setCreationTime(Ljava/nio/file/attribute/FileTime;)Ljava/util/zip/ZipEntry;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipEntrySetCreationTime,
		}

	// getCreationTime()Ljava/nio/file/attribute/FileTime;
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.getCreationTime()Ljava/nio/file/attribute/FileTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryGetCreationTime,
		}

	// setSize(J)V
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.setSize(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipEntrySetSize,
		}

	// getSize()J
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.getSize()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryGetSize,
		}

	// getCompressedSize()J
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.getCompressedSize()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryGetCompressedSize,
		}

	// setCompressedSize(J)V
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.setCompressedSize(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipEntrySetCompressedSize,
		}

	// setCrc(J)V
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.setCrc(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipEntrySetCrc,
		}

	// getCrc()J
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.getCrc()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryGetCrc,
		}

	// setMethod(I)V
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.setMethod(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipEntrySetMethod,
		}

	// getMethod()I
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.getMethod()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryGetMethod,
		}

	// setExtra([B)V
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.setExtra([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipEntrySetExtra,
		}

	// getExtra()[B
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.getExtra()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryGetExtra,
		}

	// setComment(Ljava/lang/String;)V
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.setComment(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipEntrySetComment,
		}

	// getComment()Ljava/lang/String;
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.getComment()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryGetComment,
		}

	// isDirectory()Z
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.isDirectory()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryIsDirectory,
		}

	// toString()Ljava/lang/String;
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryToString,
		}

	// hashCode()I
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryHashCode,
		}

	// clone()Ljava/lang/Object;
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.clone()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipEntryClone,
		}

	// equals(Ljava/lang/Object;)Z
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipEntryEquals,
		}

	// Standard Object lifecycle & concurrency methods
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.getClass()Ljava/lang/Class;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.finalize()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.notify()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.notifyAll()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.wait()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.wait(J)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/ZipEntry.wait(JI)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
}

func zipEntryClinit(params []any) any {
	_ = statics.AddStatic(zipEntryClassName+".STORED", statics.Static{
		Type:  types.Int,
		Value: ZipEntryStored,
	})
	_ = statics.AddStatic(zipEntryClassName+".DEFLATED", statics.Static{
		Type:  types.Int,
		Value: ZipEntryDeflated,
	})
	return nil
}

// CreateZipEntryObject creates an initialized ZipEntry instance in Go.
func CreateZipEntryObject(name string) (*object.Object, *ghelpers.GErrBlk) {
	if len([]byte(name)) > 0xFFFF {
		return nil, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "entry name is too long")
	}
	obj := object.MakeEmptyObjectWithClassName(&zipEntryClassName)
	strObj := object.StringObjectFromGoString(name)
	initZipEntryFields(obj, strObj)
	return obj, nil
}

func initZipEntryFields(this *object.Object, nameObj *object.Object) {
	this.FieldTable["name"] = object.Field{Ftype: types.StringClassRef, Fvalue: nameObj}
	this.FieldTable["time"] = object.Field{Ftype: types.Long, Fvalue: int64(-1)}
	this.FieldTable["mtime"] = object.Field{Ftype: "Ljava/nio/file/attribute/FileTime;", Fvalue: object.Null}
	this.FieldTable["atime"] = object.Field{Ftype: "Ljava/nio/file/attribute/FileTime;", Fvalue: object.Null}
	this.FieldTable["ctime"] = object.Field{Ftype: "Ljava/nio/file/attribute/FileTime;", Fvalue: object.Null}
	this.FieldTable["crc"] = object.Field{Ftype: types.Long, Fvalue: int64(-1)}
	this.FieldTable["size"] = object.Field{Ftype: types.Long, Fvalue: int64(-1)}
	this.FieldTable["csize"] = object.Field{Ftype: types.Long, Fvalue: int64(-1)}
	this.FieldTable["method"] = object.Field{Ftype: types.Int, Fvalue: int64(-1)}
	this.FieldTable["extra"] = object.Field{Ftype: "[B", Fvalue: object.Null}
	this.FieldTable["comment"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.Null}
	this.FieldTable["flag"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
}

func zipEntryInitString(params []any) any {
	if len(params) < 2 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.<init>: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.<init>: this is null")
	}

	if params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "entry name is null")
	}
	nameObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(nameObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "entry name is null")
	}

	nameStr := object.GoStringFromStringObject(nameObj)
	if len([]byte(nameStr)) > 0xFFFF {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "entry name is too long")
	}

	initZipEntryFields(this, nameObj)
	return nil
}

func zipEntryInitCopy(params []any) any {
	if len(params) < 2 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.<init>: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.<init>: this is null")
	}

	if params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "entry is null")
	}
	srcObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(srcObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "entry is null")
	}

	// Copy fields from srcObj
	for k, v := range srcObj.FieldTable {
		if k == "extra" {
			if v.Fvalue != nil && !object.IsNull(v.Fvalue) {
				if extraObj, ok := v.Fvalue.(*object.Object); ok {
					this.FieldTable[k] = object.Field{Ftype: "[B", Fvalue: cloneByteArrayObject(extraObj)}
					continue
				}
			}
			this.FieldTable[k] = object.Field{Ftype: "[B", Fvalue: object.Null}
		} else {
			this.FieldTable[k] = v
		}
	}

	return nil
}

func zipEntryGetName(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getName: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getName: this is null")
	}

	if fld, exists := this.FieldTable["name"]; exists && fld.Fvalue != nil {
		return fld.Fvalue
	}
	return object.StringObjectFromGoString("")
}

func zipEntrySetTime(params []any) any {
	if len(params) < 2 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setTime: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setTime: this is null")
	}

	timeVal, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid time value")
	}

	this.FieldTable["time"] = object.Field{Ftype: types.Long, Fvalue: timeVal}
	this.FieldTable["mtime"] = object.Field{Ftype: "Ljava/nio/file/attribute/FileTime;", Fvalue: object.Null}
	return nil
}

func zipEntryGetTime(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getTime: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getTime: this is null")
	}

	if mtimeFld, exists := this.FieldTable["mtime"]; exists && mtimeFld.Fvalue != nil && !object.IsNull(mtimeFld.Fvalue) {
		if ftObj, ok := mtimeFld.Fvalue.(*object.Object); ok {
			if ms, ok := fileTimeToMillisHelper(ftObj); ok {
				return ms
			}
		}
	}

	if timeFld, exists := this.FieldTable["time"]; exists && timeFld.Fvalue != nil {
		return timeFld.Fvalue.(int64)
	}

	return int64(-1)
}

func zipEntrySetTimeLocal(params []any) any {
	if len(params) < 2 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setTimeLocal: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setTimeLocal: this is null")
	}

	if params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "time is null")
	}
	ldtObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(ldtObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "time is null")
	}

	ms, ok := millisFromLocalDateTime(ldtObj)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid LocalDateTime")
	}

	ft := makeFileTimeFromMillis(ms)
	this.FieldTable["mtime"] = object.Field{Ftype: "Ljava/nio/file/attribute/FileTime;", Fvalue: ft}
	this.FieldTable["time"] = object.Field{Ftype: types.Long, Fvalue: ms}
	return nil
}

func zipEntryGetTimeLocal(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getTimeLocal: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getTimeLocal: this is null")
	}

	var ms int64 = -1
	if mtimeFld, exists := this.FieldTable["mtime"]; exists && mtimeFld.Fvalue != nil && !object.IsNull(mtimeFld.Fvalue) {
		if ftObj, ok := mtimeFld.Fvalue.(*object.Object); ok {
			if m, ok := fileTimeToMillisHelper(ftObj); ok {
				ms = m
			}
		}
	}

	if ms == -1 {
		if timeFld, exists := this.FieldTable["time"]; exists && timeFld.Fvalue != nil {
			ms = timeFld.Fvalue.(int64)
		}
	}

	if ms == -1 {
		return object.Null
	}

	return localDateTimeFromMillis(ms)
}

func zipEntrySetLastModifiedTime(params []any) any {
	if len(params) < 2 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setLastModifiedTime: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setLastModifiedTime: this is null")
	}

	if params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "lastModifiedTime is null")
	}
	ftObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(ftObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "lastModifiedTime is null")
	}

	this.FieldTable["mtime"] = object.Field{Ftype: "Ljava/nio/file/attribute/FileTime;", Fvalue: ftObj}
	if ms, ok := fileTimeToMillisHelper(ftObj); ok {
		this.FieldTable["time"] = object.Field{Ftype: types.Long, Fvalue: ms}
	}
	return this
}

func zipEntryGetLastModifiedTime(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getLastModifiedTime: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getLastModifiedTime: this is null")
	}

	if mtimeFld, exists := this.FieldTable["mtime"]; exists && mtimeFld.Fvalue != nil && !object.IsNull(mtimeFld.Fvalue) {
		return mtimeFld.Fvalue
	}

	if timeFld, exists := this.FieldTable["time"]; exists && timeFld.Fvalue != nil {
		ms := timeFld.Fvalue.(int64)
		if ms != -1 {
			return makeFileTimeFromMillis(ms)
		}
	}

	return object.Null
}

func zipEntrySetLastAccessTime(params []any) any {
	if len(params) < 2 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setLastAccessTime: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setLastAccessTime: this is null")
	}

	if params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "lastAccessTime is null")
	}
	ftObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(ftObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "lastAccessTime is null")
	}

	this.FieldTable["atime"] = object.Field{Ftype: "Ljava/nio/file/attribute/FileTime;", Fvalue: ftObj}
	return this
}

func zipEntryGetLastAccessTime(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getLastAccessTime: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getLastAccessTime: this is null")
	}

	if atimeFld, exists := this.FieldTable["atime"]; exists && atimeFld.Fvalue != nil && !object.IsNull(atimeFld.Fvalue) {
		return atimeFld.Fvalue
	}
	return object.Null
}

func zipEntrySetCreationTime(params []any) any {
	if len(params) < 2 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setCreationTime: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setCreationTime: this is null")
	}

	if params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "creationTime is null")
	}
	ftObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(ftObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "creationTime is null")
	}

	this.FieldTable["ctime"] = object.Field{Ftype: "Ljava/nio/file/attribute/FileTime;", Fvalue: ftObj}
	return this
}

func zipEntryGetCreationTime(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getCreationTime: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getCreationTime: this is null")
	}

	if ctimeFld, exists := this.FieldTable["ctime"]; exists && ctimeFld.Fvalue != nil && !object.IsNull(ctimeFld.Fvalue) {
		return ctimeFld.Fvalue
	}
	return object.Null
}

func zipEntrySetSize(params []any) any {
	if len(params) < 2 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setSize: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setSize: this is null")
	}

	sizeVal, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid size value")
	}
	if sizeVal < 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid entry size")
	}

	this.FieldTable["size"] = object.Field{Ftype: types.Long, Fvalue: sizeVal}
	return nil
}

func zipEntryGetSize(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getSize: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getSize: this is null")
	}

	if sizeFld, exists := this.FieldTable["size"]; exists && sizeFld.Fvalue != nil {
		return sizeFld.Fvalue.(int64)
	}
	return int64(-1)
}

func zipEntrySetCompressedSize(params []any) any {
	if len(params) < 2 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setCompressedSize: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setCompressedSize: this is null")
	}

	csizeVal, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid csize value")
	}

	this.FieldTable["csize"] = object.Field{Ftype: types.Long, Fvalue: csizeVal}
	return nil
}

func zipEntryGetCompressedSize(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getCompressedSize: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getCompressedSize: this is null")
	}

	if csizeFld, exists := this.FieldTable["csize"]; exists && csizeFld.Fvalue != nil {
		return csizeFld.Fvalue.(int64)
	}
	return int64(-1)
}

func zipEntrySetCrc(params []any) any {
	if len(params) < 2 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setCrc: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setCrc: this is null")
	}

	crcVal, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid crc value")
	}
	if crcVal < 0 || crcVal > 0xFFFFFFFF {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid entry crc-32")
	}

	this.FieldTable["crc"] = object.Field{Ftype: types.Long, Fvalue: crcVal}
	return nil
}

func zipEntryGetCrc(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getCrc: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getCrc: this is null")
	}

	if crcFld, exists := this.FieldTable["crc"]; exists && crcFld.Fvalue != nil {
		return crcFld.Fvalue.(int64)
	}
	return int64(-1)
}

func zipEntrySetMethod(params []any) any {
	if len(params) < 2 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setMethod: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setMethod: this is null")
	}

	methodVal, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid method value")
	}
	if methodVal != ZipEntryStored && methodVal != ZipEntryDeflated {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid compression method")
	}

	this.FieldTable["method"] = object.Field{Ftype: types.Int, Fvalue: methodVal}
	return nil
}

func zipEntryGetMethod(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getMethod: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getMethod: this is null")
	}

	if methodFld, exists := this.FieldTable["method"]; exists && methodFld.Fvalue != nil {
		return methodFld.Fvalue.(int64)
	}
	return int64(-1)
}

func zipEntrySetExtra(params []any) any {
	if len(params) < 2 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setExtra: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setExtra: this is null")
	}

	if params[1] == nil || object.IsNull(params[1]) {
		this.FieldTable["extra"] = object.Field{Ftype: "[B", Fvalue: object.Null}
		return nil
	}

	extraObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(extraObj) {
		this.FieldTable["extra"] = object.Field{Ftype: "[B", Fvalue: object.Null}
		return nil
	}

	bytes, ok := getByteArrayFromObject(extraObj)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid extra field")
	}
	if len(bytes) > 0xFFFF {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid extra field length")
	}

	this.FieldTable["extra"] = object.Field{Ftype: "[B", Fvalue: cloneByteArrayObject(extraObj)}
	return nil
}

func zipEntryGetExtra(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getExtra: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getExtra: this is null")
	}

	if extraFld, exists := this.FieldTable["extra"]; exists && extraFld.Fvalue != nil && !object.IsNull(extraFld.Fvalue) {
		if extraObj, ok := extraFld.Fvalue.(*object.Object); ok {
			return cloneByteArrayObject(extraObj)
		}
	}
	return object.Null
}

func zipEntrySetComment(params []any) any {
	if len(params) < 2 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setComment: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.setComment: this is null")
	}

	if params[1] == nil || object.IsNull(params[1]) {
		this.FieldTable["comment"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.Null}
		return nil
	}

	commentObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(commentObj) {
		this.FieldTable["comment"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.Null}
		return nil
	}

	commentStr := object.GoStringFromStringObject(commentObj)
	if len([]byte(commentStr)) > 0xFFFF {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid entry comment length")
	}

	this.FieldTable["comment"] = object.Field{Ftype: types.StringClassRef, Fvalue: commentObj}
	return nil
}

func zipEntryGetComment(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getComment: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.getComment: this is null")
	}

	if commFld, exists := this.FieldTable["comment"]; exists && commFld.Fvalue != nil && !object.IsNull(commFld.Fvalue) {
		return commFld.Fvalue
	}
	return object.Null
}

func zipEntryIsDirectory(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.isDirectory: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.isDirectory: this is null")
	}

	if nameFld, exists := this.FieldTable["name"]; exists && nameFld.Fvalue != nil && !object.IsNull(nameFld.Fvalue) {
		if strObj, ok := nameFld.Fvalue.(*object.Object); ok {
			nameStr := object.GoStringFromStringObject(strObj)
			if strings.HasSuffix(nameStr, "/") {
				return types.JavaBoolTrue
			}
		}
	}
	return types.JavaBoolFalse
}

func zipEntryToString(params []any) any {
	return zipEntryGetName(params)
}

func zipEntryHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.hashCode: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.hashCode: this is null")
	}

	if nameFld, exists := this.FieldTable["name"]; exists && nameFld.Fvalue != nil && !object.IsNull(nameFld.Fvalue) {
		if strObj, ok := nameFld.Fvalue.(*object.Object); ok {
			nameStr := object.GoStringFromStringObject(strObj)
			return int64(javaStringHashCode(nameStr))
		}
	}
	return int64(0)
}

func zipEntryClone(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.clone: this is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok || object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipEntry.clone: this is null")
	}

	cloneObj := object.MakeEmptyObjectWithClassName(&zipEntryClassName)
	for k, v := range this.FieldTable {
		if k == "extra" {
			if v.Fvalue != nil && !object.IsNull(v.Fvalue) {
				if extraObj, ok := v.Fvalue.(*object.Object); ok {
					cloneObj.FieldTable[k] = object.Field{Ftype: "[B", Fvalue: cloneByteArrayObject(extraObj)}
					continue
				}
			}
			cloneObj.FieldTable[k] = object.Field{Ftype: "[B", Fvalue: object.Null}
		} else {
			cloneObj.FieldTable[k] = v
		}
	}

	return cloneObj
}

func zipEntryEquals(params []any) any {
	if len(params) < 2 || params[0] == nil {
		return types.JavaBoolFalse
	}
	if params[0] == params[1] {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Helpers

func getByteArrayFromObject(obj *object.Object) ([]types.JavaByte, bool) {
	if obj == nil || object.IsNull(obj) {
		return nil, false
	}
	fld, ok := obj.FieldTable["value"]
	if !ok {
		return nil, false
	}
	switch v := fld.Fvalue.(type) {
	case []types.JavaByte:
		return v, true
	case []uint8:
		return object.JavaByteArrayFromGoByteArray(v), true
	default:
		return nil, false
	}
}

func cloneByteArrayObject(obj *object.Object) *object.Object {
	if obj == nil || object.IsNull(obj) {
		return object.Null
	}
	bytes, ok := getByteArrayFromObject(obj)
	if !ok {
		return object.Null
	}
	cloned := make([]types.JavaByte, len(bytes))
	copy(cloned, bytes)
	return object.MakePrimitiveObject("[B", types.JavaByteArray, cloned)
}

func makeFileTimeFromMillis(ms int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&fileTimeClassName)
	obj.FieldTable["value"] = object.Field{Ftype: types.Long, Fvalue: ms}
	obj.FieldTable["unit"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.StringObjectFromGoString("MILLISECONDS")}
	obj.FieldTable["unitName"] = object.Field{Ftype: "string", Fvalue: "MILLISECONDS"}
	obj.FieldTable["fromInstant"] = object.Field{Ftype: types.Bool, Fvalue: types.JavaBoolFalse}
	return obj
}

func fileTimeToMillisHelper(ft *object.Object) (int64, bool) {
	if ft == nil || object.IsNull(ft) {
		return -1, false
	}
	if fld, ok := ft.FieldTable["fromInstant"]; ok && fld.Fvalue.(int64) == types.JavaBoolTrue {
		if instFld, ok := ft.FieldTable["instant"]; ok {
			if instObj, ok := instFld.Fvalue.(*object.Object); ok && instObj != nil && !object.IsNull(instObj) {
				secFld := instObj.FieldTable["seconds"]
				nanoFld := instObj.FieldTable["nanos"]
				sec := secFld.Fvalue.(int64)
				nano := nanoFld.Fvalue.(int64)
				return sec*1000 + nano/1000000, true
			}
		}
	}
	valFld, okVal := ft.FieldTable["value"]
	if !okVal {
		return -1, false
	}
	val := valFld.Fvalue.(int64)
	unitName := "MILLISECONDS"
	if uFld, ok := ft.FieldTable["unitName"]; ok {
		if s, ok := uFld.Fvalue.(string); ok && s != "" {
			unitName = s
		}
	} else if uFld, ok := ft.FieldTable["unit"]; ok {
		if strObj, ok := uFld.Fvalue.(*object.Object); ok && strObj != nil && !object.IsNull(strObj) {
			unitName = object.GoStringFromStringObject(strObj)
		}
	}

	switch unitName {
	case "NANOSECONDS":
		return val / 1000000, true
	case "MICROSECONDS":
		return val / 1000, true
	case "MILLISECONDS":
		return val, true
	case "SECONDS":
		return val * 1000, true
	case "MINUTES":
		return val * 60000, true
	case "HOURS":
		return val * 3600000, true
	case "DAYS":
		return val * 86400000, true
	default:
		return val, true
	}
}

var localDateTimeClassName = "java/time/LocalDateTime"
var localDateClassName = "java/time/LocalDate"
var localTimeClassName = "java/time/LocalTime"

func localDateTimeFromMillis(ms int64) *object.Object {
	t := time.UnixMilli(ms)
	d := object.MakeEmptyObjectWithClassName(&localDateClassName)
	d.FieldTable["year"] = object.Field{Ftype: types.Int, Fvalue: int64(t.Year())}
	d.FieldTable["month"] = object.Field{Ftype: types.Short, Fvalue: int64(t.Month())}
	d.FieldTable["day"] = object.Field{Ftype: types.Short, Fvalue: int64(t.Day())}

	tm := object.MakeEmptyObjectWithClassName(&localTimeClassName)
	tm.FieldTable["hour"] = object.Field{Ftype: types.Byte, Fvalue: int64(t.Hour())}
	tm.FieldTable["minute"] = object.Field{Ftype: types.Byte, Fvalue: int64(t.Minute())}
	tm.FieldTable["second"] = object.Field{Ftype: types.Byte, Fvalue: int64(t.Second())}
	tm.FieldTable["nano"] = object.Field{Ftype: types.Int, Fvalue: int64(t.Nanosecond())}

	ldt := object.MakeEmptyObjectWithClassName(&localDateTimeClassName)
	ldt.FieldTable["date"] = object.Field{Ftype: types.Ref, Fvalue: d}
	ldt.FieldTable["time"] = object.Field{Ftype: types.Ref, Fvalue: tm}
	return ldt
}

func millisFromLocalDateTime(ldt *object.Object) (int64, bool) {
	if ldt == nil || object.IsNull(ldt) {
		return -1, false
	}
	var y, m, d, hour, min, sec, nano int64 = 1970, 1, 1, 0, 0, 0, 0
	if dFld, ok := ldt.FieldTable["date"]; ok {
		if dObj, ok2 := dFld.Fvalue.(*object.Object); ok2 && dObj != nil {
			if yFld, ok3 := dObj.FieldTable["year"]; ok3 {
				y = yFld.Fvalue.(int64)
			}
			if mFld, ok3 := dObj.FieldTable["month"]; ok3 {
				m = mFld.Fvalue.(int64)
			}
			if dayFld, ok3 := dObj.FieldTable["day"]; ok3 {
				d = dayFld.Fvalue.(int64)
			}
		}
	}
	if tFld, ok := ldt.FieldTable["time"]; ok {
		if tObj, ok2 := tFld.Fvalue.(*object.Object); ok2 && tObj != nil {
			if hFld, ok3 := tObj.FieldTable["hour"]; ok3 {
				hour = hFld.Fvalue.(int64)
			}
			if minFld, ok3 := tObj.FieldTable["minute"]; ok3 {
				min = minFld.Fvalue.(int64)
			}
			if sFld, ok3 := tObj.FieldTable["second"]; ok3 {
				sec = sFld.Fvalue.(int64)
			}
			if nFld, ok3 := tObj.FieldTable["nano"]; ok3 {
				nano = nFld.Fvalue.(int64)
			}
		}
	}
	gt := time.Date(int(y), time.Month(m), int(d), int(hour), int(min), int(sec), int(nano), time.Local)
	return gt.UnixMilli(), true
}
