package base64

import (
	"encoding/base64"

	"github.com/lab47/lace/core"
)

func Setup(env *core.Env) error {
	b := core.NewNSBuilder(env, "lace.base64")

	b.Defn(&core.DefnInfo{
		Name: "decode-string",
		Doc:  "Returns the bytes represented by the base64 string s.",
		Tag:  "String",
		Args: []string{"s"},
		Fn:   base64.StdEncoding.DecodeString,
	})

	b.Defn(&core.DefnInfo{
		Name: "encode-string",
		Doc:  "Returns the base64 encoding of s.",
		Tag:  "String",
		Args: []string{"s"},
		Fn:   base64.StdEncoding.EncodeToString,
	})

	return nil
}

func init() {
	core.AddNativeNamespace("lace.base64", Setup)
}
