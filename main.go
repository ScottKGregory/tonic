package tonic

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/scottkgregory/tonic/pkg/backends"
	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/dependencies"
	"github.com/scottkgregory/tonic/pkg/handlers"
	"github.com/scottkgregory/tonic/pkg/middleware"
	"github.com/scottkgregory/tonic/pkg/models"
)

// Init sets up tonic
func Init(opt models.Options) (*gin.Engine, *gin.RouterGroup, error) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(middleware.Zerologger(opt.Log))
	log := dependencies.GetLogger()

	backend := backends.NewMongoBackend(&opt.Backend)
	homeHandler := handlers.NewHomeHandler(opt.PageHeader)
	errorHandler := handlers.NewErrorHandler(opt.PageHeader)
	probeHandler := handlers.NewProbeHandler(backend)
	userHandler := handlers.NewUserHandler(backend)
	authHandler := handlers.NewAuthHandler(backend, &opt.Auth)
	permissionHandler := handlers.NewPermissionsHandler()

	if !opt.DisableHomepage {
		router.GET("/", homeHandler.Home())
	}

	if !opt.DisableErrorPages {
		router.GET("/error/:code", errorHandler.Error(0))
		router.NoRoute(errorHandler.Error(http.StatusNotFound))
	}

	if !opt.DisableHealthProbes {
		router.GET("/health", probeHandler.Health())
		router.GET("/liveliness", probeHandler.Liveliness())
		router.GET("/readiness", probeHandler.Readiness())
	}

	auth := router.Group("/auth")
	{
		auth.GET("/login", authHandler.Login())
		auth.GET("/callback", authHandler.Callback())
		auth.GET("/logout", authHandler.Logout())
	}

	api := router.Group("/api")
	api.Use(middleware.Authed(backend, &opt.Cookie, &opt.JWT, &opt.Auth))
	{
		users := api.Group("/users")
		{
			users.POST("/", middleware.HasAny("users:create:*"), userHandler.CreateUser())
			users.PUT(id(), middleware.HasAny(id("users:update:")), userHandler.UpdateUser())
			users.DELETE(id(), middleware.HasAny(id("users:delete:")), userHandler.DeleteUser())
			users.GET(id(), middleware.HasAny(id("users:get:")), userHandler.GetUser())
			users.GET("/", middleware.HasAny("users:list:*"), userHandler.ListUsers())
		}

		auth := api.Group("/auth")
		{
			auth.GET("/token", middleware.HasAny("token:get:*"), authHandler.Token())
		}

		permissions := api.Group("/permissions")
		{
			permissions.GET("/", middleware.HasAny("permissions:list:*"), permissionHandler.ListPermissions())
		}
	}

	log.Trace().Msg("Tonic setup complete")
	return router, api, nil
}

// GetLogger returns a zerolog logger with context
func GetLogger(c ...*gin.Context) *zerolog.Logger {
	return dependencies.GetLogger(c...)
}

func id(path ...string) string {
	p := ""
	if len(path) > 0 {
		p = path[0]
	}

	if strings.HasSuffix(p, ":") {
		return p + constants.IDParam
	}

	return strings.TrimSuffix(p, "/") + "/:" + constants.IDParam
}
