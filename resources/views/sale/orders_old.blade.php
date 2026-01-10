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
                                        <td><a href="{{route('sale.order.detail.delete', $order->id)}}"
                                               class="btn btn-secondary waves-effect waves-light btn-submit  "
                                               style="    color: #fff !important;">Delete </a>


                                            <a href="{{route('sale.order.complete', $order->id)}}"
                                               class="btn btn-light waves-effect waves-light btn-submit mt-3">Complete</a>
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


    </div>

@endsection
