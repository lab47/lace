package main

import (
	"flag"
	"fmt"
	"go/ast"
	"os"

	"github.com/lab47/lace/pkg/pkgreflect"
)

var (
	stdout bool
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

func typeName(x ast.Expr) string {
	for {
		switch sv := x.(type) {
		case *ast.StarExpr:
			x = sv.X
		case *ast.Ident:
			return sv.Name
		case *ast.SelectorExpr:
			return typeName(sv.X) + "." + sv.Sel.Name
		case *ast.ArrayType:
			return "[]" + typeName(sv.Elt)
		case *ast.InterfaceType:
			if len(sv.Methods.List) == 0 {
				return "any"
			}

			return "interface"
		default:
			return "Unknown"
		}
	}
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
