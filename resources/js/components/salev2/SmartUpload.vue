<template>
  <div class="px-6 py-8">
    <div class="mb-8 mobile-hide">
      <h1 class="text-3xl font-bold text-gray-900 tracking-tight mb-2">Batch Image Upload</h1>
      <p class="text-slate-500 italic">Queue-based processing for massive product catalogs</p>
    </div>

    <!-- Legacy Style Controls -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
      <div class="bg-white rounded-2xl shadow-sm border border-slate-100 p-6">
        <div class="space-y-6">
          <div>
            <label class="text-xs font-black text-slate-400 uppercase tracking-widest block mb-2">Product Type</label>
            <el-select v-model="productType" placeholder="Select Type" class="w-full shadow-sm">
              <el-option v-for="type in productTypes" :key="type" :label="type" :value="type"></el-option>
            </el-select>
          </div>

          <div>
            <label class="text-xs font-black text-slate-400 uppercase tracking-widest block mb-2">Select Images (Multiple)</label>
            <div class="flex items-center gap-4">
              <el-button type="info" plain icon="el-icon-files" @click="triggerFileSelect">
                {{ selectedFiles.length ? `${selectedFiles.length} files selected` : 'Choose Files' }}
              </el-button>
              <input type="file" ref="fileInput" multiple hidden @change="handleFileChange" accept="image/*">
              <el-button type="primary" :disabled="!selectedFiles.length || uploading" 
                         :loading="uploading" icon="el-icon-upload" @click="startUpload">
                Start Uploading
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <!-- Progress Card -->
      <div class="bg-slate-900 rounded-2xl p-6 shadow-xl text-white relative overflow-hidden">
        <div class="relative z-10 flex flex-col h-full justify-between">
          <div>
            <p class="text-[10px] font-black opacity-40 uppercase tracking-[4px] mb-4">Upload Progress</p>
            <div class="flex items-end gap-2 mb-2">
                <span class="text-4xl font-black">{{ progress }}%</span>
                <span class="text-xs opacity-60 mb-1">Session: {{ sessionId || '---' }}</span>
            </div>
            <el-progress :percentage="progress" :show-text="false" color="#3b82f6" class="mb-6"></el-progress>
          </div>
          
          <div class="grid grid-cols-4 gap-4 text-center">
            <div><p class="text-xl font-black text-blue-400">{{ totalCount }}</p><p class="text-[10px] opacity-40 font-bold uppercase">Total</p></div>
            <div><p class="text-xl font-black text-amber-400">{{ processingCount }}</p><p class="text-[10px] opacity-40 font-bold uppercase">Proc</p></div>
            <div><p class="text-xl font-black text-emerald-400">{{ successCount }}</p><p class="text-[10px] opacity-40 font-bold uppercase">Done</p></div>
            <div><p class="text-xl font-black text-rose-400">{{ failedCount }}</p><p class="text-[10px] opacity-40 font-bold uppercase">Fail</p></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Logs Split View -->
    <div class="grid grid-cols-1 xl:grid-cols-2 gap-8">
      <!-- Upload Log -->
      <div class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
        <div class="bg-slate-50 px-4 py-3 border-b border-slate-100 flex justify-between items-center">
            <h3 class="text-xs font-black text-slate-400 uppercase tracking-widest">Network Upload Log</h3>
            <span class="text-[10px] bg-blue-100 text-blue-600 px-2 py-0.5 rounded-full font-bold">{{ uploadLog.length }} Files</span>
        </div>
        <div class="max-h-[500px] overflow-y-auto">
            <el-table :data="uploadLog" size="mini" :row-class-name="logRowClass">
                <el-table-column prop="index" label="#" width="50" align="center"></el-table-column>
                <el-table-column prop="filename" label="Filename" min-width="150" show-overflow-tooltip></el-table-column>
                <el-table-column label="Status" width="100" align="center">
                    <template slot-scope="scope">
                        <el-tag :type="scope.row.status === 'success' ? 'success' : (scope.row.status === 'uploading' ? 'info' : 'danger')" 
                                size="mini" effect="dark" class="uppercase font-black text-[8px] rounded-pill">
                            {{ scope.row.status }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="message" label="Message" min-width="200" show-overflow-tooltip></el-table-column>
            </el-table>
        </div>
      </div>

      <!-- Processing Log -->
      <div class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
        <div class="bg-slate-50 px-4 py-3 border-b border-slate-100 flex justify-between items-center">
            <h3 class="text-xs font-black text-slate-400 uppercase tracking-widest">Backend Processing Log</h3>
            <span class="animate-pulse text-[10px] bg-amber-100 text-amber-600 px-2 py-0.5 rounded-full font-bold" v-if="processingCount > 0">Processing...</span>
        </div>
        <div class="max-h-[500px] overflow-y-auto">
            <el-table :data="processingLog" size="mini">
                <el-table-column type="index" label="#" width="50" align="center"></el-table-column>
                <el-table-column prop="file" label="File" min-width="150" show-overflow-tooltip></el-table-column>
                <el-table-column label="Status" width="100" align="center">
                    <template slot-scope="scope">
                        <el-tag :type="getProcStatusType(scope.row.status)" size="mini" effect="dark" class="uppercase font-black text-[8px] rounded-pill">
                            {{ scope.row.status }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="code" label="Code" width="100" align="center"></el-table-column>
                <el-table-column prop="reason" label="Result/Error" min-width="200" show-overflow-tooltip></el-table-column>
            </el-table>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'SmartUpload',
  data() {
    return {
      productType: 'TROUSERS',
      productTypes: ['TROUSERS', 'JACKET', 'SHIRT', 'BLAZER', 'TIE', 'GILE', 'BELT', 'POLO SHIRT', 'HAT', 'SUIT'],
      selectedFiles: [],
      uploading: false,
      progress: 0,
      sessionId: null,
      totalCount: 0,
      processingCount: 0,
      successCount: 0,
      failedCount: 0,
      uploadLog: [],
      processingLog: [],
      statusCheckInterval: null
    }
  },
  destroyed() {
    if (this.statusCheckInterval) clearInterval(this.statusCheckInterval);
  },
  methods: {
    triggerFileSelect() { this.$refs.fileInput.click(); },
    handleFileChange(e) {
      this.selectedFiles = Array.from(e.target.files);
      this.totalCount = this.selectedFiles.length;
    },
    generateSessionId() {
      return 'upl_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
    },
    async startUpload() {
      this.sessionId = this.generateSessionId();
      this.uploading = true;
      this.progress = 0;
      this.uploadLog = [];
      this.processingLog = [];
      this.successCount = 0;
      this.failedCount = 0;
      this.processingCount = 0;

      this.startStatusCheck();

      for (let i = 0; i < this.selectedFiles.length; i++) {
        await this.uploadSingleFile(this.selectedFiles[i], i + 1);
        this.progress = Math.round(((i + 1) / this.selectedFiles.length) * 100);
      }
      
      this.uploading = false;
    },
    async uploadSingleFile(file, index) {
      const formData = new FormData();
      formData.append('image', file);
      formData.append('type', this.productType);
      formData.append('session_id', this.sessionId);

      const logItem = { index, filename: file.name, status: 'uploading', message: 'Sending file...' };
      this.uploadLog.unshift(logItem);

      try {
        const res = await axios.post('/sale/product-batch-upload/single', formData, {
          headers: { 'Content-Type': 'multipart/form-data' }
        });
        
        if (res.data.success) {
          logItem.status = 'success';
          logItem.message = 'Queued for processing';
        } else {
          logItem.status = 'failed';
          logItem.message = res.data.message || 'Upload failed';
        }
      } catch (err) {
        logItem.status = 'failed';
        logItem.message = err.response?.data?.message || err.message;
      }
    },
    startStatusCheck() {
      if (this.statusCheckInterval) clearInterval(this.statusCheckInterval);
      
      this.statusCheckInterval = setInterval(async () => {
        try {
          const res = await axios.get('/sale/product-batch-upload/status', {
            params: { session_id: this.sessionId }
          });
          
          this.successCount = res.data.success;
          this.failedCount = res.data.failed;
          this.processingCount = res.data.processing;
          this.processingLog = res.data.logs;

          if (this.processingCount === 0 && res.data.total >= this.totalCount && !this.uploading) {
            clearInterval(this.statusCheckInterval);
          }
        } catch (err) {
          console.error('Status check failed:', err);
        }
      }, 2000);
    },
    logRowClass({row}) {
        return row.status === 'success' ? 'success-row' : (row.status === 'failed' ? 'danger-row' : '');
    },
    getProcStatusType(s) {
        return { 'SUCCESS': 'success', 'FAILED': 'danger', 'ERROR': 'danger', 'PENDING': 'warning' }[s] || 'info';
    }
  }
}
</script>

<style>
.el-table .success-row { background: #f0f9eb; }
.el-table .danger-row { background: #fef0f0; }
</style>
