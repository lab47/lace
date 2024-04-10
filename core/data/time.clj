(ns lace.time
  (:require
    [time :as gt]))

(defn now 
  "Return a time value for the current moment." 
  [] (gt/Now))

(defn since 
  "Calculate the elapse time since t."
  [ts] (gt/Since ts))

(defn sleep
  "Sleep the current go routine for specified amount of time."
  [dur] (gt/Sleep dur))

(defn parse-duration
  "Parse the given string as a duration value"
  [str] (gt/ParseDuration str))

