package filepath

import (
	"os"
	"path/filepath"

	. "github.com/lab47/lace/core"
)

func fileSeq(env *Env, root string) (*Vector, error) {
	res := EmptyVector()
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		m := FileInfoMap(env, path, info)
		res, err = res.Conjoin(m)
		return err
	})
	return res, err
}
