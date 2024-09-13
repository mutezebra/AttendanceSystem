package utils

import (
	"math/rand"
	"strings"
)

var sb = strings.Builder{}

func GenerateCode(length int) string {
	defer sb.Reset()
	codes := [10]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	for i := 0; i < length; i++ {
		sb.WriteByte(codes[rand.Intn(len(codes))])
	}
	return sb.String()
}
