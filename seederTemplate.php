<?php

namespace Database\Seeders;

use Carbon\CarbonImmutable;
use Illuminate\Database\Seeder;

class {{CLASS_NAME}} extends Seeder
{

    /**
     * Auto generated seed file
     *
     * @return void
     */
    public function run()
    {
        \DB::table('{{TABLE_NAME}}')->insert(array(
{{SEEDER_RECORD}}
        ));
    }
}
