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

    const verifyBtn = document.createElement("button");
    verifyBtn.textContent = "Verify";
    verifyBtn.title = "Highlight tin nhắn Extension đang đọc để đối chiếu với Facebook";
    verifyBtn.style.cssText = [
      "background:#334155",
      "color:#f8fafc",
      "border:none",
      "border-radius:6px",
      "padding:4px 8px",
      "font-size:11px",
      "font-weight:600",
      "cursor:pointer"
    ].join(";");
    verifyBtn.addEventListener("click", () => verifyVisualOverlay());

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
    badgeEl.appendChild(verifyBtn);
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
    if (!txt || txt.length < 2 || txt.length > 55) return true;
    const lower = txt.trim().toLowerCase();
    const exactBlacklist = /^(hộp thư|inbox|tất cả|messenger|meta business|trả lời|reply|tìm kiếm|search|chưa đọc|unread|xem trang cá nhân|view profile|bình luận|comments|tin nhắn|messages|người gửi|chỉ định cho tôi|giới thiệu|about|chi tiết|details|hoạt động|activity|xem bài viết|view post)$/i;
    if (exactBlacklist.test(lower)) return true;
    const partialBlacklist = [
      "bình luận trên facebook",
      "tin nhắn trên",
      "facebook comment",
      "messenger",
      "instagram",
      "meta business",
      "chỉ định",
      "tất cả tin nhắn",
      "đã trả lời",
      "xem bài viết",
      "đã bình luận",
      "replied to",
      "view post",
      "commented on",
      "trang cá nhân",
      "đã trả lời một quảng cáo",
      "mới nhất"
    ];
    return partialBlacklist.some((p) => lower.includes(p));
  }

  /**
   * Chat header name: the biggest text directly above the message list region
   * (above "Chỉ định cuộc trò chuyện này"). Ties resolve to the topmost element.
   */
  function extractHeaderName(region) {
    const rr = region.getBoundingClientRect();
    let best = "";
    let bestFs = 0;
    let bestTop = Infinity;
    const walker = document.createTreeWalker(document.body, NodeFilter.SHOW_TEXT, null);
    let node;
    while ((node = walker.nextNode())) {
      const txt = (node.nodeValue || "").trim();
      if (!txt) continue;
      const el = node.parentElement;
      if (!el || (badgeEl && badgeEl.contains(el))) continue;
      if (el.closest("button, [role='button'], [role='tab'], [role='region']")) continue;
      const r = el.getBoundingClientRect();
      if (r.width === 0 || r.height === 0) continue;
      if (r.left < rr.left - 10 || r.left > rr.left + rr.width * 0.6) continue;
      if (r.top < rr.top - 220 || r.bottom > rr.top) continue;
      if (isGarbageName(txt)) continue;
      const fs = parseFloat(getComputedStyle(el).fontSize || "12");
      if (fs < 15) continue;
      if (fs > bestFs || (fs === bestFs && r.top < bestTop)) {
        bestFs = fs;
        bestTop = r.top;
        best = txt;
      }
    }
    return best;
  }

  function extractCustomerName(chatPane) {
    const region = findMessageRegion();
    if (region) {
      const headerName = extractHeaderName(region);
      if (headerName) return headerName;
    }

    // 1. Right-hand Customer Info Panel (Profile header)
    const rightEls = document.querySelectorAll("h1, h2, h3, [role='heading'], a, strong, span");
    let rightBest = "";
    let rightMaxFs = 0;
    for (const el of rightEls) {
      if (badgeEl && badgeEl.contains(el)) continue;
      if (el.children.length > 1) continue;
      const r = el.getBoundingClientRect();
      if (r.width === 0 || r.height === 0) continue;
      if (r.left < window.innerWidth * 0.72) continue;
      if (r.top < 45 || r.top > 320) continue;

      const txt = (el.textContent || "").trim();
      if (isGarbageName(txt)) continue;

      const fs = parseFloat(getComputedStyle(el).fontSize || "12");
      if (fs >= 14 && fs > rightMaxFs) {
        rightMaxFs = fs;
        rightBest = txt;
      }
    }
    if (rightBest) return rightBest;

    // 2. Center chat pane header bar
    const paneRect = chatPane.getBoundingClientRect();
    const minLeft = Math.max(paneRect.left, window.innerWidth * 0.18);
    const maxRight = Math.min(paneRect.right, window.innerWidth * 0.75);

    const headerEls = document.querySelectorAll("h1, h2, h3, [role='heading'], a[role='link'] span, strong, span");
    let bestName = "";
    let bestFontSize = 0;

    for (const el of headerEls) {
      if (badgeEl && badgeEl.contains(el)) continue;
      if (el.children.length > 1) continue;
      const r = el.getBoundingClientRect();
      if (r.width === 0 || r.height === 0) continue;
      if (r.left < minLeft || r.right > maxRight) continue;
      if (r.top < 40 || r.top > 180) continue;

      const txt = (el.textContent || "").trim();
      if (isGarbageName(txt)) continue;

      const fs = parseFloat(getComputedStyle(el).fontSize || "12");
      if (fs >= 13 && fs >= bestFontSize) {
        bestFontSize = fs;
        bestName = txt;
      }
    }
    if (bestName) return bestName;

    // 3. Extract from "<Name> đã trả lời về một bài viết" if present
    const allTexts = chatPane.querySelectorAll("span, strong");
    for (const el of allTexts) {
      const txt = (el.textContent || "").trim();
      const m = txt.match(/^(.{2,45}?)\s+(?:đã trả lời|đã bình luận|replied to|commented on)\b/i);
      if (m && m[1] && !isGarbageName(m[1])) {
        return m[1].trim();
      }
    }

    // 4. Active selected conversation in left sidebar
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

  /**
   * Detect Meta Business Suite bottom banners (e.g. "Tin trả lời gợi ý", "Bạn có muốn được chỉ định...")
   * so we stop reading messages right above them.
   */
  function findBottomNoiseBoundary(chatPane, composer) {
    let bound = composer ? composer.getBoundingClientRect().top - 6 : chatPane.getBoundingClientRect().bottom - 50;
    const noisePhrases = [
      "tin trả lời gợi ý",
      "nhấp để chèn",
      "suggested replies",
      "bạn có muốn được chỉ định",
      "chỉ định cho tôi",
      "assign to me"
    ];

    const walker = document.createTreeWalker(chatPane, NodeFilter.SHOW_TEXT, null);
    let node;
    while ((node = walker.nextNode())) {
      const txt = normalizeText(node.nodeValue);
      if (!txt) continue;
      if (noisePhrases.some((p) => txt.includes(p))) {
        const el = node.parentElement;
        if (el) {
          const r = el.getBoundingClientRect();
          if (r.top > window.innerHeight * 0.35 && r.top < bound) {
            bound = r.top - 4;
          }
        }
      }
    }
    return bound;
  }

  /**
   * Collect containers that represent Meta's "Người gửi: <Staff Name>" metadata rows
   * so we skip both the label and the duplicated staff name inside them.
   */
  function findSenderLabelContainers(chatPane) {
    const blocked = new Set();
    const walker = document.createTreeWalker(chatPane, NodeFilter.SHOW_TEXT, null);
    let node;
    while ((node = walker.nextNode())) {
      const txt = normalizeText(node.nodeValue);
      if (txt === "người gửi:" || txt.startsWith("người gửi:") || txt === "sent by:" || txt.startsWith("sent by:")) {
        // Walk up 2-3 levels to cover the metadata row containing "Người gửi: <Name>"
        let el = node.parentElement;
        for (let i = 0; i < 3 && el && el !== chatPane; i++) {
          const full = normalizeText(el.innerText || el.textContent || "");
          if (full.length < 120 && el.getBoundingClientRect().height < 48) {
            blocked.add(el);
            el = el.parentElement;
          } else {
            break;
          }
        }
      }
    }
    return Array.from(blocked);
  }

  function isInsideAny(el, containers) {
    for (const c of containers) {
      if (c === el || c.contains(el)) return true;
    }
    return false;
  }

  // ---------------------------------------------------------------------------
  // Meta Business Suite structure (verified 2026-11):
  //  - Messenger/IG DM: messages live in [role=region][aria-label^="Vùng chứa danh sách tin nhắn"].
  //    Every real text bubble is a [dir=auto] element. Everything else inside the region
  //    (timestamps, "X đã trả lời một quảng cáo.", "Người gửi: <staff>", "Xem bài viết",
  //    "Hệ thống đã tự động tạo...") is NOT inside [dir=auto] and is system noise.
  // ---------------------------------------------------------------------------
  const LIKE_LABEL_RE = /^(👍|thích|like|biểu tượng thích|nút thích|thumbs up|\(y\))$/i;
  const EMOJI_IMG_RE = /static\.[a-z.]*fbcdn\.net\/images\/emoji|\/emoji\.php\//i;

  function findMessageRegion() {
    const regions = document.querySelectorAll('[role="region"][aria-label]');
    for (const r of regions) {
      const al = normalizeText(r.getAttribute("aria-label"));
      if (al.includes("danh sách tin nhắn") || al.includes("message list") || al.includes("messages list")) {
        const rect = r.getBoundingClientRect();
        if (rect.width > 200 && rect.height > 100) return r;
      }
    }
    return null;
  }

  function pickLikeItem(el, rect, centerX) {
    const label = normalizeText(
      el.getAttribute("aria-label") || el.getAttribute("alt") || el.querySelector?.("title")?.textContent || ""
    );
    if (!LIKE_LABEL_RE.test(label)) return null;
    if (rect.width < 24 || rect.height < 24 || rect.width > 140) return null;
    if (el.closest("button, [role='button']")) return null;
    return {
      sender: rect.left + rect.width / 2 > centerX ? "shop" : "customer",
      type: "text",
      text: "👍",
      imageUrl: "",
      el,
      top: rect.top
    };
  }

  function isContentImage(img, rect) {
    if (rect.width < 55 || rect.height < 55) return false;
    const src = img.currentSrc || img.src || "";
    if (!src || src.startsWith("data:image/svg")) return false;
    if (EMOJI_IMG_RE.test(src)) return false;
    if (img.naturalWidth > 0 && img.naturalWidth < 65) return false;
    return true;
  }

  function outermostDirAuto(root) {
    return Array.from(root.querySelectorAll('[dir="auto"]')).filter((el) => {
      const parentDir = el.parentElement && el.parentElement.closest('[dir="auto"]');
      return !parentDir || !root.contains(parentDir);
    });
  }

  function extractRegionMessages(region, composer) {
    const rRect = region.getBoundingClientRect();
    const centerX = rRect.left + rRect.width / 2;
    const items = [];
    const seenImg = new Set();

    for (const el of outermostDirAuto(region)) {
      if (composer && composer.contains(el)) continue;
      if (el.closest("button, [role='button']")) continue;
      const r = el.getBoundingClientRect();
      if (r.width === 0 || r.height === 0) continue;
      const text = (el.innerText || el.textContent || "").trim();
      if (!text) continue;
      items.push({
        sender: r.left + r.width / 2 > centerX ? "shop" : "customer",
        type: "text",
        text,
        imageUrl: "",
        el,
        top: r.top
      });
    }

    for (const img of region.querySelectorAll("img")) {
      const r = img.getBoundingClientRect();
      const like = pickLikeItem(img, r, centerX);
      if (like) {
        items.push(like);
        continue;
      }
      if (!isContentImage(img, r)) continue;
      const src = img.currentSrc || img.src;
      if (seenImg.has(src)) continue;
      seenImg.add(src);
      items.push({
        sender: r.left + r.width / 2 > centerX ? "shop" : "customer",
        type: "image",
        text: "",
        imageUrl: src,
        el: img,
        top: r.top
      });
    }

    // Big Like (👍) sent as a message: a [data-message-id] row with no text and no photo,
    // rendered as a CSS sprite <i data-visualcompletion="css-img"> (~36x36).
    for (const msg of region.querySelectorAll("[data-message-id]")) {
      if (msg.querySelector('[dir="auto"]')) continue;
      if (Array.from(msg.querySelectorAll("img")).some((img) => isContentImage(img, img.getBoundingClientRect()))) continue;
      const sprite = Array.from(msg.querySelectorAll('i[data-visualcompletion="css-img"], [role="img"]')).find((s) => {
        const r = s.getBoundingClientRect();
        return r.width >= 24 && r.height >= 24 && r.width <= 140 && !s.closest("button, [role='button']");
      });
      if (!sprite) continue;
      const r = sprite.getBoundingClientRect();
      items.push({
        sender: r.left + r.width / 2 > centerX ? "shop" : "customer",
        type: "text",
        text: "👍",
        imageUrl: "",
        el: sprite,
        top: r.top
      });
    }

    const seenLike = new Set();
    for (const el of region.querySelectorAll("svg[aria-label], [role='img'][aria-label]")) {
      if (el.closest("[data-message-id]") && !el.closest("[data-message-id]").querySelector('[dir="auto"]')) continue;
      if (el.tagName === "IMG") continue;
      const r = el.getBoundingClientRect();
      const like = pickLikeItem(el, r, centerX);
      if (!like) continue;
      const key = `${Math.round(r.top)}:${Math.round(r.left)}`;
      if (seenLike.has(key)) continue;
      seenLike.add(key);
      items.push(like);
    }

    items.sort((a, b) => a.top - b.top);
    return items;
  }

  function extractRawMessages(chatPane, composer) {
    // Only chat (Messenger/IG DM) is synced; comment threads produce junk names like "Mới nhất".
    if (/thread_type=(FB_PAGE_POST|IG_COMMENT|IG_MEDIA)/i.test(location.search)) return [];
    const region = findMessageRegion();
    if (region) {
      return extractRegionMessages(region, composer);
    }
    return extractRawMessagesLegacy(chatPane, composer);
  }

  function extractRawMessagesLegacy(chatPane, composer) {
    const paneRect = chatPane.getBoundingClientRect();
    const centerX = paneRect.left + paneRect.width / 2;
    const topBound = paneRect.top + 48;
    const bottomBound = findBottomNoiseBoundary(chatPane, composer);
    const senderLabelContainers = findSenderLabelContainers(chatPane);

    const items = [];
    const seenTexts = new Set();
    const seenImgSrcs = new Set();

    // 1. Extract product images
    const imgs = chatPane.querySelectorAll("img");
    for (const img of imgs) {
      const r = img.getBoundingClientRect();
      if (r.width < 55 || r.height < 55) continue;
      if (r.top < topBound || r.bottom > bottomBound + 10) continue;
      if (isInsideAny(img, senderLabelContainers)) continue;

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
        el: img,
        top: r.top
      });
    }

    // 2. Extract Like / Reaction sticker SVGs (e.g. customer sends 👍 Like icon)
    const svgs = chatPane.querySelectorAll("svg");
    const seenSvgPos = new Set();
    for (const svg of svgs) {
      const r = svg.getBoundingClientRect();
      if (r.width < 24 || r.width > 70 || r.height < 24 || r.height > 70) continue;
      if (r.top < topBound || r.bottom > bottomBound) continue;
      if (svg.closest("button, [role='button']")) continue;

      const label = normalizeText(svg.getAttribute("aria-label") || svg.querySelector("title")?.textContent || "");
      const isLike = label.includes("thích") || label.includes("like") || label.includes("thumbs");
      if (!isLike) continue;

      const key = `${Math.round(r.top)}:${Math.round(r.left)}`;
      if (seenSvgPos.has(key)) continue;
      seenSvgPos.add(key);

      const sender = r.left + r.width / 2 > centerX ? "shop" : "customer";
      items.push({
        sender,
        type: "text",
        text: "👍",
        imageUrl: "",
        el: svg,
        top: r.top
      });
    }

    // 3. Extract text messages
    const walker = document.createTreeWalker(chatPane, NodeFilter.SHOW_TEXT, null);
    let textNode;
    while ((textNode = walker.nextNode())) {
      const raw = (textNode.nodeValue || "").trim();
      if (!raw) continue;
      const parent = textNode.parentElement;
      if (!parent) continue;
      if (composer && composer.contains(parent)) continue;
      if (parent.closest("button, [role='button'], script, style, svg")) continue;
      if (isInsideAny(parent, senderLabelContainers)) continue;

      const r = parent.getBoundingClientRect();
      if (r.width === 0 || r.height === 0) continue;
      if (r.top < topBound || r.bottom > bottomBound) continue;

      // Skip small metadata text (font-size <= 11.5px like timestamps / "Đã gửi" / "Người gửi")
      const fontSize = parseFloat(getComputedStyle(parent).fontSize || "14");
      if (fontSize <= 11.5 && raw.length < 40) continue;

      if (
        /^(đã xem|seen|sent|đã gửi|đã nhận|delivered|nhập tin nhắn|type a message|trả lời|reply|bạn đã gửi|chưa đọc|unread|người gửi:?|\d{1,2}:\d{2})$/i.test(
          raw
        )
      ) {
        continue;
      }

      const distFromCenter = Math.abs(r.left + r.width / 2 - centerX);
      if (distFromCenter < 45 && r.width < 170 && raw.length < 30) {
        continue;
      }

      const bubbleEl = parent.closest('[dir="auto"]') || parent;
      if (isInsideAny(bubbleEl, senderLabelContainers)) continue;
      const bRect = bubbleEl.getBoundingClientRect();
      if (bRect.bottom > bottomBound) continue;

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
        el: bubbleEl,
        top: bRect.top
      });
    }

    items.sort((a, b) => a.top - b.top);
    return items;
  }

  function verifyVisualOverlay() {
    const { chatPane, composer } = findChatPaneAndComposer();
    const items = extractRawMessages(chatPane, composer);
    if (items.length === 0) {
      updateBadgeStatus("Verify: 0 tin nhắn", "#fbbf24");
      return;
    }

    const restored = [];
    for (const item of items) {
      if (!item.el) continue;
      const prevOutline = item.el.style.outline;
      const prevOffset = item.el.style.outlineOffset;
      const color = item.type === "image" ? "#eab308" : item.sender === "shop" ? "#3b82f6" : "#22c55e";
      item.el.style.outline = `3px solid ${color}`;
      item.el.style.outlineOffset = "2px";
      restored.push(() => {
        item.el.style.outline = prevOutline;
        item.el.style.outlineOffset = prevOffset;
      });
    }

    updateBadgeStatus(`Verify: ${items.length} mục (Xanh dương=Shop, Lá=Khách, Vàng=Ảnh)`, "#38bdf8");
    setTimeout(() => {
      restored.forEach((fn) => fn());
    }, 4500);
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
