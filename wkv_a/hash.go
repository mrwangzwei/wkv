package wkv_a

import (
	"errors"
	"reflect"
)

const (
	// offset64 FNVa offset basis. See https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function#FNV-1a_hash
	offset64 = 14695981039346656037
	// prime64 FNVa prime value. See https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function#FNV-1a_hash
	prime64 = 1099511628211
)

func sum64(key string) uint64 {

	var hash uint64 = offset64
	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= prime64
	}
	return hash

}

func tranUint(key string) (uint64, error) {

	var index uint64
	switch reflect.TypeOf(key).Kind() {
	case reflect.String:
		index = sum64(key)
		break
	default:
		return 0, errors.New("illegal key type")
	}
	return index, nil

}
