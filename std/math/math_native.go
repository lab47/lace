package math

import (
	"math"

	. "github.com/candid82/joker/core"
)

func modf(x float64) (Object, error) {
	i, f := math.Modf(x)
	res := EmptyVector()
	res, err := res.Conjoin(MakeDouble(i))
	if err != nil {
		return nil, err
	}
	return res.Conjoin(MakeDouble(f))
}
