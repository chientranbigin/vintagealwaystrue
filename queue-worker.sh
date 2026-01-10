#!/bin/bash

# Script để chạy Laravel Queue Worker cho batch upload
# Sử dụng: ./queue-worker.sh

cd "$(dirname "$0")"

echo "========================================="
echo "Laravel Queue Worker - Batch Upload"
echo "========================================="
echo ""

# Kiểm tra PHP
if ! command -v php &> /dev/null; then
    echo "❌ PHP không được cài đặt hoặc không có trong PATH"
    exit 1
fi

echo "✓ PHP version: $(php -v | head -n 1)"
echo ""

# Kiểm tra database connection
echo "Checking database connection..."
php artisan migrate:status 2>/dev/null
if [ $? -ne 0 ]; then
    echo "⚠️  Warning: Không thể kết nối database"
    echo "   Vui lòng kiểm tra cấu hình .env và đảm bảo MySQL đang chạy"
    echo ""
    read -p "Bạn có muốn tiếp tục chạy queue worker? (y/n) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Tạo bảng jobs nếu chưa có
echo ""
echo "Checking jobs table..."
php artisan migrate --path=database/migrations/2026_01_06_135152_create_jobs_table.php --force 2>/dev/null

echo ""
echo "========================================="
echo "Starting Queue Worker..."
echo "========================================="
echo ""
echo "Queue Connection: $(grep QUEUE_CONNECTION .env | cut -d '=' -f2)"
echo "Tries: 2"
echo "Timeout: 120 seconds"
echo ""
echo "Press Ctrl+C to stop"
echo ""

# Chạy queue worker
php artisan queue:restart
php artisan queue:work --tries=2 --timeout=120 --verbose
