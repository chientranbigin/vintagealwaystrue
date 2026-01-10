<!DOCTYPE html>
<html lang="en-US">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TadiMedia</title>
    <!-- Google Fonts -->
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Cormorant:ital,wght@0,400;0,500;0,600;0,700;1,400;1,700&family=Open+Sans:ital,wght@0,400;0,600;0,700;1,400;1,700&display=swap" rel="stylesheet">
    <!-- Bootstrap CSS -->
    <link rel='stylesheet' id='bootstrap-css'  href='{{asset('app/css/bootstrap/css/bootstrap.css')}}' type='text/css' media='all' />
    <!-- Font Awesome Icons CSS -->
    <link rel='stylesheet' id='font-awesome'  href='{{asset('app/css/fontawesome/css/all.css')}}' type='text/css' media='all' />
    <!-- Owl Carousel -->
    <link rel='stylesheet' id='owl-carousel'  href='{{asset('app/js/owl-carousel/owl.carousel.css')}}' type='text/css' media='all' />
    <!-- Main CSS File -->
    <link rel='stylesheet' id='gleam-style-css'  href='{{asset('app/style.css')}}' type='text/css' media='all' />
    <!-- favicons -->
    <link rel="icon" href="{{asset(array_get($elements, 'LOGO.Icon tab 32'))}}" sizes="32x32" />
    <link rel="icon" href="{{asset(array_get($elements, 'LOGO.Icon tab 192'))}}" sizes="192x192" />
    <link rel="apple-touch-icon-precomposed" href="{{asset('app/images/icons/favicon-180x180.png')}}" />

    <link rel='stylesheet' id='gleam-style-custom-css'  href='{{asset('app/custom.css')}}' type='text/css' media='all' />

</head>
<?php
$routeName = request()->route()->getName();

?>
<body class="home {{ in_array($routeName, ['about']) ? 'no-top-image': '' }}">
<div class="menu-mask"></div>
<!-- MOBILE MENU HOLDER -->
<div class="mobile-menu-holder">
    <div class="modal-menu-container">
        <div class="exit-mobile"> <span class="icon-bar1"></span> <span class="icon-bar2"></span></div>
        <ul class="menu-mobile">
            <li class="menu-item {{in_array($routeName, ['home', 'packageDetail']) ? 'current-menu-item' : ''}}">
                <a href="{{route('home')}}">Trang chủ</a>
            </li>
            <li class="menu-item {{in_array($routeName, ['album', 'albumDetail'])  ? 'current-menu-item' : ''}}">
                <a href="{{route('album')}}">Albums</a>
            </li>


            <li class="menu-item {{in_array($routeName, ['film']) ? 'current-menu-item' : ''}} ">
                <a href="{{route('film')}}">Films</a>
            </li>
            <li class="menu-item {{in_array($routeName, ['about']) ? 'current-menu-item' : ''}} ">
                <a href="{{route('about')}}">Giới thiệu</a>
            </li>

        </ul>
    </div>
    <!-- /modal-menu-container -->
    <div class="menu-contact">
        <ul>
            <li><i class="fas fa-map-marker-alt"></i><span>{{array_get($elements, 'ABOUT_ME.Address')}}</span></li>
            <li><i class="fas fa-mobile-alt"></i><span>{{array_get($elements, 'ABOUT_ME.Phone')}}</span></li>
            <li><i class="far fa-envelope"></i><span>{{array_get($elements, 'ABOUT_ME.Email')}}</span></li>
        </ul>
        <ul class="social-media">
            <li><a class="social-facebook" href="{{array_get($elements, 'ABOUT_ME.Facebook')}}" target="_blank"><i class="fab fa-facebook-f"></i></a></li>
            <li><a class="social-instagram" href="{{array_get($elements, 'ABOUT_ME.Instagram')}}" target="_blank"><i class="fab fa-instagram"></i></a></li>
            <li><a class="social-instagram" href="{{array_get($elements, 'ABOUT_ME.Tiktok')}}" target="_blank"><i class="fab fa-brands fa-tiktok"></i></a></li>

        </ul>
    </div>
