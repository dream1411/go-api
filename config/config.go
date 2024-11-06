package config

import (
	"os"
)

var JWT_SECRET = []byte(getJWTSecret())

// getJWTSecret retrieves the JWT secret key from environment variables, or uses a default value.
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	print(secret)
	if secret == "" {
		secret = "concept" // ควรเปลี่ยนเป็นค่าอื่นที่ปลอดภัย
	}
	return secret
}
