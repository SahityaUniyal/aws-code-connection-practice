package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var snsClient *sns.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1")) // Replace "us-east-1" with your AWS region
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}

	// Create SNS client
	snsClient = sns.NewFromConfig(cfg)
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	snsTopicArn := os.Getenv("SNS_TOPIC")
	log.Printf("SNS ARN :: %v", snsTopicArn)
	for _, message := range sqsEvent.Records {
		// Process each SQS message
		fmt.Printf("Processing message ID: %s, Body: %s\n", message.MessageId, message.Body)

		// Publish the message body to the SNS topic
		input := &sns.PublishInput{
			Message:  aws.String(message.Body),
			TopicArn: aws.String(snsTopicArn),
		}

		_, err := snsClient.Publish(context.Background(), input)
		if err != nil {
			log.Printf("L1 :: Failed to publish message to SNS topic: %v", err)
			return err
		}

		log.Printf("L1 :: Message successfully published to SNS topic: %s", snsTopicArn)
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
