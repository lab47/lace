(def *logger* (default))

(defmacro log [& args]
  (let [level (first args)
        msg (second args)
        kvs (rest (rest args))]
    `(emit *logger* ~level ~msg (list ~@kvs))))

(defmacro trace [& args] `(log :trace ~@args))
(defmacro debug [& args] `(log :debug ~@args))
(defmacro info [& args] `(log :info ~@args))
(defmacro warn [& args] `(log :warn ~@args))
(defmacro error [& args] `(log :error ~@args))

(defn set-level
  "Set the level of the logger"
  ([level] (set-logger-level *logger* level))
  ([logger level] (set-logger-level logger level)))
