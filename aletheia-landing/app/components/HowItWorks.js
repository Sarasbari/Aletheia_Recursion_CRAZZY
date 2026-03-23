"use client";

import { useEffect, useRef } from "react";

const STEPS = [
  { num: "01", name: "Upload", desc: "Drag & drop your image with metadata — location, author, license, event context.", hash: null },
  { num: "02", name: "Hash", desc: "Generate SHA-512 content hash and pHash perceptual fingerprint.", hash: "sha512: 3a7f2b..." },
  { num: "03", name: "Sign", desc: "ECDSA digital signature binds your identity to the hash.", hash: "sig: 0x4e8c..." },
  { num: "04", name: "Store", desc: "Image + metadata pinned to IPFS. Content-addressed, permanent.", hash: "cid: Qm7xR..." },
  { num: "05", name: "Anchor", desc: "Merkle root batch-anchored to Polygon. On-chain, forever.", hash: "tx: 0xa3f1...polygon" },
];

const ICONS = {
  Upload: <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>,
  Hash: <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round"><rect x="2" y="2" width="20" height="8" rx="2" ry="2"/><rect x="2" y="14" width="20" height="8" rx="2" ry="2"/><line x1="6" y1="6" x2="6.01" y2="6"/><line x1="6" y1="18" x2="6.01" y2="18"/></svg>,
  Sign: <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/><path d="M9 12l2 2 4-4"/></svg>,
  Store: <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/></svg>,
  Anchor: <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round"><rect x="1" y="4" width="22" height="16" rx="2" ry="2"/><line x1="1" y1="10" x2="23" y2="10"/></svg>,
};

export default function HowItWorks() {
  const sectionRef = useRef(null);

  useEffect(() => {
    const section = sectionRef.current;
    if (!section) return;

    const elements = section.querySelectorAll(".reveal, .step-item");
    const obs = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            const el = entry.target;
            const delay = el.dataset.step ? (parseInt(el.dataset.step) - 1) * 140 : 0;
            setTimeout(() => el.classList.add("visible"), delay);
            obs.unobserve(el);
          }
        });
      },
      { threshold: 0.12, rootMargin: "0px 0px -40px 0px" }
    );
    elements.forEach((el) => obs.observe(el));
    return () => obs.disconnect();
  }, []);

  return (
    <section className="how-it-works" id="how-it-works" ref={sectionRef}>
      <div className="container">
        <span className="section-label">how it works ↓</span>
        <h2 className="section-title reveal">Five Steps to Immutable Proof</h2>
        <p className="section-subtitle reveal">From raw image to on-chain anchor — every step is cryptographically verifiable.</p>
        <div className="steps-grid">
          {STEPS.map((step, i) => (
            <div className="step-item" data-step={step.num} key={step.name}>
              <span className="step-number">{step.num}</span>
              <div className="step-icon">
                <div className="icon-glow" />
                {ICONS[step.name]}
              </div>
              <h3 className="step-name">{step.name}</h3>
              <p className="step-desc">{step.desc}</p>
              {step.hash && <span className="step-hash">{step.hash}</span>}
              {i < STEPS.length - 1 && <div className="step-connector">›</div>}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
