package org.olzhas.catalogsvc.service;

import org.olzhas.catalogsvc.dto.CategoryDto;
import org.olzhas.catalogsvc.dto.PageResponse;
import org.olzhas.catalogsvc.dto.ProductDto;
import org.springframework.data.domain.Pageable;
import org.springframework.web.bind.annotation.PathVariable;

import java.util.List;
import java.util.UUID;

public interface CategoryService {

    PageResponse<ProductDto> products(@PathVariable UUID id, Pageable p);
    List<CategoryDto> all();
}
