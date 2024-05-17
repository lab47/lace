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
		Fn: func() logger.Logger {
			return log
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

func convertObj(env *core.Env, obj core.Object) (any, error) {
	switch sv := obj.(type) {
	case core.Keyword:
		return sv.Name(), nil
	case core.Number:
		return sv.NativeNumber(), nil
	case core.String:
		return sv.S(), nil
	case core.Map:
		ret := map[any]any{}

		slice, err := core.ToSlice(env, sv.Keys())
		if err != nil {
			return nil, err
		}

		for _, k := range slice {
			ok, v, err := sv.Get(env, k)
			if err != nil {
				return nil, err
			}

			if !ok {
				return core.ToNative(env, obj)
			}

			ko, err := convertObj(env, k)
			if err != nil {
				return nil, err
			}

			vo, err := convertObj(env, v)
			if err != nil {
				return nil, err
			}

			ret[ko] = vo
		}

		return ret, nil
	case core.Seqable:
		var ret []any
		slice, err := core.ToSlice(env, sv.Seq())
		if err != nil {
			return nil, err
		}
		for _, o := range slice {
			co, err := convertObj(env, o)
			if err != nil {
				return nil, err
			}
			ret = append(ret, co)
		}
		return ret, nil
	default:
		return core.ToNative(env, obj)
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

	vals, err := core.ToSlice(env, seq.Seq())
	if err != nil {
		return err
	}

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
				args = append(args, sv.S())
				expectKey = false
			case core.Seqable:
				objs, err := core.ToSlice(env, sv.Seq())
				if err != nil {
					return err
				}
				if len(objs) == 2 {
					a, err := convertObj(env, objs[0])
					if err != nil {
						return err
					}
					b, err := convertObj(env, objs[1])
					if err != nil {
						return err
					}
					args = append(args, a, b)
				} else {
					return env.NewError("key value must be keyword, symbol, or string")
				}
				expectKey = true
			default:
				return env.NewError("key value must be keyword, symbol, or string")
			}
		} else {
			co, err := convertObj(env, r)
			if err != nil {
				return err
			}
			args = append(args, co)
			expectKey = true
		}
	}

	log.Log(context.Background(), sl, message, args...)

	return nil
}
