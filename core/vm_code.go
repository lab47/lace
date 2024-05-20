package core

import (
	"bytes"
	"fmt"
	"regexp"
	"sync/atomic"

	"github.com/davecgh/go-spew/spew"
	"github.com/fxamacker/cbor/v2"
)

var nextFnId atomic.Int64

type Code struct {
	fnId        int64
	numBindings int

	importUpvals int
	totalUpvals  int

	filename string
	lines    []int

	macroLines []int

	files      []string
	fileFromIp []int

	importBindings []Symbol

	stackSize uint32

	data CodeData
}

func (c *Code) Position() string {
	return fmt.Sprintf("%s:%d", c.filename, c.lineForIp(0))
}

func (c *Code) lineForIp(ip int) int {
	for i := 0; i < len(c.lines); i += 2 {
		if c.lines[i] <= ip && c.lines[i+2] > ip {
			return c.lines[i+1]
		}
	}

	return -1
}

func (c *Code) macroLineForIp(ip int) int {
	for i := 0; i < len(c.macroLines); i += 2 {
		if c.macroLines[i] <= ip && c.macroLines[i+2] > ip {
			return c.macroLines[i+1]
		}
	}

	return -1
}

func (c *Code) fileForIp(ip int) string {
	for i := 0; i < len(c.fileFromIp); i += 2 {
		if c.fileFromIp[i] <= ip && c.fileFromIp[i+2] > ip {
			return c.files[c.fileFromIp[i+1]]
		}
	}

	return c.filename
}

type CodePosition struct {
	StartLine   int    `json:"start_line" cbor:"1,keyasint"`
	EndLine     int    `json:"end_line" cbor:"2,keyasint"`
	StartColumn int    `json:"start_column" cbor:"3,keyasint"`
	EndColumn   int    `json:"end_column" cbor:"4,keyasint"`
	Filename    string `json:"filename" cbor:"5,keyasint"`
}

func (c *CodePosition) Position() Position {
	pos := Position{
		startLine:   c.StartLine,
		endLine:     c.EndLine,
		startColumn: c.StartColumn,
		endColumn:   c.EndColumn,
	}

	if c.Filename != "" {
		pos.filename = c.Filename
	}

	return pos
}

func (c *CodePosition) Set(pos Position) {
	c.StartLine = pos.startLine
	c.EndLine = pos.endLine
	c.StartColumn = pos.startColumn
	c.EndColumn = pos.endColumn
	c.Filename = pos.filename
}

type CodeSymbol struct {
	Position  *CodePosition
	Namespace string `json:"ns" cbor:"1,keyasint"`
	Name      string `json:"name" cbor:"2,keyasint"`
}

func (c *CodeSymbol) Symbol() Symbol {
	sym := AssembleSymbol(c.Namespace, c.Name)
	if c.Position != nil {
		sym = SymbolSetInfo(sym, &ObjectInfo{
			Position: c.Position.Position(),
		})
	}

	return sym
}

func (c *CodeSymbol) Set(sym Symbol) {
	info := GetInfo(sym)
	if info != nil {
		c.Position = &CodePosition{}
		c.Position.Set(info.Position)
	}
	c.Name = sym.Name()
	c.Namespace = sym.Namespace()
}

type CodeKeyword struct {
	Position  *CodePosition
	Namespace string `json:"ns" cbor:"1,keyasint"`
	Name      string `json:"name" cbor:"2,keyasint"`
}

func (c *CodeKeyword) Keyword() Keyword {
	kw := NewKeyword(c.Namespace, c.Name)
	if c.Position != nil {
		kw = SetInfo(kw, &ObjectInfo{
			Position: c.Position.Position(),
		}).(Keyword)
	}

	return kw
}

func (c *CodeKeyword) Set(kw Keyword) {
	info := GetInfo(kw)
	if info != nil {
		c.Position = &CodePosition{}
		c.Position.Set(info.Position)
	}
	c.Name = kw.Name()
	c.Namespace = kw.Namespace()
}

type CodeVar struct {
	Name CodeSymbol `json:"name" cbor:"1,keyasint"`
}

type CodeType struct {
	Name string `json:"name" cbor:"1,keyasint"`
}

type CodeEncoded struct {
	Data []byte `json:"data" cbor:"1,keyasint"`
}

type CodeString struct {
	String string `json:"string" cbor:"1,keyasint"`
}

type CodeMapPair struct {
	Key   CodeLiteral `json:"key" cbor:"1,keyasint"`
	Value CodeLiteral `json:"value" cbor:"2,keyasint"`
}

