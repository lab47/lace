package pkgreflect

import (
	"bytes"
	"cmp"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/token"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

func codeName(t types.Type, curpkg string) string {
	switch sv := t.(type) {
	case *types.Pointer:
		return "*" + codeName(sv.Elem(), curpkg)
	case *types.Named:
		n := sv.Obj()
		if n.Pkg() == nil || n.Pkg().Name() == curpkg {
			return n.Name()
		}

		return n.Pkg().Name() + "." + n.Name()
	default:
		return t.String()
	}
}

func doesMatch(ident string, doc *ast.CommentGroup, match *Match) bool {
	if match.Directive {
		if doc == nil {
			return false
		}

		for _, c := range doc.List {
			if strings.Contains(c.Text, "//lace:export") {
				return true
			}
		}
	}

	if len(match.Patterns) == 0 {
		return true
	}

	for _, m := range match.Patterns {
		if ident == m {
			return true
		}
	}

	return false
}

type Match struct {
	Patterns  []string
	Directive bool
}

type GenOptions struct {
	Specialized bool
	InCore      bool
}

func typeSpecs(f *ast.File) []*ast.TypeSpec {
	var ret []*ast.TypeSpec

	for _, d := range f.Decls {
		if gd, ok := d.(*ast.GenDecl); ok {
			if gd.Tok == token.TYPE {
				for _, spec := range gd.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						if len(gd.Specs) == 1 {
							ts.Doc = gd.Doc
						}

						ret = append(ret, ts)
					}
				}
			}
		}
	}

	slices.SortFunc(ret, func(a, b *ast.TypeSpec) int {
		return cmp.Compare(a.Name.Name, b.Name.Name)
	})

	return ret
}

