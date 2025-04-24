package driver

import (
	"crypto/rand"
	"fmt"
)

func GenerateIdentity(length *int64) string {
	var randBytes []byte

	randBytes = make([]byte, *length/2)
	_, _ = rand.Read(randBytes)
	
	return fmt.Sprintf("%x", randBytes)
}
