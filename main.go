package tonic

import (
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
	userHandler := handlers.NewUserHandler(backend)

	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/", userHandler.CreateUser())
			users.PUT(id(), userHandler.UpdateUser())
			users.DELETE(id(), userHandler.DeleteUser())
			users.GET(id(), userHandler.GetUser())
			users.GET("/", userHandler.ListUsers())
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
