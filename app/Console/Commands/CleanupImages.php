<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Facades\DB;

class CleanupImages extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'images:cleanup';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Cleanup product images from storage that are no longer needed.';

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
        $this->info('Starting image cleanup...');

        // 1. Get all images to keep
        $activeProductThumbs = \App\Product::where('status', 'AVAILABLE')
            ->orWhere('status', 'ON_HOLD')
            ->pluck('path_thumb')
            ->filter()
            ->toArray();

        // SOLD products in last 2 months
        $recentlySoldThumbs = \App\Product::where('status', 'SOLD')
            ->where('updated_at', '>=', now()->subMonths(2))
            ->pluck('path_thumb')
            ->filter()
            ->toArray();

        // Support for Detail Images
        $activeDetailImages = \App\ProductImage::whereHas('product', function($q) {
            $q->where('status', 'AVAILABLE')
              ->orWhere('status', 'ON_HOLD')
              ->orWhere(function($sq) {
                  $sq->where('status', 'SOLD')
                     ->where('updated_at', '>=', now()->subMonths(2));
              });
        })->pluck('file_path')->toArray();

        $imagesToKeep = array_unique(array_merge($activeProductThumbs, $recentlySoldThumbs, $activeDetailImages));
        
        // Normalize paths (remove leading slash if any, ensuring it matches storage directory structure)
        $imagesToKeep = array_map(function($path) {
            return str_replace('storage/', '', $path);
        }, $imagesToKeep);

        $this->info('Finding files in storage/app/public/sale...');
        $allFiles = \Storage::disk('public')->allFiles('sale');
        
        $deletedCount = 0;
        foreach ($allFiles as $file) {
            if (!in_array($file, $imagesToKeep)) {
                $this->info("Deleting: $file");
                \Storage::disk('public')->delete($file);
                $deletedCount++;
            }
        }

        $this->info("Cleanup finished. Deleted $deletedCount files.");
    }
}
