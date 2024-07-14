package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"recrem/forms"
	"recrem/gpt/openai"
	"recrem/models"

	"github.com/gin-gonic/gin"
)

type QueryHandler struct {
}

func (q *QueryHandler) QueryByQuestion(ctx *gin.Context) {
	// 1. 获取问题 (user, question)
	queryInfoForm := forms.QueryInfoForm{}
	if err := ctx.ShouldBindJSON(&queryInfoForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// userID := queryInfoForm.UserID
	question := queryInfoForm.Question
	log.Println("Received UserID:", queryInfoForm.UserID)
	log.Println("Received Question:", queryInfoForm.Question)
	ctx.String(http.StatusOK, "Data received successfully")
	// 2. 向量化
	prompt := models.EmbeddingRequest{
		Input:          string(question),
		Model:          "text-embedding-3-small",
		EncodingFormat: "float",
	}
	resp, err := openai.O.CallEmbeddingAPI(&prompt)
	if err != nil {
		fmt.Println("call embedding api error", err)
	}
	// store the embedding
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	// Unmarshal the JSON data into the struct
	var embeddingResp models.EmbeddingResponse
	err = json.Unmarshal(body, &embeddingResp)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}
	fmt.Printf("Retrieved embedding response: %+v\n", embeddingResp.Data[0].Object)
	// 3. 相似度校验

	// 4. 返回结果
}
