/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"compress/flate"
	"container/list"
	"encoding/binary"
	"fmt"
	"hash"
	"hash/crc32"
	"jacobin/src/util"
	"sync"
	"time"
	"unsafe"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/stringPool"
	"jacobin/src/types"
)

var zipOutputStreamClassName = "java/util/zip/ZipOutputStream"

// ZIP header signatures
const (
	zipLocSig uint32 = 0x04034b50
	zipExtSig uint32 = 0x08074b50
	zipCenSig uint32 = 0x02014b50
	zipEndSig uint32 = 0x06054b50
)

type zipEntryMeta struct {
	name      string
	extra     []byte
	comment   string
	method    uint16
	flag      uint16
	dosTime   uint16
	dosDate   uint16
	crc       uint32
	csize     uint64
	size      uint64
	locOffset uint64
}

type zipOutStreamWriter struct {
	out     *object.Object
	written int64
	err     error
}

func (w *zipOutStreamWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	res := writeToOutputStream(w.out, p)
	if errBlk, ok := asGErrBlk(res); ok {
		w.err = fmt.Errorf("%s", errBlk.ErrMsg)
		return 0, w.err
	}
	w.written += int64(len(p))
	return len(p), nil
}

type zipCountingWriter struct {
	parent *zipOutStreamWriter
	count  int64
}

func (w *zipCountingWriter) Write(p []byte) (int, error) {
	n, err := w.parent.Write(p)
	w.count += int64(n)
	return n, err
}

type zipOutStreamState struct {
	out                           *object.Object
	writer                        *zipOutStreamWriter
	countingWriter                zipCountingWriter
	comment                       string
	defaultMethod                 int64
	defaultLevel                  int
	currentEntry                  *object.Object
	currentEntryDeflater          *flate.Writer
	currentEntryCrc               hash.Hash32
	currentEntryWritten           int64
	currentEntryCompressedWritten int64
	currentEntryMethod            int64
	currentEntryFlag              uint16
	currentEntryLocOffset         int64
	currentEntryDosTime           uint16
	currentEntryDosDate           uint16
	entries                       []zipEntryMeta
	names                         map[string]struct{}
	closed                        bool
	finished                      bool
}

var (
	zipOutStreamStatesMu sync.RWMutex
	zipOutStreamStates   = make(map[*object.Object]*zipOutStreamState)
)

func getZipOutStreamState(obj *object.Object) *zipOutStreamState {
	zipOutStreamStatesMu.RLock()
	defer zipOutStreamStatesMu.RUnlock()
	return zipOutStreamStates[obj]
}

func setZipOutStreamState(obj *object.Object, state *zipOutStreamState) {
	zipOutStreamStatesMu.Lock()
	defer zipOutStreamStatesMu.Unlock()
	zipOutStreamStates[obj] = state
}

// Load_Util_Zip_ZipOutputStream registers all method signatures for java.util.zip.ZipOutputStream.
func Load_Util_Zip_ZipOutputStream() {
	// <clinit>()V
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipOutputStreamClinit,
		}

	// <init>(Ljava/io/OutputStream;)V
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.<init>(Ljava/io/OutputStream;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipOutputStreamInit,
		}

	// <init>(Ljava/io/OutputStream;Ljava/nio/charset/Charset;)V
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.<init>(Ljava/io/OutputStream;Ljava/nio/charset/Charset;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// setComment(Ljava/lang/String;)V
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.setComment(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipOutputStreamSetComment,
		}

	// setMethod(I)V
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.setMethod(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipOutputStreamSetMethod,
		}

	// setLevel(I)V
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.setLevel(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipOutputStreamSetLevel,
		}

	// putNextEntry(Ljava/util/zip/ZipEntry;)V
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.putNextEntry(Ljava/util/zip/ZipEntry;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipOutputStreamPutNextEntry,
		}

	// closeEntry()V
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.closeEntry()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipOutputStreamCloseEntry,
		}

	// write([BII)V
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.write([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  zipOutputStreamWriteRange,
		}

	// finish()V
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.finish()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipOutputStreamFinish,
		}

	// close()V
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipOutputStreamClose,
		}

	// Override member functions
	// deflate()V
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.deflate()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipOutputStreamDeflate,
		}

	// flush()V
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.flush()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipOutputStreamFlush,
		}

	// write(I)V
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.write(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipOutputStreamWriteByte,
		}

	// write([B)V
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.write([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipOutputStreamWriteByteArray,
		}

	// toString()Ljava/lang/String;
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipOutputStreamToString,
		}

	// Standard Object lifecycle & concurrency methods
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.getClass()Ljava/lang/Class;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.finalize()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.notify()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.notifyAll()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.wait()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.wait(J)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.wait(JI)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: zipOutputStreamEquals}
	ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.hashCode()I"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: zipOutputStreamHashCode}
}

