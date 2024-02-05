package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"webapp/controller"
	"webapp/db"
	"webapp/middleware"
	"webapp/service"
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
	userService := service.NewUserService(postgresRepo)
	userController := controller.NewUserController(userService)
	healthController := controller.NewHealthController(healthService)

	router.GET("/healthz", healthController.GetHealth)
	router.Use(func(context *gin.Context) {
		if context.Request.URL.Path == "/healthz" && context.Request.Method != http.MethodGet {
			context.Status(http.StatusMethodNotAllowed)
			context.Abort()
		}
	})

	router.GET("/v1/user", middleware.BasicAuth(userService), userController.GetUser)
	router.PUT("/v1/user/self", middleware.BasicAuth(userService), userController.UpdateUser)

	router.POST("/v1/user", userController.CreateUser)
	router.NoRoute(func(context *gin.Context) {
		context.Data(http.StatusNotFound, "text/plain", []byte{})
		context.Abort()
	})

	return router
}
