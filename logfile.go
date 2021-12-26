package logging

import (
	"errors"
	"os"
)

func init() {

}

func openFileSafe(filePath string) (*os.File, error) {
	if filePath == "" {
		return nil, errors.New("filePath  must noe be empty")
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			newFile, err := os.Create(filePath)
			if err != nil {
				return nil, err
			}
			return newFile, nil
		} else {
			return nil, err
		}
	}
	if fileInfo.IsDir() {
		return nil, errors.New(filePath + " is dir")
	}
	return os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
}
