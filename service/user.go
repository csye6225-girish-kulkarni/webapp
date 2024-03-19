package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"webapp/apperror"
	"webapp/repository"
	"webapp/types"
)

type userService struct {
	repo repository.UserRepo
}

type UserService interface {
	CreateUser(ctx *gin.Context, userRequest types.UserRequest) (types.UserResponse, error)
	ValidateUser(ctx *gin.Context, username, password string) (bool, types.User, error)
	GetUserByUsername(ctx *gin.Context, username string) (types.User, error)
	UpdateUser(ctx *gin.Context, userRequest types.UpdateUserRequest) (types.UserResponse, error)
}

func NewUserService(repo repository.UserRepo) UserService {
	return &userService{
		repo: repo,
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

	log.Debug().Msg("User Created successfully")
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
