<template>
  <div class="p-4 md:p-6 max-w-[1600px] mx-auto">
    <!-- Top Bar -->
    <div class="flex flex-wrap items-center justify-between gap-3 mb-4">
      <div>
        <h1 class="text-xl md:text-2xl font-extrabold text-slate-900 m-0 flex items-center gap-2.5">
          <span>AI Sales Copilot</span>
          <span
            class="text-[11px] font-bold px-2.5 py-0.5 rounded-full inline-flex items-center gap-1.5"
            :class="sseConnected ? 'bg-emerald-100 text-emerald-700' : 'bg-amber-100 text-amber-700'"
          >
            <span class="w-2 h-2 rounded-full" :class="sseConnected ? 'bg-emerald-500' : 'bg-amber-500'"></span>
            {{ sseConnected ? 'Real-time đang bật' : 'Đang kết nối...' }}
          </span>
        </h1>
        <p class="text-xs text-slate-500 mt-1 mb-0">
          Đồng bộ trực tiếp qua Facebook Webhook & Graph API • Nhận diện mã sản phẩm chuẩn Smart Upload • Tư vấn size & lên đơn
        </p>
      </div>
      <div class="flex items-center gap-2">
        <el-button type="primary" size="small" icon="el-icon-download" :loading="syncingFB" @click="syncFromFacebook">
          Đồng bộ hộp thư FB
        </el-button>
        <el-button size="small" icon="el-icon-refresh" @click="fetchSessions">Làm mới</el-button>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="sessions.length === 0" class="bg-white rounded-2xl border border-slate-200 p-12 text-center shadow-sm">
      <i class="el-icon-chat-dot-round text-5xl text-slate-300 mb-3 block"></i>
      <h3 class="text-lg font-bold text-slate-700 mb-1">Chưa có cuộc hội thoại nào được đồng bộ</h3>
      <p class="text-sm text-slate-500 max-w-md mx-auto mb-4">
        Bấm nút <strong>Đồng bộ hộp thư FB</strong> để kéo hội thoại gần nhất từ Fanpage, hoặc nhắn tin thử vào Page khi đã gắn Webhook.
      </p>
      <el-button type="primary" size="small" icon="el-icon-download" :loading="syncingFB" @click="syncFromFacebook">
        Đồng bộ hộp thư FB ngay
      </el-button>
    </div>

    <!-- 3-Column Copilot Workspace -->
    <div v-else class="grid grid-cols-1 lg:grid-cols-12 gap-4 items-stretch">
      <!-- Column 1: Session List -->
      <div class="lg:col-span-3 bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden flex flex-col h-[78vh]">
        <div class="px-4 py-3 border-b border-slate-100 font-bold text-[11px] uppercase tracking-wider text-slate-500 bg-slate-50/70">
          Hội thoại gần đây ({{ sessions.length }})
        </div>
        <div class="divide-y divide-slate-100 overflow-y-auto flex-1">
          <button
            v-for="s in sessions"
            :key="s.thread_id"
            @click="selectSession(s.thread_id)"
            class="w-full text-left px-4 py-3 transition-colors hover:bg-slate-50 focus:outline-none flex items-start gap-3"
            :class="activeThreadId === s.thread_id ? 'bg-blue-50/80 border-l-4 border-blue-600' : 'border-l-4 border-transparent'"
          >
            <img
              v-if="s.avatar_url"
              :src="s.avatar_url"
              class="w-9 h-9 rounded-full object-cover shrink-0 border border-slate-200"
            />
            <div v-else class="w-9 h-9 rounded-full bg-slate-200 text-slate-700 font-bold text-sm flex items-center justify-center shrink-0 uppercase">
              {{ (s.customer_name || 'K').charAt(0) }}
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between gap-2 mb-1">
                <span class="font-bold text-sm text-slate-800 truncate">{{ s.customer_name }}</span>
                <span class="text-[10px] text-slate-400 shrink-0">{{ formatTime(s.updated_at) }}</span>
              </div>
              <div class="flex flex-wrap gap-1 mb-1">
                <span
                  v-if="s.created_order"
                  class="text-[10px] font-bold px-1.5 py-0.5 rounded bg-emerald-100 text-emerald-700"
                >
                  ĐƠN #{{ s.created_order.id }}
                </span>
                <span
                  v-if="s.needs_human"
                  class="text-[10px] font-bold px-1.5 py-0.5 rounded bg-rose-100 text-rose-700"
                >
                  CẦN NGƯỜI
                </span>
                <span
                  v-for="p in (s.detected_products || []).slice(0, 3)"
                  :key="p.id"
                  class="text-[10px] font-bold px-1.5 py-0.5 rounded"
                  :class="p.status === 'AVAILABLE' ? 'bg-blue-100 text-blue-700' : 'bg-slate-100 text-slate-500'"
                >
                  {{ p.code }}
                </span>
              </div>
              <p class="text-xs text-slate-500 truncate m-0">{{ lastMessagePreview(s) }}</p>
            </div>
          </button>
        </div>
      </div>

      <!-- Column 2: Full Conversation (Shop + Customer) -->
      <div v-if="activeSession" class="lg:col-span-5 bg-white rounded-2xl border border-slate-200 shadow-sm flex flex-col h-[78vh] overflow-hidden">
        <div class="px-4 py-3 border-b border-slate-100 flex items-center justify-between bg-slate-50/70">
          <div class="flex items-center gap-2.5 min-w-0">
            <div class="w-8 h-8 rounded-full bg-blue-600 text-white font-bold text-xs flex items-center justify-center shrink-0 uppercase">
              {{ (activeSession.customer_name || 'K').charAt(0) }}
            </div>
            <div class="min-w-0">
              <h2 class="text-sm font-bold text-slate-900 m-0 truncate">{{ activeSession.customer_name }}</h2>
              <span class="text-[11px] text-slate-400 block">Hội thoại 2 chiều (sau mốc chốt đơn gần nhất)</span>
            </div>
          </div>
          <span class="text-xs font-semibold px-2 py-1 rounded-md bg-slate-100 text-slate-600 shrink-0">
            {{ (activeSession.messages || []).length }} tin
          </span>
        </div>

        <div ref="chatContainer" class="flex-1 overflow-y-auto p-4 space-y-2.5 bg-[#f4f6f9]">
          <div
            v-for="(msg, idx) in activeSession.messages"
            :key="idx"
            class="flex"
            :class="msg.sender === 'shop' ? 'justify-end' : 'justify-start'"
          >
            <div
              class="max-w-[80%] rounded-2xl px-3.5 py-2 text-[13px] shadow-sm"
              :class="msg.sender === 'shop'
                ? 'bg-blue-600 text-white rounded-br-sm'
                : 'bg-white text-slate-800 border border-slate-200/80 rounded-bl-sm'"
            >
              <div v-if="msg.type === 'image'">
                <img
                  v-if="(msg.image_base64 || msg.image_url) && !brokenImages[msg.image_base64 || msg.image_url]"
                  :src="msg.image_base64 || msg.image_url"
                  referrerpolicy="no-referrer"
                  loading="lazy"
                  @click="previewImage(msg.image_base64 || msg.image_url)"
                  @error="markImageBroken(msg.image_base64 || msg.image_url)"
                  class="max-w-[230px] max-h-[260px] rounded-xl object-cover cursor-pointer hover:opacity-95 transition-opacity block"
                />
                <div v-else class="text-xs italic opacity-75 py-1">
                  [Ảnh]
                </div>
                <div v-if="msg.ocr_text" class="mt-1.5 text-[11px] opacity-85 bg-black/10 rounded-lg px-2 py-1 break-words">
                  OCR: {{ msg.ocr_text }}
                </div>
              </div>
              <div v-else class="whitespace-pre-wrap break-words leading-relaxed">{{ msg.text }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Column 3: AI Suggestion, Detected Products & Receipt -->
      <div v-if="activeSession" class="lg:col-span-4 space-y-3 h-[78vh] overflow-y-auto pr-0.5">
        <!-- Human Alert -->
        <div v-if="activeSession.needs_human" class="bg-rose-50 border border-rose-200 rounded-2xl p-3.5 shadow-sm">
          <div class="font-bold text-rose-700 text-xs flex items-center gap-1.5 mb-1">
            <i class="el-icon-warning"></i>
            <span>Cần người vào xử lý trực tiếp</span>
          </div>
          <p class="text-xs text-rose-600 m-0">{{ activeSession.human_reason || 'Khách hỏi trường hợp đặc biệt (giảm giá / đổi trả).' }}</p>
        </div>

        <!-- Suggested Reply -->
        <div class="bg-white rounded-2xl border border-slate-200 p-4 shadow-sm">
          <div class="flex items-center justify-between mb-2">
            <h3 class="text-sm font-bold text-slate-800 m-0">Tin nhắn Suggestion</h3>
            <el-button type="primary" size="mini" icon="el-icon-document-copy" @click="copySuggestedReply">
              Copy Text
            </el-button>
          </div>
          <el-input
            v-model="editableReply"
            type="textarea"
            :rows="4"
            placeholder="Nội dung gợi ý..."
          />
          <div v-if="activeSession.size_advice" class="mt-2.5 bg-blue-50 border border-blue-100 rounded-xl p-2.5 text-xs text-blue-800">
            <strong>Tư vấn size:</strong> {{ activeSession.size_advice }}
          </div>
        </div>

        <!-- Detected Products -->
        <div class="bg-white rounded-2xl border border-slate-200 p-4 shadow-sm">
          <h3 class="text-sm font-bold text-slate-800 mb-2.5 mt-0">
            Sản phẩm nhận diện ({{ (activeSession.detected_products || []).length }})
          </h3>
          <div v-if="!activeSession.detected_products || activeSession.detected_products.length === 0" class="text-xs text-slate-400">
            Chưa phát hiện mã sản phẩm trong ảnh/tin nhắn.
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="p in activeSession.detected_products"
              :key="p.id"
              class="flex items-start gap-3 p-2.5 rounded-xl border border-slate-100 bg-slate-50/80"
            >
              <img
                v-if="p.thumb_url"
                :src="p.thumb_url"
                class="w-14 h-14 rounded-lg object-cover border border-slate-200 shrink-0"
              />
              <div class="flex-1 min-w-0">
                <div class="flex items-center justify-between gap-2">
                  <router-link :to="`/salev2/product/${p.id}/edit`" class="font-bold text-sm text-blue-600 hover:underline">
                    {{ p.code }}
                  </router-link>
                  <span
                    class="text-[10px] font-bold px-2 py-0.5 rounded-full"
                    :class="p.status === 'AVAILABLE' ? 'bg-emerald-100 text-emerald-700' : 'bg-slate-200 text-slate-600'"
                  >
                    {{ p.status === 'AVAILABLE' ? 'CÒN HÀNG' : 'ĐÃ BÁN' }}
                  </span>
                </div>
                <div class="text-xs font-bold text-slate-700 mt-0.5">{{ formatPrice(p.price) }}</div>
                <div v-if="p.sizes && p.sizes.length" class="text-[11px] text-slate-500 mt-1 flex flex-wrap gap-x-2 gap-y-0.5">
                  <span v-for="sz in p.sizes" :key="sz.name">{{ sz.name }}: <b>{{ sz.value }}</b></span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Created Order & Receipt -->
        <div v-if="activeSession.created_order" class="bg-emerald-50/70 rounded-2xl border border-emerald-200 p-4 shadow-sm">
          <div class="flex items-center justify-between mb-2">
            <div>
              <h3 class="text-sm font-bold text-emerald-800 m-0">
                Đã lên đơn #{{ activeSession.created_order.id }}
              </h3>
              <span class="text-xs text-emerald-700 font-semibold">
                Tổng: {{ formatPrice(activeSession.created_order.final_amount || activeSession.created_order.total_amount) }}
              </span>
            </div>
            <router-link
              :to="`/salev2/order/${activeSession.created_order.id}/edit`"
              class="px-3 py-1.5 rounded-lg bg-white border border-emerald-300 text-emerald-700 text-xs font-bold hover:bg-emerald-100"
            >
              Sửa đơn #{{ activeSession.created_order.id }}
            </router-link>
          </div>

          <div class="flex gap-2 mt-3">
            <el-button type="success" size="small" icon="el-icon-picture-outline" :loading="renderingReceipt" @click="copyReceiptImage">
              Copy Ảnh Receipt
            </el-button>
            <el-button
              v-if="receiptImg && /^\d+$/.test(activeSession.thread_id || '')"
              type="primary"
              size="small"
              icon="el-icon-s-promotion"
              :loading="sendingReceipt"
              @click="sendReceiptToCustomer"
            >
              Gửi cho khách
            </el-button>
            <el-button v-if="receiptImg" size="small" icon="el-icon-download" @click="downloadReceipt">
              Tải ảnh
            </el-button>
          </div>

          <div v-if="receiptImg" class="mt-3 rounded-xl overflow-hidden border border-emerald-200 bg-white">
            <img :src="receiptImg" class="w-full h-auto block" />
          </div>
        </div>

        <!-- Manual Create Order Button (if not created yet) -->
        <div v-else-if="hasAvailableDraftProducts" class="bg-white rounded-2xl border border-slate-200 p-4 shadow-sm">
          <div class="flex items-center justify-between mb-2">
            <h3 class="text-sm font-bold text-slate-800 m-0">Lên đơn nhanh</h3>
            <span class="text-[11px] text-slate-400">Hoặc nhắn: "OK ạ, để em lên đơn ngay"</span>
          </div>
          <div class="text-xs text-slate-600 space-y-1 mb-3">
            <div>Khách: <b>{{ activeSession.order_draft?.name || activeSession.customer_name }}</b></div>
            <div>SĐT: <b>{{ activeSession.order_draft?.phone || '(chưa có)' }}</b></div>
            <div>Địa chỉ: <b>{{ activeSession.order_draft?.address || '(chưa có)' }}</b></div>
            <div>Mã: <b>{{ (activeSession.order_draft?.product_codes || []).join(', ') }}</b></div>
          </div>
          <el-button type="primary" size="small" class="w-full" :loading="creatingOrder" @click="manualCreateOrder">
            Tạo đơn & Xuất Receipt ngay
          </el-button>
        </div>
      </div>
    </div>

    <!-- Image Preview Modal -->
    <el-dialog :visible.sync="previewVisible" width="520px" append-to-body>
      <img v-if="previewSrc" :src="previewSrc" class="w-full h-auto rounded-lg block" />
    </el-dialog>

    <!-- Hidden Receipt Template -->
    <div style="position: absolute; left: -9999px; top: 0;">
      <div ref="receiptTemplate" id="copilot-order-receipt" class="p-8 bg-white w-[450px] font-sans text-slate-800">
        <div class="text-center mb-6">
          <h1 class="text-4xl font-extrabold tracking-[4px] m-0 text-slate-900">VINTAGE</h1>
          <p class="text-sm font-bold tracking-[2px] text-blue-600 mt-1 uppercase">Always True</p>
        </div>
        <div class="bg-slate-50 rounded-2xl p-6 mb-5">
          <div class="flex items-center mb-3">
            <div class="bg-blue-600 text-white w-10 h-10 rounded-full flex items-center justify-center font-bold mr-3">P</div>
            <div>
              <span class="text-[10px] text-slate-400 uppercase font-bold tracking-wider">KHÁCH HÀNG</span>
              <h2 class="text-lg font-bold m-0">{{ currentOrder?.name }}</h2>
            </div>
          </div>
          <div class="text-sm text-slate-600 space-y-1">
            <p><i class="el-icon-phone mr-2"></i>{{ currentOrder?.phone }}</p>
            <p><i class="el-icon-location mr-2"></i>{{ currentOrder?.address }}</p>
          </div>
        </div>
        <div v-if="currentOrder?.note" class="bg-yellow-50 border border-yellow-100 rounded-xl p-4 mb-6">
          <p class="text-sm font-medium">{{ currentOrder.note }}</p>
        </div>
        <h3 class="text-[10px] font-black text-slate-300 uppercase tracking-widest border-b border-slate-100 pb-2 mb-4">ĐƠN HÀNG CỦA BẠN</h3>
        <div class="flex flex-wrap gap-2.5 mb-6">
          <div v-for="p in currentOrder?.products" :key="p.id" class="w-[31%] text-center mb-3">
            <div class="h-[140px] bg-slate-50 rounded-xl mb-1.5 border border-slate-100 flex items-center justify-center overflow-hidden">
              <img :src="resolveProductThumb(p)" class="max-w-full max-h-full object-contain" crossorigin="anonymous" />
            </div>
            <div class="font-bold text-slate-800 text-xs leading-normal pt-1 pb-1 px-1 break-words">{{ p.name }}</div>
          </div>
        </div>
        <div class="bg-slate-50 rounded-2xl p-6 mb-8 border border-slate-100 font-bold">
          <div class="flex justify-between mb-3 text-sm text-slate-500">
            <span>Tổng đơn</span>
            <span>{{ formatPrice(currentOrder?.total_amount) }}</span>
          </div>
          <div class="pt-4 border-t border-slate-200 flex justify-between items-center">
            <div>
              <span class="text-lg block leading-tight">Viettel Amount</span>
              <span class="text-[11px] font-bold text-slate-400">+ phí ship</span>
            </div>
            <span class="text-2xl text-blue-600">{{ formatPrice(currentOrder?.final_amount || currentOrder?.total_amount) }}</span>
          </div>
        </div>
        <div class="text-center border-t border-dashed border-slate-200 pt-4 text-slate-500">
          <h3 class="font-bold text-sm mb-1 text-slate-900">Thank you for your order!</h3>
          <p class="text-[10px] text-slate-400">Vintage Always True • Vintage Menswear</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'CopilotView',
  data() {
    return {
      sessions: [],
      activeThreadId: null,
      editableReply: '',
      lastServerReply: '',
      sseConnected: false,
      eventSource: null,
      receiptImg: '',
      receiptCanvas: null,
      renderingReceipt: false,
      creatingOrder: false,
      sendingReceipt: false,
      syncingFB: false,
      previewVisible: false,
      previewSrc: '',
      brokenImages: {}
    };
  },
  computed: {
    activeSession() {
      if (!this.activeThreadId) return this.sessions[0] || null;
      return this.sessions.find((s) => s.thread_id === this.activeThreadId) || this.sessions[0] || null;
    },
    currentOrder() {
      return this.activeSession ? this.activeSession.created_order : null;
    },
    hasAvailableDraftProducts() {
      return Boolean(
        this.activeSession &&
        this.activeSession.order_draft &&
        this.activeSession.order_draft.product_codes &&
        this.activeSession.order_draft.product_codes.length > 0
      );
    }
  },
  mounted() {
    this.fetchSessions();
    this.connectSSE();
  },
  beforeDestroy() {
    if (this.eventSource) {
      this.eventSource.close();
    }
  },
  methods: {
    async fetchSessions() {
      try {
        const res = await axios.get('/salev2/api/copilot/sessions');
        this.sessions = res.data.sessions || [];
        if (!this.activeThreadId && this.sessions.length > 0) {
          this.selectSession(this.sessions[0].thread_id);
        } else if (this.activeSession) {
          this.syncActiveState(this.activeSession, true);
        }
      } catch (e) {
        console.error('Failed to fetch copilot sessions', e);
      }
    },
    async syncFromFacebook() {
      this.syncingFB = true;
      try {
        const res = await axios.post('/salev2/api/copilot/sync-fb');
        if (res.data && res.data.sessions) {
          this.sessions = res.data.sessions;
          if (!this.activeThreadId && this.sessions.length > 0) {
            this.selectSession(this.sessions[0].thread_id);
          }
          this.$message.success(`Đã đồng bộ ${res.data.synced || 0} cuộc hội thoại từ Facebook!`);
        }
      } catch (e) {
        const msg = (e.response && e.response.data && e.response.data.error) || 'Không thể đồng bộ từ Facebook (kiểm tra quyền pages_messaging / token)';
        this.$message.error(msg);
      } finally {
        this.syncingFB = false;
      }
    },
    connectSSE() {
      if (this.eventSource) this.eventSource.close();
      this.eventSource = new EventSource('/salev2/api/copilot/stream');
      this.eventSource.onopen = () => {
        this.sseConnected = true;
      };
      this.eventSource.onerror = () => {
        this.sseConnected = false;
      };
      this.eventSource.addEventListener('session', (event) => {
        try {
          const sess = JSON.parse(event.data);
          this.upsertSession(sess);
        } catch (e) {
          console.error('Invalid SSE session payload', e);
        }
      });
    },
    upsertSession(sess) {
      const idx = this.sessions.findIndex((s) => s.thread_id === sess.thread_id);
      const prevCount = idx !== -1 && this.sessions[idx].messages ? this.sessions[idx].messages.length : 0;
      const newCount = sess.messages ? sess.messages.length : 0;

      if (idx !== -1) {
        // Update in-place to avoid jumping/flickering the list
        this.$set(this.sessions, idx, sess);
      } else {
        this.sessions.unshift(sess);
        if (!this.activeThreadId) {
          this.activeThreadId = sess.thread_id;
        }
      }

      if (this.activeThreadId === sess.thread_id) {
        this.syncActiveState(sess, newCount > prevCount);
      }
    },
    selectSession(threadId) {
      this.activeThreadId = threadId;
      if (this.activeSession) {
        this.syncActiveState(this.activeSession, true);
      }
    },
    syncActiveState(sess, shouldScroll) {
      const nextReply = sess.suggested_reply || '';
      if (nextReply !== this.lastServerReply) {
        this.editableReply = nextReply;
        this.lastServerReply = nextReply;
      }
      if (shouldScroll) {
        this.$nextTick(() => this.scrollChatToBottom());
      }
      if (sess.created_order && (!this.currentOrder || this.currentOrder.id !== sess.created_order.id || !this.receiptImg)) {
        this.renderReceipt(sess.created_order);
      } else if (!sess.created_order) {
        this.receiptImg = '';
        this.receiptCanvas = null;
      }
    },
    scrollChatToBottom() {
      const el = this.$refs.chatContainer;
      if (el) {
        el.scrollTop = el.scrollHeight;
      }
    },
    markImageBroken(src) {
      this.$set(this.brokenImages, src, true);
    },
    previewImage(src) {
      this.previewSrc = src;
      this.previewVisible = true;
    },
    async copySuggestedReply() {
      try {
        await navigator.clipboard.writeText(this.editableReply || '');
        this.$message.success('Đã copy tin nhắn gợi ý!');
      } catch (e) {
        this.$message.error('Không thể copy tự động');
      }
    },
    async renderReceipt(order) {
      if (!order) return;
      this.renderingReceipt = true;
      await this.$nextTick();
      try {
        const el = this.$refs.receiptTemplate;
        if (el) {
          const imgs = Array.from(el.querySelectorAll('img'));
          await Promise.race([
            Promise.all(
              imgs.map((img) => {
                if (img.complete && img.naturalHeight !== 0) return Promise.resolve();
                return new Promise((resolve) => {
                  img.onload = resolve;
                  img.onerror = resolve;
                });
              })
            ),
            new Promise((resolve) => setTimeout(resolve, 2000))
          ]);
          const canvas = await html2canvas(el, { useCORS: true, scale: 2 });
          this.receiptCanvas = canvas;
          this.receiptImg = canvas.toDataURL('image/png');
        }
      } catch (e) {
        console.error('Render receipt error', e);
      } finally {
        this.renderingReceipt = false;
      }
    },
    async copyReceiptImage() {
      if (!this.receiptCanvas && this.currentOrder) {
        await this.renderReceipt(this.currentOrder);
      }
      if (!this.receiptCanvas) return;
      try {
        this.receiptCanvas.toBlob(async (blob) => {
          if (!blob) return;
          await navigator.clipboard.write([new ClipboardItem({ 'image/png': blob })]);
          this.$message.success('Đã copy ảnh Receipt! Qua Facebook Cmd+V để dán ngay.');
        }, 'image/png');
      } catch (e) {
        this.$message.error('Trình duyệt không hỗ trợ copy ảnh trực tiếp, vui lòng bấm Tải ảnh.');
      }
    },
    async sendReceiptToCustomer() {
      if (!this.receiptImg || !this.activeSession) return;
      try {
        await this.$confirm(`Gửi ảnh receipt đơn #${this.currentOrder ? this.currentOrder.id : ''} cho ${this.activeSession.customer_name} qua Messenger?`, 'Gửi receipt', {
          confirmButtonText: 'Gửi',
          cancelButtonText: 'Huỷ',
          type: 'info'
        });
      } catch (e) {
        return;
      }
      this.sendingReceipt = true;
      try {
        await axios.post('/salev2/api/copilot/send-receipt', {
          thread_id: this.activeSession.thread_id,
          image_base64: this.receiptImg
        });
        this.$message.success('Đã gửi receipt cho khách qua Messenger!');
      } catch (e) {
        const msg = (e.response && e.response.data && e.response.data.error) || e.message;
        this.$message.error('Gửi receipt thất bại: ' + msg);
      } finally {
        this.sendingReceipt = false;
      }
    },
    downloadReceipt() {
      if (!this.receiptImg || !this.currentOrder) return;
      const link = document.createElement('a');
      link.download = `Order_${this.currentOrder.id}_${this.currentOrder.name || 'Receipt'}.png`;
      link.href = this.receiptImg;
      link.click();
    },
    async manualCreateOrder() {
      if (!this.activeSession || !this.activeSession.order_draft) return;
      this.creatingOrder = true;
      try {
        const res = await axios.post('/salev2/api/copilot/order', {
          thread_id: this.activeSession.thread_id,
          ...this.activeSession.order_draft
        });
        if (res.data && res.data.order) {
          this.$set(this.activeSession, 'created_order', res.data.order);
          this.$message.success(`Đã tạo đơn #${res.data.order.id}`);
          this.renderReceipt(res.data.order);
        }
      } catch (e) {
        const msg = (e.response && e.response.data && e.response.data.error) || 'Không thể tạo đơn';
        this.$message.error(msg);
      } finally {
        this.creatingOrder = false;
      }
    },
    resolveProductThumb(p) {
      const raw = p.image_thumb_scale_url || p.image_thumb_scale || p.path_thumb || '';
      if (!raw) return '';
      if (raw.startsWith('http')) return raw;
      return '/' + raw.replace(/^\/+/, '');
    },
    lastMessagePreview(s) {
      const msgs = s.messages || [];
      if (msgs.length === 0) return '';
      const last = msgs[msgs.length - 1];
      return last.type === 'image' ? '[Hình ảnh]' : last.text;
    },
    formatPrice(val) {
      const n = Number(val || 0);
      return n.toLocaleString('vi-VN') + 'đ';
    },
    formatTime(ts) {
      if (!ts) return '';
      const d = new Date(ts);
      return d.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' });
    }
  }
};
</script>
