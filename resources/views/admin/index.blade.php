@extends('admin.master')
@section('content')
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">

                @foreach(\App\Contracts\ElementConfig::ELEMENTS as $groupName => $group)
                    <div class="page-title-box d-flex align-items-center justify-content-between mt-5">
                        <h4 class="mb-0 font-size-18">{{$groupName}}</h4>
                    </div>
                    @foreach($group as $element)
                        <form method="POST" action="{{route('admin.elements.store')}}" enctype="multipart/form-data">
                            @csrf
                            <input type="hidden" name="name" value="{{$element['name']}}">
                            <input type="hidden" name="type" value="{{$element['type']}}">
                            <input type="hidden" name="group" value="{{$groupName}}">

                            <div class="d-flex w-100 align-items-end">

                                <div class="form-group w-100">

                                    <label for="exampleFormControlTextarea1">{{$element['name']}}</label>
                                    @if($element['type'] == \App\Contracts\ElementConfig::TYPE_IMAGE && array_get($elementData, "{$groupName}.{$element['name']}"))
                                        <div class="py-3">
                                            <img src="{{asset(array_get($elementData, "{$groupName}.{$element['name']}"))}}"  height="300"/>
                                        </div>

                                    @endif

                                    @if($element['type'] == \App\Contracts\ElementConfig::TYPE_TEXT)
                                        <textarea class="form-control" id="exampleFormControlTextarea1" rows="4" name="value" >{{array_get($elementData, "{$groupName}.{$element['name']}")}}</textarea>
                                    @else
                                        <div class="custom-file">
                                            <input type="file" class="custom-file-input" id="customFile" name="value">
                                            <label class="custom-file-label" for="customFile">Choose file</label>
                                        </div>

                                    @endif
                                </div>
                                <button type="submit" class="btn btn-secondary waves-effect waves-light btn-submit ml-3 mb-3" >Save</button>



                            </div>
                        </form>


                    @endforeach
                @endforeach

            </div> <!-- container-fluid -->
        </div>
        <!-- End Page-content -->



    </div>


@endsection
