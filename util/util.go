package util

import (
	"fmt"
	"strconv"
	"strings"
)

func HexToInt(s string) (int64, error) {
	return strconv.ParseInt(strings.Replace(s, "0x", "", -1), 16, 64)
}

func IntToHex(num int64) string {
	return fmt.Sprintf("0x%x", num)
}
