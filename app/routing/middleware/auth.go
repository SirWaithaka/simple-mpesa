package middleware

import (
	"strings"

	"simple-wallet/app/auth"
	"simple-wallet/app/errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func AuthByBearerToken(secret string) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		// check that the header is actually set
		header := ctx.Get("Authorization")
		if header == "" {
			return errors.Unauthorized{Message: "authorization header not set"}
		}

		// check that the token value in header is set
		bearer := strings.Split(header, " ")
		if len(bearer) < 2 || bearer[1] == "" {
			return errors.Unauthorized{Message: "authentication token not set"}
		}

		var claims auth.TokenClaims
		token, err := auth.ParseToken(bearer[1], secret, &claims)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return errors.Unauthorized{Message: "invalid signature on token"}
			}

			return errors.Unauthorized{Message: "token has expired or is invalid"}
		}
		if valid := auth.ValidateToken(token); !valid {
			return errors.Unauthorized{Message: "invalid token"}
		}

		userDetails := map[string]string{
			"userId": claims.User.UserId,
			"email":  claims.User.Email,
		}
		ctx.Locals("userDetails", userDetails)

		return ctx.Next()
	}
}
