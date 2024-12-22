package orderfileupload

import "time"

func NewGetDateForOrderFileFunc(baseHour int) func(now time.Time) (string, error) {
	return func(now time.Time) (string, error) {
		taipeiLocation, err := time.LoadLocation("Asia/Taipei")
		if err != nil {
			return "", err
		}

		startTime := time.Date(now.Year(), now.Month(), now.Day(), baseHour, 0, 0, 0, taipeiLocation)
		if now.Before(startTime) {
			return now.Add(-time.Hour * 24).Format(time.RFC3339), nil
		}

		return now.Format(time.DateOnly), nil
	}
}
