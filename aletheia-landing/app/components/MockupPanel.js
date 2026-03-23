"use client";

import { useEffect, useState } from "react";

const STEPS = [
  { label: "SHA-512", value: "3a7f2b9e...c41d", initial: "done" },
  { label: "pHash", value: "a8b53831...", initial: "done" },
  { label: "ECDSA", value: "0x4e8c93...", initial: "done" },
  { label: "IPFS", value: "Qm7xR3...", initial: "pending" },
  { label: "Anchor", value: "0xa3f1...", initial: "waiting" },
];

export default function MockupPanel() {
  const [stepStates, setStepStates] = useState(STEPS.map((s) => s.initial));

  useEffect(() => {
    let idx = 0;
    const interval = setInterval(() => {
      setStepStates(
        STEPS.map((_, i) => {
          if (i <= idx) return "done";
          if (i === idx + 1) return "pending";
          return "waiting";
        })
      );
      idx = (idx + 1) % (STEPS.length + 1);
    }, 1800);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="hero-mockup">
      <div className="mockup-wrapper">
        <div className="mockup-bar">
          <span className="mockup-dot" />
          <span className="mockup-dot" />
          <span className="mockup-dot" />
          <span className="bar-title">aletheia.app</span>
        </div>
        <div className="mockup-body">
          <div className="mockup-panel">
            <div className="mockup-panel-title">input source</div>
            <div className="drop-zone">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                <polyline points="17 8 12 3 7 8" />
                <line x1="12" y1="3" x2="12" y2="15" />
              </svg>
              drop image here
            </div>
            <div className="meta-fields">
              <div className="meta-field"><span className="field-label">loc:</span> 28.6139° N, 77.2090° E</div>
              <div className="meta-field"><span className="field-label">author:</span> d.sharma</div>
              <div className="meta-field"><span className="field-label">license:</span> CC BY-SA 4.0</div>
              <div className="meta-field"><span className="field-label">event:</span> field_survey_03</div>
            </div>
          </div>
          <div className="mockup-panel">
            <div className="mockup-panel-title">crypto pipeline</div>
            <div className="pipeline-steps">
              {STEPS.map((step, i) => (
                <div key={step.label}>
                  <div className={`pipe-step ${stepStates[i]}`}>
                    <span className="status-dot" />
                    <span className="step-label">{step.label}</span>
                    <span className="step-value">
                      {stepStates[i] === "done" ? step.value : stepStates[i] === "pending" ? "processing..." : "waiting"}
                    </span>
                  </div>
                  {i < STEPS.length - 1 && <div className="pipe-connector" />}
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
