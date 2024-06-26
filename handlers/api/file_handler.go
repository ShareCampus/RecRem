package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
