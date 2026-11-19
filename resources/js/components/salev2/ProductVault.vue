<template>
  <div class="px-4 md:px-6 py-4 md:py-8 pt-20 md:pt-8">
    <!-- Mobile Fixed Search Bar on Top -->
    <div 
      :class="showMobileSearch ? 'translate-y-0' : '-translate-y-full'"
      class="md:hidden fixed top-0 left-0 right-0 z-40 bg-white/95 backdrop-blur px-3 py-2.5 border-b border-slate-200 shadow-sm flex items-center gap-2 transition-transform duration-300 ease-in-out">
      <el-input 
        v-model="search" 
        placeholder="Search Name, code..." 
        prefix-icon="el-icon-search" 
        clearable 
        class="flex-1"
        @focus="isSearchFocused = true"
        @blur="isSearchFocused = false"
        @input="debouncedFetch">
      </el-input>
      <el-button type="primary" size="small" icon="el-icon-plus" circle @click="$router.push('/salev2/product/create')"></el-button>
    </div>

    <div class="flex justify-between items-center mb-6 mobile-hide">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 tracking-tight">Product Vault</h1>
        <p class="text-slate-500 mt-1">Browse and share available inventory</p>
      </div>
      <div class="flex gap-3">
        <el-button type="primary" plain icon="el-icon-s-promotion" @click="$router.push('/salev2/facebook-post')">Đăng FB Page chính</el-button>
        <el-button @click="$router.push('/salev2/product/create')" icon="el-icon-plus" plain>Create</el-button>
        <el-button type="primary" icon="el-icon-upload" @click="$router.push('/salev2/upload')">Bulk Upload</el-button>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <div class="mb-6">
      <el-tabs v-model="activeTab" @tab-click="handleTabClick" type="card" class="premium-tabs">
        <el-tab-pane label="Available" name="AVAILABLE"></el-tab-pane>
        <el-tab-pane label="All Products" name="ALL"></el-tab-pane>
      </el-tabs>
    </div>

    <!-- Dynamic Filters -->
    <div class="bg-white p-4 md:p-6 rounded-2xl shadow-sm border border-slate-100 mb-8 space-y-4">
      <!-- Mobile Toggle Header with Filters -->
      <div class="flex justify-between items-center md:hidden pb-2 mb-2 border-b border-slate-100 cursor-pointer" @click="mobileFiltersOpen = !mobileFiltersOpen">
          <div class="flex items-center gap-2">
              <i class="el-icon-s-operation text-blue-600"></i>
              <span class="font-bold text-slate-700 text-sm">Filters</span>
              <i :class="mobileFiltersOpen ? 'el-icon-arrow-up' : 'el-icon-arrow-down'" class="text-slate-400 font-bold"></i>
          </div>
      </div>

      <!-- Collapsible Content -->
      <div :class="mobileFiltersOpen ? 'block' : 'hidden md:block'">
        <div class="flex flex-col lg:flex-row gap-6 items-start">
            <!-- Search & Type Group -->
            <div class="w-full lg:w-64 flex-shrink-0 space-y-3">
                <div class="hidden md:block">
                    <label class="text-[10px] font-black uppercase text-slate-400 mb-1 block tracking-wider">Search</label>
                    <el-input v-model="search" placeholder="Name, code..." prefix-icon="el-icon-search" class="w-full" @input="debouncedFetch" clearable></el-input>
                </div>
                <div>
                    <label class="text-[10px] font-black uppercase text-slate-400 mb-1 block tracking-wider">Product Type</label>
                    <el-select v-model="filterType" placeholder="Any Type" clearable class="w-full" @change="fetchProducts(1)">
                        <el-option v-for="t in productTypes" :key="t" :label="t" :value="t"></el-option>
                    </el-select>
                </div>
            </div>
            
            <!-- Dynamic Size Filters (Responsive Flow) -->
            <div class="flex-1 w-full">
            <div class="flex flex-col md:flex-row flex-wrap items-end gap-3">
                <!-- Filter Blocks -->
                <div v-for="(filter, index) in sizeFilters" :key="index" class="flex flex-wrap md:flex-nowrap items-end gap-2 p-3 bg-slate-50/50 rounded-xl border border-dashed border-slate-200">
                    <!-- Key (Mobile: Order 1, Desktop: Order 1) -->
                    <div class="w-[calc(100%-2.5rem)] md:w-32 order-1 md:order-1">
                        <label class="text-[10px] font-black uppercase text-slate-400 mb-1 block tracking-wider">Key</label>
                        <el-select v-model="filter.key" placeholder="Key" clearable class="w-full" size="mini" @change="debouncedFetch">
                        <el-option v-for="s in sizes" :key="s" :label="s" :value="s"></el-option>
                        </el-select>
                    </div>
                    
                    <!-- From (Mobile: Order 3, Desktop: Order 2) -->
                    <div class="w-[48%] md:w-28 order-3 md:order-2">
                        <label class="text-[10px] font-black uppercase text-slate-400 mb-1 block tracking-wider">From</label>
                        <el-input-number v-model="filter.from" :precision="1" :step="0.5" placeholder="From" class="w-full" controls-position="right" size="mini" @change="debouncedFetch"></el-input-number>
                    </div>
                    
                    <!-- To (Mobile: Order 4, Desktop: Order 3) -->
                    <div class="w-[48%] md:w-28 order-4 md:order-3">
                            <label class="text-[10px] font-black uppercase text-slate-400 mb-1 block tracking-wider">To</label>
                            <el-input-number v-model="filter.to" :precision="1" :step="0.5" placeholder="To" class="w-full" controls-position="right" size="mini" @change="debouncedFetch"></el-input-number>
                    </div>
                    
                    <!-- Remove Button (Mobile: Order 2, Desktop: Order 4) -->
                    <div class="w-8 md:w-auto flex justify-end md:block mb-0.5 md:mb-0 order-2 md:order-4">
                        <el-button v-if="sizeFilters.length > 1" type="danger" icon="el-icon-minus" circle size="mini" @click="removeSizeFilter(index)" plain></el-button>
                        <div v-else class="w-7 h-7 hidden md:block"></div>
                    </div>
                </div>
                
                <!-- Add Button (Inline on Desktop, Bottom on Mobile) -->
                <div class="pb-3 md:pb-3.5">
                        <el-button v-if="sizeFilters.length < 3" type="primary" icon="el-icon-plus" circle size="small" @click="addSizeFilter" plain></el-button>
                </div>
            </div>
            </div>
        </div>

        <div class="flex flex-col md:flex-row justify-between items-stretch md:items-center gap-3 pt-4 border-t border-slate-100 mt-2">
            <!-- Left: Selection Actions -->
            <div class="flex flex-col md:flex-row items-stretch md:items-center gap-2.5 w-full md:w-auto">
                <!-- Sort by upload & Mobile Clear -->
                <div class="flex items-center justify-between gap-2 px-0.5 w-full md:w-auto">
                    <el-checkbox v-model="sortByUpload" @change="fetchProducts(1)" class="mr-0">
                        <span class="text-xs md:text-sm font-bold text-slate-700">Sort by upload</span>
                    </el-checkbox>
                    <!-- Clear button visible on mobile if items selected -->
                    <el-button size="mini" type="danger" plain @click="clearSelection" icon="el-icon-delete" v-if="selectedProducts.length" class="md:hidden text-xs py-1.5 px-2.5">
                        Clear ({{ selectedProducts.length }})
                    </el-button>
                </div>

                <!-- Action buttons: Grid on mobile, inline on desktop -->
                <div class="grid grid-cols-3 md:flex md:items-center gap-2 w-full md:w-auto">
                    <el-button size="small" @click="selectAllPage" icon="el-icon-check" plain class="w-full md:w-auto text-xs md:text-sm px-1.5 md:px-3 justify-center">Select Page</el-button>
                    <el-button size="small" type="success" :disabled="!selectedProducts.length" icon="el-icon-share" @click="shareSelected" plain class="w-full md:w-auto text-xs md:text-sm px-1.5 md:px-3 justify-center">
                        Share <span v-if="selectedProducts.length">({{ selectedProducts.length }})</span>
                    </el-button>
                    <el-button size="small" type="info" icon="el-icon-download" plain @click="downloadAll" :loading="downloading" class="w-full md:w-auto text-xs md:text-sm px-1.5 md:px-3 justify-center">Download</el-button>
                    <!-- Desktop Clear button -->
                    <el-button size="small" type="danger" plain @click="clearSelection" icon="el-icon-delete" v-if="selectedProducts.length" class="hidden md:inline-flex px-3">Clear</el-button>
                </div>
            </div>

            <!-- Right: Apply/Reset Actions -->
            <div class="grid grid-cols-2 md:flex gap-2 w-full md:w-auto mt-1 md:mt-0">
                <el-button size="medium" @click="resetFilters" class="w-full md:w-auto">Reset</el-button>
                <el-button type="primary" size="medium" @click="fetchProducts(1)" class="w-full md:w-auto px-4 md:px-6 shadow-md shadow-blue-200">Apply Filters</el-button>
            </div>
        </div>
      </div>
    </div>

    <!-- Product Grid -->
    <div v-loading="loading">
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-3">
        <div v-for="product in products" :key="product.id" 
             class="group relative bg-white rounded-xl overflow-hidden border border-slate-100 shadow-sm hover:shadow-xl transition-all duration-300"
             :class="{'ring-2 ring-blue-500': isSelected(product.id)}">
          
          <!-- Image Area -->
          <div class="aspect-[3/4] overflow-hidden bg-slate-50 relative group">
            <img v-if="product.image_thumb_scale_url || product.path_thumb" :src="product.image_thumb_scale_url || product.path_thumb" class="w-full h-full object-contain hover:scale-105 transition-transform duration-500 p-1" @error="e => e.target.src = product.path_thumb">
            <div v-else class="w-full h-full flex items-center justify-center text-slate-300"><i class="el-icon-picture text-4xl"></i></div>
            
            <!-- Quick Actions Overlay (Desktop Only) -->
            <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity hidden md:flex flex-col justify-center items-center gap-2 p-6 z-10"
                 @click.stop="editProduct(product)">
                <el-button type="primary" size="medium" class="w-40 shadow-xl transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300 font-bold" 
                           @click.stop="editProduct(product)">DETAIL</el-button>
                <el-button type="success" size="medium" class="w-40 shadow-xl transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300 delay-75 font-bold" 
                           icon="el-icon-document-copy" @click.stop="copyProductInfo(product)">COPY SIZE</el-button>
                <el-button type="info" size="medium" class="w-40 shadow-xl transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300 delay-100 font-bold" 
                           icon="el-icon-camera" @click.stop="openShareModal(product, 'main')">SHARE MAIN</el-button>
                <el-button type="warning" size="medium" class="w-40 shadow-xl transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300 delay-125 font-bold" 
                           icon="el-icon-picture-outline" @click.stop="openShareModal(product, 'detail')">SHARE DETAIL</el-button>
                <el-button type="danger" size="medium" class="w-40 shadow-xl transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300 delay-150 font-bold" 
                           icon="el-icon-delete" @click.stop="confirmDelete(product)">DELETE</el-button>
            </div>
          </div>

          <!-- Product Info -->
          <div class="p-3">
            <div class="text-center mb-2">
                <h4 class="font-bold text-slate-800 text-sm uppercase leading-tight truncate">
                    {{ product.name }}
                </h4>
                <p class="text-blue-600 font-bold text-sm">{{ formatPrice(product.price) }}</p>
                <p class="text-[10px] text-slate-400 mt-0.5">{{ product.latest_upload ? formatDate(product.latest_upload) : 'undefined' }}</p>
            </div>
            
            <div class="text-center mb-3 px-1">
                <span class="text-[10px] font-black bg-blue-50 text-blue-600 px-2 py-1 rounded leading-tight uppercase inline-block">
                  {{ formatSizes(product).join(' - ') }}
                </span>
            </div>
            
            <!-- Footer: Mobile[Check-Status-Actions] | Desktop[Status-Check] -->
            <div class="flex justify-between items-center pt-2 border-t border-slate-50 gap-2">
                <!-- Checkbox: Order 1 (Mobile Left) / Order 2 (Desktop Right) -->
                <div class="order-1 md:order-2 flex items-center">
                    <el-checkbox :value="isSelected(product.id)" @change="toggleSelection(product)" class="transform scale-125"></el-checkbox>
                </div>

                <!-- Status: Order 2 (Mobile Center) / Order 1 (Desktop Left) -->
                <div class="order-2 md:order-1 flex-1 text-center md:text-left">
                    <el-tag :type="getStatusTagType(product.status)" size="mini" class="text-[8px] font-bold tracking-tighter uppercase rounded-md px-1">
                        {{ product.status }}
                    </el-tag>
                </div>

                <!-- Mobile Actions: Order 3 (Right) - Desktop Hidden -->
                <div class="order-3 md:hidden flex justify-end">
                    <el-dropdown trigger="click" @command="(cmd) => handleMobileAction(cmd, product)">
                        <span class="el-dropdown-link p-2">
                            <i class="el-icon-more transform rotate-90 text-lg font-bold text-slate-600"></i>
                        </span>
                        <el-dropdown-menu slot="dropdown">
                            <el-dropdown-item command="edit" icon="el-icon-edit">Detail</el-dropdown-item>
                            <el-dropdown-item command="copy" icon="el-icon-document-copy">Copy Size</el-dropdown-item>
                            <el-dropdown-item command="copy_main" icon="el-icon-camera">Share Main</el-dropdown-item>
                            <el-dropdown-item command="copy_detail" icon="el-icon-picture-outline">Share Detail</el-dropdown-item>
                            <el-dropdown-item command="delete" icon="el-icon-delete" divided class="text-red-500">Delete</el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="!products.length && !loading" class="text-center py-20 bg-white rounded-3xl border border-dashed border-slate-200">
          <i class="el-icon-search text-4xl text-slate-200 mb-4 block"></i>
          <p class="text-slate-400 font-medium">No products found</p>
      </div>

      <div class="mt-8 flex justify-center">
        <el-pagination background layout="prev, pager, next" :current-page.sync="page" :total="total" :page-size="60" @current-change="handlePageChange"></el-pagination>
      </div>
    </div>

    <!-- Share to Copilot / Copy Modal -->
    <div v-if="shareModalVisible && shareProduct"
         class="fixed inset-0 z-[80] flex items-end md:items-center justify-center bg-black/40 p-0 md:p-4"
         @click="closeShareModal">
      <div class="bg-white w-full md:max-w-md rounded-t-3xl md:rounded-2xl shadow-2xl flex flex-col max-h-[85vh] overflow-hidden"
           @click.stop>
        <!-- Header -->
        <div class="shrink-0 px-4 pt-4 pb-3 border-b border-slate-100">
          <div class="flex items-center gap-3">
            <img v-if="shareProduct.image_thumb_scale_url || shareProduct.path_thumb"
                 :src="shareProduct.image_thumb_scale_url || shareProduct.path_thumb"
                 class="w-11 h-11 rounded-lg object-cover bg-slate-100 shrink-0 border border-slate-200" />
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <span class="font-bold text-slate-900 text-sm uppercase truncate">{{ shareProduct.name }}</span>
                <span class="text-[10px] font-bold px-2 py-0.5 rounded-full shrink-0"
                      :class="shareKind === 'main' ? 'bg-blue-50 text-blue-600' : 'bg-amber-50 text-amber-700'">
                  {{ shareKind === 'main' ? 'Ảnh chính' : 'Main + Detail' }}
                </span>
              </div>
              <p class="text-xs text-slate-500 mt-0.5 m-0">Chọn khách Copilot để gửi thẳng qua Facebook</p>
            </div>
            <button type="button"
                    class="w-8 h-8 rounded-full bg-slate-100 hover:bg-slate-200 text-slate-600 flex items-center justify-center shrink-0 cursor-pointer"
                    @click="closeShareModal">
              <i class="el-icon-close"></i>
            </button>
          </div>
          <div class="mt-3 flex items-center gap-2 bg-slate-100 rounded-xl px-3 h-9">
            <i class="el-icon-search text-slate-400"></i>
            <input v-model="copilotSearch"
                   type="search"
                   placeholder="Tìm khách hàng Copilot..."
                   class="flex-1 bg-transparent outline-none text-sm text-slate-800" />
          </div>
        </div>

        <!-- Scrollable Copilot Customer List -->
        <div class="flex-1 overflow-y-auto divide-y divide-slate-100 min-h-[200px]">
          <div v-if="copilotLoading && copilotSessions.length === 0" class="py-12 text-center text-slate-400 text-sm">
            <i class="el-icon-loading text-2xl block mb-2"></i>
            Đang tải danh sách khách...
          </div>
          <div v-else-if="filteredCopilotSessions.length === 0" class="py-12 text-center text-slate-400 text-sm">
            <i class="el-icon-chat-dot-round text-3xl block mb-2 text-slate-300"></i>
            {{ copilotSearch ? 'Không tìm thấy khách phù hợp' : 'Chưa có khách nào bên Copilot' }}
          </div>
          <button v-for="s in filteredCopilotSessions"
                  :key="s.thread_id"
                  type="button"
                  :disabled="!!sendingThreadId"
                  @click="sendToCopilotCustomer(s)"
                  class="w-full text-left flex items-center gap-3 px-4 py-3 hover:bg-slate-50 active:bg-slate-100 transition-colors cursor-pointer disabled:opacity-60">
            <img v-if="s.avatar_url" :src="s.avatar_url" class="w-10 h-10 rounded-full object-cover shrink-0" />
            <div v-else class="w-10 h-10 rounded-full bg-slate-200 text-slate-600 font-bold flex items-center justify-center uppercase shrink-0">
              {{ (s.customer_name || 'K').charAt(0) }}
            </div>
            <div class="flex-1 min-w-0">
              <div class="font-semibold text-sm text-slate-900 truncate">{{ s.customer_name || 'Khách hàng' }}</div>
              <div class="text-xs text-slate-400 truncate">{{ copilotLastPreview(s) }}</div>
            </div>
            <span class="shrink-0 text-xs font-bold px-2.5 py-1 rounded-lg"
                  :class="sendingThreadId === s.thread_id ? 'bg-blue-100 text-blue-600' : 'bg-emerald-50 text-emerald-600'">
              <i v-if="sendingThreadId === s.thread_id" class="el-icon-loading mr-1"></i>
              {{ sendingThreadId === s.thread_id ? 'Đang gửi...' : 'Gửi FB' }}
            </span>
          </button>
        </div>

        <!-- Fixed Bottom Footer: Copy to clipboard/share like before -->
        <div class="shrink-0 border-t border-slate-200 bg-slate-50 px-4 py-3 flex items-center gap-2"
             style="padding-bottom: calc(env(safe-area-inset-bottom) + 12px)">
          <el-button type="primary"
                     class="w-full font-bold rounded-xl"
                     icon="el-icon-document-copy"
                     @click="copyFromShareModal">
            Copy {{ shareKind === 'main' ? 'Main' : 'Detail' }} (vào bộ nhớ)
          </el-button>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
