package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"time"
)

func (app *Config) GetQueueURL(c context.Context, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return app.SQS.GetQueueUrl(c, input)
}

func (app *Config) SendSQSMessage(ctx context.Context, input *sqs.SendMessageInput) error {
	_, err := app.SQS.SendMessage(ctx, input)

	return err
}

func (app *Config) GetUTCTimestampNow() string {
	t := time.Now().UTC()
	return t.Format("2006-01-02T15:04:05.000Z")
}
