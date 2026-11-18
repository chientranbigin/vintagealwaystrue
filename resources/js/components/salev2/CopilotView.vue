<template>
  <div class="copilot-root max-w-[1600px] mx-auto md:p-6 md:h-[calc(100vh-64px)] md:flex md:flex-col md:overflow-hidden">
    <!-- ===================== MOBILE (< md): Messenger-style ===================== -->
    <div class="md:hidden">
      <!-- Screen 1: conversation list -->
      <div class="bg-white min-h-screen">
        <div class="sticky top-0 z-30 bg-white border-b border-slate-100 px-4 pt-3 pb-2">
          <div class="flex items-center gap-2">
            <h1 class="text-xl font-extrabold text-slate-900 m-0 flex-1">Copilot</h1>
            <span
              class="text-[11px] font-bold px-2 py-0.5 rounded-full inline-flex items-center gap-1"
              :class="sseConnected ? 'bg-emerald-50 text-emerald-700' : 'bg-amber-50 text-amber-700'"
            >
              <span class="w-1.5 h-1.5 rounded-full" :class="sseConnected ? 'bg-emerald-500' : 'bg-amber-500'"></span>
              {{ sseConnected ? 'Live' : 'Đang nối' }}
            </span>
            <button
              class="w-9 h-9 rounded-full bg-slate-100 text-slate-600 flex items-center justify-center active:bg-slate-200"
              aria-label="Đồng bộ hộp thư Facebook"
              @click="syncFromFacebook"
            >
              <i :class="syncingFB ? 'el-icon-loading' : 'el-icon-refresh'"></i>
            </button>
          </div>
          <div class="mt-2 flex items-center gap-2 bg-slate-100 rounded-full px-3 h-9">
            <i class="el-icon-search text-slate-400"></i>
            <input
              v-model="search"
              type="search"
              placeholder="Tìm khách hoặc mã SP"
              class="flex-1 bg-transparent outline-none text-[16px] text-slate-800"
            />
          </div>
        </div>

        <div v-if="sortedSessions.length === 0" class="p-10 text-center text-slate-500 text-sm">
          <i class="el-icon-chat-dot-round text-4xl text-slate-300 block mb-2"></i>
          {{ search ? 'Không tìm thấy hội thoại phù hợp.' : 'Chưa có hội thoại. Bấm nút đồng bộ ở góc trên.' }}
        </div>

        <button
          v-for="s in sortedSessions"
          :key="'m-' + s.thread_id"
          class="w-full text-left flex items-center gap-3 px-4 py-2.5 active:bg-slate-100 transition-colors duration-700"
          :class="flashIds[s.thread_id] ? 'bg-blue-50' : 'bg-white'"
          @click="openChat(s.thread_id)"
        >
          <div class="relative shrink-0">
            <img v-if="s.avatar_url" :src="s.avatar_url" class="w-12 h-12 rounded-full object-cover" />
            <div v-else class="w-12 h-12 rounded-full bg-slate-200 text-slate-600 font-bold flex items-center justify-center uppercase">
              {{ (s.customer_name || 'K').charAt(0) }}
            </div>
          </div>
          <div class="flex-1 min-w-0 border-b border-slate-100 pb-2.5">
            <div class="flex items-center gap-2">
              <span class="flex-1 truncate text-[15px]" :class="isUnread(s) ? 'font-bold text-slate-900' : 'font-semibold text-slate-800'">
                {{ s.customer_name }}
              </span>
              <span class="text-xs shrink-0" :class="isUnread(s) ? 'text-blue-600 font-bold' : 'text-slate-400'">
                {{ formatTime(sortTime(s)) }}
              </span>
            </div>
            <div class="flex items-center gap-2">
              <p class="flex-1 truncate text-[13px] m-0" :class="isUnread(s) ? 'text-slate-900 font-semibold' : 'text-slate-500'">
                {{ lastMessagePreview(s) }}
              </p>
              <span v-if="isUnread(s)" class="w-2.5 h-2.5 rounded-full bg-blue-600 shrink-0"></span>
            </div>
            <div v-if="s.created_order || s.needs_human || (s.detected_products || []).length" class="flex flex-wrap gap-1 mt-1">
              <span v-if="s.created_order" class="text-[10px] font-bold px-1.5 py-0.5 rounded bg-emerald-100 text-emerald-700">
                Đơn #{{ s.created_order.id }}
              </span>
              <span v-if="s.needs_human" class="text-[10px] font-bold px-1.5 py-0.5 rounded bg-amber-100 text-amber-800">
                Cần NV
              </span>
              <span
                v-for="p in (s.detected_products || []).slice(0, 3)"
                :key="p.id"
                class="text-[10px] font-bold px-1.5 py-0.5 rounded"
                :class="p.status === 'AVAILABLE' ? 'bg-blue-50 text-blue-700' : 'bg-slate-100 text-slate-400 line-through'"
              >
                {{ p.code }}
              </span>
            </div>
          </div>
        </button>
      </div>

      <!-- Screen 2: chat (full screen, covers the app bottom nav) -->
      <div v-if="mobileView === 'chat' && activeSession" class="fixed inset-0 z-[60] bg-white flex flex-col">
        <div class="flex items-center gap-2 px-2 h-14 border-b border-slate-100 shrink-0" style="padding-top: env(safe-area-inset-top)">
          <button class="w-10 h-10 flex items-center justify-center text-blue-600 text-xl" aria-label="Quay lại" @click="closeChat">
            <i class="el-icon-arrow-left"></i>
          </button>
          <img v-if="activeSession.avatar_url" :src="activeSession.avatar_url" class="w-9 h-9 rounded-full object-cover" />
          <div v-else class="w-9 h-9 rounded-full bg-slate-200 text-slate-600 font-bold flex items-center justify-center uppercase">
            {{ (activeSession.customer_name || 'K').charAt(0) }}
          </div>
          <div class="flex-1 min-w-0">
            <div class="font-bold text-[15px] text-slate-900 truncate">{{ activeSession.customer_name }}</div>
            <div class="text-[11px] text-slate-400">{{ formatTime(sortTime(activeSession)) }}</div>
          </div>
          <button
            class="h-9 px-3 rounded-full bg-slate-100 text-slate-700 text-sm font-semibold flex items-center gap-1 active:bg-slate-200"
            aria-label="Sản phẩm nhận diện"
            @click="sheet = 'products'"
          >
            <i class="el-icon-goods"></i>{{ (activeSession.detected_products || []).length }}
          </button>
          <button
            v-if="!activeSession.created_order"
            class="h-9 px-3 rounded-full bg-emerald-600 text-white text-sm font-bold flex items-center gap-1 active:bg-emerald-700 disabled:opacity-60"
            :disabled="goingToOrder"
            @click="goCreateOrder"
          >
            <i :class="goingToOrder ? 'el-icon-loading' : 'el-icon-shopping-cart-full'"></i> Lên đơn
          </button>
        </div>

        <button
          v-if="activeSession.created_order"
          class="shrink-0 flex items-center gap-2 px-4 py-2 bg-emerald-50 text-emerald-800 text-[13px] font-semibold border-b border-emerald-100 text-left"
          @click="openOrderSheet"
        >
          <i class="el-icon-success"></i>
          <span class="flex-1">Đã lên đơn #{{ activeSession.created_order.id }} · {{ formatPrice(activeSession.created_order.final_amount || activeSession.created_order.total_amount) }}</span>
          <span class="text-emerald-600">Receipt ›</span>
        </button>
        <div v-if="activeSession.needs_human" class="shrink-0 px-4 py-1.5 bg-amber-50 text-amber-800 text-[12px] border-b border-amber-100">
          <i class="el-icon-warning-outline"></i> Cần nhân viên: {{ activeSession.human_reason || 'khách hỏi thông tin chưa có dữ liệu' }}
        </div>

        <div ref="mobileChat" class="flex-1 overflow-y-auto px-3 py-3 space-y-1.5 bg-white" @scroll="onChatScroll">
          <div
            v-for="(msg, idx) in activeSession.messages"
            :key="'mm-' + idx"
            class="flex"
            :class="msg.sender === 'shop' ? 'justify-end' : 'justify-start'"
          >
            <div v-if="msg.type === 'image'" class="max-w-[70%]">
              <img
                v-if="(msg.image_base64 || msg.image_url) && !brokenImages[msg.image_base64 || msg.image_url]"
                :src="msg.image_base64 || msg.image_url"
                referrerpolicy="no-referrer"
                loading="lazy"
                class="max-h-[260px] rounded-2xl object-cover block border border-slate-100"
                @load="onChatImageLoad"
                @click="previewImage(msg.image_base64 || msg.image_url)"
                @error="markImageBroken(msg.image_base64 || msg.image_url)"
              />
              <div v-else class="px-3 py-2 rounded-2xl bg-slate-100 text-xs italic text-slate-500">[Ảnh]</div>
            </div>
            <div
              v-else
              class="max-w-[78%] px-3.5 py-2 text-[15px] leading-snug whitespace-pre-wrap break-words rounded-[18px]"
              :class="[
                msg.sender === 'shop' ? 'bg-blue-600 text-white' : 'bg-slate-100 text-slate-900',
                msg.pending ? 'opacity-60' : ''
              ]"
            >{{ msg.text }}</div>
          </div>
        </div>

        <!-- Composer: AI suggestion pre-filled, editable, clear + send -->
        <div class="shrink-0 border-t border-slate-200 bg-white px-3 pt-2" style="padding-bottom: calc(env(safe-area-inset-bottom) + 8px)">
          <div v-if="pendingSend" class="flex items-center gap-2 bg-slate-900 text-white rounded-xl px-3 py-2 mb-2 text-[13px]">
            <i class="el-icon-loading"></i>
            <span class="flex-1">Đang gửi sau {{ pendingSend.left }}s…</span>
            <button class="font-bold text-amber-300" @click="undoSend">Hoàn tác</button>
          </div>
          <div class="flex items-center gap-1.5 text-[11px] font-bold mb-1" :class="replyDirty ? 'text-slate-500' : 'text-blue-700'">
            <i :class="reanalyzing[activeSession.thread_id] ? 'el-icon-loading' : (replyDirty ? 'el-icon-edit' : 'el-icon-magic-stick')"></i>
            <span class="flex-1">{{ reanalyzing[activeSession.thread_id] ? 'AI đang phân tích lại…' : (replyDirty ? 'Tin nhắn của bạn' : 'Gợi ý AI — sửa trực tiếp hoặc gửi luôn') }}</span>
            <button v-if="replyDirty && lastServerReply && editableReply !== lastServerReply" class="text-blue-600 font-semibold" @click="useSuggestion">Dùng gợi ý AI</button>
          </div>
          <div class="flex items-end gap-2">
            <div class="relative flex-1">
              <textarea
                v-model="editableReply"
                rows="3"
                placeholder="Nhập tin nhắn…"
                class="block w-full rounded-2xl border border-slate-300 bg-slate-50 pl-3 pr-10 py-2 text-[16px] leading-snug outline-none focus:border-blue-500 focus:bg-white resize-none"
                @input="replyDirty = true"
              ></textarea>
              <button
                v-if="editableReply"
                class="absolute top-1.5 right-1.5 w-7 h-7 rounded-full bg-slate-200 text-slate-600 flex items-center justify-center active:bg-slate-300"
                aria-label="Xoá nội dung để tự nhập"
                @click="clearReply"
              ><i class="el-icon-close"></i></button>
            </div>
            <button
              class="w-11 h-11 rounded-full bg-blue-600 text-white flex items-center justify-center text-lg active:bg-blue-700 disabled:opacity-40 shrink-0"
              aria-label="Gửi cho khách"
              :disabled="!canSendToCustomer || !editableReply.trim() || !!pendingSend"
              @click="queueSend"
            ><i class="el-icon-s-promotion"></i></button>
          </div>
        </div>
      </div>

      <!-- Screen 3: bottom sheets (products / order + receipt) -->
      <div v-if="sheet && activeSession" class="fixed inset-0 z-[70] flex flex-col justify-end">
        <div class="absolute inset-0 bg-black/40" @click="sheet = null"></div>
        <div class="relative bg-white rounded-t-3xl max-h-[88vh] overflow-y-auto px-4 pt-2" style="padding-bottom: calc(env(safe-area-inset-bottom) + 16px)">
          <div class="w-10 h-1.5 rounded-full bg-slate-300 mx-auto mb-3"></div>
          <button
            class="absolute top-3 right-3 w-9 h-9 rounded-full bg-slate-100 text-slate-600 flex items-center justify-center text-lg active:bg-slate-200"
            aria-label="Đóng"
            @click="sheet = null"
          ><i class="el-icon-close"></i></button>

          <template v-if="sheet === 'products'">
            <h3 class="text-base font-bold text-slate-900 mt-0 mb-3 pr-12">Sản phẩm nhận diện ({{ (activeSession.detected_products || []).length }})</h3>
            <div v-if="!(activeSession.detected_products || []).length" class="text-sm text-slate-400 pb-4">
              Chưa thấy mã sản phẩm trong ảnh/tin nhắn.
            </div>
            <div v-for="p in activeSession.detected_products" :key="'sp-' + p.id" class="flex gap-3 py-2.5 border-b border-slate-100">
              <img v-if="p.thumb_url" :src="p.thumb_url" class="w-16 h-20 rounded-lg object-cover border border-slate-100 shrink-0" @click="previewImage(p.thumb_url)" />
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2">
                  <span class="font-bold text-[15px] text-slate-900 flex-1">{{ p.code }}</span>
                  <span class="text-[11px] font-bold px-2 py-0.5 rounded-full" :class="p.status === 'AVAILABLE' ? 'bg-emerald-100 text-emerald-700' : 'bg-slate-200 text-slate-500'">
                    {{ p.status === 'AVAILABLE' ? 'Còn hàng' : 'Đã bán' }}
                  </span>
                </div>
                <div class="text-sm font-bold text-slate-700">{{ formatPrice(p.price) }}</div>
                <div v-if="p.sizes && p.sizes.length" class="text-[12px] text-slate-500 flex flex-wrap gap-x-2">
                  <span v-for="sz in p.sizes" :key="sz.name">{{ sz.name }} <b>{{ sz.value }}</b></span>
                </div>
              </div>
            </div>
            <div v-if="activeSession.size_advice" class="mt-3 text-[13px] text-blue-800 bg-blue-50 rounded-xl px-3 py-2">
              <b>Tư vấn size:</b> {{ activeSession.size_advice }}
            </div>
            <button
              v-if="!activeSession.created_order"
              class="w-full h-11 mt-3 rounded-xl bg-emerald-600 text-white font-bold disabled:opacity-60"
              :disabled="goingToOrder"
              @click="goCreateOrder"
            >
              <i :class="goingToOrder ? 'el-icon-loading' : 'el-icon-shopping-cart-full'"></i> Lên đơn
            </button>
          </template>

          <template v-else-if="sheet === 'order' && activeSession.created_order">
            <!-- Order exists: receipt -->
              <div class="flex items-center gap-2 mb-3 pr-12">
                <h3 class="text-base font-bold text-emerald-800 m-0 flex-1">Đã lên đơn #{{ activeSession.created_order.id }}</h3>
                <router-link :to="`/salev2/order/${activeSession.created_order.id}/edit`" class="text-sm font-semibold text-blue-600">Sửa đơn</router-link>
              </div>
              <div class="rounded-2xl overflow-hidden border border-slate-200 bg-slate-50 mb-3 min-h-[120px] flex items-center justify-center">
                <img v-if="receiptImg" :src="receiptImg" class="w-full h-auto block" />
                <span v-else class="text-sm text-slate-400"><i class="el-icon-loading"></i> Đang tạo receipt…</span>
              </div>
              <div class="flex gap-2">
                <button
                  class="flex-[2] h-11 rounded-xl bg-blue-600 text-white font-bold disabled:opacity-50"
                  :disabled="!receiptImg || !canSendToCustomer || sendingReceipt"
                  @click="sendReceiptToCustomer"
                >
                  <i :class="sendingReceipt ? 'el-icon-loading' : 'el-icon-s-promotion'"></i> Gửi receipt
                </button>
                <button class="flex-1 h-11 rounded-xl border border-slate-300 font-semibold text-slate-700" @click="copyReceiptImage">Copy ảnh</button>
                <button class="w-11 h-11 rounded-xl border border-slate-300 text-slate-700" aria-label="Tải ảnh receipt" @click="downloadReceipt">
                  <i class="el-icon-download"></i>
                </button>
              </div>
          </template>
        </div>
      </div>
    </div>

    <!-- ===================== DESKTOP (md+): 3 columns ===================== -->
    <div class="hidden md:flex md:flex-col md:flex-1 md:min-h-0">
      <div class="flex flex-wrap items-center justify-between gap-3 mb-4 shrink-0">
        <div>
          <h1 class="text-2xl font-extrabold text-slate-900 m-0 flex items-center gap-2.5">
            <span>AI Sales Copilot</span>
            <span
              class="text-[11px] font-bold px-2.5 py-0.5 rounded-full inline-flex items-center gap-1.5"
              :class="sseConnected ? 'bg-emerald-100 text-emerald-700' : 'bg-amber-100 text-amber-700'"
            >
              <span class="w-2 h-2 rounded-full" :class="sseConnected ? 'bg-emerald-500' : 'bg-amber-500'"></span>
              {{ sseConnected ? 'Real-time đang bật' : 'Đang kết nối...' }}
            </span>
          </h1>
        </div>
        <div class="flex items-center gap-2">
          <el-input v-model="search" size="small" prefix-icon="el-icon-search" placeholder="Tìm khách / mã SP" clearable style="width: 220px" />
          <el-button type="primary" size="small" icon="el-icon-download" :loading="syncingFB" @click="syncFromFacebook">Đồng bộ hộp thư FB</el-button>
          <el-button size="small" icon="el-icon-refresh" @click="fetchSessions">Làm mới</el-button>
        </div>
      </div>

      <div v-if="sessions.length === 0" class="bg-white rounded-2xl border border-slate-200 p-12 text-center shadow-sm">
        <i class="el-icon-chat-dot-round text-5xl text-slate-300 mb-3 block"></i>
        <h3 class="text-lg font-bold text-slate-700 mb-1">Chưa có cuộc hội thoại nào được đồng bộ</h3>
        <el-button type="primary" size="small" icon="el-icon-download" :loading="syncingFB" @click="syncFromFacebook">Đồng bộ hộp thư FB ngay</el-button>
      </div>

      <div v-else class="grid grid-cols-12 gap-4 flex-1 min-h-0">
        <!-- Column 1: list -->
        <div class="col-span-3 bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden flex flex-col h-full min-h-0">
          <div class="px-4 py-3 border-b border-slate-100 font-bold text-[11px] uppercase tracking-wider text-slate-500 bg-slate-50/70">
            Hội thoại gần đây ({{ sortedSessions.length }})
          </div>
          <div class="overflow-y-auto flex-1">
            <button
              v-for="s in sortedSessions"
              :key="'d-' + s.thread_id"
              class="w-full text-left px-4 py-3 flex items-start gap-3 border-b border-slate-100 border-l-4 transition-colors duration-700"
              :class="[
                activeThreadId === s.thread_id ? 'bg-blue-50/80 border-l-blue-600' : 'border-l-transparent hover:bg-slate-50',
                flashIds[s.thread_id] ? 'bg-blue-50' : ''
              ]"
              @click="selectSession(s.thread_id)"
            >
              <img v-if="s.avatar_url" :src="s.avatar_url" class="w-9 h-9 rounded-full object-cover shrink-0 border border-slate-200" />
              <div v-else class="w-9 h-9 rounded-full bg-slate-200 text-slate-700 font-bold text-sm flex items-center justify-center shrink-0 uppercase">
                {{ (s.customer_name || 'K').charAt(0) }}
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center justify-between gap-2 mb-0.5">
                  <span class="text-sm truncate" :class="isUnread(s) ? 'font-extrabold text-slate-900' : 'font-bold text-slate-800'">{{ s.customer_name }}</span>
                  <span class="text-[10px] shrink-0" :class="isUnread(s) ? 'text-blue-600 font-bold' : 'text-slate-400'">{{ formatTime(sortTime(s)) }}</span>
                </div>
                <div class="flex flex-wrap gap-1 mb-0.5">
                  <span v-if="s.created_order" class="text-[10px] font-bold px-1.5 py-0.5 rounded bg-emerald-100 text-emerald-700">ĐƠN #{{ s.created_order.id }}</span>
                  <span v-if="s.needs_human" class="text-[10px] font-bold px-1.5 py-0.5 rounded bg-amber-100 text-amber-800">CẦN NV</span>
                  <span
                    v-for="p in (s.detected_products || []).slice(0, 3)"
                    :key="p.id"
                    class="text-[10px] font-bold px-1.5 py-0.5 rounded"
                    :class="p.status === 'AVAILABLE' ? 'bg-blue-100 text-blue-700' : 'bg-slate-100 text-slate-500'"
                  >{{ p.code }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <p class="text-xs truncate m-0 flex-1" :class="isUnread(s) ? 'text-slate-900 font-semibold' : 'text-slate-500'">{{ lastMessagePreview(s) }}</p>
                  <span v-if="isUnread(s)" class="w-2 h-2 rounded-full bg-blue-600 shrink-0"></span>
                </div>
              </div>
            </button>
          </div>
        </div>

        <!-- Column 2: chat -->
        <div v-if="activeSession" class="col-span-5 bg-white rounded-2xl border border-slate-200 shadow-sm flex flex-col h-full min-h-0 overflow-hidden">
          <div class="px-4 py-3 border-b border-slate-100 flex items-center justify-between bg-slate-50/70">
            <div class="flex items-center gap-2.5 min-w-0">
              <img v-if="activeSession.avatar_url" :src="activeSession.avatar_url" class="w-8 h-8 rounded-full object-cover" />
              <div v-else class="w-8 h-8 rounded-full bg-blue-600 text-white font-bold text-xs flex items-center justify-center shrink-0 uppercase">
                {{ (activeSession.customer_name || 'K').charAt(0) }}
              </div>
              <h2 class="text-sm font-bold text-slate-900 m-0 truncate">{{ activeSession.customer_name }}</h2>
            </div>
            <span class="text-xs font-semibold px-2 py-1 rounded-md bg-slate-100 text-slate-600 shrink-0">{{ (activeSession.messages || []).length }} tin</span>
          </div>
          <div ref="chatContainer" class="flex-1 overflow-y-auto p-4 space-y-2 bg-[#f4f6f9]" @scroll="onChatScroll">
            <div
              v-for="(msg, idx) in activeSession.messages"
              :key="'dm-' + idx"
              class="flex"
              :class="msg.sender === 'shop' ? 'justify-end' : 'justify-start'"
            >
              <div
                class="max-w-[80%] rounded-2xl px-3.5 py-2 text-[13px] shadow-sm"
                :class="[
                  msg.sender === 'shop' ? 'bg-blue-600 text-white rounded-br-sm' : 'bg-white text-slate-800 border border-slate-200/80 rounded-bl-sm',
                  msg.pending ? 'opacity-60' : ''
                ]"
              >
                <div v-if="msg.type === 'image'">
                  <img
                    v-if="(msg.image_base64 || msg.image_url) && !brokenImages[msg.image_base64 || msg.image_url]"
                    :src="msg.image_base64 || msg.image_url"
                    referrerpolicy="no-referrer"
                    loading="lazy"
                    class="max-w-[230px] max-h-[260px] rounded-xl object-cover cursor-pointer block"
                    @load="onChatImageLoad"
                    @click="previewImage(msg.image_base64 || msg.image_url)"
                    @error="markImageBroken(msg.image_base64 || msg.image_url)"
                  />
                  <div v-else class="text-xs italic opacity-75 py-1">[Ảnh]</div>
                </div>
                <div v-else class="whitespace-pre-wrap break-words leading-relaxed">{{ msg.text }}</div>
              </div>
            </div>
          </div>
          <!-- Composer: AI suggestion pre-filled, editable, clear + send -->
        <div class="shrink-0 border-t border-slate-200 bg-white px-3 pt-2 pb-3" >
          <div v-if="pendingSend" class="flex items-center gap-2 bg-slate-900 text-white rounded-xl px-3 py-2 mb-2 text-[13px]">
            <i class="el-icon-loading"></i>
            <span class="flex-1">Đang gửi sau {{ pendingSend.left }}s…</span>
            <button class="font-bold text-amber-300" @click="undoSend">Hoàn tác</button>
          </div>
          <div class="flex items-center gap-1.5 text-[11px] font-bold mb-1" :class="replyDirty ? 'text-slate-500' : 'text-blue-700'">
            <i :class="reanalyzing[activeSession.thread_id] ? 'el-icon-loading' : (replyDirty ? 'el-icon-edit' : 'el-icon-magic-stick')"></i>
            <span class="flex-1">{{ reanalyzing[activeSession.thread_id] ? 'AI đang phân tích lại…' : (replyDirty ? 'Tin nhắn của bạn' : 'Gợi ý AI — sửa trực tiếp hoặc gửi luôn') }}</span>
            <button v-if="replyDirty && lastServerReply && editableReply !== lastServerReply" class="text-blue-600 font-semibold" @click="useSuggestion">Dùng gợi ý AI</button>
          </div>
          <div class="flex items-end gap-2">
            <div class="relative flex-1">
              <textarea
                v-model="editableReply"
                rows="3"
                placeholder="Nhập tin nhắn…"
                class="block w-full rounded-2xl border border-slate-300 bg-slate-50 pl-3 pr-10 py-2 text-[16px] leading-snug outline-none focus:border-blue-500 focus:bg-white resize-none"
                @input="replyDirty = true"
              ></textarea>
              <button
                v-if="editableReply"
                class="absolute top-1.5 right-1.5 w-7 h-7 rounded-full bg-slate-200 text-slate-600 flex items-center justify-center active:bg-slate-300"
                aria-label="Xoá nội dung để tự nhập"
                @click="clearReply"
              ><i class="el-icon-close"></i></button>
            </div>
            <button
              class="w-11 h-11 rounded-full bg-blue-600 text-white flex items-center justify-center text-lg active:bg-blue-700 disabled:opacity-40 shrink-0"
              aria-label="Gửi cho khách"
              :disabled="!canSendToCustomer || !editableReply.trim() || !!pendingSend"
              @click="queueSend"
            ><i class="el-icon-s-promotion"></i></button>
          </div>
        </div>
        </div>

        <!-- Column 3: AI + products + order -->
        <div v-if="activeSession" class="col-span-4 space-y-3 h-full min-h-0 overflow-y-auto pr-0.5">
          <div v-if="activeSession.needs_human" class="bg-amber-50 border border-amber-200 rounded-2xl p-3 text-xs text-amber-800">
            <b><i class="el-icon-warning-outline"></i> Cần nhân viên:</b> {{ activeSession.human_reason || 'khách hỏi thông tin chưa có dữ liệu' }}
          </div>

          <div class="bg-white rounded-2xl border border-slate-200 p-4 shadow-sm">
            <div class="flex items-center justify-between mb-2.5">
              <h3 class="text-sm font-bold text-slate-800 m-0">Sản phẩm nhận diện ({{ (activeSession.detected_products || []).length }})</h3>
              <el-button
                v-if="!activeSession.created_order"
                type="success"
                size="mini"
                icon="el-icon-shopping-cart-full"
                :loading="goingToOrder"
                @click="goCreateOrder"
              >Lên đơn</el-button>
            </div>
            <div v-if="activeSession.size_advice" class="mb-2.5 bg-blue-50 border border-blue-100 rounded-xl p-2.5 text-xs text-blue-800">
              <strong>Tư vấn size:</strong> {{ activeSession.size_advice }}
            </div>
            <div v-if="!(activeSession.detected_products || []).length" class="text-xs text-slate-400">Chưa phát hiện mã sản phẩm trong ảnh/tin nhắn.</div>
            <div v-else class="space-y-2">
              <div v-for="p in activeSession.detected_products" :key="'dp-' + p.id" class="flex items-start gap-3 p-2.5 rounded-xl border border-slate-100 bg-slate-50/80">
                <img v-if="p.thumb_url" :src="p.thumb_url" class="w-14 h-14 rounded-lg object-cover border border-slate-200 shrink-0" />
                <div class="flex-1 min-w-0">
                  <div class="flex items-center justify-between gap-2">
                    <router-link :to="`/salev2/product/${p.id}/edit`" class="font-bold text-sm text-blue-600 hover:underline">{{ p.code }}</router-link>
                    <span class="text-[10px] font-bold px-2 py-0.5 rounded-full" :class="p.status === 'AVAILABLE' ? 'bg-emerald-100 text-emerald-700' : 'bg-slate-200 text-slate-600'">
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

          <div v-if="activeSession.created_order" class="bg-emerald-50/70 rounded-2xl border border-emerald-200 p-4 shadow-sm">
            <div class="flex items-center justify-between mb-2">
              <div>
                <h3 class="text-sm font-bold text-emerald-800 m-0">Đã lên đơn #{{ activeSession.created_order.id }}</h3>
                <span class="text-xs text-emerald-700 font-semibold">Tổng: {{ formatPrice(activeSession.created_order.final_amount || activeSession.created_order.total_amount) }}</span>
              </div>
              <router-link
                :to="`/salev2/order/${activeSession.created_order.id}/edit`"
                class="px-3 py-1.5 rounded-lg bg-white border border-emerald-300 text-emerald-700 text-xs font-bold hover:bg-emerald-100"
              >Sửa đơn</router-link>
            </div>
            <div class="flex gap-2 mt-3">
              <el-button
                type="primary"
                size="small"
                icon="el-icon-s-promotion"
                :disabled="!receiptImg || !canSendToCustomer"
                :loading="sendingReceipt"
                @click="sendReceiptToCustomer"
              >Gửi receipt</el-button>
              <el-button size="small" icon="el-icon-picture-outline" :loading="renderingReceipt" @click="copyReceiptImage">Copy ảnh</el-button>
              <el-button v-if="receiptImg" size="small" icon="el-icon-download" @click="downloadReceipt">Tải</el-button>
            </div>
            <div v-if="receiptImg" class="mt-3 rounded-xl overflow-hidden border border-emerald-200 bg-white">
              <img :src="receiptImg" class="w-full h-auto block" />
            </div>
          </div>

        </div>
      </div>
    </div>

    <!-- Image Preview Modal -->
    <el-dialog :visible.sync="previewVisible" :width="isMobile ? '94%' : '520px'" append-to-body custom-class="copilot-preview">
      <img v-if="previewSrc" :src="previewSrc" class="w-full h-auto rounded-lg block" />
    </el-dialog>

    <!-- Hidden Receipt Template -->
    <!-- fixed (not absolute) so this off-screen render target never adds to the page height -->
    <div style="position: fixed; left: -10000px; top: 0; pointer-events: none;" aria-hidden="true">
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
const SEEN_KEY = 'copilot_seen_v1';
const UNDO_SECONDS = 5;

