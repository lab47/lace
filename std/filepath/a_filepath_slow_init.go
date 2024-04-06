// This file is generated by generate-std.clj script. Do not edit manually!

package filepath

import (
	"fmt"
	. "github.com/lab47/lace/core"
	"os"
)

func InternsOrThunks() {
	if VerbosityLevel > 0 {
		fmt.Fprintln(os.Stderr, "Lazily running slow version of filepath.InternsOrThunks().")
	}
	filepathNamespace.ResetMeta(MakeMeta(nil, `Implements utility routines for manipulating filename paths.`, "1.0"))

	filepathNamespace.InternVar("list-separator", list_separator_,
		MakeMeta(
			nil,
			`OS-specific path list separator.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	filepathNamespace.InternVar("separator", separator_,
		MakeMeta(
			nil,
			`OS-specific path separator.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	filepathNamespace.InternVar("abs", abs_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("path"))),
			`Returns an absolute representation of path. If the path is not absolute it will be
  joined with the current working directory to turn it into an absolute path.
  The absolute path name for a given file is not guaranteed to be unique.
  Calls clean on the result.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	filepathNamespace.InternVar("abs?", isabs_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("path"))),
			`Reports whether the path is absolute.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Boolean"}))

	filepathNamespace.InternVar("base", base_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("path"))),
			`Returns the last element of path. Trailing path separators are removed before
  extracting the last element. If the path is empty, returns ".". If the path consists
  entirely of separators, returns a single separator.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	filepathNamespace.InternVar("clean", clean_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("path"))),
			`Returns the shortest path name equivalent to path by purely lexical processing.
  Applies the following rules iteratively until no further processing can be done:

1. Replace multiple separator elements with a single one.
2. Eliminate each . path name element (the current directory).
3. Eliminate each inner .. path name element (the parent directory)
   along with the non-.. element that precedes it.
4. Eliminate .. elements that begin a rooted path:
   that is, replace "/.." by "/" at the beginning of a path,
   assuming separator is '/'.
The returned path ends in a slash only if it represents a root directory, such as "/" on Unix or `+"`"+`C:\`+"`"+` on Windows.

Finally, any occurrences of slash are replaced by separator.

If the result of this process is an empty string, returns the string ".".`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	filepathNamespace.InternVar("dir", dir_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("path"))),
			`Returns all but the last element of path, typically the path's directory.
  After dropping the final element, calls clean on the path and trailing slashes are removed.
  If the path is empty, returns ".". If the path consists entirely of separators,
  returns a single separator. The returned path does not end in a separator unless it is the root directory.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	filepathNamespace.InternVar("eval-symlinks", eval_symlinks_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("path"))),
			`Returns the path name after the evaluation of any symbolic links. If path is relative the result will be
  relative to the current directory, unless one of the components is an absolute symbolic link.
  Calls clean on the result.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	filepathNamespace.InternVar("ext", ext_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("path"))),
			`Returns the file name extension used by path. The extension is the suffix beginning at the final dot
  in the final element of path; it is empty if there is no dot.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	filepathNamespace.InternVar("file-seq", file_seq_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("root"))),
			`Returns a seq of maps with info about files or directories under root.`, "1.0"))

	filepathNamespace.InternVar("from-slash", from_slash_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("path"))),
			`Returns the result of replacing each slash ('/') character in path with a separator character.
  Multiple slashes are replaced by multiple separators.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	filepathNamespace.InternVar("glob", glob_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("pattern"))),
			`Returns the names of all files matching pattern or nil if there is no matching file.
  The syntax of patterns is the same as in Match. The pattern may describe hierarchical
  names such as /usr/*/bin/ed (assuming the separator is '/').

  Ignores file system errors such as I/O errors reading directories.
  Throws exception when pattern is malformed.`, "1.0").Plus(MakeKeyword("tag"), String{S: "[String]"}))

	filepathNamespace.InternVar("join", join_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("&"), MakeSymbol("elems"))),
			`Joins any number of path elements into a single path, adding a separator if necessary.
  Calls clean on the result; in particular, all empty strings are ignored. On Windows,
  the result is a UNC path if and only if the first path element is a UNC path.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	filepathNamespace.InternVar("matches?", ismatches_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("pattern"), MakeSymbol("name"))),
			`Reports whether name matches the shell file name pattern.
  Requires pattern to match all of name, not just a substring.
  Throws exception if pattern is malformed.
  On Windows, escaping is disabled. Instead, '\' is treated as path separator.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Boolean"}))

	filepathNamespace.InternVar("rel", rel_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("basepath"), MakeSymbol("targpath"))),
			`Returns a relative path that is lexically equivalent to targpath when joined to basepath
  with an intervening separator. On success, the returned path will always be relative to basepath,
  even if basepath and targpath share no elements. An exception is thrown if targpath can't be made
  relative to basepath or if knowing the current working directory would be necessary to compute it.
  Calls clean on the result.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	filepathNamespace.InternVar("split", split_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("path"))),
			`Splits path immediately following the final separator, separating it into a directory and file name component.
  If there is no separator in path, returns an empty dir and file set to path. The returned values have
  the property that path = dir+file.`, "1.0"))

	filepathNamespace.InternVar("split-list", split_list_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("path"))),
			`Splits a list of paths joined by the OS-specific list-separator, usually found in PATH or GOPATH environment variables.
  Returns an empty slice when passed an empty string.`, "1.0").Plus(MakeKeyword("tag"), String{S: "[String]"}))

	filepathNamespace.InternVar("to-slash", to_slash_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("path"))),
			`Returns the result of replacing each separator character in path with a slash ('/') character.
  Multiple separators are replaced by multiple slashes.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

	filepathNamespace.InternVar("volume-name", volume_name_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("path"))),
			`Returns leading volume name. Given "C:\foo\bar" it returns "C:" on Windows. Given "\\host\share\foo"
  returns "\\host\share". On other platforms it returns "".`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

}
