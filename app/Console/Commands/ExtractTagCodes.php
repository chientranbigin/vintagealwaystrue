<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use App\Contracts\SaleConfig;
use App\Product;

class ExtractTagCodes extends Command
{
    protected $signature = 'tags:extract {type}';
    protected $description = 'Extract product codes from tag images and bulk hold others';

    protected $apiKey = 'AIzaSyBCQwssnHoUNbueD6DGggilZou-Sd0NjU0';
    protected $cacheFile;

    public function handle()
    {
        $type = strtoupper($this->argument('type'));
        $prefix = SaleConfig::PRODUCT_KEY_TYPES[$type] ?? null;

        if (!$prefix) {
            $this->error("Unknown type: $type. Available: " . implode(', ', array_keys(SaleConfig::PRODUCT_KEY_TYPES)));
            return 1;
        }

        $dir = storage_path('app/public/sale/tags');
        if (!is_dir($dir)) {
            $this->error("Directory not found: $dir");
            return 1;
        }

        $this->cacheFile = storage_path('app/tag_vision_cache.json');
        $cache = file_exists($this->cacheFile) ? json_decode(file_get_contents($this->cacheFile), true) : [];

        $files = array_filter(scandir($dir), fn($f) => preg_match('/\.(jpe?g|png|heic)$/i', $f));
        $files = array_values($files);

        $codes = [];
        $failed = [];

        foreach ($files as $filename) {
            $filepath = "$dir/$filename";

            if (isset($cache[$filename])) {
                $text = $cache[$filename];
                $this->line("<comment>[cached]</comment>      $filename");
            } else {
                $this->line("<comment>[vision] </comment>     $filename → calling API...");
                $text = $this->callVision($filepath);
                $cache[$filename] = $text ?? '';
                file_put_contents($this->cacheFile, json_encode($cache, JSON_PRETTY_PRINT));
            }

            if ($text && preg_match('/\b(\d{4})\b/', $text, $m)) {
                $code = $prefix . $m[1];
                $codes[] = $code;
                $this->line("  <info>→ $code</info>");
            } else {
                $failed[] = $filename;
                $this->line("  <fg=red>✗ no 4-digit number found</fg=red>");
            }
        }

        $uniqueCodes = array_values(array_unique($codes));
        sort($uniqueCodes);

        $this->newLine();
        $this->info('✓ Extracted ' . count($uniqueCodes) . ' unique codes from ' . count($files) . ' images');

        if ($failed) {
            $this->warn('✗ Failed (' . count($failed) . '): ' . implode(', ', $failed));
        }

        $this->line('Codes: ' . json_encode($uniqueCodes));
        $this->newLine();

        // Auto run hold-except
        $this->call('products:hold-except', [
            'type'  => $type,
            'codes' => implode(',', $uniqueCodes),
        ]);

        return 0;
    }

    protected function callVision(string $filepath): ?string
    {
        $content = base64_encode(file_get_contents($filepath));

        $client = new \GuzzleHttp\Client(['timeout' => 30]);
        try {
            $response = $client->post(
                "https://vision.googleapis.com/v1/images:annotate?key={$this->apiKey}",
                [
                    'json' => [
                        'requests' => [[
                            'image'    => ['content' => $content],
                            'features' => [['type' => 'TEXT_DETECTION']],
                        ]],
                    ],
                ]
            );
            $data = json_decode($response->getBody(), true);
            return $data['responses'][0]['fullTextAnnotation']['text'] ?? null;
        } catch (\Exception $e) {
            $this->error("Vision API error: " . $e->getMessage());
            return null;
        }
    }
}
