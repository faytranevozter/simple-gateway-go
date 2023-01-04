package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway-api/domain"
	"gateway-api/helpers"
	"gateway-api/helpers/response"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jeremywohl/flatten"
)

func (u *gatewayUsecase) Dynamic(ctx context.Context, options domain.DefaultPayload) response.Base {
	_, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// queryparam
	newQuery := make(map[string]string)
	for k, v := range options.Query {
		newQuery[k] = v[0]
	}

	// headers
	newHeaders := options.Request.Header

	// admin
	method := options.RouteService.Method
	url := options.RouteService.BaseURL + options.RouteService.Path

	// refresh map data
	options.AuthData = helpers.ConvertMap(options.AuthData)

	availableData := make(map[string]interface{})
	availableDataStr := make(map[string]string)

	// assign param data
	for _, param := range options.Params {
		availableData["param."+param.Key] = param.Value
		availableDataStr["param."+param.Key] = param.Value
	}

	// assign auth data
	for key, value := range options.AuthData {
		availableData["auth."+key] = value
		availableDataStr["auth."+key] = fmt.Sprintf("%v", value) // convert to string
	}

	// replacer (check if url contains {} ex: "/user/{id}")
	if regexp.MustCompile(`\{([^\{\}]*)\}`).MatchString(options.RouteService.Path) {
		url = helpers.StringReplacer(url, availableDataStr)
	}

	body, statusCode, err := helpers.SimpleRequest(method, url, domain.RequestPayload{
		Data:    options.Payload,
		Query:   newQuery,
		Headers: newHeaders,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	res := make(map[string]interface{})
	json.Unmarshal(body, &res)

	// check login success
	if statusCode != http.StatusOK {
		return response.Base{
			Status: statusCode,
			Data:   res,
		}
	}

	// generate token when login
	if options.RouteService.Name == "auth.login" && len(options.RouteService.TokenField) > 0 {
		dataFlat, _ := flatten.Flatten(res, "", flatten.DotStyle)

		tokenData := make(map[string]interface{})

		for flatKey, toKey := range options.RouteService.TokenField {
			tokenData[toKey] = dataFlat[flatKey]
		}

		cfg := u.config.Token
		token := generateToken(cfg.Secret, cfg.Issuer, cfg.Audience, tokenData, cfg.Expires)

		return response.Success(map[string]interface{}{
			"token":   token,
			"message": "success",
		})
	}

	return response.Success(res)
}

func generateToken(secretKey, issuer, audience string, data map[string]interface{}, expiredSeconds int) string {
	timeNow := time.Now()

	sign := jwt.New(jwt.GetSigningMethod("HS256"))

	claims := sign.Claims.(jwt.MapClaims)
	claims["iat"] = timeNow.Unix()
	if expiredSeconds != 0 {
		duration := time.Duration(expiredSeconds) * time.Second
		claims["exp"] = timeNow.Add(duration).Unix()
	}

	if issuer != "" {
		claims["iss"] = issuer
	}

	if audience != "" {
		claims["aud"] = audience
	}

	for k, v := range data {
		claims[k] = v
	}

	jwtToken, err := sign.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Println("GenerateToken Error: ", err.Error())
	}
	return jwtToken
}