func zipOutputStreamClinit(params []any) any {
	_ = statics.AddStatic(zipOutputStreamClassName+".STORED", statics.Static{
		Type:  types.Int,
		Value: ZipEntryStored,
	})
	_ = statics.AddStatic(zipOutputStreamClassName+".DEFLATED", statics.Static{
		Type:  types.Int,
		Value: ZipEntryDeflated,
	})
	return nil
}

func zipOutputStreamInit(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipOutputStream.<init>: null receiver")
	}
	self := params[0].(*object.Object)
	outObj, ok := params[1].(*object.Object)
	if !ok || outObj == nil || object.IsNull(outObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipOutputStream.<init>: out is null")
	}

	if self.FieldTable == nil {
		self.FieldTable = make(map[string]object.Field)
	}
	self.FieldTable["out"] = object.Field{Ftype: "Ljava/io/OutputStream;", Fvalue: outObj}
	self.FieldTable["closed"] = object.Field{Ftype: "Z", Fvalue: types.JavaBoolFalse}
	self.FieldTable["entry"] = object.Field{Ftype: "Ljava/util/zip/ZipEntry;", Fvalue: object.Null}
	self.FieldTable["comment"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.Null}
	self.FieldTable["method"] = object.Field{Ftype: types.Int, Fvalue: ZipEntryDeflated}
	self.FieldTable["level"] = object.Field{Ftype: types.Int, Fvalue: int64(flate.DefaultCompression)}

	writer := &zipOutStreamWriter{out: outObj}
	state := &zipOutStreamState{
		out:            outObj,
		writer:         writer,
		countingWriter: zipCountingWriter{parent: writer},
		defaultMethod:  ZipEntryDeflated,
		defaultLevel:   flate.DefaultCompression,
		names:          make(map[string]struct{}),
		closed:         false,
		finished:       false,
	}
	setZipOutStreamState(self, state)
	return nil
}

func zipOutputStreamSetComment(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipOutputStream.setComment: null receiver")
	}
	self := params[0].(*object.Object)
	state := getZipOutStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	if len(params) < 2 || params[1] == nil || object.IsNull(params[1]) {
		state.comment = ""
		self.FieldTable["comment"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.Null}
		return nil
	}

	strObj, ok := params[1].(*object.Object)
	if !ok || strObj == nil || object.IsNull(strObj) {
		state.comment = ""
		self.FieldTable["comment"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.Null}
		return nil
	}

	commentStr := object.GoStringFromStringObject(strObj)
	if len([]byte(commentStr)) > 0xFFFF {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZIP file comment too long")
	}

	state.comment = commentStr
	self.FieldTable["comment"] = object.Field{Ftype: types.StringClassRef, Fvalue: strObj}
	return nil
}

func zipOutputStreamSetMethod(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipOutputStream.setMethod: null receiver")
	}
	self := params[0].(*object.Object)
	state := getZipOutStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid compression method")
	}
	methodVal, ok := ghelpers.GetInt64Val(params[1])
	if !ok || (methodVal != ZipEntryStored && methodVal != ZipEntryDeflated) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid compression method")
	}

	state.defaultMethod = methodVal
	self.FieldTable["method"] = object.Field{Ftype: types.Int, Fvalue: methodVal}
	return nil
}

func zipOutputStreamSetLevel(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipOutputStream.setLevel: null receiver")
	}
	self := params[0].(*object.Object)
	state := getZipOutStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid compression level")
	}
	levelVal, ok := ghelpers.GetInt64Val(params[1])
	if !ok || levelVal < -1 || levelVal > 9 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid compression level")
	}

	state.defaultLevel = int(levelVal)
	self.FieldTable["level"] = object.Field{Ftype: types.Int, Fvalue: levelVal}
	return nil
}

