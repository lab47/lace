// Generated by gen_data. Don't modify manually!

//go:build !gen_data
// +build !gen_data

//
package core

import "encoding/base64"
import _ "embed"

//go:embed a_core_data.data
var coreData []byte

func coreSetup(env *Env) error {
	ns := env.EnsureNamespace(MakeSymbol("lace.core"))
	raw, err := base64.StdEncoding.AppendDecode(nil, coreData)
	if err != nil {
		return err
	}
	return processInEnvInNS(env, ns, raw)
}

func init() {
	builtinNSSetup["lace.core"] = coreSetup
}
