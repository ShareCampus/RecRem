package openai

import (
	"fmt"
	"io"
	"net/http"
	"recrem/mock"
	"recrem/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallEmbeddingAPI(t *testing.T) {
	mockOpenAI := new(mock.MockOpenAI)
	prompt := &models.EmbeddingRequest{
		Input:          "The food was delicious and the waiter...",
		Model:          "text-embedding-3-small",
		EncodingFormat: "float",
	}

	// 设置 GetToken 方法的 mock 返回值
	mockOpenAI.On("GetToken").Return("**", nil)

	// 创建 OpenAI 实例并调用 CallEmbeddingAPI 方法
	openai := &OpenAI{}
	resp, err := openai.CallEmbeddingAPI(prompt)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	fmt.Println("Response Body:", string(body))

	mockOpenAI.AssertExpectations(t)
}
