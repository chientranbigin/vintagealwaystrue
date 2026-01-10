@extends('master')
@section('content')
    <!-- PAGE TOP -->
    <div class="top-single-bkg " style="background-image:url({{asset($album->thumb)}});">
        <div class="inner-desc">
            <div class="container">
                <h1 class="post-title single-post-title">{{$album->name}}</h1>
                <h3 class="post-title single-post-title" style="font-size: 2.5rem">{{$album->type}}</h3>
            </div>
        </div>
    </div>
    <!-- /PAGE TOP -->
    <!-- PAGE CONTENT -->
    <div class="page-holder custom-page-template page-full fullscreen-page clearfix">

        <!-- SECTION 2-->
        <section class="section-holder">
            <div class="gallery-container">
                <div class="container fs-gallery">
                    <div class="row gallery-holder gallery-grid3cols">
                        @foreach($album->items as $item)
                            <div class="col-sm-6 col-lg-4 gallery-post">
                                @if($item->type == \App\Contracts\ElementConfig::TYPE_IMAGE)
                                    <a href="{{ asset($item->path)}}"
                                       class="lightbox" title="">
                                        <img class="img-fluid"
                                             src="{{ asset($item->path)}}"
                                             alt="" style="height: auto !important;">
                                    </a>
                                @else
                                    <iframe width="800" height="800"
                                            src="{{$item->path}}">
                                    </iframe>
                                @endif

                            </div>
                        @endforeach

                    </div>
                </div>
                <!-- /container -->
            </div>
        </section>
        <!-- /SECTION 2-->
    </div>
    <!-- page-content -->
@endsection
