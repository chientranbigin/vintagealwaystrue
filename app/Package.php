<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class Package extends Model
{
    protected $fillable = [
        'name',
        'thumb',
    ];

    /**
     * Relationship with mover
     *
     * @return mixed
     */
    public function details()
    {
        return $this->hasMany(DetailPackage::class, 'package_id');
    }
}
