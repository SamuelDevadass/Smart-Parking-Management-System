package services

import (
	"math"
)

func CalculateBillAmount(totalSeconds float64, vehicleType string) float64 {
	// Calculate total hours using ceiling division and cast to int
	totalHours := int(math.Ceil(totalSeconds / 3600))

	// Determine base amount
	baseAmount := 20.0
	if vehicleType == "Four Wheeler" {
		baseAmount = 40.0
	}

	// Calculate extra cost based on hours
	var extraCost float64
	switch {
	case totalHours <= 3:
		extraCost = 0
	case totalHours == 4:
		extraCost = 0.30 * float64(totalHours-3)
	case totalHours < 5:
		extraCost = 0.30 + (0.40 * float64(totalHours-4))
	default:
		extraCost = 0.30 + 0.40 + (0.50 * float64(totalHours-5))
	}

	return baseAmount + extraCost
}