func Generate(name, laceName string, base, output, outputPkg string, match *Match, opts GenOptions) error {
	var pcfg packages.Config
	pcfg.Mode = packages.NeedSyntax | packages.NeedFiles | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedTypesSizes
	pkgs, err := packages.Load(&pcfg, name)
	if err != nil {
		panic(err)
	}

	pkg, err := build.Import(name, base, build.ImportComment)
	if err != nil {
		return err
	}

	/*
		var files []*ast.File

		ts := token.NewFileSet()
		sort.Strings(pkg.GoFiles)

		for _, path := range pkg.GoFiles {
			f, err := parser.ParseFile(ts, filepath.Join(pkg.Dir, path), nil, parser.ParseComments)
			if err != nil {
				panic(err)
			}

			files = append(files, f)
		}

		info := types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Defs:  make(map[*ast.Ident]types.Object),
			Uses:  make(map[*ast.Ident]types.Object),
		}

		var conf types.Config
		conf.Importer = importer.ForCompiler(ts, runtime.Compiler, func(path string) (io.ReadCloser, error) {
			spew.Dump(path)
			panic("no")
		})
		tpkg, err := conf.Check("p", ts, files, &info)
		if err != nil {
			log.Fatal(err)
		}
	*/

	tpkg := pkgs[0].Types
	files := pkgs[0].Syntax
	info := pkgs[0].TypesInfo

	var (
		same    bool
		pkgName string
	)

	if outputPkg == "" {
		same = true
		outputPkg = tpkg.Name()
	} else {
		pkgName = tpkg.Name() + "."
	}

	var buf bytes.Buffer
	fmt.Fprintln(&buf, "// Code generated by github.com/lab47/lace/pkg/pkgreflect DO NOT EDIT.")
	fmt.Fprintln(&buf, "package", outputPkg)
	fmt.Fprintln(&buf, "")
	fmt.Fprintln(&buf, `import "reflect"`)
	fmt.Fprintln(&buf, `import "github.com/lab47/lace/pkg/pkgreflect"`)
	if !same {
		fmt.Fprintf(&buf, "import %s \"%s\"\n", tpkg.Name(), pkg.ImportPath)
	}
	fmt.Fprintln(&buf, "")

	var (
		typesForMethods []string
		synthStruct     []string
	)

	for _, f := range files {
		tss := typeSpecs(f)

		for _, ts := range tss {
			if !ast.IsExported(ts.Name.Name) {
				continue
			}

			if !doesMatch(ts.Name.Name, ts.Doc, match) {
				continue
			}

			typesForMethods = append(typesForMethods, ts.Name.Name)

			if it, ok := ts.Type.(*ast.InterfaceType); ok {
				itt := info.Types[it].Type.(*types.Interface)

				fmt.Fprintf(&buf, "type %sImpl struct {\n", ts.Name.Name)
				synthStruct = append(synthStruct, ts.Name.Name+"Impl")

				for i := 0; i < itt.NumMethods(); i++ {
					fn := itt.Method(i)
					sig := fn.Type().(*types.Signature)
					var resStr string
					if sig.Results().Len() != 0 {
						var ary []string

						for i := 0; i < sig.Results().Len(); i++ {
							vr := sig.Results().At(i)
							ary = append(ary, codeName(vr.Type(), outputPkg))
						}
						resStr = "(" + strings.Join(ary, ", ") + ")"
					}

					params := sig.Params()

					var pStrs []string

					for i := 0; i < params.Len(); i++ {
						vr := params.At(i)
						t := vr.Type()

						pStrs = append(pStrs, codeName(t, outputPkg))
					}
					fmt.Fprintf(&buf, "  %sFn func(%s) %s\n", fn.Name(), strings.Join(pStrs, ", "), resStr)
				}
				fmt.Fprintf(&buf, "}\n")

				for i := 0; i < itt.NumMethods(); i++ {
					fn := itt.Method(i)
					sig := fn.Type().(*types.Signature)

					var args []string
					var cs []string

					for j := 0; j < sig.Params().Len(); j++ {
						e := sig.Params().At(j)
						args = append(args, fmt.Sprintf("a%d %s", j, codeName(e.Type(), outputPkg)))
						cs = append(cs, fmt.Sprintf("a%d", j))
					}

					var resStr string
					if sig.Results().Len() != 0 {
						var ary []string

						for i := 0; i < sig.Results().Len(); i++ {
							vr := sig.Results().At(i)
							ary = append(ary, codeName(vr.Type(), outputPkg))
						}
						resStr = "(" + strings.Join(ary, ", ") + ")"
					}

					fmt.Fprintf(&buf, "func (s *%sImpl) %s(%s) %s {\n", ts.Name.Name, fn.Name(), strings.Join(args, ", "), resStr)
					if sig.Results().Len() == 0 {
						fmt.Fprintf(&buf, "s.%sFn(%s)\n", fn.Name(), strings.Join(cs, ", "))
					} else {
						fmt.Fprintf(&buf, "return s.%sFn(%s)\n", fn.Name(), strings.Join(cs, ", "))
					}
					fmt.Fprintln(&buf, "}")
				}
			}
		}
	}

	fmt.Fprintln(&buf, "")
	fmt.Fprintln(&buf, "func init() {")

	for _, t := range typesForMethods {
		fmt.Fprintf(&buf, "%s_methods := map[string]pkgreflect.Func{}\n", t)
	}

	for _, f := range files {
		for _, d := range f.Decls {
			if fn, ok := d.(*ast.FuncDecl); ok {
				if !ast.IsExported(fn.Name.Name) {
					continue
				}
				if !doesMatch(fn.Name.Name, fn.Doc, match) {
					continue
				}

				fn.Doc.Text()

				if fn.Recv != nil && len(fn.Recv.List) == 1 {
					rt := typeName(fn.Recv.List[0].Type)
					if !ast.IsExported(rt) {
						continue
					}

					var args []string
					for _, f := range fn.Type.Params.List {
						for _, n := range f.Names {
							args = append(args, fmt.Sprintf("{Name: \"%s\", Tag: \"%s\"}", n.Name, typeName(f.Type)))
						}
					}

					tag := "any"

					if fn.Type.Results != nil && len(fn.Type.Results.List) == 1 {
						tag = typeName(fn.Type.Results.List[0].Type)
					}

					arity := strings.Join(args, `,`)
					fmt.Fprintf(&buf, "%s_methods[%q] = pkgreflect.Func{Args: []pkgreflect.Arg{%s}, Tag: \"%s\", Doc: %q}\n", rt, fn.Name.Name, arity, tag, strings.TrimSpace(fn.Doc.Text()))
				}
			}
		}
	}

	fmt.Fprintf(&buf, "pkgreflect.AddPackage(%q, &pkgreflect.Package{\n", laceName)
	fmt.Fprintf(&buf, "Doc: %q,\n", pkg.Doc)

	// Types
	fmt.Fprintln(&buf, "Types: map[string]pkgreflect.Type{")
	tprint(&buf, pkgName, files, "\t%q: {Doc: %q, Value: reflect.TypeOf((*%s%s)(nil)).Elem(), Methods: %[1]s_methods},\n", match)
	for _, ss := range synthStruct {
		fmt.Fprintf(&buf, "\t%q: {Doc: `Struct version of interface %s for implementation`, Value: reflect.TypeFor[%[1]s]()},\n", ss, ss[:len(ss)-4])
	}
	fmt.Fprintln(&buf, "},")
	fmt.Fprintln(&buf, "")

	// Functions
	fmt.Fprintln(&buf, "Functions: map[string]pkgreflect.FuncValue{")
	fnprint(&buf, pkgName, files, ast.Fun, match, opts)
	fmt.Fprintln(&buf, "},")
	fmt.Fprintln(&buf, "")

	// Addresses of variables
	fmt.Fprintln(&buf, "Variables: map[string]pkgreflect.Value{")
	print(&buf, pkgName, files, ast.Var, "\t%q: {Doc: %q, Value: reflect.ValueOf(&%s%s)},\n", match)
	fmt.Fprintln(&buf, "},")
	fmt.Fprintln(&buf, "")

	// Addresses of consts
	fmt.Fprintln(&buf, "Consts: map[string]pkgreflect.Value{")
	print(&buf, pkgName, files, ast.Con, "\t%q: {Doc: %q, Value: reflect.ValueOf(%s%s)},\n", match)
	fmt.Fprintln(&buf, "},")
	fmt.Fprintln(&buf, "")

	fmt.Fprintln(&buf, "})")
	fmt.Fprintln(&buf, "}")

	data, err := imports.Process("out.go", buf.Bytes(), &imports.Options{})
	if err != nil {
		os.Stderr.Write(buf.Bytes())
		return err
	}

	formatted, err := format.Source(data)
	if err != nil {
		os.Stderr.Write(buf.Bytes())
		return err
	}

	if output == "-" {
		os.Stdout.Write(formatted)
	} else {
		filename := filepath.Join(output)
		oldFileData, _ := os.ReadFile(filename)
		if !bytes.Equal(formatted, oldFileData) {
			err = os.WriteFile(filename, formatted, 0660)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
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

func getDoc(object *ast.Object) *ast.CommentGroup {
	switch v := object.Decl.(type) {
	case *ast.ValueSpec:
		return v.Doc
	case *ast.FuncDecl:
		return v.Doc
	case *ast.TypeSpec:
		return v.Doc
	default:
		return nil
	}
}

func tprint(w io.Writer, pkgName string, files []*ast.File, format string, match *Match) {
	type ent struct {
		name, doc string
	}

	names := []ent{}
	for _, f := range files {
		tss := typeSpecs(f)

		for _, ts := range tss {
			doc := ts.Doc
			name := ts.Name.Name

			if !ast.IsExported(name) {
				continue
			}

			if !doesMatch(name, doc, match) {
				continue
			}

			docStr := ""
			if doc != nil {
				docStr = strings.TrimSpace(doc.Text())
			}

			names = append(names, ent{name, docStr})
		}
	}

	for _, name := range names {
		fmt.Fprintf(w, format, name.name, name.doc, pkgName, name.name)
	}
}

func print(w io.Writer, pkgName string, files []*ast.File, kind ast.ObjKind, format string, match *Match) {
	type ent struct {
		name, doc string
	}

	names := []ent{}
	for _, f := range files {
		for name, object := range f.Scope.Objects {
			if object.Kind != kind || !ast.IsExported(name) {
				continue
			}

			doc := getDoc(object)

			if !doesMatch(name, doc, match) {
				continue
			}

			docStr := ""
			if doc != nil {
				docStr = strings.TrimSpace(doc.Text())
			}

			names = append(names, ent{name, docStr})
		}
	}

	sort.Slice(names, func(i, j int) bool {
		return names[i].name < names[j].name
	})

	for _, name := range names {
		fmt.Fprintf(w, format, name.name, name.doc, pkgName, name.name)
	}
}

func exportName(fn *ast.FuncDecl) string {
	if fn.Doc != nil {
		for _, c := range fn.Doc.List {
			if strings.HasPrefix(c.Text, "//lace:export ") {
				name := c.Text[len("//lace:export "):]
				name = strings.TrimSpace(name)
				if name != "" {
					return name
				}
			}
		}
	}

	return fn.Name.Name
}

func fnprint(w io.Writer, pkgName string, files []*ast.File, kind ast.ObjKind, match *Match, opts GenOptions) {
	var fns []string

	for _, f := range files {
		for name, object := range f.Scope.Objects {
			if object.Kind == kind && ast.IsExported(name) {
				fn := object.Decl.(*ast.FuncDecl)

				if !doesMatch(name, fn.Doc, match) {
					continue
				}
				var args []string
				var params int
				for _, f := range fn.Type.Params.List {
					for _, n := range f.Names {
						params++
						tn := typeName(f.Type)
						if tn == "core.Env" || tn == "Env" {
							continue
						}

						tn = strings.TrimPrefix(tn, "core.")

						args = append(args, fmt.Sprintf("{Name: \"%s\", Tag: \"%s\"}", n.Name, tn))
					}
				}

				rt := "any"

				if fn.Type.Results != nil && len(fn.Type.Results.List) == 1 {
					rt = typeName(fn.Type.Results.List[0].Type)
				}

				var rets int
				if fn.Type.Results != nil {
					for _, f := range fn.Type.Results.List {
						if f.Names != nil {
							rets += len(f.Names)
						} else {
							rets++
						}
					}
				}

				arity := strings.Join(args, `,`)
				if opts.Specialized && params <= 3 && rets <= 2 {
					corePkg := "core."
					if opts.InCore {
						corePkg = ""
					}

					fns = append(fns,
						fmt.Sprintf("\t\"%s\": {Doc: %q, Args: []pkgreflect.Arg{%s}, Tag: \"%s\", Value: reflect.ValueOf(%sWrapToProc%d_%d(%s%s))},\n",
							exportName(fn), strings.TrimSpace(fn.Doc.Text()), arity, rt,
							corePkg,
							params,
							len(fn.Type.Results.List),
							pkgName, fn.Name.Name),
					)
				} else {
					fns = append(fns,
						fmt.Sprintf("\t\"%s\": {Doc: %q, Args: []pkgreflect.Arg{%s}, Tag: \"%s\", Value: reflect.ValueOf(%s%s)},\n",
							exportName(fn), strings.TrimSpace(fn.Doc.Text()), arity, rt,
							pkgName, fn.Name.Name),
					)
				}
			}
		}
	}

	sort.Strings(fns)

	for _, f := range fns {
		fmt.Fprintln(w, f)
	}
}
