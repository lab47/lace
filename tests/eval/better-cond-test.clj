(ns lace.better-cond-test
  (:refer-clojure :exclude [cond if-let when-let if-some when-some])
  (:require [lace.test :refer [deftest are]]
            [lace.better-cond :refer [cond]]))

(deftest better-cond
  (are [x y] (= x y)
    2 (cond (even? 3) 5
            (odd? 3) 2)
    2 (cond (even? 3) 5
            :else 2)
    2 (cond
        :let [x 2]
        x)
    2 (cond
        :when-let [x 2]
        x)
    2 (cond
        :when-some [x 2]
        x)
    nil (cond
          :when-let [x false]
          2)
    2 (cond
        :when-let [x true]
        2)
    nil (cond
          :when-let [x nil]
          2)
    2 (cond
        :when-some [x false]
        2)
    2 (cond
        :when (even? 4)
        2)
    nil (cond
          :when (even? 3)
          2)))
