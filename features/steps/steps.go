package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"strconv"

	"github.com/cucumber/godog"
	"github.com/nimitsarup/restserver/model"
)

func (c *UsersApiComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I have these users in the database:$`, c.iHaveTheseUsers)
	ctx.Step(`^I GET "([^"]*)"$`, c.IGet)
	ctx.Step(`^I should receive the following model response with status "([^"]*)":$`, c.IShouldReceiveTheFollowingModelResponse)
	ctx.Step(`^the HTTP status code should be "([^"]*)"$`, c.TheHTTPStatusCodeShouldBe)
}

func (c *UsersApiComponent) iHaveTheseUsers(datasetsJson *godog.DocString) error {
	var componentTestData []struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.Unmarshal([]byte(datasetsJson.Content), &componentTestData)
	if err != nil {
		return err
	}

	for _, testData := range componentTestData {
		if err := c.putUserInDatabase(model.User{Id: testData.Id, Name: testData.Name, Email: testData.Email, Password: testData.Password}); err != nil {
			return err
		}
	}

	return nil
}

func (c *UsersApiComponent) IShouldReceiveTheFollowingModelResponse(expectedCodeStr string, expectedAPIResponse *godog.DocString) error {
	var expected, actual model.User
	err := c.toModel(expectedCodeStr, expectedAPIResponse, &expected, &actual)
	if err != nil {
		return err
	}

	if expected.Id != actual.Id {
		return fmt.Errorf("expecting %v, got %v", expected, actual)
	}

	return nil
}

func (c *UsersApiComponent) IGet(path string) error {
	return c.makeRequest("GET", path, nil)
}

func (c *UsersApiComponent) toModel(expectedCodeStr string, expectedAPIResponse *godog.DocString, expected, actual interface{}) error {
	if err := c.TheHTTPStatusCodeShouldBe(expectedCodeStr); err != nil {
		return err
	}

	err := json.Unmarshal([]byte(expectedAPIResponse.Content), expected)
	if err != nil {
		return err
	}

	responseBody := c.HttpResponse.Body
	body, _ := io.ReadAll(responseBody)
	err = json.Unmarshal(body, actual)
	if err != nil {
		return err
	}

	return nil
}

func (c *UsersApiComponent) TheHTTPStatusCodeShouldBe(expectedCodeStr string) error {
	expectedCode, err := strconv.Atoi(expectedCodeStr)
	if err != nil {
		return err
	}
	if expectedCode != c.HttpResponse.StatusCode {
		return fmt.Errorf("expected %d, recieved %d", expectedCode, c.HttpResponse.StatusCode)
	}
	return nil
}

func (c *UsersApiComponent) putUserInDatabase(user model.User) error {
	return c.DB.AddUser(user)
}

func (c *UsersApiComponent) makeRequest(method, path string, data []byte) error {
	handler, err := c.Initialiser()
	if err != nil {
		return err
	}
	req := httptest.NewRequest(method, "http://foo"+path, bytes.NewReader(data))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	c.HttpResponse = w.Result()
	return nil
}
