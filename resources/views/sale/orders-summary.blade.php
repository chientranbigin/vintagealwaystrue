@extends('sale.master')

@section('content')

    <style>
        /* ===== Order Dashboard ===== */
        .order-toolbar .btn {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            line-height: 1;
        }

        .order-stats .card {
            border-radius: 12px;
        }

        @media (min-width: 768px) {
            .col-stats {
                flex: 0 0 70%;
                max-width: 70%;
            }
            .col-toolbar {
                flex: 0 0 30%;
                max-width: 30%;
            }
        }

        .order-stats h1, .order-stats h6, .order-stats h3 {
            color: #fff !important;
        }

        .order-stats h1, .order-stats h3 {
            font-weight: 700;
        }

        @media (max-width: 767px) {
            .order-stats {
                margin-left: -3px;
                margin-right: -3px;
            }
            .order-stats > [class*="col-"] {
                padding-left: 3px;
                padding-right: 3px;
            }
            .order-stats .card-body {
                padding: 0.5rem 0.5rem;
                flex-direction: column;
                align-items: flex-start !important;
            }
            .order-stats h6 {
                font-size: 11px;
                margin-bottom: 2px !important;
                white-space: nowrap;
            }
            .order-stats h3 {
                font-size: 13px;
                margin-bottom: 0 !important;
            }
            .order-stats h1 {
                font-size: 20px;
                align-self: flex-end;
                margin-top: -22px;
            }
        }


        .order-toolbar input,
        .order-toolbar .btn {
            min-height: 44px;
        }

        /* ===== Order Card ===== */
        .order-card {
            border-radius: 12px;
            transition: all 0.2s ease;
            height: 100%;
        }

        .order-card:hover {
            transform: translateY(-3px);
            box-shadow: 0 6px 20px rgba(0, 0, 0, .08);
        }

        .order-header .dropdown-toggle::after {
            display: none;
        }

        .order-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 8px;
        }

        .order-info p {
            font-size: 14px;
            margin-bottom: 4px;
        }

        .order-products {
            display: flex;
            flex-wrap: wrap;
            gap: 8px;
        }

        .order-products img {
            border-radius: 6px;
            border: 1px solid #eee;
        }

        .order-summary {
            font-size: 14px;
            border-top: 1px dashed #ddd;
            padding-top: 8px;
            margin-top: 8px;
        }

        .order-actions {
            display: flex;
            justify-content: space-between;
            gap: 6px;
            margin-top: 10px;
        }

        @media (max-width: 768px) {
            .order-toolbar {
                flex-direction: column;
            }
        }

        .order-toolbar {
            display: flex;
            justify-content: space-between;
            gap: 12px;
            margin-bottom: 16px;
        }

        /* LEFT / RIGHT group */
        .toolbar-left,
        .toolbar-right {
            display: flex;
            gap: 10px;
            flex-wrap: wrap;
        }

        /* Search input width desktop */
        .toolbar-left input {
            min-width: 260px;
        }

        /* Button center fix */
        .order-toolbar .btn {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            min-height: 44px;
        }

        /* ===== Order Actions Layout ===== */
        .order-actions {
            display: flex;
            gap: 6px;
        }

        .order-actions .btn {
            flex: 1;
            display: inline-flex;
            justify-content: center;
            align-items: center;
            min-height: 40px;
        }

        .order-footer {
            margin-top: auto;
        }

        .order-summary {
            font-size: 14px;
            border-top: 1px dashed #ddd;
            padding-top: 8px;
            margin-top: 8px;
        }

        /* ===== Order Products Grid ===== */
        .order-products {
            display: grid;
            grid-template-columns: repeat(3, 1fr); /* 1 hàng 2 hình */
            gap: 8px;
            margin-top: 15px;
        }

        .order-product-item {
            display: block;
        }

        .order-product-item img {
            width: 100%;
            height: auto; /* to hơn */
            object-fit: cover; /* không méo hình */
            border-radius: 8px;
            border: 1px solid #eee;
            transition: all 0.2s ease;
        }

        .order-product-item img:hover {
            transform: scale(1.03);
        }

        /* ===== Responsive ===== */
        @media (max-width: 768px) {
            /*.order-actions {*/
            /*    flex-direction: column;*/
            /*}*/

            .order-actions .btn {
                width: 100%;
            }
        }

        /* ===== Responsive ===== */
        @media (max-width: 768px) {
            .order-toolbar {
                flex-direction: column;
            }

            .toolbar-left,
            .toolbar-right {
                width: 100%;
            }

            /* Input full width */
            .toolbar-left input {
                width: 100%;
            }

            /* 1 row = 2 buttons */
            .toolbar-left .btn {
                width: calc(50% - 5px);
            }

            .toolbar-right,
            .toolbar-right .dropdown,
            .toolbar-right .btn {
                width: 100%;
            }
        }

        /* ===== Order Actions ===== */
        .order-actions {
            display: flex;
            gap: 6px;
            flex-wrap: wrap;

        }

        .order-actions .btn {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            white-space: nowrap;
            flex: 0 0 calc(50% - 3px); /* 2 nút / hàng */

        }

        /*!* ===== Create Order highlight ===== *!*/
        /*.btn-create-order {*/
        /*    background: linear-gradient(135deg, #17a2b8, #0dcaf0);*/
        /*    border: none;*/
        /*}*/

        /*.btn-create-order:hover {*/
        /*    background: linear-gradient(135deg, #138496, #0bbcd6);*/
        /*}*/

        /* ===== Mobile: 2 buttons per row ===== */
        @media (max-width: 768px) {
            .order-actions {
                flex-wrap: wrap;
            }

            .order-actions .btn {
                flex: 0 0 calc(50% - 3px); /* 2 nút / hàng */
            }
        }

    </style>

    <div class="main-content">
        <div class="page-content">
            <div class="container-fluid pr-0 pl-0">
                <div class="row" style="    flex-direction: row-reverse;">
                    <div class="col-stats col-sm-12">
                        {{-- ===== STATISTICS ===== --}}

                        <div class="row order-stats mb-4">

                            <div class="col-4">
                                <div class="card bg-primary text-white mb-0" style="color: white">
                                    <div class="card-body d-flex justify-content-between align-items-center">
                                        <div>
                                            <h6 class="mb-1">Today</h6>
                                            <h3 class="mb-0">{{ number_format($totalSoldTodayPrice) }}đ</h3>
                                        </div>
                                        <h1 class="mb-0">{{ $totalSoldToday }}</h1>
                                    </div>
                                </div>
                            </div>

                            <div class="col-4">
                                <div class="card bg-info text-white mb-0" style="color: white">
                                    <div class="card-body d-flex justify-content-between align-items-center">
                                        <div>
                                            <h6 class="mb-1">Yesterday</h6>
                                            <h3 class="mb-0">{{ number_format($totalSoldYesterdayPrice) }}đ</h3>
                                        </div>
                                        <h1 class="mb-0">{{ $totalSoldYesterday }}</h1>
                                    </div>
                                </div>
                            </div>

                            <div class="col-4">
                                <div class="card bg-success text-white  mb-0">
                                    <div class="card-body d-flex justify-content-between align-items-center">
                                        <div>
                                            <h6 class="mb-1">Current Month</h6>
                                            <h3 class="mb-0">{{ number_format($currentTotalFinalPrice) }}đ</h3>
                                        </div>
                                        <h1 class="mb-0">{{ $currentTotalOrder }}</h1>
                                    </div>
                                </div>
                            </div>
                        </div>

                    </div>
                    <div class="col-toolbar col-sm-12">
                        {{-- ===== SEARCH & ACTION ===== --}}
                        <form method="POST" action="{{ route('sale.searchOrder') }}">
                            @csrf
                            <input type="hidden" name="urlSearch"
                                   value="{{ \Illuminate\Support\Facades\Route::currentRouteName() }}">
                            <input class="form-control"
                                   name="search"
                                   value="{{ request()->search }}"
                                   placeholder="Phone, Name, Product"/>
                            <div class="order-toolbar mt-2">


                                <!-- LEFT -->
                                <div class="toolbar-left">

                                    <button class="btn btn-secondary">Search</button>

                                    <a href="{{ route(\Illuminate\Support\Facades\Route::currentRouteName(), ['status' => 'NEW']) }}"
                                       class="btn btn-light">Clear</a>
                                </div>

                                <!-- RIGHT -->
                                <div class="toolbar-right">
                                    <div class="dropdown">
                                        <button class="btn btn-warning text-white dropdown-toggle" type="button"
                                                id="createOrderDropdown" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                                            Create
                                        </button>
                                        <div class="dropdown-menu dropdown-menu-right" aria-labelledby="createOrderDropdown">
                                            <a class="dropdown-item" href="{{ route('sale.order.create') }}">Manual</a>
                                            <a class="dropdown-item" href="{{ route('sale.order.smart-create.get') }}">Smart (AI)</a>
                                        </div>
                                    </div>
                                </div>
                            </div>

                        </form>
                    </div>
                </div>



                <div class="mb-3">
                    <b>Total {{ $orders->count() }} entries</b>
                </div>

                {{-- ===== ORDER LIST ===== --}}
                <div class="row">
                    @foreach($orders as $order)
                        <div class="col-xl-2 col-lg-3 col-md-4 col-sm-6 mb-4">
                            <div class="card order-card h-100">
                                <div class="card-body d-flex flex-column">
                                    <div class="order-content">
                                        {{-- Header --}}
                                            <div class="order-header">
                                                <div class="d-flex align-items-center">
                                                    @if($order->status == 'NEW')
                                                        <span class="badge badge-light mr-2">{{ $order->status }}</span>
                                                    @elseif($order->status == 'ON HOLD')
                                                        <span class="badge badge-warning mr-2">{{ $order->status }}</span>
                                                    @elseif($order->status == 'READY FOR SHIPPING' || $order->status == 'A SHIP NOW')
                                                        <span class="badge badge-info mr-2">{{ $order->status }}</span>
                                                    @else
                                                        <span class="badge badge-success mr-2">{{ $order->status }}</span>
                                                    @endif
                                                </div>

                                                <div class="dropdown">
                                                    <a class="text-muted dropdown-toggle font-size-18" role="button" data-toggle="dropdown" aria-haspopup="true">
                                                        <i class="mdi mdi-dots-vertical"></i>
                                                    </a>
                                                    <div class="dropdown-menu dropdown-menu-right">
                                                        <a class="dropdown-item" href="{{ route('sale.order.detail', $order->id) }}">
                                                            <i class="mdi mdi-eye-outline mr-2"></i> Detail
                                                        </a>
                                                        <a class="dropdown-item btn-capture" href="javascript:void(0);" 
                                                           data-order="{{ json_encode([
                                                               'id' => $order->id,
                                                               'name' => $order->name,
                                                               'phone' => $order->phone,
                                                               'address' => $order->address,
                                                               'note' => $order->note,
                                                               'total' => $order->total_amount,
                                                               'deposit' => $order->deposit_amount,
                                                               'final' => $order->is_paid ? 0 : ($order->total_amount - $order->deposit_amount),
                                                               'is_paid' => $order->is_paid,
                                                               'products' => $order->products->map(function($p) {
                                                                   return ['name' => $p->name, 'path' => asset($p->path_thumb)];
                                                               })
                                                           ]) }}">
                                                            <i class="mdi mdi-camera-outline mr-2"></i> Capture Order
                                                        </a>
                                                        <a class="dropdown-item" href="{{ route('sale.order.create', ['consumer_id' => $order->consumer_id]) }}">
                                                            <i class="mdi mdi-plus-circle-outline mr-2"></i> Create Similar Order
                                                        </a>
                                                        @if($order->status == 'NEW' || $order->status == 'ON HOLD')
                                                            <a class="dropdown-item" href="{{ route('sale.order.ship', $order->id) }}">
                                                                <i class="mdi mdi-truck mr-2"></i> Ready to Ship
                                                            </a>
                                                        @endif
                                                        @if($order->status == 'NEW')
                                                            <a class="dropdown-item" href="{{ route('sale.order.hold', $order->id) }}">
                                                                <i class="mdi mdi-pause-circle-outline mr-2"></i> Mark On Hold
                                                            </a>
                                                        @endif
                                                        @if($order->status != 'COMPLETED')
                                                            <a class="dropdown-item text-success" href="{{ route('sale.order.complete', $order->id) }}"
                                                               onclick="return confirm('Hoàn tất đơn hàng này?')">
                                                                <i class="mdi mdi-check-circle-outline mr-2"></i> Complete
                                                            </a>
                                                        @endif
                                                        <div class="dropdown-divider"></div>
                                                        <a class="dropdown-item text-danger" href="{{ route('sale.order.detail.delete', $order->id) }}"
                                                           onclick="return confirm('Xóa đơn hàng này?')">
                                                            <i class="mdi mdi-trash-can-outline mr-2"></i> Delete
                                                        </a>
                                                    </div>
                                                </div>
                                            </div>
                                            <div class="mb-2">
                                                <small class="text-muted">
                                                    {{ \Carbon\Carbon::parse($order->created_at)->diffForHumans() }}
                                                </small>
                                            </div>

                                        {{-- Info --}}
                                       
                                            <div class="order-info">
                                                <p><b>Họ tên:</b> {{ $order->name }}</p>
                                                <p><b>SDT:</b> {{ $order->phone }}</p>
                                                <p><b>Địa chỉ:</b> {{ $order->address }}</p>

                                                @if($order->note)
                                                    <p class="text-success">Note: {{ $order->note }}</p>
                                                @endif

                                                @if($order->desc)
                                                    <p class="text-danger">
                                                        Note sửa đồ:<br>{!! nl2br($order->desc) !!}
                                                    </p>
                                                @endif
                                            </div>
                     
                                        {{-- Products --}}
                                        <div class="order-products">
                                            @foreach($order->products as $product)
                                                <a href="{{ route('sale.product.detail', $product->id) }}"
                                                   class="order-product-item">
                                                    {{$product->name}}
                                                    <img src="{{ asset($product->path_thumb) }}">
                                                </a>
                                            @endforeach
                                        </div>


                                    </div>
                                    <div class="order-footer mt-auto">

                                        {{-- Summary --}}
                                        <div class="order-summary">
                                            @if(intval($order->total_amount))
                                                <div class="d-flex justify-content-between">
                                                    <span>Tổng đơn</span>
                                                    <span>{{ number_format($order->total_amount) }}đ</span>
                                                </div>
                                            @endif

                                            @if(intval($order->freeship_amount))
                                                <div class="d-flex justify-content-between">
                                                    <span>Giảm giá</span>
                                                    <span>{{ number_format($order->freeship_amount) }}đ</span>
                                                </div>
                                            @endif

                                            @if($order->is_paid || intval($order->deposit_amount))
                                                <div class="d-flex justify-content-between">
                                                    <span>Đã thanh toán</span>
                                                    <span>
                        {{ number_format($order->is_paid ? $order->total_amount : $order->deposit_amount) }}đ
                    </span>
                                                </div>
                                            @endif

                                            <div class="d-flex justify-content-between">
                                                <span>Đơn Viettel</span>
                                                <span>
                    {{ number_format($order->is_paid ? 0 : $order->total_amount - $order->deposit_amount) }}đ
                    <br>
                    <small>{{ $order->is_paid ? 'Freeship' : '+ Phí ship' }}</small>
                </span>
                                            </div>
                                        </div>

                                        <!-- Actions moved to header dropdown -->


                                    </div>


                                </div>
                            </div>
                        </div>
                    @endforeach
                </div>


        {{-- Hidden Capture Template --}}
        <div id="order-receipt-template">
            <div class="c-header">
                <h1>VINTAGE</h1>
                <p>Always True</p>
            </div>

            <div class="c-customer-card">
                <div class="c-customer-header">
                    <div class="c-avatar font-italic">P</div>
                    <div class="c-cust-info">
                        <span>KHÁCH HÀNG</span>
                        <h2 class="c-name">---</h2>
                    </div>
                </div>
                <div class="c-detail-item">
                    <i class="mdi mdi-phone"></i>
                    <span class="c-phone">---</span>
                </div>
                <div class="c-detail-item">
                    <i class="mdi mdi-map-marker"></i>
                    <span class="c-address">---</span>
                </div>
            </div>

            <div class="c-note-box">
                <label>GHI CHÚ:</label>
                <div class="c-note-text">---</div>
            </div>

            <div class="c-section-title">ĐƠN HÀNG CỦA BẠN</div>
            
            <div class="c-products">
                <!-- Dynamic items -->
            </div>

            <div class="c-amount-box">
                <div class="c-amount-row">
                    <div class="c-amount-label">Tổng đơn</div>
                    <div class="c-amount-value c-total-amount">0đ</div>
                </div>
                <div class="c-amount-row final">
                    <div class="c-amount-label">Đơn Viettel</div>
                    <div class="c-amount-value c-final-amount">0đ</div>
                </div>
                <div class="c-ship-hint">+ Phí ship</div>
            </div>

            <div class="c-footer-info">
                <div class="c-info-col">
                    <h4>🚚 GIAO HÀNG</h4>
                    <ul>
                        <li>• Viettel Post: HCM 1-2 ngày, tỉnh 3-5 ngày.</li>
                        <li>• Phí ship đồng giá: 20k toàn quốc.</li>
                        <li>• Hàng được giặt sấy, là ủi và xếp gọn tỉ mỉ.</li>
                    </ul>
                </div>
                <div class="c-info-col">
                    <h4>🔄 ĐỔI TRẢ</h4>
                    <ul>
                        <li>• Đổi trả nếu hàng không đúng mô tả, sai số đo >2cm.</li>
                        <li>• Không nhận trả hàng với lý do cá nhân (không thích, không hợp).</li>
                    </ul>
                </div>
            </div>

            <div class="c-thank-you">
                <h3>Thank you for your order!</h3>
                <p>Vintage Always True • Vintage Menswear</p>
            </div>
        </div>

        {{-- Modal Preview Capture --}}
        <div class="modal fade" id="captureModal" tabindex="-1" role="dialog" aria-hidden="true">
            <div class="modal-dialog modal-dialog-centered modal-lg" role="document">
                <div class="modal-content border-0 shadow-lg" style="border-radius: 20px; overflow: hidden;">
                    <div class="modal-header bg-white border-0 py-3 px-4">
                        <h5 class="modal-title font-weight-bold" style="color: #333; font-family: 'Outfit', sans-serif;">Đã tạo ảnh đơn hàng</h5>
                        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                            <span aria-hidden="true" style="font-size: 28px;">&times;</span>
                        </button>
                    </div>
                    <div class="modal-body text-center p-4" style="background: #f8f9fa;">
                        <p class="text-muted mb-3 small" style="font-family: 'Outfit', sans-serif;">Bạn có thể nhấn giữ vào ảnh để Lưu hoặc Chia sẻ trực tiếp</p>
                        <div class="shadow-sm rounded-lg overflow-hidden border mb-1" style="background: white;">
                            <img id="previewCapturedImage" src="" class="img-fluid" alt="Order Capture">
                        </div>
                    </div>
                    <div class="modal-footer border-0 justify-content-center pb-4" style="background: #f8f9fa;">
                        <button type="button" class="btn btn-light px-4 py-2 border" data-dismiss="modal" style="border-radius: 10px; font-weight: 600;">Đóng</button>
                        <button type="button" id="btnModalCopy" class="btn btn-primary px-4 py-2" style="border-radius: 10px; font-weight: 600; background: #0d6efd;">Copy ảnh</button>
                        <button type="button" id="btnModalDownload" class="btn btn-success px-4 py-2" style="border-radius: 10px; font-weight: 600;">Tải về</button>
                    </div>
                </div>
            </div>
        </div>
    @endsection

