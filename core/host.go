package core

import "fmt"

func ExtractTagFromMeta(obj Object) (Symbol, bool) {
	if m := GetMeta(obj); m != nil {
		if ok, typeName := m.GetEqu(criticalKeywords.tag); ok {
			if typeSym, ok := typeName.(Symbol); ok {
				return typeSym, true
			}
		}
	}

	return nil, false
}

func (env *Env) ResolveTypeFroMeta(obj Object) Object {
	sym, ok := ExtractTagFromMeta(obj)
	if !ok {
		return nil
	}

	return env.ResolveType(sym)
}

func (env *Env) ResolveType(sym Symbol) Object {
	if sym.Namespace() == "" {
		if vr := env.LangNamespace.Resolve(sym.Name()); vr != nil {
			obj := vr.Resolve(env)
			switch obj.(type) {
			case *Type:
				fmt.Printf("Found in lace.lang: %s: %#v\n", sym.Name(), obj)
				return obj
			case *ReflectType:
				fmt.Printf("Found in lace.lang: %s: %#v\n", sym.Name(), obj)
				return obj
			}
		}
	}

	if t := TYPES[sym.Name()]; t != nil {
		return t
	}

	return nil
}
