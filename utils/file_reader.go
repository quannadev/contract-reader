package utils

import (
	"io"
	"os"
	"path/filepath"
)

func ReadFile(path string) ([]byte, error) {
	//get the root directory of the project
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	gp := filepath.Dir(dir)

	filePath, err := os.Open(gp + "/" + path)
	if err != nil {
		return nil, err
	}
	defer filePath.Close()
	fileContent, err := io.ReadAll(filePath)
	if err != nil {
		return nil, err
	}
	return fileContent, nil

}
