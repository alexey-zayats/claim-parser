package services

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"time"
)

// EventService ...
type EventService struct {
	claimSvc *ClaimService
	issSvc   *IssuedService
	fileSvc  *FileService
}

// EventServiceDI ...
type EventServiceDI struct {
	dig.In
	ClaimSvc *ClaimService
	IssSvc   *IssuedService
	FileSvc  *FileService
}

// NewEventService ...
func NewEventService(di EventServiceDI) *EventService {
	return &EventService{
		claimSvc: di.ClaimSvc,
		fileSvc:  di.FileSvc,
		issSvc:   di.IssSvc,
	}
}

// StoreClaim ...
func (s *EventService) StoreClaim(claim *model.Claim) {

	logrus.WithFields(logrus.Fields{
		"company":  claim.Company.Title,
		"district": claim.Event.District,
	}).Debug("Claim")

	event := claim.Event

	if err := s.claimSvc.SaveRecord(context.Background(), claim); err != nil {
		logrus.WithFields(logrus.Fields{"reason": err}).Error("unable save registry record")
		s.UpdateFile(event.FileID, 3, err.Error(), claim.Source)
	}

	s.UpdateFile(event.FileID, 0, "", claim.Source)
}

// StoreRegistry ...
func (s *EventService) StoreRegistry(record *model.Registry) {

	logrus.WithFields(logrus.Fields{
		"company": record.CompanyName,
		"car":     record.CompanyCar,
		"pass":    record.PassNumber,
	}).Debug("Registry")

	event := record.Event

	if err := s.issSvc.SaveRecord(context.Background(), record); err != nil {
		logrus.WithFields(logrus.Fields{"reason": err}).Error("unable save registry record")
		s.UpdateFile(event.FileID, 3, err.Error(), "")
	}

	s.UpdateFile(event.FileID, 0, "", "")
}

// UpdateFile ...
func (s *EventService) UpdateFile(id int64, status int, log string, source string) {

	f := &model.File{
		ID:        id,
		Status:    status,
		Log:       log,
		CreatedAt: time.Now(),
		Source:    source,
	}

	if err := s.fileSvc.UpdateState(f); err != nil {
		logrus.WithFields(logrus.Fields{"reason": err}).Error("unable update state")
	}
}
