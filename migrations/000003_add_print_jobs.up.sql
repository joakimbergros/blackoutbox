-- This Source Code Form is subject to the terms of the Mozilla Public
-- License, v. 2.0. If a copy of the MPL was not distributed with this
-- file, You can obtain one at https://mozilla.org/MPL/2.0/.

-- Create print_jobs table
CREATE TABLE print_jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_id INTEGER NOT NULL,
    cups_job_id TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    submitted_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    completed_at INTEGER NULL,
    error_message TEXT,
    FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE
);

-- Create index on document_id for faster lookups
CREATE INDEX idx_print_jobs_document_id ON print_jobs(document_id);

-- Create index on status for filtering
CREATE INDEX idx_print_jobs_status ON print_jobs(status);