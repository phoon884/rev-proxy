package repositories

import (
	"os"
)

type FileRepository struct {
	path string
}

func NewFileRepository(path string) *FileRepository {
	return &FileRepository{
		path: path,
	}
}

func (j *FileRepository) GetConfig() (*[]byte, error) {
	config, err := os.ReadFile(j.path)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
