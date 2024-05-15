package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/lab47/lace/core"
	_ "github.com/lab47/lace/std-ng/all"
	_ "github.com/lab47/lace/std/csv"
	_ "github.com/lab47/lace/std/filepath"
	_ "github.com/lab47/lace/std/io"
	_ "github.com/lab47/lace/std/math"
	_ "github.com/lab47/lace/std/os"
	_ "github.com/lab47/lace/std/strconv"
	_ "github.com/lab47/lace/std/time"
	_ "github.com/lab47/lace/std/url"
	_ "github.com/lab47/lace/std/uuid"
	"github.com/spf13/pflag"
)

type (
	ReplContext struct {
		first  *core.Var
		second *core.Var
		third  *core.Var
		exc    *core.Var
	}
)

func NewReplContext(env *core.Env) *ReplContext {
	first, _ := env.Resolve(core.MakeSymbol("lace.core/*1"))
	second, _ := env.Resolve(core.MakeSymbol("lace.core/*2"))
	third, _ := env.Resolve(core.MakeSymbol("lace.core/*3"))
	exc, _ := env.Resolve(core.MakeSymbol("lace.core/*e"))
	first.SetStatic(core.NIL)
	second.SetStatic(core.NIL)
	third.SetStatic(core.NIL)
	exc.SetStatic(core.NIL)
	return &ReplContext{
		first:  first,
		second: second,
		third:  third,
		exc:    exc,
	}
}

func (ctx *ReplContext) PushValue(obj core.Object) {
	ctx.third.SetStatic(ctx.second.GetStatic())
	ctx.second.SetStatic(ctx.first.GetStatic())
	ctx.first.SetStatic(obj)
}

func (ctx *ReplContext) PushException(exc core.Object) {
	ctx.exc.SetStatic(exc)
}

func processFile(env *core.Env, filename string) error {
	var reader *core.Reader
	if filename == "-" {
		reader = core.NewReader(bufio.NewReader(core.Stdin), "<stdin>")
		filename = ""
	} else {
		var err error
		reader, err = core.NewReaderFromFile(filename)
		if err != nil {
			return err
		}
	}
	if filename != "" {
		f, err := filepath.Abs(filename)
		if err != nil {
			return err
		}
		env.SetMainFilename(f)
	}
	_, err := core.ProcessReader(env, reader, filename)
	return err
}

func skipRestOfLine(reader *core.Reader) error {
	for {
		c, err := reader.Get()
		if err != nil {
			return err
		}
		switch c {
		case core.EOF, '\n':
			return nil
		}
	}
}

func processReplCommand(env *core.Env, reader *core.Reader, parseContext *core.ParseContext, replContext *ReplContext) (bool, error) {

	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case *core.ParseError:
				replContext.PushException(r)
				fmt.Fprintln(core.Stderr, r)
			case *core.EvalError:
				replContext.PushException(r)
				fmt.Fprintln(core.Stderr, r)
			case core.Error:
				replContext.PushException(r)
				fmt.Fprintln(core.Stderr, r)
				// case *runtime.TypeAssertionError:
				// 	fmt.Fprintln(Stderr, r)
			default:
				panic(r)
			}
		}
	}()

	obj, err := core.TryRead(env, reader)
	if err == io.EOF {
		return true, nil
	}
	if err != nil {
		fmt.Fprintln(core.Stderr, err)
		err = skipRestOfLine(reader)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return false, nil
		}
		return false, nil
	}

	expr, err := core.Parse(obj, parseContext)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return false, nil
	}

	/*
		fn, err := core.Compile(env, []core.Expr{expr})
		if err != nil {
			fmt.Printf("error compiling: %s\n", err)
		} else {
			obj, err := core.EngineRun(env, fn)
			if err != nil {
				fmt.Printf("error running bytecode: %s\n", err)
			} else {
				spew.Dump(obj)
			}
		}
	*/

	res, err := core.Eval(env, expr, nil)
	if err != nil {
		if _, ok := err.(*core.ExitError); ok {
			return true, err
		}
		fmt.Printf("error: %s\n", err)
		return false, nil
	}
	replContext.PushValue(res)
	core.PrintObject(env, res, core.Stdout)
	fmt.Fprintln(core.Stdout, "")
	return false, nil
}

func makeDialectKeyword(dialect core.Dialect) core.Keyword {
	switch dialect {
	case core.EDN:
		return core.MakeKeyword("clj")
	case core.CLJ:
		return core.MakeKeyword("clj")
	case core.CLJS:
		return core.MakeKeyword("cljs")
	default:
		return core.MakeKeyword("lace")
	}
}

