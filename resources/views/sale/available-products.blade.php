@extends('sale.master')
@section('content')
    @php
        $total = $products->total();
           $currentPage = $products->currentPage();
           $perPage = $products->perPage();

           $from = ($currentPage - 1) * $perPage + 1;
           $to = min($currentPage * $perPage, $total);
           @endphp
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">
                <form method="POST" action="{{ route('sale.searchAvailableProducts') }}" enctype="multipart/form-data">
                    @csrf

                    <div class="row mb-3">


                        <div class="col-sm-12 col-md-6">
                            <input id="input-search" class="form-control" name="search" value="{{request()->search}}" placeholder="Name Products">

                        </div>


                        <div class=" col-sm-12 col-md-6 d-flex ">
                            <select class="custom-select rounded-0" id="exampleSelectRounded0" name="type" style="width: 300px">
                                <option {{!request()->type ? 'selected' : ''}}></option>

                                @foreach(\App\Contracts\SaleConfig::PRODUCT_TYPES as $type)
                                    <option {{request()->type == $type ? 'selected' : ''}}>{{$type}}</option>
                                @endforeach
                            </select>

                        </div>
                        <div class="col-sm-12 col-md-6  ">
                            <div class="row mt-3">
                                <div class="col d-flex">

                                    <select class="custom-select rounded-0" id="exampleSelectRounded0" name="from_search_key" placeholder="From">
                                        <option {{!request()->from_search_key ? 'selected' : ''}}></option>

                                        @foreach(\App\Contracts\SaleConfig::SIZES as $size)
                                            <option {{request()->from_search_key == $size ? 'selected' : ''}}>{{$size}}</option>
                                        @endforeach
                                    </select>
                                    <input id="from_search_value" class="form-control" name="from_search_value" placeholder="From Value"
                                           value="{{request()->from_search_value}}">
                                    <input id="to_search_value" class="form-control" name="to_search_value" placeholder="To Value"
                                           value="{{request()->to_search_value}}">
                                </div>

                            </div>

                        </div>
                        <div class="col-sm-12 col-md-6  ">
                            <div class="row mt-3">
                                <div class="col d-flex">

                                    <select class="custom-select rounded-0" id="exampleSelectRounded0" name="from_search_key_2" placeholder="From">
                                        <option {{!request()->from_search_key_2 ? 'selected' : ''}}></option>

                                        @foreach(\App\Contracts\SaleConfig::SIZES as $size)
                                            <option {{request()->from_search_key_2 == $size ? 'selected' : ''}}>{{$size}}</option>
                                        @endforeach
                                    </select>
                                    <input id="from_search_value" class="form-control" name="from_search_value_2" placeholder="From Value"
                                           value="{{request()->from_search_value_2}}">
                                    <input id="to_search_value" class="form-control" name="to_search_value_2" placeholder="To Value"
                                           value="{{request()->to_search_value_2}}">
                                </div>

                            </div>

                        </div>
                        <div class="col-12">

                            <button type="submit" class="btn btn-secondary waves-effect waves-light btn-submit mt-3 mr-3">
                                Search
                            </button>
                            <a href="{{route('sale.availableProducts')}}"  class="btn btn-light waves-effect waves-light btn-submit mt-3">
                                Clear
                            </a>

                            <span id="download-all-images" class="btn btn-light waves-effect waves-light btn-submit mt-3 float-right d-none d-md-inline-block">
                                Download
                            </span>


                            <a href="{{route('sale.availableProducts', ['showAll' => true])}}"  class="btn btn-light waves-effect waves-light btn-submit mt-3 float-right mr-3 d-none d-md-inline-block">
                                Show All
                            </a>

                        </div>
                    </div>

                </form>
                <div class="py-3">
                    <div class="d-flex mb-3 align-items-center flex-wrap" style="gap: 10px">
                        <span id="select-all-page" class="btn btn-outline-primary waves-effect waves-light">
                            <i class="mdi mdi-check-all mr-1"></i> Select All Page
                        </span>

                        <span id="native-share-btn" class="btn btn-success waves-effect waves-light">
                            <i class="mdi mdi-share-variant mr-1"></i> Share (<span id="selected-count">0</span>)
                        </span>

                        <button onclick="localStorage.removeItem('selectedProductImages'); updateSelectedCount(); alert('Đã xoá danh sách đã chọn!'); location.reload();" class="btn btn-outline-danger waves-effect waves-light">
                            <i class="mdi mdi-trash-can-outline mr-1"></i> Clear
                        </button>
                    </div>
                    <div>
                        <b>Showing {{$from}} to {{$to}} of {{$total}} entries</b>
                    </div>

                    @if(session('failedFiles'))
                        <div class="alert alert-warning">
                            <strong>Các file không xử lý được:</strong>
                            <ul>
                                @foreach(session('failedFiles') as $failed)
                                    <li>
                                        {{ $failed['file'] }} - {{ $failed['reason'] }}
                                    </li>
                                @endforeach
                            </ul>
                        </div>
                    @endif
                </div>
                {!! $products->appends(request()->query())->links() !!}

                <div class="row">
                    @foreach($products as $product)
                        <div class=" col-lg-2 col-md-4 col-sm-6 card">
                            <div class="card-body text-center pl-0 pr-0">

                                <h4 class="card-title">
                                    <input type="checkbox"
                                           class="product-checkbox"
                                           data-id="{{ $product->id }}"
                                           data-img="{{ asset($product->path_thumb) }}"
                                           style="position: absolute; right: 0px; top: 0px;"
                                    >
                                    {{$product->name}} - {{number_format($product->price)}}đ <br>

                                    <span class="right badge badge-info" style="font-size: 100%">{!! implode(' - ',$product->sizes->map(function($size) use ($product) {
     if ($size->name == 'DÀI QUẦN') {
                                                    $value = array_get(\App\Contracts\SaleConfig::MAPPING_SIZES, $size->name, $size->name) . '' . intval($size->value);

                                                    if ($duLai = $product->sizes->first(function ($item) { return $item->name == 'DƯ LAI';})) {
                                                        $value .= " (+". intval($duLai->value) .")";
                                                    }
                                                    return $value;

                                                }
                                                if ($size->name == 'DƯ LAI') {
                                                    return;
                                                }
                                                return isset(\App\Contracts\SaleConfig::MAPPING_SIZES[$size->name]) ? \App\Contracts\SaleConfig::MAPPING_SIZES[$size->name] . '' . intval($size->value) : null;
                                                })->filter()->toArray()) !!} </span><br>

                                </h4>


                                @if($product->path_thumb)
                                <a href="{{route('sale.product.detail', $product->id)}}">

                                    <img src="{{asset($product->path_thumb)}}" data-name="{{$product->name}}" alt="" style="width:100%">

                                </a>
                                @else
                                <a href="{{route('sale.product.detail', $product->id)}}">
                                    View
                                </a>
                                @endif

                            </div>

                            <!-- end card-body-->
                        </div>

                    @endforeach


                </div>
                {!! $products->appends(request()->query())->links() !!}


            </div> <!-- container-fluid -->
        </div>
        <!-- End Page-content -->


    </div>

