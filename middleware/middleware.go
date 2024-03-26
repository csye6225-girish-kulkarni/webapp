package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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
					log.Error().Err(err).Msg("Incorrect Password")
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				}
				if errors.Is(err, apperror.ErrEmailNotVerified) {
					if gin.Mode() == gin.TestMode {
						c.Set("user", fetchedUser)
						c.Next()
						return
					}
					log.Error().Err(err).Msg("Email not verified")
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				}
				log.Error().Err(err).Msg("Error validating the user")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			if isValid {
				log.Info().Msg("User validated successfully")
				// Storing the user details in the context to avoid redundant calls to the database
				c.Set("user", fetchedUser)
				c.Next()
				return
			}
		}

		log.Error().Msg("Error validating the user")
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
