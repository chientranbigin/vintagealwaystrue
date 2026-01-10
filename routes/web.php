<?php

/*
|--------------------------------------------------------------------------
| Web Routes
|--------------------------------------------------------------------------
|
| Here is where you can register web routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| contains the "web" middleware group. Now create something great!
|
*/
use Illuminate\Routing\Router;


Route::get('/', 'WebController@home')->name('home');
Route::get('/album', 'WebController@album')->name('album');
Route::get('/film', 'WebController@film')->name('film');
Route::get('/album/{id}', 'WebController@albumDetail')->name('albumDetail');
Route::get('/package/{id}', 'WebController@packageDetail')->name('packageDetail');
Route::get('/about', 'WebController@about')->name('about');


Route::get('/products', 'CustomerController@availableProducts')->name('customer.products');


Route::group([
    'prefix' => 'admin'
], function (Router $router) {
    Route::get('/', 'AdminController@home')->name('admin.home');
    Route::get('/albums', 'AdminController@albums')->name('admin.albums');
    Route::get('/albums/{id}', 'AdminController@albumDetail')->name('admin.albums.detail');
    Route::get('/albums/delete/{id}', 'AdminController@albumDelete')->name('admin.albums.detail.delete');
    Route::get('/albums-create', 'AdminController@albumCreate')->name('admin.albums.create');
    Route::post('/albums-create', 'AdminController@albumCreate')->name('admin.albums.create.post');
    Route::post('/albums/{id}', 'AdminController@albumUpdate')->name('admin.albums.detail');
    Route::post('/albums-type', 'AdminController@albumTypeUpdate')->name('admin.albumType.update');
    Route::get('/packages', 'AdminController@packages')->name('admin.packages');
    Route::get('/packages/{id}', 'AdminController@packageDetail')->name('admin.packages.detail');
    Route::post('/packages/{id}', 'AdminController@packageUpdate')->name('admin.packages.post');
    Route::post('/packages-detail/{id}', 'AdminController@packageDetailUpdate')->name('admin.packages.detail.post');
    Route::get('/packages-detail/delete/{id}', 'AdminController@packageDetailDelete')->name('admin.packages.detail.delete');
    Route::get('/packages/{id}/packages-create', 'AdminController@packageCreate')->name('admin.packages.create');
    Route::post('/packages/{id}/packages-create', 'AdminController@packageCreate')->name('admin.packages.create.post');
    Route::post('elements', 'AdminController@storeElement')->name('admin.elements.store');

});

Route::group([
    'prefix' => 'sale'
], function (Router $router) {
    Route::get('/', 'SaleController@home')->name('sale.home');
    Route::get('/product', 'SaleController@products')->name('sale.product');
    Route::get('/available-product', 'SaleController@availableProducts')->name('sale.availableProducts');
    Route::post('/available-product', 'SaleController@searchAvailableProducts')->name('sale.searchAvailableProducts');
    Route::get('/product/{id}', 'SaleController@productDetail')->name('sale.product.detail');
    Route::get('/product/delete/{id}', 'SaleController@productDelete')->name('sale.product.detail.delete');
    Route::get('/product-create', 'SaleController@productCreate')->name('sale.product.create');
    Route::post('/product-create', 'SaleController@productCreate')->name('sale.product.create.post');
    Route::get('/product-upload', 'SaleController@productUpload')->name('sale.product.upload');
    Route::post('/product-upload', 'SaleController@productUpload')->name('sale.product.upload.post');
    
    // Batch upload with queue
    Route::get('/product-batch-upload', 'SaleController@productBatchUpload')->name('sale.product.batch.upload');
    Route::post('/product-batch-upload/single', 'SaleController@productBatchUploadSingle')->name('sale.product.batch.upload.single');
    Route::get('/product-batch-upload/status', 'SaleController@productBatchUploadStatus')->name('sale.product.batch.upload.status');
    
    Route::post('/product/{id}', 'SaleController@productUpdate')->name('sale.product.detail');

    Route::get('/order', 'SaleController@orders')->name('sale.order');
    Route::post('/order', 'SaleController@searchOrders')->name('sale.searchOrder');
    Route::get('/consumer', 'SaleController@consumers')->name('sale.consumer');
    Route::get('/order-summary', 'SaleController@summaryOrders')->name('sale.summaryOrders');
    Route::get('/order/{id}', 'SaleController@orderDetail')->name('sale.order.detail');
    Route::get('/order/delete/{id}', 'SaleController@orderDelete')->name('sale.order.detail.delete');
    Route::get('/order-create', 'SaleController@orderCreate')->name('sale.order.create');
    Route::post('/order-create', 'SaleController@orderCreate')->name('sale.order.create.post');
    Route::post('/order/{id}', 'SaleController@orderUpdate')->name('sale.order.detail');
    Route::post('/smart-order-create', 'SaleController@orderSmartCreate')->name('sale.order.smart-create.post');
    Route::get('/smart-order-create', 'SaleController@orderSmartCreate')->name('sale.order.smart-create.get');


    Route::get('/complete/{id}', 'SaleController@orderComplete')->name('sale.order.complete');
    Route::get('/ship/{id}', 'SaleController@orderShip')->name('sale.order.ship');
    Route::get('/hold/{id}', 'SaleController@orderHold')->name('sale.order.hold');


});
