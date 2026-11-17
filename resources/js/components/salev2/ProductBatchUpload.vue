<template>
  <div class="px-4 md:px-6 py-6 max-w-6xl mx-auto h-[calc(100vh-60px)] flex flex-col">
    <!-- Header -->
    <div class="flex justify-between items-center mb-4 flex-shrink-0">
        <div class="flex items-center gap-3">
             <el-button icon="el-icon-back" circle @click="$router.push('/salev2/products')"></el-button>
             <h1 class="text-2xl font-bold text-slate-800">Smart Batch Upload</h1>
        </div>
        <div class="flex gap-3">
             <input type="file" ref="fileInput" multiple accept="image/*" class="hidden" @change="handleFiles">
             <el-button type="primary" icon="el-icon-upload" @click="$refs.fileInput.click()">Add Images</el-button>
             <el-button type="success" :disabled="!groups.length" :loading="uploading" @click="uploadAll">
                 Upload {{ groups.length }} Products
             </el-button>
        </div>
    </div>

    <!-- Facebook Album Resync Bar (Mới nhất lên đầu) -->
    <div v-if="albumCategories.length" class="mb-4 flex-shrink-0 bg-white px-4 py-3 rounded-xl border border-slate-200 shadow-sm space-y-2.5">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div class="flex flex-wrap items-center gap-2.5">
          <div class="flex items-center gap-1.5 text-slate-700 font-bold text-xs md:text-sm">
            <i class="el-icon-sort text-amber-500 text-base"></i>
            <span>Sắp xếp Album FB (Mới nhất lên đầu):</span>
          </div>
          <el-select v-model="resyncCategoryCode" placeholder="Chọn loại SP..." size="mini" class="w-48">
            <el-option
              v-for="cat in albumCategories"
              :key="cat.code"
              :label="cat.name + ' (' + cat.available_count + ' SP)'"
              :value="cat.code">
            </el-option>
          </el-select>
          <el-button
            size="mini"
            type="warning"
            plain
            icon="el-icon-refresh"
            :loading="resyncLoading"
            :disabled="!resyncCategoryCode || isResyncActiveForCurrent"
            @click="handleResyncAlbum">
            Re-sync Album
          </el-button>
          <a v-if="currentAlbumCategory && currentAlbumCategory.album_url"
             :href="currentAlbumCategory.album_url"
             target="_blank"
             class="text-xs text-blue-600 hover:underline flex items-center gap-1">
            <i class="el-icon-link"></i> Xem Album
          </a>
        </div>
        <span class="text-[11px] text-slate-400 italic">Tự động chạy ngầm lúc 01:00 sáng (GMT+7) nếu có hàng mới</span>
      </div>

      <!-- Live Resync Queue List -->
      <div v-if="visibleResyncJobs.length" class="space-y-1.5 max-h-36 overflow-y-auto pr-1 custom-scrollbar">
        <div v-for="job in visibleResyncJobs" :key="job.id" class="p-2.5 rounded-lg border text-xs space-y-1.5"
          :class="{
            'bg-blue-50/70 border-blue-200 text-blue-900': ['PENDING','DELETING','UPLOADING'].includes(job.status),
            'bg-amber-50/80 border-amber-200 text-amber-900': job.status === 'PAUSED_LIMIT',
            'bg-emerald-50/70 border-emerald-200 text-emerald-900': job.status === 'COMPLETED',
            'bg-slate-50 border-slate-200 text-slate-600': ['CANCELLED','FAILED'].includes(job.status)
          }">
          <div class="flex flex-wrap items-center justify-between gap-2">
            <div class="font-bold flex items-center gap-1.5">
              <i v-if="['DELETING','UPLOADING'].includes(job.status)" class="el-icon-loading"></i>
              <i v-else-if="job.status === 'PENDING'" class="el-icon-more-outline"></i>
              <i v-else-if="job.status === 'PAUSED_LIMIT'" class="el-icon-time"></i>
              <i v-else-if="job.status === 'COMPLETED'" class="el-icon-circle-check"></i>
              <i v-else class="el-icon-info"></i>
              <span class="px-1.5 py-0.5 rounded bg-white/80 text-[10px] font-black uppercase tracking-wide border border-black/5">#{{ job.id }} {{ job.category_name || job.category_code }}</span>
              <span v-if="job.status === 'PENDING'">Đang chờ chạy (Xóa: {{ job.total_delete }}, Tải: {{ job.total_upload }})</span>
              <span v-else-if="job.status === 'DELETING'">Đang xóa ảnh cũ ({{ job.deleted_count }}/{{ job.total_delete }})</span>
              <span v-else-if="job.status === 'UPLOADING'">Đang tải lên Mới → Cũ ({{ job.uploaded_count }}/{{ job.total_upload }})</span>
              <span v-else-if="job.status === 'PAUSED_LIMIT'">Tạm nghỉ tránh giới hạn FB (Chạy tiếp lúc {{ formatResyncTime(job.resume_at) }})</span>
              <span v-else-if="job.status === 'COMPLETED'">Hoàn tất (Xóa {{ job.deleted_count }}/{{ job.total_delete }}, Tải {{ job.uploaded_count }}/{{ job.total_upload }})</span>
              <span v-else>{{ job.status }}</span>
            </div>
            <div class="flex items-center gap-1.5" v-if="['PENDING','DELETING','UPLOADING','PAUSED_LIMIT'].includes(job.status)">
              <el-button v-if="job.status === 'PAUSED_LIMIT'" size="mini" type="primary" plain @click="handleResumeResync(job.id)">Chạy tiếp ngay</el-button>
              <el-button size="mini" type="danger" plain @click="handleCancelResync(job.id)">Hủy</el-button>
            </div>
          </div>
          <el-progress
            v-if="['DELETING','UPLOADING','PAUSED_LIMIT'].includes(job.status)"
            :percentage="getJobProgressPct(job)"
            :status="job.status === 'PAUSED_LIMIT' ? 'warning' : undefined"
            :stroke-width="5">
          </el-progress>
          <div v-if="job.last_error" class="text-[11px] opacity-80">
            {{ job.last_error }}
          </div>
        </div>
      </div>
    </div>

    <div class="flex-1 flex gap-6 overflow-hidden min-h-0">
        <!-- Staging Area (Left) -->
        <div class="w-1/2 flex flex-col bg-white rounded-xl shadow-sm border border-slate-200">
            <div class="p-3 border-b border-slate-100 flex justify-between items-center bg-slate-50 rounded-t-xl">
                 <h3 class="font-bold text-slate-700 text-sm">Staging Area ({{ stagingFiles.length }})</h3>
                 <div v-if="selectedIds.length">
                     <span class="text-xs font-bold text-blue-600 mr-2">{{ selectedIds.length }} selected</span>
                     <el-button type="primary" size="mini" @click="createGroupFromSelection">Group as Product</el-button>
                 </div>
            </div>
            
            <div class="flex-1 overflow-y-auto p-4 custom-scrollbar">
                 <div v-if="!stagingFiles.length" class="h-full flex flex-col items-center justify-center text-slate-300">
                     <i class="el-icon-picture-outline text-6xl mb-2"></i>
                     <p>Drop images here or click Add Images</p>
                 </div>
                 
                 <div class="grid grid-cols-4 gap-2">
                     <div v-for="file in visibleStagingFiles" :key="file.id" 
                          class="aspect-square relative group cursor-pointer border rounded-lg overflow-hidden transition-all duration-200"
                          :class="selectedIds.includes(file.id) ? 'border-4 border-blue-500 ring-2 ring-blue-200' : 'border-slate-200 hover:border-slate-400'"
                          @click.exact="toggleSelect(file.id)"
                          @click.shift="rangeSelect(file.id)">
                         
                         <img :src="file.url" class="w-full h-full object-cover" loading="lazy">
                         
                         <!-- Selection Indicator -->
                         <div class="absolute top-1 right-1 w-5 h-5 rounded-full border-2 flex items-center justify-center"
                              :class="selectedIds.includes(file.id) ? 'bg-blue-500 border-blue-500' : 'bg-white/50 border-white'">
                              <i v-if="selectedIds.includes(file.id)" class="el-icon-check text-white text-xs font-bold"></i>
                         </div>
                     </div>
                 </div>
                 
                 <div v-if="hasMoreFiles" class="py-4 text-center">
                     <el-button @click="loadMore" size="small" round>Load More ({{ stagingFiles.length - renderLimit }} remaining)</el-button>
                 </div>
            </div>
        </div>

        <!-- Grouped Products (Right) -->
        <div class="w-1/2 flex flex-col bg-slate-100 rounded-xl border border-slate-200">
             <div class="p-3 border-b border-slate-200 flex justify-between items-center bg-white rounded-t-xl">
                 <div class="flex items-center gap-4">
                    <h3 class="font-bold text-slate-700 text-sm">Products to Create ({{ groups.length }})</h3>
                    <el-select v-model="batchType" size="mini" placeholder="Select Type" class="w-32">
                        <el-option label="SHIRT" value="SHIRT"></el-option>
                        <el-option label="TROUSERS" value="TROUSERS"></el-option>
                        <el-option label="BLAZER" value="BLAZER"></el-option>
                        <el-option label="GILE" value="GILE"></el-option>
                        <el-option label="SUIT" value="SUIT"></el-option>
                        <el-option label="POLO SHIRT" value="POLO SHIRT"></el-option>
                        <el-option label="JACKET" value="JACKET"></el-option>
                    </el-select>
                 </div>
                 <el-button type="text" class="text-red-500" @click="clearGroups" v-if="groups.length">Clear All</el-button>
             </div>

             <div v-if="batchComplete || successGroupCount > 0 || errorGroupCount > 0" class="px-4 py-2.5 bg-white border-b border-slate-200 flex flex-wrap items-center justify-between gap-2 text-sm font-bold">
                 <div class="flex items-center gap-3">
                     <span class="text-slate-600">Kết quả:</span>
                     <span class="text-green-600"><i class="el-icon-check"></i> {{ successGroupCount }} thành công</span>
                     <span class="text-red-500"><i class="el-icon-warning"></i> {{ errorGroupCount }} lỗi</span>
                 </div>
                 <div class="flex flex-wrap items-center gap-2">
                     <el-button v-if="batchComplete" size="mini" type="primary" plain icon="el-icon-s-grid" @click="summaryModalVisible = true">
                         Xem Summary
                     </el-button>
                     <el-button v-if="successGroupCount > 0" size="mini" type="info" plain icon="el-icon-delete" @click="clearSuccessGroups">
                         Clear Success ({{ successGroupCount }})
                     </el-button>
                     <el-button v-if="errorGroupCount > 0" type="danger" size="mini" :loading="uploading" @click="retryErrors">
                         Tải lại {{ errorGroupCount }} lỗi
                     </el-button>
                     <el-button
                         v-if="batchComplete && successGroupCount > 0 && canResyncBatchType"
                         type="warning"
                         size="mini"
                         icon="el-icon-sort"
                         :loading="resyncLoading"
                         @click="handleResyncBatchType">
                         Re-sync Album {{ batchType || resyncCategoryCode }}
                     </el-button>
                 </div>
             </div>

             <div class="flex-1 overflow-y-auto p-4 space-y-4 custom-scrollbar">
                 <div v-for="(group, idx) in groups" :key="idx" class="bg-white p-3 rounded-lg shadow-sm border border-slate-200 relative group-card">
                     <!-- Main Image Header -->
                     <div class="flex gap-3 mb-3">
                         <div class="w-20 h-20 flex-shrink-0 bg-slate-100 rounded border border-slate-200 relative overflow-hidden">
                             <img :src="group.main.url" class="w-full h-full object-cover">
                             <div class="absolute bottom-0 left-0 w-full bg-blue-600 text-white text-[10px] font-bold text-center py-0.5">MAIN</div>
                             
                             <!-- Status Overlay -->
                             <div v-if="group.status && group.status !== 'PENDING'" class="absolute inset-0 bg-white/90 flex flex-col items-center justify-center p-1 text-center">
                                 <i v-if="group.status === 'PROCESSING' || group.status === 'UPLOADING'" class="el-icon-loading text-blue-500 font-bold"></i>
                                 <i v-if="group.status === 'SUCCESS'" class="el-icon-check text-green-500 font-bold text-lg"></i>
                                 <i v-if="group.status === 'ERROR'" class="el-icon-warning text-red-500 font-bold text-lg"></i>
                             </div>
                         </div>
                         <div class="flex-1 min-w-0">
                             <div class="flex justify-between items-start">
                                 <h4 class="font-bold text-slate-800 text-sm truncate">
                                     {{ group.productName ? group.productName : `Product #${idx + 1}` }}
                                 </h4>
                                 <el-button v-if="!group.status || group.status === 'PENDING'" type="text" icon="el-icon-close" class="text-slate-400 -mt-2 -mr-2" @click="disbandGroup(idx)"></el-button>
                             </div>
                             
                             <!-- Status Text -->
                             <div class="mt-1">
                                 <p v-if="group.status === 'ERROR'" class="text-xs text-red-600 font-bold leading-tight">{{ group.error }}</p>
                                 <p v-else-if="group.status === 'SUCCESS'" class="text-xs text-green-600 font-bold">Uploaded Successfully</p>
                                 <p v-else-if="group.status === 'PROCESSING'" class="text-xs text-blue-500 font-bold">Processing...</p>
                                 <p v-else class="text-xs text-slate-500">Will extract info from Main Image.</p>
                             </div>

                             <div class="mt-2 text-xs font-bold text-slate-600 bg-slate-50 inline-block px-2 py-1 rounded">
                                 + {{ group.details.length }} Detail Images
                             </div>
                         </div>
                     </div>

                     <!-- Detail Images Strip -->
                     <div class="flex gap-2 overflow-x-auto pb-2 scrollbar-thin">
                         <div v-for="detail in group.details" :key="detail.id" 
                              class="w-12 h-12 flex-shrink-0 relative cursor-pointer border border-transparent hover:border-blue-400 rounded overflow-hidden"
                              title="Click to set as Main"
                              @click="setAsMain(idx, detail)">
                             <img :src="detail.url" class="w-full h-full object-cover opacity-70 hover:opacity-100">
                         </div>
                     </div>
                 </div>
             </div>
        </div>
    </div>

    <!-- Upload Summary Modal -->
    <el-dialog
      :visible.sync="summaryModalVisible"
      title="Kết quả Upload sản phẩm"
      width="860px"
      top="5vh"
      custom-class="rounded-2xl">
      <div class="space-y-4">
        <!-- Summary Header Stats + Actions -->
        <div class="flex flex-wrap items-center justify-between gap-3 p-3 bg-slate-50 rounded-xl border border-slate-200">
          <div class="flex items-center gap-4 text-sm font-bold">
            <span class="text-emerald-600 flex items-center gap-1">
              <i class="el-icon-circle-check text-base"></i> Thành công: {{ successGroupCount }}
            </span>
            <span class="text-red-500 flex items-center gap-1">
              <i class="el-icon-warning text-base"></i> Thất bại: {{ errorGroupCount }}
            </span>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <el-button v-if="errorGroupCount > 0" type="danger" size="small" :loading="uploading" @click="summaryModalVisible = false; retryErrors()">
              Tải lại {{ errorGroupCount }} lỗi
            </el-button>
            <el-button
              v-if="successGroupCount > 0 && canResyncBatchType"
              type="warning"
              size="small"
              icon="el-icon-sort"
              :loading="resyncLoading"
              @click="handleResyncBatchType">
              Re-sync Album {{ batchType || resyncCategoryCode }}
            </el-button>
          </div>
        </div>

        <div class="max-h-[60vh] overflow-y-auto space-y-5 pr-1 custom-scrollbar">
          <!-- Failed Products Grid -->
          <div v-if="errorGroups.length" class="space-y-2">
            <div class="text-xs font-black uppercase tracking-wider text-red-600 flex items-center gap-1.5">
              <i class="el-icon-warning"></i> Sản phẩm lỗi ({{ errorGroups.length }})
            </div>
            <div class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 gap-2.5">
              <div v-for="(group, idx) in errorGroups" :key="'m-err-' + idx"
                   class="bg-white rounded-lg border-2 border-red-200 overflow-hidden shadow-sm flex flex-col">
                <div class="aspect-square relative bg-slate-100">
                  <img :src="group.main.url" class="w-full h-full object-cover">
                  <span class="absolute top-1 right-1 px-1.5 py-0.5 rounded bg-red-600 text-white text-[10px] font-bold">LỖI</span>
                </div>
                <div class="p-2 space-y-0.5">
                  <div class="font-bold text-xs text-slate-800 truncate">{{ group.productName || 'Chưa rõ mã' }}</div>
                  <div class="text-[10px] text-red-600 leading-tight line-clamp-2" :title="group.error">{{ group.error }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Successful Products Grid -->
          <div v-if="successGroups.length" class="space-y-2">
            <div class="text-xs font-black uppercase tracking-wider text-emerald-700 flex items-center gap-1.5">
              <i class="el-icon-circle-check"></i> Đã tạo thành công ({{ successGroups.length }})
            </div>
            <div class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 gap-2.5">
              <div v-for="(group, idx) in successGroups" :key="'m-ok-' + idx"
                   class="bg-white rounded-lg border border-emerald-200 overflow-hidden shadow-sm flex flex-col">
                <div class="aspect-square relative bg-slate-100">
                  <img :src="group.main.url" class="w-full h-full object-cover">
                  <span class="absolute top-1 right-1 w-5 h-5 rounded-full bg-emerald-500 text-white flex items-center justify-center text-xs font-bold shadow">
                    <i class="el-icon-check"></i>
                  </span>
                  <span v-if="group.details.length" class="absolute bottom-1 left-1 px-1.5 py-0.5 rounded bg-black/60 text-white text-[10px] font-bold">
                    +{{ group.details.length }} ảnh
                  </span>
                </div>
                <div class="p-2">
                  <div class="font-bold text-xs text-slate-800 truncate text-center">{{ group.productName || 'Đã lưu' }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <span slot="footer" class="flex justify-between items-center">
        <el-button v-if="successGroupCount > 0" type="info" plain size="small" icon="el-icon-delete" @click="clearSuccessGroups(); summaryModalVisible = false">
          Clear Success ({{ successGroupCount }}) & Đóng
        </el-button>
        <div v-else></div>
        <el-button size="small" @click="summaryModalVisible = false">Đóng</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
export default {
    data() {
        return {
            stagingFiles: [],
            selectedIds: [],
            groups: [], // Array of { main, details, sessionId, status, error, productName }
            uploading: false,
            batchType: '',
            lastSelectedId: null,
            renderLimit: 48,
            pollInterval: null,
            albumCategories: [],
            resyncCategoryCode: '',
            resyncJobs: [],
            resyncLoading: false,
            resyncPollTimer: null,
            summaryModalVisible: false
        }
    },
    computed: {
        currentAlbumCategory() {
            return this.albumCategories.find(c => c.code === this.resyncCategoryCode) || null;
        },
        visibleResyncJobs() {
            if (!this.resyncJobs.length) return [];
            const active = this.resyncJobs
                .filter(j => ['PENDING', 'DELETING', 'UPLOADING', 'PAUSED_LIMIT'].includes(j.status))
                .slice()
                .sort((a, b) => a.id - b.id);
            if (active.length > 0) return active;
            // If no active jobs, show the latest completed/cancelled job for selected category (or most recent overall)
            const selected = this.resyncCategoryCode
                ? this.resyncJobs.find(j => j.category_code === this.resyncCategoryCode)
                : null;
            return selected ? [selected] : [this.resyncJobs[0]];
        },
        isResyncActiveForCurrent() {
            return this.resyncJobs.some(j =>
                j.category_code === this.resyncCategoryCode &&
                ['PENDING', 'DELETING', 'UPLOADING', 'PAUSED_LIMIT'].includes(j.status)
            );
        },
        visibleStagingFiles() {
            return this.stagingFiles.slice(0, this.renderLimit);
        },
        hasMoreFiles() {
            return this.renderLimit < this.stagingFiles.length;
        },
        batchComplete() {
            return this.groups.length > 0 && !this.pollInterval &&
                this.groups.some(g => Boolean(g.sessionId)) &&
                this.groups.every(g => !g.sessionId || g.status === 'SUCCESS' || g.status === 'ERROR');
        },
        successGroups() {
            return this.groups.filter(g => g.status === 'SUCCESS');
        },
        errorGroups() {
            return this.groups.filter(g => g.status === 'ERROR');
        },
        successGroupCount() {
            return this.successGroups.length;
        },
        errorGroupCount() {
            return this.errorGroups.length;
        },
        canResyncBatchType() {
            const target = this.batchType || this.resyncCategoryCode;
            return this.albumCategories.some(c => c.code === target);
        }
    },
    mounted() {
        this.fetchAlbumCategories();
        this.fetchResyncStatus();
    },
    beforeDestroy() {
        if (this.pollInterval) clearInterval(this.pollInterval);
        if (this.resyncPollTimer) {
            clearInterval(this.resyncPollTimer);
            this.resyncPollTimer = null;
        }
    },
    methods: {
        async fetchAlbumCategories() {
            try {
                const res = await axios.get('/salev2/api/facebook/post-init');
                this.albumCategories = (res.data.categories || []).filter(c => c.fb_album_id);
                if (this.albumCategories.length > 0 && !this.resyncCategoryCode) {
                    this.resyncCategoryCode = this.albumCategories[0].code;
                }
            } catch (e) {
                // silent fail
            }
        },
        async fetchResyncStatus() {
            try {
                const res = await axios.get('/salev2/api/facebook/albums/resync-status');
                this.resyncJobs = res.data.jobs || [];
                const hasActive = this.resyncJobs.some(j => ['PENDING', 'DELETING', 'UPLOADING', 'PAUSED_LIMIT'].includes(j.status));
                if (hasActive && !this.resyncPollTimer) {
                    this.resyncPollTimer = setInterval(() => this.fetchResyncStatus(), 5000);
                } else if (!hasActive && this.resyncPollTimer) {
                    clearInterval(this.resyncPollTimer);
                    this.resyncPollTimer = null;
                }
            } catch (e) {
                // silent fail
            }
        },
        async handleResyncBatchType() {
            const target = this.batchType || this.resyncCategoryCode;
            if (!target) return;
            this.resyncCategoryCode = target;
            await this.$nextTick();
            await this.handleResyncAlbum();
        },
        async handleResyncAlbum() {
            if (!this.resyncCategoryCode) return;
            try {
                await this.$confirm(
                    `Hệ thống sẽ xóa toàn bộ ảnh hiện tại trong Album "${this.currentAlbumCategory?.name}" và tải lại theo thứ tự Mới nhất → Cũ nhất.\n\nNếu chạm giới hạn API của Facebook, tiến trình sẽ tự động tạm nghỉ 45 phút rồi chạy tiếp cho tới khi hoàn tất.`,
                    'Xác nhận Sắp xếp lại Album',
                    { confirmButtonText: 'Bắt đầu chạy', cancelButtonText: 'Đóng', type: 'warning' }
                );
            } catch (e) {
                return;
            }
            this.resyncLoading = true;
            try {
                await axios.post('/salev2/api/facebook/albums/resync', {
                    category_code: this.resyncCategoryCode
                });
                this.$message.success('Đã kích hoạt tiến trình sắp xếp lại Album!');
                await this.fetchResyncStatus();
            } catch (err) {
                this.$message.error('Lỗi kích hoạt: ' + (err.response?.data?.error || err.message));
            } finally {
                this.resyncLoading = false;
            }
        },
        async handleCancelResync(jobId) {
            try {
                await axios.post(`/salev2/api/facebook/albums/resync/${jobId}/cancel`);
                this.$message.info('Đã hủy tiến trình sắp xếp Album');
                await this.fetchResyncStatus();
            } catch (err) {
                this.$message.error('Lỗi khi hủy: ' + (err.response?.data?.error || err.message));
            }
        },
        async handleResumeResync(jobId) {
            try {
                await axios.post(`/salev2/api/facebook/albums/resync/${jobId}/resume`);
                this.$message.success('Đã kích hoạt chạy tiếp ngay!');
                await this.fetchResyncStatus();
            } catch (err) {
                this.$message.error('Lỗi khi chạy tiếp: ' + (err.response?.data?.error || err.message));
            }
        },
        getJobProgressPct(job) {
            if (!job) return 0;
            const total = (job.total_delete || 0) + (job.total_upload || 0);
            if (total === 0) return 0;
            const done = (job.deleted_count || 0) + (job.uploaded_count || 0);
            return Math.min(100, Math.round((done / total) * 100));
        },
        formatResyncTime(iso) {
            if (!iso) return '...';
            const d = new Date(iso);
            return d.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' });
        },
        loadMore() {
            this.renderLimit += 48;
        },
        handleFiles(e) {
            const files = Array.from(e.target.files);
            const baseTime = Date.now();
            
            files.forEach((file, index) => {
                const item = {
                    id: Math.random().toString(36).substr(2, 9),
                    file: file,
                    url: '', 
                    originalIndex: baseTime + index,
                    loading: true
                };
                this.stagingFiles.push(item);
                this.createThumbnail(file).then(thumbUrl => {
                    item.url = thumbUrl;
                    item.loading = false;
                });
            });
            e.target.value = '';
        },
        createThumbnail(file) {
            return new Promise((resolve) => {
                const img = new Image();
                const url = URL.createObjectURL(file);
                img.src = url;
                img.onload = () => {
                    const canvas = document.createElement('canvas');
                    const ctx = canvas.getContext('2d');
                    const MAX_WIDTH = 300;
                    let width = img.width;
                    let height = img.height;
                    
                    if (width > MAX_WIDTH) {
                        height *= MAX_WIDTH / width;
                        width = MAX_WIDTH;
                    }
                    
                    canvas.width = width;
                    canvas.height = height;
                    ctx.drawImage(img, 0, 0, width, height);
                    
                    URL.revokeObjectURL(url);
                    resolve(canvas.toDataURL('image/jpeg', 0.7));
                };
                img.onerror = () => {
                    resolve(url);
                };
            });
        },
        toggleSelect(id) {
            const idx = this.selectedIds.indexOf(id);
            if (idx > -1) {
                this.selectedIds.splice(idx, 1);
            } else {
                this.selectedIds.push(id);
            }
            this.lastSelectedId = id;
        },
        rangeSelect(id) {
            if (!this.lastSelectedId || this.stagingFiles.length === 0) {
                this.toggleSelect(id);
                return;
            }
            const startIdx = this.stagingFiles.findIndex(f => f.id === this.lastSelectedId);
            const endIdx = this.stagingFiles.findIndex(f => f.id === id);
            if (startIdx === -1 || endIdx === -1) return;
            
            const min = Math.min(startIdx, endIdx);
            const max = Math.max(startIdx, endIdx);
            const rangeIds = this.stagingFiles.slice(min, max + 1).map(f => f.id);
            this.selectedIds = [...new Set([...this.selectedIds, ...rangeIds])];
        },
        createGroupFromSelection() {
            if (!this.selectedIds.length) return;
            const selectedFiles = this.stagingFiles.filter(f => this.selectedIds.includes(f.id));
            if (selectedFiles.length === 0) return;

            const main = selectedFiles[0];
            const details = selectedFiles.slice(1);
            
            this.groups.push({ main, details, status: 'PENDING' });
            this.stagingFiles = this.stagingFiles.filter(f => !this.selectedIds.includes(f.id));
            this.selectedIds = [];
            this.lastSelectedId = null;
        },
        disbandGroup(idx) {
            const group = this.groups[idx];
            this.stagingFiles.push(group.main, ...group.details);
            this.stagingFiles.sort((a, b) => a.originalIndex - b.originalIndex);
            this.groups.splice(idx, 1);
        },
        setAsMain(groupIdx, newMain) {
            const group = this.groups[groupIdx];
            const oldMain = group.main;
            group.details = group.details.filter(d => d.id !== newMain.id);
            group.details.push(oldMain);
            group.main = newMain;
        },
        clearGroups() {
            this.groups.forEach((g, idx) => this.disbandGroup(idx));
        },
        clearSuccessGroups() {
            const count = this.successGroupCount;
            this.groups = this.groups.filter(g => g.status !== 'SUCCESS');
            if (count > 0) {
                this.$message.success(`Đã dọn ${count} sản phẩm thành công khỏi danh sách`);
            }
        },
        startPolling() {
            if (this.pollInterval) return;
            this.pollInterval = setInterval(this.checkStatus, 2000);
        },
        detectBatchDuplicates() {
            const nameCounts = {};
            this.groups.forEach(g => {
                if (g.productName) {
                    nameCounts[g.productName] = (nameCounts[g.productName] || 0) + 1;
                }
            });

            this.groups.forEach(g => {
                if (g.productName && nameCounts[g.productName] > 1) {
                    if (g.status === 'SUCCESS') {
                        g.status = 'ERROR';
                        g.error = `Batch Duplicate: ${g.productName}`;
                    }
                }
            });
        },
        async checkStatus() {
            const pendingGroups = this.groups.filter(g => g.sessionId && g.status !== 'SUCCESS' && g.status !== 'ERROR');
            
            if (pendingGroups.length) {
                const sessionIds = pendingGroups.map(g => g.sessionId);
                 try {
                    const res = await axios.get('/salev2/api/upload-logs', { params: { session_ids: sessionIds.join(',') } });
                    const logs = res.data;
                    
                    logs.forEach(log => {
                        const group = this.groups.find(g => g.sessionId === log.session_id && g.main.file.name === log.file_name);
                        if (group) {
                            group.status = (log.status === 'FAILED' ? 'ERROR' : log.status); 
                            group.error = log.message;
                            group.productName = log.product_code;
                            group.productId = log.product_id;
                        }
                    });
                } catch (e) {
                    console.error("Polling error", e);
                }
            }

            // Detect completion
            const stillPending = this.groups.filter(g => g.sessionId && g.status !== 'SUCCESS' && g.status !== 'ERROR');
            if (!stillPending.length) {
                clearInterval(this.pollInterval);
                this.pollInterval = null;
                // Run post-upload verification
                this.detectBatchDuplicates();
                this.summaryModalVisible = true;
            }
        },
        async uploadAll() {
            if (!this.groups.length) return;
            if (!this.batchType) {
                this.$message.error('Please select a Product Type first');
                return;
            }
            this.uploading = true;
            
            // Unified ID for the batch
            const batchId = 'BATCH_' + Date.now() + '_' + Math.random().toString(36).substr(2, 5);
            
            try {
                for (let i = 0; i < this.groups.length; i++) {
                     if (this.groups[i].status === 'SUCCESS') continue;

                     let gData = new FormData();
                     gData.append('main', this.groups[i].main.file);
                     this.groups[i].details.forEach((d, j) => {
                         gData.append(`details[${j}]`, d.file);
                     });
                     
                     // Pass batch ID and selected type
                     gData.append('session_id', batchId);
                     gData.append('type', this.batchType);
                     
                     this.groups[i].status = 'UPLOADING';

                     try {
                         const res = await axios.post('/salev2/api/upload-smart-group', gData, {
                             headers: { 'Content-Type': 'multipart/form-data' }
                         });
                         // Save individual group tracking ID
                         this.$set(this.groups[i], 'sessionId', res.data.session_id);
                         this.$set(this.groups[i], 'status', 'PROCESSING');
                     } catch (e) {
                         this.groups[i].status = 'ERROR';
                         this.groups[i].error = e.response && e.response.data && e.response.data.message 
                             ? e.response.data.message 
                             : (e.message || 'Network Error');
                     }

                     // Small pause between sequential uploads to avoid overwhelming memory and network
                     await new Promise(r => setTimeout(r, 400));
                }
                this.startPolling();
            } catch (e) {
                console.error(e);
            } finally {
                this.uploading = false; 
            }
        },
        retryErrors() {
            this.groups.forEach(g => {
                if (g.status === 'ERROR') {
                    g.status = 'PENDING';
                    g.error = null;
                    g.sessionId = null;
                }
            });
            this.uploadAll();
        }
    }
}
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
    width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
    background: #f1f5f9; 
}
.custom-scrollbar::-webkit-scrollbar-thumb {
    background: #cbd5e1; 
    border-radius: 3px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: #94a3b8; 
}
</style>