function readSeen() {
  try {
    return JSON.parse(localStorage.getItem(SEEN_KEY) || '{}') || {};
  } catch (e) {
    return {};
  }
}

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
      brokenImages: {},
      // mobile / list state
      isMobile: typeof window !== 'undefined' ? window.innerWidth < 768 : false,
      mobileView: 'list',
      sheet: null,
      replyDirty: false,
      goingToOrder: false,
      search: '',
      seen: readSeen(),
      flashIds: {},
      // send with undo
      pendingSend: null,
      renderedOrderId: null,
      reanalyzing: {},
      stickToBottom: true
    };
  },
  computed: {
    sortedSessions() {
      const q = this.search.trim().toLowerCase();
      const list = this.sessions.filter((s) => {
        if (!q) return true;
        if ((s.customer_name || '').toLowerCase().includes(q)) return true;
        return (s.detected_products || []).some((p) => (p.code || '').toLowerCase().includes(q));
      });
      return list.slice().sort((a, b) => this.sortTime(b) - this.sortTime(a));
    },
    activeSession() {
      if (!this.activeThreadId) return this.sortedSessions[0] || null;
      return this.sessions.find((s) => s.thread_id === this.activeThreadId) || null;
    },
    currentOrder() {
      return this.activeSession ? this.activeSession.created_order : null;
    },
    availableProducts() {
      return ((this.activeSession && this.activeSession.detected_products) || []).filter((p) => p.status === 'AVAILABLE');
    },
    canSendToCustomer() {
      return Boolean(this.activeSession && /^\d+$/.test(this.activeSession.thread_id || ''));
    }
  },
  mounted() {
    this.fetchSessions();
    this.connectSSE();
    window.addEventListener('resize', this.onResize);
    window.addEventListener('popstate', this.onPopState);
  },
  beforeDestroy() {
    if (this.eventSource) this.eventSource.close();
    window.removeEventListener('resize', this.onResize);
    window.removeEventListener('popstate', this.onPopState);
    this.clearPendingSend();
  },
  methods: {
    onResize() {
      this.isMobile = window.innerWidth < 768;
    },
    // ---------- data ----------
    async fetchSessions() {
      try {
        const res = await axios.get('/salev2/api/copilot/sessions');
        this.sessions = res.data.sessions || [];
        if (this.openFromRoute()) return;
        if (!this.activeThreadId && this.sortedSessions.length > 0 && !this.isMobile) {
          this.selectSession(this.sortedSessions[0].thread_id);
        } else if (this.activeSession) {
          this.syncActiveState(this.activeSession, true);
        }
      } catch (e) {
        console.error('Failed to fetch copilot sessions', e);
      }
    },
    async syncFromFacebook() {
      if (this.syncingFB) return;
      this.syncingFB = true;
      try {
        const res = await axios.post('/salev2/api/copilot/sync-fb');
        if (res.data && res.data.sessions) {
          this.sessions = res.data.sessions;
          this.$message.success(`Đã đồng bộ ${res.data.synced || 0} hội thoại`);
        }
      } catch (e) {
        const msg = (e.response && e.response.data && e.response.data.error) || 'Không thể đồng bộ từ Facebook';
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
          this.upsertSession(JSON.parse(event.data));
        } catch (e) {
          console.error('Invalid SSE session payload', e);
        }
      });
    },
    upsertSession(sess) {
      const idx = this.sessions.findIndex((s) => s.thread_id === sess.thread_id);
      const prev = idx !== -1 ? this.sessions[idx] : null;
      const prevCount = prev && prev.messages ? prev.messages.length : 0;
      const newCount = sess.messages ? sess.messages.length : 0;
      const isNewer = !prev || this.sortTime(sess) > this.sortTime(prev);

      if (idx !== -1) {
        this.$set(this.sessions, idx, sess);
      } else {
        this.sessions.push(sess);
      }
      // List is sorted by last message time (computed), so a new message moves it to the top.
      if (isNewer) this.flash(sess.thread_id);

      if (this.activeThreadId === sess.thread_id) {
        if (!this.isMobile || this.mobileView === 'chat') this.markSeen(sess);
        this.syncActiveState(sess, newCount > prevCount);
      }
    },
    flash(threadId) {
      this.$set(this.flashIds, threadId, true);
      setTimeout(() => this.$delete(this.flashIds, threadId), 2000);
    },
    // ---------- unread / time ----------
    sortTime(s) {
      if (!s) return 0;
      const last = s.last_message_at && !String(s.last_message_at).startsWith('0001') ? s.last_message_at : s.updated_at;
      const t = last ? new Date(last).getTime() : 0;
      return isNaN(t) ? 0 : t;
    },
    isUnread(s) {
      const msgs = s.messages || [];
      if (!msgs.length || msgs[msgs.length - 1].sender !== 'customer') return false;
      return this.sortTime(s) > (this.seen[s.thread_id] || 0);
    },
    markSeen(s) {
      if (!s) return;
      this.$set(this.seen, s.thread_id, Math.max(this.sortTime(s), Date.now()));
      try {
        localStorage.setItem(SEEN_KEY, JSON.stringify(this.seen));
      } catch (e) {
        /* storage unavailable: unread state is per page view only */
      }
    },
    // ---------- navigation ----------
    selectSession(threadId) {
      if (this.activeThreadId !== threadId) {
        this.clearPendingSend();
        this.replyDirty = false;
        this.lastServerReply = '';
        this.receiptImg = '';
        this.receiptCanvas = null;
      }
      this.activeThreadId = threadId;
      this.stickToBottom = true;
      if (this.activeSession) {
        this.markSeen(this.activeSession);
        this.syncActiveState(this.activeSession, true);
        if (this.activeSession.ai_stale) this.reanalyze(threadId);
      }
    },
    async reanalyze(threadId) {
      if (this.reanalyzing[threadId]) return;
      this.$set(this.reanalyzing, threadId, true);
      try {
        const res = await axios.post('/salev2/api/copilot/reanalyze', { thread_id: threadId });
        if (res.data && res.data.session) this.upsertSession(res.data.session);
      } catch (e) {
        console.error('Reanalyze failed', e);
      } finally {
        this.$delete(this.reanalyzing, threadId);
      }
    },
    onChatScroll(e) {
      const el = e.target;
      this.stickToBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 80;
    },
    onChatImageLoad() {
      if (this.stickToBottom) this.scrollChatToBottom();
    },
    openChat(threadId) {
      this.selectSession(threadId);
      this.mobileView = 'chat';
      this.sheet = null;
      // Let the phone's back gesture / button close the chat instead of leaving the page.
      try {
        window.history.pushState({ copilotChat: true }, '');
      } catch (e) {
        /* ignore */
      }
      this.$nextTick(() => this.scrollChatToBottom());
    },
    closeChat() {
      if (window.history.state && window.history.state.copilotChat) {
        window.history.back();
      } else {
        this.mobileView = 'list';
        this.sheet = null;
      }
    },
    onPopState() {
      if (this.sheet) {
        this.sheet = null;
      }
      if (this.mobileView === 'chat') {
        this.mobileView = 'list';
        this.clearPendingSend();
      }
    },
    syncActiveState(sess, shouldScroll) {
      const nextReply = sess.suggested_reply || '';
      if (nextReply !== this.lastServerReply) {
        this.lastServerReply = nextReply;
        // Never overwrite what the user is typing; they can pull the new suggestion back in.
        if (!this.replyDirty) this.editableReply = nextReply;
      }
      if (shouldScroll) this.$nextTick(() => this.scrollChatToBottom());
      if (sess.created_order && (!this.receiptImg || !this.renderedOrderId || this.renderedOrderId !== sess.created_order.id)) {
        this.renderReceipt(sess.created_order);
      } else if (!sess.created_order) {
        this.receiptImg = '';
        this.receiptCanvas = null;
      }
    },
    // Back from Create Order (?thread=...&receipt=1): reopen that chat with its receipt.
    openFromRoute() {
      const q = this.$route.query || {};
      const threadId = q.thread ? String(q.thread) : '';
      if (!threadId || !this.sessions.some((s) => s.thread_id === threadId)) return false;
      if (this.isMobile) {
        this.openChat(threadId);
        if (q.receipt && this.activeSession && this.activeSession.created_order) this.openOrderSheet();
      } else {
        this.selectSession(threadId);
      }
      this.$router.replace({ path: this.$route.path }).catch(() => {});
      return true;
    },
    // "Lên đơn": open the regular Create Order page pre-filled with customer + products + prices.
    async goCreateOrder() {
      const sess = this.activeSession;
      if (!sess || this.goingToOrder) return;
      this.goingToOrder = true;
      try {
        const draft = sess.order_draft || {};
        const avail = this.availableProducts.map((p) => p.code);
        let codes = (draft.product_codes || []).filter((c) => avail.includes(c));
        if (!codes.length && avail.length === 1) codes = avail;
        let matched = [];
        if (codes.length) {
          const res = await axios.get('/salev2/api/products', {
            params: { search: codes.join(','), status: 'AVAILABLE', limit: 100, include_empty_thumb: 'true' }
          });
          const list = (res.data && res.data.data) || [];
          matched = codes
            .map((c) => list.find((p) => (p.name || '').toUpperCase() === c.toUpperCase()))
            .filter(Boolean);
        }
        sessionStorage.setItem(
          'smart_order_prefill',
          JSON.stringify({
            parsed: {
              name: draft.name || sess.customer_name || '',
              phone: draft.phone || '',
              address: draft.address || '',
              note: draft.note || ''
            },
            matched_products: matched,
            copilot_thread: sess.thread_id
          })
        );
        this.$router.push({ path: '/salev2/order/create', query: { from_smart: '1' } });
      } catch (e) {
        this.$message.error('Không chuẩn bị được đơn: ' + (e.message || ''));
      } finally {
        this.goingToOrder = false;
      }
    },
    clearReply() {
      this.editableReply = '';
      this.replyDirty = true;
    },
    useSuggestion() {
      this.editableReply = this.lastServerReply;
      this.replyDirty = false;
    },
    openOrderSheet() {
      this.sheet = 'order';
      if (this.activeSession && this.activeSession.created_order && !this.receiptImg) {
        this.renderReceipt(this.activeSession.created_order);
      }
    },
    scrollChatToBottom() {
      this.stickToBottom = true;
      ['chatContainer', 'mobileChat'].forEach((ref) => {
        const el = this.$refs[ref];
        if (el) el.scrollTop = el.scrollHeight;
      });
    },
    markImageBroken(src) {
      this.$set(this.brokenImages, src, true);
    },
    previewImage(src) {
      this.previewSrc = src;
      this.previewVisible = true;
    },
    queueSend() {
      const text = (this.editableReply || '').trim();
      if (!text || !this.canSendToCustomer || this.pendingSend) return;
      const threadId = this.activeSession.thread_id;
      const pending = { text, threadId, left: UNDO_SECONDS, timer: null };
      pending.timer = setInterval(() => {
        pending.left -= 1;
        if (pending.left <= 0) {
          this.clearPendingSend();
          this.doSend(threadId, text);
        }
      }, 1000);
      this.pendingSend = pending;
    },
    undoSend() {
      this.clearPendingSend();
      this.$message.info('Đã huỷ gửi');
    },
    clearPendingSend() {
      if (this.pendingSend && this.pendingSend.timer) clearInterval(this.pendingSend.timer);
      this.pendingSend = null;
    },
    async doSend(threadId, text) {
      const sess = this.sessions.find((s) => s.thread_id === threadId);
      const optimistic = { sender: 'shop', type: 'text', text, pending: true };
      if (sess) {
        if (!sess.messages) this.$set(sess, 'messages', []);
        sess.messages.push(optimistic);
        this.$nextTick(() => this.scrollChatToBottom());
      }
      try {
        await axios.post('/salev2/api/copilot/send-text', { thread_id: threadId, text });
        this.$message.success('Đã gửi cho khách');
        if (this.activeThreadId === threadId && this.editableReply.trim() === text) {
          this.editableReply = '';
          this.replyDirty = false;
        }
      } catch (e) {
        if (sess && sess.messages) {
          const i = sess.messages.indexOf(optimistic);
          if (i !== -1) sess.messages.splice(i, 1);
        }
        const msg = (e.response && e.response.data && e.response.data.error) || e.message;
        this.$message.error('Gửi thất bại: ' + msg);
      }
    },
    // ---------- receipt ----------
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
            new Promise((resolve) => setTimeout(resolve, 2500))
          ]);
          const canvas = await html2canvas(el, { useCORS: true, scale: 2 });
          this.receiptCanvas = canvas;
          this.receiptImg = canvas.toDataURL('image/png');
          this.renderedOrderId = order.id;
        }
      } catch (e) {
        console.error('Render receipt error', e);
      } finally {
        this.renderingReceipt = false;
      }
    },
    async copyReceiptImage() {
      if (!this.receiptCanvas && this.currentOrder) await this.renderReceipt(this.currentOrder);
      if (!this.receiptCanvas) return;
      try {
        this.receiptCanvas.toBlob(async (blob) => {
          if (!blob) return;
          try {
            await navigator.clipboard.write([new ClipboardItem({ 'image/png': blob })]);
            this.$message.success('Đã copy ảnh receipt');
          } catch (err) {
            this.$message.error('Máy không hỗ trợ copy ảnh, hãy bấm Tải hoặc Gửi receipt');
          }
        }, 'image/png');
      } catch (e) {
        this.$message.error('Máy không hỗ trợ copy ảnh, hãy bấm Tải hoặc Gửi receipt');
      }
    },
    async sendReceiptToCustomer() {
      if (!this.receiptImg || !this.activeSession || this.sendingReceipt) return;
      try {
        await this.$confirm(`Gửi receipt đơn #${this.currentOrder ? this.currentOrder.id : ''} cho ${this.activeSession.customer_name}?`, 'Gửi receipt', {
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
        this.$message.success('Đã gửi receipt cho khách');
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
    // ---------- order ----------
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
      const body = last.type === 'image' ? 'Đã gửi một ảnh' : last.text;
      return last.sender === 'shop' ? `Bạn: ${body}` : body;
    },
    formatPrice(val) {
      const n = Number(val || 0);
      return n.toLocaleString('vi-VN') + 'đ';
    },
    formatTime(ts) {
      if (!ts) return '';
      const d = new Date(ts);
      if (isNaN(d.getTime())) return '';
      const now = new Date();
      if (d.toDateString() === now.toDateString()) {
        return d.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' });
      }
      return d.toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit' });
    }
  }
};
</script>

<style>
.copilot-preview .el-dialog__body {
  padding: 8px 12px 16px;
}
</style>
