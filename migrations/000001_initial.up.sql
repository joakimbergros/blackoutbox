-- This Source Code Form is subject to the terms of the Mozilla Public
-- License, v. 2.0. If a copy of the MPL was not distributed with this
-- file, You can obtain one at https://mozilla.org/MPL/2.0/.

CREATE TABLE documents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    system_id INTEGER NOT NULL,
    file_id TEXT NOT NULL,
    file_path TEXT NOT NULL,
    print_at INTEGER NULL,
    last_printed_at INTEGER NULL,
    tags TEXT NULL,
    updated_at INTEGER NULL,
    deleted_at INTEGER NULL,
    UNIQUE(system_id, file_id),
    FOREIGN KEY (system_id) REFERENCES systems(id) ON DELETE CASCADE
);

CREATE INDEX idx_documents_system_id ON documents(system_id);
CREATE INDEX idx_documents_file_id ON documents(file_id);
