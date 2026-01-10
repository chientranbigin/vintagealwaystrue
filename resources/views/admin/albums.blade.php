@extends('admin.master')
@section('content')
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">


                <div class="card">
                    <div class="card-body">
                        <h4 class="card-title">Album Type</h4>
                        <div class="table-responsive">
                            <table class="table mb-0">
                                <thead>
                                <tr>
                                    <th>Name</th>
                                    <th>Description</th>
                                    <th>Action</th>
                                </tr>
                                </thead>
                                <tbody>
                                @foreach($albumTypes as $type )
                                    <form method="POST" action="{{route('admin.albumType.update')}}" enctype="multipart/form-data">
                                        @csrf
                                        <input type="hidden" name="id" value="{{$type->id}}">
                                    <tr>
                                        <td>

                                            <input class="form-control" name="name" value="{{$type->name}}">
                                        </td>
                                        <td>

                                            <input class="form-control" name="desc" value="{{$type->desc}}">
                                        </td>
                                        <td>
                                            <button type="submit" class="btn btn-secondary waves-effect waves-light btn-submit ml-3 mb-3" >Save</button>
                                        </td>
                                    </tr>
                                    </form>
                                @endforeach


                                </tbody>
                            </table>
                        </div>

                    </div>
                    <!-- end card-body-->
                </div>
                <div class="card">
                    <div class="card-body">
                        <h4 class="card-title">Albums</h4>
                        <a  href="{{route('admin.albums.create')}}" class="btn btn-secondary waves-effect waves-light btn-submit  mb-3"  style="    color: #fff !important;">Create </a>

                        <div class="table-responsive">
                            <table class="table mb-0">
                                <thead>
                                <tr>
                                    <th>#</th>
                                    <th>Name</th>
                                    <th>Type</th>
                                    <th>Created At</th>
                                    <th>Updated At</th>
                                    <th>Action</th>
                                </tr>
                                </thead>
                                <tbody>
                                @foreach($albums as $album)
                                    <tr>
                                        <td>{{$album->id}}</td>
                                        <td>
                                            <a href="{{route('admin.albums.detail', $album->id)}}">{{$album->name}}</a>
                                        </td>
                                        <td>{{$album->albumType->name}}</td>
                                        <td>{{$album->created_at}}</td>
                                        <td>{{$album->updated_at}}</td>
                                        <td><a  href="{{route('admin.albums.detail.delete', $album->id)}}" class="btn btn-secondary waves-effect waves-light btn-submit  "  style="    color: #fff !important;">Delete </a></td>
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
