<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="utf-8"/>
    <title>Admin Vintage Always True</title>
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta content="Premium Multipurpose Admin & Dashboard Template" name="description"/>
    <meta content="MyraStudio" name="author"/>
    <meta http-equiv="X-UA-Compatible" content="IE=edge"/>

    <!-- App favicon -->
    <link rel="shortcut icon" href="{{ asset('admin-theme/assets/images/favicon.ico') }}">

    <!-- App css -->
    <link href="{{ asset('admin-theme/assets/css/bootstrap.min.css') }}" rel="stylesheet" type="text/css"/>
    <link href="{{ asset('admin-theme/assets/css/icons.min.css') }}" rel="stylesheet" type="text/css"/>
    <link href="{{ asset('admin-theme/assets/css/theme.min.css') }}" rel="stylesheet" type="text/css"/>
    
    <!-- Modern Font -->
    <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;500;600;700&display=swap" rel="stylesheet">

    <style>
        body {
            font-family: 'Outfit', sans-serif !important;
            background-color: #f4f7f6;
        }
        h1, h2, h3, h4, h5, h6, .card-title, .nav-link, button {
            font-family: 'Outfit', sans-serif !important;
        }
        .btn-submit {
            height: 40px;
            min-width: 90px;
            align-content: center;
        }

        /*.main-content {*/
        /*    max-width: 800px;*/
        /*}*/
        .page-content {
            padding-top: 20px
        }

        .main-content {
            margin-left: 0px
        }
        .product-checkbox {
            width: 15px;
            height: 15px;
            transform: scale(1.5);
            margin: 8px;
            accent-color: #0d6efd; /* Đổi màu nếu cần (blue Bootstrap) */
        }
        /* ===== Base ===== */
        .mobile-nav,
        .mobile-nav-more {
            display: flex;
            gap: 10px;
        }

        .navbar-nav .nav-link {
            padding: 10px 14px;
            font-weight: 500;
        }

        .navbar-nav {
            flex-direction: row;              /* ⬅️ NGANG */

        }

        /* ===== Mobile ===== */
        @media (max-width: 768px) {
            .mobile-nav {
                display: flex;
                flex-wrap: nowrap;
                overflow-x: auto;
                -webkit-overflow-scrolling: touch;
                gap: 8px;
                flex-direction:row
            }

            .mobile-nav::-webkit-scrollbar {
                display: none;
            }

            .mobile-nav-more {
                display: none;
                flex-direction: row;              /* ⬅️ NGANG */
                flex-wrap: nowrap;
                gap: 8px;
                margin-top: 8px;
                padding: 8px;
                background: #fff;
                border-radius: 10px;
                box-shadow: 0 6px 16px rgba(0,0,0,.12);
                overflow-x: auto;                 /* ⬅️ SCROLL NGANG */
                -webkit-overflow-scrolling: touch;
            }

            .mobile-nav-more.show {
                display: flex;
            }

            .mobile-nav-more::-webkit-scrollbar {
                display: none;
            }

            .mobile-nav-more .nav-item {
                flex: 0 0 auto;                   /* ⬅️ không wrap */
                white-space: nowrap;
            }

            .nav-more {
                margin-left: auto;
            }
            .desktop-nav {
                display: none;
            }
        }


        /* ===== Desktop ===== */
        @media (min-width: 769px) {
            .nav-more,
            .mobile-nav-more {
                display: none;
            }

            .mobile-nav {
                justify-content: center;
            }
            .desktop-nav {
                display: flex;
            }
        }
        /* ===== Navbar base polish ===== */
        /* ===== Layout adjustments for fixed nav ===== */
        body {
            padding-bottom: 80px; /* Space for mobile bottom nav */
        }
        @media (min-width: 769px) {
            body {
                padding-top: 80px; /* Space for desktop top nav */
                padding-bottom: 0;
            }
        }

        /* ===== Header Light Design (Glassmorphism) ===== */
        .main-header {
            background: rgba(255, 255, 255, 0.9) !important;
            backdrop-filter: blur(10px);
            -webkit-backdrop-filter: blur(10px);
            border-bottom: 1px solid rgba(0,0,0,0.05) !important;
            position: fixed;
            width: 100%;
            left: 0;
            z-index: 1050;
            box-shadow: 0 2px 15px rgba(0,0,0,0.03);
        }

        /* Desktop top, Mobile bottom */
        @media (min-width: 769px) {
            .main-header {
                top: 0;
            }
        }
        @media (max-width: 768px) {
            .main-header {
                bottom: 0;
                box-shadow: 0 -2px 15px rgba(0,0,0,0.05);
                border-top: 1px solid rgba(0,0,0,0.05);
                border-bottom: none !important;
            }
        }

        /* Menu text - Dark for light theme */
        .main-header .nav-link {
            color: #555 !important;
            opacity: 0.7;
            font-size: 11px;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            padding: 10px 4px !important;
            transition: all 0.3s ease;
        }

        /* Icon Colors */
        .nav-icon {
            font-size: 1.5rem;
            margin-bottom: 2px;
            color: #666;
            transition: color 0.3s ease;
        }

        /* Active State - Premium primary color */
        .main-header .nav-item.active .nav-link {
            opacity: 1;
            font-weight: 700;
            color: #0d6efd !important;
        }
        .main-header .nav-item.active .nav-icon {
            color: #0d6efd !important;
        }

        /* Mobile grid */
        @media (max-width: 768px) {
            .mobile-nav {
                display: grid !important;
                grid-template-columns: repeat(4, 1fr);
                width: 100%;
                margin: 0;
                padding: 0;
            }
            .mobile-nav .nav-item { border: none; }
            
            /* More menu slide up */
            .mobile-nav-more {
                position: fixed;
                bottom: 85px;
                left: 15px;
                right: 15px;
                background: #fff !important;
                border-radius: 20px;
                box-shadow: 0 10px 30px rgba(0,0,0,0.1);
                padding: 15px !important;
                flex-direction: column !important;
                gap: 5px !important;
                display: none;
                z-index: 1060;
                border: 1px solid rgba(0,0,0,0.05);
            }
            .mobile-nav-more.show { display: flex !important; }
            .mobile-nav-more .nav-link {
                color: #444 !important;
                flex-direction: row !important;
                justify-content: flex-start !important;
                font-size: 15px !important;
                padding: 14px 20px !important;
                border-radius: 12px;
                opacity: 1;
            }
            .mobile-nav-more .nav-link i {
                font-size: 1.2rem;
                width: 30px;
                color: #0d6efd;
            }
            .mobile-nav-more .nav-link:hover { background: #f0f7ff; color: #0d6efd !important; }
        }

        @media (min-width: 769px) {
            .main-header .nav-link {
                flex-direction: row !important;
                font-size: 14px;
                padding: 15px 22px !important;
            }
            .nav-icon { display: none; }
        }
        </style>
    @yield('style')
</head>

<body>
<?php $routeName = request()->route()->getName(); ?>

        <!-- Begin page -->
<div id="layout-wrapper">

    <nav class="main-header navbar navbar-light navbar-white">
        <div class="container">

            <ul class="navbar-nav mobile-nav">
                <li class="nav-item {{ in_array($routeName, ['sale.summaryOrders']) ? 'active' : '' }}">
                    <a href="{{ route('sale.summaryOrders', ['status' => 'NEW']) }}" class="nav-link">
                        <i class="mdi mdi-view-dashboard-outline nav-icon"></i>
                        Summary
                    </a>
                </li>
                <li class="nav-item {{ in_array($routeName, ['sale.availableProducts']) ? 'active' : '' }}">
                    <a href="{{ route('sale.availableProducts') }}" class="nav-link">
                        <i class="mdi mdi-format-list-checkbox nav-icon"></i>
                        Available
                    </a>
                </li>
                <li class="nav-item {{ in_array($routeName, ['sale.home']) ? 'active' : '' }}">
                    <a href="{{ route('sale.home') }}" class="nav-link">
                        <i class="mdi mdi-chart-areaspline nav-icon"></i>
                        Reports
                    </a>
                </li>

                {{-- TOGGLE --}}
                <li class="nav-item nav-more">
                    <a href="#" class="nav-link" onclick="toggleMoreMenu(event)">
                        <i class="mdi mdi-menu nav-icon"></i>
                        More
                    </a>
                </li>

                <div class="desktop-nav">
                    {{-- ALWAYS SHOW --}}
                    

                    <li class="nav-item {{ in_array($routeName, ['sale.order']) ? 'active' : '' }}">
                        <a href="{{ route('sale.order', ['status' => 'NEW']) }}" class="nav-link">Orders</a>
                    </li>



                    <li class="nav-item {{ in_array($routeName, ['sale.product']) ? 'active' : '' }}">
                        <a href="{{ route('sale.product') }}" class="nav-link">Products</a>
                    </li>



                    <li class="nav-item {{ in_array($routeName, ['sale.product.batch.upload']) ? 'active' : '' }}">
                        <a href="{{ route('sale.product.batch.upload') }}" class="nav-link">
                            Upload
                        </a>
                    </li>

                    <li class="nav-item">
                        <a href="/salev2" class="nav-link text-primary font-weight-bold">
                            <i class="mdi mdi-swap-horizontal mr-1"></i> Switch to V2
                        </a>
                    </li>
                </div>
            </ul>

            {{-- HIDDEN MENU --}}
            <ul class="navbar-nav mobile-nav-more" id="mobileMoreMenu">

                {{-- ALWAYS SHOW --}}
              
                <li class="nav-item {{ in_array($routeName, ['sale.order']) ? 'active' : '' }}">
                    <a href="{{ route('sale.order', ['status' => 'NEW']) }}" class="nav-link">
                        <i class="mdi mdi-clipboard-text-outline mr-2"></i> Orders
                    </a>
                </li>
                <li class="nav-item {{ in_array($routeName, ['sale.product']) ? 'active' : '' }}">
                    <a href="{{ route('sale.product') }}" class="nav-link">
                        <i class="mdi mdi-package-variant-closed mr-2"></i> Products
                    </a>
                </li>
                <li class="nav-item {{ in_array($routeName, ['sale.product.batch.upload']) ? 'active' : '' }}">
                    <a href="{{ route('sale.product.batch.upload') }}" class="nav-link">
                        <i class="mdi mdi-cloud-upload-outline mr-2"></i> Upload Batch
                    </a>
                </li>
                <li class="nav-item {{ in_array($routeName, ['sale.product.create']) ? 'active' : '' }}">
                    <a href="{{ route('sale.product.create') }}" class="nav-link">
                        <i class="mdi mdi-plus-box-outline mr-2"></i> Create Product
                    </a>
                </li>
                <li class="nav-item">
                    <a href="/salev2" class="nav-link text-primary">
                        <i class="mdi mdi-swap-horizontal mr-2"></i> Switch to V2
                    </a>
                </li>
            </ul>

        </div>
    </nav>


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
<script src="{{ asset('admin-theme/assets/js/jquery.min.js') }}"></script>
<script src="{{ asset('admin-theme/assets/js/bootstrap.bundle.min.js') }}"></script>
<script src="{{ asset('admin-theme/assets/js/metismenu.min.js') }}"></script>
<script src="{{ asset('admin-theme/assets/js/waves.js') }}"></script>
<script src="{{ asset('admin-theme/assets/js/simplebar.min.js') }}"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/axios/1.6.8/axios.min.js" integrity="sha512-PJa3oQSLWRB7wHZ7GQ/g+qyv6r4mbuhmiDb8BjSFZ8NZ2a42oTtAq5n0ucWAwcQDlikAtkub+tPVCw4np27WCg==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
<!-- Morris Js-->
<script src="{{ asset('admin-theme/plugins/morris-js/morris.min.js') }}"></script>
<!-- Raphael Js-->
<script src="{{ asset('admin-theme/plugins/raphael/raphael.min.js') }}"></script>

<!-- Morris Custom Js-->
<script src="{{ asset('admin-theme/assets/pages/dashboard-demo.js') }}"></script>

<!-- App js -->
<script src="{{ asset('admin-theme/assets/js/theme.js') }}"></script>
<script>
    function toggleMoreMenu(e) {
        e.preventDefault();
        document.getElementById('mobileMoreMenu').classList.toggle('show');
    }
</script>

<script>

    const route = "{{route('sale.availableProducts')}}"

    const downloadAs = (url, name) => {
        axios.get(url, {
            headers: {
                "Content-Type": "application/octet-stream"
            },
            responseType: "blob"
        })
            .then(response => {
                const a = document.createElement("a");
                const url = window.URL.createObjectURL(response.data);
                a.href = url;
                a.download = name;
                a.click();
            })
            .catch(err => {
                console.log("error", err);
            });
    };

    $(document).ready(function () {


        $('#download-all-images').on('click', () => {
            var images = document.getElementsByTagName('img');
            var srcList = [];
            var i = 0;

            setInterval(function(){
                if(images.length > i){
                    downloadAs(images[i].src, $(images[i]).data('name'))
                    i++;
                }
            },200);
        })
    })

   const images = [
       "https://tadimedia.vn/storage/sale/xDh1csQ1TsRcta600hVPJ6vX5edqWDK0Q90oKm0J.jpg",
    "https://tadimedia.vn/storage/sale/q2dc9MJYfYBtebVovgeGsPVUUQDtNWOwuA1pS6jE.jpg",
   ]


    let blodImages = []
    function setCanvasImage(path,func){

        let img = new Image
        let c = document.createElement('canvas')
        let ctx = c.getContext('2d')

        img.onload = function(){
            c.width = this.naturalWidth
            c.height = this.naturalHeight
            ctx.drawImage(this,0,0)
            c.toBlob(blob=>{
                func(blob)
            },'image/png')
        }
        img.src = path
    }
    //
    // images.forEach(image => {
    //     setCanvasImage(image,(imgBlob)=>{
    //         console.log('doing it!')
    //         blodImages.push(imgBlob)
    //         console.log('blodImages', blodImages)
    //         if (blodImages.length === images.length) {
    //             navigator.clipboard.write(
    //                 blodImages.map(blod => {
    //                     return new ClipboardItem({'image/png': blod})
    //                 })
    //
    //             )
    //                 .then(e=>{console.log('Image copied to clipboard')})
    //                 .catch(e=>{console.log(e)})
    //         }
    //
    //     })
    // })

</script>
<script>
    const STORAGE_KEY = 'selectedProductImages';

    document.addEventListener('DOMContentLoaded', () => {
        const stored = JSON.parse(localStorage.getItem(STORAGE_KEY)) || [];

        // Auto check các checkbox đã lưu
        document.querySelectorAll('.product-checkbox').forEach(cb => {
            const id = cb.dataset.id;
            if (stored.find(p => p.id == id)) {
                cb.checked = true;
            }
        });
    });

    function toggleSelection(checkbox) {
        let stored = JSON.parse(localStorage.getItem(STORAGE_KEY)) || [];
        const id = checkbox.dataset.id;
        const img = checkbox.dataset.img;

        if (checkbox.checked) {
            if (!stored.find(p => p.id == id)) {
                stored.push({ id, img });
            }
        } else {
            stored = stored.filter(p => p.id != id);
        }

        localStorage.setItem(STORAGE_KEY, JSON.stringify(stored));
    }

    const copyImagesBtn = document.getElementById('copy-images');
    if (copyImagesBtn) {
        copyImagesBtn.addEventListener('click', () => {
            const stored = JSON.parse(localStorage.getItem(STORAGE_KEY)) || [];
            if (stored.length === 0) {
                alert("Bạn chưa chọn ảnh nào!");
                return;
            }

            const urls = stored.map(p => p.img).join('\n');
            navigator.clipboard.writeText(urls).then(() => {
                alert(`Đã copy ${stored.length} ảnh vào clipboard!`);
            }).catch(err => {
                alert("Lỗi khi copy vào clipboard: " + err);
            });
        });
    }
</script>
    @yield('script')
</body>

</html>
