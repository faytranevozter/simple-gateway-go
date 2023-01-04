package http

import (
	"gateway-api/domain"
	"gateway-api/helpers/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) DynamicRoute(c *gin.Context, RouteID string) {
	ctx := c.Request.Context()

	payload := make(map[string]interface{})
	c.ShouldBindJSON(&payload)

	route := h.routesMap[RouteID]

	authData := make(map[string]interface{})
	rawData, ok := c.Get("JWTDATA")
	if route.Auth {
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "token not valid"))
			return
		} else {
			authData = rawData.(map[string]interface{})
		}
	}

	resp := h.usecase.Dynamic(ctx, domain.DefaultPayload{
		Request:      c.Request,
		Query:        c.Request.URL.Query(),
		Params:       c.Params,
		Payload:      payload,
		AuthData:     authData,
		RouteService: route,
	})

	c.JSON(resp.Status, resp.Data)
}
