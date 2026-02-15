# BlackoutBox

A resilient, self-contained document management system designed for emergency scenarios where digital infrastructure becomes unavailable. Built for public sector organizations, particularly elderly care departments, to ensure critical information remains accessible during internet outages, cyber attacks, or infrastructure failures.

## 🎯 Vision

BlackoutBox serves as a digital "black box" that stores essential documents and automatically prints them when normal systems fail. Think of it as an emergency backup system that ensures continuity of care and access to vital information when "shit hits the fan."

## 🚀 Current Status

**Hackathon Project - Proof of Concept**

This is an early-stage implementation developed during a hackathon to showcase the potential of offline-first emergency document management. The core functionality is implemented and tested, but production deployment is not yet complete.

## ✨ Key Features

- **Offline-First Architecture**: Complete independence from internet connectivity
- **Automatic Emergency Printing**: Trigger document printing based on scheduled timestamps
- **Health Check Monitoring**: Monitor external system health via HTTP endpoints
- **Automatic Trigger-Based Printing**: Print documents when systems fail for extended periods
- **Print Job Tracking**: Monitor CUPS print jobs with status tracking
- **Stuck Job Detection**: Alert on print jobs that have been pending too long
- **Multi-System Support**: Organize documents by system (e.g., different care facilities, departments)
- **Tag-Based Organization**: Categorize documents for quick retrieval
- **Soft Delete**: Preserve document history with deletion tracking
- **RESTful API**: Simple, standard HTTP interface for integration
- **Self-Contained Deployment**: Runs on minimal hardware

## 🏗️ Architecture

### Tech Stack

- **Go 1.25.6** - Lightweight, efficient backend
- **SQLite3** - Embedded database, no external dependencies
- **Standard Library** - Minimal external dependencies for reliability

### Deployment Targets

- **Primary**: Raspberry Pi with PCI RAID storage + wired printer
- **Alternative**: Decommissioned laptops

The system is designed to be completely self-contained with minimal resource requirements.

## 📋 API Endpoints

### Documents

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/documents` | List all documents or filter by `system-id` or `file-id` |
| `GET` | `/documents/{id}` | Get a specific document by ID |
| `POST` | `/documents` | Upload a new document |
| `PATCH` | `/documents` | Update a document (placeholder) |

### Systems

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/systems` | List all systems |
| `GET` | `/systems/{id}` | Get a specific system by ID |
| `POST` | `/systems` | Create a new system |
| `PUT` | `/systems/{id}` | Update a system |
| `DELETE` | `/systems/{id}` | Delete a system (soft delete) |

### System CRUD Examples

**Create a System:**
```bash
curl -X POST http://localhost:3000/systems \
  -H "Content-Type: application/json" \
  -d '{
    "id": "care-facility-1",
    "name": "Elderly Care Facility Alpha",
    "description": "Primary care facility in downtown"
  }'
```

**Get All Systems:**
```bash
curl -X GET http://localhost:3000/systems
```

**Get Specific System:**
```bash
curl -X GET http://localhost:3000/systems/care-facility-1
```

**Update a System:**
```bash
curl -X PUT http://localhost:3000/systems/care-facility-1 \
  -H "Content-Type: application/json" \
  -d '{
    "id": "care-facility-1",
    "name": "Elderly Care Facility Alpha - Updated",
    "description": "Primary care facility with expanded capacity"
  }'
```

**Delete a System:**
```bash
curl -X DELETE http://localhost:3000/systems/care-facility-1
```

### System Bulk Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/systems/{id}/sync` | Mirror system storage state with request data |

### Templates

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/templates` | List all documents or filter by `system-id` or `file-id` |
| `POST` | `/templates` | Upload a new document |
| `DELETE` | `/templates` | Remove a document |

### Triggers

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/triggers` | List all health check triggers |
| `GET` | `/triggers/{id}` | Get a specific trigger by ID |
| `POST` | `/triggers` | Create a new health check trigger |
| `DELETE` | `/triggers/{id}` | Delete a trigger |

