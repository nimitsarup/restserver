package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nimitsarup/restserver/handlers"
	hMock "github.com/nimitsarup/restserver/handlers/mock"
	"github.com/stretchr/testify/assert"
)

func TestAPI_AddUser(t *testing.T) {
	type fields struct {
		Handlers handlers.HandlersInterface
	}
	type args struct {
		resp *httptest.ResponseRecorder
		req  *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
	}{
		{
			name:   "WHENNoBodyTHEN400",
			fields: fields{},
			args: args{req: httptest.NewRequest(http.MethodPost, "/users", nil),
				resp: httptest.NewRecorder()},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "WHENBodyInvalidTHENCompareStatus",
			fields: fields{Handlers: &hMock.HandlersInterfaceMock{
				AddUserFunc: func(body string) (int, error) { return http.StatusInternalServerError, errors.New("test error") }},
			},
			args: args{req: httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"name", "email":"email", "password":"passwd"}`)),
				resp: httptest.NewRecorder()},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "WHENAllGoodTHEN200",
			fields: fields{Handlers: &hMock.HandlersInterfaceMock{
				AddUserFunc: func(body string) (int, error) { return http.StatusOK, nil }},
			},
			args: args{req: httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"name", "email":"email", "password":"passwd"}`)),
				resp: httptest.NewRecorder()},
			expectedStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{
				Handlers: tt.fields.Handlers,
			}
			a.AddUser(tt.args.resp, tt.args.req)

			assert.Equal(t, tt.expectedStatus, tt.args.resp.Code)
		})
	}
}
