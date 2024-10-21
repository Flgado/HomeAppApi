// handler_test.go
package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

func generateValidToken(sub string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
	})
	secret := []byte("your_secret_key")
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return ""
	}
	return tokenString
}

func TestHandler_Success(t *testing.T) {
	// assert
	mockS3Client := new(MockS3Client)

	mockS3Client.On("GetObject", mock.Anything, mock.AnythingOfType("*s3.GetObjectInput")).Return(&s3.GetObjectOutput{
		Body: io.NopCloser(strings.NewReader("fake_image_data")),
	}, nil)

	validToken := generateValidToken("user123")
	request := events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Authorization": "Bearer " + validToken,
		},
	}

	// act
	response, err := Handler(mockS3Client, request)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody map[string]string
	err = json.Unmarshal([]byte(response.Body), &responseBody)
	assert.NoError(t, err)

	expectedImage := base64.StdEncoding.EncodeToString([]byte("fake_image_data"))
	assert.Equal(t, expectedImage, responseBody["image"])

	mockS3Client.AssertExpectations(t)
}
