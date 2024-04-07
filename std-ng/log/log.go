package log

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/lab47/lace/core"
	"github.com/lab47/lsvd/logger"
)

var levelMap map[string]slog.Level

func Setup(env *core.Env) error {
	b := core.NewNSBuilder(env, "lace.log")

	log := logger.New(logger.Info)

	b.Defn(&core.DefnInfo{
		Name: "default",
		Doc:  "Returns the default logger object",
		Fn: func() core.Object {
			return core.MakeOpaque(log)
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "set-logger-level",
		Doc:  "Set the level of the logger.",
		Args: []string{"logger", "level"},
		Fn:   setLevel,
	})

	b.Defn(&core.DefnInfo{
		Name: "emit",
		Doc:  "Log the given keys and values to the default logger",
		Args: []string{"logger", "level", "message", "pairs"},
		Fn:   emit,
	})

	levelMap = map[string]slog.Level{
		"trace": logger.Trace,
		"debug": logger.Debug,
		"info":  logger.Info,
		"warn":  logger.Warn,
		"error": logger.Error,
	}

	return b.Run(code)
}

//go:embed log.clj
var code []byte

func init() {
	core.AddNativeNamespace("lace.log", Setup)
}

func convertObj(obj core.Object) any {
	switch sv := obj.(type) {
	case core.Keyword:
		return sv.Name()
	case core.Int:
		return sv.I
	case core.String:
		return sv.S
	case core.Map:
		ret := map[any]any{}

		for _, k := range core.ToSlice(sv.Keys()) {
			ok, v := sv.Get(k)
			if !ok {
				return core.ToNative(obj)
			}

			ret[convertObj(k)] = convertObj(v)
		}

		return ret
	case core.Seqable:
		var ret []any
		for _, o := range core.ToSlice(sv.Seq()) {
			ret = append(ret, convertObj(o))
		}
		return ret
	default:
		return core.ToNative(obj)
	}
}

func setLevel(env *core.Env, logobj core.Object, level core.Keyword) error {
	var log logger.Logger

	err := core.ExtractOpaque(env, logobj, &log)
	if err != nil {
		return err
	}

	sl, ok := levelMap[level.Name()]
	if !ok {
		return fmt.Errorf("unknown level: %s", level.Name())
	}

	log.SetLevel(sl)
	return nil
}

func emit(env *core.Env, logobj core.Object, level core.Keyword, message string, seq core.Seqable) error {
	var log logger.Logger

	err := core.ExtractOpaque(env, logobj, &log)
	if err != nil {
		return err
	}

	sl, ok := levelMap[level.Name()]
	if !ok {
		sl = slog.LevelWarn
	}

	vals := core.ToSlice(seq.Seq())

	var args []any
	expectKey := true

	for _, r := range vals {
		if expectKey {
			switch sv := r.(type) {
			case core.Symbol:
				args = append(args, sv.Name())
				expectKey = false
			case core.Keyword:
				args = append(args, sv.Name())
				expectKey = false
			case core.String:
				args = append(args, sv.S)
				expectKey = false
			case core.Seqable:
				objs := core.ToSlice(sv.Seq())
				if len(objs) == 2 {
					args = append(args, convertObj(objs[0]), convertObj(objs[1]))
				} else {
					return env.RT.NewError("key value must be keyword, symbol, or string")
				}
				expectKey = true
			default:
				return env.RT.NewError("key value must be keyword, symbol, or string")
			}
		} else {
			args = append(args, convertObj(r))
			expectKey = true
		}
	}

	log.Log(context.Background(), sl, message, args...)

	return nil
}
