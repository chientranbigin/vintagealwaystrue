@extends('sale.master')
@section('content')
    <div class="main-content">

        <div class="page-content">

            <div class="container-fluid">


                <div class="card">
                    <div class="card-body">

                        <form id="productCreateForm" method="POST" action="{{route('sale.product.create.post')}}"
                              enctype="multipart/form-data">
                            @csrf
                            <div class="row">
                                <div class="col-md-6">
                                    <div class="form-group">
                                        <label for="productType">Type</label>
                                        <select class="form-control" id="productType" name="type">
                                            @foreach(App\Contracts\SaleConfig::PRODUCT_TYPES as $type)
                                                <option value="{{$type}}">{{$type}}</option>
                                            @endforeach
                                        </select>
                                    </div>

                                    <div class="form-group" id="quickInputGroup">
                                        <label>Quick Input (<span id="quickLabels" class="text-primary"></span>)</label>
                                        <input type="text" class="form-control" id="quickInput" placeholder="VD: 47 54 76 60">
                                        <small class="text-muted">Nhập số đo cách nhau bằng dấu cách</small>
                                    </div>


                                    <div class="form-group w-100">
                                        <label for="customFile">Thumb</label>
                                        <div class="custom-file">
                                            <input type="file" class="custom-file-input" id="customFile" name="banner">
                                            <label class="custom-file-label" for="customFile">Choose file</label>
                                        </div>
                                    </div>
                                    <div class="form-group w-100">
                                        <label>Price</label>
                                        <input type="number" class="form-control" name="price" id="productPrice">
                                    </div>

                                    <div class="form-group w-100">
                                        <label>Description</label>
                                        <textarea class="form-control" rows="5" name="description" id="productDesc"></textarea>
                                    </div>
                                </div>
                                <div class="col-md-6">
                                    <h3 class="mb-3">Sizes</h3>
                                    
                                    <!-- ÁO Section -->
                                    <div class="card mb-3" id="section-ao">
                                        <div class="card-header bg-light d-flex justify-content-between align-items-center" style="cursor: pointer" data-toggle="collapse" data-target="#collapseAo">
                                            <h5 class="mb-0">ÁO (Blazer, Jacket, Shirt...)</h5>
                                            <i class="mdi mdi-chevron-down"></i>
                                        </div>
                                        <div id="collapseAo" class="collapse show">
                                            <div class="card-body">
                                                <div class="row">
                                                    @foreach(App\Contracts\SaleConfig::SIZES as $key => $size)
                                                        @if($key < 6)
                                                            <div class="col-6 mb-2">
                                                                <label class="small mb-1">{{$size}}</label>
                                                                <input type="hidden" name="sizes[{{$key}}][name]" value="{{$size}}">
                                                                <input type="number" step="0.1" class="form-control form-control-sm size-input" data-index="{{$key}}" name="sizes[{{$key}}][value]">
                                                            </div>
                                                        @endif
                                                    @endforeach
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                    <!-- QUẦN Section -->
                                    <div class="card" id="section-quan">
                                        <div class="card-header bg-light d-flex justify-content-between align-items-center" style="cursor: pointer" data-toggle="collapse" data-target="#collapseQuan">
                                            <h5 class="mb-0">QUẦN (Trousers)</h5>
                                            <i class="mdi mdi-chevron-down"></i>
                                        </div>
                                        <div id="collapseQuan" class="collapse">
                                            <div class="card-body">
                                                <div class="row">
                                                    @foreach(App\Contracts\SaleConfig::SIZES as $key => $size)
                                                        @if($key >= 6)
                                                            <div class="col-6 mb-2">
                                                                <label class="small mb-1">{{$size}}</label>
                                                                <input type="hidden" name="sizes[{{$key}}][name]" value="{{$size}}">
                                                                <input type="number" step="0.1" class="form-control form-control-sm size-input" data-index="{{$key}}" name="sizes[{{$key}}][value]">
                                                            </div>
                                                        @endif
                                                    @endforeach
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div class="mt-4">
                                <button type="submit" class="btn btn-primary btn-lg px-5">Save Product</button>
                            </div>
                        </form>

                        <!-- Success Modal -->
                        <div class="modal fade" id="successModal" tabindex="-1" role="dialog" aria-hidden="true">
                            <div class="modal-dialog modal-dialog-centered" role="document">
                                <div class="modal-content border-0 shadow">
                                    <div class="modal-header bg-success text-white">
                                        <h5 class="modal-title">Sản phẩm đã tạo thành công!</h5>
                                        <button type="button" class="close text-white" data-dismiss="modal" aria-label="Close">
                                            <span aria-hidden="true">&times;</span>
                                        </button>
                                    </div>
                                    <div class="modal-body p-4 text-center">
                                        <div class="mb-3">
                                            <i class="mdi mdi-check-circle-outline text-success" style="font-size: 64px;"></i>
                                        </div>
                                        <h2 class="text-primary font-weight-bold" id="resProductCode">---</h2>
                                        <div class="mt-3 p-3 bg-light rounded text-left">
                                            <strong>Số đo:</strong>
                                            <p id="resProductSizes" class="mb-0 mt-1"></p>
                                        </div>
                                    </div>
                                    <div class="modal-footer">
                                        <button type="button" class="btn btn-secondary" data-dismiss="modal">Đóng (Tiếp tục tạo)</button>
                                    </div>
                                </div>
                            </div>
                        </div>

                    </div>
                </div>
            </div>
        </div>
    </div>
