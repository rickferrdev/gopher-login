package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rickferrdev/gopher-login/internal/core/domain"
	ports_errors "github.com/rickferrdev/gopher-login/internal/core/ports/errors"
	"github.com/rickferrdev/gopher-login/internal/core/services/user/auth"
	ports_mocks "github.com/rickferrdev/gopher-login/internal/tests/mocks/ports"
	"github.com/stretchr/testify/suite"
)

type test struct {
	suite.Suite

	ctx     context.Context
	ctrl    *gomock.Controller
	storage *ports_mocks.MockUser
	token   *ports_mocks.MockTokenizer
	hasher  *ports_mocks.MockHasher
	service *auth.Service
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(test))
}

func (suite *test) SetupTest() {
	suite.ctx = context.Background()
	suite.ctrl = gomock.NewController(suite.T())

	suite.storage = ports_mocks.NewMockUser(suite.ctrl)
	suite.token = ports_mocks.NewMockTokenizer(suite.ctrl)
	suite.hasher = ports_mocks.NewMockHasher(suite.ctrl)

	service, err := auth.New(auth.Params{
		UserStorage:    suite.storage,
		JwtPlatform:    suite.token,
		HasherPlatform: suite.hasher,
	})

	suite.Require().NoError(err)

	authService, ok := service.(*auth.Service)
	suite.Require().True(ok)

	suite.service = authService
}

func (suite *test) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *test) RequireErrorCode(err error, code ports_errors.Code) {
	suite.Require().Error(err)
	suite.True(
		ports_errors.IsCode(err, code),
		"expected error code %s, got %v",
		code,
		err,
	)
}

func (suite *test) TestTableLogin() {
	user := &domain.User{
		ID:        "user-id",
		Username:  "john",
		Email:     "john@example.com",
		Password:  "hashed-password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	table := []struct {
		name   string
		email  string
		pass   string
		setup  func()
		expect func(outputToken string, err error)
	}{
		{
			name:  "Success",
			email: " John@Example.com ",
			pass:  "secure123",
			setup: func() {
				suite.storage.EXPECT().
					FindByEmail(gomock.Any(), "john@example.com").
					Times(1).
					Return(user, nil)

				suite.hasher.EXPECT().
					Validate([]byte("hashed-password"), []byte("secure123")).
					Times(1).
					Return(nil)

				suite.token.EXPECT().
					GenerateUserToken("user-id").
					Times(1).
					Return("jwt-token", nil)
			},
			expect: func(outputToken string, err error) {
				suite.Require().NoError(err)
				suite.Equal("jwt-token", outputToken)
			},
		},
		{
			name:  "FailEmptyEmail",
			email: "",
			pass:  "secure123",
			setup: func() {},
			expect: func(outputToken string, err error) {
				suite.Empty(outputToken)
				suite.RequireErrorCode(err, ports_errors.CodeHttpBadRequest)
			},
		},
		{
			name:  "FailEmptyPassword",
			email: "john@example.com",
			pass:  "",
			setup: func() {},
			expect: func(outputToken string, err error) {
				suite.Empty(outputToken)
				suite.RequireErrorCode(err, ports_errors.CodeHttpBadRequest)
			},
		},
		{
			name:  "FailUserNotFound",
			email: "john@example.com",
			pass:  "secure123",
			setup: func() {
				suite.storage.EXPECT().
					FindByEmail(gomock.Any(), "john@example.com").
					Times(1).
					Return(nil, ports_errors.NewDatabaseNotFound(nil))
			},
			expect: func(outputToken string, err error) {
				suite.Empty(outputToken)
				suite.RequireErrorCode(err, ports_errors.CodeInvalidCredentials)
			},
		},
		{
			name:  "FailStorageInternal",
			email: "john@example.com",
			pass:  "secure123",
			setup: func() {
				suite.storage.EXPECT().
					FindByEmail(gomock.Any(), "john@example.com").
					Times(1).
					Return(nil, ports_errors.NewDatabaseInternal(errors.New("database failed")))
			},
			expect: func(outputToken string, err error) {
				suite.Empty(outputToken)
				suite.RequireErrorCode(err, ports_errors.CodeHttpInternalServer)
			},
		},
		{
			name:  "FailInvalidPassword",
			email: "john@example.com",
			pass:  "wrong-password",
			setup: func() {
				suite.storage.EXPECT().
					FindByEmail(gomock.Any(), "john@example.com").
					Times(1).
					Return(user, nil)

				suite.hasher.EXPECT().
					Validate([]byte("hashed-password"), []byte("wrong-password")).
					Times(1).
					Return(ports_errors.NewHasherValidateFailed(nil))
			},
			expect: func(outputToken string, err error) {
				suite.Empty(outputToken)
				suite.RequireErrorCode(err, ports_errors.CodeInvalidCredentials)
			},
		},
		{
			name:  "FailTokenGeneration",
			email: "john@example.com",
			pass:  "secure123",
			setup: func() {
				suite.storage.EXPECT().
					FindByEmail(gomock.Any(), "john@example.com").
					Times(1).
					Return(user, nil)

				suite.hasher.EXPECT().
					Validate([]byte("hashed-password"), []byte("secure123")).
					Times(1).
					Return(nil)

				suite.token.EXPECT().
					GenerateUserToken("user-id").
					Times(1).
					Return("", ports_errors.NewJwtGenerateFailed(nil))
			},
			expect: func(outputToken string, err error) {
				suite.Empty(outputToken)
				suite.RequireErrorCode(err, ports_errors.CodeHttpInternalServer)
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			tt.setup()

			output, err := suite.service.Login(suite.ctx, tt.email, tt.pass)

			if output == nil {
				tt.expect("", err)
				return
			}

			tt.expect(output.Token, err)
			suite.False(output.CreatedAt.IsZero())
		})
	}
}