### Print Jobs

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/print_jobs` | List all print jobs |
| `GET` | `/print_jobs/{id}` | Get a specific print job by ID |
| `GET` | `/print_jobs/stuck` | Get stuck print jobs (>5 min) |

### Query Parameters

- `system-id` - Filter documents by system identifier
- `file-id` - Filter documents by file identifier
- `threshold` - For `/print_jobs/stuck`, time in seconds (default: 300)

### Request/Response Format

All endpoints use JSON for request and response bodies.

**Timestamp Format:**
All timestamps are represented as Unix epoch integers (seconds since 1970-01-01 00:00:00 UTC).

- Example: `1738581234` represents February 4, 2026 10:00:00 AM UTC
- Null values indicate the timestamp is not set
- Use `Date.now() / 1000` in JavaScript or `time.Now().Unix()` in Go to generate current timestamps

**Document Model:**
```json
{
  "id": 1,
  "system_id": "care-facility-1",
  "file_id": "emergency-protocol-001",
  "file_path": "uploads/care-facility-1/1738581234_protocol.pdf",
  "print_at": 1738581234,
  "last_printed_at": null,
  "tags": ["emergency", "protocol", "high-priority"],
  "updated_at": 1738581234,
  "deleted_at": null
}
```

**Trigger Model:**
```json
{
  "id": 1,
  "system_id": "care-facility-1",
  "url": "https://api.example.com/health",
  "last_failed_at": null,
  "buffer_seconds": 300,
  "status": "ok",
  "last_checked_at": 1738581234,
  "retry_count": 0,
  "created_at": 1738581234,
  "updated_at": 1738581234
}
```

**Print Job Model:**
```json
{
  "id": 1,
  "document_id": 1,
  "cups_job_id": "123",
  "status": "printing",
  "submitted_at": 1738581234,
  "completed_at": null,
  "error_message": null
}
```

**System Model:**
```json
{
  "id": "care-facility-1",
  "name": "Elderly Care Facility Alpha",
  "description": "Primary care facility in downtown",
  "created_at": 1738581234,
  "updated_at": 1738581234,
  "deleted_at": null
}
```

## 🛠️ Getting Started

### Prerequisites

- Go 1.25.6 or higher
- SQLite3
- CUPS (Common Unix Printing System) for printing functionality
- (Optional) `migrate` CLI tool for database migrations

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd blackoutbox

# Install dependencies
go mod download

# Run database migrations
just migrate-up

# Start the server
go run main.go
```

The server will start on `http://localhost:3000`

### Testing

```bash
# Run all tests
just test

# Run tests with coverage
just test-coverage
```

### Available Commands

```bash
just migrate-up      # Apply database migrations
just migrate-down    # Rollback database migrations
just test            # Run all tests
just test-coverage   # Run tests with coverage report
```

## 📄 Document Upload

Upload documents using multipart form data:

```bash
curl -X POST http://localhost:3000/documents \
  -F "system_id=care-facility-1" \
  -F "file_id=emergency-protocol-001" \
  -F "file=@protocol.pdf" \
  -F "tags=[\"emergency\",\"protocol\"]" \
  -F "print_at=1738581234"
```

### Required Fields

- `system_id` - System/department identifier
- `file_id` - Unique file identifier
- `file` - The document file (max 10MB)

### Optional Fields

- `tags` - JSON array of tags for categorization
- `print_at` - Unix timestamp for automatic printing

## 🎯 Health Check Triggers

Create health check triggers to monitor external systems and automatically print documents when they fail:

```bash
curl -X POST http://localhost:3000/triggers \
  -H "Content-Type: application/json" \
  -d '{
    "system_id": "care-facility-1",
    "url": "https://api.example.com/health",
    "buffer_seconds": 300
  }'
```

### How It Works

