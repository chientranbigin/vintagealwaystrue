<template>
  <div class="px-4 md:px-6 py-8">
    <div class="flex flex-col md:flex-row justify-between items-center mb-6 gap-4">
      <div class="mobile-hide">
        <h1 class="text-3xl font-bold text-gray-900 tracking-tight">Order Center</h1>
        <p class="text-slate-500 mt-1">Manage sales, shipping, and receipt generation</p>
      </div>
      <div class="flex gap-2 w-full md:w-auto">
        <el-input v-model="search" placeholder="Phone, Name..." prefix-icon="el-icon-search" class="w-full md:w-64" @input="debouncedFetch"></el-input>
        <el-button type="danger" plain icon="el-icon-money" @click="$router.push('/salev2/expenses')">Costs</el-button>
        <el-button type="primary" icon="el-icon-plus" @click="$router.push('/salev2/order/create')">New</el-button>
        <el-button icon="el-icon-refresh" @click="fetchOrders(1)"></el-button>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <div class="mb-6 flex justify-between items-center">
      <el-tabs v-model="activeTab" @tab-click="handleTabClick" type="card" class="premium-tabs">
        <el-tab-pane label="New Orders" name="NEW"></el-tab-pane>
        <el-tab-pane label="All Orders" name="ALL"></el-tab-pane>
      </el-tabs>
      
      <div class="flex gap-2">
          <el-radio-group v-model="viewMode" size="small">
            <el-radio-button label="grid"><i class="el-icon-menu"></i></el-radio-button>
            <el-radio-button label="table"><i class="el-icon-tickets"></i></el-radio-button>
          </el-radio-group>
      </div>
    </div>

    <!-- Grid View Mode -->
    <div v-if="viewMode === 'grid'" v-loading="loading">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-6 gap-4">
        <div v-for="order in orders" :key="order.id" 
             class="bg-white rounded-2xl border border-slate-100 shadow-sm hover:shadow-xl transition-all duration-300 flex flex-col overflow-hidden w-full md:max-w-[280px] mx-auto h-full">
          
          <!-- Card Header: Consumer Info & Actions -->
          <div class="p-4 border-b border-slate-50 bg-slate-50/30">
            <div class="flex justify-between items-start mb-2">
                <div>
                    <h5 class="font-black text-slate-900 leading-none mb-1 uppercase tracking-tight">{{ order.name }}</h5>
                    <div class="flex items-center gap-2">
                        <p class="text-blue-600 text-xs font-bold">{{ order.phone }}</p>
                        <span class="text-[10px] text-slate-400 font-medium italic">• {{ order.created_at_human }}</span>
                    </div>
                </div>
                <div class="flex items-center gap-1">
                    <el-tag :type="getStatusType(order.status)" size="mini" effect="dark" class="rounded-pill uppercase font-bold text-[8px]">
                        {{ order.status }}
                    </el-tag>
                    <el-dropdown trigger="click">
                        <span class="cursor-pointer p-1 rounded hover:bg-slate-100 transition-colors">
                            <i class="el-icon-more text-slate-400"></i>
                        </span>
                        <el-dropdown-menu slot="dropdown">
                            <el-dropdown-item icon="el-icon-camera" @click.native="captureOrder(order)">Capture Receipt</el-dropdown-item>
                            <el-dropdown-item icon="el-icon-edit" @click.native="$router.push({name: 'order-edit', params: {id: order.id}})">Edit Order</el-dropdown-item>
                            <el-dropdown-item icon="el-icon-user" @click.native="viewCustomerHistory(order)">Customer History</el-dropdown-item>
                            <el-divider class="my-1"></el-divider>
                            <el-dropdown-item v-if="order.status !== 'A SHIP NOW'" icon="el-icon-truck" @click.native="updateStatus(order, 'A SHIP NOW')" class="text-primary font-bold">Ready to Ship</el-dropdown-item>
                            <el-dropdown-item v-if="order.status !== 'ON HOLD'" icon="el-icon-video-pause" @click.native="updateStatus(order, 'ON HOLD')" class="text-warning font-bold">Mark On Hold</el-dropdown-item>
                            <el-dropdown-item v-if="order.status !== 'COMPLETED'" icon="el-icon-check" @click.native="updateStatus(order, 'COMPLETED')" class="text-success font-bold">Complete</el-dropdown-item>
                            <el-dropdown-item icon="el-icon-delete" @click.native="confirmDelete(order.id)" class="text-danger">Delete</el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                </div>
            </div>
            <p class="text-[12px] text-slate-500 break-words leading-snug"><i class="el-icon-location-outline mr-1"></i>{{ order.address }}</p>
            <div v-if="order.note" class="mt-2 text-[10px] bg-red-50 text-red-600 p-2 rounded-lg font-bold">
                <i class="el-icon-warning mr-1"></i>{{ order.note }}
            </div>
          </div>

          <!-- Card Body: Summary & Products -->
          <div class="p-4 flex-1 space-y-4">
            <!-- Product Preview -->
            <div class="grid grid-cols-2 gap-2 pb-2">
                <img v-for="p in order.products" :key="p.id" :src="p.path_thumb" 
                     class="w-full h-auto object-contain bg-slate-100 rounded shadow-sm border border-slate-100">
            </div>

            <!-- Summary Block (Matches Legacy Layout) -->
            <div class="bg-gray-50 rounded-xl p-3 text-[10px] space-y-1">
                <div v-for="p in order.products" :key="p.id" class="text-slate-600 block mb-1">
                   <div class="font-bold flex justify-between">
                       <span>{{ p.type }} {{ p.name }}</span>
                       <span class="text-slate-400">{{ formatPrice(p.price) }}</span>
                   </div>
                   <div class="flex flex-wrap gap-1 mt-0.5">
                       <span v-for="size in formatSizes(p)" :key="size" 
                             class="text-[8px] bg-slate-200 text-slate-600 px-1 py-0.5 rounded leading-none uppercase">
                         {{ size }}
                       </span>
                   </div>
                </div>
                
                <div v-if="order.additional_amount > 0" class="text-blue-600 font-bold border-t border-slate-100 pt-2 flex justify-between">
                    <span>Phí cộng thêm:</span>
                    <span>+{{ formatPrice(order.additional_amount) }}</span>
                </div>
                <div v-if="order.is_paid || order.deposit_amount > 0" class="text-emerald-600 font-bold flex justify-between">
                    <span>{{ order.is_paid ? 'Đã thanh toán:' : 'Cọc:' }}</span>
                    <span>-{{ formatPrice(order.is_paid ? order.total_amount : order.deposit_amount) }}</span>
                </div>
            </div>
          </div>

          <!-- Card Footer: Financials -->
          <div class="p-4 pt-0">
            <div class="flex justify-between items-end">
                <div>
                    <span class="text-[9px] text-slate-400 font-black uppercase block leading-none mb-1 ">Viettel Amount</span>
                    <div class="flex items-baseline gap-1">
                        <span class="text-xl font-black text-slate-900">{{ formatPrice(order.viettel_amount) }}</span>
                        <span class="text-[10px] font-bold text-slate-400">+ phí ship</span>
                    </div>
                </div>
                <div v-if="order.is_old" class="animate-pulse bg-red-500 text-white text-[8px] font-black px-2 py-0.5 rounded-pill">DELAY</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Table View Mode (Legacy Sync) -->
    <div v-else class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden" v-loading="loading">
      <el-table :data="orders" style="width: 100%">
        <el-table-column label="Consumer" width="200">
          <template slot-scope="scope">
            <div class="flex flex-col text-xs">
              <span class="font-bold text-slate-800">{{ scope.row.name }}</span>
              <span class="text-slate-500">{{ scope.row.phone }}</span>
              <span class="text-[10px] text-slate-400 italic">{{ scope.row.created_at_human }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Summary" min-width="250">
          <template slot-scope="scope">
            <div class="text-[10px] space-y-0.5 text-slate-500 italic">
                <span v-for="p in scope.row.products" :key="p.id" class="block">
                    {{ p.type }} {{ p.name }}: {{ formatPrice(p.price) }}
                </span>
                <span class="font-bold text-blue-600 block pt-1 border-t border-slate-50">
                    Total: {{ formatPrice(scope.row.is_paid ? 0 : scope.row.total_amount - scope.row.deposit_amount) }} 
                    {{ scope.row.is_paid ? '(PAID)' : '+ phí ship' }}
                </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Additionals" width="120" align="center">
            <template slot-scope="scope">
                <p class="text-[10px]"><span class="text-blue-600">Phí thêm:</span> {{ formatPrice(scope.row.additional_amount) }}</p>
                <p class="text-[10px]"><span class="text-emerald-600">Cọc:</span> {{ formatPrice(scope.row.deposit_amount) }}</p>
            </template>
        </el-table-column>
        <el-table-column label="Viettel Amt" width="120" align="center">
            <template slot-scope="scope">
                <span class="bg-gray-800 text-white px-1.5 py-0.5 rounded text-[10px] font-bold">{{ formatPrice(scope.row.viettel_amount) }}</span>
            </template>
        </el-table-column>
        <el-table-column label="Status" width="120">
          <template slot-scope="scope">
            <el-tag :type="getStatusType(scope.row.status)" size="mini" effect="dark" class="rounded-pill uppercase font-bold text-[10px]">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Action" align="right" width="100">
          <template slot-scope="scope">
            <el-dropdown trigger="click">
              <el-button type="text" icon="el-icon-more"></el-button>
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item icon="el-icon-camera" @click.native="captureOrder(scope.row)">Capture Order</el-dropdown-item>
                <el-dropdown-item icon="el-icon-view" @click.native="viewOrder(scope.row)">View Detail</el-dropdown-item>
                <el-dropdown-item icon="el-icon-delete" @click.native="confirmDelete(scope.row.id)" class="text-danger">Delete</el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- Pagination -->
    <div v-if="orders.length" class="mt-8 flex justify-center">
      <el-pagination background layout="prev, pager, next" :current-page.sync="page" :total="total" :page-size="10" @current-change="handlePageChange"></el-pagination>
    </div>

    <!-- Hidden Capture Logic -->
    <div style="position: absolute; left: -9999px; top: 0;">
        <div ref="receiptTemplate" id="order-receipt-v2" class="p-8 bg-white w-[450px] font-sans text-slate-800">
            <div class="text-center mb-6">
                <h1 class="text-4xl font-extrabold tracking-[4px] m-0 text-slate-900">VINTAGE</h1>
                <p class="text-sm font-bold tracking-[2px] text-blue-600 mt-1 uppercase">Always True</p>
            </div>
            <div class="bg-slate-50 rounded-2xl p-6 mb-5">
                <div class="flex items-center mb-3">
                    <div class="bg-blue-600 text-white w-10 h-10 rounded-full flex items-center justify-center font-bold mr-3">P</div>
                    <div>
                        <span class="text-[10px] text-slate-400 uppercase font-bold tracking-wider">KHÁCH HÀNG</span>
                        <h2 class="text-lg font-bold m-0">{{ currentCapture?.name }}</h2>
                    </div>
                </div>
                <div class="text-sm text-slate-600 space-y-1">
                    <p><i class="el-icon-phone mr-2"></i>{{ currentCapture?.phone }}</p>
                    <p><i class="el-icon-location mr-2"></i>{{ currentCapture?.address }}</p>
                </div>
            </div>
            <div v-if="currentCapture?.note" class="bg-yellow-50 border border-yellow-100 rounded-xl p-4 mb-6"><p class="text-sm font-medium">{{ currentCapture.note }}</p></div>
            <h3 class="text-[10px] font-black text-slate-300 uppercase tracking-widest border-b border-slate-100 pb-2 mb-4">ĐƠN HÀNG CỦA BẠN</h3>
            <div class="flex flex-wrap gap-4 mb-6">
                <div v-for="p in currentCapture?.products" :key="p.id" class="w-[calc(50%-8px)] text-center">
                    <img :src="p.path_thumb" class="w-full aspect-[3/4] object-cover rounded-xl bg-slate-50 mb-2 border border-slate-100">
                    <span class="text-sm font-bold">{{ p.name }}</span>
                </div>
            </div>
            <div class="bg-slate-50 rounded-2xl p-6 mb-8 border border-slate-100 font-bold">
                <div class="flex justify-between mb-3 text-sm text-slate-500"><span>Tổng đơn</span><span>{{ formatPrice(currentCapture?.total_amount) }}</span></div>
                <div class="pt-4 border-t border-slate-200 flex justify-between items-center"><span class="text-lg">Viettel Amount</span><span class="text-2xl text-blue-600">{{ formatPrice(currentCapture?.viettel_amount) }}</span></div>
            </div>
            <div class="pt-6 border-t border-dashed border-slate-200 text-slate-500">
                <div class="flex gap-4 mb-4">
                    <div class="flex-1">
                        <h4 class="font-black text-slate-800 uppercase mb-2 text-[10px]">🚚 GIAO HÀNG</h4>
                        <ul class="list-none p-0 text-[10px] space-y-1">
                            <li>• Viettel Post: HCM 1-2 ngày, tỉnh 3-5 ngày.</li>
                            <li>• Phí ship đồng giá: 20k toàn quốc.</li>
                            <li>• Hàng được giặt sấy, là ủi và xếp gọn tỉ mỉ.</li>
                        </ul>
                    </div>
                    <div class="flex-1">
                        <h4 class="font-black text-slate-800 uppercase mb-2 text-[10px]">🔄 ĐỔI TRẢ</h4>
                        <ul class="list-none p-0 text-[10px] space-y-1">
                           <li>• Đổi trả nếu hàng không đúng mô tả, sai số đo >2cm.</li>
                           <li>• Không nhận trả hàng với lý do cá nhân (không thích, không hợp).</li>
                        </ul>
                    </div>
                </div>
                <div class="text-center border-t border-dashed border-slate-200 pt-4">
                    <h3 class="font-bold text-sm mb-1 text-slate-900">Thank you for your order!</h3>
                    <p class="text-[10px] text-slate-400">Vintage Always True • Vintage Menswear</p>
                </div>
            </div>
        </div>
    </div>

    <!-- Capture Preview -->
    <el-dialog title="Order Receipt Generated" :visible.sync="captureVisible" width="500px" custom-class="mobile-dialog">
        <div class="text-center"><img :src="captureImg" class="w-full max-h-[70vh] w-auto mx-auto rounded-2xl border mb-6 shadow-sm"><el-button type="primary" size="medium" icon="el-icon-download" @click="downloadCapture">Download Image</el-button></div>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'OrderManager',
  data() {
    return {
      activeTab: 'NEW',
      viewMode: 'grid',
      orders: [], loading: false, total: 0, page: 1, search: '',
      currentCapture: null, captureVisible: false, captureImg: '',
      sizeMapping: {
        'VAI': 'V', 'NGỰC': 'N', 'EO': 'E', 'DÀI ÁO': 'D', 'DÀI ÁO SAU': 'DS', 
        'DÀI TAY': 'DT', 'EO QUẦN': 'E', 'ĐÁY': 'ĐA', 'ĐÙI': 'Đ', 'DÀI QUẦN': 'D', 
        'ỐNG': 'Ô', 'DƯ LAI': 'L'
      }
    }
  },
  mounted() { this.fetchOrders(); },
  methods: {
    formatPrice(val) { return val !== undefined ? new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND', maximumFractionDigits: 0 }).format(val) : '0đ'; },
    getStatusType(s) { return { 'NEW': 'info', 'COMPLETED': 'success', 'A SHIP NOW': 'primary', 'ON HOLD': 'warning' }[s] || 'info'; },
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
    async fetchOrders(page = 1) {
      this.page = page; this.loading = true;
      try {
        const params = { page, search: this.search };
        if (this.activeTab === 'NEW') params.status = 'NEW';
        
        const res = await axios.get('/salev2/api/orders', { params });
        this.orders = res.data.data; this.total = res.data.total;
      } catch (err) { console.error(err); } finally { this.loading = false; }
    },
    handleTabClick() { this.fetchOrders(1); },
    debouncedFetch: _.debounce(function() { this.fetchOrders(1); }, 500),
    handlePageChange(p) { this.fetchOrders(p); },
    async captureOrder(order) {
        this.currentCapture = order;
        this.$nextTick(async () => {
            const canvas = await html2canvas(this.$refs.receiptTemplate, { useCORS: true, scale: 2 });
            this.captureImg = canvas.toDataURL('image/jpeg', 0.9); this.captureVisible = true;
        });
    },
    downloadCapture() {
        const link = document.createElement('a'); link.download = `Order_${this.currentCapture.name}.jpg`; link.href = this.captureImg; link.click();
    },
    async confirmDelete(id) {
        if (confirm('Delete this order?')) {
            try { await axios.post(`/salev2/api/order/${id}/delete`); this.fetchOrders(this.page); } catch(e) {}
        }
    },
    async completeOrder(id) {
        try { await axios.post(`/salev2/api/order/${id}`, { status: 'COMPLETED' }); this.fetchOrders(this.page); } catch(e) {}
    },

    async updateStatus(order, status) {
        try { 
            await axios.post(`/salev2/api/order/${order.id}`, { status });
            this.$message.success(`Order marked as ${status}`);
            this.fetchOrders(this.page); 
        } catch(e) {
            this.$message.error('Update failed');
        }
    },
    viewCustomerHistory(order) {
        // Redirect to Orders page filtered by this phone
        if(this.$route.path === '/salev2/orders') {
            this.search = order.phone;
            this.fetchOrders(1);
        } else {
            this.$router.push({ path: '/salev2/orders', query: { search: order.phone } });
        }
    }
  },
  mounted() {
      if(this.$route.query.search) {
          this.search = this.$route.query.search;
      }
      this.fetchOrders();
  }
}
</script>

<style>
.scrollbar-hide::-webkit-scrollbar { display: none; }
.scrollbar-hide { -ms-overflow-style: none; scrollbar-width: none; }
@media (max-width: 768px) {
    .mobile-dialog { width: 95% !important; margin-top: 5vh !important; }
    .mobile-dialog .el-dialog__body { padding: 10px 15px; }
}
</style>
