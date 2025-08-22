package utils

import (
	"crypto/rand"
	"fmt"
)

func NewUuid() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		panic("no se pudo generar UUID")
	}
	// Formatear según RFC 4122
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // versión 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // variante 10

	uuidStr := fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])

	return uuidStr
}
