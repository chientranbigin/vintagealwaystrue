# TÀI LIỆU HỆ THỐNG, HẠ TẦNG & HƯỚNG DẪN BÀN GIAO (SYSTEM & INFRASTRUCTURE HANDOFF DOCS)

> **Dự án**: Vintage Always True - Sale Management System  
> **Cập nhật lần cuối**: 21/10/2026  
> **Mục đích**: Tài liệu lưu trữ toàn diện về kiến trúc hệ thống, hạ tầng đám mây, các thông tin bảo mật/khóa bí mật (Secrets & Keys), thông tin SSH VPS cũ, các tính năng/lỗi đã xử lý và hướng dẫn tiếp tục phát triển cho các phiên làm việc tiếp theo.

---

## 1. TỔNG QUAN HỆ THỐNG & KIẾN TRÚC HIỆN TẠI

Hệ thống đã được chuyển dịch (migrate) và tối ưu từ kiến trúc PHP/Laravel nguyên khối (chạy trên cPanel VPS cũ) sang kiến trúc **Golang Microservice (Gin + GORM) kết hợp Vue 2**, chạy containerized trên đám mây.

- **Mã nguồn Golang chính**: `/Applications/MAMP/htdocs/vintagealwaystrue-go`
  - Git Repo: `https://github.com/chientranbigin/vintagealwaystrue.git`
  - Nhánh hoạt động: `go-backend`
- **Mã nguồn Laravel/PHP cũ & tài nguyên gốc**: `/Applications/MAMP/htdocs/sale`
- **Địa chỉ Production**: `https://vintage-sale-go.onrender.com`
- **Phân hệ ứng dụng**:
  - **Sale V1** (`/sale`): Giao diện Blade server-side rendering gốc từ Laravel được tái hiện hoàn toàn trên Go handler + template engine HTML thuần, tương thích toàn bộ route cũ.
  - **Sale V2** (`/salev2`): Giao diện Single Page Application (SPA) hiện đại viết bằng **Vue 2 + Element UI + TailwindCSS**, load qua `public/js/app.js` và điều hướng qua Vue Router (`/salev2/orders`, `/salev2/order/create`, `/salev2/order/:id/edit`, `/salev2/products`).

---

## 2. THÔNG TIN HẠ TẦNG ĐÁM MÂY & BÍ MẬT HỆ THỐNG (SECRETS & KEYS)

### 2.1. Render (Cloud Web Service)
- **Tên dịch vụ**: `vintage-sale-go`
- **Production URL**: `https://vintage-sale-go.onrender.com`
- **Service ID**: `srv-dbbmmu7avr4c7384dcd0`
- **Render API Key**: `rnd_R0LMUnydejKZUCl4HmUBIZJ9MVj8`
- **Lệnh kích hoạt Deploy tự động qua API**:
  ```bash
  curl -s -X POST https://api.render.com/v1/services/srv-dbbmmu7avr4c7384dcd0/deploys \
    -H "Authorization: Bearer rnd_R0LMUnydejKZUCl4HmUBIZJ9MVj8" \
    -H "Content-Type: application/json" \
    -d '{"clearCache": "do_not_clear"}'
  ```

### 2.2. Cơ sở dữ liệu Cloud (TiDB Cloud Serverless MySQL)
- **Engine**: MySQL 8.0 compatible (TiDB Serverless Cloud)
- **Host**: `gateway01.ap-southeast-1.prod.aws.tidbcloud.com`
- **Port**: `4000`
- **Database**: `sale_db`
- **Username**: `24jHtiebmCYdexH.root`
- **Password**: `SaleDbPassword@2026`
- **Full DSN Connection String (Go / GORM)**:
  ```text
  24jHtiebmCYdexH.root:SaleDbPassword@2026@tcp(gateway01.ap-southeast-1.prod.aws.tidbcloud.com:4000)/sale_db?tls=true&charset=utf8mb4&parseTime=True&loc=Asia%2FHo_Chi_Minh&time_zone=%27%2B07%3A00%27
  ```

