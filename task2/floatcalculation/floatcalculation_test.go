package floatcalculation

import (
	"log"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func DecimalFromString(str string) decimal.Decimal {
	num, err := decimal.NewFromString(str)
	if err != nil {
		log.Fatal(err.Error())
	}
	return num
}

func TestFloatCalculation(t *testing.T) {
	calc := FloatCalculator{}

	cases := []struct {
		name       string
		X1, X2, X3 decimal.Decimal
		Y1, Y2, Y3 decimal.Decimal
		E          int32
		X          decimal.Decimal
		Y          decimal.Decimal
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
			Err:     "",
		},
		{
			name: "Положительная точность, не равны",
			X1:   DecimalFromString("1.5"), X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("4.5"), Y2: DecimalFromString("6.0"), Y3: DecimalFromString("2.0"),
			E: 3,
			X: DecimalFromString("2.250"), Y: DecimalFromString("1.500"),
			IsEqual: "F",
			Err:     "",
		},
		{
			name: "Нулевая точность, равны",
			X1:   DecimalFromString("3.4"), X2: DecimalFromString("2.0"), X3: DecimalFromString("1.5"),
			Y1: DecimalFromString("6.8"), Y2: DecimalFromString("4.0"), Y3: DecimalFromString("1.5"),
			E: 0,
			X: DecimalFromString("3"), Y: DecimalFromString("3"),
			IsEqual: "T",
			Err:     "",
		},
		{
			name: "Нулевая точность, не равны",
			X1:   DecimalFromString("3.4"), X2: DecimalFromString("2.0"), X3: DecimalFromString("1.5"),
			Y1: DecimalFromString("12"), Y2: DecimalFromString("4.0"), Y3: DecimalFromString("1.5"),
			E: 0,
			X: DecimalFromString("3"), Y: DecimalFromString("5"),
			IsEqual: "F",
			Err:     "",
		},
		{
			name: "Отрицательная точность, равны",
			X1:   DecimalFromString("15.5"), X2: DecimalFromString("3.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("31.0"), Y2: DecimalFromString("6.0"), Y3: DecimalFromString("3.1"),
			E: -1,
			X: DecimalFromString("20"), Y: DecimalFromString("20"),
			IsEqual: "T",
			Err:     "",
		},
		{
			name: "Отрицательная точность, не равны",
			X1:   DecimalFromString("15.5"), X2: DecimalFromString("6.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("31.0"), Y2: DecimalFromString("6.0"), Y3: DecimalFromString("3.1"),
			E: -1,
			X: DecimalFromString("10"), Y: DecimalFromString("20"),
			IsEqual: "F",
			Err:     "",
		},
		{
			name: "Отрицательные числа, равны",
			X1:   DecimalFromString("-1.5"), X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("-4.5"), Y2: DecimalFromString("6.0"), Y3: DecimalFromString("3.0"),
			E: 3,
			X: DecimalFromString("-2.250"), Y: DecimalFromString("-2.250"),
			IsEqual: "T",
			Err:     "",
		},
		{
			name: "Отрицательные числа, не равны",
			X1:   DecimalFromString("-1.5"), X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("-4.5"), Y2: DecimalFromString("6.0"), Y3: DecimalFromString("2.0"),
			E: 3,
			X: DecimalFromString("-2.250"), Y: DecimalFromString("-1.500"),
			IsEqual: "F",
			Err:     "",
		},
		{
			name: "Очень большая отрицательная точность, очень большие положительные числа, равны",
			X1:   DecimalFromString("1.0005e20"), X2: DecimalFromString("1e10"), X3: DecimalFromString("1e10"),
			Y1: DecimalFromString("1.004e20"), Y2: DecimalFromString("1e10"), Y3: DecimalFromString("1e10"),
			E: -19,
			X: DecimalFromString("1e20"), Y: DecimalFromString("1e20"),
			IsEqual: "T",
			Err:     "",
		},
		{
			name: "Очень большая отрицательная точность, очень большие положительные числа, не равны",
			X1:   DecimalFromString("1.00005e20"), X2: DecimalFromString("1e10"), X3: DecimalFromString("1e10"),
			Y1: DecimalFromString("3.004e20"), Y2: DecimalFromString("1e10"), Y3: DecimalFromString("1e10"),
			E: -19,
			X: DecimalFromString("1e20"), Y: DecimalFromString("3e20"),
			IsEqual: "F",
			Err:     "",
		},
		{
			name: "Очень большая отрицательная точность, очень большие отрицательные числа, равны",
			X1:   DecimalFromString("-1.0005e20"), X2: DecimalFromString("1e10"), X3: DecimalFromString("1e10"),
			Y1: DecimalFromString("-1.004e20"), Y2: DecimalFromString("1e10"), Y3: DecimalFromString("1e10"),
			E: -19,
			X: DecimalFromString("-1e20"), Y: DecimalFromString("-1e20"),
			IsEqual: "T",
			Err:     "",
		},
		{
			name: "Очень большая отрицательная точность, очень большие отрицательные числа, не равны",
			X1:   DecimalFromString("-1.0005e20"), X2: DecimalFromString("1e10"), X3: DecimalFromString("1e10"),
			Y1: DecimalFromString("-3.004e20"), Y2: DecimalFromString("1e10"), Y3: DecimalFromString("1e10"),
			E: -19,
			X: DecimalFromString("-1e20"), Y: DecimalFromString("-3e20"),
			IsEqual: "F",
			Err:     "",
		},
		{
			name: "Очень большая положительная точность, очень малые положительные числа, равны",
			X1:   DecimalFromString("1.00005e-20"), X2: DecimalFromString("1e-10"), X3: DecimalFromString("1e-10"),
			Y1: DecimalFromString("1.004e-20"), Y2: DecimalFromString("1e-10"), Y3: DecimalFromString("1e-10"),
			E: 21,
			X: DecimalFromString("1e-20"), Y: DecimalFromString("1e-20"),
			IsEqual: "T",
			Err:     "",
		},
		{
			name: "Очень большая положительная точность, очень малые положительные числа, не равны",
			X1:   DecimalFromString("1.00005e-20"), X2: DecimalFromString("1e-10"), X3: DecimalFromString("1e-10"),
			Y1: DecimalFromString("3.004e-20"), Y2: DecimalFromString("1e-10"), Y3: DecimalFromString("1e-10"),
			E: 21,
			X: DecimalFromString("1e-20"), Y: DecimalFromString("3e-20"),
			IsEqual: "F",
			Err:     "",
		},
		{
			name: "Очень большая положительная точность, очень малые отрицательные числа, равны",
			X1:   DecimalFromString("-1.00005e-20"), X2: DecimalFromString("1e-10"), X3: DecimalFromString("1e-10"),
			Y1: DecimalFromString("-1.004e-20"), Y2: DecimalFromString("1e-10"), Y3: DecimalFromString("1e-10"),
			E: 21,
			X: DecimalFromString("-1e-20"), Y: DecimalFromString("-1e-20"),
			IsEqual: "T",
			Err:     "",
		},
		{
			name: "Очень большая положительная точность, очень малые отрицательные числа, не равны",
			X1:   DecimalFromString("-1.00005e-20"), X2: DecimalFromString("1e-10"), X3: DecimalFromString("1e-10"),
			Y1: DecimalFromString("-3.004e-20"), Y2: DecimalFromString("1e-10"), Y3: DecimalFromString("1e-10"),
			E: 21,
			X: DecimalFromString("-1e-20"), Y: DecimalFromString("-3e-20"),
			IsEqual: "F",
			Err:     "",
		},
		{
			name: "X1 равен нулю",
			X1:   decimal.Zero, X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("4.0"), Y2: DecimalFromString("5.0"), Y3: DecimalFromString("6.0"),
			E: 3,
			X: DecimalFromString("0.000"), Y: DecimalFromString("4.800"),
			IsEqual: "F",
			Err:     "",
		},
		{
			name: "Y1 равен нулю",
			X1:   DecimalFromString("1.0"), X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: decimal.Zero, Y2: DecimalFromString("5.0"), Y3: DecimalFromString("6.0"),
			E: 3,
			X: DecimalFromString("1.500"), Y: DecimalFromString("0.000"),
			IsEqual: "F",
			Err:     "",
		},
		{
			name: "Оба X1 и Y1 равны нулю",
			X1:   decimal.Zero, X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: decimal.Zero, Y2: DecimalFromString("5.0"), Y3: DecimalFromString("6.0"),
			E: 3,
			X: DecimalFromString("0.000"), Y: DecimalFromString("0.000"),
			IsEqual: "T",
			Err:     "",
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
			Err:     "",
		},
		{
			name: "Y3 равен нулю",
			X1:   DecimalFromString("1.0"), X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: DecimalFromString("4.0"), Y2: DecimalFromString("5.0"), Y3: decimal.Zero,
			E: 3,
			X: DecimalFromString("1.500"), Y: decimal.Zero,
			IsEqual: "F",
			Err:     "",
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
			Err:     "",
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
			Err:     "",
		},
		{
			name: "Y1 и Y3 равны нулю",
			X1:   DecimalFromString("1.0"), X2: DecimalFromString("2.0"), X3: DecimalFromString("3.0"),
			Y1: decimal.Zero, Y2: DecimalFromString("5.0"), Y3: decimal.Zero,
			E: 3,
			X: DecimalFromString("1.500"), Y: decimal.Zero,
			IsEqual: "F",
			Err:     "",
		},
	}

	for _, test_case := range cases {
		test_case := test_case
		t.Run(test_case.name, func(t *testing.T) {
			t.Parallel()
			X, Y, IsEqual, err := calc.FloatCalculation(test_case.X1, test_case.X2, test_case.X3, test_case.Y1, test_case.Y2, test_case.Y3, test_case.E)
			assert.True(t, test_case.X.Equal(X), "ожидаемый X %v, полученный %v", test_case.X, X)
			assert.True(t, test_case.Y.Equal(Y), "ожидаемый Y %v, полученный %v", test_case.Y, Y)
			assert.Equal(t, test_case.IsEqual, IsEqual, "ожидаемый IsEqual %v, полученный %v", test_case.IsEqual, IsEqual)
			if test_case.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test_case.Err)
			}
		})
	}
}
