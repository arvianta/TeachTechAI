package utils

import (
	"fmt"
	"time"
)

const (
	timezone    = "Asia/Bangkok"
	otpCooldown = 2 * time.Minute
	otpExpiry   = 5 * time.Minute
)

func GetCurrentTime() (time.Time, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load location: %w", err)
	}
	return time.Now().In(location), nil
}

func CalculateRemainingCooldown(createdAt time.Time) (int, int, error) {
	now, err := GetCurrentTime()
	if err != nil {
		return 0, 0, err
	}

	timeSinceCreation := now.Sub(createdAt.In(now.Location()))
	remainingCooldown := otpCooldown - timeSinceCreation
	if remainingCooldown > 0 {
		remainingCooldownSeconds := int(remainingCooldown.Seconds())
		minutes := remainingCooldownSeconds / 60
		seconds := remainingCooldownSeconds % 60
		return minutes, seconds, nil
	}
	return 0, 0, nil
}

func GetExpiryTime() (time.Time, error) {
	now, err := GetCurrentTime()
	if err != nil {
		return time.Time{}, err
	}
	return now.Add(otpExpiry), nil
}
