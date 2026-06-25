package auth_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang/mock/gomock"
	ports_errors "github.com/rickferrdev/gopher-login/internal/core/ports/errors"
	ports_services "github.com/rickferrdev/gopher-login/internal/core/ports/services"
	"github.com/rickferrdev/gopher-login/internal/inbound/http/rest/handlers/auth"
	"github.com/rickferrdev/gopher-login/internal/infra/server"
	"github.com/rickferrdev/gopher-login/internal/infra/struct_validator"
	ports_mocks "github.com/rickferrdev/gopher-login/internal/tests/mocks/ports"
	"github.com/stretchr/testify/suite"
)

type test struct {
	suite.Suite

	app     *fiber.App
	ctrl    *gomock.Controller
	service *ports_mocks.MockAuth
	handler *auth.Handler
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(test))
}

func (suite *test) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.service = ports_mocks.NewMockAuth(suite.ctrl)

	validator, err := struct_validator.New()
	suite.Require().NoError(err)

	app, group, err := server.New(validator)
	suite.Require().NoError(err)

	suite.app = app

	handler, err := auth.New(auth.Params{
		Router:      group,
		AuthService: suite.service,
	})
	suite.Require().NoError(err)

	suite.handler = handler
}

func (suite *test) TearDownTest() {
	suite.ctrl.Finish()

	if suite.app != nil {
		suite.Require().NoError(suite.app.Shutdown())
	}
}

func (suite *test) ExpectError(body io.ReadCloser, expect error) {
	raw, err := io.ReadAll(body)
	suite.Require().NoError(err)

	var got ports_errors.Error
	suite.Require().NoError(json.Unmarshal(raw, &got))

	var expected *ports_errors.Error
	suite.Require().ErrorAs(expect, &expected)

	suite.Equal(expected.Code, got.Code)
	suite.Equal(expected.Message, got.Message)
	suite.Equal(expected.Status, got.Status)

}

func (suite *test) do(method, path string, payload string) *http.Response {
	request := httptest.NewRequest(method, path, bytes.NewBufferString(payload))
	request.Header.Set("Content-Type", "application/json")

	response, err := suite.app.Test(request)
	suite.Require().NoError(err)

	return response
}

func (suite *test) TestTableAuthLogin() {
	table := []struct {
		name   string
		body   string
		setup  func()
		expect func(*http.Response)
	}{
		{
			name: "Success",
			body: `{
				"email": "john@example.com",
				"password": "secure123"
			}`,
			setup: func() {
				suite.service.EXPECT().
					Login(gomock.Any(), "john@example.com", "secure123").
					Times(1).
					Return(&ports_services.LoginOutput{
						Token:     "jwt-token",
						CreatedAt: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
					}, nil)
			},
			expect: func(response *http.Response) {
				suite.Equal(http.StatusOK, response.StatusCode)

				var got auth.ResponseLoginDTO
				suite.Require().NoError(json.NewDecoder(response.Body).Decode(&got))

				suite.Equal("jwt-token", got.Token)
				suite.Equal(time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC), got.CreatedAt)
			},
		},
		{
			name: "FailInternalServer",
			body: `{
				"email": "john@example.com",
				"password": "secure123"
			}`,
			setup: func() {
				suite.service.EXPECT().
					Login(gomock.Any(), "john@example.com", "secure123").
					Times(1).
					Return(nil, ports_errors.NewHttpInternalServer(nil))
			},
			expect: func(response *http.Response) {
				suite.Equal(http.StatusInternalServerError, response.StatusCode)
				suite.ExpectError(response.Body, ports_errors.NewHttpInternalServer(nil))
			},
		},
		{
			name: "FailInvalidCredentials",
			body: `{
				"email": "john@example.com",
				"password": "wrong-password"
			}`,
			setup: func() {
				suite.service.EXPECT().
					Login(gomock.Any(), "john@example.com", "wrong-password").
					Times(1).
					Return(nil, ports_errors.NewInvalidCredentials(nil))
			},
			expect: func(response *http.Response) {
				suite.Equal(http.StatusUnauthorized, response.StatusCode)
				suite.ExpectError(response.Body, ports_errors.NewInvalidCredentials(nil))
			},
		},
		{
			name: "FailInvalidBody",
			body: `{
				"email": "invalid-email",
				"password": "123"
			}`,
			setup: func() {
				suite.service.EXPECT().
					Login(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			expect: func(response *http.Response) {
				suite.Equal(http.StatusBadRequest, response.StatusCode)
				suite.ExpectError(response.Body, ports_errors.NewHttpBadRequest(nil))
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			tt.setup()

			response := suite.do(
				http.MethodPost,
				"/api/v1/auth/login",
				tt.body,
			)

			tt.expect(response)
		})
	}
}

func (suite *test) TestTableAuthRegister() {
	table := []struct {
		name   string
		body   string
		setup  func()
		expect func(*http.Response)
	}{
		{
			name: "Success",
			body: `{
				"username": "john",
				"email": "john@example.com",
				"password": "secure123"
			}`,
			setup: func() {
				suite.service.EXPECT().
					Register(gomock.Any(), "john", "john@example.com", "secure123").
					Times(1).
					Return(&ports_services.RegisterOutput{
						UserID:    "user-id",
						CreatedAt: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
					}, nil)
			},
			expect: func(response *http.Response) {
				suite.Equal(http.StatusCreated, response.StatusCode)

				var got auth.ResponseRegisterDTO
				suite.Require().NoError(json.NewDecoder(response.Body).Decode(&got))

				suite.Equal("user-id", got.ID)
				suite.Equal(time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC), got.CreatedAt)
			},
		},
		{
			name: "FailUserAlreadyExists",
			body: `{
				"username": "john",
				"email": "john@example.com",
				"password": "secure123"
			}`,
			setup: func() {
				suite.service.EXPECT().
					Register(gomock.Any(), "john", "john@example.com", "secure123").
					Times(1).
					Return(nil, ports_errors.NewUserAlreadyExists(nil))
			},
			expect: func(response *http.Response) {
				suite.Equal(http.StatusConflict, response.StatusCode)
				suite.ExpectError(response.Body, ports_errors.NewUserAlreadyExists(nil))
			},
		},
		{
			name: "FailInvalidBody",
			body: `{
				"username": "jo",
				"email": "invalid-email",
				"password": "123"
			}`,
			setup: func() {
				suite.service.EXPECT().
					Register(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			expect: func(response *http.Response) {
				suite.Equal(http.StatusBadRequest, response.StatusCode)
				suite.ExpectError(response.Body, ports_errors.NewHttpBadRequest(nil))
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			tt.setup()

			response := suite.do(
				http.MethodPost,
				"/api/v1/auth/register",
				tt.body,
			)

			tt.expect(response)
		})
	}
}
