(ns lace.test-lace.os
  (:require [lace.os :as os]
            [lace.test :refer [deftest is]]))

(deftest exec-pipe
  (if (= (get (os/env) "TTY_TESTS") "1")
    (is (= 0 (:exit (os/exec "stty" {:args ["echo"] :stdin *in*}))))
    (println "Skipping tty tests (STDIN is not a tty)")))
