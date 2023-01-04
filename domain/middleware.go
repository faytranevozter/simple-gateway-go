package domain

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// default param
type DefaultPayload struct {
	Query        url.Values
	Params       gin.Params
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
