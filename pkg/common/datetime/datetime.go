package datetime

import "time"

func ToUTC(time *time.Time) {
	if time != nil {
		*time = time.UTC()
		return
	}

	return
}
