package os

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/lab47/lace/core"
)

func env(env *Env) (Object, error) {
	res := EmptyArrayMap()
	for _, v := range os.Environ() {
		parts := strings.Split(v, "=")
		res.Add(env, String{S: parts[0]}, String{S: parts[1]})
	}
	return res, nil
}

func setEnv(key string, value string) (Object, error) {
	err := os.Setenv(key, value)
	if err != nil {
		return nil, err
	}
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
	if ok, dirObj := opts.GetEqu(MakeKeyword("dir")); ok {
		dirv, err := AssertString(env, dirObj, "dir must be a string")
		if err != nil {
			return nil, err
		}

		dir = dirv.S
	}
	if ok, argsObj := opts.GetEqu(MakeKeyword("args")); ok {
		sv, err := AssertSeqable(env, argsObj, "args must be Seqable")
		if err != nil {
			return nil, err
		}

		s := sv.Seq()
		for !s.IsEmpty() {
			f, err := s.First(env)
			if err != nil {
				return nil, err
			}
			so, err := AssertString(env, f, "args must be strings")
			if err != nil {
				return nil, err
			}
			args = append(args, so.S)
			s = s.Rest()
		}
	}
	if ok, stdinObj := opts.GetEqu(MakeKeyword("stdin")); ok {
		// Check if the intent was to pipe stdin into the program being called and
		// use Stdin directly rather than env.stdin.Value, which is a buffered wrapper.
		// TODO: this won't work correctly if env.stdin is bound to something other than Stdin
		if env.IsStdIn(stdinObj) {
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
				return nil, env.RT.NewError("stdin option must be either an IOReader or a string, got " + stdinObj.GetType().Name())
			}
		}
	}
	if ok, stdoutObj := opts.GetEqu(MakeKeyword("stdout")); ok {
		switch s := stdoutObj.(type) {
		case Nil:
		case *IOWriter:
			stdout = s.Writer
		case io.Writer:
			stdout = s
		default:
			return nil, env.RT.NewError("stdout option must be an IOWriter, got " + stdoutObj.GetType().Name())
		}
	}
	if ok, stderrObj := opts.GetEqu(MakeKeyword("stderr")); ok {
		switch s := stderrObj.(type) {
		case Nil:
		case *IOWriter:
			stderr = s.Writer
		case io.Writer:
			stderr = s
		default:
			return nil, env.RT.NewError("stderr option must be an IOWriter, got " + stderrObj.GetType().Name())
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
	if err != nil {
		return nil, err
	}
	res := EmptyVector()
	name := MakeKeyword("name")
	size := MakeKeyword("size")
	mode := MakeKeyword("mode")
	isDir := MakeKeyword("dir?")
	modTime := MakeKeyword("modtime")
	for _, f := range files {
		m := EmptyArrayMap()
		m.AddEqu(name, MakeString(f.Name()))
		m.AddEqu(size, MakeInt(int(f.Size())))
		m.AddEqu(mode, MakeInt(int(f.Mode())))
		m.AddEqu(isDir, MakeBoolean(f.IsDir()))
		m.AddEqu(modTime, MakeInt(int(f.ModTime().Unix())))
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

func stat(env *Env, filename string) (Object, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	return FileInfoMap(env, info.Name(), info), nil
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
