<?php

namespace App\Services;

//
use GuzzleHttp\Client;

class GoogleVisionService
{
    public function extractTextFromImage(string $imageUrl): ?string
    {

        $client = new Client([
            'timeout' => 120, // giây
            'connect_timeout' => 10,
        ]);
        $response = $client->post('https://vision.googleapis.com/v1/images:annotate?key=AIzaSyBCQwssnHoUNbueD6DGggilZou-Sd0NjU0', [
            'json' => [
                "requests" => [
                    [
                        "image" => [
                            "source" => ["imageUri" => asset($imageUrl)]
                        ],
                        "features" => [
                            ["type" => "TEXT_DETECTION"]
                        ]
                    ]
                ]
            ]
        ]);

        $result = json_decode($response->getBody(), true);

        return $result['responses'][0]['fullTextAnnotation']['text'] ?? null;
    }
    function extractPrice($text)
    {
        // Regex: số + optional khoảng trắng + chữ "k"
        preg_match('/\d+\s*k/i', $text, $matches);

        if (!empty($matches)) {
            // Trả về giá tiền (ví dụ "320k")
            return trim($matches[0]);
        }

        return null; // không có giá
    }


}
