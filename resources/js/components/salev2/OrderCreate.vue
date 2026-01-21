<template>
  <div class="px-4 md:px-6 py-8">
     <div class="flex justify-between items-center mb-6">
       <div>
         <h1 class="text-3xl font-bold text-gray-900 tracking-tight">Create Order</h1>
         <p class="text-slate-500 mt-1">Fill customer details and select products to create a new order</p>
       </div>
       <el-button @click="$router.push('/salev2/orders')" icon="el-icon-back" circle></el-button>
     </div>

     <div v-loading="loading">
         <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
             
             <!-- Left: Customer & Status -->
             <div class="space-y-6">
                 <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
                     <h2 class="font-bold text-slate-800 text-lg mb-4">Customer Info</h2>
                     <div class="space-y-4">
                         <div>
                             <label class="block text-slate-400 font-bold text-xs uppercase mb-1">Name</label>
                             <el-input v-model="order.name"></el-input>
                         </div>
                         <div>
                             <label class="block text-slate-400 font-bold text-xs uppercase mb-1">Phone</label>
                             <el-input v-model="order.phone"></el-input>
                         </div>
                         <div>
                             <label class="block text-slate-400 font-bold text-xs uppercase mb-1">Address</label>
                             <el-input type="textarea" v-model="order.address" :rows="2"></el-input>
                         </div>
                         <div>
                             <label class="block text-slate-400 font-bold text-xs uppercase mb-1">Note</label>
                             <el-input type="textarea" v-model="order.note" :rows="2" placeholder="Internal note..."></el-input>
                         </div>
                         <div>
                             <label class="block text-slate-400 font-bold text-xs uppercase mb-1">Public Description</label>
                             <el-input type="textarea" v-model="order.desc" :rows="2" placeholder="Public note..."></el-input>
                         </div>
                     </div>
                 </div>

                 <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
                     <h2 class="font-bold text-slate-800 text-lg mb-4">Order Settings</h2>
                     <div class="space-y-4">
                         <div>
                             <label class="block text-slate-400 font-bold text-xs uppercase mb-1">Status</label>
                             <el-select v-model="order.status" class="w-full">
                                 <el-option value="NEW" label="NEW"></el-option>
                                 <el-option value="ON HOLD" label="ON HOLD"></el-option>
                                 <el-option value="A SHIP NOW" label="A SHIP NOW"></el-option>
                                 <el-option value="COMPLETED" label="COMPLETED"></el-option>
                                 <el-option value="CANCELED" label="CANCELED"></el-option>
                             </el-select>
                         </div>
                         <div class="flex items-center gap-2">
                             <el-switch v-model="order.is_freeship" active-color="#13ce66"></el-switch>
                             <span class="text-sm font-medium">Free Shipping (20k)</span>
                         </div>
                         <div class="flex items-center gap-2">
                             <el-switch v-model="order.is_paid" active-color="#13ce66"></el-switch>
                             <span class="text-sm font-medium">Paid in Full</span>
                         </div>
                     </div>
                 </div>
             </div>

             <!-- Middle: Products -->
             <div class="lg:col-span-2 space-y-6">
                 <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
                     <div class="flex justify-between items-center mb-4">
                        <h2 class="font-bold text-slate-800 text-lg">Products</h2>
                        <el-button type="text" icon="el-icon-plus" @click="showProductSelector = true">Add Product</el-button>
                     </div>
                     
                     <!-- Selected Products List -->
                     <div class="space-y-3">
                         <div v-for="(p, index) in order.products" :key="p.id" class="flex gap-4 p-3 bg-slate-50 rounded-xl border border-slate-100 relative group">
                             <img :src="p.path_thumb" class="w-16 h-20 object-contain bg-white rounded border border-slate-200">
                             <div class="flex-1">
                                 <div class="flex justify-between items-start">
                                     <div>
                                         <h4 class="font-bold text-slate-800">{{ p.name }}</h4>
                                         <span class="text-[10px] text-slate-500 font-bold uppercase">{{ p.type }}</span>
                                     </div>
                                     <el-input-number v-model="p.price" :step="10000" size="mini" controls-position="right" class="w-28"></el-input-number>
                                 </div>
                                 <div class="mt-2 text-xs text-slate-500">
                                     <span v-for="s in p.sizes" :key="s.name" class="mr-2">{{ s.name }}: {{ Math.floor(s.value) }}</span>
                                 </div>
                             </div>
                             <button @click="removeProduct(index)" class="absolute -top-2 -right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center shadow-sm opacity-0 group-hover:opacity-100 transition-opacity">
                                 <i class="el-icon-close font-bold text-xs"></i>
                             </button>
                         </div>
                         <div v-if="!order.products.length" class="text-center py-8 text-slate-400 italic">No products selected</div>
                     </div>

                     <!-- Financials -->
                      <div class="mt-6 pt-6 border-t border-slate-100 space-y-3">
                         <div class="flex justify-between text-sm">
                             <span class="text-slate-500">Subtotal</span>
                             <span class="font-bold">{{ formatPrice(subtotal) }}</span>
                         </div>
                         <div class="flex justify-between text-sm items-center">
                             <span class="text-slate-500">Additional Fee</span>
                             <el-input-number v-model="order.additional_amount" :step="10000" size="mini" controls-position="right" class="w-32"></el-input-number>
                         </div>
                         <div class="flex justify-between text-sm items-center">
                              <span class="text-slate-500">Freeship/Discount</span>
                              <el-input-number v-model="order.freeship_amount" :step="10000" size="mini" controls-position="right" class="w-32"></el-input-number>
                          </div>
                         <div class="flex justify-between text-sm items-center">
                             <span class="text-slate-500">Deposit</span>
                             <el-input-number v-model="order.deposit_amount" :step="10000" size="mini" controls-position="right" class="w-32"></el-input-number>
                         </div>
                         <!-- Shipping Fee hidden as requested -->
                         <div class="flex justify-between text-lg font-black text-slate-900 pt-3 border-t border-dashed border-slate-200">
                             <span>Total</span>
                             <span>{{ formatPrice(totalAmount) }}</span>
                         </div>
                         <div class="flex justify-between text-base font-bold text-blue-600">
                             <span>Viettel Amount</span>
                             <span>{{ formatPrice(viettelAmount) }}</span>
                         </div>
                      </div>
                 </div>

                 <!-- Action -->
                  <div class="flex justify-end gap-3 pt-4">
                     <el-button @click="$router.push('/salev2/orders')">Cancel</el-button>
                     <el-button type="primary" :loading="submitting" @click="submitOrder" class="px-8 shadow-lg shadow-green-200">Complete Order</el-button>
                 </div>
             </div>
         </div>
     </div>

     <!-- Product Selector Dialog -->
    <el-dialog title="Add Products" :visible.sync="showProductSelector" width="90%" custom-class="mobile-dialog">
        <div class="mb-4">
            <el-input v-model="productSearch" placeholder="Search by name..." prefix-icon="el-icon-search" @input="debouncedSearch"></el-input>
        </div>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-3 max-h-[60vh] overflow-y-auto">
            <div v-for="p in availableProducts" :key="p.id" 
                 class="border rounded-lg p-2 cursor-pointer hover:border-blue-500 transition-colors relative"
                 @click="toggleProductSelect(p)">
                <img :src="p.path_thumb" class="w-full aspect-[3/4] object-contain bg-slate-50 mb-2">
                <p class="font-bold text-xs text-center">{{ p.name }}</p>
                <div v-if="isProductSelected(p.id)" class="absolute inset-0 bg-blue-500/20 border-2 border-blue-500 rounded-lg flex items-center justify-center">
                    <i class="el-icon-check text-white bg-blue-500 rounded-full p-1"></i>
                </div>
            </div>
        </div>
        <span slot="footer">
            <el-button @click="showProductSelector = false">Done</el-button>
        </span>
    </el-dialog>
  </div>
