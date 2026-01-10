@section('script')
    <script src="https://cdnjs.cloudflare.com/ajax/libs/html2canvas/1.4.1/html2canvas.min.js"></script>
    <script>
        $(document).ready(function() {
            console.log('Capture script ready');
            $(document).on('click', '.btn-capture', function() {
                const $btn = $(this);
                const orderData = $btn.data('order');
                const $receipt = $('#order-receipt-template');
                const originalHtml = $btn.html();
                
                // Immediate feedback
                $btn.html('<i class="mdi mdi-loading mdi-spin"></i> Processing...').prop('disabled', true);

                // Populate data (sync)
                $receipt.find('.c-name').text(orderData.name);
                $receipt.find('.c-phone').text(orderData.phone);
                $receipt.find('.c-address').text(orderData.address);
                $receipt.find('.c-note-text').text(orderData.note || 'Không có');
                
                if(!orderData.note) {
                    $receipt.find('.c-note-box').hide();
                } else {
                    $receipt.find('.c-note-box').show();
                }

                // Populate products
                const $prodContainer = $receipt.find('.c-products');
                $prodContainer.empty();
                orderData.products.forEach(p => {
                    $prodContainer.append(`
                        <div class="c-prod-item">
                            <img src="${p.path}" crossorigin="anonymous">
                            <div class="c-prod-code">${p.name}</div>
                        </div>
                    `);
                });

                // Format money
                const formatMoney = (val) => new Intl.NumberFormat('vi-VN').format(val) + 'đ';
                $receipt.find('.c-total-amount').text(formatMoney(orderData.total));
                $receipt.find('.c-final-amount').text(formatMoney(orderData.final));
                
                // Show template off-screen
                $receipt.css({
                    'display': 'block',
                    'position': 'fixed',
                    'top': '0',
                    'left': '0',
                    'opacity': '0',
                    'pointer-events': 'none',
                    'background': 'white',
                    'z-index': '9999'
                });

                // Small delay to ensure DOM updates and image refs are ready
                setTimeout(() => {
                    html2canvas($receipt[0], {
                        useCORS: true,
                        scale: 1.5, // Further reduced for speed/size
                        backgroundColor: '#ffffff',
                        logging: false
                    }).then(async canvas => {
                        const blob = await new Promise(resolve => canvas.toBlob(resolve, 'image/jpeg', 0.7)); // 0.7 quality
                        if (!blob) {
                            alert('Lỗi: Không thể tạo file ảnh.');
                            $receipt.hide();
                            $btn.html(originalHtml).prop('disabled', false);
                            return;
                        }

                        const fileName = `Order_${orderData.name}_${orderData.id}.jpg`;
                        const file = new File([blob], fileName, { type: 'image/jpeg' });

                        // Reset UI
                        $btn.html(originalHtml).prop('disabled', false);
                        $receipt.hide();

                        // 1. Native Share (iOS/Mobile)
                        if (navigator.share && navigator.canShare && navigator.canShare({ files: [file] })) {
                            try {
                                await navigator.share({
                                    files: [file],
                                    title: 'Đơn hàng ' + orderData.name,
                                });
                            } catch (err) {
                                if (err.name !== 'AbortError') {
                                    console.error('Share failed:', err);
                                    // Fallback to clipboard if share cancelled/failed
                                    copyToClipboard(blob, canvas, fileName);
                                }
                            }
                        } 
                        // 2. Clipboard/Download (Desktop/Fallback)
                        else {
                            copyToClipboard(blob, canvas, fileName);
                        }
                    }).catch(err => {
                        console.error('Capture failed:', err);
                        alert('Lỗi khi tạo ảnh: ' + err.message);
                        $receipt.hide();
                        $btn.html(originalHtml).prop('disabled', false);
                    });
                }, 100);
            });

            async function copyToClipboard(blob, canvas, fileName) {
                if (typeof ClipboardItem !== 'undefined') {
                    try {
                        const item = new ClipboardItem({ 'image/jpeg': blob });
                        await navigator.clipboard.write([item]);
                        alert('Đã copy vào Clipboard!');
                    } catch (err) {
                        downloadImage(canvas, fileName);
                        alert('Đã tải ảnh về máy (không hỗ trợ copy).');
                    }
                } else {
                    downloadImage(canvas, fileName);
                    alert('Đã tải ảnh về máy.');
                }
            }

            function downloadImage(canvas, fileName) {
                const link = document.createElement('a');
                link.download = fileName;
                link.href = canvas.toDataURL('image/jpeg', 0.7);
                link.click();
            }
        });
    </script>
