<?php

namespace App\Http\Controllers;

use App\Album;
use App\AlbumType;
use App\Consumer;
use App\Contracts\ElementConfig;
use App\Contracts\SaleConfig;
use App\DetailPackage;
use App\HtmlElement;
use App\Import;
use App\Package;
use App\Product;
use App\OrderProduct;
use App\Order;
use App\Services\GoogleVisionService;
use App\Size;
use Carbon\Carbon;
use Illuminate\Foundation\Bus\DispatchesJobs;
use Illuminate\Http\Request;
use Illuminate\Routing\Controller as BaseController;
use App\ProductUploadLog;
use Illuminate\Foundation\Validation\ValidatesRequests;
use Illuminate\Foundation\Auth\Access\AuthorizesRequests;
use Illuminate\Support\Facades\DB;
use GuzzleHttp\Client;
use Illuminate\Support\Facades\Storage;

class SaleController extends Controller
{
    public function home()
    {
        $orders = Order::with('products')->get();
        $products = Product::all();
        $totalImportPrice = $products->sum('import_price');
        $totalAllImportPrice = Import::all()->sum('price');
        $totalProduct = Product::all()->count();

        $totalPrice = $orders->sum('total_amount');
        $totalFinalPrice = $orders->sum('final_amount');
        $totalOrder = $orders->count();
        $typeProduct = \DB::table('products')
            ->select([
                \DB::raw('count(id) as number'),
                \DB::raw('sum(price) as price'),
                'type',
                'status',
            ])
            ->whereNotNull('path_thumb')
            ->groupBy('type', 'status')
            ->get();

        $orderNew = Order::where('status', 'NEW')->get();
        $totalOrderNew = $orderNew->count();
        $totalOrderNewPrice = $orderNew->sum('final_amount');

        $currentOrders = Order::with('products')->whereBetween('created_at', [
            Carbon::now()->startOfMonth(),
            Carbon::now()->endOfMonth()
        ])->get();

        $currentTotalFinalPrice = $currentOrders->sum('final_amount');
        $currentTotalOrder = $currentOrders->count();

        $previousOrders = Order::with('products')->whereBetween('created_at', [
            Carbon::now()->subMonth()->startOfMonth(),
            Carbon::now()->subMonth()->endOfMonth()
        ])->get();

        $previousTotalFinalPrice = $previousOrders->sum('final_amount');
        $previousTotalOrder = $previousOrders->count();

        $currentProducts = Product::join('order_products', 'products.id', '=', 'order_products.product_id')
            ->join('orders', 'orders.id', '=', 'order_products.order_id')
            ->whereBetween('orders.created_at', [
                Carbon::now()->startOfMonth(),
                Carbon::now()->endOfMonth()
            ])
            ->select([
                \DB::raw('count(products.id) as number'),
                \DB::raw('sum(price) as price'),
                'type',
                'products.status',
            ])
            ->groupBy('type', 'status')
            ->get();

        $previousProducts = Product::join('order_products', 'products.id', '=', 'order_products.product_id')
            ->join('orders', 'orders.id', '=', 'order_products.order_id')
            ->whereBetween('orders.created_at', [
                Carbon::now()->subMonth()->startOfMonth(),
                Carbon::now()->subMonth()->endOfMonth()
            ])
            ->select([
                \DB::raw('count(products.id) as number'),
                \DB::raw('sum(price) as price'),
                'type',
                'products.status',
            ])
            ->groupBy('type', 'status')
            ->get();

        $orderByMonth = Order::select([
            \DB::raw('count(id) as number'),
            \DB::raw('sum(final_amount) as final_amount'),
            \DB::raw(" DATE_FORMAT(created_at, '%Y-%m') AS MONTH"),
        ])
            ->groupBy('MONTH')
            ->get();

        return view('sale.index', compact('totalImportPrice',
            'totalProduct',
            'totalPrice',
            'typeProduct',
            'totalFinalPrice',
            'totalOrder',
            'totalOrderNew',
            'totalAllImportPrice',
            'currentOrders',
            'currentTotalOrder',
            'currentTotalFinalPrice',
            'totalOrderNewPrice',
            'previousTotalOrder',
            'previousTotalFinalPrice',
            'previousProducts',
            'currentProducts',
            'orderByMonth'
        ));
    }


