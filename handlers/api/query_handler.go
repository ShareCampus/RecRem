package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"recrem/config/etcd"
	"recrem/forms"
	"recrem/gpt/openai"
	"recrem/models"
	utils "recrem/utils/similarity"

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
	userID := queryInfoForm.UserID
	question := queryInfoForm.Question
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

	embeddingFileArrs, err := etcd.GetVectorsThroughPrefix(userID)
	if err != nil {
		fmt.Println("Error getting all user files embedding through userid:", err)
		return
	}

	mostSimilarIndex, maxSimilarity := utils.FindMostSimilarVector(embeddingResp.Data[0].Embedding, embeddingFileArrs)
	fmt.Printf("Most similar vector index: %d, similarity: %f\n", mostSimilarIndex, maxSimilarity)
	return
}

func (q *QueryHandler) QueryByFile(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	fmt.Println(userID)
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 10<<20)

	file, handler, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("Error retrieving the file: %v", err))
		return
	}
	defer file.Close()

	// Calculate the hash of the file content directly from the HTTP request
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error hashing the file: %v", err))
		return
	}
	fileHash := hex.EncodeToString(hasher.Sum(nil))

	// Use the hash as the file name
	handler.Filename = userID + fileHash + ".md"
	filePath := filepath.Join(TEMPFILESTORAGEPATH, handler.Filename)

	// Rewind the file reader since io.Copy above has exhausted it
	file.Seek(0, 0)

	destFile, err := os.Create(filePath)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error creating the file: %v", err))
		return
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, file)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error saving the file: %v", err))
		return
	}

	// run the script and get the file content
	pythonPath := "python3"
	scriptPath := "./script/extract.py"
	cmd := exec.Command(pythonPath, scriptPath, filePath) // ignore_security_alert RCE
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing script:", err)
		return
	}

	// TODO summerize the content
	// vector embedding
	prompt := models.EmbeddingRequest{
		Input:          string(output),
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

	var embeddingResp models.EmbeddingResponse
	err = json.Unmarshal(body, &embeddingResp)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	embeddingFileArrs, err := etcd.GetVectorsThroughPrefix(userID)
	if err != nil {
		fmt.Println("Error getting all user files embedding through userid:", err)
		return
	}

	mostSimilarIndex, maxSimilarity := utils.FindMostSimilarVector(embeddingResp.Data[0].Embedding, embeddingFileArrs)
	fmt.Printf("Most similar vector index: %d, similarity: %f\n", mostSimilarIndex, maxSimilarity)
	return
}
