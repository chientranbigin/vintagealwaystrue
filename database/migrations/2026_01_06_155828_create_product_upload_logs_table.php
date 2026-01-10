<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateProductUploadLogsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('product_upload_logs', function (Blueprint $table) {
            $table->increments('id');
            $table->string('session_id')->index();
            $table->string('file_name');
            $table->string('file_path');
            $table->string('product_code')->nullable()->index();
            $table->string('status')->default('PENDING'); // PENDING, SUCCESS, FAILED
            $table->text('message')->nullable();
            $table->text('detected_text')->nullable();
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
        Schema::dropIfExists('product_upload_logs');
    }
}