func configureLinterMode(env *core.Env, dialect core.Dialect, filename string, workingDir string) error {
	if err := core.ProcessLinterFiles(env, dialect, filename, workingDir); err != nil {
		return err
	}

	core.LINTER_MODE = true
	core.DIALECT = dialect
	lm, _ := env.Resolve(core.MakeSymbol("lace.core/*linter-mode*"))
	lm.SetStatic(core.Boolean(true))
	mf, err := env.Features.Disjoin(env, core.MakeKeyword("lace"))
	if err != nil {
		return err
	}
	f, err := mf.Conj(env, makeDialectKeyword(dialect))
	if err != nil {
		return err
	}
	env.Features = f.(core.Set)
	return core.ProcessLinterData(env, dialect)
}

func detectDialect(filename string) core.Dialect {
	switch {
	case strings.HasSuffix(filename, ".edn"):
		return core.EDN
	case strings.HasSuffix(filename, ".cljs"):
		return core.CLJS
	case strings.HasSuffix(filename, ".clj"):
		return core.LACE
	}
	return core.CLJ
}

func matchesDialect(path string, dialect core.Dialect) bool {
	ext := ".clj"
	switch dialect {
	case core.CLJS:
		ext = ".cljs"
	case core.LACE:
		ext = ".clj"
	case core.EDN:
		ext = ".edn"
	}
	return strings.HasSuffix(path, ext)
}

func isIgnored(path string) bool {
	for _, r := range core.WARNINGS.IgnoredFileRegexes {
		m := r.FindStringSubmatchIndex(path)
		if len(m) > 0 {
			if m[1]-m[0] == len(path) {
				return true
			}
		}
	}
	return false
}

func dialectFromArg(arg string) core.Dialect {
	switch strings.ToLower(arg) {
	case "clj":
		return core.CLJ
	case "cljs":
		return core.CLJS
	case "lace":
		return core.LACE
	case "edn":
		return core.EDN
	}
	return core.UNKNOWN
}

