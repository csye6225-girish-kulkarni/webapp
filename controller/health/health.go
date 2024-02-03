package health

import (
	"Health-Check/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Controller interface {
	GetHealth(context *gin.Context)
}

type healthContoller struct {
	healthService *service.HealthService
}

func NewHealthController(hs *service.HealthService) Controller {
	return &healthContoller{
		healthService: hs,
	}
}

func (c *healthContoller) GetHealth(ctx *gin.Context) {
	ctx.Header("cache-control", "no-cache")
	// Request Payload validation
	if ctx.Request.ContentLength > 0 {
		log.Println("Request has a payload")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Request query params validation
	if len(ctx.Request.URL.RawQuery) > 0 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := c.healthService.Ping(ctx)
	if err != nil {
		log.Println("Unable to Ping to DB err : %v", err)
		ctx.Status(http.StatusServiceUnavailable)
		return
	}
	log.Println("Database Successfully pinged")
	ctx.Status(http.StatusOK)
	return
}
