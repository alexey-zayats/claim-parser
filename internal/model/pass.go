package model

import "time"

// Pass ...
type Pass struct {
	ID int `db:"id"`
	// Вид деятельности
	CompanyBranch string `db:"company_branch"`
	// Код ОКВЭД
	CompanyOkved string `db:"company_okved"`
	// ИНН
	CompanyInn string `db:"company_inn"`
	// Наименование компании
	CompanyName string `db:"company_name"`
	// Адрес компании
	CompanyAddress string `db:"company_address"`
	// Телефон директора компании
	CompanyCeoPhone string `db:"company_ceo_phone"`
	// E-mai директора компании
	CompanyCeoEmail string `db:"company_ceo_email"`
	// Фамилия директора компании
	CompanyLastname string `db:"company_lastname"`
	// Имя директора компании
	CompanyFirstname string `db:"company_firstname"`
	// Отчество директора компании
	CompanyPatrname string `db:"company_patrname"`
	// Фамилия сотрудника компании
	EmployeeLastname string `db:"employee_lastname"`
	// Имя сотрудника компании
	EmployeeFirstname string `db:"employee_firstname"`
	// Отчество сотрудника компании
	EmployeePatrname string `db:"employee_patrname"`
	// Номер автомобиля сотрудника компании
	EmployeeCar string `db:"employee_car"`
	// Согласие обработки
	EmployeeAgree int `db:"employee_agree"`
	// Подтверждение актуальности
	EmployeeConfirm int `db:"employee_confirm"`

	// excel || fromstruct
	Source int `db:"source"`

	// Муниципальный округ
	District int `db:"district"`
	// Тип пропуска
	PassType int `db:"type"`
	// Номер пропуска
	PassNumber string `db:"number"`
	// Должность согласователя
	AlighnerPost string `db:"alighner_post"`
	// ФИО согласователя
	AlighnerName string `db:"alighner_name"`
	// Способ направления
	SendType string `db:"send_type"`
	// Статус
	Status int `db:"status"`
	// ID из таблицы files файла заявки
	FileID int `db:"file_id"`
	// Дата создания
	CreatedAt time.Time `db:"created_at"`
	// Пользователь
	CreatedBy int `db:"created_by"`
}
