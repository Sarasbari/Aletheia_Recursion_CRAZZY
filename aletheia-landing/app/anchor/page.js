"use client";

import { useState, useRef, useCallback, useEffect } from "react";
import { uploadImage, uploadVideo } from "../lib/api";
import { useAppKitAccount } from "@reown/appkit/react";
import { useSendTransaction } from "wagmi";
import { parseEther, stringToHex } from "viem";
import Navbar from "../components/Navbar";
import ProcessingModal from "../components/ProcessingModal";

const PIPELINE_STEPS = [
  { key: "upload", label: "Upload to server" },
  { key: "hash", label: "SHA-256 hash" },
  { key: "phash", label: "Perceptual hash" },
  { key: "sign", label: "ECDSA signature" },
  { key: "store", label: "IPFS storage" },
  { key: "anchor", label: "On-chain anchor" },
];

const ANCHOR_ADDRESS = "0x000000000000000000000000000000000000dEaD";

export default function AnchorPage() {
  const [file, setFile] = useState(null);
  const [preview, setPreview] = useState(null);
  const [mediaType, setMediaType] = useState("image");
  const [metadata, setMetadata] = useState({
    location: "",
    deviceInfo: "",
    uploaderId: "",
  });
  const [status, setStatus] = useState("idle");
  const [currentStep, setCurrentStep] = useState(-1);
  const [result, setResult] = useState(null);
  const [txHash, setTxHash] = useState(null);
  const [error, setError] = useState(null);
  const [inputMode, setInputMode] = useState("file");
  const [cameraActive, setCameraActive] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const inputRef = useRef(null);
  const videoRef = useRef(null);
  const streamRef = useRef(null);

  const { address, isConnected } = useAppKitAccount();
  const { sendTransactionAsync } = useSendTransaction();

  const handleFile = useCallback((f) => {
    if (!f) return;
    setFile(f);
    setResult(null);
    setError(null);
    setTxHash(null);
    setStatus("idle");
    setCurrentStep(-1);
    const reader = new FileReader();
    reader.onload = (e) => setPreview(e.target.result);
    reader.readAsDataURL(f);
  }, []);

  const handleDrop = useCallback(
    (e) => {
      e.preventDefault();
      e.currentTarget.classList.remove("drag-over");
      const f = e.dataTransfer?.files?.[0];
      if (!f) return;
      if (mediaType === "video" && f.type.startsWith("video/")) handleFile(f);
      else if (mediaType === "image" && f.type.startsWith("image/")) handleFile(f);
    },
    [handleFile, mediaType]
  );

  const handleDragOver = (e) => {
    e.preventDefault();
    e.currentTarget.classList.add("drag-over");
  };

  const handleDragLeave = (e) => {
    e.currentTarget.classList.remove("drag-over");
  };

  // Camera functions
  const startCamera = useCallback(async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: "environment", width: { ideal: 1280 }, height: { ideal: 720 } },
        audio: false,
      });
      streamRef.current = stream;
      setCameraActive(true);
    } catch {
      setError("Camera access denied. Please allow camera permissions.");
    }
  }, []);

  // Attach stream to video element after render
  useEffect(() => {
    if (cameraActive && videoRef.current && streamRef.current) {
      videoRef.current.srcObject = streamRef.current;
      videoRef.current.play().catch(() => {});
    }
  }, [cameraActive]);

  const capturePhoto = useCallback(() => {
    const video = videoRef.current;
    if (!video || !video.videoWidth) return;
    const canvas = document.createElement("canvas");
    canvas.width = video.videoWidth;
    canvas.height = video.videoHeight;
    canvas.getContext("2d").drawImage(video, 0, 0);
    canvas.toBlob(
      (blob) => {
        if (!blob) return;
        const f = new File([blob], `capture_${Date.now()}.jpg`, { type: "image/jpeg" });
        handleFile(f);
        stopCamera();
      },
      "image/jpeg",
      0.92
    );
  }, [handleFile]);

  const stopCamera = useCallback(() => {
    if (streamRef.current) {
      streamRef.current.getTracks().forEach((t) => t.stop());
      streamRef.current = null;
    }
    setCameraActive(false);
  }, []);

  const handleSubmit = async () => {
    if (!file || !isConnected) return;
    setShowModal(true);
    setStatus("uploading");
    setError(null);
    setResult(null);
    setTxHash(null);

    for (let i = 0; i < PIPELINE_STEPS.length - 1; i++) {
      setCurrentStep(i);
      await new Promise((r) => setTimeout(r, 600));
    }

    try {
      let data;
      if (mediaType === "video") {
        data = await uploadVideo({ videoFile: file });
      } else {
        data = await uploadImage({
          imageFile: file,
          devicePublicKey: "",
          deviceSignature: "",
          captureTimestamp: new Date().toISOString(),
          location: metadata.location || undefined,
          deviceInfo: metadata.deviceInfo || undefined,
          uploaderId: metadata.uploaderId || undefined,
        });
      }
      setResult(data);
      setCurrentStep(PIPELINE_STEPS.length - 1);
      setStatus("anchoring");

      const proofData = stringToHex(
        `aletheia:${data.proofId || ""}:${data.sha256 || ""}`
      );
      const hash = await sendTransactionAsync({
        to: ANCHOR_ADDRESS,
        value: parseEther("0.0001"),
        data: proofData,
        gasPrice: 30n * 10n ** 9n, // 30 Gwei
      });
      setTxHash(hash);
      setStatus("done");
    } catch (err) {
      setError(err.message);
      setStatus("error");
    }
  };

  const truncate = (h) => (h ? `${h.slice(0, 10)}…${h.slice(-6)}` : "—");

  return (
    <>
      <Navbar />
      <main className="app-page">
        <div className="app-container">
          {/* Left panel — inputs */}
          <section className="app-panel input-panel">
            <div className="panel-header">
              <div className="panel-badge">
                <span className="badge-dot" />
                anchor
              </div>
              <h1>Upload & Anchor</h1>
              <p>Upload or capture an image or video. We hash, sign, and anchor it on Polygon Amoy.</p>
            </div>

            {/* Wallet chip */}
            <div className={`wallet-chip ${isConnected ? "connected" : ""}`}>
              <span className="wc-dot" />
              {isConnected ? (
                <span>Connected — <span className="wc-addr">{address?.slice(0, 6)}…{address?.slice(-4)}</span></span>
              ) : (
                <span>Wallet not connected — connect to anchor</span>
              )}
            </div>

            {/* Media type toggle */}
            <div className="media-type-toggle">
              <button className={mediaType === "image" ? "active" : ""} onClick={() => { setMediaType("image"); setFile(null); setPreview(null); setResult(null); setError(null); }}>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
                Image
              </button>
              <button className={mediaType === "video" ? "active" : ""} onClick={() => { setMediaType("video"); setInputMode("file"); setFile(null); setPreview(null); setResult(null); setError(null); stopCamera(); }}>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polygon points="23 7 16 12 23 17 23 7"/><rect x="1" y="5" width="15" height="14" rx="2" ry="2"/></svg>
                Video
              </button>
            </div>

            {/* Input switcher */}
            <div className="input-tabs">
              <button className={inputMode === "file" ? "active" : ""} onClick={() => { setInputMode("file"); stopCamera(); }}>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
                Upload
              </button>
              {mediaType === "image" && (
                <button className={inputMode === "camera" ? "active" : ""} onClick={() => setInputMode("camera")}>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg>
                  Camera
                </button>
              )}
            </div>

            {/* Drop zone */}
            {inputMode === "file" && (
              <div
                className={`drop-area${file ? " has-file" : ""}`}
                onDrop={handleDrop}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onClick={() => inputRef.current?.click()}
              >
                {preview && mediaType === "image" ? (
                  <img src={preview} alt="Preview" className="drop-preview" />
                ) : preview && mediaType === "video" ? (
                  <video src={preview} className="drop-preview video-preview" controls muted />
                ) : (
                  <div className="drop-placeholder">
                    <div className="drop-icon">
                      {mediaType === "video" ? (
                        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><polygon points="23 7 16 12 23 17 23 7"/><rect x="1" y="5" width="15" height="14" rx="2" ry="2"/></svg>
                      ) : (
                        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
                      )}
                    </div>
                    <p className="drop-text">{mediaType === "video" ? "Drag & drop video here" : "Drag & drop image here"}</p>
                    <p className="drop-subtext">{mediaType === "video" ? "or click to browse • MP4, MOV, WEBM" : "or click to browse • PNG, JPG, WEBP"}</p>
                  </div>
                )}
                <input ref={inputRef} type="file" accept={mediaType === "video" ? "video/*" : "image/*"} hidden onChange={(e) => handleFile(e.target.files?.[0])} />
              </div>
            )}

            {/* Camera */}
            {inputMode === "camera" && (
              <div className="camera-area">
                {cameraActive ? (
                  <>
                    <video ref={videoRef} autoPlay playsInline muted className="camera-video" />
                    <div className="camera-actions">
                      <button className="capture-btn" onClick={capturePhoto} aria-label="Capture photo">
                        <span className="capture-ring" />
                      </button>
                      <button className="camera-cancel" onClick={stopCamera}>Cancel</button>
                    </div>
                  </>
                ) : preview ? (
                  <div className="drop-area has-file" onClick={() => { setPreview(null); setFile(null); }}>
                    <img src={preview} alt="Captured" className="drop-preview" />
                    <span className="retake-badge">tap to retake</span>
                  </div>
                ) : (
                  <div className="drop-area camera-prompt" onClick={startCamera}>
                    <div className="drop-placeholder">
                      <div className="drop-icon camera-icon-wrapper">
                        <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" strokeLinejoin="round"><path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg>
                      </div>
                      <p className="drop-text">Tap to open camera</p>
                      <p className="drop-subtext">Uses rear camera by default</p>
                    </div>
                  </div>
                )}
              </div>
            )}



            {file && (
              <div className="file-badge">
                <span className="file-name">{file.name}</span>
                <span className="file-size">{(file.size / 1024).toFixed(0)} KB</span>
              </div>
            )}

            {/* Metadata */}
            <div className="meta-grid">
              <label className="meta-input">
                <span>Location</span>
                <input type="text" placeholder="28.6139° N, 77.2090° E" value={metadata.location} onChange={(e) => setMetadata((m) => ({ ...m, location: e.target.value }))} />
              </label>
              <label className="meta-input">
                <span>Device</span>
                <input type="text" placeholder="iPhone 15 Pro" value={metadata.deviceInfo} onChange={(e) => setMetadata((m) => ({ ...m, deviceInfo: e.target.value }))} />
              </label>
              <label className="meta-input full">
                <span>Uploader ID</span>
                <input type="text" placeholder="d.sharma" value={metadata.uploaderId} onChange={(e) => setMetadata((m) => ({ ...m, uploaderId: e.target.value }))} />
              </label>
            </div>

            <button className="submit-btn" onClick={handleSubmit} disabled={!file || !isConnected || status === "uploading" || status === "anchoring"}>
              {!isConnected ? "Connect Wallet First" : !file ? (mediaType === "video" ? "Select a Video" : "Select an Image") : "Anchor on Polygon (~0.0001 POL)"}
            </button>
          </section>

          {/* Right panel — results */}
          <section className="app-panel result-panel">
            <div className="panel-header">
              <div className="panel-badge result-badge">
                <span className="badge-dot result-dot" />
                pipeline
              </div>
              <h2>Crypto Pipeline</h2>
            </div>

            <div className="pipeline-card">
              {PIPELINE_STEPS.map((step, i) => {
                const isDone = (status === "done" && currentStep >= i) || currentStep > i;
                const isActive = currentStep === i && status !== "done";
                return (
                  <div key={step.key} className={`pl-step ${isDone ? "done" : isActive ? "active" : "idle"}`}>
                    <div className="pl-indicator">
                      {isDone ? (
                        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
                      ) : isActive ? (
                        <div className="pl-spinner" />
                      ) : (
                        <span>{String(i + 1).padStart(2, "0")}</span>
                      )}
                    </div>
                    <span className="pl-label">{step.label}</span>
                    {i < PIPELINE_STEPS.length - 1 && <div className="pl-line" />}
                  </div>
                );
              })}
            </div>

            {result && (
              <div className="result-card">
                <div className={`result-badge-large ${result.status === "DUPLICATE" ? "warn" : "success"}`}>
                  {result.status === "DUPLICATE" ? "⚠ Duplicate" : "✓ Anchored"}
                </div>
                <div className="result-grid">
                  {[
                    { k: "Proof ID", v: result.proofId },
                    { k: "SHA-256", v: result.sha256 },
                    result.phash && { k: "pHash", v: result.phash },
                    result.merkleRoot && { k: "Merkle Root", v: result.merkleRoot },
                    result.storageRef && { k: "Storage", v: result.storageRef },
                  ]
                    .filter(Boolean)
                    .map((r) => (
                      <div className="rg-row" key={r.k}>
                        <span className="rg-key">{r.k}</span>
                        <span className="rg-val">{truncate(r.v)}</span>
                      </div>
                    ))}
                  {txHash && (
                    <div className="rg-row">
                      <span className="rg-key">Tx Hash</span>
                      <a className="rg-val rg-link" href={`https://amoy.polygonscan.com/tx/${txHash}`} target="_blank" rel="noopener noreferrer">
                        {truncate(txHash)} ↗
                      </a>
                    </div>
                  )}
                </div>
              </div>
            )}

            {error && <div className="error-card"><span>✕</span> {error}</div>}

            {!result && !error && (
              <div className="empty-state">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="0.8" strokeLinecap="round" strokeLinejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/><path d="M9 12l2 2 4-4"/></svg>
                <p>Upload an image or video to see the pipeline in action</p>
                <span className="empty-hint">Gas fee: ~0.0001 POL (Amoy testnet)</span>
              </div>
            )}
          </section>
        </div>
      </main>

      <ProcessingModal
        isOpen={showModal && (status === "uploading" || status === "anchoring" || status === "done" || status === "error")}
        steps={PIPELINE_STEPS}
        currentStep={currentStep}
        status={status}
        onClose={() => setShowModal(false)}
        mediaType={mediaType}
      />
    </>
  );
}
