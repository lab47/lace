;;;; See DEVELOPER.md for information on this script.

;;; Load all core namespaces. This pulls in any std namespaces upon which
;;; they depend. Then, set 'preloaded' to those std namespaces, forcing
;;; complete evaluation so lazy evaluation wouldn't find lace.os due to
;;; the subsequent require.
(doseq [ns (remove #(= % 'user) lace.core/*core-namespaces*)] (require ns))
(def preloaded (set (remove lace.core/*core-namespaces* (filter #(lace.core/ns-initialized? %) (map #(symbol (str %)) (all-ns))))))

(require '[lace.string :as s]
         '[lace.os :as os])

(def rpl s/replace)

;;; Discover namespaces dynamically by finding *.clj files.
(def namespaces
  (vec (->> (os/ls ".")
            (remove :dir?)
            (map :name)
            (remove #(= "generate-std.clj" %))
            (filter #(s/ends-with? % ".clj"))
            (map #(rpl % #"[.]clj$" ""))
            (map symbol))))

(apply require :reload namespaces) ; :reload in case namespaces become 'lace.base64 etc.

(def fn-template
  (slurp "fn.tmpl"))

(def arity-template
  (slurp "arity.tmpl"))

(def package-template
  (slurp "package.tmpl"))

(def package-slow-template
  (slurp "package-slow.tmpl"))

(def package-fast-template
  (slurp "package-fast.tmpl"))

(def intern-template
  (slurp "intern.tmpl"))

(def addmeta-template
  (s/trim-right (slurp "addmeta.tmpl")))

(defn q
  [s]
  (str "\"" s "\""))

(defn raw-quoted-string
  "Returns a Go-style backtick-quoted string with backticks handled by appending double-quoted backticks"
  [s]
  (str "`" (rpl s "`" "` + \"`\" + `") "`"))
  

(defn go-name
  "Convert Clojure-style function name to unique Go-style name suitable as its internal implementation."
  [fn-name]
  (let [n (-> fn-name
              (rpl "-" "_")
              (rpl "?" "")
              (str "_"))]
    (if (s/ends-with? fn-name "?")
      (str "is" n)
      n)))

(defn extract-args
  "Generate code to extract the arguments in the list, assigning the extracted values to variable names corresponding to the argument names."
  [args]
  (s/join
   "\n\t\t"
   (map-indexed
    (fn [i arg]
      (let [m (meta arg)
            t (cond-> (:tag m)
                (:varargs m) (str "s"))]
        (str arg ", err := Extract" t "(_env, _args, " (str i) "); if err != nil { return nil, err }")))
    args)))

(defn handle-varargs
  [args]
  (let [c (count args)]
    (if (and (> c 1)
             (= '& (nth args (- c 2))))
      (let [vargs (vary-meta (last args) assoc :varargs true)]
        (conj (subvec args 0 (- c 2)) vargs))
      args)))

(defn ^:private type-name
  [tag]
  (if (vector? tag)
    (str (first tag) "Vector")
    (str tag)))

(defn generate-arity
  [args go tag]
  (let [handle-args (handle-varargs args)
        cnt (count handle-args)
        varargs? (< cnt (count args))
        go-expr (cond
                  (string? go) go
                  varargs? (:varargs go)
                  :else (get go cnt))
        go-res (if (s/starts-with? go-expr "!")
                 (subs go-expr 1)
                 (str "_res, err := " go-expr))]
    (-> arity-template
        (rpl "{arity}" (if varargs? "true" (str "_c == " (count args))))
        (rpl "{arityCheck}" (if varargs?
                              (str "if err := CheckArity(_env, _args, " (dec cnt) ", " 999 "); err != nil { return nil, err }")
                              "{blank}"))
        (rpl "{args}" (if (empty? args)
                        "{blank}"
                        (extract-args handle-args)))
        (rpl "{goExpr}" (rpl go-res "; " "\n\t\t"))
        (rpl "{return}"
             (if tag
               (str "return Make" (type-name tag) "(_res), err")
               "return _res, err")))))

(defn generate-arglist
  [args]
  (str "NewVectorFrom("
       (s/join ", " (for [arg args]
                      (str "MakeSymbol(" (q (str arg)) ")")))
       ")"))

(defn make-value
  "Returns code to make the Joker object representing the given value.

  E.g. 'String{S: \"value\"}'. Except for integers, the values are treated as strings (for now)."
  [v]
  (condp = (str (type v))
    "Int" (str "Int{I: " v "}")
    (str "String{S: " (q v) "}")))

(defn add-other-meta
  "Append meta tags other than what are normally present or irrelevant (:go)."
  [m]
  (let [m (dissoc m :doc :added :arglists :ns :name :file :line :column :go)]
    (s/join "" (map #(-> addmeta-template
                      (rpl "{key}" (s/replace-first (str (key %)) ":" ""))
                      (rpl "{value}" (make-value (val %)))) m))))

(defn generate-fn-decl
  [ns-name ns-name-final k v]
  (let [m (meta v)
        arglists (:arglists m)
        go-fn-name (go-name (str k))
        arities (s/join "\n\t" (map #(generate-arity % (:go m) (:tag m)) arglists))
        fn-str (-> fn-template
                   (rpl "{goName}" go-fn-name)
                   (rpl "{pkg}" ns-name)
                   (rpl "{fnName}" (str k))
                   (rpl "{arities}" arities))
        intern-str (-> intern-template
                       (rpl "{nsFullName}" ns-name)
                       (rpl "{nsName}" ns-name-final)
                       (rpl "{fnName}" (str k))
                       (rpl "{goName}" go-fn-name)
                       (rpl "{fnDocstring}" (raw-quoted-string (:doc m)))
                       (rpl "{added}" (:added m))
                       (rpl "{moreMeta}" (add-other-meta m))
                       (rpl "{args}"
                            (str "NewListFrom("
                                 (s/join ", " (for [args arglists]
                                                (generate-arglist args)))
                                 ")")))]
    [fn-str intern-str]))

(defn go-return-type
  "Returns the return type of the Make<t>() function. Would be unnecessary if Go code could declare a var as having 'the type returned by <func>'."
  [t]
  (condp = t
    "BigInt" "*BigInt"
    "BigIntU" "*BigInt"
    "Number" "*BigInt"
    "StringVector" "*Vector"
    "Error" "String"
    t))

(defn generate-const-or-var-decl
  [name m]
  (let [type (type-name (:tag m))]
    (if (= type "Var")
      (format "var %s *GoVar" name)  ; Not yet supported by this version of Joker (see https://github.com/jcburley/lace/)
      (format "var %s %s" name (go-return-type type)))))

(defn generate-const-or-var-init
  [name m]
  (let [type (type-name (:tag m))]
    (if (= type "Var")
      (format "\t%s = &GoVar{Value: &%s}"  ; Get pointer to the actual var, not a copy of the var
              name
              (:go m))
      (format "\t%s = Make%s(%s)"
              name
              type
              (:go m)))))

(defn generate-non-fn-decl
  [ns-name ns-name-final k v]
  (let [m (meta v)
        go-non-fn-name (go-name (str k))
        non-fn-str (generate-const-or-var-decl go-non-fn-name m)
        intern-str (-> intern-template
                       (rpl "{nsFullName}" ns-name)
                       (rpl "{nsName}" ns-name-final)
                       (rpl "{fnName}" (str k))
                       (rpl "{goName}" go-non-fn-name)
                       (rpl "{fnDocstring}" (raw-quoted-string (:doc m)))
                       (rpl "{added}" (:added m))
                       (rpl "{moreMeta}" (add-other-meta m))
                       (rpl "{args}" "nil"))]
    [non-fn-str intern-str]))

(defn generate-non-fn-init
  [ns-name-final k v]
  (let [m (meta v)
        go-non-fn-name (go-name (str k))
        non-fn-str (generate-const-or-var-init go-non-fn-name m)]
    non-fn-str))

(defn comment-out
  [s]
  (-> s
      (rpl "\n// " "\n")
      (rpl "\n" "\n//")
      (rpl "\n// package" "\npackage")))

(defn compare-imports
  [^String l ^String r]
  (cond
    (s/starts-with? l ". ") (if (s/starts-with? r ". ")
                              (compare l r)
                              -1)
    (s/starts-with? r ". ") 1
    :else (compare l r)))

(defn- ns-public-fns
  "Return only publics that are functions."
  [ns]
  (remove #(:macro (meta (val %)))
          (filter #(:arglists (meta (val %))) (ns-publics ns))))

(defn- ns-public-non-fns
  "Return only publics that are not functions."
  [ns]
  (remove #(:macro (meta (val %)))
          (remove #(:arglists (meta (val %))) (ns-publics ns))))

(defn- ns-public-go-fns
  "Return only publics that are functions and have additional Go-specific metadata."
  [ns]
  (filter #(:go (meta (val %))) (ns-public-fns ns)))

(defn- ns-public-go-non-fns
  "Return only publics that are not functions and have additional Go-specific metadata."
  [ns]
  (filter #(:go (meta (val %))) (ns-public-non-fns ns)))

(defn- warn-about-skipped-publics
  [skipped]
  (when (> (count skipped) 0)
    (println-err "WARNING: skipping publics that are not Go-calling functions or vars:" skipped)))

(defn generate-ns
  [ns-sym ns-name ns-name-final]
  (let [ns (find-ns ns-sym)
        m (meta ns)
        go-non-fns (sort-by first (ns-public-go-non-fns ns))
        go-fns (sort-by first (ns-public-go-fns ns))
        fn-decls (for [[k v] go-fns]
                  (generate-fn-decl ns-name ns-name-final k v))
        non-fn-decls (for [[k v] go-non-fns]
                       (generate-non-fn-decl ns-name ns-name-final k v))
        non-fn-inits (for [[k v] go-non-fns]
                       (generate-non-fn-init ns-name-final k v))
        res (-> package-template
                (rpl "{nsFullName}" ns-name)
                (rpl "{nsName}" ns-name-final)
                (rpl "{imports}"
                     (s/join "\n\t" (sort compare-imports (conj
                                                           (mapv q (:go-imports m))
                                                           ". \"github.com/lab47/lace/core\""))))
                (rpl "{non-fn-decls}" (s/join "\n" (map first non-fn-decls)))
                (rpl "{non-fn-inits}" (s/join "\n" non-fn-inits))
                (rpl "{fn-decls}" (s/join "\n" (map first fn-decls))))
        res (if (:empty m)
              (comment-out res)
              res)]
    (warn-about-skipped-publics (remove (set (concat (map key go-fns)
                                                     (map key go-non-fns)))
                                        (map #(key %) (ns-publics ns))))
    res))

(defn generate-ns-slow-init
  [ns-sym ns-name ns-name-final pre]
  (let [ns (find-ns ns-sym)
        m (meta ns)
        go-non-fns (sort-by first (ns-public-go-non-fns ns))
        go-fns (sort-by first (ns-public-go-fns ns))
        fn-decls (for [[k v] go-fns]
                  (generate-fn-decl ns-name ns-name-final k v))
        non-fn-decls (for [[k v] go-non-fns]
                       (generate-non-fn-decl ns-name ns-name-final k v))
        non-fn-inits (for [[k v] go-non-fns]
                       (generate-non-fn-init ns-name-final k v))
        res (-> package-slow-template
                (rpl "{maybeSlowOnly}" (if pre "// +build !fast_init\n" ""))
                (rpl "{nsFullName}" ns-name)
                (rpl "{nsName}" ns-name-final)
                (rpl "{imports}" ". \"github.com/lab47/lace/core\"")
                (rpl "{nsDocstring}" (raw-quoted-string (:doc m)))
                (rpl "{non-fn-interns}" (s/join "\n\t" (map second non-fn-decls)))
                (rpl "{fn-interns}" (s/join "\n\t" (map second fn-decls))))
        res (if (:empty m)
              (comment-out res)
              res)]
    res))

(defn generate-thunk
  [ns-name ns-name-final k v]
  (let [fn (go-name (str k))]
    (format "STD_thunk_%s_%s_var = __%s" ns-name fn fn)))

(defn generate-ns-fast-init
  [ns-sym ns-name ns-name-final]
  (let [ns (find-ns ns-sym)
        m (meta ns)
        thunks (for [[k v] (sort-by first (ns-public-go-fns ns))]
                 (generate-thunk ns-name ns-name-final k v))
        res (-> package-fast-template
                (rpl "{nsName}" ns-name-final)
                (rpl "{imports}" ". \"github.com/lab47/lace/core\"")
                (rpl "{thunks}" (s/join "\n\t" thunks)))
        res (if (:empty m)
              (comment-out res)
              res)]
    res))

(defn ns-file-name
  [dir ns-name-final]
  (str dir "/a_" ns-name-final ".go"))

(defn ns-file-name-slow-init
  [dir ns-name-final]
  (str dir "/a_" ns-name-final "_slow_init.go"))

(defn ns-file-name-fast-init
  [dir ns-name-final]
  (str dir "/a_" ns-name-final "_fast_init.go"))

(defn remove-blanky-lines
  [s]
  (-> s
      (rpl #"[[:space:]]*{blank}" "")))

(defn lace-ize
  "Convert a single-element namespace symbol to lace.namespace"
  ^Symbol [^Symbol s]
  (symbol (str "lace." s)))

(doseq [ns-sym namespaces]
  (let [ns-name (str ns-sym)
        dir (rpl ns-name "." "/")
        ns-name-final (rpl ns-name #".*[.]" "")
        pre (preloaded (lace-ize ns-sym))]
    (spit (ns-file-name dir ns-name-final)
          (remove-blanky-lines (generate-ns ns-sym ns-name ns-name-final)))
    (spit (ns-file-name-slow-init dir ns-name-final)
          (remove-blanky-lines (generate-ns-slow-init ns-sym ns-name ns-name-final pre)))
    (when pre
      (spit (ns-file-name-fast-init dir ns-name-final)
           (remove-blanky-lines (generate-ns-fast-init ns-sym ns-name ns-name-final))))))
