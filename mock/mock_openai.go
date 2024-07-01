package mock

import (
	"net/http"
	"recrem/models"

	"github.com/stretchr/testify/mock"
)

type MockOpenAI struct {
	mock.Mock
}

func (m *MockOpenAI) GetToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockOpenAI) CallEmbeddingAPI(prompt *models.EmbeddingRequest) (*http.Response, error) {
	args := m.Called(prompt)
	return args.Get(0).(*http.Response), args.Error(1)
}
