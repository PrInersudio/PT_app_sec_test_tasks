package hanlefloatcalculation

import (
	"float_service/response"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

type Request struct {
	X1 decimal.Decimal `json:"X1" validate:"required"`
	X2 decimal.Decimal `json:"X2" validate:"required"`
	X3 decimal.Decimal `json:"X3" validate:"required"`
	Y1 decimal.Decimal `json:"Y1" validate:"required"`
	Y2 decimal.Decimal `json:"Y2" validate:"required"`
	Y3 decimal.Decimal `json:"Y3" validate:"required"`
	E  int32           `json:"E" validate:"required"`
}

type Response struct {
	response.Response
	X       decimal.Decimal `json:"X"`
	Y       decimal.Decimal `json:"Y"`
	IsEqual string          `json:"IsEqual"`
}

type FloatCalculator interface {
	floatCalculation( // входные переменные
		X1 decimal.Decimal,
		X2 decimal.Decimal,
		X3 decimal.Decimal,
		Y1 decimal.Decimal,
		Y2 decimal.Decimal,
		Y3 decimal.Decimal,
		E int32,
	) ( // выходные переменные
		X decimal.Decimal,
		Y decimal.Decimal,
		IsEqual string,
	)
}

// создание нового обработчика запроса
func New(log *slog.Logger, floatCalculator FloatCalculator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.handlefloatcalculation.New"
		// добавляем в логи имя функции и ID запроса
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		// чтение запроса
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Ошибка декодирования тела запроса", err)
			render.JSON(w, r, response.Error("Ошибка декодирования запроса"))
			return
		}
		log.Info("Декодировано тело запроса", slog.Any("request", req))
		// валидация запроса
		if err := validator.New().Struct(req); err != nil {
			log.Error("Некорректный запрос", err)
			render.JSON(w, r, response.Error("Некорректный запрос"))
			return
		}
	}
}
