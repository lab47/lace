package core

import (
	"fmt"
)

func AssertStringable(obj Object, msg string) (String, error) {
	switch c := obj.(type) {
	case String:
		return c, nil
	case Char:
		return String{S: string(c.Ch)}, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Stringable", obj.GetType().ToString(false))
		}
		return String{}, StubNewError(msg)
	}
}

func EnsureStringable(args []Object, index int) (String, error) {
	switch c := args[index].(type) {
	case String:
		return c, nil
	case Char:
		return String{S: string(c.Ch)}, nil
	default:
		return String{}, StubNewArgTypeError(index, c, "Stringable")
	}
}
