package org.olzhas.catalogsvc.service;

import org.olzhas.catalogsvc.dto.*;
import org.springframework.data.domain.Pageable;
import org.springframework.web.bind.annotation.PathVariable;

import java.util.List;
import java.util.UUID;

public interface CategoryService {

    PageResponse<ProductDto> products(@PathVariable UUID id, Pageable p);
    List<CategoryDto> all();
    CategoryDto create(CategoryCreateReq categoryDto);
    CategoryDto rename(UUID id, CategoryUpdateReq categoryDto);
    void delete(UUID id);
    void deleteSlug(String slug);
}
