# BlackoutBox

Ett robust, fristående dokumenthanteringssystem designat för nödsituationer där digital infrastruktur blir otillgänglig. Byggt för offentlig sektor, särskilt äldreomsorgsverksamheter, för att säkerställa att kritisk information förblir tillgänglig vid internetavbrott, cyberattacker eller infrastrukturfel.

## 🎯 Vision

BlackoutBox fungerar som en digital "svart låda" som lagrar viktiga dokument och automatiskt skriver ut dem när normalsystemen fallerar. Tänk på det som ett reservsystem som säkerställer kontinuitet i vården och tillgång till vital information när "skiten träffar fläkten".

## 🚀 Nuvarande status

**Hackathon-projekt - Konceptbevis**

Detta är ett tidigt skede av implementation utvecklad under en hackathon för att visa potentialen i offline-först akut dokumenthantering. Kärnfunktionaliteten är implementerad och testad, men produktionssättning är inte än komplett.

## ✨ Nyckelfunktioner

- **Offline-först-arkitektur**: Fullständigt oberoende av internetanslutning
- **Automatisk akut utskrift**: Utlös dokumentutskrift baserat på schemalagda tidsstämplar
- **Hälsokontrollövervakning**: Övervaka externa systems hälsa via HTTP-slutpunkter
- **Automatisk utlösarbaserad utskrift**: Skriv ut dokument när system misslyckas under längre perioder
- **Utskriftsjobbspårning**: Övervaka CUPS-utskriftsjobb med statusspårning
- **Detektion av fastnande jobb**: Avisera om utskriftsjobb som har väntat för länge
- **Stöd för flera system**: Organisera dokument efter system (t.ex. olika vårdinrättningar, avdelningar)
- **Taggbaserad organisering**: Kategorisera dokument för snabb hämtning
- **Mjuk borttagning**: Bevara dokumenthistorik med borttagningshantering
- **RESTful API**: Enkelt, standardiserat HTTP-gränssnitt för integration
- **Fristående distribution**: Körs på minimal hårdvara

## 🏗️ Arkitektur

### Teknikstack

- **Go 1.25.6** - Lättviktig, effektiv backend
- **SQLite3** - Inbyggd databas, inga externa beroenden
- **Standardbibliotek** - Minimala externa beroenden för tillförlitlighet

### Distribution

- **Primärt**: Raspberry Pi med PCI RAID-lagring + trådbunden skrivare
- **Alternativ**: Avskaffade bärbara datorer

Systemet är designat för att vara helt fristående med minimala resurskrav.

## 📋 API-slutpunkter

### Dokument

| Metod | Slutpunkt | Beskrivning |
|--------|-----------|-------------|
| `GET` | `/documents` | Lista alla dokument eller filtrera efter `system-id` eller `file-id` |
| `GET` | `/documents/{id}` | Hämta ett specifikt dokument efter ID |
| `POST` | `/documents` | Ladda upp ett nytt dokument |
| `PATCH` | `/documents` | Uppdatera ett dokument (platshållare) |

### System

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/systems/{id}/sync` | Spegla lagring på enhet med data från förfrågan |
| `DELETE` | `/systems/{id}` | Ta bort alla dokument relaterat till systemet |

### Templates

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/templates` | List all documents or filter by `system-id` or `file-id` |
| `POST` | `/templates` | Upload a new document |
| `DELETE` | `/templates` | Remove a document |

### Utlösare (Triggers)

| Metod | Slutpunkt | Beskrivning |
|--------|-----------|-------------|
| `GET` | `/triggers` | Lista alla hälsokontrollutlösare |
| `GET` | `/triggers/{id}` | Hämta en specifik utlösare efter ID |
| `POST` | `/triggers` | Skapa en ny hälsokontrollutlösare |
| `DELETE` | `/triggers/{id}` | Ta bort en utlösare |

### Utskriftsjobb

| Metod | Slutpunkt | Beskrivning |
|--------|-----------|-------------|
| `GET` | `/print_jobs` | Lista alla utskriftsjobb |
| `GET` | `/print_jobs/{id}` | Hämta ett specifikt utskriftsjobb efter ID |
| `GET` | `/print_jobs/stuck` | Hämta fastnade utskriftsjobb (>5 min) |

### Frågeparametrar

- `system-id` - Filtrera dokument efter systemidentifierare
- `file-id` - Filtrera dokument efter filidentifierare
- `threshold` - För `/print_jobs/stuck`, tid i sekunder (standard: 300)

### Format för begäran/svar

Alla slutpunkter använder JSON för begäran- och svarstexter.

**Dokumentmodell:**
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

**Utlösarmodell (Trigger Model):**
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
  "created_at": "2026-02-04T10:00:00Z",
  "updated_at": "2026-02-04T10:00:00Z"
}
```

**Utskriftsjobbsmodell (Print Job Model):**
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

## 🛠️ Kom igång

### Förutsättningar

- Go 1.25.6 eller högre
- SQLite3
- CUPS (Common Unix Printing System) för utskriftsfunktionalitet
- (Valfritt) `migrate` CLI-verktyg för databasmigreringar

### Installation

```bash
# Klona arkivet
git clone <repository-url>
cd blackoutbox

# Installera beroenden
go mod download

# Kör databasmigreringar
just migrate-up

# Starta servern
go run main.go
```

Servern startar på `http://localhost:3000`

### Testning

```bash
# Kör alla tester
just test

# Kör tester med täckning
just test-coverage
```

### Tillgängliga kommandon

