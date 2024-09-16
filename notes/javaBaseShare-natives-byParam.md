# in: void


### ()V
java.base/share/classes/java/lang/Shutdown.java:141:    native void beforeHalt();
java.base/share/classes/java/lang/ref/Reference.java:405:    native void clear0();
java.base/share/classes/java/lang/Thread.java:2946:    native void clearInterruptEvent();
java.base/share/classes/java/lang/StackFrameInfo.java:100:    native void expandStackFrameInfo();
java.base/share/classes/jdk/internal/misc/Unsafe.java:3430: native void fullFence();
java.base/share/classes/java/lang/Runtime.java:758:    native void gc();
java.base/share/classes/java/net/Inet4Address.java:383:    native void init();
java.base/share/classes/java/net/Inet6Address.java:932:    native void init();
java.base/share/classes/java/net/InetAddress.java:1961:    native void init();
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:198: native void init();
java.base/share/classes/java/net/NetworkInterface.java:640: native void init();
java.base/share/classes/jdk/internal/misc/VM.java:467:    native void initialize();
java.base/share/classes/java/io/FileDescriptor.java:222:    native void initIDs();
java.base/share/classes/java/io/FileInputStream.java:584:    native void initIDs();
java.base/share/classes/java/io/FileOutputStream.java:478:    native void initIDs();
java.base/share/classes/java/io/RandomAccessFile.java:1218: native void initIDs();
java.base/share/classes/java/lang/ProcessHandleImpl.java:569: native void initIDs();
java.base/share/classes/java/util/zip/Inflater.java:731:    native void initIDs();
java.base/share/classes/sun/nio/ch/Net.java:795:    native void initIDs();
java.base/share/classes/sun/nio/ch/IOUtil.java:602:    native void initIDs();
java.base/share/classes/java/io/ObjectStreamClass.java:232: native void initNative();
java.base/share/classes/java/lang/ProcessHandleImpl.java:82: native void initNative();
java.base/share/classes/java/lang/Thread.java:2945:    native void interrupt0();
java.base/share/classes/java/lang/Object.java:302:    native void notify();
java.base/share/classes/java/lang/Object.java:327:    native void notifyAll();
java.base/share/classes/java/lang/VirtualThread.java:1084:    native void notifyJvmtiStart();
java.base/share/classes/java/lang/VirtualThread.java:1088:    native void notifyJvmtiEnd
java.base/share/classes/jdk/internal/vm/Continuation.java:430: native void pin();
java.base/share/classes/java/lang/Class.java:236:    native void registerNatives();
java.base/share/classes/java/lang/Thread.java:221:    native void registerNatives();
java.base/share/classes/java/lang/System.java:118:    native void registerNatives();
java.base/share/classes/java/lang/VirtualThread.java:1102:    native void registerNatives();
java.base/share/classes/jdk/internal/perf/Perf.java:432:    native void registerNatives();
java.base/share/classes/jdk/internal/foreign/abi/NativeEntryPoint.java:108: native void registerNatives();
java.base/share/classes/jdk/internal/misc/Unsafe.java:58:    native void registerNatives();
java.base/share/classes/jdk/internal/vm/Continuation.java:485: native void registerNatives();
java.base/share/classes/jdk/internal/foreign/abi/UpcallLinker.java:217: native void registerNatives();
java.base/share/classes/jdk/internal/foreign/abi/UpcallStubs.java:48: native void registerNatives();
java.base/share/classes/java/lang/Thread.java:1550:    native void start0();
java.base/share/classes/jdk/internal/vm/Continuation.java:437: native void unpin()
java.base/share/classes/java/lang/ref/Reference.java:226:    native void waitForReferencePendingList();
java.base/share/classes/jdk/internal/misc/Unsafe.java:1032: native void writebackPostSync0();
java.base/share/classes/jdk/internal/misc/Unsafe.java:1026: native void writebackPreSync0();
java.base/share/classes/java/lang/Thread.java:449:    native void yield0();

#### ()V + throws exception
java.base/share/classes/java/io/FileDescriptor.java:311:    native void close0() throws IOException
java.base/share/classes/java/io/FileDescriptor.java:219:    native void sync0() throws SyncFailedException;


### ()Z
java.base/share/classes/sun/nio/ch/Net.java:521:    native boolean canIPv6SocketJoinIPv4Group0();
java.base/share/classes/sun/nio/ch/Net.java:523:    native boolean canJoin6WithIPv4Group0();
java.base/share/classes/sun/nio/ch/Net.java:525:    native boolean canUseIPv6OptionsWithIPv4LocalAddress0();
java.base/share/classes/java/lang/StackStreamFactory.java:994: native boolean checkStackWalkModes();
java.base/share/classes/java/lang/ref/Reference.java:221:    native boolean hasReferencePendingList();
java.base/share/classes/sun/net/spi/DefaultProxySelector.java:354: native boolean init();
java.base/share/classes/sun/nio/ch/UnixDomainSockets.java:173: native boolean init();
java.base/share/classes/java/lang/Class.java:855:    native boolean isArray();
java.base/share/classes/java/lang/StringUTF16.java:1503:    native boolean isBigEndian();
java.base/share/classes/jdk/internal/misc/CDS.java:75:    native boolean isDumpingArchive0();
java.base/share/classes/jdk/internal/misc/CDS.java:74:    native boolean isDumpingClassList0();
java.base/share/classes/java/lang/ref/Finalizer.java:66:    native boolean isFinalizationEnabled();
java.base/share/classes/java/lang/Class.java:4695:    native boolean isHidden();
java.base/share/classes/java/lang/Class.java:844:    native boolean isInterface();
java.base/share/classes/java/net/InetAddress.java:442:    native boolean isIPv4Available();
java.base/share/classes/sun/nio/ch/Net.java:510:    native boolean isIPv6Available0();
java.base/share/classes/java/net/InetAddress.java:445:    native boolean isIPv6Supported();
java.base/share/classes/java/lang/Class.java:897:    native boolean isPrimitive();
java.base/share/classes/jdk/internal/misc/PreviewFeatures.java:54: native boolean isPreviewEnabled();
java.base/share/classes/java/lang/Class.java:3873:    native boolean isRecord0();
java.base/share/classes/sun/nio/ch/Net.java:512:    native boolean isReusePortAvailable0();
java.base/share/classes/jdk/internal/misc/CDS.java:76:    native boolean isSharingEnabled0();
java.base/share/classes/jdk/internal/vm/ContinuationSupport.java:53: native boolean isSupported0();
java.base/share/classes/jdk/internal/vm/ForeignLinkerSupport.java:43: native boolean isSupported0();
java.base/share/classes/java/io/Console.java:439:    native boolean istty();
java.base/share/classes/sun/nio/ch/Net.java:519:    native boolean shouldSetBothIPv4AndIPv6Options0();
java.base/share/classes/java/util/concurrent/atomic/AtomicLong.java:70: native boolean VMSupportsCS8();

