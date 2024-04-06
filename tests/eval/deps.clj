(ns deps
  (:require [lace.os :as os]
            [lace.filepath :as fp]
            [lace.string :as str]))

(def lib-dir
  (-> *main-file*
       (str/split fp/separator)
       (butlast)
       (concat ["lib"])
       ((fn [x] (apply str (interpose fp/separator x))))))

(ns-sources
  {"test-local.*" {:url lib-dir}})