@section('script')
    <script src="https://cdnjs.cloudflare.com/ajax/libs/html2canvas/1.4.1/html2canvas.min.js"></script>
    <script>
        $(document).ready(function() {
            $(document).on('click', '.btn-capture', function() {
                const $btn = $(this);
                const orderData = $btn.data('order');
                const $receipt = $('#order-receipt-template');
                const originalHtml = $btn.html();
                
                // Feedback
                $btn.html('<i class="mdi mdi-loading mdi-spin"></i> Processing...').prop('disabled', true);

                // Populate template
                $receipt.find('.c-name').text(orderData.name);
                $receipt.find('.c-phone').text(orderData.phone);
                $receipt.find('.c-address').text(orderData.address);
                $receipt.find('.c-note-text').text(orderData.note || 'Không có');
                $receipt.find('.c-note-box').toggle(!!orderData.note);

                const $prodContainer = $receipt.find('.c-products').empty();
                orderData.products.forEach(p => {
                    $prodContainer.append(`
                        <div class="c-prod-item">
                            <img src="${p.path}" crossorigin="anonymous">
                            <div class="c-prod-code">${p.name}</div>
                        </div>
                    `);
                });

                const formatMoney = (val) => new Intl.NumberFormat('vi-VN').format(val) + 'đ';
                $receipt.find('.c-total-amount').text(formatMoney(orderData.total));
                $receipt.find('.c-final-amount').text(formatMoney(orderData.final));
                
                $receipt.css({
                    'display': 'block', 'position': 'fixed', 'top': '0', 'left': '-9999px',
                    'opacity': '1', 'pointer-events': 'none', 'z-index': '9999', 'background': 'white',
                    'width': '450px' // Fix width to ensure correct rendering
                });

                let currentCanvas = null;
                let currentBlob = null;
                let currentFileName = "";

                html2canvas($receipt[0], {
                    useCORS: true, scale: 2, backgroundColor: '#ffffff', logging: false
                }).then(async canvas => {
                    currentCanvas = canvas;
                    const jpegBlob = await new Promise(resolve => canvas.toBlob(resolve, 'image/jpeg', 0.8));
                    currentBlob = jpegBlob;
                    
                    $btn.html(originalHtml).prop('disabled', false);
                    $receipt.hide();

                    if (!jpegBlob) {
                        alert('Lỗi tạo ảnh.');
                        return;
                    }

                    currentFileName = `Order_${orderData.name}_${orderData.id}.jpg`;
                    const dataUrl = canvas.toDataURL('image/jpeg', 0.8);
                    
                    // Show Modal instead of direct share
                    $('#previewCapturedImage').attr('src', dataUrl);
                    $('#captureModal').modal('show');
                    
                }).catch(err => {
                    alert('Lỗi: ' + err.message);
                    $btn.html(originalHtml).prop('disabled', false);
                    $receipt.hide();
                });

                // Modal Buttons logic
                $('#btnModalCopy').off('click').on('click', async function() {
                    if (typeof ClipboardItem !== 'undefined') {
                        try {
                            const pngBlob = await new Promise(resolve => currentCanvas.toBlob(resolve, 'image/png'));
                            const data = [new ClipboardItem({ 'image/png': pngBlob })];
                            await navigator.clipboard.write(data);
                            alert('Đã copy vào Clipboard!');
                        } catch (err) {
                            alert('Không thể copy, vui lòng nhấn giữ ảnh để lưu.');
                        }
                    } else {
                        alert('Trình duyệt không hỗ trợ copy ảnh.');
                    }
                });

                $('#btnModalDownload').off('click').on('click', function() {
                    const link = document.createElement('a');
                    link.download = currentFileName;
                    link.href = currentCanvas.toDataURL('image/jpeg', 0.8);
                    link.click();
                });
            });
        });
    </script>