### ()I
java.base/share/classes/sun/nio/ch/NativeSocketAddress.java:330: native int AFINET();
java.base/share/classes/sun/nio/ch/NativeSocketAddress.java:331: native int AFINET6();
java.base/share/classes/java/lang/Runtime.java:696:    native int availableProcessors();
java.base/share/classes/sun/net/sdp/SdpSupport.java:71:    native int create0() throws IOException;
java.base/share/classes/jdk/internal/vm/Continuation.java:299: native int doYield();
java.base/share/classes/sun/nio/ch/IOUtil.java:596:    native int fdLimit();
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:210: native int ffi_default_abi();
java.base/share/classes/java/lang/Class.java:4824:    native int getClassAccessFlagsRaw0();
java.base/share/classes/java/lang/Class.java:4810:    native int getClassFileVersion0();
java.base/share/classes/java/lang/Class.java:1422:    native int getModifiers();
java.base/share/classes/java/lang/Object.java:107:    native int hashCode();
java.base/share/classes/sun/nio/ch/IOUtil.java:598:    native int iovMax();
java.base/share/classes/sun/nio/ch/Net.java:517:    native int isExclusiveBindAvailable();
java.base/share/classes/jdk/internal/vm/vector/VectorSupport.java:737: native int registerNatives();
java.base/share/classes/sun/nio/ch/NativeSocketAddress.java:335: native int offsetFamily();
java.base/share/classes/sun/nio/ch/NativeSocketAddress.java:337: native int offsetSin4Addr();
java.base/share/classes/sun/nio/ch/NativeSocketAddress.java:339: native int offsetSin6Addr();
java.base/share/classes/sun/nio/ch/NativeSocketAddress.java:336: native int offsetSin4Port();
java.base/share/classes/sun/nio/ch/NativeSocketAddress.java:338: native int offsetSin6Port();
java.base/share/classes/sun/nio/ch/NativeSocketAddress.java:340: native int offsetSin6ScopeId();
java.base/share/classes/sun/nio/ch/NativeSocketAddress.java:341: native int offsetSin6FlowInfo();
java.base/share/classes/sun/nio/ch/NativeSocketAddress.java:334: native int sizeofFamily();
java.base/share/classes/sun/nio/ch/NativeSocketAddress.java:332: native int sizeofSockAddr4();
java.base/share/classes/sun/nio/ch/NativeSocketAddress.java:333: native int sizeofSockAddr6();
java.base/share/classes/sun/nio/ch/UnixDomainSockets.java:175: native int socket0() throws IOException;

#### ()I + throws exception
java.base/share/classes/java/io/FileInputStream.java:479:    native int available0() throws IOException;
java.base/share/classes/java/io/FileInputStream.java:237:    native int read0() throws IOException;
java.base/share/classes/java/io/RandomAccessFile.java:386:    native int read0() throws IOException;
java.base/share/classes/sun/nio/ch/UnixDomainSockets.java:175: native int socket0() throws IOException;

### ()S
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:211: native short ffi_type_struct();
java.base/share/classes/sun/nio/ch/Net.java:808:    native short pollinValue();
java.base/share/classes/sun/nio/ch/Net.java:809:    native short polloutValue();
java.base/share/classes/sun/nio/ch/Net.java:810:    native short pollerrValue();
java.base/share/classes/sun/nio/ch/Net.java:811:    native short pollhupValue();
java.base/share/classes/sun/nio/ch/Net.java:812:    native short pollnvalValue();
java.base/share/classes/sun/nio/ch/Net.java:813:    native short pollconnValue();

### ()J
java.base/share/classes/java/lang/System.java:528:    native long currentTimeMillis();
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:213: native long ffi_type_void();
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:214: native long ffi_type_uint8();
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:215: native long ffi_type_sint8();
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:216: native long ffi_type_uint16();
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:217: native long ffi_type_sint16();
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:218: native long ffi_type_uint32();
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:219: native long ffi_type_sint32();
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:220: native long ffi_type_uint64();
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:221: native long ffi_type_sint64();
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:222: native long ffi_type_float();
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:223: native long ffi_type_double();
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:224: native long ffi_type_pointer();
java.base/share/classes/java/lang/Runtime.java:707:    native long freeMemory();
java.base/share/classes/java/lang/ProcessHandleImpl.java:315: native long getCurrentPid0();
java.base/share/classes/jdk/internal/misc/VM.java:417:    native long getegid();
java.base/share/classes/jdk/internal/misc/VM.java:405:    native long geteuid();
java.base/share/classes/java/lang/Thread.java:2950:    native long getNextThreadIdOffset();
java.base/share/classes/jdk/internal/misc/CDS.java:101:    native long getRandomSeedForDumping();
java.base/share/classes/jdk/internal/misc/VM.java:399:    native long getuid();
java.base/share/classes/jdk/internal/misc/VM.java:411:    native long getgid();
java.base/share/classes/jdk/internal/perf/Perf.java:417:    native long highResCounter();
java.base/share/classes/jdk/internal/perf/Perf.java:430:    native long highResFrequency();
java.base/share/classes/java/lang/Runtime.java:731:    native long maxMemory();
java.base/share/classes/java/lang/System.java:572:    native long nanoTime();
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:200: native long sizeofCif();
java.base/share/classes/java/lang/Runtime.java:720:    native long totalMemory();
java.base/share/classes/sun/nio/ch/IOUtil.java:600:    native long writevMax();

