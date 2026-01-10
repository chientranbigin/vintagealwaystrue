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
                        <el-option label="T-SHIRT" value="T-SHIRT"></el-option>
                        <el-option label="TROUSERS" value="TROUSERS"></el-option>
                        <el-option label="BLAZER" value="BLAZER"></el-option>
                        <el-option label="GILE" value="GILE"></el-option>
                    </el-select>
                 </div>
                 <el-button type="text" class="text-red-500" @click="clearGroups" v-if="groups.length">Clear All</el-button>
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
            batchType: 'T-SHIRT',
            lastSelectedId: null,
            renderLimit: 48,
            pollInterval: null
        }
    },
    computed: {
        visibleStagingFiles() {
            return this.stagingFiles.slice(0, this.renderLimit);
        },
        hasMoreFiles() {
            return this.renderLimit < this.stagingFiles.length;
        }
    },
    beforeDestroy() {
        if (this.pollInterval) clearInterval(this.pollInterval);
    },
    methods: {
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
                        const group = this.groups.find(g => g.sessionId === log.session_id);
                        if (group) {
                            group.status = log.status; // PROCESSING, SUCCESS, ERROR
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
            }
        },
        async uploadAll() {
            if (!this.groups.length) return;
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
                         this.groups[i].error = 'Network Error';
                     }
                }
                this.startPolling();
            } catch (e) {
                console.error(e);
            } finally {
                this.uploading = false; 
            }
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
