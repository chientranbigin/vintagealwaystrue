

@extends('consumer.master')
@section('content')

    <!-- Breadcrumb Start -->
    <div class="container-fluid">
        <div class="row px-xl-5">
            <div class="col-12">
                <nav class="breadcrumb bg-light mb-30">
                    <p class="text-primary mb-0"><b> Thông tin và trạng thái sản phẩm sẽ được cập nhật liên tục chính xác đến 99%. Để đảm bảo cập nhật thông tin mới nhất, vui lòng refresh page (F5) để thông tin được cập nhật. </b></p>
                </nav>

            </div>
        </div>
    </div>
    <!-- Breadcrumb End -->


    <!-- Shop Start -->
    <div class="container-fluid">
        <div class="row px-xl-5">
            <!-- Shop Sidebar Start -->
            <div class="col-lg-3 col-md-4">
{{--                <h5 class="section-title position-relative text-uppercase mb-3"><span class="bg-secondary pr-3">Filter by type</span></h5>--}}
                <div class="bg-light p-4 mb-30">
                    <p><b>Nhập mã sản phẩm, loại sản phẩm hoặc thông số cụ thể tại đây để dễ dàng tìm sản phẩm. <br>
                        Các mẫu hết hàng sẽ không hiện trên hệ thống.</b></p>
                    <form method="POST" action="{{ route('sale.searchAvailableProducts', ['urlSearch' => 'customer.products']) }}" enctype="multipart/form-data">
                        @csrf

                        <div class="row mb-3 mt-3">


                            <div class="col-sm-12 col-md-12 form-group">
                                <label>Mã sản phẩm, có thể tìm nhiều sản phẩm 1 lần (vdu TS345, BZ91, TI23)</label>
                                <input id="input-search" class="form-control" name="search" value="{{request()->search}}" placeholder="Name Products">

                            </div>
{{--                            <div class=" col-sm-12 col-md-12   form-group">--}}
{{--                                <label>Trạng thái sản phẩm (sẵn hàng, đã bán¸...)</label>--}}

{{--                                <select class="custom-select rounded-0" id="exampleSelectRounded0" name="status" >--}}
{{--                                    <option {{!request()->status ? 'selected' : ''}} value="">Tất cả</option>--}}

{{--                                    <option {{request()->status == 'AVAILABLE' ? 'selected' : ''}} value="AVAILABLE">Sẵn hàng</option>--}}
{{--                                    <option {{request()->status == 'SOLD' ? 'selected' : ''}} value="SOLD">Đã bán</option>--}}
{{--                                </select>--}}

{{--                            </div>--}}

                            <div class=" col-sm-12 col-md-12   form-group">
                                <label>Loại sản phẩm (trousers, shirt, ...)</label>

                                <select class="custom-select rounded-0" id="exampleSelectRounded0"  name="type" >
                                    <option {{!request()->type ? 'selected' : ''}} value="">Tất cả</option>


                                @foreach(\App\Contracts\SaleConfig::PRODUCT_TYPES as $type)
                                        <option {{request()->type == $type ? 'selected' : ''}}>{{$type}}</option>
                                    @endforeach
                                </select>

                            </div>
{{--                            <hr style="width: 100%; margin-top: 30px">--}}

                            <div class="col-sm-12 col-md-12  mt-2">
                                <h5><b>Số đo </b></h5>
                                <label>Vdu: eo quần từ from 43 - to 45, dài quần from 95 - to 100 </label>

                                <div class="row ">

                                    <div class="col-12 ">
{{--                                        <label>From</label>--}}
                                        <div class="d-flex">
                                            <select class="custom-select rounded-0" id="exampleSelectRounded0" name="from_search_key" placeholder="From">
                                                <option {{!request()->from_search_key ? 'selected' : ''}}></option>

                                                @foreach(\App\Contracts\SaleConfig::SIZES as $size)
                                                    <option {{request()->from_search_key == $size ? 'selected' : ''}}>{{$size}}</option>
                                                @endforeach
                                            </select>
                                            <input id="from_search_value" class="form-control" name="from_search_value" placeholder="From"
                                                   value="{{request()->from_search_value}}">

                                            <input id="to_search_value" class="form-control" name="to_search_value" placeholder="To"
                                                   value="{{request()->to_search_value}}">
                                        </div>

                                    </div>

