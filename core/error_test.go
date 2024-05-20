package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStacktracePath(t *testing.T) {
	path := "/home/evanphx/go/pkg/mod/github.com/lab47/lace@v0.0.0-20240519031230-1a12db59080d/core/fn.go"

	mod := extractMod(path)

	r := require.New(t)

	r.Equal("<github.com/lab47/lace>/core/fn.go", mod)
}
