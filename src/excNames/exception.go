/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package excNames

// List of Java exceptions (as of Java 17)

// -----------------------------------------------------------------------//
// IMPORTANT! Do not modify this list unless you modify the corresponding
// entries in JVMExceptionNames in exactly the same way.
// This list and that table must be kept strictly in sync.
const (
	Unknown = iota

	// runtime exceptions
	AccessDeniedException
	AnnotationTypeMismatchException
	ArithmeticException
	ArrayIndexOutOfBoundsException // added (from Java 8 iastore spec)
	ArrayStoreException
	AtomicMoveNotSupportedException
	BufferOverflowException
	BufferUnderflowException
	CannotRedoException
	CannotUndoException
	CatalogException
	ClassCastException
	ClassNotFoundException
	ClassNotPreparedException
	CMMException
	CompletionException
	ConcurrentModificationException
	DateTimeException
	DOMException
	DuplicateFormatFlagsException
	DuplicateRequestException
	EmptyStackException
	EnumConstantNotPresentException
	EventException
	FileNotFoundException
	FileSystemAlreadyExistsException
	FileSystemNotFoundException
	FindException
	FormatFlagsConversionMismatchException
	FormatterClosedException
	IllegalAccessException
	IllegalArgumentException
	IllegalCallerException
	IllegalFormatCodePointException
	IllegalFormatConversionException
	IllegalMonitorStateException
	IllegalPathStateException
	IllegalStateException
	IllformedLocaleException
	ImagingOpException
	InaccessibleObjectException
	IncompleteAnnotationException
	InconsistentDebugInfoException
	IndexOutOfBoundsException
	InstantiationException
	InternalException
	InvalidCodeIndexException
	InvalidLineNumberException
	InvalidMarkException
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
	UserPrincipalNotFoundException
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
	IllegalThreadStateException
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
	NoSuchFieldException
	NoSuchMethodException
	NotBoundException
	NotOwnerException
	ParseException
	ParserConfigurationException
	PrinterException
	PrintException
	PrivilegedActionException
	PropertyVetoException
	ReadOnlyBufferException
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

	// Java errors
	AnnotationFormatError
	AssertionError
	AWTError
	CoderMalfunctionError
	FactoryConfigurationError
	InternalError
	IOError
	LinkageError
	NoClassDefFoundError
	NoSuchFieldError
	NoSuchMethodError
	OutOfMemoryError
	SchemaFactoryConfigurationError
	ServiceConfigurationError
	ThreadDeath
	TransformerFactoryConfigurationError
	UnsatisfiedLinkError
	UnsupportedClassVersionError
	VerifyError
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
	"",                                    // no exception (present because list of consts begins at 1)
	"java.nio.file.AccessDeniedException", // VERIFIED
	"java.lang.annotation.AnnotationTypeMismatchException", // VERIFIED
	"java.lang.ArithmeticException",                        // VERIFIED
	"java.lang.ArrayIndexOutOfBoundsException",             // VERIFIED
	"java.lang.ArrayStoreException",                        // VERIFIED
	"java.nio.file.AtomicMoveNotSupportedException",        // VERIFIED
	"java.nio.BufferOverflowException",                     // VERIFIED
	"java.nio.BufferUnderflowException",                    // VERIFIED
	"javax.swing.undo.CannotRedoException",                 // VERIFIED
	"javax.swing.undo.CannotUndoException",                 // VERIFIED
	"java.lang.CatalogException",
	"java.lang.ClassCastException",
	"java.lang.ClassNotFoundException", // verified
	"java.lang.ClassNotPreparedException",
	"java.lang.CMMException",
	"java.lang.CompletionException",
	"java.util.ConcurrentModificationException", // VERIFIED
	"java.lang.DateTimeException",
	"java.lang.DOMException",
	"java.util.DuplicateFormatFlagsException", // VERIFIED
	"java.lang.DuplicateRequestException",
	"java.util.EmptyStackException",             // VERIFIED
	"java.lang.EnumConstantNotPresentException", // VERIFIED
	"java.lang.EventException",
	"java.io.FileNotFoundException",
	"java.lang.FileSystemAlreadyExistsException",
	"java.lang.FileSystemNotFoundException",
	"java.lang.FindException",
	"java.util.FormatFlagsConversionMismatchException", // VERIFIED
	"java.util.FormatterClosedException",               // VERIFIED
	"java.lang.IllegalAccessException",                 // VERIFIED
	"java.lang.IllegalArgumentException",               // VERIFIED
	"java.lang.IllegalCallerException",                 // VERIFIED
	"java.util.IllegalFormatCodePointException",        // VERIFIED
	"java.util.IllegalFormatConversionException",       // VERIFIED **
	"java.lang.IllegalMonitorStateException",           // VERIFIED
	"java.lang.IllegalPathStateException",
	"java.lang.IllegalStateException", // VERIFIED
	"java.lang.IllformedLocaleException",
	"java.lang.ImagingOpException",
	"java.lang.reflect.InaccessibleObjectException",     // VERIFIED
	"java.lang.annotaion.IncompleteAnnotationException", // VERIFIED
	"java.lang.InconsistentDebugInfoException",
	"java.lang.IndexOutOfBoundsException", // VERIFIED
	"java.lang.InstantiationException",    // VERIFIED
	"java.lang.InternalException",
	"java.lang.InvalidCodeIndexException",
	"java.lang.InvalidLineNumberException",
	"java.nio.InvalidMarkException",                     // VERIFIED
	"java.lang.module.InvalidModuleDescriptorException", // VERIFIED
	"java.lang.InvalidModuleException",
	"java.lang.InvalidRequestStateException",
	"java.lang.InvalidStackFrameException",
	"java.lang.JarSignerException",
	"java.lang.JMRuntimeException",
	"java.lang.JSException",
	"java.lang.LayerInstantiationException", // VERIFIED
	"java.lang.LSException",
	"java.lang.reflect.MalformedParameterizedTypeException", // VERIFIED
	"java.lang.reflect.MalformedParametersException",        // VERIFIED
	"java.lang.MirroredTypesException",
	"java.lang.MissingResourceException",
	"java.lang.NativeMethodException",
	"java.lang.NegativeArraySizeException", // VERIFIED
	"java.lang.NoSuchDynamicMethodException",
	"java.lang.NoSuchElementException",
	"java.lang.NoSuchMechanismException",
	"java.lang.NullPointerException",  // VERIFIED
	"java.lang.NumberFormatException", // VERIFIED
	"java.lang.ObjectCollectedException",
	"java.lang.ProfileDataException",
	"java.lang.ProviderException",
	"java.lang.ProviderNotFoundException",
	"java.lang.RangeException",
	"java.lang.RasterFormatException",
	"java.lang.RejectedExecutionException",
	"java.lang.module.ResolutionException", // VERIFIED
	"java.lang.SecurityException",          // VERIFIED
	"java.lang.SPIResolutionException",
	"java.lang.TypeNotPresentException", // VERIFIED
	"java.lang.UncheckedIOException",
	"java.lang.reflect.UndeclaredThrowableException", // VERIFIED
	"java.lang.UnknownEntityException",
	"java.lang.UnmodifiableModuleException",
	"java.lang.UnmodifiableSetException",
	"java.lang.UnsupportedOperationException",                // VERIFIED
	"java.nio.file.attribute.UserPrincipalNotFoundException", // VERIFIED
	"java.lang.VMDisconnectedException",
	"java.lang.VMMismatchException",
	"java.lang.VMOutOfMemoryException",
	"java.lang.invoke.WrongMethodTypeException", // verified
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
	"java.lang.DatatypeConfigurationException", // VERIFIED
	"java.lang.DestroyFailedException",
	"java.lang.ExecutionControlException",
	"java.lang.ExecutionException",
	"java.lang.ExpandVetoException",
	"java.lang.FontFormatException",
	"java.lang.GeneralSecurityException",
	"java.lang.GSSException",
	"java.lang.instrument.IllegalClassFormatException", // verified
	"java.lang.IllegalConnectorArgumentsException",
	"java.lang.IllegalThreadStateException", // VERIFIED
	"java.lang.IncompatibleThreadStateException",
	"java.lang.InterruptedException", // VERIFIED
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
	"java.lang.invoke.LambdaConversionException", // VERIFIED
	"java.lang.LastOwnerException",
	"java.lang.LineUnavailableException",
	"java.lang.MarshalException",
	"java.lang.MidiUnavailableException",
	"java.lang.MimeTypeParseException",
	"java.lang.NamingException",
	"java.lang.NoninvertibleTransformException",
	"java.lang.NoSuchFieldException",  // VERIFIED
	"java.lang.NoSuchMethodException", // VERIFIED
	"java.lang.NotBoundException",
	"java.lang.NotOwnerException",
	"java.lang.ParseException",
	"java.lang.ParserConfigurationException",
	"java.lang.PrinterException",
	"java.lang.PrintException",
	"java.lang.PrivilegedActionException",
	"java.lang.PropertyVetoException",
	"java.nio.ReadOnlyBufferException",       // VERIFIED
	"java.lang.ReflectiveOperationException", // VERIFIED
	"java.lang.RefreshFailedException",
	"java.lang.RuntimeException", // VERIFIED
	"java.lang.SAXException",
	"java.lang.ScriptException",
	"java.lang.ServerNotActiveException",
	"java.lang.SQLException",
	"java.lang.StringConcatException",
	"java.lang.StringIndexOutOfBoundsException", // VERIFIED
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

	// Java errors
	"java.lang.AnnotationFormatError",
	"java.lang.AssertionError",
	"java.lang.AWTError",
	"java.lang.CoderMalfunctionError",
	"java.lang.FactoryConfigurationError",
	"java.lang.InternalError", // VERIFIED
	"java.lang.IOError",
	"java.lang.LinkageError",         // VERIFIED
	"java.lang.NoClassDefFoundError", // VERIFIED
	"java.lang.NoSuchFieldError",     // VERIFIED
	"java.lang.NoSuchMethodError",    // VERIFIED
	"java.lang.OutOfMemoryError",     // VERIFIED
	"java.lang.SchemaFactoryConfigurationError",
	"java.lang.ServiceConfigurationError",
	"java.lang.ThreadDeath", // VERIFIED --verify that this is an error/exception
	"java.lang.TransformerFactoryConfigurationError",
	"java.lang.UnsatisfiedLinkError",         // VERIFIED
	"java.lang.UnsupportedClassVersionError", // VERIFIED
	"java.lang.VerifyError",                  // VERIFIED
	"java.lang.VirtualMachineError",          // VERIFIED

	// charset exceptions
	"java.lang.ChangedCharSetException",
	"java.lang.CharacterCodingException",
	"java.lang.CharConversionException",
	"java.lang.UnsupportedEncodingException",
	"java.lang.UTFDataFormatException",
}
