package tonic

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/scottkgregory/tonic/internal/api"
	"github.com/scottkgregory/tonic/internal/api/errors"
	"github.com/scottkgregory/tonic/internal/backends"
	"github.com/scottkgregory/tonic/internal/constants"
	"github.com/scottkgregory/tonic/internal/dependencies"
	"github.com/scottkgregory/tonic/internal/middleware"
	"github.com/scottkgregory/tonic/internal/models"
	"github.com/scottkgregory/tonic/pkg/handlers"
	pkgModels "github.com/scottkgregory/tonic/pkg/models"
)

// Init sets up tonic
func Init(opt models.Options) (*gin.Engine, *gin.RouterGroup, error) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(middleware.Zerologger(opt.Log))
	log := dependencies.GetLogger()

	var backend backends.Backend
	var err error
	if opt.Backend.InMemory {
		backend = backends.NewMemoryBackend(&opt.Backend)
	} else {
		backend, err = backends.NewMongoBackend(context.Background(), &opt.Backend)
	}
	if err != nil {
		return nil, nil, err
	}

	homeHandler := handlers.NewHomeHandler(opt.PageHeader)
	errorHandler := handlers.NewErrorHandler(opt.PageHeader)
	probeHandler := handlers.NewProbeHandler(backend)
	userHandler := handlers.NewUserHandler(backend)
	authHandler := handlers.NewAuthHandler(backend, &opt.Auth, &opt.Permissions)
	permissionHandler := handlers.NewPermissionsHandler(&opt.Permissions)

	router.Use(middleware.Authed(backend, &opt.Auth.Cookie, &opt.Auth.JWT, &opt.Auth, &opt.Permissions, false))

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
	api.Use(middleware.Authed(backend, &opt.Auth.Cookie, &opt.Auth.JWT, &opt.Auth, &opt.Permissions, true))
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
	return router, api, nil
}

// GetLogger returns a zerolog logger with context
func GetLogger(c ...*gin.Context) *zerolog.Logger {
	return dependencies.GetLogger(c...)
}

// GetUser returns the current user from context
func GetUser(c *gin.Context) (user *pkgModels.User, ok bool) {
	return dependencies.GetUser(c)
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

// Surface errors
type ForbiddenErr = errors.ForbiddenErr
type NotFoundErr = errors.NotFoundErr
type UnauthorisedErr = errors.UnauthorisedErr
type ValidationErr = errors.ValidationErr
type GenericErr = errors.GenericErr

// NewForbiddenError creates a new ForbiddenErr
func NewForbiddenError(required ...string) *ForbiddenErr {
	return errors.NewForbiddenError(required...)
}

// NewNotFoundError creates a new NotFoundErr
func NewNotFoundError(id string) *NotFoundErr {
	return errors.NewNotFoundError(id)
}

// NewUnauthorisedError creates a new UnauthorisedErr
func NewUnauthorisedError() *UnauthorisedErr {
	return errors.NewUnauthorisedError()
}

// NewValidationError creates a new ValidationErr
func NewValidationError(validation ...map[string]string) *ValidationErr {
	return errors.NewValidationError(validation...)
}

// NewGenericError creates a new GenericErr
func NewGenericError(err ...error) *GenericErr {
	return errors.NewGenericError(err...)
}

// Surface structs
type Options = models.Options
type LogOptions = models.LogOptions
type AuthOptions = models.AuthOptions
type PermissionsOptions = models.PermissionsOptions
type JWTOptions = models.JWTOptions
type OIDCOptions = models.OIDCOptions
type CookieOptions = models.CookieOptions
type BackendOptions = models.BackendOptions

// Surface middleware

// Authed is a middleware that requires a user be authenticated.
// If cancel is set to false the request will be allowed to continue, but tonic.GetUser will return false
func Authed(backend backends.Backend,
	cookieOptions *models.CookieOptions,
	jwtOptions *models.JWTOptions,
	authOptions *models.AuthOptions,
	permissionOptions *models.PermissionsOptions,
	cancel bool) gin.HandlerFunc {
	return middleware.Authed(backend,
		cookieOptions,
		jwtOptions,
		authOptions,
		permissionOptions,
		cancel,
	)
}

// HasAny is a middleware that requires the user has any of the provided permissions to pass
func HasAny(required ...string) gin.HandlerFunc {
	return middleware.HasAny(required...)
}

// HasAny is a middleware that requires the user has all of the provided permissions to pass
func HasAll(required ...string) gin.HandlerFunc {
	return middleware.HasAll(required...)
}

// Surface API responses

type ResponseModel = api.ResponseModel

//SmartResponse returns a response object appropriate to the supplied error
func SmartResponse(c *gin.Context, data interface{}, err error) {
	api.SmartResponse(c, data, err)
}

//NoContentResponse returns a new NoContentResponse, use at the end of a request
func NoContentResponse(c *gin.Context) {
	api.NoContentResponse(c)
}

//SuccessResponse returns a new SuccessResponse, use at the end of a request
func SuccessResponse(c *gin.Context, data interface{}) {
	api.SuccessResponse(c, data)
}

//UnauthorisedResponse returns a new UnauthorisedResponse, use at the end of a request
func UnauthorisedResponse(c *gin.Context, errs ...*UnauthorisedErr) {
	api.UnauthorisedResponse(c, errs...)
}

//ForbiddenResponse returns a new ForbiddenResponse, use at the end of a request
func ForbiddenResponse(c *gin.Context, errs ...*ForbiddenErr) {
	api.ForbiddenResponse(c, errs...)
}

//NotFoundResponse returns a new NotFoundResponse, use at the end of a request
func NotFoundResponse(c *gin.Context, errs ...*NotFoundErr) {
	api.NotFoundResponse(c, errs...)
}

//ValidationErrorResponse returns a new ValidationErrorResponse, use at the end of a request
func ValidationErrorResponse(c *gin.Context, errs ...*ValidationErr) {
	api.ValidationErrorResponse(c, errs...)
}

//ErrorResponse returns a new ErrorResponse, use at the end of a request
func ErrorResponse(c *gin.Context, code int, err error, validation ...map[string]string) {
	api.ErrorResponse(c, code, err, validation...)
}

//GlobalKey is a const for global validation errors
const GlobalKey = constants.GlobalKey
