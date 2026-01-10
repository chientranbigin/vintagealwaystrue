<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class DetailPackage extends Model
{
    protected $fillable = [
        'package_id',
        'name',
        'price',
        'content',
    ];
}
