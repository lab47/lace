// This file is generated by generate-std.joke script. Do not edit manually!

package time

import (
	"fmt"
	"os"

	. "github.com/lab47/lace/core"
)

func InternsOrThunks() {
	if VerbosityLevel > 0 {
		fmt.Fprintln(os.Stderr, "Lazily running slow version of time.InternsOrThunks().")
	}
	timeNamespace.ResetMeta(MakeMeta(nil, `Provides functionality for measuring and displaying time.`, "1.0"))

	timeNamespace.InternVar("ansi-c", ansi_c_,
		MakeMeta(
			nil,
			`Mon Jan _2 15:04:05 2006`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("hour", hour_,
		MakeMeta(
			nil,
			`Number of nanoseconds in 1 hour`, "1.0").Plus(MakeKeyword("tag"), String{S: "BigInt"}))

	timeNamespace.InternVar("kitchen", kitchen_,
		MakeMeta(
			nil,
			`3:04PM`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("microsecond", microsecond_,
		MakeMeta(
			nil,
			`Number of nanoseconds in 1 microsecond`, "1.0").Plus(MakeKeyword("tag"), String{S: "Int"}))

	timeNamespace.InternVar("millisecond", millisecond_,
		MakeMeta(
			nil,
			`Number of nanoseconds in 1 millisecond`, "1.0").Plus(MakeKeyword("tag"), String{S: "Int"}))

	timeNamespace.InternVar("minute", minute_,
		MakeMeta(
			nil,
			`Number of nanoseconds in 1 minute`, "1.0").Plus(MakeKeyword("tag"), String{S: "BigInt"}))

	timeNamespace.InternVar("nanosecond", nanosecond_,
		MakeMeta(
			nil,
			`Number of nanoseconds in 1 nanosecond`, "1.0").Plus(MakeKeyword("tag"), String{S: "Int"}))

	timeNamespace.InternVar("rfc1123", rfc1123_,
		MakeMeta(
			nil,
			`Mon, 02 Jan 2006 15:04:05 MST`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("rfc1123-z", rfc1123_z_,
		MakeMeta(
			nil,
			`Mon, 02 Jan 2006 15:04:05 -0700`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("rfc3339", rfc3339_,
		MakeMeta(
			nil,
			`2006-01-02T15:04:05Z07:00`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("rfc3339-nano", rfc3339_nano_,
		MakeMeta(
			nil,
			`2006-01-02T15:04:05.999999999Z07:00`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("rfc822", rfc822_,
		MakeMeta(
			nil,
			`02 Jan 06 15:04 MST`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("rfc822-z", rfc822_z_,
		MakeMeta(
			nil,
			`02 Jan 06 15:04 -0700`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("rfc850", rfc850_,
		MakeMeta(
			nil,
			`Monday, 02-Jan-06 15:04:05 MST`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("ruby-date", ruby_date_,
		MakeMeta(
			nil,
			`Mon Jan 02 15:04:05 -0700 2006`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("second", second_,
		MakeMeta(
			nil,
			`Number of nanoseconds in 1 second`, "1.0").Plus(MakeKeyword("tag"), String{S: "Int"}))

	timeNamespace.InternVar("stamp", stamp_,
		MakeMeta(
			nil,
			`Jan _2 15:04:05`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("stamp-micro", stamp_micro_,
		MakeMeta(
			nil,
			`Jan _2 15:04:05.000000`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("stamp-milli", stamp_milli_,
		MakeMeta(
			nil,
			`Jan _2 15:04:05.000`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("stamp-nano", stamp_nano_,
		MakeMeta(
			nil,
			`Jan _2 15:04:05.000000000`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("unix-date", unix_date_,
		MakeMeta(
			nil,
			`Mon Jan _2 15:04:05 MST 2006`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("add", add_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("t"), MakeSymbol("d"))),
			`Returns the time t+d.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Time"}))

	timeNamespace.InternVar("add-date", add_date_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("t"), MakeSymbol("years"), MakeSymbol("months"), MakeSymbol("days"))),
			`Returns the time t + (years, months, days).`, "1.0").Plus(MakeKeyword("tag"), String{S: "Time"}))

	timeNamespace.InternVar("format", format_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("t"), MakeSymbol("layout"))),
			`Returns a textual representation of the time value formatted according to layout,
  which defines the format by showing how the reference time, defined to be
  Mon Jan 2 15:04:05 -0700 MST 2006
  would be displayed if it were the value; it serves as an example of the desired output.
  The same display rules will then be applied to the time value..`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("from-unix", from_unix_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("sec"), MakeSymbol("nsec"))),
			`Returns the local Time corresponding to the given Unix time, sec seconds and
  nsec nanoseconds since January 1, 1970 UTC. It is valid to pass nsec outside the range [0, 999999999].`, "1.0").Plus(MakeKeyword("tag"), String{S: "Time"}))

	timeNamespace.InternVar("hours", hours_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("d"))),
			`Returns the duration (passed as a number of nanoseconds) as a floating point number of hours.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	timeNamespace.InternVar("minutes", minutes_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("d"))),
			`Returns the duration (passed as a number of nanoseconds) as a floating point number of minutes.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	timeNamespace.InternVar("now", now_,
		MakeMeta(
			NewListFrom(NewVectorFrom()),
			`Returns the current local time.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Time"}))

	timeNamespace.InternVar("parse", parse_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("layout"), MakeSymbol("value"))),
			`Parses a time string.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Time"}))

	timeNamespace.InternVar("parse-duration", parse_duration_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("s"))),
			`Parses a duration string. A duration string is a possibly signed sequence of decimal numbers,
  each with optional fraction and a unit suffix, such as 300ms, -1.5h or 2h45m. Valid time units are
  ns, us (or µs), ms, s, m, h.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Int"}))

	timeNamespace.InternVar("round", round_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("d"), MakeSymbol("m"))),
			`Returns the result of rounding d to the nearest multiple of m. d and m represent time durations in nanoseconds.
  The rounding behavior for halfway values is to round away from zero. If m <= 0, returns d unchanged.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Int"}))

	timeNamespace.InternVar("seconds", seconds_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("d"))),
			`Returns the duration (passed as a number of nanoseconds) as a floating point number of seconds.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	timeNamespace.InternVar("since", since_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("t"))),
			`Returns the time in nanoseconds elapsed since t.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Int"}))

	timeNamespace.InternVar("sleep", sleep_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("d"))),
			`Pauses the execution thread for at least the duration d (expressed in nanoseconds).
  A negative or zero duration causes sleep to return immediately.`, "1.0"))

	timeNamespace.InternVar("string", string_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("d"))),
			`Returns a string representing the duration in the form 72h3m0.5s.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	timeNamespace.InternVar("sub", sub_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("t"), MakeSymbol("u"))),
			`Returns the duration t-u in nanoseconds.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Int"}))

	timeNamespace.InternVar("truncate", truncate_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("d"), MakeSymbol("m"))),
			`Returns the result of rounding d toward zero to a multiple of m. If m <= 0, returns d unchanged.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Int"}))

	timeNamespace.InternVar("unix", unix_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("t"))),
			`Returns t as a Unix time, the number of seconds elapsed since January 1, 1970 UTC.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Int"}))

	timeNamespace.InternVar("until", until_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("t"))),
			`Returns the duration in nanoseconds until t.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Int"}))

}
