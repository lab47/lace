(def exit-code 0)

(defn clean
  [output]
  (let [output (if (nil? output) "" output)]
    (->> output
         (lace.string/split-lines)
         (remove #(lace.string/starts-with? % "  "))
         (remove #(= % ""))
         (lace.string/join "\n"))))

(defn test-flags
  [out description flags expected]
  (let [pwd (get (lace.os/env) "PWD")
        flag-parts (lace.string/split flags #"\s+")
        stdin? (some #(= "<" %) flag-parts)
        cmd (str pwd "/lace")
        output (if stdin?
                 (out (lace.os/sh "sh" "-c" (str cmd " " flags)))
                 (out (apply (partial lace.os/sh cmd) flag-parts)))
        output (clean output)]
    (when-not (= output expected)
      (println "FAILED: testing" description "(" flags ")")
      (println "EXPECTED")
      (println expected)
      (println "ACTUAL")
      (println output)
      (println "")
      (var-set #'exit-code 1))))

(defn testing
  [out description & tests]
  (let [tests (partition 2 tests)]
    (doseq [[flags expected] tests]
      (test-flags out description flags expected))))

(testing :err "auto detect dialect from filename"
  "--lint tests/flags/input.clj"
  ""
  "--lint tests/flags/input-warning.clj"
  "tests/flags/input-warning.clj:1:7: Parse warning: unused binding: a")

(testing :err "specify lace dialect"
  "--lintlace tests/flags/input.clj"
  ""

  "--lint --dialect lace tests/flags/input.clj"
  ""

  "--lintlace tests/flags/input.clj"
  "tests/flags/input.clj:1:2: Parse error: Unable to resolve symbol: clojure.string/split"

  "--lint --dialect lace tests/flags/input.cljs"
  "tests/flags/input.cljs:1:7: Parse error: Unable to resolve symbol: js/console")

(testing :err "specify clj dialect"
  "--lintclj tests/flags/input.clj"
  ""

  "--lint --dialect clj tests/flags/input.clj"
  ""

  "--lintclj tests/flags/input.edn"
  "tests/flags/input.edn:1:17: Read warning: No reader function for tag foo/bar"

  "--lint --dialect clj tests/flags/input.cljs"
  "tests/flags/input.cljs:1:7: Parse error: Unable to resolve symbol: js/console")

(testing :err "specify cljs dialect"
  "--lintcljs tests/flags/input.cljs"
  ""

  "--lint --dialect cljs tests/flags/input.cljs"
  ""

  "--lintcljs tests/flags/input.edn"
  "tests/flags/input.edn:1:17: Read warning: No reader function for tag foo/bar"

  "--lintcljs tests/flags/input.clj"
  "tests/flags/input.clj:1:2: Parse error: Unable to resolve symbol: clojure.string/split")

(testing :err "reading from stdin"
  "--lint --dialect edn - < tests/flags/input.edn"
  ""

  "--lint --dialect clj - < tests/flags/input.edn"
  "<stdin>:1:17: Read warning: No reader function for tag foo/bar"

  "--lint --dialect clj - < tests/flags/input.clj"
  ""

  "--lint --dialect cljs - < tests/flags/input.clj"
  "<stdin>:1:2: Parse error: Unable to resolve symbol: clojure.string/split"

  "--lint --dialect cljs - < tests/flags/input.cljs"
  ""

  "--lint --dialect lace - < tests/flags/input.cljs"
  "<stdin>:1:7: Parse error: Unable to resolve symbol: js/console"

  "--lint --dialect lace - < tests/flags/input.clj"
  "")

(testing :err "working directory override"
  "--lint --dialect clj - < tests/flags/macro.clj"
  "<stdin>:4:11: Parse error: Unable to resolve symbol: something"

  "--lint --dialect clj --working-dir tests/flags/config - < tests/flags/macro.clj"
  "")

(testing :out "script args don't cause errors"
  "tests/flags/script-flags.clj -go-style-flag -otherflag"
  "[-go-style-flag -otherflag]"

  "tests/flags/script-flags.clj --unix-style-flag --otherflag"
  "[--unix-style-flag --otherflag]"

  "tests/flags/script-flags.clj --unix-style-flag-with-value foobar"
  "[--unix-style-flag-with-value foobar]"

  "tests/flags/script-flags.clj --unix-style-flag-with-value=foobar"
  "[--unix-style-flag-with-value=foobar]"

  "tests/flags/script-flags.clj --short-hand-flags -a -b -c"
  "[--short-hand-flags -a -b -c]"

  "tests/flags/script-flags.clj --short-hand-flags -abc"
  "[--short-hand-flags -abc]"

  "tests/flags/script-flags.clj -- something that is not a flag"
  "[-- something that is not a flag]")

(testing :err "negative numbers parsed correctly"
         "--hashmap-threshold -1 tests/flags/input.clj"
         "")

(lace.os/exit exit-code)
