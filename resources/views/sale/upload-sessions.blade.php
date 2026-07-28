@extends('sale.master')
@section('content')
<div class="main-content">
    <div class="page-content">
        <div class="container-fluid">

            <h3 class="mb-3">Upload Sessions</h3>

            @php
                $lastWeek = null;
            @endphp

            @forelse($sessions as $session)
                @php
                    $carbon = \Carbon\Carbon::parse($session['created_at']);
                    $week = $carbon->isoWeek() . '/' . $carbon->isoWeekYear();
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
    </div>
</div>

@section('script')
<script>
$(document).ready(function() {
    $(document).on('click', '.session-header', function() {
        const $card = $(this).closest('.session-card');
        const $products = $card.find('.session-products');

        if ($products.is(':visible')) {
            $products.slideUp();
            return;
        }

        if ($products.data('loaded')) {
            $products.slideDown();
            return;
        }

        const sessionId = $card.data('session-id');
        $products.html('<div class="p-2">Đang tải...</div>').slideDown();

        $.get('/sale/upload-sessions/' + sessionId + '/products', function(res) {
            const logs = res.logs || [];
            let html = '';

            logs.forEach(function(log) {
                const p = log.product;
                if (!p) {
                    html += '<div class="col-lg-2 col-md-4 col-sm-6 mt-2">' +
                        '<div class="text-muted small">' + (log.file_name || '') + ' — ' + log.status + '</div>' +
                        '</div>';
                    return;
                }

                const statusBadge = p.status === 'SOLD' ? 'badge-danger' : (p.status === 'ON_HOLD' ? 'badge-warning' : 'badge-info');

                html += '<div class="col-lg-2 col-md-4 col-sm-6 mt-2 text-center">' +
                    '<a href="/sale/product/' + p.id + '">' +
                        '<img src="' + (p.path_thumb || '') + '" style="width:100%">' +
                    '</a>' +
                    '<div>' + p.name + ' - ' + Number(p.price).toLocaleString() + 'đ</div>' +
                    '<span class="badge ' + statusBadge + '">' + p.status + '</span>' +
                    '</div>';
            });

            $products.html(html || '<div class="p-2 text-muted">Không có sản phẩm.</div>').data('loaded', true);
        });
    });
});
</script>
@endsection
