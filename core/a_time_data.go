// Generated by gen_data. Don't modify manually!

//go:build !gen_data
// +build !gen_data
package core

import _ "embed"

//go:embed a_time_data.data
var timeData []byte

func timeSetup(env *Env) error {
	ns := env.EnsureNamespace(MakeSymbol("lace.time"))
	return processInEnvInNS(env, ns, timeData)
}

func init() {
	builtinNSSetup["lace.time"] = timeSetup
}


