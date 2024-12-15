package queue

import (
	"encoding/json"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/services"
	"github.com/sirupsen/logrus"
)

func StartConsumer() error {
	err := Ch.Qos(1, 0, false)
	if err != nil {
		return err
	}

	messages, err := Ch.Consume(Q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for message := range messages {
			var payload EmailVerificationPayload
			err := json.Unmarshal(message.Body, &payload)
			if err != nil {
				logrus.Errorf("Error unmarshalling the message: %v", err)
			}

			err = services.SendEmail(payload.Email, payload.Username, payload.Code)
			if err != nil {
				logrus.Errorf("Error sending email: %v", err)
			} else {
				logrus.Infof("Email sent to %s", payload.Email)
			}

			message.Ack(false)
		}
	}()

	return nil
}
