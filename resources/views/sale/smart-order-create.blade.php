@extends('sale.master')
@section('content')
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">


                <div class="card" style="max-width: 500px">
                    <div class="card-body">

                        <form method="POST" action="{{ route('sale.order.smart-create.post') }}"
                              enctype="multipart/form-data">
                            @csrf

                            <input type="hidden" name="backRoute" value="{{url()->previous()}}">
                            <input type="hidden" name="consumer_id" value="{{request()->consumer_id}}">

                            <div class="form-group w-100">

                                <label for="exampleFormControlTextarea1">Description</label>
                                <textarea class="form-control" type="text" name="description" rows="5"> </textarea>

                            </div>
                            <button type="submit"
                                    class="btn btn-secondary waves-effect waves-light btn-submit mb-5">Save
                            </button>
                        </form>

                    </div>
                    <!-- end card-body-->
                </div>


            </div> <!-- container-fluid -->
        </div>
        <!-- End Page-content -->


    </div>
@endsection
