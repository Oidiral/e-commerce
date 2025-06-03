package org.olzhas.catalogsvc.service.impl;

import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.dto.CategoryDto;
import org.olzhas.catalogsvc.dto.PageResponse;
import org.olzhas.catalogsvc.dto.ProductDto;
import org.olzhas.catalogsvc.exceptionHandler.NotFoundException;
import org.olzhas.catalogsvc.mapper.CategoryMapper;
import org.olzhas.catalogsvc.mapper.ProductMapper;
import org.olzhas.catalogsvc.model.Product;
import org.olzhas.catalogsvc.repository.CategoryRepository;
import org.olzhas.catalogsvc.repository.ProductRepository;
import org.olzhas.catalogsvc.service.CategoryService;
import org.olzhas.catalogsvc.utils.PageConverter;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class CategoryServiceImpl implements CategoryService {

    private final ProductRepository productRepository;
    private final ProductMapper productMapper;
    private final CategoryRepository categoryRepository;
    private final CategoryMapper categoryMapper;

    @Override
    public PageResponse<ProductDto> products(UUID id, Pageable p) {
        Page<Product> page = productRepository.findByCategoryId(id, p)
                .orElseThrow(() -> new NotFoundException("Category not found with id: " + id));
        List<ProductDto> dtos = page.map(productMapper::toDto).getContent();
        return PageConverter.toPageResponse(dtos, page);
    }

    @Override
    public List<CategoryDto> all() {
        return categoryRepository.findAll()
                .stream()
                .map(categoryMapper::toDto)
                .toList();
    }
}
