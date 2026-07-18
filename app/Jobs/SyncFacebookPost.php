<?php

namespace App\Jobs;

use App\Product;
use App\Services\FacebookService;
use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use Illuminate\Support\Facades\Log;

class SyncFacebookPost implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    public int $tries = 3;
    public int $backoff = 60;

    protected $product;
    protected $action; // create | sold | revert

    public function __construct(Product $product, string $action)
    {
        $this->product = $product;
        $this->action  = $action;
    }

    public function handle(FacebookService $fb)
    {
        if (!config('services.facebook.page_token')) {
            Log::channel('product_upload')->info("SyncFacebookPost: no FB token configured, skipping.");
            return;
        }

        try {
            if ($this->action === 'create') {
                $postId = $fb->createPost($this->product);
                if ($postId) {
                    $this->product->update(['fb_post_id' => $postId]);
                    Log::channel('product_upload')->info("FB post created: {$postId} for {$this->product->name}");
                }
            } elseif ($this->action === 'sold') {
                if (!$this->product->fb_post_id) return;
                $fb->markSold($this->product->fb_post_id, $this->product);
                Log::channel('product_upload')->info("FB post marked SOLD: {$this->product->fb_post_id}");
            } elseif ($this->action === 'revert') {
                if (!$this->product->fb_post_id) return;
                $fb->revertAvailable($this->product->fb_post_id, $this->product);
                Log::channel('product_upload')->info("FB post reverted AVAILABLE: {$this->product->fb_post_id}");
            }
        } catch (\Exception $e) {
            Log::channel('product_upload')->error("SyncFacebookPost [{$this->action}] error: " . $e->getMessage());
            throw $e;
        }
    }
}
