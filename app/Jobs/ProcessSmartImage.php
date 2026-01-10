<?php

namespace App\Jobs;

use App\Product;
use App\ProductImage;
use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use Illuminate\Support\Facades\Storage;
use Google\Cloud\Vision\V1\ImageAnnotatorClient;

class ProcessSmartImage implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    protected $mainImagePath;
    protected $detailPaths;
    protected $type;
    protected $sessionId;
    protected $originalName;

    public function __construct($mainImagePath, $detailPaths, $type, $sessionId, $originalName)
    {
        $this->mainImagePath = $mainImagePath;
        $this->detailPaths = $detailPaths;
        $this->type = $type;
        $this->sessionId = $sessionId;
        $this->originalName = $originalName;
    }

    public function handle()
    {
        // 1. Init Log (Step 1 start)
        \App\ProductUploadLog::create([
            'session_id' => $this->sessionId,
            'file_name' => $this->originalName,
            'file_path' => $this->mainImagePath,
            'status' => 'PROCESSING',
            'message' => 'Step 1: Recognition & Product Search...'
        ]);

        try {
            // STEP 1: RECOGNITION & FIND PRODUCT
            $vision = new \App\Services\GoogleVisionService();
            $detectedText = '';
            try {
                $detectedText = $vision->extractTextFromImage($this->mainImagePath);
            } catch (\Exception $e) {
                throw new \Exception("Không thể đọc được chữ từ ảnh (Google Vision Error): " . $e->getMessage());
            }
            \Illuminate\Support\Facades\Log::channel('product_upload')->info("Detected Text: " . $detectedText);

            
            // A. Extract Code
            $prefix = \App\Contracts\SaleConfig::PRODUCT_KEY_TYPES[$this->type] ?? '';
            $pattern = '/' . preg_quote($prefix, '/') . '\d+/i';
            
            preg_match($pattern, $detectedText, $matches);
            $candidateName = isset($matches[0]) ? strtoupper($matches[0]) : null;

            
    

            if (!$candidateName) {
                throw new \Exception("Không detect được mã sản phẩm ($prefix...) từ ảnh hoặc tên file.");
            }

            // B. Extract Price
            $price = 0;
            if ($detectedText && preg_match('/(\d+)[kK]/', $detectedText, $pMatches)) {
                $price = $pMatches[1] * 1000;
            }

            // Apply Default Pricing Mapping if price is 0
            if ($price == 0) {
                $categoryMap = [
                    'TROUSERS' => 180000,
                    'BLAZER'   => 280000,
                    'GILE'     => 120000
                ];
                $price = $categoryMap[$this->type] ?? 0;
            }

            // C. Search Existing Product
            $existing = \App\Product::where('type', $this->type)
                ->where('name', $candidateName)
                ->first();
            
            if (!$existing) {
                throw new \Exception("Sản phẩm mã ($candidateName) không tồn tại. Vui lòng tạo sp trước.");
            }

            // Update record: Price only
            if ($price > 0) {
                $existing->update(['price' => $price]);
            }

            // Update Log with detected info
            \App\ProductUploadLog::where('session_id', $this->sessionId)
                ->where('file_path', $this->mainImagePath)
                ->update([
                    'product_code' => $candidateName,
                    'product_id' => $existing->id,
                    'amount' => $price,
                    'detected_text' => $detectedText
                ]);

            // STEP 2: ATTACHMENT
            \App\ProductUploadLog::where('session_id', $this->sessionId)
                ->where('file_path', $this->mainImagePath)
                ->update(['message' => 'Step 2: Attaching images...']);

            // Attach Main Image as a detail
            $this->processDetailImage($existing, $this->mainImagePath);

            // Attach Detail Images
            foreach ($this->detailPaths as $path) {
                $this->processDetailImage($existing, $path);
            }

            // Log Success
            \App\ProductUploadLog::where('session_id', $this->sessionId)
                ->where('file_path', $this->mainImagePath)
                ->update([
                    'status' => 'SUCCESS',
                    'message' => 'Xử lý thành công (2 bước)'
                ]);

        } catch (\Exception $e) {
             \App\ProductUploadLog::where('session_id', $this->sessionId)
                ->where('file_path', $this->mainImagePath)
                ->update([
                    'status' => 'ERROR',
                    'product_code' => $candidateName ?? 'UNKNOWN',
                    'message' => $e->getMessage()
                ]);
            
            \Illuminate\Support\Facades\Log::channel('product_upload')->error($e);
        }
    }

    protected function processDetailImage($product, $relativePath)
    {
        $fullPath = storage_path('app/public/' . str_replace('storage/', '', $relativePath));
        \Illuminate\Support\Facades\Log::channel('product_upload')->info($fullPath);

        if (!file_exists($fullPath)) return;

        // 1. Load Image using GD
        list($origWidth, $origHeight, $type) = getimagesize($fullPath);
        $src = null;
        switch ($type) {
            case IMAGETYPE_JPEG: $src = imagecreatefromjpeg($fullPath); break;
            case IMAGETYPE_PNG:  $src = imagecreatefrompng($fullPath);  break;
            case IMAGETYPE_GIF:  $src = imagecreatefromgif($fullPath);  break;
            case IMAGETYPE_WEBP: $src = imagecreatefromwebp($fullPath); break;
        }
        if (!$src) return;

        // 2. Resize Logic (Constraint: Max width 1024)
        $maxWidth = 1024;
        $newWidth = $origWidth;
        $newHeight = $origHeight;

        if ($origWidth > $maxWidth) {
            $ratio = $maxWidth / $origWidth;
            $newWidth = $maxWidth;
            $newHeight = round($origHeight * $ratio);
        }

        $dst = imagecreatetruecolor($newWidth, $newHeight);
        
        // Handle transparency for PNG/GIF/WEBP
        if ($type == IMAGETYPE_PNG || $type == IMAGETYPE_GIF || $type == IMAGETYPE_WEBP) {
            imagealphablending($dst, false);
            imagesavealpha($dst, true);
            $transparent = imagecolorallocatealpha($dst, 255, 255, 255, 127);
            imagefilledrectangle($dst, 0, 0, $newWidth, $newHeight, $transparent);
        }

        imagecopyresampled($dst, $src, 0, 0, 0, 0, $newWidth, $newHeight, $origWidth, $origHeight);

        // 3. Save Image (Fallback to JPEG if WebP not supported)
        if (function_exists('imagewebp')) {
            $newPath = preg_replace('/\.[^.]+$/', '.webp', $relativePath);
            $targetFullPath = storage_path('app/public/' . $newPath);
            imagewebp($dst, $targetFullPath, 60);
        } else {
            $newPath = preg_replace('/\.[^.]+$/', '.jpg', $relativePath);
            $targetFullPath = storage_path('app/public/' . $newPath);
            imagejpeg($dst, $targetFullPath, 75);
        }

        // 4. Cleanup
        imagedestroy($src);
        imagedestroy($dst);

        if ($newPath !== $relativePath) {
            @unlink($fullPath);
        }

        // Create Database Record (Store relative path, controller handles asset() mapping)
        ProductImage::create([
            'product_id' => $product->id,
            'file_path' => $newPath
        ]);
    }
}
