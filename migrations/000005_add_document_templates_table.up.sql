CREATE TABLE templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    system_id TEXT NOT NULL,
    file_id TEXT NOT NULL,
    template_path TEXT NOT NULL,
    description TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    UNIQUE(system_id, file_id, template_path)
)
