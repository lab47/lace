// This file is generated by generate-std.clj script. Do not edit manually!


package url

import (
	. "github.com/lab47/lace/core"
	"fmt"
	"os"
)

func InternsOrThunks(env *Env, ns *Namespace) {
	if VerbosityLevel > 0 {
		fmt.Fprintln(os.Stderr, "Lazily running slow version of url.InternsOrThunks().")
	}
	ns.ResetMeta(MakeMeta(nil, `Parses URLs and implements query escaping.`, "1.0"))

	
	ns.InternVar(env, "path-escape", path_escape_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Escapes the string so it can be safely placed inside a URL path segment.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar(env, "path-unescape", path_unescape_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Does the inverse transformation of path-escape, converting each 3-byte encoded
  substring of the form "%AB" into the hex-decoded byte 0xAB. It also converts
  '+' into ' ' (space). It returns an error if any % is not followed by two hexadecimal digits.

  PathUnescape is identical to QueryUnescape except that it does not unescape '+' to ' ' (space).`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar(env, "query-escape", query_escape_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Escapes the string so it can be safely placed inside a URL query.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar(env, "query-unescape", query_unescape_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Does the inverse transformation of query-escape, converting each 3-byte encoded
  substring of the form "%AB" into the hex-decoded byte 0xAB. It also converts
  '+' into ' ' (space). It returns an error if any % is not followed by two hexadecimal digits.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

}
