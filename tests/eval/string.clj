(ns lace.test-lace.string
  (:require [lace.string :as str]
            [lace.test :refer [deftest is]]))

(deftest quoted
  (is (= (str #"a\.\(b\.c\)") (str (str/re-quote "a.(b.c)")))))

(deftest string
  (is (= "a.(b.c)" (str #"a.(b.c)"))))

(deftest regex
  (is (= (str #"a.b.c") (str #"a.b.c"))))

(deftest split-of-regex
  (is (= ["a" "c" "ef"] (str/split "abcdef" #"(b|d)"))))

(deftest split-of-string
  (is (= ["ab" "def"] (str/split "abcdef" "c")))
  (is (= ["abcdef"] (str/split "abcdef" "(b|d)"))))
