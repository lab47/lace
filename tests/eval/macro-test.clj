(ns lace.macro-test
  (:require [lace.test :refer [deftest is]]
            [lace.string :as s]))

(defmacro try-macro [ & body ] `(try ~@body (catch Error)))
(def try-macro-expand (macroexpand '(try-macro)))

(deftest try-log-test
  (is (= '(try (catch Error))
          try-macro-expand)
      "should properly syntax-quote types"))

(defmacro try-return [ & body ] `(try ~@body (catch Error t# t#)))
(deftest try-expanding-typename
  (is (s/includes? (str (try-return (throw (ex-info "Ouch" {})))) "Exception: Ouch")))
