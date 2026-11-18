# Vintage Always True - Chrome Extension (Sales Copilot)

Extension nội bộ đọc tin nhắn và ảnh sản phẩm từ **Meta Business Suite Inbox** gửi về Server Go để phân tích gợi ý trả lời và lên đơn tự động khi có câu lệnh chốt đơn.

## Cách cài đặt vào Chrome (30 giây)

1. Mở Chrome, truy cập `chrome://extensions/`
2. Bật công tắc **Developer mode** (Chế độ dành cho nhà phát triển) ở góc trên bên phải.
3. Bấm **Load unpacked** (Tải tiện ích đã giải nén) và chọn thư mục:
   `/Applications/MAMP/htdocs/vintagealwaystrue-go/chrome-extension`
4. Ghim icon **Vintage Sales Copilot** lên thanh công cụ Chrome, bấm vào để kiểm tra cấu hình:
   - **Server URL:** `https://vintage-sale-go.onrender.com` (hoặc `http://localhost:8080` khi test local)
   - **Câu lệnh chốt đơn:** `OK ạ, để em lên đơn ngay`

## Cách hoạt động trên Meta Business Suite

- Mở trang `https://business.facebook.com/latest/inbox` và chọn bất kỳ cuộc hội thoại nào.
- Ở góc dưới bên phải sẽ xuất hiện badge `Copilot: Sẵn sàng` kèm nút **Sync ngay** (hoặc dùng phím tắt `Alt + C`).
- Extension tự động:
  1. Cắt ngữ cảnh hội thoại **sau lần chốt đơn gần nhất** (không gom nhầm đồ của tuần trước).
  2. Chuyển ảnh sản phẩm sang Base64 ngay trên trình duyệt.
  3. Gửi về `POST /api/copilot/ingest` trên Server Go.
  4. Nếu tin nhắn mới nhất của Shop là `OK ạ, để em lên đơn ngay`, Extension gắn cờ `action: "create_order"`.
