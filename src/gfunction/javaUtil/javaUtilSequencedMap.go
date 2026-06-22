/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Util_SequencedMap() {

	ghelpers.MethodSignatures["java/util/SequencedMap.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/SequencedMap.firstEntry()Ljava/util/Map$Entry;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/SequencedMap.lastEntry()Ljava/util/Map$Entry;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/SequencedMap.pollFirstEntry()Ljava/util/Map$Entry;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/SequencedMap.pollLastEntry()Ljava/util/Map$Entry;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/SequencedMap.putFirst(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/SequencedMap.putLast(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/SequencedMap.reversed()Ljava/util/SequencedMap;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/SequencedMap.sequencedEntrySet()Ljava/util/SequencedSet;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/SequencedMap.sequencedKeySet()Ljava/util/SequencedSet;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/SequencedMap.sequencedValues()Ljava/util/SequencedCollection;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}
}
