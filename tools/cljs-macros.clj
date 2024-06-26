(load-file "../core/data/linter_cljx.clj")
(load-file "../core/data/linter_cljs.clj")
(def interns (ns-interns 'lace.core))

(defn exists?
  [line]
  (let [parts (lace.string/split (lace.string/trim-space line) #" ")
        name (second (rest parts))]
    (if name
      (get interns (symbol name))
      false)))

(let [input (slurp "cljs-macros.input")
      lines (lace.string/split-lines input)]
  (doseq [line (remove exists? lines)]
    (println line)))
