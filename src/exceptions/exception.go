/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-4 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package exceptions

import (
	"fmt"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/shutdown"
	"jacobin/statics"
	"jacobin/thread"
	"runtime/debug"
)

// List of Java exceptions (as of Java 17)

// -----------------------------------------------------------------------//
// IMPORTANT! Do not modify this list unless you modify the corresponding
// entries in JVMExceptionNames in exactly the same way.
// This list and that table must be kept strictly in sync.
const (
	Unknown = iota

	// runtime exceptions
	AnnotationTypeMismatchException
	ArithmeticException
	ArrayIndexOutOfBoundsException // added (from Java 8 iastore spec)
	ArrayStoreException
	BufferOverflowException
	BufferUnderflowException
	CannotRedoException
	CannotUndoException
	CatalogException
	ClassCastException
	ClassNotPreparedException
	CMMException
	CompletionException
	ConcurrentModificationException
	DateTimeException
	DOMException
	DuplicateRequestException
	EmptyStackException
	EnumConstantNotPresentException
	EventException
	FileNotFoundException
	FileSystemAlreadyExistsException
	FileSystemNotFoundException
	FindException
	IllegalArgumentException
	IllegalCallerException
	IllegalMonitorStateException
	IllegalPathStateException
	IllegalStateException
	IllformedLocaleException
	ImagingOpException
	InaccessibleObjectException
	IncompleteAnnotationException
	InconsistentDebugInfoException
	IndexOutOfBoundsException
	InternalException
	InvalidCodeIndexException
	InvalidLineNumberException
	InvalidModuleDescriptorException
	InvalidModuleException
	InvalidRequestStateException
	InvalidStackFrameException
	JarSignerException
	JMRuntimeException
	JSException
	LayerInstantiationException
	LSException
	MalformedParameterizedTypeException
	MalformedParametersException
	MirroredTypesException
	MissingResourceException
	NativeMethodException
	NegativeArraySizeException
	NoSuchDynamicMethodException
	NoSuchElementException
	NoSuchMechanismException
	NullPointerException
	NumberFormatException
	ObjectCollectedException
	ProfileDataException
	ProviderException
	ProviderNotFoundException
	RangeException
	RasterFormatException
	RejectedExecutionException
	ResolutionException
	SecurityException
	SPIResolutionException
	TypeNotPresentException
	UncheckedIOException
	UndeclaredThrowableException
	UnknownEntityException
	UnmodifiableModuleException
	UnmodifiableSetException
	UnsupportedOperationException
	VMDisconnectedException
	VMMismatchException
	VMOutOfMemoryException
	WrongMethodTypeException
	XPathException

	// non-runtime exceptions
	AbsentInformationException
	AclNotFoundException
	ActivationException
	AgentInitializationException
	AgentLoadException
	AlreadyBoundException
	AttachNotSupportedException
	AWTException
	BackingStoreException
	BadAttributeValueExpException
	BadBinaryOpValueExpException
	BadLocationException
	BadStringOperationException
	BrokenBarrierException
	CardException
	CertificateException
	ClassNotLoadedException
	CloneNotSupportedException
	DataFormatException
	DatatypeConfigurationException
	DestroyFailedException
	ExecutionControlException
	ExecutionException
	ExpandVetoException
	FontFormatException
	GeneralSecurityException
	GSSException
	IllegalClassFormatException
	IllegalConnectorArgumentsException
	IncompatibleThreadStateException
	InterruptedException
	IntrospectionException
	InvalidApplicationException
	InvalidMidiDataException
	InvalidPreferencesFormatException
	InvalidTargetObjectTypeException
	InvalidTypeException
	InvocationException
	IOException
	JMException
	JShellException
	KeySelectorException
	LambdaConversionException
	LastOwnerException
	LineUnavailableException
	MarshalException
	MidiUnavailableException
	MimeTypeParseException
	NamingException
	NoninvertibleTransformException
	NotBoundException
	NotOwnerException
	ParseException
	ParserConfigurationException
	PrinterException
	PrintException
	PrivilegedActionException
	PropertyVetoException
	ReflectiveOperationException
	RefreshFailedException
	RuntimeException
	SAXException
	ScriptException
	ServerNotActiveException
	SQLException
	StringConcatException
	StringIndexOutOfBoundsException
	TimeoutException
	TooManyListenersException
	TransformerException
	TransformException
	UnmodifiableClassException
	UnsupportedAudioFileException
	UnsupportedCallbackException
	UnsupportedFlavorException
	UnsupportedLookAndFeelException
	URIReferenceException
	URISyntaxException
	VMStartException
	XAException
	XMLParseException
	XMLSignatureException
	XMLStreamException

	// Java exceptions
	AnnotationFormatError
	AssertionError
	AWTError
	CoderMalfunctionError
	FactoryConfigurationError
	IOError
	LinkageError
	SchemaFactoryConfigurationError
	ServiceConfigurationError
	ThreadDeath
	TransformerFactoryConfigurationError
	VirtualMachineError

	// Character set exceptions
	ChangedCharSetException
	CharacterCodingException
	CharConversionException
	UnsupportedEncodingException
	UTFDataFormatException
)

