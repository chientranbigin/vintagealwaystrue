<template>
  <div class="px-4 md:px-6 py-8">
    <div class="flex justify-between items-center mb-6 mobile-hide">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 tracking-tight">Product Vault</h1>
        <p class="text-slate-500 mt-1">Browse and share available inventory</p>
      </div>
      <div class="flex gap-3">
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
      <!-- Mobile Toggle Header with Filters and Create Action -->
      <div class="flex justify-between items-center md:hidden pb-2 mb-2 border-b border-slate-100 cursor-pointer">
          <div class="flex items-center gap-2" @click="mobileFiltersOpen = !mobileFiltersOpen">
              <i class="el-icon-s-operation text-blue-600"></i>
              <span class="font-bold text-slate-700 text-sm">Filters</span>
              <i :class="mobileFiltersOpen ? 'el-icon-arrow-up' : 'el-icon-arrow-down'" class="text-slate-400 font-bold"></i>
          </div>
          <el-button type="primary" size="mini" icon="el-icon-plus" circle @click="$router.push('/salev2/product/create')"></el-button>
      </div>

      <!-- Collapsible Content -->
      <div :class="mobileFiltersOpen ? 'block' : 'hidden md:block'">
        <div class="flex flex-col lg:flex-row gap-6 items-start">
            <!-- Search & Type Group -->
            <div class="w-full lg:w-64 flex-shrink-0 space-y-3">
                <div>
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
                        <el-select v-model="filter.key" placeholder="Key" clearable class="w-full" size="mini">
                        <el-option v-for="s in sizes" :key="s" :label="s" :value="s"></el-option>
                        </el-select>
                    </div>
                    
                    <!-- From (Mobile: Order 3, Desktop: Order 2) -->
                    <div class="w-[48%] md:w-28 order-3 md:order-2">
                        <label class="text-[10px] font-black uppercase text-slate-400 mb-1 block tracking-wider">From</label>
                        <el-input-number v-model="filter.from" :precision="1" :step="0.5" placeholder="From" class="w-full" controls-position="right" size="mini"></el-input-number>
                    </div>
                    
                    <!-- To (Mobile: Order 4, Desktop: Order 3) -->
                    <div class="w-[48%] md:w-28 order-4 md:order-3">
                            <label class="text-[10px] font-black uppercase text-slate-400 mb-1 block tracking-wider">To</label>
                            <el-input-number v-model="filter.to" :precision="1" :step="0.5" placeholder="To" class="w-full" controls-position="right" size="mini"></el-input-number>
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

        <div class="flex flex-col md:flex-row justify-between items-center gap-4 pt-4 border-t border-slate-100 mt-2">
            <!-- Left: Selection Actions (No Scroll, Compact) -->
            <div class="flex flex-nowrap w-full md:w-auto gap-2">
                <el-button size="small" @click="selectAllPage" icon="el-icon-check" plain class="flex-1 md:flex-none px-2">Select Page</el-button>
                <el-button size="small" type="success" :disabled="!selectedProducts.length" icon="el-icon-share" @click="shareSelected" plain class="flex-1 md:flex-none px-2">
                    Share <span v-if="selectedProducts.length">({{ selectedProducts.length }})</span>
                </el-button>
                <el-button size="small" type="danger" plain @click="clearSelection" icon="el-icon-delete" v-if="selectedProducts.length" class="flex-none px-3">Clear</el-button>
            </div>
            <!-- Right: Apply/Reset Actions -->
            <div class="flex gap-3 w-full md:w-auto">
                <el-button size="medium" @click="resetFilters" class="flex-1 md:flex-none">Reset</el-button>
                <el-button type="primary" size="medium" @click="fetchProducts(1)" class="flex-1 md:flex-none px-6 shadow-md shadow-blue-200">Apply Filters</el-button>
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
            <img v-if="product.path_thumb" :src="product.path_thumb" class="w-full h-full object-contain hover:scale-105 transition-transform duration-500 p-1">
            <div v-else class="w-full h-full flex items-center justify-center text-slate-300"><i class="el-icon-picture text-4xl"></i></div>
            
            <!-- Quick Actions Overlay -->
            <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col justify-center items-center gap-3 p-4">
                <el-button type="primary" size="small" class="w-28 shadow-lg transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300" 
                           @click="editProduct(product)">Edit</el-button>
                <el-button type="success" size="small" class="w-28 shadow-lg transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300 delay-75" 
                           icon="el-icon-document-copy" @click="copyProductInfo(product)">Copy</el-button>
            </div>
          </div>

          <!-- Product Info -->
          <div class="p-3">
            <div class="text-center mb-2">
                <h4 class="font-bold text-slate-800 text-sm uppercase leading-tight truncate">
                    {{ product.name }}
                </h4>
                <p class="text-blue-600 font-bold text-sm">{{ formatPrice(product.price) }}</p>
            </div>
            
            <div class="text-center mb-3 px-1">
                <span class="text-[10px] font-black bg-blue-50 text-blue-600 px-2 py-1 rounded leading-tight uppercase inline-block">
                  {{ formatSizes(product).join(' - ') }}
                </span>
            </div>
            
            <!-- Footer: Status Left, Checkbox Right -->
            <div class="flex justify-between items-center pt-2 border-t border-slate-50">
                <el-tag :type="getStatusTagType(product.status)" size="mini" class="text-[8px] font-bold tracking-tighter uppercase rounded-md px-1">
                    {{ product.status }}
                </el-tag>
                <el-checkbox :value="isSelected(product.id)" @change="toggleSelection(product)" class="transform scale-125"></el-checkbox>
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
      productTypes: ['TROUSERS', 'JACKET', 'SHIRT', 'BLAZER', 'TIE', 'GILE', 'BELT', 'POLO SHIRT', 'HAT', 'SUIT'],
      sizes: ['VAI', 'NGỰC', 'EO', 'DÀI ÁO', 'DÀI ÁO SAU', 'DÀI TAY', 'EO QUẦN', 'ĐÁY', 'ĐÙI', 'DÀI QUẦN', 'ỐNG', 'DƯ LAI'],
      sizeMapping: {
        'VAI': 'V', 'NGỰC': 'N', 'EO': 'E', 'DÀI ÁO': 'D', 'DÀI ÁO SAU': 'DS', 
        'DÀI TAY': 'DT', 'EO QUẦN': 'E', 'ĐÁY': 'ĐA', 'ĐÙI': 'Đ', 'DÀI QUẦN': 'D', 
        'ỐNG': 'Ô', 'DƯ LAI': 'L'
      },
      // Edit Modal Data
      editVisible: false,
      editLoading: false,
      editQuickInput: '',
      editForm: { id: null, type: '', price: 0, description: '', sizes: [], image: null, imagePreview: '' },
      mappings: {
            'BLAZER': { labels: 'Vai - Ngực - Dài trước - Dài tay', indices: [0, 1, 3, 5] },
            'JACKET': { labels: 'Vai - Ngực - Dài - Dài tay', indices: [0, 1, 3, 5] },
            'TROUSERS': { labels: 'Eo - Đáy - Đùi - Dài quần - Ống - Dư lai', indices: [6, 7, 8, 9, 10, 11] },
            'SHIRT': { labels: 'Vai - Ngực - Dài - Dài tay', indices: [0, 1, 3, 5] },
            'GILE': { labels: 'Vai - Ngực - Dài - Dài tay', indices: [0, 1, 3, 5] },
            'POLO SHIRT': { labels: 'Vai - Ngực - Dài - Dài tay', indices: [0, 1, 3, 5] },
            'SUIT': { labels: 'Vai - Ngực - Dài - Dài tay', indices: [0, 1, 3, 5] },
      },
      allSizes: ['VAI', 'NGỰC', 'EO', 'DÀI ÁO', 'DÀI ÁO SAU', 'DÀI TAY', 'EO QUẦN', 'ĐÁY', 'ĐÙI', 'DÀI QUẦN', 'ỐNG', 'DƯ LAI']
    }
  },
  mounted() {
    this.fetchProducts();
    const stored = localStorage.getItem('selectedProducts');
    if (stored) this.selectedProducts = JSON.parse(stored);
  },
  methods: {
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
              results.push(`${label}${value}(+${Math.floor(duLai.value)})`);
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
        };
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
            this.selectedProducts.push({ id: product.id, img: product.path_thumb });
        }
        localStorage.setItem('selectedProducts', JSON.stringify(this.selectedProducts));
    },
    selectAllPage() {
        this.products.forEach(p => {
            if (!this.isSelected(p.id)) {
                this.selectedProducts.push({ id: p.id, img: p.path_thumb });
            }
        });
        localStorage.setItem('selectedProducts', JSON.stringify(this.selectedProducts));
    },
    clearSelection() {
        this.selectedProducts = [];
        localStorage.removeItem('selectedProducts');
    },
    async shareSelected() {
        if (!navigator.share) {
            this.$message.warning('Browser does not support native sharing');
            return;
        }
        const loading = this.$loading({ lock: true, text: 'Preparing images...', background: 'rgba(0,0,0,0.7)' });
        try {
            const files = [];
            for (const item of this.selectedProducts) {
                const response = await fetch(item.img);
                const blob = await response.blob();
                files.push(new File([blob], `${item.id}.jpg`, { type: 'image/jpeg' }));
            }
            await navigator.share({
                files: files,
                title: 'Vintage Always True',
                text: 'Check out these items!',
            });
        } catch (err) {
            if (err.name !== 'AbortError') this.$message.error('Share failed: ' + err.message);
        } finally {
            loading.close();
        }
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
    editProduct(product) {
        this.editForm = {
            id: product.id,
            type: product.type,
            price: product.price,
            description: product.description,
            sizes: this.allSizes.map(name => {
                const existing = product.sizes.find(s => s.name === name);
                return { name, value: existing ? existing.value : '' };
            }),
            image: null,
            imagePreview: product.path_thumb
        };
        this.editVisible = true;
    },
    handleEditFileChange(e) {
        const file = e.target.files[0];
        if(file) {
            this.editForm.image = file;
            this.editForm.imagePreview = URL.createObjectURL(file);
        }
    },
    parseEditQuickInput(val) {
        if(!val) return;
        const nums = val.trim().split(/\s+/);
        const config = this.mappings[this.editForm.type];
        if(config && config.indices) {
            config.indices.forEach((idx, i) => {
                if(nums[i]) this.editForm.sizes[idx].value = nums[i];
            });
        }
    },
    async saveProduct() {
        this.editLoading = true;
        try {
            let formData = new FormData();
            formData.append('type', this.editForm.type);
            formData.append('price', this.editForm.price);
            formData.append('description', this.editForm.description || '');
            if(this.editForm.image) {
                formData.append('banner', this.editForm.image);
            }
            this.editForm.sizes.forEach((s, i) => {
                formData.append(`sizes[${i}][name]`, s.name);
                formData.append(`sizes[${i}][value]`, s.value);
            });

            await axios.post(`/salev2/api/product/${this.editForm.id}/update`, formData, {
                 headers: { 'Content-Type': 'multipart/form-data' }
            });
            
            this.$message.success('Product Updated');
            this.editVisible = false;
            this.fetchProducts(this.page);
        } catch(e) {
            this.$message.error('Update failed');
        } finally {
            this.editLoading = false;
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
