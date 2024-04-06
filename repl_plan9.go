package main

import (
	"bufio"
	"fmt"
	"io"

	. "github.com/lab47/lace/core"
)

func repl(env *Env, phase Phase) {
	ProcessReplData()
	env.FindNamespace(MakeSymbol("user")).ReferAll(env.FindNamespace(MakeSymbol("lace.repl")))
	fmt.Printf("Welcome to lace %s. Use EOF (Ctrl-D) or SIGINT (Ctrl-C) to exit.\n", VERSION)
	parseContext := &ParseContext{Env: env}
	replContext := NewReplContext(parseContext.Env)

	var runeReader io.RuneReader
	runeReader = bufio.NewReader(Stdin)
	reader := NewReader(runeReader, "<repl>")

	for {
		print(env.CurrentNamespace().Name.ToString(false) + "=> ")
		if processReplCommand(env, reader, phase, parseContext, replContext) {
			return
		}
	}
}
