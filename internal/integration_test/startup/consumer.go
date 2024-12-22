package startup

import (
	"bitbucket.org/starlinglabs/cst-wstyle-integration/config"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/event"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/event/orderfile/aws/upload_consumer"
	"context"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func InitSQS() *sqs.Client {
	awsRegion := config.Get().AWSRegion

	cfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}
	sqsClient := sqs.NewFromConfig(cfg)

	return sqsClient
}

func NewConsumers(orderFileBatchConsumer *upload_consumer.OrderFileEventConsumer) map[string]event.Consumer {
	consumerMap := make(map[string]event.Consumer)
	consumerMap["article"] = orderFileBatchConsumer
	return consumerMap
}
