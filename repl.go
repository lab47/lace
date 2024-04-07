//go:build !plan9
// +build !plan9

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/candid82/liner"
	"github.com/lab47/lace/core"
)

/*
func xrepl(env *core.Env, phase core.Phase) {
	core.ProcessReplData()
	env.FindNamespace(core.MakeSymbol("user")).ReferAll(env.FindNamespace(core.MakeSymbol("lace.repl")))
	fmt.Printf("Welcome to lace %s. Use EOF (Ctrl-D) or SIGINT (Ctrl-C) to exit.\n", core.VERSION)
	parseContext := &core.ParseContext{Env: env}
	replContext := NewReplContext(parseContext.Env)

	var runeReader io.RuneReader
	var rl *readline.Instance
	var err error
	if noReadline {
		runeReader = bufio.NewReader(core.Stdin)
	} else {
		rl, err = readline.New("")
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		defer rl.Close()
		runeReader = core.NewLineRuneReader(rl)
		for _, line := range strings.Split(string(dataRead), "\n") {
			rl.SaveHistory(line)
		}
		dataRead = []rune{}
	}

	reader := core.NewReader(runeReader, "<repl>")

	for {
		if noReadline {
			print(env.CurrentNamespace().Name.ToString(false) + "=> ")
		} else {
			rl.SetPrompt(env.CurrentNamespace().Name.ToString(false) + "=> ")
		}
		if processReplCommand(env, reader, phase, parseContext, replContext) {
			return
		}
	}
}
*/

var qualifiedSymbolRe *regexp.Regexp = regexp.MustCompile(`([0-9A-Za-z_\-\+\*\'\.]+)/([0-9A-Za-z_\-\+\*\']*$)`)
var callRe *regexp.Regexp = regexp.MustCompile(`\(\s*([0-9A-Za-z_\-\+\*\'\.]*$)`)

func makeCompleter(env *core.Env) func(line string, pos int) (head string, c []string, tail string) {
	return func(line string, pos int) (head string, c []string, tail string) {
		head = line[:pos]
		tail = line[pos:]
		var match []string
		var prefix string
		var ns *core.Namespace
		var addNamespaces bool
		if match = qualifiedSymbolRe.FindStringSubmatch(head); match != nil {
			nsName := match[1]
			prefix = match[2]
			ns = env.NamespaceFor(env.CurrentNamespace(), core.MakeSymbol(nsName+"/"+prefix))
		} else if match = callRe.FindStringSubmatch(head); match != nil {
			prefix = match[1]
			ns = env.CurrentNamespace()
			addNamespaces = true
		}
		if ns == nil {
			return
		}

		for k, _ := range ns.Mappings() {
			if strings.HasPrefix(*k, prefix) {
				c = append(c, *k)
			}
		}
		if addNamespaces {
			for k, _ := range env.Namespaces {
				if strings.HasPrefix(*k, prefix) {
					c = append(c, *k)
				}
			}
			for k, _ := range ns.Aliases() {
				if strings.HasPrefix(*k, prefix) {
					c = append(c, *k)
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

func saveReplHistory(rl *liner.State, filename string) {
	if filename == "" {
		return
	}
	if f, err := os.Create(filename); err == nil {
		rl.WriteHistory(f)
		f.Close()
	}
}

func repl(env *core.Env, phase core.Phase) error {
	core.ProcessReplData()
	env.FindNamespace(core.MakeSymbol("user")).ReferAll(env.FindNamespace(core.MakeSymbol("lace.repl")))
	fmt.Printf("Welcome to lace %s. Use '(exit)', %s to exit.\n", core.VERSION, "Contrl-D")
	parseContext := &core.ParseContext{Env: env}
	replContext := NewReplContext(env)

	var runeReader io.RuneReader
	var rl *liner.State
	var historyFilename string
	if noReadline {
		runeReader = bufio.NewReader(core.Stdin)
	} else {
		rl = liner.NewLiner()
		defer rl.Close()
		rl.SetCtrlCAborts(true)
		rl.SetWordCompleter(makeCompleter(env))
		rl.SetTabCompletionStyle(liner.TabPrints)

		runeReader = core.NewLineRuneReader(rl)

		for _, line := range strings.Split(string(dataRead), "\n") {
			if strings.TrimSpace(line) != "" {
				rl.AppendHistory(line)
			}
		}
		dataRead = []rune{}
	}

	reader := core.NewReader(runeReader, "<repl>")

	for {
		namespace := env.CurrentNamespace().Name.ToString(false)
		if noReadline {
			print(namespace + "=> ")
		} else {
			runeReader.(*core.LineRuneReader).Prompt = (namespace + "=> ")
		}
		done, err := processReplCommand(env, reader, phase, parseContext, replContext)

		if err != nil {
			return err
		}

		if done {
			saveReplHistory(rl, historyFilename)
			return nil
		}
	}
}
