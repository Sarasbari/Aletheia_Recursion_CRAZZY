const emptyStateEl = document.getElementById("emptyState");
const resultPanelEl = document.getElementById("resultPanel");

const previewEl = document.getElementById("preview");
const aiProbabilityEl = document.getElementById("aiProbability");
const aiStatusEl = document.getElementById("aiStatus");
const verifyStatusEl = document.getElementById("verifyStatus");
const timestampEl = document.getElementById("timestamp");
const anchorBtn = document.getElementById("anchorBtn");
const statusMessageEl = document.getElementById("statusMessage");

let currentResult = null;

function classifyAi(probability) {
  if (probability < 30) {
    return { text: "Likely Real", className: "real" };
  }
  if (probability <= 70) {
    return { text: "Uncertain", className: "uncertain" };
  }
  return { text: "Likely AI-generated", className: "ai" };
}

function renderEmpty() {
  emptyStateEl.classList.remove("hidden");
  resultPanelEl.classList.add("hidden");
}

function renderResult(result) {
  currentResult = result;

  const ai = classifyAi(result.aiProbability);

  previewEl.src = result.imageUrl;
  aiProbabilityEl.textContent = `${result.aiProbability.toFixed(1)}%`;
  aiStatusEl.textContent = ai.text;
  aiStatusEl.className = `badge ${ai.className}`;
  verifyStatusEl.textContent = result.verified ? "✅ Verified" : "❌ Not Found";
  timestampEl.textContent = result.timestamp ? new Date(result.timestamp).toLocaleString() : "N/A";
  statusMessageEl.textContent = result.cached ? "Loaded from cache." : "";

  emptyStateEl.classList.add("hidden");
  resultPanelEl.classList.remove("hidden");
}

async function loadLastResult() {
  try {
    const response = await chrome.runtime.sendMessage({ type: "ALETHEIA_GET_LAST_RESULT" });
    if (!response?.ok || !response.result) {
      renderEmpty();
      return;
    }
    renderResult(response.result);
  } catch {
    renderEmpty();
  }
}

anchorBtn.addEventListener("click", async () => {
  if (!currentResult) {
    return;
  }

  anchorBtn.disabled = true;
  statusMessageEl.textContent = "Anchoring image...";

  try {
    const response = await chrome.runtime.sendMessage({
      type: "ALETHEIA_ANCHOR_IMAGE",
      imageUrl: currentResult.imageUrl,
      imageHash: currentResult.imageHash
    });

    if (!response?.ok || !response?.data?.success) {
      throw new Error(response?.error || "Anchor failed");
    }

    statusMessageEl.textContent = "Image anchored successfully.";
  } catch (error) {
    statusMessageEl.textContent = `Anchor failed: ${error.message}`;
  } finally {
    anchorBtn.disabled = false;
  }
});

loadLastResult();
