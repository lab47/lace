package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnv(t *testing.T) {
	t.Run("can run code directly", func(t *testing.T) {
		r := require.New(t)

		e, err := NewEnv()
		r.NoError(err)

		obj, err := e.Eval("(+ 3 4)")
		r.NoError(err)

		r.True(obj.Equals(e, MakeInt(7)))
	})
}
