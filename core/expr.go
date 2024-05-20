package core

import "reflect"

func dumpPosition(p Position) Map {
	res := EmptyArrayMap()
	res.AddEqu(criticalKeywords.startLine, MakeInt(p.startLine))
	res.AddEqu(criticalKeywords.endLine, MakeInt(p.endLine))
	res.AddEqu(criticalKeywords.startColumn, MakeInt(p.startColumn))
	res.AddEqu(criticalKeywords.endColumn, MakeInt(p.endColumn))
	res.AddEqu(criticalKeywords.filename, MakeString(p.Filename()))
	return res
}

func exprArrayMap(expr Expr, exprType string, pos bool) *ArrayMap {
	res := EmptyArrayMap()
	res.AddEqu(criticalKeywords.type_, MakeKeyword(exprType))
	if pos {
		res.AddEqu(criticalKeywords.pos, dumpPosition(expr.Pos()))
	}
	return res
}

func addVector(res *ArrayMap, body []Expr, name string, pos bool) {
	b := EmptyVector()
	for _, e := range body {
		b, _ = b.Conjoin(e.Dump(pos))
	}
	res.AddEqu(MakeKeyword(name), b)
}

func (expr *LiteralExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "literal", pos)
	res.AddEqu(criticalKeywords.object, expr.obj)
	return res
}

func (expr *VectorExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "vector", pos)
	addVector(res, expr.v, "vector", pos)
	return res
}

func (expr *MapExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "map", pos)
	addVector(res, expr.keys, "keys", pos)
	addVector(res, expr.values, "values", pos)
	return res
}

func (expr *SetExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "set", pos)
	addVector(res, expr.elements, "set", pos)
	return res
}

func (expr *IfExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "if", pos)
	res.AddEqu(MakeKeyword("condition"), expr.cond.Dump(pos))
	res.AddEqu(MakeKeyword("positive"), expr.positive.Dump(pos))
	res.AddEqu(MakeKeyword("negative"), expr.negative.Dump(pos))
	return res
}

func (expr *DefExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "def", pos)
	res.AddEqu(criticalKeywords.var_, expr.vr)
	res.AddEqu(criticalKeywords.name, expr.name)
	if expr.value != nil {
		res.AddEqu(criticalKeywords.value, expr.value.Dump(pos))
	}
	if expr.meta != nil {
		res.AddEqu(criticalKeywords.meta, expr.meta.Dump(pos))
	}
	return res
}

func (expr *CallExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "call", pos)
	res.AddEqu(MakeKeyword("name"), MakeString(expr.Name()))
	res.AddEqu(MakeKeyword("callable"), expr.callable.Dump(pos))
	addVector(res, expr.args, "args", pos)
	return res
}

type MethodExpr struct {
	Position
	name   Symbol
	method string

	obj  Expr
	args []Expr

	// inline cache
	lastType reflect.Type
	lastFn   ProcFn
}

var _ Expr = &MethodExpr{}

func (expr *MethodExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "call", pos)
	res.AddEqu(MakeKeyword("method"), MakeString(expr.method))
	res.AddEqu(MakeKeyword("object"), expr.obj.Dump(pos))
	addVector(res, expr.args, "args", pos)
	return res
}

func (expr *MacroCallExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "macro-call", pos)
	res.AddEqu(MakeKeyword("name"), MakeString(expr.name))
	args := EmptyVector()
	for _, arg := range expr.args {
		args, _ = args.Conjoin(arg)
	}
	res.AddEqu(MakeKeyword("args"), args)
	return res
}

func (expr *RecurExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "recur", pos)
	addVector(res, expr.args, "args", pos)
	return res
}

func (expr *VarRefExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "var-ref", pos)
	res.AddEqu(criticalKeywords.var_, expr.vr)
	return res
}

func (expr *SetMacroExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "set-macro", pos)
	res.AddEqu(criticalKeywords.var_, expr.vr)
	return res
}

func (expr *BindingExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "binding", pos)
	res.AddEqu(MakeKeyword("name"), expr.binding.name)
	return res
}

func (expr *MetaExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "meta", pos)
	res.AddEqu(criticalKeywords.meta, expr.meta.Dump(pos))
	res.AddEqu(MakeKeyword("expr"), expr.expr.Dump(pos))
	return res
}

func (expr *DoExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "do", pos)
	addVector(res, expr.body, "body", pos)
	return res
}

func (expr *FnArityExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "arity", pos)
	args := EmptyVector()
	for _, arg := range expr.args {
		args, _ = args.Conjoin(arg)
	}
	res.AddEqu(MakeKeyword("args"), args)
	addVector(res, expr.body, "body", pos)
	return res
}

func (expr *FnExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "fn", pos)
	if expr.self.Name() != "" {
		res.AddEqu(MakeKeyword("self"), expr.self)
	}
	if expr.variadic != nil {
		res.AddEqu(MakeKeyword("variadic"), expr.variadic.Dump(pos))
	}
	arities := EmptyVector()
	for _, a := range expr.arities {
		arities, _ = arities.Conjoin(a.Dump(pos))
	}
	res.AddEqu(MakeKeyword("arities"), arities)
	return res
}

func (expr *LetExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "let", pos)
	names := EmptyVector()
	for _, name := range expr.names {
		names, _ = names.Conjoin(name)
	}
	addVector(res, expr.values, "values", pos)
	addVector(res, expr.body, "body", pos)
	return res
}

func (expr *LoopExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "loop", pos)
	names := EmptyVector()
	for _, name := range expr.names {
		names, _ = names.Conjoin(name)
	}
	addVector(res, expr.values, "values", pos)
	addVector(res, expr.body, "body", pos)
	return res
}

func (expr *ThrowExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "throw", pos)
	res.AddEqu(MakeKeyword("expr"), expr.e.Dump(pos))
	return res
}

func (expr *CatchExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "catch", pos)
	res.AddEqu(MakeKeyword("error-type"), expr.excType)
	res.AddEqu(MakeKeyword("error-symbol"), expr.excSymbol)
	addVector(res, expr.body, "body", pos)
	return res
}

func (expr *TryExpr) Dump(pos bool) Map {
	res := exprArrayMap(expr, "try", pos)
	addVector(res, expr.body, "body", pos)
	addVector(res, expr.finallyExpr, "finally", pos)
	catches := EmptyVector()
	for _, c := range expr.catches {
		catches, _ = catches.Conjoin(c.Dump(pos))
	}
	res.AddEqu(MakeKeyword("catches"), catches)
	return res
}
