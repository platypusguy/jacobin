/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import "jacobin/src/gfunction/ghelpers"

// TODO: The entire SecurityManager class is deprecated and will be removed in the future.

var classNameSecurityManager = "java/lang/SecurityManager"

func Load_Lang_SecurityManager() {
	ghelpers.MethodSignatures[classNameSecurityManager+".<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkAccept(Ljava/lang/String;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkAccess(Ljava/lang/Thread;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkAccess(Ljava/lang/ThreadGroup;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkAwtEventQueueAccess()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkConnect(Ljava/lang/String;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkConnect(Ljava/lang/String;ILjava/lang/Object;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkCreateClassLoader()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkDelete(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkExec(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkExit(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkLink(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkListen(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkMemberAccess(Ljava/lang/Class;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkMulticast(Ljava/net/InetAddress;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkMulticast(Ljava/net/InetAddress;B)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkPackageAccess(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkPackageDefinition(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkPermission(Ljava/security/Permission;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkPermission(Ljava/security/Permission;Ljava/lang/Object;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkPrintJobAccess()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkPropertiesAccess()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkPropertyAccess(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkRead(Ljava/io/FileDescriptor;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkRead(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkRead(Ljava/lang/String;Ljava/lang/Object;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkSecurityAccess(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkSetFactory()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkSystemClipboardAccess()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkTopLevelWindow(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkWrite(Ljava/io/FileDescriptor;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".checkWrite(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".classDepth(Ljava/lang/String;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".classLoaderDepth()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".currentClassLoader()Ljava/lang/ClassLoader;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".currentLoadedClass()Ljava/lang/Class;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".getClassContext()[Ljava/lang/Class;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".getInCheck()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".getSecurityContext()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures[classNameSecurityManager+".getThreadGroup()Ljava/lang/ThreadGroup;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapDeprecated,
		}
}
