<?php

namespace App\Console\Commands;

use App\Product;
use App\Services\ImageService;
use Illuminate\Console\Command;

class GenerateProductThumbs extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'products:generate-thumbs {--force : Force regenerate all thumbnails} {--limit= : Limit the number of products to process}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Generate scaled thumbnails for products';

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
        $query = Product::whereNotNull('path_thumb')->where('path_thumb', '!=', '');
        
        if (!$this->option('force')) {
            $query->where(function($q) {
                $q->whereNull('image_thumb_scale')->orWhere('image_thumb_scale', '');
            });
        }

        if ($limit = $this->option('limit')) {
            $query->limit($limit);
        }

        $products = $query->get();
        $total = $products->count();

        if ($total === 0) {
            $this->info('No products need thumbnail generation.');
            return;
        }

        $this->info("Generating thumbnails for $total products...");
        $bar = $this->output->createProgressBar($total);
        $bar->start();

        $imageService = new ImageService();
        $success = 0;
        $failed = 0;

        foreach ($products as $product) {
            try {
                $scaledPath = $imageService->createScaledThumbnail($product->path_thumb);
                if ($scaledPath) {
                    $product->update(['image_thumb_scale' => $scaledPath]);
                    $success++;
                } else {
                    $failed++;
                    $this->error("Failed to scale: {$product->path_thumb}");
                }
            } catch (\Exception $e) {
                $failed++;
                $this->error("Error for ID {$product->id}: " . $e->getMessage());
            }
            $bar->advance();
        }

        $bar->finish();
        $this->info('');
        $this->info("Done! Success: $success, Failed: $failed");
    }
}