#### ()J + throws exception
java.base/share/classes/java/io/RandomAccessFile.java:612:    native long getFilePointer() throws IOException;
java.base/share/classes/java/io/FileInputStream.java:404:    native long length0() throws IOException;
java.base/share/classes/java/io/RandomAccessFile.java:657:    native long length0() throws IOException;
java.base/share/classes/java/io/FileInputStream.java:414:    native long position0() throws IOException;

### ()[B
java.base/share/classes/java/lang/Class.java:3490:    native byte[] getRawAnnotations();
java.base/share/classes/java/lang/Class.java:3492:    native byte[] getRawTypeAnnotations();
java.base/share/classes/java/lang/reflect/Executable.java:461: native byte[] getTypeAnnotationBytes0();
java.base/share/classes/java/lang/reflect/Field.java:1294:    native byte[] getTypeAnnotationBytes0();

### ()Ljava/lang/String;
java.base/share/classes/java/io/Console.java:369:    native String encoding();
java.base/share/classes/java/lang/NullPointerException.java:130: native String getExtendedNPEMessage();
java.base/share/classes/java/lang/Class.java:3462:    native String getGenericSignature0();
java.base/share/classes/java/lang/Class.java:2009:    native String getSimpleBinaryName0();
java.base/share/classes/java/util/TimeZone.java:649:    native String getSystemGMTOffsetID();
java.base/share/classes/jdk/internal/vm/VMSupport.java:113: native String getVMTemporaryDirectory();
java.base/share/classes/java/lang/Class.java:1002:    native String initClassName();
java.base/share/classes/java/lang/String.java:4640:    native String intern();

### ()Ljava/lang/String; + throws exception
java.base/share/classes/java/net/Inet4AddressImpl.java:37:    native String getLocalHostName() throws UnknownHostException;
java.base/share/classes/java/net/Inet6AddressImpl.java:48:    native String getLocalHostName() throws UnknownHostException;

### ()Ljava/lang/String;
java.base/share/classes/jdk/internal/misc/VM.java:462:    native String[] getRuntimeArguments();
java.base/share/classes/jdk/internal/loader/BootLoader.java:336: native String[] getSystemPackageNames();
java.base/share/classes/jdk/internal/util/SystemProps.java:320: native String[] platformProperties();
java.base/share/classes/jdk/internal/util/SystemProps.java:310: native String[] vmProperties();

### ()Ljava/lang/Object;
java.base/share/classes/java/lang/Thread.java:310:    native Object findScopedValueBindings();
java.base/share/classes/java/lang/Thread.java:2441:    native Object getStackTrace0();

### ()Ljava/lang/Object; + throws exception
java.base/share/classes/java/lang/Object.java:236:    native Object clone() throws CloneNotSupportedException;

###  ()[Ljava/lang/Object;
java.base/share/classes/java/lang/Class.java:1470:    native Object[] getSigners();
java.base/share/classes/java/lang/Class.java:1569:    native Object[] getEnclosingMethod0();
java.base/share/classes/java/lang/Thread.java:417:    native Object[] scopedValueCache();

###  ()Ljava/lang/Class;
java.base/share/classes/jdk/internal/reflect/Reflection.java:74: native Class<?> getCallerClass();
java.base/share/classes/java/lang/Object.java:67:    native Class<?> getClass();
java.base/share/classes/java/lang/SecurityManager.java:365: native Class<?>[] getClassContext();
java.base/share/classes/java/lang/Class.java:1752:    native Class<?> getDeclaringClass0();
java.base/share/classes/java/lang/Class.java:4407:    native Class<?> getNestHost0();
java.base/share/classes/java/lang/Class.java:1115:    native Class<? super T> getSuperclass();

###  ()[Ljava/lang/Class;
java.base/share/classes/java/lang/Class.java:3863:    native Class<?>[]    getDeclaredClasses0();
java.base/share/classes/java/lang/Class.java:1291:    native Class<?>[] getInterfaces0();
java.base/share/classes/java/lang/Class.java:4487:    native Class<?>[] getNestMembers0();
java.base/share/classes/java/lang/Class.java:4794:    native Class<?>[] getPermittedSubclasses0();

### ()Ljava/security/AccessControlContext;
java.base/share/classes/java/security/AccessController.java:1003: native AccessControlContext getInheritedAccessControlContext();
java.base/share/classes/java/security/AccessController.java:993: native AccessControlContext getStackAccessControlContext();

### ()Ljava/lang/ClassLoader;
java.base/share/classes/jdk/internal/misc/VM.java:382:    native ClassLoader latestUserDefinedLoader0();

### ()Ljdk/internal/reflect/ConstantPool;
java.base/share/classes/java/lang/Class.java:3497:    native ConstantPool getConstantPool();

### ()[Ljava/net/NetworkInterface;
java.base/share/classes/java/net/NetworkInterface.java:431: native NetworkInterface[] getAll()

