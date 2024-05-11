package core

import (
	"fmt"
	"io"
	"reflect"
)

func AssertComparable(env *Env, obj Object, msg string) (Comparable, error) {
	switch c := obj.(type) {
	case Comparable:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Comparable", obj.GetType().Name())
		}
		var v Comparable
		return v, env.NewError(msg)
	}
}

func EnsureComparable(env *Env, args []Object, index int) (Comparable, error) {
	if len(args) <= index {
		var t Comparable
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Comparable:
		return c, nil
	default:
		var v Comparable
		return v, env.NewArgTypeError(index, c, "Comparable")
	}
}

func AssertVector(env *Env, obj Object, msg string) (*Vector, error) {
	switch c := obj.(type) {
	case *Vector:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Vector", obj.GetType().Name())
		}
		var v *Vector
		return v, env.NewError(msg)
	}
}

func EnsureVector(env *Env, args []Object, index int) (*Vector, error) {
	if len(args) <= index {
		var t *Vector
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case *Vector:
		return c, nil
	default:
		var v *Vector
		return v, env.NewArgTypeError(index, c, "Vector")
	}
}

func AssertChar(env *Env, obj Object, msg string) (Char, error) {
	switch c := obj.(type) {
	case Char:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Char", obj.GetType().Name())
		}
		var v Char
		return v, env.NewError(msg)
	}
}

func EnsureChar(env *Env, args []Object, index int) (Char, error) {
	if len(args) <= index {
		var t Char
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Char:
		return c, nil
	default:
		var v Char
		return v, env.NewArgTypeError(index, c, "Char")
	}
}

func AssertString(env *Env, obj Object, msg string) (String, error) {
	switch c := obj.(type) {
	case String:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "String", obj.GetType().Name())
		}
		var v String
		return v, env.NewError(msg)
	}
}

func EnsureString(env *Env, args []Object, index int) (String, error) {
	if len(args) <= index {
		var t String
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case String:
		return c, nil
	default:
		var v String
		return v, env.NewArgTypeError(index, c, "String")
	}
}

func AssertSymbol(env *Env, obj Object, msg string) (Symbol, error) {
	switch c := obj.(type) {
	case Symbol:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Symbol", obj.GetType().Name())
		}
		var v Symbol
		return v, env.NewError(msg)
	}
}

func EnsureSymbol(env *Env, args []Object, index int) (Symbol, error) {
	if len(args) <= index {
		var t Symbol
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Symbol:
		return c, nil
	default:
		var v Symbol
		return v, env.NewArgTypeError(index, c, "Symbol")
	}
}

func AssertKeyword(env *Env, obj Object, msg string) (Keyword, error) {
	switch c := obj.(type) {
	case Keyword:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Keyword", obj.GetType().Name())
		}
		var v Keyword
		return v, env.NewError(msg)
	}
}

func EnsureKeyword(env *Env, args []Object, index int) (Keyword, error) {
	if len(args) <= index {
		var t Keyword
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Keyword:
		return c, nil
	default:
		var v Keyword
		return v, env.NewArgTypeError(index, c, "Keyword")
	}
}

func AssertRegex(env *Env, obj Object, msg string) (*Regex, error) {
	switch c := obj.(type) {
	case *Regex:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Regex", obj.GetType().Name())
		}
		var v *Regex
		return v, env.NewError(msg)
	}
}

func EnsureRegex(env *Env, args []Object, index int) (*Regex, error) {
	if len(args) <= index {
		var t *Regex
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case *Regex:
		return c, nil
	default:
		var v *Regex
		return v, env.NewArgTypeError(index, c, "Regex")
	}
}

func AssertBoolean(env *Env, obj Object, msg string) (Boolean, error) {
	switch c := obj.(type) {
	case Boolean:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Boolean", obj.GetType().Name())
		}
		var v Boolean
		return v, env.NewError(msg)
	}
}

func EnsureBoolean(env *Env, args []Object, index int) (Boolean, error) {
	if len(args) <= index {
		var t Boolean
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Boolean:
		return c, nil
	default:
		var v Boolean
		return v, env.NewArgTypeError(index, c, "Boolean")
	}
}