func (suite *test) TestTableRegister() {
	table := []struct {
		name     string
		username string
		email    string
		pass     string
		setup    func()
		expect   func(outputUserID string, err error)
	}{
		{
			name:     "Success",
			username: " john ",
			email:    " John@Example.com ",
			pass:     "secure123",
			setup: func() {
				suite.hasher.EXPECT().
					Generate([]byte("secure123")).
					Times(1).
					Return([]byte("hashed-password"), nil)

				suite.storage.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Times(1).
					DoAndReturn(func(ctx context.Context, user domain.User) error {
						suite.NotEmpty(user.ID)
						suite.Equal("john", user.Username)
						suite.Equal("john@example.com", user.Email)
						suite.Equal("hashed-password", user.Password)
						suite.False(user.CreatedAt.IsZero())
						suite.False(user.UpdatedAt.IsZero())

						return nil
					})
			},
			expect: func(outputUserID string, err error) {
				suite.Require().NoError(err)
				suite.NotEmpty(outputUserID)
			},
		},
		{
			name:     "FailEmptyUsername",
			username: "",
			email:    "john@example.com",
			pass:     "secure123",
			setup:    func() {},
			expect: func(outputUserID string, err error) {
				suite.Empty(outputUserID)
				suite.RequireErrorCode(err, ports_errors.CodeHttpBadRequest)
			},
		},
		{
			name:     "FailEmptyEmail",
			username: "john",
			email:    "",
			pass:     "secure123",
			setup:    func() {},
			expect: func(outputUserID string, err error) {
				suite.Empty(outputUserID)
				suite.RequireErrorCode(err, ports_errors.CodeHttpBadRequest)
			},
		},
		{
			name:     "FailEmptyPassword",
			username: "john",
			email:    "john@example.com",
			pass:     "",
			setup:    func() {},
			expect: func(outputUserID string, err error) {
				suite.Empty(outputUserID)
				suite.RequireErrorCode(err, ports_errors.CodeHttpBadRequest)
			},
		},
		{
			name:     "FailHasherGenerate",
			username: "john",
			email:    "john@example.com",
			pass:     "secure123",
			setup: func() {
				suite.hasher.EXPECT().
					Generate([]byte("secure123")).
					Times(1).
					Return(nil, ports_errors.NewHasherGenerateFailed(nil))
			},
			expect: func(outputUserID string, err error) {
				suite.Empty(outputUserID)
				suite.RequireErrorCode(err, ports_errors.CodeHttpInternalServer)
			},
		},
		{
			name:     "FailUserAlreadyExists",
			username: "john",
			email:    "john@example.com",
			pass:     "secure123",
			setup: func() {
				suite.hasher.EXPECT().
					Generate([]byte("secure123")).
					Times(1).
					Return([]byte("hashed-password"), nil)

				suite.storage.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Times(1).
					Return(ports_errors.NewDatabaseDuplicateKey(nil))
			},
			expect: func(outputUserID string, err error) {
				suite.Empty(outputUserID)
				suite.RequireErrorCode(err, ports_errors.CodeUserAlreadyExists)
			},
		},
		{
			name:     "FailStorageInternal",
			username: "john",
			email:    "john@example.com",
			pass:     "secure123",
			setup: func() {
				suite.hasher.EXPECT().
					Generate([]byte("secure123")).
					Times(1).
					Return([]byte("hashed-password"), nil)

				suite.storage.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Times(1).
					Return(ports_errors.NewDatabaseInternal(errors.New("insert failed")))
			},
			expect: func(outputUserID string, err error) {
				suite.Empty(outputUserID)
				suite.RequireErrorCode(err, ports_errors.CodeHttpInternalServer)
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			tt.setup()

			output, err := suite.service.Register(
				suite.ctx,
				tt.username,
				tt.email,
				tt.pass,
			)

			if output == nil {
				tt.expect("", err)
				return
			}

			tt.expect(output.UserID, err)
			suite.False(output.CreatedAt.IsZero())
		})
	}
}
