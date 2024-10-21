package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iotdataplane"
)

type iotTopic struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var toPublish iotTopic

	err := json.Unmarshal([]byte(request.Body), &toPublish)
	if err != nil {
		log.Printf("Failed to decode JSON: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, nil
	}

	err = publishMessage(toPublish.Topic, toPublish.Message)
	if err != nil {
		log.Printf("failed to publish message: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func publishMessage(topic, message string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := iotdataplane.NewFromConfig(cfg)

	payload := map[string]string{
		"message": message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	input := &iotdataplane.PublishInput{
		Topic:   &topic,
		Payload: []byte(jsonPayload),
		Qos:     0,
	}

	_, err = client.Publish(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
