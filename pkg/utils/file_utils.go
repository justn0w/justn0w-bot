package utils

import (
	"io"
	"mime/multipart"
)

func GetFileContent(uploadFile *multipart.FileHeader) (string, error) {
	fileContent, err := uploadFile.Open()
	if err != nil {
		return "", err
	}
	defer fileContent.Close()

	content, err := io.ReadAll(fileContent)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
