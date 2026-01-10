@extends('sale.master')
@section('content')
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">


                <div class="card">
                    <div class="card-body">

                        <form method="POST" action="{{route('sale.order.detail', $order->id)}}"
                              enctype="multipart/form-data">
                            @csrf
                            <input type="hidden" name="id" value="{{$order->id}}">
                            <input type="hidden" name="backRoute" value="{{url()->previous()}}">

                            <div class="row">
                                <div class="col">
                                    <button type="submit"
                                            class="btn btn-secondary waves-effect waves-light btn-submit mb-5">Save
                                    </button>
                                    <div class="form-group">
                                        <label for="exampleFormControlSelect1">Type</label>
                                        <select class="form-control" id="exampleFormControlSelect1" name="status">
                                            @foreach(App\Contracts\SaleConfig::ORDER_STATUSES as $status)
                                                <option {{$status == $order->status ? 'selected' : ''}}  value="{{$status}}">{{$status}}</option>

                                            @endforeach
                                        </select>
                                    </div>
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Name</label>
                                        <input class="form-control" name="name" value="{{$order->name}}">

                                    </div>
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Phone</label>
                                        <input class="form-control" name="phone" value="{{$order->phone}}">

                                    </div>
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Address</label>
                                        <input class="form-control" name="address" value="{{$order->address}}">

                                    </div>
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Note</label>
                                        <input class="form-control" name="note" value="{{$order->note}}">

                                    </div>
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Description</label>
                                        <textarea class="form-control" id="exampleFormControlTextarea1" rows="7" name="desc" >{{$order->desc}}</textarea>


                                    </div>

                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Total Amount</label>
                                        <input type="number" class="form-control" name="additional_amount"
                                               value="{{intval($order->total_amount)}}">

                                    </div>
                                    {{--<div class="form-group w-100 ml-4">--}}

                                    {{--<input class="form-check-input " type="checkbox" name="is_freeship" {{$order->is_freeship ? 'checked' : ''}}>--}}
                                    {{--Is Freeship--}}
                                    {{--</div>--}}
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Sale Amount</label>
                                        <input type="number" class="form-control" name="freeship_amount"
                                               value="{{intval($order->freeship_amount)}}">

                                    </div>
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Deposit Amount</label>
                                        <input type="number" class="form-control" name="deposit_amount"
                                               value="{{intval($order->deposit_amount)}}">

                                    </div>
                                    <div class="form-group w-100 ml-4">

                                        <input class="form-check-input " type="checkbox"
                                               name="is_paid" {{$order->is_paid ? 'checked' : ''}} >
                                        Is Paid Full
                                    </div>
                                </div>
                                <div class="col">


                                    <table class="table mb-0">
                                        <thead>
                                        <tr>
                                            <th>Checked</th>
                                            <th>Type</th>
                                            <th>Name</th>
                                            {{--                                                <th>Thumb</th>--}}
                                            <th>Price</th>
                                            {{--                                                <th>Sizes</th>--}}


                                        </tr>
                                        </thead>
                                        <tbody>
                                        @foreach ($order->products as $product)
                                            <tr>
                                                <td><input class="form-check-input ml-4" type="checkbox"
                                                           name="products[]"
                                                           value="{{ $product->id }}" checked></td>

                                                <td>{{ $product->type }}</td>
                                                <td>
                                                    <a
                                                            href="{{ route('sale.product.detail', $product->id) }}">{{ $product->name }}</a>
                                                </td>
                                                {{--                                                    <td><img src="{{ asset($product->path_thumb) }}" alt=""--}}
                                                {{--                                                            height="100"></td>--}}
                                                <td>{{ number_format($product->price) }}đ</td>
                                                {{--                                                    <td>{!! implode(--}}
                                                {{--                                                        '<br>',--}}
                                                {{--                                                        $product->sizes->map(function ($size) {--}}
                                                {{--                                                                return $size->name . ': ' . intval($size->value);--}}
                                                {{--                                                            })->toArray(),--}}
                                                {{--                                                    ) !!}</td>--}}
                                            </tr>
                                        @endforeach
                                        @foreach ($products as $product)
                                            <tr>
                                                <td><input class="form-check-input ml-4" type="checkbox"
                                                           name="products[]"
                                                           value="{{ $product->id }}"></td>

                                                <td>{{ $product->type }}</td>
                                                <td>
                                                    <a
                                                            href="{{ route('sale.product.detail', $product->id) }}">{{ $product->name }}</a>
                                                </td>
                                                {{--                                                    <td><img src="{{ asset($product->path_thumb) }}" alt=""--}}
                                                {{--                                                            height="100"></td>--}}
                                                <td>{{ number_format($product->price) }}đ</td>
                                                {{--                                                    <td>{!! implode(--}}
                                                {{--                                                        '<br>',--}}
                                                {{--                                                        $product->sizes->map(function ($size) {--}}
                                                {{--                                                                return $size->name . ': ' . intval($size->value);--}}
                                                {{--                                                            })->toArray(),--}}
                                                {{--                                                    ) !!}</td>--}}
                                            </tr>
                                        @endforeach


                                        </tbody>
                                    </table>
                                </div>
                            </div>

{{--                            <button type="submit" class="btn btn-secondary waves-effect waves-light btn-submit" >Save</button>--}}

                        </form>

                    </div>
                    <!-- end card-body-->
                </div>


            </div> <!-- container-fluid -->
        </div>
        <!-- End Page-content -->



    </div>


@endsection
