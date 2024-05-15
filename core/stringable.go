package core

import (
	"fmt"
)

func AssertStringable(obj Object, msg string) (String, error) {
	switch c := obj.(type) {
	case String:
		return c, nil
	case Char:
		return MakeString(string(c.Ch)), nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Stringable", obj.GetType().Name())
		}
		return "", StubNewError(msg)
	}
}

func EnsureStringable(args []Object, index int) (String, error) {
	switch c := args[index].(type) {
	case String:
		return c, nil
	case Char:
		return MakeString(string(c.Ch)), nil
	default:
		return "", StubNewArgTypeError(index, c, "Stringable")
	}
}
