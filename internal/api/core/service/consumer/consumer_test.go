package consumer

import (
	"context"
	"log/slog"
	"testing"

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
	table := []struct {
		name         string
		input        string
		mockReturn   *domain.Consumer
		mockErr      error
		expectErr    error
		expectOutput *ports.ConsumerPayload
	}{
		{
			name:         "Success",
			input:        "rickferrdev",
			mockReturn:   &domain.Consumer{Username: "rickferrdev"},
			mockErr:      nil,
			expectErr:    nil,
			expectOutput: &ports.ConsumerPayload{Username: "rickferrdev"},
		},
		{
			name:         "Failed",
			input:        "rickferrdev",
			mockReturn:   nil,
			mockErr:      nil,
			expectErr:    ports.ErrConsumerNotFound,
			expectOutput: nil,
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			suite.repo.EXPECT().FindByUsername(gomock.Any(), tt.input).Times(1).Return(tt.mockReturn, tt.mockErr)

			service := New(suite.repo, slog.New(slog.DiscardHandler))

			obtained, err := service.FindByUsername(context.Background(), tt.input)
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
