@extends('master')
@section('content')
    <!-- PAGE TOP -->
    <div class="top-single-bkg  topnoimg ">
        <div class="inner-desc">
            <div class="container">
                <h1 class="post-title single-post-title">Giới thiệu</h1>
            </div>
        </div>
    </div>
    <!-- /PAGE TOP -->
    <!-- PAGE CONTENT -->
    <div class="page-holder custom-page-template page-full fullscreen-page home-page-content clearfix">
        <!-- SECTION 1-->
        <section class="section-holder section-contact41">
            <div class="container">
                <div class="row">
                    <div class="col-md-4">
                        <div class="box-holder alignc margin-bm54">
                            <div class="box-icon2">
                                <div class="icon-box-3-inner">
                                    <i class="fas fa-map-marker-alt"></i>
                                </div>
                            </div>
                            <h5 class="white">Địa chỉ</h5>
                            <p class="white">{{array_get($elements, 'ABOUT_ME.Address')}}</p>
                        </div>
                    </div>
                    <!-- /col-md-4 -->
                    <div class="col-md-4">
                        <div class="box-holder alignc margin-bm54">
                            <div class="box-icon2">
                                <div class="icon-box-3-inner">
                                    <i class="fas fa-phone-alt"></i>
                                </div>
                            </div>
                            <h5 class="white">Phone </h5>
                            <p class="white">{{array_get($elements, 'ABOUT_ME.Phone')}}</p>
                        </div>
                    </div>
                    <!-- /col-md-4 -->
                    <div class="col-md-4">
                        <div class="box-holder alignc">
                            <div class="box-icon2">
                                <div class="icon-box-3-inner">
                                    <i class="far fa-envelope"></i>
                                </div>
                            </div>
                            <h5 class="white">Email</h5>
                            <p class="white">{{array_get($elements, 'ABOUT_ME.Email')}}</p>
                        </div>
                    </div>
                    <!-- /col-md-4 -->
                </div>
                <!-- /row -->
            </div>
            <!-- /container -->
        </section>
        <section class="section-holder">
            <div class="container">
                <div class="row">
                    <div class="col-lg-6 margin-bm32 aligncm">
                        <img class="img-fluid" src="{{asset(array_get($elements, 'ABOUT_ME.Avatar'))}}" alt="">
                    </div>
                    <div class="col-lg-6 aligncm">
                        <h2 class="section-title margin-b32">{{array_get($elements, 'ABOUT_ME.About page title')}}</h2>
                        <p>{{array_get($elements, 'ABOUT_ME.About page content')}}</p>
                        <ul class="social-media margin-t32">
                            <li><a class="social-facebook" href="{{array_get($elements, 'ABOUT_ME.Facebook')}}" target="_blank"><i class="fab fa-facebook-f"></i></a></li>
                            <li><a class="social-instagram" href="{{array_get($elements, 'ABOUT_ME.Instagram')}}" target="_blank"><i class="fab fa-instagram"></i></a></li>
                            <li><a class="social-instagram" href="{{array_get($elements, 'ABOUT_ME.Tiktok')}}" target="_blank"><i class="fab fa-brands fa-tiktok"></i></a></li>

                        </ul>
                    </div>
                </div>
                <!-- /row -->
            </div>
            <!-- /container -->
        </section>
        <!-- /SECTION 1-->

    </div>
    <!-- home-page-content -->
    <!-- /PAGE CONTENT -->
@endsection
