package javaLang

import "jacobin/src/gfunction/ghelpers"

func Load_Lang_Record() {

	ghelpers.MethodSignatures["java/lang/Record.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapClass,
		}

	ghelpers.MethodSignatures["java/lang/Record.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

}
