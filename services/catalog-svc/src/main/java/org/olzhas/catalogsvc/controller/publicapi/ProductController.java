package org.olzhas.catalogsvc.controller.publicapi;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.dto.PageResponse;
import org.olzhas.catalogsvc.dto.ProductDto;
import org.olzhas.catalogsvc.dto.ProductFilter;
import org.olzhas.catalogsvc.dto.ProductImageDto;
import org.olzhas.catalogsvc.service.ProductService;
import org.olzhas.catalogsvc.service.image.ProductImageService;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.multipart.MultipartFile;
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
    private final ProductImageService productImageService;

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

    @PostMapping("/{id}/images")
    public ProductImageDto upload(@PathVariable UUID id,
                                 @RequestParam("file") MultipartFile file,
                                 @RequestParam(value = "primary", defaultValue = "false") boolean primary) {
        return productImageService.upload(id, file, primary);
    }

    @GetMapping("/images/{imageId}")
    public ResponseEntity<Resource> download(@PathVariable UUID imageId) {
        Resource resource = productImageService.download(imageId);
        return ResponseEntity.ok()
                .header(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=" + imageId)
                .body(resource);
    }
}
