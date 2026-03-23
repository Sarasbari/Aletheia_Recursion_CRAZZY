# Basic Browser Extension Setup

This is a minimal Manifest V3 browser extension starter.

## Included files

- `manifest.json` - Extension configuration
- `popup.html`, `popup.css`, `popup.js` - Popup UI and logic
- `background.js` - Service worker background script
- `content.js` - Content script injected into pages

## Run in Chrome/Edge

1. Open extensions page:
   - Chrome: `chrome://extensions`
   - Edge: `edge://extensions`
2. Enable **Developer mode**.
3. Click **Load unpacked**.
4. Select this folder: `browser-extension`.
5. Click the extension icon and press **Highlight This Page**.

## Next steps

- Add custom icons in an `icons/` folder and register them in `manifest.json`.
- Replace sample highlight behavior with your real feature logic.
