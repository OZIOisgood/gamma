<img src="assets/poster.png" width="100%" alt="Gamma Poster" />

# Gamma

Gamma is a distributed video processing platform (a Mux-like) designed to handle video ingestion, processing, and delivery.

<img src="assets/screen-recording.gif" width="100%" alt="Gamma Demo" />

## How to start

### Prerequisites
- Docker & Docker Compose
- Go 1.23+
- Node.js & pnpm (for the dashboard)

### Quick Start

1. **Start Infrastructure** (PostgreSQL, NATS, MinIO):
   ```bash
   make docker-up
   ```

2. **Run Migrations**:
   ```bash
   make migrate-up
   ```

3. **Run Services** (in separate terminals):
   ```bash
   make run-api
   make run-worker
   ```

4. **Start Dashboard**:
   ```bash
   make dashboard-start
   ```
   Access the dashboard at `http://localhost:4200`.

## How does it work?

## System Architecture

```mermaid
flowchart TB
    subgraph Client ["Client Side"]
        direction TB
        User["👤 User"] -- Uses --> Dashboard["💻 Dashboard<br>(Angular)"]
    end

    Dashboard -- "HTTP / WebSocket" --> API["⚙️ API Service<br>(Go)"]
    API -- "SQL" --> DB[("🐘 PostgreSQL<br>(Database)")]
    Dashboard -- "Direct Upload" --> MinIO[("🗄️ MinIO<br>(S3 Storage)")]
    MinIO -. "Events" .-> NATS["📨 NATS<br>(JetStream)"]
    
    subgraph WorkerPool ["⚡ Scalable Worker Pool"]
        direction LR
        Worker["🛠️ Worker 1<br>(Go)"]
        Worker2["🛠️ Worker 2..N<br>(Go)"]
    end

    NATS -- "Jobs" --> Worker
    NATS -.- Worker2

    Worker -- "Process / Delete" --> MinIO
    Worker -- "Update Status" --> DB
    Worker -- "Events" --> NATS
    NATS -- "Notify" --> API

     User:::user
     Dashboard:::angular
     API:::go
     Worker:::go
     Worker2:::go
     MinIO:::storage
     NATS:::messaging
     DB:::db

     style WorkerPool fill:transparent,stroke:#00bcd4,stroke-width:2px,stroke-dasharray: 5 5,color:#fff
     style Client fill:transparent,stroke:#90a4ae,stroke-width:2px,stroke-dasharray: 5 5,color:#fff
    
    classDef user fill:#37474f,stroke:#90a4ae,stroke-width:2px,color:#fff
    classDef angular fill:#880e4f,stroke:#f50057,stroke-width:2px,color:#fff
    classDef go fill:#006064,stroke:#00bcd4,stroke-width:2px,color:#fff
    classDef storage fill:#b71c1c,stroke:#ff5252,stroke-width:2px,color:#fff
    classDef messaging fill:#1b5e20,stroke:#66bb6a,stroke-width:2px,color:#fff
    classDef db fill:#1a237e,stroke:#7986cb,stroke-width:2px,color:#fff
```

## Detailed Flows

### 1. Upload & Processing Flow

```mermaid
sequenceDiagram
    actor User
    participant Dash as Dashboard
    participant API
    participant DB as PostgreSQL
    participant S3 as MinIO
    participant NATS
    participant Worker

    User->>Dash: Select Video File
    Dash->>API: POST /uploads
    API->>S3: Generate Presigned URL
    API->>DB: Create Upload (pending)
    API-->>Dash: Return UploadID, URL
    
    Dash->>S3: PUT File (Direct Upload)
    S3->>NATS: Event: s3:ObjectCreated
    
    NATS->>Worker: Consume Upload Event
    Worker->>S3: Download Original
    Worker->>Worker: Transcode (FFmpeg)
    Worker->>S3: Upload HLS Segments
    Worker->>DB: Create Asset (ready)
    Worker->>DB: Update Upload (ready)
    
    Worker->>NATS: Event: asset_processed
    NATS->>API: Consume Event
    API-->>Dash: WebSocket: asset_processed
    Dash->>User: Update UI (Ready)
```

### 2. Deletion Flow

```mermaid
sequenceDiagram
    actor User
    participant Dash as Dashboard
    participant API
    participant DB as PostgreSQL
    participant NATS
    participant Worker
    participant S3 as MinIO

    User->>Dash: Click Delete
    Dash->>API: DELETE /assets/{id}
    API->>DB: Soft Delete Asset
    API->>DB: Soft Delete Upload
    API->>NATS: Event: delete_asset
    API-->>Dash: 204 No Content
    Dash->>User: Remove from List
    
    NATS->>Worker: Consume Delete Event
    Worker->>S3: Delete Original File
    Worker->>S3: Delete HLS Folder
```

Gamma is built using a microservices architecture:

### Microservices
- **API (`cmd/api`)**: Handles HTTP requests, file uploads, and serves data to the frontend.
- **Worker (`cmd/worker`)**: Consumes jobs from NATS to process videos (transcoding, etc.) asynchronously.

### Technologies
- **Backend**: Go
- **Frontend**: Angular
- **Database**: PostgreSQL
- **Messaging**: NATS
- **Storage**: S3-compatible object storage (MinIO for local development)

## Roadmap

### Implemented
- Basic video ingestion and upload flow
- Asynchronous worker processing
- Basic Dashboard UI
- Multi-quality transcoding (ABR)

### To Do
See [ISSUES.md](ISSUES.md) for the full roadmap and todo list.