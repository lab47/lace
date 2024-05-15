package core

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/candid82/liner"
)

type (
	ReplContext struct {
		first  *Var
		second *Var
		third  *Var
		exc    *Var
	}
)

func NewReplContext(env *Env) *ReplContext {
	first, _ := env.Resolve(MakeSymbol("lace.core/*1"))
	second, _ := env.Resolve(MakeSymbol("lace.core/*2"))
	third, _ := env.Resolve(MakeSymbol("lace.core/*3"))
	exc, _ := env.Resolve(MakeSymbol("lace.core/*e"))
	first.SetStatic(NIL)
	second.SetStatic(NIL)
	third.SetStatic(NIL)
	exc.SetStatic(NIL)
	return &ReplContext{
		first:  first,
		second: second,
		third:  third,
		exc:    exc,
	}
}

func (ctx *ReplContext) PushValue(obj Object) {
	ctx.third.SetStatic(ctx.second.GetStatic())
	ctx.second.SetStatic(ctx.first.GetStatic())
	ctx.first.SetStatic(obj)
}

func (ctx *ReplContext) PushException(exc Object) {
	ctx.exc.SetStatic(exc)
}

func (env *Env) REPL(in io.Reader, out io.Writer) error {
	env.FindNamespace(MakeSymbol("user")).ReferAll(env.FindNamespace(MakeSymbol("lace.repl")), true)
	fmt.Printf("Welcome to lace %s. Use '(exit)', %s to exit.\n", VERSION, "Contrl-D")
	parseContext := &ParseContext{Env: env}
	replContext := NewReplContext(env)

	var runeReader io.RuneReader

	rl := liner.NewLiner()
	defer rl.Close()
	rl.SetCtrlCAborts(true)
	rl.SetWordCompleter(makeCompleter(env))
	rl.SetTabCompletionStyle(liner.TabPrints)

	runeReader = NewLineRuneReader(rl)

	reader := NewReader(runeReader, "<repl>")

	for {
		namespace := env.CurrentNamespace().Name.String()
		runeReader.(*LineRuneReader).Prompt = (namespace + "=> ")
		done, err := processReplCommand(env, reader, parseContext, replContext)

		if err != nil {
			return err
		}

		if done {
			return nil
		}
	}

}

func skipRestOfLine(reader *Reader) error {
	for {
		c, err := reader.Get()
		if err != nil {
			return err
		}
		switch c {
		case EOF, '\n':
			return nil
		}
	}
}

func processReplCommand(env *Env, reader *Reader, parseContext *ParseContext, replContext *ReplContext) (bool, error) {
	obj, err := TryRead(env, reader)
	if err == io.EOF {
		return true, nil
	}
	if err != nil {
		fmt.Fprintln(Stderr, err)
		err = skipRestOfLine(reader)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return false, nil
		}
		return false, nil
	}

	expr, err := Parse(obj, parseContext)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return false, nil
	}

	res, err := Eval(env, expr, nil)
	if err != nil {
		switch r := err.(type) {
		case *ExitError:
			return true, err
		case *ParseError:
			replContext.PushException(r)
			fmt.Fprintln(Stderr, r)
		case *EvalError:
			replContext.PushException(r)
			DisplayError(env, r)
		case Error:
			replContext.PushException(r)
			fmt.Fprintln(Stderr, r)
		default:
			fmt.Printf("error: %s\n", err)
		}
		return false, nil
	}

	replContext.PushValue(res)
	PrintObject(env, res, Stdout)
	fmt.Fprintln(Stdout, "")
	return false, nil
}

var qualifiedSymbolRe *regexp.Regexp = regexp.MustCompile(`([0-9A-Za-z_\-\+\*\'\.]+)/([0-9A-Za-z_\-\+\*\']*$)`)
var callRe *regexp.Regexp = regexp.MustCompile(`\(\s*([0-9A-Za-z_\-\+\*\'\.]*$)`)

func makeCompleter(env *Env) func(line string, pos int) (head string, c []string, tail string) {
	return func(line string, pos int) (head string, c []string, tail string) {
		head = line[:pos]
		tail = line[pos:]
		var match []string
		var prefix string
		var ns *Namespace
		var addNamespaces bool
		if match = qualifiedSymbolRe.FindStringSubmatch(head); match != nil {
			nsName := match[1]
			prefix = match[2]
			ns = env.NamespaceFor(env.CurrentNamespace(), MakeSymbol(nsName+"/"+prefix))
		} else if match = callRe.FindStringSubmatch(head); match != nil {
			prefix = match[1]
			ns = env.CurrentNamespace()
			addNamespaces = true
		}
		if ns == nil {
			return
		}

		for _, k := range ns.VarNames() {
			if strings.HasPrefix(k, prefix) {
				c = append(c, k)
			}
		}
		if addNamespaces {
			for _, k := range env.AllNamespaces() {
				if strings.HasPrefix(k, prefix) {
					c = append(c, k)
				}
			}
			for _, k := range ns.AliasNames() {
				if strings.HasPrefix(k, prefix) {
					c = append(c, k)
				}
			}
		}
		if len(c) > 0 {
			head = head[:len(head)-len(prefix)]
		}
		sort.Strings(c)
		return
	}
}
