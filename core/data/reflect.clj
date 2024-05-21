(lace.core/in-ns 'lace.reflect)
; Because the sequencing is all weird with loading reflect 
; (as it needs to load before core with lace.lang to resolve symbobls)
; we do this refer fixup here. If we do a full refer here, it will clobber our
; existing vars (like get).
(lace.core/refer 'lace.core :only '(defmacro let loop gensym seq first second nnext list apply vector))

(defmacro build
  "Create a new struct value and populate it's fields"
  [typ args]
  (let [id (gensym)]
    (loop [x nil
           args (seq args)]
      (if args
        (let [k (first args)
              v (second args)
              body `(~@x (set! ~id ~k ~v))]
          (recur body (nnext args)))
        `(let [~id (~typ)] ~@x ~id)))))
