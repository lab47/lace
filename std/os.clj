(ns
  ^{:go-imports ["os"]
    :doc "Provides a platform-independent interface to operating system functionality."}
  os)

(defn env
  "Returns a map representing the environment."
  {:added "1.0"
   :go "env()"}
  [])

(defn set-env
  "Sets the specified key to the specified value in the environment."
  {:added "1.0"
   :go "setEnv(key, value)"}
  [^String key ^String value])

(defn get-env
  "Returns the value of the environment variable named by the key or nil if the variable is not present in the environment."
  {:added "1.0"
   :go "getEnv(key)"}
  [^String key])

(defn args
  "Returns a sequence of the command line arguments, starting with the program name (normally, lace)."
  {:added "1.0"
   :go "commandArgs()"}
  [])

(defn exit
  "Causes the current program to exit with the given status code."
  {:added "1.0"
   :go "NIL, nil; ExitJoker(code)"}
  [^Int code])

(defn sh
  "Executes the named program with the given arguments. Returns a map with the following keys:
      :success - whether or not the execution was successful,
      :err-msg (present iff :success if false) - string capturing error object returned by Go runtime
      :exit - exit code of program (or attempt to execute it),
      :out - string capturing stdout of the program,
      :err - string capturing stderr of the program."
  {:added "1.0"
   :go "sh(\"\", nil, nil, nil, name, arguments)"}
  [^String name & ^String arguments])

(defn sh-from
  "Executes the named program with the given arguments and working directory set to dir.
  Returns a map with the following keys:
      :success - whether or not the execution was successful,
      :err-msg (present iff :success if false) - string capturing error object returned by Go runtime
      :exit - exit code of program (or attempt to execute it),
      :out - string capturing stdout of the program,
      :err - string capturing stderr of the program."
  {:added "1.0"
   :go "sh(dir, nil, nil, nil, name, arguments)"}
  [^String dir ^String name & ^String arguments])

(defn exec
  "Executes the named program with the given arguments. opts is a map with the following keys (all optional):
  :args - vector of arguments (all arguments must be strings),
  :dir - if specified, working directory will be set to this value before executing the program,
  :stdin - if specified, provides stdin for the program. Can be either a string or an IOReader.
  If it's a string, the string's content will serve as stdin for the program. IOReader can be, for example,
  *in* (in which case Joker's stdin will be redirected to the program's stdin) or the value returned by (lace.os/open).
  :stdout - if specified, must be an IOWriter. It can be, for example, *out* (in which case the program's stdout will be redirected
  to Joker's stdout) or the value returned by (lace.os/create).
  :stderr - the same as :stdout, but for stderr.
  Returns a map with the following keys:
  :success - whether or not the execution was successful,
  :err-msg (present iff :success if false) - string capturing error object returned by Go runtime
  :exit - exit code of program (or attempt to execute it),
  :out - string capturing stdout of the program (unless :stdout option was passed)
  :err - string capturing stderr of the program (unless :stderr option was passed)."
  {:added "1.0"
   :go "execute(_env, name, opts)"}
  [^String name ^Map opts])

(defn mkdir
  "Creates a new directory with the specified name and permission bits."
  {:added "1.0"
   :go "mkdir(name, perm)"}
  [^String name ^Int perm])

(defn ls
  "Reads the directory named by dirname and returns a list of directory entries sorted by filename."
  {:added "1.0"
   :go "readDir(dirname)"}
  [^String dirname])

(defn ^String cwd
  "Returns a rooted path name corresponding to the current directory. If the current directory can
  be reached via multiple paths (due to symbolic links), cwd may return any one of them."
  {:added "1.0"
   :go "getwd()"}
  [])

(defn chdir
  "Chdir changes the current working directory to the named directory. If there is an error, an exception will be thrown. Returns nil."
  {:added "1.0"
   :go "chdir(dirname)"}
  [^String dirname])

(defn stat
  "Returns a map describing the named file. The info map has the following attributes:
  :name - base name of the file
  :size - length in bytes for regular files; system-dependent for others
  :mode - file mode bits
  :modtime - modification time
  :dir? - true if file is a directory"
  {:added "1.0"
   :go "stat(filename)"}
  [^String filename])

(defn ^Boolean exists?
  "Returns true if file or directory with the given path exists. Otherwise returns false."
  {:added "1.0"
   :go "exists(path)"}
  [^String path])

(defn ^File open
  "Opens the named file for reading. If successful, the file can be used for reading;
  the associated file descriptor has mode O_RDONLY."
  {:added "1.0"
   :go "os.Open(name)"}
  [^String name])

(defn ^File create
  "Creates the named file with mode 0666 (before umask), truncating it if it already exists."
  {:added "1.0"
   :go "os.Create(name)"}
  [^String name])

(defn close
  "Closes the file, rendering it unusable for I/O."
  {:added "1.0"
   :go "NIL, f.Close()"}
  [^File f])

(defn remove
  "Removes the named file or (empty) directory."
  {:added "1.0"
   :go "NIL, os.Remove(name)"}
  [^String name])

(defn remove-all
  "Removes path and any children it contains.

  It removes everything it can, then panics with the first error (if
  any) it encountered."
  {:added "1.0"
   :go "NIL, os.RemoveAll(path)"}
  [^String path])