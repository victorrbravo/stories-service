package pg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ConvertToDaysHours(durationStr string) (string, error) {
	// Regular expression to match and extract the number of hours

	re := regexp.MustCompile(`(\d+)h`)

	// Check if the input matches the expected pattern
	matches := re.FindStringSubmatch(durationStr)
	if matches == nil {
		return durationStr, nil
	}
	newDuration := re.ReplaceAllString(durationStr, "")

	// Extract the number of hours as a string and convert it to an integer
	hoursStr := matches[1]
	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		return "", fmt.Errorf("invalid number of hours: %s", hoursStr)
	}

	// Calculate days and remaining hours
	days := hours / 24
	remainingHours := hours % 24

	// Build the result string
	var result strings.Builder
	if days > 0 {
		result.WriteString(fmt.Sprintf("%dd", days))
	}
	if remainingHours > 0 || days == 0 {
		if days > 0 {
			result.WriteString(" ")
		}
		result.WriteString(fmt.Sprintf("%dh", remainingHours))
	}
	newDuration = result.String() + newDuration
	return newDuration, nil
}

// formatMoney formats a float64 number as a string with currency format "$0,000.00"
func FormatToMoney(amount float64) string {
	// Convert the float to a string with 2 decimal places
	amountStr := fmt.Sprintf("%.2f", amount)

	// Split the string into whole and fractional parts
	parts := strings.Split(amountStr, ".")
	whole := parts[0]
	fractional := parts[1]

	// Add commas as thousand separators
	var result strings.Builder
	length := len(whole)
	for i, digit := range whole {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(digit)
	}

	// Concatenate the formatted whole part, decimal point, and fractional part
	formatted := result.String() + "." + fractional

	// Prepend the dollar sign
	return "$" + formatted
}

func ParseDollarAmount(input string) (float64, error) {
	// Check for and remove the leading '$' sign
	lenInput := len(input)
	multiply := 1.0
	if lenInput == 0 {
		return 0, fmt.Errorf("parse dollar failed: string length is 0")

	}
	amountStr := input
	if input[0] == '-' {
		amountStr = strings.TrimPrefix(input, "-")
		multiply = -1.0
	}

	if len(amountStr) == 0 || amountStr[0] != '$' {
		return 0, fmt.Errorf("invalid format: %s", input)
	}

	amountStr = strings.TrimPrefix(amountStr, "$")

	amountStr = strings.ReplaceAll(amountStr, ",", "")
	// Convert the remaining string to a float
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount: %v", err)
	}
	amount = amount * multiply
	return amount, nil
}

func ParsePostgresInterval(input string) (time.Duration, error) {
	// Regular expressions for various parts
	re := regexp.MustCompile(`(?:(\d+)\s*years?)?\s*(?:(\d+)\s*months?)?\s*(?:(\d+)\s*days?)?\s*(?:(\d+):(\d+):(\d+))?`)

	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return 0, fmt.Errorf("invalid interval format: %s", input)
	}

	var totalDuration time.Duration

	// Parse and accumulate years
	if matches[1] != "" {
		years, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, fmt.Errorf("invalid years: %v", err)
		}
		totalDuration += time.Duration(years) * 365 * 24 * time.Hour // Simplified year to 365 days
	}

	// Parse and accumulate months
	if matches[2] != "" {
		months, err := strconv.Atoi(matches[2])
		if err != nil {
			return 0, fmt.Errorf("invalid months: %v", err)
		}
		totalDuration += time.Duration(months) * 30 * 24 * time.Hour // Simplified month to 30 days
	}

	// Parse and accumulate days
	if matches[3] != "" {
		days, err := strconv.Atoi(matches[3])
		if err != nil {
			return 0, fmt.Errorf("invalid days: %v", err)
		}
		totalDuration += time.Duration(days) * 24 * time.Hour
	}

	// Parse and accumulate hours
	if matches[4] != "" {
		hours, err := strconv.Atoi(matches[4])
		if err != nil {
			return 0, fmt.Errorf("invalid hours: %v", err)
		}
		totalDuration += time.Duration(hours) * time.Hour
	}

	// Parse and accumulate minutes
	if matches[5] != "" {
		minutes, err := strconv.Atoi(matches[5])
		if err != nil {
			return 0, fmt.Errorf("invalid minutes: %v", err)
		}
		totalDuration += time.Duration(minutes) * time.Minute
	}

	// Parse and accumulate seconds
	if matches[6] != "" {
		seconds, err := strconv.Atoi(matches[6])
		if err != nil {
			return 0, fmt.Errorf("invalid seconds: %v", err)
		}
		totalDuration += time.Duration(seconds) * time.Second
	}

	return totalDuration, nil
}
