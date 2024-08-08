package helper

import (
	"time"

	"math/rand"
)

// Function to generate a random price within a given range
func GenerateRandomPrice(min, max int) int {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	return randGen.Intn(max-min+1) + min
}
