package api

import (
	"encoding/json"
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

	"github.com/gin-gonic/gin"
)

const (
	FILESTORAGEPATH = "./storage/"
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

	filePath := filepath.Join(FILESTORAGEPATH, handler.Filename)

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
	scriptPath := "./utils/extract.py"

	cmd := exec.Command(pythonPath, scriptPath, filePath) // ignore_security_alert RCE
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing script:", err)
		return
	}
	fmt.Println("Output from Python script:")
	fmt.Println(string(output))

	// TODO summerize the content
	// vertor embedding
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

	// etcd.EtcdIns.Put(ctx, "test_key", string(body))
	// value, err := etcd.EtcdIns.Get(ctx, "test_key")
	// embeddingResp := models.EmbeddingResponse{}
	// err = json.Unmarshal(value, &embeddingResp)
	// if err != nil {
	// 	fmt.Println("Error unmarshalling response:", err)
	// 	return
	// }

	// Put the JSON data into etcd
	_, err = etcd.EtcdIns.Put(ctx, "test_key", string(body))
	if err != nil {
		log.Fatalf("Failed to put data into etcd: %v", err)
	}
	// Get the data back from etcd
	value, err := etcd.EtcdIns.Get(ctx, "test_key")
	if err != nil {
		log.Fatalf("Failed to get data from etcd: %v", err)
	}

	// Unmarshal the JSON data into the struct
	var embeddingResp models.EmbeddingResponse
	if len(value.Kvs) > 0 {
		err = json.Unmarshal(value.Kvs[0].Value, &embeddingResp)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			return
		}
		fmt.Printf("Retrieved embedding response: %+v\n", embeddingResp.Data[0].Object)
	} else {
		fmt.Println("No value found for key: test_key")
	}
	ctx.String(http.StatusOK, fmt.Sprintf("File uploaded successfully: %s", handler.Filename))
}

func (f *FileHandler) DeleteFile(ctx *gin.Context) {
	fileName := ctx.Query("filename")
	log.Println("delte filename", fileName)
	filePath := filepath.Join("./", fileName)
	err := os.Remove(filePath) // ignore_security_alert
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
