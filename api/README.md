# Aletheia Scalable Backend (Microservices + Celery)

This folder contains a microservices-ready backend architecture for Aletheia using:

- Go + Gin for external/internal APIs
- Celery for async processing
- Redis for Celery broker/backend and verification cache
- Cassandra for logs/analytics
- BigchainDB for metadata
- IPFS for image storage
- Polygon anchoring pipeline

## Services

- `services/api-gateway` (Gin, API key + rate limiting)
- `services/upload-service` (Gin, receives image and enqueues Celery task)
- `services/verification-service` (Gin + Redis cache)
- `services/hash-service` (internal hashing endpoint)
- `services/merkle-service` (internal merkle endpoint)
- `services/trust-engine` (internal trust scoring endpoint)

## Async Worker Layer (Celery)

- `workers/celery-worker/task_api.py` (FastAPI task enqueue/status API)
- `workers/celery-worker/tasks.py` (Celery tasks for hash, pHash, merkle, IPFS, BigchainDB, Polygon placeholder)
- `workers/celery-worker/celery_app.py` (Celery app config)

## External API

- `POST /api/v1/images/upload` (Gateway -> Upload Service -> Celery)
- `POST /api/v1/verify` (Gateway -> Verification Service)
- `POST /api/v1/verify/url` (single-service API under `cmd/api`, URL-based verification)
- `POST /api/v1/video/upload` (single-service API under `cmd/api`)
- `POST /api/v1/video/verify` (single-service API under `cmd/api`)

## Swagger / OpenAPI

- OpenAPI spec file: `docs/openapi.yaml`
- Covers all single-service endpoints under `/api/v1/*`

Quick preview with Swagger UI (Docker):

```bash
docker run --rm -p 8088:8080 -e SWAGGER_JSON=/spec/openapi.yaml -v ${PWD}/docs:/spec swaggerapi/swagger-ui
```

Then open `http://localhost:8088`.

## Trust Infrastructure Enhancements

- Tamper-proof timestamping
	- proof record stores `hash`, anchored `timestamp`, and `blockNumber`
	- verification validates timestamp consistency and returns `timestampValid`
- URL-based verification
	- new endpoint: `POST /api/v1/verify/url`
	- payload: `{ "imageUrl": "https://..." }`
	- service fetches bytes with size/content-type checks, then runs normal verify flow
- Context awareness
	- upload metadata supports `captureTimestamp`, `location` (optional), and `deviceInfo` (optional)
	- verification detects context mismatch and emits warnings/suspicious flags
- Fraud and abuse detection
	- duplicate upload detection by exact SHA-256
	- upload spam detection per uploader/device identity window
	- replay detection for same media with altered metadata
- Proof report generation
	- verify response now embeds a complete `proofReport` certificate object
	- includes proof id, hash, CID, timestamp, block number, verdict, trust score, breakdown, and flags

## Verify Response (Current)

```json
{
	"verdict": "AUTHENTIC | TAMPERED | SIMILAR | UNKNOWN",
	"trustScore": 92,
	"timestampValid": true,
	"sourceType": "UPLOAD | URL",
	"proofReport": {
		"proofId": "...",
		"sha256": "...",
		"CID": "...",
		"timestamp": "2026-03-23T10:20:30Z",
		"blockNumber": 123456,
		"verdict": "AUTHENTIC",
		"trustScore": 92,
		"breakdown": {
			"hash": "MATCH",
			"similarity": 0.91,
			"signature": "VALID"
		},
		"flags": {
			"duplicate": false,
			"spam": false,
			"replay": false,
			"suspicious": false
		}
	},
	"warnings": []
}
```

## Video Authenticity Pipeline

- Accept video upload (`multipart/form-data`, field: `video`)
- Extract audio via FFmpeg (`ffmpeg -i input -q:a 0 -map a output.mp3`)
- Compute SHA-256 for video and extracted audio
- Upload both to IPFS
- Store metadata in BigchainDB/in-memory repository
- Verify tampering verdicts:
	- `AUTHENTIC`
	- `AUDIO_TAMPERED`
	- `VIDEO_TAMPERED`
	- `TAMPERED`

## Trust Score

- SHA match: `+35`
- pHash similarity: `+20`
- Metadata integrity: `+10`
- Witness placeholder: `+15`

## Infrastructure

```bash
docker compose up -d
```

This starts Redis, Cassandra, IPFS, Celery Task API, Celery Worker, and Flower.

## Run Go Services

```bash
go run ./services/api-gateway
go run ./services/upload-service
go run ./services/verification-service
```

## Check Celery Worker

1. Flower dashboard:
	- `http://localhost:5555`
2. Task API status endpoint:
	- `GET http://localhost:8090/tasks/{taskId}`
3. Docker logs:
	- `docker compose logs -f celery-worker`

## Notes

- Kafka has been removed from this architecture.
- Legacy single-service MVP under `cmd/` and `internal/` remains for backward compatibility.
- FFmpeg must be installed and available in PATH for video endpoints.
