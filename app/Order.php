<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class Order extends Model
{
    protected $fillable = [
        'name',
        'phone',
        'address',
        'note',
        'desc',
        'additional_amount',
        'freeship_amount',
        'total_amount',
        'final_amount',
        'deposit_amount',
    
        'is_freeship',
        'is_paid',
        'status',
        'consumer_id',
        'expenses'
    ];

    protected $casts = [
        'expenses' => 'array'
    ];

    /**
     * Relationship with mover
     *
     * @return mixed
     */
    public function products()
    {
        return $this->belongsToMany(Product::class, 'order_products');
    }
}
