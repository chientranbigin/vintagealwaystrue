<template>
  <div class="px-4 md:px-6 py-8">
    <div class="flex justify-between items-center mb-6 mobile-hide">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 tracking-tight">Customer Network</h1>
        <p class="text-slate-500 mt-1">Manage profiles and purchase history</p>
      </div>
      <div class="flex gap-3">
          <el-button type="primary" icon="el-icon-plus">Add Customer</el-button>
      </div>
    </div>

    <!-- Mobile Header -->
    <div class="md:hidden flex justify-between items-center mb-6">
        <h1 class="text-2xl font-bold text-slate-800">Customers</h1>
        <el-button size="small" type="primary" icon="el-icon-plus" circle></el-button>
    </div>

    <!-- Search & Filters -->
    <div class="bg-white p-4 rounded-2xl shadow-sm border border-slate-100 mb-6">
        <el-input v-model="search" placeholder="Search by Name, Phone, Address..." prefix-icon="el-icon-search" @input="debouncedFetch" clearable class="w-full md:w-96">
        </el-input>
    </div>

    <!-- Customer List (Grid on Desktop, Card on Mobile) -->
    <div v-loading="loading">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            <div v-for="customer in consumers" :key="customer.id" class="bg-white rounded-xl p-5 border border-slate-100 shadow-sm hover:shadow-md transition-shadow relative group">
                <div class="flex justify-between items-start mb-3">
                    <div class="flex items-center gap-3">
                        <div class="w-10 h-10 rounded-full bg-slate-100 flex items-center justify-center text-slate-500 font-bold">
                            {{ customer.name.charAt(0).toUpperCase() }}
                        </div>
                        <div>
                            <h4 class="font-bold text-slate-800 text-sm">{{ customer.name }}</h4>
                            <p class="text-xs text-slate-500">{{ customer.phone }}</p>
                        </div>
                    </div>
                    <div class="bg-blue-50 text-blue-600 px-2 py-1 rounded text-[10px] font-bold uppercase tracking-wider">
                        {{ customer.orders_count }} Orders
                    </div>
                </div>
                
                <div class="mb-4">
                    <p class="text-xs text-slate-500 line-clamp-2 min-h-[2.5em] flex items-center">
                        <i class="el-icon-location-outline mr-1"></i> 
                        {{ customer.address || 'No address provided' }}
                    </p>
                </div>

                <div class="flex gap-2">
                    <el-button type="primary" plain size="small" class="flex-1" icon="el-icon-shopping-cart-full" @click="createOrder(customer)">Create Order</el-button>
                    <el-button size="small" icon="el-icon-time" circle @click="$message.info('History View Coming Soon')"></el-button>
                </div>
            </div>
        </div>

           <div v-if="!consumers.length && !loading" class="text-center py-20 bg-white rounded-3xl border border-dashed border-slate-200">
          <i class="el-icon-user text-4xl text-slate-200 mb-4 block"></i>
          <p class="text-slate-400 font-medium">No customers found</p>
      </div>

      <div class="mt-8 flex justify-center">
        <el-pagination background layout="prev, pager, next" :current-page.sync="page" :total="total" :page-size="20" @current-change="fetchConsumers"></el-pagination>
      </div>
    </div>
  </div>
</template>

<script>
export default {
    name: 'CustomerManager',
    data() {
        return {
            consumers: [],
            loading: false,
            search: '',
            page: 1,
            total: 0
        }
    },
    mounted() {
        this.fetchConsumers();
    },
    methods: {
        async fetchConsumers(page = 1) {
            this.page = page;
            this.loading = true;
            try {
                const res = await axios.get('/salev2/api/consumers', {
                    params: { page, search: this.search }
                });
                this.consumers = res.data.data;
                this.total = res.data.total;
            } catch (err) {
                console.error(err);
                this.$message.error('Failed to load consumers');
            } finally {
                this.loading = false;
            }
        },
        debouncedFetch: _.debounce(function() { this.fetchConsumers(1); }, 500),
        createOrder(customer) {
            this.$router.push({ 
                name: 'order-create', 
                query: { consumer_id: customer.id } 
            });
        }
    }
}
</script>
