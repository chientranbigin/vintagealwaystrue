<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SaleV2 - Enterprise Edition</title>
    <meta name="csrf-token" content="{{ csrf_token() }}">
    
    <!-- CDNs for Reliability & Speed -->
    <script src="https://cdn.tailwindcss.com"></script>
    <link rel="stylesheet" href="https://unpkg.com/element-ui/lib/theme-chalk/index.css">
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
    
    <style>
        body { font-family: 'Inter', sans-serif; background-color: #f8fafc; padding-bottom: 70px; }
        @media (min-width: 768px) {
            body { padding-bottom: 0; }
        }
        /* Fix Element UI conflicts with Tailwind */
        .el-button [class*="el-icon-"] + span { margin-left: 5px; }
        .router-link-active { background-color: #eff6ff; color: #2563eb; }
        .router-link-active i { color: #2563eb; }
        
        /* Mobile specific Hides */
        @media (max-width: 767px) {
            .mobile-hide { display: none !important; }
        }
        
        /* POS Button Active Override */
        .router-link-active.pos-link { background-color: transparent !important; }
        .router-link-active.pos-link i { color: white !important; }
    </style>
</head>
<body>
    <div id="app">
        <div class="min-h-screen bg-slate-50 flex flex-col">
            <!-- Desktop Header Navigation -->
            <header class="hidden md:flex w-full bg-white border-b border-slate-200 h-16 items-center justify-between px-6 shadow-sm z-30 relative">
                <!-- Logo -->
                <div class="flex items-center gap-3">
                    <div class="w-8 h-8 bg-gradient-to-br from-blue-600 to-indigo-700 rounded-lg flex items-center justify-center text-white font-bold shadow-md text-sm">V2</div>
                    <span class="font-bold text-slate-800 text-lg tracking-tight">SaleV2</span>
                </div>

                <!-- Horizontal Nav -->
                <nav class="flex items-center gap-1 h-full">
                    <router-link to="/salev2" exact class="flex items-center gap-2 px-4 h-10 rounded-lg transition-all hover:bg-slate-50 text-slate-600 hover:text-blue-600 font-medium text-sm">
                        <i class="el-icon-odometer text-lg"></i>
                        <span>Dashboard</span>
                    </router-link>
                    <router-link to="/salev2/order/create" class="flex items-center gap-2 px-4 h-10 rounded-lg transition-all hover:bg-slate-50 text-slate-600 hover:text-blue-600 font-medium text-sm">
                        <i class="el-icon-plus text-lg font-bold"></i>
                         <span>POS</span>
                    </router-link>
                    <router-link to="/salev2/orders" class="flex items-center gap-2 px-4 h-10 rounded-lg transition-all hover:bg-slate-50 text-slate-600 hover:text-blue-600 font-medium text-sm">
                        <i class="el-icon-shopping-cart-2 text-lg"></i>
                        <span>Orders</span>
                    </router-link>
                    <router-link to="/salev2/products" class="flex items-center gap-2 px-4 h-10 rounded-lg transition-all hover:bg-slate-50 text-slate-600 hover:text-blue-600 font-medium text-sm">
                        <i class="el-icon-box text-lg"></i>
                        <span>Products</span>
                    </router-link>
                     <router-link to="/salev2/customers" class="flex items-center gap-2 px-4 h-10 rounded-lg transition-all hover:bg-slate-50 text-slate-600 hover:text-blue-600 font-medium text-sm">
                        <i class="el-icon-user text-lg"></i>
                        <span>Customers</span>
                    </router-link>
                    <div class="w-px h-6 bg-slate-200 mx-2"></div>
                    <router-link to="/salev2/upload" class="flex items-center gap-2 px-4 h-10 rounded-lg transition-all hover:bg-slate-50 text-slate-600 hover:text-blue-600 font-medium text-sm">
                        <i class="el-icon-upload text-lg"></i>
                        <span>Upload</span>
                    </router-link>
                </nav>

                <!-- Profile -->
                <div class="flex items-center gap-3 pl-4 border-l border-slate-100">
                    <div class="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 font-bold text-xs ring-2 ring-white shadow-sm">A</div>
                    <div class="hidden lg:block">
                        <p class="text-xs font-bold text-slate-800">Admin</p>
                    </div>
                </div>
            </header>

            <!-- Main Content Area -->
            <main class="flex-1 overflow-y-auto">
                <router-view></router-view>
            </main>

            <!-- Bottom Navigation (Mobile) -->
            <nav class="md:hidden fixed bottom-0 left-0 right-0 bg-white border-t border-slate-200 flex justify-around items-center h-16 z-50 px-2 shadow-[0_-2px_10px_rgba(0,0,0,0.05)]">
                <router-link to="/salev2" exact class="flex flex-col items-center justify-center w-full h-full text-slate-400">
                    <i class="el-icon-odometer text-xl"></i>
                    <span class="text-[10px] font-bold uppercase mt-1">Dash</span>
                </router-link>
                <router-link to="/salev2/orders" class="flex flex-col items-center justify-center w-full h-full text-slate-400">
                    <i class="el-icon-shopping-cart-2 text-xl"></i>
                    <span class="text-[10px] font-bold uppercase mt-1">Orders</span>
                </router-link>
                <router-link to="/salev2/order/create" class="flex flex-col items-center justify-center w-full h-full text-slate-400 pos-link">
                     <div class="w-10 h-10 bg-blue-600 rounded-full flex items-center justify-center text-white shadow-lg shadow-blue-200 mb-1">
                        <i class="el-icon-plus text-xl font-bold"></i>
                     </div>
                </router-link>
                 <router-link to="/salev2/products" class="flex flex-col items-center justify-center w-full h-full text-slate-400">
                    <i class="el-icon-box text-xl"></i>
                    <span class="text-[10px] font-bold uppercase mt-1">Products</span>
                </router-link>
                <router-link to="/salev2/customers" class="flex flex-col items-center justify-center w-full h-full text-slate-400">
                    <i class="el-icon-user text-xl"></i>
                    <span class="text-[10px] font-bold uppercase mt-1">Customers</span>
                </router-link>
            </nav>
        </div>
    </div>

    <!-- Scripts -->
    <script src="https://unpkg.com/vue@2.6.14/dist/vue.js"></script>
    <script src="https://unpkg.com/vue-router@3.5.3/dist/vue-router.js"></script>
    <script src="https://unpkg.com/element-ui/lib/index.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/apexcharts"></script>
    <script src="https://cdn.jsdelivr.net/npm/vue-apexcharts"></script>
    <script src="https://cdn.jsdelivr.net/npm/axios/dist/axios.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/html2canvas/1.4.1/html2canvas.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/lodash@4.17.21/lodash.min.js"></script>

    <script src="{{ mix('js/app.js') }}"></script>
</body>
</html>
