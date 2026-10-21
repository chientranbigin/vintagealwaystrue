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

## 5. QUY TRÌNH BUILD & TRIỂN KHAI (DEPLOY WORKFLOW)

Khi thực hiện sửa đổi mã nguồn trong các task tiếp theo, thực hiện theo các bước sau:

1. **Nếu sửa Vue Component trong Sale V2 (`resources/js/components/salev2/`)**:
   - Cập nhật file `.vue`.
   - Cập nhật bundle tương ứng trong `public/js/app.js` (hoặc build lại bằng Webpack/Mix).
   - Đồng bộ sang cả thư mục `/Applications/MAMP/htdocs/sale/public/js/app.js`.

2. **Kiểm tra biên dịch Go**:
   ```bash
   cd /Applications/MAMP/htdocs/vintagealwaystrue-go
   go build -o /dev/null ./cmd/server
   ```

3. **Commit và đẩy lên GitHub**:
   ```bash
   cd /Applications/MAMP/htdocs/vintagealwaystrue-go
   git add .
   git commit -m "Nội dung thay đổi"
   git push origin go-backend
   ```

4. **Kích hoạt Deploy trên Render**:
   ```bash
   curl -s -X POST https://api.render.com/v1/services/srv-dbbmmu7avr4c7384dcd0/deploys \
     -H "Authorization: Bearer rnd_R0LMUnydejKZUCl4HmUBIZJ9MVj8" \
     -H "Content-Type: application/json" \
     -d '{"clearCache": "do_not_clear"}'
   ```
   Sau 1 - 2 phút, kiểm tra trạng thái deploy:
   ```bash
   curl -s https://api.render.com/v1/services/srv-dbbmmu7avr4c7384dcd0/deploys?limit=1 \
     -H "Authorization: Bearer rnd_R0LMUnydejKZUCl4HmUBIZJ9MVj8"
   ```
   Khi trường `"status": "live"` là hệ thống đã cập nhật thành công lên Production.

---

## 6. HƯỚNG DẪN KHI MỞ CONVERSATION MỚI
Khi bạn mở phiên làm việc mới với Assistant, bạn chỉ cần gửi yêu cầu kèm nhắc nhở:
> *"Hãy đọc file `INFRASTRUCTURE_AND_SYSTEM_DOCS.md` trong repo để nắm toàn bộ thông tin hệ thống, hạ tầng, credentials và tiếp tục làm task..."*

Assistant mới sẽ đọc file này và có thể bắt tay vào làm việc ngay lập tức mà không cần hỏi lại các thông tin hạ tầng hay lịch sử phát triển trước đó.
