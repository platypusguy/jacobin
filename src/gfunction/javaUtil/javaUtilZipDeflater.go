/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"
	"hash/adler32"
	"sync"
	"unsafe"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
)

var deflaterClassName = "java/util/zip/Deflater"

// Deflater constants
const (
	DeflaterDeflated           int64 = 8
	DeflaterNoCompression      int64 = 0
	DeflaterBestSpeed          int64 = 1
	DeflaterBestCompression    int64 = 9
	DeflaterDefaultCompression int64 = -1
	DeflaterFiltered           int64 = 1
	DeflaterHuffmanOnly        int64 = 2
	DeflaterDefaultStrategy    int64 = 0
	DeflaterNoFlush            int64 = 0
	DeflaterSyncFlush          int64 = 2
	DeflaterFullFlush          int64 = 3
)

type deflaterState struct {
	mu            sync.Mutex
	level         int
	strategy      int
	nowrap        bool
	dict          []byte
	dictSet       bool
	headerWritten bool

	input         []byte
	outputBuf     []byte

	flateWriter   *flate.Writer
	outCollector  *bytes.Buffer

	finish        bool
	closedWriter  bool
	finished      bool
	ended         bool

	bytesRead     int64
	bytesWritten  int64
	adler         uint32
}

var (
	deflaterStatesMu sync.RWMutex
	deflaterStates   = make(map[*object.Object]*deflaterState)
)

func getDeflaterState(obj *object.Object) *deflaterState {
	deflaterStatesMu.RLock()
	defer deflaterStatesMu.RUnlock()
	return deflaterStates[obj]
}

func setDeflaterState(obj *object.Object, state *deflaterState) {
	deflaterStatesMu.Lock()
	defer deflaterStatesMu.Unlock()
	deflaterStates[obj] = state
}

func (s *deflaterState) writeZlibHeader() {
	if s.nowrap || s.headerWritten {
		return
	}
	s.headerWritten = true
	cmf := byte(0x78) // CM = 8 (deflate), CINFO = 7 (32K window)
	var flevel byte
	switch {
	case s.level == 1:
		flevel = 0 // fastest
	case s.level >= 2 && s.level <= 5:
		flevel = 1 // fast
	case s.level >= 7:
		flevel = 3 // max compression
	default:
		flevel = 2 // default
	}
	fdict := byte(0)
	if s.dictSet {
		fdict = 1
	}
	baseFLG := (flevel << 6) | (fdict << 5)
	fcheck := (31 - ((uint16(cmf)*256 + uint16(baseFLG)) % 31)) % 31
	flg := baseFLG | byte(fcheck)

	s.outCollector.WriteByte(cmf)
	s.outCollector.WriteByte(flg)

	if s.dictSet {
		var dictID [4]byte
		binary.BigEndian.PutUint32(dictID[:], adler32.Checksum(s.dict))
		s.outCollector.Write(dictID[:])
	}
}

func (s *deflaterState) ensureWriter() error {
	if s.flateWriter != nil {
		return nil
	}
	s.writeZlibHeader()
	flateLevel := s.level
	if s.strategy == int(DeflaterHuffmanOnly) {
		flateLevel = flate.HuffmanOnly
	}
	var fw *flate.Writer
	var err error
	if s.dictSet && len(s.dict) > 0 {
		fw, err = flate.NewWriterDict(s.outCollector, flateLevel, s.dict)
	} else {
		fw, err = flate.NewWriter(s.outCollector, flateLevel)
	}
	if err != nil {
		return err
	}
	s.flateWriter = fw
	return nil
}

