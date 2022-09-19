package util

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var Currencies = [...]string{
	"USD", "EUR", "JPY",
	"VND", "CNY", "RUB",
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
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
	return Currencies[rand.Intn(len(Currencies))]
}

func RandomCountryCode() sql.NullInt32 {
	return sql.NullInt32{rand.Int31n(999), true}
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomStr(8))
}
