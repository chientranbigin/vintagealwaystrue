@extends('master')
@section('content')

    <!-- /PAGE TOP -->
    <!-- PAGE CONTENT -->
    <div class="page-holder custom-page-template page-full fullscreen-page clearfix">
        <section class="parallax section-holder section-home41" style="background-image:url({{asset($package->thumb)}});">
            <div class="container parallax-content alignc">
                <div class="row">
                    <div class="col-lg-8 offset-lg-2">
                        <h1 class="display-3 white">{{$package->name}}</h1>
                    </div>
                </div>
                <!-- /row -->
            </div>
            <!-- /container -->
        </section>
        @foreach($package->details as $detail)
            <!-- SECTION 1-->
            <section class="section-holder margin-b72 aligncm">
                <div class="container">
                    <div class="row">
                        <div class="col-lg-12">
                            <h2 class="section-title margin-b32">{{$detail->name}}</h2>
                        </div>

                        <div class="col-lg-12">
                            <p>{{$detail->content}}</p>
                        </div>
                        <div class="col-lg-12">
                            <p><b>{{$detail->price}}</b></p>
                        </div>
                    </div>
                    <!-- /row -->
                </div>
                <!-- /container -->
            </section>
        @endforeach


        <!-- /SECTION 2-->
    </div>
    <!-- page-content -->
@endsection