func zipOutputStreamPutNextEntry(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipOutputStream.putNextEntry: null receiver")
	}
	self := params[0].(*object.Object)
	state := getZipOutStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	if len(params) < 2 || params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "entry is null")
	}
	entryObj, ok := params[1].(*object.Object)
	if !ok || entryObj == nil || object.IsNull(entryObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "entry is null")
	}

	// Close open entry if any
	if state.currentEntry != nil {
		res := zipOutputStreamCloseEntry([]any{self})
		if errBlk, ok := asGErrBlk(res); ok {
			return errBlk
		}
	}

	nameObj := zipEntryGetName([]any{entryObj})
	nameStr := ""
	if nameObj != nil && !object.IsNull(nameObj) {
		if sObj, ok := nameObj.(*object.Object); ok {
			nameStr = object.GoStringFromStringObject(sObj)
		}
	}
	if len([]byte(nameStr)) > 0xFFFF {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "entry name is too long")
	}

	if _, exists := state.names[nameStr]; exists {
		return ghelpers.GetGErrBlk(excNames.ZipException, fmt.Sprintf("duplicate entry: %s", nameStr))
	}
	state.names[nameStr] = struct{}{}

	methodVal := zipEntryGetMethod([]any{entryObj})
	method := int64(-1)
	if m, ok := methodVal.(int64); ok {
		method = m
	}
	if method == -1 {
		method = state.defaultMethod
		zipEntrySetMethod([]any{entryObj, method})
	}
	if method != ZipEntryStored && method != ZipEntryDeflated {
		return ghelpers.GetGErrBlk(excNames.ZipException, "unsupported compression method")
	}

	timeVal := zipEntryGetTime([]any{entryObj})
	tMs := int64(-1)
	if t, ok := timeVal.(int64); ok {
		tMs = t
	}
	if tMs == -1 {
		tMs = time.Now().UnixMilli()
		zipEntrySetTime([]any{entryObj, tMs})
	}

	var flag uint16
	var crcVal, csizeVal, sizeVal uint32
	if method == ZipEntryStored {
		size := zipEntryGetSize([]any{entryObj}).(int64)
		csize := zipEntryGetCompressedSize([]any{entryObj}).(int64)
		crc := zipEntryGetCrc([]any{entryObj}).(int64)
		if size == -1 || crc == -1 {
			return ghelpers.GetGErrBlk(excNames.ZipException, "STORED entry missing size, compressed size, or crc-32")
		}
		if csize == -1 {
			csize = size
			zipEntrySetCompressedSize([]any{entryObj, csize})
		} else if csize != size {
			return ghelpers.GetGErrBlk(excNames.ZipException, "STORED entry invalid size")
		}
		flag = 0x0800 // UTF-8 filename flag
		crcVal = uint32(crc)
		csizeVal = uint32(csize)
		sizeVal = uint32(size)
	} else {
		// DEFLATED: bit 3 (data descriptor) + bit 11 (UTF-8)
		flag = 0x08 | 0x0800
	}

	extraBytes := getExtraBytesHelper(entryObj)
	nameBytes := []byte(nameStr)
	dosTime, dosDate := javaToDosTime(tMs)
	locOffset := state.writer.written

	versionNeeded := uint16(10)
	if method == ZipEntryDeflated {
		versionNeeded = 20
	}

	// Write Local File Header (30 bytes + len(nameBytes) + len(extraBytes))
	var lfh [30]byte
	binary.LittleEndian.PutUint32(lfh[0:4], zipLocSig)
	binary.LittleEndian.PutUint16(lfh[4:6], versionNeeded)
	binary.LittleEndian.PutUint16(lfh[6:8], flag)
	binary.LittleEndian.PutUint16(lfh[8:10], uint16(method))
	binary.LittleEndian.PutUint16(lfh[10:12], dosTime)
	binary.LittleEndian.PutUint16(lfh[12:14], dosDate)
	binary.LittleEndian.PutUint32(lfh[14:18], crcVal)
	binary.LittleEndian.PutUint32(lfh[18:22], csizeVal)
	binary.LittleEndian.PutUint32(lfh[22:26], sizeVal)
	binary.LittleEndian.PutUint16(lfh[26:28], uint16(len(nameBytes)))
	binary.LittleEndian.PutUint16(lfh[28:30], uint16(len(extraBytes)))

	if _, err := state.writer.Write(lfh[:]); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
	}
	if len(nameBytes) > 0 {
		if _, err := state.writer.Write(nameBytes); err != nil {
			return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
		}
	}
	if len(extraBytes) > 0 {
		if _, err := state.writer.Write(extraBytes); err != nil {
			return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
		}
	}

	state.currentEntry = entryObj
	state.currentEntryWritten = 0
	state.currentEntryCompressedWritten = 0
	state.currentEntryCrc = crc32.NewIEEE()
	state.currentEntryMethod = method
	state.currentEntryFlag = flag
	state.currentEntryLocOffset = locOffset
	state.currentEntryDosTime = dosTime
	state.currentEntryDosDate = dosDate

	if method == ZipEntryDeflated {
		state.countingWriter.count = 0
		fw, err := flate.NewWriter(&state.countingWriter, state.defaultLevel)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
		}
		state.currentEntryDeflater = fw
	}

	self.FieldTable["entry"] = object.Field{Ftype: "Ljava/util/zip/ZipEntry;", Fvalue: entryObj}
	return nil
}

