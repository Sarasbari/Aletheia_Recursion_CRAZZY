"use client";

import { useEffect, useRef } from "react";

const FEATURES = [
  "SHA-512 + pHash generation",
  "ECDSA signing",
  "IPFS decentralized storage",
  "Polygon on-chain anchor",
  "Merkle proof generation",
  "Verification API access",
];

export default function Pricing() {
  const sectionRef = useRef(null);

  useEffect(() => {
    const section = sectionRef.current;
    if (!section) return;

    const elements = section.querySelectorAll(".reveal");
    const obs = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            entry.target.classList.add("visible");
            obs.unobserve(entry.target);
          }
        });
      },
      { threshold: 0.12 }
    );
    elements.forEach((el) => obs.observe(el));
    return () => obs.disconnect();
  }, []);

  return (
    <section className="pricing" id="pricing" ref={sectionRef}>
      <div className="container">
        <span className="section-label">what it costs</span>
        <h2 className="section-title reveal">Gas-Optimized Anchoring</h2>
        <p className="section-subtitle reveal">
          Merkle tree batching keeps costs negligible.
        </p>

        <div className="pricing-card reveal">
          <div className="pricing-amount">
            <span className="symbol">~</span>
            <span className="value">0.0001</span>
            <span className="unit">MATIC / image</span>
          </div>
          <p className="pricing-label">
            Batch anchoring via Merkle tree keeps costs minimal
          </p>

          <div className="pricing-features">
            {FEATURES.map((f) => (
              <div className="pricing-feature" key={f}>
                <span className="check">✓</span>
                {f}
              </div>
            ))}
          </div>

          <div className="pricing-cta">
            <a href="#" className="btn-primary">
              Start Anchoring
            </a>
          </div>
          <p className="pricing-note">
            Gas costs at current Polygon rates. Actual costs may vary.
          </p>
        </div>

        <div className="footnote">
          <p>
            <span className="fn-marker">*</span> Built because we got tired of
            fake images winning arguments.
          </p>
        </div>
      </div>
    </section>
  );
}
