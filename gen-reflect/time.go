package reflect

import (
	"reflect"
	time "time"

	"github.com/lab47/lace/pkg/pkgreflect"
)

func init() {
	ParseError_methods := map[string]pkgreflect.Func{}
	Timer_methods := map[string]pkgreflect.Func{}
	Ticker_methods := map[string]pkgreflect.Func{}
	Duration_methods := map[string]pkgreflect.Func{}
	Month_methods := map[string]pkgreflect.Func{}
	Time_methods := map[string]pkgreflect.Func{}
	Weekday_methods := map[string]pkgreflect.Func{}
	Location_methods := map[string]pkgreflect.Func{}
	Time_methods["String"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: "String returns the time formatted using the format string\n\n\t\"2006-01-02 15:04:05.999999999 -0700 MST\"\n\nIf the time has a monotonic clock reading, the returned string\nincludes a final field \"m=±<value>\", where value is the monotonic\nclock reading formatted as a decimal number of seconds.\n\nThe returned string is meant for debugging; for a stable serialized\nrepresentation, use t.MarshalText, t.MarshalBinary, or t.Format\nwith an explicit format string."}
	Time_methods["GoString"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: "GoString implements fmt.GoStringer and formats t to be printed in Go source\ncode."}
	Time_methods["Format"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "layout", Tag: "string"}}, Tag: "string", Doc: "Format returns a textual representation of the time value formatted according\nto the layout defined by the argument. See the documentation for the\nconstant called Layout to see how to represent the layout format.\n\nThe executable example for Time.Format demonstrates the working\nof the layout string in detail and is a good reference."}
	Time_methods["AppendFormat"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}, {Name: "layout", Tag: "string"}}, Tag: "[]byte", Doc: "AppendFormat is like Format but appends the textual\nrepresentation to b and returns the extended buffer."}
	ParseError_methods["Error"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: "Error returns the string representation of a ParseError."}
	Timer_methods["Stop"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "bool", Doc: "Stop prevents the Timer from firing.\nIt returns true if the call stops the timer, false if the timer has already\nexpired or been stopped.\nStop does not close the channel, to prevent a read from the channel succeeding\nincorrectly.\n\nTo ensure the channel is empty after a call to Stop, check the\nreturn value and drain the channel.\nFor example, assuming the program has not received from t.C already:\n\n\tif !t.Stop() {\n\t\t<-t.C\n\t}\n\nThis cannot be done concurrent to other receives from the Timer's\nchannel or other calls to the Timer's Stop method.\n\nFor a timer created with AfterFunc(d, f), if t.Stop returns false, then the timer\nhas already expired and the function f has been started in its own goroutine;\nStop does not wait for f to complete before returning.\nIf the caller needs to know whether f is completed, it must coordinate\nwith f explicitly."}
	Timer_methods["Reset"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "d", Tag: "Duration"}}, Tag: "bool", Doc: "Reset changes the timer to expire after duration d.\nIt returns true if the timer had been active, false if the timer had\nexpired or been stopped.\n\nFor a Timer created with NewTimer, Reset should be invoked only on\nstopped or expired timers with drained channels.\n\nIf a program has already received a value from t.C, the timer is known\nto have expired and the channel drained, so t.Reset can be used directly.\nIf a program has not yet received a value from t.C, however,\nthe timer must be stopped and—if Stop reports that the timer expired\nbefore being stopped—the channel explicitly drained:\n\n\tif !t.Stop() {\n\t\t<-t.C\n\t}\n\tt.Reset(d)\n\nThis should not be done concurrent to other receives from the Timer's\nchannel.\n\nNote that it is not possible to use Reset's return value correctly, as there\nis a race condition between draining the channel and the new timer expiring.\nReset should always be invoked on stopped or expired channels, as described above.\nThe return value exists to preserve compatibility with existing programs.\n\nFor a Timer created with AfterFunc(d, f), Reset either reschedules\nwhen f will run, in which case Reset returns true, or schedules f\nto run again, in which case it returns false.\nWhen Reset returns false, Reset neither waits for the prior f to\ncomplete before returning nor does it guarantee that the subsequent\ngoroutine running f does not run concurrently with the prior\none. If the caller needs to know whether the prior execution of\nf is completed, it must coordinate with f explicitly."}
	Ticker_methods["Stop"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "Stop turns off a ticker. After Stop, no more ticks will be sent.\nStop does not close the channel, to prevent a concurrent goroutine\nreading from the channel from seeing an erroneous \"tick\"."}
	Ticker_methods["Reset"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "d", Tag: "Duration"}}, Tag: "any", Doc: "Reset stops a ticker and resets its period to the specified duration.\nThe next tick will arrive after the new period elapses. The duration d\nmust be greater than zero; if not, Reset will panic."}
	Time_methods["After"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "u", Tag: "Time"}}, Tag: "bool", Doc: "After reports whether the time instant t is after u."}
	Time_methods["Before"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "u", Tag: "Time"}}, Tag: "bool", Doc: "Before reports whether the time instant t is before u."}
	Time_methods["Compare"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "u", Tag: "Time"}}, Tag: "int", Doc: "Compare compares the time instant t with u. If t is before u, it returns -1;\nif t is after u, it returns +1; if they're the same, it returns 0."}
	Time_methods["Equal"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "u", Tag: "Time"}}, Tag: "bool", Doc: "Equal reports whether t and u represent the same time instant.\nTwo times can be equal even if they are in different locations.\nFor example, 6:00 +0200 and 4:00 UTC are Equal.\nSee the documentation on the Time type for the pitfalls of using == with\nTime values; most code should use Equal instead."}
	Month_methods["String"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: "String returns the English name of the month (\"January\", \"February\", ...)."}
	Weekday_methods["String"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: "String returns the English name of the day (\"Sunday\", \"Monday\", ...)."}
	Time_methods["IsZero"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "bool", Doc: "IsZero reports whether t represents the zero time instant,\nJanuary 1, year 1, 00:00:00 UTC."}
	Time_methods["Date"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "Date returns the year, month, and day in which t occurs."}
	Time_methods["Year"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int", Doc: "Year returns the year in which t occurs."}
	Time_methods["Month"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "Month", Doc: "Month returns the month of the year specified by t."}
	Time_methods["Day"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int", Doc: "Day returns the day of the month specified by t."}
	Time_methods["Weekday"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "Weekday", Doc: "Weekday returns the day of the week specified by t."}
	Time_methods["ISOWeek"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int", Doc: "ISOWeek returns the ISO 8601 year and week number in which t occurs.\nWeek ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to\nweek 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1\nof year n+1."}
	Time_methods["Clock"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int", Doc: "Clock returns the hour, minute, and second within the day specified by t."}
	Time_methods["Hour"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int", Doc: "Hour returns the hour within the day specified by t, in the range [0, 23]."}
	Time_methods["Minute"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int", Doc: "Minute returns the minute offset within the hour specified by t, in the range [0, 59]."}
	Time_methods["Second"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int", Doc: "Second returns the second offset within the minute specified by t, in the range [0, 59]."}
	Time_methods["Nanosecond"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int", Doc: "Nanosecond returns the nanosecond offset within the second specified by t,\nin the range [0, 999999999]."}
	Time_methods["YearDay"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int", Doc: "YearDay returns the day of the year specified by t, in the range [1,365] for non-leap years,\nand [1,366] in leap years."}
	Duration_methods["String"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: "String returns a string representing the duration in the form \"72h3m0.5s\".\nLeading zero units are omitted. As a special case, durations less than one\nsecond format use a smaller unit (milli-, micro-, or nanoseconds) to ensure\nthat the leading digit is non-zero. The zero duration formats as 0s."}
	Duration_methods["Nanoseconds"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int64", Doc: "Nanoseconds returns the duration as an integer nanosecond count."}
	Duration_methods["Microseconds"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int64", Doc: "Microseconds returns the duration as an integer microsecond count."}
	Duration_methods["Milliseconds"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int64", Doc: "Milliseconds returns the duration as an integer millisecond count."}
	Duration_methods["Seconds"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "float64", Doc: "Seconds returns the duration as a floating point number of seconds."}
	Duration_methods["Minutes"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "float64", Doc: "Minutes returns the duration as a floating point number of minutes."}
	Duration_methods["Hours"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "float64", Doc: "Hours returns the duration as a floating point number of hours."}
	Duration_methods["Truncate"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "m", Tag: "Duration"}}, Tag: "Duration", Doc: "Truncate returns the result of rounding d toward zero to a multiple of m.\nIf m <= 0, Truncate returns d unchanged."}
	Duration_methods["Round"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "m", Tag: "Duration"}}, Tag: "Duration", Doc: "Round returns the result of rounding d to the nearest multiple of m.\nThe rounding behavior for halfway values is to round away from zero.\nIf the result exceeds the maximum (or minimum)\nvalue that can be stored in a Duration,\nRound returns the maximum (or minimum) duration.\nIf m <= 0, Round returns d unchanged."}
	Duration_methods["Abs"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "Duration", Doc: "Abs returns the absolute value of d.\nAs a special case, math.MinInt64 is converted to math.MaxInt64."}
	Time_methods["Add"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "d", Tag: "Duration"}}, Tag: "Time", Doc: "Add returns the time t+d."}
	Time_methods["Sub"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "u", Tag: "Time"}}, Tag: "Duration", Doc: "Sub returns the duration t-u. If the result exceeds the maximum (or minimum)\nvalue that can be stored in a Duration, the maximum (or minimum) duration\nwill be returned.\nTo compute t-d for a duration d, use t.Add(-d)."}
	Time_methods["AddDate"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "years", Tag: "int"}, {Name: "months", Tag: "int"}, {Name: "days", Tag: "int"}}, Tag: "Time", Doc: "AddDate returns the time corresponding to adding the\ngiven number of years, months, and days to t.\nFor example, AddDate(-1, 2, 3) applied to January 1, 2011\nreturns March 4, 2010.\n\nNote that dates are fundamentally coupled to timezones, and calendrical\nperiods like days don't have fixed durations. AddDate uses the Location of\nthe Time value to determine these durations. That means that the same\nAddDate arguments can produce a different shift in absolute time depending on\nthe base Time value and its Location. For example, AddDate(0, 0, 1) applied\nto 12:00 on March 27 always returns 12:00 on March 28. At some locations and\nin some years this is a 24 hour shift. In others it's a 23 hour shift due to\ndaylight savings time transitions.\n\nAddDate normalizes its result in the same way that Date does,\nso, for example, adding one month to October 31 yields\nDecember 1, the normalized form for November 31."}
	Time_methods["UTC"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "Time", Doc: "UTC returns t with the location set to UTC."}
	Time_methods["Local"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "Time", Doc: "Local returns t with the location set to local time."}
	Time_methods["In"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "loc", Tag: "Location"}}, Tag: "Time", Doc: "In returns a copy of t representing the same time instant, but\nwith the copy's location information set to loc for display\npurposes.\n\nIn panics if loc is nil."}
	Time_methods["Location"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "Location", Doc: "Location returns the time zone information associated with t."}
	Time_methods["Zone"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "Zone computes the time zone in effect at time t, returning the abbreviated\nname of the zone (such as \"CET\") and its offset in seconds east of UTC."}
	Time_methods["ZoneBounds"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "Time", Doc: "ZoneBounds returns the bounds of the time zone in effect at time t.\nThe zone begins at start and the next zone begins at end.\nIf the zone begins at the beginning of time, start will be returned as a zero Time.\nIf the zone goes on forever, end will be returned as a zero Time.\nThe Location of the returned times will be the same as t."}
	Time_methods["Unix"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int64", Doc: "Unix returns t as a Unix time, the number of seconds elapsed\nsince January 1, 1970 UTC. The result does not depend on the\nlocation associated with t.\nUnix-like operating systems often record time as a 32-bit\ncount of seconds, but since the method here returns a 64-bit\nvalue it is valid for billions of years into the past or future."}
	Time_methods["UnixMilli"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int64", Doc: "UnixMilli returns t as a Unix time, the number of milliseconds elapsed since\nJanuary 1, 1970 UTC. The result is undefined if the Unix time in\nmilliseconds cannot be represented by an int64 (a date more than 292 million\nyears before or after 1970). The result does not depend on the\nlocation associated with t."}
	Time_methods["UnixMicro"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int64", Doc: "UnixMicro returns t as a Unix time, the number of microseconds elapsed since\nJanuary 1, 1970 UTC. The result is undefined if the Unix time in\nmicroseconds cannot be represented by an int64 (a date before year -290307 or\nafter year 294246). The result does not depend on the location associated\nwith t."}
	Time_methods["UnixNano"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int64", Doc: "UnixNano returns t as a Unix time, the number of nanoseconds elapsed\nsince January 1, 1970 UTC. The result is undefined if the Unix time\nin nanoseconds cannot be represented by an int64 (a date before the year\n1678 or after 2262). Note that this means the result of calling UnixNano\non the zero Time is undefined. The result does not depend on the\nlocation associated with t."}
	Time_methods["MarshalBinary"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "MarshalBinary implements the encoding.BinaryMarshaler interface."}
	Time_methods["UnmarshalBinary"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "data", Tag: "[]byte"}}, Tag: "error", Doc: "UnmarshalBinary implements the encoding.BinaryUnmarshaler interface."}
	Time_methods["GobEncode"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "GobEncode implements the gob.GobEncoder interface."}
	Time_methods["GobDecode"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "data", Tag: "[]byte"}}, Tag: "error", Doc: "GobDecode implements the gob.GobDecoder interface."}
	Time_methods["MarshalJSON"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "MarshalJSON implements the json.Marshaler interface.\nThe time is a quoted string in the RFC 3339 format with sub-second precision.\nIf the timestamp cannot be represented as valid RFC 3339\n(e.g., the year is out of range), then an error is reported."}
	Time_methods["UnmarshalJSON"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "data", Tag: "[]byte"}}, Tag: "error", Doc: "UnmarshalJSON implements the json.Unmarshaler interface.\nThe time must be a quoted string in the RFC 3339 format."}
	Time_methods["MarshalText"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "MarshalText implements the encoding.TextMarshaler interface.\nThe time is formatted in RFC 3339 format with sub-second precision.\nIf the timestamp cannot be represented as valid RFC 3339\n(e.g., the year is out of range), then an error is reported."}
	Time_methods["UnmarshalText"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "data", Tag: "[]byte"}}, Tag: "error", Doc: "UnmarshalText implements the encoding.TextUnmarshaler interface.\nThe time must be in the RFC 3339 format."}
	Time_methods["IsDST"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "bool", Doc: "IsDST reports whether the time in the configured location is in Daylight Savings Time."}
	Time_methods["Truncate"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "d", Tag: "Duration"}}, Tag: "Time", Doc: "Truncate returns the result of rounding t down to a multiple of d (since the zero time).\nIf d <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.\n\nTruncate operates on the time as an absolute duration since the\nzero time; it does not operate on the presentation form of the\ntime. Thus, Truncate(Hour) may return a time with a non-zero\nminute, depending on the time's Location."}
	Time_methods["Round"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "d", Tag: "Duration"}}, Tag: "Time", Doc: "Round returns the result of rounding t to the nearest multiple of d (since the zero time).\nThe rounding behavior for halfway values is to round up.\nIf d <= 0, Round returns t stripped of any monotonic clock reading but otherwise unchanged.\n\nRound operates on the time as an absolute duration since the\nzero time; it does not operate on the presentation form of the\ntime. Thus, Round(Hour) may return a time with a non-zero\nminute, depending on the time's Location."}
	Location_methods["String"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: "String returns a descriptive name for the time zone information,\ncorresponding to the name argument to LoadLocation or FixedZone."}
	pkgreflect.AddPackage("time", &pkgreflect.Package{
		Doc: "Package time provides functionality for measuring and displaying time.",
		Types: map[string]pkgreflect.Type{
			"Duration":   {Doc: "", Value: reflect.TypeOf((*time.Duration)(nil)).Elem(), Methods: Duration_methods},
			"Location":   {Doc: "", Value: reflect.TypeOf((*time.Location)(nil)).Elem(), Methods: Location_methods},
			"Month":      {Doc: "", Value: reflect.TypeOf((*time.Month)(nil)).Elem(), Methods: Month_methods},
			"ParseError": {Doc: "", Value: reflect.TypeOf((*time.ParseError)(nil)).Elem(), Methods: ParseError_methods},
			"Ticker":     {Doc: "", Value: reflect.TypeOf((*time.Ticker)(nil)).Elem(), Methods: Ticker_methods},
			"Time":       {Doc: "", Value: reflect.TypeOf((*time.Time)(nil)).Elem(), Methods: Time_methods},
			"Timer":      {Doc: "", Value: reflect.TypeOf((*time.Timer)(nil)).Elem(), Methods: Timer_methods},
			"Weekday":    {Doc: "", Value: reflect.TypeOf((*time.Weekday)(nil)).Elem(), Methods: Weekday_methods},
		},

		Functions: map[string]pkgreflect.FuncValue{
			"After": {Doc: "After waits for the duration to elapse and then sends the current time\non the returned channel.\nIt is equivalent to NewTimer(d).C.\nThe underlying Timer is not recovered by the garbage collector\nuntil the timer fires. If efficiency is a concern, use NewTimer\ninstead and call Timer.Stop if the timer is no longer needed.", Args: []pkgreflect.Arg{{Name: "d", Tag: "Duration"}}, Tag: "Unknown", Value: reflect.ValueOf(time.After)},

			"AfterFunc": {Doc: "AfterFunc waits for the duration to elapse and then calls f\nin its own goroutine. It returns a Timer that can\nbe used to cancel the call using its Stop method.\nThe returned Timer's C field is not used and will be nil.", Args: []pkgreflect.Arg{{Name: "d", Tag: "Duration"}, {Name: "f", Tag: "Unknown"}}, Tag: "Timer", Value: reflect.ValueOf(time.AfterFunc)},

			"Date": {Doc: "Date returns the Time corresponding to\n\n\tyyyy-mm-dd hh:mm:ss + nsec nanoseconds\n\nin the appropriate zone for that time in the given location.\n\nThe month, day, hour, min, sec, and nsec values may be outside\ntheir usual ranges and will be normalized during the conversion.\nFor example, October 32 converts to November 1.\n\nA daylight savings time transition skips or repeats times.\nFor example, in the United States, March 13, 2011 2:15am never occurred,\nwhile November 6, 2011 1:15am occurred twice. In such cases, the\nchoice of time zone, and therefore the time, is not well-defined.\nDate returns a time that is correct in one of the two zones involved\nin the transition, but it does not guarantee which.\n\nDate panics if loc is nil.", Args: []pkgreflect.Arg{{Name: "year", Tag: "int"}, {Name: "month", Tag: "Month"}, {Name: "day", Tag: "int"}, {Name: "hour", Tag: "int"}, {Name: "min", Tag: "int"}, {Name: "sec", Tag: "int"}, {Name: "nsec", Tag: "int"}, {Name: "loc", Tag: "Location"}}, Tag: "Time", Value: reflect.ValueOf(time.Date)},

			"FixedZone": {Doc: "FixedZone returns a Location that always uses\nthe given zone name and offset (seconds east of UTC).", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}, {Name: "offset", Tag: "int"}}, Tag: "Location", Value: reflect.ValueOf(time.FixedZone)},

			"LoadLocation": {Doc: "LoadLocation returns the Location with the given name.\n\nIf the name is \"\" or \"UTC\", LoadLocation returns UTC.\nIf the name is \"Local\", LoadLocation returns Local.\n\nOtherwise, the name is taken to be a location name corresponding to a file\nin the IANA Time Zone database, such as \"America/New_York\".\n\nLoadLocation looks for the IANA Time Zone database in the following\nlocations in order:\n\n  - the directory or uncompressed zip file named by the ZONEINFO environment variable\n  - on a Unix system, the system standard installation location\n  - $GOROOT/lib/time/zoneinfo.zip\n  - the time/tzdata package, if it was imported", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}}, Tag: "any", Value: reflect.ValueOf(time.LoadLocation)},

			"LoadLocationFromTZData": {Doc: "LoadLocationFromTZData returns a Location with the given name\ninitialized from the IANA Time Zone database-formatted data.\nThe data should be in the format of a standard IANA time zone file\n(for example, the content of /etc/localtime on Unix systems).", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}, {Name: "data", Tag: "[]byte"}}, Tag: "any", Value: reflect.ValueOf(time.LoadLocationFromTZData)},

			"NewTicker": {Doc: "NewTicker returns a new Ticker containing a channel that will send\nthe current time on the channel after each tick. The period of the\nticks is specified by the duration argument. The ticker will adjust\nthe time interval or drop ticks to make up for slow receivers.\nThe duration d must be greater than zero; if not, NewTicker will\npanic. Stop the ticker to release associated resources.", Args: []pkgreflect.Arg{{Name: "d", Tag: "Duration"}}, Tag: "Ticker", Value: reflect.ValueOf(time.NewTicker)},

			"NewTimer": {Doc: "NewTimer creates a new Timer that will send\nthe current time on its channel after at least duration d.", Args: []pkgreflect.Arg{{Name: "d", Tag: "Duration"}}, Tag: "Timer", Value: reflect.ValueOf(time.NewTimer)},

			"Now": {Doc: "Now returns the current local time.", Args: []pkgreflect.Arg{}, Tag: "Time", Value: reflect.ValueOf(time.Now)},

			"Parse": {Doc: "Parse parses a formatted string and returns the time value it represents.\nSee the documentation for the constant called Layout to see how to\nrepresent the format. The second argument must be parseable using\nthe format string (layout) provided as the first argument.\n\nThe example for Time.Format demonstrates the working of the layout string\nin detail and is a good reference.\n\nWhen parsing (only), the input may contain a fractional second\nfield immediately after the seconds field, even if the layout does not\nsignify its presence. In that case either a comma or a decimal point\nfollowed by a maximal series of digits is parsed as a fractional second.\nFractional seconds are truncated to nanosecond precision.\n\nElements omitted from the layout are assumed to be zero or, when\nzero is impossible, one, so parsing \"3:04pm\" returns the time\ncorresponding to Jan 1, year 0, 15:04:00 UTC (note that because the year is\n0, this time is before the zero Time).\nYears must be in the range 0000..9999. The day of the week is checked\nfor syntax but it is otherwise ignored.\n\nFor layouts specifying the two-digit year 06, a value NN >= 69 will be treated\nas 19NN and a value NN < 69 will be treated as 20NN.\n\nThe remainder of this comment describes the handling of time zones.\n\nIn the absence of a time zone indicator, Parse returns a time in UTC.\n\nWhen parsing a time with a zone offset like -0700, if the offset corresponds\nto a time zone used by the current location (Local), then Parse uses that\nlocation and zone in the returned time. Otherwise it records the time as\nbeing in a fabricated location with time fixed at the given zone offset.\n\nWhen parsing a time with a zone abbreviation like MST, if the zone abbreviation\nhas a defined offset in the current location, then that offset is used.\nThe zone abbreviation \"UTC\" is recognized as UTC regardless of location.\nIf the zone abbreviation is unknown, Parse records the time as being\nin a fabricated location with the given zone abbreviation and a zero offset.\nThis choice means that such a time can be parsed and reformatted with the\nsame layout losslessly, but the exact instant used in the representation will\ndiffer by the actual zone offset. To avoid such problems, prefer time layouts\nthat use a numeric zone offset, or use ParseInLocation.", Args: []pkgreflect.Arg{{Name: "layout", Tag: "string"}, {Name: "value", Tag: "string"}}, Tag: "any", Value: reflect.ValueOf(time.Parse)},

			"ParseDuration": {Doc: "ParseDuration parses a duration string.\nA duration string is a possibly signed sequence of\ndecimal numbers, each with optional fraction and a unit suffix,\nsuch as \"300ms\", \"-1.5h\" or \"2h45m\".\nValid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".", Args: []pkgreflect.Arg{{Name: "s", Tag: "string"}}, Tag: "any", Value: reflect.ValueOf(time.ParseDuration)},

			"ParseInLocation": {Doc: "ParseInLocation is like Parse but differs in two important ways.\nFirst, in the absence of time zone information, Parse interprets a time as UTC;\nParseInLocation interprets the time as in the given location.\nSecond, when given a zone offset or abbreviation, Parse tries to match it\nagainst the Local location; ParseInLocation uses the given location.", Args: []pkgreflect.Arg{{Name: "layout", Tag: "string"}, {Name: "value", Tag: "string"}, {Name: "loc", Tag: "Location"}}, Tag: "any", Value: reflect.ValueOf(time.ParseInLocation)},

			"Since": {Doc: "Since returns the time elapsed since t.\nIt is shorthand for time.Now().Sub(t).", Args: []pkgreflect.Arg{{Name: "t", Tag: "Time"}}, Tag: "Duration", Value: reflect.ValueOf(time.Since)},

			"Sleep": {Doc: "Sleep pauses the current goroutine for at least the duration d.\nA negative or zero duration causes Sleep to return immediately.", Args: []pkgreflect.Arg{{Name: "d", Tag: "Duration"}}, Tag: "any", Value: reflect.ValueOf(time.Sleep)},

			"Tick": {Doc: "Tick is a convenience wrapper for NewTicker providing access to the ticking\nchannel only. While Tick is useful for clients that have no need to shut down\nthe Ticker, be aware that without a way to shut it down the underlying\nTicker cannot be recovered by the garbage collector; it \"leaks\".\nUnlike NewTicker, Tick will return nil if d <= 0.", Args: []pkgreflect.Arg{{Name: "d", Tag: "Duration"}}, Tag: "Unknown", Value: reflect.ValueOf(time.Tick)},

			"Unix": {Doc: "Unix returns the local Time corresponding to the given Unix time,\nsec seconds and nsec nanoseconds since January 1, 1970 UTC.\nIt is valid to pass nsec outside the range [0, 999999999].\nNot all sec values have a corresponding time value. One such\nvalue is 1<<63-1 (the largest int64 value).", Args: []pkgreflect.Arg{{Name: "sec", Tag: "int64"}, {Name: "nsec", Tag: "int64"}}, Tag: "Time", Value: reflect.ValueOf(time.Unix)},

			"UnixMicro": {Doc: "UnixMicro returns the local Time corresponding to the given Unix time,\nusec microseconds since January 1, 1970 UTC.", Args: []pkgreflect.Arg{{Name: "usec", Tag: "int64"}}, Tag: "Time", Value: reflect.ValueOf(time.UnixMicro)},

			"UnixMilli": {Doc: "UnixMilli returns the local Time corresponding to the given Unix time,\nmsec milliseconds since January 1, 1970 UTC.", Args: []pkgreflect.Arg{{Name: "msec", Tag: "int64"}}, Tag: "Time", Value: reflect.ValueOf(time.UnixMilli)},

			"Until": {Doc: "Until returns the duration until t.\nIt is shorthand for t.Sub(time.Now()).", Args: []pkgreflect.Arg{{Name: "t", Tag: "Time"}}, Tag: "Duration", Value: reflect.ValueOf(time.Until)},
		},

		Variables: map[string]pkgreflect.Value{
			"Local": {Doc: "", Value: reflect.ValueOf(&time.Local)},
			"UTC":   {Doc: "", Value: reflect.ValueOf(&time.UTC)},
		},

		Consts: map[string]pkgreflect.Value{
			"ANSIC":       {Doc: "", Value: reflect.ValueOf(time.ANSIC)},
			"April":       {Doc: "", Value: reflect.ValueOf(time.April)},
			"August":      {Doc: "", Value: reflect.ValueOf(time.August)},
			"DateOnly":    {Doc: "", Value: reflect.ValueOf(time.DateOnly)},
			"DateTime":    {Doc: "", Value: reflect.ValueOf(time.DateTime)},
			"December":    {Doc: "", Value: reflect.ValueOf(time.December)},
			"February":    {Doc: "", Value: reflect.ValueOf(time.February)},
			"Friday":      {Doc: "", Value: reflect.ValueOf(time.Friday)},
			"Hour":        {Doc: "", Value: reflect.ValueOf(time.Hour)},
			"January":     {Doc: "", Value: reflect.ValueOf(time.January)},
			"July":        {Doc: "", Value: reflect.ValueOf(time.July)},
			"June":        {Doc: "", Value: reflect.ValueOf(time.June)},
			"Kitchen":     {Doc: "", Value: reflect.ValueOf(time.Kitchen)},
			"Layout":      {Doc: "", Value: reflect.ValueOf(time.Layout)},
			"March":       {Doc: "", Value: reflect.ValueOf(time.March)},
			"May":         {Doc: "", Value: reflect.ValueOf(time.May)},
			"Microsecond": {Doc: "", Value: reflect.ValueOf(time.Microsecond)},
			"Millisecond": {Doc: "", Value: reflect.ValueOf(time.Millisecond)},
			"Minute":      {Doc: "", Value: reflect.ValueOf(time.Minute)},
			"Monday":      {Doc: "", Value: reflect.ValueOf(time.Monday)},
			"Nanosecond":  {Doc: "", Value: reflect.ValueOf(time.Nanosecond)},
			"November":    {Doc: "", Value: reflect.ValueOf(time.November)},
			"October":     {Doc: "", Value: reflect.ValueOf(time.October)},
			"RFC1123":     {Doc: "", Value: reflect.ValueOf(time.RFC1123)},
			"RFC1123Z":    {Doc: "", Value: reflect.ValueOf(time.RFC1123Z)},
			"RFC3339":     {Doc: "", Value: reflect.ValueOf(time.RFC3339)},
			"RFC3339Nano": {Doc: "", Value: reflect.ValueOf(time.RFC3339Nano)},
			"RFC822":      {Doc: "", Value: reflect.ValueOf(time.RFC822)},
			"RFC822Z":     {Doc: "", Value: reflect.ValueOf(time.RFC822Z)},
			"RFC850":      {Doc: "", Value: reflect.ValueOf(time.RFC850)},
			"RubyDate":    {Doc: "", Value: reflect.ValueOf(time.RubyDate)},
			"Saturday":    {Doc: "", Value: reflect.ValueOf(time.Saturday)},
			"Second":      {Doc: "", Value: reflect.ValueOf(time.Second)},
			"September":   {Doc: "", Value: reflect.ValueOf(time.September)},
			"Stamp":       {Doc: "Handy time stamps.", Value: reflect.ValueOf(time.Stamp)},
			"StampMicro":  {Doc: "", Value: reflect.ValueOf(time.StampMicro)},
			"StampMilli":  {Doc: "", Value: reflect.ValueOf(time.StampMilli)},
			"StampNano":   {Doc: "", Value: reflect.ValueOf(time.StampNano)},
			"Sunday":      {Doc: "", Value: reflect.ValueOf(time.Sunday)},
			"Thursday":    {Doc: "", Value: reflect.ValueOf(time.Thursday)},
			"TimeOnly":    {Doc: "", Value: reflect.ValueOf(time.TimeOnly)},
			"Tuesday":     {Doc: "", Value: reflect.ValueOf(time.Tuesday)},
			"UnixDate":    {Doc: "", Value: reflect.ValueOf(time.UnixDate)},
			"Wednesday":   {Doc: "", Value: reflect.ValueOf(time.Wednesday)},
		},
	})
}