func zipOutputStreamCloseEntry(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipOutputStream.closeEntry: null receiver")
	}
	self := params[0].(*object.Object)
	state := getZipOutStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	if state.currentEntry == nil {
		return nil
	}

	entryObj := state.currentEntry
	if state.currentEntryMethod == ZipEntryDeflated {
		if state.currentEntryDeflater != nil {
			if err := state.currentEntryDeflater.Close(); err != nil {
				return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
			}
			state.currentEntryCompressedWritten = state.countingWriter.count
			state.currentEntryDeflater = nil
		}

		if (state.currentEntryFlag & 0x08) != 0 {
			// Write Data Descriptor (16 bytes)
			var ext [16]byte
			binary.LittleEndian.PutUint32(ext[0:4], zipExtSig)
			binary.LittleEndian.PutUint32(ext[4:8], state.currentEntryCrc.Sum32())
			binary.LittleEndian.PutUint32(ext[8:12], uint32(state.currentEntryCompressedWritten))
			binary.LittleEndian.PutUint32(ext[12:16], uint32(state.currentEntryWritten))
			if _, err := state.writer.Write(ext[:]); err != nil {
				return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
			}
		}

		// Validate if entry had pre-set values
		preSize := zipEntryGetSize([]any{entryObj}).(int64)
		if preSize != -1 && preSize != state.currentEntryWritten {
			return ghelpers.GetGErrBlk(excNames.ZipException, fmt.Sprintf("invalid entry size (expected %d but got %d bytes)", preSize, state.currentEntryWritten))
		}
		preCsize := zipEntryGetCompressedSize([]any{entryObj}).(int64)
		if preCsize != -1 && preCsize != state.currentEntryCompressedWritten {
			return ghelpers.GetGErrBlk(excNames.ZipException, fmt.Sprintf("invalid entry compressed size (expected %d but got %d bytes)", preCsize, state.currentEntryCompressedWritten))
		}
		preCrc := zipEntryGetCrc([]any{entryObj}).(int64)
		if preCrc != -1 && uint32(preCrc) != state.currentEntryCrc.Sum32() {
			return ghelpers.GetGErrBlk(excNames.ZipException, fmt.Sprintf("invalid entry CRC-32 (expected 0x%x but got 0x%x)", preCrc, state.currentEntryCrc.Sum32()))
		}
	} else if state.currentEntryMethod == ZipEntryStored {
		expectedSize := zipEntryGetSize([]any{entryObj}).(int64)
		if state.currentEntryWritten != expectedSize {
			return ghelpers.GetGErrBlk(excNames.ZipException, fmt.Sprintf("invalid entry size (expected %d but got %d bytes)", expectedSize, state.currentEntryWritten))
		}
		expectedCrc := zipEntryGetCrc([]any{entryObj}).(int64)
		if uint32(expectedCrc) != state.currentEntryCrc.Sum32() {
			return ghelpers.GetGErrBlk(excNames.ZipException, fmt.Sprintf("invalid entry CRC-32 (expected 0x%x but got 0x%x)", expectedCrc, state.currentEntryCrc.Sum32()))
		}
	}

	zipEntrySetSize([]any{entryObj, state.currentEntryWritten})
	zipEntrySetCompressedSize([]any{entryObj, state.currentEntryCompressedWritten})
	zipEntrySetCrc([]any{entryObj, int64(state.currentEntryCrc.Sum32())})

	nameObj := zipEntryGetName([]any{entryObj})
	nameStr := ""
	if nameObj != nil && !object.IsNull(nameObj) {
		if sObj, ok := nameObj.(*object.Object); ok {
			nameStr = object.GoStringFromStringObject(sObj)
		}
	}

	commentObj := zipEntryGetComment([]any{entryObj})
	commentStr := ""
	if commentObj != nil && !object.IsNull(commentObj) {
		if cObj, ok := commentObj.(*object.Object); ok {
			commentStr = object.GoStringFromStringObject(cObj)
		}
	}

	meta := zipEntryMeta{
		name:      nameStr,
		extra:     getExtraBytesHelper(entryObj),
		comment:   commentStr,
		method:    uint16(state.currentEntryMethod),
		flag:      state.currentEntryFlag,
		dosTime:   state.currentEntryDosTime,
		dosDate:   state.currentEntryDosDate,
		crc:       state.currentEntryCrc.Sum32(),
		csize:     uint64(state.currentEntryCompressedWritten),
		size:      uint64(state.currentEntryWritten),
		locOffset: uint64(state.currentEntryLocOffset),
	}
	state.entries = append(state.entries, meta)

	state.currentEntry = nil
	state.currentEntryDeflater = nil
	state.currentEntryCrc = nil
	state.currentEntryWritten = 0
	state.currentEntryCompressedWritten = 0
	self.FieldTable["entry"] = object.Field{Ftype: "Ljava/util/zip/ZipEntry;", Fvalue: object.Null}

	return nil
}

