// Generated by gen_data. Don't modify manually!

//go:build !gen_data
// +build !gen_data

//
package core

import "encoding/base64"
import _ "embed"

//go:embed a_hiccup_data.data
var hiccupData []byte

func hiccupSetup(env *Env) error {
	ns := env.EnsureNamespace(MakeSymbol("lace.hiccup"))
	raw, err := base64.StdEncoding.AppendDecode(nil, hiccupData)
	if err != nil {
		return err
	}
	return processInEnvInNS(env, ns, raw)
}

func init() {
	builtinNSSetup["lace.hiccup"] = hiccupSetup
}
