;   Copyright (c) Rich Hickey. All rights reserved.
;   The use and distribution terms for this software are covered by the
;   Eclipse Public License 1.0 (http://opensource.org/licenses/eclipse-1.0.php)
;   which can be found in the file epl-v10.html at the root of this distribution.
;   By using this software in any fashion, you are agreeing to be bound by
;   the terms of this license.
;   You must not remove this notice, or any other, from this software.

; Author: Frantisek Sodomka, Robert Lachlan

(ns lace.test-lace.multimethods
  (:require [lace.test :refer [deftest is are testing]]))

; http://clojure.org/multimethods

; defmulti
; defmethod
; remove-method
; prefer-method
; methods
; prefers

(defmacro for-all
  [& args]
  `(dorun (for ~@args)))

(deftest basic-multimethod-test
  (testing "Check basic dispatch"
    (defmulti too-simple identity)
    (defmethod too-simple :a [x] :a)
    (defmethod too-simple :b [x] :b)
    (defmethod too-simple :default [x] :default)
    (is (= :a (too-simple :a)))
    (is (= :b (too-simple :b)))
    (is (= :default (too-simple :c))))
  ;; (testing "Remove a method works"
  ;;   (remove-method too-simple :a)
  ;;   (is (= :default (too-simple :a))))
  (testing "Add another method works"
    (defmethod too-simple :d [x] :d)
    (is (= :d (too-simple :d)))))

(deftest remove-all-methods-test
  (testing "Core function remove-all-methods works"
    (defmulti simple1 identity)
    (defmethod simple1 :a [x] :a)
    (defmethod simple1 :b [x] :b)
    (is (= {} (methods (remove-all-methods simple1))))))

(deftest methods-test
  (testing "Core function methods works"
    (defmulti simple2 identity)
    (defmethod simple2 :a [x] :a)
    (defmethod simple2 :b [x] :b)
    (is (= #{:a :b} (into #{} (keys (methods simple2)))))
    (is (= :a ((:a (methods simple2)) 1)))
    (defmethod simple2 :c [x] :c)
    (is (= #{:a :b :c} (into #{} (keys (methods simple2)))))
;;    (remove-method simple2 :a)
;;    (is (= #{:b :c} (into #{} (keys (methods simple2)))))
    ))

(deftest get-method-test
  (testing "Core function get-method works"
    (defmulti simple3 identity)
    (defmethod simple3 :a [x] :a)
    (defmethod simple3 :b [x] :b)
    (is (fn? (get-method simple3 :a)))
    (is (= :a ((get-method simple3 :a) 1)))
    (is (fn? (get-method simple3 :b)))
    (is (= :b ((get-method simple3 :b) 1)))
    (is (nil? (get-method simple3 :c)))))
