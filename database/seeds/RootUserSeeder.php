<?php

use Illuminate\Database\Seeder;

class RootUserSeeder extends Seeder
{
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        \App\User::updateOrCreate(
            ['email' => 'root@alwaystrue.com'],
            [
                'name' => 'Root Admin',
                'password' => bcrypt('root_password'),
            ]
        );
    }
}
