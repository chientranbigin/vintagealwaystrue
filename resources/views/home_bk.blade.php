@extends('master')

@section('content')
    <!-- PAGE CONTENT -->
    <div class="page-holder custom-page-template page-full fullscreen-page home-page-content clearfix">
        <!-- SECTION 1-->
        <section class="parallax section-holder section-home41"
                 style="background-image:url('{{asset(array_get($elements, 'HOME_BANNER.Banner'))}}');">
            <div class="container parallax-content alignc">
                <div class="row">
                    <div class="col-lg-8 offset-lg-2">
                        <h1 class="display-3 white">{{array_get($elements, 'HOME_BANNER.Text')}}</h1>
                    </div>
                </div>
                <!-- /row -->
            </div>
            <!-- /container -->
        </section>
        <!-- /SECTION 1-->
        <section class="section-holder">
            <div class="container">
                <div class="short-portfolio-grid-holder-v3">
                    <ul class="portfolio-filter">
                        <li><a class="view-more current" href="#" data-filter="*">Tất cả</a></li>
                        @foreach($albumTypes as $type )
                            <li><a class="view-more" href="#" data-filter=".type{{$type->id}}">{{$type->name}}</a></li>

                        @endforeach
                    </ul>
                    <div class="port-masonry port-masonry-3cols">
                        <div class="row layout-masonry">
                            @foreach($albums as $album)

                            <div class="col-sm-6 col-lg-4 blog-item-masonry port-item-masonry type{{$album->album_type_id}}">
                                <div class="port-item-htext">
                                    <a href="{{route('albumDetail', $album->id)}}">
                                        <div class="port-img"><img src="{{asset($album->thumb)}}" class="img-fluid" alt="{{$album->name}}"></div>
                                        <div class="album-arrow"></div>
                                    </a>
                                </div>
                                <div class="port-caption-masonry">
                                    <h3><a href="{{route('albumDetail', $album->id)}}">{{$album->name}}</a></h3>
                                    <ul class="port-categ-masonry">
                                        <li><a href="#">{{$album->type}}</a></li>
                                    </ul>
                                </div>
                            </div>
                            @endforeach

                        </div>
                    </div>
                    <!-- /port-grid-v3 -->
                </div>
                <!-- /short-portfolio-grid-holder-v3 -->

            </div>
            <!-- /container -->
        </section>


        <section class="section-holder">
            <div class="container">
                <div class="row">
                    @foreach($packages as $package)
                        <div class="col-lg-4 margin-bm54">
                            <div class="price-box" style="    height: 100%;">
                                <h4 class="price-title">{{$package->name}}</h4>
                                <span class="price-divider"></span>
                                <ul class="price-features" style="height: 200px">
                                    @foreach($package->details as $detail)
                                        <li>{{$detail->name}}: {{$detail->price}}</li>

                                    @endforeach
                                </ul>
                                <p><a class="read-more" href="{{route('packageDetail', $package->id)}}">Xem chi tiết</a></p>
                            </div>
                        </div>
                    @endforeach


                </div>
                <!-- /row -->
            </div>
            <!-- /container -->
        </section>


    </div>
    <!-- home-page-content -->
@endsection
