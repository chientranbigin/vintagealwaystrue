<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="utf-8">
    <title>Vintage Always True Shop Online</title>
    <meta content="width=device-width, initial-scale=1.0" name="viewport">
    <meta content="Vintage Always True Shop" name="keywords">
    <meta content="Vintage Always True Shop" name="description">

    <!-- Favicon -->
    <link href="{{asset('customer-theme/img/favicon.ico')}}" rel="icon">

    <!-- Google Web Fonts -->
    <link rel="preconnect" href="https://fonts.gstatic.com">
    <link href="https://fonts.googleapis.com/css2?family=Roboto:wght@400;500;700&display=swap" rel="stylesheet">

    <!-- Font Awesome -->
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.10.0/css/all.min.css" rel="stylesheet">

    <!-- Libraries Stylesheet -->
    <link href="{{asset('customer-theme/lib/animate/animate.min.css')}}" rel="stylesheet">
    <link href="{{asset('customer-theme/lib/owlcarousel/assets/owl.carousel.min.css')}}" rel="stylesheet">

    <!-- Customized Bootstrap Stylesheet -->
    <link href="{{asset('customer-theme/css/style.css')}}" rel="stylesheet">
    <style>
        .btn-submit {
            height: 40px;
            width: 90px;
            align-content: center;
        }
        .nav-item {
            padding: 20px 10px !important;
        }
        .badge {
            line-height: 1.5;
        }
    </style>
    <!-- Paste this right before your closing </head> tag -->
    <script type="text/javascript">
        (function (f, b) { if (!b.__SV) { var e, g, i, h; window.mixpanel = b; b._i = []; b.init = function (e, f, c) { function g(a, d) { var b = d.split("."); 2 == b.length && ((a = a[b[0]]), (d = b[1])); a[d] = function () { a.push([d].concat(Array.prototype.slice.call(arguments, 0))); }; } var a = b; "undefined" !== typeof c ? (a = b[c] = []) : (c = "mixpanel"); a.people = a.people || []; a.toString = function (a) { var d = "mixpanel"; "mixpanel" !== c && (d += "." + c); a || (d += " (stub)"); return d; }; a.people.toString = function () { return a.toString(1) + ".people (stub)"; }; i = "disable time_event track track_pageview track_links track_forms track_with_groups add_group set_group remove_group register register_once alias unregister identify name_tag set_config reset opt_in_tracking opt_out_tracking has_opted_in_tracking has_opted_out_tracking clear_opt_in_out_tracking start_batch_senders people.set people.set_once people.unset people.increment people.append people.union people.track_charge people.clear_charges people.delete_user people.remove".split( " "); for (h = 0; h < i.length; h++) g(a, i[h]); var j = "set set_once union unset remove delete".split(" "); a.get_group = function () { function b(c) { d[c] = function () { call2_args = arguments; call2 = [c].concat(Array.prototype.slice.call(call2_args, 0)); a.push([e, call2]); }; } for ( var d = {}, e = ["get_group"].concat( Array.prototype.slice.call(arguments, 0)), c = 0; c < j.length; c++) b(j[c]); return d; }; b._i.push([e, f, c]); }; b.__SV = 1.2; e = f.createElement("script"); e.type = "text/javascript"; e.async = !0; e.src = "undefined" !== typeof MIXPANEL_CUSTOM_LIB_URL ? MIXPANEL_CUSTOM_LIB_URL : "file:" === f.location.protocol && "//cdn.mxpnl.com/libs/mixpanel-2-latest.min.js".match(/^\/\//) ? "https://cdn.mxpnl.com/libs/mixpanel-2-latest.min.js" : "//cdn.mxpnl.com/libs/mixpanel-2-latest.min.js"; g = f.getElementsByTagName("script")[0]; g.parentNode.insertBefore(e, g); } })(document, window.mixpanel || []);

        mixpanel.init('6ae32e0b040359048fe7b7a9391c9675', {autocapture: true,  track_pageview: true });


    </script>
</head>

<body>
<!-- Topbar Start -->
<div class="container-fluid">

    <div class="row align-items-center bg-light py-3 px-xl-5 d-none d-lg-flex">
        <div class="col-lg-4">
            <a href="{{route('home')}}" class="text-decoration-none">
                <span class="h1 text-uppercase text-primary bg-dark px-2">VINTAGE</span>
                <span class="h1 text-uppercase text-dark bg-primary px-2 ml-n1">ALWAYS TRUE</span>
            </a>
        </div>
        <div class="col-lg-4 col-6 text-left">
            <form action="">
                <div class="input-group">
                    <input type="text" class="form-control" placeholder="Search for products">
                    <div class="input-group-append">
                            <span class="input-group-text bg-transparent text-primary">
                                <i class="fa fa-search"></i>
                            </span>
                    </div>
                </div>
            </form>
        </div>
        <div class="col-lg-4 col-6 text-right">
            <p class="m-0">Contact me</p>
            <h5 class="m-0"><a href="https://www.fb.com/vintage.always.true/">fb.com/vintage.always.true</a></h5>
        </div>
    </div>
</div>
<!-- Topbar End -->
<?php $routeName = request()->route()->getName(); ?>

<div class="container-fluid bg-dark mb-30">
    <div class="row px-xl-5">
        <div class="col-lg-3 d-none d-lg-block">
            <a class="btn d-flex align-items-center justify-content-between bg-primary w-100" data-toggle="collapse" href="#navbar-vertical" style="height: 65px; padding: 0 30px;">
                <h6 class="text-dark m-0"><i class="fa fa-bars mr-2"></i>Sản phẩm</h6>
                <i class="fa fa-angle-down text-dark"></i>
            </a>
            <nav class="collapse position-absolute navbar navbar-vertical navbar-light align-items-start p-0 bg-light" id="navbar-vertical" style="width: calc(100% - 30px); z-index: 999;">
                <div class="navbar-nav w-100">
                   @foreach(\App\Contracts\SaleConfig::PRODUCT_TYPES as $type)
                        <a href="{{route('customer.products', compact('type'))}}" class="nav-item nav-link">{{ucfirst(strtolower($type))}}</a>
                   @endforeach
                </div>
            </nav>
        </div>
        <div class="col-lg-9">
            <nav class="navbar navbar-expand-lg bg-dark navbar-dark py-3 py-lg-0 px-0">
                <a href="" class="text-decoration-none d-block d-lg-none">
                    <span class="h1 text-uppercase text-dark bg-light px-2">VINTAGE</span>
                    <span class="h1 text-uppercase text-light bg-primary px-2 ml-n1">ALWAYS TRUE</span>
                </a>
{{--                <button type="button" class="navbar-toggler" data-toggle="collapse" data-target="#navbarCollapse">--}}
{{--                    <span class="navbar-toggler-icon"></span>--}}
{{--                </button>--}}
                <div class=" navbar-collapse justify-content-between" id="navbarCollapse" style="        display: flex !important
;
        flex-basis: auto;">
                    <div class="navbar-nav mr-auto py-0" style="flex-direction: row !important;" >
                        <a href="{{url('/')}}" class="nav-item nav-link {{ in_array($routeName, ['home']) ? 'active' : '' }}">Trang chủ</a>
                        <a href="{{route('customer.products')}}" class="nav-item nav-link {{ in_array($routeName, ['customer.products']) ? 'active' : '' }}">Sản phẩm</a>
{{--                        <a href="contact.html" class="nav-item nav-link">Liên hệ</a>--}}
                    </div>
{{--                    <div class="navbar-nav ml-auto py-0 d-none d-lg-block">--}}

{{--                        <a href="" class="btn px-0 ml-3">--}}
{{--                            <i class="fas fa-shopping-cart text-primary"></i>--}}
{{--                            <span class="badge text-secondary border border-secondary rounded-circle" style="padding-bottom: 2px;">0</span>--}}
{{--                        </a>--}}
{{--                    </div>--}}
                </div>
            </nav>
        </div>
    </div>
</div>
@yield('content')


<!-- Back to Top -->
<a href="#" class="btn btn-primary back-to-top"><i class="fa fa-angle-double-up"></i></a>


<!-- JavaScript Libraries -->
<script src="https://code.jquery.com/jquery-3.4.1.min.js"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/js/bootstrap.bundle.min.js"></script>
<script src="{{asset('customer-theme/lib/easing/easing.min.js')}}"></script>
<script src="{{asset('customer-theme/lib/owlcarousel/owl.carousel.min.js')}}"></script>

<!-- Contact Javascript File -->
<script src="{{asset('customer-theme/mail/jqBootstrapValidation.min.js')}}"></script>
<script src="{{asset('customer-theme/mail/contact.js')}}"></script>

<!-- Template Javascript -->
<script src="{{asset('customer-theme/js/main.js')}}"></script>
</body>

</html>
