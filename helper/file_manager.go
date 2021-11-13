package helper

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveFile(location string, file *multipart.FileHeader) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fileLocation := filepath.Join(dir, location, file.Filename)
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
