package types

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"time"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type UserRequest struct {
	Username  string `json:"username" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=20"`
	FirstName string `json:"firstName" validate:"required,alpha"`
	LastName  string `json:"lastName" validate:"required,alpha"`
}

func (ur *UserRequest) Validate() error {
	return validate.Struct(ur)
}

type User struct {
	ID                    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
	Username              string    `gorm:"type:varchar(255);unique"`
	Password              string    `gorm:"type:varchar(255)"`
	FirstName             string    `gorm:"type:varchar(255)"`
	LastName              string    `gorm:"type:varchar(255)"`
	EmailVerificationUUID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	IsEmailVerified       bool      `gorm:"type:boolean;default:false"`
	IsEmailLinkExpired    bool      `gorm:"type:boolean;default:false"`
}

type UserResponse struct {
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

func (ur *UpdateUserRequest) Validate() error {
	return validate.Struct(ur)
}

type EmailVerification struct {
	EmailVerificationUUID uuid.UUID `json:"emailVerificationUUID"`
	VerificationLink      string    `json:"verificationLink"`
	Username              string    `json:"username"`
	FirstName             string    `json:"firstName"`
	LastName              string    `json:"lastName"`
}
