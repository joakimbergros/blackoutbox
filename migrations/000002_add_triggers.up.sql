-- This Source Code Form is subject to the terms of the Mozilla Public
-- License, v. 2.0. If a copy of the MPL was not distributed with this
-- file, You can obtain one at https://mozilla.org/MPL/2.0/.

-- Create triggers table
CREATE TABLE triggers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    system_id TEXT NOT NULL,
    url TEXT NOT NULL,
    last_failed_at INTEGER NULL,
    buffer_seconds INTEGER NOT NULL DEFAULT 300,
    status TEXT NOT NULL DEFAULT 'ok',
    last_checked_at INTEGER NULL,
    retry_count INTEGER NOT NULL DEFAULT 0,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

-- Create index on system_id for faster lookups
CREATE INDEX idx_triggers_system_id ON triggers(system_id);

-- Create index on status for filtering
CREATE INDEX idx_triggers_status ON triggers(status);