-- This Source Code Form is subject to the terms of the Mozilla Public
-- License, v. 2.0. If a copy of the MPL was not distributed with this
-- file, You can obtain one at https://mozilla.org/MPL/2.0/.

-- Drop updated table
DROP TABLE IF EXISTS documents;

-- Recreate original table structure
CREATE TABLE documents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ext_id STRING NOT NULL,
    is_system BOOLEAN NOT NULL DEFAULT 0,
    file_path TEXT NOT NULL,
    updated_at DATETIME NULL,
    deleted_at DATETIME NULL
);