"use client";

import MockupPanel from "./MockupPanel";
import PipelineTicker from "./PipelineTicker";

export default function Hero() {
  return (
    <section className="hero" id="hero">
      <div className="hero-glow green" />
      <div className="hero-glow amber" />
      <div className="container">
        <div className="hero-content">
          <div className="hero-text">
            <div className="badge">
              <span className="dot" /> built on polygon
            </div>
            <h1>
              <span className="thin">Every image</span>{" "}
              <span className="heavy">tells a truth.</span>
              <span className="line2">
                <span className="highlight">We make it permanent.</span>
              </span>
            </h1>
            <p className="hero-description">
              We hash it. We sign it. We put it on-chain.
              <br />
              Nobody can argue with math.
            </p>
            <div className="hero-cta-group">
              <a href="/anchor" className="btn-primary" id="ctaAnchor">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                  <polyline points="17 8 12 3 7 8" />
                  <line x1="12" y1="3" x2="12" y2="15" />
                </svg>
                Anchor an Image
              </a>
              <a href="#how-it-works" className="btn-secondary">
                See How It Works
              </a>
            </div>
            <div className="hero-stats">
              <div className="hero-stat">
                <span className="value">SHA-512</span>
                <span className="label">Content Hash</span>
              </div>
              <div className="hero-stat">
                <span className="value">ECDSA</span>
                <span className="label">Digital Signature</span>
              </div>
              <div className="hero-stat">
                <span className="value">Polygon</span>
                <span className="label">On-chain Anchor</span>
              </div>
            </div>
          </div>
          <MockupPanel />
        </div>
        <PipelineTicker />
      </div>
    </section>
  );
}
