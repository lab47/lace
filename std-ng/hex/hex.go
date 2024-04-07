package hex

import (
	"encoding/hex"

	"github.com/lab47/lace/core"
)

func Setup(env *core.Env) error {
	b := core.NewNSBuilder(env, "lace.hex")

	b.Defn(&core.DefnInfo{
		Name: "decode-string",
		Doc:  "Returns the bytes represented by the hexadecimal string s.",
		Args: []string{"s"},
		Tag:  "String",
		Fn:   hex.DecodeString,
	})

	b.Defn(&core.DefnInfo{
		Name: "encode-string",
		Doc:  "Returns the hexadecimal encoding of s.",
		Args: []string{"s"},
		Tag:  "String",
		Fn:   hex.EncodeToString,
	})

	return nil
}

func init() {
	core.AddNativeNamespace("lace.hex", Setup)
}