func AssertTime(env *Env, obj Object, msg string) (Time, error) {
	switch c := obj.(type) {
	case Time:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Time", obj.GetType().Name())
		}
		var v Time
		return v, env.NewError(msg)
	}
}

func EnsureTime(env *Env, args []Object, index int) (Time, error) {
	if len(args) <= index {
		var t Time
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Time:
		return c, nil
	default:
		var v Time
		return v, env.NewArgTypeError(index, c, "Time")
	}
}

func AssertNumber(env *Env, obj Object, msg string) (Number, error) {
	switch c := obj.(type) {
	case Number:
		return c, nil
	case *ReflectValue:
		switch c.val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return MakeInt(int(c.val.Int())), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return MakeInt(int(c.val.Uint())), nil
		default:
			if msg == "" {
				msg = fmt.Sprintf("Expected %s, got %s", "Number", obj.GetType().Name())
			}
			var v Number
			return v, env.NewError(msg)
		}
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Number", obj.GetType().Name())
		}
		var v Number
		return v, env.NewError(msg)
	}
}

func EnsureNumber(env *Env, args []Object, index int) (Number, error) {
	if len(args) <= index {
		var t Number
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Number:
		return c, nil
	default:
		var v Number
		return v, env.NewArgTypeError(index, c, "Number")
	}
}

func AssertSeqable(env *Env, obj Object, msg string) (Seqable, error) {
	switch c := obj.(type) {
	case Seqable:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Seqable", obj.GetType().Name())
		}
		var v Seqable
		return v, env.NewError(msg)
	}
}

func EnsureSeqable(env *Env, args []Object, index int) (Seqable, error) {
	if len(args) <= index {
		var t Seqable
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Seqable:
		return c, nil
	default:
		var v Seqable
		return v, env.NewArgTypeError(index, c, "Seqable")
	}
}

func AssertCallable(env *Env, obj Object, msg string) (Callable, error) {
	switch c := obj.(type) {
	case Callable:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Callable", obj.GetType().Name())
		}
		var v Callable
		return v, env.NewError(msg)
	}
}

func EnsureCallable(env *Env, args []Object, index int) (Callable, error) {
	if len(args) <= index {
		var t Callable
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Callable:
		return c, nil
	default:
		var v Callable
		return v, env.NewArgTypeError(index, c, "Callable")
	}
}

func AssertType(env *Env, obj Object, msg string) (*Type, error) {
	switch c := obj.(type) {
	case *Type:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Type", obj.GetType().Name())
		}
		var v *Type
		return v, env.NewError(msg)
	}
}

func EnsureType(env *Env, args []Object, index int) (*Type, error) {
	if len(args) <= index {
		var t *Type
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case *Type:
		return c, nil
	default:
		var v *Type
		return v, env.NewArgTypeError(index, c, "Type")
	}
}

func AssertMeta(env *Env, obj Object, msg string) (Meta, error) {
	switch c := obj.(type) {
	case Meta:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Meta", obj.GetType().Name())
		}
		var v Meta
		return v, env.NewError(msg)
	}
}

func EnsureMeta(env *Env, args []Object, index int) (Meta, error) {
	if len(args) <= index {
		var t Meta
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Meta:
		return c, nil
	default:
		var v Meta
		return v, env.NewArgTypeError(index, c, "Meta")
	}
}

func AssertInt(env *Env, obj Object, msg string) (Int, error) {
	switch c := obj.(type) {
	case Int:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Int", obj.GetType().Name())
		}
		var v Int
		return v, env.NewError(msg)
	}
}

func EnsureInt(env *Env, args []Object, index int) (Int, error) {
	if len(args) <= index {
		var t Int
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Int:
		return c, nil
	default:
		var v Int
		return v, env.NewArgTypeError(index, c, "Int")
	}
}

func AssertDouble(env *Env, obj Object, msg string) (Double, error) {
	switch c := obj.(type) {
	case Double:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Double", obj.GetType().Name())
		}
		var v Double
		return v, env.NewError(msg)
	}
}

