(ns
  ^{:go-imports ["time"]
    :doc "Provides functionality for measuring and displaying time."}
  time)

(defn sleep
  "Pauses the execution thread for at least the duration d (expressed in nanoseconds).
  A negative or zero duration causes sleep to return immediately."
  {:added "1.0"
   :go "! time.Sleep(time.Duration(d)); _res, err := NIL, nil"}
  [^Int d])

(defn ^Time now
  "Returns the current local time."
  {:added "1.0"
   :go "time.Now(), nil"}
  [])

(defn ^Time from-unix
  "Returns the local Time corresponding to the given Unix time, sec seconds and
  nsec nanoseconds since January 1, 1970 UTC. It is valid to pass nsec outside the range [0, 999999999]."
  {:added "1.0"
   :go "time.Unix(int64(sec), int64(nsec)), nil"}
  [^Int sec ^Int nsec])

(defn ^Int unix
  "Returns t as a Unix time, the number of seconds elapsed since January 1, 1970 UTC."
  {:added "1.0"
   :go "int(t.Unix()), nil"}
  [^Time t])

(defn ^Int sub
  "Returns the duration t-u in nanoseconds."
  {:added "1.0"
  ;; TODO: 32-bit issue
   :go "int(t.Sub(u)), nil"}
  [^Time t ^Time u])

(defn ^Time add
  "Returns the time t+d."
  {:added "1.0"
   :go "t.Add(time.Duration(d)), nil"}
  [^Time t ^Int d])

(defn ^Time add-date
  "Returns the time t + (years, months, days)."
  {:added "1.0"
   :go "t.AddDate(years, months, days), nil"}
  [^Time t ^Int years ^Int months ^Int days])

(defn ^Time parse
  "Parses a time string."
  {:added "1.0"
   :go "time.Parse(layout, value)"}
  [^String layout ^String value])

(defn ^Int parse-duration
  "Parses a duration string. A duration string is a possibly signed sequence of decimal numbers,
  each with optional fraction and a unit suffix, such as 300ms, -1.5h or 2h45m. Valid time units are
  ns, us (or µs), ms, s, m, h."
  {:added "1.0"
  ;; TODO: 32-bit issue
   :go "! t, err := time.ParseDuration(s); _res := int(t);"}
  [^String s])

(defn ^Int since
  "Returns the time in nanoseconds elapsed since t."
  {:added "1.0"
  ;; TODO: 32-bit issue
   :go "int(time.Since(t)), nil"}
  [^Time t])

(defn ^Int until
  "Returns the duration in nanoseconds until t."
  {:added "1.0"
  ;; TODO: 32-bit issue
   :go "int(time.Until(t)), nil"}
  [^Time t])

(defn ^String format
  "Returns a textual representation of the time value formatted according to layout,
  which defines the format by showing how the reference time, defined to be
  Mon Jan 2 15:04:05 -0700 MST 2006
  would be displayed if it were the value; it serves as an example of the desired output.
  The same display rules will then be applied to the time value.."
  {:added "1.0"
   :go "t.Format(layout), nil"}
  [^Time t ^String layout])

(defn ^Double hours
  "Returns the duration (passed as a number of nanoseconds) as a floating point number of hours."
  {:added "1.0"
   :go "time.Duration(d).Hours(), nil"}
  [^Int d])

(defn ^Double minutes
  "Returns the duration (passed as a number of nanoseconds) as a floating point number of minutes."
  {:added "1.0"
   :go "time.Duration(d).Minutes(), nil"}
  [^Int d])

(defn ^Int round
  "Returns the result of rounding d to the nearest multiple of m. d and m represent time durations in nanoseconds.
  The rounding behavior for halfway values is to round away from zero. If m <= 0, returns d unchanged."
  {:added "1.0"
  ;; TODO: 32-bit issue
   :go "int(time.Duration(d).Round(time.Duration(m))), nil"}
  [^Int d ^Int m])

(defn ^Double seconds
  "Returns the duration (passed as a number of nanoseconds) as a floating point number of seconds."
  {:added "1.0"
   :go "time.Duration(d).Seconds(), nil"}
  [^Int d])

