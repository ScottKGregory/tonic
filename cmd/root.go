package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/rs/zerolog"
	_ "github.com/rs/zerolog/log"
	"github.com/scottkgregory/mamba"
	"github.com/scottkgregory/tonic/pkg/backends"
	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/dependencies"
	"github.com/scottkgregory/tonic/pkg/handlers"
	"github.com/scottkgregory/tonic/pkg/middleware"
	"github.com/scottkgregory/tonic/pkg/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg models.Config

var rootCmd = &cobra.Command{
	Use:   "webapi",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg.Log.IgnoreRoutes = []string{"/health", "/liveliness", "/readiness"}
		gin.SetMode(gin.ReleaseMode)

		router := gin.New()
		router.Use(middleware.Zerologger(cfg.Log))
		log := dependencies.GetLogger()

		var backend backends.Backend
		var err error
		if cfg.Backend.InMemory {
			backend = backends.NewMemoryBackend(&cfg.Backend)
		} else {
			backend, err = backends.NewMongoBackend(context.Background(), &cfg.Backend)
		}
		if err != nil {
			panic(err)
		}

		homeHandler := handlers.NewHomeHandler(cfg.PageHeader)
		errorHandler := handlers.NewErrorHandler(cfg.PageHeader)
		probeHandler := handlers.NewProbeHandler(backend)
		userHandler := handlers.NewUserHandler(backend)
		authHandler := handlers.NewAuthHandler(backend, &cfg.Auth, &cfg.Permissions)
		permissionHandler := handlers.NewPermissionsHandler(&cfg.Permissions)

		router.Use(middleware.Authed(backend, &cfg.Auth.Cookie, &cfg.Auth.JWT, &cfg.Auth, &cfg.Permissions, false))

		if !cfg.DisableHomepage {
			router.GET("/", homeHandler.Home())
		}

		if !cfg.DisableErrorPages {
			router.GET("/error/:code", errorHandler.Error(0))
			router.NoRoute(errorHandler.Error(http.StatusNotFound))
		}

		if !cfg.DisableHealthProbes {
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
		api.Use(middleware.Authed(backend, &cfg.Auth.Cookie, &cfg.Auth.JWT, &cfg.Auth, &cfg.Permissions, true))
		{
			users := api.Group("/users")
			{
				users.POST("/", middleware.HasAny("users:create:*"), userHandler.CreateUser())
				users.PUT(IDPath(), middleware.HasAny(IDPath("users:update:")), userHandler.UpdateUser())
				users.DELETE(IDPath(), middleware.HasAny(IDPath("users:delete:")), userHandler.DeleteUser())
				users.GET(IDPath(), middleware.HasAny(IDPath("users:get:")), userHandler.GetUser())
				users.GET("/", middleware.HasAny("users:list:*"), userHandler.ListUsers())
			}

			api.GET("/me", userHandler.Me())

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

		logger := dependencies.GetLogger()
		logger.Info().Int("port", cfg.Port).Msg("Starting listener")
		err = router.Run(fmt.Sprintf(":%d", cfg.Port))
		if err != nil {
			logger.Fatal().Err(err).Msg("Error starting listener")
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	mamba.MustBind(&models.Config{}, rootCmd, &mamba.Options{PrefixEmbedded: false})

	rootCmd.AddCommand(certsCmd)
}

func initConfig() {
	cfgFile := viper.GetString("configfile")

	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("Error unmarshalling config: %s", err))
	}
}

// IDPath will add the standard ID param to the given path
func IDPath(path ...string) string {
	p := ""
	if len(path) > 0 {
		p = path[0]
	}

	if strings.HasSuffix(p, ":") {
		return p + constants.IDParam
	}

	return strings.TrimSuffix(p, "/") + "/:" + constants.IDParam
}

// GetID gets the id from the path, setup using tonic.IDPath
func GetID(c *gin.Context) (id string) {
	return c.Param(constants.IDParam)
}
