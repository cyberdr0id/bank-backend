package util

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

const alphabet = "abcdefghikjlkmnopqrstuvwxyz"

func RandString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]

		sb.WriteByte(c)
	}

	return sb.String()
}

func RandOwner() string {
	return RandString(6)
}

func RandMoney() int64 {
	return RandInt(0, 1000)
}

func RandCurrency() string {
	currencies := []string{"USD", "EUR", "ILS", "GBP", "JPY", "CNY", "CRC", "KZT", "GEL", "PLN", "LAK"}
	n := len(currencies)

	return currencies[rand.Intn(n)]
}
