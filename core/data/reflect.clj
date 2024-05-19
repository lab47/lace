(lace.core/in-ns 'lace.reflect)
(lace.core/refer 'lace.core)

(defmacro build
  "Create a new struct value and populate it's fields"
  [typ args]
  (let [id (gensym)]
    (loop [x nil
           args (seq args)]
      (if args
        (let [k (first args)
              v (second args)
              body `(~@x (put ~id ~k ~v))]
          (recur body (nnext args)))
        `(let [~id (~typ)] ~@x ~id)))))
