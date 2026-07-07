package random

import "crypto/rand"

func GenerateRandomState() string {
	return rand.Text()
}
