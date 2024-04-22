package core

import (
	"bytes"
	"fmt"
)

type Code struct {
	numBindings int

	importUpvals int
	totalUpvals  int

	filename string
	lines    []int

	importBindings []Symbol

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

type CodeSymbol struct {
	Namespace string `json:"ns" cbor:"1,keyasint"`
	Name      string `json:"name" cbor:"2,keyasint"`
}

func (c *CodeSymbol) Symbol() Symbol {
	return AssembleSymbol(c.Namespace, c.Name)
}

type CodeVar struct {
	Name CodeSymbol `json:"name" cbor:"1,keyasint"`
}

type CodeType struct {
	Name CodeSymbol `json:"name" cbor:"1,keyasint"`
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
	VarNames       []string       `json:"var_names" cbor:"7,keyasint,omitempty"`
	DefVarNames    []string       `json:"def_var_names" cbor:"8,keyasint,omitempty"`
	Literals       []*CodeLiteral `json:"literals" cbor:"9,keyasint,omitempty"`
	Codes          []*CodeAsData  `json:"codes" cbor:"10,keyasint,omitempty"`
	Methods        []CodeMethod   `json:"methods" cbor:"11,keyasint,omitempty"`
	Instructions   []uint32       `json:"instructions" cbor:"12,keyasint,omitempty"`
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
		cad.VarNames = append(cad.VarNames, sym.String())
	}

	for _, sym := range c.data.defVarNames {
		cad.DefVarNames = append(cad.DefVarNames, sym.Name())
	}

	for _, lit := range c.data.literals {
		var cl CodeLiteral

		switch o := lit.(type) {
		case Symbol:
			cl.Symbol = &CodeSymbol{
				Namespace: o.Namespace(),
				Name:      o.Name(),
			}
		case *Var:
			ns := o.ns.Name.Name()
			name := o.name.Name()

			cl.Var = &CodeVar{
				Name: CodeSymbol{
					Namespace: ns,
					Name:      name,
				},
			}
		case *Type:
			name := o.Name()

			cl.Var = &CodeVar{
				Name: CodeSymbol{
					Name: name,
				},
			}

		default:
			var buf bytes.Buffer
			PrintObject(env, o, &buf)

			cl.Encoded = &CodeEncoded{
				Data: buf.Bytes(),
			}
		}
	}

	for _, sub := range c.data.codes {
		cad.Codes = append(cad.Codes, sub.AsData(env))
	}

	for _, meth := range c.data.methods {
		cad.Methods = append(cad.Methods, CodeMethod{
			Method: meth.Method,
			Arity:  meth.Arity,
		})
	}

	cad.Instructions = c.data.insns

	return cad
}

func (cad *CodeAsData) AsCode(env *Env) (*Code, error) {
	c := &Code{
		numBindings:  cad.NumBindings,
		importUpvals: cad.ImportUpvals,
		totalUpvals:  cad.TotalUpvals,
		filename:     cad.Filename,
		lines:        cad.Lines,
	}

	for _, str := range cad.ImportBindings {
		c.importBindings = append(c.importBindings, MakeSymbol(str))
	}

	for _, str := range cad.VarNames {
		sym := MakeSymbol(str)
		c.data.varNames = append(c.data.varNames, sym)

		vr, ok := env.Resolve(sym)
		if !ok {
			return nil, fmt.Errorf("missing var: %s", sym.String())
		}

		c.data.vars = append(c.data.vars, vr)
	}

	for _, str := range cad.VarNames {
		sym := MakeSymbol(str)
		c.data.varNames = append(c.data.varNames, sym)

		ns := env.CurrentNamespace()
		vr, err := ns.Intern(env, sym)
		if err != nil {
			return nil, err
		}

		c.data.defVars = append(c.data.defVars, vr)
	}

	for _, lit := range cad.Literals {
		switch {
		case lit.Symbol != nil:
			c.data.literals = append(c.data.literals, lit.Symbol.Symbol())
		case lit.Var != nil:
			vr, ok := env.Resolve(lit.Var.Name.Symbol())
			if !ok {
				return nil, fmt.Errorf("missing var: %s", lit.Var.Name.Symbol())
			}

			c.data.literals = append(c.data.literals, vr)
		case lit.Type == nil:
			sym := lit.Type.Name.Symbol()
			c.data.literals = append(c.data.literals, TYPES[*&sym.name])
		case lit.Encoded != nil:
			obj, err := readFromReader(env, bytes.NewReader(lit.Encoded.Data))
			if err != nil {
				return nil, err
			}
			c.data.literals = append(c.data.literals, obj)
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
		c.data.methods = append(c.data.methods, MethodSite{
			Method: meth.Method,
			Arity:  meth.Arity,
		})
	}

	return c, nil
}
