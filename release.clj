(alias 's 'lace.string)
(alias 'os 'lace.os)

(def ^:private version-pattern "version\\s+\"(.*)\"")
(def ^:private shasum-pattern "sha256\\s+\"(.*)\"")

(defn- exit-err
  [& args]
  (apply println-err args)
  (lace.os/exit 1))

(defn- exec
  [& args]
  (let [res (apply os/sh args)]
    (when-not (:success res)
      (exit-err (str "'" (s/join " " args) "'") "failed:" (:out res) (:err res)))
    (:out res)))

(defn- exec-from
  [dir & args]
  (let [res (apply os/sh-from dir args)]
    (when-not (:success res)
      (exit-err (str "'[" dir "] " (s/join " " args) "'") "failed:" (:out res) (:err res)))
    (:out res)))

(defn- shasum
  [filename]
  (let [res (exec "shasum" "-a" "256" filename)]
    (first (s/split res #"\s"))))

(defn- update-brew-formula
  [version]
  (println "Updating brew formula...")
  (let [version' (subs version 1)
        formula-file (get (os/env) "JOKER_FORMULA_PATH")
        formula (slurp formula-file)
        [_ current-version] (re-find (re-pattern version-pattern) formula)
        linux-sha (shasum (str "lace-" version' "-linux-amd64.zip"))
        mac-sha (shasum (str "lace-" version' "-mac-amd64.zip"))
        shasum-rex (re-pattern shasum-pattern)]
    (when-not current-version
      (exit-err "Could not find version number"))
    (spit formula-file (-> formula
                           (s/replace current-version version')
                           (s/replace shasum-rex (str "sha256 \"" mac-sha "\""))
                           (s/replace-first shasum-rex (str "sha256 \"" linux-sha "\""))))
    (let [res (os/sh-from (s/replace formula-file "/lace.rb" "")
                          "git" "commit" "-a" "-m" (str "lace-" version))]
      (when-not (:success res)
        (exit-err "git commit failed:" (:out res) (:err res))))))

(defn- update-version
  [version]
  (println (str "Updating version to " version "..."))
  (let [procs (slurp "core/procs.go")
        procs-lines (s/split-lines procs)
        version-line (first (filter #(re-find #"const VERSION = .*" %) procs-lines))]
    (when-not procs
      (exit-err "Could not find version line"))
    (spit "core/procs.go" (s/replace procs version-line (str "const VERSION = \"" version "\""))))
  (exec "git" "commit" "-a" "-m" version)
  (exec "git" "tag" version)
  (exec "git" "push")
  (exec "git" "push" "--tags"))

(defn- build-binaries
  [version]
  (println "Building binaries...")
  (exec "./build-all.sh" (subs version 1))
  (exec "rm" "-f" "lace.exe")
  (exec "rm" "-f" "lace"))

(defn- generate-docs
  []
  (exec-from "docs" "../lace" "generate-docs.clj"))

(defn- main
  [args]
  (let [[_ _ version] args]
    (when-not version
      (exit-err "No version provided"))
    (when-not (re-matches #"v\d\.\d{1,2}\.\d{1,2}" version)
      (exit-err "Invalid version provided:" version))
    (generate-docs)
    (update-version version)
    (build-binaries version)
    (update-brew-formula version)))

(main (lace.os/args))