func zipOutputStreamWriteRange(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipOutputStream.write: null receiver")
	}
	self := params[0].(*object.Object)
	state := getZipOutStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	if state.currentEntry == nil {
		return ghelpers.GetGErrBlk(excNames.ZipException, "no current ZIP entry")
	}

	if len(params) < 2 || params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "buffer is null")
	}
	bufObj, ok := params[1].(*object.Object)
	if !ok || bufObj == nil || object.IsNull(bufObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "buffer is null")
	}

	valFld, ok := bufObj.FieldTable["value"]
	if !ok || valFld.Fvalue == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "buffer lacks value field")
	}
	javaBytes, ok := valFld.Fvalue.([]types.JavaByte)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "buffer value is not JavaByte array")
	}

	if len(params) < 4 {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "missing offset or length")
	}
	offVal, ok1 := ghelpers.GetInt64Val(params[2])
	lenVal, ok2 := ghelpers.GetInt64Val(params[3])
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "invalid offset or length")
	}

	if offVal < 0 || lenVal < 0 || offVal > int64(len(javaBytes))-lenVal {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("offset=%d length=%d arrayLength=%d", offVal, lenVal, len(javaBytes)))
	}
	if lenVal == 0 {
		return nil
	}

	rawBytes := object.GoByteArrayFromJavaByteArray(javaBytes[offVal : offVal+lenVal])
	state.currentEntryCrc.Write(rawBytes)

	if state.currentEntryMethod == ZipEntryStored {
		expectedSize := zipEntryGetSize([]any{state.currentEntry}).(int64)
		if state.currentEntryWritten+lenVal > expectedSize {
			return ghelpers.GetGErrBlk(excNames.ZipException, fmt.Sprintf("invalid entry size (expected %d but got %d bytes)", expectedSize, state.currentEntryWritten+lenVal))
		}
		if _, err := state.writer.Write(rawBytes); err != nil {
			return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
		}
		state.currentEntryWritten += lenVal
		state.currentEntryCompressedWritten += lenVal
	} else {
		if _, err := state.currentEntryDeflater.Write(rawBytes); err != nil {
			return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
		}
		state.currentEntryWritten += lenVal
	}

	return nil
}

