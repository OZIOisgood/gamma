<img src="assets/poster.png" width="100%" alt="Gamma Poster" />

# Gamma

Gamma is a distributed video processing platform (a Mux-like) designed to handle video ingestion, processing, and delivery.

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

### To Do
- [ ] Delete assets
- [ ] Search (Elasticsearch)
- [ ] i18n for dashboard
- [ ] Security improvements
- [ ] Quality selection
- [ ] Subtitles selection
- [ ] Documentation
- [ ] CONTRIBUTING file
- [ ] Client player
- [ ] Preview GIFs
- [ ] Thumbnails
- [ ] Environment selection in dashboard (dev, qa, prod)
- [ ] SDKs for popular languages
- [ ] CORS configuration
- [ ] Quality optimization on asset creation
- [ ] Webhooks
- [ ] DRM
- [ ] Test flag on asset (watermark + 10s limit + auto-delete after 24h)
- [ ] Metadata support: title, creator_id, external_id