package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func HexToInt(hex string) int64 {
	parsed, _ := strconv.ParseInt(strings.Replace(hex, "0x", "", -1), 16, 64)
	return parsed
}

func IntToHex(i int64) string {
	return fmt.Sprintf("0x%x", i)
}