    public function products(Request $request)
    {
        $products = Product::with('sizes')->orderBy('status')->select('products.*')->distinct()->orderBy('products.id', 'DESC');
        if ($request->type) {
            $products->where('type', $request->type);
        }


        if ($request->search) {
            $names = array_map(
                function ($item) {
                    return trim($item);
                },
                explode(',', $request->search)
            );
            $products->where(function ($query) use ($names) {
                foreach ($names as $name) {
                    $query->orWhere('name', 'like', "%{$name}%");

                }
            });


        }

        if ($request->sold_from) {
            $from = Carbon::parse($request->sold_from)->startOfDay(); // 00:00:00
            $to = Carbon::parse($request->sold_to)->endOfDay();       // 23:59:59

            $products->join('order_products', 'products.id', '=', 'order_products.product_id')
                ->join('orders', 'orders.id', '=', 'order_products.order_id')
                ->whereBetween('orders.created_at', [$from, $to]);
        }

        if ($request->type) {
            $products->where('type', $request->type);
        }

        if ($request->status) {
            $products->where('status', $request->status);
        }


        if ($request->from_search_key || $request->to_search_key) {

            $products->join('sizes', function ($query) use ($request) {
                $query->on('sizes.product_id', 'products.id');
                if ($request->from_search_value) {
                    $query
                        ->where('sizes.name', $request->from_search_key)
                        ->where('sizes.value', '>=', $request->from_search_value);
                }
                if ($request->to_search_value) {
                    $query
                        ->where('sizes.name', $request->from_search_key)
                        ->where('sizes.value', '<=', $request->to_search_value);
                }

            });

        }

        if ($request->from_search_key_2 || $request->to_search_key_2) {

            $products->join(DB::raw('sizes sizes_2'), function ($query) use ($request) {
                $query->on('sizes_2.product_id', 'products.id');

                if ($request->from_search_value_2) {
                    $query
                        ->where('sizes_2.name', $request->from_search_key_2)
                        ->where('sizes_2.value', '>=', $request->from_search_value_2);
                }
                if ($request->to_search_value_2) {
                    $query
                        ->where('sizes_2.name', $request->from_search_key_2)
                        ->where('sizes_2.value', '<=', $request->to_search_value_2);
                }
            });

        }


        $products = $products->paginate(40);

        return view('sale.products', compact('products'));
    }

    public function availableProducts(Request $request)
    {


        $products = Product::with('sizes')->where('status', 'AVAILABLE')->whereNotNull('path_thumb')->select('products.*')->distinct()->orderBy('products.id', 'DESC');

        if ($request->search) {
            $names = array_map(
                function ($item) {
                    return trim($item);
                },
                explode(',', $request->search)
            );
            $products->where(function ($query) use ($names) {
                foreach ($names as $name) {
                    $query->orWhere('name', 'like', "%{$name}%");

                }
            });


        }

        if ($request->type) {
            $products->where('type', $request->type);
        }

        if ($request->status) {
            $products->where('status', $request->status);
        }

        if ($request->from_search_key || $request->to_search_key) {

            $products->join('sizes', function ($query) use ($request) {
                $query->on('sizes.product_id', 'products.id');
                if ($request->from_search_value) {
                    $query
                        ->where('sizes.name', $request->from_search_key)
                        ->where('sizes.value', '>=', $request->from_search_value);
                }
                if ($request->to_search_value) {
                    $query
                        ->where('sizes.name', $request->from_search_key)
                        ->where('sizes.value', '<=', $request->to_search_value);
                }

            });

        }

        if ($request->from_search_key_2 || $request->to_search_key_2) {

            $products->join(DB::raw('sizes sizes_2'), function ($query) use ($request) {
                $query->on('sizes_2.product_id', 'products.id');

                if ($request->from_search_value_2) {
                    $query
                        ->where('sizes_2.name', $request->from_search_key_2)
                        ->where('sizes_2.value', '>=', $request->from_search_value_2);
                }
                if ($request->to_search_value_2) {
                    $query
                        ->where('sizes_2.name', $request->from_search_key_2)
                        ->where('sizes_2.value', '<=', $request->to_search_value_2);
                }
            });

        }
//        dd(
//            vsprintf(
//                str_replace(array('?'),array('\'%s\''),$products->toSql()),
//                $products->getBindings()
//            ),
//            $request->all()
//        );
        $products = $products->paginate(60);


        return view('sale.available-products', compact('products'));
    }

