package router

import (
	"github.com/labstack/echo/v4"
	"github.com/nuninnih/service_marketplace/app/api/controller/job"
	"github.com/nuninnih/service_marketplace/app/api/controller/project"
	"github.com/nuninnih/service_marketplace/app/api/controller/proposal"
	"github.com/nuninnih/service_marketplace/app/api/controller/user"
	"github.com/nuninnih/service_marketplace/app/api/middleware"
)

func RegisterPath(
	e *echo.Echo,
	jwtSecret string,
	ctrlUser *user.Controller,
	ctrlJob *job.Controller,
	ctrlProp *proposal.Controller,
	ctrlProj *project.Controller,

) {
	jwtMiddleware := middleware.JWTMiddleware(jwtSecret)

	allAccess := middleware.ACLMiddleware(map[string]bool{
		"client":     true,
		"freelancer": true,
	})

	clientAccess := middleware.ACLMiddleware(map[string]bool{
		"client": true,
	})

	freelancerAccess := middleware.ACLMiddleware(map[string]bool{
		"freelancer": true,
	})

	// FREE ROUTE
	// dashboard endpoint -- no need login
	e.GET("/freelancers", ctrlUser.GetAllFreelancer)
	e.GET("/jobs", ctrlJob.GetAllJobs)
	e.POST("/webhook", ctrlJob.WebhookHandler)

	userEndpoint := e.Group("/users")
	userEndpoint.POST("/register", ctrlUser.Register)
	userEndpoint.POST("/login", ctrlUser.Login)

	jobEndpoint := e.Group("/jobs", jwtMiddleware)
	jobEndpoint.POST("", ctrlJob.CreateJob, clientAccess)                           //client
	jobEndpoint.GET("/my", ctrlJob.GetMyJob, clientAccess)                          //client
	jobEndpoint.GET("/:id/proposals", ctrlProp.GetJobProposalPerUser, clientAccess) //client
	jobEndpoint.POST("/:id/proposals", ctrlProp.CreateProposal, freelancerAccess)   //freelancer
	jobEndpoint.GET("/:id", ctrlJob.GetJobById, allAccess)                          //freelancer

	proposalEndpoint := e.Group("/proposals", jwtMiddleware)
	proposalEndpoint.PATCH("/:id/accept", ctrlProp.ApproveProposal, clientAccess) //client

	projectEndpoint := e.Group("/projects", jwtMiddleware)
	// projectEndpoint.GET("/my")           // freelancer
	projectEndpoint.PATCH("/:id/submit", ctrlProj.UpdateStatusProject, freelancerAccess) // freelancer
	// projectEndpoint.PACTH("/:id/pay")     //client
}
