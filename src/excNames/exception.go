/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
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
	EmptyStackException             // in HotSpot, used by Stack class; in Jacobin, for all stack underflows
	EnumConstantNotPresentException // typically, used in annotation processing
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
	LayerInstantiationException
	LSException
	MalformedParameterizedTypeException
	MalformedParametersException // for HotSpot reflection: param count wrong, CP index invalid, illegal flag combo
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
	WrongMethodTypeException // used here in many places; in HotSpot, it's mostly for method handles
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
	InvalidApplicationException // MBean exception in JMX, rarely shown to user
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
	BootstrapMethodError  // invokedynamic
	ClassCircularityError // for circularity in superclass hierarchy
	ClassFormatError
	CoderMalfunctionError
	ExceptionInInitializerError // for exceptions in static initalizers
	FactoryConfigurationError
	IncompatibleClassChangeError // if class has changed unexpectedly
	InternalError
	IOError
	LinkageError
	NoClassDefFoundError
	NoSuchFieldError
	NoSuchMethodError
	OutOfMemoryError
	SchemaFactoryConfigurationError
	ServiceConfigurationError
	StackOverflowError
	ThreadDeath
	TransformerFactoryConfigurationError
	UnknownError
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
	"java.lang.annotation.AnnotationTypeMismatchException",   // VERIFIED
	"java.lang.ArithmeticException",                          // VERIFIED
	"java.lang.ArrayIndexOutOfBoundsException",               // VERIFIED
	"java.lang.ArrayStoreException",                          // VERIFIED
	"java.nio.file.AtomicMoveNotSupportedException",          // VERIFIED
	"java.nio.BufferOverflowException",                       // VERIFIED
	"java.nio.BufferUnderflowException",                      // VERIFIED
	"javax.swing.undo.CannotRedoException",                   // VERIFIED
	"javax.swing.undo.CannotUndoException",                   // VERIFIED
	"javax.xml.catalog.CatalogException",                     // VERIFIED
	"java.lang.ClassCastException",                           // VERIFIED
	"java.lang.ClassNotFoundException",                       // VERIFIED
	"com.sun.jdi.ClassNotPreparedException",                  // VERIFIED
	"java.awt.color.CMMException",                            // VERIFIED
	"java.util.concurrent.CompletionException",               // VERIFIED
	"java.util.ConcurrentModificationException",              // VERIFIED
	"java.time.DateTimeException",                            // VERIFIED
	"org.w3c.dom.DOMException",                               // VERIFIED
	"java.util.DuplicateFormatFlagsException",                // VERIFIED
	"com.sun.jdi.request.DuplicateRequestException",          // VERIFIED
	"java.util.EmptyStackException",                          // VERIFIED
	"java.lang.EnumConstantNotPresentException",              // VERIFIED
	"org.w3c.dom.events.EventException",                      // VERIFIED
	"java.io.FileNotFoundException",                          // VERiFIED
	"java.nio.file.FileSystemAlreadyExistsException",         // VERIFIED
	"java.nio.file.FileSystemNotFoundException",              // VERIFIED
	"java.lang.module.FindException",                         // VERIFIED
	"java.util.FormatFlagsConversionMismatchException",       // VERIFIED
	"java.util.FormatterClosedException",                     // VERIFIED
	"java.lang.IllegalAccessException",                       // VERIFIED
	"java.lang.IllegalArgumentException",                     // VERIFIED
	"java.lang.IllegalCallerException",                       // VERIFIED
	"java.util.IllegalFormatCodePointException",              // VERIFIED
	"java.util.IllegalFormatConversionException",             // VERIFIED ** got this far in java.util
	"java.lang.IllegalMonitorStateException",                 // VERIFIED
	"java.awt.geom.IllegalPathStateException",                // VERIFIED
	"java.lang.IllegalStateException",                        // VERIFIED
	"java.util.IllformedLocaleException",                     // VERIFIED
	"java.awt.image.ImagingOpException",                      // VERIFIED
	"java.lang.reflect.InaccessibleObjectException",          // VERIFIED
	"java.lang.annotaion.IncompleteAnnotationException",      // VERIFIED
	"com.sun.jdi.InconsistentDebugInfoException",             // VERIFIED
	"java.lang.IndexOutOfBoundsException",                    // VERIFIED
	"java.lang.InstantiationException",                       // VERIFIED
	"com.sun.jdi.InternalException",                          // VERIFIED
	"com.sun.jdi.InvalidCodeIndexException",                  // VERIFIED
	"com.sun.jdi.InvalidLineNumberException",                 // VERIFIED
	"java.nio.InvalidMarkException",                          // VERIFIED
	"java.lang.module.InvalidModuleDescriptorException",      // VERIFIED
	"com.sun.jdi.InvalidModuleException",                     // VERIFIED
	"com.sun.jdi.request.InvalidRequestStateException",       // VERIFIED
	"com.sun.jdi.InvalidStackFrameException",                 // VERIFIED
	"jdk.security.jarsigner.JarSignerException",              // VERIFIED
	"jjavax.management.JMRuntimeException",                   // VERIFIED
	"java.lang.LayerInstantiationException",                  // VERIFIED
	"org.w3c.dom.ls.LSException",                             // VERIFIED
	"java.lang.reflect.MalformedParameterizedTypeException",  // VERIFIED
	"java.lang.reflect.MalformedParametersException",         // VERIFIED
	"javax.lang.model.type.MirroredTypesException",           // VERIFIED
	"java.util.MissingResourceException",                     // VERIFIED
	"com.sun.jdi.NativeMethodException",                      // VERIFIED
	"java.lang.NegativeArraySizeException",                   // VERIFIED
	"jdk.dynalink.NoSuchDynamicMethodException",              // VERIFIED
	"java.util.NoSuchElementException",                       // VERIFIED
	"javax.xml.crypto.NoSuchMechanismException",              // VERIFIED
	"java.lang.NullPointerException",                         // VERIFIED
	"java.lang.NumberFormatException",                        // VERIFIED
	"com.sun.jdi.ObjectCollectedException",                   // VERIFIED
	"java.awt.color.ProfileDataException",                    // VERIFIED
	"java.security.ProviderException",                        // VERIFIED
	"java.nio.file.ProviderNotFoundException",                // VERIFIED
	"org.w3c.dom.ranges.RangeException",                      // VERIFIED
	"java.awt.image.RasterFormatException",                   // VERIFIED
	"java.util.concurrent.RejectedExecutionException",        // VERIFIED
	"java.lang.module.ResolutionException",                   // VERIFIED
	"java.lang.SecurityException",                            // VERIFIED
	"jdk.jshell.spi.SPIResolutionException",                  // VERIFIED
	"java.lang.TypeNotPresentException",                      // VERIFIED
	"java.io.UncheckedIOException",                           // VERIFIED
	"java.lang.reflect.UndeclaredThrowableException",         // VERIFIED
	"javax.lang.model.UnknownEntityException",                // VERIFIED
	"java.lang.instrument.UnmodifiableModuleException",       // VERIFIED
	"javax.print.attribute.UnmodifiableSetException",         // VERIFIED
	"java.lang.UnsupportedOperationException",                // VERIFIED
	"java.nio.file.attribute.UserPrincipalNotFoundException", // VERIFIED
	"com.sun.jdi.VMDisconnectedException",                    // VERIFIED
	"com.sun.jdi.VMMismatchException",                        // VERIFIED
	"com.sun.jdi.VMOutOfMemoryException",                     // VERIFIED
	"java.lang.invoke.WrongMethodTypeException",              // VERIFIED
	"javax.xml.xpath.XPathException",                         // VERIFIED

	// non-runtime exceptions
	"com.sun.jdi.AbsentInformationException",                   // VERIFIED
	"java.security.acl.AclNotFoundException",                   // VERIFIED might not be part of JDK 17
	"java.rmi.activation.ActivationException",                  // VERIFIED might not be part of JDK 17
	"com.sun.tools.attach.AgentInitializationException",        // VERIFIED
	"com.sun.tools.attach.AgentLoadException",                  // VERIFIED
	"java.rmi.AlreadyBoundException",                           // VERIFIED
	"com.sun.tools.attach.AttachNotSupportedException",         // VERIFIED
	"java.awt.AWTException",                                    // VERIFIED
	"java.util.prefs.BackingStoreException",                    // VERIFIED
	"javax.management.BadAttributeValueExpException",           // VERIFIED
	"javax.management.BadBinaryOpValueExpException",            // VERIFIED
	"javax.swing.text.BadLocationException",                    // VERIFIED
	"javax.management.BadStringOperationException",             // VERIFIED
	"java.util.concurrent.BrokenBarrierException",              // VERIFIED
	"javax.smartcardio.CardException",                          // VERIFIED
	"java.security.cert.CertificateException",                  // VERIFIED
	"com.sun.jdi.ClassNotLoadedException",                      // VERIFIED
	"java.lang.CloneNotSupportedException",                     // VERIFIED
	"java.util.zip.DataFormatException",                        // VERIFIED
	"java.lang.DatatypeConfigurationException",                 // VERIFIED
	"javax.security.auth.DestroyFailedException",               // VERIFIED
	"dk.jshell.spi.ExecutionControl.ExecutionControlException", // VERIFIED
	"java.util.concurrent.ExecutionException",                  // VERIFIED
	"javax.swing.tree.ExpandVetoException",                     // VERIFIED
	"java.awt.FontFormatException",                             // VERIFIED
	"java.security.GeneralSecurityException",                   // VERIFIED
	"org.ietf.jgss.GSSException",                               // VERIFIED
	"java.lang.instrument.IllegalClassFormatException",         // VERIFIED
	"com.sun.jdi.connect.IllegalConnectorArgumentsException",   // VERIFIED
	"java.lang.IllegalThreadStateException",                    // VERIFIED
	"com.sun.jdi.IncompatibleThreadStateException",             // VERIFIED
	"java.lang.InterruptedException",                           // VERIFIED
	"javax.management.IntrospectionException",                  // VERIFIED
	"javax.management.InvalidApplicationException",             // VERIFIED
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
	"javax.print.PrintException", // VERIFIED
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
	"java.lang.instrument.UnmodifiableClassException", // VERIFIED
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
	"java.lang.annotation.AnnotationFormatError",               // VERIFIED
	"java.lang.AssertionError",                                 // VERIFIED
	"java.awt.AWTError",                                        // VERIFIED
	"java.lang.BootstrapMethodError",                           // VERIFIED
	"java.lang.ClassCircularityError",                          // VERIFIED
	"java.lang.ClassFormatError",                               // VERIFIED
	"java.nio.charset.CoderMalfunctionError",                   // VERIFIED
	"java.lang.ExceptionInInitializerError",                    // VERIFIED
	"javax.xml.parsers.FactoryConfigurationError",              // VERIFIED
	"java.lang.IncompatibleClassChangeError",                   // VERIFIED
	"java.lang.InternalError",                                  // VERIFIED
	"java.io.IOError",                                          // VERIFIED
	"java.lang.LinkageError",                                   // VERIFIED
	"java.lang.NoClassDefFoundError",                           // VERIFIED
	"java.lang.NoSuchFieldError",                               // VERIFIED
	"java.lang.NoSuchMethodError",                              // VERIFIED
	"java.lang.OutOfMemoryError",                               // VERIFIED
	"javax.xml.validation.SchemaFactoryConfigurationError",     // VERIFIED
	"java.util.ServiceConfigurationError",                      // VERIFIED
	"java.lang.StackOverflowError",                             // VERIFIED
	"java.lang.ThreadDeath",                                    // VERIFIED -- this really is a Java error
	"javax.xml.transform.TransformerFactoryConfigurationError", // VERIFIED
	"java.lang.UnknownError",                                   // VERIFIED
	"java.lang.UnsatisfiedLinkError",                           // VERIFIED
	"java.lang.UnsupportedClassVersionError",                   // VERIFIED
	"java.lang.VerifyError",                                    // VERIFIED
	"java.lang.VirtualMachineError",                            // VERIFIED

	// charset exceptions (but note java.nio.charset.CoderMalfunctionError in the error section above)
	"javax.swing.text.ChangedCharSetException",  // VERIFIED
	"java.nio.charset.CharacterCodingException", // VERIFIED
	"java.io.CharConversionException",           // VERIFIED
	"java.io.UnsupportedEncodingException",      // VERIFIED
	"java.io.UTFDataFormatException",            // VERIFIED
}