func zipOutputStreamFinish(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipOutputStream.finish: null receiver")
	}
	self := params[0].(*object.Object)
	state := getZipOutStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	if state.finished {
		return nil
	}

	if state.currentEntry != nil {
		res := zipOutputStreamCloseEntry([]any{self})
		if errBlk, ok := asGErrBlk(res); ok {
			return errBlk
		}
	}

	cenStartOffset := state.writer.written

	// Write Central Directory Headers
	for _, meta := range state.entries {
		nameBytes := []byte(meta.name)
		extraBytes := meta.extra
		commentBytes := []byte(meta.comment)

		versionNeeded := uint16(10)
		if meta.method == uint16(ZipEntryDeflated) {
			versionNeeded = 20
		}

		var cdh [46]byte
		binary.LittleEndian.PutUint32(cdh[0:4], zipCenSig)
		binary.LittleEndian.PutUint16(cdh[4:6], 20) // version made by
		binary.LittleEndian.PutUint16(cdh[6:8], versionNeeded)
		binary.LittleEndian.PutUint16(cdh[8:10], meta.flag)
		binary.LittleEndian.PutUint16(cdh[10:12], meta.method)
		binary.LittleEndian.PutUint16(cdh[12:14], meta.dosTime)
		binary.LittleEndian.PutUint16(cdh[14:16], meta.dosDate)
		binary.LittleEndian.PutUint32(cdh[16:20], meta.crc)
		binary.LittleEndian.PutUint32(cdh[20:24], uint32(meta.csize))
		binary.LittleEndian.PutUint32(cdh[24:28], uint32(meta.size))
		binary.LittleEndian.PutUint16(cdh[28:30], uint16(len(nameBytes)))
		binary.LittleEndian.PutUint16(cdh[30:32], uint16(len(extraBytes)))
		binary.LittleEndian.PutUint16(cdh[32:34], uint16(len(commentBytes)))
		binary.LittleEndian.PutUint16(cdh[34:36], 0) // disk number start
		binary.LittleEndian.PutUint16(cdh[36:38], 0) // internal file attributes
		binary.LittleEndian.PutUint32(cdh[38:42], 0) // external file attributes
		binary.LittleEndian.PutUint32(cdh[42:46], uint32(meta.locOffset))

		if _, err := state.writer.Write(cdh[:]); err != nil {
			return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
		}
		if len(nameBytes) > 0 {
			if _, err := state.writer.Write(nameBytes); err != nil {
				return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
			}
		}
		if len(extraBytes) > 0 {
			if _, err := state.writer.Write(extraBytes); err != nil {
				return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
			}
		}
		if len(commentBytes) > 0 {
			if _, err := state.writer.Write(commentBytes); err != nil {
				return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
			}
		}
	}

	cenEndOffset := state.writer.written
	cenSize := cenEndOffset - cenStartOffset

	// Write End of Central Directory Record (22 bytes + comment)
	commentBytes := []byte(state.comment)
	var eocd [22]byte
	binary.LittleEndian.PutUint32(eocd[0:4], zipEndSig)
	binary.LittleEndian.PutUint16(eocd[4:6], 0)                            // number of this disk
	binary.LittleEndian.PutUint16(eocd[6:8], 0)                            // disk where central directory starts
	binary.LittleEndian.PutUint16(eocd[8:10], uint16(len(state.entries)))  // number of central directory records on this disk
	binary.LittleEndian.PutUint16(eocd[10:12], uint16(len(state.entries))) // total number of central directory records
	binary.LittleEndian.PutUint32(eocd[12:16], uint32(cenSize))
	binary.LittleEndian.PutUint32(eocd[16:20], uint32(cenStartOffset))
	binary.LittleEndian.PutUint16(eocd[20:22], uint16(len(commentBytes)))

	if _, err := state.writer.Write(eocd[:]); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
	}
	if len(commentBytes) > 0 {
		if _, err := state.writer.Write(commentBytes); err != nil {
			return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
		}
	}

	state.finished = true
	return nil
}

