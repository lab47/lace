// This file is generated by generate-std.joke script. Do not edit manually!

package hex

import (
	"fmt"
	"os"

	. "github.com/lab47/lace/core"
)

func InternsOrThunks() {
	if VerbosityLevel > 0 {
		fmt.Fprintln(os.Stderr, "Lazily running slow version of hex.InternsOrThunks().")
	}
	hexNamespace.ResetMeta(MakeMeta(nil, `Implements hexadecimal encoding and decoding.`, "1.0"))

	hexNamespace.InternVar("decode-string", decode_string_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Returns the bytes represented by the hexadecimal string s.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	hexNamespace.InternVar("encode-string", encode_string_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Returns the hexadecimal encoding of s.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

}