func (s *deflaterState) deflate(dest []byte, flush int64) (int, error) {
	if s.ended {
		return 0, fmt.Errorf("Deflater has been closed")
	}
	if len(dest) == 0 {
		return 0, nil
	}
	if s.finished {
		return 0, nil
	}

	if err := s.ensureWriter(); err != nil {
		return 0, err
	}

	// Feed any uncompressed input into flate writer
	if len(s.input) > 0 {
		n, err := s.flateWriter.Write(s.input)
		if n > 0 {
			s.adler = updateAdler32(s.adler, s.input[:n])
			s.bytesRead += int64(n)
			s.input = s.input[n:]
		}
		if err != nil {
			return 0, err
		}
	}

	// Handle finish or flush
	if s.finish {
		if !s.closedWriter && len(s.input) == 0 {
			if err := s.flateWriter.Close(); err != nil {
				return 0, err
			}
			s.closedWriter = true
			if !s.nowrap {
				var trailer [4]byte
				binary.BigEndian.PutUint32(trailer[:], s.adler)
				s.outCollector.Write(trailer[:])
			}
		}
	} else if flush == DeflaterSyncFlush || flush == DeflaterFullFlush {
		if err := s.flateWriter.Flush(); err != nil {
			return 0, err
		}
	}

	// Collect any compressed bytes generated into s.outputBuf
	if s.outCollector.Len() > 0 {
		s.outputBuf = append(s.outputBuf, s.outCollector.Bytes()...)
		s.outCollector.Reset()
	}

	toCopy := len(dest)
	if len(s.outputBuf) < toCopy {
		toCopy = len(s.outputBuf)
	}
	if toCopy > 0 {
		copy(dest[:toCopy], s.outputBuf[:toCopy])
		s.outputBuf = s.outputBuf[toCopy:]
		s.bytesWritten += int64(toCopy)
	}

	if s.finish && s.closedWriter && len(s.outputBuf) == 0 && len(s.input) == 0 {
		s.finished = true
	}

	return toCopy, nil
}

func (s *deflaterState) reset() {
	s.input = nil
	s.outputBuf = nil
	if s.outCollector != nil {
		s.outCollector.Reset()
	}
	s.flateWriter = nil
	s.closedWriter = false
	s.finish = false
	s.finished = false
	s.headerWritten = false
	s.dict = nil
	s.dictSet = false
	s.bytesRead = 0
	s.bytesWritten = 0
	s.adler = 1
}

func (s *deflaterState) end() {
	s.ended = true
	s.input = nil
	s.outputBuf = nil
	s.outCollector = nil
	s.flateWriter = nil
	s.dict = nil
}