func EnsureDouble(env *Env, args []Object, index int) (Double, error) {
	if len(args) <= index {
		var t Double
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Double:
		return c, nil
	default:
		var v Double
		return v, env.NewArgTypeError(index, c, "Double")
	}
}

func AssertStack(env *Env, obj Object, msg string) (Stack, error) {
	switch c := obj.(type) {
	case Stack:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Stack", obj.GetType().Name())
		}
		var v Stack
		return v, env.NewError(msg)
	}
}

func EnsureStack(env *Env, args []Object, index int) (Stack, error) {
	if len(args) <= index {
		var t Stack
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Stack:
		return c, nil
	default:
		var v Stack
		return v, env.NewArgTypeError(index, c, "Stack")
	}
}

func AssertMap(env *Env, obj Object, msg string) (Map, error) {
	switch c := obj.(type) {
	case Map:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Map", obj.GetType().Name())
		}
		var v Map
		return v, env.NewError(msg)
	}
}

func EnsureMap(env *Env, args []Object, index int) (Map, error) {
	if len(args) <= index {
		var t Map
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Map:
		return c, nil
	default:
		var v Map
		return v, env.NewArgTypeError(index, c, "Map")
	}
}

func AssertSet(env *Env, obj Object, msg string) (Set, error) {
	switch c := obj.(type) {
	case Set:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Set", obj.GetType().Name())
		}
		var v Set
		return v, env.NewError(msg)
	}
}

func EnsureSet(env *Env, args []Object, index int) (Set, error) {
	if len(args) <= index {
		var t Set
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Set:
		return c, nil
	default:
		var v Set
		return v, env.NewArgTypeError(index, c, "Set")
	}
}

func AssertAssociative(env *Env, obj Object, msg string) (Associative, error) {
	switch c := obj.(type) {
	case Associative:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Associative", obj.GetType().Name())
		}
		var v Associative
		return v, env.NewError(msg)
	}
}

func EnsureAssociative(env *Env, args []Object, index int) (Associative, error) {
	if len(args) <= index {
		var t Associative
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Associative:
		return c, nil
	default:
		var v Associative
		return v, env.NewArgTypeError(index, c, "Associative")
	}
}

func AssertReversible(env *Env, obj Object, msg string) (Reversible, error) {
	switch c := obj.(type) {
	case Reversible:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Reversible", obj.GetType().Name())
		}
		var v Reversible
		return v, env.NewError(msg)
	}
}

func EnsureReversible(env *Env, args []Object, index int) (Reversible, error) {
	if len(args) <= index {
		var t Reversible
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Reversible:
		return c, nil
	default:
		var v Reversible
		return v, env.NewArgTypeError(index, c, "Reversible")
	}
}

func AssertNamed(env *Env, obj Object, msg string) (Named, error) {
	switch c := obj.(type) {
	case Named:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Named", obj.GetType().Name())
		}
		var v Named
		return v, env.NewError(msg)
	}
}

func EnsureNamed(env *Env, args []Object, index int) (Named, error) {
	if len(args) <= index {
		var t Named
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Named:
		return c, nil
	default:
		var v Named
		return v, env.NewArgTypeError(index, c, "Named")
	}
}

func AssertComparator(env *Env, obj Object, msg string) (Comparator, error) {
	switch c := obj.(type) {
	case Comparator:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Comparator", obj.GetType().Name())
		}
		var v Comparator
		return v, env.NewError(msg)
	}
}

func EnsureComparator(env *Env, args []Object, index int) (Comparator, error) {
	if len(args) <= index {
		var t Comparator
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Comparator:
		return c, nil
	default:
		var v Comparator
		return v, env.NewArgTypeError(index, c, "Comparator")
	}
}

func AssertRatio(env *Env, obj Object, msg string) (*Ratio, error) {
	switch c := obj.(type) {
	case *Ratio:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Ratio", obj.GetType().Name())
		}
		var v *Ratio
		return v, env.NewError(msg)
	}
}

