package services

import (
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"strings"
	"time"
)

// EventService ...
type EventService struct {
	vcs     *VehicleClaimService
	vis     *VehicleIssuedService
	pcs     *PeopleClaimService
	fileSvc *FileService
}

// EventServiceDI ...
type EventServiceDI struct {
	dig.In
	VCS     *VehicleClaimService
	VIS     *VehicleIssuedService
	PCS     *PeopleClaimService
	FileSvc *FileService
}

// NewEventService ...
func NewEventService(di EventServiceDI) *EventService {
	return &EventService{
		vcs:     di.VCS,
		vis:     di.VIS,
		fileSvc: di.FileSvc,
		pcs:     di.PCS,
	}
}

// Store ...
func (s *EventService) Store(out *model.Out) {
	switch out.Kind {
	case model.OutVehicleClaim:

		claim := out.Value.(*model.VehicleClaim)
		s.vehicleClaim(out.Event, claim)

	case model.OutVehicleRegistry:

		record := out.Value.(*model.VehicleRegistry)
		s.vehicleRegistry(out.Event, record)

	case model.OutPeopleClaim:

		record := out.Value.(*model.PeopleClaim)
		s.peopleClaim(out.Event, record)

	}
}

func (s *EventService) peopleClaim(event *model.Event, claim *model.PeopleClaim) {

	logrus.WithFields(logrus.Fields{
		"company":  claim.Company.Title,
		"district": event.DistrictID,
		"phone":    claim.Company.HeadPhone,
		"email":    claim.Company.HeadEmail,
	}).Debug("People.Claim")

	rec := fmt.Sprintf("%s;%d;%s", claim.Created, claim.Company.TIN, claim.Company.Title)

	if claim.Success {
		if err := s.pcs.SaveRecord(event, claim); err != nil {
			logrus.WithFields(logrus.Fields{"reason": err}).Error("unable save people claim record")
			log := rec + ";sql: " + err.Error()
			s.UpdateFile(event.FileID, 3, log, claim.Source)
			return
		}
	} else {
		log := rec + ";parse: " + strings.Join(claim.Reason, ", ") + "\n"
		s.UpdateFile(event.FileID, 1, log, claim.Source)
	}
}

func (s *EventService) vehicleClaim(event *model.Event, claim *model.VehicleClaim) {

	logrus.WithFields(logrus.Fields{
		"company":  claim.Company.Title,
		"district": event.DistrictID,
	}).Debug("Vehicle.Claim")

	rec := fmt.Sprintf("%s;%d;%s", claim.Created, claim.Company.TIN, claim.Company.Title)

	if claim.Success {
		if err := s.vcs.SaveRecord(event, claim); err != nil {
			logrus.WithFields(logrus.Fields{"reason": err}).Error("unable save people claim record")
			log := rec + ";sql: " + err.Error()
			s.UpdateFile(event.FileID, 3, log, claim.Source)
			return
		}
	} else {
		log := rec + ";parse: " + strings.Join(claim.Reason, ", ") + "\n"
		s.UpdateFile(event.FileID, 1, log, "")
	}
}

// StoreRegistry ...
func (s *EventService) vehicleRegistry(event *model.Event, record *model.VehicleRegistry) {

	rec := fmt.Sprintf("%s;%d;%s", record.IssuedAt, record.CompanyInn, record.CompanyName)

	if record.Success {

		logrus.WithFields(logrus.Fields{
			"company":  record.CompanyName,
			"district": event.DistrictID,
			"car":      record.CompanyCar,
			"pass":     record.PassNumber,
		}).Debug("Vehicle.Registry")

		if err := s.vis.SaveRecord(event, record); err != nil {
			logrus.WithFields(logrus.Fields{"reason": err}).Error("unable save registry record")
			log := rec + ";sql: " + err.Error()
			s.UpdateFile(event.FileID, 3, log, "")
			return
		}
	} else {
		log := rec + ";parse: " + strings.Join(record.Reason, ", ") + "\n"
		s.UpdateFile(event.FileID, 3, log, "")
	}
}

// UpdateFile ...
func (s *EventService) UpdateFile(id int64, status int, log string, source string) {

	f := &entity.File{
		ID:        id,
		Status:    status,
		Log:       log + "\n",
		CreatedAt: time.Now(),
		Source:    source,
	}

	if err := s.fileSvc.UpdateState(f); err != nil {
		logrus.WithFields(logrus.Fields{"reason": err}).Error("unable update file state")
	}
}
