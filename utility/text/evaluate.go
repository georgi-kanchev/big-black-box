package text

import (
	"big-black-box/internal"
	"big-black-box/utility/number"
	"unicode"
)

var values = make([]float32, 0, 8)
var operators = make([]rune, 0, 8)

func Evaluate(mathExpression string, variables func(string) float32) float32 {
	mathExpression = Remove(mathExpression, " ")
	values = values[:0]
	operators = operators[:0]
	var bracketCountOpen, bracketCountClose int

	for i := 0; i < len(mathExpression); i++ {
		var c = rune(mathExpression[i])

		if unicode.IsDigit(c) || c == '.' {
			values = append(values, calcGetNumber(mathExpression, &i))
		} else if c == '(' {
			operators = append(operators, c)
			bracketCountOpen++
		} else if c == ')' {
			bracketCountClose++
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				if calcProcess() {
					return number.NaN()
				}
			}
			if len(operators) > 0 {
				operators = operators[:len(operators)-1]
			}
		} else if calcIsOperator(c) {
			// Check for unary minus or plus
			if (i == 0) || mathExpression[i-1] == '(' || calcIsOperator(rune(mathExpression[i-1])) {
				i++
				if i >= len(mathExpression) {
					return number.NaN()
				}
				var val = calcGetNumber(mathExpression, &i)
				if c == '-' {
					val = -val
				}
				values = append(values, val)
			} else {
				// Normal binary operator
				for len(operators) > 0 && calcPriority(operators[len(operators)-1]) >= calcPriority(c) {
					if calcProcess() {
						return number.NaN()
					}
				}
				operators = append(operators, c)
			}
		} else if unicode.IsLetter(c) {
			var start = i
			for i < len(mathExpression) && (unicode.IsLetter(rune(mathExpression[i])) || unicode.IsDigit(rune(mathExpression[i]))) {
				i++
			}
			var name = mathExpression[start:i]
			i--
			if variables == nil {
				return number.NaN()
			}
			var v = variables(name)
			if number.IsNaN(v) {
				return number.NaN()
			}
			values = append(values, v)
		}

		if bracketCountClose > bracketCountOpen {
			return number.NaN()
		}
	}

	if bracketCountOpen != bracketCountClose {
		return number.NaN()
	}

	for len(operators) > 0 {
		if calcProcess() {
			return number.NaN()
		}
	}

	if len(values) == 0 {
		return number.NaN()
	}
	return values[len(values)-1]
}

//=================================================================
// private

func repeatPad(padStr string, totalRunes int) string {
	if padStr == "" {
		return ""
	}
	internal.BuilderPush()
	var padRunes = []rune(padStr)
	var count = 0
	for count < totalRunes {
		for _, r := range padRunes {
			internal.BuilderWriteRune(r)
			count++
			if count >= totalRunes {
				var res = internal.BuilderResult()
				internal.BuilderPop()
				return truncateToRunes(res, totalRunes)
			}
		}
	}
	var res = internal.BuilderResult()
	internal.BuilderPop()
	return truncateToRunes(res, totalRunes)
}
func truncateToRunes(s string, maxRunes int) string {
	internal.BuilderPush()
	var count = 0
	for _, r := range s {
		if count >= maxRunes {
			break
		}
		internal.BuilderWriteRune(r)
		count++
	}
	var res = internal.BuilderResult()
	internal.BuilderPop()
	return res
}
func isSeparator(r rune) bool {
	return unicode.IsSpace(r) || r == '_' || r == '-' || r == '/' || r == '.'
}

func calcIsOperator(c rune) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '^' || c == '%'
}
func calcPriority(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/', '%':
		return 2
	case '^':
		return 3
	default:
		return 0
	}
}
func calcApplyOp(val1, val2 float32, op rune) float32 {
	switch op {
	case '+':
		return val1 + val2
	case '-':
		return val1 - val2
	case '*':
		return val1 * val2
	case '/':
		if val2 != 0 {
			return val1 / val2
		}
	case '%':
		if val2 != 0 {
			return number.DivisionRemainder(val1, val2)
		}
	case '^':
		return number.Power(val1, val2)
	}
	return number.NaN()
}
func calcProcess() bool {
	if len(values) < 2 || len(operators) < 1 {
		return true
	}
	var val2 = values[len(values)-1]
	values = values[:len(values)-1]
	var val1 = values[len(values)-1]
	values = values[:len(values)-1]
	var op = operators[len(operators)-1]
	operators = operators[:len(operators)-1]
	values = append(values, calcApplyOp(val1, val2, op))
	return false
}
func calcGetNumber(expr string, i *int) float32 {
	var start = *i
	for *i < len(expr) && (unicode.IsDigit(rune(expr[*i])) || expr[*i] == '.') {
		(*i)++
	}
	var numStr = expr[start:*i]
	(*i)--
	return ToNumber[float32](numStr)
}
