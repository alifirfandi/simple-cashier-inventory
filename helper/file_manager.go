package helper

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveFile(file *multipart.FileHeader, location, filename string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fileLocation := filepath.Join(dir, location, filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	fileData, err := file.Open()
	if err != nil {
		return err
	}
	defer fileData.Close()

	if _, err := io.Copy(targetFile, fileData); err != nil {
		return err
	}

	return nil
}
