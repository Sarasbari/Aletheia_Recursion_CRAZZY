import base64
import hashlib
import io
import math
import os
from datetime import datetime, timezone

import imagehash
import requests
from PIL import Image

from celery_app import celery_app


def _sha256_hex(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def _phash_hex(data: bytes) -> str:
    img = Image.open(io.BytesIO(data)).convert("RGB")
    return str(imagehash.phash(img))


def _merkle_root_16x16(data: bytes) -> str:
    img = Image.open(io.BytesIO(data)).convert("RGBA")
    w, h = img.size
    block = 16
    leaves = []

    for y in range(0, h, block):
        for x in range(0, w, block):
            chunk = img.crop((x, y, min(x + block, w), min(y + block, h))).tobytes()
            leaves.append(hashlib.sha256(chunk).digest())

    if not leaves:
        return hashlib.sha256(b"").hexdigest()

    layer = leaves
    while len(layer) > 1:
        nxt = []
        for i in range(0, len(layer), 2):
            left = layer[i]
            right = layer[i + 1] if i + 1 < len(layer) else left
            nxt.append(hashlib.sha256(left + right).digest())
        layer = nxt

    return layer[0].hex()


def _upload_to_ipfs(data: bytes) -> str:
    endpoint = os.getenv("IPFS_ENDPOINT", "")
    if not endpoint:
        return "local-" + _sha256_hex(data)[:32]

    base = endpoint.strip()
    if not base.startswith("http://") and not base.startswith("https://"):
        base = f"http://{base}"
    url = base.rstrip("/") + "/api/v0/add?pin=true"

    files = {"file": ("image.bin", data)}
    resp = requests.post(url, files=files, timeout=30)
    resp.raise_for_status()
    payload = resp.json()
    return payload.get("Hash") or payload.get("Cid") or ("local-" + _sha256_hex(data)[:32])


def _anchor_polygon(record: dict) -> str:
    rpc = os.getenv("POLYGON_RPC_URL", "")
    priv = os.getenv("POLYGON_PRIVATE_KEY", "")
    contract = os.getenv("POLYGON_CONTRACT_ADDRESS", "")
    if not rpc or not priv or not contract:
        return "mock-" + _sha256_hex(str(record).encode("utf-8"))
    # Placeholder for real signed tx path. Keep deterministic fallback in this scaffold.
    return "mock-" + _sha256_hex(str(record).encode("utf-8"))


def _save_bigchaindb(record: dict) -> None:
    base = os.getenv("BIGCHAINDB_URL", "")
    if not base:
        return

    url = base.rstrip("/") + "/metadata"
    headers = {"Content-Type": "application/json"}
    api_key = os.getenv("BIGCHAINDB_API_KEY", "")
    if api_key:
        headers["Authorization"] = f"Bearer {api_key}"

    resp = requests.post(url, json=record, headers=headers, timeout=20)
    resp.raise_for_status()


@celery_app.task(name="aletheia.process_upload")
def process_upload(job: dict) -> dict:
    image_b64 = job.get("imageBase64", "")
    image_bytes = base64.b64decode(image_b64)

    sha = _sha256_hex(image_bytes)
    ph = _phash_hex(image_bytes)
    merkle_root = _merkle_root_16x16(image_bytes)
    cid = _upload_to_ipfs(image_bytes)

    proof = {
        "jobId": job.get("jobId"),
        "fileName": job.get("fileName"),
        "sha256": sha,
        "phash": ph,
        "merkleRoot": merkle_root,
        "ipfsCID": cid,
        "timestamp": datetime.now(timezone.utc).isoformat(),
        "status": "ANCHORED",
    }
    proof["txHash"] = _anchor_polygon(proof)
    _save_bigchaindb(proof)

    return proof