    public function searchAvailableProducts(Request $request)
    {

        return redirect()->route($request->urlSearch ?: 'sale.availableProducts', $request->except(['_token']));
    }

    public function searchOrders(Request $request)
    {

        return redirect()->route($request->urlSearch ?: 'sale.order', $request->except(['_token']));
    }

    public function productDetail($id)
    {

        $product = Product::with('sizes')->where('id', $id)->first();


        return view('sale.product-detail', compact('product'));
    }

    public function productDelete($id)
    {
        Product::find($id)->delete();
        return redirect()->back();
    }


    public function orderComplete($id)
    {
        Order::find($id)->update([
            'status' => 'COMPLETED'
        ]);
        return redirect()->back();
    }

    public function orderShip($id)
    {
        Order::find($id)->update([
            'status' => 'A SHIP NOW',

        ]);
        return redirect()->back();
    }

    public function orderHold($id)
    {
        Order::find($id)->update([
            'status' => 'ON HOLD'
        ]);
        return redirect()->back();
    }

    public function productCreate(Request $request)
    {
        if ($request->isMethod('POST')) {

            if ($file = $request->file('banner')) {
                $request = $request->merge([
                    'path_thumb' => str_replace('public', 'storage', $file->store('public/sale'))
                ]);

            }
            $latestId = (int)filter_var((Product::where('type', $request->type)->latest()->first()->name ?? 0), FILTER_SANITIZE_NUMBER_INT) + 1;
            $request = $request->merge([
                'name' => \App\Contracts\SaleConfig::PRODUCT_KEY_TYPES[$request->type] . $latestId,
                'status' => 'AVAILABLE'
            ]);
            $product = Product::create($request->all());


            foreach ($request->sizes as $size) {
                if ($size['value']) {
                    Size::create([
                        'product_id' => $product->id,
                        'value' => $size['value'],
                        'name' => $size['name'],
                    ]);
                }

            }
            if ($request->ajax()) {
                $measurements = collect($request->sizes)
                    ->filter(fn($s) => !empty($s['value']))
                    ->map(fn($s) => "{$s['name']}: {$s['value']}")
                    ->implode(', ');

                return response()->json([
                    'success' => true,
                    'product_code' => $product->name,
                    'measurements' => $measurements
                ]);
            }

            return redirect()->route('sale.product');
        }


        return view('sale.product-create');
    }


//    public function productUpload(Request $request)
//    {
//        if ($request->isMethod('POST')) {
//            if ($files = $request->file('banner')) {
//                $requests    = [];
//                $fileMap     = [];
//                $products    = [];
//                $failedFiles = [];
//                $client      = new Client(['timeout' => 120]);
//
//                $remainingFiles = count($files);
//
//                foreach ($files as $file) {
//                    $remainingFiles--;
//
//                    // Lưu file vào storage
//                    $originalName = $file->getClientOriginalName();
//                    $storedPath   = str_replace('public', 'storage', $file->store('public/sale'));
//
//                    $requests[] = [
//                        "image" => [
//                            "source" => ["imageUri" => asset($storedPath)]
//                        ],
//                        "features" => [["type" => "TEXT_DETECTION"]],
//                    ];
//
//                    // lưu cả tên file gốc và path
//                    $fileMap[] = [
//                        'stored'   => $storedPath,
//                        'original' => $originalName,
//                    ];
//
//                    // Gửi batch mỗi 16 file hoặc khi hết file
//                    if (count($fileMap) == 16 || $remainingFiles == 0) {
//                        $response = $client->post(
//                            'https://vision.googleapis.com/v1/images:annotate?key=' . env('GOOGLE_VISION_KEY'),
//                            ['json' => ["requests" => $requests]]
//                        );
//                        $results = json_decode($response->getBody(), true);
//
//                        // Xử lý kết quả OCR
//                        foreach ($results['responses'] as $i => $result) {
//                            $currentFile = $fileMap[$i];
//                            $detectedText = $result['fullTextAnnotation']['text'] ?? null;
//
//                            if ($detectedText) {
//                                $prefix  = SaleConfig::PRODUCT_KEY_TYPES[$request->type] ?? '';
//                                $pattern = '/'.preg_quote($prefix, '/').'\d+/';
//
//                                preg_match($pattern, $detectedText, $matches);
//                                $code = $matches[0] ?? null;
//
//                                if ($code) {
//                                    $product = Product::where([
//                                        'type' => $request->type,
//                                        'name' => $code,
//                                    ])->first();
//
//                                    if ($product) {
//                                        $product->update([
//                                            'path_thumb' => $currentFile['stored']
//                                        ]);
//                                        $products[] = $code;
//                                    } else {
//                                        // Detect code nhưng không có product
//                                        $failedFiles[] = [
//                                            'file'   => $currentFile['original'],
//                                            'path'   => $currentFile['stored'],
//                                            'reason' => "Không tìm thấy product với code: {$code}"
//                                        ];
//                                    }
//                                } else {
//                                    // Có text nhưng không ra code
//                                    $failedFiles[] = [
//                                        'file'   => $currentFile['original'],
//                                        'path'   => $currentFile['stored'],
//                                        'reason' => "Không detect được mã code trong text"
//                                    ];
//                                }
//                            } else {
//                                // Không detect ra text
//                                $failedFiles[] = [
//                                    'file'   => $currentFile['original'],
//                                    'path'   => $currentFile['stored'],
//                                    'reason' => "Không detect được text"
//                                ];
//                            }
//                        }
//
//                        // Reset batch
//                        $requests = [];
//                        $fileMap  = [];
//                    }
//                }
//
//                // Redirect kèm danh sách failedFiles (flash session)
//                return redirect()
//                    ->route('sale.availableProducts', [
//                        'search' => implode(',', $products)
//                    ])
//                    ->with('failedFiles', $failedFiles);
//            }
//        }
//
//        return view('sale.product-upload');
//    }
//
//
    private function processBatch($client, $requests, $fileMap, $type, &$products, &$failedFiles)
    {
        $response = $client->post(
            'https://vision.googleapis.com/v1/images:annotate?key=' . env('GOOGLE_VISION_KEY'),
            ['json' => ["requests" => $requests]]
        );
        $results = json_decode($response->getBody(), true);

        foreach ($results['responses'] as $i => $result) {
            $currentFile  = $fileMap[$i];
            $detectedText = $result['fullTextAnnotation']['text'] ?? null;

            if ($detectedText) {
                $prefix  = SaleConfig::PRODUCT_KEY_TYPES[$type] ?? '';
                $pattern = '/'.preg_quote($prefix, '/').'\d+/';

                preg_match($pattern, $detectedText, $matches);
                $code = $matches[0] ?? null;

                if ($code) {
                    $product = Product::where([
                        'type' => $type,
                        'name' => $code,
                    ])->first();

                    if ($product) {
                        $product->update([
                            'path_thumb' => $currentFile['stored']
                        ]);
                        $products[] = $code;
                    } else {
                        $failedFiles[] = [
                            'file'   => $currentFile['original'],
                            'path'   => $currentFile['stored'],
                            'reason' => "Không tìm thấy product với code: {$code}"
                        ];
                    }
                } else {
                    $failedFiles[] = [
                        'file'   => $currentFile['original'],
                        'path'   => $currentFile['stored'],
                        'reason' => "Không detect được mã code trong text"
                    ];
                }
            } else {
                $failedFiles[] = [
                    'file'   => $currentFile['original'],
                    'path'   => $currentFile['stored'],
                    'reason' => "Không detect được text"
                ];
            }
        }
    }

