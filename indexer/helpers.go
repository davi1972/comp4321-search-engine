package Indexer

import (
	"encoding/binary"
)

func uint64ToByte(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func byteToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}