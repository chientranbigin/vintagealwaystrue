<?php

namespace App\Http\Controllers;

use App\Album;
use App\AlbumType;
use App\Contracts\ElementConfig;
use App\DetailPackage;
use App\HtmlElement;
use App\Package;
use Illuminate\Foundation\Bus\DispatchesJobs;
use Illuminate\Http\Request;
use Illuminate\Routing\Controller as BaseController;
use Illuminate\Foundation\Validation\ValidatesRequests;
use Illuminate\Foundation\Auth\Access\AuthorizesRequests;

class AdminController extends Controller
{
    public function home()
    {
        $elements = HtmlElement::all();

        $elementData = array_combine(
            $elements
                ->map(function ($elements) {
                    return "{$elements->group}.{$elements->name}";
                })
                ->toArray(),
            $elements
                ->map(function ($elements) {
                    return $elements->value;
                })
                ->toArray()
        );

        return view('admin.index', compact('elementData'));
    }

    public function storeElement(Request $request)
    {

        if ($file = $request->file('value')) {
            $value = str_replace('public', 'storage', $file->store('public/elements'));
        } else {
            $value = $request->value;
        }


        HtmlElement::updateOrCreate(
            [
                'group' => $request->group,
                'name' => $request->name,
            ],
            [
                'value' => $value,
                'type' => $request->type,

            ]
        );
        return redirect()->back();
    }

    public function albums()
    {
        $albumTypes = AlbumType::all();
        $albums  = Album::with('albumType')->get();
        return view('admin.albums' , compact('albumTypes', 'albums'));
    }

    public function packages()
    {

        $packages = Package::all();

        return view('admin.packages', compact('packages'));
    }
    public function packageDetail($id)
    {

        $package = Package::with('details')->where('id', $id)->first();

        return view('admin.packages-detail', compact('package'));
    }

    public function packageCreate($id, Request $request)
    {
        if ($request->isMethod('POST')) {
            DetailPackage::create($request->merge([
                'package_id' => $id
            ])->all());

            return redirect()->route('admin.packages.detail', $id);
        }


        return view('admin.packages-detail-create', compact('id'));
    }

    public function packageUpdate($id, Request $request)
    {
        if ($file = $request->file('banner')) {
            $request = $request->merge([
                'thumb' => str_replace('public', 'storage', $file->store('public/package'))
            ]);

        }
        Package::find($id)->update($request->all());

        return redirect()->back();
    }

    public function packageDetailUpdate($id, Request $request)
    {

        DetailPackage::find($id)->update($request->all());

        return redirect()->back();
    }


    public function albumTypeUpdate(Request $request)
    {
        AlbumType::find($request->id)->update([
            'name' => $request->name,
            'desc' => $request->desc,
        ]);
        return redirect()->back();

    }

    public function albumDetail($id)
    {
        $albumTypes = AlbumType::all();
        $album  = Album::with('albumType', 'items')->where('id', $id)->first();

        $video = $album->items->first(function ($item) {
           return $item->type == ElementConfig::TYPE_VIDEO;
        });
        return view('admin.album-detail' , compact('albumTypes', 'album', 'video'));
    }

    public function albumDelete($id)
    {
        Album::find($id)->delete();
        return redirect()->back();
    }
    public function packageDetailDelete($id)
    {
        DetailPackage::find($id)->delete();
        return redirect()->back();
    }
    public function albumCreate(Request $request)
    {
        if ($request->isMethod('POST')) {
            if ($file = $request->file('banner')) {
                $request = $request->merge([
                    'thumb' => str_replace('public', 'storage', $file->store('public/album'))
                ]);

            }

            $album = Album::create($request->all());


            if ($items = $request->file('items')) {
                foreach ($items as $item) {
                    $album->items()->create([
                        'path' => str_replace('public', 'storage', $item->store('public/album')),
                        'type' => ElementConfig::TYPE_IMAGE
                    ]);
                }
            }

            $videoNew = $request->video;

            if ($videoNew) {
                $album->items()->create([
                    'path' => $videoNew,
                    'type' => ElementConfig::TYPE_VIDEO
                ]);
            }

            return redirect()->route('admin.albums');
        }


        $albumTypes = AlbumType::all();
        return view('admin.album-create' , compact('albumTypes'));
    }



    public function albumUpdate(Request $request)
    {

        $album = Album::find($request->id);
        if ($file = $request->file('banner')) {
            $request = $request->merge([
                'thumb' => str_replace('public', 'storage', $file->store('public/album'))
            ]);

        }

        if ($items = $request->file('items')) {
            $album->items()->delete();
            foreach ($items as $item) {
                $album->items()->create([
                    'path' => str_replace('public', 'storage', $item->store('public/album')),
                    'type' => ElementConfig::TYPE_IMAGE
                ]);
            }
        }

        $videoNew = $request->video;

        $video = $album->items->first(function ($item) {
            return $item->type == ElementConfig::TYPE_VIDEO;
        });

        if ($video) {
            $album->items()->where('id', $video->id)->delete();
        }

        if ($videoNew) {
            $album->items()->create([
                'path' => $videoNew,
                'type' => ElementConfig::TYPE_VIDEO
            ]);
        }

        $album->update($request->all());
        return redirect()->back();

    }
}
