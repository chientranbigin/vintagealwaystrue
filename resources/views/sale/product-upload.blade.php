@extends('sale.master')
@section('content')
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">


                <div class="card">
                    <div class="card-body">

                        <form method="POST" action="{{route('sale.product.upload.post')}}"
                              enctype="multipart/form-data">
                            @csrf
                            <div class="row">
                                <div class="col">

                                    <div class="form-group">
                                        <label for="exampleFormControlSelect1">Type</label>
                                        <select class="form-control" id="exampleFormControlSelect1" name="type">
                                            @foreach(App\Contracts\SaleConfig::PRODUCT_TYPES as $type)
                                                <option value="{{$type}}">{{$type}}</option>

                                            @endforeach
                                        </select>
                                    </div>
                                    <div class="form-group w-100">

                                        <label for="exampleFormControlTextarea1">Thumb</label>


                                        <div class="custom-file">
                                            <input multiple type="file" class="custom-file-input" id="customFile" name="banner[]">
                                            <label class="custom-file-label" for="customFile">Choose file</label>
                                        </div>
                                    </div>


                                </div>
                            </div>

                            <button type="submit" class="btn btn-secondary waves-effect waves-light btn-submit">Save
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
