package core

import (
	"fmt"
)

func AssertStringable(obj Object, msg string) String {
	switch c := obj.(type) {
	case String:
		return c
	case Char:
		return String{S: string(c.Ch)}
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Stringable", obj.GetType().ToString(false))
		}
		panic(StubNewError(msg))
	}
}

func EnsureStringable(args []Object, index int) String {
	switch c := args[index].(type) {
	case String:
		return c
	case Char:
		return String{S: string(c.Ch)}
	default:
		panic(StubNewArgTypeError(index, c, "Stringable"))
	}
}
