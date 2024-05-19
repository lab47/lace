// Generated by gen_data. Don't modify manually!

//go:build !gen_data
// +build !gen_data
package core

import _ "embed"

//go:embed a_test_data.data
var testData []byte

func testSetup(env *Env) error {
	ns := env.EnsureNamespace(MakeSymbol("lace.test"))
	return processInEnvInNS(env, ns, testData)
}

func init() {
	builtinNSSetup["lace.test"] = testSetup
}



