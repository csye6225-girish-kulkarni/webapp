package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"webapp/apperror"
	"webapp/service"
)

func BasicAuth(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, hasAuth := c.Request.BasicAuth()

		if hasAuth {
			isValid, fetchedUser, err := userService.ValidateUser(c, username, password)
			if err != nil {
				if errors.Is(err, apperror.ErrIncorrectPassword) {
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				}
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			if isValid {
				// Storing the user details in the context to avoid redundant calls to the database
				c.Set("user", fetchedUser)
				c.Next()
				return
			}
		}

		c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func SetNoCacheHeader() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Header("cache-control", "no-cache")
		context.Next()
	}
}

func CheckNoAuthEndpoints() gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.GetHeader("Authorization") != "" {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}
		context.Next()
	}
}
