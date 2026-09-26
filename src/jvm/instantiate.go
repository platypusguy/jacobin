/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-4 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) Consult jacobin.org.
 */

package jvm

import (
	"container/list"
	"errors"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/gfunction/javaLang"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"jacobin/src/util"
	"strings"
	"sync"
	"sync/atomic"
)

// instantiating an object is a two-part process (except for arrays, which are handled
// by special bytecodes):
//
//  1. the class needs to be loaded, so that its details and its methods are knowable
//
//  2. the class fields (if static) and instance fields (if non-static) are allocated.
//     Details for this second step appear in the loop that drives createField().
//
//     NOTE: The "any" type returned is always *object.Object.
//     This is being done to avoid a golang circularity error when the caller
//     is one of the native 'G' functions.
type initState int

const (
	initNone initState = iota
	initInProgress
	initDone
	initFailed
)

type classInit struct {
	done  atomic.Bool // lock-free fast path
	state initState   // guarded by initMu
	owner *list.List  // frameStack of the initializing thread; guarded by initMu
}

var (
	initTable sync.Map // classname -> *classInit
	initMu    sync.Mutex
	initCond  = sync.NewCond(&initMu)
)

func getClassInit(name string) *classInit {
	v, _ := initTable.LoadOrStore(name, &classInit{})
	return v.(*classInit)
}

// ensureInitialized implements the JVMS 5.5 procedure: exactly one thread runs
// initFn; others wait; the initializing thread re-entering (e.g. from within
// its own <clinit>) proceeds without blocking.
func ensureInitialized(classname string, fs *list.List, initFn func() error) error {
	ci := getClassInit(classname)
	// If the class was previously initialized, but the class area was reset
	// in the meantime (e.g. classloader.Init() called again, as happens
	// between unit tests), the class won't actually be present anymore.
	// In that case, don't trust the stale "done" state--reinitialize it.
	if ci.done.Load() && classloader.MethAreaFetch(classname) != nil {
		return nil
	}

	initMu.Lock()
	for {
		switch ci.state {
		case initDone:
			if classloader.MethAreaFetch(classname) != nil {
				initMu.Unlock()
				return nil
			}
			ci.state = initNone // stale: class area was reset since init
			ci.done.Store(false)
			continue
		case initFailed:
			initMu.Unlock()
			return fmt.Errorf("could not initialize class %s (previous initialization failed)", classname)
		case initInProgress:
			if ci.owner == fs {
				initMu.Unlock()
				return nil
			}
			initCond.Wait()
		default: // initNone: this goroutine wins the right to initialize
			ci.state = initInProgress
			ci.owner = fs
			initMu.Unlock() // never hold initMu while running <clinit>
			return runGuarded(ci, initFn)
		}
	}
}

func runGuarded(ci *classInit, initFn func() error) error {
	ok := false
	defer func() { // also runs on panic, so waiters are never stranded
		initMu.Lock()
		ci.owner = nil
		if ok {
			ci.state = initDone
			ci.done.Store(true)
		} else {
			ci.state = initFailed
		}
		initMu.Unlock()
		initCond.Broadcast()
	}()
	if err := initFn(); err != nil {
		return err
	}
	ok = true
	return nil
}

// superclassChain returns k's superclass names, nearest first, stopping
// before java/lang/Object. Each superclass is loaded on the way up.
func superclassChain(k *classloader.Klass, classname string) ([]string, error) {
	superclasses := []string{}
	superclassNamePtr := stringPool.GetStringPointer(k.Data.SuperclassIndex)
	for {
		if classname == types.ObjectClassName || *superclassNamePtr == types.ObjectClassName {
			break
		}
		if err := loadThisClass(*superclassNamePtr); err != nil {
			return nil, err // error message will have been displayed
		}
		superclasses = append(superclasses, *superclassNamePtr)

		loadedSuperclass := classloader.MethAreaFetch(*superclassNamePtr)
		if loadedSuperclass == nil || loadedSuperclass.Data == nil {
			return nil, fmt.Errorf("superclassChain: %s is nil after loading", *superclassNamePtr)
		}
		superclassNamePtr = stringPool.GetStringPointer(loadedSuperclass.Data.SuperclassIndex)
	}
	return superclasses, nil
}

