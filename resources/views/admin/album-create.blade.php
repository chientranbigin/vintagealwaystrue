@extends('admin.master')
@section('content')
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">


                <div class="card">
                    <div class="card-body">
                        <h4 class="card-title">Album Info</h4>
                        <form method="POST" action="{{route('admin.albums.create.post')}}" enctype="multipart/form-data">
                            @csrf

                            <div class="form-group w-100">

                                <label for="exampleFormControlTextarea1">Name</label>
                                <input class="form-control" name="name" >

                            </div>
                            <div class="form-group">
                                <label for="exampleFormControlSelect1">Album type</label>
                                <select class="form-control" id="exampleFormControlSelect1" name="album_type_id">
                                    @foreach($albumTypes as $type)
                                        <option  value="{{$type->id}}">{{$type->name}}</option>

                                    @endforeach
                                </select>
                            </div>
                            <div class="form-group w-100">

                                <label for="exampleFormControlTextarea1">Banner</label>


                                <div class="custom-file">
                                    <input type="file" class="custom-file-input" id="customFile" name="banner">
                                    <label class="custom-file-label" for="customFile">Choose file</label>
                                </div>
                            </div>


                            <div class="form-group w-100">

                                <label for="exampleFormControlTextarea1">Album Photos</label>


                                <div class="custom-file">
                                    <input type="file" class="custom-file-input" id="customFile" name="items[]" multiple>
                                    <label class="custom-file-label" for="customFile">Choose file</label>
                                </div>
                            </div>
                            <div class="form-group w-100">

                                <label for="exampleFormControlTextarea1">Album video</label>
                                <input class="form-control" name="video" >

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