// Load_Util_Zip_Deflater registers all method signatures for java.util.zip.Deflater.
func Load_Util_Zip_Deflater() {
	// <clinit>()V
	ghelpers.MethodSignatures["java/util/zip/Deflater.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  deflaterClinit,
		}

	// <init>()V
	ghelpers.MethodSignatures["java/util/zip/Deflater.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  deflaterInitDefault,
		}

	// <init>(I)V
	ghelpers.MethodSignatures["java/util/zip/Deflater.<init>(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  deflaterInitLevel,
		}

	// <init>(IZ)V
	ghelpers.MethodSignatures["java/util/zip/Deflater.<init>(IZ)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  deflaterInitLevelNowrap,
		}

	// setInput([BII)V
	ghelpers.MethodSignatures["java/util/zip/Deflater.setInput([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  deflaterSetInputRange,
		}

	// setInput([B)V
	ghelpers.MethodSignatures["java/util/zip/Deflater.setInput([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  deflaterSetInputArray,
		}

	// setInput(Ljava/nio/ByteBuffer;)V - Trapped
	ghelpers.MethodSignatures["java/util/zip/Deflater.setInput(Ljava/nio/ByteBuffer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	// setDictionary([BII)V
	ghelpers.MethodSignatures["java/util/zip/Deflater.setDictionary([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  deflaterSetDictionaryRange,
		}

	// setDictionary([B)V
	ghelpers.MethodSignatures["java/util/zip/Deflater.setDictionary([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  deflaterSetDictionaryArray,
		}

	// setDictionary(Ljava/nio/ByteBuffer;)V - Trapped
	ghelpers.MethodSignatures["java/util/zip/Deflater.setDictionary(Ljava/nio/ByteBuffer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	// setStrategy(I)V
	ghelpers.MethodSignatures["java/util/zip/Deflater.setStrategy(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  deflaterSetStrategy,
		}

	// setLevel(I)V
	ghelpers.MethodSignatures["java/util/zip/Deflater.setLevel(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  deflaterSetLevel,
		}

	// needsInput()Z
	ghelpers.MethodSignatures["java/util/zip/Deflater.needsInput()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  deflaterNeedsInput,
		}

	// finish()V
	ghelpers.MethodSignatures["java/util/zip/Deflater.finish()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  deflaterFinish,
		}

	// finished()Z
	ghelpers.MethodSignatures["java/util/zip/Deflater.finished()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  deflaterFinished,
		}

	// deflate([BII)I
	ghelpers.MethodSignatures["java/util/zip/Deflater.deflate([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  deflaterDeflateRange,
		}

	// deflate([B)I
	ghelpers.MethodSignatures["java/util/zip/Deflater.deflate([B)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  deflaterDeflateArray,
		}

	// deflate([BIII)I
	ghelpers.MethodSignatures["java/util/zip/Deflater.deflate([BIII)I"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  deflaterDeflateRangeFlush,
		}

	// deflate(Ljava/nio/ByteBuffer;)I - Trapped
	ghelpers.MethodSignatures["java/util/zip/Deflater.deflate(Ljava/nio/ByteBuffer;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	// deflate(Ljava/nio/ByteBuffer;I)I - Trapped
	ghelpers.MethodSignatures["java/util/zip/Deflater.deflate(Ljava/nio/ByteBuffer;I)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// getAdler()I
	ghelpers.MethodSignatures["java/util/zip/Deflater.getAdler()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  deflaterGetAdler,
		}

	// getTotalIn()I
	ghelpers.MethodSignatures["java/util/zip/Deflater.getTotalIn()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  deflaterGetTotalIn,
		}

	// getBytesRead()J
	ghelpers.MethodSignatures["java/util/zip/Deflater.getBytesRead()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  deflaterGetBytesRead,
		}

	// getTotalOut()I
	ghelpers.MethodSignatures["java/util/zip/Deflater.getTotalOut()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  deflaterGetTotalOut,
		}

	// getBytesWritten()J
	ghelpers.MethodSignatures["java/util/zip/Deflater.getBytesWritten()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  deflaterGetBytesWritten,
		}

	// reset()V
	ghelpers.MethodSignatures["java/util/zip/Deflater.reset()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  deflaterReset,
		}

	// end()V
	ghelpers.MethodSignatures["java/util/zip/Deflater.end()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  deflaterEnd,
		}

	// finalize()V
	ghelpers.MethodSignatures["java/util/zip/Deflater.finalize()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  deflaterFinalize,
		}

	// toString()Ljava/lang/String;
	ghelpers.MethodSignatures["java/util/zip/Deflater.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  deflaterToString,
		}

	// equals(Ljava/lang/Object;)Z
	ghelpers.MethodSignatures["java/util/zip/Deflater.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  deflaterEquals,
		}

	// hashCode()I
	ghelpers.MethodSignatures["java/util/zip/Deflater.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  deflaterHashCode,
		}

	// Standard Object lifecycle & concurrency methods
	ghelpers.MethodSignatures["java/util/zip/Deflater.getClass()Ljava/lang/Class;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/Deflater.notify()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/Deflater.notifyAll()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/Deflater.wait()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/Deflater.wait(J)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/zip/Deflater.wait(JI)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
}

func deflaterClinit(params []any) any {
	_ = statics.AddStatic(deflaterClassName+".DEFLATED", statics.Static{
		Type:  types.Int,
		Value: DeflaterDeflated,
	})
	_ = statics.AddStatic(deflaterClassName+".NO_COMPRESSION", statics.Static{
		Type:  types.Int,
		Value: DeflaterNoCompression,
	})
	_ = statics.AddStatic(deflaterClassName+".BEST_SPEED", statics.Static{
		Type:  types.Int,
		Value: DeflaterBestSpeed,
	})
	_ = statics.AddStatic(deflaterClassName+".BEST_COMPRESSION", statics.Static{
		Type:  types.Int,
		Value: DeflaterBestCompression,
	})
	_ = statics.AddStatic(deflaterClassName+".DEFAULT_COMPRESSION", statics.Static{
		Type:  types.Int,
		Value: DeflaterDefaultCompression,
	})
	_ = statics.AddStatic(deflaterClassName+".FILTERED", statics.Static{
		Type:  types.Int,
		Value: DeflaterFiltered,
	})
	_ = statics.AddStatic(deflaterClassName+".HUFFMAN_ONLY", statics.Static{
		Type:  types.Int,
		Value: DeflaterHuffmanOnly,
	})
	_ = statics.AddStatic(deflaterClassName+".DEFAULT_STRATEGY", statics.Static{
		Type:  types.Int,
		Value: DeflaterDefaultStrategy,
	})
	_ = statics.AddStatic(deflaterClassName+".NO_FLUSH", statics.Static{
		Type:  types.Int,
		Value: DeflaterNoFlush,
	})
	_ = statics.AddStatic(deflaterClassName+".SYNC_FLUSH", statics.Static{
		Type:  types.Int,
		Value: DeflaterSyncFlush,
	})
	_ = statics.AddStatic(deflaterClassName+".FULL_FLUSH", statics.Static{
		Type:  types.Int,
		Value: DeflaterFullFlush,
	})
	return nil
}

func deflaterInitDefault(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.<init>: null receiver")
	}
	return deflaterInitLevelNowrap([]any{params[0], DeflaterDefaultCompression, types.JavaBoolFalse})
}

func deflaterInitLevel(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.<init>: null receiver")
	}
	level := DeflaterDefaultCompression
	if len(params) >= 2 {
		if lvl, ok := ghelpers.GetInt64Val(params[1]); ok {
			level = lvl
		}
	}
	return deflaterInitLevelNowrap([]any{params[0], level, types.JavaBoolFalse})
}

func deflaterInitLevelNowrap(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.<init>: null receiver")
	}
	self := params[0].(*object.Object)

	level := DeflaterDefaultCompression
	if len(params) >= 2 {
		if lvl, ok := ghelpers.GetInt64Val(params[1]); ok {
			level = lvl
		}
	}
	if level < -1 || level > 9 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid compression level")
	}

	nowrap := false
	if len(params) >= 3 {
		if nw, ok := ghelpers.GetInt64Val(params[2]); ok {
			nowrap = (nw != 0)
		}
	}

	if self.FieldTable == nil {
		self.FieldTable = make(map[string]object.Field)
	}
	self.FieldTable["level"] = object.Field{Ftype: types.Int, Fvalue: level}
	self.FieldTable["strategy"] = object.Field{Ftype: types.Int, Fvalue: DeflaterDefaultStrategy}
	self.FieldTable["nowrap"] = object.Field{Ftype: "Z", Fvalue: types.ConvertGoBoolToJavaBool(nowrap)}
	self.FieldTable["finished"] = object.Field{Ftype: "Z", Fvalue: types.JavaBoolFalse}

	state := &deflaterState{
		level:        int(level),
		strategy:     int(DeflaterDefaultStrategy),
		nowrap:       nowrap,
		outCollector: new(bytes.Buffer),
		adler:        1,
	}
	setDeflaterState(self, state)
	return nil
}

func deflaterSetInputRange(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.setInput: null receiver")
	}
	self := params[0].(*object.Object)
	state := getDeflaterState(self)
	if state == nil || state.ended {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater has been closed")
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
		return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "missing offset or length")
	}
	offVal, ok1 := ghelpers.GetInt64Val(params[2])
	lenVal, ok2 := ghelpers.GetInt64Val(params[3])
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "invalid offset or length")
	}

	if offVal < 0 || lenVal < 0 || offVal > int64(len(javaBytes))-lenVal {
		return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, fmt.Sprintf("offset=%d length=%d arrayLength=%d", offVal, lenVal, len(javaBytes)))
	}

	rawBytes := object.GoByteArrayFromJavaByteArray(javaBytes[offVal : offVal+lenVal])

	state.mu.Lock()
	defer state.mu.Unlock()
	state.input = append([]byte(nil), rawBytes...)
	return nil
}

