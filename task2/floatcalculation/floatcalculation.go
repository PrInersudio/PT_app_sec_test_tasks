package floatcalculation

import (
	"errors"

	"github.com/shopspring/decimal"
)

type FloatCalculator struct{}

/*
Вычисления параметров:
-	X = X1 / X2 * X3 (значение возвращаем с точностью E);
-	Y = Y1 / Y2 * Y3 (значение возвращаем с точностью E);
-	IsEqual = “T” - если выводимые значения равны, и “F” в противном случае.
*/
func (c *FloatCalculator) FloatCalculation(
	X1, X2, X3 decimal.Decimal,
	Y1, Y2, Y3 decimal.Decimal,
	E int32,
) (
	X, Y decimal.Decimal,
	IsEqual string,
	err error,
) {
	if X2.IsZero() || Y2.IsZero() {
		X, Y, IsEqual, err = FloatCalculationError("деление на нуль")
		return
	}
	if E < 0 {
		X, Y, IsEqual, err = FloatCalculationError("отрицательная точность")
		return
	}
	err = nil
	X = X1.Mul(X3).DivRound(X2, E)
	Y = Y1.Mul(Y3).DivRound(Y2, E)
	if X.Equal(Y) {
		IsEqual = "T"
	} else {
		IsEqual = "F"
	}
	return
}

func FloatCalculationError(msg string) (
	X, Y decimal.Decimal,
	IsEqual string,
	err error,
) {
	X = decimal.Zero
	Y = decimal.Zero
	IsEqual = "F"
	err = errors.New(msg)
	return
}
