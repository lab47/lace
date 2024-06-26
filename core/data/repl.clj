(ns
  ^{:doc "Utilities meant to be used interactively at the REPL."
    :added "1.0"}
  lace.repl
  (:require [lace.string :as s]))
  

(def ^:private special-doc-map
  '{def {:forms [(def symbol doc-string? init?)]
         :doc "Creates and interns a global var with the name
  of symbol in the current namespace (*ns*) or locates such a var if
  it already exists.  If init is supplied, it is evaluated, and the
  root binding of the var is set to the resulting value.  If init is
  not supplied, the root binding of the var is unaffected."}
    do {:forms [(do exprs*)]
        :doc "Evaluates the expressions in order and returns the value of
  the last. If no expressions are supplied, returns nil."}
    if {:forms [(if test then else?)]
        :doc "Evaluates test. If not the singular values nil or false,
  evaluates and yields then, otherwise, evaluates and yields else. If
  else is not supplied it defaults to nil."}
    quote {:forms [(quote form)]
           :doc "Yields the unevaluated form."}
    recur {:forms [(recur exprs*)]
           :doc "Evaluates the exprs in order, then, in parallel, rebinds
  the bindings of the recursion point to the values of the exprs.
  Execution then jumps back to the recursion point, a loop or fn method."}
    throw {:forms [(throw expr)]
           :doc "The expr is evaluated and thrown, therefore it should yield an Error object.
  User code should normally use (ex-info) function to create new Error objects."}
    try {:forms [(try expr* catch-clause* finally-clause?)]
         :doc "catch-clause => (catch type name expr*)
  finally-clause => (finally expr*)

  Catches and handles errors.
  User code should normally use (ex-info) function to create new Error objects."}
    var {:forms [(var symbol)]
         :doc "The symbol must resolve to a var, and the Var object
  itself (not its value) is returned. The reader macro #'x expands to (var x)."}})

(defn- special-doc [name-symbol]
  (assoc (or (special-doc-map name-symbol) (meta (resolve name-symbol)))
         :name name-symbol
         :special-form true))

(defn- namespace-doc [nspace]
  (assoc (meta nspace) :name (ns-name nspace)))

(defn- print-doc [{n :ns
                   nm :name
                   :keys [forms arglists special-form doc url macro spec]
                   :as m}]
  (println "-------------------------")
  (println (or spec (str (when n (str (ns-name n) "/")) nm)))
  (when forms
    (doseq [f forms]
      (print "  ")
      (prn f)))
  (when arglists
    (prn arglists))
  (cond
    special-form
    (do
      (println "Special Form")
      (println " " doc)
      (if (contains? m :url)
        (when url
          (println (str "\n  Please see http://clojure.org/" url)))
        (println (str "\n  Please see http://clojure.org/special_forms#" nm))))
    macro
    (println "Macro")
    spec
    (println "Spec"))
  (when doc (println " " doc)))

(defmacro doc
  "Prints documentation for a var, type, or special form given its name,
  or for a spec if given a keyword"
  {:added "1.0"}
  [name]
  (if-let [special-name ('{& fn catch try finally try} name)]
    `(#'print-doc (#'special-doc '~special-name))
    (cond
      (special-doc-map name) `(#'print-doc (#'special-doc '~name))
      (keyword? name) (println "Keywords (spec) not yet supported")
      (find-ns name) `(#'print-doc (#'namespace-doc (find-ns '~name)))
      (resolve name) (let [x# (resolve name)
                           v# (if (= (type x#) Type)
                                name
                                `(var ~name))]
                       `(#'print-doc (meta ~v#))))))

;; ----------------------------------------------------------------------
;; Examine Clojure functions (Vars, really)

(defn apropos
  "Given a regular expression or stringable thing, return a seq of all
public definitions in all currently-loaded namespaces that match the
str-or-pattern."
  {:added "1.0"}
  [str-or-pattern]
  (let [matches? (if (instance? Regex str-or-pattern)
                   #(re-find str-or-pattern (str %))
                   #(lace.string/includes? (str %) (str str-or-pattern)))]
    (sort (mapcat (fn [ns]
                    (let [ns-name (str ns)]
                      (map #(symbol ns-name (str %))
                           (filter matches? (keys (ns-publics ns))))))
                  (all-ns)))))

(defn dir-fn
  "Returns a sorted seq of symbols naming public vars in
  a namespace or namespace alias. Looks for aliases in *ns*"
  {:added "1.0"}
  [ns]
  (sort (map first (ns-publics (the-ns (get (ns-aliases *ns*) ns ns))))))

(defmacro dir
  "Prints a sorted directory of public vars in a namespace"
  {:added "1.0"}
  [nsname]
  `(doseq [v# (dir-fn '~nsname)]
     (println v#)))
