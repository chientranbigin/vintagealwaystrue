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

        // SOLD products where order is NOT COMPLETED
        $soldThumbsToKeep = \App\Product::where('status', 'SOLD')
            ->where(function($q) {
                $q->whereDoesntHave('orders')
                  ->orWhereHas('orders', function($sq) {
                      $sq->where('status', '!=', 'COMPLETED');
                  });
            })
            ->pluck('path_thumb')
            ->filter()
            ->toArray();

        // Support for Detail Images
        $activeDetailImages = \App\ProductImage::whereHas('product', function($q) {
            $q->whereIn('status', ['AVAILABLE', 'ON_HOLD'])
              ->orWhere(function($sq) {
                  $sq->where('status', 'SOLD')
                     ->where(function($ssq) {
                         $ssq->whereDoesntHave('orders')
                             ->orWhereHas('orders', function($sssq) {
                                 $sssq->where('status', '!=', 'COMPLETED');
                             });
                     });
              });
        })->pluck('file_path')->toArray();

        $imagesToKeep = array_unique(array_merge($activeProductThumbs, $soldThumbsToKeep, $activeDetailImages));
        
        // Normalize paths (remove leading slash, 'storage/', 'public/' if any)
        $normalize = function($path) {
            $path = ltrim($path, '/');
            if (strpos($path, 'storage/') === 0) {
                $path = substr($path, 8);
            }
            if (strpos($path, 'public/') === 0) {
                $path = substr($path, 7);
            }
            return ltrim($path, '/');
        };

        $imagesToKeep = array_map($normalize, $imagesToKeep);
        
        $this->info('Finding files in storage/app/public/sale...');
        $allFiles = \Storage::disk('public')->allFiles('sale');
        
        $deletedCount = 0;
        $skippedRecent = 0;
        $now = time();
        $safeThreshold = 48 * 3600; // 48 hours for extra safety

        foreach ($allFiles as $file) {
            // Normalize current file path to check against DB
            $normalizedFile = $normalize($file);
            
            if (!in_array($normalizedFile, $imagesToKeep)) {
                // AGE CHECK: Never delete recent files to avoid race conditions with uploads/jobs
                $lastModified = \Storage::disk('public')->lastModified($file);
                if ($now - $lastModified < $safeThreshold) {
                    $skippedRecent++;
                    continue;
                }

                $this->info("Deleting orphaned: $file");
                \Storage::disk('public')->delete($file);
                $deletedCount++;
            }
        }

        $this->info("Cleanup finished.");
        $this->info("- Deleted: $deletedCount files");
        $this->info("- Skipped (Recent): $skippedRecent files");
    }
}
