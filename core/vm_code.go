package core

import (
	"bytes"
	"fmt"
	"sync/atomic"

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
		sym.info = &ObjectInfo{
			Position: c.Position.Position(),
		}
	}

	return sym
}

func (c *CodeSymbol) Set(sym Symbol) {
	if sym.info != nil {
		c.Position = &CodePosition{}
		c.Position.Set(sym.info.Position)
	}
	c.Name = sym.Name()
	c.Namespace = sym.Namespace()
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

type CodeLiteral struct {
	Symbol  *CodeSymbol  `json:"symbol" cbor:"1,keyasint,omitempty"`
	Var     *CodeVar     `json:"var" cbor:"2,keyasint,omitempty"`
	Type    *CodeType    `json:"type" cbor:"3,keyasint,omitempty"`
	Encoded *CodeEncoded `json:"encoded" cbor:"4,keyasint,omitempty"`
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

func (c *Code) AsData(env *Env) *CodeAsData {
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

		switch o := lit.(type) {
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
		default:
			var buf bytes.Buffer
			PrintObject(env, o, &buf)

			cl.Encoded = &CodeEncoded{
				Data: buf.Bytes(),
			}
		}

		cad.Literals = append(cad.Literals, &cl)
	}

	for _, sub := range c.data.codes {
		cad.Codes = append(cad.Codes, sub.AsData(env))
	}

	for _, meth := range c.data.methods {
		cad.Methods = append(cad.Methods, CodeMethod(meth))
	}

	cad.Instructions = c.data.insns

	cad.StackSize = c.stackSize

	return cad
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
		switch {
		case lit.Symbol != nil:
			c.data.literals = append(c.data.literals, lit.Symbol.Symbol())
		case lit.Var != nil:
			ns := env.FindNamespace(MakeSymbol(lit.Var.Name.Namespace))
			if ns == nil {
				panic("unknown ns")
			}

			vr := ns.Resolve(lit.Var.Name.Name)
			if vr == nil {
				return nil, fmt.Errorf("missing var 2: %s", lit.Var.Name.Symbol().String())
			}

			c.data.literals = append(c.data.literals, vr)
		case lit.Type != nil:
			name := lit.Type.Name
			c.data.literals = append(c.data.literals, TYPES[name])
		case lit.Encoded != nil:
			obj, err := readFromReader(env, bytes.NewReader(lit.Encoded.Data))
			if err != nil {
				return nil, err
			}
			c.data.literals = append(c.data.literals, obj)
		default:
			panic("bad literal")
		}
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
	data := code.AsData(env)
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
