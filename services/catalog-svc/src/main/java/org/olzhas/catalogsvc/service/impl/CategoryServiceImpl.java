package org.olzhas.catalogsvc.service.impl;

import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.dto.*;
import org.olzhas.catalogsvc.exceptionHandler.NotFoundException;
import org.olzhas.catalogsvc.mapper.CategoryMapper;
import org.olzhas.catalogsvc.mapper.ProductMapper;
import org.olzhas.catalogsvc.model.Category;
import org.olzhas.catalogsvc.model.Product;
import org.olzhas.catalogsvc.repository.CategoryRepository;
import org.olzhas.catalogsvc.repository.ProductRepository;
import org.olzhas.catalogsvc.service.CategoryService;
import org.olzhas.catalogsvc.utils.PageConverter;
import org.olzhas.catalogsvc.utils.SlugUtil;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class CategoryServiceImpl implements CategoryService {

    private final ProductRepository productRepository;
    private final ProductMapper productMapper;
    private final CategoryRepository categoryRepository;
    private final CategoryMapper categoryMapper;

    @Override
    public PageResponse<ProductDto> products(UUID id, Pageable p) {
        Page<Product> page = productRepository.findByCategoryId(id, p);
        if (page.isEmpty()) {
            throw new NotFoundException("No products found for category with id: " + id);
        }
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

    @Override
    @Transactional
    public CategoryDto create(CategoryCreateReq categoryDto) {
        Category category = categoryMapper.toEntity(categoryDto);
        String baseSlug = SlugUtil.toSlug(categoryDto.getName());
        String candidate = baseSlug;
        int counter = 1;

        while (categoryRepository.existsBySlug(candidate)) {
            candidate = baseSlug + "-" + counter++;
        }

        category.setSlug(candidate);
        return categoryMapper.toDto(categoryRepository.save(category));
    }

    @Override
    @Transactional
    public CategoryDto rename(UUID id, CategoryUpdateReq categoryDto) {
        Category category = categoryRepository.findById(id)
                .orElseThrow(() -> new NotFoundException("Category with id: " + id + " not found"));

        categoryMapper.partialUpdate(categoryDto, category);

        String baseSlug = SlugUtil.toSlug(category.getName());
        String candidate = baseSlug;
        int counter = 1;

        while (categoryRepository.existsBySlug(candidate) && !candidate.equals(category.getSlug())) {
            candidate = baseSlug + "-" + counter++;
        }

        category.setSlug(candidate);
        return categoryMapper.toDto(categoryRepository.save(category));
    }

    @Override
    public void delete(UUID id) {
        if (!categoryRepository.existsById(id)) {
            throw new NotFoundException("Category with id: " + id + " not found");
        }
        categoryRepository.deleteById(id);
    }

    @Override
    public void deleteSlug(String slug) {
        if (!categoryRepository.existsBySlug(slug)) {
            throw new NotFoundException("Category with slug: " + slug + " not found");
        }
        categoryRepository.deleteBySlug(slug);
    }
}
