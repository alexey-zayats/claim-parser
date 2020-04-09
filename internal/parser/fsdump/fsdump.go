package fsdump

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/parser/formstruct"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"regexp"
	"strings"
)

// Parser ...
type Parser struct {
}

// Name ...
const Name = "fsdump"

// Register ...
func Register() {
	parser.Instance().Add(Name, NewParser)
}

// NewParser ...
func NewParser() (parser.Backend, error) {
	return &Parser{}, nil
}

// Parse ...
func (p *Parser) Parse(ctx context.Context, param *dict.Dict, out chan interface{}) error {

	var path string

	if iface, ok := param.Get("path"); ok {
		path = iface.(string)
	} else {
		return fmt.Errorf("not found 'path' in param dict")
	}

	var event = &model.Event{
		CreatedBy: 1,
	}
	if iface, ok := param.Get("event"); ok {
		event = iface.(*model.Event)
	}

	logrus.WithFields(logrus.Fields{"name": Name, "path": path, "event": event}).Debug("fsdump.Parse")

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "unable read file %s", path)
	}

	head := regexp.MustCompile(`/\* \d+ createdAt:(.[^\*]+)\*/`)
	nl := regexp.MustCompile(`\r?\n`)
	space := regexp.MustCompile(`^\s+(?:.+)\s+$`)
	object := regexp.MustCompile(`ObjectId\("(.[^"]+)"\)`)

	inClaim := false

	lines := make([]string, 0)
	var created string

	for _, line := range nl.Split(string(data), -1) {

		select {
		case <-ctx.Done():
			return fmt.Errorf("canceled")
		default:

			line = space.ReplaceAllString(line, "")

			if head.MatchString(line) {
				inClaim = true

				m := head.FindAllStringSubmatch(line, -1)
				if len(m) > 0 {
					created = fmt.Sprintf("\t\"createdAt\" : \"%s\",", m[0][1])
				}

				continue
			}

			if len(line) == 0 {

				inClaim = false

				claim, err := p.makeClaim(event, created, lines)
				if err != nil {
					return errors.Wrap(err, "unable unmarshal json")
				}

				logrus.WithFields(logrus.Fields{"[company]": claim.Company.Title}).Debug("claim")

				lines = make([]string, 0)

				out <- claim
				continue
			}

			if inClaim == false {
				continue
			}

			m := object.FindAllStringSubmatch(line, -1)
			if len(m) > 0 {
				line = fmt.Sprintf("\t\"_id\" : \"%s\",", m[0][1])
			}

			lines = append(lines, line)
		}
	}

	claim, err := p.makeClaim(event, created, lines)
	if err != nil {
		return errors.Wrap(err, "unable unmarshal json")
	}
	out <- claim

	out <- nil

	return nil
}

func (p *Parser) makeClaim(event *model.Event, created string, lines []string) (*model.Claim, error) {

	last := len(lines) - 1
	lines = append(lines, lines[last])
	copy(lines[3:], lines[2:last])
	lines[2] = created

	last = len(lines) - 1
	lines[last] = strings.ReplaceAll(lines[last], ",", "")

	var form Form
	var s = strings.Join(lines, "\n")
	if err := json.Unmarshal([]byte(s), &form); err != nil {
		return nil, errors.Wrap(err, "unable unmarshal json")
	}

	claim := &model.Claim{
		Code:       form.ID,
		Created:    form.Created.Time,
		DistrictID: Districts[form.FormID].ID,
		District:   Districts[form.FormID].Title,
		Event:      event,
	}

	for _, f := range form.Data {

		value := f.Value[0]

		switch Forms[form.FormID][f.FID] {
		case formstruct.StateKind:
			claim.Company.Activity = value
		case formstruct.StateName:
			claim.Company.Title = value
		case formstruct.StateAddress:
			claim.Company.Address = value
		case formstruct.StateINN:
			re := regexp.MustCompile(`\D`)
			claim.Company.INN = re.ReplaceAllString(value, "")
		case formstruct.StateFIO:
			fio := strings.Split(value, " ")

			if len(fio) < 3 {
				claim.Valid = false
				reason := "Нет данных по ФИО руководителя"
				claim.Reason = &reason
			} else {
				claim.Company.Head = model.Person{
					FIO: model.FIO{
						Surname:    fio[0],
						Name:       fio[1],
						Patronymic: fio[2],
					},
				}
			}
		case formstruct.StatePhone:
			claim.Company.Head.Contact.Phone = value
		case formstruct.StateEMail:
			claim.Company.Head.Contact.EMail = value
		case formstruct.StateCars:
			claim.Source = value
			claim.Cars = parser.ParseCars(value)
		case formstruct.StateAgreement:
			claim.Agreement = value
		case formstruct.StateReliability:
			claim.Reliability = value
		}
	}

	return claim, nil
}
