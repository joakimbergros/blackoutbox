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

### Query Parameters

- `system-id` - Filter documents by system identifier
- `file-id` - Filter documents by file identifier

### Request/Response Format

All endpoints use JSON for request and response bodies.

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
  "updated_at": "2026-02-04T10:00:00Z",
  "deleted_at": null
}
```

## 🛠️ Getting Started

### Prerequisites

- Go 1.25.6 or higher
- SQLite3
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

## 🗄️ Database Schema

```sql
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
```

Indexes are created on `system_id` and `file_id` for fast lookups.

## 🚧 Roadmap

### Planned Features

- [ ] **Template Documents**: Printable forms for manual data entry
- [ ] **Scanning Integration**: Scan filled forms back into the system
- [ ] **LLM-Based Parsing**: Extract handwritten information using AI
- [ ] **Export to Source Systems**: Sync parsed data back to primary systems
- [ ] **Print Queue Management**: Better control over emergency printing
- [ ] **Web Interface**: User-friendly UI for document management
- [ ] **Authentication & Authorization**: Secure access control
- [ ] **Backup & Recovery**: Automated backup strategies
- [ ] **Monitoring & Alerts**: System health monitoring

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

[Specify your license here]

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
