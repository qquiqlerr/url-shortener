package random

import (
	"math/rand"
	"time"
)

func RandomString(stringLenght int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var out string
	var asciidigit int
	for range stringLenght {
		asciidigit = r.Intn(26) + 97
		out += string(rune(asciidigit))
	}
	return out
}
