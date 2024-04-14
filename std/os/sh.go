//go:build !plan9
// +build !plan9

package os

import (
	"bytes"
	"io"
	"os/exec"
	"syscall"

	. "github.com/lab47/lace/core"
)

func sh(dir string, stdin io.Reader, stdout io.Writer, stderr io.Writer, name string, args []string) (Object, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdin = stdin

	var stdoutBuffer, stderrBuffer bytes.Buffer
	if stdout != nil {
		cmd.Stdout = stdout
	} else {
		cmd.Stdout = &stdoutBuffer
	}
	if stderr != nil {
		cmd.Stderr = stderr
	} else {
		cmd.Stderr = &stderrBuffer
	}

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()

	res := EmptyArrayMap()
	res.AddEqu(MakeKeyword("success"), Boolean{B: err == nil})

	var exitCode int
	if err != nil {
		res.AddEqu(MakeKeyword("err-msg"), String{S: err.Error()})
		if exiterr, ok := err.(*exec.ExitError); ok {
			ws := exiterr.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			exitCode = defaultFailedCode
		}
	} else {
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	res.AddEqu(MakeKeyword("exit"), Int{I: exitCode})
	if stdout == nil {
		res.AddEqu(MakeKeyword("out"), String{S: stdoutBuffer.String()})
	}
	if stderr == nil {
		res.AddEqu(MakeKeyword("err"), String{S: stderrBuffer.String()})
	}
	return res, nil
}
