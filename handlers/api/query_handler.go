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
	"recrem/utils"
	"recrem/utils/similarity"

	"github.com/gin-gonic/gin"
)

type QueryHandler struct {
}

func (q *QueryHandler) QueryByQuestion(ctx *gin.Context) {
	response := utils.Result{}
	queryInfoForm := forms.QueryInfoForm{}

	if err := ctx.ShouldBindJSON(&queryInfoForm); err != nil {
		response.Msg = fmt.Sprintf("binding info form: %s", err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	userID := queryInfoForm.UserID
	question := queryInfoForm.Question

	prompt := models.EmbeddingRequest{
		Input:          string(question),
		Model:          models.TEXTEMBEDDINGSMALL,
		EncodingFormat: "float",
	}

	resp, err := openai.O.CallEmbeddingAPI(&prompt)
	if err != nil {
		response.Msg = fmt.Sprintf("call embedding api error: %s", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// store the embedding
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Msg = fmt.Sprintf("Error reading response body: %s", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	// Unmarshal the JSON data into the struct
	var embeddingResp models.EmbeddingResponse
	err = json.Unmarshal(body, &embeddingResp)
	if err != nil {
		response.Msg = fmt.Sprintf("Error unmarshalling response: %s", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	embeddingFileArrs, err := etcd.GetVectorsThroughPrefix(userID)
	if err != nil {
		response.Msg = fmt.Sprintf("Error getting all user files embedding through userid: %s", err)
		fmt.Println("Error getting all user files embedding through userid:", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	mostSimilarIndex, maxSimilarity := similarity.FindMostSimilarVector(embeddingResp.Data[0].Embedding, embeddingFileArrs)
	response.Msg = fmt.Sprintf("Most similar vector index: %d, similarity: %f\n", mostSimilarIndex, maxSimilarity)
	ctx.JSON(http.StatusOK, response)
}

func (q *QueryHandler) QueryByFile(ctx *gin.Context) {
	response := utils.Result{}
	userID := ctx.Query("user_id")
	fmt.Println(userID)
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 10<<20)

	file, handler, err := ctx.Request.FormFile("filename")
	if err != nil {
		response.Msg = fmt.Sprintf("Error retrieving the file: %v", err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	defer file.Close()

	// Calculate the hash of the file content directly from the HTTP request
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		response.Msg = fmt.Sprintf("Error hashing the file: %v", err)
		ctx.JSON(http.StatusInternalServerError, response)
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
		response.Msg = fmt.Sprintf("Error creating the file: %v", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, file)
	if err != nil {
		response.Msg = fmt.Sprintf("Error saving the file: %v", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// run the script and get the file content
	pythonPath := "python3"
	scriptPath := "./script/extract.py"
	cmd := exec.Command(pythonPath, scriptPath, filePath) // ignore_security_alert RCE
	output, err := cmd.Output()
	if err != nil {
		response.Msg = fmt.Sprintf("Error executing script: %v", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// TODO summerize the content
	// vector embedding
	prompt := models.EmbeddingRequest{
		Input:          string(output),
		Model:          models.TEXTEMBEDDINGSMALL,
		EncodingFormat: "float",
	}
	resp, err := openai.O.CallEmbeddingAPI(&prompt)
	if err != nil {
		response.Msg = fmt.Sprintf("Call embedding API error: %v", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// store the embedding
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Msg = fmt.Sprintf("Error reading response body: %v", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	var embeddingResp models.EmbeddingResponse
	err = json.Unmarshal(body, &embeddingResp)
	if err != nil {
		response.Msg = fmt.Sprintf("Error unmarshalling response: %v", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	embeddingFileArrs, err := etcd.GetVectorsThroughPrefix(userID)
	if err != nil {
		response.Msg = fmt.Sprintf("Error getting all user files embedding through userID: %v", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	mostSimilarIndex, maxSimilarity := similarity.FindMostSimilarVector(embeddingResp.Data[0].Embedding, embeddingFileArrs)
	response.Msg = fmt.Sprintf("Most similar vector index: %d, similarity: %f\n", mostSimilarIndex, maxSimilarity)
	response.Data = map[string]interface{}{
		"mostSimilarIndex": mostSimilarIndex,
		"maxSimilarity":    maxSimilarity,
	}
	ctx.JSON(http.StatusOK, response)
}
