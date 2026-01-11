/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-4 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/misc"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// We don't run String's static initializer block because the initialization
// is already handled in String creation

func Load_Lang_String() {

	// === OBJECT INSTANTIATION ===

	// String instantiation without parameters i.e. String string = new String();
	// need to replace eventually by enabling the Java initializer to run
	ghelpers.MethodSignatures["java/lang/String.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringClinit,
		}

	// Instantiate an empty String
	ghelpers.MethodSignatures["java/lang/String.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  newEmptyString,
		}

	// String(byte[] bytes) - instantiate a String from a byte array
	ghelpers.MethodSignatures["java/lang/String.<init>([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  newStringFromBytes,
		}

	// String(byte[] ascii, int hibyte) *** DEPRECATED
	ghelpers.MethodSignatures["java/lang/String.<init>([BI)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapDeprecated,
		}

	// String(byte[] bytes, int offset, int length)	- instantiate a String from a subset of a byte array
	ghelpers.MethodSignatures["java/lang/String.<init>([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  newStringFromBytesSubset,
		}

	// String(byte[] ascii, int hibyte, int offset, int count) *** DEPRECATED
	ghelpers.MethodSignatures["java/lang/String.<init>([BIII)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapDeprecated,
		}

	// TODO: String(byte[] bytes, int offset, int length, String charsetName) *********** CHARSET
	ghelpers.MethodSignatures["java/lang/String.<init>([BIILjava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	// TODO: String(byte[] bytes, int offset, int length, Charset charset) ************** CHARSET
	ghelpers.MethodSignatures["java/lang/String.<init>([BIILjava/nio/charset/Charset;)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	// TODO: String(byte[] bytes, String charsetName) *********************************** CHARSET
	ghelpers.MethodSignatures["java/lang/String.<init>([BLjava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// TODO: String(byte[] bytes, Charset charset) ************************************** CHARSET
	ghelpers.MethodSignatures["java/lang/String.<init>([BLjava/nio/charset/Charset;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// Instantiate a String from a character array
	ghelpers.MethodSignatures["java/lang/String.<init>([C)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  newStringFromChars,
		}

	// Instantiate a String from a subset of a character array
	ghelpers.MethodSignatures["java/lang/String.<init>([CII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  newStringFromChars,
		}

	// TODO: String(int[] codePoints, int offset, int count) ************************ CODEPOINTS
	ghelpers.MethodSignatures["java/lang/String.<init>([III)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	// String(String original) -- instantiate a String from another String.
	ghelpers.MethodSignatures["java/lang/String.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  newStringFromString,
		}

	// String(StringBuffer buffer) -- instantiate a String from a StringBuffer.
	ghelpers.MethodSignatures["java/lang/String.<init>(Ljava/lang/StringBuffer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  newStringFromString,
		}

	// Not in API: Is the String Latin1?
	ghelpers.MethodSignatures["java/lang/String.isLatin1()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringIsLatin1,
		}

	// String(StringBuilder builder) -- instantiate a String from a StringBuilder.
	ghelpers.MethodSignatures["java/lang/String.<init>(Ljava/lang/StringBuilder;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  newStringFromString,
		}

	// ==== METHOD FUNCTIONS (in alphabetical order by their function names) ====

	// Returns the char value at the specified index.
	ghelpers.MethodSignatures["java/lang/String.charAt(I)C"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringCharAt,
		}

	// TODO: Returns a stream of int zero-extending the char values from this sequence.
	ghelpers.MethodSignatures["java/lang/String.chars()Ljava/util/stream/IntStream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	// Internal boundary-checker - not in the API.
	ghelpers.MethodSignatures["java/lang/String.checkBoundsBeginEnd(III)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  stringCheckBoundsBeginEnd,
		}

	// Internal boundary-checker - not in the API.
	ghelpers.MethodSignatures["java/lang/String.checkBoundsOffCount(III)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  stringCheckBoundsOffCount,
		}

	// TODO: Returns the character (Unicode code point) at the specified index.
	ghelpers.MethodSignatures["java/lang/String.codePointAt(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	// TODO: Returns the character (Unicode code point) before the specified index.
	ghelpers.MethodSignatures["java/lang/String.codePointBefore(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	// TODO: Returns the number of Unicode code points in the specified text range of this String.
	ghelpers.MethodSignatures["java/lang/String.codePointCount(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// TODO: Returns a stream of code point values from this sequence.
	ghelpers.MethodSignatures["java/lang/String.codePoints()Ljava/util/stream/IntStream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	// Compare 2 strings lexicographically, case-sensitive (upper/lower).
	// The return value is a negative integer, zero, or a positive integer
	// as the String argument is greater than, equal to, or less than this String,
	// case-sensitive.
	ghelpers.MethodSignatures["java/lang/String.compareTo(Ljava/lang/String;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringCompareToCaseSensitive,
		}

	// Compare 2 strings lexicographically, ignoring case (upper/lower).
	// The return value is a negative integer, zero, or a positive integer
	// as the String argument is greater than, equal to, or less than this String,
	// ignoring case considerations.
	ghelpers.MethodSignatures["java/lang/String.compareToIgnoreCase(Ljava/lang/String;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringCompareToIgnoreCase,
		}

	// Concatenates the specified string to the end of this string.
	ghelpers.MethodSignatures["java/lang/String.concat(Ljava/lang/String;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringConcat,
		}

	// Returns true if and only if this string contains the specified sequence of char values.
	ghelpers.MethodSignatures["java/lang/String.contains(Ljava/lang/CharSequence;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringContains,
		}

	// Compares this string to the specified CharSequence.
	ghelpers.MethodSignatures["java/lang/String.contentEquals(Ljava/lang/CharSequence;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  javaLangStringContentEquals,
		}

	// Compares this string to the specified StringBuffer.
	ghelpers.MethodSignatures["java/lang/String.contentEquals(Ljava/lang/StringBuffer;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  javaLangStringContentEquals,
		}

	// Return a string representing a char array.
	ghelpers.MethodSignatures["java/lang/String.copyValueOf([C)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  valueOfCharArray,
		}

	// Return a string representing a char subarray.
	ghelpers.MethodSignatures["java/lang/String.copyValueOf([CII)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  newStringFromCharsSubset,
		}

	// TODO: Returns an Optional containing the nominal descriptor for this instance, which is the instance itself.
	ghelpers.MethodSignatures["java/lang/String.describeConstable()Ljava/util/Optional;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	// OpenJDK JVM "java/lang/String.endsWith(Ljava/lang/String;)Z" works with the jacobin String object.
	// Does the base string end with the specified suffix argument?

	// Compares this string to the specified object.
	ghelpers.MethodSignatures["java/lang/String.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringEquals,
		}

	// Compares this String to another String, ignoring case considerations.
	ghelpers.MethodSignatures["java/lang/String.equalsIgnoreCase(Ljava/lang/String;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringEqualsIgnoreCase,
		}

	// Return a formatted string using the reference object string as the format string
	// and the supplied arguments as input object arguments.
	// E.g. String string = String.format("%s %i", "ABC", 42);
	ghelpers.MethodSignatures["java/lang/String.format(Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  sprintf,
		}

	// TODO: Return a formatted string using the specified locale, format string, and arguments.
	ghelpers.MethodSignatures["java/lang/String.format(Ljava/util/Locale;Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	// This method is equivalent to String.format(this, args).
	ghelpers.MethodSignatures["java/lang/String.formatted([Ljava/lang/Object;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  sprintf,
		}

	// Encodes this String into a sequence of bytes using the default charset, storing the result into a new byte array.
	ghelpers.MethodSignatures["java/lang/String.getBytes()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  getBytesFromString,
		}

	// void getBytes(int srcBegin, int srcEnd, byte[] dst, int dstBegin)  ********************* DEPRECATED
	ghelpers.MethodSignatures["java/lang/String.getBytes(II[BI)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapDeprecated,
		}

	// TODO: Encodes this String into a sequence of bytes using the given charset, storing the result into a new byte array.
	ghelpers.MethodSignatures["java/lang/String.getBytes(Ljava/nio/charset/Charset;)[B"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	// Encodes this String into a sequence of bytes using the named charset, storing the result into a new byte array. ************************ CHARSET
	ghelpers.MethodSignatures["java/lang/String.getBytes(Ljava/lang/String;)[B"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	// Not in API: getBytes([BIIBI)V
	// original Java source: https://gist.github.com/platypusguy/03c1a9e3acb1cb2cfc2d821aa2dd4490
	ghelpers.MethodSignatures["java/lang/String.getBytes([BIIBI)V"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  stringGetBytesBIIBI,
		}

	// Copies characters from this string into the destination character array.
	ghelpers.MethodSignatures["java/lang/String.getChars(II[CI)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  stringGetChars,
		}

	// Compute the Java String.hashCode() value.
	ghelpers.MethodSignatures["java/lang/String.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringHashCode,
		}

	// TODO: Adjusts the indentation of each line of this string based on the value of n, and normalizes line termination characters.
	ghelpers.MethodSignatures["java/lang/String.indent(I)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	// Returns the index within this string of the first occurrence of the specified character.
	ghelpers.MethodSignatures["java/lang/String.indexOf(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringIndexOfCh,
		}

	// Returns the index within this string of the first occurrence of the specified character,
	// starting the search at the specified index.
	ghelpers.MethodSignatures["java/lang/String.indexOf(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringIndexOfCh,
		}

	// Returns the index within this string of the first occurrence of the specified character,
	// starting the search at beginIndex and stopping before endIndex.
	ghelpers.MethodSignatures["java/lang/String.indexOf(III)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  stringIndexOfCh,
		}

	// Returns the index within this string of the first occurrence of the specified substring.
	ghelpers.MethodSignatures["java/lang/String.indexOf(Ljava/lang/String;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringIndexOfString,
		}

	// Returns the index within this string of the first occurrence of the specified substring,
	// starting at the specified index.
	ghelpers.MethodSignatures["java/lang/String.indexOf(Ljava/lang/String;I)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringIndexOfString,
		}

	// Returns the index of the first occurrence of the specified substring within the specified index range of this string.
	ghelpers.MethodSignatures["java/lang/String.indexOf(Ljava/lang/String;II)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  stringIndexOfString,
		}

	/*
		When the intern method is invoked, if the pool already contains a string equal to this String object as determined
		by the equals(Object) method, then the string from the pool is returned.
		Otherwise, this String object is added to the pool and a reference to this String object is returned.
	*/
	ghelpers.MethodSignatures["java/lang/String.intern()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringIntern,
		}

	// Is the base string whitespace?
	ghelpers.MethodSignatures["java/lang/String.isBlank()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringIsBlank,
		}

	// Is the base string empty?
	ghelpers.MethodSignatures["java/lang/String.isEmpty()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringIsEmpty,
		}

	// TODO: Returns a new String composed of copies of the CharSequence elements joined together with a copy of the specified delimiter.
	ghelpers.MethodSignatures["java/lang/String.join(Ljava/lang/CharSequence;[Ljava/lang/CharSequence;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// TODO: Returns a new String composed of copies of the CharSequence elements joined together with a copy of the specified delimiter.
	ghelpers.MethodSignatures["java/lang/String.join(Ljava/lang/CharSequence;[Ljava/lang/Iterable;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// Returns the index within this string of the last occurrence of the specified character.
	ghelpers.MethodSignatures["java/lang/String.lastIndexOf(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  lastIndexOfCharacter,
		}

	// Returns the index within this string of the last occurrence of the specified character, searching backward starting at the specified index.
	ghelpers.MethodSignatures["java/lang/String.lastIndexOf(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  lastIndexOfCharacter,
		}

	// Returns the index within this string of the last occurrence of the specified substring.
	ghelpers.MethodSignatures["java/lang/String.lastIndexOf(Ljava/lang/String;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  lastIndexOfString,
		}

	// Returns the index within this string of the last occurrence of the specified substring, searching backward starting at the specified index.
	ghelpers.MethodSignatures["java/lang/String.lastIndexOf(Ljava/lang/String;I)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  lastIndexOfString,
		}

	// Return the length of the base String.
	ghelpers.MethodSignatures["java/lang/String.length()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringLength,
		}

	// TODO: Returns a stream of lines extracted from this string, separated by line terminators.
	ghelpers.MethodSignatures["java/lang/String.lines()Ljava/util/stream/Stream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	// Tells whether this string matches the given regular expression or not.
	ghelpers.MethodSignatures["java/lang/String.matches(Ljava/lang/String;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringMatches,
		}

	// TODO: Returns the index within this String that is offset from the given index by codePointOffset code points.
	ghelpers.MethodSignatures["java/lang/String.offsetByCodePoints(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// Tests if two string regions are equal.
	// Pass a flag indicating whether to ignore case or not.
	ghelpers.MethodSignatures["java/lang/String.regionMatches(ZILjava/lang/String;II)Z"] = // Has an ignoreCase flag
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  stringRegionMatches,
		}

	// Tests if two string regions are equal, case-sensitive.
	ghelpers.MethodSignatures["java/lang/String.regionMatches(ILjava/lang/String;II)Z"] = // Does not have an ignoreCase flag
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  stringRegionMatches,
		}

	// Returns a string whose value is the concatenation of this string repeated the specified number of times.
	ghelpers.MethodSignatures["java/lang/String.repeat(I)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringRepeat,
		}

	// Replace a single character by another in the given string.
	ghelpers.MethodSignatures["java/lang/String.replace(CC)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringReplaceCC,
		}

	// Replace a character sequence by another in the given string.
	ghelpers.MethodSignatures["java/lang/String.replace(Ljava/lang/CharSequence;Ljava/lang/CharSequence;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringReplaceLiteral,
		}

	// Replaces each substring of this string that matches the given regular expression with the given replacement.
	ghelpers.MethodSignatures["java/lang/String.replaceAll(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringReplaceAllRegex,
		}

	// Replaces the first substring of this string that matches the given regular expression with the given replacement.
	ghelpers.MethodSignatures["java/lang/String.replaceFirst(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringReplaceFirstRegex,
		}

	// TODO: Resolves this instance as a ConstantDesc, the result of which is the instance itself.
	ghelpers.MethodSignatures["java/lang/String.resolveConstantDesc(Ljava/lang/invoke/MethodHandles/Lookup;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	// Split the base string into an array of strings.
	ghelpers.MethodSignatures["java/lang/String.split(Ljava/lang/String;)[Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringSplit,
		}

	// Split the base string into an array of strings with a specified limit.
	ghelpers.MethodSignatures["java/lang/String.split(Ljava/lang/String;I)[Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringSplitLimit,
		}

	// TODO: Split the base string around matches of the given regular expression and returns both the strings and the matching delimiters.
	ghelpers.MethodSignatures["java/lang/String.splitWithDelimiters(Ljava/lang/String;I)[Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// Tests if this string starts with the specified prefix.
	ghelpers.MethodSignatures["java/lang/String.startsWith(Ljava/lang/String;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringStartsWith,
		}

	// Tests if the substring of this string beginning at the specified index starts with the specified prefix.
	ghelpers.MethodSignatures["java/lang/String.startsWith(Ljava/lang/String;I)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringStartsWith,
		}

	// Returns a string whose value is this string, with all leading and trailing white space removed.
	ghelpers.MethodSignatures["java/lang/String.strip()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringStrip,
		}

	// TODO: Returns a string whose value is this string, with incidental white space removed from the beginning and end of every line.
	ghelpers.MethodSignatures["java/lang/String.stripIndent()[Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	// Returns a string whose value is the base string with all leading white space removed.
	ghelpers.MethodSignatures["java/lang/String.stripLeading()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringStripLeading,
		}

	// Returns a string whose value is the base string with all trailing white space removed.
	ghelpers.MethodSignatures["java/lang/String.stripTrailing()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringStripTrailing,
		}

	// TODO: Returns a character sequence that is a subsequence of this sequence.
	ghelpers.MethodSignatures["java/lang/String.subSequence(II)Ljava/lang/CharSequence;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// Return a substring starting at the given index of the byte array.
	ghelpers.MethodSignatures["java/lang/String.substring(I)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  substringToTheEnd,
		}

	// Return a substring starting at the given index of the byte array of the given length.
	ghelpers.MethodSignatures["java/lang/String.substring(II)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  substringStartEnd,
		}

	// Return a string in all lower case, using the reference object string as input.
	ghelpers.MethodSignatures["java/lang/String.toCharArray()[C"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  toCharArray,
		}

	// Return a string in all lower case, using the reference object string as input.
	ghelpers.MethodSignatures["java/lang/String.toLowerCase()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  toLowerCase,
		}

	// TODO: Converts all of the characters in this String to lower case using the rules of the given Locale.
	ghelpers.MethodSignatures["java/lang/String.toLowerCase(Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  toLowerCase, // TODO: Locale processing
		}

	// Return the base string as-is.
	ghelpers.MethodSignatures["java/lang/String.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringToString,
		}

	// Return a string in all upper case, using the reference object string as input.
	ghelpers.MethodSignatures["java/lang/String.toUpperCase()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  toUpperCase,
		}

	// TODO: Converts all of the characters in this String to upper case using the rules of the given Locale.
	ghelpers.MethodSignatures["java/lang/String.toUpperCase(Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	// TODO: What should we do with this? <R> R transform(Function<? super String,? extends R> f)

	// TODO: Return a string whose value is the base string with escape sequences translated as if in a string literal.
	ghelpers.MethodSignatures["java/lang/String.translateEscapes()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	// Return a string trimmed of leading and trailing whitespace.
	ghelpers.MethodSignatures["java/lang/String.trim()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  trimString,
		}

	// Return a string representing a boolean value.
	ghelpers.MethodSignatures["java/lang/String.valueOf(Z)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  valueOfBoolean,
		}

	// Return a string representing a char value.
	ghelpers.MethodSignatures["java/lang/String.valueOf(C)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  valueOfChar,
		}

	// Return a string representing a char array.
	ghelpers.MethodSignatures["java/lang/String.valueOf([C)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  valueOfCharArray,
		}

	// Return a string representing a char subarray.
	ghelpers.MethodSignatures["java/lang/String.valueOf([CII)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  newStringFromCharsSubset,
		}

	// Return a string representing a double value.
	ghelpers.MethodSignatures["java/lang/String.valueOf(D)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  valueOfDouble,
		}

	// Return a string representing a float value.
	ghelpers.MethodSignatures["java/lang/String.valueOf(F)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  valueOfFloat,
		}

	// Return a string representing an int value.
	ghelpers.MethodSignatures["java/lang/String.valueOf(I)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  valueOfInt,
		}

	// Return a string representing an int value.
	ghelpers.MethodSignatures["java/lang/String.valueOf(J)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  valueOfLong,
		}

	// Return a string representing the value of an Object.
	ghelpers.MethodSignatures["java/lang/String.valueOf(Ljava/lang/Object;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  valueOfObject,
		}

}

// ==== INSTANTIATION AND INITIALIZATION FUNCTIONS ====

// "java/lang/String.<clinit>()V" -- String class initialisation
func stringClinit([]interface{}) interface{} {
	klass := classloader.MethAreaFetch(types.StringClassName)
	if klass == nil || klass.Data == nil {
		errMsg := fmt.Sprintf("stringClinit: Could not find class %s in the MethodArea", types.StringClassName)
		return ghelpers.GetGErrBlk(excNames.ClassNotLoadedException, errMsg)
	}
	klass.Data.ClInit = types.ClInitRun // just mark that String.<clinit>() has been run
	return nil
}

// Instantiate a new empty string - "java/lang/String.<init>()V"
func newEmptyString(params []interface{}) interface{} {
	// params[0] = target object for string (updated)
	obj := params[0].(*object.Object)
	bytes := make([]types.JavaByte, 0)
	object.UpdateValueFieldFromJavaBytes(obj, bytes)
	return nil
}

// Instantiate a new string object from a Go byte array.
// "java/lang/String.<init>([B)V"
func newStringFromBytes(params []interface{}) interface{} {
	// params[0] = reference string (to be updated with byte array)
	// params[1] = byte array object
	if params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "newStringFromBytes: null parameter")
	}
	obj := params[0].(*object.Object)
	fld, ok := params[1].(*object.Object).FieldTable["value"]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "newStringFromBytes: missing value field")
	}
	switch fld.Fvalue.(type) {
	case []byte:
		bytes := object.JavaByteArrayFromGoByteArray(fld.Fvalue.([]byte))
		object.UpdateValueFieldFromJavaBytes(obj, bytes)
	case []types.JavaByte:
		bytes := fld.Fvalue.([]types.JavaByte)
		object.UpdateValueFieldFromJavaBytes(obj, bytes)
	}
	return nil
}

// Construct a string object from a subset of a JavaByte array.
// "java/lang/String.<init>([BII)V"
func newStringFromBytesSubset(params []interface{}) interface{} {
	// params[0] = reference string (to be updated with byte array)
	// params[1] = byte array object
	// params[2] = start offset
	// params[3] = length
	if params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "newStringFromBytesSubset: null parameter")
	}
	obj := params[0].(*object.Object)
	fld, ok := params[1].(*object.Object).FieldTable["value"]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "newStringFromBytesSubset: missing value field")
	}
	var bytes []types.JavaByte
	switch fld.Fvalue.(type) {
	case []byte:
		bytes = object.JavaByteArrayFromGoByteArray(fld.Fvalue.([]byte))
	case []types.JavaByte:
		bytes = fld.Fvalue.([]types.JavaByte)
	}

	// Get substring start and length
	ssStart := params[2].(int64)
	ssLen := params[3].(int64)

	// Validate boundaries.
	totalLength := int64(len(bytes))
	if ssStart < 0 || ssLen < 0 || (ssStart+ssLen) > totalLength {
		errMsg1 := "newStringFromBytesSubset: Invalid offset or length"
		errMsg2 := fmt.Sprintf("\n\twholelen=%d, offset=%d, length=%d\n\n", totalLength, ssStart, ssLen)
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg1+errMsg2)
	}

	// Compute subarray and update params[0].
	bytes = bytes[ssStart : ssStart+ssLen]
	object.UpdateValueFieldFromJavaBytes(obj, bytes)
	return nil
}

// Instantiate a new string object from a Go int64 array (Java char array).
// "java/lang/String.<init>([C)V"
func newStringFromChars(params []interface{}) interface{} {
	// params[0] = reference string (to be updated with byte array)
	// params[1] = char array object
	if params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "newStringFromChars: null parameter")
	}
	obj := params[0].(*object.Object)
	fld, ok := params[1].(*object.Object).FieldTable["value"]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "newStringFromChars: missing value field")
	}
	ints, ok := fld.Fvalue.([]int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "newStringFromChars: value field is not []int64")
	}

	var bytes []types.JavaByte
	for _, ii := range ints {
		bytes = append(bytes, types.JavaByte(ii&0xFF))
	}
	object.UpdateValueFieldFromJavaBytes(obj, bytes)
	return nil
}

// Construct a string object from a subset of a character array.
// "java/lang/String.valueOf([CII)Ljava/lang/String;"
func newStringFromCharsSubset(params []interface{}) interface{} {
	// params[0] = character array object
	// params[1] = start offset
	// params[2] = length
	// Return the string.
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "newStringFromCharsSubset: null parameter")
	}
	fld, ok := params[0].(*object.Object).FieldTable["value"]
	if !ok {
		errMsg := fmt.Sprintf("newStringFromCharsSubset: Missing value field in character array object")
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	iarray, ok := fld.Fvalue.([]int64)
	if !ok {
		errMsg := fmt.Sprintf("newStringFromCharsSubset: Invalid value field type (%s : %T) in character array object", fld.Ftype, fld.Fvalue)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get substring start and length
	ssStart := params[1].(int64)
	ssLen := params[2].(int64)

	// Validate boundaries.
	totalLength := int64(len(iarray))
	if ssStart < 0 || ssLen < 0 || (ssStart+ssLen) > totalLength {
		errMsg1 := "newStringFromCharsSubset: Invalid offset or length"
		errMsg2 := fmt.Sprintf("\n\twholelen=%d, offset=%d, length=%d\n\n", totalLength, ssStart, ssLen)
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg1+errMsg2)
	}

	// Compute subarray.
	iarray = iarray[ssStart : ssStart+ssLen]
	var bytes []types.JavaByte
	for _, ii := range iarray {
		bytes = append(bytes, types.JavaByte(ii&0xFF))
	}
	obj := object.StringObjectFromJavaByteArray(bytes)
	return obj

}

// New String (consisting of JavaBytes) from String, StringBuilder, or StringBuffer.
func newStringFromString(params []interface{}) interface{} {
	// params[0] = reference string (to be updated with byte array)
	// params[1] = String, StringBuilder, or StringBuffer object
	var javaBytes []types.JavaByte
	switch params[1].(*object.Object).FieldTable["value"].Fvalue.(type) {
	case []byte:
		bytes := params[1].(*object.Object).FieldTable["value"].Fvalue.([]byte)
		javaBytes = object.JavaByteArrayFromGoByteArray(bytes)
	case []types.JavaByte:
		javaBytes = params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	}

	object.UpdateValueFieldFromJavaBytes(params[0].(*object.Object), javaBytes)
	return nil
}

// ==== METHODS FOR STRING ACTIVITIES ====

// Get character at the given index.
// "java/lang/String.charAt(I)C"
func stringCharAt(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringCharAt: null parameter")
	}
	// Unpack the reference string and convert it to a rune array.
	ptrObj := params[0].(*object.Object)
	str := object.GoStringFromStringObject(ptrObj)
	runeArray := []rune(str)

	// Get index.
	index := params[1].(int64)

	if index < 0 || int(index) >= len(runeArray) {
		errMsg := fmt.Sprintf("stringCharAt: Index out of bounds: %d, length: %d", index, len(runeArray))
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
	}

	// Return indexed character.
	runeValue := runeArray[index]
	return int64(runeValue)
}

// "java/lang/String.compareTo(Ljava/lang/String;)I"
func stringCompareToCaseSensitive(params []interface{}) interface{} {
	if params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringCompareToCaseSensitive: null parameter")
	}
	str1 := object.GoStringFromStringObject(params[0].(*object.Object))
	str2 := object.GoStringFromStringObject(params[1].(*object.Object))
	if str2 == str1 {
		return int64(0)
	}
	if str1 < str2 {
		return int64(-1)
	}
	return int64(1)
}

// "java/lang/String.compareToIgnoreCase(Ljava/lang/String;)I"
func stringCompareToIgnoreCase(params []interface{}) interface{} {
	if params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringCompareToIgnoreCase: null parameter")
	}
	str1 := strings.ToLower(object.GoStringFromStringObject(params[0].(*object.Object)))
	str2 := strings.ToLower(object.GoStringFromStringObject(params[1].(*object.Object)))
	if str2 == str1 {
		return int64(0)
	}
	if str1 < str2 {
		return int64(-1)
	}
	return int64(1)
}

// "java/lang/String.concat(Ljava/lang/String;)Ljava/lang/String;"
func stringConcat(params []interface{}) interface{} {
	if params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringConcat: null parameter")
	}

	str1 := object.GoStringFromStringObject(params[0].(*object.Object))
	str2 := object.GoStringFromStringObject(params[1].(*object.Object))

	str := str1 + str2
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.contains(Ljava/lang/CharSequence;)Z"
// charSequence is an interface, generally implemented via String or array of chars
// Here, we assume one of those two options.
func stringContains(params []interface{}) interface{} {
	if params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringContains: null parameter")
	}

	targetString := object.GoStringFromStringObject(params[0].(*object.Object))
	searchString := object.GoStringFromStringObject(params[1].(*object.Object))

	if strings.Contains(targetString, searchString) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func javaLangStringContentEquals(params []interface{}) interface{} {
	if params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "javaLangStringContentEquals: null parameter")
	}

	str1 := object.GoStringFromStringObject(params[0].(*object.Object))
	str2 := object.GoStringFromStringObject(params[1].(*object.Object))

	// Are they equal in value?
	if str1 == str2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Are 2 strings equal?
// "java/lang/String.equals(Ljava/lang/Object;)Z"
func stringEquals(params []interface{}) interface{} {
	// params[0]: reference string object
	// params[1]: compare-to string Object
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringEquals: null parameter")
	}

	str1 := object.GoStringFromStringObject(params[0].(*object.Object))

	if params[1] == nil {
		return types.JavaBoolFalse
	}
	// In Java, equals(Object) should check if it's a String.
	// Jacobin's object.GoStringFromStringObject handles non-string objects by returning "",
	// but we should ideally check if it's actually a String object.
	if !object.IsStringObject(params[1]) {
		return types.JavaBoolFalse
	}

	str2 := object.GoStringFromStringObject(params[1].(*object.Object))

	// Are they equal in value?
	if str1 == str2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Are 2 strings equal, ignoring case?
// "java/lang/String.equalsIgnoreCase(Ljava/lang/String;)Z"
func stringEqualsIgnoreCase(params []interface{}) interface{} {
	// params[0]: reference string object
	// params[1]: compare-to string Object
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringEqualsIgnoreCase: null parameter")
	}
	if params[1] == nil {
		return types.JavaBoolFalse
	}

	str1 := object.GoStringFromStringObject(params[0].(*object.Object))
	str2 := object.GoStringFromStringObject(params[1].(*object.Object))

	// Are they equal in value?
	if strings.EqualFold(str1, str2) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// "java/lang/String.format(Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/String;"
// "java/lang/String.formatted([Ljava/lang/Object;)Ljava/lang/String;"
func sprintf(params []interface{}) interface{} {
	// params[0]: format string
	// params[1]: argument slice (array of object pointers)
	return misc.StringFormatter(params)
}

// java/lang/String.getBytes()[B
func getBytesFromString(params []interface{}) interface{} {
	// params[0] = reference string with byte array to be returned
	bytes := object.JavaByteArrayFromStringObject(params[0].(*object.Object))
	return object.MakePrimitiveObject("[B", types.ByteArray, bytes)
}

// java/lang/String.getBytes([BIIBI)V
// JDK17 Java source: https://gist.github.com/platypusguy/03c1a9e3acb1cb2cfc2d821aa2dd4490
func stringGetBytesBIIBI(params []any) any {
	fmt.Fprintln(os.Stderr, "java/lang/String.getBytes([BIIBI)V *****************")
	return nil
}

// java/lang/String.lastIndex(char)
// java/lang/String.lastIndex(char, beginIndex)
// Finds the last instance of the search character in the base string.
// Returns an index if the character is found or -1 if the character is not found
func lastIndexOfCharacter(params []any) any {
	// Get base string.
	baseStringObject := params[0].(*object.Object)
	baseString := object.GoStringFromStringObject(baseStringObject)

	// Get search string argument.
	searchByte := byte(params[1].(int64))

	// Get index starting point.
	var beginIndex int64
	if len(params) > 2 {
		beginIndex = params[2].(int64)
	} else {
		beginIndex = int64(len(baseString))
	}

	// Find search argument in base string if it is there.
	lastIndex := strings.LastIndexByte(baseString[:beginIndex], searchByte)

	// Return success (index value) or failure (-1).
	return int64(lastIndex)
}

// java/lang/String.lastIndex(string)
// java/lang/String.lastIndex(string, beginIndex)
// finds the last instance of the search string in the base string. Returns an
// index to the first character if the string is found, -1 if the string is not found
func lastIndexOfString(params []any) any {
	// Get base string.
	baseStringObject := params[0].(*object.Object)
	baseString := object.GoStringFromStringObject(baseStringObject)

	// Get search string argument.
	searchStringObject := params[1].(*object.Object)
	searchString := object.GoStringFromStringObject(searchStringObject)

	// Get indes starting point.
	var beginIndex int64
	if len(params) > 2 {
		beginIndex = params[2].(int64)
	} else {
		beginIndex = int64(len(baseString))
	}

	// Find search argument in base string if it is there.
	lastIndex := strings.LastIndex(baseString[:beginIndex], searchString)

	// Return success (index value) or failure (-1).
	return int64(lastIndex)
}

// "java/lang/String.isLatin1()Z"
func stringIsLatin1(params []interface{}) interface{} {
	// TODO: Someday, the answer might be false.
	return types.JavaBoolTrue // true
}

// "java/lang/String.length()I"
func stringLength(params []interface{}) interface{} {
	// params[0] = string object whose string length is to be measured
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringLength: null parameter")
	}
	obj := params[0].(*object.Object)
	bytes := object.JavaByteArrayFromStringObject(obj)
	return int64(len(bytes))
}

// java/lang/String.matches(Ljava/lang/String;)Z
// is the string in params[0] a match for the regex in params[1]?
func stringMatches(params []any) any {
	if len(params) != 2 {
		errMsg := fmt.Sprintf("stringMatches: Expected a string and a regular expression")
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	if params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringMatches: null parameter")
	}
	baseStringObject := params[0].(*object.Object)
	baseString := object.GoStringFromStringObject(baseStringObject)

	regexStringObject := params[1].(*object.Object)
	regexString := object.GoStringFromStringObject(regexStringObject)

	regex, err := regexp.Compile(regexString)
	if err != nil {
		errMsg := fmt.Sprintf("stringMatches: Invalid regular expression: %s", regexString)
		return ghelpers.GetGErrBlk(excNames.PatternSyntaxException, errMsg)
	}
	if regex.MatchString(baseString) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// do two regions in a string match?
// https://docs.oracle.com/en/java/javase/17/docs/api/java.base/java/lang/String.html#regionMatches(boolean,int,java.lang.String,int,int)
func stringRegionMatches(params []any) any {
	// param[0] the base string
	// param[1] offset of region in base string
	// param[2] pointer to second string
	// param[3] offset in second string
	// param[4] length of region to compare
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringRegionMatches: null parameter")
	}
	baseStringObject := params[0].(*object.Object)
	baseByteArray := object.JavaByteArrayFromStringObject(baseStringObject)

	// If this call includes boolean ignoreCase, then the parameters are shifted in the params array.
	ignoreCase := false
	pix := 1 // Assume no boolean ignoreCase parameter is present.
	if len(params) > 5 {
		pix = 2                                                // The boolean ignoreCase parameter is present.
		ignoreCase = (params[1].(int64) == types.JavaBoolTrue) // Get the flag value.
	}

	baseOffset := params[pix].(int64)

	if params[pix+1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringRegionMatches: null parameter")
	}
	compareStringObject := params[pix+1].(*object.Object)
	compareByteArray := object.JavaByteArrayFromStringObject(compareStringObject)
	compareOffset := params[pix+2].(int64)

	if baseOffset < 0 || compareOffset < 0 { // in the JDK, this is the indicated response, rather than an exception(!)
		return types.JavaBoolFalse
	}

	regionLength := params[pix+3].(int64)
	if baseOffset+regionLength > int64(len(baseByteArray)) || // again, erroneous values simply return false
		compareOffset+regionLength > int64(len(compareByteArray)) {
		return types.JavaBoolFalse
	}

	section1 := baseByteArray[baseOffset : baseOffset+regionLength]
	section2 := compareByteArray[compareOffset : compareOffset+regionLength]
	if ignoreCase {
		if object.JavaByteArrayEqualsIgnoreCase(section1, section2) {
			return types.JavaBoolTrue
		}
	} else {
		if object.JavaByteArrayEquals(section1, section2) { // case-sensitive equal
			return types.JavaBoolTrue
		}
	}
	return types.JavaBoolFalse
}

// "java/lang/String.repeat(I)Ljava/lang/String;"
func stringRepeat(params []interface{}) interface{} {
	// params[0] = base string
	// params[1] = int64 repetition factor
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringRepeat: null parameter")
	}
	oldStr := object.GoStringFromStringObject(params[0].(*object.Object))
	var newStr string
	count := params[1].(int64)
	for ii := int64(0); ii < count; ii++ {
		newStr = newStr + oldStr
	}

	// Return new string in an object.
	obj := object.StringObjectFromGoString(newStr)
	return obj

}

// "java/lang/String.replace(CC)Ljava/lang/String;"
func stringReplaceCC(params []interface{}) interface{} {
	// params[0] = base string
	// params[1] = character to be replaced
	// params[2] = replacement character
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringReplaceCC: null parameter")
	}
	str := object.GoStringFromStringObject(params[0].(*object.Object))
	oldChar := byte((params[1].(int64)) & 0xFF)
	newChar := byte((params[2].(int64)) & 0xFF)
	newStr := strings.ReplaceAll(str, string(oldChar), string(newChar))

	// Return final string in an object.
	obj := object.StringObjectFromGoString(newStr)
	return obj
}

// "java/lang/String.substring(I)Ljava/lang/String;"
func substringToTheEnd(params []interface{}) interface{} {
	// params[0] = base string
	// params[1] = start offset
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "substringToTheEnd: null parameter")
	}
	str := object.GoStringFromStringObject(params[0].(*object.Object))

	// Get substring start offset and compute end offset
	ssStart := params[1].(int64)
	ssEnd := int64(len(str))

	// Validate boundaries.
	totalLength := int64(len(str))
	if ssStart < 0 || ssStart > totalLength {
		errMsg1 := "substringToTheEnd: Invalid substring offset"
		errMsg2 := fmt.Sprintf("\n\twhole='%s' wholelen=%d, offset=%d\n\n", str, totalLength, ssStart)
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg1+errMsg2)
	}

	// Compute substring.
	str = str[ssStart:ssEnd]

	// Return new string in an object.
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.substring(II)Ljava/lang/String;"
func substringStartEnd(params []interface{}) interface{} {
	// params[0] = base string
	// params[1] = start offset
	// params[2] = end offset

	obj, ok := params[0].(*object.Object)
	if !ok {
		if object.IsNull(params[0]) {
			errMsg := "substringStartEnd: params[0] is null"
			return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
		}
		errMsg := "substringStartEnd: params[0] is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	str := object.GoStringFromStringObject(obj)

	// Get substring start and end offset
	ssStart := params[1].(int64)
	ssEnd := params[2].(int64)

	// Validate boundaries.
	totalLength := int64(len(str))
	if ssStart < 0 || ssEnd > totalLength || ssStart > ssEnd {
		errMsg1 := "substringStartEnd: Invalid substring range"
		errMsg2 := fmt.Sprintf("\n\twhole='%s' wholelen=%d, start=%d, end=%d\n\n", str, totalLength, ssStart, ssEnd)
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg1+errMsg2)
	}

	// Compute substring.
	str = str[ssStart:ssEnd]

	// Return new string in an object.
	obj = object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.toCharArray()[C"
func toCharArray(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "toCharArray: null parameter")
	}
	// params[0]: input string
	str := object.GoStringFromStringObject(params[0].(*object.Object))
	runes := []rune(str)
	var iArray []int64
	for _, r := range runes {
		iArray = append(iArray, int64(r))
	}
	return object.MakePrimitiveObject("[C", types.CharArray, iArray)
}

