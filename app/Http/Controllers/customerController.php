<?php

namespace App\Http\Controllers;

use App\Album;
use App\AlbumType;
use App\Contracts\ElementConfig;
use App\DetailPackage;
use App\HtmlElement;
use App\Package;
use App\Product;
use Illuminate\Foundation\Bus\DispatchesJobs;
use Illuminate\Http\Request;
use Illuminate\Routing\Controller as BaseController;
use Illuminate\Foundation\Validation\ValidatesRequests;
use Illuminate\Foundation\Auth\Access\AuthorizesRequests;
use Illuminate\Support\Facades\DB;

class CustomerController extends Controller
{

    public function availableProducts(Request $request)
    {



        $products = Product::with('sizes')->where('status', 'AVAILABLE')->whereNotNull('path_thumb')->select('products.*')->distinct()->orderBy('created_at', 'DESC');

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

//        if ($request->status) {
//            $products->where('status', $request->status);
//        }

        if ($request->from_search_key || $request->to_search_key ){

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

        if ( $request->from_search_key_2 || $request->to_search_key_2){

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
//            )
//        );
        $products = $products->paginate(60);


        return view('consumer.products', compact('products'));
    }


}
