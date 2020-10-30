package auth

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
)

// TokenParsingError ...
type TokenParsingError struct {
	message string
}

func (err TokenParsingError) Error() string {
	return err.message
}

type TokenDetails struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
}

type TokenClaims struct {
	User TokenDetails `json:"user"`

	jwt.StandardClaims
}

func generateToken(userID uuid.UUID) *jwt.Token {

	issuedAt := time.Now().Unix()
	expirationTime := time.Now().Add(6 * time.Hour).Unix()

	claims := TokenClaims{
		User: TokenDetails{
			UserId: userID.String(),
		},

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			IssuedAt:  issuedAt,
		},
	}

	// We build a token, we give it and expiry of 6 hours.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token
}

// GetTokenString generates a jwt access token for a user
func GetTokenString(userID uuid.UUID, secret string) (string, error) {
	token := generateToken(userID)

	str, err := token.SignedString([]byte(secret))
	if err != nil { // we have an error generating the token i.e. "500"
		log.Println(err)
		return "", TokenParsingError{message: err.Error()}
	}
	return str, nil
}

func ParseToken(token, secret string, claims *TokenClaims) (*jwt.Token, error) {
	tok, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	return tok, err
}

func ValidateToken(tok *jwt.Token) bool {
	if !tok.Valid {
		return false
	}
	return true
}
