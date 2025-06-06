package org.olzhas.catalogsvc.service;

import org.olzhas.catalogsvc.dto.*;
import org.springframework.data.domain.Pageable;
import org.springframework.web.bind.annotation.PathVariable;

import java.util.List;
import java.util.UUID;

public interface ProductService {
    PageResponse<ProductDto> findAll(Pageable pageable);
    ProductDto getById(UUID id);
    InternalProductDto getWithQuantity(UUID id);
    PageResponse<ProductDto> search(ProductFilter filter, Pageable pageable);
    List<ProductImageDto> getImages(UUID productId);
    void reserve(UUID id, int quantity);
    void release(UUID id, int qty);
    ProductDto create(ProductCreateReq req);
    void delete(@PathVariable UUID id);
    ProductDto updateHard(UUID id, ProductUpdateReq req);
    ProductDto updateSoft(UUID id, ProductPatchReq request);
}
