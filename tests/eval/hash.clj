(ns lace.test-lace.hash
  (:require
   [lace.test :refer [deftest is testing]]))

(deftest stable-hashes
  (testing "stable string hashes"
    (is (= (hash "hey") 4290027229))
    (is (= (hash "there") 3102463325)))
  (testing "stable symbol hashes"
    (is (= (hash 'hey) 3793824397))
    (is (= (hash 'there) 2940266537))
    (is (= (hash 'lace.core/cond) 3232247079))
    (is (= (hash 'lace.repl/doc) 3494663759))
    (is (= (hash 'user/foo) 2980260858)))
  (testing "stable keyword hashes"
    (is (= (hash :hey) 819820356))
    (is (= (hash :there) 1648208352))
    (is (= (hash ::you) 3944753178))
    (is (= (hash :user/foo) 1616868817))))
