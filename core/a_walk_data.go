// Generated by gen_data. Don't modify manually!

//go:build !gen_data
// +build !gen_data
package core

import _ "embed"

//go:embed a_walk_data.data
var walkData []byte

func walkSetup(env *Env) error {
	ns := env.EnsureNamespace(MakeSymbol("lace.walk"))
	return processInEnvInNS(env, ns, walkData)
}

func init() {
	builtinNSSetup["lace.walk"] = walkSetup
}