{{--                                    <div class="col-6">--}}
{{--                                        <label class="">To</label>--}}
{{--                                        <div class="d-flex">--}}
{{--                                            <select class="custom-select rounded-0" id="exampleSelectRounded0" name="to_search_key" placeholder="To">--}}
{{--                                                <option {{!request()->to_search_key ? 'selected' : ''}}></option>--}}

{{--                                                @foreach(\App\Contracts\SaleConfig::SIZES as $size)--}}
{{--                                                    <option {{request()->to_search_key == $size ? 'selected' : ''}}>{{$size}}</option>--}}
{{--                                                @endforeach--}}
{{--                                            </select>--}}
{{--                                            <input id="to_search_value" class="form-control" name="to_search_value" placeholder="To Value"--}}
{{--                                                   value="{{request()->to_search_value}}">--}}
{{--                                        </div>--}}

{{--                                    </div>--}}
                                </div>

                            </div>
                            <div class="col-sm-12 col-md-12 mt-3">

                                <div class="row ">

                                    <div class="col-12 ">
{{--                                        <label>From</label>--}}
                                        <div class="d-flex">
                                            <select class="custom-select rounded-0" id="exampleSelectRounded0" name="from_search_key_2" placeholder="From">
                                                <option {{!request()->from_search_key_2 ? 'selected' : ''}}></option>

                                                @foreach(\App\Contracts\SaleConfig::SIZES as $size)
                                                    <option {{request()->from_search_key_2 == $size ? 'selected' : ''}}>{{$size}}</option>
                                                @endforeach
                                            </select>
                                            <input id="from_search_value_2" class="form-control" name="from_search_value_2" placeholder="From"
                                                   value="{{request()->from_search_value_2}}">
                                            <input id="to_search_value_2" class="form-control" name="to_search_value_2" placeholder="To"
                                                   value="{{request()->to_search_value_2}}">
                                        </div>

                                    </div>

{{--                                    <div class="col-6 ">--}}
{{--                                        <label class="">To</label>--}}
{{--                                        <div class="d-flex">--}}
{{--                                            <select class="custom-select rounded-0" id="exampleSelectRounded0" name="to_search_key_2" placeholder="To">--}}
{{--                                                <option {{!request()->to_search_key_2 ? 'selected' : ''}}></option>--}}

{{--                                                @foreach(\App\Contracts\SaleConfig::SIZES as $size)--}}
{{--                                                    <option {{request()->to_search_key_2 == $size ? 'selected' : ''}}>{{$size}}</option>--}}
{{--                                                @endforeach--}}
{{--                                            </select>--}}
{{--                                            <input id="to_search_value_2" class="form-control" name="to_search_value_2" placeholder="To Value"--}}
{{--                                                   value="{{request()->to_search_value_2}}">--}}
{{--                                        </div>--}}

{{--                                    </div>--}}
                                </div>

                            </div>
{{--                            <div class="col-sm-12 col-md-6  ">--}}
{{--                                <div class="row mt-3">--}}
{{--                                    <div class="col d-flex">--}}

{{--                                        <select class="custom-select rounded-0" id="exampleSelectRounded0" name="from_search_key_2" placeholder="From">--}}
{{--                                            <option {{!request()->from_search_key_2 ? 'selected' : ''}}></option>--}}

{{--                                            @foreach(\App\Contracts\SaleConfig::SIZES as $size)--}}
{{--                                                <option {{request()->from_search_key_2 == $size ? 'selected' : ''}}>{{$size}}</option>--}}
{{--                                            @endforeach--}}
{{--                                        </select>--}}
{{--                                        <input id="from_search_value" class="form-control" name="from_search_value_2" placeholder="From Value"--}}
{{--                                               value="{{request()->from_search_value_2}}">--}}
{{--                                    </div>--}}
{{--                                    <div class="col d-flex">--}}

{{--                                        <select class="custom-select rounded-0" id="exampleSelectRounded0" name="to_search_key_2" placeholder="To">--}}
{{--                                            <option {{!request()->to_search_key_2 ? 'selected' : ''}}></option>--}}

