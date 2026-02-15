// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package models

type Template struct {
	Id          int    `json:"id"`
	SystemId    string `json:"system_id"`
	FileId      string `json:"file_id"`
	FilePath    string `json:"file_path"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	DeletedAt   *int64 `json:"deleted_at"`
}
