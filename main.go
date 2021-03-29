package tonic

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/pkg/backends"
	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/dependencies"
	"github.com/scottkgregory/tonic/pkg/handlers"
	"github.com/scottkgregory/tonic/pkg/middleware"
	"github.com/scottkgregory/tonic/pkg/models"
)

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

	router.GET("/", homeHandler.Home())
	router.GET("/error/:code", errorHandler.Error(0))
	router.NoRoute(errorHandler.Error(http.StatusNotFound))

	router.GET("/health", probeHandler.Health())
	router.GET("/liveliness", probeHandler.Liveliness())
	router.GET("/readiness", probeHandler.Readiness())

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
			users.POST("/", middleware.HasAny(backend, "users:create:*"), userHandler.CreateUser())
			users.PUT(id(), middleware.HasAny(backend, "users:update:*"), userHandler.UpdateUser())
			users.DELETE(id(), middleware.HasAny(backend, "users:delete:*"), userHandler.DeleteUser())
			users.GET(id(), middleware.HasAny(backend, "users:get:*"), userHandler.GetUser())
			users.GET("/", middleware.HasAny(backend, "users:list:*"), userHandler.ListUsers())
		}

		auth := api.Group("/auth")
		{
			auth.GET("/token", middleware.HasAny(backend, "token:get:*"), authHandler.Token())
		}

		permissions := api.Group("/permissions")
		{
			permissions.GET("/", permissionHandler.ListPermissions())
		}
	}

	log.Trace().Msg("Tonic setup complete")
	return router, nil, nil
}

func id(path ...string) string {
	p := ""
	if len(path) > 0 {
		p = path[0]
	}

	return strings.TrimSuffix(p, "/") + "/:" + constants.IDParam
}
