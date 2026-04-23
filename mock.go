package gemini

import (
	"context"
	"io"
	"os"
	"time"
)

type MockClient struct {
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (c *MockClient) Generate(_ context.Context, _, _, _ string) ([]byte, error) {
	time.Sleep(10 * time.Second)
	return os.ReadFile("internal/gemini/mockData/img.png")
}

func (c *MockClient) Edit(_ context.Context, _, _, _ string, _ io.Reader) ([]byte, error) {
	time.Sleep(7 * time.Second)
	return os.ReadFile("internal/gemini/mockData/img_edited.png")
}