func doStaticDefaults(k *classloader.Klass, classname string) {
	for i := 0; i < len(k.Data.Fields); i++ {
		fld := k.Data.Fields[i]
		if !fld.IsStatic {
			continue
		}
		fldName := k.Data.CP.Utf8Refs[fld.Name]
		fldType := []byte(k.Data.CP.Utf8Refs[fld.Desc])

		var fldValue any
		if fld.ConstValue != nil {
			// A final static initialized with a compile-time constant expression
			// carries its value in the ConstantValue attribute; javac generally
			// emits no <clinit> PUTSTATIC for these, so this is the only place
			// the real value is ever installed. Mirrors createField's handling.
			fldValue = fld.ConstValue
		} else {
			switch fldType[0] {
			case 'B', 'C', 'S', 'I', 'J', 'Z':
				fldValue = int64(0)
			case 'F', 'D':
				fldValue = float64(0.00)
			case 'L', '[':
				fldValue = object.Null
			}
		}
		statics.AddStaticIfAbsent(classname+"."+fldName,
			statics.Static{Type: string(fldType[0]), Value: fldValue})
	}
}

// initializeClass performs the JVMS 5.5 class-initialization work exactly
// once per class: code validity check, static-field defaults, then <clinit>.
func initializeClass(k *classloader.Klass, classname string, superclasses []string, fs *list.List) error {
	// (a) code check — skip for JDK classes, as before
	if !util.IsFilePartOfJDK(&classname) {
		for _, m := range k.Data.MethodTable {
			code := m.CodeAttr.Code
			var err error
			if globals.TraceCodeCheck {
				methName := k.Data.CP.Utf8Refs[m.Name]
				methDesc := k.Data.CP.Utf8Refs[m.Desc]
				fullMethodName := fmt.Sprintf("%s.%s%s", k.Data.Name, methName, methDesc)
				err = classloader.CheckCodeValidity(
					&code, &k.Data.CP, m.CodeAttr.MaxStack, m.CodeAttr.MaxLocals, k.Data.Access, &fullMethodName)
			} else {
				err = classloader.CheckCodeValidity(
					&code, &k.Data.CP, m.CodeAttr.MaxStack, m.CodeAttr.MaxLocals, k.Data.Access, nil)
			}
			if err != nil {
				methName := k.Data.CP.Utf8Refs[m.Name]
				methDesc := k.Data.CP.Utf8Refs[m.Desc]
				errMsg := fmt.Sprintf("InitializeClass: CheckCodeValidity failed in %s.%s%s: %s",
					classname, methName, methDesc, err.Error())
				status := exceptions.ThrowEx(excNames.ClassFormatError, errMsg, nil)
				if status != exceptions.Caught {
					return errors.New(errMsg)
				}
			}
		}
		k.CodeChecked = true
		classloader.MethAreaInsert(classname, k)
		if globals.TraceCloadi {
			trace.Trace("InitializeClass: Code checked for class: " + classname)
		}
	}

	// (b) static defaults — now covers classes with superclasses too,
	//     not only the no-superclass case the original code handled.
	doStaticDefaults(k, classname)
	for _, superclassName := range superclasses {
		sk := classloader.MethAreaFetch(superclassName)
		if sk != nil && sk.Data != nil {
			doStaticDefaults(sk, superclassName)
		}
	}

	// (c) <clinit>
	if _, ok := k.Data.MethodTable["<clinit>()V"]; ok {
		// runInitializationBlock treats a non-empty chain as complete and will
		// not add classname itself, so we must prepend it here — mirroring
		// what the original InstantiateClass did before this was split out.
		// An empty chain is passed through unchanged so runInitializationBlock's
		// own rebuild-and-filter path (for classes with no superclasses) still runs.
		initChain := superclasses
		if len(initChain) > 0 {
			initChain = append([]string{classname}, initChain...)
		}
		if err := runInitializationBlock(k, initChain, fs); err != nil {
			errMsg := fmt.Sprintf("InitializeClass: runInitializationBlock failed with %s.<clinit>()V", classname)
			trace.Error(errMsg)
			return err
		}
	}
	return nil
}

