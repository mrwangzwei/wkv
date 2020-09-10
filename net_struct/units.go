package net_struct

import (
	"fmt"
	"math"
)

//单位bit(位)
type Bits bool
type Bytes uint

const ByteBit uint = 8
const BitBase uint = 2

func ByteLength(bits []Bits) Bytes {
	length := len(bits)
	aa := float64(length) / float64(ByteBit)
	var byteLen Bytes
	byteLen = Bytes(math.Ceil(aa))
	return byteLen
}

func BitsVal(bits []Bits) uint {
	bitsBuf := UintBuf(bits)
	blen := len(bits)
	var length float64 = 0
	for loc, val := range bitsBuf {
		length += float64(val) * math.Pow(float64(BitBase), float64(blen-loc-1))
	}
	return uint(length)
}

func UintBuf(bits []Bits) []uint {
	blen := len(bits)
	byteBuf := make([]uint, blen, blen)
	for index, i := range bits {
		byteBuf[index] = BoolToUint(i)
	}
	return byteBuf
}

func BoolToUint(bit Bits) uint {
	switch bit {
	case true:
		return 1
	default:
		return 0
	}
}

func IntToBool(val uint) Bits {
	switch val {
	case 0:
		return false
	default:
		return true
	}
}

func UintToBits(val uint) []Bits {
	var bits []Bits
	for ; val > 0; val /= 2 {
		lsb := val % 2
		bits = append(bits, IntToBool(lsb))
	}
	InvertBits(bits)
	fmt.Println(len(bits))
	return bits
}

func InvertUint(data []uint) []uint {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}

func InvertBits(data []Bits) []Bits {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}
