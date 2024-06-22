package tests

import (
	"FloatService/handlers/handlefloatcalculation"
	"FloatService/response"
	"encoding/json"
	"log"
	"net/url"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

const (
	host         = "localhost:8081"
	limit        = 50
	interval     = time.Second
	limitRateMsg = "Слишком много запросов."
	clearance    = 5 // значение, на которое может отличаться количество запросов от лимита
)

// проверяем, нужные ли ответы присылает обработчик
func TestFloatService_HanlderResponses(t *testing.T) {
	cases := []struct {
		name     string
		request  string
		response string
	}{
		{
			name:     "Успех",
			request:  `{"X1":"1", "X2":"2", "X3":"3","Y1":"1","Y2":"2","Y3":"3","E":5}`,
			response: `{"status":"OK","X":"1.5","Y":"1.5","IsEqual":"T"}`,
		},

		{
			name:     "Деление на нуль",
			request:  `{"X1":"1", "X2":"0", "X3":"3","Y1":"1","Y2":"2","Y3":"3","E":5}`,
			response: `{"status": "Error", "error": "деление на нуль"}`,
		},

		{
			name:     "Некорректный запрос",
			request:  "}{",
			response: `{"status":"Error","error":"Ошибка декодирования запроса."}`,
		},

		{
			name:     "Ошибка валидации: нет X1",
			request:  `{"X2":"2", "X3":"3","Y1":"1","Y2":"2","Y3":"3","E":5}`,
			response: `{"status": "Error", "error": "Некорректный запрос"}`,
		},

		{
			name:     "Ошибка валидации: нет X2",
			request:  `{"X1":"1", "X3":"3","Y1":"1","Y2":"2","Y3":"3","E":5}`,
			response: `{"status": "Error", "error": "Некорректный запрос"}`,
		},

		{
			name:     "Ошибка валидации: нет X3",
			request:  `{"X1":"1", "X2":"2","Y1":"1","Y2":"2","Y3":"3","E":5}`,
			response: `{"status": "Error", "error": "Некорректный запрос"}`,
		},

		{
			name:     "Ошибка валидации: нет Y1",
			request:  `{"X1":"1", "X2":"2", "X3":"3","Y2":"2","Y3":"3","E":5}`,
			response: `{"status": "Error", "error": "Некорректный запрос"}`,
		},

		{
			name:     "Ошибка валидации: нет Y2",
			request:  `{"X1":"1", "X2":"2", "X3":"3","Y1":"1","Y3":"3","E":5}`,
			response: `{"status": "Error", "error": "Некорректный запрос"}`,
		},

		{
			name:     "Ошибка валидации: нет Y3",
			request:  `{"X1":"1", "X2":"2", "X3":"3","Y1":"1","Y2":"2","E":5}`,
			response: `{"status": "Error", "error": "Некорректный запрос"}`,
		},

		{
			name:     "Ошибка валидации: нет E",
			request:  `{"X1":"1", "X2":"2", "X3":"3","Y1":"1","Y2":"2","Y3":"3"}`,
			response: `{"status": "Error", "error": "Некорректный запрос"}`,
		},

		{
			name:     "Отправка целых",
			request:  `{"X1":1, "X2":2, "X3":3,"Y1":1,"Y2":2,"Y3":3,"E":5}`,
			response: `{"status":"OK","X":"1.5","Y":"1.5","IsEqual":"T"}`,
		},

		{
			name:     "Отправка чисел с плавающей запятой",
			request:  `{"X1":1.0, "X2":2.0, "X3":3.0,"Y1":1.0,"Y2":2.0,"Y3":3.0,"E":5}`,
			response: `{"status":"OK","X":"1.5","Y":"1.5","IsEqual":"T"}`,
		},
	}

	for _, test_case := range cases {
		test_case := test_case
		t.Run(test_case.name, func(t *testing.T) {
			t.Parallel()
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}
			var response map[string]interface{}
			json.Unmarshal([]byte(test_case.response), &response)
			e := httpexpect.Default(t, u.String())
			e.GET("/").
				WithText(test_case.request).
				Expect().
				Status(200).
				JSON().Object().IsEqual(response)
		})
	}
}

