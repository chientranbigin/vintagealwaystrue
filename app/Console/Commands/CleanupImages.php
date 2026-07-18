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

        // 1. Get all images to keep (AVAILABLE only, not ON_HOLD)
        $activeProductThumbs = \App\Product::where('status', 'AVAILABLE')
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
            $q->where('status', 'AVAILABLE')
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
TS3856,TS3896,TS4006,TS3256,TS3852,TS3669,TS3793,TS3988,TS3995,TS4008,TS3915,TS3898,TS3934,TS3718,TS3748,TS3749,TS3818,TS3805,TS3730,TS3720,TS3778,TS3652,TS3632,TS3665,TS3853,TS3812,TS2994,TS3888,TS3473,TS3982,TS4014,TS3599,TS4015,TS3356,TS3371,TS3897,TS2749,TS4004,TS3589,TS2702,TS3729,TS3762,TS3807,TS3782,TS3725,TS3685,TS3855,TS3633,TS3777,TS2591,TS3631,TS3827,TS3734,TS3815,TS3828,TS3472,TS3932,TS3944,TS3929,TS3891,TS3563,TS3927,TS3911,TS3955,TS3926,TS3584,TS3935,TS3874,TS3617,TS3880,TS3719,TS3732,TS3701,TS3671,TS3902,TS4016,TS3808,TS3779,TS4000,TS3990,TS3984,TS3973,TS3173,TS3435,TS1993,TS3281,TS3978,TS3903
