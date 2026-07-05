<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use App\Product;

class BulkHoldExcept extends Command
{
    protected $signature = 'products:hold-except {type} {codes}';
    protected $description = 'Set ON_HOLD all AVAILABLE products of a type NOT in the given codes list';

    public function handle()
    {
        $type  = strtoupper($this->argument('type'));
        $codes = array_filter(array_map('trim', explode(',', $this->argument('codes'))));

        if (empty($codes)) {
            $this->error('No codes provided.');
            return 1;
        }

        $query = Product::where('type', $type)
            ->where('status', 'AVAILABLE')
            ->whereNotIn('name', $codes);

        $count = $query->count();

        if ($count === 0) {
            $this->info('No products to hold.');
            return 0;
        }

        $names = $query->pluck('name')->toArray();
        $query->update(['status' => 'ON_HOLD']);

        $this->info("Updated $count products to ON_HOLD:");
        foreach ($names as $name) {
            $this->line("  - $name");
        }

        return 0;
    }
}
