package middleware

import (
	"Health-Check/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BasicAuth(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, hasAuth := c.Request.BasicAuth()

		if hasAuth {
			isValid, fetchedUser, err := userService.ValidateUser(c, username, password)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
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
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}
