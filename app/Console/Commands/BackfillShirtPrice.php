<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use App\Product;
use Carbon\Carbon;

class BackfillShirtPrice extends Command
{
    protected $signature = 'products:backfill-shirt-price';
    protected $description = 'Set price to 150000 for SHIRT products with price 0: AVAILABLE with an image, or SOLD within the last 30 days';

    public function handle()
    {
        $availableQuery = Product::where('type', 'SHIRT')
            ->where('status', 'AVAILABLE')
            ->whereNotNull('path_thumb')
            ->where(function ($q) {
                $q->where('price', 0)->orWhereNull('price');
            });

        $soldQuery = Product::where('type', 'SHIRT')
            ->where('status', 'SOLD')
            ->where(function ($q) {
                $q->where('price', 0)->orWhereNull('price');
            })
            ->whereHas('orders', function ($q) {
                $q->where('orders.created_at', '>=', Carbon::now()->subDays(30));
            });

        $total = 0;

        foreach (['AVAILABLE' => $availableQuery, 'SOLD (30 ngày gần nhất)' => $soldQuery] as $label => $query) {
            $count = $query->count();

            if ($count === 0) {
                $this->info("No SHIRT products to backfill for: $label");
                continue;
            }

            $names = $query->pluck('name')->toArray();
            $query->update(['price' => 150000]);

            $this->info("Updated $count SHIRT products ($label) to price 150000:");
            foreach ($names as $name) {
                $this->line("  - $name");
            }
            $total += $count;
        }

        $this->info("Total updated: $total");

        return 0;
    }
}
