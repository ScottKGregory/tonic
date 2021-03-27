package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/pkg/backends"
	"github.com/scottkgregory/tonic/pkg/helpers"
	"github.com/scottkgregory/tonic/pkg/models"
)

type StatusService struct {
	backend *backends.Mongo
}

func NewStatusService(backendOptions *models.Backend) *StatusService {
	return &StatusService{
		backend: backends.NewMongoBackend(backendOptions),
	}
}

func (s *StatusService) Health(c *gin.Context) {
	if err := s.backend.Ping(helpers.GetLogger(c)); err != nil {
		helpers.APIErrorResponse(c, http.StatusInternalServerError, "Error connecting to backend")
		return
	}

	helpers.APISuccessResponse(c, "Healthy")
}

func (s *StatusService) Readiness(c *gin.Context) {
	if err := s.backend.Ping(helpers.GetLogger(c)); err != nil {
		helpers.APIErrorResponse(c, http.StatusInternalServerError, "Error connecting to backend")
		return
	}

	helpers.APISuccessResponse(c, "Ready")
}

func (s *StatusService) Liveliness(c *gin.Context) {
	if err := s.backend.Ping(helpers.GetLogger(c)); err != nil {
		helpers.APIErrorResponse(c, http.StatusInternalServerError, "Error connecting to backend")
		return
	}

	helpers.APISuccessResponse(c, "Alive")
}
