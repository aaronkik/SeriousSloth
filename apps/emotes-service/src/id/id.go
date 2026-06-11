package id

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// New returns a new id formatted as "<prefix>_<hex>".
func New(prefix string) string {
	bytes := make([]byte, 12)
	rand.Read(bytes)
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(bytes))
}