@endsection

@section('script')
<script>
$(document).ready(function() {
    const mappings = {
        'BLAZER': {
            labels: 'Vai - Ngực - Dài trước - Dài tay',
            indices: [0, 1, 3, 5]
        },
        'JACKET': {
            labels: 'Vai - Ngực - Dài - Dài tay',
            indices: [0, 1, 3, 5]
        },
        'TROUSERS': {
            labels: 'Eo - Đáy - Đùi - Dài quần - Ống - Dư lai',
            indices: [6, 7, 8, 9, 10, 11]
        }
    };

    function updateTypeLogic() {
        const type = $('#productType').val();
        const mapping = mappings[type] || { labels: 'Dữ liệu chưa hỗ trợ quick input', indices: [] };
        
        $('#quickLabels').text(mapping.labels);
        
        // Auto toggle sections
        if (type === 'BLAZER' || type === 'JACKET' || type === 'SHIRT' || type === 'GILE') {
            $('#collapseAo').collapse('show');
            $('#collapseQuan').collapse('hide');
        } else if (type === 'TROUSERS') {
            $('#collapseAo').collapse('hide');
            $('#collapseQuan').collapse('show');
        }
    }

    $('#productType').on('change', updateTypeLogic);
    updateTypeLogic();

    // Map quick input
    $('#quickInput').on('input', function() {
        const type = $('#productType').val();
        const mapping = mappings[type];
        if (!mapping || mapping.indices.length === 0) return;

        const values = $(this).val().trim().split(/\s+/);
        
        // Clear related inputs first
        // mapping.indices.forEach(idx => $(`.size-input[data-index="${idx}"]`).val(''));

        values.forEach((val, i) => {
            if (i < mapping.indices.length) {
                const idx = mapping.indices[i];
                $(`.size-input[data-index="${idx}"]`).val(val);
            }
        });
    });

    // AJAX Submit
    $('#productCreateForm').on('submit', function(e) {
        e.preventDefault();
        
        const formData = new FormData(this);
        const submitBtn = $(this).find('button[type="submit"]');
        submitBtn.prop('disabled', true).html('<span class="spinner-border spinner-border-sm mr-2"></span> Đang lưu...');

        $.ajax({
            url: $(this).attr('action'),
            method: 'POST',
            data: formData,
            processData: false,
            contentType: false,
            success: function(res) {
                $('#resProductCode').text(res.product_code);
                $('#resProductSizes').text(res.measurements);
                $('#successModal').modal('show');

                // Clear fields but keep type
                $('#quickInput').val('');
                $('#productPrice').val('');
                $('#productDesc').val('');
                $('#customFile').val('');
                $('.custom-file-label').text('Choose file');
                $('.size-input').val('');
                
                submitBtn.prop('disabled', false).text('Save Product');
            },
            error: function(xhr) {
                alert('Có lỗi xảy ra: ' + (xhr.responseJSON?.message || 'Vui lòng thử lại'));
                submitBtn.prop('disabled', false).text('Save Product');
            }
        });
    });

    // File input label update
    $('.custom-file-input').on('change', function() {
        var fileName = $(this).val().split('\\').pop();
        $(this).next('.custom-file-label').addClass("selected").html(fileName);
    });
});
</script>
<style>
    .card-header h5 { font-size: 15px; font-weight: 700; color: #555; }
    .size-input { border-radius: 6px; }
    .modal-content { border-radius: 15px; }
    #section-ao .card-header { border-left: 4px solid #007bff; }
    #section-quan .card-header { border-left: 4px solid #28a745; }
</style>

                    </div>
                    <!-- end card-body-->
                </div>


            </div> <!-- container-fluid -->
        </div>
        <!-- End Page-content -->


    </div>

@endsection
