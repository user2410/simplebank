package util

import (
	"database/sql"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomStr Generate a random string of length n
func RandomStr(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomStr(6)
}

func RandomMoney() int64 {
	return rand.Int63n(1000)
}

func RandomCurrency() string {
	var currencies = []string{
		"USD", "EUR", "JPY",
		"VND", "CNY", "RUB",
	}
	return currencies[rand.Intn(len(currencies))]
}

func RandomCountryCode() sql.NullInt32 {
	return sql.NullInt32{rand.Int31n(999), true}
}
