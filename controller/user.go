package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"net/http"
	"webapp/apperror"
	"webapp/service"
	"webapp/types"
)

type UserController interface {
	CreateUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	VerifyEmail(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (uc *userController) CreateUser(ctx *gin.Context) {
	var (
		request types.UserRequest
	)
	// Validating If Request Payload has any unknown fields
	j := json.NewDecoder(ctx.Request.Body)
	j.DisallowUnknownFields()
	err := j.Decode(&request)
	if err != nil {
		log.Error().Err(err).Msg("Bad Request")
		ctx.Status(http.StatusBadRequest)
		return
	}
	// Validating the Content of the Request
	err = request.Validate()
	if err != nil {
		log.Error().Err(err).Msg("Bad Request")
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = request.Validate()
	if err != nil {
		log.Error().Err(err).Msg("Bad Request")
		ctx.Status(http.StatusBadRequest)
		return
	}

	response, err := uc.userService.CreateUser(ctx, request)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23505" {
				log.Error().Msg("Username Already Exists")
				ctx.Status(http.StatusBadRequest)
				return
			}
		}
		log.Error().Err(err).Msg("Failed to create user")
		ctx.Status(http.StatusInternalServerError)
		return
	}

	log.Info().Msg("User created successfully")
	ctx.JSON(http.StatusCreated, response)
	return
}

func (uc *userController) GetUser(ctx *gin.Context) {
	var (
		response types.UserResponse
	)
	if ctx.Request.Body != http.NoBody {
		log.Error().Msg("Bad Request with error : Request has a payload")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Request query params validation
	if len(ctx.Request.URL.RawQuery) > 0 {
		log.Error().Msg("Bad Request with error : Request has query parameters")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, ok := ctx.Get("user")
	if !ok {
		log.Error().Msg("Unauthorized")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userDetails := user.(types.User)
	response = types.UserResponse{
		Username:  userDetails.Username,
		FirstName: userDetails.FirstName,
		LastName:  userDetails.LastName,
		CreatedAt: userDetails.CreatedAt,
		UpdatedAt: userDetails.UpdatedAt,
		ID:        userDetails.ID.String(),
	}

	log.Info().Msg("User details fetched successfully")
	ctx.JSON(http.StatusOK, response)
	return
}

func (uc *userController) UpdateUser(ctx *gin.Context) {
	var (
		request types.UpdateUserRequest
	)
	// Validating If Request Payload has any unknown fields
	j := json.NewDecoder(ctx.Request.Body)
	j.DisallowUnknownFields()
	err := j.Decode(&request)
	if err != nil {
		log.Error().Err(err).Msg("Bad Request")
		ctx.Status(http.StatusBadRequest)
		return
	}
	// Validating the Content of the Request
	err = request.Validate()
	if err != nil {
		log.Error().Err(err).Msg("Bad Request")
		ctx.Status(http.StatusBadRequest)
		return
	}

	_, err = uc.userService.UpdateUser(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Info().Msg("User updated successfully")
	ctx.Status(http.StatusNoContent)
	return
}

func (uc *userController) VerifyEmail(ctx *gin.Context) {
	emailUUID := ctx.Query("uuid")
	if emailUUID == "" {
		log.Error().Msg("Bad Request with error : Email UUID is empty")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := uc.userService.VerifyEmail(ctx, emailUUID)
	if err != nil {
		if errors.Is(err, apperror.ErrLinkExpired) {
			log.Error().Msg("Link has expired")
			ctx.Status(http.StatusGone)
			return
		}
		log.Error().Err(err).Msg("Failed to verify email")
		ctx.Status(http.StatusInternalServerError)
		return
	}

	log.Info().Msg("Email verified successfully")
	ctx.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
	return
}
