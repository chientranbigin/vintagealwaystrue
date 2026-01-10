<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class Album extends Model
{
    protected $fillable = [
        'album_type_id',
        'name',
        'thumb',
        'type',
    ];


    /**
     * Relationship with mover
     *
     * @return mixed
     */
    public function albumType()
    {
        return $this->belongsTo(AlbumType::class);
    }


    /**
     * Relationship with mover
     *
     * @return mixed
     */
    public function items()
    {
        return $this->hasMany(AlbumItem::class);
    }
}
