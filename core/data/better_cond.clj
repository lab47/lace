(ns lace.better-cond
  "A collection of variations on Clojure's core macros."
  {:author "Christophe Grand and Mark Engelberg"
   :added "1.0"}
  (:refer-clojure :exclude [cond when-let if-let when-some if-some]))

(defmacro if-let
  "A variation on if-let where all the exprs in the bindings vector must be true.
   Also supports :let."
  {:added "1.0"}
  ([bindings then]
   `(if-let ~bindings ~then nil))
  ([bindings then else]
   (if (seq bindings)
     (if (or (= :let (bindings 0)) (= 'let (bindings 0)))
       `(let ~(bindings 1)
          (if-let ~(subvec bindings 2) ~then ~else))
       `(let [test# ~(bindings 1)]
          (if test#
            (let [~(bindings 0) test#]
              (if-let ~(subvec bindings 2) ~then ~else))
            ~else)))
     then)))

(defmacro when-let
  "A variation on when-let where all the exprs in the bindings vector must be true.
   Also supports :let."
  {:added "1.0"}
  [bindings & body]
  `(if-let ~bindings (do ~@body)))

(defmacro if-some
  "A variation on if-some where all the exprs in the bindings vector must be non-nil.
   Also supports :let."
  {:added "1.0"}
  ([bindings then]
   `(if-some ~bindings ~then nil))
  ([bindings then else]
   (if (seq bindings)
     (if (or (= :let (bindings 0)) (= 'let (bindings 0)))
       `(let ~(bindings 1)
          (if-some ~(subvec bindings 2) ~then ~else))
       `(let [test# ~(bindings 1)]
          (if (nil? test#)
            ~else
            (let [~(bindings 0) test#]
              (if-some ~(subvec bindings 2) ~then ~else)))))
     then)))

(defmacro when-some
  "A variation on when-some where all the exprs in the bindings vector must be non-nil.
   Also supports :let."
  {:added "1.0"}
  [bindings & body]
  `(if-some ~bindings (do ~@body)))

(defmacro cond
  "A variation on cond which sports let bindings, do and implicit else:
     (cond 
       (odd? a) 1
       :do (println a)
       :let [a (quot a 2)]
       (odd? a) 2
       3).
   Also supports :when-let and :when-some. 
   :let, :when-let, :when-some and :do do not need to be written as keywords."
  {:added "1.0"}
  [& clauses]
  (when-let [[test expr & more-clauses] (seq clauses)]
    (if (next clauses)
      (if (or (= :do test) (= 'do test))
        `(do ~expr (cond ~@more-clauses))
        (if (or (= :let test) (= 'let test))
          `(let ~expr (cond ~@more-clauses))
          (if (or (= :when test) (= 'when test))
            `(when ~expr (cond ~@more-clauses))
            (if (or (= :when-let test) (= 'when-let test))
              `(when-let ~expr (cond ~@more-clauses))
              (if (or (= :when-some test) (= 'when-some test))
                `(when-some ~expr (cond ~@more-clauses))
                `(if ~test ~expr (cond ~@more-clauses)))))))
      test)))

