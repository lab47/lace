(ns lace.test.run-eval-tests
  (require [lace.test :refer [run-tests]]
           [lace.os :as os]))

(let [res (os/exec "test" {:stdin *in* :args ["-t" "0"]})]
  (when (:success res)
    (os/set-env "TTY_TESTS" "1")))

(defn- slurp-or
  [path alt]
  (try
    (slurp path)
    (catch Error e
      alt)))

(defn- have-option
  "Quick and dirty options parser."
  [opt]
  (some #(= opt %) *command-line-args*))

(defn- run-forked-test
  [lace-cmd test-dir verbose?]
  (when verbose?
    (println (str "Running test in subdirectory " test-dir)))
  (let [dir (str "tests/eval/" test-dir "/")
        filename "input.clj"
        stdin (slurp-or (str dir "stdin.txt") *in*)
        res (lace.os/exec lace-cmd {:dir dir :args [filename] :stdin stdin})
        out (:out res)
        err (:err res)
        rc (:exit res)
        expected-out (slurp-or (str dir "stdout.txt") "")
        expected-err (slurp-or (str dir "stderr.txt") "")
        expected-rc (if-let [rc (slurp-or (str dir "rc.txt") false)]
                      (int (bigint (with-in-str rc (read-line))))
                      0)]
    (if (and (= expected-out out) (= expected-err err) (= expected-rc rc))
      true
      (do
        (println "FAILED:" test-dir)
        (when-not (= expected-out out)
          (println "EXPECTED STDOUT:")
          (println "----------------")
          (println expected-out)
          (println "----------------")
          (newline)
          (println "ACTUAL STDOUT:")
          (println "----------------")
          (println out)
          (println "----------------"))
        (when-not (= expected-err err)
          (println "EXPECTED STDERR:")
          (println "----------------")
          (println expected-err)
          (println "----------------")
          (newline)
          (println "ACTUAL STDERR:")
          (println "----------------")
          (println err)
          (println "----------------"))
        (when-not (= expected-rc rc)
          (println "EXPECTED RC:" expected-rc)
          (println "ACTUAL RC:" rc))
        false))))

(defn- run-forked-tests
  "Directories have input.clj files that are passed as input to a
  fork'ed instance of Joker.  Ths test harness is rudimentary but
  handles more complicated cases than lace.test and can also catch
  test failures in that and its dependencies (such as defmulti,
  lace.template, and lace.walk) before a huge deluge of failures is
  reported by the per-file tests that are performed after these."
  [verbose?]
  (let [test-dirs (->> (lace.os/ls "tests/eval")
                       (filter :dir?)
                       (map :name)
                       (remove #(lace.string/starts-with? % ".")))
        pwd (get (lace.os/env) "PWD")
        exe (str pwd "/lace")
        failures (->> test-dirs
                      (remove #(run-forked-test exe % verbose?))
                      (count))]
    failures))

(defn- run-internal-tests
  ([]
   ;; Run tests in all tests/eval/*.clj files:
   (run-internal-tests (->> (lace.os/ls "tests/eval")
                            (remove :dir?)
                            (map :name)
                            (filter #(lace.string/ends-with? % ".clj"))
                            (remove #(lace.string/starts-with? % ".")))))
  ([test-files]
   (let [namespaces (for [test-file test-files]
                      (binding [*ns* *ns*]
                        (load-file (str "tests/eval/" test-file))
                        *ns*))]
     (let [res (apply run-tests namespaces)
           failures (+ (:fail res) (:error res))]
       failures))))

(defn- main
  []
  (let [verbose? (not (have-option "--no-verbose"))
        filename (first (filter #(lace.string/ends-with? % ".clj") *command-line-args*))
        forked-failures (if filename 0 (run-forked-tests verbose?))
        internal-failures (if filename
                            (run-internal-tests [filename])
                            (run-internal-tests))]
    (when (pos? forked-failures)
      (println (str "There were " forked-failures " failures and/or errors in forked tests; returning exit code 1")))
    (when (pos? internal-failures)
      (println "There were failures and/or errors in namespace tests; returning exit code 1"))
    (when (or (pos? forked-failures)
              (pos? internal-failures))
      (lace.os/exit 1))))

(main)
