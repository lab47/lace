// This file is generated by generate-std.clj script. Do not edit manually!

package io

import (
	"fmt"
	. "github.com/lab47/lace/core"
	"os"
)

func InternsOrThunks() {
	if VerbosityLevel > 0 {
		fmt.Fprintln(os.Stderr, "Lazily running slow version of io.InternsOrThunks().")
	}
	ioNamespace.ResetMeta(MakeMeta(nil, `Provides basic interfaces to I/O primitives.`, "1.0"))

	ioNamespace.InternVar("close", close_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("f"))),
			`Closes f (IOWriter, IOReader, or File) if possible. Otherwise throws an error.`, "1.0"))

	ioNamespace.InternVar("copy", copy_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("dst"), MakeSymbol("src"))),
			`Copies from src to dst until either EOF is reached on src or an error occurs.
  Returns the number of bytes copied or throws an error.
  src must be IOReader, e.g. as returned by lace.os/open.
  dst must be IOWriter, e.g. as returned by lace.os/create.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Int"}))

	ioNamespace.InternVar("pipe", pipe_,
		MakeMeta(
			NewListFrom(NewVectorFrom()),
			`Pipe creates a synchronous in-memory pipe. It can be used to connect code expecting an IOReader
  with code expecting an IOWriter.
  Returns a vector [reader, writer].`, "1.0"))

}
