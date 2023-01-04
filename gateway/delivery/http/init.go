package http

import (
	"gateway-api/domain"
	"gateway-api/gateway/delivery/http/middleware"
	"gateway-api/gateway/usecase"
	"gateway-api/helpers"
	"gateway-api/helpers/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type handler struct {
	usecase   usecase.GatewayUsecase
	config    domain.Config
	routesMap map[string]domain.RouteService
}

func NewHandler(r *gin.Engine, cu usecase.GatewayUsecase) {
	config := cu.GetConfig()
	handler := &handler{
		usecase:   cu,
		config:    config,
		routesMap: helpers.ConfigToMap(config),
	}

	mdl := middleware.InitMiddleware(config.Token.Secret)

	for _, s := range handler.config.Services {
		for _, configRoute := range s.Routes {

			for _, path := range configRoute.Path {
				dynamicRouteHandler := make([]gin.HandlerFunc, 0)

				if configRoute.Auth && !utils.Strings(configRoute.Middlewares).Include("auth.jwt") {
					configRoute.Middlewares = append(configRoute.Middlewares, "auth.jwt")
				}

				for _, routeMiddleware := range configRoute.Middlewares {
					if routeMiddleware == "auth.jwt" {
						dynamicRouteHandler = append(dynamicRouteHandler, mdl.AuthUser())
					}
				}

				dynamicRouteHandler = append(dynamicRouteHandler, func(ctx *gin.Context) {
					handler.DynamicRoute(ctx, path)
				})

				// handling per route
				if strings.ToUpper(configRoute.Method) == "GET" {
					r.GET(path, dynamicRouteHandler...)
				} else if strings.ToUpper(configRoute.Method) == "POST" {
					r.POST(path, dynamicRouteHandler...)
				}
			}
		}
	}

}