### ()[Ljava/lang/reflect/Parameter;
java.base/share/classes/java/lang/reflect/Executable.java:460: native Parameter[] getParameters0();

### ()[Ljava/lang/reflect/RecordComponent;
java.base/share/classes/java/lang/Class.java:3872:    native RecordComponent[] getRecordComponents0();

### ()Ljava/lang/ref/Reference;
java.base/share/classes/java/lang/ref/Reference.java:216:    native Reference<?> getAndClearReferencePendingList();

### ()Ljava/security/ProtectionDomain;
java.base/share/classes/java/lang/Class.java:3240:    native java.security.ProtectionDomain getProtectionDomain0();

### ()Ljava/lang/Thread;
java.base/share/classes/java/lang/Thread.java:406:    native Thread currentThread();
java.base/share/classes/java/lang/Thread.java:399:    native Thread currentCarrierThread();

### ()[Ljava/lang/Thread;
java.base/share/classes/java/lang/Thread.java:2566:    native Thread[] getThreads();

# in: boolean (14 total methods)


### (Z)V
java.base/share/classes/java/lang/VirtualThread.java:1100:    native void notifyJvmtiHideFrames(boolean hide);
java.base/share/classes/java/lang/VirtualThread.java:1092:    native void notifyJvmtiMount(boolean hide);
java.base/share/classes/java/lang/VirtualThread.java:1096:    native void notifyJvmtiUnmount(boolean hide);

### (Z)J
java.base/share/classes/java/util/zip/Inflater.java:732:    native long init(boolean nowrap);

### (Z)J + throws IOException
java.base/share/classes/sun/nio/ch/IOUtil.java:573:    native long makePipe(boolean blocking) throws IOException;

### (Z)Ljava/lang/reflect/Constructor;
java.base/share/classes/java/lang/Class.java:3862:    native Constructor<T>[] getDeclaredConstructors0(boolean publicOnly);

### (Z)Ljava/lang/reflect/Field;
java.base/share/classes/java/lang/Class.java:3860:    native Field[]       getDeclaredFields0(boolean publicOnly);

### (Z)Ljava/lang/reflect/Method;
java.base/share/classes/java/lang/Class.java:3861:    native Method[]      getDeclaredMethods0(boolean publicOnly);


# in: boolean, miscellaneous


### (ZZZZ)I
java.base/share/classes/sun/nio/ch/Net.java:548:    native int socket0(boolean preferIPv6, boolean stream, boolean reuse, boolean fastLoopback);

### (ZLjava/io/FileDescriptor;III)
java.base/share/classes/sun/nio/ch/Net.java:720:    native int joinOrDrop4(boolean join, FileDescriptor fd, int group, int interf, int source) throws IOException;

### (ZLjava/io/FileDescriptor;III) + throws IOException
java.base/share/classes/sun/nio/ch/Net.java:741:    native int blockOrUnblock4(boolean block, FileDescriptor fd, int group, int interf, int source) throws IOException;
java.base/share/classes/sun/nio/ch/Net.java:763:    native int joinOrDrop6(boolean join, FileDescriptor fd, byte[] group, int index, byte[] source) throws IOException;
java.base/share/classes/sun/nio/ch/Net.java:784:    native int blockOrUnblock6(boolean block, FileDescriptor fd, byte[] group, int index, byte[] source) throws IOException;

### (ZJ)V
java.base/share/classes/jdk/internal/misc/Unsafe.java:2405: native void park(boolean isAbsolute, long time);


# in: int


### (I)V
java.base/share/classes/java/lang/Shutdown.java:153:    native void halt0(int status);
java.base/share/classes/jdk/internal/misc/Signal.java:237:    native void raise0(int sig);
java.base/share/classes/java/lang/Thread.java:2944:    native void setPriority0(int newPriority);

### (I)V + throws IOException
java.base/share/classes/sun/net/sdp/SdpSupport.java:73:    native void convert0(int fd) throws IOException;
java.base/share/classes/java/io/RandomAccessFile.java:557:    native void write0(int b) throws IOException;

### (I)Z
java.base/share/classes/java/io/FileDescriptor.java:232:    native boolean getAppend(int fd);

### (I)Z + throws IOException
java.base/share/classes/sun/nio/ch/IOUtil.java:580:    native boolean drain(int fd) throws IOException;

### (I)I + throws IOException
java.base/share/classes/sun/nio/ch/IOUtil.java:586:    native int drain1(int fd) throws IOException;

### (I)J
java.base/share/classes/java/io/FileDescriptor.java:227:    native long getHandle(int d);

### (I)F
java.base/share/classes/java/lang/Float.java:978:    native float intBitsToFloat(int bits);

### (I)Ljava/nio/ByteBuffer;
java.base/share/classes/jdk/internal/perf/Perf.java:244:    native ByteBuffer attach0(int lvmid) throws IOException;

### (I)Ljava/net/NetworkInterface;
java.base/share/classes/java/net/NetworkInterface.java:437: native NetworkInterface getByIndex0(int index)

### (I)Ljava/lang/Throwable;
java.base/share/classes/java/lang/Throwable.java:826:    native Throwable fillInStackTrace(int dummy);

### (II)I
java.base/share/classes/java/util/zip/Adler32.java:136:    native int update(int adler, int b);
java.base/share/classes/java/util/zip/CRC32.java:136:    native int update(int crc, int b);

### (IJ)V
java.base/share/classes/java/io/FileCleanable.java:56:    native void cleanupClose0(int fd, long handle) throws IOException;

### (IJ)J
java.base/share/classes/jdk/internal/misc/Signal.java:235:    native long handle0(int sig, long nativeH);

### (IZ)V + throws IOException
java.base/share/classes/java/io/FileOutputStream.java:302:    native void write(int b, boolean append) throws IOException;

