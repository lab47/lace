(ns lace.time
  (:require
    [time :as gt]
    [lace.reflect :as r]))

(refer 'time 'Time 'Duration)

(defn ^Time now 
  "Return a time value for the current moment." 
  [] (gt/Now))

(defn ^Time from-unix
  "Returns the local Time corresponding to the given Unix time, sec seconds and
  nsec nanoseconds since January 1, 1970 UTC. It is valid to pass nsec outside the range [0, 999999999]."
  [^Int sec ^Int nsec]
  (gt/Unix sec nsec))

(defn ^Int unix
  "Returns t as a Unix time, the number of seconds elapsed since January 1, 1970 UTC."
  [ts]
  (.unix ts))

(defn ^Int sub
  "Returns the duration t-u in nanoseconds."
  [^Time t ^Time u]
  (.sub t u))

(defn ^Int add
  "Returns a new time by adding the given duration to the given t."
  [^Time t ^Duration u]
  (.add t u))

(defn ^Time add-date
  "Returns the time t + (years, months, days)."
  {:added "1.0"}
  [^Time t ^Int years ^Int months ^Int days]
  (.add-date t years months days))

(defn ^Time parse
  "Parses a time string."
  {:added "1.0"}
  [^String layout ^String value]
  (gt/Parse layout value))
  
(defn ^Duration parse-duration
  "Parse the given string as a duration value"
  [str] (gt/ParseDuration str))

(defn ^Duration since 
  "Calculate the elapse time since t."
  [ts] (gt/Since ts))

(defn sleep
  "Sleep the current go routine for specified amount of time."
  [dur] 
  (if (string? dur) 
    (gt/Sleep (r/cast gt/Duration (parse-duration dur)))
    (gt/Sleep (r/cast gt/Duration dur))))

(defn ^Duration in-seconds
  "Convert a given values to a duration of seconds"
  [sec] (r/cast Duration (* gt/Second sec)))

(defn ^Duration in-milliseconds
  "Convert a given values to a duration of seconds"
  [sec] (r/cast Duration (* gt/Millisecond sec)))

(defn ^Duration in-microseconds
  "Convert a given values to a duration of seconds"
  [sec] (r/cast Duration (* gt/Microsecond sec)))

(defn ^Int until
  "Returns the duration in nanoseconds until t."
  {:added "1.0"}
  [^Time t]
  (gt/Until t))

(defn ^String format
  "Returns a textual representation of the time value formatted according to layout,
  which defines the format by showing how the reference time, defined to be
  Mon Jan 2 15:04:05 -0700 MST 2006
  would be displayed if it were the value; it serves as an example of the desired output.
  The same display rules will then be applied to the time value.."
  {:added "1.0"}
  [^Time t ^String layout]
  (.format t layout))

(defn ^Double hours
  "Returns the duration (passed as a number of nanoseconds) as a floating point number of hours."
  {:added "1.0"}
  [^Duration d]
  (.hours d)) 

(defn ^Double minutes
  "Returns the duration (passed as a number of nanoseconds) as a floating point number of minutes."
  {:added "1.0"}
  [^Int d]
  (.minutes d))

(defn ^Int round
  "Returns the result of rounding d to the nearest multiple of m. d and m represent time durations in nanoseconds.
  The rounding behavior for halfway values is to round away from zero. If m <= 0, returns d unchanged."
  {:added "1.0"}
  [^Duration d ^Duration m]
  (.round d m))

(defn ^Double seconds
  "Returns the duration (passed as a number of nanoseconds) as a floating point number of seconds."
  {:added "1.0"}
  [^Duration d]
  (.seconds d))

(defn ^String string
  "Returns a string representing the duration in the form 72h3m0.5s."
  {:added "1.0"}
  [d]
  (.string d))

(defn ^Int truncate
  "Returns the result of rounding d toward zero to a multiple of m. If m <= 0, returns d unchanged."
  {:added "1.0"}
  [^Duration d ^Duration m]
  (.truncate d m))
