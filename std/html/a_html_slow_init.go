// This file is generated by generate-std.clj script. Do not edit manually!

//go:build !fast_init
// +build !fast_init

package html

import (
	"fmt"
	"os"

	. "github.com/lab47/lace/core"
)

func InternsOrThunks() {
	if VerbosityLevel > 0 {
		fmt.Fprintln(os.Stderr, "Lazily running slow version of html.InternsOrThunks().")
	}
	htmlNamespace.ResetMeta(MakeMeta(nil, `Provides functions for escaping and unescaping HTML text.`, "1.0"))

	htmlNamespace.InternVar("escape", escape_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Escapes special characters like < to become &lt;. It escapes only five such characters: <, >, &, ' and ".`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	htmlNamespace.InternVar("unescape", unescape_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Unescapes entities like &lt; to become <.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

}
