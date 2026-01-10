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
                    <input type="hidden" class="form-control"
                           name="urlSearch" value="{{\Illuminate\Support\Facades\Route::currentRouteName()}}">
                    <div class="row mb-3">


                        <div class="col-sm-12 col-md-6">
                            <input id="input-search" class="form-control" name="search" value="{{request()->search}}"
                                   placeholder="Name Products">

                        </div>


                        <div class=" col-sm-12 col-md-3 d-flex ">
                            <select class="custom-select rounded-0" id="exampleSelectRounded0" name="type"
                                    style="width: 300px">
                                <option {{!request()->type ? 'selected' : ''}}></option>

                                @foreach(\App\Contracts\SaleConfig::PRODUCT_TYPES as $type)
                                    <option {{request()->type == $type ? 'selected' : ''}}>{{$type}}</option>
                                @endforeach
                            </select>

                        </div>
                        <div class=" col-sm-12 col-md-3 d-flex ">
                            <select class="custom-select rounded-0" id="exampleSelectRounded0" name="status"
                                    style="width: 300px">
                                <option {{!request()->type ? 'selected' : ''}}></option>

                                <option {{request()->type == 'AVAILABLE' ? 'selected' : ''}}>AVAILABLE</option>
                                <option {{request()->type == 'SOLD' ? 'selected' : ''}}>SOLD</option>
                                <option {{request()->type == 'ON_HOLD' ? 'selected' : ''}}>ON_HOLD</option>
                            </select>

                        </div>
                        <div class="col-sm-12 col-md-6  ">
                            <div class="row mt-3">
                                <div class="col d-flex">

                                    <select class="custom-select rounded-0" id="exampleSelectRounded0"
                                            name="from_search_key" placeholder="From">
                                        <option {{!request()->from_search_key ? 'selected' : ''}}></option>

                                        @foreach(\App\Contracts\SaleConfig::SIZES as $size)
                                            <option {{request()->from_search_key == $size ? 'selected' : ''}}>{{$size}}</option>
                                        @endforeach
                                    </select>
                                    <input id="from_search_value" class="form-control" name="from_search_value"
                                           placeholder="From Value"
                                           value="{{request()->from_search_value}}">
                                    <input id="to_search_value" class="form-control" name="to_search_value"
                                           placeholder="To Value"
                                           value="{{request()->to_search_value}}">
                                </div>
                                <div class="col ">
                                    <div class="form-group ">
                                        <div class="d-flex">
                                            <input type="date" class="form-control" value="{{request()->sold_from}}"
                                                   placeholder="From" name="sold_from">
                                            <input type="date" class="form-control" value="{{request()->sold_to}}"
                                                   placeholder="To" name="sold_to">
                                        </div>

                                    </div>
                                    {{--                                    <select class="custom-select rounded-0" id="exampleSelectRounded0" name="to_search_key" placeholder="To">--}}
                                    {{--                                        <option {{!request()->to_search_key ? 'selected' : ''}}></option>--}}

                                    {{--                                        @foreach(\App\Contracts\SaleConfig::SIZES as $size)--}}
                                    {{--                                            <option {{request()->to_search_key == $size ? 'selected' : ''}}>{{$size}}</option>--}}
                                    {{--                                        @endforeach--}}
                                    {{--                                    </select>--}}

                                </div>
                            </div>

                        </div>
                        <div class="col-12">

                            <button type="submit"
                                    class="btn btn-secondary waves-effect waves-light btn-submit mt-3 mr-3">
                                Search
                            </button>
                            <a href="{{route('sale.product')}}"
                               class="btn btn-light waves-effect waves-light btn-submit mt-3">
                                Clear
                            </a>

                            <span id="download-all-images"
                                  class="btn btn-light waves-effect waves-light btn-submit mt-3 float-right">
                                Download
                            </span>

                            <a href="{{route('sale.product', ['showAll' => true])}}"
                               class="btn btn-light waves-effect waves-light btn-submit mt-3 float-right mr-3">
                                Show All
                            </a>
                        </div>
                    </div>

                </form>
                <div class="card">
                    <div class="card-body">
                        <h4 class="card-title">Products</h4>
                        <a href="{{route('sale.product.create')}}"
                           class="btn btn-secondary waves-effect waves-light btn-submit  mb-3"
                           style="    color: #fff !important;">Create </a>

                        <div class="table-responsive">
                            <div class="py-3">
                                <b>Showing {{$from}} to {{$to}} of {{$total}} entries</b>
                            </div>
                            {!! $products->appends(request()->query())->links() !!}

                            <table class="table mb-0">
                                <thead>
                                <tr>

                                    <th>Type</th>
                                    <th>Name</th>
                                    <th>Thumb</th>
                                    <th>Price</th>
                                    {{--                                    <th>Sizes</th>--}}
                                    <th>Detail</th>
                                    <th>Status</th>
                                    <th>Action</th>
                                </tr>
                                </thead>
                                <tbody>
                                @foreach($products as $product)
                                    <tr>
                                        <td>{{$product->type}}</td>
                                        <td>
                                            <a href="{{route('sale.product.detail', $product->id)}}">{{$product->name}}</a>
                                        </td>
                                        <td><img src="{{asset($product->path_thumb)}}" alt=""
                                                 data-name="{{$product->name}}" style="width:200px"></td>
                                        <td>{{number_format($product->price)}}đ</td>
                                        {{--                                        <td>{!! implode('<br>',$product->sizes->map(function($size) {--}}
                                        {{--                                            return $size->name . ': ' . intval($size->value) ;--}}
                                        {{--                                            })->toArray()) !!}</td>--}}
                                        <td>
                                            {{$product->name}} {!! implode(' - ',$product->sizes->map(function($size) use ($product) {
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

                                                if ($size->name == 'EO QUẦN') {
                                                                                                    return isset(\App\Contracts\SaleConfig::MAPPING_SIZES[$size->name]) ? \App\Contracts\SaleConfig::MAPPING_SIZES[$size->name] . '' . intval($size->value *2) : null;

                                                }
                                                return isset(\App\Contracts\SaleConfig::MAPPING_SIZES[$size->name]) ? \App\Contracts\SaleConfig::MAPPING_SIZES[$size->name] . '' . intval($size->value) : null;
                                                })->filter()->toArray()) !!}<br>


                                            @if($product->description)
                                                <br> {{$product->description}}
                                            @endif
                                            <br>
                                            {{$product->name}} {!! implode(' - ',$product->sizes->map(function($size) use ($product) {
    return array_get(\App\Contracts\SaleConfig::MAPPING_NAME_SIZE, $size->name, $size->name) . ' ' . intval($size->value);

                                                })->filter()->toArray()) !!} <br>
                                            {{$product->price/1000}}k


                                            {{--                                            @if($product->type == 'TROUSERS')--}}
                                            {{--                                            {{strtoupper($product->title ?: $product->color)}} TROUSERS<br>--}}
                                            {{--                                            -----------------------------------------<br>--}}
                                            {{--                                                @if($product->description)--}}
                                            {{--                                                  {!! nl2br($product->description) !!}   <br>--}}
                                            {{--                                                    -----------------------------------------<br>--}}
                                            {{--                                                @endif--}}
                                            {{--                                            Thông số: {!! implode(' - ',$product->sizes->map(function($size) use ($product) {--}}
                                            {{--    return array_get(\App\Contracts\SaleConfig::MAPPING_NAME_SIZE, $size->name, $size->name) . ' ' . intval($size->value);--}}

                                            {{--                                                })->filter()->toArray()) !!}<br>--}}
                                            {{--                                            Màu: {{$product->color}}<br>--}}
                                            {{--                                            Chất liệu vải: {{$product->fabric}}<br>--}}
                                            {{--                                            Độ dày vải: {{$product->detail_fabric}}<br>--}}
                                            {{--                                            Độ mới: mới cao, không lỗi<br>--}}
                                            {{--                                            Mã số: {{$product->name}}<br>--}}
                                            {{--                                            ----------------------------------------<br>--}}
                                            {{--                                            Giá: {{$product->price/1000}}k<br><br>--}}

                                            {{--                                            #vintageshopall {{array_get(\App\Contracts\SaleConfig::MAPPING_HASH_TAG, $product->type)}}--}}
                                            {{--                                                #size{{ceil((($product->sizes->first(function($size){ return $size->name == 'EO QUẦN';})->value ?? 0) * 2) / 2.54)}}vintageshoptrousers--}}
                                            {{--                                            @endif--}}
                                        </td>
                                        <td>
                                            @if($product->status == 'AVAILABLE')
                                                <span class="right badge badge-info">{{$product->status}}</span>
                                            @elseif($product->status == 'ON_HOLD')
                                                <span class="right badge badge-warning">{{$product->status}}</span>
                                            @else
                                                <span class="right badge badge-danger">{{$product->status}}</span>
                                            @endif
                                        </td>
                                        <td><a href="{{route('sale.product.detail.delete', $product->id)}}"
                                               class="btn btn-secondary waves-effect waves-light btn-submit  "
                                               style="    color: #fff !important;">Delete </a></td>
                                    </tr>
                                @endforeach


                                </tbody>
                            </table>
                            {!! $products->appends(request()->query())->links() !!}
                        </div>

                    </div>
                    <!-- end card-body-->
                </div>

            </div> <!-- container-fluid -->
        </div>
        <!-- End Page-content -->



    </div>


@endsection
