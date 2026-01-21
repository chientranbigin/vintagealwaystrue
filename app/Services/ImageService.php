<?php

namespace App\Services;

class ImageService
{
    /**
     * Scale and compress an image using PHP GD
     * 
     * @param string $sourcePath Path relative to storage/app/public
     * @param int $maxWidth
     * @param int $quality
     * @return string|null Path to the scaled image relative to storage/app/public
     */
    public function createScaledThumbnail($sourcePath, $maxWidth = 400, $quality = 75)
    {
        // Handle absolute URLs by extracting the internal path
        $baseUrl = config('app.url');
        $cleanPath = $sourcePath;
        
        if (strpos($cleanPath, $baseUrl) === 0) {
            $cleanPath = substr($cleanPath, strlen($baseUrl));
        }

        // Clean up common prefixes like '/storage/' or 'storage/'
        $cleanPath = ltrim($cleanPath, '/');
        if (strpos($cleanPath, 'storage/') === 0) {
            $cleanPath = substr($cleanPath, 8);
        }

        $fullSourcePath = storage_path('app/public/' . $cleanPath);
        
        if (!file_exists($fullSourcePath) || is_dir($fullSourcePath)) {
            \Log::warning("ImageService: File not found or is directory: {$fullSourcePath}");
            return null;
        }

        $info = @getimagesize($fullSourcePath);
        if (!$info) {
            \Log::error("ImageService: getimagesize failed for {$fullSourcePath}");
            return null;
        }

        list($width, $height, $type) = $info;
        
        // Calculate new dimensions
        if ($width <= $maxWidth) {
            $newWidth = $width;
            $newHeight = $height;
        } else {
            $newWidth = $maxWidth;
            $newHeight = round($height * ($maxWidth / $width));
        }

        // Load source image
        switch ($type) {
            case IMAGETYPE_JPEG:
                $source = imagecreatefromjpeg($fullSourcePath);
                break;
            case IMAGETYPE_PNG:
                $source = imagecreatefrompng($fullSourcePath);
                break;
            case IMAGETYPE_GIF:
                $source = imagecreatefromgif($fullSourcePath);
                break;
            case IMAGETYPE_WEBP:
                if (function_exists('imagecreatefromwebp')) {
                    $source = imagecreatefromwebp($fullSourcePath);
                } else {
                    return null;
                }
                break;
            default:
                return null;
        }

        if (!$source) {
            \Log::error("ImageService: Failed to create image resource for {$fullSourcePath}. Type: {$type}");
            return null;
        }

        // Create new image
        $thumb = imagecreatetruecolor($newWidth, $newHeight);

        // Handle transparency for PNG/GIF
        if ($type == IMAGETYPE_PNG || $type == IMAGETYPE_GIF) {
            imagealphablending($thumb, false);
            imagesavealpha($thumb, true);
            $transparent = imagecolorallocatealpha($thumb, 255, 255, 255, 127);
            imagefilledrectangle($thumb, 0, 0, $newWidth, $newHeight, $transparent);
        }

        // Resize
        imagecopyresampled($thumb, $source, 0, 0, 0, 0, $newWidth, $newHeight, $width, $height);

        // Generate thumb path
        $pathInfo = pathinfo($cleanPath);
        $thumbName = $pathInfo['filename'] . '_scaled.jpg'; // We convert to JPG for better compression
        $thumbRelativePath = $pathInfo['dirname'] . '/' . $thumbName;
        $fullThumbPath = storage_path('app/public/' . ltrim($thumbRelativePath, '/'));

        // Save as JPEG
        if (imagejpeg($thumb, $fullThumbPath, $quality)) {
            imagedestroy($source);
            imagedestroy($thumb);
            return $thumbRelativePath;
        }

        imagedestroy($source);
        imagedestroy($thumb);
        return null;
    }
}
