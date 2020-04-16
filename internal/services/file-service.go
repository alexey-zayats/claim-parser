package services

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
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
func (s *FileService) Create(data *entity.File) error {
	return s.repo.Create(data)
}

// UpdateState ...
func (s *FileService) UpdateState(data *entity.File) error {
	return s.repo.UpdateState(data)
}

// Delete ...
func (s *FileService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Read ...
func (s *FileService) Read(id int64) (*entity.File, error) {
	return s.repo.Read(id)
}
