(() => {
  const trackedImages = new WeakSet();
  // WeakSet prevents duplicate registrations while allowing garbage collection.

  let activeImage = null;
  let hideOverlayTimer = null;

  const overlayActions = document.createElement("div");
  overlayActions.className = "aletheia-overlay-actions";

  const overlayButton = document.createElement("button");
  overlayButton.className = "truthlens-overlay-btn truthlens-verify-btn";
  overlayButton.type = "button";
  overlayButton.textContent = "🔍";
  overlayButton.title = "Verify image authenticity";

  const anchorOverlayButton = document.createElement("button");
  anchorOverlayButton.className = "truthlens-overlay-btn truthlens-anchor-btn";
  anchorOverlayButton.type = "button";
  anchorOverlayButton.textContent = "+";
  anchorOverlayButton.title = "Add image to blockchain";

  const resultCard = document.createElement("aside");
  resultCard.className = "truthlens-card truthlens-hidden";

  const style = document.createElement("style");
  style.textContent = `
    .aletheia-overlay-actions {
      position: fixed;
      display: none;
      gap: 6px;
      z-index: 2147483646;
    }

    .truthlens-overlay-btn {
      all: initial;
      position: static;
      width: 28px;
      height: 28px;
      border: 0;
      border-radius: 999px;
      background: #111827;
      color: #ffffff;
      cursor: pointer;
      box-shadow: 0 8px 16px rgba(0, 0, 0, 0.25);
      font-size: 14px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-family: "Segoe UI", Tahoma, sans-serif;
      line-height: 1;
      text-align: center;
    }

    .truthlens-overlay-btn:hover {
      transform: scale(1.05);
    }

    .truthlens-overlay-btn.truthlens-anchor-btn {
      background: #0f766e;
      font-weight: 700;
      font-size: 18px;
    }

    .truthlens-card {
      position: fixed;
      width: min(320px, calc(100vw - 24px));
      background: #ffffff;
      color: #111827;
      border: 1px solid #e5e7eb;
      border-radius: 12px;
      box-shadow: 0 20px 45px rgba(17, 24, 39, 0.2);
      z-index: 2147483647;
      font-family: "Segoe UI", Tahoma, sans-serif;
      overflow: hidden;
      animation: truthlens-fade-in 0.16s ease-out;
    }

    .truthlens-hidden {
      display: none;
    }

    .truthlens-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 10px 12px;
      background: #f8fafc;
      border-bottom: 1px solid #e5e7eb;
      font-size: 13px;
      font-weight: 600;
    }

    .truthlens-close {
      border: 0;
      background: transparent;
      cursor: pointer;
      font-size: 16px;
      color: #6b7280;
    }

    .truthlens-body {
      padding: 12px;
      font-size: 13px;
      line-height: 1.35;
    }

    .truthlens-preview {
      width: 100%;
      max-height: 120px;
      object-fit: cover;
      border-radius: 8px;
      margin-bottom: 10px;
      background: #f3f4f6;
    }

    .truthlens-row {
      margin: 8px 0;
      color: #374151;
    }

    .truthlens-pill {
      display: inline-flex;
      align-items: center;
      border-radius: 999px;
      padding: 2px 8px;
      font-size: 12px;
      font-weight: 600;
      margin-left: 6px;
    }

    .truthlens-pill.real {
      background: #dcfce7;
      color: #166534;
    }

    .truthlens-pill.uncertain {
      background: #fef3c7;
      color: #92400e;
    }

    .truthlens-pill.ai {
      background: #fee2e2;
      color: #991b1b;
    }

    .truthlens-message {
      margin-top: 8px;
      font-size: 12px;
      color: #334155;
      min-height: 16px;
    }

    .truthlens-spinner {
      width: 18px;
      height: 18px;
      border: 2px solid #d1d5db;
      border-top-color: #0f766e;
      border-radius: 50%;
      animation: truthlens-spin 0.8s linear infinite;
      margin-right: 8px;
    }

    .truthlens-loading {
      display: flex;
      align-items: center;
      font-size: 13px;
      color: #334155;
      padding: 10px 0;
    }

    @keyframes truthlens-spin {
      to {
        transform: rotate(360deg);
      }
    }

    @keyframes truthlens-fade-in {
      from {
        opacity: 0;
        transform: translateY(4px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }
  `;

  document.documentElement.appendChild(style);
  overlayActions.appendChild(overlayButton);
  overlayActions.appendChild(anchorOverlayButton);

  document.documentElement.appendChild(overlayActions);
  document.documentElement.appendChild(resultCard);

  function classifyAiProbability(value) {
    if (value < 30) {
      return { label: "Likely Real", tone: "real" };
    }
    if (value <= 70) {
      return { label: "Uncertain", tone: "uncertain" };
    }
    return { label: "Likely AI-generated", tone: "ai" };
  }

  function getImageUrl(img) {
    return img?.currentSrc || img?.src || "";
  }

  function setOverlayPosition(img) {
    const rect = img.getBoundingClientRect();
    const x = Math.min(window.innerWidth - 68, Math.max(6, rect.right - 62));
    const y = Math.min(window.innerHeight - 34, Math.max(6, rect.top + 6));

    overlayActions.style.left = `${x}px`;
    overlayActions.style.top = `${y}px`;
    overlayActions.style.display = "flex";
  }

  function hideOverlayButtons() {
    overlayActions.style.display = "none";
  }

  function setCardPosition(referenceRect) {
    const baseX = referenceRect ? referenceRect.right + 10 : window.innerWidth - 336;
    const baseY = referenceRect ? referenceRect.top : 12;

    const maxX = Math.max(12, window.innerWidth - 336);
    const maxY = Math.max(12, window.innerHeight - 260);

    resultCard.style.left = `${Math.max(12, Math.min(baseX, maxX))}px`;
    resultCard.style.top = `${Math.max(12, Math.min(baseY, maxY))}px`;
  }

  function hideOverlaySoon() {
    clearTimeout(hideOverlayTimer);
    hideOverlayTimer = setTimeout(() => {
      const isOverlayFocused =
        document.activeElement === overlayButton || document.activeElement === anchorOverlayButton;
      if (!isOverlayFocused) {
        hideOverlayButtons();
      }
    }, 160);
  }

  function hideCard() {
    resultCard.classList.add("truthlens-hidden");
  }

  function showLoadingCard(imageUrl, referenceRect) {
    setCardPosition(referenceRect);
    resultCard.innerHTML = `
      <div class="truthlens-header">
        <span>Aletheia Analysis</span>
        <button class="truthlens-close" type="button" aria-label="Close">×</button>
      </div>
      <div class="truthlens-body">
        <img class="truthlens-preview" src="${imageUrl}" alt="Selected image preview" />
        <div class="truthlens-loading">
          <span class="truthlens-spinner"></span>
          <span>Analyzing image authenticity...</span>
        </div>
      </div>
    `;
    resultCard.classList.remove("truthlens-hidden");
    resultCard.querySelector(".truthlens-close")?.addEventListener("click", hideCard);
  }

  function showAnchoringCard(imageUrl, referenceRect) {
    setCardPosition(referenceRect);
    resultCard.innerHTML = `
      <div class="truthlens-header">
        <span>Aletheia Anchor</span>
        <button class="truthlens-close" type="button" aria-label="Close">×</button>
      </div>
      <div class="truthlens-body">
        <img class="truthlens-preview" src="${imageUrl}" alt="Selected image preview" />
        <div class="truthlens-loading">
          <span class="truthlens-spinner"></span>
          <span>Anchoring image to blockchain...</span>
        </div>
      </div>
    `;
    resultCard.classList.remove("truthlens-hidden");
    resultCard.querySelector(".truthlens-close")?.addEventListener("click", hideCard);
  }

  function renderResultCard(result, referenceRect) {
    const ai = classifyAiProbability(result.aiProbability);
    const verifiedText = result.verified ? "✅ Verified" : "❌ Not Found";
    const verifyTimestamp = result.timestamp
      ? new Date(result.timestamp).toLocaleString()
      : "N/A";
    const cacheText = result.cached ? " (cached)" : "";

    setCardPosition(referenceRect);
    resultCard.innerHTML = `
      <div class="truthlens-header">
        <span>Aletheia Result${cacheText}</span>
        <button class="truthlens-close" type="button" aria-label="Close">×</button>
      </div>
      <div class="truthlens-body">
        <img class="truthlens-preview" src="${result.imageUrl}" alt="Analyzed image preview" />

        <div class="truthlens-row">
          <strong>AI Probability:</strong> ${result.aiProbability.toFixed(1)}%
          <span class="truthlens-pill ${ai.tone}">${ai.label}</span>
        </div>

        <div class="truthlens-row">
          <strong>Blockchain Status:</strong> ${verifiedText}
        </div>

        <div class="truthlens-row">
          <strong>Timestamp:</strong> ${verifyTimestamp}
        </div>
      </div>
    `;
    resultCard.classList.remove("truthlens-hidden");

    const closeBtn = resultCard.querySelector(".truthlens-close");

    closeBtn?.addEventListener("click", hideCard);
  }

  function renderAnchorResultCard(result, referenceRect) {
    const timestampText = result.timestamp ? new Date(result.timestamp).toLocaleString() : "N/A";

    setCardPosition(referenceRect);
    resultCard.innerHTML = `
      <div class="truthlens-header">
        <span>Aletheia Anchor</span>
        <button class="truthlens-close" type="button" aria-label="Close">×</button>
      </div>
      <div class="truthlens-body">
        <img class="truthlens-preview" src="${result.imageUrl}" alt="Anchored image preview" />
        <div class="truthlens-row"><strong>Status:</strong> Image anchored successfully.</div>
        <div class="truthlens-row"><strong>Hash:</strong> ${result.imageHash.slice(0, 16)}...</div>
        <div class="truthlens-row"><strong>Timestamp:</strong> ${timestampText}</div>
      </div>
    `;
    resultCard.classList.remove("truthlens-hidden");
    resultCard.querySelector(".truthlens-close")?.addEventListener("click", hideCard);
  }

  function renderErrorCard(imageUrl, errorMessage, referenceRect) {
    setCardPosition(referenceRect);
    resultCard.innerHTML = `
      <div class="truthlens-header">
        <span>Aletheia Result</span>
        <button class="truthlens-close" type="button" aria-label="Close">×</button>
      </div>
      <div class="truthlens-body">
        ${imageUrl ? `<img class="truthlens-preview" src="${imageUrl}" alt="Image preview" />` : ""}
        <div class="truthlens-row"><strong>Error:</strong> ${errorMessage}</div>
      </div>
    `;
    resultCard.classList.remove("truthlens-hidden");
    resultCard.querySelector(".truthlens-close")?.addEventListener("click", hideCard);
  }

  async function analyzeImageByUrl(imageUrl, referenceRect = null) {
    if (!imageUrl) {
      return;
    }

    showLoadingCard(imageUrl, referenceRect);

    try {
      const response = await chrome.runtime.sendMessage({
        type: "ALETHEIA_ANALYZE_IMAGE",
        imageUrl
      });

      if (!response?.ok || !response?.result) {
        throw new Error(response?.error || "Analysis failed");
      }

      renderResultCard(response.result, referenceRect);
    } catch (error) {
      renderErrorCard(imageUrl, error.message, referenceRect);
    }
  }

  async function anchorImageByUrl(imageUrl, referenceRect = null) {
    if (!imageUrl) {
      return;
    }

    showAnchoringCard(imageUrl, referenceRect);

    try {
      const response = await chrome.runtime.sendMessage({
        type: "ALETHEIA_ANCHOR_IMAGE_URL",
        imageUrl
      });

      if (!response?.ok || !response?.result?.success) {
        throw new Error(response?.error || "Anchor failed");
      }

      renderAnchorResultCard(response.result, referenceRect);
    } catch (error) {
      renderErrorCard(imageUrl, error.message, referenceRect);
    }
  }

  function isCandidateImage(img) {
    if (!(img instanceof HTMLImageElement)) {
      return false;
    }

    const src = getImageUrl(img);
    if (!src) {
      return false;
    }

    const rect = img.getBoundingClientRect();
    return rect.width >= 48 && rect.height >= 48;
  }

  function registerImage(img) {
    if (!isCandidateImage(img) || trackedImages.has(img)) {
      return;
    }

    trackedImages.add(img);
    img.dataset.truthlensReady = "1";
  }

  function scanImages(root = document) {
    const images = root.querySelectorAll ? root.querySelectorAll("img") : [];
    images.forEach(registerImage);
  }

  scanImages(document);

  // Track images added after initial load (infinite scroll, client-side rendering, etc.).
  const observer = new MutationObserver((mutations) => {
    for (const mutation of mutations) {
      mutation.addedNodes.forEach((node) => {
        if (!(node instanceof Element)) {
          return;
        }

        if (node.tagName === "IMG") {
          registerImage(node);
        }

        scanImages(node);
      });
    }
  });

  observer.observe(document.documentElement, {
    childList: true,
    subtree: true
  });

  document.addEventListener(
    "mouseover",
    (event) => {
      const img = event.target?.closest?.("img");
      if (!img || !trackedImages.has(img)) {
        return;
      }

      activeImage = img;
      setOverlayPosition(img);
    },
    true
  );

  document.addEventListener(
    "scroll",
    () => {
      if (activeImage && overlayActions.style.display === "flex") {
        setOverlayPosition(activeImage);
      }
    },
    true
  );

  window.addEventListener("resize", () => {
    if (activeImage && overlayActions.style.display === "flex") {
      setOverlayPosition(activeImage);
    }
  });

  document.addEventListener(
    "mouseout",
    (event) => {
      const target = event.target;
      const related = event.relatedTarget;
      const isMovingToOverlay =
        related === overlayActions || overlayActions.contains(related);

      if (target instanceof HTMLImageElement && target === activeImage) {
        if (isMovingToOverlay) {
          return;
        }
        hideOverlaySoon();
      }
    },
    true
  );

  overlayActions.addEventListener("mouseenter", () => {
    clearTimeout(hideOverlayTimer);
  });

  overlayActions.addEventListener("mouseleave", hideOverlaySoon);

  overlayButton.addEventListener("click", () => {
    if (!activeImage) {
      return;
    }

    const imageUrl = getImageUrl(activeImage);
    const rect = activeImage.getBoundingClientRect();
    analyzeImageByUrl(imageUrl, rect);
  });

  anchorOverlayButton.addEventListener("click", () => {
    if (!activeImage) {
      return;
    }

    const imageUrl = getImageUrl(activeImage);
    const rect = activeImage.getBoundingClientRect();
    anchorImageByUrl(imageUrl, rect);
  });

  chrome.runtime.onMessage.addListener((message) => {
    if (message?.type !== "ALETHEIA_VERIFY_FROM_CONTEXT_MENU") {
      return;
    }

    const imageUrl = message.imageUrl;
    const referenceRect = activeImage ? activeImage.getBoundingClientRect() : null;
    analyzeImageByUrl(imageUrl, referenceRect);
  });
})();
