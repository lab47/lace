(in-ns 'tests.lib1)
(lace.core/refer 'lace.core)
(println "in lib1")

(require ['tests.lib2])

(def v1 1)

(defn f1
  [s]
  (count s))
