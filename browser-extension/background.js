const API_BASE_URL = "https://api.example.com";
const CONTEXT_MENU_ID = "aletheia-verify-image";

// In-memory cache keyed by image URL to avoid repeated analyze/verify requests.
const imageResultCache = new Map();
const imageHashCache = new Map();

function createImageContextMenu() {
  chrome.contextMenus.removeAll(() => {
    chrome.contextMenus.create({
      id: CONTEXT_MENU_ID,
      title: "Verify Image",
      contexts: ["image"]
    });
  });
}

function toHex(buffer) {
  const bytes = new Uint8Array(buffer);
  return Array.from(bytes)
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
}

async function fetchJson(path, payload) {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(payload)
  });

  if (!response.ok) {
    throw new Error(`API request failed: ${path} (${response.status})`);
  }

  return response.json();
}

async function hashImageFromUrl(imageUrl) {
  // Fetch + blob + SHA-256 happens in service worker to avoid page CORS limits.
  const imageResponse = await fetch(imageUrl);
  if (!imageResponse.ok) {
    throw new Error(`Image fetch failed (${imageResponse.status})`);
  }

  const blob = await imageResponse.blob();
  const buffer = await blob.arrayBuffer();
  const digest = await crypto.subtle.digest("SHA-256", buffer);
  return toHex(digest);
}

async function getImageHash(imageUrl) {
  if (imageHashCache.has(imageUrl)) {
    return imageHashCache.get(imageUrl);
  }

  const imageHash = await hashImageFromUrl(imageUrl);
  imageHashCache.set(imageUrl, imageHash);
  return imageHash;
}

async function analyzeImage(imageUrl) {
  if (!imageUrl) {
    throw new Error("Missing image URL");
  }

  if (imageResultCache.has(imageUrl)) {
    return {
      ...imageResultCache.get(imageUrl),
      cached: true
    };
  }

  const imageHash = await getImageHash(imageUrl);

  const [analyzeData, verifyData] = await Promise.all([
    fetchJson("/analyze", { image_url: imageUrl }),
    fetchJson("/verify", { image_hash: imageHash })
  ]);

  const result = {
    imageUrl,
    imageHash,
    aiProbability: Number(analyzeData?.ai_probability ?? 0),
    verified: Boolean(verifyData?.verified),
    timestamp: verifyData?.timestamp || null,
    analyzedAt: new Date().toISOString(),
    cached: false
  };

  imageResultCache.set(imageUrl, result);
  await chrome.storage.local.set({ aletheiaLastResult: result });
  return result;
}

async function anchorImage(imageUrl, imageHash) {
  if (!imageUrl || !imageHash) {
    throw new Error("Missing image URL or hash");
  }

  return fetchJson("/anchor", {
    image_url: imageUrl,
    image_hash: imageHash
  });
}

async function anchorImageByUrl(imageUrl) {
  if (!imageUrl) {
    throw new Error("Missing image URL");
  }

  const imageHash = await getImageHash(imageUrl);
  const data = await anchorImage(imageUrl, imageHash);

  return {
    imageUrl,
    imageHash,
    success: Boolean(data?.success),
    timestamp: new Date().toISOString()
  };
}

chrome.runtime.onInstalled.addListener(() => {
  createImageContextMenu();
});

chrome.runtime.onStartup.addListener(() => {
  createImageContextMenu();
});

chrome.contextMenus.onClicked.addListener((info, tab) => {
  if (info.menuItemId !== CONTEXT_MENU_ID || !info.srcUrl || !tab?.id) {
    return;
  }

  // Delegate UI rendering to content script in the clicked tab.
  chrome.tabs.sendMessage(tab.id, {
    type: "ALETHEIA_VERIFY_FROM_CONTEXT_MENU",
    imageUrl: info.srcUrl
  });
});

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  const type = message?.type;

  if (type === "ALETHEIA_ANALYZE_IMAGE") {
    analyzeImage(message.imageUrl)
      .then((result) => sendResponse({ ok: true, result }))
      .catch((error) => sendResponse({ ok: false, error: error.message }));
    return true;
  }

  if (type === "ALETHEIA_ANCHOR_IMAGE") {
    anchorImage(message.imageUrl, message.imageHash)
      .then((data) => sendResponse({ ok: true, data }))
      .catch((error) => sendResponse({ ok: false, error: error.message }));
    return true;
  }

  if (type === "ALETHEIA_ANCHOR_IMAGE_URL") {
    anchorImageByUrl(message.imageUrl)
      .then((result) => sendResponse({ ok: true, result }))
      .catch((error) => sendResponse({ ok: false, error: error.message }));
    return true;
  }

  if (type === "ALETHEIA_GET_LAST_RESULT") {
    chrome.storage.local
      .get(["aletheiaLastResult"])
      .then((data) => sendResponse({ ok: true, result: data.aletheiaLastResult || null }))
      .catch((error) => sendResponse({ ok: false, error: error.message }));
    return true;
  }

  return false;
});
