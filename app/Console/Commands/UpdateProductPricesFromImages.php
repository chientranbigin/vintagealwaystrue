<?php

namespace App\Console\Commands;

use App\Product;
use App\Services\GoogleVisionService;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class UpdateProductPricesFromImages extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'product:update-prices {--force : Update prices even if they are already set}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Update prices for all available products with images using OCR or default mapping';

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function handle()
    {
        $query = Product::where('status', 'AVAILABLE')
            ->whereNotNull('path_thumb');

        if (!$this->option('force')) {
            $query->where(function($q) {
                $q->whereNull('price')->orWhere('price', 0);
            });
        }

        $products = $query->get();
        $total = $products->count();

        if ($total === 0) {
            $this->info('No products found to update.');
            return;
        }

        $this->info("Found {$total} products to process.");
        $bar = $this->output->createProgressBar($total);
        $bar->start();

        $vision = new GoogleVisionService();
        $defaultPriceMapping = [
            'TROUSERS' => 180,
        ];

        $updatedCount = 0;
        $failedCount = 0;

        foreach ($products as $product) {
            try {
                $detectedText = $vision->extractTextFromImage($product->path_thumb);
                $ocrPrice = null;
                
                if ($detectedText) {
                    $ocrPrice = $vision->extractPrice($detectedText);
                }

                if ($ocrPrice) {
                    $price = $this->parsePrice($ocrPrice);
                    $source = 'OCR';
                } else {
                    $defaultPrice = $defaultPriceMapping[$product->type] ?? 0;
                    $price = $this->parsePrice($defaultPrice);
                    $source = 'Default Mapping';
                }

                $product->update(['price' => $price]);
                $updatedCount++;
                
                // Log progress for large batches
                if ($total > 50) {
                    Log::info("Command: Updated product #{$product->id} price to {$price} via {$source}");
                }

            } catch (\Exception $e) {
                $this->error("\nError processing product #{$product->id}: {$e->getMessage()}");
                $failedCount++;
            }

            $bar->advance();
        }

        $bar->finish();
        $this->info("\n\nFinished processing.");
        $this->info("Successfully updated: {$updatedCount}");
        $this->info("Failed: {$failedCount}");
    }

    /**
     * Parse giá từ chuỗi (ví dụ "320k" -> 320000)
     */
    private function parsePrice($priceStr) {
        if (!$priceStr) return 0;

        // Bỏ hết ký tự không phải số hoặc dấu chấm
        $number = preg_replace('/[^0-9.]/', '', $priceStr);

        if ($number === '') return 0;

        return (int) $number * 1000; // vì "k" = 1000
    }
}
