(ns lace.test-lace.csv
  (:require [lace.csv :as csv]
            [lace.test :refer [deftest is]]))

(deftest test-csv-seq
  (is (= (csv/csv-seq "a,b,c\nd,e,f") '(["a" "b" "c"] ["d" "e" "f"]))))
