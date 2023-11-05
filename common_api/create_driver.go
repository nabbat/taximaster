package common_api

import "github.com/nabbat/taximaster/validator"

type (
	CreateDriverRequest struct {
		// ФИО водителя
		Name string `json:"name" validate:"required"`
		// ИД основного автомобиля
		CarID int `json:"car_id" validate:"required"`

		// Пароль (обязательное поле, если используется сервер связи с водителями)
		Password string `json:"password" validate:"omitempty"`
		// ИД службы ЕДС (обязательное поле, если используется ЕДС, иначе можно не указывать)
		UdsID int `json:"uds_id,omitempty" validate:"omitempty"`
		// Паспортные данные
		Passport string `json:"passport,omitempty" validate:"omitempty"`
		// Водительское удостоверение
		DriverLicense string `json:"driver_license,omitempty" validate:"omitempty"`
		// Тип работника (0 - работник компании, 1 - частник)
		EmployeeType int `json:"employee_type" validate:"omitempty,eq=0|eq=1"`
		// День рождения
		Birthday string `json:"birthday,omitempty" validate:"omitempty,datetime=20060102150405"`
		// Табельный номер
		Number string `json:"number,omitempty" validate:"omitempty"`
		// Удостоверение
		License string `json:"license,omitempty" validate:"omitempty"`
		// Дата приема на работу
		StartDate string `json:"start_date,omitempty" validate:"omitempty,datetime=20060102150405"`
		// Дата окончания договора
		LicDate string `json:"lic_date,omitempty" validate:"omitempty,datetime=20060102150405"`
		// Терминальный аккаунт (если не указан, будет сгенерирован автоматически), должен состоять из 5 цифр
		TermAccount string `json:"term_account,omitempty" validate:"omitempty"`
		// Описание
		Comment string `json:"comment,omitempty" validate:"omitempty"`
		// Фотография водителя
		DriverPhoto string `json:"driver_photo,omitempty" validate:"omitempty,base64"`
		// Имя для TaxoPhone
		NameForTaxophone string `json:"name_for_taxophone,omitempty" validate:"omitempty"`
		// Массив телефонов водителя
		Phones *[]Phone `json:"phones,omitempty" validate:"omitempty"`
		// Массив параметров водителя. Устарело. Рекомендуется использовать параметр attribute_values
		OrderParams *[]int `json:"order_params,omitempty" validate:"omitempty"`
		// Массив значений атрибутов
		AttributeValues *[]AttributeValue `json:"attribute_values,omitempty" validate:"omitempty"`
	}

	CreateDriverResponse struct {
		DriverID int `json:"driver_id"`
	}
)

// Создание водителя
func (cl *Client) CreateDriver(req CreateDriverRequest) (CreateDriverResponse, error) {
	var response = CreateDriverResponse{}

	err := validator.Validate(req)
	if err != nil {
		return response, err
	}

	/*
		100 Автомобиль с ИД=ID не найден
		101 Служба ЕДС с ИД=ID не найдена
		102 Параметр с ИД=ID не найден или не может быть привязан к водителю
		103 Терминальный аккаунт не уникален
		104 Некорректный терминальный аккаунт
		107 Основной телефон может быть только один
		108 Водитель должен иметь основной телефон
	*/
	e := errorMap{
		100: ErrCarNotFound,
		101: ErrUdsNotFound,
		102: ErrParameterNotFoundOrCannotBeBoundDriver,
		103: ErrDriverConflictByTerminalAccount,
		104: ErrTerminalAccountIncorrect,
		107: ErrConflictByPrimaryPhone,
		108: ErrDriverRequiredPrimaryPhone,
	}

	err = cl.PostJson("create_driver", e, req, &response)

	return response, err
}
