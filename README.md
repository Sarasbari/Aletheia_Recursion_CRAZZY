# 🌐 Aletheia Verifier — Browser Extension

> **Real-time Image Authenticity Verification on the Web**

---

## 🚨 Problem

While browsing news or social media, users encounter images that may be:

* ❌ AI-generated (deepfakes)
* ❌ Manipulated or edited
* ❌ Misused to spread misinformation

There is **no instant way to verify authenticity directly on the web**.

---

## 💡 Solution

**Aletheia Verifier** is a browser extension that allows users to:

* 🔍 Verify any image directly on a webpage
* 🤖 Detect AI-generated or manipulated content
* 🔗 Check blockchain-backed authenticity
* 📊 View a real-time **Trust Score (0–100)**

---

## ⚙️ How It Works

1. Extension scans webpage for images
2. User clicks **“Verify”** on any image
3. Image is sent to backend
4. System performs:

   * SHA-256 hashing
   * Perceptual hashing (pHash)
   * AI deepfake detection
   * Blockchain verification
5. Result displayed instantly

---

## 🔄 Flow Diagram

```text
Web Image → Extension → Backend API
        → Hash + AI + Blockchain Check
        → Trust Score → Result Display
```

---

## 🔑 Features

### 🔍 One-Click Verification

Verify any image directly from browser

### 📊 Trust Score System

* 0–100 authenticity score
* Combines cryptography + AI

### 🤖 AI Deepfake Detection

Detects AI-generated or synthetic images

### 🔗 Blockchain Verification

Checks immutable proof stored on-chain

### 🟥 Visual Indicators

* Green → Authentic
* Red → Tampered
* Yellow → Suspicious

### 📦 Optional Upload

Register new images to blockchain

---

## 🏗️ Architecture

```text
Browser Extension
   ↓
Content Script (Image Detection)
   ↓
Background Script (API Calls)
   ↓
Backend (Hash + AI + Blockchain)
   ↓
Result → Popup UI / Overlay
```

---

## 🛠️ Tech Stack

| Layer       | Technology                      |
| ----------- | ------------------------------- |
| Extension   | JavaScript (Chrome Manifest v3) |
| Frontend UI | HTML, CSS, JS                   |
| Backend API | Go / Node.js                    |
| AI Service  | Python (FastAPI)                |
| Storage     | IPFS                            |
| Blockchain  | Polygon                         |
| Hashing     | SHA-256, pHash                  |

---

## 📁 Project Structure



---

## 🚀 Installation

1. Clone the repository:

```bash
git clone https://github.com/your-username/aletheia-extension
cd aletheia-extension
```

2. Open Chrome → Extensions
3. Enable **Developer Mode**
4. Click **Load Unpacked**
5. Select the `extension/` folder

---

## 🧪 Usage

* Open any website (news/social media)
* Click **Verify** button on an image
* View:

  * Trust Score
  * Authenticity status
  * AI detection result

---

## 🎯 Use Cases

* 📰 Journalists verifying news images
* 🌐 Users checking viral content
* 🛡️ Preventing misinformation
* 📱 Real-time media validation

---

## ⚡ Limitations

* Requires backend for verification
* AI detection is probabilistic (not 100%)
* New images must be registered for blockchain proof

---

## 🚀 Future Scope

* 🔄 Auto-scan all images on page
* 🌐 Social media integration
* 📱 Mobile browser support
* 🔐 Zero-knowledge proof verification
* 🧠 Improved AI detection models

---

## 🏆 Why This Matters

> “Seeing is no longer believing — verification is.”

Aletheia Verifier brings **trust directly
