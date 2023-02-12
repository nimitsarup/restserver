package main

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nimitsarup/restserver/config"
	"github.com/nimitsarup/restserver/handlers"
	"github.com/nimitsarup/restserver/service"
)

const LambdaTimeoutDuration = 15 * time.Second

var svc, _ = service.NewServices(&config.Config{IsRunningInCloud: true}, nil)
var handler = handlers.Handlers{DB: svc.GetDB()}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ApiResponse := events.APIGatewayProxyResponse{}
	switch request.HTTPMethod {
	case "GET":
		id := request.PathParameters["id"]

		if id != "" {
			// get by id
			status, user := handler.GetUser(id)
			ApiResponse = events.APIGatewayProxyResponse{Body: marshalToString(user), StatusCode: status}
		} else {
			// get all
			status, users := handler.GetAllUsers()
			ApiResponse = events.APIGatewayProxyResponse{Body: marshalToString(users), StatusCode: status}
		}

	case "POST":
		status, err := handler.AddUser(request.Body)
		body := ""
		if err != nil {
			body = "error adding user" + err.Error()

		}

		ApiResponse = events.APIGatewayProxyResponse{Body: body, StatusCode: status}
	}
	// Response
	return ApiResponse, nil
}

func marshalToString(resp interface{}) string {
	b, _ := json.Marshal(resp)
	return string(b)
}

func main() {
	lambda.Start(HandleRequest)
}