// проверяем сами вычисления
func TestFloatService_Сalculations(t *testing.T) {
	cases := []struct {
		name       string
		X1, X2, X3 decimal.Decimal
		Y1, Y2, Y3 decimal.Decimal
		E          int32
		X, Y       decimal.Decimal
		IsEqual    string
		Err        string
	}{
		{
			name: "Положительная точность, равны",
			X1:   DecimalFromString("1.5"), X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("4.5"), Y2: DecimalFromString("6.0"), Y3: DecimalFromString("3.0"),
			E: 3,
			X: DecimalFromString("2.250"), Y: DecimalFromString("2.250"),
			IsEqual: "T",
		},
		{
			name: "Положительная точность, не равны",
			X1:   DecimalFromString("1.5"), X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("4.5"), Y2: DecimalFromString("6.0"), Y3: DecimalFromString("2.0"),
			E: 3,
			X: DecimalFromString("2.250"), Y: DecimalFromString("1.500"),
			IsEqual: "F",
		},
		{
			name: "Нулевая точность, равны",
			X1:   DecimalFromString("3.4"), X2: DecimalFromString("2.0"), X3: DecimalFromString("1.5"),
			Y1: DecimalFromString("6.8"), Y2: DecimalFromString("4.0"), Y3: DecimalFromString("1.5"),
			E: 0,
			X: DecimalFromString("3"), Y: DecimalFromString("3"),
			IsEqual: "T",
		},
		{
			name: "Нулевая точность, не равны",
			X1:   DecimalFromString("3.4"), X2: DecimalFromString("2.0"), X3: DecimalFromString("1.5"),
			Y1: DecimalFromString("12"), Y2: DecimalFromString("4.0"), Y3: DecimalFromString("1.5"),
			E: 0,
			X: DecimalFromString("3"), Y: DecimalFromString("5"),
			IsEqual: "F",
		},
		{
			name: "Отрицательная точность, равны",
			X1:   DecimalFromString("15.5"), X2: DecimalFromString("3.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("31.0"), Y2: DecimalFromString("6.0"), Y3: DecimalFromString("3.1"),
			E: -1,
			X: DecimalFromString("20"), Y: DecimalFromString("20"),
			IsEqual: "T",
		},
		{
			name: "Отрицательная точность, не равны",
			X1:   DecimalFromString("15.5"), X2: DecimalFromString("6.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("31.0"), Y2: DecimalFromString("6.0"), Y3: DecimalFromString("3.1"),
			E: -1,
			X: DecimalFromString("10"), Y: DecimalFromString("20"),
			IsEqual: "F",
		},
		{
			name: "Отрицательные числа, равны",
			X1:   DecimalFromString("-1.5"), X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("-4.5"), Y2: DecimalFromString("6.0"), Y3: DecimalFromString("3.0"),
			E: 3,
			X: DecimalFromString("-2.250"), Y: DecimalFromString("-2.250"),
			IsEqual: "T",
		},
		{
			name: "Отрицательные числа, не равны",
			X1:   DecimalFromString("-1.5"), X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("-4.5"), Y2: DecimalFromString("6.0"), Y3: DecimalFromString("2.0"),
			E: 3,
			X: DecimalFromString("-2.250"), Y: DecimalFromString("-1.500"),
			IsEqual: "F",
		},
		{
			name: "Очень большая отрицательная точность, очень большие положительные числа, равны",
			X1:   DecimalFromString("1.0005e20"), X2: DecimalFromString("1e10"), X3: DecimalFromString("1e10"),
			Y1: DecimalFromString("1.004e20"), Y2: DecimalFromString("1e10"), Y3: DecimalFromString("1e10"),
			E: -19,
			X: DecimalFromString("1e20"), Y: DecimalFromString("1e20"),
			IsEqual: "T",
		},
		{
			name: "Очень большая отрицательная точность, очень большие положительные числа, не равны",
			X1:   DecimalFromString("1.00005e20"), X2: DecimalFromString("1e10"), X3: DecimalFromString("1e10"),
			Y1: DecimalFromString("3.004e20"), Y2: DecimalFromString("1e10"), Y3: DecimalFromString("1e10"),
			E: -19,
			X: DecimalFromString("1e20"), Y: DecimalFromString("3e20"),
			IsEqual: "F",
		},
		{
			name: "Очень большая отрицательная точность, очень большие отрицательные числа, равны",
			X1:   DecimalFromString("-1.0005e20"), X2: DecimalFromString("1e10"), X3: DecimalFromString("1e10"),
			Y1: DecimalFromString("-1.004e20"), Y2: DecimalFromString("1e10"), Y3: DecimalFromString("1e10"),
			E: -19,
			X: DecimalFromString("-1e20"), Y: DecimalFromString("-1e20"),
			IsEqual: "T",
		},
		{
			name: "Очень большая отрицательная точность, очень большие отрицательные числа, не равны",
			X1:   DecimalFromString("-1.0005e20"), X2: DecimalFromString("1e10"), X3: DecimalFromString("1e10"),
			Y1: DecimalFromString("-3.004e20"), Y2: DecimalFromString("1e10"), Y3: DecimalFromString("1e10"),
			E: -19,
			X: DecimalFromString("-1e20"), Y: DecimalFromString("-3e20"),
			IsEqual: "F",
		},
		{
			name: "Очень большая положительная точность, очень малые положительные числа, равны",
			X1:   DecimalFromString("1.00005e-20"), X2: DecimalFromString("1e-10"), X3: DecimalFromString("1e-10"),
			Y1: DecimalFromString("1.004e-20"), Y2: DecimalFromString("1e-10"), Y3: DecimalFromString("1e-10"),
			E: 21,
			X: DecimalFromString("1e-20"), Y: DecimalFromString("1e-20"),
			IsEqual: "T",
		},
		{
			name: "Очень большая положительная точность, очень малые положительные числа, не равны",
			X1:   DecimalFromString("1.00005e-20"), X2: DecimalFromString("1e-10"), X3: DecimalFromString("1e-10"),
			Y1: DecimalFromString("3.004e-20"), Y2: DecimalFromString("1e-10"), Y3: DecimalFromString("1e-10"),
			E: 21,
			X: DecimalFromString("1e-20"), Y: DecimalFromString("3e-20"),
			IsEqual: "F",
		},
		{
			name: "Очень большая положительная точность, очень малые отрицательные числа, равны",
			X1:   DecimalFromString("-1.00005e-20"), X2: DecimalFromString("1e-10"), X3: DecimalFromString("1e-10"),
			Y1: DecimalFromString("-1.004e-20"), Y2: DecimalFromString("1e-10"), Y3: DecimalFromString("1e-10"),
			E: 21,
			X: DecimalFromString("-1e-20"), Y: DecimalFromString("-1e-20"),
			IsEqual: "T",
		},
		{
			name: "Очень большая положительная точность, очень малые отрицательные числа, не равны",
			X1:   DecimalFromString("-1.00005e-20"), X2: DecimalFromString("1e-10"), X3: DecimalFromString("1e-10"),
			Y1: DecimalFromString("-3.004e-20"), Y2: DecimalFromString("1e-10"), Y3: DecimalFromString("1e-10"),
			E: 21,
			X: DecimalFromString("-1e-20"), Y: DecimalFromString("-3e-20"),
			IsEqual: "F",
		},
		{
			name: "X1 равен нулю",
			X1:   decimal.Zero, X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("4.0"), Y2: DecimalFromString("5.0"), Y3: DecimalFromString("6.0"),
			E: 3,
			X: DecimalFromString("0.000"), Y: DecimalFromString("4.800"),
			IsEqual: "F",
		},
		{
			name: "Y1 равен нулю",
			X1:   DecimalFromString("1.0"), X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: decimal.Zero, Y2: DecimalFromString("5.0"), Y3: DecimalFromString("6.0"),
			E: 3,
			X: DecimalFromString("1.500"), Y: DecimalFromString("0.000"),
			IsEqual: "F",
		},
		{
			name: "Оба X1 и Y1 равны нулю",
			X1:   decimal.Zero, X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: decimal.Zero, Y2: DecimalFromString("5.0"), Y3: DecimalFromString("6.0"),
			E: 3,
			X: DecimalFromString("0.000"), Y: DecimalFromString("0.000"),
			IsEqual: "T",
		},
		{
			name: "X2 равен нулю",
			X1:   DecimalFromString("1.0"), X2: decimal.Zero, X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("4.0"), Y2: DecimalFromString("5.0"), Y3: DecimalFromString("6.0"),
			E: 3,
			X: decimal.Zero, Y: decimal.Zero,
			IsEqual: "F",
			Err:     "деление на нуль",
		},
		{
			name: "Y2 равен нулю",
			X1:   DecimalFromString("1.0"), X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("4.0"), Y2: decimal.Zero, Y3: DecimalFromString("6.0"),
			E: 3,
			X: decimal.Zero, Y: decimal.Zero,
			IsEqual: "F",
			Err:     "деление на нуль",
		},
		{
			name: "X3 равен нулю",
			X1:   DecimalFromString("1.0"), X2: DecimalFromString("2.0"), X3: decimal.Zero,
			Y1: DecimalFromString("4.0"), Y2: DecimalFromString("5.0"), Y3: DecimalFromString("6.0"),
			E: 3,
			X: decimal.Zero, Y: DecimalFromString("4.800"),
			IsEqual: "F",
		},
		{
			name: "Y3 равен нулю",
			X1:   DecimalFromString("1.0"), X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("4.0"), Y2: DecimalFromString("5.0"), Y3: decimal.Zero,
			E: 3,
			X: DecimalFromString("1.500"), Y: decimal.Zero,
			IsEqual: "F",
		},
		{
			name: "X2 и Y2 равны нулю",
			X1:   DecimalFromString("1.0"), X2: decimal.Zero, X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("4.0"), Y2: decimal.Zero, Y3: DecimalFromString("6.0"),
			E: 3,
			X: decimal.Zero, Y: decimal.Zero,
			IsEqual: "F",
			Err:     "деление на нуль",
		},
		{
			name: "X3 и Y3 равны нулю",
			X1:   DecimalFromString("1.0"), X2: DecimalFromString("2.0"), X3: decimal.Zero,
			Y1: DecimalFromString("4.0"), Y2: DecimalFromString("5.0"), Y3: decimal.Zero,
			E: 3,
			X: decimal.Zero, Y: decimal.Zero,
			IsEqual: "T",
		},
		{
			name: "Все X параметры равны нулю",
			X1:   decimal.Zero, X2: decimal.Zero, X3: decimal.Zero,
			Y1: DecimalFromString("4.0"), Y2: DecimalFromString("5.0"), Y3: DecimalFromString("6.0"),
			E: 3,
			X: decimal.Zero, Y: decimal.Zero,
			IsEqual: "F",
			Err:     "деление на нуль",
		},
		{
			name: "Все Y параметры равны нулю",
			X1:   DecimalFromString("1.0"), X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: decimal.Zero, Y2: decimal.Zero, Y3: decimal.Zero,
			E: 3,
			X: decimal.Zero, Y: decimal.Zero,
			IsEqual: "F",
			Err:     "деление на нуль",
		},
		{
			name: "Все X и Y параметры равны нулю",
			X1:   decimal.Zero, X2: decimal.Zero, X3: decimal.Zero,
			Y1: decimal.Zero, Y2: decimal.Zero, Y3: decimal.Zero,
			E: 3,
			X: decimal.Zero, Y: decimal.Zero,
			IsEqual: "F",
			Err:     "деление на нуль",
		},
		{
			name: "X1 и X3 равны нулю",
			X1:   decimal.Zero, X2: DecimalFromString("2.0"), X3: decimal.Zero,
			Y1: DecimalFromString("4.0"), Y2: DecimalFromString("5.0"), Y3: DecimalFromString("6.0"),
			E: 3,
			X: decimal.Zero, Y: DecimalFromString("4.800"),
			IsEqual: "F",
		},
		{
			name: "Y1 и Y3 равны нулю",
			X1:   DecimalFromString("1.0"), X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: decimal.Zero, Y2: DecimalFromString("5.0"), Y3: decimal.Zero,
			E: 3,
			X: DecimalFromString("1.500"), Y: decimal.Zero,
			IsEqual: "F",
		},
	}

	for _, test_case := range cases {
		test_case := test_case
		t.Run(test_case.name, func(t *testing.T) {
			t.Parallel()
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}
			var resp interface{}
			if test_case.Err == "" {
				resp = handlefloatcalculation.Response{
					Response: response.OK(),
					X:        test_case.X,
					Y:        test_case.Y,
					IsEqual:  test_case.IsEqual,
				}
			} else {
				resp = response.Error(test_case.Err)
			}
			e := httpexpect.Default(t, u.String())
			e.GET("/").
				WithJSON(handlefloatcalculation.Request{
					X1: test_case.X1, X2: test_case.X2, X3: test_case.X3,
					Y1: test_case.Y1, Y2: test_case.Y2, Y3: test_case.Y3,
					E: &test_case.E,
				}).
				Expect().
				Status(200).
				JSON().Object().IsEqual(resp)
		})
	}
}

func DecimalFromString(str string) decimal.Decimal {
	num, err := decimal.NewFromString(str)
	if err != nil {
		log.Fatal(err.Error())
	}
	return num
}

// проверяем лимит на количество запросов
func TestFloatService_RateLimit(t *testing.T) {
	// спим, чтобы лимит не подействовал раньше из-за других тестов
	time.Sleep(2 * interval)
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	E := int32(5)
	request := handlefloatcalculation.Request{
		X1: DecimalFromString("1"), X2: DecimalFromString("2"), X3: DecimalFromString("3"),
		Y1: DecimalFromString("1"), Y2: DecimalFromString("2"), Y3: DecimalFromString("3"),
		E: &E,
	}
	e := httpexpect.Default(t, u.String())
	// запросы в пределах лимита должны получить статус 200
	statusCode := 0
	var numRequests int
	for numRequests = 0; numRequests < limit+clearance; numRequests++ {
		resp := e.GET("/").WithJSON(request).Expect()
		statusCode = resp.Raw().StatusCode
		if statusCode == 402 {
			resp.JSON().IsEqual(response.Error(limitRateMsg))
			break
		}
	}
	require.Equal(t, 402, statusCode)
	log.Printf("Лимит: %d, интервал: %v, количество посланных запросов до статуса 402: %d", limit, interval, numRequests)
	require.LessOrEqual(t, absInt(limit-numRequests), clearance)
}

func absInt(num int) int {
	if num < 0 {
		return -num
	}
	return num
}
