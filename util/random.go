package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {

	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min int64, max int64) int64 {

	return min + rand.Int63n(max-min+1)

}

func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomName() string {
	return RandomString(6)
}

func RandomRole() int64 {

	return RandomInt(1, 3)
}

func RandomRoleString() string {
	return ROLES[RandomRole()]
}

func RandomBudget() decimal.Decimal {
	min := float64(1000)
	max := float64(100000)
	budget := decimal.NewFromFloat(min + rand.Float64()*(max-min)).Round(int32(2))
	//fmt.Println("Subtotal:", budget)
	return budget
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomAddress() string {
	return RandomString(1000)
}

func RandomFloatInInterval(interval int64) float64 {
	return float64(rand.Int63n(2*interval) - interval)
}

func RandomLongitude() float64 {
	return RandomFloatInInterval(180)
}

func RandomLatitude() float64 {
	return RandomFloatInInterval(90)
}