func deflaterSetInputArray(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.setInput: null receiver")
	}
	self := params[0].(*object.Object)
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
	return deflaterSetInputRange([]any{self, bufObj, int64(0), int64(len(javaBytes))})
}

func deflaterSetDictionaryRange(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.setDictionary: null receiver")
	}
	self := params[0].(*object.Object)
	state := getDeflaterState(self)
	if state == nil || state.ended {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater has been closed")
	}

	if len(params) < 2 || params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "dictionary is null")
	}
	bufObj, ok := params[1].(*object.Object)
	if !ok || bufObj == nil || object.IsNull(bufObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "dictionary is null")
	}

	valFld, ok := bufObj.FieldTable["value"]
	if !ok || valFld.Fvalue == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "dictionary lacks value field")
	}
	javaBytes, ok := valFld.Fvalue.([]types.JavaByte)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "dictionary value is not JavaByte array")
	}

	if len(params) < 4 {
		return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "missing offset or length")
	}
	offVal, ok1 := ghelpers.GetInt64Val(params[2])
	lenVal, ok2 := ghelpers.GetInt64Val(params[3])
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "invalid offset or length")
	}

	if offVal < 0 || lenVal < 0 || offVal > int64(len(javaBytes))-lenVal {
		return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, fmt.Sprintf("offset=%d length=%d arrayLength=%d", offVal, lenVal, len(javaBytes)))
	}

	rawBytes := object.GoByteArrayFromJavaByteArray(javaBytes[offVal : offVal+lenVal])

	state.mu.Lock()
	defer state.mu.Unlock()
	state.dict = append([]byte(nil), rawBytes...)
	state.dictSet = true
	state.adler = adler32.Checksum(state.dict)
	return nil
}