{{--                                            @foreach(\App\Contracts\SaleConfig::SIZES as $size)--}}
{{--                                                <option {{request()->to_search_key_2 == $size ? 'selected' : ''}}>{{$size}}</option>--}}
{{--                                            @endforeach--}}
{{--                                        </select>--}}
{{--                                        <input id="to_search_value" class="form-control" name="to_search_value_2" placeholder="To Value"--}}
{{--                                               value="{{request()->to_search_value_2}}">--}}
{{--                                    </div>--}}
{{--                                </div>--}}

{{--                            </div>--}}
                            <div class="col-12">

                                <button type="submit" class="btn btn-primary  waves-effect waves-light  mt-3 mr-3 my-3 py-3 btn-block">
                                    Search
                                </button>
                                <a href="{{route('customer.products')}}"  class="btn btn-secondary waves-effect waves-light  mt-3  mt-3 mr-3 my-3 py-3 btn-block">
                                    Clear
                                </a>

{{--                                <span id="download-all-images" class="btn btn-light waves-effect waves-light btn-submit mt-3 float-right">--}}
{{--                                Download--}}
{{--                            </span>--}}

{{--                                <a href="{{route('sale.availableProducts', ['showAll' => true])}}"  class="btn btn-light waves-effect waves-light btn-submit mt-3 float-right mr-3">--}}
{{--                                    Show All--}}
{{--                                </a>--}}
                            </div>
                        </div>

                    </form>

                </div>
                <!-- Price Start -->
{{--                <h5 class="section-title position-relative text-uppercase mb-3"><span class="bg-secondary pr-3">Filter by price</span></h5>--}}
{{--                <div class="bg-light p-4 mb-30">--}}
{{--                    <form>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between mb-3">--}}
{{--                            <input type="checkbox" class="custom-control-input" checked id="price-all">--}}
{{--                            <label class="custom-control-label" for="price-all">All Price</label>--}}
{{--                            <span class="badge border font-weight-normal">1000</span>--}}
{{--                        </div>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between mb-3">--}}
{{--                            <input type="checkbox" class="custom-control-input" id="price-1">--}}
{{--                            <label class="custom-control-label" for="price-1">$0 - $100</label>--}}
{{--                            <span class="badge border font-weight-normal">150</span>--}}
{{--                        </div>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between mb-3">--}}
{{--                            <input type="checkbox" class="custom-control-input" id="price-2">--}}
{{--                            <label class="custom-control-label" for="price-2">$100 - $200</label>--}}
{{--                            <span class="badge border font-weight-normal">295</span>--}}
{{--                        </div>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between mb-3">--}}
{{--                            <input type="checkbox" class="custom-control-input" id="price-3">--}}
{{--                            <label class="custom-control-label" for="price-3">$200 - $300</label>--}}
{{--                            <span class="badge border font-weight-normal">246</span>--}}
{{--                        </div>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between mb-3">--}}
{{--                            <input type="checkbox" class="custom-control-input" id="price-4">--}}
{{--                            <label class="custom-control-label" for="price-4">$300 - $400</label>--}}
{{--                            <span class="badge border font-weight-normal">145</span>--}}
{{--                        </div>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between">--}}
{{--                            <input type="checkbox" class="custom-control-input" id="price-5">--}}
{{--                            <label class="custom-control-label" for="price-5">$400 - $500</label>--}}
{{--                            <span class="badge border font-weight-normal">168</span>--}}
{{--                        </div>--}}
{{--                    </form>--}}
{{--                </div>--}}
{{--                <!-- Price End -->--}}

{{--                <!-- Color Start -->--}}
{{--                <h5 class="section-title position-relative text-uppercase mb-3"><span class="bg-secondary pr-3">Filter by color</span></h5>--}}
{{--                <div class="bg-light p-4 mb-30">--}}
{{--                    <form>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between mb-3">--}}
{{--                            <input type="checkbox" class="custom-control-input" checked id="color-all">--}}
{{--                            <label class="custom-control-label" for="price-all">All Color</label>--}}
{{--                            <span class="badge border font-weight-normal">1000</span>--}}
{{--                        </div>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between mb-3">--}}
{{--                            <input type="checkbox" class="custom-control-input" id="color-1">--}}
{{--                            <label class="custom-control-label" for="color-1">Black</label>--}}
{{--                            <span class="badge border font-weight-normal">150</span>--}}
{{--                        </div>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between mb-3">--}}
{{--                            <input type="checkbox" class="custom-control-input" id="color-2">--}}
{{--                            <label class="custom-control-label" for="color-2">White</label>--}}
{{--                            <span class="badge border font-weight-normal">295</span>--}}
{{--                        </div>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between mb-3">--}}
{{--                            <input type="checkbox" class="custom-control-input" id="color-3">--}}
{{--                            <label class="custom-control-label" for="color-3">Red</label>--}}
{{--                            <span class="badge border font-weight-normal">246</span>--}}
{{--                        </div>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between mb-3">--}}
{{--                            <input type="checkbox" class="custom-control-input" id="color-4">--}}
{{--                            <label class="custom-control-label" for="color-4">Blue</label>--}}
{{--                            <span class="badge border font-weight-normal">145</span>--}}
{{--                        </div>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between">--}}
{{--                            <input type="checkbox" class="custom-control-input" id="color-5">--}}
{{--                            <label class="custom-control-label" for="color-5">Green</label>--}}
{{--                            <span class="badge border font-weight-normal">168</span>--}}
{{--                        </div>--}}
{{--                    </form>--}}
{{--                </div>--}}
{{--                <!-- Color End -->--}}

{{--                <!-- Size Start -->--}}
{{--                <h5 class="section-title position-relative text-uppercase mb-3"><span class="bg-secondary pr-3">Filter by size</span></h5>--}}
{{--                <div class="bg-light p-4 mb-30">--}}
{{--                    <form>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between mb-3">--}}
{{--                            <input type="checkbox" class="custom-control-input" checked id="size-all">--}}
{{--                            <label class="custom-control-label" for="size-all">All Size</label>--}}
{{--                            <span class="badge border font-weight-normal">1000</span>--}}
{{--                        </div>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between mb-3">--}}
{{--                            <input type="checkbox" class="custom-control-input" id="size-1">--}}
{{--                            <label class="custom-control-label" for="size-1">XS</label>--}}
{{--                            <span class="badge border font-weight-normal">150</span>--}}
{{--                        </div>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between mb-3">--}}
{{--                            <input type="checkbox" class="custom-control-input" id="size-2">--}}
{{--                            <label class="custom-control-label" for="size-2">S</label>--}}
{{--                            <span class="badge border font-weight-normal">295</span>--}}
{{--                        </div>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between mb-3">--}}
{{--                            <input type="checkbox" class="custom-control-input" id="size-3">--}}
{{--                            <label class="custom-control-label" for="size-3">M</label>--}}
{{--                            <span class="badge border font-weight-normal">246</span>--}}
{{--                        </div>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between mb-3">--}}
{{--                            <input type="checkbox" class="custom-control-input" id="size-4">--}}
{{--                            <label class="custom-control-label" for="size-4">L</label>--}}
{{--                            <span class="badge border font-weight-normal">145</span>--}}
{{--                        </div>--}}
{{--                        <div class="custom-control custom-checkbox d-flex align-items-center justify-content-between">--}}
{{--                            <input type="checkbox" class="custom-control-input" id="size-5">--}}
{{--                            <label class="custom-control-label" for="size-5">XL</label>--}}
{{--                            <span class="badge border font-weight-normal">168</span>--}}
{{--                        </div>--}}
{{--                    </form>--}}
{{--                </div>--}}
                <!-- Size End -->
            </div>
            <!-- Shop Sidebar End -->


            <!-- Shop Product Start -->
            <div class="col">
                {!! $products->appends(request()->query())->links() !!}

                <div class="row pb-3">
{{--                    <div class="col-12 pb-1">--}}
{{--                        <div class="d-flex align-items-center justify-content-between mb-4">--}}
{{--                            <div>--}}
{{--                                <button class="btn btn-sm btn-light"><i class="fa fa-th-large"></i></button>--}}
{{--                                <button class="btn btn-sm btn-light ml-2"><i class="fa fa-bars"></i></button>--}}
{{--                            </div>--}}
{{--                            <div class="ml-2">--}}
{{--                                <div class="btn-group">--}}
{{--                                    <button type="button" class="btn btn-sm btn-light dropdown-toggle" data-toggle="dropdown">Sorting</button>--}}
{{--                                    <div class="dropdown-menu dropdown-menu-right">--}}
{{--                                        <a class="dropdown-item" href="#">Latest</a>--}}
{{--                                        <a class="dropdown-item" href="#">Popularity</a>--}}
{{--                                        <a class="dropdown-item" href="#">Best Rating</a>--}}
{{--                                    </div>--}}
{{--                                </div>--}}
{{--                                <div class="btn-group ml-2">--}}
{{--                                    <button type="button" class="btn btn-sm btn-light dropdown-toggle" data-toggle="dropdown">Showing</button>--}}
{{--                                    <div class="dropdown-menu dropdown-menu-right">--}}
{{--                                        <a class="dropdown-item" href="#">10</a>--}}
{{--                                        <a class="dropdown-item" href="#">20</a>--}}
{{--                                        <a class="dropdown-item" href="#">30</a>--}}
{{--                                    </div>--}}
{{--                                </div>--}}
{{--                            </div>--}}
{{--                        </div>--}}
{{--                    </div>--}}
                    @if($products->isEmpty())
                        <p style="width: 100%; text-align: center">Không tìm được sản phẩm phù hợp. <br>
                            Vui lòng bấm nút "Clear" để hiện tất cả sản phẩm.</p>
                    @endif
                    @foreach($products as $product)
                        <div class="col-lg-3 col-md-4 col-sm-6 pb-1">
                            <div class="product-item bg-light mb-4">
                                <div class="product-img position-relative overflow-hidden">
                                    <img class="img-fluid w-100" src="{{asset($product->path_thumb)}}" alt="">
{{--                                    <div class="product-action">--}}
{{--                                        <a class="btn btn-outline-dark btn-square" href=""><i class="fa fa-shopping-cart"></i></a>--}}
{{--                                        <a class="btn btn-outline-dark btn-square" href=""><i class="far fa-heart"></i></a>--}}
{{--                                        <a class="btn btn-outline-dark btn-square" href=""><i class="fa fa-sync-alt"></i></a>--}}
{{--                                        <a class="btn btn-outline-dark btn-square" href=""><i class="fa fa-search"></i></a>--}}
{{--                                    </div>--}}
                                </div>
                                <div style="width: 100%"  @if($product->status == 'AVAILABLE') class=" bg-primary text-center text-white py-1" @else class=" bg-danger text-center text-primary py-1" @endif >
                                    <span>{{$product->status == 'AVAILABLE' ? 'Sẵn hàng' : 'Đã bán'}}</span>


                                </div>
                                <div class="text-center py-3">
                                    <a class="h6 text-decoration-none text-truncate mb-0" >{{ucfirst(strtolower($product->type))}} {{$product->name}}</a>

                                    <div class="d-flex align-items-center justify-content-center mt-2">
                                       <p style="font-size: 14px; margin-bottom: 0">{!! implode(' - ',$product->sizes->map(function($size) use ($product) {
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
                                                })->filter()->toArray()) !!} </p>
                                    </div>
{{--                                    <div class="d-flex align-items-center justify-content-center mb-1">--}}
{{--                                        <small class="fa fa-star text-primary mr-1"></small>--}}
{{--                                        <small class="fa fa-star text-primary mr-1"></small>--}}
{{--                                        <small class="fa fa-star text-primary mr-1"></small>--}}
{{--                                        <small class="fa fa-star text-primary mr-1"></small>--}}
{{--                                        <small class="fa fa-star text-primary mr-1"></small>--}}
{{--                                        <small>(99)</small>--}}
{{--                                    </div>--}}
                                </div>
                            </div>
                        </div>
                    @endforeach


{{--                    <div class="col-12">--}}
{{--                        <nav>--}}
{{--                            <ul class="pagination justify-content-center">--}}
{{--                                <li class="page-item disabled"><a class="page-link" href="#">Previous</span></a></li>--}}
{{--                                <li class="page-item active"><a class="page-link" href="#">1</a></li>--}}
{{--                                <li class="page-item"><a class="page-link" href="#">2</a></li>--}}
{{--                                <li class="page-item"><a class="page-link" href="#">3</a></li>--}}
{{--                                <li class="page-item"><a class="page-link" href="#">Next</a></li>--}}
{{--                            </ul>--}}
{{--                        </nav>--}}
{{--                    </div>--}}
                </div>
                {!! $products->appends(request()->query())->links() !!}

            </div>
            <!-- Shop Product End -->
        </div>
    </div>

@endsection