func isNumber(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func notOption(arg string) bool {
	return arg == "-" || !strings.HasPrefix(arg, "-") || isNumber(arg[1:])
}

var runningProfile interface {
	Stop()
}

func MainIn(nsName string) {
	env, err := core.NewEnv()
	if err != nil {
		fmt.Printf("unable to initialize environment: %s", err)
		os.Exit(1)
	}

	env.InitEnv(core.Stdin, core.Stdout, core.Stderr, os.Args[1:])

	fs := pflag.NewFlagSet("lace", pflag.ExitOnError)
	version := fs.BoolP("version", "v", false, "report the version number")
	cpuProfile := fs.String("cpuprofile", "", "Write CPU profile info to the specified path")
	cpuProfileRate := fs.Int("cpuprofile-rate", 100, "Specify the sampling rate of the cpu profiler")
	memProfile := fs.String("memprofile", "", "Write Memory profile info to the specified path")
	debugBytecode := fs.Bool("debug-bytecode", false, "Display bytecode for functions are it is generated")

	if err := fs.Parse(os.Args); err != nil {
		fmt.Printf("error parsing arguments: %s\n", err)
		os.Exit(1)
	}

	env.SetEnvArgs(fs.Args()[1:])

	env.SetClassPath(".")

	if *version {
		println(core.VERSION)
		return
	}

	/* Set up profiling. */

	cpuProfileName := *cpuProfile
	memProfileName := *memProfile

	var teardown []func()

	core.SetExit(func(code int) {
		for _, x := range teardown {
			x()
		}

		finish(memProfileName)
		os.Exit(code)
	})

	env.DebugBytecode = *debugBytecode

	if cpuProfileName != "" {
		f, err := os.Create(cpuProfileName)
		if err != nil {
			fmt.Fprintf(core.Stderr, "Error: Could not create CPU profile `%s': %v\n",
				cpuProfileName, err)
			cpuProfileName = ""
			core.Exit(96)
		}
		defer f.Close()
		err = pprof.StartCPUProfile(f)
		runtime.SetCPUProfileRate(*cpuProfileRate)
		if err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
		teardown = append(teardown, pprof.StopCPUProfile)
		fmt.Fprintf(core.Stderr, "Profiling started at rate=%d. See file `%s'.\n",
			*cpuProfileRate, cpuProfileName)
	} else if memProfileName != "" {
		defer finish(memProfileName)
	}

	_, err = core.Load(env, nsName)
	if err != nil {
		core.DisplayError(env, err)
		os.Exit(1)
	}

	ns := env.FindNamespace(core.MakeSymbol(nsName))
	if ns == nil {
		fmt.Fprintf(core.Stderr, "Unable to find namespace to executed main: %s", nsName)
		os.Exit(1)
	}

	vr, err := ns.Intern(env, core.MakeSymbol("main"))
	if err != nil {
		fmt.Fprintf(core.Stderr, "Unable to find %s/main", nsName)
		os.Exit(1)
	}

	if !vr.Set() {
		fmt.Fprintf(core.Stderr, "%s/main is nil\n", nsName)
		os.Exit(1)
	}

	cl, ok := vr.GetStatic().(core.Callable)
	if !ok {
		fmt.Fprintf(core.Stderr, "%s/main is not callable", nsName)
		os.Exit(1)
	}

	_, err = cl.Call(env, []core.Object{})
	if err != nil {
		core.DisplayError(env, err)
		os.Exit(1)
	}

	fmt.Println("here")
}

func Main() {
	env, err := core.NewEnv()
	if err != nil {
		fmt.Printf("unable to initialize environment: %s", err)
		os.Exit(1)
	}

	env.InitEnv(core.Stdin, core.Stdout, core.Stderr, os.Args[1:])

	fs := pflag.NewFlagSet("lace", pflag.ExitOnError)
	version := fs.BoolP("version", "v", false, "report the version number")
	cpuProfile := fs.String("cpuprofile", "", "Write CPU profile info to the specified path")
	cpuProfileRate := fs.Int("cpuprofile-rate", 100, "Specify the sampling rate of the cpu profiler")
	memProfile := fs.String("memprofile", "", "Write Memory profile info to the specified path")
	debugBytecode := fs.Bool("debug-bytecode", false, "Display bytecode for functions are it is generated")

	if err := fs.Parse(os.Args); err != nil {
		fmt.Printf("error parsing arguments: %s\n", err)
		os.Exit(1)
	}

	var filename string

	if fs.NArg() >= 2 {
		filename = fs.Arg(1)
		env.SetEnvArgs(fs.Args()[2:])
	}

	//env.SetClassPath(classPath)

	if *version {
		println(core.VERSION)
		return
	}

	/* Set up profiling. */

	cpuProfileName := *cpuProfile
	memProfileName := *memProfile

	var teardown []func()

	core.SetExit(func(code int) {
		for _, x := range teardown {
			x()
		}

		finish(memProfileName)
		os.Exit(code)
	})

	env.DebugBytecode = *debugBytecode

	if cpuProfileName != "" {
		f, err := os.Create(cpuProfileName)
		if err != nil {
			fmt.Fprintf(core.Stderr, "Error: Could not create CPU profile `%s': %v\n",
				cpuProfileName, err)
			cpuProfileName = ""
			core.Exit(96)
		}
		defer f.Close()
		err = pprof.StartCPUProfile(f)
		runtime.SetCPUProfileRate(*cpuProfileRate)

		if err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
		teardown = append(teardown, pprof.StopCPUProfile)
		fmt.Fprintf(core.Stderr, "Profiling started at rate=%d. See file `%s'.\n",
			*cpuProfileRate, cpuProfileName)
		defer finish(memProfileName)
	} else if memProfileName != "" {
		defer finish(memProfileName)
	}

	if filename != "" {
		if err := processFile(env, filename); err != nil {
			if ee, ok := err.(*core.ExitError); ok {
				core.Exit(ee.Code)
			}

			core.Exit(1)
		} else {
			return
		}
	}

	env.REPL(os.Stdin, os.Stdout)
}

func finish(memProfileName string) {
	if runningProfile != nil {
		runningProfile.Stop()
		runningProfile = nil
	}

	if memProfileName != "" {
		f, err := os.Create(memProfileName)
		if err != nil {
			fmt.Fprintf(core.Stderr, "Error: Could not create memory profile `%s': %v\n",
				memProfileName, err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Fprintf(core.Stderr, "Error: Could not write memory profile `%s': %v\n",
				memProfileName, err)
		}
		f.Close()
		fmt.Fprintf(core.Stderr, "Memory profile rate=%d written to `%s'.\n",
			runtime.MemProfileRate, memProfileName)
		memProfileName = ""
	}
}
