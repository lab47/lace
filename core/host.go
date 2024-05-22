package core

func ExtractTagFromMeta(obj any) (Symbol, bool) {
	if m := GetMeta(obj); m != nil {
		if ok, typeName := m.GetEqu(criticalKeywords.tag); ok {
			if typeSym, ok := typeName.(Symbol); ok {
				return typeSym, true
			}
		}
	}

	return nil, false
}

func (env *Env) ResolveTypeFroMeta(obj any) any {
	sym, ok := ExtractTagFromMeta(obj)
	if !ok {
		return nil
	}

	return env.ResolveType(sym)
}

func (env *Env) ResolveType(sym Symbol) any {
	if sym.Namespace() == "" {
		if vr := env.LangNamespace.Resolve(sym.Name()); vr != nil {
			obj := vr.Resolve(env)
			switch obj.(type) {
			case Type:
				return obj
			}
		}
	}

	return nil
}
