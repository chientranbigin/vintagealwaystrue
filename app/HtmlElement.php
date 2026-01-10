<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class HtmlElement extends Model
{
    protected $fillable = [
        'group',
        'name',
        'type',
        'value',
    ];
}
