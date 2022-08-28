package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var (
	jwtSignedMethod = jwt.SigningMethodHS256
)

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			SECRET_KEY := os.Getenv("SECRET_JWT")
			signature := strings.Split(c.Request().Header.Get("Authorization"), " ")
			if len(signature) < 2 {
				return c.JSON(http.StatusForbidden, "Invalid token")
			}

			if signature[0] != "Bearer" {
				return c.JSON(http.StatusForbidden, "Invalid token")
			}

			claim := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(signature[1], claim, func(t *jwt.Token) (interface{}, error) {
				return []byte(SECRET_KEY), nil
			})
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"messages": "wrong signed",
				})
			}
			if claim["Role"] != "customer" && claim["Role"] != "ppn" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"messages": "role not authorized",
				})
			}

			method, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok || method != jwtSignedMethod {
				return c.JSON(http.StatusForbidden, "Invalid token")
			}

			return next(c)
		}
	}
}
