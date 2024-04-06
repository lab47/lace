package main

import (
	"bufio"
	"fmt"
	"io"

	. "github.com/lab47/lace/core"
)

func repl(phase Phase) {
	ProcessReplData()
	GLOBAL_ENV.FindNamespace(MakeSymbol("user")).ReferAll(GLOBAL_ENV.FindNamespace(MakeSymbol("lace.repl")))
	fmt.Printf("Welcome to lace %s. Use EOF (Ctrl-D) or SIGINT (Ctrl-C) to exit.\n", VERSION)
	parseContext := &ParseContext{GlobalEnv: GLOBAL_ENV}
	replContext := NewReplContext(parseContext.GlobalEnv)

	var runeReader io.RuneReader
	runeReader = bufio.NewReader(Stdin)
	reader := NewReader(runeReader, "<repl>")

	for {
		print(GLOBAL_ENV.CurrentNamespace().Name.ToString(false) + "=> ")
		if processReplCommand(reader, phase, parseContext, replContext) {
			return
		}
	}
}