    public function productUpload(Request $request)
    {
        if ($request->isMethod('POST')) {
            if ($files = $request->file('banner')) {
                $client      = new Client(['timeout' => 120]);
                $requests    = [];
                $fileMap     = [];
                $products    = [];
                $failedFiles = [];

                $remainingFiles = count($files);

                foreach ($files as $file) {
                    $remainingFiles--;

                    $originalName = $file->getClientOriginalName();
                    $storedPath   = str_replace('public', 'storage', $file->store('public/sale'));

                    $requests[] = [
                        "image" => [
                            "source" => ["imageUri" => asset($storedPath)]
                        ],
                        "features" => [["type" => "TEXT_DETECTION"]],
                    ];

                    $fileMap[] = [
                        'stored'   => $storedPath,
                        'original' => $originalName,
                    ];

                    // gửi batch mỗi 16 file
                    if (count($fileMap) == 16 || $remainingFiles == 0) {
                        $this->processBatch($client, $requests, $fileMap, $request->type, $products, $failedFiles);
                        $requests = [];
                        $fileMap  = [];
                    }
                }

                /**
                 * Retry: gom các file fail, thử OCR lại 1 lần nữa
                 */
                if (!empty($failedFiles)) {
                    $retryRequests = [];
                    $retryFileMap  = [];
                    $retryFails    = []; // giữ fail sau khi retry

                    foreach ($failedFiles as $failed) {
                        $retryRequests[] = [
                            "image" => [
                                "source" => ["imageUri" => asset($failed['path'])]
                            ],
                            "features" => [["type" => "TEXT_DETECTION"]],
                        ];

                        $retryFileMap[] = [
                            'stored'   => $failed['path'],
                            'original' => $failed['file'],
                        ];
                    }

                    // clear failedFiles để ghi lại kết quả retry
                    $failedFiles = [];

                    // gọi OCR retry cho tất cả fail
                    $this->processBatch($client, $retryRequests, $retryFileMap, $request->type, $products, $failedFiles);
                }



                return redirect()
                    ->route('sale.availableProducts', [
                        'search' => implode(',', $products)
                    ])
                    ->with('failedFiles', $failedFiles);
            }
        }

        return view('sale.product-upload');
    }

