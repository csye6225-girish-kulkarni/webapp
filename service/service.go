package service

import (
	"Health-Check/repository"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Ping(ctx *gin.Context) error
}

type HealthService struct {
	repo repository.UserRepo
}

func NewHealthService(repo repository.UserRepo) *HealthService {
	return &HealthService{
		repo: repo,
	}
}

func (hs *HealthService) Ping(ctx *gin.Context) error {
	err := hs.repo.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}
