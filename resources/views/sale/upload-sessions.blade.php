@extends('sale.master')
@section('content')
<div class="main-content">
    <div class="page-content">
        <div class="container-fluid">

            <h3 class="mb-3">Upload Sessions</h3>

            <ul class="nav nav-tabs mb-3" id="uploadSessionsTabs">
                <li class="nav-item">
                    <a class="nav-link active" data-tab="sessions" href="#">Sessions</a>
                </li>
                <li class="nav-item">
                    <a class="nav-link" data-tab="latest" href="#">Latest Uploads</a>
                </li>
                <li class="nav-item">
                    <a class="nav-link" data-tab="sold" href="#">Sold by Date</a>
                </li>
            </ul>

            {{-- ===== Tab: Sessions ===== --}}
            <div id="tab-sessions" class="upload-tab-pane">
                @php $lastWeek = null; @endphp

                @forelse($sessions as $session)
                    @php
                        $carbon = \Carbon\Carbon::parse($session['created_at']);
                        $week = $carbon->format('W') . '/' . $carbon->format('o');
                    @endphp

                    @if($week !== $lastWeek)
                        <div class="week-separator" style="margin: 20px 0 10px; font-weight: bold; color: #888; border-bottom: 1px solid #ddd; padding-bottom: 4px;">
                            Tuần {{ $week }}
                        </div>
                        @php $lastWeek = $week; @endphp
                    @endif

                    <div class="card mb-2 session-card" data-session-id="{{ $session['session_id'] }}" style="cursor: pointer;">
                        <div class="card-body d-flex justify-content-between align-items-center session-header">
                            <div>
                                <strong>{{ $carbon->format('d/m/Y H:i') }}</strong>
                                <span class="text-muted ml-2">({{ $session['session_id'] }})</span>
                            </div>
                            <div>
                                <span class="badge badge-secondary">{{ $session['total_products'] }} sp</span>
                                <span class="badge badge-danger">{{ $session['sold_products'] }} sold</span>
                                @foreach($session['type_breakdown'] as $type => $count)
                                    <span class="badge badge-info">{{ $type }}: {{ $count }}</span>
                                @endforeach
                            </div>
                        </div>
                        <div class="session-products row px-3 pb-3" style="display:none;"></div>
                    </div>
                @empty
                    <p class="text-muted">Chưa có session upload nào.</p>
                @endforelse
            </div>

            {{-- ===== Tab: Latest Uploads ===== --}}
            <div id="tab-latest" class="upload-tab-pane" style="display:none;">
                <div class="card mb-3">
                    <div class="card-body d-flex flex-wrap align-items-center" style="gap: 15px;">
                        <div class="form-check">
                            <input type="checkbox" class="form-check-input" id="onlyAvailable" checked>
                            <label class="form-check-label font-weight-bold" for="onlyAvailable">Only Available</label>
                        </div>
                        <select id="latestTypeFilter" class="form-control" style="width: 180px;">
                            <option value="">All Types</option>
                            @foreach(\App\Contracts\SaleConfig::PRODUCT_TYPES as $type)
                                <option value="{{ $type }}">{{ $type }}</option>
                            @endforeach
                        </select>
                        <div class="flex-grow-1"></div>
                        <button id="bulkHoldBtn" class="btn btn-warning" disabled>
                            Set ON_HOLD <span id="latestSelectedCount"></span>
                        </button>
                    </div>
                </div>

                <div class="row" id="latestProductsGrid"></div>

                <div class="mt-3 d-flex justify-content-center" id="latestPagination"></div>
            </div>

            {{-- ===== Tab: Sold by Date ===== --}}
            <div id="tab-sold" class="upload-tab-pane" style="display:none;"></div>

        </div>
    </div>
</div>

