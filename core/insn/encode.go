package insn

import (
	"fmt"
	"math"
)

const (
	MaxOpCode = math.MaxInt8 - 1
	MaxA      = (1 << 24) - 1

	OpShift = 24
	OpMask  = 0xff
	AShift  = 0
	AMask   = MaxA
)

func MakeA(op uint8, a uint) (uint32, error) {
	if op > MaxOpCode {
		return 0, fmt.Errorf("unable to encoding instruction, opcode too large")
	}

	if a > MaxA {
		return 0, fmt.Errorf("unable to encoding instruction, opcode too large")
	}

	i := (uint32(op&OpMask) << OpShift) | (uint32(a) << AShift)

	return i, nil
}

func Decode(i uint32) (uint8, uint) {
	op := uint8(i >> OpShift & OpMask)
	a := uint(i>>AShift) & AMask

	return op, a
}
