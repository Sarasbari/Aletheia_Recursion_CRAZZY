let isHighlighted = false;

chrome.runtime.onMessage.addListener((message) => {
  if (message?.type !== "TOGGLE_HIGHLIGHT") {
    return;
  }

  isHighlighted = !isHighlighted;
  if (isHighlighted) {
    document.documentElement.style.outline = "4px solid #0ea5e9";
    document.documentElement.style.outlineOffset = "-4px";
  } else {
    document.documentElement.style.outline = "";
    document.documentElement.style.outlineOffset = "";
  }
});
