<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login - SaleV2 Enterprise</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
    <style>
        body { font-family: 'Inter', sans-serif; }
    </style>
</head>
<body class="bg-slate-50 min-h-screen flex items-center justify-center p-6">
    <div class="max-w-md w-full">
        <!-- Logo / Brand -->
        <div class="text-center mb-10">
            <div class="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-blue-600 to-indigo-700 rounded-2xl shadow-xl text-white text-2xl font-bold mb-4">V2</div>
            <h1 class="text-3xl font-bold text-slate-900 tracking-tight">Welcome Back</h1>
            <p class="text-slate-500 mt-2 text-lg">Sign in to access SaleV2 Dashboard</p>
        </div>

        <!-- Login Card -->
        <div class="bg-white rounded-3xl shadow-xl shadow-slate-200/50 border border-slate-100 p-8 md:p-10">
            <form method="POST" action="{{ route('login') }}" class="space-y-6">
                @csrf
                
                <div>
                    <label for="email" class="block text-sm font-bold text-slate-700 mb-2 uppercase tracking-wide">Email Address</label>
                    <input id="email" type="email" name="email" value="{{ old('email') }}" required autofocus
                           class="w-full px-4 py-4 rounded-xl border border-slate-200 focus:ring-4 focus:ring-blue-100 focus:border-blue-500 transition-all outline-none text-slate-800"
                           placeholder="admin@alwaystrue.com">
                    @if ($errors->has('email'))
                        <p class="text-red-500 text-xs mt-2 font-medium">{{ $errors->first('email') }}</p>
                    @endif
                </div>

                <div>
                    <div class="flex justify-between mb-2">
                        <label for="password" class="block text-sm font-bold text-slate-700 uppercase tracking-wide">Password</label>
                    </div>
                    <input id="password" type="password" name="password" required
                           class="w-full px-4 py-4 rounded-xl border border-slate-200 focus:ring-4 focus:ring-blue-100 focus:border-blue-500 transition-all outline-none text-slate-800"
                           placeholder="••••••••">
                    @if ($errors->has('password'))
                        <p class="text-red-500 text-xs mt-2 font-medium">{{ $errors->first('password') }}</p>
                    @endif
                </div>

                <div class="flex items-center">
                    <input id="remember" type="checkbox" name="remember" {{ old('remember') ? 'checked' : '' }}
                           class="w-5 h-5 rounded border-slate-300 text-blue-600 focus:ring-blue-500">
                    <label for="remember" class="ml-3 text-sm font-medium text-slate-600">Keep me signed in</label>
                </div>

                <button type="submit" 
                        class="w-full bg-slate-900 hover:bg-black text-white font-bold py-4 rounded-xl shadow-lg transition-all transform active:scale-[0.98]">
                    Sign In
                </button>
            </form>
        </div>

        <p class="text-center mt-8 text-slate-400 text-sm">
            &copy; 2026 Vintage Always True. All Rights Reserved.
        </p>
    </div>
</body>
</html>
