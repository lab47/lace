package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lab47/lace/pkg/pkgreflect"
)

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("pkgreflect <package name> <output file>")
		os.Exit(1)
	}

	pkg := flag.Arg(0)
	output := "-"

	if flag.NArg() == 2 {
		output = flag.Arg(1)
	}

	parseDir(pkg, output)
}

func parseDir(name, output string) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	err = pkgreflect.Generate(name, wd, output)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
}
