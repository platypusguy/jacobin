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
	"io"
	"math"
	"sync"
	"time"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
)

// Load_Util_Zip_ZipInputStream registers all method signatures for java.util.zip.ZipInputStream.
func Load_Util_Zip_ZipInputStream() {
	// <clinit>()V
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	// <init>(Ljava/io/InputStream;)V
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.<init>(Ljava/io/InputStream;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipInputStreamInit,
		}

	// <init>(Ljava/io/InputStream;Ljava/nio/charset/Charset;)V
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.<init>(Ljava/io/InputStream;Ljava/nio/charset/Charset;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// available()I
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.available()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipInputStreamAvailable,
		}

	// close()V
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipInputStreamClose,
		}

	// closeEntry()V
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.closeEntry()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipInputStreamCloseEntry,
		}

	// createZipEntry(Ljava/lang/String;)Ljava/util/zip/ZipEntry;
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.createZipEntry(Ljava/lang/String;)Ljava/util/zip/ZipEntry;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipInputStreamCreateZipEntry,
		}

	// getNextEntry()Ljava/util/zip/ZipEntry;
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.getNextEntry()Ljava/util/zip/ZipEntry;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipInputStreamGetNextEntry,
		}

	// mark(I)V
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.mark(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	// markSupported()Z
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.markSupported()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnFalse,
		}

	// read()I
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.read()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipInputStreamReadByte,
		}

	// read([B)I
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.read([B)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipInputStreamReadByteArray,
		}

	// read([BII)I
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.read([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  zipInputStreamReadRange,
		}

	// readAllBytes()[B
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.readAllBytes()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipInputStreamReadAllBytes,
		}

	// readNBytes(I)[B
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.readNBytes(I)[B"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipInputStreamReadNBytes,
		}

	// reset()V
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.reset()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zipInputStreamReset,
		}

	// skip(J)J
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.skip(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipInputStreamSkip,
		}

	// transferTo(Ljava/io/OutputStream;)J
	ghelpers.MethodSignatures["java/util/zip/ZipInputStream.transferTo(Ljava/io/OutputStream;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zipInputStreamTransferTo,
		}
}

var zipInputStreamClassName = "java/util/zip/ZipInputStream"

type zipByteReader struct {
	inObj *object.Object
	buf   []byte
	pos   int
	limit int
	eof   bool
}

type zipStreamState struct {
	reader         *zipByteReader
	currentEntry   *object.Object
	entryEOF       bool
	method         int64
	flag           uint16
	expectedCrc    uint32
	expectedSize   int64
	crcCalc        hash.Hash32
	flateReader    io.ReadCloser
	entryBytesRem  int64
	entryBytesRead int64
	closed         bool
}

var (
	zipStreamStatesMu sync.RWMutex
	zipStreamStates   = make(map[*object.Object]*zipStreamState)
)

func asGErrBlk(val any) (*ghelpers.GErrBlk, bool) {
	if val == nil {
		return nil, false
	}
	switch v := val.(type) {
	case *ghelpers.GErrBlk:
		return v, true
	case error:
		return ghelpers.GetGErrBlk(excNames.IOException, v.Error()), true
	}
	return nil, false
}

func newZipByteReader(inObj *object.Object) *zipByteReader {
	return &zipByteReader{
		inObj: inObj,
		buf:   make([]byte, 8192),
	}
}

