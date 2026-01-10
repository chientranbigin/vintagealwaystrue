@extends('sale.master')
@section('content')
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">

                <div class="card">
                    <div class="card-body">
                        <h4 class="card-title">Summary</h4>
                        <div class="row">


                            <div class="col">
                                <h2>All</h2> <br>

{{--                                Total All Import Price: {{number_format($totalAllImportPrice)}}đ <br>--}}
                                {{--                        Total Import Price: {{number_format($totalImportPrice)}}đ <br>--}}
                                Total Order: {{$totalOrder}} <br>
                                Total Order Price: {{number_format($totalFinalPrice)}}đ <br>
                                {{--                        Total Order Final Price: {{number_format($totalFinalPrice)}}đ <br>--}}

{{--                                Total Product: {{$totalProduct}} <br>--}}

                                <table class="table mt-4">
                                    <thead>
                                    <tr>

                                        <th>Type</th>
                                        <th>Sold</th>
                                        <th>Available</th>
                                    </tr>
                                    </thead>
                                    <tbody>
                                    @foreach(\App\Contracts\SaleConfig::PRODUCT_TYPES as $type)
                                        <tr>
                                            <td>{{$type}}</td>
                                            <td>
                                                {{array_first($typeProduct, function ($item) use ($type) { return $item->type == $type && $item->status == 'SOLD';})->number ?? 0}}
                                            </td>
                                            <td>
                                                {{array_first($typeProduct, function ($item) use ($type) { return $item->type == $type && $item->status == 'AVAILABLE';})->number ?? 0}}
                                            </td>
                                        </tr>
                                    @endforeach
                                    <tr>
                                        <td>Total</td>
                                        <td>
                                            {{collect($typeProduct)->filter(function ($item) { return $item->status == 'SOLD';})->sum('number')}}

                                        </td>
                                        <td>{{collect($typeProduct)->filter(function ($item) { return $item->status == 'AVAILABLE';})->sum('number')}}

                                        </td>
                                    </tr>
                                    </tbody>
                                </table>
                            </div>
                            <div class="col">
                                <h2>Previous Month</h2> <br>
                                Total Order: {{$previousTotalOrder}} <br>
                                Total Order Price: {{number_format($previousTotalFinalPrice)}}đ <br>
                                <table class="table mt-4">
                                    <thead>
                                    <tr>

                                        <th>Type</th>
                                        <th>Sold</th>
                                    </tr>
                                    </thead>
                                    <tbody>
                                @foreach(\App\Contracts\SaleConfig::PRODUCT_TYPES as $type)
                                    <tr>
                                        <td>{{$type}}</td>
                                        <td>
                                            {{array_first($previousProducts, function ($item) use ($type) { return $item->type == $type && $item->status == 'SOLD';})->number ?? 0}}
                                        </td>

                                    </tr>
                                @endforeach
                                    </tbody>
                                </table>
                            </div>
                            <div class="col">
                                <h2>Current Month</h2> <br>
                                Total Order: {{$currentTotalOrder}} <br>
                                Total Order Price: {{number_format($currentTotalFinalPrice)}}đ <br>
                                <table class="table mt-4">
                                    <thead>
                                    <tr>

                                        <th>Type</th>
                                        <th>Sold</th>
                                    </tr>
                                    </thead>
                                    <tbody>
                                @foreach(\App\Contracts\SaleConfig::PRODUCT_TYPES as $type)
                                    <tr>
                                        <td>{{$type}}</td>
                                        <td>
                                            {{array_first($currentProducts, function ($item) use ($type) { return $item->type == $type && $item->status == 'SOLD';})->number ?? 0}}
                                        </td>

                                    </tr>
                                @endforeach
                                    </tbody>
                                </table>
                            </div>

                            <div class="col">
                                <h2>Orders by month</h2> <br>

                                <table class="table mt-4">
                                    <thead>
                                    <tr>

                                        <th>Month</th>
                                        <th>Count</th>
                                        <th>Total</th>
                                    </tr>
                                    </thead>
                                    <tbody>
                                    @foreach($orderByMonth as $orderByMonth)
                                        <tr>
                                            <td>{{$orderByMonth->MONTH}}</td>
                                            <td>{{$orderByMonth->number}}</td>
                                            <td>{{number_format($orderByMonth->final_amount)}}đ</td>

                                        </tr>
                                    @endforeach

                                    </tbody>
                                </table>
                            </div>

                        </div>
                    </div>
                </div>


            </div> <!-- container-fluid -->
        </div>
        <!-- End Page-content -->


    </div>

@endsection
