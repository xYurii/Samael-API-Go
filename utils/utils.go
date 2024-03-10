package utils

import (
	"fmt"
	"time"
)

func ParseDuration(hour int) int64 {
	return int64((time.Duration(hour) * time.Hour) / time.Millisecond)
}

func InCommandCooldown(lastUsageTime int64, cooldownHours int) bool {
	cooldown := time.Duration(cooldownHours) * time.Hour / time.Millisecond
	currentTimeMilliseconds := time.Now().UnixNano() / int64(time.Millisecond)

	return currentTimeMilliseconds-lastUsageTime <= int64(cooldown)
}

func FormatTime(duration time.Duration, maxUnits int) string {
	days := int(duration.Hours() / 24)
	duration -= time.Duration(days) * 24 * time.Hour
	hours := int(duration.Hours())
	duration -= time.Duration(hours) * time.Hour
	minutes := int(duration.Minutes())
	duration -= time.Duration(minutes) * time.Minute
	seconds := int(duration.Seconds())

	units := []struct {
		value int
		label string
	}{
		{days, "dia"},
		{hours, "hora"},
		{minutes, "minuto"},
		{seconds, "segundo"},
	}

	result := ""
	count := 0
	for _, unit := range units {
		if unit.value > 0 && count < maxUnits {
			if result != "" {
				if count > 0 {
					if count == maxUnits-1 {
						result += " e "
					} else {
						result += ", "
					}
				}
			}
			result += fmt.Sprintf("%d %s", unit.value, unit.label)
			if unit.value != 1 {
				result += "s" // plural
			}
			count++
		}
	}

	return result
}
