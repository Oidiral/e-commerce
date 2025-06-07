package org.olzhas.catalogsvc.service.impl;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.olzhas.catalogsvc.dto.CategoryCreateReq;
import org.olzhas.catalogsvc.dto.CategoryDto;
import org.olzhas.catalogsvc.mapper.CategoryMapper;
import org.olzhas.catalogsvc.mapper.ProductMapper;
import org.olzhas.catalogsvc.model.Category;
import org.olzhas.catalogsvc.repository.CategoryRepository;
import org.olzhas.catalogsvc.repository.ProductRepository;

import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CategoryServiceImplTest {

    @Mock
    private ProductRepository productRepository;
    @Mock
    private ProductMapper productMapper;
    @Mock
    private CategoryRepository categoryRepository;
    @Mock
    private CategoryMapper categoryMapper;

    @InjectMocks
    private CategoryServiceImpl categoryService;

    @Test
    void createGeneratesUniqueSlug() {
        CategoryCreateReq req = new CategoryCreateReq("Some Name");
        Category entity = new Category();
        entity.setName(req.getName());

        when(categoryMapper.toEntity(req)).thenReturn(entity);
        when(categoryRepository.existsBySlug("some-name")).thenReturn(true);
        when(categoryRepository.existsBySlug("some-name-1")).thenReturn(false);
        when(categoryRepository.save(any(Category.class))).thenAnswer(inv -> inv.getArgument(0));
        when(categoryMapper.toDto(any(Category.class))).thenAnswer(inv -> {
            Category c = inv.getArgument(0);
            return CategoryDto.builder().id(UUID.randomUUID()).name(c.getName()).slug(c.getSlug()).build();
        });

        CategoryDto dto = categoryService.create(req);

        assertEquals("some-name-1", dto.getSlug());
        verify(categoryRepository).existsBySlug("some-name");
        verify(categoryRepository).existsBySlug("some-name-1");
    }
}