// -----------------------------------------------------------------------//
// IMPORTANT! Do not modify this list unless you modify the corresponding
// entries in the preceding list of constants in exactly the same way.
// This table and that list must be kept strictly in sync.
var JVMexceptionNames = []string{
	"", // no exception (present because list of consts begins at 1)
	"java.lang.AnnotationTypeMismatchException",
	"java.lang.ArithmeticException",
	"java.lang.ArrayIndexOutOfBoundsException", // added (from Java 8 iastore spec)
	"java.lang.ArrayStoreException",
	"java.lang.BufferOverflowException",
	"java.lang.BufferUnderflowException",
	"java.lang.CannotRedoException",
	"java.lang.CannotUndoException",
	"java.lang.CatalogException",
	"java.lang.ClassCastException",
	"java.lang.ClassNotPreparedException",
	"java.lang.CMMException",
	"java.lang.CompletionException",
	"java.lang.ConcurrentModificationException",
	"java.lang.DateTimeException",
	"java.lang.DOMException",
	"java.lang.DuplicateRequestException",
	"java.lang.EmptyStackException",
	"java.lang.EnumConstantNotPresentException",
	"java.lang.EventException",
	"java.io.FileNotFoundException",
	"java.lang.FileSystemAlreadyExistsException",
	"java.lang.FileSystemNotFoundException",
	"java.lang.FindException",
	"java.lang.IllegalArgumentException",
	"java.lang.IllegalCallerException",
	"java.lang.IllegalMonitorStateException",
	"java.lang.IllegalPathStateException",
	"java.lang.IllegalStateException",
	"java.lang.IllformedLocaleException",
	"java.lang.ImagingOpException",
	"java.lang.InaccessibleObjectException",
	"java.lang.IncompleteAnnotationException",
	"java.lang.InconsistentDebugInfoException",
	"java.lang.IndexOutOfBoundsException",
	"java.lang.InternalException",
	"java.lang.InvalidCodeIndexException",
	"java.lang.InvalidLineNumberException",
	"java.lang.InvalidModuleDescriptorException",
	"java.lang.InvalidModuleException",
	"java.lang.InvalidRequestStateException",
	"java.lang.InvalidStackFrameException",
	"java.lang.JarSignerException",
	"java.lang.JMRuntimeException",
	"java.lang.JSException",
	"java.lang.LayerInstantiationException",
	"java.lang.LSException",
	"java.lang.MalformedParameterizedTypeException",
	"java.lang.MalformedParametersException",
	"java.lang.MirroredTypesException",
	"java.lang.MissingResourceException",
	"java.lang.NativeMethodException",
	"java.lang.NegativeArraySizeException",
	"java.lang.NoSuchDynamicMethodException",
	"java.lang.NoSuchElementException",
	"java.lang.NoSuchMechanismException",
	"java.lang.NullPointerException",
	"java.lang.NumberFormatException",
	"java.lang.ObjectCollectedException",
	"java.lang.ProfileDataException",
	"java.lang.ProviderException",
	"java.lang.ProviderNotFoundException",
	"java.lang.RangeException",
	"java.lang.RasterFormatException",
	"java.lang.RejectedExecutionException",
	"java.lang.ResolutionException",
	"java.lang.SecurityException",
	"java.lang.SPIResolutionException",
	"java.lang.TypeNotPresentException",
	"java.lang.UncheckedIOException",
	"java.lang.UndeclaredThrowableException",
	"java.lang.UnknownEntityException",
	"java.lang.UnmodifiableModuleException",
	"java.lang.UnmodifiableSetException",
	"java.lang.UnsupportedOperationException",
	"java.lang.VMDisconnectedException",
	"java.lang.VMMismatchException",
	"java.lang.VMOutOfMemoryException",
	"java.lang.WrongMethodTypeException",
	"java.lang.XPathException",

	// non-runtime exceptions
	"java.lang.AbsentInformationException",
	"java.lang.AclNotFoundException",
	"java.lang.ActivationException",
	"java.lang.AgentInitializationException",
	"java.lang.AgentLoadException",
	"java.lang.AlreadyBoundException",
	"java.lang.AttachNotSupportedException",
	"java.lang.AWTException",
	"java.lang.BackingStoreException",
	"java.lang.BadAttributeValueExpException",
	"java.lang.BadBinaryOpValueExpException",
	"java.lang.BadLocationException",
	"java.lang.BadStringOperationException",
	"java.lang.BrokenBarrierException",
	"java.lang.CardException",
	"java.lang.CertificateException",
	"java.lang.ClassNotLoadedException",
	"java.lang.CloneNotSupportedException",
	"java.lang.DataFormatException",
	"java.lang.DatatypeConfigurationException",
	"java.lang.DestroyFailedException",
	"java.lang.ExecutionControlException",
	"java.lang.ExecutionException",
	"java.lang.ExpandVetoException",
	"java.lang.FontFormatException",
	"java.lang.GeneralSecurityException",
	"java.lang.GSSException",
	"java.lang.IllegalClassFormatException",
	"java.lang.IllegalConnectorArgumentsException",
	"java.lang.IncompatibleThreadStateException",
	"java.lang.InterruptedException",
	"java.lang.IntrospectionException",
	"java.lang.InvalidApplicationException",
	"java.lang.InvalidMidiDataException",
	"java.lang.InvalidPreferencesFormatException",
	"java.lang.InvalidTargetObjectTypeException",
	"com.sun.jdi.InvalidTypeException",
	"java.lang.InvocationException",
	"java.io.IOException",
	"java.lang.JMException",
	"java.lang.JShellException",
	"java.lang.KeySelectorException",
	"java.lang.LambdaConversionException",
	"java.lang.LastOwnerException",
	"java.lang.LineUnavailableException",
	"java.lang.MarshalException",
	"java.lang.MidiUnavailableException",
	"java.lang.MimeTypeParseException",
	"java.lang.NamingException",
	"java.lang.NoninvertibleTransformException",
	"java.lang.NotBoundException",
	"java.lang.NotOwnerException",
	"java.lang.ParseException",
	"java.lang.ParserConfigurationException",
	"java.lang.PrinterException",
	"java.lang.PrintException",
	"java.lang.PrivilegedActionException",
	"java.lang.PropertyVetoException",
	"java.lang.ReflectiveOperationException",
	"java.lang.RefreshFailedException",
	"java.lang.RuntimeException",
	"java.lang.SAXException",
	"java.lang.ScriptException",
	"java.lang.ServerNotActiveException",
	"java.lang.SQLException",
	"java.lang.StringConcatException",
	"java.lang.StringIndexOutOfBoundsException",
	"java.lang.TimeoutException",
	"java.lang.TooManyListenersException",
	"java.lang.TransformerException",
	"java.lang.TransformException",
	"java.lang.UnmodifiableClassException",
	"java.lang.UnsupportedAudioFileException",
	"java.lang.UnsupportedCallbackException",
	"java.lang.UnsupportedFlavorException",
	"java.lang.UnsupportedLookAndFeelException",
	"java.lang.URIReferenceException",
	"java.lang.URISyntaxException",
	"java.lang.VMStartException",
	"java.lang.XAException",
	"java.lang.XMLParseException",
	"java.lang.XMLSignatureException",
	"java.lang.XMLStreamException",

	// Java exceptions
	"java.lang.AnnotationFormatError",
	"java.lang.AssertionError",
	"java.lang.AWTError",
	"java.lang.CoderMalfunctionError",
	"java.lang.FactoryConfigurationError",
	"java.lang.IOError",
	"java.lang.LinkageError",
	"java.lang.SchemaFactoryConfigurationError",
	"java.lang.ServiceConfigurationError",
	"java.lang.ThreadDeath",
	"java.lang.TransformerFactoryConfigurationError",
	"java.lang.VirtualMachineError",

	"java.lang.ChangedCharSetException",
	"java.lang.CharacterCodingException",
	"java.lang.CharConversionException",
	"java.lang.UnsupportedEncodingException",
	"java.lang.UTFDataFormatException",
}

