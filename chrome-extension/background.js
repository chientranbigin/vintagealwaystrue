const DEFAULT_CONFIG = {
  serverUrl: "https://vintage-sale-go.onrender.com",
  secretKey: "",
  orderTrigger: "OK ạ, để em lên đơn ngay",
  autoSync: true
};

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.type === "COPILOT_INGEST") {
    handleIngest(request.payload)
      .then((data) => sendResponse({ ok: true, data }))
      .catch((err) => sendResponse({ ok: false, error: err.message || String(err) }));
    return true; // keep channel open for async response
  }

  if (request.type === "FETCH_IMAGE_BASE64") {
    fetchImageAsDataUrl(request.url)
      .then((dataUrl) => sendResponse({ ok: true, dataUrl }))
      .catch((err) => sendResponse({ ok: false, error: err.message || String(err) }));
    return true;
  }
});

async function getConfig() {
  const stored = await chrome.storage.local.get(Object.keys(DEFAULT_CONFIG));
  return { ...DEFAULT_CONFIG, ...stored };
}

async function handleIngest(payload) {
  const config = await getConfig();
  const baseUrl = (config.serverUrl || DEFAULT_CONFIG.serverUrl).replace(/\/+$/, "");
  const endpoint = `${baseUrl}/api/copilot/ingest`;

  const headers = {
    "Content-Type": "application/json"
  };
  if (config.secretKey) {
    headers["X-Copilot-Secret"] = config.secretKey;
  }

  const resp = await fetch(endpoint, {
    method: "POST",
    headers,
    body: JSON.stringify(payload)
  });

  if (!resp.ok) {
    const errText = await resp.text().catch(() => "");
    throw new Error(`HTTP ${resp.status}: ${errText.slice(0, 200)}`);
  }

  return resp.json().catch(() => ({ status: "ok" }));
}

async function fetchImageAsDataUrl(url) {
  let resp = await fetch(url);
  if (!resp.ok) {
    resp = await fetch(url, { credentials: "include" });
  }
  if (!resp.ok) {
    throw new Error(`Image fetch failed: ${resp.status}`);
  }
  const blob = await resp.blob();
  const buffer = await blob.arrayBuffer();
  const bytes = new Uint8Array(buffer);
  let binary = "";
  const chunkSize = 0x8000;
  for (let i = 0; i < bytes.length; i += chunkSize) {
    binary += String.fromCharCode.apply(null, bytes.subarray(i, i + chunkSize));
  }
  const base64 = btoa(binary);
  const mime = blob.type || "image/jpeg";
  return `data:${mime};base64,${base64}`;
}
