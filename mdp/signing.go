package mdp

import (
	"crypto/md5"
	"fmt"
	"math/rand"
)

// Letters contains a list of runes that are considered valid for the nonce
// parameter.
// TODO(arandall): verify support for additional values.
var Letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

// Hex contains a list of runes that are considered valid for the messageID
// parameter.
// TODO(arandall): verify support for additional values.
var HEX = []rune("abcdef0123456789")

// RandSeq can be used to create a random nonce.
func RandSeq(format []rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = format[rand.Intn(len(format))]
	}
	return string(b)
}

// GenerateSignature creates an md5 signing string from the concatenation of the relevant strings.
func GenerateSignature(strings ...string) string {
	d := []byte{}
	for _, s := range strings {
		d = append(d, []byte(s)...)
	}
	return fmt.Sprintf("%x", md5.Sum(d))
}