// Throw duplicates the exception mechanism in Java. Right now, it displays the
// exceptions message. Will add: catch logic, stack trace, and halt of execution
// TODO: use ThreadNum to find the right thread
func Throw(exceptionType int, msg string) {
	/* // This code should be moved to the interpreter and the info pushed to this function.
	   func Throw(excType int, clName string, threadNum int, methName string, cp int) {
	   	thd := globals.GetGlobalRef().Threads.ThreadsList.Front().Value.(*thread.ExecThread)
	   	frameStack := thd.Stack
	   	f := frames.PeekFrame(frameStack, 0)
	   	fmt.Println("class name: " + f.ClName)
	   	msg := fmt.Sprintf(
	   		"%s%sin %s, in%s, at bytecode[]: %d", JacobinRuntimeErrLiterals[excType], ": ", clName, methName, cp)
	*/
	helloMsg := fmt.Sprintf("[Throw] Arrived, which: %d, msg: %s", exceptionType, msg)
	log.Log(helloMsg, log.SEVERE)

	// TODO: Temporary until error/exception processing is complete.
	glob := globals.GetGlobalRef()
	if glob.JacobinName == "test" {
		return
	}
	var stack string
	bytes := debug.Stack()
	if len(bytes) > 0 {
		stack = string(bytes)
	} else {
		stack = ""
	}
	glob.ErrorGoStack = stack
	ShowPanicCause(msg)
	ShowFrameStack(&thread.ExecThread{})
	ShowGoStackTrace(nil)
	statics.DumpStatics()
	_ = shutdown.Exit(shutdown.APP_EXCEPTION)
}

// JVMexception reports runtime exceptions occurring in the JVM (rather than in the app)
// such as invalid JAR files, and the like. For the moment, it prints out the exceptions msg
// only. Eventually, it will print out considerably more info depending on the setting of
// globals.JVMstrict. NOTE: this function calls Shutdown(), as all JVM runtime exceptions
// are fatal.
func JVMexception(excType int, msg string) {
	_ = log.Log(msg, log.SEVERE)
	shutdown.Exit(shutdown.JVM_EXCEPTION)
}
