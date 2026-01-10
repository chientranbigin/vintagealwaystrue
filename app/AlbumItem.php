<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class AlbumItem extends Model
{
    protected $fillable = [
        'album_id',
        'type',
        'path',
    ];
}
