package router

import (
	"github.com/labstack/echo/v4"
	"github.com/nuninnih/service_marketplace/app/api/controller/job"
	"github.com/nuninnih/service_marketplace/app/api/controller/user"
	"github.com/nuninnih/service_marketplace/app/api/middleware"
)

func RegisterPath(
	e *echo.Echo,
	jwtSecret string,
	ctrlUser *user.Controller,
	ctrlJob *job.Controller,

) {
	jwtMiddleware := middleware.JWTMiddleware(jwtSecret)

	clientAccess := middleware.ACLMiddleware(map[string]bool{
		"client": true,
	})

	// freelancerAccess := middleware.ACLMiddleware(map[string]bool{
	// 	"freelancer": true,
	// })

	// FREE ROUTE
	e.GET("/freelancers", ctrlUser.GetAllFreelancer)
	e.GET("/jobs", ctrlJob.GetAllJobs)

	userEndpoint := e.Group("/users")
	userEndpoint.POST("/register", ctrlUser.Register)
	userEndpoint.POST("/login", ctrlUser.Login)

	// dashboard endpoint -- no need login
	jobEndpoint := e.Group("/jobs", jwtMiddleware)
	jobEndpoint.POST("", ctrlJob.CreateJob, clientAccess)

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
