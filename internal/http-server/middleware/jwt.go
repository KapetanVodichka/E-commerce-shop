package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Токен не передан", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]

		key := secretKey
		if envKey := os.Getenv("JWT_SECRET"); envKey != "" {
			key = []byte(envKey)
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return key, nil
		})
		if err != nil {
			http.Error(w, "Невалидный токен", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp, _ := claims["exp"].(float64)
			if float64(time.Now().Unix()) > exp {
				http.Error(w, "Срок действия токена истёк", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Невалидный токен", http.StatusUnauthorized)
			return
		}
	})
}
