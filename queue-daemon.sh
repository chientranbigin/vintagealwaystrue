#!/bin/bash
while true; do
    php artisan queue:work --tries=3
    sleep 2
done
