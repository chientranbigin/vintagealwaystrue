@extends('admin.master')
@section('content')
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">

                    <div class="card">
                        <div class="card-body">
                            <h4 class="card-title">Package Types</h4>

                            <div class="table-responsive">
                                <table class="table mb-0">
                                    <thead>
                                    <tr>
                                        <th>#</th>
                                        <th>Name</th>
                                        <th>Created At</th>
                                        <th>Updated At</th>
                                    </tr>
                                    </thead>
                                    <tbody>
                                    @foreach($packages as $package)
                                        <tr>
                                            <td>{{$package->id}}</td>
                                            <td>
                                                <a href="{{route('admin.packages.detail', $package->id)}}">{{$package->name}}</a>
                                            </td>
                                            <td>{{$package->created_at}}</td>
                                            <td>{{$package->updated_at}}</td>
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
