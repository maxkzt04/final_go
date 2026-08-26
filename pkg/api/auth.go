package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// создание JWT токена
func createToken(password string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	hash := sha256.Sum256([]byte(password))
	hashStr := fmt.Sprintf("%x", hash)
	claims["hash"] = hashStr

	// устанавливаем срок действия токена на 8 часов
	//exp := time.Now().Add(8 * time.Hour)
	//claims["exp"] = exp.Unix()

	tokenString, err := token.SignedString([]byte(password))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		//pass := os.Getenv("TODO_PASSWORD")
		pass := password
		if len(pass) > 0 {
			var jwtToken string // JWT-токен из куки

			cookie, err := r.Cookie("token")
			if err == nil {
				jwtToken = cookie.Value
			}
			var valid bool

			token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
				return []byte(pass), nil
			})

			if err == nil && token.Valid {
				claims := token.Claims.(jwt.MapClaims)

				hash := sha256.Sum256([]byte(pass))
				hashStr := fmt.Sprintf("%x", hash)

				valid = claims["hash"] == hashStr
			}

			if !valid {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