func (r *zipByteReader) ensure(n int) error {
	for r.limit-r.pos < n && !r.eof {
		if r.pos > 0 {
			remaining := copy(r.buf, r.buf[r.pos:r.limit])
			r.pos = 0
			r.limit = remaining
		}
		if r.limit >= len(r.buf) {
			newBuf := make([]byte, len(r.buf)*2)
			copy(newBuf, r.buf[r.pos:r.limit])
			r.buf = newBuf
		}

		javaBytes := make([]types.JavaByte, 4096)
		bufWrapper := object.MakePrimitiveObject("[B", types.JavaByteArray, javaBytes)

		inClassName := stringPool.GetStringPointer(r.inObj.KlassName)
		if inClassName == nil {
			r.eof = true
			return io.EOF
		}

		method3 := fmt.Sprintf("%s.read([BII)I", *inClassName)
		res := ghelpers.Invoke(method3, []any{r.inObj, bufWrapper, int64(0), int64(4096)})

		if errBlk, ok := asGErrBlk(res); ok {
			if errBlk.ExceptionType == excNames.NoSuchMethodException {
				method1 := fmt.Sprintf("%s.read()I", *inClassName)
				res1 := ghelpers.Invoke(method1, []any{r.inObj})
				if errBlk1, ok1 := asGErrBlk(res1); ok1 {
					return fmt.Errorf("read error: %s", errBlk1.ErrMsg)
				}
				val, ok1 := res1.(int64)
				if !ok1 || val < 0 {
					r.eof = true
					return io.EOF
				}
				r.buf[r.limit] = byte(val)
				r.limit++
				continue
			}
			return fmt.Errorf("read error: %s", errBlk.ErrMsg)
		}

		readN, ok := res.(int64)
		if !ok || readN < 0 {
			r.eof = true
			return io.EOF
		}
		if readN == 0 {
			continue
		}

		actualBytes := bufWrapper.FieldTable["value"].Fvalue.([]types.JavaByte)
		for i := 0; i < int(readN); i++ {
			r.buf[r.limit+i] = byte(actualBytes[i])
		}
		r.limit += int(readN)
	}
	if r.limit-r.pos < n {
		return io.EOF
	}
	return nil
}

func (r *zipByteReader) peek4() (uint32, bool) {
	if err := r.ensure(4); err != nil && r.limit-r.pos < 4 {
		return 0, false
	}
	return binary.LittleEndian.Uint32(r.buf[r.pos : r.pos+4]), true
}

func (r *zipByteReader) ReadByte() (byte, error) {
	if err := r.ensure(1); err != nil {
		return 0, err
	}
	b := r.buf[r.pos]
	r.pos++
	return b, nil
}

func (r *zipByteReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if r.pos >= r.limit {
		if err := r.ensure(1); err != nil {
			return 0, err
		}
	}
	n := copy(p, r.buf[r.pos:r.limit])
	r.pos += n
	return n, nil
}

func (r *zipByteReader) ReadFull(p []byte) error {
	need := len(p)
	read := 0
	for read < need {
		if r.pos >= r.limit {
			if err := r.ensure(1); err != nil {
				return err
			}
		}
		n := copy(p[read:], r.buf[r.pos:r.limit])
		r.pos += n
		read += n
	}
	return nil
}

func getZipStreamState(obj *object.Object) *zipStreamState {
	zipStreamStatesMu.RLock()
	defer zipStreamStatesMu.RUnlock()
	return zipStreamStates[obj]
}

func setZipStreamState(obj *object.Object, state *zipStreamState) {
	zipStreamStatesMu.Lock()
	defer zipStreamStatesMu.Unlock()
	zipStreamStates[obj] = state
}

func zipInputStreamInit(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.<init>: null receiver")
	}
	self := params[0].(*object.Object)
	inObj, ok := params[1].(*object.Object)
	if !ok || inObj == nil || object.IsNull(inObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.<init>: null in")
	}

	var charsetObj = object.Null
	if len(params) >= 3 {
		cs, ok := params[2].(*object.Object)
		if !ok || cs == nil || object.IsNull(cs) {
			return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.<init>: null charset")
		}
		charsetObj = cs
	}

	if self.FieldTable == nil {
		self.FieldTable = make(map[string]object.Field)
	}
	self.FieldTable["in"] = object.Field{Ftype: "Ljava/io/InputStream;", Fvalue: inObj}
	self.FieldTable["closed"] = object.Field{Ftype: "Z", Fvalue: types.JavaBoolFalse}
	self.FieldTable["entryClosed"] = object.Field{Ftype: "Z", Fvalue: types.JavaBoolTrue}
	self.FieldTable["entry"] = object.Field{Ftype: "Ljava/util/zip/ZipEntry;", Fvalue: object.Null}
	if !object.IsNull(charsetObj) {
		self.FieldTable["charset"] = object.Field{Ftype: "Ljava/nio/charset/Charset;", Fvalue: charsetObj}
	}

	state := &zipStreamState{
		reader: newZipByteReader(inObj),
		closed: false,
	}
	setZipStreamState(self, state)

	return nil
}

