@extends('sale.master')
@section('content')
    @php
        $total = $orders->total();
           $currentPage = $orders->currentPage();
           $perPage = $orders->perPage();

           $from = ($currentPage - 1) * $perPage + 1;
           $to = min($currentPage * $perPage, $total);
    @endphp
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">


                <div class="card">
                    <div class="card-body">
                        <div class="row">
                            <div class="card col bg-primary border-primary mx-2">
                                <div class="card-body">
                                    <div class=" d-flex  align-items-center" style="justify-content: space-between">
                                        <div><h5 class="card-title mb-0 text-white">New Orders</h5></div>
                                        <div><h2 class=" mb-0 text-white"> {{number_format($totalOrderNewPrice)}}đ</h2>
                                        </div>
                                        <div><h2 class=" mb-0 text-white">{{$totalOrderNew}}</h2></div>
                                    </div>
                                </div>
                            </div>
                            <div class="card col bg-success border-success mx-2">
                                <div class="card-body">
                                    <div class=" d-flex  align-items-center" style="justify-content: space-between">
                                        <div><h5 class="card-title mb-0 text-white">Current Month Orders</h5></div>
                                        <div><h2 class=" mb-0 text-white"> {{number_format($currentTotalFinalPrice)}}
                                                đ</h2></div>
                                        <div><h2 class=" mb-0 text-white">{{$currentTotalOrder}}</h2></div>
                                    </div>
                                </div>
                            </div>


                        </div>
                        <form method="POST" action="{{ route('sale.searchOrder') }}"
                              enctype="multipart/form-data">
                            @csrf
                            <input type="hidden" class="form-control"
                                   name="urlSearch" value="{{\Illuminate\Support\Facades\Route::currentRouteName()}}">
                            <div class="row mb-3">


                                <div class="col-sm-12 col-md-6">
                                    <input id="input-search" class="form-control mt-3" name="search"
                                           value="{{request()->search}}" placeholder="Phone, Name, Product">

                                </div>


                                <div class="col-6">

                                    <button type="submit"
                                            class="btn btn-secondary waves-effect waves-light btn-submit mt-3 mr-3">
                                        Search
                                    </button>
                                    <a href="{{route('sale.order')}}"
                                       class="btn btn-light waves-effect waves-light btn-submit mt-3">
                                        Clear
                                    </a>

                                </div>
                            </div>

                        </form>
                        <a href="{{route('sale.order.create')}}"
                           class="btn btn-secondary waves-effect waves-light btn-submit  mb-3"
                           style="    color: #fff !important;">Create </a>
                        <a href="{{route('sale.order.smart-create.get')}}"
                           class="btn btn-info waves-effect waves-light btn-submit  mb-3"
                           style="    color: #fff !important;">Create Smart</a>
                        <a href="{{ route('sale.order') }}"
                           class="btn btn-secondary waves-effect waves-light btn-submit  mb-3"
                           style="    color: #fff !important;">All Orders </a>

                        <div class="table-responsive">
                            <div class="py-3">
                                <b>Showing {{$from}} to {{$to}} of {{$total}} entries</b>
                            </div>
                            {!! $orders->appends(request()->query())->links() !!}

                            <table class="table mb-0">
                                <thead>
                                <tr>

                                    <th>Consumer</th>
                                    <th>Summary</th>

                                    <th>Products</th>
{{--                                    <th>Description</th>--}}
                                    <th>Additional Amount</th>
                                    {{--                                    <th>Total Amount</th>--}}
                                    {{--                                    <th>Is Paid Full</th>--}}
                                    {{--                                    <th>Is Freeship</th>--}}
                                    <th>Deposit Amount</th>
                                                                        <th>Viettel Order Amount</th>
                                    <th>Final Amount</th>
                                    <th>Status</th>
                                    <th>Action</th>
                                </tr>
                                </thead>
                                <tbody>
                                @foreach($orders as $order)
                                    <tr>
                                        <td><a href="{{route('sale.order.detail', $order->id)}}">{!! implode('<br>', [
                                            $order->name,
                                            $order->phone,
                                            $order->address,

                                            ]) !!}</a>
                                            @if($order->note)
                                                <br><br>
                                                <span class="right badge badge-danger">{{$order->note}}</span>
                                            @endif
                                        </td>
                                        <td>
                                            <p>
                                                Họ Tên: {{$order->name}}<br>
                                                SDT: {{$order->phone}}<br>
                                                Địa chỉ: {{$order->address}}
                                                @if($order->note)
                                                    <br>
                                                    <span class="right badge badge-danger">Note: {{$order->note}}</span>
                                                @endif
                                                @if($order->desc)
                                                    <br>
                                                    <span class="right badge badge-danger">Note sửa đồ: {!! nl2br($order->desc) !!}</span>
                                                @endif

                                            </p>
                                            @foreach($order->products as $product)
                                                <p>{{$product->type}} {{$product->name}}
                                                    : {{number_format($product->price)}}đ</p>
                                            @endforeach
                                            @if(intval($order->additional_amount))
                                                <p>Phí sửa đồ: {{number_format($order->additional_amount)}}đ</p>
                                            @endif
                                            @if(intval($order->freeship_amount))
                                                <p>Giảm giá: {{number_format($order->freeship_amount)}}đ </p>
                                            @endif
                                            @if($order->is_paid || intval($order->deposit_amount))
                                                <p>Đã thanh toán
                                                    trước: {{number_format($order->is_paid ? $order->total_amount : $order->deposit_amount)}}
                                                    đ </p>
                                            @endif
                                            <div class="d-flex">
                                                <p>Tổng
                                                    đơn: {{number_format($order->is_paid ? 0 : $order->total_amount - $order->deposit_amount)}}
                                                    đ {{$order->is_paid ? 'Freeship' : '+ Phí ship'}}</p>
                                            </div>


                                        </td>
                                        <td>
                                            <div class="d-flex" style="max-width: 500px; flex-wrap: wrap">
                                                @foreach($order->products as $product)
                                                    <div class="mt-1">
                                                        <a href="{{route('sale.product.detail', $product->id)}}">
                                                            {{$product->name}} <br/>
                                                            <img src="{{asset($product->path_thumb)}}" alt=""
                                                                 height="100" class="ml-1">
                                                        </a>

                                                    </div>
                                                    <br>
                                                @endforeach
                                            </div>

                                        </td>
{{--                                        <td>{{$order->desc}}</td>--}}
                                        <td>{{number_format($order->additional_amount)}}đ</td>
                                        {{--                                        <td><span >{{number_format($order->total_amount)}}đ</span></td>--}}
                                        {{--                                        <td>{{$order->is_paid ? 'Yes' : 'No'}}</td>--}}
                                        {{--                                        <td>{{$order->is_freeship ? 'Yes' : 'No'}}</td>--}}
                                        <td><span>{{number_format($order->is_paid ? $order->total_amount - $order->deposit_amount : $order->deposit_amount)}}đ</span>
                                        </td>
                                        <td><span style="background-color:grey; color:white; padding: 5px">{{number_format($order->is_paid ? 0 : $order->total_amount - $order->deposit_amount)}}đ</span>
                                        </td>
                                        <td>{{number_format($order->final_amount)}}đ</td>

                                        <td>
                                            @if($order->status == 'NEW')
                                                <span class="right badge badge-light">{{$order->status}}</span>
                                            @elseif($order->status == 'ON HOLD')
                                                <span class="right badge badge-warning">{{$order->status}}</span>
                                            @elseif($order->status == 'READY FOR SHIPPING')
                                                <span class="right badge badge-info">{{$order->status}}</span>
                                            @else
                                                <span class="right badge badge-success">{{$order->status}}</span>
                                            @endif
                                        </td>
                                        <td>
                                            <div class="dropdown">
                                                <a class="text-muted dropdown-toggle font-size-20" role="button" data-toggle="dropdown" aria-haspopup="true">
                                                    <i class="mdi mdi-dots-vertical"></i>
                                                </a>
                                                <div class="dropdown-menu dropdown-menu-right">
                                                    <a class="dropdown-item" href="{{route('sale.order.detail', $order->id)}}">
                                                        <i class="mdi mdi-eye-outline mr-2"></i> View Detail
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
                                                    @if($order->status != 'COMPLETED')
                                                        <a class="dropdown-item" href="{{route('sale.order.complete', $order->id)}}">
                                                            <i class="mdi mdi-check-circle-outline mr-2 text-success"></i> Complete
                                                        </a>
                                                    @endif
                                                    <div class="dropdown-divider"></div>
                                                    <a class="dropdown-item text-danger" href="{{route('sale.order.detail.delete', $order->id)}}" 
                                                       onclick="return confirm('Are you sure you want to delete this order?')">
                                                        <i class="mdi mdi-trash-can-outline mr-2"></i> Delete
                                                    </a>
                                                </div>
                                            </div>
                                        </td>
                                    </tr>
                                @endforeach


                                </tbody>
                            </table>

                            {!! $orders->appends(request()->query())->links() !!}

                        </div>

                    </div>
                    <!-- end card-body-->
                </div>

            </div> <!-- container-fluid -->
        </div>
        <!-- End Page-content -->



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
                
                $btn.html('<i class="mdi mdi-loading mdi-spin"></i> Processing...').prop('disabled', true);

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
                    'width': '450px' 
                });

                let currentCanvas = null;
                let currentFileName = "";

                html2canvas($receipt[0], {
                    useCORS: true, scale: 2, backgroundColor: '#ffffff', logging: false
                }).then(async canvas => {
                    currentCanvas = canvas;
                    const jpegBlob = await new Promise(resolve => canvas.toBlob(resolve, 'image/jpeg', 0.8));
                    
                    $btn.html(originalHtml).prop('disabled', false);
                    $receipt.hide();

                    if (!jpegBlob) {
                        alert('Lỗi tạo ảnh.');
                        return;
                    }

                    currentFileName = `Order_${orderData.name}_${orderData.id}.jpg`;
                    const dataUrl = canvas.toDataURL('image/jpeg', 0.8);
                    
                    $('#previewCapturedImage').attr('src', dataUrl);
                    $('#captureModal').modal('show');

                }).catch(err => {
                    alert('Lỗi: ' + err.message);
                    $btn.html(originalHtml).prop('disabled', false);
                    $receipt.hide();
                });

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