@endsection
@section('style')
    <style>
        #order-receipt-template {
            width: 450px;
            padding: 30px;
            font-family: 'Inter', sans-serif;
            color: #333;
            display: none;
        }
        .c-header { text-align: center; margin-bottom: 30px; }
        .c-header h1 { font-size: 32px; font-weight: 800; letter-spacing: 5px; margin: 0; color: #111; }
        .c-header p { font-size: 14px; font-weight: 600; letter-spacing: 3px; color: #0066cc; margin-top: 5px; text-transform: uppercase; }
        
        .c-customer-card {
            background: #f8f9fa;
            border-radius: 20px;
            padding: 25px;
            margin-bottom: 20px;
        }
        .c-customer-header { display: flex; align-items: center; margin-bottom: 15px; }
        .c-avatar { background: #007bff; color: white; width: 45px; height: 45px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-weight: bold; margin-right: 15px; font-size: 18px; }
        .c-cust-info h2 { font-size: 20px; font-weight: 700; margin: 0; }
        .c-cust-info span { font-size: 12px; color: #888; text-transform: uppercase; }
        
        .c-detail-item { display: flex; align-items: flex-start; margin-bottom: 10px; font-size: 14px; color: #555; }
        .c-detail-item i { width: 25px; color: #777; font-size: 16px; margin-top: 2px; }

        .c-note-box {
            background: #fff9e6;
            border-radius: 15px;
            padding: 15px;
            margin-bottom: 25px;
        }
        .c-note-box label { font-weight: 700; font-size: 12px; color: #b08d00; display: block; margin-bottom: 5px; }
        .c-note-text { font-size: 14px; font-weight: 500; }

        .c-section-title { font-size: 12px; font-weight: 700; color: #bbb; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 15px; border-bottom: 1px solid #eee; padding-bottom: 10px; }

        .c-products { display: flex; flex-wrap: wrap; gap: 15px; margin-bottom: 25px; }
        .c-prod-item { width: calc(50% - 8px); text-align: center; }
        .c-prod-item img { 
            max-width: 100%; 
            max-height: 220px; 
            width: auto;
            height: auto;
            object-fit: contain; 
            background: #f8f9fa;
            border-radius: 12px; 
            margin-bottom: 8px;
            display: block;
            margin-left: auto;
            margin-right: auto;
        }
        .c-prod-code { font-weight: 700; font-size: 14px; }

        .c-amount-box {
            background: #fdfdfd;
            border: 1px solid #f0f0f0;
            border-radius: 20px;
            padding: 25px;
            margin-bottom: 30px;
        }
        .c-amount-row { display: flex; justify-content: space-between; align-items: center; margin-bottom: 15px; }
        .c-amount-label { font-size: 16px; font-weight: 500; color: #777; }
        .c-amount-value { font-size: 18px; font-weight: 700; }
        .c-amount-row.final { margin-top: 20px; padding-top: 20px; border-top: 1px solid #eee; }
        .c-amount-row.final .c-amount-label { color: #111; font-weight: 700; font-size: 18px; }
        .c-amount-row.final .c-amount-value { color: #0066cc; font-size: 26px; }
        .c-ship-hint { font-size: 12px; color: #aaa; text-align: right; margin-top: -10px; }

        .c-footer-info { display: flex; gap: 20px; margin-bottom: 40px; }
        .c-info-col { flex: 1; }
        .c-info-col h4 { font-size: 12px; font-weight: 800; text-transform: uppercase; margin-bottom: 10px; display: flex; align-items: center; gap: 5px; }
        .c-info-col ul { padding: 0; margin: 0; list-style: none; }
        .c-info-col li { font-size: 11px; color: #777; line-height: 1.5; margin-bottom: 5px; }

        .c-thank-you { text-align: center; border-top: 1px dashed #ddd; padding-top: 25px; }
        .c-thank-you h3 { font-size: 16px; font-weight: 700; margin-bottom: 5px; }
        .c-thank-you p { font-size: 11px; color: #aaa; }

        /* Modal Preview Styles */
        #previewCapturedImage {
            max-height: 65vh;
            width: auto;
            margin: 0 auto;
            display: block;
            border-radius: 8px;
        }
        #captureModal .modal-body {
            max-height: 75vh;
            overflow-y: auto;
        }
        @media (max-width: 768px) {
            #captureModal .modal-dialog {
                margin: 0.5rem;
            }
            #previewCapturedImage {
                max-height: 60vh;
            }
        }
    </style>
@endsection
