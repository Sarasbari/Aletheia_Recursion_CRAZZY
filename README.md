# Aletheia Media Authenticity Protocol

A comprehensive, blockchain-powered ecosystem designed to anchor and verify media authenticity using cryptographic proofs on the Polygon Amoy Testnet. Aletheia eliminates deepfakes and media tampering by establishing an immutable chain of custody for digital media.

## 🚀 Ecosystem Architecture

Aletheia is composed of three primary components:

### 1. The Browser Extension (`/browser-extension`)
A lightweight, zero-friction Chrome extension that allows end-users to interact with the blockchain seamlessly.
- **Anchor (`+`)**: Uploads an image directly to the blockchain, computing exact SHA256 hashes and Merkle Roots while auto-paying Polygon gas fees via MetaMask.
- **Verify (`🔍`)**: Scans any image on the web and instantaneously matches its cryptographic hash against the blockchain registry to guarantee if it is authentic and un-tampered.

### 2. The Developer SDK & CLI (`/sdk`)
A powerful Node.js Command Line Interface for developers to integrate Aletheia programmatically.
- **WalletConnect Integration**: Pay gas fees seamlessly by scanning QR codes natively via your mobile wallet.
- **Direct Wallet Keys**: Fully automated CI/CD-friendly testing using standard `ethers` private keys via the `--wallet-key` flag.
- **Programmatic Verifiability**: Access Aletheia's hashing engine, upload tools, and authenticity checks locally.

### 3. The Protocol Backend (`/api`)
The core infrastructure routing algorithms spanning hardware-verified uploads and decentralized tracking.
- **Cryptography Engine**: Performs real-time validation of exact `SHA256` matching and fuzzy `pHash` (Perceptual Hashing) similarity searching.
- **Context & Temporal Tamper Detection**: Rejects forged capture timestamps and mismatched geographical signatures.

---

## 🛠️ Getting Started

### Running the Browser Extension
1. Open Google Chrome and navigate to `chrome://extensions`.
2. Enable **Developer mode** in the top right corner.
3. Click **Load unpacked** and select the `/browser-extension` folder in this repository.
4. Open any webpage with images, hover slightly, and you'll see the Aletheia overlay tools appear! Make sure you have the [MetaMask extension](https://metamask.io/) installed to pay testnet gas fees.

### Using the CLI SDK
The Aletheia SDK can be run locally via Node.

```bash
cd sdk
npm install

# Anchor an image using a direct Private Key (bypassing WalletConnect UI)
node cli.js upload-pay "C:/path/to/image.jpg" --base-url http://192.168.0.177:8080 --wallet-key <YOUR_PRIVATE_KEY>

# Verify an image using your mobile wallet via WalletConnect
node cli.js verify-pay "C:/path/to/image.jpg" --base-url http://192.168.0.177:8080
```

### Backend API Configuration
Ensure the core backend is running (typically binding to port `8080`). All extension and SDK connections point to this active server to establish blockchain persistence.

---

## 🔐 Security & Blockchain Specifications

- **Network**: Polygon Amoy Testnet (Chain ID: `80002`)
- **Native Currency**: `POL`
- **Tamper Evidence**: Evaluates `authenticity_score`, `merkle_root` integrity, and detects malicious replay transactions.

## 📄 License
MIT License.
