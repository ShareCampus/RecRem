package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"recrem/config/etcd"
	"recrem/forms"
	"recrem/gpt/openai"
	"recrem/models"

	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
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

	// 通过userid获取所有的文件向量
	userFilesEmbedding, err := etcd.EtcdIns.Get(context.Background(), userID, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("Error getting all user files embedding:", err)
		return
	}

	for _, kv := range userFilesEmbedding.Kvs {
		fmt.Printf("the key is %s\n", kv.Key)
	}
	// 计算和问题最相近的向量
	embeddingArrs := make([][]float64, 0)
	for _, kv := range userFilesEmbedding.Kvs {
		embeddingResp := models.EmbeddingResponse{}
		err = json.Unmarshal(kv.Value, &embeddingResp)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			return
		}
		embeddingArrs = append(embeddingArrs, embeddingResp.Data[0].Embedding)
	}

	mostSimilarIndex, maxSimilarity := findMostSimilarVector(embeddingResp.Data[0].Embedding, embeddingArrs)
	fmt.Printf("Most similar vector index: %d, similarity: %f\n", mostSimilarIndex, maxSimilarity)
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
	fmt.Println("Output from Python script:")
	fmt.Println(string(output))

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

	// Unmarshal the JSON data into the struct
	var embeddingResp models.EmbeddingResponse
	err = json.Unmarshal(body, &embeddingResp)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}
	fmt.Printf("Retrieved embedding response: %+v\n", embeddingResp.Data[0].Object)

	// 通过userid获取所有的文件向量
	userFilesEmbedding, err := etcd.EtcdIns.Get(context.Background(), userID, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("Error getting all user files embedding:", err)
		return
	}

	for _, kv := range userFilesEmbedding.Kvs {
		fmt.Printf("the key is %s\n", kv.Key)
	}
	// 计算和问题最相近的向量
	embeddingArrs := make([][]float64, 0)
	for _, kv := range userFilesEmbedding.Kvs {
		embeddingResp := models.EmbeddingResponse{}
		err = json.Unmarshal(kv.Value, &embeddingResp)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			return
		}
		embeddingArrs = append(embeddingArrs, embeddingResp.Data[0].Embedding)
	}

	mostSimilarIndex, maxSimilarity := findMostSimilarVector(embeddingResp.Data[0].Embedding, embeddingArrs)
	fmt.Printf("Most similar vector index: %d, similarity: %f\n", mostSimilarIndex, maxSimilarity)

}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		panic("Vectors must be the same length")
	}
	var dotProduct, magnitudeA, magnitudeB float64
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		magnitudeA += a[i] * a[i]
		magnitudeB += b[i] * b[i]
	}
	magnitudeA = math.Sqrt(magnitudeA)
	magnitudeB = math.Sqrt(magnitudeB)
	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}
	return dotProduct / (magnitudeA * magnitudeB)
}

func findMostSimilarVector(questionVector []float64, embeddingArrs [][]float64) (int, float64) {
	var maxSimilarity float64
	var mostSimilarIndex int
	for i, vector := range embeddingArrs {
		similarity := cosineSimilarity(questionVector, vector)
		if similarity > maxSimilarity || i == 0 {
			maxSimilarity = similarity
			mostSimilarIndex = i
		}
	}
	return mostSimilarIndex, maxSimilarity
}
