// This file is generated by generate-std.clj script. Do not edit manually!

package base64

import (
	"fmt"
	. "github.com/lab47/lace/core"
	"os"
)

func InternsOrThunks() {
	if VerbosityLevel > 0 {
		fmt.Fprintln(os.Stderr, "Lazily running slow version of base64.InternsOrThunks().")
	}
	base64Namespace.ResetMeta(MakeMeta(nil, `Implements base64 encoding as specified by RFC 4648.`, "1.0"))

	base64Namespace.InternVar("decode-string", decode_string_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Returns the bytes represented by the base64 string s.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	base64Namespace.InternVar("encode-string", encode_string_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Returns the base64 encoding of s.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

}
