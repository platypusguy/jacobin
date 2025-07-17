/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

func Load_Lang_Process_Builder() {

	MethodSignatures["javaLangProcessBuilder.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.<init>(Ljava/util/List;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.command(Ljava/util/List;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.command([Ljava/lang/String;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.directory()Ljava/io/File;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.directory(Ljava/io/File;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.environment()Ljava/util/Map;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.inheritIO()Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.redirectError()Ljava/lang/ProcessBuilder$Redirect;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.redirectError(Ljava/lang/ProcessBuilder$Redirect;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.redirectErrorStream()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.redirectErrorStream(Z)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.redirectInput()Ljava/lang/ProcessBuilder$Redirect;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.redirectInput(Ljava/io/File;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.redirectInput(Ljava/lang/ProcessBuilder$Redirect;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.redirectOutput()Ljava/lang/ProcessBuilder$Redirect;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.redirectOutput(Ljava/io/File;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.redirectOutput(Ljava/lang/ProcessBuilder$Redirect;)Ljava/lang/ProcessBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["javaLangProcessBuilder.start()Ljava/lang/Process;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}
}
