package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"

	"github.com/lab47/lace/core"
	"golang.org/x/crypto/blake2b"
)

func Setup(env *core.Env) error {
	b := core.NewNSBuilder(env, "lace.crypto")

	b.Defn(&core.DefnInfo{
		Name: "hmac",
		Doc:  "Returns HMAC signature for message and key using specified algorithm. Algorithm is one of the following: :sha1, :sha224, :sha256, :sha384, :sha512.",
		Args: []string{"algorithm", "message", "key"},
		Tag:  "String",
		Fn:   hmacSum,
	})

	b.Defn(&core.DefnInfo{
		Name: "sha256",
		Doc:  "Returns the SHA256 checksum of the data.",
		Args: []string{"data"},
		Tag:  "String",
		Fn: func(data []byte) []byte {
			ary := sha256.Sum256(data)
			return ary[:]
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "sha224",
		Doc:  "Returns the SHA224 checksum of the data.",
		Args: []string{"data"},
		Tag:  "String",
		Fn: func(data []byte) []byte {
			ary := sha256.Sum224(data)
			return ary[:]
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "sha386",
		Doc:  "Returns the SHA386 checksum of the data.",
		Args: []string{"data"},
		Tag:  "String",
		Fn: func(data []byte) []byte {
			ary := sha512.Sum384(data)
			return ary[:]
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "sha512",
		Doc:  "Returns the SHA512 checksum of the data.",
		Args: []string{"data"},
		Tag:  "String",
		Fn: func(data []byte) []byte {
			ary := sha512.Sum512(data)
			return ary[:]
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "sha512-224",
		Doc:  "Returns the SHA512/224 checksum of the data.",
		Args: []string{"data"},
		Tag:  "String",
		Fn: func(data []byte) []byte {
			ary := sha512.Sum512_224(data)
			return ary[:]
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "sha512-224",
		Doc:  "Returns the SHA512/256 checksum of the data.",
		Args: []string{"data"},
		Tag:  "String",
		Fn: func(data []byte) []byte {
			ary := sha512.Sum512_256(data)
			return ary[:]
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "md5",
		Doc:  "Returns the MD5 checksum of the data.",
		Args: []string{"data"},
		Tag:  "String",
		Fn: func(data []byte) []byte {
			ary := md5.Sum(data)
			return ary[:]
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "sha1",
		Doc:  "Returns the SHA1 checksum of the data.",
		Args: []string{"data"},
		Tag:  "String",
		Fn: func(data []byte) []byte {
			ary := sha1.Sum(data)
			return ary[:]
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "blake2b",
		Doc:  "Returns the Blake2b 256 checksum of the data.",
		Args: []string{"data"},
		Tag:  "String",
		Fn: func(data []byte) []byte {
			ary := blake2b.Sum256(data)
			return ary[:]
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "blake2b-512",
		Doc:  "Returns the Blake2b 512 checksum of the data.",
		Args: []string{"data"},
		Tag:  "String",
		Fn: func(data []byte) []byte {
			ary := blake2b.Sum512(data)
			return ary[:]
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "blake2b-512",
		Doc:  "Returns the Blake2b 384 checksum of the data.",
		Args: []string{"data"},
		Tag:  "String",
		Fn: func(data []byte) []byte {
			ary := blake2b.Sum384(data)
			return ary[:]
		},
	})

	return nil
}

func init() {
	core.AddNativeNamespace("lace.crypto", Setup)
}

func hmacSum(algorithm core.Keyword, message, key string) (string, error) {
	var h func() hash.Hash
	switch algorithm.Name() {
	case "sha1":
		h = sha1.New
	case "sha224":
		h = sha256.New224
	case "sha256":
		h = sha256.New
	case "sha384":
		h = sha512.New384
	case "sha512":
		h = sha512.New
	default:
		return "", core.StubNewError("unsupported algorithm " + algorithm.Name() +
			". Supported algorithms are: :sha1, :sha224, :sha256, :sha384, :sha512")
	}
	mac := hmac.New(h, []byte(key))
	mac.Write([]byte(message))
	return string(mac.Sum(nil)), nil
}
