package org.olzhas.catalogsvc.service;

import org.olzhas.catalogsvc.dto.PageResponse;
import org.olzhas.catalogsvc.dto.ProductDto;
import org.olzhas.catalogsvc.dto.ProductFilter;
import org.olzhas.catalogsvc.dto.ProductImageDto;
import org.springframework.data.domain.Pageable;

import java.util.List;
import java.util.UUID;

public interface ProductService {
    PageResponse<ProductDto> findAll(Pageable pageable);
    ProductDto getById(UUID id);
    PageResponse<ProductDto> search(ProductFilter filter, Pageable pageable);
    List<ProductImageDto> getImages(UUID productId);
}
