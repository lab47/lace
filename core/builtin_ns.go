package core

import "fmt"

var builtinNSSetup = map[string]func(env *Env) error{}

func runBuiltinNS(env *Env, name string) error {
	if fn, ok := builtinNSSetup[name]; ok {
		return fn(env)
	}

	return fmt.Errorf("unknown builtin ns: %s", name)
}
