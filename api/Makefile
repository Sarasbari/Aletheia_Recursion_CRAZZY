SHELL := /bin/sh

.PHONY: infra gateway upload verify hash merkle trust storage blockchain task-api celery-worker flower celery-status

infra:
	docker compose up -d

gateway:
	SERVICE_NAME=api-gateway PORT=8080 go run ./services/api-gateway

upload:
	SERVICE_NAME=upload-service PORT=8081 go run ./services/upload-service

verify:
	SERVICE_NAME=verification-service PORT=8082 go run ./services/verification-service

hash:
	SERVICE_NAME=hash-service PORT=8083 go run ./services/hash-service

merkle:
	SERVICE_NAME=merkle-service PORT=8084 go run ./services/merkle-service

trust:
	SERVICE_NAME=trust-engine PORT=8085 go run ./services/trust-engine

storage:
	SERVICE_NAME=storage-service PORT=8086 go run ./services/storage-service

blockchain:
	SERVICE_NAME=blockchain-service PORT=8087 go run ./services/blockchain-service

task-api:
	cd workers/celery-worker && uvicorn task_api:app --host 0.0.0.0 --port 8090

celery-worker:
	cd workers/celery-worker && celery -A celery_app.celery_app worker --loglevel=info

flower:
	cd workers/celery-worker && celery -A celery_app.celery_app flower --port=5555

celery-status:
	cd workers/celery-worker && celery -A celery_app.celery_app inspect active
