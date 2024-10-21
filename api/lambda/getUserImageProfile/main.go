package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang-jwt/jwt/v5"
)

type s3ClientInterface interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func Handler(s3Client s3ClientInterface, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bucketName := "jfolgado-homeapp-userprofile-v1"

	tokenString := request.Headers["Authorization"]
	if tokenString == "" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized, Body: "Missing Authorization header"}, nil
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, _, _ := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized, Body: "Invalid token claims"}, nil
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized, Body: "sub claim not found in token"}, nil
	}

	imageKey := fmt.Sprintf("%s.png", sub)
	resp, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(imageKey),
	})

	if err != nil {
		log.Printf("failed to get image: %v from BucketName: %v", imageKey, bucketName)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: fmt.Sprintf("Error fetching image from S3, err: %v", err)}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: "Error reading image data"}, nil
	}

	encodedImage := base64.StdEncoding.EncodeToString(body)
	responseBody, err := json.Marshal(map[string]string{"image": encodedImage})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to encode JSON"}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBody),
	}, nil
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return Handler(s3Client, request)
	})
}
