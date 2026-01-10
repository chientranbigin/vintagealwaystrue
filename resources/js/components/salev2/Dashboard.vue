<template>
  <div class="px-6 py-8">
    <div class="flex justify-between items-center mb-8 mobile-hide">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 tracking-tight">Sales Dashboard</h1>
        <p class="text-slate-500 mt-1">Real-time performance and inventory metrics</p>
      </div>
      <el-button type="primary" icon="el-icon-refresh" @click="fetchStats" :loading="loading">Refresh Data</el-button>
    </div>

    <!-- Legacy Style Stat Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div v-for="card in statCards" :key="card.title" 
             class="relative overflow-hidden bg-white rounded-2xl p-6 shadow-sm border border-slate-100 hover:shadow-md transition-shadow">
          <div :class="['absolute top-0 right-0 w-24 h-24 -mr-8 -mt-8 rounded-full opacity-10', card.bgClass]"></div>
          <div class="relative z-10 flex items-center justify-between">
            <div>
              <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">{{ card.title }}</p>
              <h3 class="text-2xl font-black text-slate-900">{{ card.value }}</h3>
            </div>
            <div :class="['w-12 h-12 rounded-xl flex items-center justify-center text-xl text-white shadow-lg', card.bgClass]">
              <i :class="card.icon"></i>
            </div>
          </div>
        </div>
    </div>

    <!-- Legacy Comparison Tables -->
    <div class="grid grid-cols-1 xl:grid-cols-3 gap-8 mb-8">
      <!-- All Time Summary -->
      <div class="bg-white rounded-2xl shadow-sm border border-slate-100 p-6">
        <div class="flex justify-between items-center mb-6">
            <h3 class="text-lg font-black text-slate-900">All Time</h3>
            <div class="text-right">
                <p class="text-[10px] font-bold text-slate-400 uppercase leading-none">Total Revenue</p>
                <p class="text-blue-600 font-black">{{ formatPrice(data?.all.total_revenue) }}</p>
            </div>
        </div>
        <el-table :data="allTypeStats" size="mini" style="width: 100%" class="legacy-dash-table">
          <el-table-column prop="type" label="Type"></el-table-column>
          <el-table-column label="Sold" width="70" align="center">
            <template slot-scope="scope">{{ getCount(data?.all.type_stats, scope.row.type, 'SOLD') }}</template>
          </el-table-column>
          <el-table-column label="Avail" width="70" align="center">
            <template slot-scope="scope">{{ getCount(data?.all.type_stats, scope.row.type, 'AVAILABLE') }}</template>
          </el-table-column>
        </el-table>
      </div>

      <!-- Current Month -->
      <div class="bg-white rounded-2xl shadow-sm border border-slate-100 p-6">
        <div class="flex justify-between items-center mb-6">
            <h3 class="text-lg font-black text-slate-900">Current Month</h3>
            <div class="text-right">
                <p class="text-[10px] font-bold text-slate-400 uppercase leading-none">Revenue</p>
                <p class="text-emerald-600 font-black">{{ formatPrice(data?.current_month.total_revenue) }}</p>
            </div>
        </div>
        <el-table :data="allTypeStats" size="mini" style="width: 100%" class="legacy-dash-table">
          <el-table-column prop="type" label="Type"></el-table-column>
          <el-table-column label="Sold" width="80" align="right">
            <template slot-scope="scope">{{ getCount(data?.current_month.type_stats, scope.row.type) }}</template>
          </el-table-column>
        </el-table>
      </div>

      <!-- Previous Month -->
      <div class="bg-white rounded-2xl shadow-sm border border-slate-100 p-6">
        <div class="flex justify-between items-center mb-6">
            <h3 class="text-lg font-black text-slate-900">Previous Month</h3>
            <div class="text-right">
                <p class="text-[10px] font-bold text-slate-400 uppercase leading-none">Revenue</p>
                <p class="text-slate-600 font-black">{{ formatPrice(data?.prev_month.total_revenue) }}</p>
            </div>
        </div>
        <el-table :data="allTypeStats" size="mini" style="width: 100%" class="legacy-dash-table">
          <el-table-column prop="type" label="Type"></el-table-column>
          <el-table-column label="Sold" width="80" align="right">
            <template slot-scope="scope">{{ getCount(data?.prev_month.type_stats, scope.row.type) }}</template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <!-- Monthly Trend Table -->
    <div class="bg-white rounded-2xl shadow-sm border border-slate-100 p-8">
        <h3 class="text-xl font-black text-slate-900 mb-6">Monthly Performance</h3>
        <el-table :data="data?.monthly_sales" style="width: 100%">
            <el-table-column prop="month" label="Month" width="150">
                <template slot-scope="scope">
                    <span class="font-bold text-slate-700">{{ scope.row.month }}</span>
                </template>
            </el-table-column>
            <el-table-column prop="number" label="Orders Count" width="150" align="center"></el-table-column>
            <el-table-column prop="final_amount" label="Total Revenue">
                <template slot-scope="scope">
                    <span class="text-blue-600 font-black">{{ formatPrice(scope.row.final_amount) }}</span>
                </template>
            </el-table-column>
            <el-table-column label="Trend" width="200">
                <template slot-scope="scope">
                    <div class="w-full bg-slate-100 h-2 rounded-full overflow-hidden">
                        <div class="bg-blue-500 h-full" :style="{ width: getTrendWidth(scope.row.final_amount) + '%' }"></div>
                    </div>
                </template>
            </el-table-column>
        </el-table>
    </div>
  </div>
