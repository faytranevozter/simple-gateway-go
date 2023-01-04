package middleware

import (
	"gateway-api/helpers/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Middleware represent the data-struct for middleware
type Middleware struct {
	// another stuff, may be needed by middleware
	jwtSecretUser string
}

// AuthUser for jwt authentication
func (m *Middleware) AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get header
		hAuth := c.GetHeader("Authorization")

		splitToken := strings.Split(hAuth, "Bearer ")
		if len(splitToken) != 2 {
			c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Unauthorized"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := splitToken[1]

		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.jwtSecretUser), nil
		})

		if token == nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, err.Error()))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "invalid token data"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var claimMap map[string]interface{} = claims

		c.Set("JWTDATA", claimMap)

		c.Next()
	}
}

// InitMiddleware initialize the middleware
func InitMiddleware(secret string) *Middleware {
	return &Middleware{
		jwtSecretUser: secret,
	}
}
