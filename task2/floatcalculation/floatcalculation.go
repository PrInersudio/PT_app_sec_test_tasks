package floatcalculation

import "github.com/shopspring/decimal"

/*
Вычисления параметров:
-	X = X1 / X2 * X3 (значение возвращаем с точностью E);
-	Y = Y1 / Y2 * Y3 (значение возвращаем с точностью E);
-	IsEqual = “T” - если выводимые значения равны, и “F” в противном случае.
*/
func floatCalculation( // входные переменные
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
) {
	X = X1.Mul(X3).DivRound(X2, E)
	Y = Y1.Mul(Y3).DivRound(Y2, E)
	if X.Equal(Y) {
		IsEqual = "T"
	} else {
		IsEqual = "F"
	}
	return
}
