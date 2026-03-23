from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from celery.result import AsyncResult

from celery_app import celery_app
from tasks import process_upload

app = FastAPI(title="Aletheia Task API")


class UploadTaskRequest(BaseModel):
    jobId: str
    fileName: str
    imageBase64: str


@app.get("/health")
def health():
    return {"service": "task-api", "status": "ok"}


@app.post("/tasks/upload")
def enqueue_upload(req: UploadTaskRequest):
    task = process_upload.delay(req.model_dump())
    return {"taskId": task.id, "jobId": req.jobId, "status": "QUEUED"}


@app.get("/tasks/{task_id}")
def task_status(task_id: str):
    res = AsyncResult(task_id, app=celery_app)
    payload = {
        "taskId": task_id,
        "state": res.state,
    }
    if res.successful():
        payload["result"] = res.result
    elif res.failed():
        payload["error"] = str(res.result)
    return payload