<style>
.upload-product-card { position: relative; }
.upload-product-card img { width: 100%; height: 180px; object-fit: contain; background: #f8f9fa; }
.upload-product-card .badge { font-size: 10px; }
</style>

@section('script')
<script>
$(document).ready(function() {

    // ---- Tabs ----
    $('#uploadSessionsTabs a').on('click', function(e) {
        e.preventDefault();
        $('#uploadSessionsTabs a').removeClass('active');
        $(this).addClass('active');
        $('.upload-tab-pane').hide();
        const tab = $(this).data('tab');
        $('#tab-' + tab).show();

        if (tab === 'latest' && !$('#latestProductsGrid').data('loaded')) {
            fetchLatestProducts(1);
        }
        if (tab === 'sold' && !$('#tab-sold').data('loaded')) {
            fetchSoldDates();
        }
    });

    // ================= SESSIONS =================
    $(document).on('click', '.session-header', function() {
        const $card = $(this).closest('.session-card');
        const $products = $card.find('.session-products');

        if ($products.is(':visible')) { $products.slideUp(); return; }
        if ($products.data('loaded')) { $products.slideDown(); return; }

        const sessionId = $card.data('session-id');
        $products.html('<div class="p-2">Đang tải...</div>').slideDown();

        $.get('/sale/upload-sessions/' + sessionId + '/products', function(res) {
            const logs = res.logs || [];
            let html = '';
            logs.forEach(function(log) {
                if (!log.product) {
                    html += '<div class="col-6 col-sm-6 col-md-4 col-lg-2 mt-2"><div class="text-muted small">' + (log.file_name || '') + ' — ' + log.status + '</div></div>';
                    return;
                }
                html += renderProductCard(log.product);
            });
            $products.html(html || '<div class="p-2 text-muted">Không có sản phẩm.</div>').data('loaded', true);
        });
    });

    // ================= LATEST UPLOADS =================
    let latestSelected = [];

    $('#onlyAvailable, #latestTypeFilter').on('change', function() { fetchLatestProducts(1); });

    window.fetchLatestProducts = function(page) {
        const params = { page: page || 1, limit: 60 };
        if ($('#onlyAvailable').is(':checked')) params.status = 'AVAILABLE';
        const type = $('#latestTypeFilter').val();
        if (type) params.type = type;

        $('#latestProductsGrid').html('<div class="col-12 text-center py-5 text-muted">Đang tải...</div>').data('loaded', true);

        axios.get('/salev2/api/products-by-upload', { params }).then(function(res) {
            const products = res.data.data || [];
            let html = '';
            products.forEach(function(p) {
                html += renderProductCard(p, true);
            });
            $('#latestProductsGrid').html(html || '<div class="col-12 text-center py-5 text-muted">Không có sản phẩm.</div>');
            renderPagination(res.data.current_page, res.data.last_page);
        });
    };

    function renderPagination(current, last) {
        let html = '';
        for (let i = 1; i <= last; i++) {
            html += '<button class="btn btn-sm ' + (i === current ? 'btn-primary' : 'btn-outline-secondary') + ' mx-1 latest-page-btn" data-page="' + i + '">' + i + '</button>';
        }
        $('#latestPagination').html(html);
    }

    $(document).on('click', '.latest-page-btn', function() {
        fetchLatestProducts($(this).data('page'));
    });

    $(document).on('change', '.latest-select-checkbox', function() {
        const id = parseInt($(this).data('id'));
        if (this.checked) {
            if (!latestSelected.includes(id)) latestSelected.push(id);
        } else {
            latestSelected = latestSelected.filter(x => x !== id);
        }
        $('#bulkHoldBtn').prop('disabled', !latestSelected.length);
        $('#latestSelectedCount').text(latestSelected.length ? '(' + latestSelected.length + ')' : '');
    });

    $('#bulkHoldBtn').on('click', function() {
        if (!latestSelected.length) return;
        const $btn = $(this);
        $btn.prop('disabled', true).text('Đang cập nhật...');

        Promise.all(latestSelected.map(id => axios.post('/salev2/api/product/' + id + '/update', { status: 'ON_HOLD' })))
            .then(function() {
                latestSelected = [];
                $('#latestSelectedCount').text('');
                fetchLatestProducts(1);
            })
            .finally(function() {
                $btn.text('Set ON_HOLD');
            });
    });

    // ================= SOLD BY DATE =================
    window.fetchSoldDates = function() {
        $('#tab-sold').html('<div class="text-center py-5 text-muted">Đang tải...</div>').data('loaded', true);

        axios.get('/salev2/api/sold-by-date').then(function(res) {
            const dates = res.data.dates || [];
            let html = '';
            let lastWeek = null;

            dates.forEach(function(row) {
                const d = new Date(row.date + 'T00:00:00');
                const week = getISOWeek(d);
                if (week !== lastWeek) {
                    html += '<div class="week-separator" style="margin: 20px 0 10px; font-weight: bold; color: #888; border-bottom: 1px solid #ddd; padding-bottom: 4px;">Tuần ' + week + '</div>';
                    lastWeek = week;
                }

                const dateLabel = d.toLocaleDateString('vi-VN', { weekday: 'short', day: '2-digit', month: '2-digit', year: 'numeric' });

                if (row.is_empty) {
                    html += '<div class="card mb-2" style="background:#fff5f5;">' +
                        '<div class="card-body d-flex align-items-center">' +
                        '<span class="text-danger font-weight-bold">' + dateLabel + '</span>' +
                        '<span class="text-danger ml-2 small">— không bán được</span>' +
                        '</div></div>';
                } else {
                    let badges = '';
                    Object.keys(row.type_breakdown || {}).forEach(function(type) {
                        badges += '<span class="badge badge-info">' + type + ': ' + row.type_breakdown[type] + '</span> ';
                    });
                    badges += '<span class="badge badge-primary">' + row.total + ' sp</span> ';
                    badges += '<span class="badge badge-success">' + formatPrice(row.total_revenue) + '</span>';

                    html += '<div class="card mb-2 sold-date-card" data-date="' + row.date + '" style="cursor:pointer;">' +
                        '<div class="card-body d-flex justify-content-between align-items-center sold-date-header">' +
                        '<strong>' + dateLabel + '</strong>' +
                        '<div>' + badges + '</div>' +
                        '</div>' +
                        '<div class="sold-date-products row px-3 pb-3" style="display:none;"></div>' +
                        '</div>';
                }
            });

            $('#tab-sold').html(html || '<div class="text-center py-5 text-muted">Không có dữ liệu.</div>');
        });
    };

    $(document).on('click', '.sold-date-header', function() {
        const $card = $(this).closest('.sold-date-card');
        const $products = $card.find('.sold-date-products');

        if ($products.is(':visible')) { $products.slideUp(); return; }
        if ($products.data('loaded')) { $products.slideDown(); return; }

        const date = $card.data('date');
        $products.html('<div class="p-2">Đang tải...</div>').slideDown();

        axios.get('/salev2/api/sold-by-date/' + date + '/products').then(function(res) {
            const products = res.data.products || [];
            let html = '';
            products.forEach(function(p) { html += renderProductCard(p); });
            $products.html(html || '<div class="p-2 text-muted">Không có sản phẩm.</div>').data('loaded', true);
        });
    });

    // ================= SHARED HELPERS =================
    function formatPrice(val) {
        if (!val) return '0đ';
        return new Intl.NumberFormat('vi-VN').format(val) + 'đ';
    }

    function getISOWeek(d) {
        const date = new Date(d);
        date.setDate(date.getDate() + 3 - (date.getDay() + 6) % 7);
        const w1 = new Date(date.getFullYear(), 0, 4);
        const week = 1 + Math.round(((date - w1) / 86400000 - 3 + (w1.getDay() + 6) % 7) / 7);
        return week + '/' + date.getFullYear();
    }

    function statusBadgeClass(status) {
        return { AVAILABLE: 'badge-success', SOLD: 'badge-secondary', ON_HOLD: 'badge-warning' }[status] || 'badge-info';
    }

    function renderProductCard(p, withCheckbox) {
        const img = p.image_thumb_scale_url || p.path_thumb || '';
        const checkbox = withCheckbox
            ? '<input type="checkbox" class="latest-select-checkbox" data-id="' + p.id + '" style="position:absolute; top:6px; right:6px; transform: scale(1.3);">'
            : '';

        return '<div class="col-6 col-sm-6 col-md-4 col-lg-2 mt-2">' +
            '<div class="card upload-product-card">' +
            checkbox +
            '<a href="/sale/product/' + p.id + '">' +
            (img ? '<img src="' + img + '">' : '<div class="d-flex align-items-center justify-content-center" style="height:180px;background:#f8f9fa;">No image</div>') +
            '</a>' +
            '<div class="card-body text-center p-2">' +
            '<div class="font-weight-bold text-truncate">' + p.name + '</div>' +
            '<div class="text-primary font-weight-bold">' + formatPrice(p.price) + '</div>' +
            '<span class="badge ' + statusBadgeClass(p.status) + '">' + p.status + '</span>' +
            '</div>' +
            '</div>' +
            '</div>';
    }
});
</script>
@endsection
