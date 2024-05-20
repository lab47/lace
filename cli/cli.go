package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/pprof"

	"github.com/lab47/lablog/logger"
	"github.com/lab47/lace/core"
	"github.com/lab47/lace/pkg/build"
	_ "github.com/lab47/lace/std-ng/all"
	"github.com/spf13/pflag"
	"golang.org/x/sys/unix"
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

func (ctx *ReplContext) PushValue(obj any) {
	ctx.third.SetStatic(ctx.second.GetStatic())
	ctx.second.SetStatic(ctx.first.GetStatic())
	ctx.first.SetStatic(obj)
}

func (ctx *ReplContext) PushException(exc any) {
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

	_, err = cl.Call(env, []any{})
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

	log := logger.New(logger.Info)

	//_ = os.Args[0]
	args := os.Args[1:]
	switch {
	case len(args) >= 1:
		if _, err := os.Stat(args[0]); err == nil {
			args = append([]string{"run"}, args...)
		}
	case len(args) == 0:
		args = []string{"repl"}
	}

	cmd := args[0]
	args = args[1:]

	switch cmd {
	case "run":
		if len(args) >= 1 {
			if dir := isProject(args[0]); dir != "" {
				err = os.Chdir(dir)
				if err != nil {
					log.Error("unable to change to project directory", "error", err, "dir", dir)
					os.Exit(1)
				}
				runInProject(log, dir, env, args[1:])
				return
			}
		}
		if dir := findProject(); dir == "" {
			run(env, args)
			return
		} else {
			runInProject(log, dir, env, args)
			return
		}
	case "repl":
		_ = env.REPL(os.Stdin, os.Stdout)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
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
	}
}

func isProject(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return ""
	}

	_, err = os.Stat(filepath.Join(abs, "lace.yml"))
	if err == nil {
		return abs
	}

	return ""
}

func findProject() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for dir != "/" {
		if _, err := os.Stat(filepath.Join(dir, "lace.yml")); err == nil {
			return dir
		}

		dir = filepath.Dir(dir)
	}

	return ""
}

func runInProject(log logger.Logger, dir string, env *core.Env, args []string) {
	b, err := build.LoadBuilder(log, dir)
	if err != nil {
		log.Error("error loading project builder", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	exe, err := b.Run(ctx)
	if err != nil {
		log.Error("error running build", "error", err)
		os.Exit(1)
	}

	argv := append([]string{exe}, args...)

	err = unix.Exec(exe, argv, os.Environ())

	log.Error("error executing exe", "error", err)

	cmd := exec.Command(exe, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Error("error running compiled app", "error", err)
		os.Exit(0)
	}

	os.Exit(1)
}

func run(env *core.Env, args []string) {
	fs := pflag.NewFlagSet("lace", pflag.ExitOnError)
	version := fs.BoolP("version", "v", false, "report the version number")
	cpuProfile := fs.String("cpuprofile", "", "Write CPU profile info to the specified path")
	cpuProfileRate := fs.Int("cpuprofile-rate", 100, "Specify the sampling rate of the cpu profiler")
	memProfile := fs.String("memprofile", "", "Write Memory profile info to the specified path")
	debugBytecode := fs.Bool("debug-bytecode", false, "Display bytecode for functions are it is generated")

	if err := fs.Parse(args); err != nil {
		fmt.Printf("error parsing arguments: %s\n", err)
		os.Exit(1)
	}

	var filename string

	if fs.NArg() >= 1 {
		filename = fs.Arg(0)
		args = fs.Args()[1:]
	}

	env.InitEnv(core.Stdin, core.Stdout, core.Stderr, args)

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

}