```bash
just migrate-up      # Tillämpa databasmigreringar
just migrate-down    # Återställ databasmigreringar
just test            # Kör alla tester
just test-coverage   # Kör tester med täckningsrapport
```

## 📄 Dokumentuppladdning

Ladda upp dokument med multipart-formulärdata:

```bash
curl -X POST http://localhost:3000/documents \
  -F "system_id=care-facility-1" \
  -F "file_id=emergency-protocol-001" \
  -F "file=@protocol.pdf" \
  -F "tags=[\"emergency\",\"protocol\"]" \
  -F "print_at=1738581234"
```

### Obligatoriska fält

- `system_id` - System-/avdelningsidentifierare
- `file_id` - Unik filidentifierare
- `file` - Dokumentfilen (max 10MB)

### Valfria fält

- `tags` - JSON-array med taggar för kategorisering
- `print_at` - Unix-tidsstämpel för automatisk utskrift

## 🎯 Hälsokontrollutlösare

Skapa hälsokontrollutlösare för att övervaka externa system och automatiskt skriva ut dokument när de misslyckas:

```bash
curl -X POST http://localhost:3000/triggers \
  -H "Content-Type: application/json" \
  -d '{
    "system_id": "care-facility-1",
    "url": "https://api.example.com/health",
    "buffer_seconds": 300
  }'
```

### Hur det fungerar

1. **Hälsokontroller**: Bakgrundsarbetaren kontrollerar utlösar-URL:er var 30:e sekund
2. **Nagios-stil logik**:
   - OK (200-299): Återställ antal försök
   - Fel (400+ eller timeout): Öka antal försök
   - Efter 3 på varandra följande misslyckanden + buffertid: Utlös utskriftsjobb
3. **Automatisk utskrift**: Alla dokument kopplade till system_id skrivs ut
4. **Statusspårning**: Utlösare har statusar: `ok`, `error`, `triggered`

### Utlösarfält

- `system_id` (obligatoriskt) - System-/avdelningsidentifierare att övervaka
- `url` (obligatoriskt) - Hälsokontrollslutpunktens URL
- `buffer_seconds` (valfritt) - Tid att vänta innan utlösning (standard: 300)

## 🗄️ Databasschema

### Dokumenttabell (Documents Table)

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

### Malltabell
```sql
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
```

### Utlösartabell (Triggers Table)

```sql
CREATE TABLE triggers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    system_id TEXT NOT NULL,
    url TEXT NOT NULL,
    last_failed_at INTEGER,
    buffer_seconds INTEGER NOT NULL DEFAULT 300,
    status TEXT NOT NULL DEFAULT 'ok',
    last_checked_at INTEGER,
    retry_count INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Utskriftsjobbstabell (Print Jobs Table)

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

Index skapas på `system_id` och `file_id` för snabba uppslagningar.

## 🚧 Färdplan

### Planerade funktioner

- [x] **Maldokument**: Utskrivbara formulär för manuell datainmatning
- [ ] **Skannerintegration**: Skanna ifyllda formulär tillbaka till systemet
- [ ] **LLM-baserad tolkning**: Extrahera handskriven information med hjälp av AI
- [ ] **Export till källsystem**: Synkronisera tolkad data tillbaka till primära system
- [ ] **Utskriftsköhantering**: Bättre kontroll över akuta utskrifter
- [ ] **Webbgränssnitt**: Användarvänligt UI för dokumenthantering
- [ ] **Autentisering och auktorisering**: Säker åtkomstkontroll
- [ ] **Säkerhetskopiering och återställning**: Automatiserade säkerhetskopieringsstrategier
- [~] **Övervakning och aviseringar**: Hälsomonitorering av systemet (delvis implementerad - hälsokontroller med utskrift, webhooks planerade)

### Framtida förbättringar

- Stöd för ytterligare dokumentformat
- Flerspråksstöd
- Avancerad sökning och filtrering
- Dokumentversionshantering
- Integration med befintliga vårdledningssystem

## 🏥 Användningsfall: Äldreomsorg

På äldreomsorgsanläggningar måste kritisk information förbli tillgänglig vid nödsituationer:

- **Nödprotokoll**: Steg-för-steg-procedurer för medicinska nödsituationer
- **Patientinformation**: Essentiella journaler och vårdinstruktioner
- **Kontaktlistor**: Nödkontakter och personalregister
- **Medicineringsscheman**: Kritiska medicinadministrationsguider
- **Anläggningskartor**: Evakueringsvägar och säkra zoner

När infrastrukturen fallerar skriver BlackoutBox automatiskt ut dessa dokument, vilket säkerställer att personalen kan fortsätta vården utan avbrott.

## 🤝 Bidrag

Detta är ett hackathon-projekt, men bidrag är välkomna! Du är välkommen att:

- Rapportera buggar
- Föreslå nya funktioner
- Skicka pull requests
- Förbättra dokumentationen

## 📝 Licens

MPL2

## 👥 Team

Utvecklat under [KLIRR-hack 3–4 februari 2026](https://www.klirr-hack.se) av:
- Joakim Bergros ([@joakimbergros](https://github.com/joakimbergros))
- Ammar Kasem ([@Ammar-Kasem](https://github.com/Ammar-Kasem))
- Gustav Fröjdlund ([@gustavfrojdlund](https://github.com/gustavfrojdlund))

## 🙏 Tack

- Byggt för äldreomsorgsverksamheter för att säkerställa kontinuitet i vården
- Inspirerat av behovet av resilient infrastruktur i kritiska offentliga tjänster

---

**Obs**: Detta är en proof-of-concept-implementation. Produktionssättning kräver ytterligare säkerhetshärdning, testning och infrastrukturplanering.
