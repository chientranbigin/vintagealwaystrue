<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class AlbumType extends Model
{
    protected $fillable = [
        'name',
        'desc',
    ];
}
