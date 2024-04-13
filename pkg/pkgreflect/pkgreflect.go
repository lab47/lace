package pkgreflect

import "reflect"

type Value struct {
	Doc   string
	Value reflect.Value
}

type Arg struct {
	Name string
	Tag  string
}

type Func struct {
	Doc  string
	Tag  string
	Args []Arg
}

type FuncValue struct {
	Doc  string
	Tag  string
	Args []Arg

	Value reflect.Value
}

type Type struct {
	Doc     string
	Value   reflect.Type
	Methods map[string]Func
}

type Package struct {
	Name      string
	Doc       string
	Types     map[string]Type
	Functions map[string]FuncValue
	Variables map[string]Value
	Consts    map[string]Value
}

var registry = map[string]*Package{}

func AddPackage(name string, pkg *Package) {
	registry[name] = pkg
}

func FindPackage(name string) *Package {
	return registry[name]
}

func Registry() map[string]*Package {
	return registry
}
