package router

import (
	"Health-Check/controller"
	"Health-Check/db"
	"Health-Check/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
	"os"
)

func InitializeRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	connString := os.Getenv("POSTGRES_CONN_STR")
	connString = connString + "?sslmode=disable"
	postgresRepo := db.NewPostgreSQL(connString)

	healthService := service.NewHealthService(postgresRepo)

	healthController := controller.NewHealthController(healthService)
	router.GET("/healthz", healthController.GetHealth)

	router.Use(func(context *gin.Context) {
		if context.Request.URL.Path == "/healthz" && context.Request.Method != http.MethodGet {
			context.Status(http.StatusMethodNotAllowed)
			context.Abort()
		}
	})

	router.NoRoute(func(context *gin.Context) {
		context.Data(http.StatusNotFound, "text/plain", []byte{})
		context.Abort()
	})

	return router
}
