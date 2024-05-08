package core

import (
	"reflect"

	"github.com/lab47/lace/pkg/pkgreflect"
)

func init() {
	pkgreflect.AddPackage("lace.lang", &pkgreflect.Package{
		Doc:   "",
		Types: map[string]pkgreflect.Type{},

		Functions: map[string]pkgreflect.FuncValue{
			"ConcatSimple": {Doc: "Concatinate N sequences together", Args: []pkgreflect.Arg{{Name: "env", Tag: "Env"}, {Name: "args", Tag: "Args"}}, Tag: "any", Value: reflect.ValueOf(ConcatSimple)},

			"Conj": {Doc: "Create a new Sequence by combine the value with the collection.", Args: []pkgreflect.Arg{{Name: "env", Tag: "Env"}, {Name: "col", Tag: "Object"}, {Name: "val", Tag: "Object"}}, Tag: "any", Value: reflect.ValueOf(Conj)},

			"Cons": {Doc: "Add an element to a Seq value, returning a new Seq", Args: []pkgreflect.Arg{{Name: "env", Tag: "Env"}, {Name: "val", Tag: "Object"}, {Name: "seq", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(Cons)},

			"Equals": {Doc: "Compare two values returning a boolean if they are equal or not", Args: []pkgreflect.Arg{{Name: "env", Tag: "Env"}, {Name: "a", Tag: "Object"}, {Name: "b", Tag: "Object"}}, Tag: "any", Value: reflect.ValueOf(Equals)},

			"First": {Doc: "Return the first element in a Seq", Args: []pkgreflect.Arg{{Name: "env", Tag: "Env"}, {Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(First)},

			"List": {Doc: "Create a new lace List from the given arguments", Args: []pkgreflect.Arg{{Name: "env", Tag: "Env"}, {Name: "args", Tag: "Args"}}, Tag: "any", Value: reflect.ValueOf(MakeList)},

			"Next": {Doc: "Return elements other than the first one in a Seq", Args: []pkgreflect.Arg{{Name: "env", Tag: "Env"}, {Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(Next)},

			"PushBindings": {Doc: "Add given bindings to the set of current Var bindings, returning\nthe original set.", Args: []pkgreflect.Arg{{Name: "env", Tag: "Env"}, {Name: "assoc", Tag: "Map"}}, Tag: "any", Value: reflect.ValueOf(PushBindings)},

			"Rest": {Doc: "Return all elements of a seq except for the first one.", Args: []pkgreflect.Arg{{Name: "env", Tag: "Env"}, {Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(Rest)},

			"Seq": {Doc: "Convert the given value to a Seq", Args: []pkgreflect.Arg{{Name: "env", Tag: "Env"}, {Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(ConvertToSeq)},

			"SetBindings": {Doc: "Reset the local var bindings to the given value.", Args: []pkgreflect.Arg{{Name: "env", Tag: "Env"}, {Name: "assoc", Tag: "Associative"}}, Tag: "any", Value: reflect.ValueOf(SetBindings)},
		},

		Variables: map[string]pkgreflect.Value{},

		Consts: map[string]pkgreflect.Value{},
	})
}
