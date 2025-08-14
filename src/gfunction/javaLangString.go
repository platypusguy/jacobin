/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-4 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
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
	MethodSignatures["java/lang/String.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringClinit,
		}

	// Instantiate an empty String
	MethodSignatures["java/lang/String.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  newEmptyString,
		}

	// String(byte[] bytes) - instantiate a String from a byte array
	MethodSignatures["java/lang/String.<init>([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  newStringFromBytes,
		}

	// String(byte[] ascii, int hibyte) *** DEPRECATED
	MethodSignatures["java/lang/String.<init>([BI)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapDeprecated,
		}

	// String(byte[] bytes, int offset, int length)	- instantiate a String from a subset of a byte array
	MethodSignatures["java/lang/String.<init>([BII)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  newStringFromBytesSubset,
		}

	// String(byte[] ascii, int hibyte, int offset, int count) *** DEPRECATED
	MethodSignatures["java/lang/String.<init>([BIII)V"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  trapDeprecated,
		}

	// TODO: String(byte[] bytes, int offset, int length, String charsetName) *********** CHARSET
	MethodSignatures["java/lang/String.<init>([BIILjava/lang/String;)V"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  trapFunction,
		}

	// TODO: String(byte[] bytes, int offset, int length, Charset charset) ************** CHARSET
	MethodSignatures["java/lang/String.<init>([BIILjava/nio/charset/Charset;)V"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  trapFunction,
		}

	// TODO: String(byte[] bytes, String charsetName) *********************************** CHARSET
	MethodSignatures["java/lang/String.<init>([BLjava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	// TODO: String(byte[] bytes, Charset charset) ************************************** CHARSET
	MethodSignatures["java/lang/String.<init>([BLjava/nio/charset/Charset;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	// Instantiate a String from a character array
	MethodSignatures["java/lang/String.<init>([C)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  newStringFromChars,
		}

	// Instantiate a String from a subset of a character array
	MethodSignatures["java/lang/String.<init>([CII)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  newStringFromChars,
		}

	// TODO: String(int[] codePoints, int offset, int count) ************************ CODEPOINTS
	MethodSignatures["java/lang/String.<init>([III)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	// String(String original) -- instantiate a String from another String.
	MethodSignatures["java/lang/String.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  newStringFromString,
		}

	// String(StringBuffer buffer) -- instantiate a String from a StringBuffer.
	MethodSignatures["java/lang/String.<init>(Ljava/lang/StringBuffer;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  newStringFromString,
		}

	// Not in API: Is the String Latin1?
	MethodSignatures["java/lang/String.isLatin1()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringIsLatin1,
		}

	// String(StringBuilder builder) -- instantiate a String from a StringBuilder.
	MethodSignatures["java/lang/String.<init>(Ljava/lang/StringBuilder;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  newStringFromString,
		}

	// ==== METHOD FUNCTIONS (in alphabetical order by their function names) ====

	// Returns the char value at the specified index.
	MethodSignatures["java/lang/String.charAt(I)C"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringCharAt,
		}

	// TODO: Returns a stream of int zero-extending the char values from this sequence.
	MethodSignatures["java/lang/String.chars()Ljava/util/stream/IntStream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	// Internal boundary-checker - not in the API.
	MethodSignatures["java/lang/String.checkBoundsBeginEnd(III)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  stringCheckBoundsBeginEnd,
		}

	// Internal boundary-checker - not in the API.
	MethodSignatures["java/lang/String.checkBoundsOffCount(III)I"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  stringCheckBoundsOffCount,
		}

	// TODO: Returns the character (Unicode code point) at the specified index.
	MethodSignatures["java/lang/String.codePointAt(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	// TODO: Returns the character (Unicode code point) before the specified index.
	MethodSignatures["java/lang/String.codePointBefore(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	// TODO: Returns the number of Unicode code points in the specified text range of this String.
	MethodSignatures["java/lang/String.codePointCount(II)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	// TODO: Returns a stream of code point values from this sequence.
	MethodSignatures["java/lang/String.codePoints()Ljava/util/stream/IntStream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	// Compare 2 strings lexicographically, case-sensitive (upper/lower).
	// The return value is a negative integer, zero, or a positive integer
	// as the String argument is greater than, equal to, or less than this String,
	// case-sensitive.
	MethodSignatures["java/lang/String.compareTo(Ljava/lang/String;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringCompareToCaseSensitive,
		}

	// Compare 2 strings lexicographically, ignoring case (upper/lower).
	// The return value is a negative integer, zero, or a positive integer
	// as the String argument is greater than, equal to, or less than this String,
	// ignoring case considerations.
	MethodSignatures["java/lang/String.compareToIgnoreCase(Ljava/lang/String;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringCompareToIgnoreCase,
		}

	// Concatenates the specified string to the end of this string.
	MethodSignatures["java/lang/String.concat(Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringConcat,
		}

	// Returns true if and only if this string contains the specified sequence of char values.
	MethodSignatures["java/lang/String.contains(Ljava/lang/CharSequence;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringContains,
		}

	// Compares this string to the specified CharSequence.
	MethodSignatures["java/lang/String.contentEquals(Ljava/lang/CharSequence;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  javaLangStringContentEquals,
		}

	// Compares this string to the specified StringBuffer.
	MethodSignatures["java/lang/String.contentEquals(Ljava/lang/StringBuffer;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  javaLangStringContentEquals,
		}

	// Return a string representing a char array.
	MethodSignatures["java/lang/String.copyValueOf([C)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  valueOfCharArray,
		}

	// Return a string representing a char subarray.
	MethodSignatures["java/lang/String.copyValueOf([CII)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  newStringFromCharsSubset,
		}

	// TODO: Returns an Optional containing the nominal descriptor for this instance, which is the instance itself.
	MethodSignatures["java/lang/String.describeConstable()Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	// OpenJDK JVM "java/lang/String.endsWith(Ljava/lang/String;)Z" works with the jacobin String object.
	// Does the base string end with the specified suffix argument?

	// Compares this string to the specified object.
	MethodSignatures["java/lang/String.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringEquals,
		}

	// Compares this String to another String, ignoring case considerations.
	MethodSignatures["java/lang/String.equalsIgnoreCase(Ljava/lang/String;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringEqualsIgnoreCase,
		}

	// Return a formatted string using the reference object string as the format string
	// and the supplied arguments as input object arguments.
	// E.g. String string = String.format("%s %i", "ABC", 42);
	MethodSignatures["java/lang/String.format(Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  sprintf,
		}

	// TODO: Return a formatted string using the specified locale, format string, and arguments.
	MethodSignatures["java/lang/String.format(Ljava/util/Locale;Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	// This method is equivalent to String.format(this, args).
	MethodSignatures["java/lang/String.formatted([Ljava/lang/Object;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  sprintf,
		}

	// Encodes this String into a sequence of bytes using the default charset, storing the result into a new byte array.
	MethodSignatures["java/lang/String.getBytes()[B"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  getBytesFromString,
		}

	// void getBytes(int srcBegin, int srcEnd, byte[] dst, int dstBegin)  ********************* DEPRECATED
	MethodSignatures["java/lang/String.getBytes(II[BI)V"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  trapDeprecated,
		}

	// TODO: Encodes this String into a sequence of bytes using the given charset, storing the result into a new byte array.
	MethodSignatures["java/lang/String.getBytes(Ljava/nio/charset/Charset;)[B"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	// Encodes this String into a sequence of bytes using the named charset, storing the result into a new byte array. ************************ CHARSET
	MethodSignatures["java/lang/String.getBytes(Ljava/lang/String;)[B"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	// Not in API: getBytes([BIIBI)V
	// original Java source: https://gist.github.com/platypusguy/03c1a9e3acb1cb2cfc2d821aa2dd4490
	MethodSignatures["java/lang/String.getBytes([BIIBI)V"] =
		GMeth{
			ParamSlots: 5,
			GFunction:  stringGetBytesBIIBI,
		}

	// Copies characters from this string into the destination character array.
	MethodSignatures["java/lang/String.getChars(II[CI)V"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  stringGetChars,
		}

	// Compute the Java String.hashCode() value.
	MethodSignatures["java/lang/String.hashCode()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringHashCode,
		}

	// TODO: Adjusts the indentation of each line of this string based on the value of n, and normalizes line termination characters.
	MethodSignatures["java/lang/String.indent(I)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	// Returns the index within this string of the first occurrence of the specified character.
	MethodSignatures["java/lang/String.indexOf(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringIndexOfCh,
		}

	// Returns the index within this string of the first occurrence of the specified character,
	// starting the search at the specified index.
	MethodSignatures["java/lang/String.indexOf(II)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringIndexOfCh,
		}

	// Returns the index within this string of the first occurrence of the specified character,
	// starting the search at beginIndex and stopping before endIndex.
	MethodSignatures["java/lang/String.indexOf(III)I"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  stringIndexOfCh,
		}

	// Returns the index within this string of the first occurrence of the specified substring.
	MethodSignatures["java/lang/String.indexOf(Ljava/lang/String;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringIndexOfString,
		}

	// Returns the index within this string of the first occurrence of the specified substring,
	// starting at the specified index.
	MethodSignatures["java/lang/String.indexOf(Ljava/lang/String;I)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringIndexOfString,
		}

	// Returns the index of the first occurrence of the specified substring within the specified index range of this string.
	MethodSignatures["java/lang/String.indexOf(Ljava/lang/String;II)I"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  stringIndexOfString,
		}

	/*
		When the intern method is invoked, if the pool already contains a string equal to this String object as determined
		by the equals(Object) method, then the string from the pool is returned.
		Otherwise, this String object is added to the pool and a reference to this String object is returned.
	*/
	MethodSignatures["java/lang/String.intern()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringIntern,
		}

	// Is the base string whitespace?
	MethodSignatures["java/lang/String.isBlank()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringIsBlank,
		}

	// Is the base string empty?
	MethodSignatures["java/lang/String.isEmpty()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringIsEmpty,
		}

	// TODO: Returns a new String composed of copies of the CharSequence elements joined together with a copy of the specified delimiter.
	MethodSignatures["java/lang/String.join(Ljava/lang/CharSequence;[Ljava/lang/CharSequence;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	// TODO: Returns a new String composed of copies of the CharSequence elements joined together with a copy of the specified delimiter.
	MethodSignatures["java/lang/String.join(Ljava/lang/CharSequence;[Ljava/lang/Iterable;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	// Returns the index within this string of the last occurrence of the specified character.
	MethodSignatures["java/lang/String.lastIndexOf(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  lastIndexOfCharacter,
		}

	// Returns the index within this string of the last occurrence of the specified character, searching backward starting at the specified index.
	MethodSignatures["java/lang/String.lastIndexOf(II)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  lastIndexOfCharacter,
		}

	// Returns the index within this string of the last occurrence of the specified substring.
	MethodSignatures["java/lang/String.lastIndexOf(Ljava/lang/String;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  lastIndexOfString,
		}

	// Returns the index within this string of the last occurrence of the specified substring, searching backward starting at the specified index.
	MethodSignatures["java/lang/String.lastIndexOf(Ljava/lang/String;I)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  lastIndexOfString,
		}

	// Return the length of the base String.
	MethodSignatures["java/lang/String.length()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringLength,
		}

	// TODO: Returns a stream of lines extracted from this string, separated by line terminators.
	MethodSignatures["java/lang/String.lines()Ljava/util/stream/Stream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	// Tells whether this string matches the given regular expression or not.
	MethodSignatures["java/lang/String.matches(Ljava/lang/String;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringMatches,
		}

	// TODO: Returns the index within this String that is offset from the given index by codePointOffset code points.
	MethodSignatures["java/lang/String.offsetByCodePoints(II)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	// Tests if two string regions are equal.
	// Pass a flag indicating whether to ignore case or not.
	MethodSignatures["java/lang/String.regionMatches(ZILjava/lang/String;II)Z"] = // Has an ignoreCase flag
		GMeth{
			ParamSlots: 5,
			GFunction:  stringRegionMatches,
		}

	// Tests if two string regions are equal, case-sensitive.
	MethodSignatures["java/lang/String.regionMatches(ILjava/lang/String;II)Z"] = // Does not have an ignoreCase flag
		GMeth{
			ParamSlots: 4,
			GFunction:  stringRegionMatches,
		}

	// Returns a string whose value is the concatenation of this string repeated the specified number of times.
	MethodSignatures["java/lang/String.repeat(I)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringRepeat,
		}

	// Replace a single character by another in the given string.
	MethodSignatures["java/lang/String.replace(CC)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringReplaceCC,
		}

	// Replace a character sequence by another in the given string.
	MethodSignatures["java/lang/String.replace(Ljava/lang/CharSequence;Ljava/lang/CharSequence;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringReplaceAllRegex,
		}

	// Replaces each substring of this string that matches the given regular expression with the given replacement.
	MethodSignatures["java/lang/String.replaceAll(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringReplaceAllRegex,
		}

	// Replaces the first substring of this string that matches the given regular expression with the given replacement.
	MethodSignatures["java/lang/String.replaceFirst(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringReplaceFirstRegex,
		}

	// TODO: Resolves this instance as a ConstantDesc, the result of which is the instance itself.
	MethodSignatures["java/lang/String.resolveConstantDesc(Ljava/lang/invoke/MethodHandles/Lookup;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	// Split the base string into an array of strings.
	MethodSignatures["java/lang/String.split(Ljava/lang/String;)[Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringSplit,
		}

	// Split the base string into an array of strings with a specified limit.
	MethodSignatures["java/lang/String.split(Ljava/lang/String;I)[Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringSplitLimit,
		}

	// TODO: Split the base string around matches of the given regular expression and returns both the strings and the matching delimiters.
	MethodSignatures["java/lang/String.splitWithDelimiters(Ljava/lang/String;I)[Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	// Tests if this string starts with the specified prefix.
	MethodSignatures["java/lang/String.startsWith(Ljava/lang/String;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringStartsWith,
		}

	// Tests if the substring of this string beginning at the specified index starts with the specified prefix.
	MethodSignatures["java/lang/String.startsWith(Ljava/lang/String;I)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringStartsWith,
		}

	// Returns a string whose value is this string, with all leading and trailing white space removed.
	MethodSignatures["java/lang/String.strip()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringStrip,
		}

	// TODO: Returns a string whose value is this string, with incidental white space removed from the beginning and end of every line.
	MethodSignatures["java/lang/String.stripIndent()[Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	// Returns a string whose value is the base string with all leading white space removed.
	MethodSignatures["java/lang/String.stripLeading()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringStripLeading,
		}

	// Returns a string whose value is the base string with all trailing white space removed.
	MethodSignatures["java/lang/String.stripTrailing()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringStripTrailing,
		}

	// TODO: Returns a character sequence that is a subsequence of this sequence.
	MethodSignatures["java/lang/String.subSequence(II)Ljava/lang/CharSequence;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	// Return a substring starting at the given index of the byte array.
	MethodSignatures["java/lang/String.substring(I)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  substringToTheEnd,
		}

	// Return a substring starting at the given index of the byte array of the given length.
	MethodSignatures["java/lang/String.substring(II)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  substringStartEnd,
		}

	// Return a string in all lower case, using the reference object string as input.
	MethodSignatures["java/lang/String.toCharArray()[C"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  toCharArray,
		}

	// Return a string in all lower case, using the reference object string as input.
	MethodSignatures["java/lang/String.toLowerCase()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  toLowerCase,
		}

	// TODO: Converts all of the characters in this String to lower case using the rules of the given Locale.
	MethodSignatures["java/lang/String.toLowerCase(Ljava/util/Locale;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  toLowerCase, // TODO: Locale processing
		}

	// Return the base string as-is.
	MethodSignatures["java/lang/String.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringToString,
		}

	// Return a string in all upper case, using the reference object string as input.
	MethodSignatures["java/lang/String.toUpperCase()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  toUpperCase,
		}

	// TODO: Converts all of the characters in this String to upper case using the rules of the given Locale.
	MethodSignatures["java/lang/String.toUpperCase(Ljava/util/Locale;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	// TODO: What should we do with this? <R> R transform(Function<? super String,? extends R> f)

	// TODO: Return a string whose value is the base string with escape sequences translated as if in a string literal.
	MethodSignatures["java/lang/String.translateEscapes()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	// Return a string trimmed of leading and trailing whitespace.
	MethodSignatures["java/lang/String.trim()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trimString,
		}

	// Return a string representing a boolean value.
	MethodSignatures["java/lang/String.valueOf(Z)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  valueOfBoolean,
		}

	// Return a string representing a char value.
	MethodSignatures["java/lang/String.valueOf(C)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  valueOfChar,
		}

	// Return a string representing a char array.
	MethodSignatures["java/lang/String.valueOf([C)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  valueOfCharArray,
		}

	// Return a string representing a char subarray.
	MethodSignatures["java/lang/String.valueOf([CII)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  newStringFromCharsSubset,
		}

	// Return a string representing a double value.
	MethodSignatures["java/lang/String.valueOf(D)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  valueOfDouble,
		}

	// Return a string representing a float value.
	MethodSignatures["java/lang/String.valueOf(F)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  valueOfFloat,
		}

	// Return a string representing an int value.
	MethodSignatures["java/lang/String.valueOf(I)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  valueOfInt,
		}

	// Return a string representing an int value.
	MethodSignatures["java/lang/String.valueOf(J)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  valueOfLong,
		}

	// Return a string representing the value of an Object.
	MethodSignatures["java/lang/String.valueOf(Ljava/lang/Object;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  valueOfObject,
		}

}

// ==== INSTANTIATION AND INITIALIZATION FUNCTIONS ====

// "java/lang/String.<clinit>()V" -- String class initialisation
func stringClinit([]interface{}) interface{} {
	klass := classloader.MethAreaFetch(types.StringClassName)
	if klass == nil {
		errMsg := fmt.Sprintf("stringClinit: Could not find class %s in the MethodArea", types.StringClassName)
		return getGErrBlk(excNames.ClassNotLoadedException, errMsg)
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
	obj := params[0].(*object.Object)
	switch params[1].(*object.Object).FieldTable["value"].Fvalue.(type) {
	case []byte:
		bytes := object.JavaByteArrayFromGoByteArray(
			params[1].(*object.Object).FieldTable["value"].Fvalue.([]byte))
		object.UpdateValueFieldFromJavaBytes(obj, bytes)
	case []types.JavaByte:
		bytes := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
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
	// params[3] = end offset
	obj := params[0].(*object.Object)
	var bytes []types.JavaByte
	switch params[1].(*object.Object).FieldTable["value"].Fvalue.(type) {
	case []byte:
		bytes = object.JavaByteArrayFromGoByteArray(params[1].(*object.Object).FieldTable["value"].Fvalue.([]byte))
	case []types.JavaByte:
		bytes = params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	}

	// Get substring start and end offset
	ssStart := params[2].(int64)
	ssEnd := params[3].(int64)

	// Validate boundaries.
	totalLength := int64(len(bytes))
	if totalLength < 1 || ssStart < 0 || ssEnd < 1 || ssStart > (totalLength-1) || (ssStart+ssEnd) > totalLength {
		errMsg1 := "newStringFromBytesSubset: Either nil input byte array, invalid substring offset, or invalid substring length"
		errMsg2 := fmt.Sprintf("\n\twhole='%s' wholelen=%d, offset=%d, sslen=%d\n\n",
			object.GoStringFromJavaByteArray(bytes), totalLength, ssStart, ssEnd)
		return getGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg1+errMsg2)
	}

	// Compute subarray and update params[0].
	bytes = bytes[ssStart : ssStart+ssEnd]
	object.UpdateValueFieldFromJavaBytes(obj, bytes)
	return nil
}

// Instantiate a new string object from a Go int64 array (Java char array).
// "java/lang/String.<init>([C)V"
func newStringFromChars(params []interface{}) interface{} {
	// params[0] = reference string (to be updated with byte array)
	// params[1] = byte array object
	obj := params[0].(*object.Object)
	ints := params[1].(*object.Object).FieldTable["value"].Fvalue.([]int64)

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
	// params[2] = end offset
	// Return the string.
	fld, ok := params[0].(*object.Object).FieldTable["value"]
	if !ok {
		errMsg := fmt.Sprintf("newStringFromCharsSubset: Missing value field in character array object")
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	iarray, ok := fld.Fvalue.([]int64)
	if !ok {
		errMsg := fmt.Sprintf("newStringFromCharsSubset: Invalid value field type (%s : %T) in character array object", fld.Ftype, fld.Fvalue)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get substring start and end offset
	ssStart := params[1].(int64)
	ssEnd := params[2].(int64)

	// Validate boundaries.
	totalLength := int64(len(iarray))
	if totalLength < 1 || ssStart < 0 || ssEnd < 1 || ssStart > (totalLength-1) || (ssStart+ssEnd) > totalLength {
		errMsg1 := "newStringFromCharsSubset: Either nil input byte array, invalid substring offset, or invalid substring length"
		errMsg2 := fmt.Sprintf("\n\twholelen=%d, offset=%d, sslen=%d\n\n", totalLength, ssStart, ssEnd)
		return getGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg1+errMsg2)
	}

	// Compute subarray and update params[0].
	iarray = iarray[ssStart : ssStart+ssEnd]
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
	// Unpack the reference string and convert it to a rune array.
	ptrObj := params[0].(*object.Object)
	str := object.GoStringFromStringObject(ptrObj)
	runeArray := []rune(str)

	// Get index.
	index := params[1].(int64)

	// Return indexed character.
	runeValue := runeArray[index]
	return int64(runeValue)
}

// "java/lang/String.compareTo(Ljava/lang/String;)I"
func stringCompareToCaseSensitive(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	str1 := object.GoStringFromStringObject(obj)
	obj = params[1].(*object.Object)
	str2 := object.GoStringFromStringObject(obj)
	if str2 == str1 {
		return types.JavaBoolFalse
	}
	if str1 < str2 {
		return int64(-1)
	}
	return types.JavaBoolTrue
}

// "java/lang/String.compareToIgnoreCase(Ljava/lang/String;)I"
func stringCompareToIgnoreCase(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	str1 := strings.ToLower(object.GoStringFromStringObject(obj))
	obj = params[1].(*object.Object)
	str2 := strings.ToLower(object.GoStringFromStringObject(obj))
	if str2 == str1 {
		return int64(0)
	}
	if str1 < str2 {
		return int64(-1)
	}
	return types.JavaBoolTrue
}

// "java/lang/String.concat(Ljava/lang/String;)Ljava/lang/String;"
func stringConcat(params []interface{}) interface{} {
	var str1, str2 string

	fld := params[0].(*object.Object).FieldTable["value"]
	switch fld.Fvalue.(type) {
	case []byte:
		str1 = string(fld.Fvalue.([]byte))
	case []types.JavaByte:
		str1 = object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
	}

	fld = params[1].(*object.Object).FieldTable["value"]
	switch fld.Fvalue.(type) {
	case []byte:
		str2 = string(fld.Fvalue.([]byte))
	case []types.JavaByte:
		str2 = object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
	}

	str := str1 + str2
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.contains(Ljava/lang/CharSequence;)Z"
// charSequence is an interface, generally implemented via String or array of chars
// Here, we assume one of those two options.
func stringContains(params []interface{}) interface{} {
	// get the search string (the string we're searching for, i.e., "foo" in "seafood")
	searchFor := params[1].(*object.Object)
	var searchString string
	switch searchFor.FieldTable["value"].Fvalue.(type) {
	case []types.JavaByte:
		searchString =
			object.GoStringFromJavaByteArray(searchFor.FieldTable["value"].Fvalue.([]types.JavaByte))
	case []uint8:
		searchString = string(searchFor.FieldTable["value"].Fvalue.([]byte))
	case string:
		searchString = searchFor.FieldTable["value"].Fvalue.(string)
	}
	searchIn := params[0].(*object.Object)

	// now get the target string (the string being searched)
	var targetString string
	switch searchIn.FieldTable["value"].Fvalue.(type) {
	case []types.JavaByte:
		targetString =
			object.GoStringFromJavaByteArray(searchIn.FieldTable["value"].Fvalue.([]types.JavaByte))
	case []uint8:
		targetString = string(searchIn.FieldTable["value"].Fvalue.([]byte))
	case string:
		targetString = searchIn.FieldTable["value"].Fvalue.(string)
	}

	if strings.Contains(targetString, searchString) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func javaLangStringContentEquals(params []interface{}) interface{} {
	var str1, str2 string
	obj := params[0].(*object.Object)
	switch obj.FieldTable["value"].Fvalue.(type) {
	case []byte:
		str1 = string(obj.FieldTable["value"].Fvalue.([]byte))
	case []types.JavaByte:
		str1 = object.GoStringFromJavaByteArray(obj.FieldTable["value"].Fvalue.([]types.JavaByte))
	}

	obj = params[1].(*object.Object)
	switch obj.FieldTable["value"].Fvalue.(type) {
	case []byte:
		str2 = string(obj.FieldTable["value"].Fvalue.([]byte))
	case []types.JavaByte:
		str2 = object.GoStringFromJavaByteArray(obj.FieldTable["value"].Fvalue.([]types.JavaByte))
	}

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
	obj := params[0].(*object.Object)
	str1 := object.GoStringFromStringObject(obj)
	obj = params[1].(*object.Object)
	str2 := object.GoStringFromStringObject(obj)

	// Are they equal in value?
	if str1 == str2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Are 2 strings equal, ignoring case?
// "java/lang/String.equalsIgnoreCase(Ljava/lang/String;)Z"
func stringEqualsIgnoreCase(params []interface{}) interface{} {
	var str1, str2 string
	// params[0]: reference string object
	// params[1]: compare-to string Object
	obj := params[0].(*object.Object)
	switch obj.FieldTable["value"].Fvalue.(type) {
	case []byte:
		str1 = string(obj.FieldTable["value"].Fvalue.([]byte))
	case []types.JavaByte:
		str1 = object.GoStringFromJavaByteArray(obj.FieldTable["value"].Fvalue.([]types.JavaByte))
	}

	obj = params[1].(*object.Object)
	switch obj.FieldTable["value"].Fvalue.(type) {
	case []byte:
		str2 = string(obj.FieldTable["value"].Fvalue.([]byte))
	case []types.JavaByte:
		str2 = object.GoStringFromJavaByteArray(obj.FieldTable["value"].Fvalue.([]types.JavaByte))
	}

	// Are they equal in value?
	upstr1 := strings.ToUpper(str1)
	upstr2 := strings.ToUpper(str2)
	if upstr1 == upstr2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// "java/lang/String.format(Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/String;"
// "java/lang/String.formatted([Ljava/lang/Object;)Ljava/lang/String;"
func sprintf(params []interface{}) interface{} {
	// params[0]: format string
	// params[1]: argument slice (array of object pointers)
	return StringFormatter(params)
}

// java/lang/String.getBytes()[B
func getBytesFromString(params []interface{}) interface{} {
	// params[0] = reference string with byte array to be returned
	bytes := object.JavaByteArrayFromStringObject(params[0].(*object.Object))
	return Populator("[B", types.ByteArray, bytes)
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
	obj := params[0].(*object.Object)
	bytes := object.JavaByteArrayFromStringObject(obj)
	return int64(len(bytes))
}

// java/lang/String.matches(Ljava/lang/String;)Z
// is the string in params[0] a match for the regex in params[1]?
func stringMatches(params []any) any {
	if len(params) != 2 {
		errMsg := fmt.Sprintf("stringMatches: Expected a string and a regular expression")
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	baseStringObject := params[0].(*object.Object)
	baseString := object.GoStringFromStringObject(baseStringObject)

	regexStringObject := params[1].(*object.Object)
	regexString := object.GoStringFromStringObject(regexStringObject)

	regex, err := regexp.Compile(regexString)
	if err != nil {
		errMsg := fmt.Sprintf("stringMatches: Invalid regular expression: %s", regexString)
		return getGErrBlk(excNames.PatternSyntaxException, errMsg)
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
	str := object.GoStringFromStringObject(params[0].(*object.Object))

	// Get substring start offset and compute end offset
	ssStart := params[1].(int64)
	ssEnd := int64(len(str))

	// Validate boundaries.
	totalLength := int64(len(str))
	if totalLength < 1 || ssStart < 0 || ssEnd < 1 || ssStart > (totalLength-1) || ssEnd > totalLength {
		errMsg1 := "substringToTheEnd: Either nil input byte array, invalid substring offset, or invalid substring length"
		errMsg2 := fmt.Sprintf("\n\twhole='%s' wholelen=%d, offset=%d, sslen=%d\n\n", str, totalLength, ssStart, ssEnd)
		return getGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg1+errMsg2)
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
	str := object.GoStringFromStringObject(params[0].(*object.Object))

	// Get substring start and end offset
	ssStart := params[1].(int64)
	ssEnd := params[2].(int64)

	// Validate boundaries.
	totalLength := int64(len(str))
	if totalLength < 1 || ssStart < 0 || ssEnd < 1 || ssStart > (totalLength-1) || ssEnd > totalLength {
		errMsg1 := "substringStartEnd: Either nil input byte array, invalid substring offset, or invalid substring length"
		errMsg2 := fmt.Sprintf("\n\twhole='%s' wholelen=%d, offset=%d, sslen=%d\n\n", str, totalLength, ssStart, ssEnd)
		return getGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg1+errMsg2)
	}

	// Compute substring.
	str = str[ssStart:ssEnd]

	// Return new string in an object.
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.toCharArray()[C"
func toCharArray(params []interface{}) interface{} {
	// params[0]: input string
	obj := params[0].(*object.Object)
	bytes := obj.FieldTable["value"].Fvalue.([]types.JavaByte)
	var iArray []int64
	for _, bb := range bytes {
		iArray = append(iArray, int64(bb))
	}
	return Populator("[C", types.CharArray, iArray)
}

// "java/lang/String.toLowerCase()Ljava/lang/String;"
func toLowerCase(params []interface{}) interface{} {
	// params[0]: input string
	str := strings.ToLower(object.GoStringFromStringObject(params[0].(*object.Object)))
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.toUpperCase()Ljava/lang/String;"
func toUpperCase(params []interface{}) interface{} {
	// params[0]: input string
	str := strings.ToUpper(object.GoStringFromStringObject(params[0].(*object.Object)))
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/String.trim()Ljava/lang/String;"
func trimString(params []interface{}) interface{} {
	// params[0]: input string
	str := strings.Trim(object.GoStringFromStringObject(params[0].(*object.Object)), " ")
	obj := object.StringObjectFromGoString(str)
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
		return getGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
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
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
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
		return getGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
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
		return getGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
	}

	return offset
}

// "java/lang/String.hashCode()I"
func stringHashCode(params []interface{}) interface{} {
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
	baseObj := params[0].(*object.Object)
	baseStr := object.GoStringFromStringObject(baseObj)
	argObj := params[1].(*object.Object)
	prefix := object.GoStringFromStringObject(argObj)
	if len(params) == 3 {
		offset := int(params[2].(int64))
		if offset < 0 || offset > len(baseStr) {
			errMsg := fmt.Sprintf("stringStartsWith: base: %s, prefix: %s, offset: %d", baseStr, prefix, offset)
			return getGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
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
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	var srcBytes []types.JavaByte
	switch srcFld.Fvalue.(type) {
	case []byte:
		srcBytes = object.JavaByteArrayFromGoByteArray(srcFld.Fvalue.([]byte))
	case []types.JavaByte:
		srcBytes = srcFld.Fvalue.([]types.JavaByte)
	default:
		errMsg := fmt.Sprintf("stringGetChars: Invalid value field type (%s : %T) in base object", srcFld.Ftype, srcFld.Fvalue)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
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
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	dstChars, ok := dstFld.Fvalue.([]int64)
	if !ok {
		errMsg := fmt.Sprintf("stringGetChars: Invalid value field type (%s : %T) in char array object",
			dstFld.Ftype, dstFld.Fvalue)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get char array start offset.
	dstBegin := params[4].(int64)

	// Compute chara array length.
	dstLength := int64(len(dstChars))

	// Validate boundaries.
	if srcBegin < 0 || srcEnd < srcBegin || srcEnd > srcLength || dstBegin < 0 || dstBegin+(srcEnd-srcBegin) > dstLength {
		errMsg1 := "stringGetChars: Either nil input byte array, invalid substring offset, or invalid substring length"
		errMsg2 := fmt.Sprintf("\n\twholelen=%d, offset=%d, sslen=%d\n\n", srcLength, srcBegin, srcEnd)
		return getGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg1+errMsg2)
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
	// Get field of base object.
	srcFld, ok := params[0].(*object.Object).FieldTable["value"]
	if !ok {
		errMsg := fmt.Sprintf("stringIndexOfCh: Missing value field in base object")
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
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
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
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
			return getGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
		}
		endIndex = params[3].(int64)
		if endIndex > lenSrcBytes || beginIndex > endIndex {
			errMsg := fmt.Sprintf("stringIndexOfCh: Base string len: %d, end index: %d", lenSrcBytes, endIndex)
			return getGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
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
			return getGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
		}
		endIndex = params[3].(int64)
		if endIndex > lenOrigBaseString || beginIndex > endIndex {
			errMsg := fmt.Sprintf("stringIndexOfString: Base string len: %d, end index: %d", lenOrigBaseString, endIndex)
			return getGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
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
	baseString := object.GoStringFromStringObject(params[0].(*object.Object))
	if len(strings.TrimSpace(baseString)) == 0 {
		return types.JavaBoolTrue
	} else {
		return types.JavaBoolFalse
	}
}

func stringIsEmpty(params []interface{}) interface{} {
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

func stringReplaceAllRegex(params []interface{}) interface{} {
	// Get 3 string arguments.
	input := object.GoStringFromStringObject(params[0].(*object.Object))
	pattern := object.GoStringFromStringObject(params[1].(*object.Object))
	replacement := object.GoStringFromStringObject(params[2].(*object.Object))

	// Compile the regular expression.
	re, err := regexp.Compile(pattern)
	if err != nil {
		errMsg := fmt.Sprintf("stringReplaceAllRegex: Invalid regular expression pattern: %s", pattern)
		return getGErrBlk(excNames.PatternSyntaxException, errMsg)
	}

	// Replace all substrings that match the pattern with the replacement string.
	result := re.ReplaceAllString(input, replacement)

	return object.StringObjectFromGoString(result)

}

func stringReplaceFirstRegex(params []interface{}) interface{} {
	// Get 3 string arguments.
	input := object.GoStringFromStringObject(params[0].(*object.Object))
	pattern := object.GoStringFromStringObject(params[1].(*object.Object))
	replacement := object.GoStringFromStringObject(params[2].(*object.Object))

	// Compile the regular expression.
	re, err := regexp.Compile(pattern)
	if err != nil {
		errMsg := fmt.Sprintf("stringReplaceFirstRegex: Invalid regular expression pattern: %s", pattern)
		return getGErrBlk(excNames.PatternSyntaxException, errMsg)
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
	// params[0] = base string
	// params[1] = regular expression in a string

	input := object.GoStringFromStringObject(params[0].(*object.Object))
	pattern := object.GoStringFromStringObject(params[1].(*object.Object))

	// Compile the regular expression.
	re, err := regexp.Compile(pattern)
	if err != nil {
		errMsg := fmt.Sprintf("stringSplit: Invalid regular expression pattern: %s", pattern)
		return getGErrBlk(excNames.PatternSyntaxException, errMsg)
	}

	// Split input based on the pattern.
	result := re.Split(input, -1) // -1 means split on all occurrences.

	// Prepare object array and return it.
	var outObjArray []*object.Object
	for ix := 0; ix < len(result); ix++ {
		outObjArray = append(outObjArray, object.StringObjectFromGoString(result[ix]))
	}
	return Populator("[Ljava/lang/String;", types.RefArray, outObjArray)

}

func stringSplitLimit(params []interface{}) interface{} {
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
		return getGErrBlk(excNames.PatternSyntaxException, errMsg)
	}

	// Split input based on the pattern.
	result := re.Split(input, int(limit))

	// Prepare object array and return it.
	var outObjArray []*object.Object
	for ix := 0; ix < len(result); ix++ {
		outObjArray = append(outObjArray, object.StringObjectFromGoString(result[ix]))
	}
	return Populator("[Ljava/lang/String;", types.RefArray, outObjArray)
}

func stringStrip(params []interface{}) interface{} {
	input := object.GoStringFromStringObject(params[0].(*object.Object))
	result := strings.TrimSpace(input)
	return object.StringObjectFromGoString(result)
}

func stringStripLeading(params []interface{}) interface{} {
	input := object.GoStringFromStringObject(params[0].(*object.Object))
	result := strings.TrimLeftFunc(input, unicode.IsSpace)
	return object.StringObjectFromGoString(result)
}

func stringStripTrailing(params []interface{}) interface{} {
	input := object.GoStringFromStringObject(params[0].(*object.Object))
	result := strings.TrimRightFunc(input, unicode.IsSpace)
	return object.StringObjectFromGoString(result)
}
