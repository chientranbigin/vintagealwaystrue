# 🤖 HƯỚNG DẪN BẮT BUỘC DÀNH CHO AI AGENT (AGENT GUIDELINES)

> ⚠️ **QUY TẮC CỐT LÕI (CRITICAL RULE)**:  
> **MỖI KHI SỬA CODE, TÍNH NĂNG HOẶC GIAO DIỆN, AGENT BẮT BUỘC PHẢI THỰC HIỆN ĐẦY ĐỦ QUY TRÌNH BUILD & DEPLOY LÊN PRODUCTION (RENDER).**  
> Tuyệt đối không dừng lại ở việc sửa file local mà quên deploy, vì trang web bán hàng thực tế và toàn bộ logic đồng bộ Facebook, đơn hàng đang chạy trực tiếp trên máy chủ Render (`https://vintage-sale-go.onrender.com`). Nếu không deploy, các thao tác thực tế của người dùng trên web sẽ chạy code cũ và gây lỗi!

---

## 🚀 QUY TRÌNH DEPLOY BẮT BUỘC SAU KHI SỬA CODE (CHECKLIST)

Mỗi khi agent hoàn tất thay đổi mã nguồn (Go backend, template HTML hoặc Vue frontend), bắt buộc thực hiện theo đúng các bước sau:

### Bước 1: Build Frontend (nếu có sửa file trong `resources/js/...`)
```bash
npm run prod
node -c public/js/app.js
cp /Applications/MAMP/htdocs/vintagealwaystrue-go/public/js/app.js /Applications/MAMP/htdocs/sale/public/js/app.js
```

### Bước 2: Kiểm tra biên dịch Go và Unit Tests cục bộ
```bash
go build -o /dev/null ./cmd/server
go test ./handlers/...
go test ./services/...
```

### Bước 3: Commit và Push lên GitHub
```bash
git add .
git commit -m "feat/fix: mô tả ngắn gọn thay đổi"
git push origin go-backend
```

### Bước 4: Kích hoạt triển khai (Trigger Deploy) trên Render
```bash
curl -s -X POST https://api.render.com/v1/services/srv-dbbmmu7avr4c7384dcd0/deploys \
  -H "Authorization: Bearer rnd_R0LMUnydejKZUCl4HmUBIZJ9MVj8" \
  -H "Content-Type: application/json" \
  -d '{"clearCache": "do_not_clear"}'
```

### Bước 5: Kiểm tra xác nhận trạng thái Deploy
Kiểm tra cho tới khi Render trả về `"status":"live"`:
```bash
curl -s https://api.render.com/v1/services/srv-dbbmmu7avr4c7384dcd0/deploys?limit=1 \
  -H "Authorization: Bearer rnd_R0LMUnydejKZUCl4HmUBIZJ9MVj8" | grep -o '"status":"[^"]*"'
```

### Bước 6: Kiểm tra Production Health Check
```bash
curl -s https://vintage-sale-go.onrender.com/healthz
# Kỳ vọng: {"engine":"golang","status":"ok"}
```
Báo cáo hoàn thành cho người dùng chỉ sau khi Render đã deploy thành công (`live`).
