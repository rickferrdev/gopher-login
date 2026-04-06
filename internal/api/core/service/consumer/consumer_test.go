package consumer

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/rickferrdev/gopher-login/internal/api/core/domain"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type test struct {
	suite.Suite
	ctrl *gomock.Controller
	repo *MockRepository
}

func (suite *test) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.repo = NewMockRepository(suite.ctrl)
}

func (suite *test) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *test) TestFindByUsername() {
	params := ServiceParams{
		Repository: suite.repo,
	}

	table := []struct {
		name  string
		input string

		mockErr    error
		mockOutput *domain.Consumer

		expectErr    error
		expectOutput *ports.ConsumerPayload
	}{
		{
			name:         "Success",
			input:        "rickferrdev",
			mockErr:      nil,
			mockOutput:   &domain.Consumer{Username: "rickferrdev"},
			expectErr:    nil,
			expectOutput: &ports.ConsumerPayload{Username: "rickferrdev"},
		},
		{
			name:         "Failed",
			input:        "rickferrdev",
			mockErr:      ports.NewError(ports.CodeUserNotFound, ports.MessageNotFound, 404, nil),
			mockOutput:   nil,
			expectErr:    &ports.GopherError{Code: ports.CodeUserNotFound},
			expectOutput: nil,
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			suite.repo.EXPECT().FindByUsername(gomock.Any(), tt.input).Times(1).Return(tt.mockOutput, tt.mockErr)

			service := New(params)

			obtained, err := service.FindByUsername(context.Background(), tt.input)
			if tt.expectErr != nil {
				suite.ErrorIs(err, tt.expectErr)
			}

			suite.Equal(obtained, tt.expectOutput)
		})
	}
}

func (suite *test) TestFindByID() {
	params := ServiceParams{
		Repository: suite.repo,
	}
	table := []struct {
		name         string
		input        string
		mockErr      error
		mockOutput   *domain.Consumer
		expectErr    error
		expectOutput *ports.ConsumerPayload
	}{
		{
			name:         "Success",
			input:        uuid.NewString(),
			mockOutput:   &domain.Consumer{Username: "rickferrdev"},
			mockErr:      nil,
			expectErr:    nil,
			expectOutput: &ports.ConsumerPayload{Username: "rickferrdev"},
		},
		{
			name:         "Failed",
			input:        uuid.NewString(),
			mockOutput:   nil,
			mockErr:      ports.NewError(ports.CodeUserNotFound, ports.MessageNotFound, 404, nil),
			expectErr:    &ports.GopherError{Code: ports.CodeUserNotFound},
			expectOutput: nil,
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			suite.repo.EXPECT().FindByID(gomock.Any(), tt.input).Times(1).Return(tt.mockOutput, tt.mockErr)

			service := New(params)

			obtained, err := service.FindByID(context.Background(), tt.input)
			if tt.expectErr != nil {
				suite.ErrorIs(err, tt.expectErr)
			}

			suite.Equal(obtained, tt.expectOutput)
		})
	}
}

func TestConsumerService(t *testing.T) {
	suite.Run(t, new(test))
}
