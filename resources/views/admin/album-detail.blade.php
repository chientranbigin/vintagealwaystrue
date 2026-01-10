@extends('admin.master')
@section('content')
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">


                <div class="card">
                    <div class="card-body">
                        <h4 class="card-title">Album Info</h4>
                        <form method="POST" action="{{route('admin.albums.detail', $album->id)}}" enctype="multipart/form-data">
                            @csrf
                            <input type="hidden" name="id" value="{{$album->id}}">

                            <div class="form-group w-100">

                                <label for="exampleFormControlTextarea1">Name</label>
                                <input class="form-control" name="name" value="{{$album->name}}">

                            </div>
                            <div class="form-group w-100">

                                <label for="exampleFormControlTextarea1">Description</label>
                                <input class="form-control" name="type" value="{{$album->type}}">

                            </div>
                            <div class="form-group">
                                <label for="exampleFormControlSelect1">Album type</label>
                                <select class="form-control" id="exampleFormControlSelect1" name="album_type_id">
                                    @foreach($albumTypes as $type)
                                        <option {{$type->id == $album->album_type_id ? 'selected' : ''}} value="{{$type->id}}">{{$type->name}}</option>

                                    @endforeach
                                </select>
                            </div>
                            <div class="form-group w-100">

                                <label for="exampleFormControlTextarea1">Banner</label>
                                @if($album->thumb)
                                    <div class="py-3">
                                        <img src="{{asset($album->thumb)}}"  height="300"/>
                                    </div>
                                @endif

                                <div class="custom-file">
                                    <input type="file" class="custom-file-input" id="customFile" name="banner">
                                    <label class="custom-file-label" for="customFile">Choose file</label>
                                </div>
                            </div>


                            <div class="form-group w-100">

                                <label for="exampleFormControlTextarea1">Album Photos</label>
                                @foreach($album->items as $item)
                                    @if($item->type == \App\Contracts\ElementConfig::TYPE_IMAGE)
                                        <div class="py-3">
                                            <img src="{{asset($item->path)}}"  height="300"/>
                                        </div>
                                    @endif

                                @endforeach


                                <div class="custom-file">
                                    <input type="file" class="custom-file-input" id="customFile" name="items[]" multiple>
                                    <label class="custom-file-label" for="customFile">Choose file</label>
                                </div>
                            </div>
                            <div class="form-group w-100">

                                <label for="exampleFormControlTextarea1">Album video</label>
                                <input class="form-control" name="video" value="{{$video->path ?? ''}}">

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