func zipInputStreamGetNextEntry(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.getNextEntry: null receiver")
	}
	self := params[0].(*object.Object)

	state := getZipStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	// Close open entry if any
	if state.currentEntry != nil && !object.IsNull(state.currentEntry) {
		res := zipInputStreamCloseEntry([]any{self})
		if errBlk, ok := asGErrBlk(res); ok {
			return errBlk
		}
	}

	state.currentEntry = nil
	state.entryEOF = false
	state.flateReader = nil
	state.crcCalc = nil
	state.entryBytesRead = 0

	// Read 4-byte signature
	var sigBuf [4]byte
	err := state.reader.ReadFull(sigBuf[:])
	if err != nil {
		return object.Null
	}

	sig := binary.LittleEndian.Uint32(sigBuf[:])
	if sig != 0x04034b50 {
		return object.Null
	}

	// Read remaining 26 bytes of local file header
	var lfh [26]byte
	if err := state.reader.ReadFull(lfh[:]); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "Truncated ZIP local file header")
	}

	flag := binary.LittleEndian.Uint16(lfh[2:4])
	method := binary.LittleEndian.Uint16(lfh[4:6])
	mtime := binary.LittleEndian.Uint16(lfh[6:8])
	mdate := binary.LittleEndian.Uint16(lfh[8:10])
	crc := binary.LittleEndian.Uint32(lfh[10:14])
	csize := binary.LittleEndian.Uint32(lfh[14:18])
	size := binary.LittleEndian.Uint32(lfh[18:22])
	nameLen := binary.LittleEndian.Uint16(lfh[22:24])
	extraLen := binary.LittleEndian.Uint16(lfh[24:26])

	// Read entry name
	nameBytes := make([]byte, nameLen)
	if nameLen > 0 {
		if err := state.reader.ReadFull(nameBytes); err != nil {
			return ghelpers.GetGErrBlk(excNames.IOException, "Truncated ZIP entry name")
		}
	}
	nameStr := string(nameBytes)

	// Read extra data
	var extraBytes []byte
	if extraLen > 0 {
		extraBytes = make([]byte, extraLen)
		if err := state.reader.ReadFull(extraBytes); err != nil {
			return ghelpers.GetGErrBlk(excNames.IOException, "Truncated ZIP entry extra field")
		}
	}

	// Create ZipEntry
	entryObj, errBlk := CreateZipEntryObject(nameStr)
	if errBlk != nil {
		return errBlk
	}

	// Populate entry
	zipEntrySetMethod([]any{entryObj, int64(method)})
	zipEntrySetTime([]any{entryObj, dosToJavaTime(mdate, mtime)})

	if (flag & 0x08) == 0 {
		zipEntrySetCrc([]any{entryObj, int64(crc)})
		zipEntrySetCompressedSize([]any{entryObj, int64(csize)})
		zipEntrySetSize([]any{entryObj, int64(size)})
	} else {
		zipEntrySetCrc([]any{entryObj, int64(-1)})
		zipEntrySetCompressedSize([]any{entryObj, int64(-1)})
		zipEntrySetSize([]any{entryObj, int64(-1)})
	}

	if len(extraBytes) > 0 {
		jBytes := make([]types.JavaByte, len(extraBytes))
		for i, b := range extraBytes {
			jBytes[i] = types.JavaByte(b)
		}
		extraObj := object.MakePrimitiveObject("[B", types.JavaByteArray, jBytes)
		zipEntrySetExtra([]any{entryObj, extraObj})
	}

	state.currentEntry = entryObj
	state.method = int64(method)
	state.flag = flag
	state.expectedCrc = crc
	state.expectedSize = int64(size)
	state.crcCalc = crc32.NewIEEE()
	state.entryBytesRead = 0
	state.entryEOF = false

	if method == uint16(ZipEntryDeflated) {
		state.flateReader = flate.NewReader(state.reader)
	} else if method == uint16(ZipEntryStored) {
		if flag&0x08 != 0 && size == 0 {
			return ghelpers.GetGErrBlk(excNames.IOException, "only DEFLATED entries can have EXT descriptor")
		}
		state.entryBytesRem = int64(size)
	}

	self.FieldTable["entry"] = object.Field{Ftype: "Ljava/util/zip/ZipEntry;", Fvalue: entryObj}
	self.FieldTable["entryClosed"] = object.Field{Ftype: "Z", Fvalue: types.JavaBoolFalse}

	return entryObj
}

