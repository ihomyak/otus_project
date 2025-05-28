package services

import "time"

const RateLimitContextKey string = "RateLimitKey"

type RateLimiter struct {
	TokenID  string
	Limit    int
	Duration time.Duration
}

func ConvertTextToTime(text string) time.Duration {
	var timeDuration time.Duration
	switch text {
	case "hour":
		timeDuration = time.Hour
	case "minute":
		timeDuration = time.Minute
	case "second":
		timeDuration = time.Second
	default:
		// todo: custom error
		panic("Unknown time format")
	}
	return timeDuration
}
