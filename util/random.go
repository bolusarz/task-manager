package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var random *rand.Rand

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	return min + random.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[random.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomPassword(n int) (string, error) {
	if n < 4 {
		return "", nil
	}

	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	symbols := "!@#$%^&*()-_=+[]{}<>?/|"
	all := lower + upper + digits + symbols

	password := make([]byte, n)
	var err error

	sets := []string{lower, upper, digits, symbols}

	for i := 0; i < 4; i++ {
		password[i], err = randomChar(sets[i])
		if err != nil {
			return "", err
		}
	}

	for i := 4; i < n; i++ {
		password[i], err = randomChar(all)
		if err != nil {
			return "", err
		}
	}

	shuffle(password)

	return string(password), nil
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}

func randomChar(set string) (byte, error) {
	return set[random.Intn(len(set))], nil
}

func shuffle(data []byte) {
	for i := range data {
		if i == len(data)-1 {
			continue
		}

		j := RandomInt(int64(i+1), int64(len(data))-1)
		data[i], data[j] = data[j], data[i]
	}
}
