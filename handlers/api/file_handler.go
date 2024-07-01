package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"recrem/gpt/openai"
	"recrem/models"

	"github.com/gin-gonic/gin"
)

type FileHandler struct{}

func (f *FileHandler) UploadFile(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 10<<20)

	file, handler, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("Error retrieving the file: %v", err))
		return
	}
	defer file.Close()

	filePath := filepath.Join("./", handler.Filename)

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

	// 执行脚本获取文本中的内容
	pythonPath := "python3"
	scriptPath := "./utils/extract.py"

	cmd := exec.Command(pythonPath, scriptPath, filePath)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing script:", err)
		return
	}
	fmt.Println("Output from Python script:")
	fmt.Println(string(output))

	// 将内容进行总结

	// 将内容进行向量化操作
	prompt := models.EmbeddingRequest{
		Input:          string(output),
		Model:          "text-embedding-3-small",
		EncodingFormat: "float",
	}
	resp, err := openai.O.CallEmbeddingAPI(&prompt)
	if err != nil {
		fmt.Println("call embedding api error", err)
	}
	// 存储到数据库中
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response body:", string(body))
	ctx.String(http.StatusOK, fmt.Sprintf("File uploaded successfully: %s", handler.Filename))
}

func (f *FileHandler) DeleteFile(ctx *gin.Context) {
	fileName := ctx.Query("filename")
	log.Println("delte filename", fileName)
	filePath := filepath.Join("./", fileName)
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.String(http.StatusNotFound, fmt.Sprintf("File not found: %s", fileName))
		} else {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error deleting the file: %v", err))
		}
		return
	}

	ctx.String(http.StatusOK, fmt.Sprintf("File deleted successfully: %s", fileName))
}
