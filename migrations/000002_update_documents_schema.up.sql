-- This Source Code Form is subject to the terms of the Mozilla Public
-- License, v. 2.0. If a copy of the MPL was not distributed with this
-- file, You can obtain one at https://mozilla.org/MPL/2.0/.

-- Drop existing table (acceptable since early development)
DROP TABLE IF EXISTS documents;

-- Create new table with updated schema
CREATE TABLE documents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    system_id STRING NOT NULL,
    file_id STRING NOT NULL,
    file_path TEXT NOT NULL,
    print_at INTEGER NULL,
    last_printed_at INTEGER NULL,
    tags JSON NULL,
    updated_at DATETIME NULL,
    deleted_at DATETIME NULL,
    UNIQUE(system_id, file_id)
);

-- Create index on system_id for faster lookups
CREATE INDEX idx_documents_system_id ON documents(system_id);

-- Create index on file_id for faster lookups
CREATE INDEX idx_documents_file_id ON documents(file_id);