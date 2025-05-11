package utils

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func GenerateRandomInRange(min, max int) int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	return r.Intn(max-min) + min
}

func GetDateTime(s string) time.Time {
	layout := time.DateTime
	t, _ := time.Parse(layout, s)
	return t
}
