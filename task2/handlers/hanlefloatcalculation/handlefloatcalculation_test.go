package hanlefloatcalculation_test

import (
	"FloatService/handlers/hanlefloatcalculation"
	"FloatService/handlers/hanlefloatcalculation/mocks"
	"FloatService/nulllogger"
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestHanleFloatCalculation(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		respError string
		mockError error
	}{
		{
			name: "Успех",
		},

		{
			name:      "Деление на нуль",
			respError: "деление на нуль",
			mockError: errors.New("деление на нуль"),
		},

		{
			name:      "Некорректный запрос",
			input:     "}{",
			respError: "Ошибка декодирования запроса.",
		},

		{
			name:      "Ошибка валидации: нет X1",
			input:     `{"X2":"2", "X3":"3","Y1":"1","Y2":"2","Y3":"3","E":5}`,
			respError: "Некорректный запрос",
		},

		{
			name:      "Ошибка валидации: нет X2",
			input:     `{"X1":"1", "X3":"3","Y1":"1","Y2":"2","Y3":"3","E":5}`,
			respError: "Некорректный запрос",
		},

		{
			name:      "Ошибка валидации: нет X3",
			input:     `{"X1":"1", "X2":"2","Y1":"1","Y2":"2","Y3":"3","E":5}`,
			respError: "Некорректный запрос",
		},

		{
			name:      "Ошибка валидации: нет Y1",
			input:     `{"X1":"1", "X2":"2", "X3":"3","Y2":"2","Y3":"3","E":5}`,
			respError: "Некорректный запрос",
		},

		{
			name:      "Ошибка валидации: нет Y2",
			input:     `{"X1":"1", "X2":"2", "X3":"3","Y1":"1","Y3":"3","E":5}`,
			respError: "Некорректный запрос",
		},

		{
			name:      "Ошибка валидации: нет Y3",
			input:     `{"X1":"1", "X2":"2", "X3":"3","Y1":"1","Y2":"2","E":5}`,
			respError: "Некорректный запрос",
		},

		{
			name:      "Ошибка валидации: нет E",
			input:     `{"X1":"1", "X2":"2", "X3":"3","Y1":"1","Y2":"2","Y3":"3"}`,
			respError: "Некорректный запрос",
		},
	}
	for _, test_case := range cases {
		test_case := test_case
		t.Run(test_case.name, func(t *testing.T) {
			t.Parallel()
			calculatorMock := mocks.NewFloatCalculatorInt(t)
			if test_case.respError == "" || test_case.mockError != nil {
				calculatorMock.On(
					"FloatCalculation",
					decimal.New(1, 0), decimal.New(2, 0), decimal.New(3, 0),
					decimal.New(1, 0), decimal.New(2, 0), decimal.New(3, 0),
					int32(5),
				).Return(
					decimal.New(5, -1),
					decimal.New(5, -1),
					"T",
					test_case.mockError,
				).Once()
			}
			handler := hanlefloatcalculation.New(slog.New(&nulllogger.NullLogger{}), calculatorMock)
			input := `{"X1":"1", "X2":"2", "X3":"3","Y1":"1","Y2":"2","Y3":"3","E":5}`
			if test_case.input != "" {
				input = test_case.input
			}
			req, err := http.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(input)))
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			require.Equal(t, rr.Code, http.StatusOK)
			body := rr.Body.String()
			var resp hanlefloatcalculation.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, test_case.respError, resp.Error)
		})
	}
}
