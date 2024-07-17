package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"recrem/config/etcd"
	"recrem/gpt/openai"
	"recrem/models"
	"recrem/utils"

	"github.com/gin-gonic/gin"
)

const (
	TEMPFILESTORAGEPATH = "./storage/temp/"
	FILESTORAGEPATH     = "./storage/"
)

type FileHandler struct{}

func (f *FileHandler) UploadFile(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 10<<20)

	response := utils.Result{}

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
	filePath := filepath.Join(FILESTORAGEPATH, handler.Filename)

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
		response.Msg = fmt.Sprintf("Error executing script: %s", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	//TODO: summerize the content

	// vector embedding
	prompt := models.EmbeddingRequest{
		Input:          string(output),
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

	// Put the JSON data into etcd
	fileEtcdKey := userID + fileHash
	_, err = etcd.EtcdIns.Put(ctx, fileEtcdKey, string(body))
	if err != nil {
		log.Fatalf("Failed to put data into etcd: %v", err)
		response.Msg = fmt.Sprintf("Failed to put data into etcd: %v", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response.Msg = fmt.Sprintf("File uploaded successfully: %s", handler.Filename)
	ctx.JSON(http.StatusOK, response)
}

func (f *FileHandler) DeleteFile(ctx *gin.Context) {
	fileName := ctx.Query("filename")
	response := utils.Result{}
	log.Println("delte filename", fileName)
	filePath := filepath.Join(FILESTORAGEPATH, fileName)
	err := os.Remove(filePath) // ignore_security_alert
	if err != nil {
		if os.IsNotExist(err) {
			response.Msg = fmt.Sprintf("File not found: %s", fileName)
			ctx.JSON(http.StatusInternalServerError, response)
		} else {
			response.Msg = fmt.Sprintf("Error deleting the file: %v", fileName)
			ctx.JSON(http.StatusInternalServerError, response)
		}
		return
	}
	response.Msg = fmt.Sprintf("File deleted successfully: %s", fileName)
	ctx.JSON(http.StatusOK, response)
}