    /**
     * Batch upload với Queue - Trang chính
     */
    public function productBatchUpload(Request $request)
    {
        if ($request->isMethod('GET')) {
            return view('sale.product-batch-upload');
        }
    }

    /**
     * Upload từng ảnh qua AJAX
     */
    public function productBatchUploadSingle(Request $request)
    {
        try {
            \Log::info('Batch upload single attempt', [
                'has_file' => $request->hasFile('image'),
                'file_error' => $request->hasFile('image') ? $request->file('image')->getError() : 'no file',
                'all_data' => $request->except('image'),
            ]);
            $request->validate([
                'image' => 'required', 
                'type' => 'required|string',
                'session_id' => 'required|string',
            ]);

            $file = $request->file('image');
            $originalName = $file->getClientOriginalName();
            $storedPath = str_replace('public', 'storage', $file->store('public/sale'));

            // Tạo record log PENDING trong DB
            ProductUploadLog::create([
                'session_id' => $request->session_id,
                'file_name' => $originalName,
                'file_path' => $storedPath,
                'status' => 'PENDING',
            ]);

            // Dispatch job vào queue
            \App\Jobs\ProcessProductImageUpload::dispatch(
                $storedPath,
                $originalName,
                $request->type,
                $request->session_id
            );

            return response()->json([
                'success' => true,
                'message' => 'File uploaded and queued for processing',
                'file' => $originalName,
                'path' => $storedPath,
            ]);

        } catch (\Illuminate\Validation\ValidationException $e) {
            return response()->json([
                'success' => false,
                'message' => 'Validation failed',
                'errors' => $e->errors(),
            ], 422);
        } catch (\Exception $e) {
            return response()->json([
                'success' => false,
                'message' => $e->getMessage(),
            ], 400);
        }
    }

