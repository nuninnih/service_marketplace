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

	// dashboard endpoint -- no need login
	// GET /jobs
	// GET /freelancer

	// client endpoint
	// POST /jobs
	// GET /my/jobs
	// GET /jobs/:id/proposals
	// POST /proposals/:id/accept -- PATCH aja kayaknya
	// POST /projects/:id/pay

	// freelancer endpoint
	// GET /jobs
	// GET /jobs/:id
	// POST /jobs/:id/proposals
	// GET /my/projects
	// PATCH /projects/:id/submit

}
