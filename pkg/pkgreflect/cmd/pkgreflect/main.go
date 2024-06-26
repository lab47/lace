package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lab47/lace/pkg/pkgreflect"
)

var laceName = flag.String("lace-name", "", "default to same as name")
var pkgName = flag.String("pkg-name", "", "name of the package to generate")
var directive = flag.Bool("honor-directive", false, "Export all entities with //lace:export directive in doc")
var match = flag.String("match", "", "if set, only match elements that have these names (comma delim list)")
var inCore = flag.Bool("in-core", false, "output code to be linked directly into core (ie lace.lang)")
var specialized = flag.Bool("specialized", false, "output code that attepmts to specialize functions")

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

	match := &pkgreflect.Match{
		Patterns:  strings.Split(*match, ","),
		Directive: *directive,
	}

	err = pkgreflect.Generate(name, ln, wd, output, *pkgName, match, pkgreflect.GenOptions{
		Specialized: *specialized,
		InCore:      *inCore,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
}