### (IB)I + throws IOException
java.base/share/classes/sun/nio/ch/IOUtil.java:575:    native int write1(int fd, byte b) throws IOException;

### (I[BII)I
java.base/share/classes/java/util/zip/CRC32.java:144:    native int updateBytes0(int crc, byte[] b, int off, int len);
java.base/share/classes/java/util/zip/Adler32.java:139:    native int updateBytes(int adler, byte[] b, int off, int len);

### (IILjdk/internal/vm/ContinuationScope;Ljdk/internal/vm/Continuation;II[*?T?*)*?R?*
java.base/share/classes/java/lang/StackStreamFactory.java:457: native R callStackWalk(int mode, int skipframes, ContinuationScope contScope, Continuation continuation, int batchSize, int startIndex, T[] frames);

### (IIZ)J
java.base/share/classes/java/util/zip/Deflater.java:915:    native long init(int level, int strategy, boolean nowrap);

### (IJII)I
java.base/share/classes/java/util/zip/Adler32.java:142:    native int updateByteBuffer(int adler, long addr, int off, int len);
java.base/share/classes/java/util/zip/CRC32.java:163:    native int updateByteBuffer0(int alder, long addr, int off, int len);

### (IJJ)I
java.base/share/classes/jdk/internal/foreign/abi/fallback/LibFallback.java:208: native int ffi_get_struct_offsets(int abi, long type, long offsets);

### (IJIIIL*?T?*)I
java.base/share/classes/java/lang/StackStreamFactory.java:474: native int fetchStackFrames(int mode, long anchor, int lastBatchFrameCount, int batchSize, int startIndex, T[] frames);

### (J[*?T?*Ljdk/internal/vm/Continuation;)V
java.base/share/classes/java/lang/StackStreamFactory.java:478: native void setContinuation(long anchor, T[] frames, Continuation cont);

# in: long

### (J)V
java.base/share/classes/java/util/zip/Deflater.java:937:    native void end(long addr);
java.base/share/classes/java/util/zip/Inflater.java:750:    native void end(long addr);
java.base/share/classes/jdk/internal/misc/Unsafe.java:3825: native void freeMemory0(long address);
java.base/share/classes/java/lang/ProcessHandleImpl.java:576: native void info0(long pid);
java.base/share/classes/java/util/zip/Deflater.java:936:    native void reset(long addr);
java.base/share/classes/java/util/zip/Inflater.java:749:    native void reset(long addr);
java.base/share/classes/jdk/internal/misc/Unsafe.java:1020: native void writeback0(long address)

### (J)V + throws IOException
java.base/share/classes/java/io/RandomAccessFile.java:640:    native void seek0(long pos) throws IOException;
java.base/share/classes/java/io/RandomAccessFile.java:687:    native void setLength0(long newLength) throws IOException;

### (J)V + throws InterruptedException
java.base/share/classes/java/lang/Thread.java:498:    native void sleepNanos0(long nanos) throws InterruptedException;
java.base/share/classes/java/lang/Object.java:387:    native void wait0(long timeoutMillis) throws InterruptedException;

### (J)Z
java.base/share/classes/jdk/internal/foreign/abi/NativeEntryPoint.java:97: native boolean freeDowncallStub0(long downcallStub);
java.base/share/classes/jdk/internal/foreign/abi/UpcallStubs.java:46: native boolean freeUpcallStub0(long addr);

### (J)I
java.base/share/classes/java/util/zip/Deflater.java:935:    native int getAdler(long addr);
java.base/share/classes/java/util/zip/Inflater.java:748:    native int getAdler(long addr);

### (J)J
java.base/share/classes/jdk/internal/misc/Unsafe.java:3823: native long allocateMemory0(long bytes);
java.base/share/classes/jdk/internal/misc/VM.java:447:    native long getNanoTimeAdjustment(long offsetInSeconds);
java.base/share/classes/java/lang/ProcessHandleImpl.java:423: native long isAlive0(long pid);

### (J)J + throws IOException
java.base/share/classes/java/io/FileInputStream.java:450:    native long skip0(long n) throws IOException;

### (J)D
java.base/share/classes/java/lang/Double.java:1213:    native double longBitsToDouble(long bits);

### (J)Ljava/lang/Object;
java.base/share/classes/jdk/internal/misc/Unsafe.java:309:    native Object getUncompressedObject(long address);


# in: Object


### (Ljava/lang/Object;)V
java.base/share/classes/java/lang/Thread.java:423:    native void ensureMaterializedForStackWalk(Object o);
java.base/share/classes/java/security/AccessController.java:744: native void ensureMaterializedForStackWalk(Object o);
java.base/share/classes/java/lang/ref/Finalizer.java:107:    native void reportComplete(Object finalizee);
java.base/share/classes/jdk/internal/misc/Unsafe.java:2391: native void unpark(Object thread);

### (Ljava/lang/Object;)Z
java.base/share/classes/java/lang/Thread.java:2362:    native boolean holdsLock(Object obj);
java.base/share/classes/java/lang/Class.java:804:    native boolean isInstance(Object obj);
java.base/share/classes/java/lang/ref/PhantomReference.java:78: native boolean refersTo0(Object o);
java.base/share/classes/java/lang/ref/Reference.java:388:    native boolean refersTo0(Object o);

### (Ljava/lang/Object;)I
java.base/share/classes/java/lang/reflect/Array.java:124:    native int getLength(Object array)
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:114: native int      getSize0 (Object constantPoolOop);
java.base/share/classes/java/lang/System.java:685:    native int identityHashCode(Object x);


# in: Object, int


