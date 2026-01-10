<template>
  <div class="px-4 md:px-6 py-6 h-[calc(100vh-80px)] overflow-hidden flex flex-col">
    <!-- Header -->
    <div class="flex justify-between items-center mb-4 flex-none">
        <div class="flex items-center gap-3">
            <el-button icon="el-icon-back" circle @click="$router.push('/salev2/customers')"></el-button>
            <h1 class="text-2xl font-bold text-slate-800">Create Order</h1>
        </div>
        <div class="flex gap-2">
            <el-button type="success" icon="el-icon-check" :loading="submitting" @click="submitOrder" class="px-6 shadow-lg shadow-green-100">Complete Order</el-button>
        </div>
    </div>

    <div class="flex flex-col lg:flex-row gap-6 h-full overflow-hidden">
        <!-- LEFT: Customer & Financials -->
        <div class="w-full lg:w-96 flex-none flex flex-col gap-4 overflow-y-auto pb-20 lg:pb-0">
            <!-- Customer Card -->
            <div class="bg-white p-4 rounded-xl border border-slate-100 shadow-sm">
                <h3 class="font-bold text-slate-700 mb-3 flex items-center gap-2"><i class="el-icon-user"></i> Customer</h3>
                <div v-if="customer" class="space-y-3">
                    <el-input v-model="customer.name" placeholder="Name" size="small">
                        <template slot="prepend">Name</template>
                    </el-input>
                    <el-input v-model="customer.phone" placeholder="Phone" size="small">
                        <template slot="prepend">Phone</template>
                    </el-input>
                    <el-input v-model="customer.address" placeholder="Address" type="textarea" :rows="2" size="small"></el-input>
                    
                    <!-- Note moved here -->
                    <div>
                        <label class="text-xs font-bold text-slate-500 uppercase">Note</label>
                        <el-input v-model="order.note" type="textarea" :rows="2" placeholder="Internal note..."></el-input>
                    </div>
                </div>
                <div v-else class="text-center py-4">
                    <el-spinner v-if="loadingCustomer"></el-spinner>
                    <p v-else class="text-slate-400 text-sm">No customer selected</p>
                </div>
            </div>

            <!-- Financials Card -->
            <div class="bg-white p-4 rounded-xl border border-slate-100 shadow-sm flex-1">
                <h3 class="font-bold text-slate-700 mb-3 flex items-center gap-2"><i class="el-icon-money"></i> Payment</h3>
                
                <div class="space-y-4">
                    <!-- Note -->

                     <div>
                        <label class="text-xs font-bold text-slate-500 uppercase">Public Description (In Receipt)</label>
                        <el-input v-model="order.desc" type="textarea" :rows="2" placeholder="Public note..."></el-input>
                    </div>

                    <div class="pt-4 border-t border-slate-100 space-y-2">
                        <div class="flex justify-between items-center text-sm">
                            <span class="text-slate-500">Subtotal ({{ selectedProducts.length }} items)</span>
                            <span class="font-bold">{{ formatPrice(subtotal) }}</span>
                        </div>
                        
                        <!-- Extra Fee -->
                        <div class="flex items-center justify-between gap-2">
                             <span class="text-slate-500 text-sm">Additional Fee</span>
                             <el-input-number v-model="order.additional_amount" size="mini" :step="10000" class="w-28"></el-input-number>
                        </div>

                         <!-- Discount/Freeship (Sale Amount) -->
                        <div class="flex items-center justify-between gap-2">
                             <span class="text-slate-500 text-sm">Discount / Sale</span>
                             <el-input-number v-model="order.freeship_amount" size="mini" :step="10000" class="w-28"></el-input-number>
                        </div>
                        
                        <!-- Deposit -->
                         <div class="flex items-center justify-between gap-2 pt-2 border-t border-dashed border-slate-100">
                             <span class="text-slate-700 text-sm font-medium">Deposit</span>
                             <el-input-number v-model="order.deposit_amount" size="mini" :step="50000" class="w-28"></el-input-number>
                        </div>

                        <!-- Paid Toggle -->
                        <div class="flex justify-between items-center py-2">
                            <span class="font-bold text-slate-700 text-sm">Paid in Full</span>
                            <el-switch v-model="order.is_paid"></el-switch>
                        </div>

                        <!-- Final Amount -->
                         <div class="flex justify-between items-center pt-3 border-t border-slate-900">
                            <span class="font-black text-slate-900 text-lg">TOTAL</span>
                            <span class="font-black text-blue-600 text-xl">{{ formatPrice(finalTotal) }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- RIGHT: Product Selection -->
        <div class="flex-1 bg-white rounded-xl border border-slate-100 shadow-sm flex flex-col overflow-hidden">
            <!-- Search Bar -->
            <div class="p-4 border-b border-slate-100 flex gap-2">
                <el-input v-model="productSearch" placeholder="Search product to add..." prefix-icon="el-icon-search" @input="debouncedSearch" class="flex-1" clearable></el-input>
                <el-button icon="el-icon-refresh" circle @click="searchProducts"></el-button>
            </div>
            
            <!-- Products Grid / Results -->
            <div class="flex-1 overflow-y-auto p-4 relative">
                <div v-if="loadingProducts" class="absolute inset-0 bg-white/80 z-10 flex items-center justify-center">
                    <el-spinner></el-spinner>
                </div>
                
                <div v-if="searchResults.length" class="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
                    <div v-for="prod in searchResults" :key="prod.id" 
                         class="border rounded-lg p-2 cursor-pointer transition-all hover:shadow-md h-full flex flex-col"
                         :class="isSelected(prod) ? 'border-blue-500 bg-blue-50 ring-1 ring-blue-500' : 'border-slate-100'"
                         @click="toggleProduct(prod)">
                        <div class="aspect-[3/4] bg-slate-100 rounded mb-2 overflow-hidden">
                             <img v-if="prod.path_thumb" :src="prod.path_thumb" class="w-full h-full object-contain">
                        </div>
                        <div class="text-center mt-auto">
                            <h5 class="font-bold text-xs truncate">{{ prod.name }}</h5>
                            <p class="text-xs text-blue-600 font-bold">{{ formatPrice(prod.price) }}</p>
                            <span class="text-[10px] text-slate-400 bg-slate-100 px-1 rounded">{{ prod.status }}</span>
                        </div>
                    </div>
                </div>
                 <div v-else-if="productSearch" class="text-center py-10 text-slate-400">
                    No products found for "{{ productSearch }}"
                </div>
                 <div v-else class="text-center py-10 text-slate-400">
                    <i class="el-icon-search text-4xl mb-2"></i>
                    <p>Type to search products</p>
                </div>
            </div>

            <!-- Selected Items Bar -->
            <div class="p-3 bg-slate-50 border-t border-slate-200 h-24 overflow-x-auto whitespace-nowrap flex items-center gap-2">
                <span class="text-xs font-bold text-slate-400 uppercase mr-2 flex-none">Selected ({{ selectedProducts.length }})</span>
                <div v-if="!selectedProducts.length" class="text-slate-400 text-sm italic">No items selected</div>
                <div v-for="item in selectedProducts" :key="item.id" class="inline-flex items-center gap-2 bg-white border border-slate-200 rounded-lg p-1.5 shadow-sm pr-3">
                     <img v-if="item.path_thumb" :src="item.path_thumb" class="w-8 h-8 object-cover rounded">
                     <div class="flex flex-col">
                         <span class="text-xs font-bold">{{ item.name }}</span>
                         <el-input-number v-model="item.price" size="mini" :controls="false" class="w-20"></el-input-number>
                     </div>
                     <i class="el-icon-close text-slate-400 hover:text-red-500 cursor-pointer p-1" @click.stop="toggleProduct(item)"></i>
                </div>
            </div>
        </div>
    </div>
  </div>
</template>

<script>
export default {
    name: 'OrderCreate',
    data() {
        return {
            customer: { name: '', phone: '', address: '' },
            order: { note: '', desc: '', additional_amount: 0, freeship_amount: 0, deposit_amount: 0, is_paid: false },
            
            productSearch: '',
            searchResults: [],
            selectedProducts: [],
            
            loadingCustomer: false,
            loadingProducts: false,
            submitting: false
        }
    },
    computed: {
        subtotal() {
            return this.selectedProducts.reduce((sum, p) => sum + Number(p.price), 0);
        },
        finalTotal() {
            let total = this.subtotal + this.order.additional_amount - this.order.freeship_amount;
            return total > 0 ? total : 0;
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
    methods: {
        formatPrice(val) {
            return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND', maximumFractionDigits: 0 }).format(val);
        },
        async fetchConsumer(id) {
            this.loadingCustomer = true;
            try {
                 const res = await axios.get(`/salev2/api/consumer/${id}`);
                 this.customer = {
                     name: res.data.name,
                     phone: res.data.phone,
                     address: res.data.address
                 };
            } catch(e) {
                this.$message.error('Failed to load customer');
            } finally {
                this.loadingCustomer = false;
            }
        },
        debouncedSearch: _.debounce(function() { this.searchProducts(); }, 500),
        async searchProducts() {
            if(!this.productSearch) return;
            this.loadingProducts = true;
            try {
                const res = await axios.get('/salev2/api/products', {
                    params: { search: this.productSearch, status: 'AVAILABLE' }
                });
                this.searchResults = res.data.data;
            } catch(e) {
                console.error(e);
            } finally {
                this.loadingProducts = false;
            }
        },
        toggleProduct(prod) {
            const idx = this.selectedProducts.findIndex(p => p.id === prod.id);
            if(idx > -1) {
                this.selectedProducts.splice(idx, 1);
            } else {
                this.selectedProducts.push(prod);
            }
        },
        isSelected(prod) {
            return this.selectedProducts.some(p => p.id === prod.id);
        },
        async submitOrder() {
            if(!this.customer.name || !this.customer.phone) return this.$message.error('Customer info required');
            if(!this.selectedProducts.length) return this.$message.error('No products selected');
            
            this.submitting = true;
            try {
                const payload = {
                    ...this.customer, // Flatten customer fields
                    consumer_id: this.$route.query.consumer_id, // Optional
                    ...this.order,
                    total_amount: this.finalTotal, // Legacy expects total_amount as Final? No, let's allow backend to calculate or send calc.
                    // Actually storeOrder expects: products (array of IDs)
                    products: this.selectedProducts.map(p => p.id),
                    products_info: this.selectedProducts.map(p => ({ id: p.id, price: p.price }))
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
