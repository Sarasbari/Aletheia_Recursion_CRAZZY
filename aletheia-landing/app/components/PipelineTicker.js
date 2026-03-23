"use client";

import { useEffect, useRef } from "react";

const NODES = [
  "IMAGE_INPUT", "EXIF_PARSE", "SHA-512", "pHASH", "ECDSA_SIGN",
  "IPFS_PIN", "MERKLE_LEAF", "BATCH_ROOT", "TX_SUBMIT", "POLYGON_ANCHOR",
  "VERIFY_OK ✓",
];

export default function PipelineTicker() {
  const trackRef = useRef(null);
  const activeRef = useRef(0);

  useEffect(() => {
    const track = trackRef.current;
    if (!track) return;

    for (let copy = 0; copy < 3; copy++) {
      NODES.forEach((label, i) => {
        const node = document.createElement("div");
        node.className = "pipeline-node";
        const span = document.createElement("span");
        span.className = "node-label";
        span.textContent = label;
        node.appendChild(span);
        if (i < NODES.length - 1) {
          const conn = document.createElement("span");
          conn.className = "connector";
          node.appendChild(conn);
        }
        track.appendChild(node);
      });
      if (copy < 2) {
        const spacer = document.createElement("div");
        spacer.className = "pipeline-node";
        const conn = document.createElement("span");
        conn.className = "connector";
        spacer.appendChild(conn);
        track.appendChild(spacer);
      }
    }

    const interval = setInterval(() => {
      const labels = track.querySelectorAll(".node-label");
      labels.forEach((n) => n.classList.remove("active"));
      if (labels[activeRef.current]) labels[activeRef.current].classList.add("active");
      activeRef.current = (activeRef.current + 1) % labels.length;
    }, 700);

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="pipeline-visual">
      <div className="pipeline-track" ref={trackRef} />
    </div>
  );
}
