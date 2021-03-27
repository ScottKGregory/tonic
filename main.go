package tonic

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/helpers"
	"github.com/scottkgregory/tonic/pkg/middleware"
	"github.com/scottkgregory/tonic/pkg/models"
	"github.com/scottkgregory/tonic/pkg/services"
)

func Init(opt models.Options) (*gin.Engine, *gin.RouterGroup, error) {
	gin.SetMode(gin.ReleaseMode)

	statusService := services.NewStatusService(&opt.Backend)
	permissionService := services.NewPermissionService(&opt.Permissions)
	authService := services.NewAuthService(&opt.Auth, permissionService, &opt.Backend)
	userService := services.NewUserService(&opt.Backend, permissionService)

	r := gin.New()
	r.Use(middleware.Zerologger(opt.Log))
	log := helpers.GetLogger()

	// Home page
	if !opt.DisableHomepage {
		log.Trace().Msg("Registering default homepage")
		r.GET("/", services.HomePage)
	}

	// Error Pages
	if !opt.DisableErrorPages {
		log.Trace().Msg("Registering default error pages")
		r.GET("/error/:code", services.ErrorPage(0))
		r.NoRoute(services.ErrorPage(404))
	}

	// Probes
	if !opt.DisableHealthProbes {
		log.Trace().Msg("Registering default health probes")
		r.GET("/health", statusService.Health)
		r.GET("/liveliness", statusService.Liveliness)
		r.GET("/readiness", statusService.Readiness)
	}

	// Auth
	var authed *gin.RouterGroup
	if !opt.Auth.Disabled {
		log.Trace().Msg("Registering default auth endpoints")
		r.GET("/login", authService.Login)
		r.GET("/logout", authService.Logout)
		r.GET("/callback", authService.Callback)

		// Guarded
		log.Trace().Msg("Registering guarded router group")
		authed = r.Group("/api")
		authed.Use(middleware.Authed(&opt.Cookie, &opt.JWT, authService))
		{
			authed.GET("/token", middleware.Any(authService.Token, "GetToken"))
			authed.GET("/me", authService.Me)

			authed.GET("/users", middleware.Any(userService.ListUsers, "ListUsers"))
			authed.PUT(fmt.Sprintf("/users/:%s", constants.IDParam), middleware.Any(userService.UpdateUser, "ManageUsers"))
			authed.GET(fmt.Sprintf("/users/:%s", constants.IDParam), middleware.Any(userService.GetUser, "DescribeUser"))
		}
	}

	log.Trace().Msg("Tonic setup complete")
	return r, authed, nil
}