### (Ljava/lang/Object;I)I
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:125: native int      getIntAt0 (Object constantPoolOop, int index);
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:117: native int      getClassRefIndexAt0 (Object constantPoolOop, int index);
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:123: native int      getNameAndTypeRefIndexAt0(Object constantPoolOop, int index);
java.base/share/classes/java/lang/reflect/Array.java:238:    native int getInt(Object array, int index)

### (Ljava/lang/Object;I)J
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:126: native long     getLongAt0 (Object constantPoolOop, int index);
java.base/share/classes/java/lang/reflect/Array.java:257:    native long getLong(Object array, int index)

### (Ljava/lang/Object;I)F
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:127: native float    getFloatAt0 (Object constantPoolOop, int index);
java.base/share/classes/java/lang/reflect/Array.java:276:    native float getFloat(Object array, int index)

### (Ljava/lang/Object;I)D
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:128: native double   getDoubleAt0 (Object constantPoolOop, int index);
java.base/share/classes/java/lang/reflect/Array.java:295:    native double getDouble(Object array, int index)

### (Ljava/lang/Object;I)Ljava/lang/String;
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:129: native String   getStringAt0        (Object constantPoolOop, int index);
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:130: native String   getUTF8At0          (Object constantPoolOop, int index);

### (Ljava/lang/Object;I)B
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:131: native byte     getTagAt0           (Object constantPoolOop, int index);
java.base/share/classes/java/lang/reflect/Array.java:181:    native byte getByte(Object array, int index)

### (Ljava/lang/Object;I)Ljava/lang/Class;
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:115: native Class<?> getClassAt0         (Object constantPoolOop, int index);
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:116: native Class<?> getClassAtIfLoaded0 (Object constantPoolOop, int index);

### (Ljava/lang/Object;I)Ljava/lang/reflect/Member;
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:118: native Member   getMethodAt0        (Object constantPoolOop, int index);
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:119: native Member   getMethodAtIfLoaded0(Object constantPoolOop, int index);

### (Ljava/lang/Object;I)Ljava/lang/reflect/Field;
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:120: native Field    getFieldAt0         (Object constantPoolOop, int index);
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:121: native Field    getFieldAtIfLoaded0 (Object constantPoolOop, int index);

### (Ljava/lang/Object;I)[Ljava/lang/String;
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:122: native String[] getMemberRefInfoAt0 (Object constantPoolOop, int index);
java.base/share/classes/jdk/internal/reflect/ConstantPool.java:124: native String[] getNameAndTypeRefInfoAt0(Object constantPoolOop, int index);

### (Ljava/lang/Object;I)Z
java.base/share/classes/java/lang/reflect/Array.java:162:    native boolean getBoolean(Object array, int index)

### (Ljava/lang/Object;I)C
java.base/share/classes/java/lang/reflect/Array.java:200:    native char getChar(Object array, int index)

### (Ljava/lang/Object;I)S
java.base/share/classes/java/lang/reflect/Array.java:219:    native short getShort(Object array, int index)

