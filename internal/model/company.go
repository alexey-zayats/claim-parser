package model

/*
	ИНН (Идентификационный Номер Налогоплательщика)
		ITN (Individual Taxpayer Number) — для физических лиц
		TIN (Taxpayer Identification Number — для юридических лиц;

	ОГРН (Основной Государственный Регистрационный Номер)
		PSRN (Primary State Registration Number)

	Источник: https://englishfull.ru/znat/inn-po-angliyski.html
*/

// Company данные о компании
type Company struct {

	// Activity вид детельности
	Activity string
	// Title название компании
	Title string
	// Address адрес ЮЛ
	Address string

	INN  string
	OGRN string

	// ФИО директора
	HeadName string

	// Телефон директора
	HeadPhone string
	// Email директора
	HeadEmail string
}
