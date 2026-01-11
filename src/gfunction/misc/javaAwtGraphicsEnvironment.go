/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package misc

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
)

func Load_Awt_Graphics_Environment() {

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.checkHeadless()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.createGraphics(Ljava/awt/image/BufferedImage;)Ljava/awt/Graphics2D;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.getAllFonts()[Ljava/awt/Font;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.getAvailableFontFamilyNames()[Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.getAvailableFontFamilyNames(Ljava/util/Locale;)[Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.getCenterPoint()Ljava/awt/Point;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.getDefaultScreenDevice()Ljava/awt/GraphicsDevice;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.getLocalGraphicsEnvironment()Ljava/awt/GraphicsEnvironment;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.getMaximumWindowBounds()Ljava/awt/Rectangle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.getScreenDevices()[Ljava/awt/GraphicsDevice;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.isHeadless()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  awtgeIsHeadless,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.isHeadlessInstance()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  awtgeIsHeadless,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.preferLocaleFonts()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.preferProportionalFonts()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/GraphicsEnvironment.registerFont(Ljava/awt/Font;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

}

// "java/awt/GraphicsEnvironment.isHeadless()Z"
func awtgeIsHeadless(params []interface{}) interface{} {
	glob := globals.GetGlobalRef()
	return glob.Headless
}
