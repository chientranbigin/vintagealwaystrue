<template>
  <div class="px-4 md:px-6 py-6 max-w-4xl mx-auto">
    <!-- Header -->
    <div class="flex justify-between items-center mb-6">
        <div class="flex items-center gap-3">
             <el-button icon="el-icon-back" circle @click="$router.push('/salev2/products')"></el-button>
             <h1 class="text-2xl font-bold text-slate-800">Create New Product</h1>
        </div>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
        <div class="p-6 grid grid-cols-1 md:grid-cols-2 gap-8">
            <!-- Left: Basic Info -->
            <div class="space-y-5">
                <div>
                    <label class="block text-sm font-bold text-slate-700 mb-2">Product Type</label>
                    <el-select v-model="form.type" placeholder="Select Type" class="w-full" @change="handleTypeChange">
                         <el-option v-for="type in types" :key="type" :label="type" :value="type"></el-option>
                    </el-select>
                </div>

                <div>
                     <label class="block text-sm font-bold text-slate-700 mb-2">
                        Quick Input <span class="text-blue-500 font-normal ml-2 text-xs">({{ sizeConfig.labels }})</span>
                     </label>
                     <el-input v-model="quickInput" placeholder="Paste measurements (e.g., 47 54 76 60)" @input="parseQuickInput">
                         <template slot="prepend"><i class="el-icon-magic-stick"></i></template>
                     </el-input>
                     <p class="text-[10px] text-slate-400 mt-1">Enter numbers separated by space to auto-fill sizes below.</p>
                </div>

                <!-- Sizes Section Moved Here -->
                <div class="bg-slate-50 rounded-lg p-5 border border-slate-200">
                    <h3 class="font-bold text-slate-800 mb-4 border-b border-slate-200 pb-2">Measurements</h3>
                    <div v-if="form.type" class="grid grid-cols-2 gap-4">
                        <div v-for="(size, idx) in activeSizes" :key="idx">
                            <label class="block text-xs font-bold text-slate-500 mb-1 uppercase">{{ size }}</label>
                            <el-input v-model="form.sizes[idx].value" size="small">
                                    <template slot="append">cm</template>
                            </el-input>
                        </div>
                    </div>
                     <div v-else class="text-center py-4 text-slate-400 text-xs">Select Product Type first</div>
                </div>

                <div>
                    <label class="block text-sm font-bold text-slate-700 mb-2">Price</label>
                    <el-input-number v-model="form.price" :step="50000" class="w-full"></el-input-number>
                </div>

                <div>
                    <label class="block text-sm font-bold text-slate-700 mb-2">Description</label>
                    <el-input type="textarea" :rows="3" v-model="form.description" placeholder="Product details..."></el-input>
                </div>
            </div>

            <!-- Right: Image -->
            <div>
                 <label class="block text-sm font-bold text-slate-700 mb-2">Thumbnail Image (Optional)</label>
                 <div class="border-2 border-dashed border-slate-200 rounded-lg p-6 text-center hover:bg-slate-50 transition-colors relative cursor-pointer h-full min-h-[300px] flex flex-col items-center justify-center bg-white" @click="$refs.fileInput.click()">
                     <input type="file" ref="fileInput" class="hidden" accept="image/*" @change="handleFileChange">
                     <div v-if="previewUrl" class="w-full h-full flex items-center justify-center">
                         <img :src="previewUrl" class="max-h-[400px] w-full object-contain">
                     </div>
                     <div v-else>
                         <i class="el-icon-upload text-4xl text-slate-300 mb-3"></i>
                         <p class="text-sm text-slate-500 font-bold">Click to upload image</p>
                     </div>
                 </div>
            </div>
        </div>

        <!-- Footer -->
        <div class="px-6 py-4 bg-slate-50 border-t border-slate-100 flex justify-end gap-3">
             <el-button @click="$router.back()">Cancel</el-button>
             <el-button type="primary" :loading="loading" @click="submitByAjax">Create Product</el-button>
        </div>
    </div>
  </div>
</template>

