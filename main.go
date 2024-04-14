package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"net"
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
	"github.com/pkg/profile"
)

var dataRead = []rune{}
var saveForRepl = true

type replayable struct {
	reader *core.Reader
}

func (r *replayable) ReadRune() (ch rune, size int, err error) {
	ch, err = r.reader.Get()
	if err != nil {
		return 0, 0, err
	}
	if ch == core.EOF {
		err = io.EOF
		size = 0
	} else {
		dataRead = append(dataRead, ch)
		size = 1
	}
	return
}

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
	first.Value = core.NIL
	second.Value = core.NIL
	third.Value = core.NIL
	exc.Value = core.NIL
	return &ReplContext{
		first:  first,
		second: second,
		third:  third,
		exc:    exc,
	}
}

func (ctx *ReplContext) PushValue(obj core.Object) {
	ctx.third.Value = ctx.second.Value
	ctx.second.Value = ctx.first.Value
	ctx.first.Value = obj
}

func (ctx *ReplContext) PushException(exc core.Object) {
	ctx.exc.Value = exc
}

func processFile(env *core.Env, filename string, phase core.Phase) error {
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
	if saveForRepl {
		reader = core.NewReader(&replayable{reader}, "<replay>")
	}
	return core.ProcessReader(env, reader, filename, phase)
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

func processReplCommand(env *core.Env, reader *core.Reader, phase core.Phase, parseContext *core.ParseContext, replContext *ReplContext) (bool, error) {

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

	if phase == core.READ {
		s, err := obj.ToString(env, true)
		if err != nil {
			return false, err
		}
		fmt.Println(s)
		return false, nil
	}

	expr, err := core.Parse(obj, parseContext)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return false, nil
	}
	if phase == core.PARSE {
		fmt.Println(expr)
		return false, nil
	}

	res, err := core.TopEval(env, expr, nil)
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

func srepl(env *core.Env, port string, phase core.Phase) error {
	core.ProcessReplData()
	env.FindNamespace(core.MakeSymbol("user")).ReferAll(env.FindNamespace(core.MakeSymbol("lace.repl")))
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Fprintf(core.Stderr, "Cannot start srepl listening on %s: %s\n",
			replSocket, err.Error())
		core.Exit(12)
	}
	defer l.Close()

	fmt.Printf("Joker repl listening at %s...\n", l.Addr())
	conn, err := l.Accept() // Wait for a single connection
	if err != nil {
		fmt.Fprintf(core.Stderr, "Cannot start repl accepting on %s: %s\n",
			l.Addr(), err.Error())
		core.Exit(13)
	}

	oldStdIn := core.Stdin
	oldStdOut := core.Stdout
	oldStdErr := core.Stderr
	oldStdinValue, oldStdoutValue, oldStderrValue := env.StdIO()
	core.Stdin = conn
	core.Stdout = conn
	core.Stderr = conn
	newIn := core.MakeBufferedReader(conn)
	newOut := core.MakeIOWriter(conn)
	env.SetStdIO(newIn, newOut, newOut)
	defer func() {
		conn.Close()
		core.Stdin = oldStdIn
		core.Stdout = oldStdOut
		core.Stderr = oldStdErr
		env.SetStdIO(oldStdinValue, oldStdoutValue, oldStderrValue)
	}()

	fmt.Printf("Joker repl accepting client at %s...\n", conn.RemoteAddr())

	runeReader := bufio.NewReader(conn)

	/* The rest of this code comes from repl(), below: */

	parseContext := &core.ParseContext{Env: env}
	replContext := NewReplContext(parseContext.Env)

	reader := core.NewReader(runeReader, "<srepl>")

	fmt.Fprintf(core.Stdout, "Welcome to lace %s, client at %s. Use '(lace.os/exit 0)', or close the connection, to exit.\n",
		core.VERSION, conn.RemoteAddr())

	for {
		fmt.Fprint(core.Stdout, env.CurrentNamespace().Name.String()+"=> ")
		done, err := processReplCommand(env, reader, phase, parseContext, replContext)
		if err != nil {
			return err
		}

		if done {
			return nil
		}
	}
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
	lm.Value = core.Boolean{B: true}
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

func lintFile(env *core.Env, filename string, dialect core.Dialect, workingDir string) error {
	phase := core.PARSE
	if dialect == core.EDN {
		phase = core.READ
	}
	err := core.ReadConfig(env, filename, workingDir)
	if err != nil {
		return err
	}
	err = configureLinterMode(env, dialect, filename, workingDir)
	if err != nil {
		return err
	}
	if processFile(env, filename, phase) == nil {
		core.WarnOnUnusedNamespaces(env)
		core.WarnOnUnusedVars(env)
	}

	return nil
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

func lintDir(env *core.Env, dirname string, dialect core.Dialect, reportGloballyUnused bool) error {
	var processErr error
	phase := core.PARSE
	if dialect == core.EDN {
		phase = core.READ
	}
	ns := env.CurrentNamespace()
	err := core.ReadConfig(env, "", dirname)
	if err != nil {
		return err
	}
	err = configureLinterMode(env, dialect, "", dirname)
	if err != nil {
		return err
	}

	err = filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintln(core.Stderr, "Error: ", err)
			return nil
		}
		if !info.IsDir() && matchesDialect(path, dialect) && !isIgnored(path) {
			env.CoreNamespace.Resolve("*loaded-libs*").Value = core.EmptySet()
			processErr = processFile(env, path, phase)
			if processErr == nil {
				core.WarnOnUnusedNamespaces(env)
				core.WarnOnUnusedVars(env)
			}
			core.ResetUsage(env)
			env.SetCurrentNamespace(ns)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if processErr == nil && reportGloballyUnused {
		core.WarnOnGloballyUnusedNamespaces(env)
		core.WarnOnGloballyUnusedVars(env)
	}

	return nil
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

func usage(out io.Writer) {
	fmt.Fprintf(out, "Joker - %s\n\n", core.VERSION)
	fmt.Fprintln(out, "Usage: lace [args] [-- <repl-args>]                starts a repl")
	fmt.Fprintln(out, "   or: lace [args] --repl [<socket>] [-- <repl-args>]")
	fmt.Fprintln(out, "                                                    starts a repl (on optional network socket)")
	fmt.Fprintln(out, "   or: lace [args] --eval <expr> [-- <expr-args>]  evaluate <expr>, print if non-nil")
	fmt.Fprintln(out, "   or: lace [args] [--file] <filename> [<script-args>]")
	fmt.Fprintln(out, "                                                    input from file")
	fmt.Fprintln(out, "   or: lace [args] --lint <filename>               lint the code in file")
	fmt.Fprintln(out, "\nNotes:")
	fmt.Fprintln(out, "  -e is a synonym for --eval.")
	fmt.Fprintln(out, "  '-' for <filename> means read from standard input (stdin).")
	fmt.Fprintln(out, "  Evaluating '(println (str *command-line-args*))' prints the arguments")
	fmt.Fprintln(out, "    in <repl-args>, <expr-args>, or <script-args> (TBD).")
	fmt.Fprintln(out, "  <socket> is passed to Go's net.Listen() function. If multiple --*repl options are specified,")
	fmt.Fprintln(out, "    the final one specified \"wins\".")

	fmt.Fprintln(out, "\nOptions (<args>):")
	fmt.Fprintln(out, "  --help, -h")
	fmt.Fprintln(out, "    Print this help message and exit.")
	fmt.Fprintln(out, "  --version, -v")
	fmt.Fprintln(out, "    Print version number and exit.")
	fmt.Fprintln(out, "  --read")
	fmt.Fprintln(out, "    Read, but do not parse nor evaluate, the input.")
	fmt.Fprintln(out, "  --parse")
	fmt.Fprintln(out, "    Read and parse, but do not evaluate, the input.")
	fmt.Fprintln(out, "  --evaluate")
	fmt.Fprintln(out, "    Read, parse, and evaluate the input (default unless --lint in effect).")
	fmt.Fprintln(out, "  --exit-to-repl [<socket>]")
	fmt.Fprintln(out, "    After successfully processing --eval or --file, drop into repl instead of exiting.")
	fmt.Fprintln(out, "  --error-to-repl [<socket>]")
	fmt.Fprintln(out, "    After failure processing --eval or --file, drop into repl instead of exiting.")
	fmt.Fprintln(out, "  --no-readline")
	fmt.Fprintln(out, "    Disable readline functionality in the repl. Useful when using rlwrap.")
	fmt.Fprintln(out, "  --working-dir <directory>")
	fmt.Fprintln(out, "    Specify directory to lint or working directory for lint configuration if linting single file (requires --lint).")
	fmt.Fprintln(out, "  --report-globally-unused")
	fmt.Fprintln(out, "    Report globally unused namespaces and public vars when linting directories (requires --lint and --working-dir).")
	fmt.Fprintln(out, "  --dialect <dialect>")
	fmt.Fprintln(out, "    Set input dialect (\"clj\", \"cljs\", \"lace\", \"edn\") for linting;")
	fmt.Fprintln(out, "    default is inferred from <filename> suffix, if any.")
	fmt.Fprintln(out, "  --hashmap-threshold <n>")
	fmt.Fprintln(out, "    Set HASHMAP_THRESHOLD accordingly (internal magic of some sort).")
	fmt.Fprintln(out, "  --profiler <type>")
	fmt.Fprintln(out, "    Specify type of profiler to use (default 'runtime/pprof' or 'pkg/profile').")
	fmt.Fprintln(out, "  --cpuprofile <name>")
	fmt.Fprintln(out, "    Write CPU profile to specified file or directory (depending on")
	fmt.Fprintln(out, "    profiler chosen).")
	fmt.Fprintln(out, "  --cpuprofile-rate <rate>")
	fmt.Fprintln(out, "    Specify rate (hz, aka samples per second) for the 'runtime/pprof' CPU")
	fmt.Fprintln(out, "    profiler to use.")
	fmt.Fprintln(out, "  --memprofile <name>")
	fmt.Fprintln(out, "    Write memory profile to specified file.")
	fmt.Fprintln(out, "  --memprofile-rate <rate>")
	fmt.Fprintln(out, "    Specify rate (one sample per <rate>) for the memory profiler to use.")
}

var (
	debugOut                 io.Writer
	helpFlag                 bool
	versionFlag              bool
	phase                    core.Phase = core.EVAL // --read, --parse, --evaluate
	workingDir               string
	lintFlag                 bool
	reportGloballyUnusedFlag bool
	dialect                  core.Dialect = core.UNKNOWN
	eval                     string
	replFlag                 bool
	replSocket               string
	classPath                string
	filename                 string
	remainingArgs            []string
	profilerType             string = "runtime/pprof"
	cpuProfileName           string
	cpuProfileRate           int
	cpuProfileRateFlag       bool
	memProfileName           string
	noReadline               bool
	exitToRepl               bool
	errorToRepl              bool
)

func isNumber(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func notOption(arg string) bool {
	return arg == "-" || !strings.HasPrefix(arg, "-") || isNumber(arg[1:])
}

func parseArgs(args []string) {
	if len(args) > 1 {
		// peek to see if the first arg is "--debug*"
		switch args[1] {
		case "--debug", "--debug=stderr":
			debugOut = core.Stderr
		case "--debug=stdout":
			debugOut = core.Stdout
		}
	}

	length := len(args)
	stop := false
	missing := false
	noFileFlag := false
	if v, ok := os.LookupEnv("JOKER_CLASSPATH"); ok {
		classPath = v
	} else {
		classPath = ""
	}
	var i int
	for i = 1; i < length; i++ { // shift
		if debugOut != nil {
			fmt.Fprintf(debugOut, "arg[%d]=%s\n", i, args[i])
		}
		switch args[i] {
		case "-": // denotes stdin
			stop = true
		case "--": // formally ends options processing
			stop = true
			noFileFlag = true
			i += 1 // do not include "--" in *command-line-args*
		case "--debug":
			debugOut = core.Stderr
		case "--debug=stderr":
			debugOut = core.Stderr
		case "--debug=stdout":
			debugOut = core.Stdout
		case "--verbose":
			if i < length-1 && notOption(args[i+1]) {
				i += 1 // shift
				verbosity, err := strconv.ParseInt(args[i], 10, 64)
				if err != nil {
					fmt.Fprintln(core.Stderr, "Error: ", err)
					return
				}
				if verbosity <= 0 {
					core.VerbosityLevel = 0
				} else {
					core.VerbosityLevel = int(verbosity)
				}
			} else {
				core.VerbosityLevel++
			}
		case "--help", "-h":
			helpFlag = true
			return // don't bother parsing anything else
		case "--version", "-v":
			versionFlag = true
		case "--read":
			phase = core.READ
		case "--parse":
			phase = core.PARSE
		case "--evaluate":
			phase = core.EVAL
		case "--working-dir":
			if i < length-1 && notOption(args[i+1]) {
				i += 1 // shift
				workingDir = args[i]
			} else {
				missing = true
			}
		case "--report-globally-unused":
			reportGloballyUnusedFlag = true
		case "--lint":
			lintFlag = true
		case "--lintclj":
			lintFlag = true
			dialect = core.CLJ
		case "--lintcljs":
			lintFlag = true
			dialect = core.CLJS
		case "--lintlace":
			lintFlag = true
			dialect = core.LACE
		case "--lintedn":
			lintFlag = true
			dialect = core.EDN
		case "--dialect":
			if i < length-1 && notOption(args[i+1]) {
				i += 1 // shift
				dialect = dialectFromArg(args[i])
			} else {
				missing = true
			}
		case "--hashmap-threshold":
			if i < length-1 && notOption(args[i+1]) {
				i += 1 // shift
				thresh, err := strconv.ParseInt(args[i], 10, 64)
				if err != nil {
					fmt.Fprintln(core.Stderr, "Error: ", err)
					return
				}
				if thresh < 0 {
					core.HASHMAP_THRESHOLD = math.MaxInt64
				} else {
					core.HASHMAP_THRESHOLD = thresh
				}
			} else {
				missing = true
			}
		case "-e", "--eval":
			if i < length-1 && notOption(args[i+1]) {
				i += 1 // shift
				eval = args[i]
				phase = core.PRINT_IF_NOT_NIL
			} else {
				missing = true
			}
		case "--repl":
			replFlag = true
			if i < length-1 && notOption(args[i+1]) {
				i += 1 // shift
				replSocket = args[i]
			}
		case "-c", "--classpath":
			if i < length-1 && notOption(args[i+1]) {
				i += 1 // shift
				classPath = args[i]
			} else {
				missing = true
			}
		case "--no-readline":
			noReadline = true
		case "--exit-to-repl":
			exitToRepl = true
			if i < length-1 && notOption(args[i+1]) {
				i += 1 // shift
				replSocket = args[i]
			}
		case "--error-to-repl":
			errorToRepl = true
			if i < length-1 && notOption(args[i+1]) {
				i += 1 // shift
				replSocket = args[i]
			}
		case "--file":
			if i < length-1 && notOption(args[i+1]) {
				i += 1 // shift
				filename = args[i]
			}
		case "--profiler":
			if i < length-1 && notOption(args[i+1]) {
				i += 1 // shift
				profilerType = args[i]
			} else {
				missing = true
			}
		case "--cpuprofile":
			if i < length-1 && notOption(args[i+1]) {
				i += 1 // shift
				cpuProfileName = args[i]
			} else {
				missing = true
			}
		case "--cpuprofile-rate":
			if i < length-1 && notOption(args[i+1]) {
				i += 1 // shift
				rate, err := strconv.Atoi(args[i])
				if err != nil {
					fmt.Fprintln(core.Stderr, "Error: ", err)
					return
				}
				if rate > 0 {
					cpuProfileRate = rate
					cpuProfileRateFlag = true
				}
			} else {
				missing = true
			}
		case "--memprofile":
			if i < length-1 && notOption(args[i+1]) {
				i += 1 // shift
				memProfileName = args[i]
			} else {
				missing = true
			}
		case "--memprofile-rate":
			if i < length-1 && notOption(args[i+1]) {
				i += 1 // shift
				rate, err := strconv.Atoi(args[i])
				if err != nil {
					fmt.Fprintln(core.Stderr, "Error: ", err)
					return
				}
				if rate > 0 {
					runtime.MemProfileRate = rate
				}
			} else {
				missing = true
			}
		default:
			if strings.HasPrefix(args[i], "-") {
				fmt.Fprintf(core.Stderr, "Error: Unrecognized option '%s'\n", args[i])
				core.Exit(2)
			}
			stop = true
		}
		if stop || missing {
			break
		}
	}
	if missing {
		fmt.Fprintf(core.Stderr, "Error: Missing argument for '%s' option\n", args[i])
		core.Exit(3)
	}
	if i < length && !noFileFlag && filename == "" {
		if debugOut != nil {
			fmt.Fprintf(debugOut, "filename=%s\n", args[i])
		}
		filename = args[i]
		i += 1 // shift
	}
	if i < length {
		if debugOut != nil {
			fmt.Fprintf(debugOut, "remaining=%v\n", args[i:])
		}
		remainingArgs = args[i:]
	}
}

var runningProfile interface {
	Stop()
}

func main() {
	core.SetExit(func(code int) {
		finish()
		os.Exit(code)
	})

	env, err := core.NewEnv()
	if err != nil {
		fmt.Printf("unable to initialize environment: %s", err)
		os.Exit(1)
	}

	env.InitEnv(core.Stdin, core.Stdout, core.Stderr, os.Args[1:])

	parseArgs(os.Args) // Do this early enough so --verbose can show lace.core being processed.

	saveForRepl = saveForRepl && (exitToRepl || errorToRepl) // don't bother saving stuff if no repl

	// SetupGlobalEnvCoreData()

	env.ReferCoreToUser()
	env.SetEnvArgs(remainingArgs)
	env.SetClassPath(classPath)

	if debugOut != nil {
		fmt.Fprintf(debugOut, "debugOut=%v\n", debugOut)
		fmt.Fprintf(debugOut, "helpFlag=%v\n", helpFlag)
		fmt.Fprintf(debugOut, "versionFlag=%v\n", versionFlag)
		fmt.Fprintf(debugOut, "phase=%v\n", phase)
		fmt.Fprintf(debugOut, "lintFlag=%v\n", lintFlag)
		fmt.Fprintf(debugOut, "reportGloballyUnusedFlag=%v\n", reportGloballyUnusedFlag)
		fmt.Fprintf(debugOut, "dialect=%v\n", dialect)
		fmt.Fprintf(debugOut, "workingDir=%v\n", workingDir)
		fmt.Fprintf(debugOut, "HASHMAP_THRESHOLD=%v\n", core.HASHMAP_THRESHOLD)
		fmt.Fprintf(debugOut, "eval=%v\n", eval)
		fmt.Fprintf(debugOut, "replFlag=%v\n", replFlag)
		fmt.Fprintf(debugOut, "replSocket=%v\n", replSocket)
		fmt.Fprintf(debugOut, "classPath=%v\n", classPath)
		fmt.Fprintf(debugOut, "noReadline=%v\n", noReadline)
		fmt.Fprintf(debugOut, "filename=%v\n", filename)
		fmt.Fprintf(debugOut, "remainingArgs=%v\n", remainingArgs)
		fmt.Fprintf(debugOut, "exitToRepl=%v\n", exitToRepl)
		fmt.Fprintf(debugOut, "errorToRepl=%v\n", errorToRepl)
		fmt.Fprintf(debugOut, "saveForRepl=%v\n", saveForRepl)
	}

	if helpFlag {
		usage(core.Stdout)
		return
	}

	if versionFlag {
		println(core.VERSION)
		return
	}

	if len(remainingArgs) > 0 {
		if lintFlag {
			fmt.Fprintf(core.Stderr, "Error: Cannot provide arguments to code while linting it.\n")
			core.Exit(4)
		}
		if phase != core.EVAL && phase != core.PRINT_IF_NOT_NIL {
			fmt.Fprintf(core.Stderr, "Error: Cannot provide arguments to code without evaluating it.\n")
			core.Exit(5)
		}
	}

	/* Set up profiling. */

	if cpuProfileName != "" {
		switch profilerType {
		case "pkg/profile":
			runningProfile = profile.Start(profile.ProfilePath(cpuProfileName))
			defer finish()
		case "runtime/pprof":
			f, err := os.Create(cpuProfileName)
			if err != nil {
				fmt.Fprintf(core.Stderr, "Error: Could not create CPU profile `%s': %v\n",
					cpuProfileName, err)
				cpuProfileName = ""
				core.Exit(96)
			}
			if cpuProfileRateFlag {
				runtime.SetCPUProfileRate(cpuProfileRate)
			}
			err = pprof.StartCPUProfile(f)
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(core.Stderr, "Profiling started at rate=%d. See file `%s'.\n",
				cpuProfileRate, cpuProfileName)
			defer finish()
		default:
			fmt.Fprintf(core.Stderr,
				"Unrecognized profiler: %s\n  Use 'pkg/profile' or 'runtime/pprof'.\n",
				profilerType)
			core.Exit(96)
		}
	} else if memProfileName != "" {
		defer finish()
	}

	if eval != "" {
		if lintFlag {
			fmt.Fprintf(core.Stderr, "Error: Cannot combine --eval/-e and --lint.\n")
			core.Exit(6)
		}
		if replFlag {
			fmt.Fprintf(core.Stderr, "Error: Cannot combine --eval/-e and --repl.\n")
			core.Exit(7)
		}
		if workingDir != "" {
			fmt.Fprintf(core.Stderr, "Error: Cannot combine --eval/-e and --working-dir.\n")
			core.Exit(8)
		}
		if reportGloballyUnusedFlag {
			fmt.Fprintf(core.Stderr, "Error: Cannot combine --eval/-e and --report-globally-unused.\n")
			core.Exit(17)
		}
		if filename != "" {
			fmt.Fprintf(core.Stderr, "Error: Cannot combine --eval/-e and a <filename> argument.\n")
			core.Exit(9)
		}
		reader := core.NewReader(strings.NewReader(eval), "<expr>")
		if saveForRepl {
			reader = core.NewReader(&replayable{reader}, "<replay>")
		}
		if err := core.ProcessReader(env, reader, "", phase); err != nil {
			if !errorToRepl {
				core.Exit(1)
			}
		} else {
			if !exitToRepl {
				return
			}
		}
	}

	if lintFlag {
		if replFlag {
			fmt.Fprintf(core.Stderr, "Error: Cannot combine --lint and --repl.\n")
			core.Exit(10)
		}
		if exitToRepl {
			fmt.Fprintf(core.Stderr, "Error: Cannot combine --lint and --exit-to-repl.\n")
			core.Exit(14)
		}
		if errorToRepl {
			fmt.Fprintf(core.Stderr, "Error: Cannot combine --lint and --error-to-repl.\n")
			core.Exit(15)
		}
		if dialect == core.UNKNOWN {
			dialect = detectDialect(filename)
		}
		if filename != "" {
			err := lintFile(env, filename, dialect, workingDir)
			if err != nil {
				fmt.Fprintf(core.Stderr, "Error linting file: %s\n", err)
			}
		} else if workingDir != "" {
			err := lintDir(env, workingDir, dialect, reportGloballyUnusedFlag)
			if err != nil {
				fmt.Fprintf(core.Stderr, "Error linting dir: %s\n", err)
			}
		} else {
			fmt.Fprintf(core.Stderr, "Error: Missing --file or --working-dir argument.\n")
			core.Exit(16)
		}
		if core.PROBLEM_COUNT > 0 {
			core.Exit(1)
		}
		return
	}

	if workingDir != "" {
		fmt.Fprintf(core.Stderr, "Error: Cannot specify --working-dir option when not linting.\n")
		core.Exit(11)
	}

	if filename != "" {
		if err := processFile(env, filename, phase); err != nil {
			if !errorToRepl {
				core.Exit(1)
			}
		} else {
			if !exitToRepl {
				return
			}
		}
	}

	if replSocket != "" {
		err = srepl(env, replSocket, phase)
	} else {
		err = repl(env, phase)
	}

	if ee, ok := err.(*core.ExitError); ok {
		os.Exit(ee.Code)
	}
}

func finish() {
	if runningProfile != nil {
		runningProfile.Stop()
		runningProfile = nil
	} else if cpuProfileName != "" {
		pprof.StopCPUProfile()
		fmt.Fprintf(core.Stderr, "Profiling stopped. See file `%s'.\n", cpuProfileName)
		cpuProfileName = ""
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
