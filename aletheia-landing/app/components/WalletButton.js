"use client";

import { useAppKit, useAppKitAccount } from "@reown/appkit/react";

export default function WalletButton() {
  const { open } = useAppKit();
  const { address, isConnected } = useAppKitAccount();

  const truncated = address
    ? `${address.slice(0, 6)}...${address.slice(-4)}`
    : "";

  if (isConnected) {
    return (
      <button className="wallet-btn connected" onClick={() => open()}>
        <span className="wallet-dot" />
        {truncated}
      </button>
    );
  }

  return (
    <button className="wallet-btn" onClick={() => open()}>
      <svg
        width="14"
        height="14"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      >
        <rect x="1" y="4" width="22" height="16" rx="2" ry="2" />
        <line x1="1" y1="10" x2="23" y2="10" />
      </svg>
      Connect Wallet
    </button>
  );
}
