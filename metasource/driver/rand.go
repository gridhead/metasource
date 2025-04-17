package driver

import (
	"crypto/rand"
	"fmt"
)

func randomHex(n int) string {
	b := make([]byte, n/2)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
