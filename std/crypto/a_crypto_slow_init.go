// This file is generated by generate-std.clj script. Do not edit manually!

package crypto

import (
	"fmt"
	. "github.com/lab47/lace/core"
	"os"
)

func InternsOrThunks() {
	if VerbosityLevel > 0 {
		fmt.Fprintln(os.Stderr, "Lazily running slow version of crypto.InternsOrThunks().")
	}
	cryptoNamespace.ResetMeta(MakeMeta(nil, `Implements common cryptographic and hash functions.`, "1.0"))

	cryptoNamespace.InternVar("hmac", hmac_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("algorithm"), MakeSymbol("message"), MakeSymbol("key"))),
			`Returns HMAC signature for message and key using specified algorithm.
  Algorithm is one of the following: :sha1, :sha224, :sha256, :sha384, :sha512.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	cryptoNamespace.InternVar("md5", md5_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the MD5 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	cryptoNamespace.InternVar("sha1", sha1_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the SHA1 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	cryptoNamespace.InternVar("sha224", sha224_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the SHA224 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	cryptoNamespace.InternVar("sha256", sha256_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the SHA256 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	cryptoNamespace.InternVar("sha384", sha384_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the SHA384 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	cryptoNamespace.InternVar("sha512", sha512_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the SHA512 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	cryptoNamespace.InternVar("sha512-224", sha512_224_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the SHA512/224 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	cryptoNamespace.InternVar("sha512-256", sha512_256_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the SHA512/256 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

}