</div>
<!-- /MOBILE MENU HOLDER -->
<!-- HEADER -->
<header class="main-header header-1">
    <div class="container">
        <div class="header-container">
            <div class="logo logo-1 logo-white"><a href="{{route('home')}}"><img class="img-fluid"  src="{{ asset(array_get($elements, 'LOGO.Logo'))}}" alt="Gleam"></a></div>
            <div class="logo logo-1 logo-dark"><a href="{{route('home')}}"><img class="img-fluid" src="{{ asset(array_get($elements, 'LOGO.Logo'))}}" alt="Gleam"></a></div>
            <nav class="nav-holder nav-holder-1">
                <ul class="menu-nav">
                    <li class="menu-item {{in_array($routeName, ['home', 'packageDetail']) ? 'current-menu-item' : ''}}">
                        <a href="{{route('home')}}">Trang chủ</a>
                    </li>
                    <li class="menu-item {{in_array($routeName, ['album', 'albumDetail'])  ? 'current-menu-item' : ''}}">
                        <a href="{{route('album')}}">Albums</a>
                    </li>
                    <li class="menu-item {{in_array($routeName, ['film']) ? 'current-menu-item' : ''}} ">
                        <a href="{{route('film')}}">Films</a>
                    </li>
                    <li class="menu-item {{in_array($routeName, ['about']) ? 'current-menu-item' : ''}} ">
                        <a href="{{route('about')}}">Giới thiệu</a>
                    </li>
                </ul>
            </nav>
            <div class="nav-button-holder"> <button type="button" class="nav-button"> <span class="icon-bar"></span> </button></div>
            <ul class="social-media header-social-1">
                <li><a class="social-facebook" href="{{array_get($elements, 'ABOUT_ME.Facebook')}}" target="_blank"><i class="fab fa-facebook-f"></i></a></li>
                <li><a class="social-instagram" href="{{array_get($elements, 'ABOUT_ME.Instagram')}}" target="_blank"><i class="fab fa-instagram"></i></a></li>
                <li><a class="social-instagram" href="{{array_get($elements, 'ABOUT_ME.Tiktok')}}" target="_blank"><i class="fab fa-brands fa-tiktok"></i></a></li>

            </ul>
        </div>
    </div>
</header>

<!-- home-page-content -->
@yield('content')
<footer>
    <div class="container">
        <div class="footer-widgets">
            <div class="row">
                <div class="col-md-4">
                    <div class="foo-block">
                        <div class="widget widget-footer widget_text">
                            <h5 class="widgettitle">Giới thiệu</h5>
                            <div class="textwidget">
                                <p>{{array_get($elements, 'ABOUT_ME.Intro Footer')}}</p>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="col-md-4">
                    <div class="foo-block">
                        <div class="widget widget-footer widget_text">
                            <h5 class="widgettitle">Liên hệ</h5>
                            <div class="textwidget">
                                <div class="textwidget">
                                    <p>{{array_get($elements, 'ABOUT_ME.Address')}}</p>
                                    <p>{{array_get($elements, 'ABOUT_ME.Phone')}}<br> {{array_get($elements, 'ABOUT_ME.Email')}}</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="col-md-4">
                    <div class="foo-block foo-last">
                        <div class="widget_text widget widget-footer widget_custom_html">
                            <h5 class="widgettitle">Menu</h5>
                            <div class="textwidget custom-html-widget">
                                <ul>
                                    <li>  <a href="{{route('home')}}">Trang chủ</a></li>
                                    <li> <a href="{{route('album')}}">Albums</a></li>
                                    <li> <a href="{{route('film')}}">Films</a></li>
                                    <li> <a href="{{route('about')}}">Giới thiệu</a></li>

                                </ul>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="copyright">
            <div class="footer-copy"> Copyright &copy; 2023, TadiMedia</div>
        </div>
    </div>
</footer>
<!-- FOOTER -->
<div class="scrollup">
    <a class="scrolltop" href="#">
        <i class="fa fa-chevron-up"></i>
    </a>
</div>

<!-- JS -->
<script src='{{ asset('app/js/jquery.js') }}'></script>
<script src='{{ asset('app/js/jquery-migrate.min.js')}}'></script>
<script src='{{ asset('app/css/bootstrap/js/popper.js')}}'></script>
<script src='{{ asset('app/css/bootstrap/js/bootstrap.js')}}'></script>
<script src='{{ asset('app/js/easing.js')}}'></script>
<script src='{{ asset('app/js/fitvids.js')}}'></script>
<script src='{{ asset('app/js/owl-carousel/owl.carousel.js')}}'></script>
<script src='{{ asset('app/js/isotope.js')}}'></script>
<script src='{{ asset('app/js/simple-lightbox.js')}}'></script>
<!-- MAIN JS -->
<script src='{{ asset('app/js/init.js')}}'></script>
<!-- CONTAT FORM JS -->
<script src='{{ asset('app/js/jquery.form.min.js')}}'></script>
<script src='{{ asset('app/js/contactform.js')}}'></script>
<script>
    // $( document ).ready(function() {
    //     console.log('ffff')
    //     $('.video-player').addEventListener('click', function(){
    //         console.log('hhi')
    //         // // Play video and go fullscreen
    //         // player.playVideo();
    //         // var playerElement = $(playerWrapperClass);
    //         // var requestFullScreen = playerElement.requestFullScreen || playerElement.mozRequestFullScreen || playerElement.webkitRequestFullScreen;
    //         // if (requestFullScreen) {
    //         //     requestFullScreen.bind(playerElement)();
    //         // }
    //     })
    // });

</script>
</body>
</html>
