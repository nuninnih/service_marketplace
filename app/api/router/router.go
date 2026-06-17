package router

import (
	"github.com/labstack/echo/v4"
	"github.com/nuninnih/service_marketplace/app/api/controller/user"
)

func RegisterPath(
	e *echo.Echo,
	jwtSecret string,
	ctrlUser *user.Controller,

) {
	// jwtMiddleware := middleware.JWTMiddleware(jwtSecret)

	// clientAccess := middleware.ACLMiddleware(map[string]bool{
	// 	"CLIENT": true,
	// })
	// freelancerAccess := middleware.ACLMiddleware(map[string]bool{
	// 	"FREELANCER": true,
	// })

	userEndpoint := e.Group("/users")
	userEndpoint.POST("/register", ctrlUser.Register)
	userEndpoint.POST("/login", ctrlUser.Login)
}