type CodeMap struct {
	Pairs []CodeMapPair `json:"pairs" cbor:"1,keyasint"`
}

type CodeInt struct {
	Int int64 `json:"int" cbor:"1,keyasint"`
}

type CodeSeq struct {
	Elements []CodeLiteral `json:"elements" cbor:"1,keyasint"`
}

type CodeVector struct {
	Elements []CodeLiteral `json:"elements" cbor:"1,keyasint"`
}

type CodeNamespace struct {
	Name CodeSymbol `json:"name" cbor:"1,keyasint"`
}

type CodeBoolean struct {
	Bool bool `json:"bool" cbor:"1,keyasint"`
}

type CodeChar struct {
	Ch rune `json:"char" cbor:"1,keyasint"`
}

type CodeRegex struct {
	Regex string `json:"regex" cbor:"1,keyasint"`
}

type CodeDouble struct {
	Double float64 `json:"double" cbor:"1,keyasint"`
}

type CodeLiteral struct {
	Symbol  *CodeSymbol    `json:"symbol,omitempty" cbor:"1,keyasint,omitempty"`
	Var     *CodeVar       `json:"var,omitempty" cbor:"2,keyasint,omitempty"`
	Type    *CodeType      `json:"type,omitempty" cbor:"3,keyasint,omitempty"`
	Encoded *CodeEncoded   `json:"encoded,omitempty" cbor:"4,keyasint,omitempty"`
	String  *CodeString    `json:"string,omitempty" cbor:"5,keyasint,omitempty"`
	Keyword *CodeKeyword   `json:"keyword,omitempty" cbor:"6,keyasint,omitempty"`
	Map     *CodeMap       `json:"map,omitempty" cbor:"7,keyasint,omitempty"`
	Int     *CodeInt       `json:"int,omitempty" cbor:"8,keyasint,omitempty"`
	Seq     *CodeSeq       `json:"seq,omitempty" cbor:"9,keyasint,omitempty"`
	Vector  *CodeVector    `json:"vector,omitempty" cbor:"10,keyasint,omitempty"`
	NS      *CodeNamespace `json:"ns,omitempty" cbor:"11,keyasint,omitempty"`
	Bool    *CodeBoolean   `json:"bool,omitempty" cbor:"12,keyasint,omitempty"`
	Char    *CodeChar      `json:"char,omitempty" cbor:"13,keyasint,omitempty"`
	Regex   *CodeRegex     `json:"regex,omitempty" cbor:"14,keyasint,omitempty"`
	Double  *CodeDouble    `json:"double,omitempty" cbor:"15,keyasint,omitempty"`
}

func (cl *CodeLiteral) Set(env *Env, lit Object) error {
	switch o := lit.(type) {
	case Boolean:
		cl.Bool = &CodeBoolean{Bool: ToBool(o)}
	case Char:
		cl.Char = &CodeChar{Ch: o.Ch()}
	case Symbol:
		cl.Symbol = &CodeSymbol{}
		cl.Symbol.Set(o)
	case *Var:
		cl.Var = &CodeVar{}

		sym := AssembleSymbol(o.ns.Name.Name(), o.name.String())

		cl.Var.Name.Set(sym)
	case *Type:
		name := o.Name()

		cl.Type = &CodeType{
			Name: name,
		}
	case String:
		cl.String = &CodeString{String: o.S()}
	case *Regex:
		cl.Regex = &CodeRegex{
			Regex: o.R.String(),
		}
	case Keyword:
		cl.Keyword = &CodeKeyword{}
		cl.Keyword.Set(o)
	case Map:
		mi := o.Iter()

		var cm CodeMap

		for mi.HasNext() {
			p := mi.Next()
			var k CodeLiteral
			err := k.Set(env, p.Key)
			if err != nil {
				return err
			}

			var v CodeLiteral
			err = v.Set(env, p.Value)
			if err != nil {
				return err
			}

			cm.Pairs = append(cm.Pairs, CodeMapPair{
				Key:   k,
				Value: v,
			})

			cl.Map = &cm
		}
	case Int:
		cl.Int = &CodeInt{
			Int: o.I64(),
		}
	case Double:
		cl.Double = &CodeDouble{
			Double: o.D,
		}
	case *Vector:
		var cs CodeVector

		for i := 0; i < o.Count(); i++ {
			obj, err := o.Nth(env, i)
			if err != nil {
				return err
			}

			var cl CodeLiteral
			err = cl.Set(env, obj)
			if err != nil {
				return err
			}

			cs.Elements = append(cs.Elements, cl)
		}

		cl.Vector = &cs
	case Seq:
		i := iter(o)

		var cs CodeSeq

		for i.HasNext(env) {
			obj, err := i.Next(env)
			if err != nil {
				return err
			}

			var cl CodeLiteral
			err = cl.Set(env, obj)
			if err != nil {
				return err
			}

			cs.Elements = append(cs.Elements, cl)
		}

		cl.Seq = &cs
	case *Namespace:
		var cs CodeSymbol
		cs.Set(o.Name)

		cl.NS = &CodeNamespace{
			Name: cs,
		}
	default:
		fmt.Printf("type as string: %T\n", o)
		var buf bytes.Buffer
		PrintObject(env, o, &buf)

		cl.Encoded = &CodeEncoded{
			Data: buf.Bytes(),
		}
	}

	return nil
}

