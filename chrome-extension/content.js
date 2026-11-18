(() => {
  const DEBOUNCE_MS = 1500;
  const MAX_MESSAGES = 25;
  const MAX_IMAGES = 4;

  let config = {
    serverUrl: "https://vintage-sale-go.onrender.com",
    secretKey: "",
    orderTrigger: "OK ạ, để em lên đơn ngay",
    autoSync: true
  };

  let debounceTimer = null;
  let lastPayloadHash = "";
  let lastTriggeredOrderKey = "";
  const imageCache = new Map(); // url -> dataUrl

  // Load config & listen for changes
  chrome.storage.local.get(["serverUrl", "secretKey", "orderTrigger", "autoSync"], (res) => {
    config = { ...config, ...res };
    initBadge();
    startObserver();
  });

  chrome.storage.onChanged.addListener((changes) => {
    for (const [k, v] of Object.entries(changes)) {
      config[k] = v.newValue;
    }
    updateBadgeStatus(config.autoSync ? "Sẵn sàng" : "Đang tạm dừng", "#64748b");
  });

  // --- UI Badge ---
  let badgeEl = null;
  let statusTextEl = null;

  function initBadge() {
    if (document.getElementById("vat-copilot-badge")) return;
    badgeEl = document.createElement("div");
    badgeEl.id = "vat-copilot-badge";
    badgeEl.style.cssText = [
      "position:fixed",
      "bottom:16px",
      "right:16px",
      "z-index:999999",
      "background:#0f172a",
      "color:#f8fafc",
      "padding:8px 12px",
      "border-radius:10px",
      "font-family:system-ui,-apple-system,sans-serif",
      "font-size:12px",
      "box-shadow:0 4px 16px rgba(0,0,0,0.25)",
      "display:flex",
      "align-items:center",
      "gap:8px",
      "border:1px solid #334155"
    ].join(";");

    statusTextEl = document.createElement("span");
    statusTextEl.textContent = "Copilot: Sẵn sàng";

    const syncBtn = document.createElement("button");
    syncBtn.textContent = "Sync ngay";
    syncBtn.style.cssText = [
      "background:#2563eb",
      "color:#fff",
      "border:none",
      "border-radius:6px",
      "padding:4px 8px",
      "font-size:11px",
      "font-weight:600",
      "cursor:pointer"
    ].join(";");
    syncBtn.addEventListener("click", () => syncCurrentConversation(true));

    badgeEl.appendChild(statusTextEl);
    badgeEl.appendChild(syncBtn);
    document.body.appendChild(badgeEl);
  }

  function updateBadgeStatus(text, color = "#38bdf8") {
    if (!statusTextEl) return;
    statusTextEl.textContent = `Copilot: ${text}`;
    statusTextEl.style.color = color;
  }

  // --- Keyboard shortcut Alt+C ---
  document.addEventListener("keydown", (e) => {
    if (e.altKey && (e.key === "c" || e.key === "C")) {
      e.preventDefault();
      syncCurrentConversation(true);
    }
  });

  // --- MutationObserver ---
  function startObserver() {
    const observer = new MutationObserver(() => {
      if (!config.autoSync) return;
      if (!location.href.includes("/inbox")) return;

      clearTimeout(debounceTimer);
      debounceTimer = setTimeout(() => {
        syncCurrentConversation(false);
      }, DEBOUNCE_MS);
    });

    observer.observe(document.body, {
      childList: true,
      subtree: true
    });
  }

  // --- DOM Extraction (Resilient to Meta CSS class changes) ---
  function normalizeText(str) {
    return (str || "")
      .replace(/\s+/g, " ")
      .trim()
      .toLowerCase();
  }

  function getMainChatContainer() {
    // Meta Business Suite Inbox usually renders the active conversation inside role="main"
    const main = document.querySelector('[role="main"]');
    if (main) return main;
    return document.body;
  }

  function extractCustomerName(mainEl) {
    // Look for prominent header text / link inside main pane
    const candidates = mainEl.querySelectorAll("h1, h2, h3, [role='heading'], a[role='link'] span");
    for (const el of candidates) {
      const txt = (el.textContent || "").trim();
      if (txt && txt.length >= 2 && txt.length <= 80 && !txt.includes("Meta Business")) {
        return txt;
      }
    }
    return "Khách Facebook";
  }

  function extractConversationThreadId() {
    const url = new URL(location.href);
    return url.searchParams.get("selected_item_id") || url.searchParams.get("thread_id") || location.pathname + location.search;
  }

  function extractRawMessages(mainEl) {
    const mainRect = mainEl.getBoundingClientRect();
    const centerX = mainRect.left + mainRect.width / 2;

    // Find message rows: in Meta Business Suite, message groups/rows live inside scrollable grids/lists
    // We inspect leaf text bubbles & product images inside the chat timeline
    const timelineRoot =
      mainEl.querySelector('[role="grid"]') ||
      mainEl.querySelector('[data-pagelet*="Inbox"]') ||
      mainEl;

    const nodes = timelineRoot.querySelectorAll('[dir="auto"], img');
    const items = [];
    const seenEls = new Set();

    for (const node of nodes) {
      if (seenEls.has(node)) continue;
      const rect = node.getBoundingClientRect();
      if (rect.width === 0 || rect.height === 0) continue;

      // Ignore elements in left sidebar or header top bar
      if (rect.top < mainRect.top + 55) continue;

      if (node.tagName === "IMG") {
        // Filter out small avatars, emojis, reaction icons
        if (rect.width < 60 || rect.height < 60) continue;
        const src = node.currentSrc || node.src || "";
        if (!src || src.startsWith("data:image/svg")) continue;
        // Skip stickers/emojis hosted on static emoji paths if tiny natural size
        if (node.naturalWidth > 0 && node.naturalWidth < 70) continue;

        const sender = rect.left + rect.width / 2 > centerX ? "shop" : "customer";
        items.push({
          sender,
          type: "image",
          text: "",
          imageUrl: src,
          top: rect.top
        });
        seenEls.add(node);
        continue;
      }

      // Text node ([dir="auto"])
      // Avoid picking up nested [dir="auto"] twice
      if (node.parentElement && node.parentElement.closest('[dir="auto"]') !== node) {
        // This node is inside another [dir="auto"]
        continue;
      }

      const text = (node.innerText || node.textContent || "").trim();
      if (!text) continue;

      // Skip UI timestamps / status labels
      if (/^(đã xem|seen|sent|đã gửi|nhập tin nhắn|type a message)$/i.test(text)) continue;

      // Heuristic: shop messages sit on the right half of main chat pane; customer messages sit on the left half
      const msgCenterX = rect.left + rect.width / 2;
      const sender = msgCenterX > centerX ? "shop" : "customer";

      items.push({
        sender,
        type: "text",
        text,
        imageUrl: "",
        top: rect.top
      });
      seenEls.add(node);
    }

    // Sort top-to-bottom by visual position
    items.sort((a, b) => a.top - b.top);
    return items.map(({ sender, type, text, imageUrl }) => ({ sender, type, text, imageUrl }));
  }

  /**
   * Apply Context Cutoff & Trigger Detection:
   * - Walk backwards from newest message.
   * - If newest shop message matches orderTrigger -> action = "create_order",
   *   and keep walking backwards until the PREVIOUS orderTrigger (or max limit).
   * - If an older orderTrigger is encountered -> stop immediately so previous week's orders are excluded.
   */
  function sliceSessionAndDetectAction(allMessages) {
    const triggerNorm = normalizeText(config.orderTrigger);
    if (allMessages.length === 0) {
      return { action: "suggest", messages: [] };
    }

    let action = "suggest";
    let startIdx = allMessages.length - 1;

    // Check if the latest message (or last 2 messages in case of echo/timestamp) is the Shop's order trigger
    const lastShopTriggerIdx = findRecentShopTrigger(allMessages, triggerNorm);
    if (lastShopTriggerIdx !== -1 && lastShopTriggerIdx >= allMessages.length - 2) {
      action = "create_order";
      startIdx = lastShopTriggerIdx - 1;
    }

    const sliced = [];
    for (let i = startIdx; i >= 0 && sliced.length < MAX_MESSAGES; i--) {
      const m = allMessages[i];
      if (m.sender === "shop" && m.type === "text" && triggerNorm && normalizeText(m.text).includes(triggerNorm)) {
        // Reached previous order boundary -> stop here!
        break;
      }
      sliced.unshift(m);
    }

    if (action === "create_order" && lastShopTriggerIdx !== -1) {
      sliced.push(allMessages[lastShopTriggerIdx]);
    }

    return { action, messages: sliced };
  }

  function findRecentShopTrigger(messages, triggerNorm) {
    if (!triggerNorm) return -1;
    for (let i = messages.length - 1; i >= Math.max(0, messages.length - 3); i--) {
      const m = messages[i];
      if (m.sender === "shop" && m.type === "text" && normalizeText(m.text).includes(triggerNorm)) {
        return i;
      }
    }
    return -1;
  }

  async function resolveImagesToBase64(messages) {
    // Collect up to MAX_IMAGES newest images (prioritize customer images)
    const imageIndices = [];
    for (let i = messages.length - 1; i >= 0; i--) {
      if (messages[i].type === "image" && messages[i].imageUrl) {
        imageIndices.push(i);
        if (imageIndices.length >= MAX_IMAGES) break;
      }
    }

    for (const idx of imageIndices) {
      const url = messages[idx].imageUrl;
      if (imageCache.has(url)) {
        messages[idx].imageBase64 = imageCache.get(url);
        continue;
      }
      try {
        const res = await chrome.runtime.sendMessage({
          type: "FETCH_IMAGE_BASE64",
          url
        });
        if (res && res.ok && res.dataUrl) {
          imageCache.set(url, res.dataUrl);
          messages[idx].imageBase64 = res.dataUrl;
        }
      } catch (err) {
        console.warn("[VAT Copilot] Image convert failed:", err);
      }
    }
    return messages;
  }

  async function syncCurrentConversation(force) {
    const mainEl = getMainChatContainer();
    const rawMessages = extractRawMessages(mainEl);
    if (rawMessages.length === 0) {
      if (force) updateBadgeStatus("Không thấy tin nhắn", "#fbbf24");
      return;
    }

    const customerName = extractCustomerName(mainEl);
    const threadId = extractConversationThreadId();
    const { action, messages } = sliceSessionAndDetectAction(rawMessages);

    if (messages.length === 0) return;

    const signature = JSON.stringify({
      threadId,
      customerName,
      action,
      msgs: messages.map((m) => `${m.sender}:${m.type}:${m.text}:${m.imageUrl}`)
    });

    if (!force && signature === lastPayloadHash) {
      return;
    }

    // Prevent duplicate create_order on the same trigger message
    const orderTriggerKey = `${threadId}::${messages.map((m) => m.text).join("|")}`;
    if (!force && action === "create_order" && orderTriggerKey === lastTriggeredOrderKey) {
      return;
    }

    updateBadgeStatus(action === "create_order" ? "Đang lên đơn..." : "Đang phân tích...", "#38bdf8");

    const enrichedMessages = await resolveImagesToBase64(messages);

    const payload = {
      thread_id: threadId,
      customer_name: customerName,
      action,
      captured_at: new Date().toISOString(),
      messages: enrichedMessages.map((m) => ({
        sender: m.sender,
        type: m.type,
        text: m.text,
        image_base64: m.imageBase64 || "",
        image_url: m.imageUrl || ""
      }))
    };

    chrome.runtime.sendMessage({ type: "COPILOT_INGEST", payload }, (resp) => {
      if (chrome.runtime.lastError) {
        updateBadgeStatus("Lỗi kết nối Ext", "#f87171");
        return;
      }
      if (resp && resp.ok) {
        lastPayloadHash = signature;
        if (action === "create_order") {
          lastTriggeredOrderKey = orderTriggerKey;
          updateBadgeStatus("Đã gửi lệnh lên đơn!", "#4ade80");
        } else {
          updateBadgeStatus(`Đã đồng bộ (${enrichedMessages.length} tin)`, "#4ade80");
        }
      } else {
        const msg = resp && resp.error ? resp.error : "Lỗi server";
        updateBadgeStatus(msg.slice(0, 30), "#f87171");
      }
    });
  }
})();
