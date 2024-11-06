package middleware

import (
	"learn-api/config"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

// JWTAuthMiddleware checks if the request has a valid JWT token
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ดึงค่า Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		// ตัด "Bearer " ออกและเก็บเฉพาะ token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// ตรวจสอบ JWT Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// ตรวจสอบว่า Signing method เป็น HS256 หรือไม่
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("Invalid signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return config.JWT_SECRET, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// ส่งต่อการทำงานไปยัง handler ถ้า token ถูกต้อง
		next.ServeHTTP(w, r)
	})
}
