<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use App\Product;
use App\Services\FacebookService;

class FbTestPost extends Command
{
    protected $signature = 'fb:test-post {product? : Product name/code, e.g. TS1235} {--type= : Product type to pick a random AVAILABLE one if no product given}';
    protected $description = 'Create a test Facebook post for a product (seeding/dev only)';

    public function handle()
    {
        $query = Product::with(['sizes', 'images']);

        if ($name = $this->argument('product')) {
            $product = $query->where('name', $name)->first();
            if (!$product) {
                $this->error("Product not found: $name");
                return 1;
            }
        } else {
            $type = $this->option('type');
            if ($type) {
                $query->where('type', strtoupper($type));
            }
            $product = $query->where('status', 'AVAILABLE')
                ->whereNotNull('path_thumb')
                ->inRandomOrder()
                ->first();

            if (!$product) {
                $this->error('No AVAILABLE product with thumb found.');
                return 1;
            }
        }

        $this->info("Product: {$product->name} (id={$product->id}, type={$product->type}, status={$product->status})");

        $fb = new FacebookService();
        $postId = $fb->createPost($product);

        if (!$postId) {
            $this->error('Failed to create post (no images or API error).');
            return 1;
        }

        $product->update(['fb_post_id' => $postId]);

        $this->info("✓ Post created: $postId");
        $this->line("Message preview:");
        $this->line($fb->buildMessage($product, 'AVAILABLE'));

        return 0;
    }
}
