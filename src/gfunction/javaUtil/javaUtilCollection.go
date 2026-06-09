/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Util_Collection() {

	ghelpers.MethodSignatures["java/util/Collection.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/Collection.add(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  hashsetAdd,
		}

	ghelpers.MethodSignatures["java/util/Collection.addAll(Ljava/util/Collection;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Collection.clear()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  hashmapInit,
		}

	ghelpers.MethodSignatures["java/util/Collection.contains(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  hashsetContains,
		}

	ghelpers.MethodSignatures["java/util/Collection.containsAll(Ljava/util/Collection;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Collection.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Collection.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Collection.isEmpty()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  hashsetIsEmpty,
		}

	ghelpers.MethodSignatures["java/util/Collection.iterator()Ljava/util/Iterator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  setIterator,
		}

	ghelpers.MethodSignatures["java/util/Collection.remove(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  hashsetRemove,
		}

	ghelpers.MethodSignatures["java/util/Collection.removeAll(Ljava/util/Collection;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Collection.retainAll(Ljava/util/Collection;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Collection.size()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  hashmapSize,
		}

	ghelpers.MethodSignatures["java/util/Collection.toArray()[Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction, // As requested
		}

	ghelpers.MethodSignatures["java/util/Collection.toArray([Ljava/lang/Object;)[Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction, // As requested
		}

	// From issue: parallelStream, removeIf, spliterator, and toArray should be trapped.

	ghelpers.MethodSignatures["java/util/Collection.parallelStream()Ljava/util/stream/Stream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Collection.removeIf(Ljava/util/function/Predicate;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Collection.spliterator()Ljava/util/Spliterator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Collection.stream()Ljava/util/stream/Stream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Collection.forEach(Ljava/util/function/Consumer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}
}