func readDataDescriptor(state *zipStreamState) error {
	var ddBuf [4]byte
	if err := state.reader.ReadFull(ddBuf[:]); err != nil {
		return err
	}
	val := binary.LittleEndian.Uint32(ddBuf[:])
	var crc uint32
	var csize uint32
	var size uint32

	if val == 0x08074b50 { // EXTSIG
		var rest [12]byte
		if err := state.reader.ReadFull(rest[:]); err != nil {
			return err
		}
		crc = binary.LittleEndian.Uint32(rest[0:4])
		csize = binary.LittleEndian.Uint32(rest[4:8])
		size = binary.LittleEndian.Uint32(rest[8:12])
	} else {
		crc = val
		var rest [8]byte
		if err := state.reader.ReadFull(rest[:]); err != nil {
			return err
		}
		csize = binary.LittleEndian.Uint32(rest[0:4])
		size = binary.LittleEndian.Uint32(rest[4:8])
	}

	state.expectedCrc = crc
	state.expectedSize = int64(size)
	if state.currentEntry != nil && !object.IsNull(state.currentEntry) {
		zipEntrySetCrc([]any{state.currentEntry, int64(crc)})
		zipEntrySetCompressedSize([]any{state.currentEntry, int64(csize)})
		zipEntrySetSize([]any{state.currentEntry, int64(size)})
	}
	return nil
}

func zipInputStreamReadRange(params []any) any {
	if len(params) < 4 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.read: null receiver")
	}
	self := params[0].(*object.Object)

	state := getZipStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	bufObj, ok := params[1].(*object.Object)
	if !ok || bufObj == nil || object.IsNull(bufObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.read: null buffer")
	}

	javaBytes, ok := bufObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.read: invalid byte array")
	}

	off, ok := ghelpers.GetInt64Val(params[2])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZipInputStream.read: invalid offset")
	}
	lenVal, ok := ghelpers.GetInt64Val(params[3])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZipInputStream.read: invalid length")
	}

	if off < 0 || lenVal < 0 || off+lenVal > int64(len(javaBytes)) {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "ZipInputStream.read: index out of bounds")
	}

	if lenVal == 0 {
		return int64(0)
	}

	if state.currentEntry == nil || object.IsNull(state.currentEntry) || state.entryEOF {
		return int64(-1)
	}

	if state.method == ZipEntryStored {
		if state.entryBytesRem == 0 {
			state.entryEOF = true
			if (state.flag & 0x08) != 0 {
				if err := readDataDescriptor(state); err != nil {
					return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
				}
			}
			if state.crcCalc != nil {
				computed := state.crcCalc.Sum32()
				if computed != state.expectedCrc {
					return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("corrupt zip file: CRC-32 mismatch (expected %08x, got %08x)", state.expectedCrc, computed))
				}
			}
			return int64(-1)
		}

		if state.entryBytesRem > 0 {
			toRead := int(lenVal)
			if int64(toRead) > state.entryBytesRem {
				toRead = int(state.entryBytesRem)
			}

			tmp := make([]byte, toRead)
			n, err := state.reader.Read(tmp)
			if n > 0 {
				for i := range n {
					javaBytes[off+int64(i)] = types.JavaByte(tmp[i])
				}
				_, err = state.crcCalc.Write(tmp[:n])
				if err != nil {
					return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
				}
				state.entryBytesRead += int64(n)
				state.entryBytesRem -= int64(n)

				if state.entryBytesRem == 0 {
					state.entryEOF = true
					if (state.flag & 0x08) != 0 {
						if err := readDataDescriptor(state); err != nil {
							return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
						}
					}
					if state.crcCalc != nil {
						computed := state.crcCalc.Sum32()
						if computed != state.expectedCrc {
							return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("corrupt zip file: CRC-32 mismatch (expected %08x, got %08x)", state.expectedCrc, computed))
						}
					}
				}
				return int64(n)
			}
			if err != nil {
				state.entryEOF = true
				return int64(-1)
			}
			return int64(0)
		}

		// entryBytesRem == -1 (streaming STORED with data descriptor)
		if sig, ok := state.reader.peek4(); ok && (sig == 0x08074b50 || sig == 0x02014b50 || sig == 0x06054b50 || sig == 0x04034b50) {
			state.entryEOF = true
			if sig == 0x08074b50 {
				if err := readDataDescriptor(state); err != nil {
					return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
				}
			}
			if state.crcCalc != nil {
				computed := state.crcCalc.Sum32()
				if computed != state.expectedCrc {
					return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("corrupt zip file: CRC-32 mismatch (expected %08x, got %08x)", state.expectedCrc, computed))
				}
			}
			return int64(-1)
		}

		b, err := state.reader.ReadByte()
		if err != nil {
			state.entryEOF = true
			return int64(-1)
		}
		javaBytes[off] = types.JavaByte(b)
		_, err = state.crcCalc.Write([]byte{b})
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
		}
		state.entryBytesRead++

		readCount := 1
		for readCount < int(lenVal) {
			if sig, ok := state.reader.peek4(); ok && (sig == 0x08074b50 || sig == 0x02014b50 || sig == 0x06054b50 || sig == 0x04034b50) {
				break
			}
			b2, err2 := state.reader.ReadByte()
			if err2 != nil {
				break
			}
			javaBytes[off+int64(readCount)] = types.JavaByte(b2)
			_, err = state.crcCalc.Write([]byte{b2})
			if err != nil {
				return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
			}
			state.entryBytesRead++
			readCount++
		}

		return int64(readCount)
	} else if state.method == ZipEntryDeflated {
		if state.flateReader == nil {
			state.entryEOF = true
			return int64(-1)
		}

		tmp := make([]byte, lenVal)
		n, err := state.flateReader.Read(tmp)
		if n > 0 {
			for i := range n {
				javaBytes[off+int64(i)] = types.JavaByte(tmp[i])
			}
			_, err = state.crcCalc.Write(tmp[:n])
			if err != nil {
				return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
			}
			state.entryBytesRead += int64(n)
		}

		if err == io.EOF {
			state.entryEOF = true
			err = state.flateReader.Close()
			if err != nil {
				return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
			}
			state.flateReader = nil

			if (state.flag & 0x08) != 0 {
				if err := readDataDescriptor(state); err != nil {
					return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
				}
			}

			if state.crcCalc != nil {
				computed := state.crcCalc.Sum32()
				if computed != state.expectedCrc {
					return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("corrupt zip file: CRC-32 mismatch (expected %08x, got %08x)", state.expectedCrc, computed))
				}
			}

			if n == 0 {
				return int64(-1)
			}
			return int64(n)
		} else if err != nil {
			return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Decompression error: %v", err))
		}

		return int64(n)
	}

	return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Unsupported compression method: %d", state.method))

}

