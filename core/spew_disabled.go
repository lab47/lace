//go:build !go_spew
// +build !go_spew

package core

var procGoSpew = func(env *Env, args []Object) (res Object) {
	return MakeBoolean(false)
}