### (Ljava/lang/Object;I)V
java.base/share/classes/java/lang/System.java:667:    native void arraycopy(Object src,  int  srcPos,


### (Ljava/lang/Object;I)Ljava/lang/Object;
java.base/share/classes/java/lang/reflect/Array.java:143:    native Object get(Object array, int index)


# in: Object, int, miscellaneous --> void


### (Ljava/lang/Object;ILjava/lang/Object;)V
java.base/share/classes/java/lang/reflect/Array.java:315:    native void set(Object array, int index, Object value)

### (Ljava/lang/Object;IZ)V
java.base/share/classes/java/lang/reflect/Array.java:335:    native void setBoolean(Object array, int index, boolean z)

### (Ljava/lang/Object;IB)V
java.base/share/classes/java/lang/reflect/Array.java:355:    native void setByte(Object array, int index, byte b)

### (Ljava/lang/Object;IC)V
java.base/share/classes/java/lang/reflect/Array.java:375:    native void setChar(Object array, int index, char c)

### (Ljava/lang/Object;IS)V
java.base/share/classes/java/lang/reflect/Array.java:395:    native void setShort(Object array, int index, short s)

### (Ljava/lang/Object;II)V
java.base/share/classes/java/lang/reflect/Array.java:415:    native void setInt(Object array, int index, int i)

### (Ljava/lang/Object;IJ)V
java.base/share/classes/java/lang/reflect/Array.java:435:    native void setLong(Object array, int index, long l)

### (Ljava/lang/Object;IF)V
java.base/share/classes/java/lang/reflect/Array.java:455:    native void setFloat(Object array, int index, float templateFunction)

### (Ljava/lang/Object;ID)V
java.base/share/classes/java/lang/reflect/Array.java:475:    native void setDouble(Object array, int index, double d)


# in: object, long


### (Ljava/lang/Object;J)Z
java.base/share/classes/jdk/internal/misc/Unsafe.java:1411: native boolean compareAndSetReference(Object o, long offset,
java.base/share/classes/jdk/internal/misc/Unsafe.java:1472: native boolean compareAndSetInt(Object o, long offset,
java.base/share/classes/jdk/internal/misc/Unsafe.java:202:    native boolean getBoolean(Object o, long offset);
java.base/share/classes/jdk/internal/misc/Unsafe.java:2094: native boolean getBooleanVolatile(Object o, long offset);

### (Ljava/lang/Object;J)Ljava/lang/Object;
java.base/share/classes/jdk/internal/misc/Unsafe.java:1416: native Object compareAndExchangeReference(Object o, long offset,
java.base/share/classes/jdk/internal/misc/Unsafe.java:185:    native Object getReference(Object o, long offset);
java.base/share/classes/jdk/internal/misc/Unsafe.java:2075: native Object getReferenceVolatile(Object o, long offset);

### (Ljava/lang/Object;J)I
java.base/share/classes/jdk/internal/misc/Unsafe.java:1477: native int compareAndExchangeInt(Object o, long offset,
java.base/share/classes/jdk/internal/misc/Unsafe.java:155:    native int getInt(Object o, long offset);
java.base/share/classes/jdk/internal/misc/Unsafe.java:2086: native int     getIntVolatile(Object o, long offset);

### (Ljava/lang/Object;J)B
java.base/share/classes/jdk/internal/misc/Unsafe.java:2102: native byte    getByteVolatile(Object o, long offset);
java.base/share/classes/jdk/internal/misc/Unsafe.java:210:    native byte    getByte(Object o, long offset);

### (Ljava/lang/Object;J)S
java.base/share/classes/jdk/internal/misc/Unsafe.java:2110: native short   getShortVolatile(Object o, long offset);
java.base/share/classes/jdk/internal/misc/Unsafe.java:218:    native short   getShort(Object o, long offset);

### (Ljava/lang/Object;J)C
java.base/share/classes/jdk/internal/misc/Unsafe.java:2118: native char    getCharVolatile(Object o, long offset);
java.base/share/classes/jdk/internal/misc/Unsafe.java:226:    native char    getChar(Object o, long offset);

### (Ljava/lang/Object;J)J
java.base/share/classes/jdk/internal/misc/Unsafe.java:2126: native long    getLongVolatile(Object o, long offset);
java.base/share/classes/jdk/internal/misc/Unsafe.java:234:    native long    getLong(Object o, long offset);

### (Ljava/lang/Object;J)F
java.base/share/classes/jdk/internal/misc/Unsafe.java:2134: native float   getFloatVolatile(Object o, long offset);
java.base/share/classes/jdk/internal/misc/Unsafe.java:242:    native float   getFloat(Object o, long offset);

### (Ljava/lang/Object;J)D
java.base/share/classes/jdk/internal/misc/Unsafe.java:2142: native double  getDoubleVolatile(Object o, long offset);
java.base/share/classes/jdk/internal/misc/Unsafe.java:250:    native double  getDouble(Object o, long offset);


# in: Object, long, miscellaneous


### (Ljava/lang/Object;JI)V
java.base/share/classes/jdk/internal/misc/Unsafe.java:178:    native void putInt(Object o, long offset, int x);
java.base/share/classes/jdk/internal/misc/Unsafe.java:2090: native void    putIntVolatile(Object o, long offset, int x);

### (Ljava/lang/Object;JLjava/lang/Object;)V
java.base/share/classes/jdk/internal/misc/Unsafe.java:198:    native void putReference(Object o, long offset, Object x);
java.base/share/classes/jdk/internal/misc/Unsafe.java:2082: native void putReferenceVolatile(Object o, long offset, Object x);

### (Ljava/lang/Object;JJJ)Z
java.base/share/classes/jdk/internal/misc/Unsafe.java:2019: native boolean compareAndSetLong(Object o, long offset, long expected, long x);
                                                  
### (Ljava/lang/Object;JJJ)J
java.base/share/classes/jdk/internal/misc/Unsafe.java:2024: native long compareAndExchangeLong(Object o, long offset, long expected, long x);

### (Ljava/lang/Object;JZ)V
java.base/share/classes/jdk/internal/misc/Unsafe.java:206:    native void    putBoolean(Object o, long offset, boolean x);
java.base/share/classes/jdk/internal/misc/Unsafe.java:2098: native void    putBooleanVolatile(Object o, long offset, boolean x);

### (Ljava/lang/Object;JB)V
java.base/share/classes/jdk/internal/misc/Unsafe.java:2106: native void    putByteVolatile(Object o, long offset, byte x);
java.base/share/classes/jdk/internal/misc/Unsafe.java:214:    native void    putByte(Object o, long offset, byte x);

### (Ljava/lang/Object;JS)V
java.base/share/classes/jdk/internal/misc/Unsafe.java:2114: native void    putShortVolatile(Object o, long offset, short x);
java.base/share/classes/jdk/internal/misc/Unsafe.java:222:    native void    putShort(Object o, long offset, short x);

### (Ljava/lang/Object;JC)V
java.base/share/classes/jdk/internal/misc/Unsafe.java:2122: native void    putCharVolatile(Object o, long offset, char x);
java.base/share/classes/jdk/internal/misc/Unsafe.java:230:    native void    putChar(Object o, long offset, char x);

### (Ljava/lang/Object;JJ)V
java.base/share/classes/jdk/internal/misc/Unsafe.java:2130: native void    putLongVolatile(Object o, long offset, long x);
java.base/share/classes/jdk/internal/misc/Unsafe.java:238:    native void    putLong(Object o, long offset, long x);

### (Ljava/lang/Object;JF)V
java.base/share/classes/jdk/internal/misc/Unsafe.java:2138: native void    putFloatVolatile(Object o, long offset, float x);
java.base/share/classes/jdk/internal/misc/Unsafe.java:246:    native void    putFloat(Object o, long offset, float x);

### (Ljava/lang/Object;JD)V
java.base/share/classes/jdk/internal/misc/Unsafe.java:2146: native void    putDoubleVolatile(Object o, long offset, double x);
java.base/share/classes/jdk/internal/misc/Unsafe.java:254:    native void    putDouble(Object o, long offset, double x);


### (Ljava/lang/Object;JJB)V
java.base/share/classes/jdk/internal/misc/Unsafe.java:3826: native void setMemory0(Object o, long offset, long bytes, byte value);

### (Ljava/lang/Object;JLjava/lang/Object;JJ)V
java.base/share/classes/jdk/internal/misc/Unsafe.java:3828: native void copyMemory0(Object srcBase, long srcOffset, Object destBase, long destOffset, long bytes);

### (Ljava/lang/Object;JLjava/lang/Object;JJJ)V
java.base/share/classes/jdk/internal/misc/Unsafe.java:3829: native void copySwapMemory0(Object srcBase, long srcOffset, Object destBase, long destOffset, long bytes, long elemSize);


# in: Object[]


### ([Ljava/lang/Object;)V
java.base/share/classes/java/lang/Thread.java:420:    native void setScopedValueCache(Object[] cache);
java.base/share/classes/java/lang/Class.java:1476:    native void setSigners(Object[] signers);


# in: String


### (Ljava/jang/String;)V
java.base/share/classes/jdk/internal/misc/CDS.java:207:    native void dumpClassList(String listFileName);
java.base/share/classes/jdk/internal/misc/CDS.java:208:    native void dumpDynamicArchive(String archiveFileName);
java.base/share/classes/jdk/internal/misc/CDS.java:77:    native void logLambdaFormInvoker(String line);
java.base/share/classes/java/lang/Thread.java:2947:    native void setNativeName(String name);

### (Ljava/jang/String;)V +  throws FileNotFoundException
java.base/share/classes/java/io/FileInputStream.java:203:    native void open0(String name) throws FileNotFoundException;

### (Ljava/jang/String;)I
java.base/share/classes/jdk/internal/misc/Signal.java:227:    native int findSignal0(String sigName);

### (Ljava/jang/String;)Ljava/jang/String;
java.base/share/classes/jdk/internal/loader/NativeLibraries.java:550: native String findBuiltinLib(String name);
java.base/share/classes/jdk/internal/loader/BootLoader.java:346: native String getSystemPackageLocation(String name);
java.base/share/classes/java/util/TimeZone.java:643:    native String getSystemTimeZoneID(String javaHome);
java.base/share/classes/java/lang/System.java:2077:    native String mapLibraryName(String libname);

### (Ljava/jang/String;)Ljava/nio/ByteBuffer;
java.base/share/classes/jdk/internal/jimage/NativeImageBuffer.java:48: native ByteBuffer getNativeMap(String imagePath);

### (Ljava/jang/String;)[Ljava/net/InetAddress;
java.base/share/classes/java/net/Inet4AddressImpl.java:45:    native InetAddress[] lookupAllHostAddr(String hostname) throws UnknownHostException;

### (Ljava/jang/String;)Ljava/net/NetworkInterface;
java.base/share/classes/java/net/NetworkInterface.java:434: native NetworkInterface getByName0(String name)

### (Ljava/jang/String;)Ljava/lang/Class;
java.base/share/classes/java/lang/Class.java:3246:    native Class<?> getPrimitiveClass(String name);


# in: String, int


### (Ljava/jang/String;I)Z + throws SocketException
java.base/share/classes/java/net/NetworkInterface.java:563: native boolean isUp0(String name, int ind) throws SocketException;
java.base/share/classes/java/net/NetworkInterface.java:564: native boolean isLoopback0(String name, int ind) throws SocketException;
java.base/share/classes/java/net/NetworkInterface.java:565: native boolean supportsMulticast0(String name, int ind) throws SocketException;
java.base/share/classes/java/net/NetworkInterface.java:566: native boolean isP2P0(String name, int ind) throws SocketException;

### (Ljava/jang/String;I)I + throws SocketException
java.base/share/classes/java/net/NetworkInterface.java:568: native int getMTU0(String name, int ind) throws SocketException;

### (Ljava/jang/String;I)[Ljava/net/InetAddress; + throws UnknownHostException
java.base/share/classes/java/net/Inet6AddressImpl.java:55:    native InetAddress[] lookupAllHostAddr(String hostname, int characteristics) throws UnknownHostException;

### (Ljava/jang/String;I)V
java.base/share/classes/java/io/RandomAccessFile.java:336:    native void open0(String name, int mode)


# in: String, miscellaneous


### (Ljava/jang/String;ZLjava/lang/ClassLoader;Ljava/lang/Class;)Ljava/lang/Class; + throws ClassNotFoundException
java.base/share/classes/java/lang/Class.java:538:    native Class<?> forName0(String name, boolean initialize, ClassLoader loader, Class<?> caller) throws ClassNotFoundException;

### (Ljava/jang/String;Z)V + throws FileNotFoundException
java.base/share/classes/java/io/FileOutputStream.java:277:    native void open0(String name, boolean append) throws FileNotFoundException;

### (Ljava/jang/String;ZJ)V
java.base/share/classes/jdk/internal/loader/NativeLibraries.java:549: native void unload(String name, boolean isBuiltin, long handle);

### (Ljava/jang/String;J)V
java.base/share/classes/jdk/internal/loader/RawNativeLibraries.java:192: native void unload0(String name, long handle);

### (Ljava/jang/String;[BIILjava/lang/ClassLoader;Ljava/security/ProtectionDomain;)Ljava/jang/Class;
java.base/share/classes/jdk/internal/misc/Unsafe.java:1333: native Class<?> defineClass0(String name, byte[] b, int off, int len, ClassLoader loader, ProtectionDomain protectionDomain);

### (Ljava/jang/String;IIJ)Ljava/nio/ByteBuffer;
java.base/share/classes/jdk/internal/perf/Perf.java:289:    native ByteBuffer createLong(String name, int variability, int units, long value);

### (Ljava/jang/String;II[BI)Ljava/nio/ByteBuffer;
java.base/share/classes/jdk/internal/perf/Perf.java:399:    native ByteBuffer createByteArray(String name, int variability, int units, byte[] value, int maxLength);


# in: byte[]


### ([BLjava/jang/String;I)[B + throws SocketException
java.base/share/classes/java/net/NetworkInterface.java:567: native byte[] getMacAddr0(byte[] inAddr, String name, int ind) throws SocketException;

