const statusEl = document.getElementById("status");
const highlightBtn = document.getElementById("highlightBtn");

async function getCurrentTab() {
  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
  return tab;
}

highlightBtn.addEventListener("click", async () => {
  try {
    const tab = await getCurrentTab();
    if (!tab || !tab.id) {
      statusEl.textContent = "No active tab found.";
      return;
    }

    await chrome.tabs.sendMessage(tab.id, { type: "TOGGLE_HIGHLIGHT" });
    statusEl.textContent = "Toggled page highlight.";
  } catch (error) {
    statusEl.textContent = "Unable to connect to this page.";
  }
});
