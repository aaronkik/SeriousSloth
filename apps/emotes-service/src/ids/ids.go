package ids

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
)

// New returns a new id formatted as "<prefix><hex>" where hex is 12 random bytes.
func New(prefix string) string {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s%s", prefix, hex.EncodeToString(bytes))
}
