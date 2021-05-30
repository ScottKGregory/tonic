package dependencies

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/scottkgregory/tonic/internal/constants"
)

var Tag string

func GetLogger(c ...*gin.Context) *zerolog.Logger {
	if len(c) > 0 {
		l, b := c[0].Get(constants.LoggerKey)
		if b {
			return l.(*zerolog.Logger)
		}
	}

	ll := log.With().Str("tag", Tag).Logger()
	return &ll
}
