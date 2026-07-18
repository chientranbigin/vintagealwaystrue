<?php

namespace App\Services;

use App\Product;
use GuzzleHttp\Client;

class FacebookService
{
    protected Client $http;
    protected string $pageId;
    protected string $token;

    // Trousers: size → eo range [min, max) (eo thực = DB × 2)
    const TROUSER_SIZE_MAP = [
        28 => [74, 76],
        29 => [76, 78],
        30 => [78, 80],
        31 => [80, 82],
        32 => [82, 84],
        33 => [84, 86],
        34 => [86, 88],
        35 => [88, 90],
        36 => [90, 92],
        37 => [92, 94],
        38 => [94, 96],
        39 => [96, 98],
        40 => [98, 100],
    ];

    public function __construct()
    {
        $this->pageId = config('services.facebook.page_id');
        $this->token  = config('services.facebook.page_token');
        $this->http   = new Client(['base_uri' => 'https://graph.facebook.com/v19.0/', 'timeout' => 30]);
    }

    public function createPost(Product $product): ?string
    {
        $product->load(['images', 'sizes']);

        // Collect image paths (thumb first, then detail images)
        $imagePaths = [];
        if ($product->path_thumb) {
            $imagePaths[] = $this->absolutePath($product->path_thumb);
        }
        foreach ($product->images as $img) {
            $p = $this->absolutePath($img->file_path);
            if ($p && !in_array($p, $imagePaths)) {
                $imagePaths[] = $p;
            }
        }

        if (empty($imagePaths)) return null;

        // Upload each image as unpublished photo
        $mediaIds = [];
        foreach ($imagePaths as $path) {
            if (!file_exists($path)) continue;
            $res  = $this->http->post("{$this->pageId}/photos", [
                'multipart' => [
                    ['name' => 'source', 'contents' => fopen($path, 'r')],
                    ['name' => 'published', 'contents' => 'false'],
                    ['name' => 'access_token', 'contents' => $this->token],
                ],
            ]);
            $data = json_decode($res->getBody(), true);
            if (!empty($data['id'])) {
                $mediaIds[] = ['media_fbid' => $data['id']];
            }
        }

        if (empty($mediaIds)) return null;

        // Create post with all photos
        $message = $this->buildMessage($product, 'AVAILABLE');
        $body = [
            'message'       => $message,
            'attached_media' => json_encode($mediaIds),
            'access_token'  => $this->token,
        ];

        $res  = $this->http->post("{$this->pageId}/feed", ['form_params' => $body]);
        $data = json_decode($res->getBody(), true);

        return $data['id'] ?? null;
    }

    public function markSold(string $fbPostId, Product $product): void
    {
        $message = $this->buildMessage($product, 'SOLD');
        $this->http->post($fbPostId, [
            'form_params' => [
                'message'      => $message,
                'access_token' => $this->token,
            ],
        ]);
    }

    public function revertAvailable(string $fbPostId, Product $product): void
    {
        $message = $this->buildMessage($product, 'AVAILABLE');
        $this->http->post($fbPostId, [
            'form_params' => [
                'message'      => $message,
                'access_token' => $this->token,
            ],
        ]);
    }

    public function buildMessage(Product $product, string $status): string
    {
        $product->loadMissing('sizes');

        $statusLine = $status === 'SOLD'
            ? "⛔ ĐÃ BÁN — {$product->name}"
            : "✅ CÒN HÀNG — {$product->name}";

        $typeLabel = $product->type;
        $price     = number_format($product->price, 0, ',', '.') . 'đ';

        // Detected text from latest upload log
        $detectedText = \App\ProductUploadLog::where('product_id', $product->id)
            ->whereNotNull('detected_text')
            ->latest()
            ->value('detected_text');

        $lines = [
            $statusLine,
            "👔 {$typeLabel} | {$price}",
        ];

        if ($detectedText) {
            $lines[] = '';
            $lines[] = trim($detectedText);
        }

        $lines[] = '';
        $lines[] = $this->buildHashtags($product, $status);

        return implode("\n", $lines);
    }

    protected function buildHashtags(Product $product, string $status): string
    {
        $tags = [];

        // Product code
        $tags[] = "#{$product->name}";

        // Type + status
        $statusTag = $status === 'SOLD' ? "{$product->type}_SOLD" : "{$product->type}_AVAILABLE";
        $tags[] = "#{$statusTag}";

        // Measurement tags based on type
        if ($product->type === 'TROUSERS') {
            $eoSize = $product->sizes->firstWhere('name', 'EO');
            if ($eoSize) {
                $eoDb     = (int) $eoSize->value;
                $eoActual = $eoDb * 2;
                $tags[]   = "#EO_{$eoActual}";
                $tags[]   = '#EO_' . $eoActual . '_' . ($eoActual + 2);
                $size     = $this->getTrouserSize($eoActual);
                if ($size) $tags[] = "#SIZE{$size}";
            }
        } elseif (in_array($product->type, ['BLAZER', 'SUIT', 'SHIRT', 'GILE'])) {
            $vaiSize = $product->sizes->firstWhere('name', 'VAI');
            if ($vaiSize) {
                $vai    = (int) $vaiSize->value;
                $tags[] = "#VAI_{$vai}";
                // Range bước 1, chẵn start
                $start  = ($vai % 2 === 0) ? $vai : $vai - 1;
                $tags[] = "#VAI_{$start}_" . ($start + 1);
            }
        }

        return implode(' ', $tags);
    }

    protected function getTrouserSize(int $eoActual): ?int
    {
        foreach (self::TROUSER_SIZE_MAP as $size => [$min, $max]) {
            if ($eoActual >= $min && $eoActual < $max) {
                return $size;
            }
        }
        return null;
    }

    protected function absolutePath(string $relativePath): string
    {
        $path = ltrim($relativePath, '/');
        if (str_starts_with($path, 'storage/')) {
            $path = substr($path, 8);
        }
        return storage_path('app/public/' . $path);
    }
}
