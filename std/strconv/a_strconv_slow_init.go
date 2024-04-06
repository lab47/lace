// This file is generated by generate-std.clj script. Do not edit manually!


package strconv

import (
	. "github.com/lab47/lace/core"
	"fmt"
	"os"
)

func InternsOrThunks(env *Env, ns *Namespace) {
	if VerbosityLevel > 0 {
		fmt.Fprintln(os.Stderr, "Lazily running slow version of strconv.InternsOrThunks().")
	}
	ns.ResetMeta(MakeMeta(nil, `Implements conversions to and from string representations of basic data types.`, "1.0"))

	
	ns.InternVar("atoi", atoi_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Equivalent to (parse-int s 10 0).`, "1.0").Plus(MakeKeyword("tag"), String{S: "Int"}))

	ns.InternVar("can-backquote?", iscan_backquote_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Reports whether the string s can be represented unchanged as a single-line backquoted string without control characters other than tab.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Boolean"}))

	ns.InternVar("format-bool", format_bool_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("b"))),
			`Returns "true" or "false" according to the value of b.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar("format-double", format_double_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("f"), MakeSymbol("fmt"), MakeSymbol("prec"), MakeSymbol("bitSize"))),
			`Converts the floating-point number f to a string, according to the format fmt and precision prec. It rounds the result assuming that the original was obtained from a floating-point value of bitSize bits (32 for float32, 64 for float64).
  The format fmt is one of 'b' (-ddddp±ddd, a binary exponent), 'e' (-d.dddde±dd, a decimal exponent), 'E' (-d.ddddE±dd, a decimal exponent), 'f' (-ddd.dddd, no exponent), 'g' ('e' for large exponents, 'f' otherwise), or 'G' ('E' for large exponents, 'f' otherwise).
  The precision prec controls the number of digits (excluding the exponent) printed by the 'e', 'E', 'f', 'g', and 'G' formats. For 'e', 'E', and 'f' it is the number of digits after the decimal point. For 'g' and 'G' it is the maximum number of significant digits (trailing zeros are removed). The special precision -1 uses the smallest number of digits necessary such that ParseFloat will return f exactly.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar("format-int", format_int_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("i"), MakeSymbol("base"))),
			`Returns the string representation of i in the given base, for 2 <= base <= 36. The result uses the lower-case letters 'a' to 'z' for digit values >= 10.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar("graphic?", isgraphic_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("c"))),
			`Reports whether the char is defined as a Graphic by Unicode. Such characters include letters, marks, numbers, punctuation, symbols, and spaces, from categories L, M, N, P, S, and Zs.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Boolean"}))

	ns.InternVar("itoa", itoa_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("i"))),
			`Equivalent to (format-int i 10).`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar("parse-bool", parse_bool_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Returns the boolean value represented by the string. It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False. Any other value returns an error.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Boolean"}))

	ns.InternVar("parse-double", parse_double_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Converts the string s to a floating-point number.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar("parse-int", parse_int_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"), MakeSymbol("base"), MakeSymbol("bitSize"))),
			`Interprets a string s in the given base (0, 2 to 36) and bit size (0 to 64) and returns the corresponding value i.
  If base == 0, the base is implied by the string's prefix: base 16 for "0x", base 8 for "0", and base 10 otherwise. For bases 1, below 0 or above 36 an error is returned.
  The bitSize argument specifies the integer type that the result must fit into. Bit sizes 0, 8, 16, 32, and 64 correspond to int, int8, int16, int32, and int64. For a bitSize below 0 or above 64 an error is returned.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Int"}))

	ns.InternVar("printable?", isprintable_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("c"))),
			`Reports whether the char is defined as printable by Joker: letters, numbers, punctuation, symbols and ASCII space.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Boolean"}))

	ns.InternVar("quote", quote_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Returns a double-quoted string literal representing s. The returned string uses escape sequences (\t, \n, \xFF, \u0100)
  for control characters and non-printable characters as defined by printable?.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar("quote-char", quote_char_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("c"))),
			`Returns a single-quoted char literal representing the character. The returned string uses escape sequences (\t, \n, \xFF, \u0100)
  for control characters and non-printable characters as defined by printable?.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar("quote-char-to-ascii", quote_char_to_ascii_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("c"))),
			`Returns a single-quoted char literal representing the character. The returned string uses escape sequences (\t, \n, \xFF, \u0100)
  for non-ASCII characters and non-printable characters as defined by printable?.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar("quote-char-to-graphic", quote_char_to_graphic_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("c"))),
			`Returns a single-quoted char literal representing the character. The returned string uses escape sequences (\t, \n, \xFF, \u0100)
  for non-ASCII characters and non-printable characters as defined by graphic?.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar("quote-to-ascii", quote_to_ascii_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Returns a double-quoted string literal representing s. The returned string uses escape sequences (\t, \n, \xFF, \u0100)
  for non-ASCII characters and non-printable characters as defined by printable?.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar("quote-to-graphic", quote_to_graphic_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Returns a double-quoted string literal representing s. The returned string uses escape sequences (\t, \n, \xFF, \u0100)
  for non-ASCII characters and non-printable characters as defined by graphic?.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar("unquote", unquote_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Interprets s as a single-quoted, double-quoted, or backquoted string literal, returning the string value that s quotes.
  (If s is single-quoted, it would be a Go character literal; Unquote returns the corresponding one-character string.)`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

}
