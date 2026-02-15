package time

import "time"

// UnixToTime converts Unix timestamp to time.Time
func UnixToTime(unixTimestamp *int64) *time.Time {
	if unixTimestamp == nil {
		return nil
	}
	t := time.Unix(*unixTimestamp, 0)
	return &t
}

// TimeToUnix converts time.Time to Unix timestamp
func TimeToUnix(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	unix := t.Unix()
	return &unix
}

// UnixToTimeNonPtr converts Unix timestamp to time.Time (non-pointer version)
func UnixToTimeNonPtr(unixTimestamp int64) time.Time {
	return time.Unix(unixTimestamp, 0)
}

// TimeToUnixNonPtr converts time.Time to Unix timestamp (non-pointer version)
func TimeToUnixNonPtr(t time.Time) int64 {
	return t.Unix()
}