1. **Health Checks**: The background worker checks trigger URLs every 30 seconds
2. **Nagios-Style Logic**:
   - OK (200-299): Reset retry count
   - Error (400+ or timeout): Increment retry count
   - After 3 consecutive failures + buffer time: Trigger print jobs
3. **Automatic Printing**: All documents associated with the system_id are printed
4. **Status Tracking**: Triggers have statuses: `ok`, `error`, `triggered`

### Trigger Fields

- `system_id` (required) - System/department identifier to monitor
- `url` (required) - Health check endpoint URL
- `buffer_seconds` (optional) - Time to wait before triggering (default: 300)

## 🗄️ Database Schema

### Systems Table

```sql
CREATE TABLE systems (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    deleted_at INTEGER NULL,
    CHECK (id != ''),
    CHECK (name != '')
);
```

### Documents Table

```sql
CREATE TABLE documents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    system_id TEXT NOT NULL,
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
```

### Templates table
```sql
CREATE TABLE templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    system_id TEXT NOT NULL,
    file_id TEXT NOT NULL,
    template_path TEXT NOT NULL,
    description TEXT,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    deleted_at INTEGER NULL,
    UNIQUE(system_id, file_id),
    FOREIGN KEY (system_id) REFERENCES systems(id) ON DELETE CASCADE
)
```

### Triggers Table

```sql
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
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (system_id) REFERENCES systems(id) ON DELETE CASCADE
);
```

### Print Jobs Table

```sql
CREATE TABLE print_jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_id INTEGER NOT NULL,
    cups_job_id TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    submitted_at INTEGER NOT NULL,
    completed_at INTEGER,
    error_message TEXT,
    FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE
);
```

Indexes are created on `system_id` and `file_id` for fast lookups.

## 🚧 Roadmap

### Planned Features

- [x] **Template Documents**: Printable forms for manual data entry
- [ ] **Scanning Integration**: Scan filled forms back into the system
- [ ] **LLM-Based Parsing**: Extract handwritten information using AI
- [ ] **Export to Source Systems**: Sync parsed data back to primary systems
- [ ] **Print Queue Management**: Better control over emergency printing
- [ ] **Web Interface**: User-friendly UI for document management
- [ ] **Authentication & Authorization**: Secure access control
- [ ] **Backup & Recovery**: Automated backup strategies
- [~] **Monitoring & Alerts**: System health monitoring (partially implemented - health checks with printing, webhooks planned)

### Future Enhancements

- Support for additional document formats
- Multi-language support
- Advanced search and filtering
- Document versioning
- Integration with existing care management systems

## 🏥 Use Case: Elderly Care

In elderly care facilities, critical information must remain accessible during emergencies:

- **Emergency Protocols**: Step-by-step procedures for medical emergencies
- **Patient Information**: Essential medical records and care instructions
- **Contact Lists**: Emergency contacts and staff directories
- **Medication Schedules**: Critical medication administration guides
- **Facility Maps**: Evacuation routes and safe zones

When infrastructure fails, BlackoutBox automatically prints these documents, ensuring staff can continue providing care without interruption.

## 🤝 Contributing

This is a hackathon project, but contributions are welcome! Feel free to:

- Report bugs
- Suggest new features
- Submit pull requests
- Improve documentation

## 📝 License

MPL-2.0

## 👥 Team

Developed during [KLIRR-hack 3–4 February 2026](https://www.klirr-hack.se) by:
- Joakim Bergros ([@joakimbergros](https://github.com/joakimbergros))
- Ammar Kasem ([@Ammar-Kasem](https://github.com/Ammar-Kasem))
- Gustav Fröjdlund ([@gustavfrojdlund](https://github.com/gustavfrojdlund))

## 🙏 Acknowledgments

- Built for elderly care departments to ensure continuity of care
- Inspired by the need for resilient infrastructure in critical public services

---

**Note**: This is a proof-of-concept implementation. Production deployment requires additional security hardening, testing, and infrastructure planning.
