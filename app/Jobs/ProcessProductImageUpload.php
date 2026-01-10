<?php

namespace App\Jobs;

use App\Product;
use App\Services\GoogleVisionService;
use App\ProductUploadLog;
use App\Contracts\SaleConfig;
use Illuminate\Bus\Queueable;
use Illuminate\Queue\SerializesModels;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Storage;

class ProcessProductImageUpload implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    protected $imagePath;
    protected $originalName;
    protected $type;
    protected $uploadSessionId;

    /**
     * Số lần retry nếu job fail
     */
    public $tries = 2;

    /**
     * Timeout cho job (seconds)
     */
    public $timeout = 120;

    /**
     * Create a new job instance.
     */
    public function __construct($imagePath, $originalName, $type, $uploadSessionId)
    {
        $this->imagePath = $imagePath;
        $this->originalName = $originalName;
        $this->type = $type;
        $this->uploadSessionId = $uploadSessionId;
    }

    /**
     * Execute the job.
     */
    public function handle()
    {
        try {
            $vision = new GoogleVisionService();
            
            // Extract text từ ảnh
            $detectedText = $vision->extractTextFromImage($this->imagePath);

            if (!$detectedText) {
                $this->logFailed('Không detect được text từ ảnh');
                return;
            }

            // Extract mã sản phẩm
            $prefix = SaleConfig::PRODUCT_KEY_TYPES[$this->type] ?? '';
            $pattern = '/' . preg_quote($prefix, '/') . '\d+/';
            
            preg_match($pattern, $detectedText, $matches);
            $code = $matches[0] ?? null;

            if (!$code) {
                $this->logFailed('Không detect được mã code trong text', $detectedText);
                return;
            }

            // Tìm product
            $product = Product::where([
                'type' => $this->type,
                'name' => $code,
            ])->first();

            if (!$product) {
                $this->logFailed("Không tìm thấy product với code: {$code}", $detectedText, $code);
                return;
            }

    
            // Update ảnh cho product
            $product->update([
                'path_thumb' => $this->imagePath
            ]);

            $this->logSuccess($code, $detectedText);

        } catch (\Exception $e) {
            $this->logFailed('Exception: ' . $e->getMessage());
            throw $e; // Re-throw để Laravel retry
        }
    }

    /**
     * Log khi thành công
     */
    protected function logSuccess($code, $detectedText)
    {
        // Update database
        ProductUploadLog::where([
            'session_id' => $this->uploadSessionId,
            'file_name' => $this->originalName,
        ])->update([
            'status' => 'SUCCESS',
            'product_code' => $code,
            'detected_text' => $detectedText,
            'message' => 'Xử lý thành công'
        ]);

        Log::channel('product_upload')->info("✓ SUCCESS", [
            'session_id' => $this->uploadSessionId,
            'file' => $this->originalName,
            'code' => $code,
            'path' => $this->imagePath,
            'text' => substr($detectedText, 0, 100),
        ]);
    }

    /**
     * Log khi thất bại
     */
    protected function logFailed($reason, $detectedText = null, $code = null)
    {
        // Update database
        ProductUploadLog::where([
            'session_id' => $this->uploadSessionId,
            'file_name' => $this->originalName,
        ])->update([
            'status' => 'FAILED',
            'message' => $reason,
            'product_code' => $code,
            'detected_text' => $detectedText
        ]);

        Log::channel('product_upload')->warning("✗ FAILED", [
            'session_id' => $this->uploadSessionId,
            'file' => $this->originalName,
            'path' => $this->imagePath,
            'reason' => $reason,
            'text' => $detectedText ? substr($detectedText, 0, 100) : null,
        ]);
    }

    /**
     * Handle a job failure.
     */
    public function failed(\Throwable $exception)
    {
        // Update database
        ProductUploadLog::where([
            'session_id' => $this->uploadSessionId,
            'file_name' => $this->originalName,
        ])->update([
            'status' => 'ERROR',
            'message' => 'Job Failed: ' . $exception->getMessage()
        ]);

        Log::channel('product_upload')->error("✗ JOB FAILED", [
            'session_id' => $this->uploadSessionId,
            'file' => $this->originalName,
            'path' => $this->imagePath,
            'exception' => $exception->getMessage(),
            'trace' => $exception->getTraceAsString(),
        ]);
    }
}
