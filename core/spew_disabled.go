//go:build !go_spew
// +build !go_spew

package core

var procGoSpew = func(env *Env, args []any) (any, error) {
	return MakeBoolean(false), nil
}
