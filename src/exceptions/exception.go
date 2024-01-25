/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-4 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package exceptions

import (
	"jacobin/globals"
	"jacobin/log"
	"jacobin/shutdown"
	"jacobin/thread"
	"runtime/debug"
)

// List of Java exceptions (as of Java 17)
// -----------------------------------------------------------------------//
// IMPORTANT: Do not modify this list unless you modify the corresponding
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
// IMPORTANT: Do not modify this list unless you modify the corresponding
// entries in the preceding list of constants in exactly the same way.
// This table and that list must be kept strictly in sync.
var JVMexceptionNames = []string{
	"java.lang.AnnotationTypeMismatchException.class",
	"java.lang.ArithmeticException.class",
	"java.lang.ArrayIndexOutOfBoundsException.class", // added (from Java 8 iastore spec)
	"java.lang.ArrayStoreException.class",
	"java.lang.BufferOverflowException.class",
	"java.lang.BufferUnderflowException.class",
	"java.lang.CannotRedoException.class",
	"java.lang.CannotUndoException.class",
	"java.lang.CatalogException.class",
	"java.lang.ClassCastException.class",
	"java.lang.ClassNotPreparedException.class",
	"java.lang.CMMException.class",
	"java.lang.CompletionException.class",
	"java.lang.ConcurrentModificationException.class",
	"java.lang.DateTimeException.class",
	"java.lang.DOMException.class",
	"java.lang.DuplicateRequestException.class",
	"java.lang.EmptyStackException.class",
	"java.lang.EnumConstantNotPresentException.class",
	"java.lang.EventException.class",
	"java.lang.FileSystemAlreadyExistsException.class",
	"java.lang.FileSystemNotFoundException.class",
	"java.lang.FindException.class",
	"java.lang.IllegalArgumentException.class",
	"java.lang.IllegalCallerException.class",
	"java.lang.IllegalMonitorStateException.class",
	"java.lang.IllegalPathStateException.class",
	"java.lang.IllegalStateException.class",
	"java.lang.IllformedLocaleException.class",
	"java.lang.ImagingOpException.class",
	"java.lang.InaccessibleObjectException.class",
	"java.lang.IncompleteAnnotationException.class",
	"java.lang.InconsistentDebugInfoException.class",
	"java.lang.IndexOutOfBoundsException.class",
	"java.lang.InternalException.class",
	"java.lang.InvalidCodeIndexException.class",
	"java.lang.InvalidLineNumberException.class",
	"java.lang.InvalidModuleDescriptorException.class",
	"java.lang.InvalidModuleException.class",
	"java.lang.InvalidRequestStateException.class",
	"java.lang.InvalidStackFrameException.class",
	"java.lang.JarSignerException.class",
	"java.lang.JMRuntimeException.class",
	"java.lang.JSException.class",
	"java.lang.LayerInstantiationException.class",
	"java.lang.LSException.class",
	"java.lang.MalformedParameterizedTypeException.class",
	"java.lang.MalformedParametersException.class",
	"java.lang.MirroredTypesException.class",
	"java.lang.MissingResourceException.class",
	"java.lang.NativeMethodException.class",
	"java.lang.NegativeArraySizeException.class",
	"java.lang.NoSuchDynamicMethodException.class",
	"java.lang.NoSuchElementException.class",
	"java.lang.NoSuchMechanismException.class",
	"java.lang.NullPointerException.class",
	"java.lang.NumberFormatException.class",
	"java.lang.ObjectCollectedException.class",
	"java.lang.ProfileDataException.class",
	"java.lang.ProviderException.class",
	"java.lang.ProviderNotFoundException.class",
	"java.lang.RangeException.class",
	"java.lang.RasterFormatException.class",
	"java.lang.RejectedExecutionException.class",
	"java.lang.ResolutionException.class",
	"java.lang.SecurityException.class",
	"java.lang.SPIResolutionException.class",
	"java.lang.TypeNotPresentException.class",
	"java.lang.UncheckedIOException.class",
	"java.lang.UndeclaredThrowableException.class",
	"java.lang.UnknownEntityException.class",
	"java.lang.UnmodifiableModuleException.class",
	"java.lang.UnmodifiableSetException.class",
	"java.lang.UnsupportedOperationException.class",
	"java.lang.VMDisconnectedException.class",
	"java.lang.VMMismatchException.class",
	"java.lang.VMOutOfMemoryException.class",
	"java.lang.WrongMethodTypeException.class",
	"java.lang.XPathException.class",

	// non-runtime exceptions
	"java.lang.AbsentInformationException.class",
	"java.lang.AclNotFoundException.class",
	"java.lang.ActivationException.class",
	"java.lang.AgentInitializationException.class",
	"java.lang.AgentLoadException.class",
	"java.lang.AlreadyBoundException.class",
	"java.lang.AttachNotSupportedException.class",
	"java.lang.AWTException.class",
	"java.lang.BackingStoreException.class",
	"java.lang.BadAttributeValueExpException.class",
	"java.lang.BadBinaryOpValueExpException.class",
	"java.lang.BadLocationException.class",
	"java.lang.BadStringOperationException.class",
	"java.lang.BrokenBarrierException.class",
	"java.lang.CardException.class",
	"java.lang.CertificateException.class",
	"java.lang.ClassNotLoadedException.class",
	"java.lang.CloneNotSupportedException.class",
	"java.lang.DataFormatException.class",
	"java.lang.DatatypeConfigurationException.class",
	"java.lang.DestroyFailedException.class",
	"java.lang.ExecutionControlException.class",
	"java.lang.ExecutionException.class",
	"java.lang.ExpandVetoException.class",
	"java.lang.FontFormatException.class",
	"java.lang.GeneralSecurityException.class",
	"java.lang.GSSException.class",
	"java.lang.IllegalClassFormatException.class",
	"java.lang.IllegalConnectorArgumentsException.class",
	"java.lang.IncompatibleThreadStateException.class",
	"java.lang.InterruptedException.class",
	"java.lang.IntrospectionException.class",
	"java.lang.InvalidApplicationException.class",
	"java.lang.InvalidMidiDataException.class",
	"java.lang.InvalidPreferencesFormatException.class",
	"java.lang.InvalidTargetObjectTypeException.class",
	"java.lang.InvalidTypeException.class",
	"java.lang.InvocationException.class",
	"java.lang.IOException.class",
	"java.lang.JMException.class",
	"java.lang.JShellException.class",
	"java.lang.KeySelectorException.class",
	"java.lang.LambdaConversionException.class",
	"java.lang.LastOwnerException.class",
	"java.lang.LineUnavailableException.class",
	"java.lang.MarshalException.class",
	"java.lang.MidiUnavailableException.class",
	"java.lang.MimeTypeParseException.class",
	"java.lang.NamingException.class",
	"java.lang.NoninvertibleTransformException.class",
	"java.lang.NotBoundException.class",
	"java.lang.NotOwnerException.class",
	"java.lang.ParseException.class",
	"java.lang.ParserConfigurationException.class",
	"java.lang.PrinterException.class",
	"java.lang.PrintException.class",
	"java.lang.PrivilegedActionException.class",
	"java.lang.PropertyVetoException.class",
	"java.lang.ReflectiveOperationException.class",
	"java.lang.RefreshFailedException.class",
	"java.lang.RuntimeException.class",
	"java.lang.SAXException.class",
	"java.lang.ScriptException.class",
	"java.lang.ServerNotActiveException.class",
	"java.lang.SQLException.class",
	"java.lang.StringConcatException.class",
	"java.lang.StringIndexOutOfBoundsException.class",
	"java.lang.TimeoutException.class",
	"java.lang.TooManyListenersException.class",
	"java.lang.TransformerException.class",
	"java.lang.TransformException.class",
	"java.lang.UnmodifiableClassException.class",
	"java.lang.UnsupportedAudioFileException.class",
	"java.lang.UnsupportedCallbackException.class",
	"java.lang.UnsupportedFlavorException.class",
	"java.lang.UnsupportedLookAndFeelException.class",
	"java.lang.URIReferenceException.class",
	"java.lang.URISyntaxException.class",
	"java.lang.VMStartException.class",
	"java.lang.XAException.class",
	"java.lang.XMLParseException.class",
	"java.lang.XMLSignatureException.class",
	"java.lang.XMLStreamException.class",

	// Java exceptions
	"java.lang.AnnotationFormatError.class",
	"java.lang.AssertionError.class",
	"java.lang.AWTError.class",
	"java.lang.CoderMalfunctionError.class",
	"java.lang.FactoryConfigurationError.class",
	"java.lang.IOError.class",
	"java.lang.LinkageError.class",
	"java.lang.SchemaFactoryConfigurationError.class",
	"java.lang.ServiceConfigurationError.class",
	"java.lang.ThreadDeath.class",
	"java.lang.TransformerFactoryConfigurationError.class",
	"java.lang.VirtualMachineError.class",

	"java.lang.ChangedCharSetException.class",
	"java.lang.CharacterCodingException.class",
	"java.lang.CharConversionException.class",
	"java.lang.UnsupportedEncodingException.class",
	"java.lang.UTFDataFormatException.class",
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
	_ = log.Log(msg, log.SEVERE)

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
