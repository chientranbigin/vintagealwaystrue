@extends('sale.master')
@section('content')
    <div class="main-content">
        <div class="page-content">
            <div class="container-fluid">
                <div class="card">
                    <div class="card-body">
                        <h4 class="card-title mb-4">Batch Upload Product Images (Queue-based)</h4>
                        
                        <div class="alert alert-info">
                            <strong>Hướng dẫn:</strong>
                            <ul class="mb-0">
                                <li>Upload được nhiều ảnh cùng lúc (không giới hạn số lượng)</li>
                                <li>Mỗi ảnh tối đa 10MB</li>
                                <li>Hệ thống sẽ xử lý từng ảnh một cách bất đồng bộ</li>
                                <li>Bạn có thể theo dõi tiến trình upload và xử lý realtime</li>
                            </ul>
                        </div>

                        <div class="row">
                            <div class="col-md-6">
                                <div class="form-group">
                                    <label for="productType">Product Type</label>
                                    <select class="form-control" id="productType">
                                        @foreach(App\Contracts\SaleConfig::PRODUCT_TYPES as $type)
                                            <option value="{{$type}}">{{$type}}</option>
                                        @endforeach
                                    </select>
                                </div>

                                <div class="form-group">
                                    <label for="fileInput">Select Images</label>
                                    <div class="custom-file">
                                        <input multiple type="file" class="custom-file-input" id="fileInput" accept="image/*">
                                        <label class="custom-file-label" for="fileInput">Choose files...</label>
                                    </div>
                                </div>

                                <button type="button" class="btn btn-primary btn-lg" id="startUpload">
                                    <i class="mdi mdi-upload"></i> Start Upload
                                </button>
                            </div>

                            <div class="col-md-6">
                                <div class="card bg-light">
                                    <div class="card-body">
                                        <h5 class="card-title">Upload Progress</h5>
                                        <div class="mb-3">
                                            <strong>Session ID:</strong> <span id="sessionId">-</span>
                                        </div>
                                        <div class="mb-3">
                                            <div class="progress" style="height: 30px;">
                                                <div class="progress-bar progress-bar-striped progress-bar-animated" 
                                                     role="progressbar" 
                                                     id="uploadProgress" 
                                                     style="width: 0%">0%</div>
                                            </div>
                                        </div>
                                        <div class="row text-center">
                                            <div class="col-3">
                                                <h3 class="text-primary" id="totalCount">0</h3>
                                                <p class="text-muted mb-0">Total</p>
                                            </div>
                                            <div class="col-3">
                                                <h3 class="text-warning" id="processingCount">0</h3>
                                                <p class="text-muted mb-0">Processing</p>
                                            </div>
                                            <div class="col-3">
                                                <h3 class="text-success" id="successCount">0</h3>
                                                <p class="text-muted mb-0">Success</p>
                                            </div>
                                            <div class="col-3">
                                                <h3 class="text-danger" id="failedCount">0</h3>
                                                <p class="text-muted mb-0">Failed</p>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <!-- Upload Log -->
                        <div class="mt-4">
                            <h5>Upload Log</h5>
                            <div class="table-responsive" style="max-height: 400px; overflow-y: auto;">
                                <table class="table table-sm table-bordered">
                                    <thead class="thead-light">
                                        <tr>
                                            <th width="5%">#</th>
                                            <th width="25%">File</th>
                                            <th width="15%">Status</th>
                                            <th width="15%">Product Code</th>
                                            <th width="40%">Message</th>
                                        </tr>
                                    </thead>
                                    <tbody id="uploadLog">
                                        <tr>
                                            <td colspan="5" class="text-center text-muted">No uploads yet</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>

                        <!-- Processing Log (from queue) -->
                        <div class="mt-4">
                            <h5>Processing Log (Real-time)</h5>
                            <div class="table-responsive" style="max-height: 400px; overflow-y: auto;">
                                <table class="table table-sm table-bordered">
                                    <thead class="thead-light">
                                        <tr>
                                            <th width="5%">#</th>
                                            <th width="25%">File</th>
                                            <th width="15%">Status</th>
                                            <th width="15%">Product Code</th>
                                            <th width="40%">Reason</th>
                                        </tr>
                                    </thead>
                                    <tbody id="processingLog">
                                        <tr>
                                            <td colspan="5" class="text-center text-muted">Waiting for processing...</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <style>
        .log-success { background-color: #d4edda; }
        .log-failed { background-color: #f8d7da; }
        .log-uploading { background-color: #d1ecf1; }
    </style>

    <script>
        let selectedFiles = [];
        let sessionId = null;
        let uploadedCount = 0;
        let totalFiles = 0;
        let statusCheckInterval = null;

        // Generate session ID
        function generateSessionId() {
            return 'upload_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
        }

        // File input change
        document.getElementById('fileInput').addEventListener('change', function(e) {
            selectedFiles = Array.from(e.target.files);
            const label = document.querySelector('.custom-file-label');
            label.textContent = selectedFiles.length + ' file(s) selected';
            document.getElementById('totalCount').textContent = selectedFiles.length;
        });

        // Start upload
        document.getElementById('startUpload').addEventListener('click', async function() {
            if (selectedFiles.length === 0) {
                alert('Please select files first!');
                return;
            }

            // Generate session ID
            sessionId = generateSessionId();
            document.getElementById('sessionId').textContent = sessionId;

            totalFiles = selectedFiles.length;
            uploadedCount = 0;

            // Clear logs
            document.getElementById('uploadLog').innerHTML = '';
            document.getElementById('processingLog').innerHTML = '';

            // Disable button
            this.disabled = true;
            this.innerHTML = '<i class="mdi mdi-loading mdi-spin"></i> Uploading...';

            // Start checking processing status immediately
            startStatusCheck();

            // Upload each file
            for (let i = 0; i < selectedFiles.length; i++) {
                await uploadFile(selectedFiles[i], i + 1);
            }
        });

        // Upload single file
        async function uploadFile(file, index) {
            const formData = new FormData();
            formData.append('image', file);
            formData.append('type', document.getElementById('productType').value);
            formData.append('session_id', sessionId);

            // Add to upload log
            addUploadLog(index, file.name, 'uploading', '-', 'Uploading...');

            try {
                const response = await fetch('{{ route("sale.product.batch.upload.single") }}', {
                    method: 'POST',
                    headers: {
                        'X-CSRF-TOKEN': '{{ csrf_token() }}'
                    },
                    body: formData
                });

                const result = await response.json();

                if (result.success) {
                    updateUploadLog(index, 'success', '-', 'Uploaded, queued for processing');
                } else {
                    updateUploadLog(index, 'failed', '-', result.message);
                }
            } catch (error) {
                updateUploadLog(index, 'failed', '-', error.message);
            }

            uploadedCount++;
            updateProgress();
        }

        // Add upload log row
        function addUploadLog(index, filename, status, code, message) {
            const tbody = document.getElementById('uploadLog');
            const row = document.createElement('tr');
            row.id = 'upload-' + index;
            row.className = 'log-' + status;
            row.innerHTML = `
                <td>${index}</td>
                <td>${filename}</td>
                <td><span class="badge badge-info">Uploading</span></td>
                <td>${code}</td>
                <td>${message}</td>
            `;
            tbody.appendChild(row);
        }

        // Update upload log row
        function updateUploadLog(index, status, code, message) {
            const row = document.getElementById('upload-' + index);
            if (row) {
                row.className = 'log-' + status;
                const badge = status === 'success' ? 'badge-success' : 'badge-danger';
                const statusText = status === 'success' ? 'Uploaded' : 'Failed';
                row.innerHTML = `
                    <td>${index}</td>
                    <td>${row.children[1].textContent}</td>
                    <td><span class="badge ${badge}">${statusText}</span></td>
                    <td>${code}</td>
                    <td>${message}</td>
                `;
            }
        }

        // Update progress bar
        function updateProgress() {
            const percent = Math.round((uploadedCount / totalFiles) * 100);
            const progressBar = document.getElementById('uploadProgress');
            progressBar.style.width = percent + '%';
            progressBar.textContent = percent + '%';
        }

        // Start checking processing status
        function startStatusCheck() {
            if (statusCheckInterval) {
                clearInterval(statusCheckInterval);
            }

            statusCheckInterval = setInterval(async () => {
                try {
                    const response = await fetch('{{ route("sale.product.batch.upload.status") }}?session_id=' + sessionId);
                    const result = await response.json();

                    // Update counts
                    document.getElementById('successCount').textContent = result.success;
                    document.getElementById('failedCount').textContent = result.failed;
                    if(document.getElementById('processingCount')) {
                        document.getElementById('processingCount').textContent = result.processing;
                    }

                    // Update processing log
                    updateProcessingLog(result.logs);

                    // Stop checking only when nothing is left in processing AND we've accounted for all files
                    if (result.processing === 0 && result.total >= totalFiles && uploadedCount >= totalFiles) {
                        clearInterval(statusCheckInterval);
                        const startBtn = document.getElementById('startUpload');
                        startBtn.disabled = false;
                        startBtn.innerHTML = '<i class="mdi mdi-upload"></i> Start Upload';
                    }
                } catch (error) {
                    console.error('Error checking status:', error);
                }
            }, 2000); // Check every 2 seconds
        }

        // Update processing log
        function updateProcessingLog(logs) {
            const tbody = document.getElementById('processingLog');
            
            if (logs.length === 0) {
                return;
            }

            tbody.innerHTML = '';

            logs.forEach((log, index) => {
                let statusClass = '';
                let badgeClass = '';
                let statusText = log.status || 'Unknown';

                switch(log.status) {
                    case 'SUCCESS':
                        statusClass = 'log-success';
                        badgeClass = 'badge-success';
                        break;
                    case 'FAILED':
                    case 'ERROR':
                        statusClass = 'log-failed';
                        badgeClass = 'badge-danger';
                        break;
                    case 'PENDING':
                        statusClass = 'log-uploading';
                        badgeClass = 'badge-warning';
                        statusText = 'Processing...';
                        break;
                    default:
                        badgeClass = 'badge-secondary';
                }

                const row = document.createElement('tr');
                row.className = statusClass;
                row.innerHTML = `
                    <td>${index + 1}</td>
                    <td>${log.file || '-'}</td>
                    <td><span class="badge ${badgeClass}">${statusText}</span></td>
                    <td>${log.code || '-'}</td>
                    <td>${log.reason || (log.text ? log.text.substring(0, 100) + '...' : '-')}</td>
                `;
                tbody.appendChild(row);
            });
        }
    </script>
@endsection
