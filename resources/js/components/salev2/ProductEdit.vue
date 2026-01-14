<template>
  <div class="px-4 md:px-6 py-8">
     <div class="flex justify-between items-center mb-6">
       <div>
         <h1 class="text-3xl font-bold text-gray-900 tracking-tight">Edit Product: {{ form.name }}</h1>
         <p class="text-slate-500 mt-1">Update product details and measurements</p>
       </div>
       <el-button @click="$router.push('/salev2/products')" icon="el-icon-back" circle></el-button>
     </div>

     <div class="grid grid-cols-1 md:grid-cols-3 gap-6" v-loading="loading">
         <!-- Left Column: Basic Info & Image -->
         <div class="space-y-6">
             <!-- Type Selection (Read-only/Limited Edit) -->
            <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
                 <label class="block text-slate-400 font-bold text-xs uppercase mb-2 tracking-wider">Product Type</label>
                 <el-select v-model="form.type" placeholder="Select Type" class="w-full" disabled>
                     <el-option v-for="t in productTypes" :key="t" :label="t" :value="t"></el-option>
                 </el-select>
            </div>

            <!-- Image Upload -->
             <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
                <label class="block text-slate-400 font-bold text-xs uppercase mb-2 tracking-wider">Thumbnail Image</label>
                <div class="relative group cursor-pointer" @click="$refs.fileInput.click()">
                    <div v-if="imagePreview" class="aspect-[3/4] rounded-lg overflow-hidden bg-slate-50 border border-slate-200">
                        <img :src="imagePreview" class="w-full h-full object-contain">
                    </div>
                    <div v-else class="aspect-[3/4] rounded-lg bg-slate-50 border-2 border-dashed border-slate-200 flex flex-col items-center justify-center text-slate-400">
                        <i class="el-icon-picture text-3xl mb-2"></i>
                        <span class="text-xs font-bold uppercase">Upload Image</span>
                    </div>
                    <!-- Overlay -->
                     <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center rounded-lg">
                         <span class="text-white font-bold uppercase tracking-widest text-xs border border-white px-4 py-2 rounded">Change</span>
                     </div>
                 </div>
                 <input type="file" ref="fileInput" class="hidden" accept="image/*" @change="handleFileChange">
            </div>

            <!-- Detail Images Gallery -->
            <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100" v-if="form.details && form.details.length">
                <div class="flex justify-between items-center mb-2">
                    <label class="block text-slate-400 font-bold text-xs uppercase tracking-wider">Detail Images ({{ form.details.length }})</label>
                    <el-button type="success" size="mini" icon="el-icon-share" plain @click="shareImages">Share All</el-button>
                </div>
                <div class="grid grid-cols-3 gap-2">
                    <div v-for="img in form.details" :key="img.id" class="aspect-square rounded overflow-hidden border border-slate-200">
                        <img :src="img.file_path" class="w-full h-full object-cover">
                    </div>
                </div>
            </div>
         </div>

         <!-- Right Column: Details & Sizes -->
         <div class="md:col-span-2 space-y-6">
             <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
                 <h2 class="font-bold text-slate-800 text-lg mb-4">Product Details</h2>
                 <!-- ... existing ... -->
                 <div class="space-y-4">
                     <!-- Quick Input -->
                     <!-- Quick Input Hidden as requested -->
                     <div v-if="false">
                         <label class="flex justify-between items-center text-slate-400 font-bold text-xs uppercase mb-2 tracking-wider">
                             <span>Quick Input (Space separated)</span>
                             <span class="text-blue-600 cursor-pointer" @click="$message.info('Format: ' + (activeMapping.labels || 'Select type first'))">
                                 <i class="el-icon-info"></i> Help
                             </span>
                         </label>
                         <el-input v-model="quickInput" placeholder="e.g. 47 54 76 60" @input="parseQuickInput">
                             <template slot="prepend"><i class="el-icon-magic-stick"></i></template>
                         </el-input>
                          <p class="text-[10px] text-slate-400 mt-2" v-if="activeMapping">
                             Order: <span class="font-bold text-slate-600">{{ activeMapping.labels }}</span>
                         </p>
                     </div>
                     
                     <div class="grid grid-cols-2 gap-4">
                        <div>
                             <label class="block text-slate-400 font-bold text-xs uppercase mb-2 tracking-wider">Price</label>
                             <el-input-number v-model="form.price" :step="10000" class="w-full" controls-position="right"></el-input-number>
                        </div>
                        <div>
                             <label class="block text-slate-400 font-bold text-xs uppercase mb-2 tracking-wider">Description</label>
                             <el-input v-model="form.description" placeholder="Optional notes..."></el-input>
                        </div>
                     </div>
                 </div>
             </div>

             <!-- Dynamic Measurement Fields -->
              <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100" v-if="form.type">
                 <h2 class="font-bold text-slate-800 text-lg mb-4">Measurements</h2>
                 <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
                     <div v-for="(sizeName, idx) in allSizes" :key="sizeName" 
                          v-show="isSizeRelevant(idx)"
                          class="bg-slate-50 p-3 rounded-lg border border-slate-100">
                         <label class="block text-[10px] font-black uppercase text-slate-400 marginBottom-1">{{ sizeName }}</label>
                         <el-input v-model="form.sizes[idx].value" size="medium" placeholder="-">
                             <template slot="append">cm</template>
                         </el-input>
                     </div>
                 </div>
             </div>

             <!-- Action Buttons -->
             <div class="fixed bottom-0 md:relative left-0 right-0 p-4 md:p-0 bg-white md:bg-transparent border-t md:border-t-0 border-slate-200 z-40 flex justify-end gap-3">
                 <el-button @click="$router.push('/salev2/products')">Cancel</el-button>
                 <el-button type="primary" :loading="submitting" @click="submit" class="w-full md:w-auto px-8 shadow-lg shadow-blue-200">
                     Save Changes
                 </el-button>
             </div>
         </div>
     </div>
  </div>
