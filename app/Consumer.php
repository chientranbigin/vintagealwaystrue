<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class Consumer extends Model
{
    protected $fillable = [
        'name',
        'address',
        'phone',
    ];

    /**
     * Relationship with mover
     *
     * @return mixed
     */
    public function orders()
    {
        return $this->hasMany(Order::class);
    }
}
