package fsdump

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/parser/fs"
	"github.com/alexey-zayats/claim-parser/internal/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"regexp"
	"strings"
)

// Parser ...
type Parser struct {
	event *model.Event
	path  string
	name  string
}

// Register ...
func Register() {
	parser.Instance().Add("form.struct.dump", NewParser)
}

// NewParser ...
func NewParser(name string) (parser.Backend, error) {
	return &Parser{
		name: name,
	}, nil
}

// WithEvent ..
func (p *Parser) WithEvent(event *model.Event) {
	p.event = event
	p.path = event.Filepath
}

// Parse ...
func (p *Parser) Parse(ctx context.Context, out chan *model.Out) error {

	logrus.WithFields(logrus.Fields{"name": p.name, "path": p.path, "event": p.event}).Debug("fsdump.Parse")

	data, err := ioutil.ReadFile(p.path)
	if err != nil {
		return errors.Wrapf(err, "unable read file %s", p.path)
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

				claim, err := p.makeClaim(created, lines)
				if err != nil {
					return errors.Wrap(err, "unable make claim")
				}

				lines = make([]string, 0)

				out <- &model.Out{
					Kind:  model.OutVehicleClaim,
					Value: claim,
					Event: p.event,
				}
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

	if len(lines) > 0 {
		claim, err := p.makeClaim(created, lines)
		if err != nil {
			return errors.Wrap(err, "unable make claim")
		}
		out <- &model.Out{
			Kind:  model.OutVehicleClaim,
			Event: p.event,
			Value: claim,
		}
	}

	out <- &model.Out{
		Kind: model.OutUnknown,
	}

	return nil
}

func (p *Parser) makeClaim(created string, lines []string) (*model.VehicleClaim, error) {

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

	claim := &model.VehicleClaim{
		Code:       form.ID,
		Created:    form.Created.Time,
		DistrictID: Districts[form.FormID].ID,
		District:   Districts[form.FormID].Title,
		Success:    true,
	}

	for _, f := range form.Data {

		ok := false

		value := f.Value[0]

		switch Forms[form.FormID][f.FID] {
		case fs.StateKind:
			claim.Company.Activity = value
		case fs.StateName:
			claim.Company.Title = value
		case fs.StateAddress:
			claim.Company.Address = value
		case fs.StateINN:

			claim.Company.INN = util.TrimSpaces(value)
			tdig := len(claim.Company.INN)
			if tdig < 10 || tdig > 12 {
				claim.Reason = append(claim.Reason, "ИНН меньше 10 или больше 12 знаков")
				claim.Success = false
			}

		case fs.StateFIO:
			claim.Company.HeadName = strings.TrimSpace(value)
		case fs.StatePhone:
			claim.Company.HeadPhone = util.TrimSpaces(value)
		case fs.StateEMail:
			claim.Company.HeadEmail = util.TrimSpaces(value)
		case fs.StateCars:
			claim.Source = value
			if claim.Passes, ok = parser.ParseVehicles(claim.Source); ok == false {
				claim.Reason = append(claim.Reason, "не удалось разобрать список номер/фио")
				claim.Success = false
			}
		case fs.StateAgreement:
			claim.Agreement = value
		case fs.StateReliability:
			claim.Reliability = value
		}
	}

	return claim, nil
}
