@extends('admin.master')
@section('content')
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">

                <div class="card">
                    <div class="card-body">
                        <h4 class="card-title">Summary</h4>
                        <form method="POST" action="{{route('admin.packages.post', $package)}}" enctype="multipart/form-data">
                            @csrf
                            <input type="hidden" name="id" value="{{$package->id}}">
                            <div class="form-group w-100">

                                <label for="exampleFormControlTextarea1">Name</label>
                                <input class="form-control" name="name" value="{{$package->name}}">

                            </div>
                            <div class="form-group w-100">

                                <label for="exampleFormControlTextarea1">Banner</label>
                                @if($package->thumb)
                                    <div class="py-3">
                                        <img src="{{asset($package->thumb)}}"  height="300"/>
                                    </div>
                                @endif


                                <div class="custom-file">
                                    <input type="file" class="custom-file-input" id="customFile" name="banner">
                                    <label class="custom-file-label" for="customFile">Choose file</label>
                                </div>
                            </div>

                            <button type="submit" class="btn btn-secondary waves-effect waves-light btn-submit" >Save</button>

                        </form>

                    </div>

                    <!-- end card-body-->
                </div>

                <a  href="{{route('admin.packages.create', $package->id)}}" class="btn btn-secondary waves-effect waves-light   my-5"  style="    color: #fff !important;">Create sub package </a>


                @foreach($package->details as $detail)
                    <div class="card">
                        <div class="card-body">
                            <h4 class="card-title">{{$detail->name}}</h4>
                            <form method="POST" action="{{route('admin.packages.detail.post', $detail->id)}}" enctype="multipart/form-data">
                                @csrf
                                <input type="hidden" name="id" value="{{$detail->id}}">
                                <div class="form-group w-100">

                                    <label for="exampleFormControlTextarea1">Name</label>
                                    <input class="form-control" name="name" value="{{$detail->name}}">

                                </div>


                                <div class="form-group w-100">

                                    <label for="exampleFormControlTextarea1">Introduce</label>
                                    <textarea class="form-control" id="exampleFormControlTextarea1" rows="4" name="content" >{{$detail->content}}</textarea>


                                </div>
                                <div class="form-group w-100">

                                    <label for="exampleFormControlTextarea1">Price</label>
                                    <input class="form-control" name="price" value="{{$detail->price}}">

                                </div>
                                <button type="submit" class="btn btn-secondary waves-effect waves-light btn-submit" >Save</button>
                                <a  href="{{route('admin.packages.detail.delete', $detail->id)}}" class="btn btn-info waves-effect waves-light btn-submit "  style="    color: #fff !important;">Delete </a>
                            </form>

                        </div>
                        <!-- end card-body-->
                    </div>
                @endforeach
            </div> <!-- container-fluid -->
        </div>
        <!-- End Page-content -->



    </div>


@endsection