func zipOutputStreamClose(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipOutputStream.close: null receiver")
	}
	self := params[0].(*object.Object)
	state := getZipOutStreamState(self)
	if state == nil || state.closed {
		return nil
	}

	if !state.finished {
		res := zipOutputStreamFinish([]any{self})
		if errBlk, ok := asGErrBlk(res); ok {
			return errBlk
		}
	}

	state.closed = true
	self.FieldTable["closed"] = object.Field{Ftype: "Z", Fvalue: types.JavaBoolTrue}

	res := closeOutputStream(state.out)
	if errBlk, ok := asGErrBlk(res); ok {
		return errBlk
	}

	return nil
}

func zipOutputStreamDeflate(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipOutputStream.deflate: null receiver")
	}
	self := params[0].(*object.Object)
	state := getZipOutStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	if state.currentEntry != nil && state.currentEntryMethod == ZipEntryDeflated && state.currentEntryDeflater != nil {
		if err := state.currentEntryDeflater.Flush(); err != nil {
			return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
		}
	}
	return nil
}

func zipOutputStreamFlush(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipOutputStream.flush: null receiver")
	}
	self := params[0].(*object.Object)
	state := getZipOutStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	if state.currentEntry != nil && state.currentEntryMethod == ZipEntryDeflated && state.currentEntryDeflater != nil {
		if err := state.currentEntryDeflater.Flush(); err != nil {
			return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
		}
	}

	res := flushOutputStream(state.out)
	if errBlk, ok := asGErrBlk(res); ok {
		return errBlk
	}
	return nil
}

func zipOutputStreamWriteByte(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipOutputStream.write: null receiver")
	}
	self := params[0].(*object.Object)
	state := getZipOutStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}
	if state.currentEntry == nil {
		return ghelpers.GetGErrBlk(excNames.ZipException, "no current ZIP entry")
	}

	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "missing byte parameter")
	}
	bVal, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid byte parameter")
	}

	jBytes := []types.JavaByte{types.JavaByte(bVal & 0xFF)}
	bufObj := object.MakePrimitiveObject("[B", types.JavaByteArray, jBytes)
	return zipOutputStreamWriteRange([]any{self, bufObj, int64(0), int64(1)})
}

func zipOutputStreamWriteByteArray(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipOutputStream.write: null receiver")
	}
	self := params[0].(*object.Object)
	state := getZipOutStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	if len(params) < 2 || params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "buffer is null")
	}
	bufObj, ok := params[1].(*object.Object)
	if !ok || bufObj == nil || object.IsNull(bufObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "buffer is null")
	}

	valFld, ok := bufObj.FieldTable["value"]
	if !ok || valFld.Fvalue == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "buffer lacks value field")
	}
	javaBytes, ok := valFld.Fvalue.([]types.JavaByte)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "buffer value is not JavaByte array")
	}

	return zipOutputStreamWriteRange([]any{self, bufObj, int64(0), int64(len(javaBytes))})
}

func zipOutputStreamToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return object.StringObjectFromGoString("null")
	}
	self := params[0].(*object.Object)
	userFormatClassName := util.ConvertInternalClassNameToUserFormat(zipOutputStreamClassName)
	return object.StringObjectFromGoString(fmt.Sprintf("%s@%x", userFormatClassName, uintptr(unsafe.Pointer(self))))
}

func zipOutputStreamEquals(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0]) {
		return types.JavaBoolFalse
	}
	if params[0] == params[1] {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func zipOutputStreamHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return int64(0)
	}
	self := params[0].(*object.Object)
	return int64(uintptr(unsafe.Pointer(self)) & 0x7FFFFFFF)
}

func getExtraBytesHelper(entryObj *object.Object) []byte {
	extraObj := zipEntryGetExtra([]any{entryObj})
	if extraObj == nil || object.IsNull(extraObj) {
		return nil
	}
	if o, ok := extraObj.(*object.Object); ok {
		if valFld, exists := o.FieldTable["value"]; exists && valFld.Fvalue != nil {
			if jBytes, ok := valFld.Fvalue.([]types.JavaByte); ok {
				return object.GoByteArrayFromJavaByteArray(jBytes)
			}
		}
	}
	return nil
}