    /**
     * Kiểm tra trạng thái upload session
     */
    public function productBatchUploadStatus(Request $request)
    {
        $sessionId = $request->input('session_id');
        
        if (!$sessionId) {
            return response()->json(['error' => 'session_id required'], 400);
        }

        // Truy vấn trực tiếp từ database
        $logs = ProductUploadLog::where('session_id', $sessionId)->get();

        $success = $logs->where('status', 'SUCCESS')->count();
        $failed = $logs->whereIn('status', ['FAILED', 'ERROR'])->count();
        $pending = $logs->where('status', 'PENDING')->count();

        return response()->json([
            'session_id' => $sessionId,
            'success' => (int)$success,
            'failed' => (int)$failed,
            'processing' => (int)$pending,
            'total' => (int)$logs->count(),
            'logs' => $logs->map(function($item) {
                return [
                    'file' => $item->file_name,
                    'status' => $item->status,
                    'code' => $item->product_code,
                    'reason' => $item->message,
                    'text' => $item->detected_text,
                ];
            })->toArray(),
        ]);
    }


    public function productUpdate(Request $request)
    {

        $product = Product::find($request->id);
        if ($file = $request->file('banner')) {
            $request = $request->merge([
                'path_thumb' => str_replace('public', 'storage', $file->store('public/sale'))
            ]);

        }

        if ($request->type != $product->type) {
            $latestId = (int)filter_var((Product::where('type', $request->type)->latest()->first()->name ?? 0), FILTER_SANITIZE_NUMBER_INT) + 1;
            $request = $request->merge([
                'name' => \App\Contracts\SaleConfig::PRODUCT_KEY_TYPES[$request->type] . $latestId,
            ]);
        }

        Size::where('product_id', $product->id)->delete();
        foreach ($request->sizes as $size) {
            if ($size['value']) {
                Size::create([
                    'product_id' => $product->id,
                    'value' => $size['value'],
                    'name' => $size['name'],
                ]);
            }

        }
        $product->update($request->all());
        if ($request->backRoute) {
            return redirect($request->backRoute);

        }
        return redirect()->back();

    }


    public function orders(Request $request)
    {
        $orders = Order::with('products')->select('orders.*')->distinct()->orderBy('orders.id', 'DESC');
        if ($request->status) {
            $orders = Order::with('products')->select('orders.*')->distinct()->orderBy('STATUS', 'ASC')->orderBy('id', 'DESC');

            $orders->whereIn('status', [
                'NEW',
                'ON HOLD',
                'A SHIP NOW',
            ]);
        }

        if ($search = $request->search) {
            $search = strtoupper($search);
            $orders->join('order_products', 'orders.id', 'order_products.order_id')
                ->join('products', 'products.id', 'order_products.product_id');
            $orders->where(function ($query) use ($search) {
                $query->where(DB::raw('upper(orders.name)'), 'like', "%{$search}%");
                $query->orWhere('phone', 'like', "%{$search}%");
                $query->orWhere(DB::raw('upper(products.name)'), 'like', "%{$search}%");
            });
        }
        $orders = $orders->select(['orders.*'])->distinct('orders.id')->paginate(20);

        $orderNew = Order::whereIn('status', [
            'NEW',
            'ON HOLD',
        ])->get();
        $totalOrderNew = $orderNew->count();
        $totalOrderNewPrice = $orderNew->sum('final_amount');
        $currentOrders = Order::with('products')->whereBetween('created_at', [
            Carbon::now()->startOfMonth(),
            Carbon::now()->endOfMonth()
        ])->get();

        $currentTotalFinalPrice = $currentOrders->sum('final_amount');
        $currentTotalOrder = $currentOrders->count();


        return view('sale.orders', compact(
            'orders',
            'orderNew',
            'totalOrderNew',
            'totalOrderNewPrice',
            'currentTotalFinalPrice',
            'currentTotalOrder'
        ));
    }

    public function consumers(Request $request)
    {
        $consumers = Consumer::orderBy('id', 'DESC')->with('orders')->get();

        return view('sale.consumers', compact(
            'consumers'

        ));
    }

