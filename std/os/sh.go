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
	res.Add(MakeKeyword("success"), Boolean{B: err == nil})

	var exitCode int
	if err != nil {
		res.Add(MakeKeyword("err-msg"), String{S: err.Error()})
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
	res.Add(MakeKeyword("exit"), Int{I: exitCode})
	if stdout == nil {
		res.Add(MakeKeyword("out"), String{S: string(stdoutBuffer.Bytes())})
	}
	if stderr == nil {
		res.Add(MakeKeyword("err"), String{S: string(stderrBuffer.Bytes())})
	}
	return res, nil
}
