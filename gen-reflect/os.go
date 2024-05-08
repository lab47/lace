package reflect

import (
	os "os"
	"reflect"

	"github.com/lab47/lace/pkg/pkgreflect"
)

type SignalImpl struct {
	SignalFn func()
	StringFn func() string
}

func (s *SignalImpl) Signal() {
	s.SignalFn()
}
func (s *SignalImpl) String() string {
	return s.StringFn()
}

func init() {
	DirEntry_methods := map[string]pkgreflect.Func{}
	PathError_methods := map[string]pkgreflect.Func{}
	SyscallError_methods := map[string]pkgreflect.Func{}
	ProcAttr_methods := map[string]pkgreflect.Func{}
	Process_methods := map[string]pkgreflect.Func{}
	Signal_methods := map[string]pkgreflect.Func{}
	ProcessState_methods := map[string]pkgreflect.Func{}
	LinkError_methods := map[string]pkgreflect.Func{}
	File_methods := map[string]pkgreflect.Func{}
	FileInfo_methods := map[string]pkgreflect.Func{}
	FileMode_methods := map[string]pkgreflect.Func{}
	File_methods["Readdir"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "n", Tag: "int"}}, Tag: "any", Doc: "Readdir reads the contents of the directory associated with file and\nreturns a slice of up to n FileInfo values, as would be returned\nby Lstat, in directory order. Subsequent calls on the same file will yield\nfurther FileInfos.\n\nIf n > 0, Readdir returns at most n FileInfo structures. In this case, if\nReaddir returns an empty slice, it will return a non-nil error\nexplaining why. At the end of a directory, the error is io.EOF.\n\nIf n <= 0, Readdir returns all the FileInfo from the directory in\na single slice. In this case, if Readdir succeeds (reads all\nthe way to the end of the directory), it returns the slice and a\nnil error. If it encounters an error before the end of the\ndirectory, Readdir returns the FileInfo read until that point\nand a non-nil error.\n\nMost clients are better served by the more efficient ReadDir method."}
	File_methods["Readdirnames"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "n", Tag: "int"}}, Tag: "any", Doc: "Readdirnames reads the contents of the directory associated with file\nand returns a slice of up to n names of files in the directory,\nin directory order. Subsequent calls on the same file will yield\nfurther names.\n\nIf n > 0, Readdirnames returns at most n names. In this case, if\nReaddirnames returns an empty slice, it will return a non-nil error\nexplaining why. At the end of a directory, the error is io.EOF.\n\nIf n <= 0, Readdirnames returns all the names from the directory in\na single slice. In this case, if Readdirnames succeeds (reads all\nthe way to the end of the directory), it returns the slice and a\nnil error. If it encounters an error before the end of the\ndirectory, Readdirnames returns the names read until that point and\na non-nil error."}
	File_methods["ReadDir"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "n", Tag: "int"}}, Tag: "any", Doc: "ReadDir reads the contents of the directory associated with the file f\nand returns a slice of DirEntry values in directory order.\nSubsequent calls on the same file will yield later DirEntry records in the directory.\n\nIf n > 0, ReadDir returns at most n DirEntry records.\nIn this case, if ReadDir returns an empty slice, it will return an error explaining why.\nAt the end of a directory, the error is io.EOF.\n\nIf n <= 0, ReadDir returns all the DirEntry records remaining in the directory.\nWhen it succeeds, it returns a nil error (not io.EOF)."}
	SyscallError_methods["Error"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: ""}
	SyscallError_methods["Unwrap"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "error", Doc: ""}
	SyscallError_methods["Timeout"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "bool", Doc: "Timeout reports whether this error represents a timeout."}
	Process_methods["Release"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "error", Doc: "Release releases any resources associated with the Process p,\nrendering it unusable in the future.\nRelease only needs to be called if Wait is not."}
	Process_methods["Kill"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "error", Doc: "Kill causes the Process to exit immediately. Kill does not wait until\nthe Process has actually exited. This only kills the Process itself,\nnot any other processes it may have started."}
	Process_methods["Wait"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "Wait waits for the Process to exit, and then returns a\nProcessState describing its status and an error, if any.\nWait releases any resources associated with the Process.\nOn most operating systems, the Process must be a child\nof the current process or an error will be returned."}
	Process_methods["Signal"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "sig", Tag: "Signal"}}, Tag: "error", Doc: "Signal sends a signal to the Process.\nSending Interrupt on Windows is not implemented."}
	ProcessState_methods["UserTime"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "time.Duration", Doc: "UserTime returns the user CPU time of the exited process and its children."}
	ProcessState_methods["SystemTime"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "time.Duration", Doc: "SystemTime returns the system CPU time of the exited process and its children."}
	ProcessState_methods["Exited"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "bool", Doc: "Exited reports whether the program has exited.\nOn Unix systems this reports true if the program exited due to calling exit,\nbut false if the program terminated due to a signal."}
	ProcessState_methods["Success"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "bool", Doc: "Success reports whether the program exited successfully,\nsuch as with exit status 0 on Unix."}
	ProcessState_methods["Sys"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "Sys returns system-dependent exit information about\nthe process. Convert it to the appropriate underlying\ntype, such as syscall.WaitStatus on Unix, to access its contents."}
	ProcessState_methods["SysUsage"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "SysUsage returns system-dependent resource usage information about\nthe exited process. Convert it to the appropriate underlying\ntype, such as *syscall.Rusage on Unix, to access its contents.\n(On Unix, *syscall.Rusage matches struct rusage as defined in the\ngetrusage(2) manual page.)"}
	ProcessState_methods["Pid"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int", Doc: "Pid returns the process id of the exited process."}
	ProcessState_methods["String"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: ""}
	ProcessState_methods["ExitCode"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int", Doc: "ExitCode returns the exit code of the exited process, or -1\nif the process hasn't exited or was terminated by a signal."}
	File_methods["Name"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: "Name returns the name of the file as presented to Open."}
	LinkError_methods["Error"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: ""}
	LinkError_methods["Unwrap"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "error", Doc: ""}
	File_methods["Read"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}}, Tag: "any", Doc: "Read reads up to len(b) bytes from the File and stores them in b.\nIt returns the number of bytes read and any error encountered.\nAt end of file, Read returns 0, io.EOF."}
	File_methods["ReadAt"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}, {Name: "off", Tag: "int64"}}, Tag: "any", Doc: "ReadAt reads len(b) bytes from the File starting at byte offset off.\nIt returns the number of bytes read and the error, if any.\nReadAt always returns a non-nil error when n < len(b).\nAt end of file, that error is io.EOF."}
	File_methods["ReadFrom"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "r", Tag: "io.Reader"}}, Tag: "any", Doc: "ReadFrom implements io.ReaderFrom."}
	File_methods["Write"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}}, Tag: "any", Doc: "Write writes len(b) bytes from b to the File.\nIt returns the number of bytes written and an error, if any.\nWrite returns a non-nil error when n != len(b)."}
	File_methods["WriteAt"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}, {Name: "off", Tag: "int64"}}, Tag: "any", Doc: "WriteAt writes len(b) bytes to the File starting at byte offset off.\nIt returns the number of bytes written and an error, if any.\nWriteAt returns a non-nil error when n != len(b).\n\nIf file was opened with the O_APPEND flag, WriteAt returns an error."}
	File_methods["WriteTo"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "w", Tag: "io.Writer"}}, Tag: "any", Doc: "WriteTo implements io.WriterTo."}
	File_methods["Seek"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "offset", Tag: "int64"}, {Name: "whence", Tag: "int"}}, Tag: "any", Doc: "Seek sets the offset for the next Read or Write on file to offset, interpreted\naccording to whence: 0 means relative to the origin of the file, 1 means\nrelative to the current offset, and 2 means relative to the end.\nIt returns the new offset and an error, if any.\nThe behavior of Seek on a file opened with O_APPEND is not specified."}
	File_methods["WriteString"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "s", Tag: "string"}}, Tag: "any", Doc: "WriteString is like Write, but writes the contents of string s rather than\na slice of bytes."}
	File_methods["Chmod"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "mode", Tag: "FileMode"}}, Tag: "error", Doc: "Chmod changes the mode of the file to mode.\nIf there is an error, it will be of type *PathError."}
	File_methods["SetDeadline"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "t", Tag: "time.Time"}}, Tag: "error", Doc: "SetDeadline sets the read and write deadlines for a File.\nIt is equivalent to calling both SetReadDeadline and SetWriteDeadline.\n\nOnly some kinds of files support setting a deadline. Calls to SetDeadline\nfor files that do not support deadlines will return ErrNoDeadline.\nOn most systems ordinary files do not support deadlines, but pipes do.\n\nA deadline is an absolute time after which I/O operations fail with an\nerror instead of blocking. The deadline applies to all future and pending\nI/O, not just the immediately following call to Read or Write.\nAfter a deadline has been exceeded, the connection can be refreshed\nby setting a deadline in the future.\n\nIf the deadline is exceeded a call to Read or Write or to other I/O\nmethods will return an error that wraps ErrDeadlineExceeded.\nThis can be tested using errors.Is(err, os.ErrDeadlineExceeded).\nThat error implements the Timeout method, and calling the Timeout\nmethod will return true, but there are other possible errors for which\nthe Timeout will return true even if the deadline has not been exceeded.\n\nAn idle timeout can be implemented by repeatedly extending\nthe deadline after successful Read or Write calls.\n\nA zero value for t means I/O operations will not time out."}
	File_methods["SetReadDeadline"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "t", Tag: "time.Time"}}, Tag: "error", Doc: "SetReadDeadline sets the deadline for future Read calls and any\ncurrently-blocked Read call.\nA zero value for t means Read will not time out.\nNot all files support setting deadlines; see SetDeadline."}
	File_methods["SetWriteDeadline"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "t", Tag: "time.Time"}}, Tag: "error", Doc: "SetWriteDeadline sets the deadline for any future Write calls and any\ncurrently-blocked Write call.\nEven if Write times out, it may return n > 0, indicating that\nsome of the data was successfully written.\nA zero value for t means Write will not time out.\nNot all files support setting deadlines; see SetDeadline."}
	File_methods["SyscallConn"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "SyscallConn returns a raw file.\nThis implements the syscall.Conn interface."}
	File_methods["Close"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "error", Doc: "Close closes the File, rendering it unusable for I/O.\nOn files that support SetDeadline, any pending I/O operations will\nbe canceled and return immediately with an ErrClosed error.\nClose will return an error if it has already been called."}
	File_methods["Chown"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "uid", Tag: "int"}, {Name: "gid", Tag: "int"}}, Tag: "error", Doc: "Chown changes the numeric uid and gid of the named file.\nIf there is an error, it will be of type *PathError.\n\nOn Windows, it always returns the syscall.EWINDOWS error, wrapped\nin *PathError."}
	File_methods["Truncate"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "size", Tag: "int64"}}, Tag: "error", Doc: "Truncate changes the size of the file.\nIt does not change the I/O offset.\nIf there is an error, it will be of type *PathError."}
	File_methods["Sync"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "error", Doc: "Sync commits the current contents of the file to stable storage.\nTypically, this means flushing the file system's in-memory copy\nof recently written data to disk."}
	File_methods["Chdir"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "error", Doc: "Chdir changes the current working directory to the file,\nwhich must be a directory.\nIf there is an error, it will be of type *PathError."}
	File_methods["Fd"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "uintptr", Doc: "Fd returns the integer Unix file descriptor referencing the open file.\nIf f is closed, the file descriptor becomes invalid.\nIf f is garbage collected, a finalizer may close the file descriptor,\nmaking it invalid; see runtime.SetFinalizer for more information on when\na finalizer might be run. On Unix systems this will cause the SetDeadline\nmethods to stop working.\nBecause file descriptors can be reused, the returned file descriptor may\nonly be closed through the Close method of f, or by its finalizer during\ngarbage collection. Otherwise, during garbage collection the finalizer\nmay close an unrelated file descriptor with the same (reused) number.\n\nAs an alternative, see the f.SyscallConn method."}
	File_methods["Stat"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "Stat returns the FileInfo structure describing file.\nIf there is an error, it will be of type *PathError."}
	pkgreflect.AddPackage("os", &pkgreflect.Package{
		Doc: "Package os provides a platform-independent interface to operating system functionality.",
		Types: map[string]pkgreflect.Type{
			"DirEntry":     {Doc: "", Value: reflect.TypeOf((*os.DirEntry)(nil)).Elem(), Methods: DirEntry_methods},
			"File":         {Doc: "", Value: reflect.TypeOf((*os.File)(nil)).Elem(), Methods: File_methods},
			"FileInfo":     {Doc: "", Value: reflect.TypeOf((*os.FileInfo)(nil)).Elem(), Methods: FileInfo_methods},
			"FileMode":     {Doc: "", Value: reflect.TypeOf((*os.FileMode)(nil)).Elem(), Methods: FileMode_methods},
			"LinkError":    {Doc: "", Value: reflect.TypeOf((*os.LinkError)(nil)).Elem(), Methods: LinkError_methods},
			"PathError":    {Doc: "", Value: reflect.TypeOf((*os.PathError)(nil)).Elem(), Methods: PathError_methods},
			"ProcAttr":     {Doc: "", Value: reflect.TypeOf((*os.ProcAttr)(nil)).Elem(), Methods: ProcAttr_methods},
			"Process":      {Doc: "", Value: reflect.TypeOf((*os.Process)(nil)).Elem(), Methods: Process_methods},
			"ProcessState": {Doc: "", Value: reflect.TypeOf((*os.ProcessState)(nil)).Elem(), Methods: ProcessState_methods},
			"Signal":       {Doc: "", Value: reflect.TypeOf((*os.Signal)(nil)).Elem(), Methods: Signal_methods},
			"SyscallError": {Doc: "", Value: reflect.TypeOf((*os.SyscallError)(nil)).Elem(), Methods: SyscallError_methods},
			"SignalImpl":   {Doc: `Struct version of interface Signal for implementation`, Value: reflect.TypeFor[SignalImpl]()},
		},

		Functions: map[string]pkgreflect.FuncValue{
			"Chdir": {Doc: "Chdir changes the current working directory to the named directory.\nIf there is an error, it will be of type *PathError.", Args: []pkgreflect.Arg{{Name: "dir", Tag: "string"}}, Tag: "error", Value: reflect.ValueOf(os.Chdir)},

			"Chmod": {Doc: "Chmod changes the mode of the named file to mode.\nIf the file is a symbolic link, it changes the mode of the link's target.\nIf there is an error, it will be of type *PathError.\n\nA different subset of the mode bits are used, depending on the\noperating system.\n\nOn Unix, the mode's permission bits, ModeSetuid, ModeSetgid, and\nModeSticky are used.\n\nOn Windows, only the 0200 bit (owner writable) of mode is used; it\ncontrols whether the file's read-only attribute is set or cleared.\nThe other bits are currently unused. For compatibility with Go 1.12\nand earlier, use a non-zero mode. Use mode 0400 for a read-only\nfile and 0600 for a readable+writable file.\n\nOn Plan 9, the mode's permission bits, ModeAppend, ModeExclusive,\nand ModeTemporary are used.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}, {Name: "mode", Tag: "FileMode"}}, Tag: "error", Value: reflect.ValueOf(os.Chmod)},

			"Chown": {Doc: "Chown changes the numeric uid and gid of the named file.\nIf the file is a symbolic link, it changes the uid and gid of the link's target.\nA uid or gid of -1 means to not change that value.\nIf there is an error, it will be of type *PathError.\n\nOn Windows or Plan 9, Chown always returns the syscall.EWINDOWS or\nEPLAN9 error, wrapped in *PathError.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}, {Name: "uid", Tag: "int"}, {Name: "gid", Tag: "int"}}, Tag: "error", Value: reflect.ValueOf(os.Chown)},

			"Chtimes": {Doc: "Chtimes changes the access and modification times of the named\nfile, similar to the Unix utime() or utimes() functions.\nA zero time.Time value will leave the corresponding file time unchanged.\n\nThe underlying filesystem may truncate or round the values to a\nless precise time unit.\nIf there is an error, it will be of type *PathError.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}, {Name: "atime", Tag: "time.Time"}, {Name: "mtime", Tag: "time.Time"}}, Tag: "error", Value: reflect.ValueOf(os.Chtimes)},

			"Clearenv": {Doc: "Clearenv deletes all environment variables.", Args: []pkgreflect.Arg{}, Tag: "any", Value: reflect.ValueOf(os.Clearenv)},

			"Create": {Doc: "Create creates or truncates the named file. If the file already exists,\nit is truncated. If the file does not exist, it is created with mode 0666\n(before umask). If successful, methods on the returned File can\nbe used for I/O; the associated file descriptor has mode O_RDWR.\nIf there is an error, it will be of type *PathError.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}}, Tag: "any", Value: reflect.ValueOf(os.Create)},

			"CreateTemp": {Doc: "CreateTemp creates a new temporary file in the directory dir,\nopens the file for reading and writing, and returns the resulting file.\nThe filename is generated by taking pattern and adding a random string to the end.\nIf pattern includes a \"*\", the random string replaces the last \"*\".\nIf dir is the empty string, CreateTemp uses the default directory for temporary files, as returned by TempDir.\nMultiple programs or goroutines calling CreateTemp simultaneously will not choose the same file.\nThe caller can use the file's Name method to find the pathname of the file.\nIt is the caller's responsibility to remove the file when it is no longer needed.", Args: []pkgreflect.Arg{{Name: "dir", Tag: "string"}, {Name: "pattern", Tag: "string"}}, Tag: "any", Value: reflect.ValueOf(os.CreateTemp)},

			"DirFS": {Doc: "DirFS returns a file system (an fs.FS) for the tree of files rooted at the directory dir.\n\nNote that DirFS(\"/prefix\") only guarantees that the Open calls it makes to the\noperating system will begin with \"/prefix\": DirFS(\"/prefix\").Open(\"file\") is the\nsame as os.Open(\"/prefix/file\"). So if /prefix/file is a symbolic link pointing outside\nthe /prefix tree, then using DirFS does not stop the access any more than using\nos.Open does. Additionally, the root of the fs.FS returned for a relative path,\nDirFS(\"prefix\"), will be affected by later calls to Chdir. DirFS is therefore not\na general substitute for a chroot-style security mechanism when the directory tree\ncontains arbitrary content.\n\nThe directory dir must not be \"\".\n\nThe result implements [io/fs.StatFS], [io/fs.ReadFileFS] and\n[io/fs.ReadDirFS].", Args: []pkgreflect.Arg{{Name: "dir", Tag: "string"}}, Tag: "fs.FS", Value: reflect.ValueOf(os.DirFS)},

			"Environ": {Doc: "Environ returns a copy of strings representing the environment,\nin the form \"key=value\".", Args: []pkgreflect.Arg{}, Tag: "[]string", Value: reflect.ValueOf(os.Environ)},

			"Executable": {Doc: "Executable returns the path name for the executable that started\nthe current process. There is no guarantee that the path is still\npointing to the correct executable. If a symlink was used to start\nthe process, depending on the operating system, the result might\nbe the symlink or the path it pointed to. If a stable result is\nneeded, path/filepath.EvalSymlinks might help.\n\nExecutable returns an absolute path unless an error occurred.\n\nThe main use case is finding resources located relative to an\nexecutable.", Args: []pkgreflect.Arg{}, Tag: "any", Value: reflect.ValueOf(os.Executable)},

			"Exit": {Doc: "Exit causes the current program to exit with the given status code.\nConventionally, code zero indicates success, non-zero an error.\nThe program terminates immediately; deferred functions are not run.\n\nFor portability, the status code should be in the range [0, 125].", Args: []pkgreflect.Arg{{Name: "code", Tag: "int"}}, Tag: "any", Value: reflect.ValueOf(os.Exit)},

			"Expand": {Doc: "Expand replaces ${var} or $var in the string based on the mapping function.\nFor example, os.ExpandEnv(s) is equivalent to os.Expand(s, os.Getenv).", Args: []pkgreflect.Arg{{Name: "s", Tag: "string"}, {Name: "mapping", Tag: "Unknown"}}, Tag: "string", Value: reflect.ValueOf(os.Expand)},

			"ExpandEnv": {Doc: "ExpandEnv replaces ${var} or $var in the string according to the values\nof the current environment variables. References to undefined\nvariables are replaced by the empty string.", Args: []pkgreflect.Arg{{Name: "s", Tag: "string"}}, Tag: "string", Value: reflect.ValueOf(os.ExpandEnv)},

			"FindProcess": {Doc: "FindProcess looks for a running process by its pid.\n\nThe Process it returns can be used to obtain information\nabout the underlying operating system process.\n\nOn Unix systems, FindProcess always succeeds and returns a Process\nfor the given pid, regardless of whether the process exists. To test whether\nthe process actually exists, see whether p.Signal(syscall.Signal(0)) reports\nan error.", Args: []pkgreflect.Arg{{Name: "pid", Tag: "int"}}, Tag: "any", Value: reflect.ValueOf(os.FindProcess)},

			"Getegid": {Doc: "Getegid returns the numeric effective group id of the caller.\n\nOn Windows, it returns -1.", Args: []pkgreflect.Arg{}, Tag: "int", Value: reflect.ValueOf(os.Getegid)},

			"Getenv": {Doc: "Getenv retrieves the value of the environment variable named by the key.\nIt returns the value, which will be empty if the variable is not present.\nTo distinguish between an empty value and an unset value, use LookupEnv.", Args: []pkgreflect.Arg{{Name: "key", Tag: "string"}}, Tag: "string", Value: reflect.ValueOf(os.Getenv)},

			"Geteuid": {Doc: "Geteuid returns the numeric effective user id of the caller.\n\nOn Windows, it returns -1.", Args: []pkgreflect.Arg{}, Tag: "int", Value: reflect.ValueOf(os.Geteuid)},

			"Getgid": {Doc: "Getgid returns the numeric group id of the caller.\n\nOn Windows, it returns -1.", Args: []pkgreflect.Arg{}, Tag: "int", Value: reflect.ValueOf(os.Getgid)},

			"Getgroups": {Doc: "Getgroups returns a list of the numeric ids of groups that the caller belongs to.\n\nOn Windows, it returns syscall.EWINDOWS. See the os/user package\nfor a possible alternative.", Args: []pkgreflect.Arg{}, Tag: "any", Value: reflect.ValueOf(os.Getgroups)},

			"Getpagesize": {Doc: "Getpagesize returns the underlying system's memory page size.", Args: []pkgreflect.Arg{}, Tag: "int", Value: reflect.ValueOf(os.Getpagesize)},

			"Getpid": {Doc: "Getpid returns the process id of the caller.", Args: []pkgreflect.Arg{}, Tag: "int", Value: reflect.ValueOf(os.Getpid)},

			"Getppid": {Doc: "Getppid returns the process id of the caller's parent.", Args: []pkgreflect.Arg{}, Tag: "int", Value: reflect.ValueOf(os.Getppid)},

			"Getuid": {Doc: "Getuid returns the numeric user id of the caller.\n\nOn Windows, it returns -1.", Args: []pkgreflect.Arg{}, Tag: "int", Value: reflect.ValueOf(os.Getuid)},

			"Getwd": {Doc: "Getwd returns a rooted path name corresponding to the\ncurrent directory. If the current directory can be\nreached via multiple paths (due to symbolic links),\nGetwd may return any one of them.", Args: []pkgreflect.Arg{}, Tag: "any", Value: reflect.ValueOf(os.Getwd)},

			"Hostname": {Doc: "Hostname returns the host name reported by the kernel.", Args: []pkgreflect.Arg{}, Tag: "any", Value: reflect.ValueOf(os.Hostname)},

			"IsExist": {Doc: "IsExist returns a boolean indicating whether the error is known to report\nthat a file or directory already exists. It is satisfied by ErrExist as\nwell as some syscall errors.\n\nThis function predates errors.Is. It only supports errors returned by\nthe os package. New code should use errors.Is(err, fs.ErrExist).", Args: []pkgreflect.Arg{{Name: "err", Tag: "error"}}, Tag: "bool", Value: reflect.ValueOf(os.IsExist)},

			"IsNotExist": {Doc: "IsNotExist returns a boolean indicating whether the error is known to\nreport that a file or directory does not exist. It is satisfied by\nErrNotExist as well as some syscall errors.\n\nThis function predates errors.Is. It only supports errors returned by\nthe os package. New code should use errors.Is(err, fs.ErrNotExist).", Args: []pkgreflect.Arg{{Name: "err", Tag: "error"}}, Tag: "bool", Value: reflect.ValueOf(os.IsNotExist)},

			"IsPathSeparator": {Doc: "IsPathSeparator reports whether c is a directory separator character.", Args: []pkgreflect.Arg{{Name: "c", Tag: "uint8"}}, Tag: "bool", Value: reflect.ValueOf(os.IsPathSeparator)},

			"IsPermission": {Doc: "IsPermission returns a boolean indicating whether the error is known to\nreport that permission is denied. It is satisfied by ErrPermission as well\nas some syscall errors.\n\nThis function predates errors.Is. It only supports errors returned by\nthe os package. New code should use errors.Is(err, fs.ErrPermission).", Args: []pkgreflect.Arg{{Name: "err", Tag: "error"}}, Tag: "bool", Value: reflect.ValueOf(os.IsPermission)},

			"IsTimeout": {Doc: "IsTimeout returns a boolean indicating whether the error is known\nto report that a timeout occurred.\n\nThis function predates errors.Is, and the notion of whether an\nerror indicates a timeout can be ambiguous. For example, the Unix\nerror EWOULDBLOCK sometimes indicates a timeout and sometimes does not.\nNew code should use errors.Is with a value appropriate to the call\nreturning the error, such as os.ErrDeadlineExceeded.", Args: []pkgreflect.Arg{{Name: "err", Tag: "error"}}, Tag: "bool", Value: reflect.ValueOf(os.IsTimeout)},

			"Lchown": {Doc: "Lchown changes the numeric uid and gid of the named file.\nIf the file is a symbolic link, it changes the uid and gid of the link itself.\nIf there is an error, it will be of type *PathError.\n\nOn Windows, it always returns the syscall.EWINDOWS error, wrapped\nin *PathError.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}, {Name: "uid", Tag: "int"}, {Name: "gid", Tag: "int"}}, Tag: "error", Value: reflect.ValueOf(os.Lchown)},

			"Link": {Doc: "Link creates newname as a hard link to the oldname file.\nIf there is an error, it will be of type *LinkError.", Args: []pkgreflect.Arg{{Name: "oldname", Tag: "string"}, {Name: "newname", Tag: "string"}}, Tag: "error", Value: reflect.ValueOf(os.Link)},

			"LookupEnv": {Doc: "LookupEnv retrieves the value of the environment variable named\nby the key. If the variable is present in the environment the\nvalue (which may be empty) is returned and the boolean is true.\nOtherwise the returned value will be empty and the boolean will\nbe false.", Args: []pkgreflect.Arg{{Name: "key", Tag: "string"}}, Tag: "any", Value: reflect.ValueOf(os.LookupEnv)},

			"Lstat": {Doc: "Lstat returns a FileInfo describing the named file.\nIf the file is a symbolic link, the returned FileInfo\ndescribes the symbolic link. Lstat makes no attempt to follow the link.\nIf there is an error, it will be of type *PathError.\n\nOn Windows, if the file is a reparse point that is a surrogate for another\nnamed entity (such as a symbolic link or mounted folder), the returned\nFileInfo describes the reparse point, and makes no attempt to resolve it.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}}, Tag: "any", Value: reflect.ValueOf(os.Lstat)},

			"Mkdir": {Doc: "Mkdir creates a new directory with the specified name and permission\nbits (before umask).\nIf there is an error, it will be of type *PathError.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}, {Name: "perm", Tag: "FileMode"}}, Tag: "error", Value: reflect.ValueOf(os.Mkdir)},

			"MkdirAll": {Doc: "MkdirAll creates a directory named path,\nalong with any necessary parents, and returns nil,\nor else returns an error.\nThe permission bits perm (before umask) are used for all\ndirectories that MkdirAll creates.\nIf path is already a directory, MkdirAll does nothing\nand returns nil.", Args: []pkgreflect.Arg{{Name: "path", Tag: "string"}, {Name: "perm", Tag: "FileMode"}}, Tag: "error", Value: reflect.ValueOf(os.MkdirAll)},

			"MkdirTemp": {Doc: "MkdirTemp creates a new temporary directory in the directory dir\nand returns the pathname of the new directory.\nThe new directory's name is generated by adding a random string to the end of pattern.\nIf pattern includes a \"*\", the random string replaces the last \"*\" instead.\nIf dir is the empty string, MkdirTemp uses the default directory for temporary files, as returned by TempDir.\nMultiple programs or goroutines calling MkdirTemp simultaneously will not choose the same directory.\nIt is the caller's responsibility to remove the directory when it is no longer needed.", Args: []pkgreflect.Arg{{Name: "dir", Tag: "string"}, {Name: "pattern", Tag: "string"}}, Tag: "any", Value: reflect.ValueOf(os.MkdirTemp)},

			"NewFile": {Doc: "NewFile returns a new File with the given file descriptor and\nname. The returned value will be nil if fd is not a valid file\ndescriptor. On Unix systems, if the file descriptor is in\nnon-blocking mode, NewFile will attempt to return a pollable File\n(one for which the SetDeadline methods work).\n\nAfter passing it to NewFile, fd may become invalid under the same\nconditions described in the comments of the Fd method, and the same\nconstraints apply.", Args: []pkgreflect.Arg{{Name: "fd", Tag: "uintptr"}, {Name: "name", Tag: "string"}}, Tag: "File", Value: reflect.ValueOf(os.NewFile)},

			"NewSyscallError": {Doc: "NewSyscallError returns, as an error, a new SyscallError\nwith the given system call name and error details.\nAs a convenience, if err is nil, NewSyscallError returns nil.", Args: []pkgreflect.Arg{{Name: "syscall", Tag: "string"}, {Name: "err", Tag: "error"}}, Tag: "error", Value: reflect.ValueOf(os.NewSyscallError)},

			"Open": {Doc: "Open opens the named file for reading. If successful, methods on\nthe returned file can be used for reading; the associated file\ndescriptor has mode O_RDONLY.\nIf there is an error, it will be of type *PathError.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}}, Tag: "any", Value: reflect.ValueOf(os.Open)},

			"OpenFile": {Doc: "OpenFile is the generalized open call; most users will use Open\nor Create instead. It opens the named file with specified flag\n(O_RDONLY etc.). If the file does not exist, and the O_CREATE flag\nis passed, it is created with mode perm (before umask). If successful,\nmethods on the returned File can be used for I/O.\nIf there is an error, it will be of type *PathError.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}, {Name: "flag", Tag: "int"}, {Name: "perm", Tag: "FileMode"}}, Tag: "any", Value: reflect.ValueOf(os.OpenFile)},

			"Pipe": {Doc: "Pipe returns a connected pair of Files; reads from r return bytes written to w.\nIt returns the files and an error, if any.", Args: []pkgreflect.Arg{}, Tag: "any", Value: reflect.ValueOf(os.Pipe)},

			"ReadDir": {Doc: "ReadDir reads the named directory,\nreturning all its directory entries sorted by filename.\nIf an error occurs reading the directory,\nReadDir returns the entries it was able to read before the error,\nalong with the error.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}}, Tag: "any", Value: reflect.ValueOf(os.ReadDir)},

			"ReadFile": {Doc: "ReadFile reads the named file and returns the contents.\nA successful call returns err == nil, not err == EOF.\nBecause ReadFile reads the whole file, it does not treat an EOF from Read\nas an error to be reported.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}}, Tag: "any", Value: reflect.ValueOf(os.ReadFile)},

			"Readlink": {Doc: "Readlink returns the destination of the named symbolic link.\nIf there is an error, it will be of type *PathError.\n\nIf the link destination is relative, Readlink returns the relative path\nwithout resolving it to an absolute one.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}}, Tag: "any", Value: reflect.ValueOf(os.Readlink)},

			"Remove": {Doc: "Remove removes the named file or (empty) directory.\nIf there is an error, it will be of type *PathError.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}}, Tag: "error", Value: reflect.ValueOf(os.Remove)},

			"RemoveAll": {Doc: "RemoveAll removes path and any children it contains.\nIt removes everything it can but returns the first error\nit encounters. If the path does not exist, RemoveAll\nreturns nil (no error).\nIf there is an error, it will be of type *PathError.", Args: []pkgreflect.Arg{{Name: "path", Tag: "string"}}, Tag: "error", Value: reflect.ValueOf(os.RemoveAll)},

			"Rename": {Doc: "Rename renames (moves) oldpath to newpath.\nIf newpath already exists and is not a directory, Rename replaces it.\nOS-specific restrictions may apply when oldpath and newpath are in different directories.\nEven within the same directory, on non-Unix platforms Rename is not an atomic operation.\nIf there is an error, it will be of type *LinkError.", Args: []pkgreflect.Arg{{Name: "oldpath", Tag: "string"}, {Name: "newpath", Tag: "string"}}, Tag: "error", Value: reflect.ValueOf(os.Rename)},

			"SameFile": {Doc: "SameFile reports whether fi1 and fi2 describe the same file.\nFor example, on Unix this means that the device and inode fields\nof the two underlying structures are identical; on other systems\nthe decision may be based on the path names.\nSameFile only applies to results returned by this package's Stat.\nIt returns false in other cases.", Args: []pkgreflect.Arg{{Name: "fi1", Tag: "FileInfo"}, {Name: "fi2", Tag: "FileInfo"}}, Tag: "bool", Value: reflect.ValueOf(os.SameFile)},

			"Setenv": {Doc: "Setenv sets the value of the environment variable named by the key.\nIt returns an error, if any.", Args: []pkgreflect.Arg{{Name: "key", Tag: "string"}, {Name: "value", Tag: "string"}}, Tag: "error", Value: reflect.ValueOf(os.Setenv)},

			"StartProcess": {Doc: "StartProcess starts a new process with the program, arguments and attributes\nspecified by name, argv and attr. The argv slice will become os.Args in the\nnew process, so it normally starts with the program name.\n\nIf the calling goroutine has locked the operating system thread\nwith runtime.LockOSThread and modified any inheritable OS-level\nthread state (for example, Linux or Plan 9 name spaces), the new\nprocess will inherit the caller's thread state.\n\nStartProcess is a low-level interface. The os/exec package provides\nhigher-level interfaces.\n\nIf there is an error, it will be of type *PathError.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}, {Name: "argv", Tag: "[]string"}, {Name: "attr", Tag: "ProcAttr"}}, Tag: "any", Value: reflect.ValueOf(os.StartProcess)},

			"Stat": {Doc: "Stat returns a FileInfo describing the named file.\nIf there is an error, it will be of type *PathError.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}}, Tag: "any", Value: reflect.ValueOf(os.Stat)},

			"Symlink": {Doc: "Symlink creates newname as a symbolic link to oldname.\nOn Windows, a symlink to a non-existent oldname creates a file symlink;\nif oldname is later created as a directory the symlink will not work.\nIf there is an error, it will be of type *LinkError.", Args: []pkgreflect.Arg{{Name: "oldname", Tag: "string"}, {Name: "newname", Tag: "string"}}, Tag: "error", Value: reflect.ValueOf(os.Symlink)},

			"TempDir": {Doc: "TempDir returns the default directory to use for temporary files.\n\nOn Unix systems, it returns $TMPDIR if non-empty, else /tmp.\nOn Windows, it uses GetTempPath, returning the first non-empty\nvalue from %TMP%, %TEMP%, %USERPROFILE%, or the Windows directory.\nOn Plan 9, it returns /tmp.\n\nThe directory is neither guaranteed to exist nor have accessible\npermissions.", Args: []pkgreflect.Arg{}, Tag: "string", Value: reflect.ValueOf(os.TempDir)},

			"Truncate": {Doc: "Truncate changes the size of the named file.\nIf the file is a symbolic link, it changes the size of the link's target.\nIf there is an error, it will be of type *PathError.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}, {Name: "size", Tag: "int64"}}, Tag: "error", Value: reflect.ValueOf(os.Truncate)},

			"Unsetenv": {Doc: "Unsetenv unsets a single environment variable.", Args: []pkgreflect.Arg{{Name: "key", Tag: "string"}}, Tag: "error", Value: reflect.ValueOf(os.Unsetenv)},

			"UserCacheDir": {Doc: "UserCacheDir returns the default root directory to use for user-specific\ncached data. Users should create their own application-specific subdirectory\nwithin this one and use that.\n\nOn Unix systems, it returns $XDG_CACHE_HOME as specified by\nhttps://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if\nnon-empty, else $HOME/.cache.\nOn Darwin, it returns $HOME/Library/Caches.\nOn Windows, it returns %LocalAppData%.\nOn Plan 9, it returns $home/lib/cache.\n\nIf the location cannot be determined (for example, $HOME is not defined),\nthen it will return an error.", Args: []pkgreflect.Arg{}, Tag: "any", Value: reflect.ValueOf(os.UserCacheDir)},

			"UserConfigDir": {Doc: "UserConfigDir returns the default root directory to use for user-specific\nconfiguration data. Users should create their own application-specific\nsubdirectory within this one and use that.\n\nOn Unix systems, it returns $XDG_CONFIG_HOME as specified by\nhttps://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if\nnon-empty, else $HOME/.config.\nOn Darwin, it returns $HOME/Library/Application Support.\nOn Windows, it returns %AppData%.\nOn Plan 9, it returns $home/lib.\n\nIf the location cannot be determined (for example, $HOME is not defined),\nthen it will return an error.", Args: []pkgreflect.Arg{}, Tag: "any", Value: reflect.ValueOf(os.UserConfigDir)},

			"UserHomeDir": {Doc: "UserHomeDir returns the current user's home directory.\n\nOn Unix, including macOS, it returns the $HOME environment variable.\nOn Windows, it returns %USERPROFILE%.\nOn Plan 9, it returns the $home environment variable.\n\nIf the expected variable is not set in the environment, UserHomeDir\nreturns either a platform-specific default value or a non-nil error.", Args: []pkgreflect.Arg{}, Tag: "any", Value: reflect.ValueOf(os.UserHomeDir)},

			"WriteFile": {Doc: "WriteFile writes data to the named file, creating it if necessary.\nIf the file does not exist, WriteFile creates it with permissions perm (before umask);\notherwise WriteFile truncates it before writing, without changing permissions.\nSince WriteFile requires multiple system calls to complete, a failure mid-operation\ncan leave the file in a partially written state.", Args: []pkgreflect.Arg{{Name: "name", Tag: "string"}, {Name: "data", Tag: "[]byte"}, {Name: "perm", Tag: "FileMode"}}, Tag: "error", Value: reflect.ValueOf(os.WriteFile)},
		},

		Variables: map[string]pkgreflect.Value{
			"Args":                {Doc: "", Value: reflect.ValueOf(&os.Args)},
			"ErrClosed":           {Doc: "", Value: reflect.ValueOf(&os.ErrClosed)},
			"ErrDeadlineExceeded": {Doc: "", Value: reflect.ValueOf(&os.ErrDeadlineExceeded)},
			"ErrExist":            {Doc: "", Value: reflect.ValueOf(&os.ErrExist)},
			"ErrInvalid":          {Doc: "ErrInvalid indicates an invalid argument.\nMethods on File will return this error when the receiver is nil.", Value: reflect.ValueOf(&os.ErrInvalid)},
			"ErrNoDeadline":       {Doc: "", Value: reflect.ValueOf(&os.ErrNoDeadline)},
			"ErrNotExist":         {Doc: "", Value: reflect.ValueOf(&os.ErrNotExist)},
			"ErrPermission":       {Doc: "", Value: reflect.ValueOf(&os.ErrPermission)},
			"ErrProcessDone":      {Doc: "", Value: reflect.ValueOf(&os.ErrProcessDone)},
			"Interrupt":           {Doc: "", Value: reflect.ValueOf(&os.Interrupt)},
			"Kill":                {Doc: "", Value: reflect.ValueOf(&os.Kill)},
			"Stderr":              {Doc: "", Value: reflect.ValueOf(&os.Stderr)},
			"Stdin":               {Doc: "", Value: reflect.ValueOf(&os.Stdin)},
			"Stdout":              {Doc: "", Value: reflect.ValueOf(&os.Stdout)},
		},

		Consts: map[string]pkgreflect.Value{
			"DevNull":           {Doc: "", Value: reflect.ValueOf(os.DevNull)},
			"ModeAppend":        {Doc: "", Value: reflect.ValueOf(os.ModeAppend)},
			"ModeCharDevice":    {Doc: "", Value: reflect.ValueOf(os.ModeCharDevice)},
			"ModeDevice":        {Doc: "", Value: reflect.ValueOf(os.ModeDevice)},
			"ModeDir":           {Doc: "The single letters are the abbreviations\nused by the String method's formatting.", Value: reflect.ValueOf(os.ModeDir)},
			"ModeExclusive":     {Doc: "", Value: reflect.ValueOf(os.ModeExclusive)},
			"ModeIrregular":     {Doc: "", Value: reflect.ValueOf(os.ModeIrregular)},
			"ModeNamedPipe":     {Doc: "", Value: reflect.ValueOf(os.ModeNamedPipe)},
			"ModePerm":          {Doc: "", Value: reflect.ValueOf(os.ModePerm)},
			"ModeSetgid":        {Doc: "", Value: reflect.ValueOf(os.ModeSetgid)},
			"ModeSetuid":        {Doc: "", Value: reflect.ValueOf(os.ModeSetuid)},
			"ModeSocket":        {Doc: "", Value: reflect.ValueOf(os.ModeSocket)},
			"ModeSticky":        {Doc: "", Value: reflect.ValueOf(os.ModeSticky)},
			"ModeSymlink":       {Doc: "", Value: reflect.ValueOf(os.ModeSymlink)},
			"ModeTemporary":     {Doc: "", Value: reflect.ValueOf(os.ModeTemporary)},
			"ModeType":          {Doc: "Mask for the type bits. For regular files, none will be set.", Value: reflect.ValueOf(os.ModeType)},
			"O_APPEND":          {Doc: "The remaining values may be or'ed in to control behavior.", Value: reflect.ValueOf(os.O_APPEND)},
			"O_CREATE":          {Doc: "", Value: reflect.ValueOf(os.O_CREATE)},
			"O_EXCL":            {Doc: "", Value: reflect.ValueOf(os.O_EXCL)},
			"O_RDONLY":          {Doc: "Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.", Value: reflect.ValueOf(os.O_RDONLY)},
			"O_RDWR":            {Doc: "", Value: reflect.ValueOf(os.O_RDWR)},
			"O_SYNC":            {Doc: "", Value: reflect.ValueOf(os.O_SYNC)},
			"O_TRUNC":           {Doc: "", Value: reflect.ValueOf(os.O_TRUNC)},
			"O_WRONLY":          {Doc: "", Value: reflect.ValueOf(os.O_WRONLY)},
			"PathListSeparator": {Doc: "", Value: reflect.ValueOf(os.PathListSeparator)},
			"PathSeparator":     {Doc: "", Value: reflect.ValueOf(os.PathSeparator)},
			"SEEK_CUR":          {Doc: "", Value: reflect.ValueOf(os.SEEK_CUR)},
			"SEEK_END":          {Doc: "", Value: reflect.ValueOf(os.SEEK_END)},
			"SEEK_SET":          {Doc: "", Value: reflect.ValueOf(os.SEEK_SET)},
		},
	})
}
