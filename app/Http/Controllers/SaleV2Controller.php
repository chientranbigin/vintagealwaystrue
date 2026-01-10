<?php

namespace App\Http\Controllers;

use App\Consumer;
use App\Contracts\SaleConfig;
use App\Import;
use App\Product;
use App\OrderProduct;
use App\Order;
use App\ProductImage;
use App\Services\GoogleVisionService;
use App\Size;
use Carbon\Carbon;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use GuzzleHttp\Client;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Facades\Log;

class SaleV2Controller extends Controller
{
    public function index()
    {
        return view('salev2.index');
    }

    /** API Endpoints for SaleV2 SPA **/

    public function products(Request $request)
    {
        $query = Product::with(['images', 'sizes'])->select('products.*')->distinct();

        if ($request->filled('search')) {
            $names = array_map(function ($item) { return trim($item); }, explode(',', $request->search));
            $query->where(function ($q) use ($names) {
                foreach ($names as $name) {
                    $q->orWhere('name', 'like', "%$name%");
                }
            });
        }

        if ($request->filled('status')) {
            $query->where('products.status', $request->status);
            if ($request->status === 'AVAILABLE') {
                $query->whereNotNull('path_thumb')->where('path_thumb', '!=', '');
            }
        }

        if ($request->filled('type')) {
            $query->where('type', $request->type);
        }

        // Legacy Size Filter 1
        if ($request->from_search_key) {
            $query->join('sizes', function ($join) use ($request) {
                $join->on('sizes.product_id', '=', 'products.id')
                    ->where('sizes.name', '=', $request->from_search_key);
                if ($request->from_search_value) {
                    $join->where('sizes.value', '>=', $request->from_search_value);
                }
                if ($request->to_search_value) {
                    $join->where('sizes.value', '<=', $request->to_search_value);
                }
            });
        }

        // Legacy Size Filter 2
        if ($request->from_search_key_2) {
            $query->join(DB::raw('sizes as sizes_2'), function ($join) use ($request) {
                $join->on('sizes_2.product_id', '=', 'products.id')
                    ->where('sizes_2.name', '=', $request->from_search_key_2);
                if ($request->from_search_value_2) {
                    $join->where('sizes_2.value', '>=', $request->from_search_value_2);
                }
                if ($request->to_search_value_2) {
                    $join->where('sizes_2.value', '<=', $request->to_search_value_2);
                }
            });
        }

        $products = $query->latest('products.id')->paginate($request->get('limit', 60));

        $products->getCollection()->transform(function($product) {
            if ($product->path_thumb) {
                $product->path_thumb = asset($product->path_thumb);
            }
            return $product;
        });

        return response()->json($products);
    }

    public function orders(Request $request)
    {
        $query = Order::with('products');

        if ($request->filled('status')) {
            if ($request->status === 'NEW') {
                $query->whereIn('status', ['NEW', 'ON HOLD']);
            } else {
                $query->where('status', $request->status);
            }
        }

        if ($request->filled('search')) {
            $search = $request->search;
            $query->where(function($q) use ($search) {
                $q->where('name', 'like', "%$search%")
                  ->orWhere('phone', 'like', "%$search%");
            });
        }

        $orders = $query->latest()->paginate($request->get('limit', 10));

        // Legacy calculation: Viettel Amount = total_amount - deposit_amount (if not is_paid)
        $orders->getCollection()->transform(function($order) {
            $order->is_old = Carbon::parse($order->created_at)->diffInDays(now()) >= 3;
            // Legacy: if is_paid then 0, else total - deposit
            $order->viettel_amount = $order->is_paid ? 0 : ($order->total_amount - $order->deposit_amount);
            
            // Transform product images to absolute URLs
            foreach ($order->products as $product) {
                if ($product->path_thumb) {
                    $product->path_thumb = asset($product->path_thumb);
                }
            }
            
            $order->created_at_human = $order->created_at->diffForHumans();
            return $order;
        });

        return response()->json($orders);
    }

    public function orderDetail($id)
    {
        $order = Order::with('products')->findOrFail($id);
        return response()->json($order);
    }

