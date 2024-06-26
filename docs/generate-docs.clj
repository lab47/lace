(doseq [ns (remove #(= % 'user) lace.core/*core-namespaces*)] (require ns))

(alias 'cli 'lace.tools.cli)
(alias 'os 'lace.os)
(alias 's 'lace.string)

(def rpl s/replace)

(def index-template
  (slurp "templates/index.html"))

(def ns-template
  (slurp "templates/ns.html"))

(def var-template
  (slurp "templates/var.html"))

(def special-form-template
  (slurp "templates/special-form.html"))

(def namespace-template
  (slurp "templates/ns-summary.html"))

(def type-template
  (slurp "templates/type-summary.html"))

(def link-item-template
  (slurp "templates/link-item.html"))

(def usage-template
  (slurp "templates/usage.html"))

(defn type-name
  [v]
  (let [m (meta v)]
    (cond
      (not (bound? v)) "Object"
      (:special-form m) "Special form"
      (:macro m) "Macro"
      (= Fn (type @v)) "Function"
      (= Proc (type @v)) "Function"
      (:tag m) (str (:tag m))
      :else (str (type @v)))))

(defn link-item-doc
  [k]
  (s/replace link-item-template "{name}" k))

(defn usage
  [k m]
  (if (:special-form m)
    (let [examples (for [form (:forms m)]
                     (s/replace usage-template "{usage}" (str form)))]
      (s/join "" examples))
    (let [examples (for [arglist (:arglists m)]
                     (s/replace usage-template "{usage}" (str (apply list k arglist))))]
      (s/join "" examples))))

(defn- source-file
  [ns]
  (s/join "_" (rest (s/split (str ns) #"\."))))

(defn var-doc
  [k v]
  (let [m (meta v)
        ns (get m :ns "<internal>")
        full-name (str ns "/" (str k))]
    (when-not (or (:added m) (:private m))
      (println "WARNING: public var without added meta key: " full-name))
    (when-not (or (:doc m) (:private m))
      (println "WARNING: public var without doc meta key: " full-name))
    (-> var-template
        (s/replace "{id}" (str k))
        (s/replace "{name}" (str k))
        (s/replace "{type}" (type-name v))
        (s/replace "{usage}" (usage k m))
        (s/replace "{docstring}" (s/replace (lace.html/escape (or (:doc m) "<<<MISSING DOCUMENTATION>>>")) "\n" "<br>\n"))
        (s/replace "{added}" (str (:added m)))
        (s/replace
         "{source}"
         (if (:line m)
           (format "<a href=\"https://github.com/lab47/lace/blob/master/core/data/%s.clj#L%s\">source</a>"
                   (source-file (:ns m))
                   (str (:line m)))
           "")))))

(defn- first-line
  [s]
  (first (s/split s #"\n")))

(defn special-form-doc
  [name meta]
  (let [usage (let [examples (for [form (:forms meta)]
                               (s/replace usage-template "{usage}" (str form)))]
                (s/join "" examples))]
    (-> special-form-template
        (s/replace "{id}" name)
        (s/replace "{name}" name)
        (s/replace "{docstring}" (s/replace (lace.html/escape (:doc meta)) "\n" "<br>\n"))
        (s/replace "{usage}" usage))))

(defn namespace-doc
  [ns-sym]
  (let [ns (find-ns ns-sym)
        k (str (ns-name ns))
        m (meta ns)]
    (when-not (:added m)
      (println "WARNING: namespace without added meta key: " k))
    (when-not (:doc m)
      (println "WARNING: namespace without doc meta key: " k))
    (-> namespace-template
        (s/replace "{id}" k)
        (s/replace "{name}" k)
        (s/replace "{docstring}" (s/replace (lace.html/escape (first-line (:doc m))) "\n" "<br>\n"))
        (s/replace "{added}" (str (:added m))))))

(defn type-doc
  [k]
  (let [m (meta (get (lace.core/types__) k))]
    (when-not (:added m)
      (println "WARNING: type without added meta key: " k))
    (when-not (:doc m)
      (println "WARNING: type without doc meta key: " k))
    (-> type-template
        (s/replace "{id}" k)
        (s/replace "{name}" k)
        (s/replace "{docstring}" (s/replace (lace.html/escape (:doc m)) "\n" "<br>\n"))
        (s/replace "{added}" (str (:added m))))))

(defn ns-doc
  [ns-sym ns-vars-fn]
  (let [ns (find-ns ns-sym)
        m (meta ns)
        vars-doc (s/join
                  ""
                  (for [[k v] (sort (ns-vars-fn ns-sym))]
                    (var-doc k v)))
        var-links-doc (s/join
                       ""
                       (for [k (sort (keys (ns-vars-fn ns-sym)))]
                         (link-item-doc (str k))))]
    (-> ns-template
        (s/replace "{name}" (name ns-sym))
        (s/replace "{added}" (str (:added m)))
        (s/replace "{docstring}" (s/replace (lace.html/escape (:doc m)) "\n" "<br>\n"))
        (s/replace "{vars}" vars-doc)
        (s/replace "{index}" var-links-doc))))

(defn index-doc
  [special-forms namespaces types]
  (let [special-forms-docs (s/join
                            ""
                            (for [sf (sort (keys special-forms))]
                              (special-form-doc (str sf) (special-forms sf))))
        special-form-links-doc (s/join
                                ""
                                (->> (sort (keys special-forms))
                                     (map #(link-item-doc (str %)))))

        namespaces-docs (s/join
                         ""
                         (for [ns-sym namespaces]
                           (namespace-doc ns-sym)))
        ns-links-doc (s/join
                      ""
                      (->> namespaces
                           (map #(link-item-doc (str %)))))
        types-docs (s/join
                    ""
                    (for [t types]
                      (type-doc t)))
        type-links-doc (s/join
                        ""
                        (->> types
                             (map #(link-item-doc (str %)))))]
    (-> index-template
        (s/replace "{index-of-special-forms}" special-form-links-doc)
        (s/replace "{special-forms}" special-forms-docs)
        (s/replace "{index-of-namespaces}" ns-links-doc)
        (s/replace "{namespaces}" namespaces-docs)
        (s/replace "{index-of-types}" type-links-doc)
        (s/replace "{types}" types-docs))))

(defn full-doc
  [ns-vars-fn]
  (let [namespaces (->> (all-ns)
                        (map ns-name)
                        (remove #(= 'user %))
                        (sort))
        types (->> (lace.core/types__)
                   (map key)
                   (sort))
        special-forms lace.repl/special-doc-map]
    (spit "index.html" (index-doc special-forms namespaces types))
    (doseq [ns namespaces]
      (spit (str ns ".html") (ns-doc ns ns-vars-fn)))))

(let [opts (cli/parse-opts *command-line-args*
                           [
                            [nil "--all" "Include private as well as public members in documentation"]
                            ["-h" "--help"]
                            ])]
  (when-let [err (or (when (:help (:options opts)) "") (:errors opts))]
    (println (s/join "\n" err))
    (println "Usage:")
    (println (:summary opts))
    (os/exit 1))
  (let [ns-vars-fn (if (:all (:options opts))
                     ns-interns
                     ns-publics)]
    (full-doc ns-vars-fn)))