</template>

<script>
export default {
    name: 'ProductEdit',
    data() {
        return {
            loading: false,
            submitting: false,
            form: {
                id: null,
                name: '',
                type: '',
                price: 0,
                description: '',
                image: null,
                details: [], // For display logic
                sizes: []
            },
            imagePreview: '',
            quickInput: '',
            productTypes: ['TROUSERS', 'JACKET', 'SHIRT', 'BLAZER', 'TIE', 'GILE', 'BELT', 'POLO SHIRT', 'HAT', 'SUIT'],
            allSizes: ['VAI', 'NGỰC', 'EO', 'DÀI ÁO', 'DÀI ÁO SAU', 'DÀI TAY', 'EO QUẦN', 'ĐÁY', 'ĐÙI', 'DÀI QUẦN', 'ỐNG', 'DƯ LAI'],
            mappings: {
                'BLAZER': { labels: 'Vai - Ngực - Dài trước - Dài tay', indices: [0, 1, 3, 5] },
                'JACKET': { labels: 'Vai - Ngực - Dài - Dài tay', indices: [0, 1, 3, 5] },
                'TROUSERS': { labels: 'Eo - Đáy - Đùi - Dài quần - Ống - Dư lai', indices: [6, 7, 8, 9, 10, 11] },
                'SHIRT': { labels: 'Vai - Ngực - Dài - Dài tay', indices: [0, 1, 3, 5] },
                'GILE': { labels: 'Vai - Ngực - Dài - Dài tay', indices: [0, 1, 3, 5] },
                'POLO SHIRT': { labels: 'Vai - Ngực - Dài - Dài tay', indices: [0, 1, 3, 5] },
                'SUIT': { labels: 'Vai - Ngực - Dài - Dài tay', indices: [0, 1, 3, 5] },
            }
        }
    },
    computed: {
        activeMapping() {
             return this.mappings[this.form.type] || null;
        }
    },
    mounted() {
        this.initializeSizes();
        this.fetchProduct();
    },
    methods: {
        initializeSizes() {
            this.form.sizes = this.allSizes.map(name => ({ name, value: '' }));
        },
        formatPrice(val) {
            if (!val) return '0';
            return new Intl.NumberFormat('vi-VN').format(val);
        },
        async fetchProduct() {
            this.loading = true;
            try {
                const id = this.$route.params.id;
                const res = await axios.get(`/salev2/api/product/${id}`);
                const product = res.data;
                
                this.form.id = product.id;
                this.form.name = product.name;
                this.form.type = product.type;
                this.form.price = product.price;
                this.form.description = product.description;
                this.form.details = product.images || []; // Populate details
                this.imagePreview = product.path_thumb;

                // Populate sizes
                if (product.sizes) {
                    product.sizes.forEach(s => {
                         const idx = this.allSizes.indexOf(s.name);
                         if (idx !== -1) this.form.sizes[idx].value = s.value;
                    });
                }
            } catch (err) {
                 this.$message.error('Failed to load product');
                 this.$router.push('/salev2/products');
            } finally {
                this.loading = false;
            }
        },
        async shareImages() {
             if (!navigator.share) {
                 alert('Browser does not support native sharing (or not HTTPS). Try Mobile Safari/Chrome.');
                 return;
             }
             
             const loading = this.$loading({ lock: true, text: 'Preparing images...', background: 'rgba(0,0,0,0.7)' });
             try {
                 const files = [];
                 
                 // Detail Images
                 if (this.form.details && this.form.details.length) {
                     for (let i = 0; i < this.form.details.length; i++) {
                         const img = this.form.details[i];
                         // FIX: Use file_path instead of path
                         const imgUrl = img.file_path || img.path; 
                         
                         if (imgUrl) {
                            try {
                                 const response = await fetch(imgUrl);
                                 if (!response.ok) throw new Error('Fetch failed');
                                 const blob = await response.blob();
                                 files.push(new File([blob], `Detail-${this.form.id}-${i+1}.jpg`, { type: 'image/jpeg' }));
                            } catch (err) {
                                console.warn('Failed to load detail image: ' + imgUrl);
                            }
                         }
                     }
                 }
                 
                 if (!files.length) {
                    alert('No valid images to share');
                    return;
                 }
                 
                 if (navigator.canShare && navigator.canShare({ files: files })) {
                     await navigator.share({
                         files: files
                     });
                 } else {
                     alert('Browser cannot share these files (limit or format issue).');
                 }
             } catch (e) {
                 console.error(e);
                 if (e.name !== 'AbortError') alert('Share failed: ' + e.message);
             } finally {
                 loading.close();
             }
        },
        handleFileChange(e) {
            const file = e.target.files[0];
            if (file) {
                this.form.image = file;
                this.imagePreview = URL.createObjectURL(file);
            }
        },
        isSizeRelevant(idx) {
            if (!this.form.type) return false;
            // Simple logic: Trousers/Shorts use bottom sizes (idx >= 6), others use top sizes (idx < 6)
            const bottomTypes = ['TROUSERS', 'SHORT'];
            const isBottom = bottomTypes.includes(this.form.type);
            return isBottom ? idx >= 6 : idx < 6;
        },
        parseQuickInput() {
             if (!this.quickInput) return;
             const nums = this.quickInput.trim().split(/\s+/);
             const config = this.activeMapping;
             
             if (config && config.indices) {
                 config.indices.forEach((sizeIdx, i) => {
                     if (nums[i]) {
                         this.form.sizes[sizeIdx].value = nums[i];
                     }
                 });
             }
        },
        async submit() {
            this.submitting = true;
            try {
                const formData = new FormData();
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

                await axios.post(`/salev2/api/product/${this.form.id}/update`, formData, {
                     headers: { 'Content-Type': 'multipart/form-data' }
                });

                this.$message.success('Product updated successfully');
                this.$router.push('/salev2/products');
            } catch (err) {
                console.error(err);
                this.$message.error('Failed to update product');
            } finally {
                this.submitting = false;
            }
        }
    }
}
</script>