func deflaterSetDictionaryArray(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.setDictionary: null receiver")
	}
	self := params[0].(*object.Object)
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "dictionary is null")
	}
	bufObj, ok := params[1].(*object.Object)
	if !ok || bufObj == nil || object.IsNull(bufObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "dictionary is null")
	}
	valFld, ok := bufObj.FieldTable["value"]
	if !ok || valFld.Fvalue == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "dictionary lacks value field")
	}
	javaBytes, ok := valFld.Fvalue.([]types.JavaByte)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "dictionary value is not JavaByte array")
	}
	return deflaterSetDictionaryRange([]any{self, bufObj, int64(0), int64(len(javaBytes))})
}

func deflaterSetStrategy(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.setStrategy: null receiver")
	}
	self := params[0].(*object.Object)
	state := getDeflaterState(self)
	if state == nil || state.ended {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater has been closed")
	}

	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid compression strategy")
	}
	stratVal, ok := ghelpers.GetInt64Val(params[1])
	if !ok || stratVal < 0 || stratVal > 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid compression strategy")
	}

	state.mu.Lock()
	defer state.mu.Unlock()
	state.strategy = int(stratVal)
	self.FieldTable["strategy"] = object.Field{Ftype: types.Int, Fvalue: stratVal}
	return nil
}

func deflaterSetLevel(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.setLevel: null receiver")
	}
	self := params[0].(*object.Object)
	state := getDeflaterState(self)
	if state == nil || state.ended {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater has been closed")
	}

	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid compression level")
	}
	levelVal, ok := ghelpers.GetInt64Val(params[1])
	if !ok || levelVal < -1 || levelVal > 9 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid compression level")
	}

	state.mu.Lock()
	defer state.mu.Unlock()
	state.level = int(levelVal)
	self.FieldTable["level"] = object.Field{Ftype: types.Int, Fvalue: levelVal}
	return nil
}

func deflaterNeedsInput(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.needsInput: null receiver")
	}
	self := params[0].(*object.Object)
	state := getDeflaterState(self)
	if state == nil || state.ended {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater has been closed")
	}

	state.mu.Lock()
	defer state.mu.Unlock()
	if len(state.input) == 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func deflaterFinish(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.finish: null receiver")
	}
	self := params[0].(*object.Object)
	state := getDeflaterState(self)
	if state == nil || state.ended {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater has been closed")
	}

	state.mu.Lock()
	defer state.mu.Unlock()
	state.finish = true
	return nil
}

