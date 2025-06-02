package org.olzhas.catalogsvc.controller.publicapi;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.dto.PageResponse;
import org.olzhas.catalogsvc.dto.ProductDto;
import org.olzhas.catalogsvc.dto.ProductFilter;
import org.olzhas.catalogsvc.dto.ProductImageDto;
import org.olzhas.catalogsvc.service.ProductService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.data.domain.Pageable;
import org.springframework.data.web.PageableDefault;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/products")
@RequiredArgsConstructor
public class ProductController {

    private final ProductService productService;

    @GetMapping
    public PageResponse<ProductDto> list(@PageableDefault(size = 25, sort = "id,desc") Pageable pageable) {
        return productService.findAll(pageable);
    }

    @GetMapping("/{id}")
    public ProductDto byId(@PathVariable UUID id) {
        return productService.getById(id);
    }

    @GetMapping("/search")
    public PageResponse<ProductDto> search(@ParameterObject @Valid ProductFilter filter,
                                           @PageableDefault(size = 25, sort = "id,desc") Pageable p) {
        return productService.search(filter, p);
    }

    @GetMapping("/{id}/images")
    public List<ProductImageDto> images(@PathVariable UUID id) {
        return productService.getImages(id);
    }
}
