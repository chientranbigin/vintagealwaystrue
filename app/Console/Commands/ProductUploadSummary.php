<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;

class ProductUploadSummary extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $description = 'Generate summary report of product upload session';
    
    protected $signature = 'products:upload-summary {session_id? : The Session ID to report}';

    public function handle()
    {
        $sessionId = $this->argument('session_id');
        $logPath = storage_path('logs/product_upload.log');

        if (!file_exists($logPath)) {
            $this->error('Log file not found.');
            return;
        }

        $logs = file($logPath, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES);
        
        // Filter by session ID if provided, otherwise check all or latest?
        // Let's filter by session if provided, or grab the latest session ID found in logs.
        if (!$sessionId) {
            // Find latest session
            $lastLog = end($logs);
             if (preg_match('/"session_id":"([^"]+)"/', $lastLog, $matches)) {
                $sessionId = $matches[1];
                $this->info("No session ID provided. Using latest session: $sessionId");
             } else {
                 $this->error("Could not determine latest session ID.");
                 return;
             }
        }

        $sessionLogs = array_filter($logs, function($line) use ($sessionId) {
            return strpos($line, $sessionId) !== false;
        });

        $success = [];
        $failed = [];
        $codes = [];
        
        foreach ($sessionLogs as $line) {
             // Look for JSON part
             if (preg_match('/(\{.*\})$/', $line, $matches)) {
                $jsonStr = $matches[1];
                $data = json_decode($jsonStr, true);
                
                // Fallback if data is null (json error)
                if (!$data) continue;

                if (strpos($line, 'SUCCESS') !== false) {
                    $success[] = $data;
                    if (isset($data['code'])) {
                        $codes[] = $data['code'];
                    }
                } elseif (strpos($line, 'FAILED') !== false) {
                    $failed[] = $data;
                }
             }
        }
        
        $this->info("=========================================");
        $this->info("REPORT SESSION: $sessionId");
        $this->info("=========================================");
        
        $this->info("SUCCESS: " . count($success));
        $headers = ['Code', 'Original File', 'Image Path'];
        $rows = [];
        foreach ($success as $item) {
            $rows[] = [$item['code'], $item['file'], $item['path']];
        }
        $this->table($headers, $rows);

        $this->info("FAILED: " . count($failed));
        $headersFail = ['Original File', 'Reason', 'Detected Text'];
        $rowsFail = [];
        foreach ($failed as $item) {
             $text = isset($item['text']) ? substr(str_replace("\n", " ", $item['text']), 0, 50) . '...' : 'N/A';
             $rowsFail[] = [$item['file'], $item['reason'], $text];
        }
        $this->table($headersFail, $rowsFail);

        // Check duplicates
        $vals = array_count_values($codes);
        $duplicates = array_filter($vals, function($count) { return $count > 1; });
        
        if (!empty($duplicates)) {
            $this->error("WARNING: DUPLICATE CODES DETECTED IN THIS SESSION");
            foreach ($duplicates as $code => $count) {
                 $this->error("- Code $code appears $count times");
            }
        } else {
            $this->info("No duplicate codes detected in this session.");
        }

        // Trigger clean up
        if ($this->confirm('Do you want to scan and move unused images to backup now?', true)) {
             $this->call('products:clean-images');
        }
    }
}
