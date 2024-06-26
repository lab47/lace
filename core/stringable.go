package core

import (
	"fmt"
)

func AssertStringable(obj any, msg string) (String, error) {
	switch c := obj.(type) {
	case String:
		return c, nil
	case Char:
		return MakeString(string(c.Ch())), nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Stringable", TypeName(obj))
		}
		return nil, StubNewError(msg)
	}
}

func EnsureStringable(args []any, index int) (String, error) {
	switch c := args[index].(type) {
	case String:
		return c, nil
	case Char:
		return MakeString(string(c.Ch())), nil
	default:
		return nil, StubNewArgTypeError(index, c, "Stringable")
	}
}
