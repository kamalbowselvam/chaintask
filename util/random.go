package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
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
	return RandomInt(0, 2)
}

func RandomRoleString() string {
	return ROLES[RandomRole()]
}

func RandomBudget() float64 {
	return float64(RandomInt(0, 100000))
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}