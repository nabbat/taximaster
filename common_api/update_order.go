package common_api

import "github.com/nabbat/taximaster/validator"

type (
	UpdateOrderRequest struct {
		// ИД заказа
		OrderID int `json:"order_id" validate:"required"`

		// Номер телефона
		Phone *string `json:"phone,omitempty" validate:"omitempty,max=30"`
		// Время подачи
		SourceTime *string `json:"source_time,omitempty" validate:"omitempty,datetime=20060102150405"`
		// Предварительный заказ
		IsPrior *bool `json:"is_prior,omitempty" validate:"omitempty"`
		// Заказчик
		Customer *string `json:"customer,omitempty" validate:"omitempty"`
		// Пассажир
		Passenger *string `json:"passenger,omitempty" validate:"omitempty"`
		// Комментарий
		Comment *string `json:"comment,omitempty" validate:"omitempty"`
		// ИД группы экипажей
		CrewGroupID *int `json:"crew_group_id,omitempty" validate:"omitempty"`
		// ИД клиента (необязателен, если phone присутствует)
		ClientID *int `json:"client_id,omitempty" validate:"omitempty"`
		// ИД службы ЕДС
		UdsID *int `json:"uds_id,omitempty" validate:"omitempty"`
		// ИД тарифа
		TariffID *int `json:"tariff_id,omitempty" validate:"omitempty"`
		// Массив адресов. Первый элемент — адрес подачи(обязательно), последний — адрес назначения, между ними — остановки
		Addresses *[]Address `json:"addresses,omitempty" validate:"omitempty"`
		// Массив параметров заказа. Устарело. Рекомендуется использовать параметр attribute_values
		OrderParams *[]int `json:"order_params,omitempty" validate:"omitempty"`
		// Массив значений атрибутов
		AttributeValues *[]AttributeValue `json:"attribute_values,omitempty" validate:"omitempty"`
		// Сумма заказа
		CostOrder *float64 `json:"cost_order,omitempty" validate:"omitempty"`
		// ИД состояния заказа
		StateID *int `json:"state_id,omitempty" validate:"omitempty"`
		// ИД скидки
		DiscountID *int `json:"discount_id,omitempty" validate:"omitempty"`
		// Автоматически подобрать скидку, если не указана явно
		AutoSelectDiscount *bool `json:"auto_select_discount,omitempty" validate:"omitempty"`
		// Автоматически подобрать тариф, если не указан явно
		AutoSelectTariff *bool `json:"auto_select_tariff,omitempty" validate:"omitempty"`
		// Автоматически пересчитать сумму заказа
		AutoRecalcCost *bool `json:"auto_recalc_cost,omitempty" validate:"omitempty"`
		// Автоматически обновить параметры заказа по клиенту и группе клиента
		AutoUpdateOrderParams *bool `json:"auto_update_order_params,omitempty" validate:"omitempty"`
		// Email для уведомлений
		Email *string `json:"email,omitempty" validate:"omitempty,email"`
		// Время перехода из предварительного в текущие заказы, мин
		PriorToCurrentBeforeMinutes *int `json:"prior_to_current_before_minutes,omitempty" validate:"omitempty"`
		// Номер рейса
		FlightNumber *string `json:"flight_number,omitempty" validate:"omitempty"`
		// Использовать специальную проверку перед изменением заказа
		NeedCustomValidate *bool `json:"need_custom_validate,omitempty" validate:"omitempty"`
	}

	UpdateOrderResponse struct {
		// Текст ошибки для пользователя
		Message string `json:"message"`
	}
)

// Изменение информации по заказу
func (cl *Client) UpdateOrder(req UpdateOrderRequest) (UpdateOrderResponse, error) {
	var response = UpdateOrderResponse{}

	err := validator.Validate(req)
	if err != nil {
		return response, err
	}

	/*
		100	Заказ не найден
		101	Состояние заказа не найдено
		102	Тариф не найден
		103	Скидка не найдена
		104	Группа экипажа не найдена
		105	Служба не найдена
		106	Клиент не найден
		107 Изменение состояния не соответствует необходимым условиям
		108	Параметр заказа не найден
		109	Атрибут не может быть привязан к заказу
		110	Ошибка специальной проверки заказа перед изменением. В ответе будет возвращаться:
		 "data": {
		   "message":"Текст ошибки для пользователя."
		 }
		111	Недостаточно средств на безналичном счете клиента в ТМ
		112	Для клиента запрещена оплата заказа наличными. Клиент должен максимально использовать в заказе безналичную оплату (оплату с основного счета)
	*/
	e := errorMap{
		100: ErrOrderNotFound,
		101: ErrOrderStateNotFound,
		102: ErrTariffNotFound,
		103: ErrDiscountNotFound,
		104: ErrCrewNotFound,
		105: ErrUdsNotFound,
		106: ErrClientNotFound,
		107: ErrStateCannotBeChanged,
		108: ErrOrderParameterNotFound,
		109: ErrAttributeCannotBeBoundOrder,
		110: ErrSpecialOrderCheck,
		111: ErrInsufficientFundsCashless,
		112: ErrCashPaymentNotAllowed,
	}

	err = cl.PostJson("update_order", e, req, &response)

	return response, err
}