@section('script')
<script>
$(document).ready(function() {
    const STORAGE_KEY = 'selectedProductImages';

    window.updateSelectedCount = function() {
        const stored = JSON.parse(localStorage.getItem(STORAGE_KEY)) || [];
        $('#selected-count').text(stored.length);
    };

    updateSelectedCount();

    // Event listener for checkbox changes
    $(document).on('change', '.product-checkbox', function() {
        const checkbox = this;
        let stored = JSON.parse(localStorage.getItem(STORAGE_KEY)) || [];
        const id = $(checkbox).data('id');
        const img = $(checkbox).data('img');

        if (checkbox.checked) {
            if (!stored.find(p => p.id == id)) {
                stored.push({ id, img });
            }
        } else {
            stored = stored.filter(p => p.id != id);
        }

        localStorage.setItem(STORAGE_KEY, JSON.stringify(stored));
        updateSelectedCount();
    });

    // SELECT ALL CURRENT PAGE
    $('#select-all-page').on('click', function() {
        $('.product-checkbox').each(function() {
            if (!this.checked) {
                this.checked = true;
                // Trigger change to update localStorage and count
                $(this).trigger('change');
            }
        });
        alert('Đã chọn tất cả sản phẩm trong trang này!');
    });
    
    // NATIVE SHARE DIRECTLY
    $('#native-share-btn').on('click', async function() {
        const stored = JSON.parse(localStorage.getItem(STORAGE_KEY)) || [];
        if (stored.length === 0) {
            alert("Bạn chưa chọn ảnh nào!");
            return;
        }

        if (!navigator.share) {
            alert('Trình duyệt của bạn không hỗ trợ tính năng chia sẻ nhanh.');
            return;
        }

        const $btn = $(this);
        const originalHtml = $btn.html();
        $btn.prop('disabled', true).html('<span class="spinner-border spinner-border-sm mr-1"></span> Đang chuẩn bị...');

        try {
            const files = [];
            for (const item of stored) {
                const response = await fetch(item.img);
                const blob = await response.blob();
                // Safari requires specific file naming and often works better with image/jpeg
                const file = new File([blob], `${item.id}.jpg`, { type: 'image/jpeg' });
                files.push(file);
            }

            if (navigator.canShare && navigator.canShare({ files: files })) {
                await navigator.share({
                    files: files,
                    title: 'Vintage Always True',
                    text: 'Gửi bạn ảnh sản phẩm',
                });
            } else {
                alert('Không thể chia sẻ nhiều file cùng lúc trên trình duyệt này.');
            }
        } catch (err) {
            console.error('Share failed:', err);
            if (err.name !== 'AbortError') {
                alert('Lỗi khi chia sẻ: ' + err.message);
            }
        } finally {
            $btn.prop('disabled', false).html(originalHtml);
        }
    });
});
</script>

<style>
    .btn-submit { min-width: 100px; }
    .spinner-border-sm { width: 1rem; height: 1rem; border-width: 0.15em; }
</style>
@endsection