(defn ^String string
  "Returns a string representing the duration in the form 72h3m0.5s."
  {:added "1.0"
   :go "time.Duration(d).String(), nil"}
  [^Int d])

(defn ^Int truncate
  "Returns the result of rounding d toward zero to a multiple of m. If m <= 0, returns d unchanged."
  {:added "1.0"
  ;; TODO: 32-bit issue
   :go "int(time.Duration(d).Truncate(time.Duration(m))), nil"}
  [^Int d ^Int m])

(def
  ^{:doc "Number of nanoseconds in 1 nanosecond"
    :added "1.0"
    :tag Int
    :go "int(time.Nanosecond)"}
  nanosecond)

(def
  ^{:doc "Number of nanoseconds in 1 microsecond"
    :added "1.0"
    :tag Int
    :go "int(time.Microsecond)"}
  microsecond)

(def
  ^{:doc "Number of nanoseconds in 1 millisecond"
    :added "1.0"
    :tag Int
    :go "int(time.Millisecond)"}
  millisecond)

(def
  ^{:doc "Number of nanoseconds in 1 second"
    :added "1.0"
    :tag Int
    :go "int(time.Second)"}
  second)

(def
  ^{:doc "Number of nanoseconds in 1 minute"
    :added "1.0"
    :tag BigInt
    :go "int64(time.Minute)"}
  minute)

(def
  ^{:doc "Number of nanoseconds in 1 hour"
    :added "1.0"
    :tag BigInt
    :go "int64(time.Hour)"}
  hour)

(def
  ^{:doc "Mon Jan _2 15:04:05 2006"
    :added "1.0"
    :tag String
    :go "time.ANSIC"}
  ansi-c)

(def
  ^{:doc "Mon Jan _2 15:04:05 MST 2006"
    :added "1.0"
    :tag String
    :go "time.UnixDate"}
  unix-date)

(def
  ^{:doc "Mon Jan 02 15:04:05 -0700 2006"
    :added "1.0"
    :tag String
    :go "time.RubyDate"}
  ruby-date)

(def
  ^{:doc "02 Jan 06 15:04 MST"
    :added "1.0"
    :tag String
    :go "time.RFC822"}
  rfc822)

(def
  ^{:doc "02 Jan 06 15:04 -0700"
    :added "1.0"
    :tag String
    :go "time.RFC822Z"}
  rfc822-z)

(def
  ^{:doc "Monday, 02-Jan-06 15:04:05 MST"
    :added "1.0"
    :tag String
    :go "time.RFC850"}
  rfc850)

(def
  ^{:doc "Mon, 02 Jan 2006 15:04:05 MST"
    :added "1.0"
    :tag String
    :go "time.RFC1123"}
  rfc1123)

(def
  ^{:doc "Mon, 02 Jan 2006 15:04:05 -0700"
    :added "1.0"
    :tag String
    :go "time.RFC1123Z"}
  rfc1123-z)

(def
  ^{:doc "2006-01-02T15:04:05Z07:00"
    :added "1.0"
    :tag String
    :go "time.RFC3339"}
  rfc3339)

(def
  ^{:doc "2006-01-02T15:04:05.999999999Z07:00"
    :added "1.0"
    :tag String
    :go "time.RFC3339Nano"}
  rfc3339-nano)

(def
  ^{:doc "3:04PM"
    :added "1.0"
    :tag String
    :go "time.Kitchen"}
  kitchen)

(def
  ^{:doc "Jan _2 15:04:05"
    :added "1.0"
    :tag String
    :go "time.Stamp"}
  stamp)

(def
  ^{:doc "Jan _2 15:04:05.000"
    :added "1.0"
    :tag String
    :go "time.StampMilli"}
  stamp-milli)

(def
  ^{:doc "Jan _2 15:04:05.000000"
    :added "1.0"
    :tag String
    :go "time.StampMicro"}
  stamp-micro)

(def
  ^{:doc "Jan _2 15:04:05.000000000"
    :added "1.0"
    :tag String
    :go "time.StampNano"}
  stamp-nano)
