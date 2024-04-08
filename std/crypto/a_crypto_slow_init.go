// This file is generated by generate-std.clj script. Do not edit manually!


package crypto

import (
	. "github.com/lab47/lace/core"
	"fmt"
	"os"
)

func InternsOrThunks(env *Env, ns *Namespace) {
	if VerbosityLevel > 0 {
		fmt.Fprintln(os.Stderr, "Lazily running slow version of crypto.InternsOrThunks().")
	}
	ns.ResetMeta(MakeMeta(nil, `Implements common cryptographic and hash functions.`, "1.0"))

	
	ns.InternVar(env, "hmac", hmac_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("algorithm"), MakeSymbol("message"), MakeSymbol("key"))),
			`Returns HMAC signature for message and key using specified algorithm.
  Algorithm is one of the following: :sha1, :sha224, :sha256, :sha384, :sha512.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar(env, "md5", md5_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the MD5 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar(env, "sha1", sha1_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the SHA1 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar(env, "sha224", sha224_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the SHA224 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar(env, "sha256", sha256_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the SHA256 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar(env, "sha384", sha384_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the SHA384 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar(env, "sha512", sha512_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the SHA512 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar(env, "sha512-224", sha512_224_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the SHA512/224 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar(env, "sha512-256", sha512_256_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data"))),
			`Returns the SHA512/256 checksum of the data.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

}
