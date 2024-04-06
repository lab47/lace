(ns lace.map-test
  (:require [lace.test :refer [deftest is]]))


(deftest hash-map-conversion
  (let [m {1 0 2 0 3 0 4 0 5 0 6 0 7 0 8 0}]
    (is (= ArrayMap (type m)))
    (is (= ArrayMap (type (merge m {1 1}))))
    (is (= ArrayMap (type (assoc m 1 2))))
    (is (= HashMap (type (merge m {9 0}))))
    (is (= HashMap (type (assoc m 9 0))))))
