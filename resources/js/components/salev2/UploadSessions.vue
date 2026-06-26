<template>
  <div class="px-4 md:px-6 py-8">
    <div class="flex justify-between items-center mb-6 mobile-hide">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 tracking-tight">Upload Sessions</h1>
        <p class="text-slate-500 mt-1">Browse products grouped by upload batch</p>
      </div>
    </div>

    <el-tabs v-model="activeTab" type="card" class="premium-tabs mb-6" @tab-click="handleTabChange">
      <el-tab-pane label="Sessions" name="sessions"></el-tab-pane>
      <el-tab-pane label="Latest Uploads" name="latest"></el-tab-pane>
      <el-tab-pane label="Sold by Date" name="sold"></el-tab-pane>
    </el-tabs>

    <!-- ===== Tab: Sessions ===== -->
    <div v-if="activeTab === 'sessions'" v-loading="loadingSessions">
      <div v-if="!sessions.length && !loadingSessions" class="text-center py-20 bg-white rounded-3xl border border-dashed border-slate-200">
        <i class="el-icon-upload text-4xl text-slate-200 mb-4 block"></i>
        <p class="text-slate-400 font-medium">No upload sessions found</p>
      </div>

      <template v-for="(group, gIdx) in sessionCalendar">
        <!-- Week separator -->
        <div v-if="gIdx === 0 || getISOWeek(group.date) !== getISOWeek(sessionCalendar[gIdx - 1].date)"
             :key="'week-s-' + group.date"
             class="flex items-center gap-3 my-5">
          <div class="flex-1 h-px bg-slate-200"></div>
          <span class="text-[11px] font-black text-slate-400 uppercase tracking-widest px-2">{{ getISOWeek(group.date) }}</span>
          <div class="flex-1 h-px bg-slate-200"></div>
        </div>

        <!-- Date group -->
        <div :key="'date-s-' + group.date" class="mb-3">
          <!-- Date label -->
          <div :class="['flex items-center gap-2 px-3 py-1.5 mb-2 rounded-lg text-sm font-bold', group.isEmpty ? 'text-red-400 bg-red-50' : 'text-slate-500 bg-slate-50']">
            <i :class="group.isEmpty ? 'el-icon-warning-outline' : 'el-icon-date'"></i>
            <span>{{ formatDateShort(group.date) }}</span>
            <span v-if="group.isEmpty" class="text-xs font-normal text-red-300">— không có upload</span>
            <span v-else class="text-xs font-normal text-slate-400">— {{ group.sessions.length }} session(s)</span>
          </div>

          <!-- Sessions within this date -->
          <div v-for="session in group.sessions" :key="session.session_id" class="mb-3 ml-4 bg-white rounded-2xl border border-slate-100 shadow-sm overflow-hidden">
            <div class="flex flex-wrap items-center gap-3 p-4 cursor-pointer hover:bg-slate-50 transition-colors" @click="toggleSession(session)">
              <div class="flex items-center gap-2 flex-1 min-w-0">
                <i :class="openSessions[session.session_id] ? 'el-icon-arrow-down' : 'el-icon-arrow-right'" class="text-slate-400 flex-shrink-0"></i>
                <div class="min-w-0">
                  <p class="text-xs text-slate-400 font-mono truncate">{{ session.session_id }}</p>
                  <p class="font-bold text-slate-800 text-sm">{{ formatDate(session.created_at) }}</p>
                </div>
              </div>
              <div class="flex flex-wrap items-center gap-2">
                <el-tag size="small" type="info">{{ session.total_products }} sp</el-tag>
                <el-tag size="small" type="success" v-if="session.sold_products">{{ session.sold_products }} sold</el-tag>
                <el-tag v-for="(count, type) in session.type_breakdown" :key="type" size="small" class="font-bold">{{ type }}: {{ count }}</el-tag>
              </div>
            </div>

            <div v-if="openSessions[session.session_id]" class="border-t border-slate-100">
              <div v-loading="loadingProducts[session.session_id]" class="p-4">
                <div v-if="sessionProducts[session.session_id]">
                  <div v-if="!sessionProducts[session.session_id].length" class="text-center py-8 text-slate-400">
                    <i class="el-icon-box text-3xl block mb-2"></i>
                    <p>No products in this session</p>
                  </div>
                  <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-3">
                    <div v-for="log in sessionProducts[session.session_id]" :key="log.id"
                         class="group relative bg-white rounded-xl overflow-hidden border border-slate-100 shadow-sm hover:shadow-xl transition-all duration-300">
                      <template v-if="log.product">
                        <div class="aspect-[3/4] overflow-hidden bg-slate-50 relative group">
                          <img v-if="log.product.image_thumb_scale_url || log.product.path_thumb" :src="log.product.image_thumb_scale_url || log.product.path_thumb" class="w-full h-full object-contain hover:scale-105 transition-transform duration-500 p-1" @error="e => e.target.src = log.product.path_thumb">
                          <div v-else class="w-full h-full flex items-center justify-center text-slate-300"><i class="el-icon-picture text-4xl"></i></div>
                          <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity hidden md:flex flex-col justify-center items-center gap-2 p-6 z-10" @click.stop="editProduct(log.product)">
                            <el-button type="primary" size="medium" class="w-40 font-bold" @click.stop="editProduct(log.product)">DETAIL</el-button>
                            <el-button type="success" size="medium" class="w-40 font-bold" icon="el-icon-document-copy" @click.stop="copyProductInfo(log.product)">COPY SIZE</el-button>
                            <el-button type="info" size="medium" class="w-40 font-bold" icon="el-icon-camera" @click.stop="copyMain(log.product)">COPY MAIN</el-button>
                            <el-button type="warning" size="medium" class="w-40 font-bold" icon="el-icon-picture-outline" @click.stop="copyDetail(log.product)">COPY DETAIL</el-button>
                          </div>
                        </div>
                        <div class="p-3">
                          <div class="text-center mb-2">
                            <h4 class="font-bold text-slate-800 text-sm uppercase leading-tight truncate">{{ log.product.name }}</h4>
                            <p class="text-blue-600 font-bold text-sm">{{ formatPrice(log.product.price) }}</p>
                            <p class="text-[10px] text-slate-400 mt-0.5">{{ formatDate(log.created_at) }}</p>
                          </div>
                          <div class="text-center mb-3 px-1">
                            <span class="text-[10px] font-black bg-blue-50 text-blue-600 px-2 py-1 rounded uppercase inline-block">{{ formatSizes(log.product).join(' - ') }}</span>
                          </div>
                          <div class="flex justify-between items-center pt-2 border-t border-slate-50 gap-2">
                            <div class="flex-1 text-center md:text-left">
                              <el-tag :type="getStatusTagType(log.product.status)" size="mini" class="text-[8px] font-bold uppercase rounded-md px-1">{{ log.product.status }}</el-tag>
                            </div>
                            <div class="md:hidden flex justify-end">
                              <el-dropdown trigger="click" @command="(cmd) => handleMobileAction(cmd, log.product)">
                                <span class="el-dropdown-link p-2"><i class="el-icon-more transform rotate-90 text-lg font-bold text-slate-600"></i></span>
                                <el-dropdown-menu slot="dropdown">
                                  <el-dropdown-item command="edit" icon="el-icon-edit">Detail</el-dropdown-item>
                                  <el-dropdown-item command="copy" icon="el-icon-document-copy">Copy Size</el-dropdown-item>
                                  <el-dropdown-item command="copy_main" icon="el-icon-camera">Copy Main</el-dropdown-item>
                                  <el-dropdown-item command="copy_detail" icon="el-icon-picture-outline">Copy Detail</el-dropdown-item>
                                </el-dropdown-menu>
                              </el-dropdown>
                            </div>
                          </div>
                        </div>
                      </template>
                      <template v-else>
                        <div class="aspect-[3/4] bg-slate-50 flex flex-col items-center justify-center gap-2 p-3">
                          <i class="el-icon-picture-outline text-3xl text-slate-300"></i>
                          <p class="text-[10px] text-slate-400 text-center truncate w-full">{{ log.file_name }}</p>
                          <el-tag :type="getLogStatusType(log.status)" size="mini" class="text-[8px] font-bold uppercase">{{ log.status }}</el-tag>
                          <p v-if="log.message" class="text-[9px] text-slate-400 text-center line-clamp-2">{{ log.message }}</p>
                        </div>
                      </template>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- ===== Tab: Latest Uploads ===== -->
    <div v-if="activeTab === 'latest'">
      <div class="bg-white p-4 rounded-2xl border border-slate-100 shadow-sm mb-6 flex flex-wrap items-center gap-4">
        <el-checkbox v-model="onlyAvailable" @change="fetchLatestProducts(1)">
          <span class="font-bold text-slate-700 text-sm">Only Available</span>
        </el-checkbox>
        <el-select v-model="latestFilterType" placeholder="All Types" clearable size="small" style="width:150px" @change="fetchLatestProducts(1)">
          <el-option v-for="t in productTypes" :key="t" :label="t" :value="t"></el-option>
        </el-select>
        <div class="flex-1"></div>
        <el-button type="warning" size="small" icon="el-icon-time" :disabled="!latestSelected.length" :loading="settingOnHold" @click="bulkSetOnHold">
          Set ON_HOLD <span v-if="latestSelected.length">({{ latestSelected.length }})</span>
        </el-button>
        <el-button size="small" type="danger" plain icon="el-icon-delete" v-if="latestSelected.length" @click="latestSelected = []">Clear</el-button>
      </div>

      <div v-loading="loadingLatest">
        <div v-if="!latestProducts.length && !loadingLatest" class="text-center py-20 bg-white rounded-3xl border border-dashed border-slate-200">
          <i class="el-icon-box text-4xl text-slate-200 mb-4 block"></i>
          <p class="text-slate-400 font-medium">No products found</p>
        </div>
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-3">
          <div v-for="product in latestProducts" :key="product.id"
               class="group relative bg-white rounded-xl overflow-hidden border border-slate-100 shadow-sm hover:shadow-xl transition-all duration-300"
               :class="{'ring-2 ring-blue-500': isLatestSelected(product.id)}">
            <div class="aspect-[3/4] overflow-hidden bg-slate-50 relative group">
              <img v-if="product.image_thumb_scale_url || product.path_thumb" :src="product.image_thumb_scale_url || product.path_thumb" class="w-full h-full object-contain hover:scale-105 transition-transform duration-500 p-1" @error="e => e.target.src = product.path_thumb">
              <div v-else class="w-full h-full flex items-center justify-center text-slate-300"><i class="el-icon-picture text-4xl"></i></div>
              <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity hidden md:flex flex-col justify-center items-center gap-2 p-6 z-10" @click.stop="editProduct(product)">
                <el-button type="primary" size="medium" class="w-40 font-bold" @click.stop="editProduct(product)">DETAIL</el-button>
                <el-button type="success" size="medium" class="w-40 font-bold" icon="el-icon-document-copy" @click.stop="copyProductInfo(product)">COPY SIZE</el-button>
                <el-button type="info" size="medium" class="w-40 font-bold" icon="el-icon-camera" @click.stop="copyMain(product)">COPY MAIN</el-button>
                <el-button type="warning" size="medium" class="w-40 font-bold" icon="el-icon-picture-outline" @click.stop="copyDetail(product)">COPY DETAIL</el-button>
              </div>
            </div>
            <div class="p-3">
              <div class="text-center mb-2">
                <h4 class="font-bold text-slate-800 text-sm uppercase leading-tight truncate">{{ product.name }}</h4>
                <p class="text-blue-600 font-bold text-sm">{{ formatPrice(product.price) }}</p>
                <p class="text-[10px] text-slate-400 mt-0.5">{{ product.latest_upload ? formatDate(product.latest_upload) : 'undefined' }}</p>
              </div>
              <div class="text-center mb-3 px-1">
                <span class="text-[10px] font-black bg-blue-50 text-blue-600 px-2 py-1 rounded uppercase inline-block">{{ formatSizes(product).join(' - ') }}</span>
              </div>
              <div class="flex justify-between items-center pt-2 border-t border-slate-50 gap-2">
                <div class="order-1 md:order-2 flex items-center">
                  <el-checkbox :value="isLatestSelected(product.id)" @change="toggleLatestSelection(product)" class="transform scale-125"></el-checkbox>
                </div>
                <div class="order-2 md:order-1 flex-1 text-center md:text-left">
                  <el-tag :type="getStatusTagType(product.status)" size="mini" class="text-[8px] font-bold uppercase rounded-md px-1">{{ product.status }}</el-tag>
                </div>
                <div class="order-3 md:hidden flex justify-end">
                  <el-dropdown trigger="click" @command="(cmd) => handleMobileAction(cmd, product)">
                    <span class="el-dropdown-link p-2"><i class="el-icon-more transform rotate-90 text-lg font-bold text-slate-600"></i></span>
                    <el-dropdown-menu slot="dropdown">
                      <el-dropdown-item command="edit" icon="el-icon-edit">Detail</el-dropdown-item>
                      <el-dropdown-item command="copy" icon="el-icon-document-copy">Copy Size</el-dropdown-item>
                      <el-dropdown-item command="copy_main" icon="el-icon-camera">Copy Main</el-dropdown-item>
                      <el-dropdown-item command="copy_detail" icon="el-icon-picture-outline">Copy Detail</el-dropdown-item>
                    </el-dropdown-menu>
                  </el-dropdown>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="mt-8 flex justify-center">
          <el-pagination background layout="prev, pager, next" :current-page.sync="latestPage" :total="latestTotal" :page-size="60" @current-change="fetchLatestProducts"></el-pagination>
        </div>
      </div>
    </div>

    <!-- ===== Tab: Sold by Date ===== -->
    <div v-if="activeTab === 'sold'" v-loading="loadingSold">
      <div v-if="!soldDates.length && !loadingSold" class="text-center py-20 bg-white rounded-3xl border border-dashed border-slate-200">
        <i class="el-icon-box text-4xl text-slate-200 mb-4 block"></i>
        <p class="text-slate-400 font-medium">No data found</p>
      </div>

      <template v-for="(row, rIdx) in soldDates">
        <!-- Week separator -->
        <div v-if="rIdx === 0 || getISOWeek(row.date) !== getISOWeek(soldDates[rIdx - 1].date)"
             :key="'week-d-' + row.date"
             class="flex items-center gap-3 my-5">
          <div class="flex-1 h-px bg-slate-200"></div>
          <span class="text-[11px] font-black text-slate-400 uppercase tracking-widest px-2">{{ getISOWeek(row.date) }}</span>
          <div class="flex-1 h-px bg-slate-200"></div>
        </div>

        <!-- Date row -->
        <div :key="'date-d-' + row.date"
             :class="['mb-3 rounded-2xl border shadow-sm overflow-hidden', row.is_empty ? 'bg-red-50 border-red-100' : 'bg-white border-slate-100']">
          <div class="flex flex-wrap items-center gap-3 p-4 cursor-pointer hover:bg-opacity-80 transition-colors"
               @click="!row.is_empty && toggleSoldDate(row)">
            <div class="flex items-center gap-2 flex-1 min-w-0">
              <i v-if="!row.is_empty" :class="openSoldDates[row.date] ? 'el-icon-arrow-down' : 'el-icon-arrow-right'" class="text-slate-400 flex-shrink-0"></i>
              <i v-else class="el-icon-warning-outline text-red-400 flex-shrink-0"></i>
              <p :class="['font-bold text-sm', row.is_empty ? 'text-red-400' : 'text-slate-800']">{{ formatDateShort(row.date) }}</p>
              <span v-if="row.is_empty" class="text-xs text-red-300 font-normal">— không bán được</span>
            </div>
            <div v-if="!row.is_empty" class="flex flex-wrap items-center gap-2">
              <el-tag v-for="(count, type) in row.type_breakdown" :key="type" size="small" class="font-bold">{{ type }}: {{ count }}</el-tag>
              <el-tag size="small" type="primary">{{ row.total }} sp</el-tag>
              <el-tag size="small" type="success">{{ formatPrice(row.total_revenue) }}</el-tag>
            </div>
          </div>

          <div v-if="!row.is_empty && openSoldDates[row.date]" class="border-t border-slate-100">
            <div v-loading="loadingSoldProducts[row.date]" class="p-4">
              <div v-if="soldDateProducts[row.date]">
                <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-3">
                  <div v-for="product in soldDateProducts[row.date]" :key="product.id"
                       class="group relative bg-white rounded-xl overflow-hidden border border-slate-100 shadow-sm hover:shadow-xl transition-all duration-300">
                    <div class="aspect-[3/4] overflow-hidden bg-slate-50 relative group">
                      <img v-if="product.image_thumb_scale_url || product.path_thumb" :src="product.image_thumb_scale_url || product.path_thumb" class="w-full h-full object-contain hover:scale-105 transition-transform duration-500 p-1" @error="e => e.target.src = product.path_thumb">
                      <div v-else class="w-full h-full flex items-center justify-center text-slate-300"><i class="el-icon-picture text-4xl"></i></div>
                      <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity hidden md:flex flex-col justify-center items-center gap-2 p-6 z-10" @click.stop="editProduct(product)">
                        <el-button type="primary" size="medium" class="w-40 font-bold" @click.stop="editProduct(product)">DETAIL</el-button>
                        <el-button type="success" size="medium" class="w-40 font-bold" icon="el-icon-document-copy" @click.stop="copyProductInfo(product)">COPY SIZE</el-button>
                        <el-button type="info" size="medium" class="w-40 font-bold" icon="el-icon-camera" @click.stop="copyMain(product)">COPY MAIN</el-button>
                        <el-button type="warning" size="medium" class="w-40 font-bold" icon="el-icon-picture-outline" @click.stop="copyDetail(product)">COPY DETAIL</el-button>
                      </div>
                    </div>
                    <div class="p-3">
                      <div class="text-center mb-2">
                        <h4 class="font-bold text-slate-800 text-sm uppercase leading-tight truncate">{{ product.name }}</h4>
                        <p class="text-blue-600 font-bold text-sm">{{ formatPrice(product.price) }}</p>
                      </div>
                      <div class="text-center mb-3 px-1">
                        <span class="text-[10px] font-black bg-blue-50 text-blue-600 px-2 py-1 rounded uppercase inline-block">{{ formatSizes(product).join(' - ') }}</span>
                      </div>
                      <div class="flex justify-between items-center pt-2 border-t border-slate-50 gap-2">
                        <div class="flex-1 text-center md:text-left">
                          <el-tag :type="getStatusTagType(product.status)" size="mini" class="text-[8px] font-bold uppercase rounded-md px-1">{{ product.status }}</el-tag>
                        </div>
                        <div class="md:hidden flex justify-end">
                          <el-dropdown trigger="click" @command="(cmd) => handleMobileAction(cmd, product)">
                            <span class="el-dropdown-link p-2"><i class="el-icon-more transform rotate-90 text-lg font-bold text-slate-600"></i></span>
                            <el-dropdown-menu slot="dropdown">
                              <el-dropdown-item command="edit" icon="el-icon-edit">Detail</el-dropdown-item>
                              <el-dropdown-item command="copy" icon="el-icon-document-copy">Copy Size</el-dropdown-item>
                              <el-dropdown-item command="copy_main" icon="el-icon-camera">Copy Main</el-dropdown-item>
                              <el-dropdown-item command="copy_detail" icon="el-icon-picture-outline">Copy Detail</el-dropdown-item>
                            </el-dropdown-menu>
                          </el-dropdown>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
