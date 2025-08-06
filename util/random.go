package util

import (
	"math/rand"
)

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	const charset = "asdfghjklqwertyuiopzxcvbnmZXCVBNMPOIUYTREWQASDFGHJKL"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func RandomOwner() string {
	return RandomString(6)
}
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
