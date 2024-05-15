package core

import (
	"fmt"
	"io"
	"os"
)

var Exit func(rc int)

func SetExit(fn func(rc int)) {
	Exit = fn
}

func writeIndent(w io.Writer, n int) error {
	space := []byte(" ")
	for i := 0; i < n; i++ {
		_, err := w.Write(space)
		if err != nil {
			return err
		}
	}

	return nil
}

func pprintObject(env *Env, obj Object, indent int, w io.Writer) (int, error) {
	switch obj := obj.(type) {
	case Pprinter:
		return obj.Pprint(env, w, indent)
	default:
		s, err := obj.ToString(env, true)
		if err != nil {
			return 0, err
		}
		fmt.Fprint(w, s)
		return indent + len(s), nil
	}
}

func FileInfoMap(env *Env, name string, info os.FileInfo) Map {
	m := EmptyArrayMap()
	m.Add(env, MakeKeyword("name"), MakeString(name))
	m.Add(env, MakeKeyword("size"), MakeInt(int(info.Size())))
	m.Add(env, MakeKeyword("mode"), MakeInt(int(info.Mode())))
	m.Add(env, MakeKeyword("modtime"), MakeTime(info.ModTime()))
	m.Add(env, MakeKeyword("dir?"), MakeBoolean(info.IsDir()))
	return m
}

func ToBool(obj Object) bool {
	switch obj := obj.(type) {
	case Nil:
		return false
	case Boolean:
		return bool(obj)
	default:
		return true
	}
}

func ToNative(env *Env, obj Object) (any, error) {
	switch sv := obj.(type) {
	case Nil:
		return nil, nil
	case Boolean:
		return bool(sv), nil
	case Int:
		return sv.Int().I(), nil
	case String:
		return sv.S(), nil
	default:
		return obj.ToString(env, false)
	}
}
