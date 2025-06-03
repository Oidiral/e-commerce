package org.olzhas.catalogsvc.controller.publicapi;

import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.dto.CategoryDto;
import org.olzhas.catalogsvc.dto.PageResponse;
import org.olzhas.catalogsvc.dto.ProductDto;
import org.olzhas.catalogsvc.model.Category;
import org.olzhas.catalogsvc.service.CategoryService;
import org.springframework.data.domain.Pageable;
import org.springframework.data.web.PageableDefault;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;
import java.util.UUID;

@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/categories")
public class CategoryController {

    private final CategoryService categoryService;

    @GetMapping
    public List<CategoryDto> all(){
        return categoryService.all();
    }

    @GetMapping("/{id}/products")
    public PageResponse<ProductDto> products(@PathVariable UUID id,@PageableDefault(size = 25, sort = "desc") Pageable p) {
        return categoryService.products(id, p);
    }
}
