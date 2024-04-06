(ns lace.test-clojure.types
  (:require [lace.test :refer [deftest is]]))

(deftest test-types
  (is (= (get (lace.core/types__) "Boolean") Boolean))
  (is (= (get (lace.core/types__) "Int") Int)))

(deftest stdin
  (is (instance? BufferedReader *in*)
      "*in* must be BufferedReader (technically it must implement StringReader interface)
      for (read-line) to work"))
