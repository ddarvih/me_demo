package conference

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestMockService_CreateLink_Err(t *testing.T) {
	m := &MockService{Err: errors.New("down")}
	_, err := m.CreateLink(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMockService_CreateLink_OK(t *testing.T) {
	m := NewMock()
	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	link, err := m.CreateLink(context.Background(), id)
	if err != nil || link == "" {
		t.Fatalf("link=%q err=%v", link, err)
	}
}
