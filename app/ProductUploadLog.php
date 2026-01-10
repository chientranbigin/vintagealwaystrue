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
        'amount',
        'status',
        'message',
        'detected_text',
    ];
}