func zipInputStreamReadByte(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.read: null receiver")
	}
	self := params[0].(*object.Object)

	javaBytes := make([]types.JavaByte, 1)
	bufObj := object.MakePrimitiveObject("[B", types.JavaByteArray, javaBytes)

	res := zipInputStreamReadRange([]any{self, bufObj, int64(0), int64(1)})
	if errBlk, ok := asGErrBlk(res); ok {
		return errBlk
	}
	n := res.(int64)
	if n <= 0 {
		return int64(-1)
	}
	return int64(byte(javaBytes[0]))
}

func zipInputStreamReadByteArray(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.read: null receiver")
	}
	self := params[0].(*object.Object)

	bufObj, ok := params[1].(*object.Object)
	if !ok || bufObj == nil || object.IsNull(bufObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.read: null buffer")
	}

	javaBytes, ok := bufObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.read: invalid byte array")
	}

	return zipInputStreamReadRange([]any{self, bufObj, int64(0), int64(len(javaBytes))})
}

func zipInputStreamReadAllBytes(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.readAllBytes: null receiver")
	}
	return zipInputStreamReadNBytes([]any{params[0], int64(math.MaxInt32)})
}

func zipInputStreamReadNBytes(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.readNBytes: null receiver")
	}
	self := params[0].(*object.Object)

	state := getZipStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	lenVal, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZipInputStream.readNBytes: invalid len")
	}
	if lenVal < 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "len < 0")
	}
	if lenVal == 0 {
		return object.MakePrimitiveObject("[B", types.JavaByteArray, []types.JavaByte{})
	}
	if lenVal > 2147483639 && lenVal != math.MaxInt32 {
		return ghelpers.GetGErrBlk(excNames.OutOfMemoryError, "Required array size too large")
	}

	if state.currentEntry == nil || object.IsNull(state.currentEntry) || state.entryEOF {
		return object.MakePrimitiveObject("[B", types.JavaByteArray, []types.JavaByte{})
	}

	var totalRead int64 = 0
	capacity := min(lenVal, int64(8192))
	result := make([]types.JavaByte, 0, capacity)

	for totalRead < lenVal {
		toRead := min(int64(8192), lenVal-totalRead)
		chunkBuf := make([]types.JavaByte, toRead)
		chunkObj := object.MakePrimitiveObject("[B", types.JavaByteArray, chunkBuf)

		res := zipInputStreamReadRange([]any{self, chunkObj, int64(0), toRead})
		if errBlk, ok := asGErrBlk(res); ok {
			return errBlk
		}
		n, ok := res.(int64)
		if !ok || n < 0 {
			break
		}
		if n == 0 {
			continue
		}

		result = append(result, chunkBuf[:n]...)
		totalRead += n
	}

	if result == nil {
		result = []types.JavaByte{}
	}

	return object.MakePrimitiveObject("[B", types.JavaByteArray, result)
}