</template>

<script>
export default {
    name: 'OrderCreate',
    data() {
        return {
            loading: false,
            submitting: false,
            order: {
                name: '',
                phone: '',
                address: '',
                note: '',
                desc: '',
                status: 'NEW',
                is_freeship: false,
                is_paid: false,
                additional_amount: 0,
                freeship_amount: 0,
                deposit_amount: 0,
                products: []
            },
            showProductSelector: false,
            productSearch: '',
            availableProducts: []
        }
    },
    computed: {
        subtotal() {
            return this.order.products.reduce((sum, p) => sum + Number(p.price), 0);
        },
        totalAmount() {
            let total = this.subtotal + (this.order.additional_amount || 0) - (this.order.freeship_amount || 0);
            return total > 0 ? total : 0;
        },
        viettelAmount() {
            if (this.order.is_paid) return 0;
            return this.totalAmount - (this.order.deposit_amount || 0);
        }
    },
    watch: {
        '$route.query.consumer_id': {
            immediate: true,
            handler(id) {
                if(id) this.fetchConsumer(id);
            }
        }
    },
    mounted() {
        this.fetchAvailableProducts();
    },
    methods: {
        formatPrice(val) {
            return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND', maximumFractionDigits: 0 }).format(val);
        },
        async fetchConsumer(id) {
            this.loading = true;
            try {
                 const res = await axios.get(`/salev2/api/consumer/${id}`);
                 this.order.name = res.data.name;
                 this.order.phone = res.data.phone;
                 this.order.address = res.data.address;
            } catch(e) {
                this.$message.error('Failed to load customer');
            } finally {
                this.loading = false;
            }
        },
        async fetchAvailableProducts() {
            try {
                const res = await axios.get('/salev2/api/products', { params: { status: 'AVAILABLE', search: this.productSearch, limit: 100 } });
                this.availableProducts = res.data.data;
            } catch(e) {}
        },
        debouncedSearch: _.debounce(function() { this.fetchAvailableProducts(); }, 500),
        isProductSelected(id) {
            return this.order.products.some(p => p.id === id);
        },
        toggleProductSelect(product) {
            if(this.isProductSelected(product.id)) {
                const idx = this.order.products.findIndex(p => p.id === product.id);
                if(idx > -1) this.order.products.splice(idx, 1);
            } else {
                this.order.products.push({...product});
            }
        },
        removeProduct(index) {
            this.order.products.splice(index, 1);
        },
        async submitOrder() {
            if(!this.order.name || !this.order.phone) return this.$message.error('Customer info required');
            if(!this.order.products.length) return this.$message.error('No products selected');
            
            this.submitting = true;
            try {
                const payload = {
                    ...this.order,
                    consumer_id: this.$route.query.consumer_id,
                    products: this.order.products.map(p => p.id),
                    products_info: this.order.products.map(p => ({ id: p.id, price: p.price })),
                    total_amount: this.totalAmount,
                    final_amount: this.totalAmount
                };
                
                await axios.post('/salev2/api/order', payload);
                this.$message.success('Order Created!');
                this.$router.push('/salev2/orders');
            } catch(e) {
                this.$message.error('Failed to create order');
            } finally {
                this.submitting = false;
            }
        }
    }
}
</script>
