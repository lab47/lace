package core

import (
	"path/filepath"
)

func RunFile(env *Env, filename string) error {
	reader, err := NewReaderFromFile(filename)
	if err != nil {
		return err
	}
	f, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	env.SetMainFilename(f)
	return ProcessReader(env, reader, filename, EVAL)
}