<script>
export default {
    name: 'ProductCreate',
    data() {
        return {
            form: {
                type: '',
                price: 0,
                description: '',
                sizes: [], 
                image: null
            },
            quickInput: '',
            previewUrl: null,
            loading: false,
            types: [
                'TROUSERS', 'JACKET', 'SHIRT', 'BLAZER', 'TIE', 'GILE', 'BELT', 'POLO SHIRT', 'HAT', 'SUIT'
            ],
            // Copied from legacy JS logic and SaleConfig
            mappings: {
                'BLAZER': { labels: 'Vai - Ngực - Dài trước - Dài tay', indices: [0, 1, 3, 5] },
                'JACKET': { labels: 'Vai - Ngực - Dài - Dài tay', indices: [0, 1, 3, 5] },
                'TROUSERS': { labels: 'Eo - Đáy - Đùi - Dài quần - Ống - Dư lai', indices: [6, 7, 8, 9, 10, 11] },
                'SHIRT': { labels: 'Vai - Ngực - Dài - Dài tay', indices: [0, 1, 3, 5] },
                'GILE': { labels: 'Vai - Ngực - Dài - Dài tay', indices: [0, 1, 3, 5] },
                'POLO SHIRT': { labels: 'Vai - Ngực - Dài - Dài tay', indices: [0, 1, 3, 5] },
                'SUIT': { labels: 'Vai - Ngực - Dài - Dài tay', indices: [0, 1, 3, 5] },
            },
            allSizes: [
                'VAI', 'NGỰC', 'EO', 'DÀI ÁO', 'DÀI ÁO SAU', 'DÀI TAY', // Top (Indices 0-5)
                'EO QUẦN', 'ĐÁY', 'ĐÙI', 'DÀI QUẦN', 'ỐNG', 'DƯ LAI'   // Bottom (Indices 6-11)
            ],
            sizeMapping: {
                'VAI': 'V', 'NGỰC': 'N', 'EO': 'E', 'DÀI ÁO': 'D', 'DÀI ÁO SAU': 'DS', 
                'DÀI TAY': 'DT', 'EO QUẦN': 'E', 'ĐÁY': 'ĐA', 'ĐÙI': 'Đ', 'DÀI QUẦN': 'D', 
                'ỐNG': 'Ô', 'DƯ LAI': 'L'
            }
        }
    },
    computed: {
        sizeConfig() {
            return this.mappings[this.form.type] || { labels: '', indices: [] };
        },
        activeSizes() {
           // We need to determine which sizes to show based on type. 
           // Legacy logic hides/shows sections: AO (< 6) vs QUAN (>= 6).
           if(['TROUSERS', 'SHORT'].includes(this.form.type)) {
               return this.allSizes.slice(6); // The bottom half
           } else {
               return this.allSizes.slice(0, 6); // The top half
           } 
        }
    },
    methods: {
        handleTypeChange() {
            // Reset sizes array to match the global structure but we only fill what's needed
            // Actually, we should initialize form.sizes with {name, value} for ALL sizes to match legacy payload structure?
            // Legacy payload: sizes[0][name]=Vai, sizes[0][value]=...
            // So we should maintain the full array of 12 items.
            this.form.sizes = this.allSizes.map(name => ({ name, value: '' }));
            this.quickInput = '';
        },
        parseQuickInput(val) {
            if(!val) return;
            const numbers = val.trim().split(/\s+/);
            const indices = this.sizeConfig.indices;
            if(!indices || !indices.length) return;

            indices.forEach((globalIdx, i) => {
                if(numbers[i] && this.form.sizes[globalIdx]) {
                    this.form.sizes[globalIdx].value = numbers[i];
                }
            });
        },
        handleFileChange(e) {
            const file = e.target.files[0];
            if(file) {
                this.form.image = file;
                this.previewUrl = URL.createObjectURL(file);
            }
        },
        async submitByAjax() {
            if(!this.form.type) return this.$message.error('Product Type is required');
            
            this.loading = true;
            try {
                let formData = new FormData();
                formData.append('type', this.form.type);
                formData.append('price', this.form.price);
                formData.append('description', this.form.description || '');
                if (this.form.image) {
                    formData.append('banner', this.form.image);
                }
                
                this.form.sizes.forEach((s, i) => {
                    formData.append(`sizes[${i}][name]`, s.name);
                    formData.append(`sizes[${i}][value]`, s.value);
                });

                const res = await axios.post('/salev2/api/product/store', formData, {
                    headers: { 'Content-Type': 'multipart/form-data' }
                });

                const product = res.data.product;
                const sizesStr = this.formatSizes(product).join(' - ');
                
                this.$alert(`Mã SP: ${product.name}<br/>Thông số: ${sizesStr}`, 'Tạo sản phẩm thành công', {
                    confirmButtonText: 'OK',
                    dangerouslyUseHTMLString: true,
                    type: 'success'
                });

                // Reset for next entry (Keep Type)
                this.form.price = 0;
                this.form.description = '';
                this.form.image = null;
                this.previewUrl = null;
                this.quickInput = '';
                this.form.sizes.forEach(s => s.value = '');
                
            } catch(e) {
                console.error(e);
                this.$message.error('Failed to create product');
            } finally {
                this.loading = false;
            }
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
              } else if (value > 0) {
                  results.push(`${label}${value}`);
              }
          });
          return results;
        }
    },
    mounted() {
        this.handleTypeChange(); // Init empty sizes
    }
}
</script>
