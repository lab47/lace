package os

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/candid82/joker/core"
)

func env() (Object, error) {
	res := EmptyArrayMap()
	for _, v := range os.Environ() {
		parts := strings.Split(v, "=")
		res.Add(String{S: parts[0]}, String{S: parts[1]})
	}
	return res, nil
}

func setEnv(key string, value string) (Object, error) {
	err := os.Setenv(key, value)
	PanicOnErr(err)
	return NIL, nil
}

func getEnv(key string) (Object, error) {
	if v, ok := os.LookupEnv(key); ok {
		return MakeString(v), nil
	}
	return NIL, nil
}

func commandArgs() (Object, error) {
	res := EmptyVector()
	var err error
	for _, arg := range os.Args {
		res, err = res.Conjoin(String{S: arg})
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

const defaultFailedCode = 127 // seen from 'sh no-such-file' on OS X and Ubuntu

func execute(env *Env, name string, opts Map) (Object, error) {
	var dir string
	var args []string
	var stdin io.Reader
	var stdout, stderr io.Writer
	if ok, dirObj := opts.Get(MakeKeyword("dir")); ok {
		dirv, err := AssertString(env, dirObj, "dir must be a string")
		if err != nil {
			return nil, err
		}

		dir = dirv.S
	}
	if ok, argsObj := opts.Get(MakeKeyword("args")); ok {
		sv, err := AssertSeqable(env, argsObj, "args must be Seqable")
		if err != nil {
			return nil, err
		}

		s := sv.Seq()
		for !s.IsEmpty() {
			so, err := AssertString(env, s.First(), "args must be strings")
			if err != nil {
				return nil, err
			}
			args = append(args, so.S)
			s = s.Rest()
		}
	}
	if ok, stdinObj := opts.Get(MakeKeyword("stdin")); ok {
		// Check if the intent was to pipe stdin into the program being called and
		// use Stdin directly rather than GLOBAL_ENV.stdin.Value, which is a buffered wrapper.
		// TODO: this won't work correctly if GLOBAL_ENV.stdin is bound to something other than Stdin
		if GLOBAL_ENV.IsStdIn(stdinObj) {
			stdin = Stdin
		} else {
			switch s := stdinObj.(type) {
			case Nil:
			case *IOReader:
				stdin = s.Reader
			case io.Reader:
				stdin = s
			case String:
				stdin = strings.NewReader(s.S)
			default:
				panic(StubNewError("stdin option must be either an IOReader or a string, got " + stdinObj.GetType().ToString(false)))
			}
		}
	}
	if ok, stdoutObj := opts.Get(MakeKeyword("stdout")); ok {
		switch s := stdoutObj.(type) {
		case Nil:
		case *IOWriter:
			stdout = s.Writer
		case io.Writer:
			stdout = s
		default:
			panic(StubNewError("stdout option must be an IOWriter, got " + stdoutObj.GetType().ToString(false)))
		}
	}
	if ok, stderrObj := opts.Get(MakeKeyword("stderr")); ok {
		switch s := stderrObj.(type) {
		case Nil:
		case *IOWriter:
			stderr = s.Writer
		case io.Writer:
			stderr = s
		default:
			panic(StubNewError("stderr option must be an IOWriter, got " + stderrObj.GetType().ToString(false)))
		}
	}
	return sh(dir, stdin, stdout, stderr, name, args)
}

func mkdir(name string, perm int) (Object, error) {
	err := os.Mkdir(name, os.FileMode(perm))
	return NIL, err
}

func readDir(dirname string) (Object, error) {
	files, err := ioutil.ReadDir(dirname)
	PanicOnErr(err)
	res := EmptyVector()
	name := MakeKeyword("name")
	size := MakeKeyword("size")
	mode := MakeKeyword("mode")
	isDir := MakeKeyword("dir?")
	modTime := MakeKeyword("modtime")
	for _, f := range files {
		m := EmptyArrayMap()
		m.Add(name, MakeString(f.Name()))
		m.Add(size, MakeInt(int(f.Size())))
		m.Add(mode, MakeInt(int(f.Mode())))
		m.Add(isDir, MakeBoolean(f.IsDir()))
		m.Add(modTime, MakeInt(int(f.ModTime().Unix())))
		res, err = res.Conjoin(m)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func getwd() (string, error) {
	res, err := os.Getwd()
	return res, err
}

func chdir(dirname string) (Object, error) {
	err := os.Chdir(dirname)
	return NIL, err
}

func stat(filename string) (Object, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	return FileInfoMap(info.Name(), info), nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, StubNewError(err.Error())
}
