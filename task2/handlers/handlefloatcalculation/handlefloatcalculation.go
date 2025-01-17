package handlefloatcalculation

import (
	"FloatService/response"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

/*
E назначен указателем для того,
чтобы конкретно валидировалось его отсуствие в запросе.
Если он просто значение, то валидатор обрасывает запросы,
где E=0.
Нашёл данное решение в issue на гитхабе.
*/
type Request struct {
	X1 decimal.Decimal `json:"X1" validate:"required"`
	X2 decimal.Decimal `json:"X2" validate:"required"`
	X3 decimal.Decimal `json:"X3" validate:"required"`
	Y1 decimal.Decimal `json:"Y1" validate:"required"`
	Y2 decimal.Decimal `json:"Y2" validate:"required"`
	Y3 decimal.Decimal `json:"Y3" validate:"required"`
	E  *int32          `json:"E" validate:"required"`
}

type Response struct {
	response.Response
	X       decimal.Decimal `json:"X"`
	Y       decimal.Decimal `json:"Y"`
	IsEqual string          `json:"IsEqual"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=FloatCalculatorInt
type FloatCalculatorInt interface {
	FloatCalculation(
		X1, X2, X3 decimal.Decimal,
		Y1, Y2, Y3 decimal.Decimal,
		E int32,
	) (
		X, Y decimal.Decimal,
		IsEqual string,
		err error,
	)
}

// создание нового обработчика запроса
func New(log *slog.Logger, calculator FloatCalculatorInt) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.handlefloatcalculation.New"
		// добавляем в логи имя функции и ID запроса
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Debug("Чтение запроса.")
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Ошибка декодирования тела запроса.", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error("Ошибка декодирования запроса."))
			return
		}
		log.Debug("Декодировано тело запроса.", slog.Any("request", req), slog.Any("E", DereferenceToString(req.E)))
		log.Debug("Валидация запроса.")
		if err := validator.New(validator.WithRequiredStructEnabled()).Struct(req); err != nil {
			log.Error("Некорректный запрос.", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error("Некорректный запрос"))
			return
		}
		log.Debug("Валидация запроса прошла успешно.")
		log.Debug("Начинаем расчёты.")
		X, Y, IsEqual, err := calculator.FloatCalculation(
			req.X1, req.X2, req.X3,
			req.Y1, req.Y2, req.Y3,
			*req.E,
		)
		if err != nil {
			log.Error("Ошибка в расчётах.", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		log.Debug("Расчёты окончены.")
		log.Debug("Отправляем ответ.")
		render.JSON(w, r, Response{
			Response: response.OK(),
			X:        X,
			Y:        Y,
			IsEqual:  IsEqual,
		})
		log.Info("Результаты отправлены.")
	}
}

func DereferenceToString(p *int32) string {
	if p != nil {
		return strconv.FormatInt(int64(*p), 10)
	}
	return "nil"
}
