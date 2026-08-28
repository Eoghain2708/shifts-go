package money

import "fmt"

func FormattedPayment(wage int32) string {
	result := float32(wage) / 100
	return fmt.Sprintf("%.2f", result)
}
