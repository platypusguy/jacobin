/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

// TODO: The entire SecurityManager class is deprecated and will be removed in the future.

var classNameSecurityManager = "java/lang/SecurityManager"

func Load_Lang_SecurityManager() {
	MethodSignatures[classNameSecurityManager+".<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkAccept(Ljava/lang/String;I)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkAccess(Ljava/lang/Thread;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkAccess(Ljava/lang/ThreadGroup;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkAwtEventQueueAccess()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkConnect(Ljava/lang/String;I)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkConnect(Ljava/lang/String;ILjava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkCreateClassLoader()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkDelete(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkExec(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkExit(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkLink(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkListen(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkMemberAccess(Ljava/lang/Class;I)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkMulticast(Ljava/net/InetAddress;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkMulticast(Ljava/net/InetAddress;B)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkPackageAccess(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkPackageDefinition(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkPermission(Ljava/security/Permission;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkPermission(Ljava/security/Permission;Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkPrintJobAccess()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkPropertiesAccess()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkPropertyAccess(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkRead(Ljava/io/FileDescriptor;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkRead(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkRead(Ljava/lang/String;Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkSecurityAccess(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkSetFactory()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkSystemClipboardAccess()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkTopLevelWindow(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapDeprecated,
		}

	MethodSignatures[classNameSecurityManager+".checkWrite(Ljava/io/FileDescriptor;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".checkWrite(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures[classNameSecurityManager+".classDepth(Ljava/lang/String;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapDeprecated,
		}

	MethodSignatures[classNameSecurityManager+".classLoaderDepth()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures[classNameSecurityManager+".currentClassLoader()Ljava/lang/ClassLoader;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures[classNameSecurityManager+".currentLoadedClass()Ljava/lang/Class;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures[classNameSecurityManager+".getClassContext()[Ljava/lang/Class;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures[classNameSecurityManager+".getInCheck()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures[classNameSecurityManager+".getSecurityContext()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures[classNameSecurityManager+".getThreadGroup()Ljava/lang/ThreadGroup;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}
}
