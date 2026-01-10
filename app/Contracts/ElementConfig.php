<?php
namespace App\Contracts;

interface ElementConfig {
    const TYPE_IMAGE = 'TYPE_IMAGE';
    const TYPE_TEXT = 'TYPE_TEXT';
    const TYPE_VIDEO = 'TYPE_VIDEO';
  const ELEMENTS= [
      'LOGO' => [
          [
              'name' => 'Logo',
              'type' => self::TYPE_IMAGE
          ],
          [
              'name' => 'Icon tab 32',
              'type' => self::TYPE_IMAGE
          ],
          [
              'name' => 'Icon tab 192',
              'type' => self::TYPE_IMAGE
          ]
      ],
      'MENU' => [
          [
              'name' => 'Home',
              'type' => self::TYPE_TEXT
          ],
//          [
//              'name' => 'Gallery',
//              'type' => self::TYPE_TEXT
//          ],
          [
              'name' => 'About Me',
              'type' => self::TYPE_TEXT
          ]
      ],

      'HOME_BANNER' => [
          [
              'name' => 'Banner',
              'type' => self::TYPE_IMAGE
          ],
          [
              'name' => 'Banner 2',
              'type' => self::TYPE_IMAGE
          ],
          [
              'name' => 'Banner 3',
              'type' => self::TYPE_IMAGE
          ],
          [
              'name' => 'Banner 4',
              'type' => self::TYPE_IMAGE
          ],
          [
              'name' => 'Text',
              'type' => self::TYPE_TEXT
          ],
          [
              'name' => 'Sub Text',
              'type' => self::TYPE_TEXT
          ],
      ],
      'ABOUT_ME' => [
          [
              'name' => 'Intro Footer',
              'type' => self::TYPE_TEXT
          ],
          [
              'name' => 'Address',
              'type' => self::TYPE_TEXT
          ],
          [
              'name' => 'Phone',
              'type' => self::TYPE_TEXT
          ],
          [
              'name' => 'Email',
              'type' => self::TYPE_TEXT
          ],
          [
              'name' => 'Avatar',
              'type' => self::TYPE_TEXT
          ],
          [
              'name' => 'About page title',
              'type' => self::TYPE_TEXT
          ],
          [
              'name' => 'About page content',
              'type' => self::TYPE_TEXT
          ],
          [
              'name' => 'Facebook',
              'type' => self::TYPE_TEXT
          ],
          [
              'name' => 'Tiktok',
              'type' => self::TYPE_TEXT
          ],
          [
              'name' => 'Instagram',
              'type' => self::TYPE_TEXT
          ],
      ],

  ];
}
