<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateTableInit extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('html_elements', function (Blueprint $table) {
            $table->increments('id');
            $table->string('group');
            $table->string('name');
            $table->string('type');
            $table->text('value');
            $table->timestamps();
        });


        Schema::create('album_types', function (Blueprint $table) {
            $table->increments('id');
            $table->string('name');
            $table->timestamps();
        });

        Schema::create('albums', function (Blueprint $table) {
            $table->increments('id');
            $table->integer('album_type_id');
            $table->string('name');
            $table->string('thumb');
            $table->timestamps();
        });

        Schema::create('album_items', function (Blueprint $table) {
            $table->increments('id');
            $table->integer('album_id');
            $table->string('type');
            $table->string('path');
            $table->timestamps();
        });


        Schema::create('packages', function (Blueprint $table) {
            $table->increments('id');
            $table->string('name');
            $table->string('thumb');
            $table->timestamps();
        });

        Schema::create('detail_packages', function (Blueprint $table) {
            $table->increments('id');
            $table->integer('package_id');
            $table->string('name');
            $table->float('price');
            $table->text('content');
            $table->timestamps();
        });

    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('html_elements');
        Schema::dropIfExists('album_types');
        Schema::dropIfExists('albums');
        Schema::dropIfExists('album_items');
        Schema::dropIfExists('packages');
        Schema::dropIfExists('detail_packages');
    }
}
