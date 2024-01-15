package randomize

import (
	"math/rand"
	"time"
)

var charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Generate random string with given input length
// Accept int and ...string as parameter will return string
// e.g.
// length : 5
// charset : "01234567890" (optional)
// return 17934
func String(length int, charsets ...string) string {
	if len(charsets) > 0 {
		charset = charsets[0]
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
