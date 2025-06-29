package services

import (
	"fmt"
	"os"
	"path/filepath"
)

func DownloadFile(fileName string) (string, error) {
	const uploadDir = "uploads"

	filePath := filepath.Join(uploadDir, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", filePath)
	}

	return filePath, nil
}
