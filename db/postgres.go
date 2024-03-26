package db

import (
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
	"webapp/types"
)

type PostgreSQL struct {
	DB *gorm.DB
}

// NewPostgreSQL constructor for PostgreSQL struct.
func NewPostgreSQL(connStr string) *PostgreSQL {
	pg := &PostgreSQL{}
	connected := make(chan bool)
	go pg.connectWithRetry(connStr, connected)

	select {
	case <-connected:
		return pg
	case <-time.After(30 * time.Second):
		log.Error().Msg("Failed to connect to the database")
		return nil
	}
}

func (p *PostgreSQL) connectWithRetry(connStr string, connected chan bool) {
	var db *gorm.DB
	var err error

	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
		if err == nil {
			p.DB = db
			log.Info().Msg("Database connection established")
			// Adding the UUID extension to the database
			db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

			err = db.AutoMigrate(&types.User{}, &types.Email{})
			if err != nil {
				log.Error().Err(err).Msg("Error while migrating the table")
				return
			}
			log.Info().Msg("User table migrated successfully")
			connected <- true
			return
		}
		log.Error().Err(err).Msgf("Attempt %d: could not connect to database", i)
		time.Sleep(time.Duration(5) * time.Second)
	}
	close(connected)
	log.Error().Err(err).Msgf("Failed to connect to the database after %d attempts", maxAttempts)
}

func (p *PostgreSQL) Ping(ctx *gin.Context) error {
	if p == nil || p.DB == nil {
		log.Debug().Msg("DB object is not initialized")
		return errors.New("DB object is not initialized")
	}

	db, err := p.DB.DB()
	if err != nil {
		log.Error().Err(err).Msg("Error getting the database object")
		return err
	}
	if err = db.Ping(); err != nil {
		log.Error().Err(err).Msg("Error pinging the database")
		return err
	}
	log.Info().Msg("Database pinged successfully")
	return nil
}

func (p *PostgreSQL) CreateUser(ctx *gin.Context, user types.User) (types.User, error) {
	if p == nil || p.DB == nil {
		log.Debug().Msg("DB object is not initialized")
		return types.User{}, errors.New("DB object is not initialized")
	}

	if err := p.DB.Create(&user).Error; err != nil {
		log.Error().Err(err).Msg("Error creating the user")
		return types.User{}, err
	}

	log.Info().Msg("User created successfully")
	return user, nil
}

func (p *PostgreSQL) GetByUsername(ctx *gin.Context, username string) (types.User, error) {
	if p == nil || p.DB == nil {
		log.Debug().Msg("DB object is not initialized")
		return types.User{}, errors.New("DB object is not initialized")
	}

	var user types.User
	if err := p.DB.Where("username = ?", username).First(&user).Error; err != nil {
		log.Error().Err(err).Msg("Error getting the user by username")
		return types.User{}, err
	}
	return user, nil
}

func (p *PostgreSQL) UpdateUser(ctx *gin.Context, user types.User) (types.User, error) {
	if p == nil || p.DB == nil {
		log.Debug().Msg("DB object is not initialized")
		return types.User{}, errors.New("DB object is not initialized")
	}
	// The user details are stored in the context in the middleware
	u, ok := ctx.Get("user")
	if !ok {
		log.Error().Msg("User not found in context")
		return types.User{}, errors.New("user not found in context")
	}
	existingUser := u.(types.User)

	if err := p.DB.Model(&existingUser).Updates(user).Error; err != nil {
		log.Error().Err(err).Msg("Error updating the user")
		return types.User{}, err
	}
	log.Info().Msg("User updated successfully")
	return user, nil
}

func (p *PostgreSQL) Close() error {
	if p == nil || p.DB == nil {
		return errors.New("DB object is not initialized")
	}
	db, err := p.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (p *PostgreSQL) MarkEmailAsVerified(ctx *gin.Context, userID string) error {
	if p == nil || p.DB == nil {
		log.Debug().Msg("DB object is not initialized")
		return errors.New("DB object is not initialized")
	}

	// Get the user by the email verification UUID
	var user types.User
	if err := p.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		log.Debug().Err(err).Msg("Error getting the user by email verification UUID")
		return err
	}

	// Check if the user's email is already verified
	if user.IsEmailVerified {
		log.Debug().Msg("Email is already verified")
		return errors.New("email is already verified")
	}

	// Mark the user's email as verified
	if err := p.DB.Model(&user).Update("is_email_verified", true).Error; err != nil {
		log.Error().Err(err).Msg("Error updating the user")
		return err
	}

	log.Info().Msg("Email verified successfully")
	return nil
}

func (p *PostgreSQL) GetByEmailVerificationUUID(ctx *gin.Context, uuid string) (types.User, types.Email, error) {
	if p == nil || p.DB == nil {
		log.Debug().Msg("DB object is not initialized")
		return types.User{}, types.Email{}, errors.New("DB object is not initialized")
	}

	var email types.Email
	if err := p.DB.Where("email_verification_uuid = ?", uuid).First(&email).Error; err != nil {
		log.Error().Err(err).Msg("Error getting the email by email verification UUID")
		return types.User{}, types.Email{}, err
	}

	var user types.User
	if err := p.DB.Where("id = ?", email.UserID).First(&user).Error; err != nil {
		log.Error().Err(err).Msg("Error getting the user by user ID")
		return types.User{}, types.Email{}, err
	}
	return user, email, nil
}

//func (p *PostgreSQL) Exec(query string, args ...interface{}) error {
//	if p == nil || p.DB == nil {
//		return errors.New("DB object is not initialized")
//	}
//	return p.DB.Exec(query, args...).Error
//}
