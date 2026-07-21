package util

import "crypto/rand"

func GenerateRandomState() string {
	return rand.Text()
}

func GenerateRandomSecret() string {
	return rand.Text()
}
