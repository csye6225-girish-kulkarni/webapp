package service

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
	"webapp/types"
)

type EmailService interface {
	SendVerificationEmailToQueue(ctx context.Context, user types.User) error
}

type RealEmailService struct{}

func NewEmailService() EmailService {
	return &RealEmailService{}
}

func (es *RealEmailService) SendVerificationEmailToQueue(ctx context.Context, user types.User) error {
	log.Info().Msg("Sending the verification email")
	client, err := pubsub.NewClient(ctx, "cloud-csye6225-dev")
	if err != nil {
		log.Error().Err(err).Msg("Error creating the pubsub client")
		return err
	}

	// Prepare the user details and verification link
	expiryTime := time.Now().Add(2 * time.Minute).Unix() // 2 minutes expiry time
	verificationLink := "http://girishkulkarni.me:8080/v1/verify-email?uuid=" +
		user.EmailVerificationUUID.String() +
		"&expiry=" + strconv.FormatInt(expiryTime, 10)
	log.Info().Str("verificationLink", verificationLink).Msg("Verification Link")

	userDetails := types.EmailVerification{
		EmailVerificationUUID: user.EmailVerificationUUID,
		VerificationLink:      verificationLink,
		Username:              user.Username,
		FirstName:             user.FirstName,
		LastName:              user.LastName,
	}

	userDetailsJson, err := json.Marshal(userDetails)
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling the user details")
		return err
	}
	// Publish the message to the queue
	topic := client.Topic("email-sender")
	result := topic.Publish(ctx, &pubsub.Message{
		Data: userDetailsJson,
	})

	// Verify if the message was published successfully
	_, err = result.Get(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error publishing the message to the queue")
		return err
	}

	log.Info().Msg("Message published successfully")
	return nil
}