func EnsureRatio(env *Env, args []Object, index int) (*Ratio, error) {
	if len(args) <= index {
		var t *Ratio
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case *Ratio:
		return c, nil
	default:
		var v *Ratio
		return v, env.NewArgTypeError(index, c, "Ratio")
	}
}

func AssertNamespace(env *Env, obj Object, msg string) (*Namespace, error) {
	switch c := obj.(type) {
	case *Namespace:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Namespace", obj.GetType().Name())
		}
		var v *Namespace
		return v, env.NewError(msg)
	}
}

func EnsureNamespace(env *Env, args []Object, index int) (*Namespace, error) {
	if len(args) <= index {
		var t *Namespace
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case *Namespace:
		return c, nil
	default:
		var v *Namespace
		return v, env.NewArgTypeError(index, c, "Namespace")
	}
}

func AssertVar(env *Env, obj Object, msg string) (*Var, error) {
	switch c := obj.(type) {
	case *Var:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Var", obj.GetType().Name())
		}
		var v *Var
		return v, env.NewError(msg)
	}
}

func EnsureVar(env *Env, args []Object, index int) (*Var, error) {
	if len(args) <= index {
		var t *Var
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case *Var:
		return c, nil
	default:
		var v *Var
		return v, env.NewArgTypeError(index, c, "Var")
	}
}

func AssertError(env *Env, obj Object, msg string) (Error, error) {
	switch c := obj.(type) {
	case Error:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Error", obj.GetType().Name())
		}
		var v Error
		return v, env.NewError(msg)
	}
}

func EnsureError(env *Env, args []Object, index int) (Error, error) {
	if len(args) <= index {
		var t Error
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Error:
		return c, nil
	default:
		var v Error
		return v, env.NewArgTypeError(index, c, "Error")
	}
}

func AssertFn(env *Env, obj Object, msg string) (*Fn, error) {
	switch c := obj.(type) {
	case *Fn:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Fn", obj.GetType().Name())
		}
		var v *Fn
		return v, env.NewError(msg)
	}
}

func EnsureFn(env *Env, args []Object, index int) (*Fn, error) {
	if len(args) <= index {
		var t *Fn
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case *Fn:
		return c, nil
	default:
		var v *Fn
		return v, env.NewArgTypeError(index, c, "Fn")
	}
}

func AssertDeref(env *Env, obj Object, msg string) (Deref, error) {
	switch c := obj.(type) {
	case Deref:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Deref", obj.GetType().Name())
		}
		var v Deref
		return v, env.NewError(msg)
	}
}

func EnsureDeref(env *Env, args []Object, index int) (Deref, error) {
	if len(args) <= index {
		var t Deref
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Deref:
		return c, nil
	default:
		var v Deref
		return v, env.NewArgTypeError(index, c, "Deref")
	}
}

func AssertAtom(env *Env, obj Object, msg string) (*Atom, error) {
	switch c := obj.(type) {
	case *Atom:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Atom", obj.GetType().Name())
		}
		var v *Atom
		return v, env.NewError(msg)
	}
}

func EnsureAtom(env *Env, args []Object, index int) (*Atom, error) {
	if len(args) <= index {
		var t *Atom
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case *Atom:
		return c, nil
	default:
		var v *Atom
		return v, env.NewArgTypeError(index, c, "Atom")
	}
}

func AssertRef(env *Env, obj Object, msg string) (Ref, error) {
	switch c := obj.(type) {
	case Ref:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Ref", obj.GetType().Name())
		}
		var v Ref
		return v, env.NewError(msg)
	}
}

func EnsureRef(env *Env, args []Object, index int) (Ref, error) {
	if len(args) <= index {
		var t Ref
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Ref:
		return c, nil
	default:
		var v Ref
		return v, env.NewArgTypeError(index, c, "Ref")
	}
}

func AssertKVReduce(env *Env, obj Object, msg string) (KVReduce, error) {
	switch c := obj.(type) {
	case KVReduce:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "KVReduce", obj.GetType().Name())
		}
		var v KVReduce
		return v, env.NewError(msg)
	}
}

func EnsureKVReduce(env *Env, args []Object, index int) (KVReduce, error) {
	if len(args) <= index {
		var t KVReduce
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case KVReduce:
		return c, nil
	default:
		var v KVReduce
		return v, env.NewArgTypeError(index, c, "KVReduce")
	}
}

