(ns deps-test
  (:require
    [lace.test :refer [deftest is are testing]]
    [deps]
    [test-local.lib :as lib-local]))

(deftest local-lib-test
  (testing "require from a local source"
    (is (= lib-local/b :b))))