    public function summaryOrders(Request $request)
    {
        $orders = Order::with('products')->orderBy('orders.id', 'DESC');
        if ($request->status) {
            $orders = Order::with('products')->orderBy('STATUS', 'ASC')->orderBy('orders.id', 'DESC');

            $orders->whereIn('status', [
                'NEW',
                'ON HOLD',
                'A SHIP NOW',

            ]);
        }

        if ($search = $request->search) {
            $search = strtoupper($search);
            $orders->join('order_products', 'orders.id', 'order_products.order_id')
                ->join('products', 'products.id', 'order_products.product_id');
            $orders->where(function ($query) use ($search) {
                $query->where(DB::raw('upper(orders.name)'), 'like', "%{$search}%");
                $query->orWhere('phone', 'like', "%{$search}%");
                $query->orWhere(DB::raw('upper(products.name)'), 'like', "%{$search}%");
            });
        }

        $orders = $orders->select(['orders.*'])->distinct('orders.id')->get();
            $orderNew = Order::whereIn('status', [
            'NEW',
            'ON HOLD',
            'A SHIP NOW',
        ])->get();
        $totalOrderNew = $orderNew->count();
        $totalOrderNewPrice = $orderNew->sum('final_amount');
        $currentOrders = Order::with('products')->whereBetween('created_at', [
            Carbon::now()->startOfMonth(),
            Carbon::now()->endOfMonth()
        ])->get();

        $currentTotalFinalPrice = $currentOrders->sum('final_amount');
        $currentTotalOrder = $currentOrders->count();

        return view('sale.orders-summary', compact(
            'orders',
            'orderNew',
            'totalOrderNew',
            'totalOrderNewPrice',
            'currentTotalFinalPrice',
            'currentTotalOrder'

        ));
    }

    public function orderDetail($id)
    {

        $order = Order::with('products')->where('id', $id)->first();
        $products = Product::where('status', 'AVAILABLE')->with('sizes')->orderBy('id', 'DESC')->get();

        return view('sale.order-detail', compact('order', 'products'));
    }

    public function orderDelete($id)
    {
        Order::find($id)->delete();
        OrderProduct::where(['order_id' => $id])->delete();

        $this->updateStatusProduct();

        return redirect()->back();
    }
    function parsePrice($priceStr) {
        if (!$priceStr) return 0;

        // Bỏ hết ký tự không phải số hoặc dấu chấm
        $number = preg_replace('/[^0-9.]/', '', $priceStr);

        if ($number === '') return 0;

        return (int) $number * 1000; // vì "k" = 1000
    }

    public function orderCreate(Request $request)
    {

        if ($request->isMethod('POST') || $request->showDetail) {

            $products = Product::whereIn('id', $request->products)->get();
            $vision = new GoogleVisionService();
            $total_amount = 0;

            $defaultPrice = [
                'TROUSERS' => 180,
            ];
            foreach ($products as $product) {
                $ocrText = null;
                $ocrPrice = null;

                if ($product->path_thumb) {
                    $ocrText = $vision->extractTextFromImage($product->path_thumb);
                    $ocrPrice = $ocrText ? $vision->extractPrice($ocrText) : null;
                    info("{$product->name} - {$ocrText} - {$ocrPrice}");
                }

                // Nếu OCR không ra giá thì fallback về DB price
                $price = $ocrPrice ?: array_get($defaultPrice, $product->type) ?: 0;
                $price = $this->parsePrice($price);

                $total_amount += $price;
            }

            info("total: ", [$total_amount]);

            $total_amount += intval($request->additional_amount);
            $total_amount -= intval($request->freeship_amount);

            $request = $request->merge([

                'status' => 'NEW',
                'total_amount' => $total_amount,
                'is_freeship' => $request->is_freeship == 'on' ? true : false,
                'is_paid' => $request->is_paid == 'on' ? true : false,
                'final_amount' => $total_amount,
                'phone' => str_replace(' ', '', $request->phone),
            ]);

            if (!$request->consumer_id) {
                $consumer = Consumer::where([
                    'name' => $request->name,
                    'phone' => $request->phone,
                    'address' => $request->address,
                ])->first();
                if (!$consumer) {
                    $consumer = Consumer::create([
                        'name' => $request->name,
                        'phone' => $request->phone,
                        'address' => $request->address,
                    ]);
                }

                $request = $request->merge([
                    'consumer_id' => $consumer->id,
                ]);
            }
            $order = Order::create($request->all());


            foreach ($products as $product) {
                OrderProduct::create([
                    'product_id' => $product->id,
                    'order_id' => $order->id,

                ]);


            }
            $this->updateStatusProduct();


            if ($request->showDetail) {
                return redirect()->route('sale.order.detail', $order->id);
            }
            if ($request->backRoute) {
                return redirect($request->backRoute);

            }

            return redirect()->back();

        }


        $products = Product::where('status', 'AVAILABLE')->with('sizes')->orderBy('id', 'DESC')->get();

        $consumer = null;
        if ($consumerId = $request->consumer_id) {
            $consumer = Consumer::find($consumerId);
        }
        return view('sale.order-create', compact('products', 'consumer'));
    }