func AssertPending(env *Env, obj Object, msg string) (Pending, error) {
	switch c := obj.(type) {
	case Pending:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Pending", obj.GetType().Name())
		}
		var v Pending
		return v, env.NewError(msg)
	}
}

func EnsurePending(env *Env, args []Object, index int) (Pending, error) {
	if len(args) <= index {
		var t Pending
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case Pending:
		return c, nil
	default:
		var v Pending
		return v, env.NewArgTypeError(index, c, "Pending")
	}
}

func AssertFile(env *Env, obj Object, msg string) (*File, error) {
	switch c := obj.(type) {
	case *File:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "File", obj.GetType().Name())
		}
		var v *File
		return v, env.NewError(msg)
	}
}

func EnsureFile(env *Env, args []Object, index int) (*File, error) {
	if len(args) <= index {
		var t *File
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case *File:
		return c, nil
	default:
		var v *File
		return v, env.NewArgTypeError(index, c, "File")
	}
}

func Assertio_Reader(env *Env, obj Object, msg string) (io.Reader, error) {
	switch c := obj.(type) {
	case io.Reader:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "io.Reader", obj.GetType().Name())
		}
		var v io.Reader
		return v, env.NewError(msg)
	}
}

func Ensureio_Reader(env *Env, args []Object, index int) (io.Reader, error) {
	if len(args) <= index {
		var t io.Reader
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case io.Reader:
		return c, nil
	default:
		var v io.Reader
		return v, env.NewArgTypeError(index, c, "io.Reader")
	}
}

func Assertio_Writer(env *Env, obj Object, msg string) (io.Writer, error) {
	switch c := obj.(type) {
	case io.Writer:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "io.Writer", obj.GetType().Name())
		}
		var v io.Writer
		return v, env.NewError(msg)
	}
}

func Ensureio_Writer(env *Env, args []Object, index int) (io.Writer, error) {
	if len(args) <= index {
		var t io.Writer
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case io.Writer:
		return c, nil
	default:
		var v io.Writer
		return v, env.NewArgTypeError(index, c, "io.Writer")
	}
}

func AssertStringReader(env *Env, obj Object, msg string) (StringReader, error) {
	switch c := obj.(type) {
	case StringReader:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "StringReader", obj.GetType().Name())
		}
		var v StringReader
		return v, env.NewError(msg)
	}
}

func EnsureStringReader(env *Env, args []Object, index int) (StringReader, error) {
	if len(args) <= index {
		var t StringReader
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case StringReader:
		return c, nil
	default:
		var v StringReader
		return v, env.NewArgTypeError(index, c, "StringReader")
	}
}

func Assertio_RuneReader(env *Env, obj Object, msg string) (io.RuneReader, error) {
	switch c := obj.(type) {
	case io.RuneReader:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "io.RuneReader", obj.GetType().Name())
		}
		var v io.RuneReader
		return v, env.NewError(msg)
	}
}

func Ensureio_RuneReader(env *Env, args []Object, index int) (io.RuneReader, error) {
	if len(args) <= index {
		var t io.RuneReader
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case io.RuneReader:
		return c, nil
	default:
		var v io.RuneReader
		return v, env.NewArgTypeError(index, c, "io.RuneReader")
	}
}

func AssertChannel(env *Env, obj Object, msg string) (*Channel, error) {
	switch c := obj.(type) {
	case *Channel:
		return c, nil
	default:
		if msg == "" {
			msg = fmt.Sprintf("Expected %s, got %s", "Channel", obj.GetType().Name())
		}
		var v *Channel
		return v, env.NewError(msg)
	}
}

func EnsureChannel(env *Env, args []Object, index int) (*Channel, error) {
	if len(args) <= index {
		var t *Channel
		return t, ErrorArity(env, index)
	}

	switch c := args[index].(type) {
	case *Channel:
		return c, nil
	default:
		var v *Channel
		return v, env.NewArgTypeError(index, c, "Channel")
	}
}