### 2.3. Cron-Job.org (Cơ chế chống ngủ đông Render Free Tier)
Do Render free tier sẽ tự động tắt (spin down) khi không có request sau 15 phút, hệ thống đã cài đặt cron job tự động ping định kỳ:
- **Provider**: [cron-job.org](https://cron-job.org)
- **API Key**: `A8vkUYjuEX4jFS8N3AvOtrV6scMN6KKJv9NHQhW/P6M=`
- **Job ID**: `8690848`
- **Tên Job**: `Vintage Sale Go - Keep Alive`
- **Target URL**: `https://vintage-sale-go.onrender.com/healthz`
- **Tần suất**: Mỗi 10 phút một lần (`*/10 * * * *`, múi giờ `Asia/Ho_Chi_Minh`)
- **Trạng thái**: Đang kích hoạt (Active)

### 2.4. Lưu trữ Media & File ảnh (CDN)
- **CDN Base URL**: `https://res.cloudinary.com/leq5n17r/image/upload/`
- **Cloud Name**: `leq5n17r`
- **Cơ chế**: Route Go `/storage/*filepath` ưu tiên phục vụ file local nếu có; nếu không có, sẽ tự động chuyển hướng 301 sang CDN Cloudinary để ảnh luôn hiển thị sắc nét với tốc độ cao.

### 2.5. Tài khoản đăng nhập hệ thống (Auth)
- **Master Admin Password**: `vintage2026` (Ngoài ra chấp nhận các mật khẩu mặc định của các tài khoản có sẵn trong DB: `admin`, `123456`, `vintagealwaystrue`).
- **Auth Token Secret**: `vintage-always-true-secret-key-2026` (HMAC SHA-256).

---

## 3. THÔNG TIN & CÁCH SSH VÀO VPS CŨ (VINAHOST)

VPS cũ là nơi trước đây chạy hệ thống PHP/Laravel và MySQL cũ (`wxcowyfv_tadimedia`).

### 3.1. Thông số kết nối
- **IP / Host**: `123.30.136.246`
- **Host domain**: `vdc-whm-cheaphosting-1112.vinahost.org`
- **Port**: `22`
- **Username**: `wxcowyfv`
- **Đường dẫn SSH Private Key trên máy Mac**: `/Users/macbook/Downloads/id_rsa.1767709933` (hoặc profile `tadimedia` trong `~/.ssh/config`)
- **Passphrase của Private Key**: `G$qFTfx6NepT7p3c`
- **Thư mục dự án trên VPS**: `/home/wxcowyfv/public_html`

### 3.2. Thông tin Database cũ trên VPS
- **DB Connection**: MySQL localhost:3306
- **DB Name**: `wxcowyfv_tadimedia`
- **DB User**: `wxcowyfv_root`
- **DB Password**: `U869eK@ndDaqiEV`

### 3.3. Các câu lệnh SSH vào VPS cũ

#### Cách 1: Chạy trực tiếp từ Terminal
```bash
ssh -i /Users/macbook/Downloads/id_rsa.1767709933 wxcowyfv@123.30.136.246
# Khi terminal hỏi: "Enter passphrase for key '/Users/macbook/Downloads/id_rsa.1767709933':"
# Nhập: G$qFTfx6NepT7p3c
```

#### Cách 2: Dùng lệnh ssh với cấu hình alias (nếu có trong `~/.ssh/config`)
```bash
ssh tadimedia
# Nhập passphrase: G$qFTfx6NepT7p3c
```

#### Cách 3: Chạy script tự động (Expect script - không cần gõ mật khẩu bằng tay)
```bash
expect -c '
set timeout 30
spawn ssh -o StrictHostKeyChecking=no -i /Users/macbook/Downloads/id_rsa.1767709933 wxcowyfv@123.30.136.246
expect {
    "Enter passphrase for key*" {
        send "G\$qFTfx6NepT7p3c\r"
        exp_continue
    }
}
interact
'
```

#### Cách 4: Chạy lệnh nhanh không tương tác (One-off remote command)
```bash
expect -c '
set timeout 30
spawn ssh -o StrictHostKeyChecking=no -i /Users/macbook/Downloads/id_rsa.1767709933 wxcowyfv@123.30.136.246 "cd /home/wxcowyfv/public_html && ls -la"
expect {
    "Enter passphrase for key*" {
        send "G\$qFTfx6NepT7p3c\r"
        exp_continue
    }
    timeout { puts "Timed out"; exit 1 }
    eof
}
'
```

---

## 4. CHI TIẾT CÁC TÍNH NĂNG & LỖI ĐÃ XỬ LÝ GẦN ĐÂY

1. **Khắc phục lỗi nhân đôi sản phẩm khi sửa đơn hàng (Fix Duplicate Products on Edit Order)**:
   - *Nguyên nhân*: Do GORM trong Go khi thực hiện `tx.Save(&order)` với struct có nạp `Products` quan hệ nhiều-nhiều sẽ tự động `INSERT` thêm các dòng mới vào bảng trung gian `order_products`. Kèm theo frontend thiếu cơ chế khử trùng (`deduplicate`).
   - *Giải pháp*: Bỏ preload khi lưu đơn, dùng `tx.Omit("Products").Save(&order)`, xoá liên kết cũ `order_products` rồi chèn lại danh sách duy nhất. Phía Vue áp dụng `new Set()` và ép kiểu so sánh số `Number(p.id)`.
2. **Hỗ trợ tìm kiếm theo Mã sản phẩm trong Sale V2**:
   - Ô tìm kiếm đơn hàng V2 hiện tìm kiếm đa năng: theo Tên khách, Số điện thoại, Địa chỉ VÀ cả **Mã / Tên sản phẩm** có trong đơn hàng (giống V1).
3. **Giữ nguyên thông tin khi tìm kiếm trên V1 (`/sale/order-summary?search=...`)**:
   - Sửa lỗi khi truyền tham số search thì đơn hàng bị trống thông tin do thiếu format và load thiếu relation.
4. **Hiển thị & Sửa khoản Giảm giá (`freeship_amount`) ở Sale V2**:
   - Thêm hiển thị `Giảm giá: -[số tiền]` màu cam trên Card grid, Table view và Hóa đơn chụp ảnh.
   - Thêm ô nhập số **Discount (Giảm giá)** trong màn hình Edit đơn hàng (`OrderEdit.vue`), tự động tính trừ vào `Total Amount`.
   - Cập nhật handler Go `UpdateOrder` nhận và lưu trường `freeship_amount`.
5. **Định dạng ngày tạo đơn thân thiện (Relative time)**:
   - Hiển thị theo phong cách Moment.js: `1 day ago`, `2 hours ago`, `30 mins ago`, `just now` trên thẻ đơn hàng ở cả V1 và V2.
6. **Lược bỏ thông số kích thước rườm rà trên thẻ đơn hàng V2**:
   - Bỏ các tag size chi tiết trong list sản phẩm của card để giao diện gọn gàng, trực quan.
7. **Khắc phục vấn đề F5 chuyển về V1**:
   - Cấu hình route SPA fallback `/salev2/*filepath` trên Go server để khi người dùng đang ở V2 reload trang thì vẫn ở lại đúng màn hình V2.

---

## 5. QUY TRÌNH BUILD & TRIỂN KHAI KHI CÓ CODE MỚI (DEPLOYMENT GUIDE)

### 5.1. Cơ chế hoạt động của hệ thống CI/CD trên Render
- **Repository**: `https://github.com/chientranbigin/vintagealwaystrue.git`
- **Nhánh triển khai**: `go-backend`
- **Cơ chế đóng gói**: Render đọc file `Dockerfile` ở thư mục gốc:
  1. Sử dụng image `golang:1.22-alpine` để biên dịch binary Go từ `./cmd/server`.
  2. Sao chép các thư mục tài nguyên tĩnh (`public/`, `static/`) vào image runtime `alpine`.
  3. Khởi chạy binary trên cổng được Render chỉ định qua biến môi trường `$PORT` (mặc định 8080 hoặc 10000).

---

### 5.2. Các bước triển khai chi tiết từng bước (Step-by-Step)

#### Bước 1: Chỉnh sửa mã nguồn
- **Nếu sửa Backend Go**:
  - Các handler nằm trong: `handlers/` (`salev1.go`, `salev2.go`, `auth.go`, v.v.).
  - Các model GORM nằm trong: `models/`.
  - Cấu hình server nằm trong: `cmd/server/main.go` và `config/config.go`.
- **Nếu sửa Giao diện Sale V2 (Vue SPA)**:
  - Các component nằm trong: `resources/js/components/salev2/` (`OrderManager.vue`, `OrderEdit.vue`, `OrderCreate.vue`, v.v.).
  - File bundle Javascript phục vụ người dùng: `public/js/app.js`.

#### Bước 2: Đồng bộ & Kiểm tra bundle Javascript (Nếu có sửa Vue)
- Đảm bảo mã cập nhật trong file `.vue` đã được phản ánh vào `public/js/app.js`.
- Kiểm tra tính hợp lệ của cú pháp file bundle:
  ```bash
  node -c public/js/app.js
  ```
- Đồng bộ file bundle sang cả thư mục dự án Laravel gốc để đồng nhất:
  ```bash
  cp /Applications/MAMP/htdocs/vintagealwaystrue-go/public/js/app.js /Applications/MAMP/htdocs/sale/public/js/app.js
  ```

#### Bước 3: Kiểm tra biên dịch Go cục bộ trước khi push
Để đảm bảo không có lỗi cú pháp hoặc thiếu package làm gãy tiến trình build trên Render:
```bash
cd /Applications/MAMP/htdocs/vintagealwaystrue-go
go build -o /dev/null ./cmd/server
```
*(Nếu lệnh kết thúc không báo lỗi gì tức là code Go đã sẵn sàng deploy)*.

#### Bước 4: Commit và đẩy lên GitHub
```bash
cd /Applications/MAMP/htdocs/vintagealwaystrue-go
git add .
git commit -m "Mô tả tính năng hoặc fix bug mới"
git push origin go-backend
```

#### Bước 5: Kích hoạt triển khai (Trigger Deploy) trên Render
Bạn có thể kích hoạt ngay từ Terminal mà không cần mở trình duyệt:
```bash
curl -s -X POST https://api.render.com/v1/services/srv-dbbmmu7avr4c7384dcd0/deploys \
  -H "Authorization: Bearer rnd_R0LMUnydejKZUCl4HmUBIZJ9MVj8" \
  -H "Content-Type: application/json" \
  -d '{"clearCache": "do_not_clear"}'
```
*(Lưu ý: Nếu bạn có thay đổi cấu hình Dockerfile hoặc thư viện Go ngoài `go.mod`, đổi `"clearCache": "clear"` để Render tải lại sạch sẽ)*.

#### Bước 6: Theo dõi tiến độ Deploy
Kiểm tra trạng thái deploy mới nhất:
```bash
curl -s https://api.render.com/v1/services/srv-dbbmmu7avr4c7384dcd0/deploys?limit=1 \
  -H "Authorization: Bearer rnd_R0LMUnydejKZUCl4HmUBIZJ9MVj8" | grep -o '"status":"[^"]*"'
```
- Khi hiển thị `"status":"build_in_progress"`: Hệ thống đang build Docker image (thường mất khoảng 30 - 45 giây).
- Khi hiển thị `"status":"live"`: Quá trình deploy đã thành công, server mới đã chính thức nhận request!

#### Bước 7: Kiểm tra xác nhận trên Production
1. Kiểm tra endpoint health check:
   ```bash
   curl -s https://vintage-sale-go.onrender.com/healthz
   # Kết quả mong muốn: {"engine":"golang","status":"ok"}
   ```
2. Truy cập trình duyệt vào trang web để kiểm tra trực quan:
   - Sale V1: `https://vintage-sale-go.onrender.com/sale`
   - Sale V2: `https://vintage-sale-go.onrender.com/salev2/orders`
   - Nhấn `Cmd + Shift + R` (hoặc `Ctrl + F5`) trên trình duyệt để xóa cache JS/CSS cũ.

---

### 5.3. Xử lý sự cố thường gặp (Troubleshooting)

| Vấn đề | Nguyên nhân | Cách xử lý |
| :--- | :--- | :--- |
| **Giao diện V2 không đổi sau khi deploy** | Trình duyệt lưu cache file `app.js` cũ | Nhấn `Cmd + Shift + R` hoặc mở DevTools > tick "Disable cache" > reload lại. |
| **Status báo `build_failed`** | Lỗi code Go hoặc thiếu file trong Dockerfile | Chạy lệnh xem log deploy:<br>`curl -s https://api.render.com/v1/services/srv-dbbmmu7avr4c7384dcd0/deploys/<deploy-id> -H "Authorization: Bearer rnd_R0LMUnydejKZUCl4HmUBIZJ9MVj8"` |
| **Cần rollback về bản deploy trước đó** | Bản code mới phát sinh lỗi nghiêm trọng | Lấy ID bản deploy ổn định gần nhất và trigger deploy lại commit đó hoặc revert git commit và push lại. |

---

## 6. HƯỚNG DẪN KHI MỞ CONVERSATION MỚI (HANDOFF PROMPT)

Khi bạn mở phiên làm việc mới với AI Assistant, bạn chỉ cần gửi đoạn tin nhắn sau:

> *"Chào bạn, dự án quản lý đơn hàng Vintage Always True đã được chuyển dịch hoàn toàn sang Golang chạy trên Render. Hãy đọc file `INFRASTRUCTURE_AND_SYSTEM_DOCS.md` trong repo để nắm toàn bộ thông tin hệ thống, các khóa bí mật hạ tầng (Render, DB, Cron-job, Cloudinary), cách SSH vào VPS cũ và quy trình deploy code mới. Sau đó hỗ trợ tôi thực hiện task tiếp theo: [Nội dung task của bạn]..."*

AI Assistant mới sẽ tự động nạp toàn bộ kiến thức dự án từ file tài liệu này và tiếp tục phát triển chuẩn xác 100%!

