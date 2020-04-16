package fs

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

// State ...
type State int

const (
	// StateInit ...
	StateInit State = iota
	// StateKind ...
	StateKind
	// StateName ...
	StateName
	// StateAddress ...
	StateAddress
	// StateINN ...
	StateINN
	// StateFIO ...
	StateFIO
	// StatePhone ...
	StatePhone
	// StateEMail ...
	StateEMail
	// StateCars ...
	StateCars
	// StateAgreement ...
	StateAgreement
	// StateReliability ...
	StateReliability
)

// Parser ...
type Parser struct {
	headers map[string]State
	path    string
	event   *model.Event
	name    string
}

// Register ...
func Register() {
	parser.Instance().Add("form.struct", NewParser)
}

// NewParser ...
func NewParser(name string) (parser.Backend, error) {

	headers := map[string]State{
		"Вид деятельности": StateKind,
		"Полное название организации, индивидуального предпринимателя": StateName,
		"Адрес, местонахождение": StateAddress,
		"ИНН":              StateINN,
		"ФИО руководителя": StateFIO,
		"Контактный телефон руководителя":                                     StatePhone,
		"Электронный адрес руководителя":                                      StateEMail,
		"Перечень автомобилей, для которых нужны пропуска на время карантина": StateCars,
		"Согласие на обработку персональных данных":                           StateAgreement,
		"Достоверность предоставляемых сведений":                              StateReliability,
	}

	return &Parser{
		headers: headers,
		name:    name,
	}, nil
}

// WithEvent ...
func (p *Parser) WithEvent(event *model.Event) {
	p.event = event
	p.path = event.Filepath
}

// Parse ...
func (p *Parser) Parse(ctx context.Context, out chan *model.Out) error {

	logrus.WithFields(logrus.Fields{"name": p.name, "path": p.path, "event": p.event}).Debug("formstruct.Parse")

	data, err := ioutil.ReadFile(p.path)
	if err != nil {
		return errors.Wrapf(err, "unable read file %s", p.path)
	}

	lines := strings.Split(string(data), "\n")

	claim := &model.VehicleClaim{
		Success: true,
	}

	state := StateInit

	key1 := "Новая запись в форме:"
	key2 := "*"

	for _, line := range lines {

		ok := false

		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, key1) {
			line = strings.ReplaceAll(line, key1, "")
			line = strings.TrimSpace(line)
			claim.District = line
			continue
		} else if strings.HasPrefix(line, key2) {
			line = strings.ReplaceAll(line, key2, "")
			line = strings.TrimSpace(line)

			if s, ok := p.headers[line]; ok {
				state = s
				continue
			}
		}

		switch state {
		case StateKind:
			claim.Company.Activity = line
		case StateName:
			claim.Company.Title = line
		case StateAddress:
			claim.Company.Address = line
		case StateINN:

			if claim.Company.TIN, ok = parser.ParseInt64(util.TrimSpaces(line)); ok == false {
				claim.Reason = append(claim.Reason, "ИНН не является числом")
				claim.Success = false
			}

			tdig := util.DigitsCount(claim.Company.TIN)
			if tdig < 10 || tdig > 12 {
				claim.Reason = append(claim.Reason, "ИНН меньше 10 или больше 12 знаков")
				claim.Success = false
			}

		case StateFIO:
			claim.Company.HeadName = strings.TrimSpace(line)
		case StatePhone:
			claim.Company.HeadPhone = util.TrimSpaces(line)
		case StateEMail:
			claim.Company.HeadEmail = util.TrimSpaces(line)
		case StateCars:
			claim.Source = line
			if claim.Passes, ok = parser.ParseVehicles(claim.Source); ok == false {
				claim.Reason = append(claim.Reason, "не удалось разобрать список номер/фио")
				claim.Success = false
			}
		case StateAgreement:
			claim.Agreement = line
		case StateReliability:
			claim.Reliability = line
		}
	}

	out <- &model.Out{
		Kind:  model.OutVehicleClaim,
		Value: claim,
		Event: p.event,
	}

	out <- &model.Out{
		Kind: model.OutUnknown,
	}

	return nil
}
