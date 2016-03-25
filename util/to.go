package util

import (
	"encoding/base64"
	"encoding/binary"
	"strconv"
)

func Itob(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)

	return b
}

func Btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func Stoi(s string) (uint64, error) {
	var i uint64
	i, err := strconv.ParseUint(s, 10, 64)

	return i, err
}

func Itoe(i uint64) string {
	b := Itob(i)
	e := base64.RawStdEncoding.EncodeToString(b)
	return e
}

func Etoi(e string) uint64 {
	b, err := base64.RawStdEncoding.DecodeString(e)
	if err != nil {
		return 0
	}

	return Btoi(b)
}