// InitializeClass loads classname if needed and runs its class-initialization
// step exactly once, regardless of how many goroutines call it concurrently.
func InitializeClass(classname string, fs *list.List) error {
	// Fast path: no loading, no locks. Note: if the class area was reset
	// since this class was initialized (e.g. classloader.Init() called
	// again, as happens between unit tests), the stale "done" state is
	// not trusted, and the class is loaded/initialized again below.
	if getClassInit(classname).done.Load() && classloader.MethAreaFetch(classname) != nil {
		return nil
	}
	if !strings.HasPrefix(classname, "[") { // arrays are not initialized
		if err := loadThisClass(classname); err != nil {
			return err
		}
	}
	k := classloader.MethAreaFetch(classname)
	if k == nil || k.Data == nil {
		return fmt.Errorf("InitializeClass: class %s unavailable after loading", classname)
	}
	superclasses, err := superclassChain(k, classname)
	if err != nil {
		return err
	}
	return ensureInitialized(classname, fs, func() error {
		return initializeClass(k, classname, superclasses, fs)
	})
}

func InstantiateClass(classname string, frameStack *list.List) (any, error) {

	// Objects that are created by Jacobin itself are instantiated separately
	// e.g. String, Thread, ThreadGroup
	switch classname {
	case types.StringClassName:
		return object.NewStringObject(), nil
	case types.ClassNameThread:
		return javaLang.ThreadCreateObject(nil), nil
	case types.ClassNameThreadGroup:
		return javaLang.MakeThreadGroup(), nil
	}

	// Class initialization (JVMS 5.5) happens before instance creation,
	// and is guaranteed to run exactly once even under concurrent callers.
	if err := InitializeClass(classname, frameStack); err != nil {
		return nil, err
	}

	k := classloader.MethAreaFetch(classname)
	if k == nil {
		errMsg := "InstantiateClass: Class is nil after loading, class: " + classname
		trace.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if k.Data == nil {
		errMsg := "InstantiateClass: class.Data is nil, class: " + classname
		trace.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	// create the object whose instantiation we're doing
	obj := object.MakeEmptyObject()
	obj.KlassName = stringPool.GetStringIndex(&classname)

	superclasses, err := superclassChain(k, classname)
	if err != nil {
		return nil, err
	}

	// handle the fields. If the object has no superclass other than Object,
	// the fields are in an array in the order they're declared in the CP.
	// If the object has a non-Object superclass, then the superclasses' fields
	// and the present object's fields are stored in a map, indexed by field name.
	if len(superclasses) == 0 {
		for i := 0; i < len(k.Data.Fields); i++ {
			fld := k.Data.Fields[i]
			fldName := k.Data.CP.Utf8Refs[fld.Name]
			fieldToAdd, err := createField(fld, k, classname)
			if err != nil {
				return nil, err
			}
			obj.FieldTable[fldName] = *fieldToAdd
		}
		return obj, nil
	}

	// in the case of superclasses, we start at the topmost superclass and
	// work our way down to the present class, adding fields to FieldTable.
	chain := append([]string{classname}, superclasses...)
	for j := len(chain) - 1; j >= 0; j-- {
		superclassName := chain[j]
		c := classloader.MethAreaFetch(superclassName)
		if c == nil {
			errMsg := fmt.Sprintf("InstantiateClass: MethAreaFetch(superclass: %s) failed", superclassName)
			trace.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		for i := 0; i < len(c.Data.Fields); i++ {
			f := c.Data.Fields[i]
			name := c.Data.CP.Utf8Refs[f.Name]
			fieldToAdd, err := createField(f, c, classname)
			if err != nil {
				return nil, err
			}
			obj.FieldTable[name] = *fieldToAdd
		}
	}

	return obj, nil
}

// creates a field for insertion into the object representation
func createField(f classloader.Field, k *classloader.Klass, classname string) (*object.Field, error) {
	desc := k.Data.CP.Utf8Refs[f.Desc]

	fieldToAdd := new(object.Field)
	fieldToAdd.Ftype = desc
	switch string(fieldToAdd.Ftype[0]) {
	case types.Ref, types.Array: // it's a reference
		fieldToAdd.Fvalue = nil
	case types.Byte:
		fieldToAdd.Fvalue = int8(0)
	case types.Char, types.Int, types.Long, types.Short, types.Bool:
		fieldToAdd.Fvalue = int64(0)
	case types.Double, types.Float:
		fieldToAdd.Fvalue = 0.0
	default:
		errMsg := fmt.Sprintf("createField: error creating field in: %s,  Invalid type: %s",
			classname, fieldToAdd.Ftype)
		trace.Error(errMsg)
		return nil, classloader.CFE(errMsg)
	}

	presentType := fieldToAdd.Ftype
	if f.IsStatic {
		// in the instantiated class, add a types.Static before the
		// type, which notifies future users that the field
		// is static and should be fetched from the Statics
		// table. TODO: This can probably be removed.
		fieldToAdd.Ftype = types.Static + presentType
		if f.ConstValue != nil { // if the field has a constant value, set it
			fieldToAdd.Fvalue = f.ConstValue
		}
	}

	/* The following code is no longer needed as it handled only the ConstValue attribute.
		Howvever, at a future point, we may want to add the ability process other field attributes
		and this code will give us a template for doing so.
	    //
		// static fields can have ConstantValue attributes,
		// which specify their initial value.
		if len(f.Attributes) > 0 {
			for j := 0; j < len(f.Attributes); j++ {
				attr := k.Data.CP.Utf8Refs[int(f.Attributes[j].AttrName)]
				if attr == "ConstantValue" && f.IsStatic { // only statics can have ConstantValue attribute
					valueIndex := int(f.Attributes[j].AttrContent[0])*256 +
						int(f.Attributes[j].AttrContent[1])
					valueType := k.Data.CP.CpIndex[valueIndex].Type
					valueSlot := k.Data.CP.CpIndex[valueIndex].Slot
					switch valueType {
					case classloader.IntConst:
						fieldToAdd.Fvalue = int64(k.Data.CP.IntConsts[valueSlot])
					case classloader.LongConst:
						fieldToAdd.Fvalue = k.Data.CP.LongConsts[valueSlot]
					case classloader.FloatConst:
						fieldToAdd.Fvalue = float64(k.Data.CP.Floats[valueSlot])
					case classloader.DoubleConst:
						fieldToAdd.Fvalue = k.Data.CP.Doubles[valueSlot]
					case classloader.StringConst:
						str := k.Data.CP.Utf8Refs[valueSlot]
						fieldToAdd.Fvalue = object.StringObjectFromGoString(str)
					default:
						errMsg := fmt.Sprintf(
							"createField: Unexpected ConstantValue type in instantiate: %d", valueType)
						trace.Error(errMsg)
						return nil, errors.New(errMsg)
					} // end of ConstantValue type switch
				} // end of ConstantValue attribute processing
			} // end of processing attributes
		} // end of search through attributes
	*/

	if f.IsStatic {
		s := statics.Static{
			Type:  presentType, // we use the type without the 'X' prefix in the statics table.
			Value: fieldToAdd.Fvalue,
		}
		// add the field to the Statics table
		fieldName := k.Data.CP.Utf8Refs[f.Name]

		_, alreadyPresent := statics.QueryStatic(classname, fieldName)
		if !alreadyPresent { // add only if the field has not been pre-loaded
			_ = statics.AddStatic(classname+"."+fieldName, s)
		}
	}
	return fieldToAdd, nil
}

// Loads the class (if it's not already loaded) and makes sure it's accessible in the method area
func loadThisClass(className string) error {
	alreadyLoaded := classloader.MethAreaFetch(className)
	if alreadyLoaded != nil { // if the class is already loaded, skip the rest of this
		return nil
	}
	// Try to load class by name
	err := classloader.LoadClassFromNameOnly(className)
	if err != nil {
		return errors.New(err.Error())
	}
	// Success in loaded by name
	if globals.TraceCloadi {
		trace.Trace("loadThisClass: Success in LoadClassFromNameOnly(" + className + ")")
	}

	// at this point the class has been loaded into the method area (MethArea). Wait for it to be ready.
	err = classloader.WaitForClassStatus(className)
	if err != nil {
		errMsg := fmt.Sprintf("loadThisClass: WaitForClassStatus(%s) failed, err: %v", className, err)
		trace.Error(errMsg)
		return errors.New(errMsg) // needed for testing, which does not shutdown on failure
	}
	return nil
}
