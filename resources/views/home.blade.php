@extends('master')

@section('content')
    <style>
        footer {
            margin-top: 0px;
        }
    </style>
    <div class="slider-container">
        <div class="owl-carousel owl-theme home-slider">
            <div class="slider-post slider-item-box-bkg">
                <div class="slider-img" style="background-image:url({{asset(array_get($elements, 'HOME_BANNER.Banner'))}});"></div>
                <div class="container slider-caption">
                    <div class="slider-text">
                        <h1>{{array_get($elements, 'HOME_BANNER.Text')}}</h1>
                        <p>{{array_get($elements, 'HOME_BANNER.Sub Text')}}</p>

                    </div>
                </div>
                <!--slider-caption-->
            </div>

            <div class="slider-post slider-item-box-bkg">
                <div class="slider-img" style="background-image:url({{asset(array_get($elements, 'HOME_BANNER.Banner 2'))}});"></div>
                <div class="container slider-caption">
                    <div class="slider-text">
                        <h1>{{array_get($elements, 'HOME_BANNER.Text')}}</h1>
                    </div>
                </div>
                <!--slider-caption-->
            </div>

            <div class="slider-post slider-item-box-bkg">
                <div class="slider-img" style="background-image:url({{asset(array_get($elements, 'HOME_BANNER.Banner 3'))}});"></div>
                <div class="container slider-caption">
                    <div class="slider-text">
                        <h1>{{array_get($elements, 'HOME_BANNER.Text')}}</h1>
                    </div>
                </div>
                <!--slider-caption-->
            </div>

            <div class="slider-post slider-item-box-bkg">
                <div class="slider-img" style="background-image:url({{asset(array_get($elements, 'HOME_BANNER.Banner 4'))}});"></div>
                <div class="container slider-caption">
                    <div class="slider-text">
                        <h1>{{array_get($elements, 'HOME_BANNER.Text')}}</h1>
                    </div>
                </div>
                <!--slider-caption-->
            </div>

            <!--slider-post-->
        </div>
    </div>
    <!-- /HOME SLIDER -->
    <!-- PAGE CONTENT -->
    <div class="page-holder custom-page-template page-full fullscreen-page home-page-content clearfix">
        <!-- SECTION 1-->
        <section class="section-holder section-home-13-1">
            <div class="container">
                <div class="row">
                    @foreach($albumTypes as $type )

                    <div class="col-lg-4">
                        <div class="icon-box-3-wrapper">
                            <div class="icon-box-3-icon icon-size-medium icon-color">
                                <i class="fas fa-camera-retro"></i>
                            </div>
                            <div class="icon-box-3-content white">
                                <h3 class="box-title box-title-3 white margin-b16">{{$type->name}}</h3>
                                <p>{{$type->desc}}</p>
                            </div>
                        </div>
                    </div>
                    @endforeach

                    <!-- /col-lg-4 -->
                    <!-- /col-lg-4 -->
                </div>
                <!-- /row -->
            </div>
            <!-- /container -->
        </section>
        <!-- /SECTION 1-->
        <!-- SECTION 2 -->
        <section class="section-holder section-home-13-2 mb-0">
            <div class="short-portfolio-grid-holder-v4">
                <div class="port-grid-v4 grid-v4-masonry">
                    <span class="grid-item-v4-init"></span>
                    @foreach($albums as $album)
                        <div class="grid-item-v4 grid-item-v4-style2">
                            <div class="grid-item-img-v4">
                                <a href="{{route('albumDetail', $album->id)}}">
                                    <div class="item-img-v4" style="background-image:url('{{asset($album->thumb)}}');"></div>
                                    <div class="port-caption-v4">
                                        <h2>{!! $album->name !!}</h2>
                                        <ul class="portfolio-categ">
                                            <li>{{$album->type}}</li>
                                        </ul>
                                        <div class="grid-img-icon"><i class="fas fa-arrow-right"></i></div>
                                    </div>
                                </a>
                            </div>
                        </div>
                    @endforeach

                    <!-- /grid-item-v4 -->
                </div>
                <!-- /port-grid-v2 -->
            </div>
            <!-- /short-portfolio-grid-holder-v2 -->
        </section>
        <!-- /SECTION 2 -->

        <!-- /SECTION 3 -->
        <!-- SECTION 4 -->
        <!-- /SECTION 4 -->
        <!-- SECTION 5 -->
        <!-- /SECTION 5 -->
        <!-- SECTION 6-->

        <!-- /SECTION 6-->
    </div>
@endsection