func zipInputStreamCloseEntry(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.closeEntry: null receiver")
	}
	self := params[0].(*object.Object)

	state := getZipStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	if state.currentEntry == nil || object.IsNull(state.currentEntry) {
		return nil
	}

	if !state.entryEOF {
		drainBytes := make([]types.JavaByte, 2048)
		drainObj := object.MakePrimitiveObject("[B", types.JavaByteArray, drainBytes)
		for !state.entryEOF {
			res := zipInputStreamReadRange([]any{self, drainObj, int64(0), int64(len(drainBytes))})
			if errBlk, ok := asGErrBlk(res); ok {
				return errBlk
			}
			n := res.(int64)
			if n < 0 {
				break
			}
		}
	}

	state.currentEntry = nil
	state.entryEOF = true
	self.FieldTable["entry"] = object.Field{Ftype: "Ljava/util/zip/ZipEntry;", Fvalue: object.Null}
	self.FieldTable["entryClosed"] = object.Field{Ftype: "Z", Fvalue: types.JavaBoolTrue}

	return nil
}

func zipInputStreamSkip(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.skip: null receiver")
	}
	self := params[0].(*object.Object)

	state := getZipStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	n, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZipInputStream.skip: invalid skip count")
	}
	if n < 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "negative skip length")
	}
	if n == 0 || state.currentEntry == nil || object.IsNull(state.currentEntry) || state.entryEOF {
		return int64(0)
	}

	maxBufSize := int64(2048)
	bufSize := min(n, maxBufSize)

	javaBytes := make([]types.JavaByte, bufSize)
	bufObj := object.MakePrimitiveObject("[B", types.JavaByteArray, javaBytes)

	totalSkipped := int64(0)
	for totalSkipped < n {
		toRead := min(n-totalSkipped, bufSize)

		res := zipInputStreamReadRange([]any{self, bufObj, int64(0), toRead})
		if errBlk, ok := asGErrBlk(res); ok {
			return errBlk
		}
		nr := res.(int64)
		if nr < 0 {
			break
		}
		totalSkipped += nr
	}

	return totalSkipped
}

func zipInputStreamTransferTo(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.transferTo: null receiver")
	}
	self := params[0].(*object.Object)

	if params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.transferTo: out is null")
	}
	out, ok := params[1].(*object.Object)
	if !ok || out == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.transferTo: out is null")
	}

	state := getZipStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	if state.currentEntry == nil || object.IsNull(state.currentEntry) || state.entryEOF {
		return int64(0)
	}

	outClassNamePtr := stringPool.GetStringPointer(out.KlassName)
	if outClassNamePtr == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.transferTo: invalid out class")
	}
	outClassName := *outClassNamePtr

	var transferred int64 = 0
	buf := make([]types.JavaByte, 8192)
	bufObj := object.MakePrimitiveObject("[B", types.JavaByteArray, buf)

	for {
		res := zipInputStreamReadRange([]any{self, bufObj, int64(0), int64(len(buf))})
		if errBlk, ok := asGErrBlk(res); ok {
			return errBlk
		}
		nRead, ok := res.(int64)
		if !ok || nRead < 0 {
			break
		}
		if nRead == 0 {
			continue
		}

		resWrite := outputStreamWriteRange(outClassName, out, bufObj, int64(0), nRead)
		if errBlk, ok := asGErrBlk(resWrite); ok {
			return errBlk
		}
		transferred += nRead
	}

	return transferred
}

