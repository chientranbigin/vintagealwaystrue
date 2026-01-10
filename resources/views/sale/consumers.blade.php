@extends('sale.master')
@section('content')
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">


                <div class="card">
                    <div class="card-body">

                        <div class="table-responsive">
                            <table class="table mb-0">
                                <thead>
                                <tr>

                                    <th>Name</th>
                                    <th>Phone</th>
                                    <th>Address</th>
                                    <th>Total Orders</th>
                                    <th>Action</th>
                                </tr>
                                </thead>
                                <tbody>
                                @foreach($consumers as $consumer)

                                    <td>
                                        {{$consumer->name}}
                                    </td>
                                    <td>
                                        {{$consumer->phone}}
                                    </td>
                                    <td>
                                        {{$consumer->address}}
                                    </td>
                                    <td>
                                        <span class="right badge badge-light">{{$consumer->orders->count()}}</span>


                                    </td>
                                    <td>
                                        <a href="{{route('sale.order.create', [
    'consumer_id' => $consumer->id
])}}"
                                           class="btn btn-secondary waves-effect waves-light   "
                                           style="    color: #fff !important;">Create Order </a>
                                    </td>
                                    </tr>
                                @endforeach


                                </tbody>
                            </table>
                        </div>

                    </div>
                    <!-- end card-body-->
                </div>

            </div> <!-- container-fluid -->
        </div>
        <!-- End Page-content -->


    </div>

@endsection
