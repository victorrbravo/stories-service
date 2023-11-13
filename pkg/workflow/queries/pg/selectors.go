package pg

import (
	"fmt"
	"stories/types"
)

const (
	MAX_YELLOW = 245
	MIN_YELLOW = 90
)

func calculateYellowRGB(brightness int, min int, max int) (int, int, int) {
	// Ensure the brightness is within the valid range (0 to 100)
	brightness = clamp(float64(brightness), float64(min), float64(max))

	red := 252
	green := 250
	blue := brightness
	return red, green, blue
}

// Helper function to clamp an integer value within a specified range
func clamp(x float64, min float64, max float64) int {
	scaleX := 0.0
	if (max - min) > 0 {
		scaleX = ((x - min) / (max - min)) * (100)
	}
	result := int(scaleX/100*155 + 90.0)
	resultInv := 90 + 155 - (result - 90)
	return resultInv
}

// Helper function to convert RGB values to a hexadecimal string
func rgbToHexString(red, green, blue int) string {
	return fmt.Sprintf("#%02X%02X%02X", red, green, blue)
}

func SelectColor(intensity types.Int, min types.Int, max types.Int) string {
	// Example intensity value for yellow

	// Calculate the RGB values for yellow based on the intensity
	red, green, blue := calculateYellowRGB(types.AsInt(intensity),
		types.AsInt(min),
		types.AsInt(max))

	// Convert RGB values to a hexadecimal string
	hexString := rgbToHexString(red, green, blue)

	// Display the result
	return hexString
}
