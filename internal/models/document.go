// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package models

type Document struct {
	Id            int64    `json:"id"`
	SystemId      int64    `json:"system_id"`
	FileReference string   `json:"file_id"`
	FilePath      string   `json:"file_path"`
	PrintAt       *int64   `json:"print_at"`
	LastPrintedAt *int64   `json:"last_printed_at"`
	Tags          []string `json:"tags"`
	UpdatedAt     *int64   `json:"updated_at"`
	DeletedAt     *int64   `json:"deleted_at"`
}
