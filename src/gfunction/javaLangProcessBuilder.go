/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

func Load_Lang_Process_Builder() {

	MethodSignatures["java/lang/ProcessBuilder.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/lang/ProcessBuilder.<init>(Ljava/util/List;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.command(Ljava/util/List;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.command([Ljava/lang/String;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.directory()Ljava/io/File;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.directory(Ljava/io/File;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.environment()Ljava/util/Map;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.inheritIO()Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.redirectError()Ljava/lang/ProcessBuilder$Redirect;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.redirectError(Ljava/lang/ProcessBuilder$Redirect;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.redirectErrorStream()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.redirectErrorStream(Z)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.redirectInput()Ljava/lang/ProcessBuilder$Redirect;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.redirectInput(Ljava/io/File;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.redirectInput(Ljava/lang/ProcessBuilder$Redirect;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.redirectOutput()Ljava/lang/ProcessBuilder$Redirect;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.redirectOutput(Ljava/io/File;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.redirectOutput(Ljava/lang/ProcessBuilder$Redirect;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessBuilder.start()Ljava/lang/Process;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}
}
