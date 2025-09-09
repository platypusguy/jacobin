/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

func Load_Traps_Java_Nio() {

	MethodSignatures["java/nio/ByteBuffer.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/nio/file/AccessMode.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/nio/file/AccessMode.valueOf(Ljava/lang/String;)Ljava/nio/AccessMode;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/nio/file/Files.copy(Ljava/io/InputStream;Ljava/nio/file/Path;[Ljava/nio/file/CopyOption;)J"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.copy(Ljava/nio/file/Path;Ljava/io/OutputStream;)J"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.copy(Ljava/nio/file/Path;Ljava/nio/file/Path;[Ljava/nio/file/CopyOption;)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.createDirectories(Ljava/nio/file/Path;[Ljava/nio/file/attribute/FileAttribute;)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.createDirectory(Ljava/nio/file/Path;[Ljava/nio/file/attribute/FileAttribute;)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.createFile(Ljava/nio/file/Path;[Ljava/nio/file/attribute/FileAttribute;)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.createLink(Ljava/nio/file/Path;Ljava/nio/file/Path;)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.createSymbolicLink(Ljava/nio/file/Path;Ljava/nio/file/Path;[Ljava/nio/file/attribute/FileAttribute;)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.delete(Ljava/nio/file/Path;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.deleteIfExists(Ljava/nio/file/Path;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.exists(Ljava/nio/file/Path;[Ljava/nio/file/LinkOption;)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.isDirectory(Ljava/nio/file/Path;[Ljava/nio/file/LinkOption;)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.isExecutable(Ljava/nio/file/Path;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.isReadable(Ljava/nio/file/Path;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.isRegularFile(Ljava/nio/file/Path;[Ljava/nio/file/LinkOption;)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.isSameFile(Ljava/nio/file/Path;Ljava/nio/file/Path;)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.isSymbolicLink(Ljava/nio/file/Path;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.isWritable(Ljava/nio/file/Path;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.move(Ljava/nio/file/Path;Ljava/nio/file/Path;[Ljava/nio/file/CopyOption;)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.newBufferedReader(Ljava/nio/file/Path;)Ljava/io/BufferedReader;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.newBufferedWriter(Ljava/nio/file/Path;[Ljava/nio/file/OpenOption;)Ljava/io/BufferedWriter;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.newDirectoryStream(Ljava/nio/file/Path;)Ljava/nio/file/DirectoryStream;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.newInputStream(Ljava/nio/file/Path;[Ljava/nio/file/OpenOption;)Ljava/io/InputStream;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.newOutputStream(Ljava/nio/file/Path;[Ljava/nio/file/OpenOption;)Ljava/io/OutputStream;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.notExists(Ljava/nio/file/Path;[Ljava/nio/file/LinkOption;)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.probeContentType(Ljava/nio/file/Path;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.readAllBytes(Ljava/nio/file/Path;)[B"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Files.size(Ljava/nio/file/Path;)J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileStore.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/nio/file/FileStore.getAttribute(Ljava/lang/String;)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileStore.getBlockSize()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileStore.getFileStoreAttributeView(Ljava/lang/Class;)Ljava/nio/file/attribute/FileStoreAttributeView;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileStore.getName()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileStore.getTotalSpace()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileStore.getUnallocatedSpace()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileStore.getUsableSpace()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileStore.isReadOnly()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileStore.supportsFileAttributeView(Ljava/lang/Class;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileStore.supportsFileAttributeView(Ljava/lang/String;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileStore.type()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileSystem.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/nio/file/FileSystem.close()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileSystem.getFileStores()Ljava/lang/Iterable;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileSystem.getPath(Ljava/lang/String;[Ljava/lang/String;)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileSystem.getPathMatcher(Ljava/lang/String;)Ljava/nio/file/PathMatcher;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileSystem.getRootDirectories()Ljava/lang/Iterable;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileSystem.getSeparator()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileSystem.isOpen()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileSystem.isReadOnly()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileSystem.newWatchService()Ljava/nio/file/WatchService;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileSystem.provider()Ljava/nio/file/spi/FileSystemProvider;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileSystem.supportedFileAttributeViews()Ljava/util/Set;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileSystems.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/nio/file/FileSystems.getDefault()Ljava/nio/file/FileSystem;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileSystems.newFileSystem(Ljava/net/URI;[Ljava/lang/Map;)Ljava/nio/file/FileSystem;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileSystems.newFileSystem(Ljava/nio/file/Path;[Ljava/lang/Map;)Ljava/nio/file/FileSystem;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/FileSystems.newFileSystem(Ljava/nio/file/Path;[Ljava/lang/Map;Ljava/lang/ClassLoader;)Ljava/nio/file/FileSystem;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/nio/file/Path.compareTo(Ljava/nio/file/Path;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.endsWith(Ljava/lang/String;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.endsWith(Ljava/nio/file/Path;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.getFileName()Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.getName(I)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.getNameCount()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.getParent()Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.getRoot()Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.isAbsolute()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.iterator()Ljava/util/Iterator;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.normalize()Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.relativize(Ljava/nio/file/Path;)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.resolve(Ljava/lang/String;)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.resolve(Ljava/nio/file/Path;)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.resolveSibling(Ljava/lang/String;)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.resolveSibling(Ljava/nio/file/Path;)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.startsWith(Ljava/lang/String;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.startsWith(Ljava/nio/file/Path;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.subpath(II)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.toAbsolutePath()Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.toFile()Ljava/io/File;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Path.toUri()Ljava/net/URI;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/PathMatcher.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/nio/file/PathMatcher.matches(Ljava/nio/file/Path;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Paths.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/nio/file/Paths.get(Ljava/lang/String;[Ljava/lang/String;)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/file/Paths.get(Ljava/net/URI;)Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/charset/StandardCharsets.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/nio/charset/StandardCharsets.ISO_8859_1()Ljava/nio/charset/Charset;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/charset/StandardCharsets.US_ASCII()Ljava/nio/charset/Charset;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/charset/StandardCharsets.UTF_16()Ljava/nio/charset/Charset;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/charset/StandardCharsets.UTF_16BE()Ljava/nio/charset/Charset;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/charset/StandardCharsets.UTF_16LE()Ljava/nio/charset/Charset;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/charset/StandardCharsets.UTF_8()Ljava/nio/charset/Charset;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/AsynchronousFileChannel.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/nio/channels/AsynchronousFileChannel.close()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/AsynchronousFileChannel.force(Z)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/AsynchronousFileChannel.isOpen()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/AsynchronousFileChannel.read(Ljava/nio/ByteBuffer;J)Ljava/util/concurrent/Future;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/AsynchronousFileChannel.read(Ljava/nio/ByteBuffer;JLcompletion/CompletionHandler;)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/AsynchronousFileChannel.size()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/AsynchronousFileChannel.truncate(J)Ljava/nio/channels/AsynchronousFileChannel;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/AsynchronousFileChannel.write(Ljava/nio/ByteBuffer;J)Ljava/util/concurrent/Future;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/AsynchronousFileChannel.write(Ljava/nio/ByteBuffer;JLcompletion/CompletionHandler;)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/AsynchronousFileChannel.open(Ljava/nio/file/Path;[Ljava/nio/file/OpenOption;)Ljava/nio/channels/AsynchronousFileChannel;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/AsynchronousFileChannel.open(Ljava/nio/file/Path;[Ljava/nio/file/OpenOption;Ljava/util/concurrent/ExecutorService;)Ljava/nio/channels/AsynchronousFileChannel;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/nio/channels/FileChannel.close()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.force(Z)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.isOpen()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.lock()Ljava/nio/channels/FileLock;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.lock(JJZ)Ljava/nio/channels/FileLock;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.map(Ljava/nio/channels/FileChannel$MapMode;JJ)Ljava/nio/MappedByteBuffer;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.position()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.position(J)Ljava/nio/channels/FileChannel;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.read(Ljava/nio/ByteBuffer;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.read(Ljava/nio/ByteBuffer;J)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.read([Ljava/nio/ByteBuffer;)J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.read([Ljava/nio/ByteBuffer;II)J"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.size()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.truncate(J)Ljava/nio/channels/FileChannel;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.transferFrom(Ljava/nio/channels/ReadableByteChannel;JJ)J"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.transferTo(JJLjava/nio/channels/WritableByteChannel;)J"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.write(Ljava/nio/ByteBuffer;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.write(Ljava/nio/ByteBuffer;J)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.write([Ljava/nio/ByteBuffer;)J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.write([Ljava/nio/ByteBuffer;II)J"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.open(Ljava/nio/file/Path;[Ljava/nio/file/OpenOption;)Ljava/nio/channels/FileChannel;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/nio/channels/FileChannel.open(Ljava/nio/file/Path;[Ljava/nio/file/OpenOption;Ljava/util/concurrent/ExecutorService;)Ljava/nio/channels/FileChannel;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}
}
