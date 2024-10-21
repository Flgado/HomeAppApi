package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type GroupConfiguration struct {
	GroupsToOverride   []string `json:"groupsToOverride"`
	IamRolesToOverride []string `json:"iamRolesToOverride"`
	PreferredRole      string   `json:"preferredRole"`
}

type UserAttributes struct {
	Status              string `json:"cognito:user_status"`
	Email               string `json:"email"`
	EmailVerified       string `json:"email_verified"`
	PhoneNumber         string `json:"phone_number"`
	PhoneNumberVerified string `json:"phone_number_verified"`
	Sub                 string `json:"sub"`
}

type CallerContext struct {
	AwsSdkVersion string `json:"awsSdkVersion"`
	ClientId      string `json:"clientId"`
}

type Request struct {
	GroupConfiguration GroupConfiguration `json:"groupConfiguration"`
	UserAttributes     UserAttributes     `json:"userAttributes"`
}

type Event struct {
	Version       string                 `json:"version"`
	TriggerSource string                 `json:"triggerSource"`
	Region        string                 `json:"region"`
	UserPoolId    string                 `json:"userPoolId"`
	UserName      string                 `json:"userName"`
	CallerContext CallerContext          `json:"callerContext"`
	Request       Request                `json:"request"`
	Response      map[string]interface{} `json:"response"`
}

func handler(ctx context.Context, event Event) (Event, error) {
	newScopes := []string{"openid", "profile", "email"}
	claimstoSuppress := []string{"aws.cognito.signin.user.admin"}
	for _, group := range event.Request.GroupConfiguration.GroupsToOverride {
		newScope := fmt.Sprintf("%s-%s", group, event.CallerContext.ClientId)
		newScopes = append(newScopes, newScope)
	}

	event.Response = map[string]interface{}{
		"claimsAndScopeOverrideDetails": map[string]interface{}{
			"accessTokenGeneration": map[string]interface{}{
				"scopesToAdd":      newScopes,
				"claimsToSuppress": claimstoSuppress,
			},
		},
	}

	return event, nil
}

func main() {
	lambda.Start(handler)
}
