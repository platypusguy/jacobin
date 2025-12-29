/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/globals"
)

func Load_Awt_Graphics_Environment() {

	MethodSignatures["java/awt/GraphicsEnvironment.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.checkHeadless()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.createGraphics(Ljava/awt/image/BufferedImage;)Ljava/awt/Graphics2D;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.getAllFonts()[Ljava/awt/Font;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.getAvailableFontFamilyNames()[Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.getAvailableFontFamilyNames(Ljava/util/Locale;)[Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.getCenterPoint()Ljava/awt/Point;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.getDefaultScreenDevice()Ljava/awt/GraphicsDevice;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.getLocalGraphicsEnvironment()Ljava/awt/GraphicsEnvironment;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.getMaximumWindowBounds()Ljava/awt/Rectangle;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.getScreenDevices()[Ljava/awt/GraphicsDevice;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.isHeadless()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  awtgeIsHeadless,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.isHeadlessInstance()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  awtgeIsHeadless,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.preferLocaleFonts()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.preferProportionalFonts()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.registerFont(Ljava/awt/Font;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

}

// "java/awt/GraphicsEnvironment.isHeadless()Z"
func awtgeIsHeadless(params []interface{}) interface{} {
	glob := globals.GetGlobalRef()
	return glob.Headless
}
