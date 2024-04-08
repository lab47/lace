package time

import (
	"time"

	"github.com/lab47/lace/core"
)

func Setup(env *core.Env) error {
	b := core.NewNSBuilder(env, "lace.time")

	b.Defn(&core.DefnInfo{
		Name:  "now",
		Doc:   "Returns the current local time.",
		Added: "1.0",
		Tag:   "Time",
		Fn:    time.Now,
	})

	b.Defn(&core.DefnInfo{
		Name:  "from-unix",
		Args:  []string{"sec", "nsec"},
		Doc:   "Returns the local Time corresponding to the given Unix time, sec seconds and nsec nanoseconds since January 1, 1970 UTC. It is valid to pass nsec outside the range [0, 999999999].",
		Added: "1.0",
		Tag:   "Time",
		Fn:    time.Unix,
	})

	b.Defn(&core.DefnInfo{
		Name: "unix",
		Args: []string{"t"},
		Doc:  "Returns t as a Unix time, the number of seconds elapsed since January 1, 1970 UTC.",
		Tag:  "Int",
		Fn: func(t time.Time) int64 {
			return t.Unix()
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "unix-nano",
		Args: []string{"t"},
		Doc:  "Returns t as a Unix time, the number of nanoseconds elapsed since January 1, 1970 UTC.",
		Tag:  "Int",
		Fn: func(t time.Time) int64 {
			return t.UnixNano()
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "nanosecond",
		Args: []string{"t"},
		Doc:  "Return the number of nanoseconds represented in t.",
		Tag:  "Int",
		Fn: func(t time.Time) int {
			return t.Nanosecond()
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "sub",
		Args: []string{"t", "t2"},
		Doc:  "Returns the duration t-t2 in nanoseconds.",
		Tag:  "Int",
		Fn: func(t, t2 time.Time) int {
			return int(t.Sub(t2))
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "add",
		Args: []string{"t", "d"},
		Doc:  "Returns the time t+d",
		Tag:  "Int",
		Fn: func(t time.Time, d int) time.Time {
			return t.Add(time.Duration(d))
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "parse",
		Args: []string{"layout", "string"},
		Doc:  "Parse a string according to layout to result in a time.",
		Tag:  "Time",
		Fn:   time.Parse,
	})

	b.Defn(&core.DefnInfo{
		Name: "format",
		Args: []string{"t", "layout"},
		Doc:  "Format the time t according to the layout as a string.",
		Tag:  "String",
		Fn: func(t time.Time, layout string) string {
			return t.Format(layout)
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "duration-string",
		Args: []string{"d"},
		Doc:  "Format the given duration as a string.",
		Tag:  "String",
		Fn: func(d int) string {
			return time.Duration(d).String()
		},
	})

	b.Defn(&core.DefnInfo{
		Name:  "sleep",
		Doc:   "Pauses the execution thread for at least the duration d (expressed in nanoseconds). A negative or zero duration causes sleep to return immediately.",
		Added: "1.0",
		Fn: func(i int) {
			time.Sleep(time.Duration(i))
		},
	})

	b.Def("rfc3339", core.MakeString(time.RFC3339))
	b.Def("rfc3339-nano", core.MakeString(time.RFC3339Nano))

	return nil
}

func init() {
	core.AddNativeNamespace("lace.time", Setup)
}