export default {
  name: 'UploadSessions',
  data() {
    return {
      activeTab: 'sessions',
      sessions: [],
      loadingSessions: false,
      openSessions: {},
      sessionProducts: {},
      loadingProducts: {},
      latestProducts: [],
      loadingLatest: false,
      latestTotal: 0,
      latestPage: 1,
      onlyAvailable: true,
      latestSelected: [],
      settingOnHold: false,
      latestFilterType: '',
      soldDates: [],
      loadingSold: false,
      openSoldDates: {},
      soldDateProducts: {},
      loadingSoldProducts: {},
      productTypes: ['TROUSERS', 'JACKET', 'SHIRT', 'BLAZER', 'TIE', 'GILE', 'BELT', 'POLO SHIRT', 'HAT', 'SUIT'],
      sizeMapping: {
        'VAI': 'V', 'NGỰC': 'N', 'EO': 'E', 'DÀI ÁO': 'D', 'DÀI ÁO SAU': 'DS',
        'DÀI TAY': 'DT', 'EO QUẦN': 'E', 'ĐÁY': 'ĐA', 'ĐÙI': 'Đ', 'DÀI QUẦN': 'D',
        'ỐNG': 'Ô', 'DƯ LAI': 'L'
      },
    };
  },
  computed: {
    sessionCalendar() {
      if (!this.sessions.length) return [];
      const byDate = {};
      this.sessions.forEach(s => {
        const date = s.created_at ? s.created_at.substring(0, 10) : null;
        if (!date) return;
        if (!byDate[date]) byDate[date] = [];
        byDate[date].push(s);
      });
      const dates = Object.keys(byDate);
      if (!dates.length) return [];
      const oldest = dates.slice().sort()[0];
      const result = [];
      const cur = new Date();
      const end = new Date(oldest + 'T00:00:00');
      while (cur >= end) {
        const d = cur.toISOString().substring(0, 10);
        result.push({ date: d, sessions: byDate[d] || [], isEmpty: !byDate[d] });
        cur.setDate(cur.getDate() - 1);
      }
      return result;
    },
  },
  mounted() {
    this.fetchSessions();
  },
  methods: {
    handleTabChange() {
      if (this.activeTab === 'latest' && !this.latestProducts.length) this.fetchLatestProducts(1);
      if (this.activeTab === 'sold' && !this.soldDates.length) this.fetchSoldDates();
    },

    getISOWeek(dateStr) {
      const d = new Date(dateStr + 'T00:00:00');
      d.setDate(d.getDate() + 3 - (d.getDay() + 6) % 7);
      const w1 = new Date(d.getFullYear(), 0, 4);
      const week = 1 + Math.round(((d - w1) / 86400000 - 3 + (w1.getDay() + 6) % 7) / 7);
      return `Tuần ${week} / ${d.getFullYear()}`;
    },

    // ---- Sessions ----
    async fetchSessions() {
      this.loadingSessions = true;
      try {
        const res = await axios.get('/salev2/api/upload-sessions');
        this.sessions = res.data.sessions;
      } catch (err) {
        this.$message.error('Failed to load sessions');
      } finally {
        this.loadingSessions = false;
      }
    },
    async toggleSession(session) {
      const id = session.session_id;
      this.$set(this.openSessions, id, !this.openSessions[id]);
      if (this.openSessions[id] && !this.sessionProducts[id]) {
        this.$set(this.loadingProducts, id, true);
        try {
          const res = await axios.get(`/salev2/api/upload-sessions/${id}/products`);
          this.$set(this.sessionProducts, id, res.data.logs);
        } catch (err) {
          this.$message.error('Failed to load session products');
        } finally {
          this.$set(this.loadingProducts, id, false);
        }
      }
    },

    // ---- Latest uploads ----
    async fetchLatestProducts(page = 1) {
      this.latestPage = page;
      this.loadingLatest = true;
      try {
        const params = { page, limit: 60 };
        if (this.onlyAvailable) params.status = 'AVAILABLE';
        if (this.latestFilterType) params.type = this.latestFilterType;
        const res = await axios.get('/salev2/api/products-by-upload', { params });
        this.latestProducts = res.data.data;
        this.latestTotal = res.data.total;
      } catch (err) {
        this.$message.error('Failed to load products');
      } finally {
        this.loadingLatest = false;
      }
    },
    isLatestSelected(id) { return this.latestSelected.some(p => p.id === id); },
    toggleLatestSelection(product) {
      const idx = this.latestSelected.findIndex(p => p.id === product.id);
      if (idx > -1) this.latestSelected.splice(idx, 1);
      else this.latestSelected.push({ id: product.id });
    },
    async bulkSetOnHold() {
      if (!this.latestSelected.length) return;
      this.settingOnHold = true;
      try {
        await Promise.all(this.latestSelected.map(p => axios.post(`/salev2/api/product/${p.id}/update`, { status: 'ON_HOLD' })));
        this.$message.success(`${this.latestSelected.length} product(s) set to ON_HOLD`);
        this.latestSelected = [];
        this.fetchLatestProducts(this.latestPage);
      } catch (err) {
        this.$message.error('Failed to update status');
      } finally {
        this.settingOnHold = false;
      }
    },

    // ---- Sold by date ----
    async fetchSoldDates() {
      this.loadingSold = true;
      try {
        const res = await axios.get('/salev2/api/sold-by-date');
        this.soldDates = res.data.dates;
        const top3 = this.soldDates.filter(r => !r.is_empty).slice(0, 3);
        await Promise.all(top3.map(row => this.toggleSoldDate(row)));
      } catch (err) {
        this.$message.error('Failed to load sold dates');
      } finally {
        this.loadingSold = false;
      }
    },
    async toggleSoldDate(row) {
      if (row.is_empty) return;
      const date = row.date;
      this.$set(this.openSoldDates, date, !this.openSoldDates[date]);
      if (this.openSoldDates[date] && !this.soldDateProducts[date]) {
        this.$set(this.loadingSoldProducts, date, true);
        try {
          const res = await axios.get(`/salev2/api/sold-by-date/${date}/products`);
          this.$set(this.soldDateProducts, date, res.data.products);
        } catch (err) {
          this.$message.error('Failed to load products');
        } finally {
          this.$set(this.loadingSoldProducts, date, false);
        }
      }
    },

    // ---- Shared helpers ----
    formatDateShort(dateStr) {
      if (!dateStr) return '';
      const d = new Date(dateStr + 'T00:00:00');
      return d.toLocaleDateString('vi-VN', { weekday: 'short', day: '2-digit', month: '2-digit', year: 'numeric' });
    },
    formatDate(dateStr) {
      if (!dateStr) return '';
      const d = new Date(dateStr);
      return d.toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' });
    },
    formatPrice(val) {
      if (!val) return '0đ';
      return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND', maximumFractionDigits: 0 }).format(val);
    },
    getStatusTagType(status) {
      return { 'AVAILABLE': 'success', 'SOLD': 'info', 'ON_HOLD': 'warning' }[status] || 'info';
    },
    getLogStatusType(status) {
      return { 'SUCCESS': 'success', 'FAILED': 'danger', 'PENDING': 'info', 'PROCESSING': 'warning', 'ERROR': 'danger' }[status] || 'info';
    },
    formatSizes(product) {
      if (!product.sizes) return [];
      const results = [];
      const duLai = product.sizes.find(s => s.name === 'DƯ LAI');
      product.sizes.forEach(size => {
        if (size.name === 'DƯ LAI') return;
        const label = this.sizeMapping[size.name] || size.name;
        const value = Math.floor(size.value);
        results.push(size.name === 'DÀI QUẦN' && duLai ? `${label}${value}(+${Math.floor(duLai.value)})` : `${label}${value}`);
      });
      return results;
    },
    editProduct(product) {
      this.$router.push({ name: 'product-edit', params: { id: product.id } });
    },
    async copyProductInfo(product) {
      try {
        await navigator.clipboard.writeText(`${product.name} ${this.formatSizes(product).join(' - ')}`);
        this.$message.success('Copied!');
      } catch { this.$message.error('Failed to copy'); }
    },
    async copyMain(product) {
      if (!navigator.share) { alert('Not supported'); return; }
      if (!product.path_thumb) { this.$message.warning('No main image.'); return; }
      const loading = this.$loading({ lock: true, text: 'Preparing...', background: 'rgba(255,255,255,0.7)' });
      try {
        const res = await fetch(product.path_thumb);
        const file = new File([await res.blob()], `Main-${product.name}.jpg`, { type: 'image/jpeg' });
        await navigator.share({ files: [file] });
      } catch (err) { if (err.name !== 'AbortError') alert('Share failed: ' + err.message); }
      finally { loading.close(); }
    },
    async copyDetail(product) {
      if (!navigator.share) { alert('Not supported'); return; }
      if (!product.images?.length) { this.$message.warning('No detail images.'); return; }
      const loading = this.$loading({ lock: true, text: 'Preparing...', background: 'rgba(255,255,255,0.7)' });
      try {
        const files = [];
        for (let i = 0; i < product.images.length; i++) {
          const r = await fetch(product.images[i].file_path);
          if (r.ok) files.push(new File([await r.blob()], `Detail-${product.id}-${i + 1}.jpg`, { type: 'image/jpeg' }));
        }
        if (!files.length) { this.$message.warning('Could not load images'); return; }
        await navigator.share({ files });
      } catch (err) { if (err.name !== 'AbortError') alert('Share failed: ' + err.message); }
      finally { loading.close(); }
    },
    handleMobileAction(cmd, product) {
      const map = { edit: 'editProduct', copy: 'copyProductInfo', copy_main: 'copyMain', copy_detail: 'copyDetail' };
      if (map[cmd]) this[map[cmd]](product);
    },
  },
};
</script>

<style>
.premium-tabs .el-tabs__item {
    font-weight: 900;
    text-transform: uppercase;
    font-size: 12px;
    letter-spacing: 1px;
}
.premium-tabs .el-tabs__header { margin-bottom: 0; }
</style>
