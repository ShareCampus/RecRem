package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type FileHandler struct{}

func (f *FileHandler) UploadFile(ctx *gin.Context) {
	// 解析表单数据，限制最大文件大小为 10MB
	err := ctx.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("Error parsing form: %s", err.Error()))
		return
	}

	// 获取上传的文件
	file, handler, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("Error retrieving the file: %s", err.Error()))
		return
	}
	defer file.Close()

	// 指定文件保存的目录和路径
	dir := "./root/user"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	filePath := filepath.Join(dir, handler.Filename)

	// 创建一个新文件，用于保存上传的文件
	newFile, err := os.Create(filePath)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error creating the file: %s", err.Error()))
		return
	}
	defer newFile.Close()

	// 将上传的文件内容拷贝到新文件中
	_, err = io.Copy(newFile, file)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error copying file: %s", err.Error()))
		return
	}

	ctx.String(http.StatusOK, fmt.Sprintf("File uploaded successfully: %s", handler.Filename))
}

func (f *FileHandler) UploadFiles(ctx *gin.Context) {
	// 解析表单数据，限制最大文件大小为 10MB
	err := ctx.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("Error parsing form: %s", err.Error()))
		return
	}

	// 获取上传的文件列表
	files := ctx.Request.MultipartForm.File["files"]

	// 指定文件保存的目录和路径
	dir := "./root/user"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	// 遍历每个文件并保存
	for _, fileHeader := range files {
		// 打开上传的文件
		file, err := fileHeader.Open()
		if err != nil {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error opening file: %s", err.Error()))
			return
		}
		defer file.Close()

		filePath := filepath.Join(dir, fileHeader.Filename)

		// 创建一个新文件，用于保存上传的文件
		newFile, err := os.Create(filePath)
		if err != nil {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error creating the file: %s", err.Error()))
			return
		}
		defer newFile.Close()

		// 将上传的文件内容拷贝到新文件中
		_, err = io.Copy(newFile, file)
		if err != nil {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error copying file: %s", err.Error()))
			return
		}
	}

	ctx.String(http.StatusOK, fmt.Sprintf("Files uploaded successfully"))
}
