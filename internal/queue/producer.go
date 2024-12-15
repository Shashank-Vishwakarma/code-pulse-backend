package queue

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailVerificationPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Code     string `json:"code"`
}

func StartProducer(payload EmailVerificationPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = Ch.PublishWithContext(ctx, "", Q.Name, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		},
	)

	if err != nil {
		return err
	}

	return nil
}