func (lit *CodeLiteral) AsValue(env *Env) (Object, error) {
	switch {
	case lit.Bool != nil:
		return MakeBoolean(lit.Bool.Bool), nil
	case lit.Symbol != nil:
		return lit.Symbol.Symbol(), nil
	case lit.Int != nil:
		return MakeInt(int(lit.Int.Int)), nil
	case lit.Double != nil:
		return MakeDouble(lit.Double.Double), nil
	case lit.Char != nil:
		return NewChar(lit.Char.Ch), nil
	case lit.Var != nil:
		ns := env.FindNamespace(MakeSymbol(lit.Var.Name.Namespace))
		if ns == nil {
			return nil, fmt.Errorf("unknown ns")
		}

		vr := ns.Resolve(lit.Var.Name.Name)
		if vr == nil {
			return nil, fmt.Errorf("missing var 2: %s", lit.Var.Name.Symbol().String())
		}

		return vr, nil
	case lit.Type != nil:
		name := lit.Type.Name
		return TYPES[name], nil
	case lit.String != nil:
		return MakeString(lit.String.String), nil
	case lit.Keyword != nil:
		return lit.Keyword.Keyword(), nil
	case lit.Regex != nil:
		re, err := regexp.Compile(lit.Regex.Regex)
		if err != nil {
			return nil, err
		}

		return MakeRegex(re), nil
	case lit.Encoded != nil:
		obj, err := readFromReader(env, bytes.NewReader(lit.Encoded.Data))
		if err != nil {
			return nil, err
		}
		return obj, nil
	case lit.Map != nil:
		var ret Associative = EmptyArrayMap()

		for _, p := range lit.Map.Pairs {
			k, err := p.Key.AsValue(env)
			if err != nil {
				return nil, err
			}

			v, err := p.Value.AsValue(env)
			if err != nil {
				return nil, err
			}

			ret, err = ret.Assoc(env, k, v)
			if err != nil {
				return nil, err
			}
		}

		return ret, nil
	case lit.Seq != nil:
		var objs []Object

		for _, cl := range lit.Seq.Elements {
			obj, err := cl.AsValue(env)
			if err != nil {
				return nil, err
			}

			objs = append(objs, obj)
		}

		if len(objs) >= 3 {
			vec := NewVectorFrom(objs...)
			return vec.Seq(), nil
		}

		return NewListFrom(objs...), nil
	case lit.Vector != nil:
		var objs []Object

		for _, cl := range lit.Vector.Elements {
			obj, err := cl.AsValue(env)
			if err != nil {
				return nil, err
			}

			objs = append(objs, obj)
		}

		return NewVectorFrom(objs...), nil
	case lit.NS != nil:
		ns := env.FindNamespace(lit.NS.Name.Symbol())
		if ns == nil {
			return nil, fmt.Errorf("unable to find existing namespace: %s", lit.NS.Name.Symbol())
		}

		return ns, nil
	default:
		spew.Dump(lit)
		return nil, fmt.Errorf("unsupported literal in code")
	}

}

type CodeMethod struct {
	Method string `json:"method" cbor:"1,keyasint"`
	Arity  uint   `json:"arity" cbor:"2,keyasint"`
}

