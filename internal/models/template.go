// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package models

import "time"

type Template struct {
	Id          int        `json:"id"`
	SystemId    string     `json:"system_id"`
	FileId      string     `json:"file_id"`
	FilePath    string     `json:"file_path"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
