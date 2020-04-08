package formstruct

import (
	"context"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"regexp"
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
	headers  map[string]State
	reNumber *regexp.Regexp
}

// Name ...
const Name = "formstruct"

// Register ...
func Register() {
	parser.Instance().Add(Name, NewParser)
}

// NewParser ...
func NewParser() (parser.Backend, error) {

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
	}, nil
}

func (p *Parser) parseCar(item string) model.Car {

	item = strings.ReplaceAll(item, "Номер: ", "")

	var numberS string
	var fioS string

	if strings.Contains(item, "-") {
		data := strings.Split(item, "-")

		numberS = data[0]
		numberS = strings.TrimSpace(numberS)
		numberS = strings.ReplaceAll(numberS, " ", "")
		numberS = strings.ToUpper(numberS)

		fioS = data[1]
		fioS = strings.TrimSpace(fioS)
	} else if strings.Contains(item, "г/н") {
		data := strings.Split(item, "г/н")
		numberS = data[1]
	} else {

		matches := p.reNumber.FindAllStringSubmatch(item, -1)
		if len(matches) > 0 {
			numberS = matches[0][1]
			fioS = strings.ReplaceAll(item, numberS, "")

		} else {
			numberS = p.reNumber.ReplaceAllString(item, "$1")
			fioS = p.reNumber.ReplaceAllString(item, "$2")

			if strings.Compare(fioS, item) == 0 {
				fioS = ""
			}
		}
	}

	numberS = strings.TrimSpace(numberS)
	numberS = strings.ReplaceAll(numberS, " ", "")
	numberS = strings.ToUpper(numberS)

	fioS = strings.TrimSpace(fioS)

	fioS = strings.ReplaceAll(fioS, ".", "")
	fio := regexp.MustCompile(`\s+`).Split(fioS, -1)

	car := model.Car{
		Number: numberS,
	}

	if len(fio) >= 3 {
		car.FIO.Surname = fio[0]
		car.FIO.Name = fio[1]
		car.FIO.Patronymic = fio[2]
		car.Valid = true
	} else {
		reason := "Нет данных по ФИО водителя"
		car.Reason = &reason
	}

	return car
}

// Parse ...
func (p *Parser) Parse(ctx context.Context, param *dict.Dict, out chan interface{}) error {

	var path string
	if iface, ok := param.Get("path"); ok {
		path = iface.(string)
	} else {
		return fmt.Errorf("not found 'path' in param dict")
	}

	var event = &model.Event{}
	if iface, ok := param.Get("event"); ok {
		event = iface.(*model.Event)
	}

	logrus.WithFields(logrus.Fields{"name": Name, "path": path, "event": event}).Debug("formstruct.Parse")

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "unable read file %s", path)
	}

	lines := strings.Split(string(data), "\n")

	claim := &model.Claim{
		Valid: true,
		Event: event,
	}

	state := StateInit

	key1 := "Новая запись в форме:"
	key2 := "*"

	for _, line := range lines {

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
			claim.Company.INN = strings.ReplaceAll(line, " ", "")
		case StateFIO:
			fio := strings.Split(line, " ")

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
		case StatePhone:
			claim.Company.Head.Contact.Phone = line
		case StateEMail:
			claim.Company.Head.Contact.EMail = line
		case StateCars:
			claim.Source = line
			claim.Cars = parser.ParseCars(line)
		case StateAgreement:
			claim.Agreement = line
		case StateReliability:
			claim.Reliability = line
		}
	}

	out <- claim

	out <- nil

	return nil
}