export default {
  name: 'ProductVault',
  data() {
    return {
      activeTab: 'AVAILABLE',
      mobileFiltersOpen: false,
      products: [],
      loading: false,
      total: 0,
      page: 1,
      search: '',
      filterType: '',
      sizeFilters: [{ key: '', from: undefined, to: undefined }],
      selectedProducts: [],
      sortByUpload: false,
      downloading: false,
      showMobileSearch: true,
      lastScrollTop: 0,
      isSearchFocused: false,
      shareModalVisible: false,
      shareProduct: null,
      shareKind: 'main',
      copilotSessions: [],
      copilotLoading: false,
      copilotSearch: '',
      sendingThreadId: '',
      productTypes: ['TROUSERS', 'JACKET', 'SHIRT', 'BLAZER', 'TIE', 'GILE', 'BELT', 'POLO SHIRT', 'HAT', 'SUIT'],
      sizes: ['VAI', 'NGỰC', 'EO', 'DÀI ÁO', 'DÀI ÁO SAU', 'DÀI TAY', 'EO QUẦN', 'ĐÁY', 'ĐÙI', 'DÀI QUẦN', 'ỐNG', 'DƯ LAI'],
      sizeMapping: {
        'VAI': 'V', 'NGỰC': 'N', 'EO': 'E', 'DÀI ÁO': 'D', 'DÀI ÁO SAU': 'DS', 
        'DÀI TAY': 'DT', 'EO QUẦN': 'E', 'ĐÁY': 'ĐA', 'ĐÙI': 'Đ', 'DÀI QUẦN': 'D', 
        'ỐNG': 'Ô', 'DƯ LAI': 'L'
      },
      allSizes: ['VAI', 'NGỰC', 'EO', 'DÀI ÁO', 'DÀI ÁO SAU', 'DÀI TAY', 'EO QUẦN', 'ĐÁY', 'ĐÙI', 'DÀI QUẦN', 'ỐNG', 'DƯ LAI']
    }
  },
  computed: {
    filteredCopilotSessions() {
      const q = (this.copilotSearch || '').trim().toLowerCase();
      return (this.copilotSessions || [])
        .filter((s) => /^\d+$/.test(s.thread_id || ''))
        .filter((s) => !q || (s.customer_name || '').toLowerCase().includes(q))
        .slice()
        .sort((a, b) => this.copilotSortTime(b) - this.copilotSortTime(a));
    }
  },
  mounted() {
    this.fetchProducts();
    const stored = localStorage.getItem('selectedProducts');
    if (stored) this.selectedProducts = JSON.parse(stored);
    window.addEventListener('scroll', this.handleScroll, { capture: true, passive: true });
  },
  beforeDestroy() {
    window.removeEventListener('scroll', this.handleScroll, { capture: true });
  },
  methods: {
    getScrollTop() {
      const mainEl = document.querySelector('main');
      const mainScroll = mainEl ? mainEl.scrollTop : 0;
      const winScroll = window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop || 0;
      return Math.max(mainScroll, winScroll);
    },
    handleScroll() {
      if (this.isSearchFocused) return;
      const currentScroll = this.getScrollTop();
      const diff = currentScroll - this.lastScrollTop;

      if (currentScroll <= 15) {
        this.showMobileSearch = true;
        this.lastScrollTop = currentScroll;
        return;
      }

      if (diff > 8) {
        this.showMobileSearch = false;
      } else if (diff < -8) {
        this.showMobileSearch = true;
      }

      this.lastScrollTop = currentScroll;
    },
    formatDate(dateStr) {
      if (!dateStr) return 'undefined';
      const d = new Date(dateStr);
      return d.toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit', year: 'numeric' });
    },
    formatPrice(val) {
      if (!val) return '0đ';
      return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND', maximumFractionDigits: 0 }).format(val);
    },
    getStatusTagType(status) {
      const maps = { 'AVAILABLE': 'success', 'SOLD': 'info', 'ON_HOLD': 'warning' };
      return maps[status] || 'info';
    },
    formatSizes(product) {
      if (!product.sizes) return [];
      const results = [];
      const duLai = product.sizes.find(s => s.name === 'DƯ LAI');
      product.sizes.forEach(size => {
          if (size.name === 'DƯ LAI') return;
          let label = this.sizeMapping[size.name] || size.name;
          let value = Math.floor(size.value);
          if (size.name === 'DÀI QUẦN' && duLai) {
              results.push(`${label}${value} (+${Math.floor(duLai.value)})`);
          } else {
              results.push(`${label}${value}`);
          }
      });
      return results;
    },
    addSizeFilter() {
        if (this.sizeFilters.length < 3) {
            this.sizeFilters.push({ key: '', from: undefined, to: undefined });
        }
    },
    removeSizeFilter(index) {
        this.sizeFilters.splice(index, 1);
        this.fetchProducts(1);
    },
    async fetchProducts(page = 1) {
      this.page = page;
      this.loading = true;
      try {
        const params = {
          page,
          search: this.search,
          type: this.filterType,
          status: this.activeTab === 'ALL' ? '' : this.activeTab,
          from_search_key: this.sizeFilters[0]?.key,
          from_search_value: this.sizeFilters[0]?.from,
          to_search_value: this.sizeFilters[0]?.to,
          from_search_key_2: this.sizeFilters[1]?.key,
          from_search_value_2: this.sizeFilters[1]?.from,
          to_search_value_2: this.sizeFilters[1]?.to,
          from_search_key_3: this.sizeFilters[2]?.key,
          from_search_value_3: this.sizeFilters[2]?.from,
          to_search_value_3: this.sizeFilters[2]?.to,
        };
        if (this.sortByUpload) params.sort = 'upload';
        if (this.activeTab === 'ALL') params.include_empty_thumb = 'true';
        const res = await axios.get('/salev2/api/products', { params });
        this.products = res.data.data;
        this.total = res.data.total;
      } catch (err) {
        console.error(err);
      } finally {
        this.loading = false;
      }
    },
    handleTabClick() {
        this.fetchProducts(1);
    },
    debouncedFetch: _.debounce(function() { this.fetchProducts(1); }, 500),
    handlePageChange(page) {
      this.fetchProducts(page);
    },
    resetFilters() {
        this.search = '';
        this.filterType = '';
        this.sizeFilters = [{ key: '', from: undefined, to: undefined }];
        this.fetchProducts(1);
    },
    isSelected(id) {
        return this.selectedProducts.some(p => p.id === id);
    },
    toggleSelection(product) {
        const idx = this.selectedProducts.findIndex(p => p.id === product.id);
        if (idx > -1) {
            this.selectedProducts.splice(idx, 1);
        } else {
            this.selectedProducts.push({ id: product.id, img: product.image_thumb_scale_url || product.path_thumb });
        }
        localStorage.setItem('selectedProducts', JSON.stringify(this.selectedProducts));
    },
    selectAllPage() {
        this.products.forEach(p => {
            if (!this.isSelected(p.id)) {
                this.selectedProducts.push({ id: p.id, img: p.image_thumb_scale_url || p.path_thumb });
            }
        });
        localStorage.setItem('selectedProducts', JSON.stringify(this.selectedProducts));
    },
    clearSelection() {
        this.selectedProducts = [];
        localStorage.removeItem('selectedProducts');
    },
    async shareSelected() {
        if (!this.selectedProducts.length) {
            this.$message.warning('Please select some products first');
            return;
        }

        if (!navigator.share) {
            alert('Browser does not support native sharing (or not HTTPS). Try Mobile Safari/Chrome.');
            return;
        }

        const loading = this.$loading({ lock: true, text: 'Preparing images...', background: 'rgba(0,0,0,0.7)' });
        
        try {
            const files = [];
            for (const item of this.selectedProducts) {
                // Ensure the image URL is fetched correctly
                const response = await fetch(item.img);
                if (!response.ok) throw new Error('Failed to load image: ' + item.img);

                const blob = await response.blob();
                files.push(new File([blob], `${item.id}.jpg`, { type: 'image/jpeg' }));
            }

            if (navigator.canShare && navigator.canShare({ files: files })) {
                await navigator.share({
                    files: files,
                    title: 'Vintage Products',
                    text: 'Here are the selected products',
                });
            } else {
                alert('Browser does not support sharing these files (check limit or format).');
            }
        } catch (err) {
            console.error('Share failed', err);
            if (err.name !== 'AbortError') {
                alert('Share failed: ' + err.message);
            }
        } finally {
            loading.close();
        }
    },
    async downloadAll() {
      const images = this.products.filter(p => p.path_thumb).map(p => ({ url: p.path_thumb, name: p.name }));
      if (!images.length) { this.$message.warning('No images to download'); return; }
      this.downloading = true;
      let i = 0;
      const interval = setInterval(async () => {
        if (i >= images.length) {
          clearInterval(interval);
          this.downloading = false;
          return;
        }
        const { url, name } = images[i++];
        try {
          const res = await axios.get(url, { responseType: 'blob' });
          const a = document.createElement('a');
          a.href = window.URL.createObjectURL(res.data);
          a.download = name + '.jpg';
          a.click();
        } catch (e) { console.error('Download failed', e); }
      }, 200);
    },
    editProduct(product) {
        this.$router.push({ name: 'product-edit', params: { id: product.id } });
    },
    async copyProductInfo(product) {
        const sizes = this.formatSizes(product).join(' - ');
        const text = `${product.name} ${sizes}`;
        try {
            await navigator.clipboard.writeText(text);
            this.$message.success('Copied to clipboard!');
        } catch (err) {
            this.$message.error('Failed to copy');
        }
    },
    async copyMain(product) {
        if (!navigator.share) {
             alert('Browser does not support native sharing (or not HTTPS).');
             return;
        }

        if (!product.path_thumb) {
            this.$message.warning('This product has no main image.');
            return;
        }

        const loading = this.$loading({ lock: true, text: 'Preparing main image...', background: 'rgba(255,255,255,0.7)' });
        try {
            const response = await fetch(product.path_thumb);
            if (!response.ok) throw new Error('Failed to load image');
            
            const blob = await response.blob();
            const file = new File([blob], `Main-${product.name}.jpg`, { type: 'image/jpeg' });
            
            await navigator.share({ files: [file] });
        } catch (err) {
            if (err.name !== 'AbortError') alert('Share failed: ' + err.message);
        } finally {
            loading.close();
        }
    },
    async copyDetail(product) {
        if (!navigator.share) {
             alert('Browser does not support native sharing (or not HTTPS).');
             return;
        }

        if (!product.images || !product.images.length) {
            this.$message.warning('This product has no detail images.');
            return;
        }

        const loading = this.$loading({ lock: true, text: 'Preparing details...', background: 'rgba(255,255,255,0.7)' });
        try {
            const files = [];
            for (let i = 0; i < product.images.length; i++) {
                 const img = product.images[i];
                 const imgUrl = img.file_path; // Already formatted by controller
                 if (imgUrl) {
                    const response = await fetch(imgUrl);
                    if (response.ok) {
                        const blob = await response.blob();
                        files.push(new File([blob], `Detail-${product.id}-${i+1}.jpg`, { type: 'image/jpeg' }));
                    }
                 }
            }

            if (files.length === 0) {
                 this.$message.warning('Could not load images');
                 return;
            }

            await navigator.share({ files: files }); // Only files, no text
        } catch (err) {
            if (err.name !== 'AbortError') alert('Share failed: ' + err.message);
        } finally {
            loading.close();
        }
    },
    confirmDelete(product) {
        this.$confirm('This will permanently delete the product and its images. Continue?', 'Warning', {
          confirmButtonText: 'OK',
          cancelButtonText: 'Cancel',
          type: 'warning'
        }).then(() => {
          this.performDelete(product);
        }).catch(() => {});
    },
    async performDelete(product) {
        try {
            await axios.delete(`/salev2/api/product/${product.id}`);
            this.$message.success('Product deleted');
            this.fetchProducts(this.page);
        } catch (err) {
            this.$message.error('Delete failed');
        }
    },
    handleMobileAction(command, product) {
        if (command === 'edit') this.editProduct(product);
        if (command === 'copy') this.copyProductInfo(product);
        if (command === 'copy_main') this.openShareModal(product, 'main');
        if (command === 'copy_detail') this.openShareModal(product, 'detail');
        if (command === 'delete') this.confirmDelete(product);
    },
    openShareModal(product, kind) {
        this.shareProduct = product;
        this.shareKind = kind;
        this.copilotSearch = '';
        this.shareModalVisible = true;
        this.fetchCopilotSessions();
    },
    closeShareModal() {
        this.shareModalVisible = false;
        this.shareProduct = null;
        this.sendingThreadId = '';
    },
    async fetchCopilotSessions() {
        this.copilotLoading = true;
        try {
            const res = await axios.get('/salev2/api/copilot/sessions');
            this.copilotSessions = (res.data && res.data.sessions) || [];
        } catch (err) {
            this.$message.error('Không tải được danh sách khách Copilot');
        } finally {
            this.copilotLoading = false;
        }
    },
    copilotSortTime(s) {
        if (!s) return 0;
        const last = s.last_message_at && !String(s.last_message_at).startsWith('0001') ? s.last_message_at : s.updated_at;
        const t = last ? new Date(last).getTime() : 0;
        return isNaN(t) ? 0 : t;
    },
    copilotLastPreview(s) {
        const msgs = s.messages || [];
        if (!msgs.length) return 'Chưa có tin nhắn';
        const m = msgs[msgs.length - 1];
        if (m.type === 'image') return '[Hình ảnh]';
        return m.text || '';
    },
    async sendToCopilotCustomer(session) {
        if (!this.shareProduct || !session || this.sendingThreadId) return;
        this.sendingThreadId = session.thread_id;
        try {
            const res = await axios.post('/salev2/api/copilot/send-product-images', {
                thread_id: session.thread_id,
                product_id: this.shareProduct.id,
                kind: this.shareKind
            });
            const n = (res.data && res.data.sent) || 0;
            const kindLabel = this.shareKind === 'main' ? 'ảnh chính' : `${n} ảnh (chính + chi tiết)`;
            this.$message.success(`Đã gửi ${kindLabel} ${this.shareProduct.name} cho ${session.customer_name || 'khách'}`);
            this.closeShareModal();
        } catch (err) {
            const msg = (err.response && err.response.data && err.response.data.error) || err.message || 'Lỗi không xác định';
            this.$message.error('Gửi ảnh thất bại: ' + msg);
        } finally {
            this.sendingThreadId = '';
        }
    },
    copyFromShareModal() {
        if (!this.shareProduct) return;
        const product = this.shareProduct;
        const kind = this.shareKind;
        if (kind === 'main') {
            this.copyMain(product);
        } else {
            this.copyDetail(product);
        }
    }
  }
}
</script>

<style>
.premium-tabs .el-tabs__item {
    font-weight: 900;
    text-transform: uppercase;
    font-size: 12px;
    letter-spacing: 1px;
}
.premium-tabs .el-tabs__header {
    margin-bottom: 0;
}
</style>
