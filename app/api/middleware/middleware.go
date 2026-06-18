package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type JwtClaims struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func unAuthorizeResponse(c echo.Context) error {
	return c.JSON(http.StatusUnauthorized, map[string]string{
		"message": "Invalid or Missing Token",
	})
}

func JWTMiddleware(jwtSign string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtSign),

		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtClaims)
		},

		ErrorHandler: func(c echo.Context, err error) error {
			return unAuthorizeResponse(c)
		},

		SuccessHandler: func(c echo.Context) {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(*JwtClaims)
			c.Set("id", claims.ID)
			c.Set("role", claims.Role)
		},
	})
}

func forbiddenResponse(c echo.Context) error {
	return c.JSON(http.StatusForbidden, map[string]interface{}{"message": http.StatusText(http.StatusForbidden)})
}

func ACLMiddleware(rolesMap map[string]bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, _ := c.Get("role").(string)

			if rolesMap[role] {
				return next(c)
			}

			return forbiddenResponse(c)
		}
	}
}