</template>

<script>
// Using global axios
export default {
  name: 'Dashboard',
  data() {
    return {
      data: null,
      loading: false,
      productTypes: ['TROUSERS', 'JACKET', 'SHIRT', 'BLAZER', 'TIE', 'GILE', 'BELT', 'POLO SHIRT', 'HAT', 'SUIT'],
      statCards: [
        { title: 'New Orders Rev', value: '0đ', icon: 'el-icon-message-solid', bgClass: 'bg-orange-500' },
        { title: 'New Orders Count', value: '0', icon: 'el-icon-bell', bgClass: 'bg-orange-400' },
        { title: 'Monthly Revenue', value: '0đ', icon: 'el-icon-money', bgClass: 'bg-emerald-500' },
        { title: 'Monthly Count', value: '0', icon: 'el-icon-shopping-cart-full', bgClass: 'bg-blue-500' }
      ]
    }
  },
  computed: {
    allTypeStats() {
      return this.productTypes.map(t => ({ type: t }));
    }
  },
  mounted() {
    this.fetchStats();
  },
  methods: {
    formatPrice(val) {
      if (!val) return '0đ';
      return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND', maximumFractionDigits: 0 }).format(val);
    },
    async fetchStats() {
      this.loading = true;
      try {
        const response = await axios.get('/salev2/api/dashboard-stats');
        this.data = response.data;
        
        // Update Stat Cards (Legacy Sync)
        this.statCards[0].value = this.formatPrice(this.data.new_orders_revenue);
        this.statCards[1].value = this.data.new_orders_count;
        this.statCards[2].value = this.formatPrice(this.data.month_orders_revenue);
        this.statCards[3].value = this.data.month_orders_count;
      } catch (error) {
        console.error('Failed to fetch stats:', error);
      } finally {
        this.loading = false;
      }
    },
    getCount(stats, type, status = 'SOLD') {
        if (!stats) return 0;
        const found = stats.find(s => s.type === type && (!s.status || s.status === status));
        return found ? found.count : 0;
    },
    getTrendWidth(amount) {
        if (!this.data || !this.data.monthly_sales.length) return 0;
        const max = Math.max(...this.data.monthly_sales.map(m => m.final_amount));
        return (amount / max) * 100;
    }
  }
}
</script>

<style>
.legacy-dash-table .el-table__header {
    background-color: #f8fafc !important;
}
.legacy-dash-table .el-table__row {
    font-size: 11px;
}
.legacy-dash-table .el-table__header th {
    font-size: 10px;
    font-weight: 900;
    text-transform: uppercase;
    color: #94a3b8;
    padding: 8px 0;
}
</style>
