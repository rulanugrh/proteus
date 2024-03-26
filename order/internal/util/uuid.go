package util

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateUUID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 6)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	uuid := fmt.Sprintf("ORDER-%s-%d", string(b), time.Now().UnixNano())

	return uuid
}