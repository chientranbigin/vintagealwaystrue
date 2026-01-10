<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="utf-8"/>
    <title>Admin Tadimedia</title>
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta content="Premium Multipurpose Admin & Dashboard Template" name="description"/>
    <meta content="MyraStudio" name="author"/>
    <meta http-equiv="X-UA-Compatible" content="IE=edge"/>

    <!-- App favicon -->
    <link rel="shortcut icon" href="{{ asset('admin-theme/assets/images/favicon.ico')}}">

    <!-- App css -->
    <link href="{{ asset('admin-theme/assets/css/bootstrap.min.css')}}" rel="stylesheet" type="text/css"/>
    <link href="{{ asset('admin-theme/assets/css/icons.min.css')}}" rel="stylesheet" type="text/css"/>
    <link href="{{ asset('admin-theme/assets/css/theme.min.css')}}" rel="stylesheet" type="text/css"/>

    <style>
        .btn-submit {
            height: 40px;
            width: 90px
        }

        /*.main-content {*/
        /*    max-width: 800px;*/
        /*}*/
    </style>
</head>

<body>
<?php $routeName = request()->route()->getName(); ?>

        <!-- Begin page -->
<div id="layout-wrapper">


    <!-- ========== Left Sidebar Start ========== -->
    <div class="vertical-menu">

        <div data-simplebar class="h-100">

            <div class="navbar-brand-box">
                <a href="/" class="logo">
                    <span>
                            TADIMEDIA
                        </span>
                </a>
            </div>

            <!--- Sidemenu -->
            <div id="sidebar-menu">
                <!-- Left Menu Start -->
                <ul class="metismenu list-unstyled" id="side-menu">
                    <li class="{{in_array($routeName, ['admin.home']) ? 'mm-active' : ''}}">
                        <a href="{{route('admin.home')}}"
                           class="waves-effect {{in_array($routeName, ['admin.home']) ? 'active' : ''}}"><span>Elements</span></a>
                    </li>
                    <li class="{{str_contains($routeName, 'albums') ? 'mm-active' : ''}}">
                        <a href="{{route('admin.albums')}}"
                           class="waves-effect {{str_contains($routeName, 'albums') ? 'active' : ''}}"><span>Albums</span></a>
                    </li>
                    <li class="{{str_contains($routeName, 'packages') ? 'mm-active' : ''}}">
                        <a href="{{route('admin.packages')}}"
                           class="waves-effect {{str_contains($routeName, 'packages') ? 'active' : ''}}"><span>Packages</span></a>
                    </li>


                </ul>
            </div>
            <!-- Sidebar -->
        </div>
    </div>
    <!-- Left Sidebar End -->

    <!-- ============================================================== -->
    <!-- Start right Content here -->
    <!-- ============================================================== -->
    @yield('content')
    <!-- end main content-->

</div>
<!-- END layout-wrapper -->

<!-- Overlay-->
<div class="menu-overlay"></div>


<!-- jQuery  -->
<script src="{{ asset('admin-theme/assets/js/jquery.min.js')}}"></script>
<script src="{{ asset('admin-teme/assets/js/bootstrap.bundle.min.js')}}"></script>
<script src="{{ asset('admin-theme/assets/js/metismenu.min.js')}}"></script>
<script src="{{ asset('admin-theme/assets/js/waves.js')}}"></script>
<script src="{{ asset('admin-theme/assets/js/simplebar.min.js')}}"></script>

<!-- Morris Js-->
<script src="{{ asset('admin-theme/plugins/morris-js/morris.min.js')}}'"></script>
<!-- Raphael Js-->
<script src="{{ asset('admin-theme/plugins/raphael/raphael.min.js')}}'"></script>

<!-- Morris Custom Js-->
<script src="{{ asset('admin-theme/assets/pages/dashboard-demo.js')}}"></script>

<!-- App js -->
<script src="{{ asset('admin-theme/assets/js/theme.js')}}"></script>

</body>

</html>
