package core

import (
	"reflect"

	"github.com/lab47/lace/pkg/pkgreflect"
)

type IntegerImpl struct {
	I64Fn         func() int64
	integerTypeFn func() string
}

func (s *IntegerImpl) I64() int64 {
	return s.I64Fn()
}
func (s *IntegerImpl) integerType() string {
	return s.integerTypeFn()
}

func init() {
	Int_methods := map[string]pkgreflect.Func{}
	Integer_methods := map[string]pkgreflect.Func{}
	pkgreflect.AddPackage("lace.lang", &pkgreflect.Package{
		Doc: "",
		Types: map[string]pkgreflect.Type{
			"Int":         {Doc: "The host int value", Value: reflect.TypeOf((*Int)(nil)).Elem(), Methods: Int_methods},
			"Integer":     {Doc: "The Common integer type (can be Int or BigInt)", Value: reflect.TypeOf((*Integer)(nil)).Elem(), Methods: Integer_methods},
			"IntegerImpl": {Doc: `Struct version of interface Integer for implementation`, Value: reflect.TypeFor[IntegerImpl]()},
		},

		Functions: map[string]pkgreflect.FuncValue{
			"CombineToString": {Doc: "Combine many values into a single string.", Args: []pkgreflect.Arg{{Name: "args", Tag: "[]Object"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(CombineToString))},

			"ConcatSimple": {Doc: "Concatinate N sequences together", Args: []pkgreflect.Arg{{Name: "args", Tag: "[]Object"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(ConcatSimple))},

			"Conj": {Doc: "Create a new Sequence by combine the value with the collection.", Args: []pkgreflect.Arg{{Name: "col", Tag: "Object"}, {Name: "val", Tag: "Object"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc3_2(Conj))},

			"Cons": {Doc: "Add an element to a Seq value, returning a new Seq", Args: []pkgreflect.Arg{{Name: "val", Tag: "Object"}, {Name: "seq", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc3_2(Cons))},

			"Equals": {Doc: "Compare two values returning a boolean if they are equal or not", Args: []pkgreflect.Arg{{Name: "a", Tag: "Object"}, {Name: "b", Tag: "Object"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc3_2(EqualsValues))},

			"First": {Doc: "Return the first element in a Seq", Args: []pkgreflect.Arg{{Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(First))},

			"List": {Doc: "Create a new lace List from the given arguments", Args: []pkgreflect.Arg{{Name: "args", Tag: "[]Object"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(MakeList))},

			"LoadLibFromPath": {Doc: "Attempt to load a given lib from a given path.", Args: []pkgreflect.Arg{{Name: "libnamev", Tag: "Symbol"}, {Name: "pathnamev", Tag: "String"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc3_2(LoadLibFromPath))},

			"NewFuture": {Doc: "NewFuture creates a new Future value and schedules the future\nto be run. Deref'ing the Future will retrieve the value (potentially\nwaiting if the value is not yet ready)", Args: []pkgreflect.Arg{{Name: "call", Tag: "Callable"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(NewFuture))},

			"Next": {Doc: "Return elements other than the first one in a Seq", Args: []pkgreflect.Arg{{Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(Next))},

			"PushBindings": {Doc: "Add given bindings to the set of current Var bindings, returning\nthe original set.", Args: []pkgreflect.Arg{{Name: "assoc", Tag: "Map"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(PushBindings))},

			"Rest": {Doc: "Return all elements of a seq except for the first one.", Args: []pkgreflect.Arg{{Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(Rest))},

			"Seq": {Doc: "Convert the given value to a Seq", Args: []pkgreflect.Arg{{Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(ConvertToSeq))},

			"SetBindings": {Doc: "Reset the local var bindings to the given value.", Args: []pkgreflect.Arg{{Name: "assoc", Tag: "Associative"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(SetBindings))},

			"StartGoRoutine": {Doc: "StartGoRoutine runs the given callable in a new goroutine, returning a channel\nthat can be used to retrieve the return value.", Args: []pkgreflect.Arg{{Name: "callable", Tag: "Callable"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(StartGoRoutine))},
		},

		Variables: map[string]pkgreflect.Value{},

		Consts: map[string]pkgreflect.Value{},
	})
}