    public function productDetail($id)
    {
        $product = Product::with(['images', 'sizes'])->findOrFail($id);
        return response()->json($product);
    }

    public function storeOrder(Request $request)
    {
        $validated = $request->validate([
            'name' => 'required',
            'phone' => 'required',
            'products' => 'required|array',
            'total_amount' => 'required|numeric',
            'expenses' => 'nullable|array',
        ]);

        return DB::transaction(function () use ($request, $validated) {
            $order = Order::create($request->all());
            $order->products()->attach($request->products);

            // Targeted status update
            Product::whereIn('id', $request->products)->update(['status' => 'SOLD']);

            // Handle Price Updates
            if ($request->has('products_info')) {
                foreach ($request->products_info as $pInfo) {
                    if (isset($pInfo['id']) && isset($pInfo['price'])) {
                        Product::where('id', $pInfo['id'])->update(['price' => $pInfo['price']]);
                    }
                }
            }

            return response()->json(['success' => true, 'order' => $order]);
        });
    }



    public function customerHistory($phone)
    {
        $orders = Order::where('phone', $phone)->with('products')->latest()->get();
        return response()->json($orders);
    }

    public function dashboardStats()
    {
        $now = Carbon::now();
        $startOfMonth = $now->copy()->startOfMonth();
        $endOfMonth = $now->copy()->endOfMonth();
        $startOfPrevMonth = $now->copy()->subMonth()->startOfMonth();
        $endOfPrevMonth = $now->copy()->subMonth()->endOfMonth();

        // Legacy "All" Summary
        $allOrders = Order::all();
        $typeStats = Product::groupBy('type', 'status')
            ->select('type', 'status', DB::raw('count(*) as count'))
            ->get();

        // New/Hold orders (Legacy logic for stat cards)
        $newOrders = Order::whereIn('status', ['NEW', 'ON HOLD'])->get();
        
        // Current Month Summary
        $currentMonthOrders = Order::whereBetween('created_at', [$startOfMonth, $endOfMonth])->get();
        $currentTypeStats = DB::table('products')
            ->join('order_products', 'products.id', '=', 'order_products.product_id')
            ->join('orders', 'orders.id', '=', 'order_products.order_id')
            ->whereBetween('orders.created_at', [$startOfMonth, $endOfMonth])
            ->where('products.status', 'SOLD')
            ->select('products.type', DB::raw('count(*) as count'))
            ->groupBy('products.type')
            ->get();

        // Previous Month Summary
        $prevMonthOrders = Order::whereBetween('created_at', [$startOfPrevMonth, $endOfPrevMonth])->get();
        $prevTypeStats = DB::table('products')
            ->join('order_products', 'products.id', '=', 'order_products.product_id')
            ->join('orders', 'orders.id', '=', 'order_products.order_id')
            ->whereBetween('orders.created_at', [$startOfPrevMonth, $endOfPrevMonth])
            ->where('products.status', 'SOLD')
            ->select('products.type', DB::raw('count(*) as count'))
            ->groupBy('products.type')
            ->get();

        return response()->json([
            'new_orders_count' => $newOrders->count(),
            'new_orders_revenue' => $newOrders->sum('final_amount'),
            'month_orders_count' => $currentMonthOrders->count(),
            'month_orders_revenue' => $currentMonthOrders->sum('final_amount'),

            'all' => [
                'total_orders' => $allOrders->count(),
                'total_revenue' => $allOrders->sum('final_amount'),
                'type_stats' => $typeStats
            ],
            'current_month' => [
                'total_orders' => $currentMonthOrders->count(),
                'total_revenue' => $currentMonthOrders->sum('final_amount'),
                'type_stats' => $currentTypeStats
            ],
            'prev_month' => [
                'total_orders' => $prevMonthOrders->count(),
                'total_revenue' => $prevMonthOrders->sum('final_amount'),
                'type_stats' => $prevTypeStats
            ],
            
            'order_counts' => Order::select('status', DB::raw('count(*) as count'))
                ->groupBy('status')->get(),
            'monthly_sales' => Order::select(DB::raw("DATE_FORMAT(created_at, '%Y-%m') as month"), DB::raw('count(*) as number'), DB::raw('sum(final_amount) as final_amount'))
                ->where('created_at', '>=', now()->subMonths(12))
                ->groupBy('month')
                ->orderBy('month', 'ASC')
                ->get(),
        ]);
    }