func deflaterFinished(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.finished: null receiver")
	}
	self := params[0].(*object.Object)
	state := getDeflaterState(self)
	if state == nil || state.ended {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater has been closed")
	}

	state.mu.Lock()
	defer state.mu.Unlock()
	if state.finished {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func deflaterDeflateRangeFlush(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.deflate: null receiver")
	}
	self := params[0].(*object.Object)
	state := getDeflaterState(self)
	if state == nil || state.ended {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater has been closed")
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
		return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "missing offset or length")
	}
	offVal, ok1 := ghelpers.GetInt64Val(params[2])
	lenVal, ok2 := ghelpers.GetInt64Val(params[3])
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "invalid offset or length")
	}

	if offVal < 0 || lenVal < 0 || offVal > int64(len(javaBytes))-lenVal {
		return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, fmt.Sprintf("offset=%d length=%d arrayLength=%d", offVal, lenVal, len(javaBytes)))
	}

	flushVal := DeflaterNoFlush
	if len(params) >= 5 {
		if f, ok := ghelpers.GetInt64Val(params[4]); ok {
			flushVal = f
		}
	}
	if flushVal != DeflaterNoFlush && flushVal != DeflaterSyncFlush && flushVal != DeflaterFullFlush {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "invalid flush mode")
	}

	if lenVal == 0 {
		return int64(0)
	}

	tempBuf := make([]byte, lenVal)

	state.mu.Lock()
	n, err := state.deflate(tempBuf, flushVal)
	state.mu.Unlock()

	if err != nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, err.Error())
	}

	for i := 0; i < n; i++ {
		javaBytes[offVal+int64(i)] = types.JavaByte(tempBuf[i])
	}

	return int64(n)
}

func deflaterDeflateRange(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.deflate: null receiver")
	}
	self := params[0].(*object.Object)
	bufObj := params[1]
	offVal := params[2]
	lenVal := params[3]
	return deflaterDeflateRangeFlush([]any{self, bufObj, offVal, lenVal, DeflaterNoFlush})
}

func deflaterDeflateArray(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.deflate: null receiver")
	}
	self := params[0].(*object.Object)
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
	return deflaterDeflateRangeFlush([]any{self, bufObj, int64(0), int64(len(javaBytes)), DeflaterNoFlush})
}

func deflaterGetAdler(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.getAdler: null receiver")
	}
	self := params[0].(*object.Object)
	state := getDeflaterState(self)
	if state == nil || state.ended {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater has been closed")
	}

	state.mu.Lock()
	defer state.mu.Unlock()
	return int64(int32(state.adler))
}

func deflaterGetTotalIn(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.getTotalIn: null receiver")
	}
	self := params[0].(*object.Object)
	state := getDeflaterState(self)
	if state == nil || state.ended {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater has been closed")
	}

	state.mu.Lock()
	defer state.mu.Unlock()
	return int64(int32(state.bytesRead))
}

func deflaterGetBytesRead(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.getBytesRead: null receiver")
	}
	self := params[0].(*object.Object)
	state := getDeflaterState(self)
	if state == nil || state.ended {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater has been closed")
	}

	state.mu.Lock()
	defer state.mu.Unlock()
	return state.bytesRead
}

func deflaterGetTotalOut(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.getTotalOut: null receiver")
	}
	self := params[0].(*object.Object)
	state := getDeflaterState(self)
	if state == nil || state.ended {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater has been closed")
	}

	state.mu.Lock()
	defer state.mu.Unlock()
	return int64(int32(state.bytesWritten))
}

func deflaterGetBytesWritten(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.getBytesWritten: null receiver")
	}
	self := params[0].(*object.Object)
	state := getDeflaterState(self)
	if state == nil || state.ended {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater has been closed")
	}

	state.mu.Lock()
	defer state.mu.Unlock()
	return state.bytesWritten
}

func deflaterReset(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater.reset: null receiver")
	}
	self := params[0].(*object.Object)
	state := getDeflaterState(self)
	if state == nil || state.ended {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Deflater has been closed")
	}

	state.mu.Lock()
	defer state.mu.Unlock()
	state.reset()
	return nil
}

func deflaterEnd(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return nil
	}
	self := params[0].(*object.Object)
	state := getDeflaterState(self)
	if state != nil {
		state.mu.Lock()
		state.end()
		state.mu.Unlock()
	}
	return nil
}

func deflaterFinalize(params []any) any {
	return deflaterEnd(params)
}

func deflaterToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return object.StringObjectFromGoString("null")
	}
	self := params[0].(*object.Object)
	return object.StringObjectFromGoString(fmt.Sprintf("%s@%x", deflaterClassName, uintptr(unsafe.Pointer(self))))
}

func deflaterEquals(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0]) {
		return types.JavaBoolFalse
	}
	if params[0] == params[1] {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func deflaterHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return int64(0)
	}
	self := params[0].(*object.Object)
	return int64(uintptr(unsafe.Pointer(self)) & 0x7FFFFFFF)
}
