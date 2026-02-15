-- This Source Code Form is subject to the terms of the Mozilla Public
-- License, v. 2.0. If a copy of the MPL was not distributed with this
-- file, You can obtain one at https://mozilla.org/MPL/2.0/.

-- Create systems table
CREATE TABLE systems (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    reference TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    deleted_at INTEGER NULL,
    CHECK (reference != ''),
    CHECK (name != '')
);

-- Create index on system ID for faster lookups
CREATE INDEX idx_systems_id ON systems(id);

-- Foreign key constraints have been added directly to the CREATE TABLE statements
-- in the respective migration files (000001_initial.up.sql, 000002_add_triggers.up.sql,
-- and 000004_add_document_templates_table.up.sql) since this is a development environment
-- and migrations will be rerun.