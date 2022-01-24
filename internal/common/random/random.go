package random

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Int(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func Ints(min, max int) string {
	return strconv.Itoa(min + rand.Intn(max-min+1))
}

func Codes(n int) string {
	var codes strings.Builder
	for i := 1; i <= n; i++ {
		c := Ints(0, 9)
		codes.WriteString(c)
	}
	return codes.String()
}
