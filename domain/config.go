package domain

import (
	"net/http"
)

type Config struct {
	Cors     CorsConfig `json:"cors"`
	Token    Token      `json:"token"`
	Services []Service  `json:"services"`
}

type Token struct {
	Secret   string `json:"secret"`
	Issuer   string `json:"issuer"`
	Audience string `json:"audience"`
	Expires  int    `json:"expires"`
}

type CorsConfig struct {
	AllowOrigins     []string `json:"allow_origins"`
	AllowMethods     []string `json:"allow_methods"`
	AllowHeaders     []string `json:"allow_headers"`
	ExposeHeaders    []string `json:"expose_headers"`
	AllowCredentials bool     `json:"allow_credentials"`
}

type Service struct {
	Name    string  `json:"name"`
	BaseURL string  `json:"base_url"`
	Routes  []Route `json:"routes"`
}

type Route struct {
	Method        string            `json:"method"`
	Name          string            `json:"name"`
	Auth          bool              `json:"auth"`
	Middlewares   []string          `json:"middlewares"`
	Path          []string          `json:"path"`
	TokenField    map[string]string `json:"token_field"`
	ServicePath   string            `json:"service_path"`
	ServiceMethod string            `json:"service_method"`
}

type RouteService struct {
	Name       string
	Auth       bool
	BaseURL    string
	Path       string
	Method     string
	TokenField map[string]string
}

type RequestPayload struct {
	Headers http.Header
	Query   map[string]string
	Data    interface{}
}
