<?php

namespace App\Contracts;

interface SaleConfig
{
    const DOMAIN_CUSTOMER = 'vintagealwaystrue.id.vn';
    const PRODUCT_TYPES = [
        'TROUSERS',
//        'HARRINGTON_JACKET',
        'JACKET',
        'SHIRT',
        'BLAZER',
        'TIE',
        'GILE',
        'BELT',
        'POLO SHIRT',
        'HAT',
        'SUIT',
    ];

    const PRODUCT_KEY_TYPES = [
        'TROUSERS' => 'TS',
        'HARRINGTON_JACKET' => 'HJ',
        'JACKET' => 'JK',
        'SHIRT' => 'ST',
        'BLAZER' => 'BZ',
        'TIE' => 'TI',
        'BELT' => 'BE',
        'GILE' => 'GL',
        'POLO SHIRT' => 'PS',
        'HAT' => "HA",
        'SUIT' => 'SU',

    ];

    const SIZES = [
        'VAI',
        'NGỰC',
        'EO',
        'DÀI ÁO',
        'DÀI ÁO SAU',
        'DÀI TAY',
        'EO QUẦN',
        'ĐÁY',
        'ĐÙI',
        'DÀI QUẦN',
        'ỐNG',
        'DƯ LAI',

    ];

    const MAPPING_SIZES = [
        'VAI' => 'V',
        'NGỰC' => 'N',
        'EO' => 'E',
        'DÀI ÁO' => 'D',
        'DÀI ÁO SAU' => 'DS',
        'DÀI TAY' => 'DT',
        'EO QUẦN' => 'E',
        'ĐÁY' => 'ĐA',
        'ĐÙI' => 'Đ',
        'DÀI QUẦN' => 'D',
        'ỐNG' => 'Ô',
        'DƯ LAI' => 'L',

    ];

    const MAPPING_NAME_SIZE = [
        'DÀI ÁO' => 'dài trước',
        'DÀI ÁO SAU' => 'DÀI SAU',
        'EO QUẦN' => 'eo',
        'DÀI QUẦN' => 'dài',
        'VAI' => 'vai',
        'NGỰC' => 'ngực',
        'EO' => 'eo',
        'DÀI TAY' => 'dài tay',

        'ĐÁY' => 'đáy',
        'ĐÙI' => 'đùi',
        'ỐNG' => 'ống',
        'DƯ LAI' => 'dư gấu',
    ];

    const ORDER_STATUSES = [
        'NEW',
        'A SHIP NOW',

        'ON HOLD',
        'READY FOR SHIPPING',
        'COMPLETED',
    ];

    const MAPPING_HASH_TAG = [
        'TROUSERS' => '#vintageshoptrousers',
        'JACKET' => '#vintageshopjacket',
        'SHIRT' => '#vintageshopshirt',
        'BLAZER' => '#vintageshopblazer',
        'TIE' => '#vintageshoptie',
        'BELT' => '#vintageshopbelt',
        'GILE' => 'vintageshopgile',

    ];



}
