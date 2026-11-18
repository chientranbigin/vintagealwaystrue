(() => {
  const DEBOUNCE_MS = 2200;
  const MAX_MESSAGES = 30;
  const MAX_IMAGES = 6;

  let config = {
    serverUrl: "https://vintage-sale-go.onrender.com",
    secretKey: "",
    orderTrigger: "OK ạ, để em lên đơn ngay",
    autoSync: true
  };

  let debounceTimer = null;
  let lastPayloadHash = "";
  let lastTriggeredOrderKey = "";
  let isSyncing = false;
  const imageCache = new Map();

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

  document.addEventListener("keydown", (e) => {
    if (e.altKey && (e.key === "c" || e.key === "C")) {
      e.preventDefault();
      syncCurrentConversation(true);
    }
  });

  function startObserver() {
    const observer = new MutationObserver((mutations) => {
      if (!config.autoSync || isSyncing) return;
      if (!location.href.includes("/inbox")) return;

      // Ignore mutations caused by our own badge
      const isOnlyBadge = mutations.every(
        (m) => badgeEl && (badgeEl.contains(m.target) || m.target === badgeEl)
      );
      if (isOnlyBadge) return;

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

  function normalizeText(str) {
    return (str || "").replace(/\s+/g, " ").trim().toLowerCase();
  }

  function findChatPaneAndComposer() {
    const composers = Array.from(
      document.querySelectorAll('[contenteditable="true"][role="textbox"], textarea, [contenteditable="true"]')
    ).filter((el) => {
      const r = el.getBoundingClientRect();
      return r.width > 120 && r.height > 15 && r.top > window.innerHeight * 0.45 && r.left > window.innerWidth * 0.15;
    });

    const composer = composers[composers.length - 1] || null;
    if (composer) {
      let cur = composer.parentElement;
      while (cur && cur !== document.body) {
        const r = cur.getBoundingClientRect();
        if (r.height > 350 && r.width > 260 && r.width < window.innerWidth * 0.78) {
          const hasHeaderOrMessages = cur.querySelectorAll("img, span, div").length > 15;
          if (hasHeaderOrMessages) {
            return { chatPane: cur, composer };
          }
        }
        cur = cur.parentElement;
      }
    }

    const main = document.querySelector('[role="main"]') || document.body;
    return { chatPane: main, composer };
  }

  function isGarbageName(txt) {
    if (!txt || txt.length < 2 || txt.length > 65) return true;
    return /^(hộp thư|inbox|tất cả|messenger|meta business|trả lời|reply|tìm kiếm|search|chưa đọc|unread|xem trang cá nhân|view profile|bình luận|comments|tin nhắn|messages)$/i.test(
      txt.trim()
    );
  }

  function extractCustomerName(chatPane) {
    const paneRect = chatPane.getBoundingClientRect();
    const minLeft = Math.max(paneRect.left, window.innerWidth * 0.18);
    const maxRight = Math.min(paneRect.right, window.innerWidth * 0.78);

    // 1. Search header bar of the center chat pane (usually top 45px..190px of viewport)
    const headerEls = document.querySelectorAll("h1, h2, h3, [role='heading'], a[role='link'] span, strong");
    let bestName = "";
    let bestFontSize = 0;

    for (const el of headerEls) {
      if (badgeEl && badgeEl.contains(el)) continue;
      const r = el.getBoundingClientRect();
      if (r.width === 0 || r.height === 0) continue;
      if (r.left < minLeft || r.right > maxRight) continue;
      if (r.top < 40 || r.top > 180) continue;

      const txt = (el.textContent || "").trim();
      if (isGarbageName(txt)) continue;

      const fs = parseFloat(getComputedStyle(el).fontSize || "12");
      if (fs >= bestFontSize) {
        bestFontSize = fs;
        bestName = txt;
      }
    }
    if (bestName) return bestName;

    // 2. Fallback: active selected conversation in left sidebar
    const selected = document.querySelector('[aria-selected="true"], [aria-current="page"], [aria-current="true"]');
    if (selected) {
      const spans = selected.querySelectorAll("span, strong");
      for (const sp of spans) {
        const txt = (sp.textContent || "").trim();
        if (!isGarbageName(txt) && !/\d{1,2}:\d{2}/.test(txt)) {
          return txt;
        }
      }
    }

    // 3. Fallback: right customer info panel heading
    const rightHeaders = document.querySelectorAll("h2, h3, [role='heading']");
    for (const el of rightHeaders) {
      const r = el.getBoundingClientRect();
      if (r.left > window.innerWidth * 0.7 && r.top > 40 && r.top < 250) {
        const txt = (el.textContent || "").trim();
        if (!isGarbageName(txt)) return txt;
      }
    }

    return "Khách Facebook";
  }

  function extractConversationThreadId() {
    const url = new URL(location.href);
    return (
      url.searchParams.get("selected_item_id") ||
      url.searchParams.get("thread_id") ||
      url.searchParams.get("asset_id") ||
      location.pathname + location.search
    );
  }

  function extractRawMessages(chatPane, composer) {
    const paneRect = chatPane.getBoundingClientRect();
    const centerX = paneRect.left + paneRect.width / 2;
    const topBound = paneRect.top + 48;
    const bottomBound = composer ? composer.getBoundingClientRect().top - 6 : paneRect.bottom - 50;

    const items = [];
    const seenTexts = new Set();
    const seenImgSrcs = new Set();

    // 1. Extract product images (both <img> and CSS background-image)
    const imgs = chatPane.querySelectorAll("img");
    for (const img of imgs) {
      const r = img.getBoundingClientRect();
      if (r.width < 55 || r.height < 55) continue;
      if (r.top < topBound || r.bottom > bottomBound + 20) continue;
      const src = img.currentSrc || img.src || "";
      if (!src || src.startsWith("data:image/svg")) continue;
      if (img.naturalWidth > 0 && img.naturalWidth < 65) continue;
      if (seenImgSrcs.has(src)) continue;
      seenImgSrcs.add(src);

      const sender = r.left + r.width / 2 > centerX ? "shop" : "customer";
      items.push({
        sender,
        type: "image",
        text: "",
        imageUrl: src,
        imgElement: img,
        top: r.top
      });
    }

    // 2. Extract text messages
    const walker = document.createTreeWalker(chatPane, NodeFilter.SHOW_TEXT, null);
    let textNode;
    while ((textNode = walker.nextNode())) {
      const raw = (textNode.nodeValue || "").trim();
      if (!raw) continue;
      const parent = textNode.parentElement;
      if (!parent) continue;
      if (composer && composer.contains(parent)) continue;
      if (parent.closest("button, [role='button'], script, style, svg")) continue;

      const r = parent.getBoundingClientRect();
      if (r.width === 0 || r.height === 0) continue;
      if (r.top < topBound || r.bottom > bottomBound) continue;

      if (
        /^(đã xem|seen|sent|đã gửi|đã nhận|delivered|nhập tin nhắn|type a message|trả lời|reply|bạn đã gửi|chưa đọc|unread|\d{1,2}:\d{2})$/i.test(
          raw
        )
      ) {
        continue;
      }

      const distFromCenter = Math.abs(r.left + r.width / 2 - centerX);
      if (distFromCenter < 35 && r.width < 150 && raw.length < 28) {
        continue;
      }

      const bubbleEl = parent.closest('[dir="auto"]') || parent;
      const bRect = bubbleEl.getBoundingClientRect();
      const fullText = (bubbleEl.innerText || bubbleEl.textContent || "").trim();
      if (!fullText) continue;

      const dedupKey = `${Math.round(bRect.top)}:${fullText}`;
      if (seenTexts.has(dedupKey)) continue;
      seenTexts.add(dedupKey);

      const sender = bRect.left + bRect.width / 2 > centerX ? "shop" : "customer";
      items.push({
        sender,
        type: "text",
        text: fullText,
        imageUrl: "",
        top: bRect.top
      });
    }

    items.sort((a, b) => a.top - b.top);
    return items.map(({ sender, type, text, imageUrl, imgElement }) => ({
      sender,
      type,
      text,
      imageUrl,
      imgElement
    }));
  }

  function sliceSessionAndDetectAction(allMessages) {
    const triggerNorm = normalizeText(config.orderTrigger);
    if (allMessages.length === 0) {
      return { action: "suggest", messages: [] };
    }

    let action = "suggest";
    let startIdx = allMessages.length - 1;

    const lastShopTriggerIdx = findRecentShopTrigger(allMessages, triggerNorm);
    if (lastShopTriggerIdx !== -1 && lastShopTriggerIdx >= allMessages.length - 2) {
      action = "create_order";
      startIdx = lastShopTriggerIdx - 1;
    }

    const sliced = [];
    for (let i = startIdx; i >= 0 && sliced.length < MAX_MESSAGES; i--) {
      const m = allMessages[i];
      if (m.sender === "shop" && m.type === "text" && triggerNorm && normalizeText(m.text).includes(triggerNorm)) {
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

  function compressDataUrl(dataUrl, maxDim = 900) {
    return new Promise((resolve) => {
      const img = new Image();
      img.onload = () => {
        const scale = Math.min(1, maxDim / Math.max(img.width, img.height));
        const w = Math.max(1, Math.round(img.width * scale));
        const h = Math.max(1, Math.round(img.height * scale));
        const canvas = document.createElement("canvas");
        canvas.width = w;
        canvas.height = h;
        const ctx = canvas.getContext("2d");
        ctx.drawImage(img, 0, 0, w, h);
        resolve(canvas.toDataURL("image/jpeg", 0.82));
      };
      img.onerror = () => resolve(dataUrl);
      img.src = dataUrl;
    });
  }

  async function convertImageToBase64(msg) {
    const url = msg.imageUrl;
    if (!url) return "";
    if (url.startsWith("data:image/")) return url;
    if (imageCache.has(url)) return imageCache.get(url);

    // Strategy 1: Fetch via background service worker (bypasses page CORS)
    try {
      const res = await chrome.runtime.sendMessage({
        type: "FETCH_IMAGE_BASE64",
        url
      });
      if (res && res.ok && res.dataUrl) {
        const compressed = await compressDataUrl(res.dataUrl);
        imageCache.set(url, compressed);
        return compressed;
      }
    } catch (e) {
      console.warn("[VAT Copilot] Background fetch failed:", e);
    }

    // Strategy 2: Direct fetch in page context (works for blob: URLs)
    try {
      const resp = await fetch(url);
      if (resp.ok) {
        const blob = await resp.blob();
        const rawDataUrl = await new Promise((resolve) => {
          const reader = new FileReader();
          reader.onloadend = () => resolve(reader.result || "");
          reader.readAsDataURL(blob);
        });
        if (rawDataUrl) {
          const compressed = await compressDataUrl(rawDataUrl);
          imageCache.set(url, compressed);
          return compressed;
        }
      }
    } catch (e) {
      console.warn("[VAT Copilot] Direct fetch failed:", e);
    }

    return "";
  }

  async function resolveImagesToBase64(messages) {
    const imageIndices = [];
    for (let i = messages.length - 1; i >= 0; i--) {
      if (messages[i].type === "image" && messages[i].imageUrl) {
        imageIndices.push(i);
        if (imageIndices.length >= MAX_IMAGES) break;
      }
    }

    for (const idx of imageIndices) {
      messages[idx].imageBase64 = await convertImageToBase64(messages[idx]);
    }
    return messages;
  }

  async function syncCurrentConversation(force) {
    if (isSyncing) return;
    const { chatPane, composer } = findChatPaneAndComposer();
    const rawMessages = extractRawMessages(chatPane, composer);
    if (rawMessages.length === 0) {
      if (force) updateBadgeStatus("Không thấy tin nhắn", "#fbbf24");
      return;
    }

    const customerName = extractCustomerName(chatPane);
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

    const orderTriggerKey = `${threadId}::${messages.map((m) => m.text).join("|")}`;
    if (!force && action === "create_order" && orderTriggerKey === lastTriggeredOrderKey) {
      return;
    }

    isSyncing = true;
    updateBadgeStatus(action === "create_order" ? "Đang lên đơn..." : "Đang phân tích...", "#38bdf8");

    try {
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
        isSyncing = false;
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
            updateBadgeStatus(`Đã đồng bộ (${customerName})`, "#4ade80");
          }
        } else {
          const msg = resp && resp.error ? resp.error : "Lỗi server";
          updateBadgeStatus(msg.slice(0, 30), "#f87171");
        }
      });
    } catch (err) {
      isSyncing = false;
      updateBadgeStatus("Lỗi đọc hội thoại", "#f87171");
    }
  }
})();
