const DEFAULT_CONFIG = {
  serverUrl: "https://vintage-sale-go.onrender.com",
  secretKey: "",
  orderTrigger: "OK ạ, để em lên đơn ngay",
  domAutoSync: false
};

document.addEventListener("DOMContentLoaded", () => {
  chrome.storage.local.get(Object.keys(DEFAULT_CONFIG), (stored) => {
    const cfg = { ...DEFAULT_CONFIG, ...stored };
    document.getElementById("serverUrl").value = cfg.serverUrl;
    document.getElementById("secretKey").value = cfg.secretKey;
    document.getElementById("orderTrigger").value = cfg.orderTrigger;
    document.getElementById("domAutoSync").checked = Boolean(cfg.domAutoSync);
  });

  document.getElementById("saveBtn").addEventListener("click", () => {
    const newConfig = {
      serverUrl: document.getElementById("serverUrl").value.trim() || DEFAULT_CONFIG.serverUrl,
      secretKey: document.getElementById("secretKey").value.trim(),
      orderTrigger: document.getElementById("orderTrigger").value.trim() || DEFAULT_CONFIG.orderTrigger,
      domAutoSync: document.getElementById("domAutoSync").checked
    };

    chrome.storage.local.set(newConfig, () => {
      const msg = document.getElementById("saved-msg");
      msg.textContent = "Đã lưu cấu hình!";
      setTimeout(() => {
        msg.textContent = "";
      }, 2000);
    });
  });

  document.getElementById("openCopilotBtn").addEventListener("click", () => {
    const baseUrl = (document.getElementById("serverUrl").value.trim() || DEFAULT_CONFIG.serverUrl).replace(/\/+$/, "");
    chrome.tabs.create({ url: `${baseUrl}/salev2/copilot` });
  });
});
