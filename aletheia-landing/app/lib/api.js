const API_BASE = "/api/v1";

/**
 * Upload an image for anchoring.
 * POST /api/v1/images/upload (multipart/form-data)
 */
export async function uploadImage({
  imageFile,
  devicePublicKey,
  deviceSignature,
  captureTimestamp,
  parentHash,
  location,
  deviceInfo,
  uploaderId,
}) {
  const form = new FormData();
  form.append("image", imageFile);
  form.append("devicePublicKey", devicePublicKey || "");
  form.append("deviceSignature", deviceSignature || "");
  form.append("captureTimestamp", captureTimestamp || new Date().toISOString());
  if (parentHash) form.append("parentHash", parentHash);
  if (location) form.append("location", location);
  if (deviceInfo) form.append("deviceInfo", deviceInfo);
  if (uploaderId) form.append("uploaderId", uploaderId);

  const res = await fetch(`${API_BASE}/images/upload`, {
    method: "POST",
    body: form,
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || `Upload failed (${res.status})`);
  }

  return res.json();
}

/**
 * Verify an image by uploading bytes.
 * POST /api/v1/verify (multipart/form-data)
 */
export async function verifyImage({
  imageFile,
  captureTimestamp,
  location,
  deviceInfo,
  devicePublicKey,
  deviceSignature,
}) {
  const form = new FormData();
  form.append("image", imageFile);
  if (captureTimestamp) form.append("captureTimestamp", captureTimestamp);
  if (location) form.append("location", location);
  if (deviceInfo) form.append("deviceInfo", deviceInfo);
  if (devicePublicKey) form.append("devicePublicKey", devicePublicKey);
  if (deviceSignature) form.append("deviceSignature", deviceSignature);

  const res = await fetch(`${API_BASE}/verify`, {
    method: "POST",
    body: form,
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || `Verification failed (${res.status})`);
  }

  return res.json();
}

/**
 * Verify an image from URL.
 * POST /api/v1/verify/url (JSON body)
 */
export async function verifyImageUrl({
  imageUrl,
  captureTimestamp,
  location,
  deviceInfo,
  devicePublicKey,
  deviceSignature,
}) {
  const body = { imageUrl };
  if (captureTimestamp) body.captureTimestamp = captureTimestamp;
  if (location) body.location = location;
  if (deviceInfo) body.deviceInfo = deviceInfo;
  if (devicePublicKey) body.devicePublicKey = devicePublicKey;
  if (deviceSignature) body.deviceSignature = deviceSignature;

  const res = await fetch(`${API_BASE}/verify/url`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || `URL verification failed (${res.status})`);
  }

  return res.json();
}

/**
 * Upload a video for anchoring.
 * POST /api/v1/video/upload (multipart/form-data)
 */
export async function uploadVideo({ videoFile }) {
  const form = new FormData();
  form.append("video", videoFile);

  const res = await fetch(`${API_BASE}/video/upload`, {
    method: "POST",
    body: form,
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || `Video upload failed (${res.status})`);
  }

  return res.json();
}

/**
 * Verify a video.
 * POST /api/v1/video/verify (multipart/form-data)
 */
export async function verifyVideo({ videoFile }) {
  const form = new FormData();
  form.append("video", videoFile);

  const res = await fetch(`${API_BASE}/video/verify`, {
    method: "POST",
    body: form,
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || `Video verification failed (${res.status})`);
  }

  return res.json();
}
