(defn f1 [])
(defn f2 [x])
(def v1 1)
;; Should PASS

(map identity)
(f1)
(f2 1)

;; Should FAIL
(map)
(*assert*)
(f1 1)
(f2)
(v1)