func writeToOutputStream(out *object.Object, p []byte) any {
	if out == nil || object.IsNull(out) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OutputStream is null")
	}

	jBytes := make([]types.JavaByte, len(p))
	for i, b := range p {
		jBytes[i] = types.JavaByte(b)
	}
	bufObj := object.MakePrimitiveObject("[B", types.JavaByteArray, jBytes)

	classNamePtr := stringPool.GetStringPointer(out.KlassName)
	if classNamePtr != nil {
		method3 := fmt.Sprintf("%s.write([BII)V", *classNamePtr)
		if gmeth, ok := ghelpers.MethodSignatures[method3]; ok && gmeth.GFunction != nil {
			res := ghelpers.Invoke(method3, []any{out, bufObj, int64(0), int64(len(p))})
			if res == nil {
				return nil
			}
			if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.NoSuchMethodException {
				return res
			}
		}

		method1 := fmt.Sprintf("%s.write(I)V", *classNamePtr)
		if gmeth, ok := ghelpers.MethodSignatures[method1]; ok && gmeth.GFunction != nil {
			for _, b := range p {
				res := ghelpers.Invoke(method1, []any{out, int64(b)})
				if res != nil {
					return res
				}
			}
			return nil
		}
	}

	if gmeth, ok := ghelpers.MethodSignatures["java/io/FilterOutputStream.write([BII)V"]; ok && gmeth.GFunction != nil {
		res := ghelpers.Invoke("java/io/FilterOutputStream.write([BII)V", []any{out, bufObj, int64(0), int64(len(p))})
		if res == nil {
			return nil
		}
		if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.NoSuchMethodException {
			return res
		}
	}

	if gmeth, ok := ghelpers.MethodSignatures["java/io/OutputStream.write([BII)V"]; ok && gmeth.GFunction != nil {
		return gmeth.GFunction([]any{list.New(), out, bufObj, int64(0), int64(len(p))})
	}

	return ghelpers.GetGErrBlk(excNames.IOException, "write: No suitable write method found for OutputStream")
}

func flushOutputStream(out *object.Object) any {
	if out == nil || object.IsNull(out) {
		return nil
	}
	classNamePtr := stringPool.GetStringPointer(out.KlassName)
	if classNamePtr != nil {
		method := fmt.Sprintf("%s.flush()V", *classNamePtr)
		if gmeth, ok := ghelpers.MethodSignatures[method]; ok && gmeth.GFunction != nil {
			return ghelpers.Invoke(method, []any{out})
		}
	}
	if gmeth, ok := ghelpers.MethodSignatures["java/io/FilterOutputStream.flush()V"]; ok && gmeth.GFunction != nil {
		return ghelpers.Invoke("java/io/FilterOutputStream.flush()V", []any{out})
	}
	return nil
}

func closeOutputStream(out *object.Object) any {
	if out == nil || object.IsNull(out) {
		return nil
	}
	classNamePtr := stringPool.GetStringPointer(out.KlassName)
	if classNamePtr != nil {
		method := fmt.Sprintf("%s.close()V", *classNamePtr)
		if gmeth, ok := ghelpers.MethodSignatures[method]; ok && gmeth.GFunction != nil {
			return ghelpers.Invoke(method, []any{out})
		}
	}
	if gmeth, ok := ghelpers.MethodSignatures["java/io/FilterOutputStream.close()V"]; ok && gmeth.GFunction != nil {
		return ghelpers.Invoke("java/io/FilterOutputStream.close()V", []any{out})
	}
	return nil
}

func javaToDosTime(milli int64) (uint16, uint16) {
	if milli < 0 {
		milli = 0
	}
	t := time.UnixMilli(milli).In(time.Local)
	year := t.Year()
	if year < 1980 {
		return 0x0000, 0x0021 // 1980-01-01 00:00:00
	}
	dosYear := uint16(year - 1980)
	dosMonth := uint16(t.Month())
	dosDay := uint16(t.Day())
	dosHour := uint16(t.Hour())
	dosMinute := uint16(t.Minute())
	dosSecond := uint16(t.Second() / 2)

	dosTime := (dosHour << 11) | (dosMinute << 5) | dosSecond
	dosDate := (dosYear << 9) | (dosMonth << 5) | dosDay
	return dosTime, dosDate
}
