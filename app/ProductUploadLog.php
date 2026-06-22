<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class ProductUploadLog extends Model
{
    protected $fillable = [
        'session_id',
        'file_name',
        'file_path',
        'product_code',
        'product_id',
        'amount',
        'status',
        'message',
        'detected_text',
    ];

    public function product()
    {
        return $this->belongsTo(Product::class);
    }
}
