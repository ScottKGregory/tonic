package middleware

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/dependencies"
	"github.com/scottkgregory/tonic/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Zerologger(config models.LogConfig) gin.HandlerFunc {
	if !config.JSON {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	dependencies.Tag = config.Tag

	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing log level, defaulting to trace")
		log.Logger = log.Level(zerolog.TraceLevel)
	} else {
		log.Logger = log.Level(level)
	}

	return func(c *gin.Context) {
		rid := c.Request.Header.Get("x-request-id")
		if rid == "" {
			rid = primitive.NewObjectID().Hex()
		}

		t := time.Now()

		l := populatedLogger(c, rid)
		if !containsStr(config.IgnoreRoutes, c.FullPath()) {
			l.Info().Msg("Requested")
		}

		c.Set(constants.LoggerKey, l)
		c.Set(constants.RequestIDKey, rid)
		c.Next()

		l = populatedLogger(c, rid)

		statusCode := c.Writer.Status()
		ll := l.Info()
		switch {
		case statusCode >= 400 && statusCode < 500:
			ll = l.Warn()
		case statusCode >= 500:
			ll = l.Error()
		}

		if !containsStr(config.IgnoreRoutes, c.FullPath()) {
			ll.
				Int("status", statusCode).
				Dur("duration-ns", time.Duration(time.Since(t).Nanoseconds())).
				Msg("Returned")
		}
	}
}

func populatedLogger(c *gin.Context, rid string) *zerolog.Logger {
	l := dependencies.GetLogger(c).
		With().
		Str("client-ip", c.ClientIP()).
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Str("request-id", rid).
		Logger()

	return &l
}

func containsStr(arr []string, s string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}

	return false
}