    public function consumers(Request $request)
    {
        $query = Consumer::withCount('orders');

        if ($request->filled('search')) {
            $search = $request->search;
            $query->where(function($q) use ($search) {
                $q->where('name', 'like', "%$search%")
                  ->orWhere('phone', 'like', "%$search%")
                  ->orWhere('address', 'like', "%$search%");
            });
        }

        return response()->json($query->with('orders')->withCount('orders')->orderBy('id', 'DESC')->paginate(20));
    }

    public function consumer($id)
    {
        return response()->json(Consumer::findOrFail($id));
    }

    public function storeProduct(Request $request)
    {
        $validated = $request->validate([
            'type' => 'required',
            'banner' => 'required|file|image',
            'sizes' => 'array'
        ]);

        return DB::transaction(function () use ($request) {
            // Upload File
            $file = $request->file('banner');
            $path = str_replace('public', 'storage', $file->store('public/sale'));

            // Auto-generate Name (Code)
            $lastProduct = Product::where('type', $request->type)->latest()->first();
            $lastId = $lastProduct ? (int)filter_var($lastProduct->name, FILTER_SANITIZE_NUMBER_INT) : 0;
            $nextId = $lastId + 1;
            $prefix = \App\Contracts\SaleConfig::PRODUCT_KEY_TYPES[$request->type] ?? $request->type;
            
            $product = Product::create([
                'type' => $request->type,
                'name' => $prefix . $nextId,
                'price' => $request->price ?? 0,
                'description' => $request->description,
                'path_thumb' => $path,
                'status' => 'AVAILABLE'
            ]);

            // Save Sizes
            if($request->sizes) {
                foreach ($request->sizes as $size) {
                    if (!empty($size['value'])) {
                        Size::create([
                            'product_id' => $product->id,
                            'name' => $size['name'],
                            'value' => $size['value'],
                        ]);
                    }
                }
            }

            return response()->json(['success' => true, 'product' => $product]);
        });
    }

    public function updateProduct(Request $request, $id)
    {
        $product = Product::findOrFail($id);
        
        return DB::transaction(function () use ($request, $product) {
            // Update Banner if new file
            if ($file = $request->file('banner')) {
                $path = str_replace('public', 'storage', $file->store('public/sale'));
                $product->path_thumb = $path;
            }
            
            // Update other fields
            $product->price = $request->price;
            $product->description = $request->description;
            // Only update Name/Type if strictly necessary, usually we don't change product code/type
            // but if user wants to fix a mistake:
            if($request->type && $request->type !== $product->type) {
                $product->type = $request->type;
                // Should we regenerate code? Probably safer NOT to change Code unless explicit.
            }
            $product->save();

            // Update Sizes (Delete all and re-create)
            if($request->sizes) {
                Size::where('product_id', $product->id)->delete();
                foreach ($request->sizes as $size) {
                    if (!empty($size['value'])) {
                        Size::create([
                            'product_id' => $product->id,
                            'name' => $size['name'],
                            'value' => $size['value'],
                        ]);
                    }
                }
            }

            return response()->json(['success' => true, 'product' => $product]);
        });
    }

    public function destroyProduct($id)
    {
        $product = Product::findOrFail($id);
        
        return DB::transaction(function() use ($product) {
            foreach ($product->images as $image) {
                $path = str_replace('/storage/', 'public/', $image->file_path);
                Storage::disk('public')->delete($path);
            }
            $product->images()->delete();
            
            if ($product->path_thumb) {
                $path = str_replace('/storage/', 'public/', $product->path_thumb);
                Storage::disk('public')->delete($path);
            }
            
            $product->sizes()->delete();
            $product->delete();
            
            return response()->json(['success' => true]);
        });
    }

