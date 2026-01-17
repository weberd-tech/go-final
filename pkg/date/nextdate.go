package date

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("empty repeat rule")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}

	parts := strings.Split(strings.TrimSpace(repeat), " ")
	if len(parts) == 0 {
		return "", errors.New("invalid repeat format")
	}

	rule := parts[0]

	switch rule {
	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if date.Format(dateFormat) > now.Format(dateFormat) {
				break
			}
		}
		return date.Format(dateFormat), nil

	case "d":
		if len(parts) < 2 {
			return "", errors.New("invalid d rule: missing number")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid number of days: %w", err)
		}

		if days > 400 {
			return "", errors.New("days exceed maximum (400)")
		}

		if days <= 0 {
			return "", errors.New("days must be positive")
		}

		for {
			date = date.AddDate(0, 0, days)
			if date.Format(dateFormat) > now.Format(dateFormat) {
				break
			}
		}
		return date.Format(dateFormat), nil

	case "w", "m":
		return "", errors.New("unsupported format")

	default:
		return "", errors.New("unknown repeat rule")
	}
}
