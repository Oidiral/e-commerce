package org.olzhas.catalogsvc.controller.admin;

import io.swagger.v3.oas.annotations.parameters.RequestBody;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.dto.*;
import org.olzhas.catalogsvc.service.ImageService;
import org.olzhas.catalogsvc.service.ProductService;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/admin/products")
@RequiredArgsConstructor
@PreAuthorize("hasRole('ADMIN')")
public class AdminProductController {
    private final ProductService productService;
    private final ImageService imageService;


    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public ProductDto create(@Valid @RequestBody ProductCreateReq req) {
        return productService.create(req);
    }

    @PutMapping("/{id}")
    public ProductDto update(
            @PathVariable UUID id,
            @Valid @RequestBody ProductUpdateReq req) {
        return productService.updateHard(id, req);
    }

    @PatchMapping("/{id}")
    public ProductDto patch(
            @PathVariable UUID id,
            @Valid @RequestBody ProductPatchReq req) {
        return productService.updateSoft(id, req);
    }

    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void delete(@PathVariable UUID id) {
        productService.delete(id);
    }

    @PostMapping(
            value = "/{id}/images",
            consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    public List<ProductImageDto> uploadImage(
            @PathVariable UUID id,
            @RequestPart MultipartFile file,
            @RequestParam(defaultValue = "false") boolean primary) {

        return imageService.upload(id, file, primary);
    }

    @GetMapping("/images/{imageId}")
    public ResponseEntity<Resource> download(@PathVariable UUID imageId) {
        Resource resource = imageService.download(imageId);
        return ResponseEntity.ok()
                .header(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=" + imageId)
                .body(resource);
    }

    @DeleteMapping("/images/{imageId}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteImage(@PathVariable UUID imageId) {
        imageService.delete(imageId);
    }
}