    public function updateOrder(Request $request, $id)
    {
        $order = Order::with('products')->findOrFail($id);
        
        return DB::transaction(function() use ($request, $order) {
            // 1. Handle Consumer linking
            if ($request->phone) {
                $consumer = Consumer::firstOrCreate(
                    ['phone' => $request->phone],
                    ['name' => $request->name, 'address' => $request->address]
                );
                // Optional: Update consumer info if changed
                $consumer->update([
                    'name' => $request->name,
                    'address' => $request->address
                ]);
                $order->consumer_id = $consumer->id;
            }

            // 2. Update Order Fields
            $order->fill($request->only([
                'name', 'phone', 'address', 'status', 
                'additional_amount', 'deposit_amount', 
                'is_freeship', 'is_paid', 'note',
                'total_amount', 'final_amount' // Trusting frontend calc or we can recalc
            ]));
            
            // Recalc final amount if not provided but dependencies are?
            // For safety, let's allow frontend to drive this for now as it has simple logic.
            
            $order->save();

            // 3. Sync Products and Update Statuses
            if ($request->has('products')) {
                $newIds = $request->products;
                $oldIds = $order->products->pluck('id')->toArray();

                // Products to become AVAILABLE (removed)
                $removedIds = array_diff($oldIds, $newIds);
                if (!empty($removedIds)) {
                    Product::whereIn('id', $removedIds)->update(['status' => 'AVAILABLE']);
                }

                // Products to become SOLD (added)
                $addedIds = array_diff($newIds, $oldIds);
                if (!empty($addedIds)) {
                    Product::whereIn('id', $addedIds)->update(['status' => 'SOLD']);
                }

                $order->products()->sync($newIds);
                
                // Handle Price Updates if provided in a separate array or structure
                // Request 'products_info' could contain [{id, price}]
                if ($request->has('products_info')) {
                    foreach ($request->products_info as $pInfo) {
                        if (isset($pInfo['id']) && isset($pInfo['price'])) {
                            Product::where('id', $pInfo['id'])->update(['price' => $pInfo['price']]);
                        }
                    }
                }
            }
            
            // 4. Handle simple status update (if payload only has status)
            // The fill() above handles it if 'status' is in request.
            
            return response()->json(['success' => true, 'order' => $order]);
        });
    }

    public function smartUpload(Request $request)
    {
        // Implementation of folder-style upload grouped by product
        // Placeholder for now, will integrate Google Vision/OpenAI logic
        return response()->json(['message' => 'Smart Upload logic initialized']);
    }
    public function uploadSmartGroup(Request $request)
    {
        $sessionId = $request->input('session_id', 'SMART_' . uniqid());
        $type = $request->input('type', 'T-SHIRT');
        
        // 1. Store Main Image
        $mainPath = '';
        $mainOriginalName = '';
        if ($request->hasFile('main')) {
            $file = $request->file('main');
            $mainOriginalName = $file->getClientOriginalName();
            // Use a unique name to avoid conflicts if same file is uploaded twice in different groups
            $filename = 'main_' . uniqid() . '_' . $mainOriginalName;
            $mainPath = str_replace('public', 'storage', $file->store('public/sale'));

        } else {
            return response()->json(['success' => false, 'message' => 'Main image is required'], 400);
        }

        // 2. Store Details
        $detailPaths = [];
        if ($request->hasFile('details')) {
            foreach ($request->file('details') as $file) {
                $filename = 'detail_' . uniqid() . '_' . $file->getClientOriginalName();
                $path = str_replace('public', 'storage', $file->store('public/sale/detail'));

                $detailPaths[] = $path;
            }
        }

        // 3. Dispatch Job (Step 1 Recognition & Step 2 Attachment)
        \App\Jobs\ProcessSmartImage::dispatch(
            $mainPath, 
            $detailPaths, 
            $type, 
            $sessionId, 
            $mainOriginalName
        );

        return response()->json([
            'success' => true, 
            'session_id' => $sessionId,
            'status' => 'DISPATCHED'
        ]);
    }

    public function uploadLogs(Request $request)
    {
        $sessionIds = $request->input('session_ids', []);
        if (is_string($sessionIds)) {
            $sessionIds = explode(',', $sessionIds);
        }

        $logs = \App\ProductUploadLog::whereIn('session_id', $sessionIds)->get();
        return response()->json($logs);
    }
}
