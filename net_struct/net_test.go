package net_struct

import (
	"fmt"
	"math"
	"testing"
)

func TestToByte(t *testing.T) {
	fmt.Println(math.Pow(float64(2), 0))
	bits := []Bits{true, true, true, true, true, true, true, true,
		true, true, true, true, true, true, true, true,
		true, true, true, true, true, true, true, true,
		true, true, true, true, true, true, true, true,
		true, true, true, true, true, true, true, true,
		true, true, true, true, true, true, true, true,
		true, true, true, true, true, true, true, true,
		true, true, true, true, true, true, true, true,
	}
	fmt.Println("ByteLength", ByteLength(bits))

}
