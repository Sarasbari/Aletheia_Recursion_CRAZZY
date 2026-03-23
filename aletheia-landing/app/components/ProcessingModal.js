"use client";

import { useEffect, useState } from "react";

const FACTS = [
  "SHA-512 produces a 128-character hexadecimal hash",
  "pHash compares visual similarity even after compression",
  "ECDSA signatures use elliptic curve cryptography",
  "IPFS content-addresses ensure data integrity",
  "Polygon processes ~65,000 transactions per second",
  "Each Merkle leaf is a unique proof of existence",
  "Your proof is permanent — it outlives any server",
  "Blockchain timestamps are cryptographically immutable",
];

export default function ProcessingModal({ isOpen, steps, currentStep, status, onClose }) {
  const [fact, setFact] = useState(0);
  const [dots, setDots] = useState("");

  useEffect(() => {
    if (!isOpen) return;
    const factInterval = setInterval(() => {
      setFact((f) => (f + 1) % FACTS.length);
    }, 3500);
    const dotInterval = setInterval(() => {
      setDots((d) => (d.length >= 3 ? "" : d + "."));
    }, 400);
    return () => {
      clearInterval(factInterval);
      clearInterval(dotInterval);
    };
  }, [isOpen]);

  if (!isOpen) return null;

  const progress = steps.length > 0
    ? Math.min(((currentStep + 1) / steps.length) * 100, 100)
    : 0;

  return (
    <div className="modal-overlay" onClick={status === "done" || status === "error" ? onClose : undefined}>
      <div className="modal-card" onClick={(e) => e.stopPropagation()}>
        {/* Top glow */}
        <div className="modal-glow" />

        {/* Header */}
        <div className="modal-header">
          <div className="modal-icon">
            {status === "done" ? (
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="var(--accent)" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
            ) : status === "error" ? (
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#ef4444" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
            ) : (
              <div className="spinner" />
            )}
          </div>
          <h3 className="modal-title">
            {status === "done"
              ? "Anchored Successfully"
              : status === "error"
              ? "Something Went Wrong"
              : "Processing" + dots}
          </h3>
          <p className="modal-subtitle">
            {status === "done"
              ? "Your proof is now on-chain and permanent."
              : status === "error"
              ? "Check the details below."
              : "Anchoring your image to the blockchain"}
          </p>
        </div>

        {/* Progress bar */}
        <div className="modal-progress-track">
          <div
            className={`modal-progress-fill ${status}`}
            style={{ width: `${status === "error" ? progress : status === "done" ? 100 : progress}%` }}
          />
        </div>

        {/* Steps */}
        <div className="modal-steps">
          {steps.map((step, i) => {
            const isDone = currentStep > i || status === "done";
            const isActive = currentStep === i && status !== "done" && status !== "error";
            const isWaiting = currentStep < i;
            return (
              <div key={step.key} className={`modal-step ${isDone ? "done" : isActive ? "active" : "waiting"}`}>
                <div className="modal-step-indicator">
                  {isDone ? (
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
                  ) : isActive ? (
                    <div className="step-spinner" />
                  ) : (
                    <span className="step-num">{String(i + 1).padStart(2, "0")}</span>
                  )}
                </div>
                <span className="modal-step-label">{step.label}</span>
              </div>
            );
          })}
        </div>

        {/* Fun fact */}
        {status !== "done" && status !== "error" && (
          <div className="modal-fact">
            <span className="fact-tag">did you know?</span>
            <p>{FACTS[fact]}</p>
          </div>
        )}

        {/* Close button */}
        {(status === "done" || status === "error") && (
          <button className="btn-primary modal-close-btn" onClick={onClose}>
            {status === "done" ? "View Results" : "Dismiss"}
          </button>
        )}
      </div>
    </div>
  );
}
