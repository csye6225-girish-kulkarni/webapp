package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"time"
	"webapp/apperror"
	"webapp/repository"
	"webapp/types"
)

type userService struct {
	repo         repository.UserRepo
	emailService EmailService
}

type UserService interface {
	CreateUser(ctx *gin.Context, userRequest types.UserRequest) (types.UserResponse, error)
	ValidateUser(ctx *gin.Context, username, password string) (bool, types.User, error)
	GetUserByUsername(ctx *gin.Context, username string) (types.User, error)
	UpdateUser(ctx *gin.Context, userRequest types.UpdateUserRequest) (types.UserResponse, error)
	VerifyEmail(ctx *gin.Context, uuid string) error
}

func NewUserService(repo repository.UserRepo, emailService EmailService) UserService {
	return &userService{
		repo:         repo,
		emailService: emailService,
	}
}

func (us *userService) CreateUser(ctx *gin.Context, userRequest types.UserRequest) (types.UserResponse, error) {
	var (
		user types.User
	)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Error hashing the password")
		return types.UserResponse{}, err
	}
	user = types.User{
		Username:  userRequest.Username,
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Password:  string(hashedPassword),
	}

	updatedUser, err := us.repo.CreateUser(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("Error creating the user")
		return types.UserResponse{}, err
	}

	// Send the verification email in a goroutine to not block the main thread
	go func() {
		newCtx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		err := us.emailService.SendVerificationEmailToQueue(newCtx, updatedUser)
		if err != nil {
			log.Error().Err(err).Msg("Error sending the verification email")
		}
	}()

	log.Info().Msg("User Created successfully")
	return types.UserResponse{
		Username:  updatedUser.Username,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		ID:        updatedUser.ID.String(),
	}, nil
}

func (us *userService) ValidateUser(ctx *gin.Context, username, password string) (bool, types.User, error) {
	user, err := us.repo.GetByUsername(ctx, username)
	if err != nil {
		log.Error().Err(err).Msg("Error getting the user by username")
		return false, types.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Error().Err(err).Msg("Incorrect Password")
		return false, types.User{}, apperror.ErrIncorrectPassword
	}

	if !user.IsEmailVerified {
		log.Error().Msg("Email not verified")
		return false, user, apperror.ErrEmailNotVerified
	}

	return true, user, nil
}

func (us *userService) GetUserByUsername(ctx *gin.Context, username string) (types.User, error) {
	user, err := us.repo.GetByUsername(ctx, username)
	if err != nil {
		log.Error().Err(err).Msg("Error getting the user by username")
		return types.User{}, err
	}
	return user, nil
}

func (us *userService) UpdateUser(ctx *gin.Context, userRequest types.UpdateUserRequest) (types.UserResponse, error) {
	var (
		user types.User
	)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Error hashing the password")
		return types.UserResponse{}, errors.New("incorrect password")
	}

	user = types.User{
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Password:  string(hashedPassword),
	}

	updatedUser, err := us.repo.UpdateUser(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("Error updating the user")
		return types.UserResponse{}, err
	}

	log.Debug().Msg("User updated successfully")
	return types.UserResponse{
		Username:  updatedUser.Username,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		ID:        updatedUser.ID.String(),
	}, nil
}

func (us *userService) VerifyEmail(ctx *gin.Context, uuid string) error {
	user, email, err := us.repo.GetByEmailVerificationUUID(ctx, uuid)
	if err != nil {
		log.Error().Err(err).Msg("Error getting the user by email verification id")
		return err
	}

	if time.Now().After(email.EmailVerificationExpiry) {
		log.Error().Msg("Link has expired")
		return apperror.ErrLinkExpired
	}

	err = us.repo.MarkEmailAsVerified(ctx, user.ID.String())
	if err != nil {
		log.Error().Err(err).Msg("Error marking the email as verified")
		return err
	}

	log.Info().Msg("Email verified successfully")
	return nil
}
