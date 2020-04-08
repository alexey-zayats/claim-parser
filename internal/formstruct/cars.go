package formstruct

import (
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"regexp"
	"strings"
)

var regs = []*regexp.Regexp{
	regexp.MustCompile(`((?:\p{L}{1})(?:\s+)?(?:\d{3})(?:\s+)?(?:\p{L}{2})(?:\s+)?(?:\d{2,3})\s?(?i:rus?)?)(?:\s+)?((?:\p{L}+)(?:\s+)?(?:\p{L}+)(?:\s+)?(?:\p{L}+)(?:\s+)?)`),
	regexp.MustCompile(`((?:\p{L}{1})(?:\s+)?(?:\d{3})(?:\s+)?(?:\p{L}{2}))(?:\s+)?((?:\p{L}+)(?:\s+)?(?:\p{L}+)(?:\s+)?(?:\p{L}+)(?:\s+)?)`),
	regexp.MustCompile(`((?:\p{L}{1})(?:\s+)?(?:\d{3})(?:\s+)?(?:\p{L}{2})(?:\s+)?(?:\d{2,3})\s?(?i:rus?)?)(?:\s+)?`),
	regexp.MustCompile(`((?:\p{L}{1})(?:\s+)?(?:\d{3})(?:\s+)?(?:\p{L}{2})(?:\s+)?(?:\d{2,3}))`),
}

// ParseCars ...
func ParseCars(line string) []model.Car {

	rm := []string{
		"номер:",
		"гос номер",
		"водитель",
		"г/н",
		"красный",
		"фольксваген туарег",
		"тойота рактис",
		"Mazda3",
		"Hyundai i30",
		"ВАЗ21140",
		"Nissan Note",
		"Toyota Corolla",
		"Opel Mokka",
		"Nissan Tiida",
		"Hyundai solaris",
		"Scoda Octavia",
		"Kia Sportage",
		"Opel Antara",
		"Chevrolet Lacetti",
		"Opel GTS",
		"Honda Civic",
		"Kia Rio",
		"Hyundai Solaris",
		"Xynday Elantra",
		"Audi A5",
		"Xynday Tocson",
		"Chevrolet cobalt",
		"Kia rio",
		"Hyundai Verna",
		"Skoda octavia",
		"Ford Focus",
		"skoda octavia",
		"Volkswagen tiguan",
		"Citroen",
		"Ford focus",
		"Ford focus",
		"Nissan Juke",
		"Mercedes Benz",
		"Volkswagen Passat 975",
		"Opel Astra",
		"Lexus is",
		"Hyundai Solaris",
		"Mitsubishi Colt",
		"ABH Toyota Crown",
		"KIA Soul",
		"Mazda 3",
		"Hyundai Solaris",
		"KIA CEED",
		"Hyundai Getz",
		"Volkswagen Jetta",
		"Hyundai sonata",
		"Лада Приора",
		"Mazda Demio",
		"KIA Optima",
		"Ford Focus",
		"Hyundai Solaris",
		"Opel Corsa",
		"Chevrolet Lacetti",
		"Volvo V-40",
		"Lexus IS250",
		"Mazda 3",
		"Kia Rio",
		"Ford Fusion",
		"Kia Rio",
		"Hyundai Getz",
		"KIA Ria",
		"peugeot 307",
		"Volkswagen Scirocco",
		"BMW 1",
		"Kia Ceed",
		"Opel Astra",
		"Honda Fit",
		"Toyota Chaser",
		"Chevrolet cruze",
		"Honda Accord",
		"Hyundai Solaris",
		"Kia Ceed",
	}

	line = strings.ReplaceAll(line, "–", "")
	line = strings.ReplaceAll(line, "—", "")
	line = strings.ReplaceAll(line, "-", "")
	line = regexp.MustCompile(`[\(\)–\,\.\r\n\t;]`).ReplaceAllString(line, " ")

	for i := range rm {
		re := regexp.MustCompile(fmt.Sprintf("(?i:%s)", rm[i]))
		line = re.ReplaceAllString(line, " ")
	}

	line = regexp.MustCompile(`\d\.`).ReplaceAllString(line, " ")
	line = regexp.MustCompile(`^\d\s`).ReplaceAllString(line, " ")
	line = regexp.MustCompile(`\s\d\s`).ReplaceAllString(line, " ")

	return parseCar(line)
}

// Pair ...
type Pair struct {
	Key   string
	Value string
}

func parseCar(item string) []model.Car {

	var number string
	var fio string

	pairs := make([]Pair, 0)

	for i := range regs {
		if regs[i].MatchString(item) {
			matches := regs[i].FindAllStringSubmatch(item, -1)
			if len(matches) > 0 {
				for _, a := range matches {

					pair := Pair{
						Key: a[1],
					}

					if len(a) > 2 {
						pair.Value = a[2]
					}

					pairs = append(pairs, pair)
				}
			}
			break
		}
	}

	cars := make([]model.Car, 0)

	re := regexp.MustCompile(`\s`)

	for _, pair := range pairs {

		number = pair.Key
		fio = pair.Value

		number = re.ReplaceAllString(number, "")
		number = strings.ToUpper(number)

		fio = regexp.MustCompile(`\w`).ReplaceAllString(fio, "")

		fio = strings.TrimSpace(fio)
		fio := regexp.MustCompile(`\s+`).Split(fio, -1)

		//fmt.Printf("<Number>: %s; <FIO>: [%#v]\n", number, strings.Join(fio, ", ") )

		car := model.Car{
			Number: number,
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

		cars = append(cars, car)
	}

	return cars
}
