# Aletheia Chrome Extension (Frontend-Only)

Aletheia is a Manifest V3 browser extension that analyzes images on webpages and shows authenticity indicators.

## Features

- Detects images on initial load and dynamically added images
- Shows hover overlay button (🔍) on images
- Computes SHA-256 hash using Web Crypto API (in background service worker)
- Calls external APIs using `fetch()`:
  - `POST /analyze`
  - `POST /verify`
  - `POST /anchor`
- Displays floating in-page result card near the image
- Includes browser action popup with latest analysis
- Bonus: right-click context menu item `Verify Image`
- Bonus: result cache to avoid repeated API calls

## File Structure

- `manifest.json`
- `content.js`
- `background.js`
- `popup.html`
- `popup.js`
- `styles.css`

## Load in Chrome

1. Open `chrome://extensions`.
2. Enable **Developer mode**.
3. Click **Load unpacked**.
4. Select this folder:
   - `browser-extension`
5. Open any webpage with images.
6. Hover an image and click the 🔍 icon.

## API Base URL

Update the API base URL in `background.js`:

- `API_BASE_URL = "https://api.example.com"`

Replace with your actual API host.
