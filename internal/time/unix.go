// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

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
