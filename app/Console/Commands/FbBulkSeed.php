<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use App\Product;
use App\Jobs\SyncFacebookPost;

class FbBulkSeed extends Command
{
    protected $signature = 'fb:bulk-seed {--type= : Only this product type} {--dry-run : List what would be posted without dispatching}';
    protected $description = 'Queue a Facebook post for every AVAILABLE product with a thumb that has not been posted yet';

    public function handle()
    {
        $query = Product::where('status', 'AVAILABLE')
            ->whereNotNull('path_thumb')
            ->whereNull('fb_post_id');

        if ($type = $this->option('type')) {
            $query->where('type', strtoupper($type));
        }

        $products = $query->orderBy('id')->get();

        if ($products->isEmpty()) {
            $this->info('Nothing to post — all AVAILABLE products already have fb_post_id (or none match).');
            return 0;
        }

        $this->info("Found {$products->count()} product(s) to post.");

        if ($this->option('dry-run')) {
            foreach ($products as $p) {
                $this->line("  {$p->name} ({$p->type})");
            }
            return 0;
        }

        $bar = $this->output->createProgressBar($products->count());
        $delaySeconds = 0;
        foreach ($products as $p) {
            SyncFacebookPost::dispatch($p, 'create')->delay(now()->addSeconds($delaySeconds));
            $delaySeconds += 8; // spread out to avoid FB API rate limiting
            $bar->advance();
        }
        $bar->finish();
        $this->line('');
        $totalMinutes = round($delaySeconds / 60, 1);
        $this->info("Queued {$products->count()} jobs, spaced 8s apart (~{$totalMinutes} min to finish). Processed by the queue worker.");

        return 0;
    }
}
