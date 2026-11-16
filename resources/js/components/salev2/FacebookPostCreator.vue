<template>
  <div class="px-4 md:px-6 py-6 max-w-7xl mx-auto space-y-6 pb-28 md:pb-24">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
      <div class="flex items-center gap-3">
        <div class="w-12 h-12 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center font-bold text-2xl shadow-inner">
          <i class="el-icon-s-promotion"></i>
        </div>
        <div>
          <h1 class="text-2xl font-black text-slate-800 tracking-tight">Đăng bài Fanpage chính (AI Content)</h1>
          <p class="text-slate-500 text-xs md:text-sm mt-0.5">Tạo bài đăng Facebook theo danh mục, tự động lồng ghép dịp lễ/mùa và sắp xếp thứ tự ảnh.</p>
        </div>
      </div>
      <div class="flex gap-2">
        <el-button icon="el-icon-back" plain @click="$router.push('/salev2/products')">Về kho hàng</el-button>
        <el-button type="primary" icon="el-icon-s-opportunity" @click="handleGeneratePreview" :loading="loadingPreview">
          Tạo preview bài đăng
        </el-button>
      </div>
    </div>

    <!-- Step 1: Chọn Danh mục & Tuỳ chọn Dịp lễ / Mùa -->
    <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 space-y-6">
      <div class="flex items-center justify-between border-b border-slate-100 pb-3">
        <div class="flex items-center gap-2">
          <span class="w-6 h-6 rounded-full bg-blue-600 text-white flex items-center justify-center text-xs font-bold">1</span>
          <h2 class="font-bold text-slate-800 text-base">Thiết lập bài đăng & Phong cách nội dung</h2>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Cột chọn loại sản phẩm -->
        <div class="space-y-3">
          <label class="text-xs font-black uppercase text-slate-400 tracking-wider block">Loại sản phẩm (Category)</label>
          <el-select v-model="selectedCategoryCode" placeholder="Chọn loại sản phẩm..." class="w-full" @change="onCategoryChange" size="medium">
            <el-option
              v-for="cat in categories"
              :key="cat.code"
              :label="cat.name + ' (' + cat.available_count + ' SP có sẵn)'"
              :value="cat.code">
              <span class="float-left font-medium">{{ cat.name }}</span>
              <span class="float-right text-slate-400 text-xs">{{ cat.available_count }} SP</span>
            </el-option>
          </el-select>
          <div v-if="currentCategory && currentCategory.album_url" class="text-xs text-slate-500 flex items-center gap-1">
            <i class="el-icon-link text-blue-500"></i>
            <span>Album liên kết:</span>
            <a :href="currentCategory.album_url" target="_blank" class="text-blue-600 hover:underline truncate max-w-xs">{{ currentCategory.album_url }}</a>
          </div>
        </div>

        <!-- Cột tuỳ chọn Dịp lễ / Sự kiện -->
        <div class="space-y-3 p-4 bg-slate-50 rounded-xl border border-slate-200/60">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <i class="el-icon-date text-blue-600"></i>
              <span class="text-xs font-black uppercase text-slate-600 tracking-wider">Viết content theo dịp lễ / mùa</span>
            </div>
            <el-switch v-model="includeEvent" active-color="#2563eb" inactive-color="#cbd5e1"></el-switch>
          </div>

          <!-- Khi bật tuỳ chọn dịp lễ -->
          <div v-if="includeEvent" class="space-y-2 pt-1">
            <div class="flex flex-wrap gap-1.5">
              <el-tag
                v-for="preset in eventPresets"
                :key="preset"
                size="small"
                :type="eventName === preset ? 'primary' : 'info'"
                class="cursor-pointer transition-all hover:scale-105"
                :effect="eventName === preset ? 'dark' : 'plain'"
                @click="eventName = preset">
                {{ preset }}
              </el-tag>
            </div>
            <el-input
              v-model="eventName"
              placeholder="Nhập dịp lễ hoặc sự kiện (VD: Giáng Sinh / Noel, Tết, Mùa lạnh...)"
              size="small"
              prefix-icon="el-icon-edit">
            </el-input>
            <p class="text-[11px] text-slate-400">AI sẽ tự động lồng ghép không khí này vào mở đầu bài viết (VD: Sắm đồ ấm diện đi chơi Noel, dự tiệc cuối năm...).</p>
          </div>

          <!-- Khi tắt tuỳ chọn dịp lễ -->
          <div v-else class="text-xs text-slate-400 italic py-1">
            Chế độ bài viết thường: AI sẽ viết theo phong cách Vintage cổ điển, chuẩn mực hằng ngày (Timeless Classic).
          </div>
        </div>
      </div>
    </div>

    <!-- Step 2: Content Preview & Soạn thảo -->
    <div v-if="hasPreviewLoaded" class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 space-y-4">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3 border-b border-slate-100 pb-3">
        <div class="flex items-center gap-2">
          <span class="w-6 h-6 rounded-full bg-blue-600 text-white flex items-center justify-center text-xs font-bold">2</span>
          <h2 class="font-bold text-slate-800 text-base">Nội dung bài viết (AI Generated)</h2>
        </div>
        <div class="flex items-center gap-2">
          <el-button size="small" icon="el-icon-refresh" type="warning" plain :loading="loadingRegen" @click="handleRegenContent">
            Viết lại với AI (Regen)
          </el-button>
          <el-button size="small" icon="el-icon-document-copy" @click="copyContent">
            Sao chép
          </el-button>
        </div>
      </div>

      <div class="space-y-2">
        <el-input
          type="textarea"
          v-model="postMessage"
          :rows="12"
          placeholder="Nội dung bài viết sẽ hiển thị ở đây..."
          class="font-sans text-sm leading-relaxed"
          resize="vertical">
        </el-input>
        <div class="flex justify-between items-center text-xs text-slate-400 px-1">
          <span>{{ postMessage.length }} ký tự</span>
          <span class="text-blue-600 font-medium">Bản thảo có thể chỉnh sửa tự do trước khi đăng</span>
        </div>
      </div>
    </div>

    <!-- Step 3: Quản lý & Sắp xếp ảnh sản phẩm -->
    <div v-if="hasPreviewLoaded" class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 space-y-4">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3 border-b border-slate-100 pb-3">
        <div class="flex items-center gap-2">
          <span class="w-6 h-6 rounded-full bg-blue-600 text-white flex items-center justify-center text-xs font-bold">3</span>
          <div>
            <h2 class="font-bold text-slate-800 text-base">Danh sách sản phẩm đăng bài</h2>
            <p class="text-slate-400 text-xs">Sắp xếp theo ngày tải lên mới nhất (Latest Upload). Kéo thả hoặc dùng nút mũi tên để đổi vị trí. 4 ảnh đầu tiên là ảnh hiển thị chính trên Feed Facebook.</p>
          </div>
        </div>

        <!-- Bộ đếm & Chọn nhanh -->
        <div class="flex flex-wrap items-center gap-2">
          <div class="px-3 py-1.5 rounded-lg border font-bold text-xs" :class="selectedCount > 120 ? 'bg-red-50 text-red-600 border-red-200' : 'bg-blue-50 text-blue-700 border-blue-200'">
            Đã chọn: {{ selectedCount }} / 120 ảnh
          </div>
          <el-button-group size="mini">
            <el-button @click="selectPreset(30)">Top 30</el-button>
            <el-button @click="selectPreset(60)">Top 60</el-button>
            <el-button @click="selectPreset(120)">Top 120</el-button>
            <el-button @click="selectAll(true)">Tất cả</el-button>
            <el-button @click="selectAll(false)">Bỏ chọn</el-button>
          </el-button-group>
        </div>
      </div>

      <!-- Cảnh báo nếu chọn quá 120 ảnh -->
      <el-alert
        v-if="selectedCount > 120"
        title="Số lượng ảnh vượt quá 120 ảnh (giới hạn tối đa của Facebook). Vui lòng bỏ chọn bớt sản phẩm."
        type="error"
        show-icon
        :closable="false">
      </el-alert>

      <!-- Grid danh sách ảnh sản phẩm -->
      <div v-if="products.length" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-3">
        <div
          v-for="(prod, index) in products"
          :key="prod.id"
          class="relative bg-slate-50 border rounded-xl overflow-hidden transition-all duration-200 flex flex-col group"
          :class="[
            prod.selected ? 'border-blue-500 shadow-sm ring-1 ring-blue-400/30' : 'border-slate-200 opacity-50 grayscale-[0.6]',
            draggedIndex === index ? 'opacity-30 border-dashed border-2 border-blue-600' : ''
          ]"
          draggable="true"
          @dragstart="onDragStart(index, $event)"
          @dragover.prevent="onDragOver(index, $event)"
          @drop="onDrop(index, $event)">

          <!-- Top Badge: Rank & Main feed indicator -->
          <div class="absolute top-1.5 left-1.5 z-10 flex gap-1 items-center">
            <span
              class="px-2 py-0.5 rounded-md text-[11px] font-black shadow-sm"
              :class="index < 4 && prod.selected ? 'bg-amber-400 text-amber-950 font-bold' : 'bg-slate-900/80 text-white'">
              #{{ index + 1 }}
            </span>
            <span v-if="index < 4 && prod.selected" class="hidden sm:inline-block px-1.5 py-0.5 rounded bg-blue-600 text-white text-[9px] font-bold shadow-sm">
              Bìa Feed
            </span>
          </div>

          <!-- Top Right: Checkbox -->
          <div class="absolute top-1.5 right-1.5 z-10">
            <el-checkbox v-model="prod.selected" class="custom-checkbox"></el-checkbox>
          </div>

          <!-- Product Image Container -->
          <div class="aspect-square w-full bg-slate-200 relative overflow-hidden cursor-move">
            <img :src="prod.thumb_url" :alt="prod.name" class="w-full h-full object-cover" loading="lazy">
          </div>

          <!-- Card Footer & Move controls -->
          <div class="p-2 bg-white flex flex-col justify-between flex-1 border-t border-slate-100">
            <div>
              <div class="flex items-center justify-between">
                <span class="font-bold text-xs text-slate-800 truncate">{{ prod.name }}</span>
                <span class="text-[11px] font-black text-blue-600">{{ formatPrice(prod.price) }}</span>
              </div>
              <div class="text-[10px] text-slate-400 truncate mt-0.5">
                {{ formatSizes(prod.sizes) }}
              </div>
            </div>

            <!-- Position Controls -->
            <div class="flex items-center justify-between gap-1 mt-2 pt-1.5 border-t border-slate-50">
              <el-button
                size="mini"
                icon="el-icon-arrow-left"
                circle
                :disabled="index === 0"
                @click.stop="moveItem(index, -1)"
                title="Dời lên">
              </el-button>
              <span class="text-[10px] text-slate-400 cursor-move">Kéo thả</span>
              <el-button
                size="mini"
                icon="el-icon-arrow-right"
                circle
                :disabled="index === products.length - 1"
                @click.stop="moveItem(index, 1)"
                title="Dời xuống">
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="text-center py-12 text-slate-400">
        <i class="el-icon-box text-5xl mb-2"></i>
        <p>Không có sản phẩm nào thuộc danh mục này đang ở trạng thái còn hàng.</p>
      </div>

      <!-- Tấm ảnh kết bài (Retro CTA End-Card) -->
      <div class="mt-6 p-4 bg-amber-50/60 rounded-xl border border-amber-200/80 flex flex-col md:flex-row items-center gap-4">
        <div class="w-24 h-36 flex-shrink-0 rounded-lg overflow-hidden border border-amber-300 shadow-sm bg-white">
          <img :src="endCardURL" alt="Retro End Card" class="w-full h-full object-cover">
        </div>
        <div class="flex-1 space-y-1 text-center md:text-left">
          <div class="flex items-center justify-center md:justify-start gap-2">
            <span class="px-2 py-0.5 rounded bg-amber-200 text-amber-900 text-xs font-black">TẤM ẢNH KẾT BÀI (CTA)</span>
            <span class="text-xs font-bold text-slate-700">Style Retro Vintage</span>
          </div>
          <h4 class="font-bold text-slate-800 text-sm">Xem thêm SP & Các mẫu còn hàng tại link trong mô tả</h4>
          <p class="text-xs text-slate-500">
            Tấm ảnh này sẽ luôn tự động được đính kèm ở <strong>vị trí cuối cùng</strong> của bài post. Khi người dùng click vào ảnh trên Facebook, mô tả (caption) của ảnh sẽ tự động hiển thị link Album của danh mục này:
            <span class="text-blue-600 font-mono font-medium block truncate mt-0.5">{{ currentCategory ? currentCategory.album_url : 'Chưa có link album' }}</span>
          </p>
        </div>
        <div class="flex-shrink-0">
          <el-switch v-model="includeEndCard" active-text="Đính kèm" active-color="#d97706"></el-switch>
        </div>
      </div>
    </div>

    <!-- Floating Footer Action Bar (Fixed) -->
    <div v-if="hasPreviewLoaded" class="fixed bottom-20 md:bottom-6 left-0 right-0 px-4 z-40 pointer-events-none flex justify-center">
      <div class="w-full max-w-5xl bg-white/95 backdrop-blur-md border border-slate-200/90 p-4 rounded-2xl shadow-2xl flex flex-col sm:flex-row items-center justify-between gap-4 pointer-events-auto transition-all duration-300">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-blue-100 text-blue-600 flex items-center justify-center font-bold shadow-inner flex-shrink-0">
            <i class="el-icon-check text-lg"></i>
          </div>
          <div>
            <div class="font-bold text-slate-800 text-sm">
              Sẵn sàng đăng: <span class="text-blue-600 font-black text-base">{{ selectedCount }}</span> sản phẩm
            </div>
            <div class="text-xs text-slate-400">Thứ tự các ảnh trên bài viết sẽ tuân thủ chính xác vị trí hiển thị ở trên.</div>
          </div>
        </div>

        <div class="flex gap-2 w-full sm:w-auto">
          <el-button
            type="primary"
            size="medium"
            icon="el-icon-upload2"
            class="w-full sm:w-auto font-bold px-8 shadow-md"
            :disabled="selectedCount === 0 || selectedCount > 120"
            :loading="publishing"
            @click="confirmPublish">
            🚀 Đăng bài lên Facebook Page chính
          </el-button>
        </div>
      </div>
    </div>

    <!-- Publishing Progress Modal -->
    <el-dialog
      :visible.sync="publishingModal"
      width="450px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="false"
      center>
      <div class="py-6 text-center space-y-4">
        <div class="w-16 h-16 rounded-full bg-blue-50 text-blue-600 flex items-center justify-center mx-auto text-3xl animate-bounce">
          <i class="el-icon-upload"></i>
        </div>
        <h3 class="text-lg font-bold text-slate-800">Đang xuất bản lên Facebook...</h3>
        <p class="text-xs text-slate-500 max-w-xs mx-auto">
          Hệ thống đang tải lần lượt {{ selectedCount }} ảnh và thiết lập bài viết theo đúng thứ tự. Quá trình này mất khoảng 30 - 60 giây.
        </p>
        <el-progress :percentage="publishProgress" :status="publishProgress === 100 ? 'success' : ''"></el-progress>
      </div>
    </el-dialog>

    <!-- Success Dialog -->
    <el-dialog
      :visible.sync="successDialog"
      title="Đăng bài thành công!"
      width="450px"
      center>
      <div class="text-center py-4 space-y-4">
        <div class="w-16 h-16 rounded-full bg-emerald-50 text-emerald-600 flex items-center justify-center mx-auto text-3xl">
          <i class="el-icon-circle-check"></i>
        </div>
        <p class="text-sm text-slate-600">
          Bài viết đã được xuất bản trực tiếp lên dòng thời gian của Fanpage chính <strong>Vintage Always True House</strong>!
        </p>
        <div v-if="publishedPostURL" class="p-3 bg-slate-50 rounded-xl border border-slate-200 text-xs break-all">
          <a :href="publishedPostURL" target="_blank" class="text-blue-600 font-bold hover:underline flex items-center justify-center gap-1">
            <span>Xem bài viết trên Facebook</span>
            <i class="el-icon-top-right"></i>
          </a>
        </div>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button @click="successDialog = false">Đóng</el-button>
        <el-button type="primary" @click="handleGeneratePreview">Tạo bài đăng tiếp theo</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'FacebookPostCreator',
  data() {
    return {
      categories: [],
      selectedCategoryCode: '',
      includeEvent: true,
      eventName: 'Giáng Sinh (Noel) & Mùa đông cuối năm',
      eventPresets: [
        'Giáng Sinh (Noel) & Mùa đông cuối năm',
        'Tết Nguyên Đán & Du xuân',
        'Tiệc tùng & Year-end party',
        'Mùa Thu Đông – Se lạnh',
        'Xả kho cuối mùa'
      ],
      hasPreviewLoaded: false,
      loadingPreview: false,
      loadingRegen: false,
      postMessage: '',
      products: [],
      draggedIndex: null,
      publishing: false,
      publishingModal: false,
      publishProgress: 0,
      publishTimer: null,
      successDialog: false,
      publishedPostURL: '',
      includeEndCard: true,
      endCardURL: 'https://res.cloudinary.com/leq5n17r/image/upload/v1794836672/sale/1794836670476969000_ummmmmmm.jpg'
    };
  },
  computed: {
    currentCategory() {
      return this.categories.find(c => c.code === this.selectedCategoryCode) || null;
    },
    selectedCount() {
      return this.products.filter(p => p.selected).length;
    }
  },
  created() {
    this.fetchInitData();
  },
  methods: {
    async fetchInitData() {
      try {
        const res = await axios.get('/salev2/api/facebook/post-init');
        this.categories = res.data.categories || [];
        if (res.data.suggested_event) {
          this.eventName = res.data.suggested_event;
        }
        if (res.data.end_card_url) {
          this.endCardURL = res.data.end_card_url;
        }
        // Tự động chọn category đầu tiên có sản phẩm
        if (this.categories.length > 0) {
          const withCount = this.categories.find(c => c.available_count > 0);
          this.selectedCategoryCode = withCount ? withCount.code : this.categories[0].code;
        }
      } catch (err) {
        this.$message.error('Lỗi tải danh mục: ' + (err.response?.data?.error || err.message));
      }
    },
    onCategoryChange() {
      this.hasPreviewLoaded = false;
      this.products = [];
      this.postMessage = '';
    },
    async handleGeneratePreview() {
      if (!this.selectedCategoryCode) {
        this.$message.warning('Vui lòng chọn một loại sản phẩm');
        return;
      }
      this.loadingPreview = true;
      try {
        // 1. Tải danh sách sản phẩm available
        const prodRes = await axios.get('/salev2/api/facebook/category-products?type=' + encodeURIComponent(this.selectedCategoryCode));
        const items = (prodRes.data.data || []).map(p => ({
          ...p,
          selected: true
        }));
        this.products = items;

        // Giới hạn mặc định tối đa 60 ảnh để người dùng dễ quan sát
        if (this.products.length > 60) {
          this.products.forEach((p, idx) => {
            p.selected = idx < 60;
          });
        }

        // 2. Tạo nội dung bài viết bằng AI
        await this.generateContent();

        this.hasPreviewLoaded = true;
      } catch (err) {
        this.$message.error('Lỗi khởi tạo preview: ' + (err.response?.data?.error || err.message));
      } finally {
        this.loadingPreview = false;
      }
    },
    async generateContent() {
      const cat = this.currentCategory;
      const res = await axios.post('/salev2/api/facebook/generate-content', {
        category_code: this.selectedCategoryCode,
        category_name: cat ? cat.name : '',
        album_url: cat ? cat.album_url : '',
        include_event: this.includeEvent,
        event_name: this.includeEvent ? this.eventName : ''
      });
      this.postMessage = res.data.content || '';
    },
    async handleRegenContent() {
      this.loadingRegen = true;
      try {
        await this.generateContent();
        this.$message.success('Đã viết lại nội dung mới thành công!');
      } catch (err) {
        this.$message.error('Lỗi tạo lại nội dung: ' + (err.response?.data?.error || err.message));
      } finally {
        this.loadingRegen = false;
      }
    },
    copyContent() {
      if (!this.postMessage) return;
      navigator.clipboard.writeText(this.postMessage).then(() => {
        this.$message.success('Đã sao chép nội dung vào khay nhớ tạm!');
      }).catch(() => {
        this.$message.info('Vui lòng bôi đen và sao chép thủ công.');
      });
    },
    selectPreset(count) {
      this.products.forEach((p, idx) => {
        p.selected = idx < count;
      });
    },
    selectAll(select) {
      this.products.forEach(p => {
        p.selected = select;
      });
    },
    moveItem(fromIndex, direction) {
      const toIndex = fromIndex + direction;
      if (toIndex < 0 || toIndex >= this.products.length) return;
      const item = this.products.splice(fromIndex, 1)[0];
      this.products.splice(toIndex, 0, item);
    },
    onDragStart(index, event) {
      this.draggedIndex = index;
      event.dataTransfer.effectAllowed = 'move';
      event.dataTransfer.setData('text/plain', index);
    },
    onDragOver(index, event) {
      event.preventDefault();
      event.dataTransfer.dropEffect = 'move';
    },
    onDrop(toIndex, event) {
      event.preventDefault();
      const fromIndex = this.draggedIndex;
      if (fromIndex !== null && fromIndex !== toIndex) {
        const item = this.products.splice(fromIndex, 1)[0];
        this.products.splice(toIndex, 0, item);
      }
      this.draggedIndex = null;
    },
    formatPrice(val) {
      if (!val) return '0đ';
      return Number(val).toLocaleString('vi-VN') + 'đ';
    },
    formatSizes(sizes) {
      if (!sizes || !sizes.length) return '';
      return sizes.map(s => (s.name || '') + ': ' + (s.value || '')).join(' | ');
    },
    confirmPublish() {
      if (this.selectedCount === 0) {
        this.$message.warning('Vui lòng chọn ít nhất 1 sản phẩm để đăng');
        return;
      }
      if (this.selectedCount > 120) {
        this.$message.error('Số lượng ảnh đã chọn vượt quá giới hạn tối đa 120 ảnh');
        return;
      }

      this.$confirm(
        `Bạn có chắc chắn muốn xuất bản bài viết gồm ${this.selectedCount} ảnh lên Fanpage chính Vintage Always True House?`,
        'Xác nhận đăng bài',
        {
          confirmButtonText: 'Đăng ngay',
          cancelButtonText: 'Xem lại',
          type: 'warning'
        }
      ).then(() => {
        this.executePublish();
      }).catch(() => {});
    },
    async executePublish() {
      this.publishing = true;
      this.publishingModal = true;
      this.publishProgress = 10;

      // Giả lập tiến trình tăng dần để người dùng an tâm
      this.publishTimer = setInterval(() => {
        if (this.publishProgress < 90) {
          this.publishProgress += 5;
        }
      }, 1000);

      try {
        const selectedIds = this.products.filter(p => p.selected).map(p => p.id);

        const res = await axios.post('/salev2/api/facebook/publish-main-post', {
          message: this.postMessage,
          product_ids: selectedIds,
          include_end_card: this.includeEndCard,
          category_code: this.selectedCategoryCode,
          album_url: this.currentCategory ? this.currentCategory.album_url : ''
        });

        this.publishProgress = 100;
        clearInterval(this.publishTimer);
        setTimeout(() => {
          this.publishingModal = false;
          this.publishedPostURL = res.data.post_url;
          this.successDialog = true;
        }, 500);

      } catch (err) {
        clearInterval(this.publishTimer);
        this.publishingModal = false;
        this.$alert(
          'Lỗi khi đăng bài: ' + (err.response?.data?.error || err.message),
          'Đăng bài thất bại',
          { type: 'error' }
        );
      } finally {
        this.publishing = false;
      }
    }
  }
};
</script>

<style scoped>
.custom-checkbox >>> .el-checkbox__inner {
  width: 20px;
  height: 20px;
  border-radius: 6px;
}
.custom-checkbox >>> .el-checkbox__inner::after {
  height: 9px;
  left: 6px;
  top: 2px;
}
</style>
