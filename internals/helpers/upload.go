package helpers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// UserRegisterFileUpload uploads single file
func MediaPicUpload(r *http.Request, field, username string) (string, error) {
	file, handler, err := r.FormFile(field)
	if err != nil {
		return "", fmt.Errorf("error in getting file: %s", err)
	}
	defer file.Close()
	fileName := handler.Filename
	if fileName == "" {
		return "", errors.New("no file found")
	}
	ext := filepath.Ext(fileName)

	path := filepath.Join("./media/user", username)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", fmt.Errorf("error in creating directory: %s", err)
	}

	newFileName := fmt.Sprintf("%s-%s%s", username, field, ext)
	relativePath := filepath.Join(path, newFileName)

	new, err := os.Create(relativePath)
	if err != nil {
		return "", fmt.Errorf("error in creating file: %s", err)
	}
	defer new.Close()
	_, err = io.Copy(new, file)
	if err != nil {
		return "", fmt.Errorf("error in copying file: %s", err)
	}
	return relativePath, nil
}

// AdminPublicUploadImage uploads single file
func AdminPublicUploadImage(r *http.Request, field, table string, id string) (string, error) {
	file, handler, err := r.FormFile(field)
	if err != nil {
		return "", fmt.Errorf("error in getting file: %s", err)
	}
	defer file.Close()
	fileName := handler.Filename
	if fileName == "" {
		return "", errors.New("no file found")
	}
	ext := filepath.Ext(fileName)

	path := filepath.Join("./public", table)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", fmt.Errorf("error in creating directory: %s", err)
	}

	newFileName := fmt.Sprintf("%s-%s-%s%s", field, table, id, ext)
	relativePath := filepath.Join(path, newFileName)

	new, err := os.Create(relativePath)
	if err != nil {
		return "", fmt.Errorf("error in creating file: %s", err)
	}
	defer new.Close()
	_, err = io.Copy(new, file)
	if err != nil {
		return "", fmt.Errorf("error in copying file: %s", err)
	}
	return relativePath, nil
}

// AdminPublicUploadImage uploads single file
func AdminPublicUploadImage2(r *http.Request, field, table string, id int64) (string, error) {
	file, handler, err := r.FormFile(field)
	if err != nil {
		return "", fmt.Errorf("error in getting file: %s", err)
	}
	defer file.Close()
	fileName := handler.Filename
	if fileName == "" {
		return "", errors.New("no file found")
	}
	ext := filepath.Ext(fileName)

	path := filepath.Join("./public", table)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", fmt.Errorf("error in creating directory: %s", err)
	}

	newFileName := fmt.Sprintf("%s-%s-%d%s", field, table, id, ext)
	relativePath := filepath.Join(path, newFileName)

	new, err := os.Create(relativePath)
	if err != nil {
		return "", fmt.Errorf("error in creating file: %s", err)
	}
	defer new.Close()
	_, err = io.Copy(new, file)
	if err != nil {
		return "", fmt.Errorf("error in copying file: %s", err)
	}
	return relativePath, nil
}
