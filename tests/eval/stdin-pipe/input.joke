(ns stdin-pipe-test
  (:require [lace.os :as os]))

(let [result (os/exec "cat" {:stdin *in*})]
  (print "|" (:out result)))
