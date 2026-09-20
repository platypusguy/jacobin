/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaNio

import (
	"jacobin/src/gfunction/ghelpers"
)

// Load_Nio_Channels_ServerSocketChannel registers trap entries for
// java.nio.channels.ServerSocketChannel. None of this class is currently
// implemented, so every method routes to TrapFunction.
func Load_Nio_Channels_ServerSocketChannel() {

	ghelpers.MethodSignatures["java/nio/channels/ServerSocketChannel.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/nio/channels/ServerSocketChannel.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/nio/channels/ServerSocketChannel.accept()Ljava/nio/channels/SocketChannel;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/nio/channels/ServerSocketChannel.bind(Ljava/net/SocketAddress;)Ljava/nio/channels/ServerSocketChannel;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/nio/channels/ServerSocketChannel.bind(Ljava/net/SocketAddress;I)Ljava/nio/channels/ServerSocketChannel;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/nio/channels/ServerSocketChannel.getLocalAddress()Ljava/net/SocketAddress;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/nio/channels/ServerSocketChannel.open()Ljava/nio/channels/ServerSocketChannel;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/nio/channels/ServerSocketChannel.setOption(Ljava/net/SocketOption;Ljava/lang/Object;)Ljava/nio/channels/ServerSocketChannel;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/nio/channels/ServerSocketChannel.socket()Ljava/net/ServerSocket;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/nio/channels/ServerSocketChannel.validOps()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}
}
