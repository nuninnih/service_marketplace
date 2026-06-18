package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	jobCTRL "github.com/nuninnih/service_marketplace/app/api/controller/job"
	proposalCRTL "github.com/nuninnih/service_marketplace/app/api/controller/proposal"
	userCTRL "github.com/nuninnih/service_marketplace/app/api/controller/user"
	"github.com/nuninnih/service_marketplace/app/api/router"
	jobRepo "github.com/nuninnih/service_marketplace/repository/job"
	proposalRepo "github.com/nuninnih/service_marketplace/repository/proposal"
	userRepo "github.com/nuninnih/service_marketplace/repository/user"
	jobSvc "github.com/nuninnih/service_marketplace/service/job"
	proposalSvc "github.com/nuninnih/service_marketplace/service/proposal"
	userSvc "github.com/nuninnih/service_marketplace/service/user"
	"github.com/nuninnih/service_marketplace/util/db"
)

func main() {
	err := os.MkdirAll("logs", os.ModePerm)
	logFile, err := os.OpenFile(
		"logs/app.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)

	if err != nil {
		fmt.Println(err)
	}

	writer := io.MultiWriter(os.Stdout, logFile)
	var loggerOption = slog.HandlerOptions{AddSource: true}
	var logger = slog.New(slog.NewJSONHandler(writer, &loggerOption))

	spew.Dump()
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.CORS())

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("incoming request",
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.Duration("latency", v.Latency),
				slog.String("method", c.Request().Method),
			)
			return nil
		},
	}))

	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"hello": "world"})
	})

	db := db.GetPostgresConnection()

	// TODO Setup Public API

	userRepository := userRepo.NewGormRepository(db)
	userService := userSvc.NewService(logger, userRepository)
	userController := userCTRL.NewController(logger, userService)

	jobRepository := jobRepo.NewGormRepository(db)
	jobService := jobSvc.NewService(logger, jobRepository)
	jobController := jobCTRL.NewController(logger, jobService)

	proposalRepository := proposalRepo.NewGormRepository(db)
	proposalService := proposalSvc.NewService(logger, proposalRepository, jobRepository)
	proposalController := proposalCRTL.NewController(logger, proposalService)

	router.RegisterPath(
		e,
		os.Getenv("JWT_SECRET"),
		userController,
		jobController,
		proposalController,
	)

	logger.Info("http server started on :" + os.Getenv("PORT"))
	err = e.Start(":" + os.Getenv("PORT"))
	if err != nil {
		logger.Error("http server error on started", slog.Any("err", err))
	}
}
