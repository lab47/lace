package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lab47/lace/pkg/pkgreflect"
)

var laceName = flag.String("lace-name", "", "default to same as name")

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

	ln := *laceName
	if ln == "" {
		ln = name
	}

	err = pkgreflect.Generate(name, ln, wd, output)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
}