    public function orderUpdate(Request $request)
    {

        $order = Order::find($request->id);
        $products = Product::whereIn('id', $request->products)->get();


        $total_amount = array_sum($products->pluck('price')->toArray()) + intval($request->additional_amount) - intval($request->freeship_amount);

        $request = $request->merge([

            'total_amount' => $total_amount,
            'is_freeship' => $request->is_freeship == 'on' ? true : false,
            'is_paid' => $request->is_paid == 'on' ? true : false,
            'final_amount' => $total_amount,
        ]);

        OrderProduct::where('order_id', $order->id)->delete();

        foreach ($products as $product) {
            OrderProduct::create([
                'product_id' => $product->id,
                'order_id' => $order->id,

            ]);


        }
        $order->update($request->all());
        $this->updateStatusProduct();

        if ($request->backRoute) {
            return redirect($request->backRoute);

        }
        return redirect()->back();

    }

    protected function updateStatusProduct()
    {
        Product::where('status', '<>', 'ON_HOLD')->update([
            'status' => 'AVAILABLE'
        ]);

        Product::whereIn('id', OrderProduct::all()->pluck('product_id')->toArray())->update([
            'status' => 'SOLD'
        ]);
    }

    public function orderSmartCreate(Request $request)
    {
        if ($request->isMethod('POST')) {
            $text = $request->input('description');

            $prompt = <<<PROMPT
Phân tích đoạn text sau và xuất ra JSON với các trường:

- "name": tên người nhận
- "phone": số điện thoại
- "address": địa chỉ giao hàng (bao gồm cả phần trong ngoặc nếu có)
- "products": danh sách mã sản phẩm (ví dụ: ts1169, ts1158...), dạng mảng chuỗi
- "note": nội dung ghi chú sau `note:`, nếu có. Nếu có nhiều cụm thì ngăn cách bằng dấu phẩy.

Yêu cầu:
- Trả về **duy nhất JSON**, không thêm chữ giải thích, không xuống dòng ngoài JSON.
- Nếu không có `note:` thì để `"note": ""`

Đoạn text:
"""$text"""
PROMPT;


            $client = new Client();

            $res = $client->post('https://api.openai.com/v1/chat/completions', [
                'headers' => [
                    'Authorization' => 'Bearer ' . env('OPENAI_API_KEY'),
                    'Content-Type'  => 'application/json',
                ],
                'json' => [
                    'model' => 'gpt-4o-2024-11-20',
                    'messages' => [
                        ['role' => 'system', 'content' => 'Bạn là một trợ lý phân tích dữ liệu giao hàng.'],
                        ['role' => 'user', 'content' => $prompt],
                    ],
                    'temperature' => 0.3,
                ]
            ]);

            $response = json_decode((string)$res->getBody(), true);
            $result = $response['choices'][0]['message']['content'] ?? '';

            // Cố gắng parse JSON nếu OpenAI trả đúng
            $parsedJson = json_decode($result, true);



            $parsedJson['phone'] = preg_replace('/\D/', '', $parsedJson['phone']);
            $consumer = Consumer::where('phone', $parsedJson['phone'])->first();
            $products = array_map(
                function ($product) {
                    return Product::where('name', 'like', '%' . $product . '%')->where('status', 'AVAILABLE')->first()->id ?? null;
                },
                $parsedJson['products']
            );
            return $this->orderCreate(new Request([
                'consumer_id' => $consumer->id ?? 0,
                'name' => $parsedJson['name'],
                'phone' => $parsedJson['phone'],
                'address' => $parsedJson['address'],
                'note' => $parsedJson['note'],
                'products' => array_filter($products),
                'showDetail' => true
            ]));
        }

        return view('sale.smart-order-create');
    }
}