// "java/lang/String.toLowerCase()Ljava/lang/String;"
func toLowerCase(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "toLowerCase: null parameter")
	}
	// params[0]: input string
	str := strings.ToLower(object.GoStringFromStringObject(params[0].(*object.Object)))
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.toUpperCase()Ljava/lang/String;"
func toUpperCase(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "toUpperCase: null parameter")
	}
	// params[0]: input string
	str := strings.ToUpper(object.GoStringFromStringObject(params[0].(*object.Object)))
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.trim()Ljava/lang/String;"
func trimString(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "trimString: null parameter")
	}
	// params[0]: input string
	str := object.GoStringFromStringObject(params[0].(*object.Object))
	// Java String.trim() removes leading and trailing characters <= \u0020
	trimmed := strings.TrimFunc(str, func(r rune) bool {
		return r <= '\u0020'
	})
	obj := object.StringObjectFromGoString(trimmed)
	return obj
}

// "java/lang/String.valueOf(Z)Ljava/lang/String;"
func valueOfBoolean(params []interface{}) interface{} {
	// params[0]: input boolean
	value := params[0].(int64)
	var str string
	if value != 0 {
		str = "true"
	} else {
		str = "false"
	}
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.valueOf(C)Ljava/lang/String;"
func valueOfChar(params []interface{}) interface{} {
	// params[0]: input char
	value := params[0].(int64)
	str := fmt.Sprintf("%c", value)
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.valueOf([C)Ljava/lang/String;"
func valueOfCharArray(params []interface{}) interface{} {
	// params[0]: input char array
	propObj := params[0].(*object.Object)
	intArray := propObj.FieldTable["value"].Fvalue.([]int64)
	var str string
	for _, ch := range intArray {
		str += fmt.Sprintf("%c", ch)
	}
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.valueOf([CII)Ljava/lang/String;"
func valueOfCharSubarray(params []interface{}) interface{} {
	// params[0]: input char array
	// params[1]: input offset
	// params[2]: input count
	propObj := params[0].(*object.Object)
	intArray := propObj.FieldTable["value"].Fvalue.([]int64)
	var wholeString string
	for _, ch := range intArray {
		wholeString += fmt.Sprintf("%c", ch)
	}
	// Get substring offset and count
	ssOffset := params[1].(int64)
	ssCount := params[2].(int64)

	// Validate boundaries.
	wholeLength := int64(len(wholeString))
	if wholeLength < 1 || ssOffset < 0 || ssCount < 1 || ssOffset > (wholeLength-1) || (ssOffset+ssCount) > wholeLength {
		errMsg := "valueOfCharSubarray: Either nil input byte array, invalid substring offset, or invalid substring length"
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
	}

	// Compute substring.
	str := wholeString[ssOffset : ssOffset+ssCount]

	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.valueOf(D)Ljava/lang/String;"
func valueOfDouble(params []interface{}) interface{} {
	// params[0]: input double
	value := params[0].(float64)
	str := strconv.FormatFloat(value, 'f', -1, 64)
	if !strings.Contains(str, ".") {
		str += ".0"
	}
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.valueOf(F)Ljava/lang/String;"
func valueOfFloat(params []interface{}) interface{} {
	// params[0]: input float
	value := params[0].(float64)
	str := strconv.FormatFloat(value, 'f', -1, 64)
	if !strings.Contains(str, ".") {
		str += ".0"
	}
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.valueOf(I)Ljava/lang/String;"
func valueOfInt(params []interface{}) interface{} {
	// params[0]: input int
	value := params[0].(int64)
	str := fmt.Sprintf("%d", value)
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.valueOf(J)Ljava/lang/String;"
func valueOfLong(params []interface{}) interface{} {
	// params[0]: input long
	value := params[0].(int64)
	str := fmt.Sprintf("%d", value)
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.valueOf(Ljava/lang/Object;)Ljava/lang/String;"
func valueOfObject(params []interface{}) interface{} {
	// params[0]: input Object or primitive
	var str string

	switch params[0].(type) {
	case *object.Object:
		inObj := params[0].(*object.Object)
		str = object.ObjectFieldToString(inObj, "value")
		if str == types.NullString {
			str = object.ObjectFieldToString(inObj, "name")
		}
	default:
		errMsg := fmt.Sprintf("valueOfObject: Unsupported parameter type: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	outObj := object.StringObjectFromGoString(str)
	return outObj
}

// "java/lang/String.intern()Ljava/lang/String;"
/**
* Returns a canonical representation for the string object.
*
* A pool of strings, initially empty, is maintained privately by the
* class {@code String}.
*
* When the intern method is invoked, if the pool already contains a
* string equal to this {@code String} object as determined by
* the {@link #equals(Object)} method, then the string from the pool is
* returned. Otherwise, this {@code String} object is added to the
* pool and a reference to this {@code String} object is returned.
*
* It follows that for any two strings {@code s} and {@code t},
* {@code s.intern() == t.intern()} is {@code true}
* if and only if {@code s.equals(t)} is {@code true}.
*
* All literal strings and string-valued constant expressions are
* interned. String literals are defined in section {@jls 3.10.5} of the
* The Java Language Specification.
*
* @return  a string that has the same contents as this string, but is
*          guaranteed to be from a pool of unique strings.

public native String intern();
*/
func stringIntern(params []interface{}) interface{} {
	// params[0]: String object
	// TODO: Need to add this to the String pool?
	obj := params[0].(*object.Object)
	return obj
}

// "java/lang/String.checkBoundsBeginEnd(III)V"
func stringCheckBoundsBeginEnd(params []interface{}) interface{} {
	begin := params[0].(int64)
	end := params[1].(int64)
	length := params[2].(int64)

	if begin < 0 || begin > end || end > length {
		errMsg := fmt.Sprintf("stringCheckBoundsBeginEnd: begin: %d, end: %d, length: %d", begin, end, length)
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
	}

	return nil
}

// "java/lang/String.checkBoundsOffCount(III)I"
func stringCheckBoundsOffCount(params []interface{}) interface{} {
	offset := params[0].(int64)
	count := params[1].(int64)
	length := params[2].(int64)

	if offset < 0 || count < 0 || offset > count || offset > (length-count) {
		errMsg := fmt.Sprintf("stringCheckBoundsOffCount: offset: %d, count: %d, length: %d", offset, count, length)
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
	}

	return offset
}

// "java/lang/String.hashCode()I"
func stringHashCode(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringHashCode: null parameter")
	}
	obj := params[0].(*object.Object)
	str := object.GoStringFromStringObject(obj)
	hash := int32(0)
	for _, wint32 := range str {
		hash = 31*hash + int32(wint32)
	}
	return int64(hash)
}

// "java/lang/String.startsWith(Ljava/lang/String;)Z"
// "java/lang/String.startsWith(Ljava/lang/String;I)Z"
func stringStartsWith(params []interface{}) interface{} {
	if params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringStartsWith: null parameter")
	}
	baseObj := params[0].(*object.Object)
	baseStr := object.GoStringFromStringObject(baseObj)
	argObj := params[1].(*object.Object)
	prefix := object.GoStringFromStringObject(argObj)
	if len(params) == 3 {
		offset := int(params[2].(int64))
		if offset < 0 || offset > len(baseStr) {
			errMsg := fmt.Sprintf("stringStartsWith: base: %s, prefix: %s, offset: %d", baseStr, prefix, offset)
			return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
		}
		if strings.HasPrefix(baseStr[offset:], prefix) {
			return types.JavaBoolTrue
		}
		return types.JavaBoolFalse
	}
	if strings.HasPrefix(baseStr, prefix) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// "java/lang/String.endsWith(Ljava/lang/String;)Z"
func stringEndsWith(params []interface{}) interface{} {
	if params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringEndsWith: null parameter")
	}
	baseObj := params[0].(*object.Object)
	baseStr := object.GoStringFromStringObject(baseObj)
	argObj := params[1].(*object.Object)
	prefix := object.GoStringFromStringObject(argObj)
	if strings.HasSuffix(baseStr, prefix) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

/*
void getChars(int srcBegin, int srcEnd, char[] dst, int dstBegin)

Copies characters from the base string into the destination character array.
* The first character to be copied is at index srcBegin.
* The last character to be copied is at index = srcEnd - 1.
* The total number of characters to be copied = srcEnd - srcBegin.
* The characters are copied into the subarray of dst starting at index dstBegin and ending at index = dstBegin + (srcEnd - srcBegin) - 1.
*
* NOTE: The API says that getChars can throw an IndexOutOfBoundsException
*       but the JVM actually throws a StringIndexOutOfBoundsException.
*/

func stringGetChars(params []interface{}) interface{} {
	// params[0] = base object (the string)
	// params[1] = srcBegin
	// params[2] = srcEnd
	// params[3] = object holding the char array
	// params[4] = dstBegin
	// Return nil
	srcFld, ok := params[0].(*object.Object).FieldTable["value"]
	if !ok {
		errMsg := fmt.Sprintf("stringGetChars: Missing value field in base object")
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	var srcBytes []types.JavaByte
	switch srcFld.Fvalue.(type) {
	case []byte:
		srcBytes = object.JavaByteArrayFromGoByteArray(srcFld.Fvalue.([]byte))
	case []types.JavaByte:
		srcBytes = srcFld.Fvalue.([]types.JavaByte)
	default:
		errMsg := fmt.Sprintf("stringGetChars: Invalid value field type (%s : %T) in base object", srcFld.Ftype, srcFld.Fvalue)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get source substring start offset and end offset.
	srcBegin := params[1].(int64)
	srcEnd := params[2].(int64)

	// Compute total length of base byte array.
	srcLength := int64(len(srcBytes))

	// Get destination char array.
	dstObj := params[3].(*object.Object)
	dstFld, ok := dstObj.FieldTable["value"]
	if !ok {
		errMsg := fmt.Sprintf("stringGetChars: Missing value field in char array object")
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	dstChars, ok := dstFld.Fvalue.([]int64)
	if !ok {
		errMsg := fmt.Sprintf("stringGetChars: Invalid value field type (%s : %T) in char array object",
			dstFld.Ftype, dstFld.Fvalue)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get char array start offset.
	dstBegin := params[4].(int64)

	// Compute chara array length.
	dstLength := int64(len(dstChars))

	// Validate boundaries.
	if srcBegin < 0 || srcEnd < srcBegin || srcEnd > srcLength {
		errMsg := fmt.Sprintf("stringGetChars: Source index out of bounds: srcBegin=%d, srcEnd=%d, length=%d", srcBegin, srcEnd, srcLength)
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
	}
	if dstBegin < 0 || dstBegin+(srcEnd-srcBegin) > dstLength {
		errMsg := fmt.Sprintf("stringGetChars: Destination index out of bounds: dstBegin=%d, count=%d, length=%d", dstBegin, srcEnd-srcBegin, dstLength)
		return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, errMsg)
	}

	// Update destination character array.
	ix := dstBegin
	for _, xbyte := range srcBytes[srcBegin:srcEnd] {
		dstChars[ix] = int64(xbyte)
		ix += 1
	}
	dstFld.Fvalue = dstChars
	dstObj.FieldTable["value"] = dstFld

	return nil

}

/*
Returns the index within this string of the first occurrence of the specified character,
starting the search at beginIndex if specified else 0,
and stopping before endIndex if specified else the length of the base string.
*/
func stringIndexOfCh(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringIndexOfCh: null parameter")
	}
	// Get field of base object.
	srcFld, ok := params[0].(*object.Object).FieldTable["value"]
	if !ok {
		errMsg := fmt.Sprintf("stringIndexOfCh: Missing value field in base object")
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get base object byte array.
	var srcBytes []types.JavaByte
	switch srcFld.Fvalue.(type) {
	case []byte:
		srcBytes = object.JavaByteArrayFromGoByteArray(srcFld.Fvalue.([]byte))
	case []types.JavaByte:
		srcBytes = srcFld.Fvalue.([]types.JavaByte)
	default:
		errMsg := fmt.Sprintf("stringIndexOfCh: Invalid value field type (%s : %T) in base object", srcFld.Ftype, srcFld.Fvalue)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get search argument and set up switch.
	arg := types.JavaByte(params[1].(int64))
	var beginIndex int64
	var endIndex int64
	lenSrcBytes := int64(len(srcBytes))

	// There are 3 slightly different functions requested.
	switch len(params) - 1 {
	case 1: // int indexOf(int ch)
		beginIndex = 0
		endIndex = lenSrcBytes
	case 2: // int indexOf(int ch, int fromIndex)
		beginIndex = params[2].(int64)
		if beginIndex < 0 {
			beginIndex = 0
		}
		if beginIndex >= lenSrcBytes {
			return int64(-1)
		}
		endIndex = lenSrcBytes
	case 3: // int indexOf(int ch, int beginIndex, int endIndex)
		beginIndex = params[2].(int64)
		if beginIndex < 0 || beginIndex >= lenSrcBytes {
			errMsg := fmt.Sprintf("stringIndexOfCh: Base string len: %d, begin index: %d", lenSrcBytes, beginIndex)
			return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
		}
		endIndex = params[3].(int64)
		if endIndex > lenSrcBytes || beginIndex > endIndex {
			errMsg := fmt.Sprintf("stringIndexOfCh: Base string len: %d, end index: %d", lenSrcBytes, endIndex)
			return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
		}
	}

	// Search for argument in base byte array.
	for ix := beginIndex; ix < endIndex; ix++ {
		if arg == srcBytes[ix] {
			return ix // Found it. Return the index.
		}
	}
	return int64(-1) // Did not find it.
}

func stringIndexOfString(params []interface{}) interface{} {
	if params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringIndexOfString: null parameter")
	}
	// Get field of base object.
	baseString := object.GoStringFromStringObject(params[0].(*object.Object))

	// Get base object byte array.
	argString := object.GoStringFromStringObject(params[1].(*object.Object))

	// Set up for switch.
	lenOrigBaseString := int64(len(baseString))
	var beginIndex int64
	var endIndex int64

	// There are 3 slightly different functions requested.
	switch len(params) - 1 {
	case 1: // int indexOf(String str)
		beginIndex = 0
		endIndex = lenOrigBaseString
	case 2: // int indexOf(String str, int fromIndex)
		beginIndex = params[2].(int64)
		if beginIndex < 0 {
			beginIndex = 0
		}
		if beginIndex >= lenOrigBaseString {
			return int64(-1)
		}
		endIndex = lenOrigBaseString
	case 3: // int indexOf(String str, int beginIndex, int endIndex)
		beginIndex = params[2].(int64)
		if beginIndex < 0 || beginIndex >= lenOrigBaseString {
			errMsg := fmt.Sprintf("stringIndexOfString: Base string len: %d, begin index: %d", lenOrigBaseString, beginIndex)
			return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
		}
		endIndex = params[3].(int64)
		if endIndex > lenOrigBaseString || beginIndex > endIndex {
			errMsg := fmt.Sprintf("stringIndexOfString: Base string len: %d, end index: %d", lenOrigBaseString, endIndex)
			return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
		}
	}

	// Search for argument in subsetted base string.
	ii := int64(strings.Index(baseString[beginIndex:endIndex], argString))
	if ii > 0 {
		ii = ii + beginIndex // relative to original base string
	}
	return ii // >= 0 if success, -1 if failure
}

func stringIsBlank(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringIsBlank: null parameter")
	}
	baseString := object.GoStringFromStringObject(params[0].(*object.Object))
	if len(strings.TrimSpace(baseString)) == 0 {
		return types.JavaBoolTrue
	} else {
		return types.JavaBoolFalse
	}
}

func stringIsEmpty(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringIsEmpty: null parameter")
	}
	baseString := object.GoStringFromStringObject(params[0].(*object.Object))
	if len(baseString) == 0 {
		return types.JavaBoolTrue
	} else {
		return types.JavaBoolFalse
	}
}

func stringToString(params []interface{}) interface{} {
	return params[0]
}

func stringReplaceLiteral(params []interface{}) interface{} {
	// If any parameters are missing or nil, throw a NullPointerException.
	if params[0] == nil || params[1] == nil || params[2] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringReplaceLiteral: null parameter")
	}

	// Get 3 arguments. These are objects that implement CharSequence.
	// In Jacobin, they should all be string objects or similar.
	input := object.GoStringFromStringObject(params[0].(*object.Object))
	target := object.GoStringFromStringObject(params[1].(*object.Object))
	replacement := object.GoStringFromStringObject(params[2].(*object.Object))

	// Replace all literal occurrences of target with replacement.
	result := strings.ReplaceAll(input, target, replacement)

	return object.StringObjectFromGoString(result)
}

func stringReplaceAllRegex(params []interface{}) interface{} {
	// If any parameters are missing or nil, throw a NullPointerException.
	if params[0] == nil || params[1] == nil || params[2] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringReplaceAllRegex: null parameter")
	}

	// Get 3 string arguments.
	input := object.GoStringFromStringObject(params[0].(*object.Object))
	pattern := object.GoStringFromStringObject(params[1].(*object.Object))
	replacement := object.GoStringFromStringObject(params[2].(*object.Object))

	// Compile the regular expression.
	re, err := regexp.Compile(pattern)
	if err != nil {
		errMsg := fmt.Sprintf("stringReplaceAllRegex: Invalid regular expression pattern: %s", pattern)
		return ghelpers.GetGErrBlk(excNames.PatternSyntaxException, errMsg)
	}

	// Replace all substrings that match the pattern with the replacement string.
	result := re.ReplaceAllString(input, replacement)

	return object.StringObjectFromGoString(result)

}

func stringReplaceFirstRegex(params []interface{}) interface{} {
	// If any parameters are missing or nil, throw a NullPointerException.
	if params[0] == nil || params[1] == nil || params[2] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringReplaceFirstRegex: null parameter")
	}

	// Get 3 string arguments.
	input := object.GoStringFromStringObject(params[0].(*object.Object))
	pattern := object.GoStringFromStringObject(params[1].(*object.Object))
	replacement := object.GoStringFromStringObject(params[2].(*object.Object))

	// Compile the regular expression.
	re, err := regexp.Compile(pattern)
	if err != nil {
		errMsg := fmt.Sprintf("stringReplaceFirstRegex: Invalid regular expression pattern: %s", pattern)
		return ghelpers.GetGErrBlk(excNames.PatternSyntaxException, errMsg)
	}

	// Find the first match for the regular expression.
	loc := re.FindStringIndex(input)
	if loc == nil {
		// No match found, return the original input string.
		return params[0]
	}

	return object.StringObjectFromGoString(input[:loc[0]] + replacement + input[loc[1]:])

}

func stringSplit(params []interface{}) interface{} {
	// If any parameters are missing or nil, throw a NullPointerException.
	if params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringSplit: null parameter")
	}

	// params[0] = base string
	// params[1] = regular expression in a string

	input := object.GoStringFromStringObject(params[0].(*object.Object))
	pattern := object.GoStringFromStringObject(params[1].(*object.Object))

	// Compile the regular expression.
	re, err := regexp.Compile(pattern)
	if err != nil {
		errMsg := fmt.Sprintf("stringSplit: Invalid regular expression pattern: %s", pattern)
		return ghelpers.GetGErrBlk(excNames.PatternSyntaxException, errMsg)
	}

	// Split input based on the pattern.
	result := re.Split(input, -1) // -1 means split on all occurrences.

	// Prepare object array and return it.
	var outObjArray []*object.Object
	for ix := 0; ix < len(result); ix++ {
		outObjArray = append(outObjArray, object.StringObjectFromGoString(result[ix]))
	}
	return object.MakePrimitiveObject("[Ljava/lang/String;", types.RefArray, outObjArray)

}

func stringSplitLimit(params []interface{}) interface{} {
	// If any parameters are missing or nil, throw a NullPointerException.
	if params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringSplitLimit: null parameter")
	}

	// params[0] = base string
	// params[1] = regular expression in a string
	// params[2] = split limit

	input := object.GoStringFromStringObject(params[0].(*object.Object))
	pattern := object.GoStringFromStringObject(params[1].(*object.Object))
	limit := params[2].(int64)
	if limit == 0 {
		limit = -1
	}

	// Compile the regular expression.
	re, err := regexp.Compile(pattern)
	if err != nil {
		errMsg := fmt.Sprintf("stringSplitLimit: Invalid regular expression pattern: %s", pattern)
		return ghelpers.GetGErrBlk(excNames.PatternSyntaxException, errMsg)
	}

	// Split input based on the pattern.
	result := re.Split(input, int(limit))

	// Prepare object array and return it.
	var outObjArray []*object.Object
	for ix := 0; ix < len(result); ix++ {
		outObjArray = append(outObjArray, object.StringObjectFromGoString(result[ix]))
	}
	return object.MakePrimitiveObject("[Ljava/lang/String;", types.RefArray, outObjArray)
}

func stringStrip(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringStrip: null parameter")
	}
	input := object.GoStringFromStringObject(params[0].(*object.Object))
	result := strings.TrimSpace(input)
	return object.StringObjectFromGoString(result)
}

func stringStripLeading(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringStripLeading: null parameter")
	}
	input := object.GoStringFromStringObject(params[0].(*object.Object))
	result := strings.TrimLeftFunc(input, unicode.IsSpace)
	return object.StringObjectFromGoString(result)
}

func stringStripTrailing(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "stringStripTrailing: null parameter")
	}
	input := object.GoStringFromStringObject(params[0].(*object.Object))
	result := strings.TrimRightFunc(input, unicode.IsSpace)
	return object.StringObjectFromGoString(result)
}
