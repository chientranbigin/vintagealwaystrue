<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateCategoriesAndAddFbAlbumPhotoId extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('categories')) {
            Schema::create('categories', function (Blueprint $table) {
                $table->increments('id');
                $table->string('code', 50)->unique();
                $table->string('name', 100);
                $table->string('fb_album_id', 100)->nullable();
                $table->boolean('is_active')->default(true);
                $table->timestamps();
            });
        }

        if (Schema::hasTable('products') && !Schema::hasColumn('products', 'fb_album_photo_id')) {
            Schema::table('products', function (Blueprint $table) {
                $table->string('fb_album_photo_id', 100)->nullable()->after('fb_post_id');
            });
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        if (Schema::hasTable('products') && Schema::hasColumn('products', 'fb_album_photo_id')) {
            Schema::table('products', function (Blueprint $table) {
                $table->dropColumn('fb_album_photo_id');
            });
        }

        Schema::dropIfExists('categories');
    }
}
