@extends('admin.master')
@section('content')
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">

                <div class="card">
                    <div class="card-body">
                        <form method="POST" action="{{route('admin.packages.create.post', $id)}}" enctype="multipart/form-data">
                            @csrf
                            <div class="form-group w-100">

                                <label for="exampleFormControlTextarea1">Name</label>
                                <input class="form-control" name="name" >

                            </div>


                            <div class="form-group w-100">

                                <label for="exampleFormControlTextarea1">Introduce</label>
                                <textarea class="form-control" id="exampleFormControlTextarea1" rows="4" name="content" ></textarea>


                            </div>
                            <div class="form-group w-100">

                                <label for="exampleFormControlTextarea1">Price</label>
                                <input class="form-control" name="price" >

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
