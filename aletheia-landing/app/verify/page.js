"use client";

import { useState, useRef, useCallback, useEffect } from "react";
import { verifyImage, verifyImageUrl } from "../lib/api";
import { useAppKitAccount } from "@reown/appkit/react";
import { useSendTransaction } from "wagmi";
import { parseEther, stringToHex } from "viem";
import Navbar from "../components/Navbar";

const VERIFY_ADDRESS = "0x000000000000000000000000000000000000dEaD";

const VERDICT_MAP = {
  AUTHENTIC: { label: "Authentic", icon: "✓", className: "authentic" },
  TAMPERED: { label: "Tampered", icon: "✕", className: "tampered" },
  SIMILAR: { label: "Similar Match", icon: "≈", className: "similar" },
  UNKNOWN: { label: "Unknown", icon: "?", className: "unknown" },
};

export default function VerifyPage() {
  const [mode, setMode] = useState("file");
  const [file, setFile] = useState(null);
  const [preview, setPreview] = useState(null);
  const [status, setStatus] = useState("idle");
  const [result, setResult] = useState(null);
  const [error, setError] = useState(null);
  const [txHash, setTxHash] = useState(null);
  const [cameraActive, setCameraActive] = useState(false);
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
    setStatus("idle");
    const reader = new FileReader();
    reader.onload = (e) => setPreview(e.target.result);
    reader.readAsDataURL(f);
  }, []);

  const handleDrop = useCallback(
    (e) => {
      e.preventDefault();
      e.currentTarget.classList.remove("drag-over");
      const f = e.dataTransfer?.files?.[0];
      if (f && f.type.startsWith("image/")) handleFile(f);
    },
    [handleFile]
  );

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

  const handleVerify = async () => {
    if (!isConnected) {
      setError("Please connect your wallet first.");
      return;
    }
    setStatus("verifying");
    setError(null);
    setResult(null);
    setTxHash(null);
    try {
      // Step 1: Send gas fee tx
      const proofData = stringToHex(
        `aletheia-verify:${Date.now()}`
      );
      const hash = await sendTransactionAsync({
        to: VERIFY_ADDRESS,
        value: parseEther("0.0001"),
        data: proofData,
        gasPrice: 30n * 10n ** 9n, // 30 Gwei
      });
      setTxHash(hash);

      // Step 2: Call backend verification
      if (!file) throw new Error("Please provide an image.");
      const data = await verifyImage({ imageFile: file });
      setResult(data);
      setStatus("done");
    } catch (err) {
      setError(err.message);
      setStatus("error");
    }
  };

  const verdict = result ? VERDICT_MAP[result.verdict] || VERDICT_MAP.UNKNOWN : null;

  return (
    <>
      <Navbar />
      <main className="app-page">
        <div className="app-container">
          {/* Left: Input */}
          <section className="app-panel input-panel">
            <div className="panel-header">
              <div className="panel-badge verify-badge">
                <span className="badge-dot verify-dot" />
                verify
              </div>
              <h1>Verify Image</h1>
              <p>Check if an image has been anchored on-chain. Upload or capture a photo.</p>
            </div>

            {/* Wallet chip */}
            <div className={`wallet-chip ${isConnected ? "connected" : ""}`}>
              <span className="wc-dot" />
              {isConnected ? (
                <span>Connected — <span className="wc-addr">{address?.slice(0, 6)}…{address?.slice(-4)}</span></span>
              ) : (
                <span>Wallet not connected — connect to verify</span>
              )}
            </div>

            <div className="input-tabs">
              <button className={mode === "file" ? "active" : ""} onClick={() => { setMode("file"); setResult(null); setError(null); stopCamera(); }}>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
                Upload
              </button>
              <button className={mode === "camera" ? "active" : ""} onClick={() => { setMode("camera"); setResult(null); setError(null); }}>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg>
                Camera
              </button>
            </div>

            {mode === "file" && (
              <div
                className={`drop-area${file ? " has-file" : ""}`}
                onDrop={handleDrop}
                onDragOver={(e) => { e.preventDefault(); e.currentTarget.classList.add("drag-over"); }}
                onDragLeave={(e) => e.currentTarget.classList.remove("drag-over")}
                onClick={() => inputRef.current?.click()}
              >
                {preview ? (
                  <img src={preview} alt="Preview" className="drop-preview" />
                ) : (
                  <div className="drop-placeholder">
                    <div className="drop-icon">
                      <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
                    </div>
                    <p className="drop-text">Drag & drop image</p>
                    <p className="drop-subtext">or click to browse</p>
                  </div>
                )}
                <input ref={inputRef} type="file" accept="image/*" hidden onChange={(e) => handleFile(e.target.files?.[0])} />
              </div>
            )}

            {mode === "camera" && (
              <div className="camera-area">
                {cameraActive ? (
                  <>
                    <video ref={videoRef} autoPlay playsInline muted className="camera-video" />
                    <div className="camera-actions">
                      <button className="capture-btn" onClick={capturePhoto} aria-label="Capture"><span className="capture-ring" /></button>
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
                      <div className="drop-icon"><svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" strokeLinejoin="round"><path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg></div>
                      <p className="drop-text">Tap to open camera</p>
                      <p className="drop-subtext">Uses rear camera</p>
                    </div>
                  </div>
                )}
              </div>
            )}



            <button
              className="submit-btn"
              onClick={handleVerify}
              disabled={status === "verifying" || !isConnected || !file}
            >
              {status === "verifying" ? (
                <><span className="btn-spinner" /> Verifying…</>
              ) : !isConnected ? (
                "Connect Wallet First"
              ) : (
                "Verify Image (~0.0001 POL)"
              )}
            </button>
          </section>

          {/* Right: Result */}
          <section className="app-panel result-panel">
            <div className="panel-header">
              <div className="panel-badge result-badge">
                <span className="badge-dot result-dot" />
                result
              </div>
              <h2>Verification Result</h2>
            </div>

            {result && (
              <>
                <div className={`verdict-card ${verdict.className}`}>
                  <span className="verdict-icon">{verdict.icon}</span>
                  <span className="verdict-text">{verdict.label}</span>
                </div>

                <div className="result-card">
                  <div className="result-grid">
                    {result.breakdown && (
                      <>
                        <div className="rg-row">
                          <span className="rg-key">Hash</span>
                          <span className={`rg-val ${result.breakdown.hashMatch ? "rg-pass" : "rg-fail"}`}>
                            {result.breakdown.hashMatch ? "MATCH" : "MISMATCH"}
                          </span>
                        </div>
                        <div className="rg-row">
                          <span className="rg-key">Signature</span>
                          <span className={`rg-val ${result.breakdown.signatureValid ? "rg-pass" : "rg-fail"}`}>
                            {result.breakdown.signatureValid ? "VALID" : "INVALID"}
                          </span>
                        </div>
                        {result.breakdown.similarity !== undefined && (
                          <div className="rg-row">
                            <span className="rg-key">Similarity</span>
                            <span className="rg-val">{(result.breakdown.similarity * 100).toFixed(1)}%</span>
                          </div>
                        )}
                      </>
                    )}
                    {result.matchedProofId && (
                      <div className="rg-row">
                        <span className="rg-key">Matched Proof</span>
                        <span className="rg-val rg-mono">{result.matchedProofId}</span>
                      </div>
                    )}
                    {result.proofReport && (
                      <div className="rg-row">
                        <span className="rg-key">SHA-256</span>
                        <span className="rg-val rg-mono">{result.proofReport.sha256?.slice(0, 16)}…</span>
                      </div>
                    )}
                    {txHash && (
                      <div className="rg-row">
                        <span className="rg-key">Tx Hash</span>
                        <a className="rg-val rg-link" href={`https://amoy.polygonscan.com/tx/${txHash}`} target="_blank" rel="noopener noreferrer">
                          {txHash.slice(0, 10)}…{txHash.slice(-6)} ↗
                        </a>
                      </div>
                    )}
                  </div>
                </div>

                {result.tamperedRegions?.length > 0 && (
                  <div className="tampered-alert">
                    <span>⚠ Tampered regions detected at:</span>
                    {result.tamperedRegions.map((r, i) => (
                      <span key={i} className="rg-mono"> ({r.x},{r.y})</span>
                    ))}
                  </div>
                )}
              </>
            )}

            {error && <div className="error-card"><span>✕</span> {error}</div>}

            {!result && !error && (
              <div className="empty-state">
                <svg width="56" height="56" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="0.7" strokeLinecap="round" strokeLinejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/><path d="M9 12l2 2 4-4"/></svg>
                <p>Upload an image or paste a URL to verify its on-chain status</p>
                <span className="empty-hint">We compare against all anchored proofs</span>
              </div>
            )}
          </section>
        </div>
      </main>
    </>
  );
}
