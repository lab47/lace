(ns lace.tests.large-forked-stdout
  (:require [lace.os :as os]
            [lace.string :as s]))

(let [exe (str (get (lace.os/env) "PWD") "/lace")
      res (os/sh exe "lots-of-stderr.clj")]
  (print (:out res))
  (let [ev (s/split-lines (:err res))]
    (println-err (ev 0))
    (println-err (ev 1))
    (println-err (ev (- (count ev) 2)))))
