; This code is a modified version of original Clojure's implementation
; (https://github.com/clojure/clojure/blob/master/src/clj/clojure/core.clj)
; which is licensed as follows:

;   Copyright (c) Rich Hickey. All rights reserved.
;   The use and distribution terms for this software are covered by the
;   Eclipse Public License 1.0 (http://opensource.org/licenses/eclipse-1.0.php)
;   which can be found in the file epl-v10.html at the root of this distribution.
;   By using this software in any fashion, you are agreeing to be bound by
;   the terms of this license.
;   You must not remove this notice, or any other, from this software.


(def
  ^{:arglists '([& items])
    :doc "Creates a new list containing the items."
    :added "1.0"
    :tag List}
  list lace.lang/List)

(def
  ^{:arglists '([x seq])
    :doc "Returns a new seq where x is the first element and seq is
         the rest."
         :added "1.0"
         :tag Seq}
  cons lace.lang/Cons)

;during bootstrap we don't have destructuring let, loop or fn, will redefine later
(def
  ^{:added "1.0"}
  let (fn* [&form &env & decl] (cons 'let* decl)))

(set-macro__ #'let)

(def
  ^{:added "1.0"}
  loop (fn* [&form &env & decl] (cons 'loop* decl)))

(set-macro__ #'loop)

(def
  ^{:added "1.0"}
  fn (fn* [&form &env & decl]
          (with-meta__
            (cons 'fn* decl)
            (meta__ &form))))

(set-macro__ #'fn)

(def
  ^{:arglists '([coll])
    :doc "Returns the first item in the collection. Calls seq on its
         argument. If coll is nil, returns nil."
         :added "1.0"}
  first lace.lang/First)

(def
  ^{:arglists '([coll])
    :doc "Returns a seq of the items after the first. Calls seq on its
         argument.  If there are no more items, returns nil."
         :added "1.0"
         :tag Seq}
  next lace.lang/Next)

(def
  ^{:arglists '([coll])
    :doc "Returns a possibly empty seq of the items after the first. Calls seq on its
         argument."
         :added "1.0"
         :tag Seq}
  rest lace.lang/Rest)

(def
  ^{:arglists '([coll x] [coll x & xs])
    :doc "conj[oin]. Returns a new collection with the xs
         'added'. (conj nil item) returns (item).  The 'addition' may
         happen at different 'places' depending on the concrete type."
         :added "1.0"}
  ; TODO: types
  conj (fn 
         (^Collection [coll x] (lace.lang/Conj coll x))
         (^Collection [coll x & xs]
          (if xs
            (recur (lace.lang/Conj coll x) (first xs) (next xs))
            (lace.lang/Conj coll x)))))

(def
  ^{:doc "Same as (first (next x))"
    :arglists '([x])
    :added "1.0"}
  second (fn [^Seqable x] (first (next x))))

(def
  ^{:doc "Same as (first (first x))"
    :arglists '([x])
    :added "1.0"}
  ffirst (fn [^Seqable x] (first (first x))))

(def
  ^{:doc "Same as (next (first x))"
    :arglists '([x])
    :added "1.0"
    :tag Seq}
  nfirst (fn [^Seqable x] (next (first x))))

(def
  ^{:doc "Same as (first (next x))"
    :arglists '([x])
    :added "1.0"}
  fnext (fn [^Seqable x] (first (next x))))

(def
  ^{:doc "Same as (next (next x))"
    :arglists '([x])
    :added "1.0"
    :tag Seq}
  nnext (fn [^Seqable x] (next (next x))))

(def
  ^{:arglists '([coll])
    :doc "Returns a seq on the collection. If the collection is
         empty, returns nil.  (seq nil) returns nil."
         :added "1.0"
         :tag Seq}
  seq lace.lang/Seq)

(def
  ^{:arglists '([c x])
    :doc "Evaluates x and tests if it is an instance of type
         c. Returns true or false"
         :added "1.0"
         :tag Boolean}
  instance? instance?__)

(def
  ^{:arglists '([x])
    :doc "Returns true if x is a sequence"
    :added "1.0"}
  seq? (fn ^Boolean [x] (instance? Seq x)))

(def
  ^{:arglists '([x])
    :doc "Returns true if x is a Char"
    :added "1.0"}
  char? (fn ^Boolean [x] (instance? Char x)))

(def
  ^{:arglists '([x])
    :doc "Returns true if x is a String"
    :added "1.0"}
  string? (fn ^Boolean [x] (instance? String x)))

(def
  ^{:arglists '([x])
    :doc "Returns true if x is a map"
    :added "1.0"}
  map? (fn ^Boolean [x] (instance? Map x)))

(def
  ^{:arglists '([x])
    :doc "Returns true if x is a vector"
    :added "1.0"}
  vector? (fn ^Boolean [x] (instance? Vector x)))

(def
  ^{:arglists '([msg map] [msg map cause])
    :doc "Create an instance of ExInfo, an Error that carries a map of additional data."
    :added "1.0"
    :tag ExInfo}
  ex-info ex-info__)

(def
  ^{:arglists '([map key val] [map key val & kvs])
    :doc "`assoc[iate]. When applied to a map, returns a new map of the
         same (hashed/sorted) type, that contains the mapping of key(s) to
         val(s). When applied to a vector, returns a new vector that
         contains val at index. Note - index must be <= (count vector)."
         :added "1.0"}
  assoc
  (fn
    (^Map [^Associative map key val] (assoc__ map key val))
    (^Map [^Associative map key val & kvs]
     (let [ret (assoc__ map key val)]
       (if kvs
         (if (next kvs)
           (recur ret (first kvs) (second kvs) (nnext kvs))
           (throw (ex-info "assoc expects even number of arguments after map/vector, found odd number" {})))
         ret)))))

(def
  ^{:arglists '([obj])
    :doc "Returns the metadata of obj, returns nil if there is no metadata."
    :added "1.0"
    :tag Map}
  meta meta__)

(def
  ^{:arglists '([obj m])
    :doc "Returns an object of the same type and value as obj, with
         map m as its metadata."
         :added "1.0"}
  with-meta with-meta__)

(def ^{:private true :dynamic true}
  assert-valid-fdecl (fn [form fdecl]))

(def
  ^{:private true}
  sigs
  (fn [form fdecl]
    (assert-valid-fdecl form fdecl)
    (let [asig
          (fn [fdecl]
            (let [arglist (first fdecl)
                  ;elide implicit macro args
                  arglist (if (=__ '&form (first arglist))
                            (subvec__ arglist 2 (count__ arglist))
                            arglist)
                  body (next fdecl)]
              (if (map? (first body))
                (if (next body)
                  (with-meta arglist (conj (if (meta arglist) (meta arglist) {}) (first body)))
                  arglist)
                arglist)))]
      (if (seq? (first fdecl))
        (loop [ret [] fdecls fdecl]
          (if fdecls
            (recur (conj ret (asig (first fdecls))) (next fdecls))
            (seq ret)))
        (list (asig fdecl))))))

(def
  ^{:arglists '([coll])
    :doc "Return the last item in coll, in linear time."
    :added "1.0"}
  last (fn [^Seqable s]
         (if (next s)
           (recur (next s))
           (first s))))

(def
  ^{:arglists '([coll])
    :doc "Return a seq of all but the last item in coll, in linear time."
    :added "1.0"}
  butlast (fn ^Seq [^Seqable s]
            (loop [ret [] s s]
              (if (next s)
                (recur (conj ret (first s)) (next s))
                (seq ret)))))

(def

  ^{:doc "Same as (def name (fn [params* ] exprs*)) or (def
         name (fn ([params* ] exprs*)+)) with any doc-string or attrs added
         to the var metadata. prepost-map defines a map with optional keys
         :pre and :post that contain collections of pre or post conditions."
         :arglists '([name doc-string? attr-map? [params*] prepost-map? body]
                     [name doc-string? attr-map? ([params*] prepost-map? body)+ attr-map?])
         :added "1.0"}
  defn (fn [&form &env name & fdecl]
         ;; Note: Cannot delegate this check to def because of the call to (with-meta name ..)
         (if (instance? Symbol name)
           nil
           (throw (ex-info "First argument to defn must be a symbol" {:form name})))
         (let [m (if (string? (first fdecl))
                   {:doc (first fdecl)}
                   {})
               fdecl (if (string? (first fdecl))
                       (next fdecl)
                       fdecl)
               m (if (map? (first fdecl))
                   (conj m (first fdecl))
                   m)
               fdecl (if (map? (first fdecl))
                       (next fdecl)
                       fdecl)
               fdecl (if (vector? (first fdecl))
                       (list fdecl)
                       fdecl)
               m (if (map? (last fdecl))
                   (conj m (last fdecl))
                   m)
               fdecl (if (map? (last fdecl))
                       (butlast fdecl)
                       fdecl)
               m (conj {:arglists (list 'quote (sigs &form fdecl))} m)
               m (conj (if (meta name) (meta name) {}) m)]
           (list 'def (with-meta name m)
                 (cons `fn fdecl)))))

(set-macro__ #'defn)

(defn cast
  "Throws an error if x is not of a type t, else returns x."
  {:added "1.0"}
  [^Type t x]
  (cast__ t x))

(def
  ^{:arglists '([coll])
    :doc "Creates a new vector containing the contents of coll."
    :added "1.0"
    :tag Vector}
  vec vec__)

(defn vector
  "Creates a new vector containing the args."
  {:added "1.0"}
  ^Vector [& args]
  (vec args))

(def
  ^{:arglists '([& keyvals])
    :doc "keyval => key val
         Returns a new hash map with supplied mappings.  If any keys are
         equal, they are handled as if by repeated uses of assoc."
         :added "1.0"
         :tag HashMap}
  hash-map hash-map__)

(def
  ^{:arglists '([& keys])
    :doc "Returns a new hash set with supplied keys.  Any equal keys are
         handled as if by repeated uses of conj."
         :added "1.0"
         :tag MapSet}
  hash-set hash-set__)

(defn nil?
  "Returns true if x is nil, false otherwise."
  {:tag Boolean
   :added "1.0"}
  [x] (=__ x nil))

(def

  ^{:doc "Like defn, but the resulting function name is declared as a
         macro and will be used as a macro by the compiler when it is
         called."
         :arglists '([name doc-string? attr-map? [params*] body]
                     [name doc-string? attr-map? ([params*] body)+ attr-map?])
         :added "1.0"}
  defmacro (fn [&form &env
                name & args]
             (let [prefix (loop [p (list name) args args]
                            (let [f (first args)]
                              (if (string? f)
                                (recur (cons f p) (next args))
                                (if (map? f)
                                  (recur (cons f p) (next args))
                                  p))))
                   fdecl (loop [fd args]
                           (if (string? (first fd))
                             (recur (next fd))
                             (if (map? (first fd))
                               (recur (next fd))
                               fd)))
                   fdecl (if (vector? (first fdecl))
                           (list fdecl)
                           fdecl)
                   add-implicit-args (fn [fd]
                                       (let [args (first fd)]
                                         (cons (vec (cons '&form (cons '&env args))) (next fd))))
                   add-args (fn [acc ds]
                              (if (nil? ds)
                                acc
                                (let [d (first ds)]
                                  (if (map? d)
                                    (conj acc d)
                                    (recur (conj acc (add-implicit-args d)) (next ds))))))
                   fdecl (seq (add-args [] fdecl))
                   decl (loop [p prefix d fdecl]
                          (if p
                            (recur (next p) (cons (first p) d))
                            d))]
               (list 'do
                     (cons `defn decl)
                     (list 'set-macro__ (list 'var name))))))

(set-macro__ #'defmacro)

(defmacro when
  "Evaluates test. If logical true, evaluates body in an implicit do."
  {:added "1.0"}
  [test & body]
  (let [b (if (>__ (count__ body) 1)
            (cons 'do body)
            (first body))]
    (list 'if test b nil)))

(defmacro when-not
  "Evaluates test. If logical false, evaluates body in an implicit do."
  {:added "1.0"}
  [test & body]
  (let [b (if (>__ (count__ body) 1)
              (cons 'do body)
              (first body))]
      (list 'if test nil b)))

(defn false?
  "Returns true if x is the value false, false otherwise."
  {:tag Boolean
   :added "1.0"}
  [x] (=__ x false))

(defn true?
  "Returns true if x is the value true, false otherwise."
  {:tag Boolean
   :added "1.0"}
  [x] (=__ x true))

(defn boolean?
  "Return true if x is a Boolean"
  {:tag Boolean
   :added "1.0"}
  [x] (instance? Boolean x))

(defn any?
  "Returns true given any argument."
  {:tag Boolean
   :added "1.0"}
  [x] true)

(defn not
  "Returns true if x is logical false, false otherwise."
  {:tag Boolean
   :added "1.0"}
  ^Boolean [x] (if x false true))

(defn some?
  "Returns true if x is not nil, false otherwise."
  {:tag Boolean
   :added "1.0"}
  [x] (not (nil? x)))

(def
  ^{:arglists '([& xs])
    :doc "With no args, returns the empty string. With one arg x, returns
         string representation of x. (str nil) returns the empty string. With more than
         one arg, returns the concatenation of the str values of the args."
         :added "1.0"
         :tag String}
  str str__)

(defn symbol?
  "Return true if x is a Symbol"
  {:added "1.0"}
  ^Boolean [x] (instance? Symbol x))

(defn keyword?
  "Return true if x is a Keyword"
  {:added "1.0"}
  ^Boolean [x] (instance? Keyword x))

(defn symbol
  "Returns a Symbol with the given namespace and name."
  {:added "1.0"}
  ;; TODO: types
  (^Symbol [name] (if (symbol? name) name (symbol__ name)))
  (^Symbol [ns name] (symbol__ ns name)))

(defn gensym
  "Returns a new symbol with a unique name. If a prefix string is
  supplied, the name is prefix# where # is some unique number. If
  prefix is not supplied, the prefix is 'G__'."
  {:added "1.0"}
  (^Symbol [] (gensym "G__"))
  (^Symbol [^String prefix-string] (gensym__ prefix-string)))

(def println-err)
(def println-linter__)

(defmacro cond
  "Takes a set of test/expr pairs. It evaluates each test one at a
  time.  If a test returns logical true, cond evaluates and returns
  the value of the corresponding expr and doesn't evaluate any of the
  other tests or exprs. (cond) returns nil."
  {:added "1.0"}
  [& clauses]
  (if clauses
    (list 'if (first clauses)
          (if (next clauses)
            (second clauses)
            (throw (ex-info "cond requires an even number of forms" {:form (first clauses)})))
          (when (next (next clauses))
            (cons 'lace.core/cond (next (next clauses)))))
    (when *linter-mode*
      (println-linter__ (ex-info "Empty cond" {:form &form :_prefix "Parse warning"})))))

(defn keyword
  "Returns a Keyword with the given namespace and name.  Do not use :
  in the keyword strings, it will be added automatically."
  {:tag Keyword
   :added "1.0"}
  ;; TODO: types
  (^Keyword [name] (cond (keyword? name) name
                     (symbol? name) (keyword__ name)
                     (string? name) (keyword__ name)))
  (^Keyword [ns name] (keyword__ ns name)))

(defn spread
  {:private true}
  [arglist]
  (cond
    (nil? arglist) nil
    (nil? (next arglist)) (seq (first arglist))
    :else (cons (first arglist) (spread (next arglist)))))

(defn list*
  "Creates a new list containing the items prepended to the rest, the
  last of which will be treated as a sequence."
  {:added "1.0"}
  (^Seq [^Seqable args] (seq args))
  (^Seq [a ^Seqable args] (cons a args))
  (^Seq [a b ^Seqable args] (cons a (cons b args)))
  (^Seq [a b c ^Seqable args] (cons a (cons b (cons c args))))
  (^Seq [a b c d & more]
   (cons a (cons b (cons c (cons d (spread more)))))))

(defn apply
  "Applies fn f to the argument list formed by prepending intervening arguments to args."
  {:added "1.0"}
  ([^Callable f ^Seqable args]
   (apply__ f (seq args)))
  ([^Callable f x ^Seqable args]
   (apply__ f (list* x args)))
  ([^Callable f x y ^Seqable args]
   (apply__ f (list* x y args)))
  ([^Callable f x y z ^Seqable args]
   (apply__ f (list* x y z args)))
  ([^Callable f a b c d & args]
   (apply__ f (cons a (cons b (cons c (cons d (spread args))))))))

(defn vary-meta
  "Returns an object of the same type and value as obj, with
  (apply f (meta obj) args) as its metadata."
  {:added "1.0"}
  [obj ^Callable f & args]
  (with-meta obj (apply f (meta obj) args)))

(defmacro lazy-seq
  "Takes a body of expressions that returns an ISeq or nil, and yields
  a Seqable object that will invoke the body only the first time seq
  is called, and will cache the result and return it on all subsequent
  seq calls. See also - realized?"
  {:added "1.0"}
  ^Seqable [& body]
  (list 'lace.core/lazy-seq__ (list* 'fn [] body)))

(defn chunked-seq?
  "Always returns false because chunked sequences are not supported"
  {:added "1.0"}
  ^Boolean [s]
  false)

(defn concat
  "Returns a lazy seq representing the concatenation of the elements in the supplied colls."
  {:added "1.0"}
  (^Seq [] (lazy-seq nil))
  (^Seq [^Seqable x] (lazy-seq x))
  (^Seq [^Seqable x ^Seqable y]
   (lazy-seq
    (let [s (seq x)]
      (if s
        (cons (first s) (concat (rest s) y))
        y))))
  (^Seq [^Seqable x ^Seqable y & zs]
   (let [cat (fn cat [xys zs]
               (lazy-seq
                (let [xys (seq xys)]
                  (if xys
                    (cons (first xys) (cat (rest xys) zs))
                    (when zs
                      (cat (first zs) (next zs)))))))]
     (cat (concat x y) zs))))


(defmacro delay
  "Takes a body of expressions and yields a Delay object that will
  invoke the body only the first time it is forced (with force or deref/@), and
  will cache the result and return it on all subsequent force
  calls. See also - realized?"
  {:added "1.0"}
  ^Delay [& body]
  (list 'lace.core/delay__ (list* 'fn [] body)))

(defn delay?
  "returns true if x is a Delay created with delay"
  {:added "1.0"}
  ^Boolean [x]
  (instance? Delay x))

(defn force
  "If x is a Delay, returns the (possibly cached) value of its expression, else returns x"
  {:added "1.0"}
  [x]
  (force__ x))

(defmacro if-not
  "Evaluates test. If logical false, evaluates and returns then expr,
  otherwise else expr, if supplied, else nil."
  {:added "1.0"}
  ([test then] `(if-not ~test ~then nil))
  ([test then else]
   `(if (not ~test) ~then ~else)))

(defn identical?
  "Tests if 2 arguments are the same object"
  {:added "1.0"}
  ^Boolean [x y]
  (identical__ x y))

(defn =
  "Equality. Returns true if x equals y, false if not. Works for nil, and compares
  numbers and collections in a type-independent manner.  Immutable data
  structures define = as a value, not an identity,
  comparison."
  {:added "1.0"}
  (^Boolean [x] true)
  (^Boolean [x y] (=__ x y))
  (^Boolean [x y & more]
   (if (=__ x y)
     (if (next more)
       (recur y (first more) (next more))
       (=__ y (first more)))
     false)))

(defn not=
  "Same as (not (= obj1 obj2))"
  {:added "1.0"}
  (^Boolean [x] false)
  (^Boolean [x y] (not (= x y)))
  (^Boolean [x y & more]
   (not (apply = x y more))))

(defn compare
  "Comparator. Returns a negative number, zero, or a positive number
  when x is logically 'less than', 'equal to', or 'greater than'
  y. Works for nil, and compares numbers and collections in a type-independent manner. x
  must implement Comparable"
  {:added "1.0"}
  ^Int [x y] (compare__ x y))

(defmacro and
  "Evaluates exprs one at a time, from left to right. If a form
  returns logical false (nil or false), and returns that value and
  doesn't evaluate any of the other expressions, otherwise it returns
  the value of the last expr. (and) returns true."
  {:added "1.0"}
  ([] true)
  ([x] x)
  ([x & next]
   `(let [and# ~x]
      (if and# (and ~@next) and#))))

(defmacro or
  "Evaluates exprs one at a time, from left to right. If a form
  returns a logical true value, or returns that value and doesn't
  evaluate any of the other expressions, otherwise it returns the
  value of the last expression. (or) returns nil."
  {:added "1.0"}
  ([] nil)
  ([x] x)
  ([x & next]
   `(let [or# ~x]
      (if or# or# (or ~@next)))))

(defn zero?
  "Returns true if x is zero, else false"
  {:added "1.0"}
  ^Boolean [^Number x] (zero?__ x))

(defn count
  "Returns the number of items in the collection. (count nil) returns
  0.  Also works on strings"
  {:added "1.0"}
  ;; TODO: types
  ^Int [coll] (count__ coll))

(defn int
  "Coerce to int"
  {:added "1.0"}
  ;; TODO: types
  ^Int [x] (int__ x))

(defn nth
  "Returns the value at the index. get returns nil if index out of
  bounds, nth throws an exception unless not-found is supplied.  nth
  also works, in O(n) time, for strings and sequences."
  {:added "1.0"}
  ;; TODO: types
  ([coll ^Number index] (nth__ coll index))
  ([coll ^Number index not-found] (nth__ coll index not-found)))

(defn <
  "Returns non-nil if nums are in monotonically increasing order,
  otherwise false."
  {:added "1.0"}
  (^Boolean [^Number x] true)
  (^Boolean [^Number x ^Number y] (<__ x y))
  (^Boolean [^Number x ^Number y & more]
   (if (< x y)
     (if (next more)
       (recur y (first more) (next more))
       (< y (first more)))
     false)))

(defn inc'
  "Returns a number one greater than num. Supports arbitrary precision.
  See also: inc"
  {:added "1.0"}
  ^Number [^Number x] (inc'__ x))

(defn inc
  "Returns a number one greater than num. Does not auto-promote
  ints, will overflow. See also: inc'"
  {:added "1.0"}
  ^Number [^Number x] (inc__ x))

(defn reduce
  "f should be a function of 2 arguments. If val is not supplied,
  returns the result of applying f to the first 2 items in coll, then
  applying f to that result and the 3rd item, etc. If coll contains no
  items, f must accept no arguments as well, and reduce returns the
  result of calling f with no arguments.  If coll has only 1 item, it
  is returned and f is not called.  If val is supplied, returns the
  result of applying f to val and the first item in coll, then
  applying f to that result and the 2nd item, etc. If coll contains no
  items, returns val and f is not called."
  {:added "1.0"}
  ([^Callable f coll]
   (let [s (seq coll)]
     (if s
       (reduce f (first s) (next s))
       (f))))
  ([^Callable f val coll]
   (let [s (seq coll)]
     (if s
       (recur f (f val (first s)) (next s))
       val))))

(defn reverse
  "Returns a seq of the items in coll in reverse order. Not lazy."
  {:added "1.0"}
  ^Collection [^Seqable coll]
  (reduce conj () coll))

(defn +'
  "Returns the sum of nums. (+) returns 0. Supports arbitrary precision.
  See also: +"
  {:added "1.0"}
  (^Number [] 0)
  (^Number [^Number x] (cast Number x))
  (^Number [^Number x ^Number y] (add'__ x y))
  (^Number [^Number x ^Number y & more]
   (reduce +' (+' x y) more)))

(defn +
  "Returns the sum of nums. (+) returns 0. Does not auto-promote
  ints, will overflow. See also: +'"
  {:added "1.0"}
  (^Number [] 0)
  (^Number [^Number x] (cast Number x))
  (^Number [^Number x ^Number y] (add__ x y))
  (^Number [^Number x ^Number y & more]
   (reduce + (+ x y) more)))

(defn *'
  "Returns the product of nums. (*) returns 1. Supports arbitrary precision.
  See also: *"
  {:added "1.0"}
  (^Number [] 1)
  (^Number [^Number x] (cast Number x))
  (^Number [^Number x ^Number y] (multiply'__ x y))
  (^Number [^Number x ^Number y & more]
   (reduce *' (*' x y) more)))

(defn *
  "Returns the product of nums. (*) returns 1. Does not auto-promote
  ints, will overflow. See also: *'"
  {:added "1.0"}
  (^Number [] 1)
  (^Number [^Number x] (cast Number x))
  (^Number [^Number x ^Number y] (multiply__ x y))
  (^Number [^Number x ^Number y & more]
   (reduce * (* x y) more)))

(defn /
  "If no denominators are supplied, returns 1/numerator,
  else returns numerator divided by all of the denominators."
  {:added "1.0"}
  (^Number [^Number x] (/ 1 x))
  (^Number [^Number x ^Number y] (divide__ x y))
  (^Number [^Number x ^Number y & more]
   (reduce / (/ x y) more)))

(defn -'
  "If no ys are supplied, returns the negation of x, else subtracts
  the ys from x and returns the result. Supports arbitrary precision.
  See also: -"
  {:added "1.0"}
  (^Number [^Number x] (subtract'__ x))
  (^Number [^Number x ^Number y] (subtract'__ x y))
  (^Number [^Number x ^Number y & more]
   (reduce -' (-' x y) more)))

(defn -
  "If no ys are supplied, returns the negation of x, else subtracts
  the ys from x and returns the result. Does not auto-promote
  ints, will overflow. See also: -'"
  {:added "1.0"}
  (^Number [^Number x] (subtract__ x))
  (^Number [^Number x ^Number y] (subtract__ x y))
  (^Number [^Number x ^Number y & more]
   (reduce - (- x y) more)))

(defn <=
  "Returns non-nil if nums are in monotonically non-decreasing order,
  otherwise false."
  {:added "1.0"}
  (^Boolean [^Number x] true)
  (^Boolean [^Number x ^Number y] (<=__ x y))
  (^Boolean [^Number x ^Number y & more]
   (if (<= x y)
     (if (next more)
       (recur y (first more) (next more))
       (<= y (first more)))
     false)))

(defn >
  "Returns non-nil if nums are in monotonically decreasing order,
  otherwise false."
  {:added "1.0"}
  (^Boolean [^Number x] true)
  (^Boolean [^Number x ^Number y] (>__ x y))
  (^Boolean [^Number x ^Number y & more]
   (if (> x y)
     (if (next more)
       (recur y (first more) (next more))
       (> y (first more)))
     false)))

(defn >=
  "Returns non-nil if nums are in monotonically non-increasing order,
  otherwise false."
  {:added "1.0"}
  (^Boolean [^Number x] true)
  (^Boolean [^Number x ^Number y] (>=__ x y))
  (^Boolean [^Number x ^Number y & more]
   (if (>= x y)
     (if (next more)
       (recur y (first more) (next more))
       (>= y (first more)))
     false)))

(defn ==
  "Returns non-nil if nums all have the equivalent
  value (type-independent), otherwise false"
  {:added "1.0"}
  (^Boolean [^Number x] true)
  (^Boolean [^Number x ^Number y] (==__ x y))
  (^Boolean [^Number x ^Number y & more]
   (if (== x y)
     (if (next more)
       (recur y (first more) (next more))
       (== y (first more)))
     false)))

(defn max
  "Returns the greatest of the nums."
  {:added "1.0"}
  (^Number [^Number x] x)
  (^Number [^Number x ^Number y] (max__ x y))
  (^Number [^Number x ^Number y & more]
   (reduce max (max x y) more)))

(defn min
  "Returns the least of the nums."
  {:added "1.0"}
  (^Number [^Number x] x)
  (^Number [^Number x ^Number y] (min__ x y))
  (^Number [^Number x ^Number y & more]
   (reduce min (min x y) more)))

(defn dec'
  "Returns a number one less than num. Supports arbitrary precision.
  See also: dec"
  {:added "1.0"}
  ^Number [^Number x] (dec'__ x))

(defn dec
  "Returns a number one less than num. Does not auto-promote
  ints, will overflow. See also: dec'"
  {:added "1.0"}
  ^Number [^Number x] (dec__ x))

(defn pos?
  "Returns true if num is greater than zero, else false"
  {:added "1.0"}
  ^Boolean [^Number x] (pos__ x))

(defn neg?
  "Returns true if num is less than zero, else false"
  {:added "1.0"}
  ^Boolean [^Number x] (neg__ x))

(defn quot
  "quot[ient] of dividing numerator by denominator."
  {:added "1.0"}
  ^Number [^Number num ^Number div]
  (quot__ num div))

(defn rem
  "remainder of dividing numerator by denominator."
  {:added "1.0"}
  ^Number [^Number num ^Number div]
  (rem__ num div))

(defn bit-not
  "Bitwise complement"
  {:added "1.0"}
  ^Int [^Int x] (bit-not__ x))

(defn bit-and
  "Bitwise and"
  {:added "1.0"}
  (^Int [^Int x ^Int y] (bit-and__ x y))
  (^Int [^Int x ^Int y & more]
   (reduce bit-and (bit-and x y) more)))

(defn bit-or
  "Bitwise or"
  {:added "1.0"}
  (^Int [^Int x ^Int y] (bit-or__ x y))
  (^Int [^Int x ^Int y & more]
   (reduce bit-or (bit-or x y) more)))

(defn bit-xor
  "Bitwise exclusive or"
  {:added "1.0"}
  (^Int [^Int x ^Int y] (bit-xor_ x y))
  (^Int [^Int x ^Int y & more]
   (reduce bit-xor (bit-xor x y) more)))

(defn bit-and-not
  "Bitwise and with complement"
  {:added "1.0"}
  (^Int [^Int x ^Int y] (bit-and-not__ x y))
  (^Int [^Int x ^Int y & more]
   (reduce bit-and-not (bit-and-not x y) more)))

(defn bit-clear
  "Clear bit at index n"
  {:added "1.0"}
  ^Int [^Int x ^Int n] (bit-clear__ x n))

(defn bit-set
  "Set bit at index n"
  {:added "1.0"}
  ^Int [^Int x ^Int n] (bit-set__ x n))

(defn bit-flip
  "Flip bit at index n"
  {:added "1.0"}
  ^Int [^Int x ^Int n] (bit-flip__ x n))

(defn bit-test
  "Test bit at index n"
  {:added "1.0"}
  ^Boolean [^Int x ^Int n] (bit-test__ x n))

(defn bit-shift-left
  "Bitwise shift left"
  {:added "1.0"}
  ^Int [^Int x ^Int n] (bit-shift-left__ x n))

(defn bit-shift-right
  "Bitwise shift right"
  {:added "1.0"}
  ^Int [^Int x ^Int n] (bit-shift-right__ x n))

(defn unsigned-bit-shift-right
  "Bitwise shift right, without sign-extension."
  {:added "1.0"}
  ^Int [^Int x ^Int n] (unsigned-bit-shift-right__ x n))

(defn integer?
  "Returns true if n is an integer"
  {:added "1.0"}
  ^Boolean [n]
  (or (instance? Int n)
      (instance? BigInt n)))

(defn even?
  "Returns true if n is even, throws an exception if n is not an integer"
  {:added "1.0"}
  ;; TODO: types (Int, BigInt)
  ^Boolean [n]
  (if (integer? n)
    (zero? (bit-and (int__ n) 1))
    (throw (ex-info (str "Argument must be an integer: " n) {}))))

(defn odd?
  "Returns true if n is odd, throws an exception if n is not an integer"
  {:added "1.0"}
  ;; TODO: types (Int, BigInt)
  ^Boolean [n] (not (even? n)))

(defn int?
  "Return true if x is a fixed precision integer"
  {:tag Boolean
   :added "1.0"}
  ^Boolean [x] (instance? Int x))

(defn pos-int?
  "Return true if x is a positive fixed precision integer"
  {:tag Boolean
   :added "1.0"}
  ^Boolean [x] (and (int? x)
                    (pos? x)))

(defn neg-int?
  "Return true if x is a negative fixed precision integer"
  {:tag Boolean
   :added "1.0"}
  ^Boolean [x] (and (int? x)
                    (neg? x)))

(defn nat-int?
  "Return true if x is a non-negative fixed precision integer"
  {:tag Boolean
   :added "1.0"}
  ^Boolean [x] (and (int? x)
                    (not (neg? x))))

(defn double?
  "Return true if x is a Double"
  {:tag Boolean
   :added "1.0"}
  ^Boolean [x] (instance? Double x))

(defn complement
  "Takes a fn f and returns a fn that takes the same arguments as f,
  has the same effects, if any, and returns the opposite truth value."
  {:added "1.0"}
  ^Fn [^Callable f]
  (fn
    ([] (not (f)))
    ([x] (not (f x)))
    ([x y] (not (f x y)))
    ([x y & zs] (not (apply f x y zs)))))

(defn constantly
  "Returns a function that takes any number of arguments and returns x."
  {:added "1.0"}
  ^Fn [x]
  (fn [& args] x))

(defn identity
  "Returns its argument."
  {:added "1.0"}
  [x] x)

(defn peek
  "For a list, same as first, for a vector, same as, but much
  more efficient than, last. If the collection is empty, returns nil."
  {:added "1.0"}
  [^Stack coll] (peek__ coll))

(defn pop
  "For a list, returns a new list without the first
  item, for a vector, returns a new vector without the last item. If
  the collection is empty, throws an exception.  Note - not the same
  as next/butlast."
  {:added "1.0"}
  [^Stack coll] (pop__ coll))

(defn contains?
  "Returns true if key is present in the given collection, otherwise
  returns false.  Note that for numerically indexed collections like
  vectors, this tests if the numeric key is within the
  range of indexes. 'contains?' operates constant or logarithmic time;
  it will not perform a linear search for a value.  See also 'some'."
  {:added "1.0"}
  ^Boolean [^Gettable coll key] (contains?__ coll key))

(defn get
  "Returns the value mapped to key, not-found or nil if key not present."
  {:added "1.0"}
  ([map key]
   (get__ map key))
  ([map key not-found]
   (get__ map key not-found)))

(defn dissoc
  "dissoc[iate]. Returns a new map of the same (hashed/sorted) type,
  that does not contain a mapping for key(s)."
  {:added "1.0"}
  (^Map [^Map map] map)
  (^Map [^Map map key]
   (dissoc__ map key))
  (^Map [^Map map key & ks]
   (let [ret (dissoc__ map key)]
     (if ks
       (recur ret (first ks) (next ks))
       ret))))

(defn disj
  "disj[oin]. Returns a new set of the same (hashed/sorted) type, that
  does not contain key(s)."
  {:added "1.0"}
  (^MapSet [^Set set] set)
  (^MapSet [^Set set key]
   (disj__ set key))
  (^MapSet [^Set set key & ks]
   (when set
     (let [ret (disj__ set key)]
       (if ks
         (recur ret (first ks) (next ks))
         ret)))))

(defn find
  "Returns the map entry for key, or nil if key not present."
  {:added "1.0"}
  [^Associative map key] (find__ map key))

(defn select-keys
  "Returns a map containing only those entries in map whose key is in keys"
  {:added "1.0"}
  ^Map [^Associative map ^Seqable keyseq]
  (loop [ret {} keys (seq keyseq)]
    (if keys
      (let [entry (find__ map (first keys))]
        (recur
         (if entry
           (conj ret entry)
           ret)
         (next keys)))
      (with-meta ret (meta map)))))

(defn keys
  "Returns a sequence of the map's keys, in the same order as (seq map)."
  {:added "1.0"}
  ^Seq [^Map map] (keys__ map))

(defn vals
  "Returns a sequence of the map's values, in the same order as (seq map)."
  {:added "1.0"}
  ^Seq [^Map map] (vals__ map))

(defn key
  "Returns the key of the map entry."
  {:added "1.0"}
  [e]
  (first e))

(defn val
  "Returns the value in the map entry."
  {:added "1.0"}
  [e]
  (second e))

(defn rseq
  "Returns, in constant time, a seq of the items in rev (which
  can be a vector or sorted-map), in reverse order. If rev is empty returns nil."
  {:added "1.0"}
  ^Seq [^Reversible rev]
  (rseq__ rev))

(defn name
  "Returns the name String of a string, symbol or keyword."
  {:tag String
   :added "1.0"}
  ;; TODO: types
  ^String [x]
  (if (string? x) x (name__ x)))

(defn namespace
  "Returns the namespace String of a symbol or keyword, or nil if not present."
  {:tag String
   :added "1.0"}
  ^String [^Named x]
  (namespace__ x))

(defn ident?
  "Return true if x is a symbol or keyword"
  {:added "1.0"}
  ^Boolean [x]
  (or (keyword? x) (symbol? x)))

(defn simple-ident?
  "Return true if x is a symbol or keyword without a namespace"
  {:added "1.0"}
  ^Boolean [x] (and (ident? x) (nil? (namespace x))))

(defn qualified-ident?
  "Return true if x is a symbol or keyword with a namespace"
  {:added "1.0"}
  ^Boolean [x] (and (ident? x) (namespace x) true))

(defn simple-symbol?
  "Return true if x is a symbol without a namespace"
  {:added "1.0"}
  ^Boolean [x] (and (symbol? x) (nil? (namespace x))))

(defn qualified-symbol?
  "Return true if x is a symbol with a namespace"
  {:added "1.0"}
  ^Boolean [x] (and (symbol? x) (namespace x) true))

(defn simple-keyword?
  "Return true if x is a keyword without a namespace"
  {:added "1.0"}
  ^Boolean [x] (and (keyword? x) (nil? (namespace x))))

(defn qualified-keyword?
  "Return true if x is a keyword with a namespace"
  {:added "1.0"}
  ^Boolean [x] (and (keyword? x) (namespace x) true))

(defmacro ->
  "Threads the expr through the forms. Inserts x as the
  second item in the first form, making a list of it if it is not a
  list already. If there are more forms, inserts the first form as the
  second item in second form, etc."
  {:added "1.0"}
  [x & forms]
  (when (and *linter-mode* (not (seq forms)) (not (false? (:no-forms-threading (:rules *linter-config*)))))
    (println-linter__ (ex-info "No forms in ->" {:form &form :_prefix "Parse warning"})))
  (loop [x x forms forms]
    (if forms
      (let [form (first forms)
            threaded (if (seq? form)
                       (with-meta `(~(first form) ~x ~@(next form)) (meta form))
                       (list form x))]
        (recur threaded (next forms)))
      x)))

(defmacro ->>
  "Threads the expr through the forms. Inserts x as the
  last item in the first form, making a list of it if it is not a
  list already. If there are more forms, inserts the first form as the
  last item in second form, etc."
  {:added "1.0"}
  [x & forms]
  (when (and *linter-mode* (not (seq forms)) (not (false? (:no-forms-threading (:rules *linter-config*)))))
    (println-linter__ (ex-info "No forms in ->>" {:form &form :_prefix "Parse warning"})))
  (loop [x x forms forms]
    (if forms
      (let [form (first forms)
            threaded (if (seq? form)
                       (with-meta `(~(first form) ~@(next form)  ~x) (meta form))
                       (list form x))]
        (recur threaded (next forms)))
      x)))

(defmacro ^{:private true} assert-args
  [& pairs]
  `(do (when-not ~(first pairs)
         (throw (ex-info
                 (str (first ~'&form) " requires " ~(second pairs))
                 {:form ~'&form})))
     ~(let [more (nnext pairs)]
        (when more
          (list* `assert-args more)))))

(defmacro if-let
  "bindings => binding-form test

  If test is true, evaluates then with binding-form bound to the value of
  test, if not, yields else"
  {:added "1.0"}
  ([bindings then]
   (assert-args
    (vector? bindings) "a vector for its binding"
    (= 2 (count bindings)) "exactly 2 forms in binding vector")
   (let [form (bindings 0) tst (bindings 1)]
     `(let [temp# ~tst]
        (if temp#
          (let [~form temp#
                t# ~then] ; to avoid "redundant do form" warning from linter if ~then is (do...)
            t#)))))
  ([bindings then else & oldform]
   (assert-args
    (vector? bindings) "a vector for its binding"
    (nil? oldform) "1 or 2 forms after binding vector"
    (= 2 (count bindings)) "exactly 2 forms in binding vector")
   (let [form (bindings 0) tst (bindings 1)]
     `(let [temp# ~tst]
        (if temp#
          (let [~form temp#
                t# ~then] ; to avoid "redundant do form" warning from linter if ~then is (do...)
            t#)
          ~else)))))

(defmacro when-let
  "bindings => binding-form test

  When test is true, evaluates body with binding-form bound to the value of test"
  {:added "1.0"}
  [bindings & body]
  (assert-args
   (vector? bindings) "a vector for its binding"
   (= 2 (count bindings)) "exactly 2 forms in binding vector")
  (let [form (bindings 0) tst (bindings 1)]
    `(let [temp# ~tst]
       (when temp#
         (let [~form temp#]
           ~@body)))))

(defmacro if-some
  "bindings => binding-form test

  If test is not nil, evaluates then with binding-form bound to the
  value of test, if not, yields else"
  {:added "1.0"}
  ([bindings then]
   `(if-some ~bindings ~then nil))
  ([bindings then else & oldform]
   (assert-args
    (vector? bindings) "a vector for its binding"
    (nil? oldform) "1 or 2 forms after binding vector"
    (= 2 (count bindings)) "exactly 2 forms in binding vector")
   (let [form (bindings 0) tst (bindings 1)]
     `(let [temp# ~tst]
        (if (nil? temp#)
          ~else
          (let [~form temp#]
            ~then))))))

(defmacro when-some
  "bindings => binding-form test

  When test is not nil, evaluates body with binding-form bound to the
  value of test"
  {:added "1.0"}
  [bindings & body]
  (assert-args
   (vector? bindings) "a vector for its binding"
   (= 2 (count bindings)) "exactly 2 forms in binding vector")
  (let [form (bindings 0) tst (bindings 1)]
    `(let [temp# ~tst]
       (if (nil? temp#)
         nil
         (let [~form temp#]
           ~@body)))))

(defn class
  "Returns the Type of x."
  {:added "1.0"}
  ^Type [x]
  (type__ x))

(defn type
  "Returns the :type metadata of x, or its Type if none"
  {:added "1.0"}
  ^Type [x]
  (or (get (meta x) :type) (type__ x)))

(defn reduce-kv
  "Reduces an associative collection. f should be a function of 3
  arguments. Returns the result of applying f to init, the first key
  and the first value in coll, then applying f to that result and the
  2nd key and value, etc. If coll contains no entries, returns init
  and f is not called. Note that reduce-kv is supported on vectors,
  where the keys will be the ordinals."
  {:added "1.0"}
  ;; TODO: types
  ([^Callable f init coll]
   (cond
     (instance? KVReduce coll)
     (reduce-kv__ f init coll)

     (instance? Map coll)
     (reduce (fn [ret kv] (f ret (first kv) (second kv))) init coll)

     :else
     (throw (ex-info (str "Cannot reduce-kv on " (type coll)) {})))))

(defn var-get
  "Gets the value in the var object"
  {:added "1.0"}
  [^Var x] (var-get__ x))

(defn var-set
  "Sets the value in the var object to val."
  {:added "1.0"}
  [^Var x val] (var-set__ x val))

(defn ^:private replace-bindings
  [binding-map]
  (reduce-kv (fn [res k v]
               (let [c (var-get k)]
                 (var-set k v)
                 (assoc res k c)))
             {}
             binding-map))

(defn with-bindings*
  "Takes a map of Var/value pairs. Sets the vars to the corresponding values.
  Then calls f with the supplied arguments. Resets the vars back to the original
  values after f returned. Returns whatever f returns."
  {:added "1.0"}
  [^Map binding-map ^Callable f & args]
  (let [existing-bindings (lace.lang/PushBindings binding-map)]
    (try
      (apply f args)
      (finally
        (lace.lang/SetBindings existing-bindings)))))

(def
  ^{:doc "The same as with-bindings*"
    :arglists '([binding-map f & args])
    :added "1.0"}
  with-redefs-fn with-bindings*)

(defmacro with-bindings
  "Takes a map of Var/value pairs. Sets the vars to the corresponding values.
  Then executes body. Resets the vars back to the original
  values after body was evaluated. Returns the value of body."
  {:added "1.0"}
  [binding-map & body]
  `(with-bindings* ~binding-map (fn [] ~@body)))

(defmacro binding
  "binding => var-symbol init-expr

  Creates new bindings for the (already-existing) vars, with the
  supplied initial values, executes the exprs in an implicit do, then
  re-establishes the bindings that existed before.  The new bindings
  are made in parallel (unlike let); all init-exprs are evaluated
  before the vars are bound to their new values."
  {:added "1.0"}
  [bindings & body]
  (assert-args
    (vector? bindings) "a vector for its binding"
    (even? (count bindings)) "an even number of forms in binding vector")
  (let [var-ize (fn [var-vals]
                  (loop [ret [] vvs (seq var-vals)]
                    (if vvs
                      (recur  (conj (conj ret `(var ~(first vvs))) (second vvs))
                             (next (next vvs)))
                      (seq ret))))]
    `(with-bindings (hash-map ~@(var-ize bindings)) ~@body)))

(defmacro with-redefs
  "The same as binding"
  {:added "1.0"}
  [bindings & body]
  `(binding ~bindings ~@body))

(defn deref
  "Also reader macro: @var/@atom/@delay. When applied to a var or atom,
  returns its current state. When applied to a delay, forces
  it if not already forced."
  {:added "1.0"}
  [^Deref ref]
  (deref__ ref))

(defn atom
  "Creates and returns an Atom with an initial value of x and zero or
  more options (in any order):

  :meta metadata-map

  If metadata-map is supplied, it will become the metadata on the
  atom."
  {:added "1.0"}
  ^Atom [x & options]
  (apply atom__ x options))

(defn swap!
  "Atomically swaps the value of atom to be:
  (apply f current-value-of-atom args).
  Returns the value that was swapped in."
  {:added "1.0"}
  [^Atom atom ^Callable f & args]
  (apply swap__ atom f args))

(defn swap-vals!
  "Atomically swaps the value of atom to be:
  (apply f current-value-of-atom args). Note that f may be called
  multiple times, and thus should be free of side effects.
  Returns [old new], the value of the atom before and after the swap."
  {:added "1.0"}
  ^Vector [^Atom atom ^Callable f & args]
  (apply swap-vals__ atom f args))

(defn reset!
  "Sets the value of atom to newval without regard for the
  current value. Returns newval."
  {:added "1.0"}
  [^Atom atom newval]
  (reset__ atom newval))

(defn reset-vals!
  "Sets the value of atom to newval. Returns [old new], the value of the
  atom before and after the reset."
  {:added "1.0"}
  ^Vector [^Atom atom newval]
  (reset-vals__ atom newval))

(defn alter-meta!
  "Atomically sets the metadata for a namespace/var/atom to be:

  (apply f its-current-meta args)

  f must be free of side-effects"
  {:added "1.0"}
  [^Ref ref ^Callable f & args]
  (apply alter-meta__ ref f args))

(defn reset-meta!
  "Atomically resets the metadata for a namespace/var/atom"
  {:added "1.0"}
  [^Ref ref ^Map metadata-map] (reset-meta__ ref metadata-map))

(defn find-var
  "Returns the global var named by the namespace-qualified symbol, or
  nil if no var with that name."
  {:added "1.0"}
  ^Var [^Symbol sym] (find-var__ sym))

(defn comp
  "Takes a set of functions and returns a fn that is the composition
  of those fns.  The returned fn takes a variable number of args,
  applies the rightmost of fns to the args, the next
  fn (right-to-left) to the result, etc."
  {:added "1.0"}
  (^Fn [] identity)
  (^Fn [^Callable f] f)
  (^Fn [^Callable f ^Callable g]
   (fn
     ([] (f (g)))
     ([x] (f (g x)))
     ([x y] (f (g x y)))
     ([x y z] (f (g x y z)))
     ([x y z & args] (f (apply g x y z args)))))
  (^Fn [^Callable f ^Callable g ^Callable h]
   (fn
     ([] (f (g (h))))
     ([x] (f (g (h x))))
     ([x y] (f (g (h x y))))
     ([x y z] (f (g (h x y z))))
     ([x y z & args] (f (g (apply h x y z args))))))
  (^Fn [^Callable f1 ^Callable f2 ^Callable f3 & fs]
   (let [fs (reverse (list* f1 f2 f3 fs))]
     (fn [& args]
       (loop [ret (apply (first fs) args) fs (next fs)]
         (if fs
           (recur ((first fs) ret) (next fs))
           ret))))))

(defn juxt
  "Takes a set of functions and returns a fn that is the juxtaposition
  of those fns.  The returned fn takes a variable number of args, and
  returns a vector containing the result of applying each fn to the
  args (left-to-right).
  ((juxt a b c) x) => [(a x) (b x) (c x)]"
  {:added "1.0"}
  (^Fn [^Callable f]
   (fn
     ([] [(f)])
     ([x] [(f x)])
     ([x y] [(f x y)])
     ([x y z] [(f x y z)])
     ([x y z & args] [(apply f x y z args)])))
  (^Fn [^Callable f ^Callable g]
   (fn
     ([] [(f) (g)])
     ([x] [(f x) (g x)])
     ([x y] [(f x y) (g x y)])
     ([x y z] [(f x y z) (g x y z)])
     ([x y z & args] [(apply f x y z args) (apply g x y z args)])))
  (^Fn [^Callable f ^Callable g ^Callable h]
   (fn
     ([] [(f) (g) (h)])
     ([x] [(f x) (g x) (h x)])
     ([x y] [(f x y) (g x y) (h x y)])
     ([x y z] [(f x y z) (g x y z) (h x y z)])
     ([x y z & args] [(apply f x y z args) (apply g x y z args) (apply h x y z args)])))
  (^Fn [^Callable f ^Callable g ^Callable h & fs]
   (let [fs (list* f g h fs)]
     (fn
       ([] (reduce #(conj %1 (%2)) [] fs))
       ([x] (reduce #(conj %1 (%2 x)) [] fs))
       ([x y] (reduce #(conj %1 (%2 x y)) [] fs))
       ([x y z] (reduce #(conj %1 (%2 x y z)) [] fs))
       ([x y z & args] (reduce #(conj %1 (apply %2 x y z args)) [] fs))))))

(defn partial
  "Takes a function f and fewer than the normal arguments to f, and
  returns a fn that takes a variable number of additional args. When
  called, the returned function calls f with args + additional args."
  {:added "1.0"}
  (^Fn [^Callable f] f)
  (^Fn [^Callable f arg1]
   (fn [& args] (apply f arg1 args)))
  (^Fn [^Callable f arg1 arg2]
   (fn [& args] (apply f arg1 arg2 args)))
  (^Fn [^Callable f arg1 arg2 arg3]
   (fn [& args] (apply f arg1 arg2 arg3 args)))
  (^Fn [^Callable f arg1 arg2 arg3 & more]
   (fn [& args] (apply f arg1 arg2 arg3 (concat more args)))))

(defn sequence
  "Coerces coll to a (possibly empty) sequence, if it is not already
  one. Will not force a lazy seq. (sequence nil) yields ()"
  {:added "1.0"}
  ;; TODO: types (Seq or Seqable)
  ^Seq [coll]
  (if (seq? coll)
    coll
    (or (seq coll) ())))

(defn every?
  "Returns true if (pred x) is logical true for every x in coll, else
  false."
  {:added "1.0"}
  ^Boolean [^Callable pred ^Seqable coll]
  (cond
    (nil? (seq coll)) true
    (pred (first coll)) (recur pred (next coll))
    :else false))

(def
  ^{:tag Boolean
    :doc "Returns false if (pred x) is logical true for every x in
         coll, else true."
         :arglists '([pred coll])
         :added "1.0"}
  not-every? (comp not every?))

(defn some
  "Returns the first logical true value of (pred x) for any x in coll,
  else nil.  One common idiom is to use a set as pred, for example
  this will return :fred if :fred is in the sequence, otherwise nil:
  (some #{:fred} coll)"
  {:added "1.0"}
  [^Callable pred ^Seqable coll]
  (when (seq coll)
    (or (pred (first coll)) (recur pred (next coll)))))

(def
  ^{:tag Boolean
    :doc "Returns false if (pred x) is logical true for any x in coll,
         else true."
         :arglists '([pred coll])
         :added "1.0"}
  not-any? (comp not some))

(defn map
  "Returns a lazy sequence consisting of the result of applying f to the
  set of first items of each coll, followed by applying f to the set
  of second items in each coll, until any one of the colls is
  exhausted.  Any remaining items in other colls are ignored. Function
  f should accept number-of-colls arguments."
  {:added "1.0"}
  (^Seq [^Callable f ^Seqable coll]
   (lazy-seq
    (when-let [s (seq coll)]
      (cons (f (first s)) (map f (rest s))))))
  (^Seq [^Callable f ^Seqable c1 ^Seqable c2]
   (lazy-seq
    (let [s1 (seq c1) s2 (seq c2)]
      (when (and s1 s2)
        (cons (f (first s1) (first s2))
              (map f (rest s1) (rest s2)))))))
  (^Seq [^Callable f ^Seqable c1 ^Seqable c2 ^Seqable c3]
   (lazy-seq
    (let [s1 (seq c1) s2 (seq c2) s3 (seq c3)]
      (when (and  s1 s2 s3)
        (cons (f (first s1) (first s2) (first s3))
              (map f (rest s1) (rest s2) (rest s3)))))))
  (^Seq [^Callable f ^Seqable c1 ^Seqable c2 ^Seqable c3 & colls]
   (let [step (fn step [cs]
                (lazy-seq
                 (let [ss (map seq cs)]
                   (when (every? identity ss)
                     (cons (map first ss) (step (map rest ss)))))))]
     (map #(apply f %) (step (conj colls c3 c2 c1))))))

(defn mapcat
  "Returns the result of applying concat to the result of applying map
  to f and colls.  Thus function f should return a collection."
  {:added "1.0"}
  ^Seq [^Callable f & colls]
  (apply concat (apply map f colls)))

(defn filter
  "Returns a lazy sequence of the items in coll for which
  (pred item) returns true. pred must be free of side-effects."
  {:added "1.0"}
  (^Seq [^Callable pred ^Seqable coll]
   (lazy-seq
    (when-let [s (seq coll)]
      (let [f (first s) r (rest s)]
        (if (pred f)
          (cons f (filter pred r))
          (filter pred r)))))))

(defn remove
  "Returns a lazy sequence of the items in coll for which
  (pred item) returns false. pred must be free of side-effects."
  {:added "1.0"}
  ^Seq [^Callable pred ^Seqable coll]
  (filter (complement pred) coll))

(defn take
  "Returns a lazy sequence of the first n items in coll, or all items if
  there are fewer than n."
  {:added "1.0"}
  ^Seq [^Number n ^Seqable coll]
  (lazy-seq
   (when (pos? n)
     (when-let [s (seq coll)]
       (cons (first s) (take (dec n) (rest s)))))))

(defn take-while
  "Returns a lazy sequence of successive items from coll while
  (pred item) returns true. pred must be free of side-effects."
  {:added "1.0"}
  ^Seq [^Callable pred ^Seqable coll]
  (lazy-seq
   (when-let [s (seq coll)]
     (when (pred (first s))
       (cons (first s) (take-while pred (rest s)))))))

(defn drop
  "Returns a lazy sequence of all but the first n items in coll."
  {:added "1.0"}
  ^Seq [^Number n ^Seqable coll]
  (let [step (fn [n coll]
               (let [s (seq coll)]
                 (if (and (pos? n) s)
                   (recur (dec n) (rest s))
                   s)))]
    (lazy-seq (step n coll))))

(defn drop-last
  "Return a lazy sequence of all but the last n (default 1) items in coll"
  {:added "1.0"}
  (^Seq [^Seqable s] (drop-last 1 s))
  (^Seq [^Number n ^Seqable s] (map (fn [x _] x) s (drop n s))))

(defn take-last
  "Returns a seq of the last n items in coll.  Depending on the type
  of coll may be no better than linear time.  For vectors, see also subvec."
  {:added "1.0"}
  ^Seq [^Number n ^Seqable coll]
  (loop [s (seq coll) lead (seq (drop n coll))]
    (if lead
      (recur (next s) (next lead))
      s)))

(defn drop-while
  "Returns a lazy sequence of the items in coll starting from the first
  item for which (pred item) returns logical false."
  {:added "1.0"}
  ^Seq [^Callable pred ^Seqable coll]
  (let [step (fn [pred coll]
               (let [s (seq coll)]
                 (if (and s (pred (first s)))
                   (recur pred (rest s))
                   s)))]
    (lazy-seq (step pred coll))))

(defn cycle
  "Returns a lazy (infinite!) sequence of repetitions of the items in coll."
  {:added "1.0"}
  ^Seq [^Seqable coll]
  (lazy-seq
   (when-let [s (seq coll)]
     (concat s (cycle s)))))

(defn split-at
  "Returns a vector of [(take n coll) (drop n coll)]"
  {:added "1.0"}
  ^Vector [^Number n ^Seqable coll]
  [(take n coll) (drop n coll)])

(defn split-with
  "Returns a vector of [(take-while pred coll) (drop-while pred coll)]"
  {:added "1.0"}
  ^Vector [^Callable pred ^Seqable coll]
  [(take-while pred coll) (drop-while pred coll)])

(defn repeat
  "Returns a lazy (infinite!, or length n if supplied) sequence of xs."
  {:added "1.0"}
  (^Seq [x] (lazy-seq (cons x (repeat x))))
  (^Seq [^Number n x] (take n (repeat x))))

(defn iterate
  "Returns a lazy sequence of x, (f x), (f (f x)) etc. f must be free of side-effects"
  {:added "1.0"}
  ^Seq [^Callable f x] (cons x (lazy-seq (iterate f (f x)))))

(defn range
  "Returns a lazy seq of nums from start (inclusive) to end
  (exclusive), by step, where start defaults to 0, step to 1, and end to
  infinity. When step is equal to 0, returns an infinite sequence of
  start. When start is equal to end, returns empty list."
  {:added "1.0"}
  (^Seq [] (iterate inc 0))
  (^Seq [^Number end] (range 0 end 1))
  (^Seq [^Number start ^Number end] (range start end 1))
  (^Seq [^Number start ^Number end ^Number step]
   (lazy-seq
    (let [comp (cond
                 (or (zero? step) (= start end)) not=
                 (pos? step) <
                 (neg? step) >)]
      (if (comp start end)
        (cons start (range (+ start step) end step))
        ())))))

(defn merge
  "Returns a map that consists of the rest of the maps conj-ed onto
  the first.  If a key occurs in more than one map, the mapping from
  the latter (left-to-right) will be the mapping in the result."
  {:added "1.0"}
  ^Map [& maps]
  (when (some identity maps)
    (reduce #(conj (or %1 {}) %2) maps)))

(defn merge-with
  "Returns a map that consists of the rest of the maps conj-ed onto
  the first.  If a key occurs in more than one map, the mapping(s)
  from the latter (left-to-right) will be combined with the mapping in
  the result by calling (f val-in-result val-in-latter)."
  {:added "1.0"}
  ^Map [^Callable f & maps]
  (when (some identity maps)
    (let [merge-entry (fn [m e]
                        (let [k (key e) v (val e)]
                          (if (contains? m k)
                            (assoc m k (f (get m k) v))
                            (assoc m k v))))
          merge2 (fn [m1 m2]
                   (reduce merge-entry (or m1 {}) (seq m2)))]
      (reduce merge2 maps))))

(defn zipmap
  "Returns a map with the keys mapped to the corresponding vals."
  {:added "1.0"}
  ^Map [^Seqable keys ^Seqable vals]
  (loop [map {}
         ks (seq keys)
         vs (seq vals)]
    (if (and ks vs)
      (recur (assoc map (first ks) (first vs))
             (next ks)
             (next vs))
      map)))

(defn ^:private line-seq*
  [rdr]
  (when-let [line (reader-read-line__ rdr)]
    (cons line (lazy-seq (line-seq* rdr)))))

(defn line-seq
  "Returns the lines of text from rdr as a lazy sequence of strings.
  rdr must be File or BufferedReader."
  {:added "1.0"}
  [rdr]
  (cond
    (instance? BufferedReader rdr)
    (line-seq* rdr)

    :else
    (line-seq* (buffered-reader__ rdr))))

(defmacro declare
  "defs the supplied var names with no bindings, useful for making forward declarations."
  {:added "1.0"}
  [& names] `(do nil nil ~@(map #(list 'def (vary-meta % assoc :declared true)) names)))

(defn sort
  "Returns a sorted sequence of the items in coll. If no comparator is
  supplied, uses compare."
  {:added "1.0"}
  (^Seq [^Seqable coll]
   (sort compare coll))
  (^Seq [^Comparator comp ^Seqable coll]
   (sort__ comp coll)))

(defn sort-by
  "Returns a sorted sequence of the items in coll, where the sort
  order is determined by comparing (keyfn item).  If no comparator is
  supplied, uses compare."
  {:added "1.0"}
  (^Seq [^Callable keyfn ^Seqable coll]
   (sort-by keyfn compare coll))
  (^Seq [^Callable keyfn ^Comparator comp ^Seqable coll]
   (sort (fn [x y] (comp (keyfn x) (keyfn y))) coll)))

(defn dorun
  "When lazy sequences are produced via functions that have side
  effects, any effects other than those needed to produce the first
  element in the seq do not occur until the seq is consumed. dorun can
  be used to force any effects. Walks through the successive nexts of
  the seq, does not retain the head and returns nil."
  {:added "1.0"}
  (^Nil [^Seqable coll]
   (when (seq coll)
     (recur (next coll))))
  (^Nil [^Number n ^Seqable coll]
   (when (and (seq coll) (pos? n))
     (recur (dec n) (next coll)))))

(defn doall
  "When lazy sequences are produced via functions that have side
  effects, any effects other than those needed to produce the first
  element in the seq do not occur until the seq is consumed. doall can
  be used to force any effects. Walks through the successive nexts of
  the seq, retains the head and returns it, thus causing the entire
  seq to reside in memory at one time."
  {:added "1.0"}
  (^Seq [^Seqable coll]
   (dorun coll)
   coll)
  (^Seq [^Number n ^Seqable coll]
   (dorun n coll)
   coll))

(defn nthnext
  "Returns the nth next of coll, (seq coll) when n is 0."
  {:added "1.0"}
  ^Seq [^Seqable coll ^Number n]
  (loop [n n xs (seq coll)]
    (if (and xs (pos? n))
      (recur (dec n) (next xs))
      xs)))

(defn nthrest
  "Returns the nth rest of coll, coll when n is 0."
  {:added "1.0"}
  ^Seq [^Seqable coll ^Number n]
  (loop [n n xs coll]
    (if (and (pos? n) (seq xs))
      (recur (dec n) (rest xs))
      xs)))

(defn partition
  "Returns a lazy sequence of lists of n items each, at offsets step
  apart. If step is not supplied, defaults to n, i.e. the partitions
  do not overlap. If a pad collection is supplied, use its elements as
  necessary to complete last partition upto n items. In case there are
  not enough padding elements, return a partition with less than n items."
  {:added "1.0"}
  (^Seq [^Number n ^Seqable coll]
   (partition n n coll))
  (^Seq [^Number n ^Number step ^Seqable coll]
   (lazy-seq
    (when-let [s (seq coll)]
      (let [p (doall (take n s))]
        (when (= n (count p))
          (cons p (partition n step (nthrest s step))))))))
  (^Seq [^Number n ^Number step ^Seqable pad ^Seqable coll]
   (lazy-seq
    (when-let [s (seq coll)]
      (let [p (doall (take n s))]
        (if (= n (count p))
          (cons p (partition n step pad (nthrest s step)))
          (list (take n (concat p pad)))))))))

(defn eval
  "Evaluates the form data structure (not text!) and returns the result."
  {:added "1.0"}
  [form] (eval__ form))

(defmacro doseq
  "Repeatedly executes body (presumably for side-effects) with
  bindings and filtering as provided by \"for\".  Does not retain
  the head of the sequence. Returns nil."
  {:added "1.0"}
  ^Nil [seq-exprs & body]
  (assert-args
   (vector? seq-exprs) "a vector for its binding"
   (even? (count seq-exprs)) "an even number of forms in binding vector")
  (let [b (if (> (count body) 1)
            `(do ~@body)
            (first body))
        step (fn step [recform exprs]
               (if-not exprs
                 [true b]
                 (let [k (first exprs)
                       v (second exprs)
                       seqsym (when-not (keyword? k) (gensym))
                       recform (if (keyword? k) recform `(recur (next ~seqsym)))
                       steppair (step recform (nnext exprs))
                       needrec (steppair 0)
                       subform (steppair 1)]
                   (cond
                     (= k :let) [needrec `(let ~v ~subform)]
                     (= k :while) [false `(when ~v
                                            ~subform
                                            ~@(when needrec [recform]))]
                     (= k :when) [false `(if ~v
                                           (do nil
                                             ~subform
                                             ~@(when needrec [recform]))
                                           ~recform)]
                     :else [true `(loop [~seqsym (seq ~v)]
                                    (when ~seqsym
                                      (let [~k (first ~seqsym)]
                                        ~subform
                                        ~@(when needrec [recform]))))]))))]
    (nth (step nil (seq seq-exprs)) 1)))

(defmacro dotimes
  "bindings => name n

  Repeatedly executes body (presumably for side-effects) with name
  bound to integers from 0 through n-1."
  {:added "1.0"}
  [bindings & body]
  (assert-args
   (vector? bindings) "a vector for its binding"
   (= 2 (count bindings)) "exactly 2 forms in binding vector")
  (let [i (first bindings)
        n (second bindings)]
    `(let [n# (int__ ~n)]
       (loop [~i 0]
         (when (< ~i n#)
           ~@body
           (recur (inc ~i)))))))

(defn num
  "Coerce to Number"
  {:added "1.0"}
  ^Number [^Number x] (num__ x))

(defn double
  "Coerce to double"
  {:added "1.0"}
  ^Double [^Number x] (double__ x))

(defn char
  "Coerce to char"
  {:added "1.0"}
  ;; TODO: types (Char or Number)
  ^Char [x] (char__ x))

(defn boolean
  "Coerce to boolean"
  {:added "1.0"}
  ^Boolean [x] (boolean__ x))

(defn number?
  "Returns true if x is a Number"
  {:added "1.0"}
  ^Boolean [x]
  (instance? Number x))

(defn mod
  "Modulus of num and div. Truncates toward negative infinity."
  {:added "1.0"}
  ^Number [^Number num ^Number div]
  (let [m (rem num div)]
    (if (or (zero? m) (= (pos? num) (pos? div)))
      m
      (+ m div))))

(defn ratio?
  "Returns true if n is a Ratio"
  {:added "1.0"}
  ^Boolean [n] (instance? Ratio n))

(defn numerator
  "Returns the numerator part of a Ratio."
  {:added "1.0"}
  ^Number [^Ratio r]
  (numerator__ r))

(defn denominator
  "Returns the denominator part of a Ratio."
  {:added "1.0"}
  ^Number [^Ratio r]
  (denominator__ r))

(defn bigfloat?
  "Returns true if n is a BigFloat"
  {:added "1.0"}
  ^Boolean [n] (instance? BigFloat n))

(defn float?
  "Returns true if n is a floating point number"
  {:added "1.0"}
  ^Boolean [n]
  (instance? Double n))

(defn rational?
  "Returns true if n is a rational number"
  {:added "1.0"}
  ^Boolean [n]
  (or (integer? n) (ratio? n)))

(defn bigint
  "Coerce to BigInt"
  {:added "1.0"}
  ^BigInt [x]
  ;; TODO: types (Number or String)
  (bigint__ x))

(defn bigfloat
  "Coerce to BigFloat"
  {:added "1.0"}
  ^BigFloat [x]
  ;; TODO: types (Number or String)
  (bigfloat__ x))

(def
  ^{:arglists '([& args])
    :tag Nil
    :doc "Prints the object(s) to the output stream that is the current value
         of *out*.  Prints the object(s), separated by spaces if there is
         more than one.  By default, pr and prn print in a way that objects
         can be read by the reader"
         :added "1.0"}
  pr pr__)

(defn pprint
  "Pretty prints x to the output stream that is the current value of *out*."
  {:added "1.0"}
  ^Nil [x]
  (pprint__ x))

(defn newline
  "Writes a platform-specific newline to *out*"
  {:added "1.0"}
  ^Nil []
  (newline__))

(defn flush
  "Flushes the output stream that is the current value of
  *out*"
  {:added "1.0"}
  ^Nil []
  (flush__ *out*)
  nil)

(def
  ^{:doc
    "When set to true, output will be flushed whenever a newline is printed.

    Defaults to true."
    :added "1.0"}
  *flush-on-newline* true)

(defn prn
  "Same as pr followed by (newline). Observes *flush-on-newline*"
  {:added "1.0"}
  ^Nil [& more]
  (apply pr more)
  (newline)
  (when *flush-on-newline*
    (flush)))

(defn print
  "Prints the object(s) to the output stream that is the current value
  of *out*.  print and println produce output for human consumption."
  {:added "1.0"}
  ^Nil [& more]
  (binding [*print-readably* nil]
    (apply pr more)))

(defn println
  "Same as print followed by (newline)"
  {:added "1.0"}
  ^Nil [& more]
  (binding [*print-readably* nil]
    (apply prn more)))

(defn read
  "Reads the next object from reader (defaults to *in*)"
  {:added "1.0"}
  ([] (read *in*))
  ([reader] (read__ reader)))

(def
  ^{:arglists '([])
    :doc "Reads the next line from *in*. Returns nil if an error (such as EOF) is detected."
    :added "1.0"
    :tag String}
  read-line read-line__)

(def
  ^{:arglists '([s])
    :doc "Reads one object from the string s."
    :added "1.0"}
  read-string read-string__)

(defn subvec
  "Returns a persistent vector of the items in vector from
  start (inclusive) to end (exclusive).  If end is not supplied,
  defaults to (count vector). This operation is O(1) and very fast, as
  the resulting vector shares structure with the original and no
  trimming is done."
  {:added "1.0"}
  (^Vector [^Vector v ^Number start]
   (subvec v start (count v)))
  (^Vector [^Vector v ^Number start ^Number end]
   (subvec__ v start end)))

(defmacro doto
  "Evaluates x then calls all of the methods and functions with the
  value of x supplied at the front of the given arguments.  The forms
  are evaluated in order.  Returns x."
  {:added "1.0"}
  [x & forms]
  (let [gx (gensym)]
    `(let [~gx ~x]
       ~@(map (fn [f]
                (if (seq? f)
                  `(~(first f) ~gx ~@(next f))
                  `(~f ~gx)))
              forms)
       ~gx)))

(defmacro time
  "Evaluates expr and prints the time it took.  Returns the value of expr."
  {:added "1.0"}
  [expr]
  `(let [start# (nano-time__)
         ret# ~expr]
     (prn (str "Elapsed time: " (/ (double (- (nano-time__) start#)) 1000000.0) " msecs"))
     ret#))

(defn macroexpand-1
  "If form represents a macro form, returns its expansion, else returns form."
  {:added "1.0"}
  [form]
  (macroexpand-1__ form))

(defn macroexpand
  "Repeatedly calls macroexpand-1 on form until it no longer
  represents a macro form, then returns it.  Note neither
  macroexpand-1 nor macroexpand expand macros in subforms."
  {:added "1.0"}
  [form]
  (let [ex (macroexpand-1 form)]
    (if (identical? ex form)
      form
      (macroexpand ex))))

(defn load-string
  "Sequentially read and evaluate the set of forms contained in the
  string"
  {:added "1.0"}
  [^String s]
  (load-string__ s))

(defn set?
  "Returns true if x implements Set"
  {:added "1.0"}
  ^Boolean [x] (instance? Set x))

(defn set
  "Returns a set of the distinct elements of coll."
  {:added "1.0"}
  ^MapSet [^Seqable coll]
  (if (set? coll)
    (with-meta coll nil)
    (reduce conj #{} coll)))

(defn ^:private filter-key
  [keyfn pred amap]
  (loop [ret {} es (seq amap)]
    (if es
      (if (pred (keyfn (first es)))
        (recur (assoc ret (key (first es)) (val (first es))) (next es))
        (recur ret (next es)))
      ret)))

(defn find-ns
  "Returns the namespace named by the symbol or nil if it doesn't exist."
  {:added "1.0"}
  ^Namespace [^Symbol sym] (find-ns__ sym))

(defn create-ns
  "Create a new namespace named by the symbol if one doesn't already
  exist, returns it or the already-existing namespace of the same
  name."
  {:added "1.0"}
  ^Namespace [^Symbol sym] (create-ns__ sym))

(defn remove-ns
  "Removes the namespace named by the symbol. Use with caution.
  Cannot be used to remove the clojure namespace."
  {:added "1.0"}
  ^Namespace [^Symbol sym] (remove-ns__ sym))

(defn all-ns
  "Returns a sequence of all namespaces."
  {:added "1.0"}
  ^Seq [] (all-ns__))

(defn the-ns
  "If passed a namespace, returns it. Else, when passed a symbol,
  returns the namespace named by it, throwing an exception if not
  found."
  {:added "1.0"}
  ^Namespace [x]
  ;; TODO: types (Namespace or Symbol)
  (if (instance? Namespace x)
    x
    (or (find-ns x)
        (if *linter-mode*
          (do
            (println-linter__ (ex-info (str "No namespace: " x " found")
                                  {:form x :_prefix "Parse warning"}))
            (create-ns__ x))
          (throw (ex-info (str "No namespace: " x " found") {:form x}))))))

(defn ns-name
  "Returns the name of the namespace, a symbol."
  {:added "1.0"}
  ^Symbol [ns]
  ;; TODO: types (Namespace or Symbol)
  (ns-name__ (the-ns ns)))

(defn ns-map
  "Returns a map of all the mappings for the namespace."
  {:added "1.0"}
  ^Map [ns]
  ;; TODO: types (Namespace or Symbol)
  (ns-map__ (the-ns ns)))

(defn ns-unmap
  "Removes the mappings for the symbol from the namespace."
  {:added "1.0"}
  ^Nil [ns ^Symbol sym]
  ;; TODO: types (Namespace or Symbol)
  (ns-unmap__ (the-ns ns) sym))

(defn ^:private public?
  [v]
  (not (:private (meta v))))

(defn ns-publics
  "Returns a map of the public intern mappings for the namespace."
  {:added "1.0"}
  ^Map [ns]
  ;; TODO: types (Namespace or Symbol)
  (let [ns (the-ns ns)]
    (filter-key val (fn [^Var v] (and (instance? Var v)
                                      (= ns (var-ns__ v))
                                      (public? v)))
                (ns-map ns))))

(defn ns-interns
  "Returns a map of the intern mappings for the namespace."
  {:added "1.0"}
  ^Map [ns]
  ;; TODO: types (Namespace or Symbol)
  (let [ns (the-ns ns)]
    (filter-key val (fn [^Var v] (and (instance? Var v)
                                      (= ns (var-ns__ v))))
                (ns-map ns))))

(defn ^:private get-refer-opt
  [opts]
  (:refer opts))

(defn intern
  "Finds or creates a var named by the symbol name in the namespace
  ns (which can be a symbol or a namespace), setting its root binding
  to val if supplied. The namespace must exist. The var will adopt any
  metadata from the name symbol.  Returns the var."
  {:added "1.0"}
  (^Var [ns ^Symbol name]
   ;; TODO: types (Namespace or Symbol)
   (let [v (intern__ (the-ns ns) name)]
     (when (meta name) (set-meta__ v (meta name)))
     v))
  (^Var [ns ^Symbol name val]
    ;; TODO: types (Namespace or Symbol)
   (let [v (intern__ (the-ns ns) name val)]
     (when (meta name) (set-meta__ v (meta name)))
     v)))

(defn refer
  "refers to all public vars of ns, subject to filters.
  filters can include at most one each of:

  :exclude list-of-symbols
  :only list-of-symbols
  :rename map-of-fromsymbol-tosymbol

  For each public interned var in the namespace named by the symbol,
  adds a mapping from the name of the var to the var to the current
  namespace.  Throws an exception if name is already mapped to
  something else in the current namespace. Filters can be used to
  select a subset, via inclusion or exclusion, or to provide a mapping
  to a symbol different from the var's name, in order to prevent
  clashes. Use :use in the ns macro in preference to calling this directly."
  {:added "1.0"}
  ^Nil [^Symbol ns-sym & filters]
  (let [ns (or (find-ns ns-sym) (throw (ex-info (str "No namespace: " ns-sym) {:form ns-sym})))
        fs (apply hash-map filters)
        nspublics (ns-publics ns)
        rename (or (:rename fs) {})
        exclude (set (:exclude fs))
        to-do (if (= :all (get-refer-opt fs))
                (keys nspublics)
                (or (get-refer-opt fs) (:only fs) (keys nspublics)))]
    (when (and to-do (not (instance? Sequential to-do)))
      (throw (ex-info ":only/:refer value must be a sequential collection of symbols" {:form ns-sym})))
    (doseq [sym to-do]
      (when-not (exclude sym)
        (let [v (nspublics sym)]
          (when-not (or *linter-mode* v)
            (throw (ex-info
                    (if (get (ns-interns ns) sym)
                      (str sym " is not public")
                      (str sym " does not exist"))
                    {:form ns-sym})))
          (refer__ *ns* (or (rename sym) sym)
                   (or v (intern-fake-var__
                          ns-sym
                          sym
                          (first (filter #(= sym %) (:refer-macros fs)))))))))))

(defn ns-refers
  "Returns a map of the refer mappings for the namespace."
  {:added "1.0"}
  ^Map [ns]
  ;; TODO: types (Namespace or Symbol)
  (let [ns (the-ns ns)]
    (filter-key val (fn [^Var v] (and (instance? Var v)
                                      (not= ns (var-ns__ v))))
                (ns-map ns))))

(defn alias
  "Add an alias in the current namespace to another
  namespace. Arguments are two symbols: the alias to be used, and
  the symbolic name of the target namespace. Use :as in the ns macro in preference
  to calling this directly."
  {:added "1.0"}
  ^Nil [^Symbol alias namespace-sym]
  ;; TODO: types (Namespace or Symbol)
  (alias__ *ns* alias (the-ns namespace-sym)))

(defn ns-aliases
  "Returns a map of the aliases for the namespace."
  {:added "1.0"}
  ^Map [ns]
  ;; TODO: types (Namespace or Symbol)
  (ns-aliases__ (the-ns ns)))

(defn ns-unalias
  "Removes the alias for the symbol from the namespace."
  {:added "1.0"}
  ^Nil [ns ^Symbol sym]
  ;; TODO: types (Namespace or Symbol)
  (ns-unalias__ (the-ns ns) sym))

(defn take-nth
  "Returns a lazy seq of every nth item in coll."
  {:added "1.0"}
  ^Seq [^Number n ^Seqable coll]
  (lazy-seq
   (when-let [s (seq coll)]
     (cons (first s) (take-nth n (drop n s))))))

(defn interleave
  "Returns a lazy seq of the first item in each coll, then the second etc."
  {:added "1.0"}
  (^Seq [] ())
  (^Seq [^Seqable c1] (lazy-seq c1))
  (^Seq [^Seqable c1 ^Seqable c2]
   (lazy-seq
    (let [s1 (seq c1) s2 (seq c2)]
      (when (and s1 s2)
        (cons (first s1) (cons (first s2)
                               (interleave (rest s1) (rest s2))))))))
  (^Seq [^Seqable c1 ^Seqable c2 & colls]
   (lazy-seq
    (let [ss (map seq (conj colls c2 c1))]
      (when (every? identity ss)
        (concat (map first ss) (apply interleave (map rest ss))))))))

(defn ns-resolve
  "Returns the var or type to which a symbol will be resolved in the
  namespace (unless found in the environment), else nil.  Note that
  if the symbol is fully qualified, the var/Type to which it resolves
  need not be present in the namespace."
  {:added "1.0"}
  ;; TODO: ns is namespace or symbol
  (^Var [ns ^Symbol sym]
    (ns-resolve ns nil sym))
  (^Var [ns ^Gettable env ^Symbol sym]
    (when-not (contains? env sym)
      (ns-resolve__ (the-ns ns) sym))))

(defn resolve
  "Same as (ns-resolve *ns* sym) or (ns-resolve *ns* env sym)"
  {:added "1.0"}
  (^Var [^Symbol sym] (ns-resolve *ns* sym))
  (^Var [^Gettable env ^Symbol sym] (ns-resolve *ns* env sym)))

(def
  ^{:arglists '([& keyvals])
    :doc "Constructs an array-map. If any keys are equal, they are handled as
         if by repeated uses of assoc."
         :added "1.0"
         :tag ArrayMap}
  array-map array-map__)

(defn ^:private make-mark-skip-unused__
  [rule-name]
  (fn [s]
    (if (and *linter-mode* (false? (rule-name (:rules *linter-config*))))
      (vary-meta s assoc :skip-unused true)
      s)))

(defn ^:private mark-skip-unused__
  [s]
  (if *linter-mode* (vary-meta s assoc :skip-unused true) s))

;redefine let and loop  with destructuring
(defn ^:private destructure [bindings]
  (let [mark-as (make-mark-skip-unused__ :unused-as)
        mark-keys (make-mark-skip-unused__ :unused-keys)
        bents (partition 2 bindings)
        pb (fn pb [bvec b v]
             (let [pvec
                   (fn [bvec b val]
                     (when (and *linter-mode* (not (seq b)))
                       (println-linter__ (ex-info "destructuring with no bindings"
                                                  {:form b :_prefix "Parse warning"})))
                     (let [gvec (gensym "vec__")
                           gseq (gensym "seq__")
                           gfirst (gensym "first__")
                           has-rest (some #{'&} b)]
                       (loop [ret (let [ret (conj bvec (mark-skip-unused__ gvec) val)]
                                    (if has-rest
                                      (conj ret gseq (list `seq gvec))
                                      ret))
                              n 0
                              bs b
                              seen-rest? false]
                         (if (seq bs)
                           (let [firstb (first bs)]
                             (cond
                               (= firstb '&) (recur (pb ret (second bs) gseq)
                                                    n
                                                    (nnext bs)
                                                    true)
                               (= firstb :as) (pb ret (mark-as (second bs)) gvec)
                               :else (if seen-rest?
                                       (throw (ex-info "Unsupported binding form, only :as can follow & parameter" {:form bindings}))
                                       (recur (pb (if has-rest
                                                    (conj ret
                                                          gfirst `(first ~gseq)
                                                          gseq `(next ~gseq))
                                                    ret)
                                                  firstb
                                                  (if has-rest
                                                    gfirst
                                                    (list `nth gvec n nil)))
                                              (inc n)
                                              (next bs)
                                              seen-rest?))))
                           ret))))
                   pmap
                   (fn [bvec b v]
                     (when (and *linter-mode* (not (seq b)))
                       (println-linter__ (ex-info "destructuring with no bindings"
                                                  {:form b :_prefix "Parse warning"})))
                     (let [gmap (gensym "map__")
                           gmapseq (with-meta gmap {:tag 'Seq})
                           defaults (:or b)]
                       (loop [ret (-> bvec (conj gmap) (conj v)
                                      (conj (mark-skip-unused__ gmap)) (conj `(if (seq? ~gmap) (apply array-map__ (seq ~gmapseq)) ~gmap))
                                      ((fn [ret]
                                         (if (:as b)
                                           (conj ret (mark-as (:as b)) gmap)
                                           ret))))
                              bes (let [transforms
                                        (reduce
                                         (fn [transforms mk]
                                           (if (keyword? mk)
                                             (let [mkns (namespace mk)
                                                   mkn (name mk)]
                                               (cond (= mkn "keys") (assoc transforms mk #(keyword (or mkns (namespace %)) (name %)))
                                                 (= mkn "syms") (assoc transforms mk #(list `quote (symbol (or mkns (namespace %)) (name %))))
                                                 (= mkn "strs") (assoc transforms mk str)
                                                 :else transforms))
                                             transforms))
                                         {}
                                         (keys b))]
                                    (reduce
                                     (fn [bes entry]
                                       (reduce #(assoc %1 %2 ((val entry) %2))
                                               (dissoc bes (key entry))
                                               ((key entry) bes)))
                                     (dissoc b :as :or)
                                     transforms))]
                         (if (seq bes)
                           (let [bb (mark-keys (key (first bes)))
                                 bk (val (first bes))
                                 local (if (instance? Named bb) (derive-info__ (with-meta (symbol nil (name bb)) (meta bb)) bb) bb)
                                 bv (if (contains? defaults local)
                                      (list `get gmap bk (defaults local))
                                      (list `get gmap bk))]
                             (recur (if (ident? bb)
                                      (-> ret (conj local bv))
                                      (pb ret bb bv))
                                    (next bes)))
                           ret))))]
               (cond
                 (symbol? b) (-> bvec (conj b) (conj v))
                 (vector? b) (pvec bvec b v)
                 (map? b) (pmap bvec b v)
                 :else (throw (ex-info (str "Unsupported binding form: " b) {:form b})))))
        process-entry (fn [bvec b] (pb bvec (first b) (second b)))]
    (if (every? symbol? (map first bents))
      bindings
      (with-meta
        (reduce process-entry [] bents)
        (meta bindings)))))

(defmacro let
  "binding => binding-form init-expr

  Evaluates the exprs in a lexical context in which the symbols in
  the binding-forms are bound to their respective init-exprs or parts
  therein."
  {:added "1.0", :special-form true, :forms '[(let [bindings*] exprs*)]}
  [bindings & body]
  (assert-args
   (vector? bindings) "a vector for its binding"
   (even? (count bindings)) "an even number of forms in binding vector")
  `(let* ~(destructure bindings) ~@body))

(defn ^{:private true}
  maybe-destructured
  [params body]
  (if (every? symbol? params)
    (cons params body)
    (loop [params params
           new-params (with-meta [] (meta params))
           lets []]
      (if params
        (if (symbol? (first params))
          (recur (next params) (conj new-params (first params)) lets)
          (let [gparam (gensym "p__")]
            (recur (next params) (conj new-params gparam)
                   (-> lets (conj (first params)) (conj gparam)))))
        `(~new-params
          (let ~lets
            ~@body))))))

;redefine fn with destructuring and pre/post conditions
(defmacro fn
  "params => positional-params* , or positional-params* & next-param
  positional-param => binding-form
  next-param => binding-form
  name => symbol

  Defines a function"
  {:added "1.0", :special-form true,
   :forms '[(fn name? [params* ] exprs*) (fn name? ([params* ] exprs*)+)]}
  [& sigs]
  (let [name (if (symbol? (first sigs)) (first sigs) nil)
        sigs (if name (next sigs) sigs)
        sigs (if (vector? (first sigs))
               (list sigs)
               (if (seq? (first sigs))
                 sigs
                 ;; Assume single arity syntax
                 (throw (ex-info
                         (if (seq sigs)
                           (str "Parameter declaration \""
                                (first sigs)
                                "\" must be a vector")
                           (str "Parameter declaration missing"))
                         {:form &form}))))
        psig (fn* [sig]
                  ;; Ensure correct type before destructuring sig
                  (when (not (seq? sig))
                    (throw (ex-info
                            (str "Invalid signature: \"" sig
                                 "\" must be a list")
                            {:form &form})))
                  (let [[params & body] sig
                        _ (when (not (vector? params))
                            (throw (ex-info
                                    (if (seq? (first sigs))
                                      (str "Parameter declaration \"" params
                                           "\" must be a vector")
                                      (str "Invalid signature: \"" sig
                                           "\" must be a list"))
                                    {:form &form})))
                        conds (when (and (next body) (map? (first body)))
                                (first body))
                        body (if conds (next body) body)
                        conds (or conds (meta params))
                        pre (:pre conds)
                        post (:post conds)
                        body (if post
                               `((let [~'% ~(if (< 1 (count body))
                                              `(do ~@body)
                                              (first body))]
                                   ~@(map (fn* [c] `(assert ~c)) post)
                                   ~'%))
                               body)
                        body (if pre
                               (concat (map (fn* [c] `(assert ~c)) pre)
                                       body)
                               body)]
                    (derive-info__ (maybe-destructured params body) sig)))
        new-sigs (map psig sigs)]
    (with-meta
      (if name
        (list* 'fn* name new-sigs)
        (cons 'fn* new-sigs))
      (meta &form))))

(defmacro loop
  "Evaluates the exprs in a lexical context in which the symbols in
  the binding-forms are bound to their respective init-exprs or parts
  therein. Acts as a recur target."
  {:added "1.0", :special-form true, :forms '[(loop [bindings*] exprs*)]}
  [bindings & body]
  (assert-args
   (vector? bindings) "a vector for its binding"
   (even? (count bindings)) "an even number of forms in binding vector")
  (let [db (destructure bindings)]
    (if (= db bindings)
      `(loop* ~bindings ~@body)
      (let [vs (take-nth 2 (drop 1 bindings))
            bs (take-nth 2 bindings)
            gs (map (fn [b] (if (symbol? b) b (gensym))) bs)
            bfs (with-meta
                  (reduce (fn [ret [b v g]]
                            (if (symbol? b)
                              (conj ret g v)
                              (conj ret g v (vary-meta b assoc :skip-unused true) g)))
                          [] (map vector bs vs gs))
                  {:skip-unused true})]
        `(let ~bfs
           (loop* ~(vec (interleave gs gs))
                  (let ~(vec (interleave bs gs))
                    ~@body)))))))

(defmacro when-first
  "bindings => x xs

  Roughly the same as (when (seq xs) (let [x (first xs)] body)) but xs is evaluated only once"
  {:added "1.0"}
  [bindings & body]
  (assert-args
   (vector? bindings) "a vector for its binding"
   (= 2 (count bindings)) "exactly 2 forms in binding vector")
  (let [[x xs] bindings]
    `(when-let [xs# (seq ~xs)]
       (let [~x (first xs#)]
         ~@body))))

(defmacro lazy-cat
  "Expands to code which yields a lazy sequence of the concatenation
  of the supplied colls.  Each coll expr is not evaluated until it is
  needed.

  (lazy-cat xs ys zs) === (concat (lazy-seq xs) (lazy-seq ys) (lazy-seq zs))"
  {:added "1.0"}
  [& colls]
  `(concat ~@(map #(list `lazy-seq %) colls)))

(defmacro for
  "List comprehension. Takes a vector of one or more
  binding-form/collection-expr pairs, each followed by zero or more
  modifiers, and yields a lazy sequence of evaluations of expr.
  Collections are iterated in a nested fashion, rightmost fastest,
  and nested coll-exprs can refer to bindings created in prior
  binding-forms.  Supported modifiers are: :let [binding-form expr ...],
  :while test, :when test.

  (take 100 (for [x (range 100000000) y (range 1000000) :while (< y x)]  [x y]))"
  {:added "1.0"}
  [seq-exprs body-expr]
  (assert-args
   (vector? seq-exprs) "a vector for its binding"
   (even? (count seq-exprs)) "an even number of forms in binding vector")
  (let [to-groups (fn [seq-exprs]
                    (reduce (fn [groups [k v]]
                              (if (keyword? k)
                                (conj (pop groups) (conj (peek groups) [k v]))
                                (conj groups [k v])))
                            [] (partition 2 seq-exprs)))
        emit-bind (fn emit-bind [[[bind expr & mod-pairs]
                                  & [[_ next-expr] :as next-groups]]]
                    (let [giter (gensym "iter__")
                          gxs (gensym "s__")
                          do-mod (fn do-mod [[[k v :as pair] & etc]]
                                   (cond
                                     (= k :let) `(let ~v ~(do-mod etc))
                                     (= k :while) `(when ~v ~(do-mod etc))
                                     (= k :when) `(if ~v
                                                    ~(do-mod etc)
                                                    (recur (rest ~gxs)))
                                     (keyword? k) (throw (ex-info (str "Invalid keyword " k " in \"for\" form") {:form k}))
                                     next-groups
                                     `(let [iterys# ~(emit-bind next-groups)
                                            fs# (seq (iterys# ~next-expr))]
                                        (if fs#
                                          (concat fs# (~giter (rest ~gxs)))
                                          (recur (rest ~gxs))))
                                     :else `(cons ~body-expr
                                                  (~giter (rest ~gxs)))))]
                      `(fn ~giter [~gxs]
                         (lazy-seq
                          (loop [~gxs ~gxs]
                            (when-first [~bind ~gxs]
                              ~(do-mod mod-pairs)))))))]
    `(let [iter# ~(emit-bind (to-groups seq-exprs))]
       (iter# ~(second seq-exprs)))))

(defmacro comment
  "Ignores body, yields nil"
  {:added "1.0"}
  [& body])

(defmacro with-out-str
  "Evaluates exprs in a context in which *out* is bound to a fresh
  Buffer.  Returns the string created by any nested printing
  calls."
  {:added "1.0"}
  [& body]
  `(binding [*out* (buffer__)]
     ~@body
     (str *out*)))

(defmacro with-in-str
  "Evaluates body in a context in which *in* is bound to a fresh
  Buffer initialized with the string s."
  {:added "1.0"}
  [s & body]
  `(binding [*in* (buffer__ ~s)]
     ~@body))

(defn pr-str
  "pr to a string, returning it"
  {:added "1.0"}
  ^String [& xs]
  (with-out-str
    (apply pr xs)))

(defn prn-str
  "prn to a string, returning it"
  {:added "1.0"}
  ^String [& xs]
  (with-out-str
    (apply prn xs)))

(defn print-str
  "print to a string, returning it"
  {:added "1.0"}
  ^String [& xs]
  (with-out-str
    (apply print xs)))

(defn println-str
  "println to a string, returning it"
  {:added "1.0"}
  ^String [& xs]
  (with-out-str
    (apply println xs)))

(defn pr-err
  "pr to *err*"
  {:added "1.0"}
  ^Nil [& xs]
  (binding [*out* *err*]
    (apply pr xs)))

(defn prn-err
  "prn to *err*"
  {:added "1.0"}
  ^Nil [& xs]
  (binding [*out* *err*]
    (apply prn xs)))

(defn print-err
  "print to *err*"
  {:added "1.0"}
  ^Nil [& xs]
  (binding [*out* *err*]
    (apply print xs)))

(defn println-err
  "println to *err*"
  {:added "1.0"}
  ^Nil [& xs]
  (binding [*out* *err*]
    (apply println xs)))

(defn ^:private println-linter__
  [& xs]
  (inc-problem-count__)
  (apply println-err xs))

(defn ex-data
  "Returns exception data (a map) if ex is an ExInfo.
  Otherwise returns nil."
  {:added "1.0"}
  ^Map [ex]
  (when (instance? ExInfo ex)
    (ex-data__ ex)))

(defn ex-cause
  "Returns the cause of ex if ex is an ExInfo.
  Otherwise returns nil."
  {:added "1.0"}
  ^Error [ex]
  (when (instance? ExInfo ex)
    (ex-cause__ ex)))

(defn ex-message
  "Returns the message attached to ex if ex is an ExInfo.
  Otherwise returns nil."
  {:added "1.0"}
  ^String [ex]
  (when (instance? Error ex)
    (ex-message__ ex)))

(defn hash
  "Returns the hash code of its argument."
  {:added "1.0"}
  ^Int [x] (hash__ x))

(defmacro assert
  "Evaluates expr and throws an exception if it does not evaluate to
  logical true."
  {:added "1.0"}
  ([x]
   (when *assert*
     `(when-not ~x
        (throw (ex-info (str "Assert failed: " '~x) {:form ~x})))))
  ([x message]
   (when *assert*
     `(when-not ~x
        (throw (ex-info (str "Assert failed: " ~message "\n" '~x) {:form ~x}))))))

(defn test
  "test [v] finds fn at key :test in var metadata and calls it,
  presuming failure will throw exception"
  {:added "1.0"}
  ^Keyword [v]
  (let [f (:test (meta v))]
    (if f
      (do (f) :ok)
      :no-test)))

(defn re-pattern
  "Returns an instance of Regex"
  {:added "1.0"}
  ;; TODO: Regex or String
  ^Regex [s]
  (if (instance? Regex s)
    s
    (regex__ s)))

(defn re-seq
  "Returns a sequence of successive matches of pattern in string"
  {:added "1.0"}
  ^Seq [^Regex re ^String s]
  (re-seq__ re s))

(defn re-find
  "Returns the leftmost regex match, if any, of string to pattern."
  {:added "1.0"}
  [^Regex re ^String s]
  (re-find__ re s))

(defn re-matches
  "Returns the match, if any, of string to pattern."
  {:added "1.0"}
  [^Regex re ^String s]
  (let [m (re-find re s)
        c (if (instance? String m)
            (count m)
            (count (first m)))]
    (when (= c (count s))
      m)))

(defn rand
  "Returns a random floating point number between 0 (inclusive) and
  n (default 1) (exclusive)."
  {:added "1.0"}
  (^Double [] (rand__))
  (^Double [^Number n] (* n (rand))))

(defn rand-int
  "Returns a random integer between 0 (inclusive) and n (exclusive)."
  {:added "1.0"}
  ^Int [^Number n] (int (rand n)))

(defmacro defn-
  "same as defn, yielding non-public def"
  {:added "1.0"}
  [name & decls]
  (list* `defn (with-meta name (assoc (meta name) :private true)) decls))

(defn tree-seq
  "Returns a lazy sequence of the nodes in a tree, via a depth-first walk.
  branch? must be a fn of one arg that returns true if passed a node
  that can have children (but may not).  children must be a fn of one
  arg that returns a sequence of the children. Will only be called on
  nodes for which branch? returns true. Root is the root node of the
  tree."
  {:added "1.0"}
  ^Seq [^Callable branch? ^Callable children root]
  (let [walk (fn walk [node]
               (lazy-seq
                (cons node
                      (when (branch? node)
                        (mapcat walk (children node))))))]
    (walk root)))

; TODO:
; (defn file-seq
;   "A tree seq on directory"
;   {:added "1.0"}
;   [dir]
;   (tree-seq
;    (fn [^java.io.File f] (. f (isDirectory)))
;    (fn [^java.io.File d] (seq (. d (listFiles))))
;    dir))

(defn xml-seq
  "A tree seq on the xml elements as per xml/parse"
  {:added "1.0"}
  ^Seq [root]
  (tree-seq
   (complement string?)
   (comp seq :content)
   root))

(defn special-symbol?
  "Returns true if s names a special form"
  {:added "1.0"}
  ^Boolean [s]
  (special-symbol?__ s))

(defn var?
  "Returns true if v is of type Var"
  {:added "1.0"}
  ^Boolean [v] (instance? Var v))

(defn subs
  "Returns the substring of s beginning at start inclusive, and ending
  at end (defaults to length of string), exclusive."
  {:added "1.0"}
  (^String [^String s ^Number start] (subs__ s start))
  (^String [^String s ^Number start ^Number end] (subs__ s start end)))

(defn max-key
  "Returns the x for which (k x), a number, is greatest."
  {:added "1.0"}
  ([^Callable k x] x)
  ([^Callable k x y] (if (> (k x) (k y)) x y))
  ([^Callable k x y & more]
   (reduce #(max-key k %1 %2) (max-key k x y) more)))

(defn min-key
  "Returns the x for which (k x), a number, is least."
  {:added "1.0"}
  ([^Callable k x] x)
  ([^Callable k x y] (if (< (k x) (k y)) x y))
  ([^Callable k x y & more]
   (reduce #(min-key k %1 %2) (min-key k x y) more)))

(defn distinct
  "Returns a lazy sequence of the elements of coll with duplicates removed."
  {:added "1.0"}
  ^Seq [^Seqable coll]
  (let [step (fn step [xs seen]
               (lazy-seq
                ((fn [[f :as xs] seen]
                   (when-let [s (seq xs)]
                     (if (contains? seen f)
                       (recur (rest s) seen)
                       (cons f (step (rest s) (conj seen f))))))
                 xs seen)))]
    (step coll #{})))

(defn replace
  "Given a map of replacement pairs and a vector/collection, returns a
  vector/seq with any elements = a key in smap replaced with the
  corresponding val in smap."
  {:added "1.0"}
  [^Associative smap ^Seqable coll]
  (if (vector? coll)
    (reduce (fn [v i]
               (if-let [e (find smap (nth v i))]
                 (assoc v i (val e))
                 v))
            coll (range (count coll)))
    (map #(if-let [e (find smap %)] (val e) %) coll)))

(defn repeatedly
  "Takes a function of no args, presumably with side effects, and
  returns an infinite (or length n if supplied) lazy sequence of calls
  to it"
  {:added "1.0"}
  (^Seq [^Callable f] (lazy-seq (cons (f) (repeatedly f))))
  (^Seq [^Number n ^Callable f] (take n (repeatedly f))))

(defn interpose
  "Returns a lazy seq of the elements of coll separated by sep.
  Returns a stateful transducer when no collection is provided."
  {:added "1.0"}
  ^Seq [sep ^Seqable coll]
  (drop 1 (interleave (repeat sep) coll)))

(defn empty
  "Returns an empty collection of the same category as coll, or nil"
  {:added "1.0"}
  ^Collection [coll]
  (empty__ coll))

(defn bound?
  "Returns true if all of the vars provided as arguments have any bound value.
  Implies that deref'ing the provided vars will succeed. Returns true if no vars are provided."
  {:added "1.0"}
  ^Boolean [& vars]
  (every? #(bound?__ ^Var %) vars))

(defn not-empty
  "If coll is empty, returns nil, else coll"
  {:added "1.0"}
  ^Seqable [^Seqable coll] (when (seq coll) coll))

(defn distinct?
  "Returns true if no two of the arguments are ="
  {:added "1.0"}
  (^Boolean [x] true)
  (^Boolean [x y] (not (= x y)))
  (^Boolean [x y & more]
   (if (not= x y)
     (loop [s #{x y} [x & etc :as xs] more]
       (if xs
         (if (contains? s x)
           false
           (recur (conj s x) etc))
         true))
     false)))

(defn format
  "Formats a string using fmt.Sprintf"
  {:added "1.0"}
  ^String [^String fmt & args]
  (apply format__ fmt args))

(defn printf
  "Prints formatted output, as per format"
  {:added "1.0"}
  ^Nil [^String fmt & args]
  (print (apply format fmt args)))

(defmacro ns
  "Sets *ns* to the namespace named by name (unevaluated), creating it
  if needed.  references can be zero or more of:
  (:require ...) (:use ...) (:load ...)
  with the syntax of require/use/load
  respectively, except the arguments are unevaluated and need not be
  quoted. Use of ns is preferred to
  individual calls to in-ns/require/use:

  (ns foo.bar
    (:require [my.lib1 :as lib1])
    (:use [my.lib2]))"
  {:arglists '([name docstring? attr-map? references*])
   :added "1.0"}
  [name & references]
  (let [process-reference
        (fn [[kname & args]]
          `(~(symbol "lace.core" (lace.core/name kname))
             ~@(map #(list 'quote %) args)))
        docstring  (when (string? (first references)) (first references))
        references (if docstring (next references) references)
        name (if docstring
               (vary-meta name assoc :doc docstring)
               name)
        metadata   (when (map? (first references)) (first references))
        references (if metadata (next references) references)
        name (if metadata
               (vary-meta name merge metadata)
               name)
        name-metadata (meta name)]
    `(do
       (lace.core/in-ns '~name)
       ~@(when name-metadata
           `((reset-meta! (lace.core/find-ns '~name) ~name-metadata)))
       ~@(when (and (not= name 'lace.core) (not-any? #(= :refer-clojure (first %)) references))
           `((lace.core/refer '~'lace.core)))
       ~@(map process-reference references)
       (if (= '~name 'lace.core)
         nil
         (do
           (var-set #'*loaded-libs* (conj *loaded-libs* '~name))
           nil)))))

(defmacro refer-clojure
  "Same as (refer 'lace.core <filters>)"
  {:added "1.0"}
  [& filters]
  `(lace.core/refer '~'lace.core ~@filters))

(defmacro defonce
  "defs name to have the value of the expr if the named var is not bound,
  else expr is unevaluated"
  {:added "1.0"}
  [name expr]
  `(let [v# (def ~name)]
     (when-not (bound? v#)
       (def ~name ~expr))))

(defonce ^:dynamic
  ^{:private true
    :doc "True while a verbose load is pending"}
  *loading-verbosely* false)

(defn in-ns
  "Sets *ns* to the namespace named by the symbol, creating it if needed."
  {:added "1.0"}
  ^Namespace [^Symbol name]
  (let [new-ns (lace.core/create-ns name)]
    (when *loading-verbosely*
      (printf "lace.core/in-ns: changing from %s to %s\n" lace.core/*ns* new-ns) (flush))
    (var-set #'lace.core/*ns* new-ns)))

(defonce ^:dynamic
  ^{:private true
    :doc "A set of symbols representing currently loaded libs"}
  *loaded-libs* #{})

(defonce ^:dynamic
  ^{:private true
    :doc "A set of symbols representing available core (not std) namespaces"}
  *core-namespaces* #{})

(defonce ^:dynamic
  ^{:private true
    :doc "A stack of paths currently being loaded"}
  *pending-paths* ())

(defonce ^:dynamic
  ^{:private true
    :doc "A vector of mappings of namespaces to root files

  Each such mapping is a two-element key/value vector. The key is a
  regular expression, matched against the namespace name; the value is
  a map specifying the source from which to load the external
  dependency's root file (currently only the :url key is supported)."}
    
  *ns-sources* [])

(defn- throw-if
  "Throws an exception with a message if pred is true"
  [pred form fmt & args]
  (when pred
    (let [^String message (apply format fmt args)
          exception (ex-info message {:form form})]
      (throw exception))))

(defn- libspec?
  "Returns true if x is a libspec"
  [x]
  (or (symbol? x)
      (and (vector? x)
           (or
            (nil? (second x))
            (keyword? (second x))))))

(defn- prependss
  "Prepends a symbol or a seq to coll"
  [x coll]
  (if (symbol? x)
    (cons x coll)
    (concat x coll)))

(def ^:declared load)

(defn ns-sources
  "Adds namespace mappings to the built-in variable *ns-sources*

  The key/value pairs in the map define how to resolve external
  dependencies. They are appended to the *ns-sources* vector in
  arbitrary order; so, use separate invocations of this function
  to add narrower keys before wider.

  Each value is itself a map containing (primarily) a :url key whose
  value is the URL of the resource. Only http:// and https:// are
  currently supported; everything else is treated as a local
  pathname. HTTP URLs are cached in $HOME/.laced/deps/."
  {:added "1.0"}
  ^Nil [^Map sources]
  (let [existing-source-keys (set (map first *ns-sources*))
        merge-fn (fn [existing-sources [k v]]
                   (concat existing-sources
                           (if-not (existing-source-keys k) [[k v]])))]
    (->> (reduce merge-fn *ns-sources* sources)
         (vec)
         (var-set #'*ns-sources*))))

(defn- load-one
  "Loads a lib given its name. If need-ns, ensures that the associated
  namespace exists after loading. If require, records the load so any
  duplicate loads can be skipped."
  [lib need-ns require]
  (load lib)
  (throw-if (and need-ns (not (find-ns lib)))
            lib
            "namespace '%s' not found after loading '%s'"
            lib (lib-path__ lib))
  (when require
    (var-set #'*loaded-libs* (conj *loaded-libs* lib))))

(defn- load-all
  "Loads a lib given its name and forces a load of any libs it directly or
  indirectly loads. If need-ns, ensures that the associated namespace
  exists after loading. If require, records the load so any duplicate loads
  can be skipped."
  [lib need-ns require]
  (let [libs (binding [*loaded-libs* #{}]
               (load-one lib need-ns require)
               *loaded-libs*)]
    (var-set #'*loaded-libs* (reduce conj *loaded-libs* libs))))

(def ^:private require-opt-keys
  [:exclude :only :rename :refer])

(defn- lib-name__
  [lib]
  lib)

(defn- load-libs-options__
  []
  #{:as :reload :reload-all :require :use :verbose :exclude :only :rename :refer})

(defn- load-lib
  "Loads a lib with options"
  [prefix lib & options]
  (throw-if (and prefix (pos? (index-of__ (name lib) \.)))
            lib
            "Found lib name '%s' containing period with prefix '%s'.  lib names inside prefix lists must not contain periods"
            (name lib) prefix)
  (let [lib (if prefix (symbol (str prefix \. lib)) lib)
        lib (lib-name__ lib)
        opts (apply hash-map options)
        unsupported (seq (remove (load-libs-options__) (keys opts)))
        _ (throw-if unsupported
                    (first unsupported)
                    (apply str "Unsupported option(s) supplied: "
                           (interpose \, unsupported)))
        {:keys [as reload reload-all require use verbose default]} opts
        loaded (contains? *loaded-libs* lib)
        load (cond reload-all
               load-all
               (or reload (not require) (not loaded))
               load-one)
        need-ns (or as use default)
        filter-opts (select-keys opts require-opt-keys)
        undefined-on-entry (not (find-ns lib))]
    (when (and *linter-mode* loaded)
      (println-linter__ (ex-info (str "duplicate require for " lib)
                            {:form lib :_prefix "Parse warning"})))
    (binding [*loading-verbosely* (or *loading-verbosely* verbose)]
      (if load
        (try
          (load lib need-ns require)
          (catch Error e
            (when undefined-on-entry
              (remove-ns lib))
            (throw e)))
        (throw-if (and need-ns (not (find-ns lib)))
                  lib
                  "namespace '%s' not found" lib))
      (when (and need-ns *loading-verbosely*)
        (printf "(lace.core/in-ns '%s)\n" (ns-name *ns*)))
      (when as
        (when *loading-verbosely*
          (printf "(lace.core/alias '%s '%s)\n" as lib))
        (alias as lib))
      (when default
        (when *loading-verbosely*
          (printf "(lace.core/alias '%s '%s)\n" default lib))
        (alias default lib))
      (when (or use (get-refer-opt filter-opts))
        (when *loading-verbosely*
          (printf "(lace.core/refer '%s" lib)
          (doseq [opt filter-opts]
            (printf " %s '%s" (key opt) (print-str (val opt))))
          (printf ")\n"))
        (apply refer lib (mapcat seq filter-opts))))))

(defn- load-libs
  "Loads libs, interpreting libspecs, prefix lists, and flags for
  forwarding to load-lib"
  [& args]
  (let [flags (filter keyword? args)
        opts (interleave flags (repeat true))
        args (filter (complement keyword?) args)]
    ; check for unsupported options
    (let [supported (load-libs-options__)
          unsupported (seq (remove supported flags))]
      (throw-if unsupported
                (first unsupported)
                (apply str "Unsupported option(s) supplied: "
                     (interpose \, unsupported))))
    ; check a load target was specified
    (throw-if (not (seq args)) args "Nothing specified to load")
    (doseq [arg args]
      (if (libspec? arg)
        (apply load-lib nil (prependss arg opts))
        (let [[prefix & args] arg]
          (throw-if (nil? prefix) arg "prefix cannot be nil")
          (doseq [arg args]
            (apply load-lib prefix (prependss arg opts))))))))

(defn- check-cyclic-dependency
  "Detects and rejects non-trivial cyclic load dependencies. The
  exception message shows the dependency chain with the cycle
  highlighted. Ignores the trivial case of a file attempting to load
  itself because."
  [path lib]
  (when (some #{path} (rest *pending-paths*))
    (let [pending (map #(if (= % path) (str "[ " % " ]") %)
                       (cons path *pending-paths*))
          chain (apply str (interpose "->" pending))]
      (throw-if true  lib "Cyclic load dependency: %s" chain))))

(defn require
  "Loads libs, skipping any that are already loaded. Each argument is
  either a libspec that identifies a lib, a prefix list that
  identifies multiple libs whose names share a common prefix, or a
  flag that modifies how all the identified libs are
  loaded. Use :require in the ns macro in preference to calling this
  directly.

  Libs

  A 'lib' is a named set of resources in *classpath* whose contents
  define a library of Clojure code. Lib names are symbols and each lib
  is associated with a Clojure namespace and a Joker package that
  share its name. A lib's name also locates its root directory within
  *classpath* using its package name to classpath-relative path
  mapping. All resources in a lib should be contained in the directory
  structure under its root directory.  All definitions a lib makes
  should be in its associated namespace.

  'require loads a lib by loading its root resource. The root resource
  path is derived from the lib name in the following manner: Consider
  a lib named by the symbol 'x.y.z; it has the root directory
  <*classpath*>/x/y/, and its root resource is
  <*classpath*>/x/y/z.clj. The root resource should contain code to
  create the lib's namespace (usually by using the ns macro) and load
  any additional lib resources.

  Libspecs

  A libspec is a lib name or a vector containing a lib name followed
  by options expressed as sequential keywords and arguments.

  Recognized options:
  :as takes a symbol as its argument and makes that symbol an alias to the
    lib's namespace in the current namespace.
  :refer takes a list of symbols to refer from the namespace or the :all
    keyword to bring in all public vars.

  Prefix Lists

  It's common for Clojure code to depend on several libs whose names have
  the same prefix. When specifying libs, prefix lists can be used to reduce
  repetition. A prefix list contains the shared prefix followed by libspecs
  with the shared prefix removed from the lib names. After removing the
  prefix, the names that remain must not contain any periods.

  Flags

  A flag is a keyword.
  Recognized flags: :reload, :reload-all, :verbose
  :reload forces loading of all the identified libs even if they are
    already loaded
  :reload-all implies :reload and also forces loading of all libs that the
    identified libs directly or indirectly load via require or use
  :verbose triggers printing information about each load, alias, and refer

  Example:

  The following would load the libraries clojure.zip and clojure.set
  abbreviated as 's'.

  (require '(clojure zip [set :as s]))"
  {:added "1.0"}

  ^Nil [& args]
  (apply load-libs :require args))

(defn requiring-resolve
  "Resolves namespace-qualified sym per 'resolve'. If initial resolve
  fails, attempts to require sym's namespace and retries."
  {:added "1.0"}
  ^Var [sym]
  (if (qualified-symbol? sym)
    (or (resolve sym)
        (do (-> sym namespace symbol require)
          (resolve sym)))
    (throw (ex-info (str "Not a qualified symbol: " sym) {}))))

(defn use
  "Like 'require, but also refers to each lib's namespace using
  lace.core/refer. Use :use in the ns macro in preference to calling
  this directly.

  'use accepts additional options in libspecs: :exclude, :only, :rename.
  The arguments and semantics for :exclude, :only, and :rename are the same
  as those documented for lace.core/refer."
  {:added "1.0"}
  ^Nil [& args] (apply load-libs :require :use args))

(defn loaded-libs
  "Returns an UNSORTED set of symbols naming the currently loaded libs"
  {:added "1.0"}
  ^MapSet []
  *loaded-libs*)

(defn load-file
  "Loads code from file f. Does not protect against recursion."
  {:added "1.0"}
  ^Nil [^String f]
  (load-file__ f))

(defn load
  "Loads code from libs, throwing error if cyclic dependency detected,
  and ignoring libs already being loaded."
  {:added "1.0"}
  ^Nil [& libs]
  (doseq [^Symbol lib libs]
    (let [^String path (lib-path__ lib)]
      (when *loading-verbosely*
        (printf "(lace.core/load %s from \"%s\")\n" lib path))
      (check-cyclic-dependency path lib)
      (when-not (= path (first *pending-paths*))
        (binding [*pending-paths* (conj *pending-paths* path)
                  *ns* *ns*]
          (if *linter-mode*
            (in-ns lib)
            (when (not (lace.core/*core-namespaces* lib))
              (lace.lang/LoadLibFromPath lib path))))))))

(defn get-in
  "Returns the value in a nested associative structure,
  where ks is a sequence of keys. Returns nil if the key
  is not present, or the not-found value if supplied."
  {:added "1.0"}
  ([m ^Seqable ks]
   (reduce get m ks))
  ([m ^Seqable ks not-found]
   (loop [sentinel {}
          m m
          ks (seq ks)]
     (if ks
       (let [m (get m (first ks) sentinel)]
         (if (identical? sentinel m)
           not-found
           (recur sentinel m (next ks))))
       m))))

(defn assoc-in
  "Associates a value in a nested associative structure, where ks is a
  sequence of keys and v is the new value and returns a new nested structure.
  If any levels do not exist, hash-maps will be created."
  {:added "1.0"}
  ^Map [^Associative m ^Seqable ks v]
  (let [[k & ks] ks]
    (if ks
      (assoc m k (assoc-in (get m k) ks v))
      (assoc m k v))))

(defn update-in
  "'Updates' a value in a nested associative structure, where ks is a
  sequence of keys and f is a function that will take the old value
  and any supplied args and return the new value, and returns a new
  nested structure.  If any levels do not exist, hash-maps will be
  created."
  {:added "1.0"}
  (^Map [^Associative m ^Seqable ks ^Callable f & args]
   (let [[k & ks] ks]
     (if ks
       (assoc m k (apply update-in (get m k) ks f args))
       (assoc m k (apply f (get m k) args))))))

(defn update
  "'Updates' a value in an associative structure, where k is a
  key and f is a function that will take the old value
  and any supplied args and return the new value, and returns a new
  structure.  If the key does not exist, nil is passed as the old value."
  {:added "1.0"}
  (^Map [^Associative m k ^Callable f]
   (assoc m k (f (get m k))))
  (^Map [^Associative m k ^Callable f x]
   (assoc m k (f (get m k) x)))
  (^Map [^Associative m k ^Callable f x y]
   (assoc m k (f (get m k) x y)))
  (^Map [^Associative m k ^Callable f x y z]
   (assoc m k (f (get m k) x y z)))
  (^Map [^Associative m k ^Callable f x y z & more]
   (assoc m k (apply f (get m k) x y z more))))

(defn coll?
  "Returns true if x implements Collection"
  {:added "1.0"}
  ^Boolean [x] (instance? Collection x))

(defn list?
  "Returns true if x is a List"
  {:added "1.0"}
  ^Boolean [x] (instance? List x))

(defn seqable?
  "Return true if the seq function is supported for x"
  {:added "1.0"}
  ^Boolean [x]
  (or (nil? x)
      (instance? Seqable x)))

(defn callable?
  "Returns true if x implements Callable. Note that many data structures
  (e.g. sets and maps) implement Callable."
  {:added "1.0"}
  ^Boolean [x] (instance? Callable x))

(defn fn?
  "Returns true if x is Fn, i.e. is an object created via fn."
  {:added "1.0"}
  ^Boolean [x] (instance? Fn x))

(defn associative?
  "Returns true if coll implements Associative"
  {:added "1.0"}
  ^Boolean [coll] (instance? Associative coll))

(defn sequential?
  "Returns true if coll implements Sequential"
  {:added "1.0"}
  ^Boolean [coll] (instance? Sequential coll))

(defn counted?
  "Returns true if coll implements count in constant time"
  {:added "1.0"}
  ^Boolean [coll] (instance? Counted coll))

(defn reversible?
  "Returns true if coll implements Reversible"
  {:added "1.0"}
  ^Boolean [coll] (instance? Reversible coll))

(defn indexed?
  "Return true if coll implements Indexed, indicating efficient lookup by index"
  {:added "1.0"}
  ^Boolean [coll] (instance? Indexed coll))

(def
  ^{:doc "bound in a repl to the most recent value printed"
    :added "1.0"}
  *1)

(def
  ^{:doc "bound in a repl to the second most recent value printed"
    :added "1.0"}
  *2)

(def
  ^{:doc "bound in a repl to the third most recent value printed"
    :added "1.0"}
  *3)

(def
  ^{:doc "bound in a repl to the most recent exception caught by the repl"
    :added "1.0"}
  *e)

(defn trampoline
  "trampoline can be used to convert algorithms requiring mutual
  recursion without stack consumption. Calls f with supplied args, if
  any. If f returns a fn, calls that fn with no arguments, and
  continues to repeat, until the return value is not a fn, then
  returns that non-fn value. Note that if you want to return a fn as a
  final value, you must wrap it in some data structure and unpack it
  after trampoline returns."
  {:added "1.0"}
  ([^Callable f]
   (let [ret (f)]
     (if (fn? ret)
       (recur ret)
       ret)))
  ([^Callable f & args]
   (trampoline #(apply f args))))

(defmacro while
  "Repeatedly executes body while test expression is true. Presumes
  some side-effect will cause test to become false/nil. Returns nil"
  {:added "1.0"}
  [test & body]
  `(loop []
     (when ~test
       ~@body
       (recur))))

(defn memoize
  "Returns a memoized version of a referentially transparent function. The
  memoized version of the function keeps a cache of the mapping from arguments
  to results and, when calls with the same arguments are repeated often, has
  higher performance at the expense of higher memory use."
  {:added "1.0"}
  ^Fn [^Callable f]
  (let [mem (atom {})]
    (fn [& args]
      (if-let [e (find @mem args)]
        (val e)
        (let [ret (apply f args)]
          (swap! mem assoc args ret)
          ret)))))

(defn empty?
  "Returns true if coll has no items - same as (not (seq coll)).
  Please use the idiom (seq x) rather than (not (empty? x))"
  {:added "1.0"}
  ^Boolean [^Seqable coll] (not (seq coll)))

(defmacro condp
  "Takes a binary predicate, an expression, and a set of clauses.
  Each clause can take the form of either:

  test-expr result-expr

  test-expr :>> result-fn

  Note :>> is an ordinary keyword.

  For each clause, (pred test-expr expr) is evaluated. If it returns
  logical true, the clause is a match. If a binary clause matches, the
  result-expr is returned, if a ternary clause matches, its result-fn,
  which must be a unary function, is called with the result of the
  predicate as its argument, the result of that call being the return
  value of condp. A single default expression can follow the clauses,
  and its value will be returned if no clause matches. If no default
  expression is provided and no clause matches, an
  exception is thrown."
  {:added "1.0"}

  [pred expr & clauses]
  (when *linter-mode*
    (when (empty? clauses)
      (println-linter__ (ex-info "condp with no clauses" {:form &form :_prefix "Parse error"})))
    (when (= 1 (count clauses))
      (println-linter__ (ex-info "condp with default expression only" {:form &form :_prefix "Parse warning"}))))
  (let [gpred (gensym "pred__")
        gexpr (gensym "expr__")
        emit (fn emit [pred expr args]
               (let [[[a b c :as clause] more]
                     (split-at (if (= :>> (second args)) 3 2) args)
                     n (count clause)]
                 (cond
                   (= 0 n) `(throw (ex-info (str "No matching clause: " ~expr) {}))
                   (= 1 n) a
                   (= 2 n) `(if (~pred ~a ~expr)
                              ~b
                              ~(emit pred expr more))
                   :else `(if-let [p# (~pred ~a ~expr)]
                            (~c p#)
                            ~(emit pred expr more)))))]
    `(let [~gpred ~pred
           ~gexpr ~expr]
       ~(emit gpred gexpr clauses))))

(defmacro add-doc-and-meta {:private true} [name docstring meta]
  `(alter-meta! (var ~name) merge (assoc ~meta :doc ~docstring)))

(add-doc-and-meta *file*
  "The path of the file being evaluated, as a String.

  When there is no file, e.g. in the REPL, the value is not defined."
  {:added "1.0"})

(add-doc-and-meta *main-file*
  "The absolute path of <filename> on the command line, as a String.

  When there is no file, e.g. in the REPL, the value is not defined."
  {:added "1.0"})

(add-doc-and-meta *command-line-args*
  "A sequence of the supplied command line arguments, or nil if
  none were supplied"
  {:added "1.0"
   :tag Seq})

(add-doc-and-meta *classpath*
  "A vector of the classpath elements as configured by --classpath or
  the JOKER_CLASSPATH environment variable.

  Use colon-delimited <cp> (semicolon-delimited on Windows) for source
  directories when loading libraries via :require and the like (but
  not load-file). An empty field denotes the directory containing the
  current file being loaded, with zero or more trailing components
  removed as determined by the number of \".\" separators in the current
  namespace; or, if no file is being loaded, the current
  directory (this is original Joker behavior); a '.' (period) by
  itself denotes solely the current directory. Defaults to the value
  of the JOKER_CLASSPATH environment variable or, if that is
  undefined, the empty string (denoting a single empty field). The
  resulting classpath is stored herein, and this variable is used (in
  lieu of command-line arguments or environment variables) for all
  pertinent subsequent operations."
  {:added "1.0"
   :private true})

(add-doc-and-meta *ns*
  "A Namespace object representing the current namespace."
  {:added "1.0"})

(add-doc-and-meta *in*
  "A BufferedReader object representing standard input for read operations.

  Defaults to stdin."
  {:added "1.0"})

(add-doc-and-meta *out*
  "A IOWriter object representing standard output for print operations.

  Defaults to stdout."
  {:added "1.0"})

(add-doc-and-meta *err*
  "A IOWriter object representing standard error for print operations.

  Defaults to stderr."
  {:added "1.0"})

(add-doc-and-meta *print-readably*
  "When set to logical false, strings and characters will be printed with
  non-alphanumeric characters converted to the appropriate escape sequences.

  Defaults to true"
  {:added "1.0"})

(defmacro letfn
  "fnspec ==> (fname [params*] exprs) or (fname ([params*] exprs)+)

  Takes a vector of function specs and a body, and generates a set of
  bindings of functions to their names. All of the names are available
  in all of the definitions of the functions, as well as the body."
  {:added "1.0",
   :forms '[(letfn [fnspecs*] exprs*)],
   :special-form true}
  [fnspecs & body]
  `(letfn* ~(vec (interleave (map first fnspecs)
                             (map #(cons `fn %) fnspecs)))
           ~@body))

(defn fnil
  "Takes a function f, and returns a function that calls f, replacing
  a nil first argument to f with the supplied value x. Higher arity
  versions can replace arguments in the second and third
  positions (y, z). Note that the function f can take any number of
  arguments, not just the one(s) being nil-patched."
  {:added "1.0"}
  (^Fn [^Callable f x]
   (fn
     ([a] (f (if (nil? a) x a)))
     ([a b] (f (if (nil? a) x a) b))
     ([a b c] (f (if (nil? a) x a) b c))
     ([a b c & ds] (apply f (if (nil? a) x a) b c ds))))
  (^Fn [^Callable f x y]
   (fn
     ([a b] (f (if (nil? a) x a) (if (nil? b) y b)))
     ([a b c] (f (if (nil? a) x a) (if (nil? b) y b) c))
     ([a b c & ds] (apply f (if (nil? a) x a) (if (nil? b) y b) c ds))))
  (^Fn [^Callable f x y z]
   (fn
     ([a b] (f (if (nil? a) x a) (if (nil? b) y b)))
     ([a b c] (f (if (nil? a) x a) (if (nil? b) y b) (if (nil? c) z c)))
     ([a b c & ds] (apply f (if (nil? a) x a) (if (nil? b) y b) (if (nil? c) z c) ds)))))

(defn partition-all
  "Returns a lazy sequence of lists like partition, but may include
  partitions with fewer than n items at the end."
  {:added "1.0"}
  (^Seq [^Number n ^Seqable coll]
   (partition-all n n coll))
  (^Seq [^Number n ^Number step ^Seqable coll]
   (lazy-seq
    (when-let [s (seq coll)]
      (let [seg (doall (take n s))]
        (cons seg (partition-all n step (nthrest s step))))))))

(defn into
  "Returns a new coll consisting of to-coll with all of the items of
  from-coll conjoined."
  {:added "1.0"}
  [to from]
  (reduce conj to from))

(defmacro case
  "Takes an expression, and a set of clauses.

  Each clause can take the form of either:

  test-expr result-expr

  (test-expr ... test-expr)  result-expr

  If the expression is equal to a value of
  test-expr, the corresponding result-expr is returned. A single
  default expression can follow the clauses, and its value will be
  returned if no clause matches. If no default expression is provided
  and no clause matches, an exception is thrown."
  {:added "1.0"}
  [expr & clauses]
  (loop [all-cases #{}
         [[test then] & more-clauses] (partition 2 clauses)]
    (when test
      (let [cases (if (list? test) (set test) (set [test]))]
        (when (some cases all-cases)
          (let [e (ex-info (str "Duplicate case test constant: " test) {:form test :_prefix "Parse error"})]
            (if *linter-mode*
              (println-linter__ e)
              (throw e))))
        (recur (into all-cases cases) more-clauses))))
  (let [parts (partition-all 2 clauses)
        setized (for [p parts
                      :let [[test then] p]]
                  (if (= 2 (count p))
                    [(if (list? test) (list 'quote (set test)) (list 'quote (set [test]))) then]
                    [test]))
        transformed-clauses (apply concat setized)]
    `(condp contains? ~expr
       ~@transformed-clauses)))

(defn mapv
  "Returns a vector consisting of the result of applying f to the
  set of first items of each coll, followed by applying f to the set
  of second items in each coll, until any one of the colls is
  exhausted.  Any remaining items in other colls are ignored. Function
  f should accept number-of-colls arguments."
  {:added "1.0"}
  (^Vector [^Callable f coll]
   (reduce (fn [v o] (conj v (f o))) [] coll))
  (^Vector [^Callable f c1 c2]
   (into [] (map f c1 c2)))
  (^Vector [^Callable f c1 c2 c3]
   (into [] (map f c1 c2 c3)))
  (^Vector [^Callable f c1 c2 c3 & colls]
   (into [] (apply map f c1 c2 c3 colls))))

(defn filterv
  "Returns a vector of the items in coll for which
  (pred item) returns true. pred must be free of side-effects."
  {:added "1.0"}
  ^Vector [^Callable pred coll]
  (reduce (fn [v o] (if (pred o) (conj v o) v))
          []
          coll))

(defn slurp
  "Opens file f and reads all its contents, returning a string."
  {:added "1.0"}
  ^String [^String f]
  (slurp__ f))

(defn spit
  "Opposite of slurp.  Opens file f, writes content, then
  closes f."
  {:added "1.0"}
  ^Nil [f content & options]
  (spit__ f content (apply hash-map options)))

(defn flatten
  "Takes any nested combination of sequential things (lists, vectors,
  etc.) and returns their contents as a single, flat sequence.
  (flatten nil) returns an empty sequence."
  {:added "1.0"}
  ^Seq [x]
  (filter (complement sequential?)
          (rest (tree-seq sequential? seq x))))

(defn group-by
  "Returns a map of the elements of coll keyed by the result of
  f on each element. The value at each key will be a vector of the
  corresponding elements, in the order they appeared in coll."
  {:added "1.0"}
  ^Map [^Callable f coll]
  (reduce
   (fn [ret x]
     (let [k (f x)]
       (assoc ret k (conj (get ret k []) x))))
   {} coll))

(defn partition-by
  "Applies f to each value in coll, splitting it each time f returns a
  new value.  Returns a lazy seq of partitions."
  {:added "1.0"}
  ^Seq [^Callable f ^Seqable coll]
  (lazy-seq
   (when-let [s (seq coll)]
     (let [fst (first s)
           fv (f fst)
           run (cons fst (take-while #(= fv (f %)) (next s)))]
       (cons run (partition-by f (seq (drop (count run) s))))))))

(defn frequencies
  "Returns a map from distinct items in coll to the number of times
  they appear."
  {:added "1.0"}
  ^Map [coll]
  (reduce (fn [counts x]
            (assoc counts x (inc (get counts x 0))))
          {} coll))

(defn reductions
  "Returns a lazy seq of the intermediate values of the reduction (as
  per reduce) of coll by f, starting with init."
  {:added "1.0"}
  (^Seq [^Callable f ^Seqable coll]
   (lazy-seq
    (if-let [s (seq coll)]
      (reductions f (first s) (rest s))
      (list (f)))))
  (^Seq [^Callable f init ^Seqable coll]
   (cons init
         (lazy-seq
          (when-let [s (seq coll)]
            (reductions f (f init (first s)) (rest s)))))))

(defn rand-nth
  "Return a random element of the (sequential) collection. Will have
  the same performance characteristics as nth for the given
  collection."
  {:added "1.0"}
  [coll]
  (nth coll (rand-int (count coll))))

(defn shuffle
  "Return a random permutation of coll"
  {:added "1.0"}
  ^Vector [coll]
  (shuffle__ coll))

(defn map-indexed
  "Returns a lazy sequence consisting of the result of applying f to 0
  and the first item of coll, followed by applying f to 1 and the second
  item in coll, etc, until coll is exhausted. Thus function f should
  accept 2 arguments, index and item."
  {:added "1.0"}
  ^Seq [^Callable f ^Seqable coll]
  (let [mapi (fn mapi [idx coll]
               (lazy-seq
                (when-let [s (seq coll)]
                  (cons (f idx (first s)) (mapi (inc idx) (rest s))))))]
    (mapi 0 coll)))

(defn keep
  "Returns a lazy sequence of the non-nil results of (f item). Note,
  this means false return values will be included.  f must be free of
  side-effects."
  {:added "1.0"}
  ^Seq [^Callable f ^Seqable coll]
  (lazy-seq
   (when-let [s (seq coll)]
     (let [x (f (first s))]
       (if (nil? x)
         (keep f (rest s))
         (cons x (keep f (rest s))))))))

(defn keep-indexed
  "Returns a lazy sequence of the non-nil results of (f index item). Note,
  this means false return values will be included.  f must be free of
  side-effects."
  {:added "1.0"}
  ^Seq [^Callable f ^Seqable coll]
  (let [keepi (fn keepi [idx coll]
                (lazy-seq
                 (when-let [s (seq coll)]
                   (let [x (f idx (first s))]
                     (if (nil? x)
                       (keepi (inc idx) (rest s))
                       (cons x (keepi (inc idx) (rest s))))))))]
    (keepi 0 coll)))

(defn bounded-count
  "If coll is counted? returns its count, else will count at most the first n
  elements of coll using its seq"
  {:added "1.0"}
  ^Int [^Number n coll]
  (if (counted? coll)
    (count coll)
    (loop [i 0 s (seq coll)]
      (if (and s (< i n))
        (recur (inc i) (next s))
        i))))

(defn every-pred
  "Takes a set of predicates and returns a function f that returns true if all of its
  composing predicates return a logical true value against all of its arguments, else it returns
  false. Note that f is short-circuiting in that it will stop execution on the first
  argument that triggers a logical false result against the original predicates."
  {:added "1.0"}
  (^Fn [^Callable p]
   (fn ep1
     ([] true)
     ([x] (boolean (p x)))
     ([x y] (boolean (and (p x) (p y))))
     ([x y z] (boolean (and (p x) (p y) (p z))))
     ([x y z & args] (boolean (and (ep1 x y z)
                                   (every? p args))))))
  (^Fn [^Callable p1 ^Callable p2]
   (fn ep2
     ([] true)
     ([x] (boolean (and (p1 x) (p2 x))))
     ([x y] (boolean (and (p1 x) (p1 y) (p2 x) (p2 y))))
     ([x y z] (boolean (and (p1 x) (p1 y) (p1 z) (p2 x) (p2 y) (p2 z))))
     ([x y z & args] (boolean (and (ep2 x y z)
                                   (every? #(and (p1 %) (p2 %)) args))))))
  (^Fn [^Callable p1 ^Callable p2 ^Callable p3]
   (fn ep3
     ([] true)
     ([x] (boolean (and (p1 x) (p2 x) (p3 x))))
     ([x y] (boolean (and (p1 x) (p2 x) (p3 x) (p1 y) (p2 y) (p3 y))))
     ([x y z] (boolean (and (p1 x) (p2 x) (p3 x) (p1 y) (p2 y) (p3 y) (p1 z) (p2 z) (p3 z))))
     ([x y z & args] (boolean (and (ep3 x y z)
                                   (every? #(and (p1 %) (p2 %) (p3 %)) args))))))
  (^Fn [^Callable p1 ^Callable p2 ^Callable p3 & ps]
   (let [ps (list* p1 p2 p3 ps)]
     (fn epn
       ([] true)
       ([x] (every? #(% x) ps))
       ([x y] (every? #(and (% x) (% y)) ps))
       ([x y z] (every? #(and (% x) (% y) (% z)) ps))
       ([x y z & args] (boolean (and (epn x y z)
                                     (every? #(every? % args) ps))))))))

(defn some-fn
  "Takes a set of predicates and returns a function f that returns the first logical true value
  returned by one of its composing predicates against any of its arguments, else it returns
  logical false. Note that f is short-circuiting in that it will stop execution on the first
  argument that triggers a logical true result against the original predicates."
  {:added "1.0"}
  (^Fn [^Callable p]
   (fn sp1
     ([] nil)
     ([x] (p x))
     ([x y] (or (p x) (p y)))
     ([x y z] (or (p x) (p y) (p z)))
     ([x y z & args] (or (sp1 x y z)
                         (some p args)))))
  (^Fn [^Callable p1 ^Callable p2]
   (fn sp2
     ([] nil)
     ([x] (or (p1 x) (p2 x)))
     ([x y] (or (p1 x) (p1 y) (p2 x) (p2 y)))
     ([x y z] (or (p1 x) (p1 y) (p1 z) (p2 x) (p2 y) (p2 z)))
     ([x y z & args] (or (sp2 x y z)
                         (some #(or (p1 %) (p2 %)) args)))))
  (^Fn [^Callable p1 ^Callable p2 ^Callable p3]
   (fn sp3
     ([] nil)
     ([x] (or (p1 x) (p2 x) (p3 x)))
     ([x y] (or (p1 x) (p2 x) (p3 x) (p1 y) (p2 y) (p3 y)))
     ([x y z] (or (p1 x) (p2 x) (p3 x) (p1 y) (p2 y) (p3 y) (p1 z) (p2 z) (p3 z)))
     ([x y z & args] (or (sp3 x y z)
                         (some #(or (p1 %) (p2 %) (p3 %)) args)))))
  (^Fn [^Callable p1 ^Callable p2 ^Callable p3 & ps]
   (let [ps (list* p1 p2 p3 ps)]
     (fn spn
       ([] nil)
       ([x] (some #(% x) ps))
       ([x y] (some #(or (% x) (% y)) ps))
       ([x y z] (some #(or (% x) (% y) (% z)) ps))
       ([x y z & args] (or (spn x y z)
                           (some #(some % args) ps)))))))

(defn- ^{:dynamic true} assert-valid-fdecl
  "A good fdecl looks like (([a] ...) ([a b] ...)) near the end of defn."
  [form fdecl]
  (when (empty? fdecl) (throw (ex-info "Parameter declaration missing" {:form form})))
  (let [argdecls (map
                  #(if (seq? %)
                     (first %)
                     (throw (ex-info
                             (if (seq? (first fdecl))
                               (str "Invalid signature: \""
                                    %
                                    "\" must be a list")
                               (str "Parameter declaration \""
                                    %
                                    "\" must be a vector"))
                             {:form form})))
                  fdecl)
        bad-args (seq (remove #(vector? %) argdecls))]
    (when bad-args
      (throw (ex-info (str "Parameter declaration \"" (first bad-args)
                           "\" must be a vector")
                      {:form form})))))

(defn realized?
  "Returns true if a value has been produced for a delay or lazy sequence."
  {:added "1.0"}
  ^Boolean [^Pending x] (realized?__ x))

(defmacro cond->
  "Takes an expression and a set of test/form pairs. Threads expr (via ->)
  through each form for which the corresponding test
  expression is true. Note that, unlike cond branching, cond-> threading does
  not short circuit after the first true test expression."
  {:added "1.0"}
  [expr & clauses]
  (if *linter-mode*
    (when-not (even? (count clauses))
      (println-linter__ (ex-info "Odd number of clauses in cond->" {:form &form :_prefix "Parse warning"})))
    (assert (even? (count clauses))))
  (when (and *linter-mode* (not (seq clauses)) (not (false? (:no-forms-threading (:rules *linter-config*)))))
    (println-linter__ (ex-info "No forms in cond->" {:form &form :_prefix "Parse warning"})))
  (let [g (gensym)
        steps (map (fn [[test step]] `(if ~test (-> ~g ~step) ~g))
                   (partition 2 clauses))]
    `(let [~g ~expr
           ~@(interleave (repeat g) (butlast steps))]
       ~(if (empty? steps)
          g
          (last steps)))))

(defmacro cond->>
  "Takes an expression and a set of test/form pairs. Threads expr (via ->>)
  through each form for which the corresponding test expression
  is true.  Note that, unlike cond branching, cond->> threading does not short circuit
  after the first true test expression."
  {:added "1.0"}
  [expr & clauses]
  (if *linter-mode*
    (when-not (even? (count clauses))
      (println-linter__ (ex-info "Odd number of clauses in cond->>" {:form &form :_prefix "Parse warning"})))
    (assert (even? (count clauses))))
  (when (and *linter-mode* (not (seq clauses)) (not (false? (:no-forms-threading (:rules *linter-config*)))))
    (println-linter__ (ex-info "No forms in cond->>" {:form &form :_prefix "Parse warning"})))
  (let [g (gensym)
        steps (map (fn [[test step]] `(if ~test (->> ~g ~step) ~g))
                   (partition 2 clauses))]
    `(let [~g ~expr
           ~@(interleave (repeat g) (butlast steps))]
       ~(if (empty? steps)
          g
          (last steps)))))

(defmacro as->
  "Binds name to expr, evaluates the first form in the lexical context
  of that binding, then binds name to that result, repeating for each
  successive form, returning the result of the last form."
  {:added "1.0"}
  [expr name & forms]
  (when (and *linter-mode* (not (seq forms)) (not (false? (:no-forms-threading (:rules *linter-config*)))))
    (println-linter__ (ex-info "No forms in as->" {:form &form :_prefix "Parse warning"})))
  `(let [~name ~expr
         ~@(interleave (repeat name) (butlast forms))]
     ~(if (empty? forms)
        name
        (last forms))))

(defmacro some->
  "When expr is not nil, threads it into the first form (via ->),
  and when that result is not nil, through the next etc."
  {:added "1.0"}
  [expr & forms]
  (when (and *linter-mode* (not (seq forms)) (not (false? (:no-forms-threading (:rules *linter-config*)))))
    (println-linter__ (ex-info "No forms in some->" {:form &form :_prefix "Parse warning"})))
  (let [g (gensym)
        steps (map (fn [step] `(if (nil? ~g) nil (-> ~g ~step)))
                   forms)]
    `(let [~g ~expr
           ~@(interleave (repeat g) (butlast steps))]
       ~(if (empty? steps)
          g
          (last steps)))))

(defmacro some->>
  "When expr is not nil, threads it into the first form (via ->>),
  and when that result is not nil, through the next etc."
  {:added "1.0"}
  [expr & forms]
  (when (and *linter-mode* (not (seq forms)) (not (false? (:no-forms-threading (:rules *linter-config*)))))
    (println-linter__ (ex-info "No forms in some->>" {:form &form :_prefix "Parse warning"})))
  (let [g (gensym)
        steps (map (fn [step] `(if (nil? ~g) nil (->> ~g ~step)))
                   forms)]
    `(let [~g ~expr
           ~@(interleave (repeat g) (butlast steps))]
       ~(if (empty? steps)
          g
          (last steps)))))

(defn dedupe
  "Returns a lazy sequence removing consecutive duplicates in coll."
  {:added "1.0"}
  ^Seq [^Seqable coll]
  (lazy-seq
   (when (seq coll)
     (cons (first coll)
           (dedupe (drop-while #(= (first coll) %) (rest coll)))))))

(defn random-sample
  "Returns items from coll with random probability of prob (0.0 -
  1.0)."
  {:added "1.0"}
  ^Seq [^Number prob ^Seqable coll]
  (filter (fn [_] (< (rand) prob)) coll))

(defn run!
  "Runs the supplied procedure (via reduce), for purposes of side
  effects, on successive items in the collection. Returns nil."
  {:added "1.0"}
  ^Nil [^Callable proc coll]
  (reduce #(proc %2) nil coll)
  nil)

(def ^{:added "1.0"} default-data-readers
  "Default map of data reader functions provided by Joker. May be
  overridden by binding *data-readers*."
  {})

(defn update-keys
  "m f => {(f k) v ...}
  Given a map m and a function f of 1-argument, returns a new map whose
  keys are the result of applying f to the keys of m, mapped to the
  corresponding values of m.
  f must return a unique key for each key of m, else the behavior is undefined."
  {:added "1.1"}
  ^Map [m ^Callable f]
  (with-meta
    (reduce-kv (fn [acc k v] (assoc acc (f k) v)) {} m)
    (meta m)))

(defn update-vals
  "m f => {k (f v) ...}
  Given a map m and a function f of 1-argument, returns a new map where the keys of m
  are mapped to result of applying f to the corresponding values of m."
  {:added "1.1"}
  ^Map [m ^Callable f]
  (with-meta
    (reduce-kv (fn [acc k v] (assoc acc k (f v))) {} m)
    (meta m)))

(defn lace-version
  "Returns lace version as a printable string."
  {:added "1.0"}
  ^String []
  (lace-version__))

(defn- check-valid-options
  "Throws an exception if the given option map contains keys not listed
  as valid, else returns nil."
  [options & valid-keys]
  (when (seq (apply disj (apply hash-set (keys options)) valid-keys))
    (throw
      (ex-info
        (apply str "Only these options are valid: "
          (first valid-keys)
          (map #(str ", " %) (rest valid-keys))) {}))))

;;multimethods

(defn- multimethod__
  [name dispatch-fn default hierarchy]
  (when hierarchy
    (throw (ex-info ":hierarchy not yet supported by lace.core/defmulti" {})))
  (let [mfatom (atom {})]
    (with-meta
      (fn [& args]
        (let [dispatch-value (apply dispatch-fn args)
              method
              (get @mfatom
                   dispatch-value
                   (get @mfatom
                        default
                        (fn [& args]
                          (throw (ex-info (format "No method in multimethod '%s' for dispatch value: %s"
                                                  name (pr-str dispatch-value)) {})))))]
          (apply method args)))
      {:dispatch-fn dispatch-fn :default default :method-table mfatom})))

(defmacro defmulti
  "Creates a new multimethod with the associated dispatch function.
  The docstring and attr-map are optional.

  Options are key-value pairs and may be one of:

  :default

  The default dispatch value, defaults to :default

  :hierarchy (UNSUPPORTED)

  The value used for hierarchical dispatch (e.g. ::square is-a ::shape)

  Hierarchies are type-like relationships that do not depend upon type
  inheritance. By default Clojure's multimethods dispatch off of a
  global hierarchy map.  However, a hierarchy relationship can be
  created with the derive function used to augment the root ancestor
  created with make-hierarchy.

  Multimethods expect the value of the hierarchy option to be supplied as
  a reference type e.g. a var (i.e. via the Var-quote dispatch macro #'
  or the var special form)."
  {:arglists '([name docstring? attr-map? dispatch-fn & options])
   :added "1.0"}
  [mm-name & options]
  (let [docstring   (if (string? (first options))
                      (first options)
                      nil)
        options     (if (string? (first options))
                      (next options)
                      options)
        m           (if (map? (first options))
                      (first options)
                      {})
        options     (if (map? (first options))
                      (next options)
                      options)
        dispatch-fn (first options)
        options     (next options)
        m           (if docstring
                      (assoc m :doc docstring)
                      m)
        m           (if (meta mm-name)
                      (conj (meta mm-name) m)
                      m)
        mm-name (with-meta mm-name m)]
    (when (= (count options) 1)
      (throw (ex-info "The syntax for defmulti has changed. Example: (defmulti name dispatch-fn :default dispatch-value)" {})))
    (let [options   (apply hash-map options)
          default   (get options :default :default)
          hierarchy (get options :hierarchy nil)]
      (check-valid-options options :default :hierarchy)
      `(let [v# (def ~mm-name)]
         (when-not (and (bound? v#)
                        (fn? (deref v#))
                        (:method-table (meta (deref v#))))
           (let [fndef# (multimethod__ ~(name mm-name) ~dispatch-fn ~default ~hierarchy)]
                (def ~mm-name fndef#)))))))

(defmacro defmethod
  "Creates and installs a new method of multimethod associated with dispatch-value. "
  {:added "1.0"}
  [multifn dispatch-val & fn-tail]
  `(do
     (swap-vals! (:method-table (meta ~multifn)) assoc ~dispatch-val (fn ~@fn-tail))
     ~multifn))

(defn remove-all-methods
  "Removes all of the methods of multimethod."
  {:added "1.0"}
  [multifn]
  (let [mfm (meta multifn)
        mfatom (:method-table mfm)]
    (reset! mfatom {}))
  multifn)

(defn remove-method
  "Removes the method of multimethod associated with dispatch-value."
  {:added "1.0"}
  [multifn dispatch-val]
  (throw (ex-info "method removal not yet supported by lace.core" {})))

(defn prefer-method
  "Causes the multimethod to prefer matches of dispatch-val-x over dispatch-val-y
   when there is a conflict"
  {:added "1.0"}
  [multifn dispatch-val-x dispatch-val-y]
  (throw (ex-info "method preference not yet supported by lace.core" {})))

(defn methods
  "Given a multimethod, returns a map of dispatch values -> dispatch fns"
  {:added "1.0"}
  ^Map [multifn]
  (let [mfm (meta multifn)
        mfatom (:method-table mfm)]
    @mfatom))

(defn get-method
  "Given a multimethod and a dispatch value, returns the dispatch fn
  that would apply to that value, or nil if none apply and no default"
  {:added "1.0"}
  ^Fn [multifn dispatch-val]
  (let [mfm (meta multifn)
        mfatom (:method-table mfm)
        default (:default mfm)]
    (get @mfatom
         dispatch-val
         (get @mfatom default))))

(defn prefers
  "Given a multimethod, returns a map of preferred value -> set of other values"
  {:added "1.0"}
  ^Map [multifn]
  (throw (ex-info "method preference not yet supported by lace.core" {})))

(def ^{:private true
       :doc "Returns currently registered types as a map."
       :added "1.0"
       :tag Map}
  types__ types__)

(defmacro go
  "Schedules the body to run inside a goroutine.
  Immediately returns a channel which will receive the result of the body when
  completed.
  If exception is thrown inside the body, it will be caught and re-thrown upon
  reading from the returned channel.

  Joker is single threaded and uses the GIL (Global Interpreter Lock) to make sure
  only one goroutine (including the root one) executes at the same time.
  However, channel operations and some I/O functions (lace.http/send, lace.os/sh*, lace.os/exec,
  and lace.time/sleep) release the GIL and allow other goroutines to run.
  So using goroutines only makes sense if you do I/O (specifically, calling the above functions)
  inside them. Also, note that a goroutine may never have a chance to run if the root goroutine
  (or another goroutine) doesn't do any I/O or channel operations (<! or >!)."
  {:added "1.0"}
  [& body]
  `(go__ (fn [] ~@body)))

(defn chan
  "Returns a new channel with an optional buffer of size n."
  {:added "1.0"}
  (^Channel [] (chan__ 0))
  (^Channel [^Int n] (chan__ n)))

(defn <!
  "Takes a value from ch.
  Returns nil if ch is closed and nothing is available on ch.
  Blocks if nothing is available on ch and ch is not closed."
  {:added "1.0"}
  [^Channel ch]
  (<!__ ch))

(defn >!
  "Puts val into ch.
  Throws an exception if val is nil.
  Blocks if ch is full (no buffer space is available).
  Returns true unless ch is already closed."
  {:added "1.0"}
  [^Channel ch val]
  (>!__ ch val))

(defn close!
  "Closes a channel. The channel will no longer accept any puts (they
  will be ignored). Data in the channel remains available for taking, until
  exhausted, after which takes will return nil. If there are any
  pending takes, they will be dispatched with nil. Closing a closed
  channel is a no-op. Returns nil.

  Logically closing happens after all puts have been delivered. Therefore, any
  blocked puts will remain blocked until a taker releases them."
  {:added "1.0"}
  [^Channel ch]
  (close!__ ch))

(defn- go-spew
  "Dump ('spew') internal Go structures for object to stderr.

  Returns true if enabled due to building Joker with the 'go_spew' build tag,
  false otherwise or if some other error occurred (which will be printed to
  stderr).

  Use the optional (map) argument to specify ConfigState
  settings. E.g. {:MaxDepth 10} specifies a maximum depth of 10
  levels. Defaults are per the default config state.

  For more info, see: https://github.com/jcburley/go-spew"
  {:added "1.0"}
  (^Boolean [o]
   (go-spew__ o))
  (^Boolean [o ^Map cfg]
   (go-spew__ o cfg)))

(defn- verbosity-level
  "Verbosity level as specified via the --verbose option to Joker."
  {:added "1.0"}
  []
  (verbosity-level__))

(defn- ns-initialized?
  "Returns whether the namespace, denoted by the symbol, has been initialized."
  {:added "1.0"}
  ^Boolean [^Symbol s]
  (ns-initialized?__ s))

(defn exit
  "Causes the current program to exit with the given status code (defaults to 0)."
  {:added "1.0"}
  ([] (exit 0))
  ([^Int code]
   (exit__ code)))
