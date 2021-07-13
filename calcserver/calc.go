package calcserver

import "errors"

// Calculate logic
func calculate(calculation CalcOperation, x, y float64) (float64, error) {
	switch calculation {
	case addCalc:
		return x + y, nil
	case subtractCalc:
		return x - y, nil
	case multiplyCalc:
		return x * y, nil
	case divideCalc:
		if y == 0 {
			return 0, errors.New("division by zero: Y param cannot be 0")
		}
		return x / y, nil
	default:
		return 0, errors.New("unknown calculation operation")
	}
}
