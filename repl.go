//go:build !plan9
// +build !plan9

package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/candid82/joker/core"
	"github.com/chzyer/readline"
)

func repl(env *core.Env, phase core.Phase) {
	core.ProcessReplData()
	env.FindNamespace(core.MakeSymbol("user")).ReferAll(env.FindNamespace(core.MakeSymbol("joker.repl")))
	fmt.Printf("Welcome to joker %s. Use EOF (Ctrl-D) or SIGINT (Ctrl-C) to exit.\n", core.VERSION)
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
