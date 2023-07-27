package jwt

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Claim struct {
	User_id uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

func ParseToken(authToken string) (*jwt.Token, error) {

	tokenArr := strings.Split(authToken, " ")

	if len(tokenArr) <= 1 {
		return nil, nil
	}
	claim := &Claim{}
	tkn, err := jwt.ParseWithClaims(tokenArr[1], claim, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	return tkn, err
}
