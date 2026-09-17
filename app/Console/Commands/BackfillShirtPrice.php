<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use App\Product;

class BackfillShirtPrice extends Command
{
    protected $signature = 'products:backfill-shirt-price';
    protected $description = 'Set price to 150000 for AVAILABLE SHIRT products that have an image but price 0';

    public function handle()
    {
        $query = Product::where('type', 'SHIRT')
            ->where('status', 'AVAILABLE')
            ->whereNotNull('path_thumb')
            ->where(function ($q) {
                $q->where('price', 0)->orWhereNull('price');
            });

        $count = $query->count();

        if ($count === 0) {
            $this->info('No SHIRT products to backfill.');
            return 0;
        }

        $names = $query->pluck('name')->toArray();
        $query->update(['price' => 150000]);

        $this->info("Updated $count SHIRT products to price 150000:");
        foreach ($names as $name) {
            $this->line("  - $name");
        }

        return 0;
    }
}
