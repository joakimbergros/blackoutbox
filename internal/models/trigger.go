// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package models

type Trigger struct {
	Id            int    `json:"id"`
	SystemId      string `json:"system_id"`
	Url           string `json:"url"`
	LastFailedAt  *int64 `json:"last_failed_at"`
	BufferSeconds int    `json:"buffer_seconds"`
	Status        string `json:"status"` // ok, error, triggered
	LastCheckedAt *int64 `json:"last_checked_at"`
	RetryCount    int    `json:"retry_count"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}
