"use client";

import { useCallback } from "react";

export default function Footer() {
  const handleCopy = useCallback((e) => {
    navigator.clipboard.writeText("0x71C7656EC7ab88b098defB751B7401B5f6d8976F");
    const el = e.currentTarget;
    el.style.borderColor = "rgba(74, 222, 128, 0.3)";
    el.style.color = "var(--accent)";
    setTimeout(() => {
      el.style.borderColor = "";
      el.style.color = "";
    }, 1500);
  }, []);

  return (
    <footer className="footer">
      <div className="container">
        <div className="footer-content">
          <div className="footer-brand">
            ALETHEIA<span>.</span>
          </div>
          <ul className="footer-links">
            <li><a href="#how-it-works">process</a></li>
            <li><a href="#why-aletheia">features</a></li>
            <li><a href="#use-cases">use cases</a></li>
            <li><a href="#">docs</a></li>
            <li><a href="#">github</a></li>
          </ul>
          <div
            className="footer-wallet"
            title="Click to copy"
            onClick={handleCopy}
          >
            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
            </svg>
            0x71C7...f3E8a9
          </div>
        </div>
        <div className="footer-bottom">
          <p>ALETHEIA — ἀλήθεια — &quot;truth, disclosure&quot; · © 2026</p>
        </div>
      </div>
    </footer>
  );
}
