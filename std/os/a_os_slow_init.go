// This file is generated by generate-std.clj script. Do not edit manually!


package os

import (
	. "github.com/lab47/lace/core"
	"fmt"
	"os"
)

func InternsOrThunks(env *Env, ns *Namespace) {
	if VerbosityLevel > 0 {
		fmt.Fprintln(os.Stderr, "Lazily running slow version of os.InternsOrThunks().")
	}
	ns.ResetMeta(MakeMeta(nil, `Provides a platform-independent interface to operating system functionality.`, "1.0"))

	
	ns.InternVar(env, "args", args_,
		MakeMeta(
			NewListFrom(NewVectorFrom()),
			`Returns a sequence of the command line arguments, starting with the program name (normally, lace).`, "1.0"))

	ns.InternVar(env, "chdir", chdir_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("dirname"))),
			`Chdir changes the current working directory to the named directory. If there is an error, an exception will be thrown. Returns nil.`, "1.0"))

	ns.InternVar(env, "close", close_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("f"))),
			`Closes the file, rendering it unusable for I/O.`, "1.0"))

	ns.InternVar(env, "create", create_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("name"))),
			`Creates the named file with mode 0666 (before umask), truncating it if it already exists.`, "1.0").Plus(env, MakeKeyword("tag"), String{S: "File"}))

	ns.InternVar(env, "cwd", cwd_,
		MakeMeta(
			NewListFrom(NewVectorFrom()),
			`Returns a rooted path name corresponding to the current directory. If the current directory can
  be reached via multiple paths (due to symbolic links), cwd may return any one of them.`, "1.0").Plus(env, MakeKeyword("tag"), String{S: "String"}))

	ns.InternVar(env, "env", env_,
		MakeMeta(
			NewListFrom(NewVectorFrom()),
			`Returns a map representing the environment.`, "1.0"))

	ns.InternVar(env, "exec", exec_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("name"), MakeSymbol("opts"))),
			`Executes the named program with the given arguments. opts is a map with the following keys (all optional):
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
  :err - string capturing stderr of the program (unless :stderr option was passed).`, "1.0"))

	ns.InternVar(env, "exists?", isexists_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("path"))),
			`Returns true if file or directory with the given path exists. Otherwise returns false.`, "1.0").Plus(env, MakeKeyword("tag"), String{S: "Boolean"}))

	ns.InternVar(env, "exit", exit_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("code"))),
			`Causes the current program to exit with the given status code.`, "1.0"))

	ns.InternVar(env, "get-env", get_env_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("key"))),
			`Returns the value of the environment variable named by the key or nil if the variable is not present in the environment.`, "1.0"))

	ns.InternVar(env, "ls", ls_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("dirname"))),
			`Reads the directory named by dirname and returns a list of directory entries sorted by filename.`, "1.0"))

	ns.InternVar(env, "mkdir", mkdir_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("name"), MakeSymbol("perm"))),
			`Creates a new directory with the specified name and permission bits.`, "1.0"))

	ns.InternVar(env, "open", open_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("name"))),
			`Opens the named file for reading. If successful, the file can be used for reading;
  the associated file descriptor has mode O_RDONLY.`, "1.0").Plus(env, MakeKeyword("tag"), String{S: "File"}))

	ns.InternVar(env, "remove", remove_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("name"))),
			`Removes the named file or (empty) directory.`, "1.0"))

	ns.InternVar(env, "remove-all", remove_all_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("path"))),
			`Removes path and any children it contains.

  It removes everything it can, then panics with the first error (if
  any) it encountered.`, "1.0"))

	ns.InternVar(env, "set-env", set_env_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("key"), MakeSymbol("value"))),
			`Sets the specified key to the specified value in the environment.`, "1.0"))

	ns.InternVar(env, "sh", sh_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("name"), MakeSymbol("&"), MakeSymbol("arguments"))),
			`Executes the named program with the given arguments. Returns a map with the following keys:
      :success - whether or not the execution was successful,
      :err-msg (present iff :success if false) - string capturing error object returned by Go runtime
      :exit - exit code of program (or attempt to execute it),
      :out - string capturing stdout of the program,
      :err - string capturing stderr of the program.`, "1.0"))

	ns.InternVar(env, "sh-from", sh_from_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("dir"), MakeSymbol("name"), MakeSymbol("&"), MakeSymbol("arguments"))),
			`Executes the named program with the given arguments and working directory set to dir.
  Returns a map with the following keys:
      :success - whether or not the execution was successful,
      :err-msg (present iff :success if false) - string capturing error object returned by Go runtime
      :exit - exit code of program (or attempt to execute it),
      :out - string capturing stdout of the program,
      :err - string capturing stderr of the program.`, "1.0"))

	ns.InternVar(env, "stat", stat_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("filename"))),
			`Returns a map describing the named file. The info map has the following attributes:
  :name - base name of the file
  :size - length in bytes for regular files; system-dependent for others
  :mode - file mode bits
  :modtime - modification time
  :dir? - true if file is a directory`, "1.0"))

}
