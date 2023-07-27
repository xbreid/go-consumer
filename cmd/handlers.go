package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-playground/validator/v10"
	"go-consumer/models"
	"log"
)

type Operation string
type Resource string

const (
	CREATE Operation = "CREATE"
	UPDATE Operation = "UPDATE"
	UPSERT Operation = "UPSERT"
	DELETE Operation = "DELETE"
)

const (
	ACCOUNT_GROUP Resource = "ACCOUNT_GROUP"
)

type MsgType struct {
	Host      string      `json:"host"`
	Resource  Resource    `json:"resource" validate:"oneof=ACCOUNT_GROUP"`
	Operation Operation   `json:"operation" validate:"oneof=CREATE UPDATE UPSERT DELETE"`
	Payload   interface{} `json:"payload"`
}

func (s *MsgType) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

func (app *Config) processSQS(ctx context.Context) (bool, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            &app.queueUrl,
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   app.visibilityTimeout,
		WaitTimeSeconds:     app.waitingTimeout, // use long polling
	}

	resp, err := app.SQS.ReceiveMessage(ctx, input)

	if err != nil {
		return false, fmt.Errorf("error receiving message %w", err)
	}

	log.Printf("received messages: %v", len(resp.Messages))
	if len(resp.Messages) == 0 {
		return false, nil
	}

	for _, msg := range resp.Messages {
		var newMsg MsgType
		id := *msg.MessageId

		err := json.Unmarshal([]byte(*msg.Body), &newMsg)
		if err != nil {
			return false, fmt.Errorf("error unmarshalling %w", err)
		}

		log.Printf("message id %s is received from SQS: %#v", id, newMsg)

		err = newMsg.Validate()
		if err != nil {
			log.Printf("message id %s has invalid format", id)
			return false, nil
		}

		payload, err := json.Marshal(newMsg.Payload)
		if err != nil {
			return false, fmt.Errorf("error marshalling %w", err)
		}

		switch newMsg.Resource {
		case ACCOUNT_GROUP:
			var data models.AccountGroup

			err = json.Unmarshal(payload, &data)
			if err != nil {
				return false, fmt.Errorf("error unmarshalling %w", err)
			}

			err = app.handleAccountGroup(ctx, data)
			if err != nil {
				return false, fmt.Errorf("error handling account group %w", err)
			}
		default:
			fmt.Printf("Unsupported resource")
		}

		_, err = app.SQS.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      &app.queueUrl,
			ReceiptHandle: msg.ReceiptHandle,
		})
		if err != nil {
			return false, fmt.Errorf("error deleting message from SQS %w", err)
		}

		log.Printf("message id %s is deleted from queue", id)
	}

	return true, nil
}

func (app *Config) handleAccountGroup(ctx context.Context, data models.AccountGroup) error {
	_, err := models.UpsertAccountGroup(ctx, app.DB, data)
	if err != nil {
		return err
	}

	return nil
}
