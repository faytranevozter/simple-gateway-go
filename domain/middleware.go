package domain

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// default param
type DefaultPayload struct {
	Context      *gin.Context
	Request      *http.Request
	ID           interface{}
	Payload      interface{}
	RouteService RouteService
	AuthData     map[string]interface{}
}

type JWTPayload struct {
	jwt.StandardClaims
}

type Map map[string]interface{}
