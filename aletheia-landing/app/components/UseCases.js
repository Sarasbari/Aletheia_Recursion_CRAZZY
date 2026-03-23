"use client";

import { useEffect, useRef } from "react";

const CASES = [
  {
    icon: "📸",
    title: "Photojournalism",
    desc: "Prove when and where a photograph was taken. Anchor raw images at the point of capture — unbroken chain of custody for editorial integrity.",
  },
  {
    icon: "⚖️",
    title: "Legal Evidence",
    desc: "Cryptographically verifiable visual evidence. Timestamps, authorship, content integrity — independently auditable in any jurisdiction.",
  },
  {
    icon: "🎨",
    title: "NFT Provenance",
    desc: "Prove your art existed at a specific time, authored by a specific key — before any marketplace gets involved.",
  },
  {
    icon: "🔬",
    title: "Scientific Research",
    desc: "Anchor microscopy, field photos, experimental data with immutable timestamps. Protect against fabrication claims.",
  },
];

export default function UseCases() {
  const sectionRef = useRef(null);

  useEffect(() => {
    const section = sectionRef.current;
    if (!section) return;

    const elements = section.querySelectorAll(".reveal, .case-card");

    const obs = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            if (entry.target.classList.contains("case-card")) {
              const cards = section.querySelectorAll(".case-card");
              cards.forEach((c, i) =>
                setTimeout(() => c.classList.add("visible"), i * 110)
              );
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
    <section className="use-cases" id="use-cases" ref={sectionRef}>
      <div className="container">
        <span className="section-label">who needs this</span>
        <h2 className="section-title reveal">Built for Proof</h2>
        <p className="section-subtitle reveal">
          From warzones to courtrooms — wherever image truth matters.
        </p>

        <div className="cases-grid">
          {CASES.map((c) => (
            <div className="case-card" key={c.title}>
              <div className="case-icon">{c.icon}</div>
              <div>
                <h3>{c.title}</h3>
                <p>{c.desc}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
