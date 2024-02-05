package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims estructura para los claims del JWT
type Claims struct {
	UserID uint
	jwt.StandardClaims
}

type Key string

const UserContextKey Key = "userID"

// GenerateToken genera un nuevo token JWT para un usuario
func GenerateToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(8 * time.Hour) // 8 horas de expiración
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken valida el token JWT y devuelve los claims
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.NewValidationError("invalid token", jwt.ValidationErrorMalformed)
	}

	return claims, nil
}

func ExtractUserIDFromRequest(r *http.Request) (uint, error) {
	userID, ok := r.Context().Value(UserContextKey).(uint) // Ajusta el tipo según sea necesario.
	if !ok {
		return 0, errors.New("could not extract user ID from context")
	}
	return userID, nil
}
