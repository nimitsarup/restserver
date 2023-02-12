package main

import (
	"encoding/json"
	"log"
	"strings"
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

func HandleRequest(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	ApiResponse := events.APIGatewayV2HTTPResponse{}
	log.Printf("came in with request %v", request)
	if strings.HasPrefix(request.RouteKey, "GET") {
		id := request.PathParameters["id"]
		if id != "" {
			// get by id
			status, user := handler.GetUser(id)
			ApiResponse = events.APIGatewayV2HTTPResponse{Body: marshalToString(user), StatusCode: status}
		} else {
			// get all
			status, users := handler.GetAllUsers()
			ApiResponse = events.APIGatewayV2HTTPResponse{Body: marshalToString(users), StatusCode: status}
		}
	} else if strings.HasPrefix(request.RouteKey, "POST") {
		status, err := handler.AddUser(request.Body)
		body := ""
		if err != nil {
			body = "error adding user" + err.Error()
		}

		ApiResponse = events.APIGatewayV2HTTPResponse{Body: body, StatusCode: status}
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
