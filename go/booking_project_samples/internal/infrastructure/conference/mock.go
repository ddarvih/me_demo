package conference

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type MockService struct {
	Err error
}

func NewMock() *MockService {
	return &MockService{}
}

func (m *MockService) CreateLink(_ context.Context, bookingID uuid.UUID) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}
	return fmt.Sprintf("https://pupupuCalls.internal/room-%s", bookingID.String()[:8]), nil
}
