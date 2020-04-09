package services

import (
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"go.uber.org/dig"
)

// FileService ...
type FileService struct {
	repo interfaces.FileRepository
}

// FileServiceDI ...
type FileServiceDI struct {
	dig.In
	Repo interfaces.FileRepository
}

// NewFileService ...
func NewFileService(di FileServiceDI) *FileService {
	return &FileService{
		repo: di.Repo,
	}
}

// Create ...
func (s *FileService) Create(data *model.File) error {
	return s.repo.Create(data)
}

// UpdateState ...
func (s *FileService) UpdateState(data *model.File) error {
	return s.repo.UpdateState(data)
}

// Delete ...
func (s *FileService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Read ...
func (s *FileService) Read(id int64) (*model.File, error) {
	return s.repo.Read(id)
}