func outputStreamWriteRange(className string, out *object.Object, bufObj *object.Object, off int64, lenVal int64) any {
	if lenVal == 0 {
		return nil
	}

	if className != "" {
		method3 := fmt.Sprintf("%s.write([BII)V", className)
		if gmeth, ok := ghelpers.MethodSignatures[method3]; ok && gmeth.GFunction != nil {
			var res any
			if gmeth.NeedsContext {
				res = gmeth.GFunction([]any{list.New(), out, bufObj, off, lenVal})
			} else {
				res = ghelpers.Invoke(method3, []any{out, bufObj, off, lenVal})
			}
			if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.NoSuchMethodException {
				return res
			}
		}

		method1 := fmt.Sprintf("%s.write(I)V", className)
		if gmeth, ok := ghelpers.MethodSignatures[method1]; ok && gmeth.GFunction != nil {
			buf := bufObj.FieldTable["value"].Fvalue.([]types.JavaByte)
			for i := int64(0); i < lenVal; i++ {
				b := int64(uint8(buf[off+i]))
				res := ghelpers.Invoke(method1, []any{out, b})
				if res != nil {
					return res
				}
			}
			return nil
		}
	}

	if gmeth, ok := ghelpers.MethodSignatures["java/io/OutputStream.write([BII)V"]; ok && gmeth.GFunction != nil {
		return gmeth.GFunction([]any{list.New(), out, bufObj, off, lenVal})
	}

	return ghelpers.GetGErrBlk(excNames.IOException, "transferTo: No suitable write method found for OutputStream")
}

func zipInputStreamAvailable(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.available: null receiver")
	}
	self := params[0].(*object.Object)

	state := getZipStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	if state.entryEOF || state.currentEntry == nil || object.IsNull(state.currentEntry) {
		return int64(0)
	}

	return int64(1)
}

func zipInputStreamClose(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.close: null receiver")
	}
	self := params[0].(*object.Object)

	state := getZipStreamState(self)
	if state != nil && !state.closed {
		state.closed = true
		if state.flateReader != nil {
			err := state.flateReader.Close()
			if err != nil {
				return ghelpers.GetGErrBlk(excNames.IOException, err.Error())
			}
			state.flateReader = nil
		}
		if state.reader != nil && state.reader.inObj != nil && !object.IsNull(state.reader.inObj) {
			inClassName := stringPool.GetStringPointer(state.reader.inObj.KlassName)
			if inClassName != nil {
				closeMethod := fmt.Sprintf("%s.close()V", *inClassName)
				ghelpers.Invoke(closeMethod, []any{state.reader.inObj})
			}
		}
		self.FieldTable["closed"] = object.Field{Ftype: "Z", Fvalue: types.JavaBoolTrue}
	}
	return nil
}

func zipInputStreamReset(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.reset: null receiver")
	}
	self := params[0].(*object.Object)

	state := getZipStreamState(self)
	if state == nil || state.closed {
		return ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}

	return ghelpers.GetGErrBlk(excNames.IOException, "mark/reset not supported")
}

func zipInputStreamCreateZipEntry(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.createZipEntry: null receiver")
	}
	nameObj, ok := params[1].(*object.Object)
	if !ok || nameObj == nil || object.IsNull(nameObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZipInputStream.createZipEntry: null name")
	}
	nameStr := object.GoStringFromStringObject(nameObj)
	entryObj, errBlk := CreateZipEntryObject(nameStr)
	if errBlk != nil {
		return errBlk
	}
	return entryObj
}

func dosToJavaTime(dosDate, dosTime uint16) int64 {
	year := int(((dosDate >> 9) & 0x7f) + 1980)
	month := time.Month((dosDate >> 5) & 0x0f)
	if month < 1 || month > 12 {
		month = 1
	}
	day := int(dosDate & 0x1f)
	if day < 1 || day > 31 {
		day = 1
	}
	hour := int((dosTime >> 11) & 0x1f)
	minute := int((dosTime >> 5) & 0x3f)
	second := int((dosTime & 0x1f) * 2)

	t := time.Date(year, month, day, hour, minute, second, 0, time.Local)
	return t.UnixMilli()
}
