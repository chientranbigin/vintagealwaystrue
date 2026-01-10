@extends('sale.master')
@section('content')
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">


                <div class="card">
                    <div class="card-body">

                        <form method="POST" action="{{route('sale.product.detail', $product->id)}}"
                              enctype="multipart/form-data">
                            @csrf
                            <input type="hidden" name="id" value="{{$product->id}}">
                            <input type="hidden" name="backRoute" value="{{url()->previous()}}">

                            <div class="row">
                                <div class="col">
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Name</label>
                                        <input class="form-control" name="name" value="{{$product->name}}">

                                    </div>

                                    <div class="form-group">
                                        <label for="exampleFormControlSelect1">Type</label>
                                        <select class="form-control" id="exampleFormControlSelect1" name="type">
                                            @foreach(App\Contracts\SaleConfig::PRODUCT_TYPES as $type)
                                                <option {{$type == $product->type ? 'selected' : ''}}  value="{{$type}}">{{$type}}</option>

                                            @endforeach
                                        </select>
                                    </div>
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Thumb</label>


                                        <div class="custom-file">
                                            <input type="file" class="custom-file-input" id="customFile" name="banner">
                                            <label class="custom-file-label" for="customFile">Choose file</label>
                                        </div>
                                        <img src="{{asset($product->path_thumb)}}" alt="" height="300" class="mt-5">
                                    </div>
{{--                                    <div class="form-group w-100">--}}

{{--                                        <label for="exampleFormControlTextarea1">Import Price</label>--}}
{{--                                        <input type="number" class="form-control" name="import_price"--}}
{{--                                               value="{{intval($product->import_price)}}">--}}

{{--                                    </div>--}}
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Price</label>
                                        <input type="number" class="form-control" name="price"
                                               value="{{intval($product->price)}}">

                                    </div>
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Description</label>
                                        <textarea class="form-control" id="exampleFormControlTextarea1" rows="7" name="description" >{{$product->description}}</textarea>


                                    </div>

                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Color</label>
                                        <textarea class="form-control" id="exampleFormControlTextarea1" rows="7" name="color" >{{$product->color}}</textarea>


                                    </div>
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Fabric</label>
                                        <textarea class="form-control" id="exampleFormControlTextarea1" rows="7" name="fabric" >{{$product->fabric}}</textarea>


                                    </div>
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Detail fabric</label>
                                        <textarea class="form-control" id="exampleFormControlTextarea1" rows="7" name="detail_fabric" >{{$product->detail_fabric}}</textarea>


                                    </div>
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Title</label>
                                        <textarea class="form-control" id="exampleFormControlTextarea1" rows="7" name="title" >{{$product->title}}</textarea>


                                    </div>

{{--                                    <div class="form-group w-100">--}}

{{--                                        <label for="exampleFormControlTextarea1">URL Facebook</label>--}}
{{--                                        <input class="form-control" name="url_detail" value="{{$product->url_detail}}">--}}

{{--                                    </div>--}}
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Status</label><br>
                                        @if($product->status == 'AVAILABLE')
                                            <span class="right badge badge-info">{{$product->status}}</span>
                                        @elseif($product->status == 'ON_HOLD')
                                            <span class="right badge badge-warning">{{$product->status}}</span>
                                        @else
                                            <span class="right badge badge-danger">{{$product->status}}</span>
                                        @endif

                                    </div>


                                </div>
                                <div class="col">
                                    <h3>Sizes</h3>

                                    {{--                           @foreach(App\Contracts\SaleConfig::SIZES  as $key => $size)--}}
                                    <div class="row">
                                        <div class="col-6">
                                            @foreach(App\Contracts\SaleConfig::SIZES  as $key => $size)
                                                @if($key < 6)
                                                    <div class="form-group w-100">
                                                        <label for="exampleFormControlTextarea1">{{$size}}</label>
                                                        <input type="hidden" class="form-control"
                                                               name="sizes[{{$key}}][name]" value="{{$size}}">
                                                        <input type="number" class="form-control"
                                                               name="sizes[{{$key}}][value]"
                                                               value="{{intval($product->sizes->first(function($item) use ($size) {return $item->name  === $size;})->value ?? null) ?: null }}">

                                                    </div>
                                                @endif
                                            @endforeach
                                        </div>


                                        <div class="col-6">
                                            @foreach(App\Contracts\SaleConfig::SIZES  as $key => $size)
                                                @if($key >= 6)
                                                    <div class="form-group w-100">
                                                        <label for="exampleFormControlTextarea1">{{$size}}</label>
                                                        <input type="hidden" class="form-control"
                                                               name="sizes[{{$key}}][name]" value="{{$size}}">
                                                        <input type="number" class="form-control"
                                                               name="sizes[{{$key}}][value]"
                                                               value="{{intval($product->sizes->first(function($item) use ($size) {return $item->name  === $size;})->value ?? null) ?: null }}">

                                                    </div>
                                                @endif
                                            @endforeach

                               </div>
{{--                                <div class="col">--}}
{{--                                    <div class="form-group w-100">--}}
{{--                                        <label for="exampleFormControlTextarea1">{{$size}}</label>--}}
{{--                                        <input type="hidden" class="form-control" name="sizes[{{$key}}][name]"  value="{{$size}}">--}}
{{--                                        <input type="number" class="form-control" name="sizes[{{$key}}][value]" value="{{intval($product->sizes->first(function($item) use ($size) {return $item->name  === $size;})->value ?? null) ?: null }}">--}}
{{--                                    </div>--}}
{{--                                </div>--}}
                            </div>

{{--                           @endforeach--}}
                                </div>
                            </div>
                   
                            <button type="submit" class="btn btn-secondary waves-effect waves-light btn-submit" >Save</button>

                        </form>

                    </div>
                    <!-- end card-body-->
                </div>


            </div> <!-- container-fluid -->
        </div>
        <!-- End Page-content -->



    </div>


@endsection
