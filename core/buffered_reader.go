package core

import (
	"bufio"
	"io"
)

// A value that can return data that has been buffered.
//
//lace:export
type BufferedReader struct {
	*bufio.Reader
	hash uint32
}

func MakeBufferedReader(rd io.Reader) *BufferedReader {
	res := &BufferedReader{bufio.NewReader(rd), 0}
	res.hash = HashPtr(res)
	return res
}
