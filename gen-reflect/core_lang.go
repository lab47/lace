// Code generated by github.com/lab47/lace/pkg/pkgreflect DO NOT EDIT.
package reflect

import "reflect"
import "github.com/lab47/lace/pkg/pkgreflect"
import lang "github.com/lab47/lace/core/lang"

func init() {
	pkgreflect.AddPackage("lace.lang", &pkgreflect.Package{
		Doc:   "This package contains implementation helpers used to implement lace itself.",
		Types: map[string]pkgreflect.Type{},

		Functions: map[string]pkgreflect.FuncValue{
			"ConcatSimple": {Doc: "", Args: []pkgreflect.Arg{{Name: "args", Tag: "Args"}}, Tag: "any", Value: reflect.ValueOf(lang.ConcatSimple)},

			"Conj": {Doc: "", Args: []pkgreflect.Arg{{Name: "col", Tag: "Object"}, {Name: "val", Tag: "Object"}}, Tag: "any", Value: reflect.ValueOf(lang.Conj)},

			"Cons": {Doc: "Add an element to a Seq value, returning a new Seq", Args: []pkgreflect.Arg{{Name: "val", Tag: "Object"}, {Name: "seq", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(lang.Cons)},

			"Equals": {Doc: "Compare two values returning a boolean if they are equal or not", Args: []pkgreflect.Arg{{Name: "a", Tag: "Object"}, {Name: "b", Tag: "Object"}}, Tag: "any", Value: reflect.ValueOf(lang.Equals)},

			"First": {Doc: "Return the first element in a Seq", Args: []pkgreflect.Arg{{Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(lang.First)},

			"List": {Doc: "Create a new lace List from the given arguments", Args: []pkgreflect.Arg{{Name: "args", Tag: "Args"}}, Tag: "any", Value: reflect.ValueOf(lang.List)},

			"Next": {Doc: "Return elements other than the first one in a Seq", Args: []pkgreflect.Arg{{Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(lang.Next)},

			"Rest": {Doc: "", Args: []pkgreflect.Arg{{Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(lang.Rest)},

			"Seq": {Doc: "", Args: []pkgreflect.Arg{{Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(lang.Seq)},
		},

		Variables: map[string]pkgreflect.Value{},

		Consts: map[string]pkgreflect.Value{},
	})
}
