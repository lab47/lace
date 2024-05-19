// Generated by gen_data. Don't modify manually!

//go:build !gen_data
// +build !gen_data

//
package core

import "encoding/base64"
import _ "embed"

//go:embed a_walk_data.data
var walkData []byte

func walkSetup(env *Env) error {
	ns := env.EnsureNamespace(MakeSymbol("lace.walk"))
	raw, err := base64.StdEncoding.AppendDecode(nil, walkData)
	if err != nil {
		return err
	}
	return processInEnvInNS(env, ns, raw)
}

func init() {
	builtinNSSetup["lace.walk"] = walkSetup
}