type CodeAsData struct {
	NumBindings    int            `json:"num_bindings" cbor:"1,keyasint,omitempty"`
	ImportUpvals   int            `json:"import_upvals" cbor:"2,keyasint,omitempty"`
	TotalUpvals    int            `json:"total_upvals" cbor:"3,keyasint,omitempty"`
	Filename       string         `json:"filename" cbor:"4,keyasint,omitempty"`
	Lines          []int          `json:"lines" cbor:"5,keyasint,omitempty"`
	ImportBindings []string       `json:"import_bindings" cbor:"6,keyasint,omitempty"`
	VarNames       []CodeSymbol   `json:"var_names" cbor:"7,keyasint,omitempty"`
	DefVarNames    []string       `json:"def_var_names" cbor:"8,keyasint,omitempty"`
	Literals       []*CodeLiteral `json:"literals" cbor:"9,keyasint,omitempty"`
	Codes          []*CodeAsData  `json:"codes" cbor:"10,keyasint,omitempty"`
	Methods        []CodeMethod   `json:"methods" cbor:"11,keyasint,omitempty"`
	Instructions   []uint32       `json:"instructions" cbor:"12,keyasint,omitempty"`
	StackSize      uint32         `json:"stack_size" cbor:"13,keyasint,omitempty"`
}

func (c *Code) AsData(env *Env) (*CodeAsData, error) {
	cad := &CodeAsData{
		NumBindings:  c.numBindings,
		ImportUpvals: c.importUpvals,
		TotalUpvals:  c.totalUpvals,
		Filename:     c.filename,
		Lines:        c.lines,
	}

	for _, sym := range c.importBindings {
		cad.ImportBindings = append(cad.ImportBindings, sym.Name())
	}

	for _, sym := range c.data.varNames {
		var cs CodeSymbol
		cs.Set(sym)

		cad.VarNames = append(cad.VarNames, cs)
	}

	for _, sym := range c.data.defVarNames {
		cad.DefVarNames = append(cad.DefVarNames, sym.Name())
	}

	for _, lit := range c.data.literals {
		var cl CodeLiteral

		err := cl.Set(env, lit)
		if err != nil {
			return nil, err
		}

		cad.Literals = append(cad.Literals, &cl)
	}

	for _, sub := range c.data.codes {
		subcad, err := sub.AsData(env)
		if err != nil {
			return nil, err
		}

		cad.Codes = append(cad.Codes, subcad)
	}

	for _, meth := range c.data.methods {
		cad.Methods = append(cad.Methods, CodeMethod(meth))
	}

	cad.Instructions = c.data.insns

	cad.StackSize = c.stackSize

	return cad, nil
}

func (cad *CodeAsData) AsCode(env *Env) (*Code, error) {
	c := &Code{
		fnId:         nextFnId.Add(1),
		numBindings:  cad.NumBindings,
		importUpvals: cad.ImportUpvals,
		totalUpvals:  cad.TotalUpvals,
		filename:     cad.Filename,
		lines:        cad.Lines,
		stackSize:    cad.StackSize,
	}

	for _, str := range cad.ImportBindings {
		c.importBindings = append(c.importBindings, MakeSymbol(str))
	}

	for _, str := range cad.DefVarNames {
		sym := MakeSymbol(str)
		c.data.defVarNames = append(c.data.defVarNames, sym)

		ns := env.CurrentNamespace()
		vr, err := ns.Intern(env, sym)
		if err != nil {
			return nil, err
		}

		c.data.defVars = append(c.data.defVars, vr)
	}

	for _, csym := range cad.VarNames {
		ns := env.EnsureNamespace(MakeSymbol(csym.Namespace))
		if ns == nil {
			panic("bad ns: " + csym.Namespace)
		}

		sym := csym.Symbol()
		c.data.varNames = append(c.data.varNames, sym)

		vr := ns.Resolve(csym.Name)
		if vr == nil {
			return nil, fmt.Errorf("missing var 1: %s", sym.String())
		}

		c.data.vars = append(c.data.vars, vr)
	}

	for _, lit := range cad.Literals {
		obj, err := lit.AsValue(env)
		if err != nil {
			return nil, err
		}

		c.data.literals = append(c.data.literals, obj)
	}

	for _, sub := range cad.Codes {
		subC, err := sub.AsCode(env)
		if err != nil {
			return nil, err
		}
		c.data.codes = append(c.data.codes, subC)
	}

	for _, meth := range cad.Methods {
		c.data.methods = append(c.data.methods, MethodSite(meth))
	}

	c.data.insns = cad.Instructions

	return c, nil
}

func MarshalCode(env *Env, code *Code) ([]byte, error) {
	data, err := code.AsData(env)
	if err != nil {
		return nil, err
	}
	return cbor.Marshal(data)
}

func UnmarshalCode(env *Env, data []byte) (*Fn, error) {
	var cad CodeAsData

	err := cbor.Unmarshal(data, &cad)
	if err != nil {
		return nil, err
	}

	code, err := cad.AsCode(env)
	if err != nil {
		return nil, err
	}

	fn := &Fn{
		code:           code,
		importedUpvals: make([]*NamedPair, code.totalUpvals),
	}

	return fn, nil
}
