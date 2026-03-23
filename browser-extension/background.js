chrome.runtime.onInstalled.addListener(() => {
  console.log("Aletheia Basic Extension installed.");
});

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message?.type === "PING") {
    sendResponse({ ok: true, from: "background" });
  }
});
