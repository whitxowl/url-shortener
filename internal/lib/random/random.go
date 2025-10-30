package random

import (
	"math/rand"
	"sync"
	"time"
)

var (
	rnd  *rand.Rand
	once sync.Once
)

func getRand() *rand.Rand {
	once.Do(func() {
		rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
	return rnd
}

func NewRandomString(size int) string {
	r := getRand()
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, size)
	for i := range b {
		b[i] = chars[r.Intn(len(chars))]
	}
	return string(b)
}
