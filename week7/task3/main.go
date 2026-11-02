package main

import (
	"strconv"
	"strings"
)

func JoinIDsSlow(ids []int) string {
	s := ""
	for _, id := range ids {
		s += strconv.Itoa(id) + ","
	}
	return s
}

func JoinIDsFast(ids []int) string {
	var sb strings.Builder
	sb.Grow(len(ids) * 7)

	var tmp [20]byte

	for _, id := range ids {
		// buf := strconv.AppendInt(nil, int64(id), 10)
		// sb.WriteString(strconv.Itoa(id))
		// sb.Write(buf)
		sb.Write(strconv.AppendInt(tmp[:0], int64(id), 10))
		sb.WriteByte(',')
	}
	return sb.String()
}

func main() {}
