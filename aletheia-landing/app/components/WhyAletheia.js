"use client";

import { useEffect, useRef } from "react";

const FEATURES = [
  { title: "Tamper-proof Hashing", desc: "Dual-hash approach: SHA-512 for bit-level integrity, pHash for perceptual similarity. Change a single pixel and the math breaks.", detail: "SHA-512 + pHash → dual verification" },
  { title: "Decentralized Storage", desc: "IPFS means no single server owns your proof. Content-addressed, censorship-resistant, globally accessible. Your data outlives us.", detail: "IPFS CID → permanent retrieval" },
  { title: "On-chain Proof", desc: "Merkle tree batching on Polygon. Immutable timestamps anyone can verify, forever, for fractions of a cent.", detail: "Polygon tx → immutable timestamp" },
];

export default function WhyAletheia() {
  const sectionRef = useRef(null);

  useEffect(() => {
    const section = sectionRef.current;
    if (!section) return;
    const elements = section.querySelectorAll(".reveal, .feature-card");
    const obs = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            if (entry.target.classList.contains("feature-card")) {
              const cards = section.querySelectorAll(".feature-card");
              cards.forEach((c, i) => setTimeout(() => c.classList.add("visible"), i * 140));
              obs.unobserve(entry.target);
            } else {
              entry.target.classList.add("visible");
              obs.unobserve(entry.target);
            }
          }
        });
      },
      { threshold: 0.1 }
    );
    elements.forEach((el) => obs.observe(el));
    return () => obs.disconnect();
  }, []);

  return (
    <section className="why-aletheia" id="why-aletheia" ref={sectionRef}>
      <div className="container">
        <span className="section-label">why this exists</span>
        <h2 className="section-title reveal">Why ALETHEIA</h2>
        <p className="section-subtitle reveal">Cryptographic guarantees at every layer — from pixel to proof.</p>
        <div className="features-grid">
          {FEATURES.map((f) => (
            <div className="feature-card" key={f.title}>
              <h3>{f.title}</h3>
              <p>{f.desc}</p>
              <div className="card-detail">{f.detail}</div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
