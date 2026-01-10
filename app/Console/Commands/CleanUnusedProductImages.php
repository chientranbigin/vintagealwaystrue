<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;

class CleanUnusedProductImages extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'products:clean-images';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Delete product images in storage/sale that are not attached to any product';

    public function handle()
    {
        $this->info('Scanning for unused product images...');

        // 1. Get all used image paths from DB
        // 1. Get all used image paths from DB
        // Filter: Keep images if Product is AVAILABLE/ON_HOLD 
        // OR if SOLD but the associated Order was created within the last 2 months.
        
        $usedPaths = \Illuminate\Support\Facades\DB::table('products')
            ->leftJoin('order_products', 'products.id', '=', 'order_products.product_id')
            ->leftJoin('orders', 'order_products.order_id', '=', 'orders.id')
            ->whereNotNull('products.path_thumb')
            ->where(function($query) {
                $query->whereIn('products.status', ['AVAILABLE', 'ON_HOLD'])
                      ->orWhere(function($q) {
                          $q->where('products.status', 'SOLD')
                            ->where('orders.created_at', '>=', \Carbon\Carbon::now()->subMonths(1));
                      });
            })
            ->distinct()
            ->pluck('products.path_thumb')
            ->map(function($path) {
                return basename($path);
            })->toArray();

        $this->info('Found ' . count($usedPaths) . ' images linked to products.');

        // 2. Scan directory
        $directory = 'public/sale'; // Storage path relative to 'app' for 'storage' disk? or 'public' disk?
        // In previous steps: $file->store('public/sale') -> this is on 'local' disk usually or 'public' disk?
        // Default filesystem usually 'local'. 'public/sale' means storage/app/public/sale.
        
        $files = \Illuminate\Support\Facades\Storage::files('public/sale');
        $movedCount = 0;
        $reclaimedSize = 0;
        
        // Create backup directory
        $backupDir = 'public/sale_backup/' . date('Y-m-d_H-i-s');
        \Illuminate\Support\Facades\Storage::makeDirectory($backupDir);
        $this->info("Backup directory created: $backupDir");

        $bar = $this->output->createProgressBar(count($files));
        $bar->start();

        foreach ($files as $file) {
            $filename = basename($file);

            if (!in_array($filename, $usedPaths)) {
                // Double check: if it's a recent file (e.g. < 1 hour), skip it to avoid race conditions with uploads in progress
                $lastModified = \Illuminate\Support\Facades\Storage::lastModified($file);
                if (time() - $lastModified < 3600) {
                     $bar->advance();
                     continue;
                }

                $size = \Illuminate\Support\Facades\Storage::size($file);
                
                // Move instead of delete
                $newPath = $backupDir . '/' . $filename;
                \Illuminate\Support\Facades\Storage::move($file, $newPath);
                
                $movedCount++;
                $reclaimedSize += $size;
            }
            $bar->advance();
        }

        $bar->finish();
        $this->line('');
        $this->info("Clean up completed.");
        $this->info("Moved: $movedCount files to backup");
        $this->info("Reclaimed main storage: " . number_format($reclaimedSize / 1024 / 1024, 2) . " MB");
    }
}
