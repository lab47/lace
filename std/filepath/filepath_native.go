package filepath

import (
	"os"
	"path/filepath"

	. "github.com/candid82/joker/core"
)

func fileSeq(root string) (*Vector, error) {
	res := EmptyVector()
	var err error
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		m := FileInfoMap(path, info)
		res, err = res.Conjoin(m)
		return err
	})
	return res, err
}
