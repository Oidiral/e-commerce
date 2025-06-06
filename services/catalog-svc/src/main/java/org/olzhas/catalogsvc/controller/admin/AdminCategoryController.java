package org.olzhas.catalogsvc.controller.admin;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.dto.CategoryCreateReq;
import org.olzhas.catalogsvc.dto.CategoryDto;
import org.olzhas.catalogsvc.dto.CategoryUpdateReq;
import org.olzhas.catalogsvc.service.CategoryService;
import org.springframework.http.HttpStatus;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@RequestMapping("/admin/categories")
@RequiredArgsConstructor
@PreAuthorize("hasRole('ADMIN')")
public class AdminCategoryController {

    private final CategoryService categoryService;

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public CategoryDto create(@Valid @RequestBody CategoryCreateReq req) {
        return categoryService.create(req);
    }

    @PutMapping("/{id}")
    public CategoryDto rename(
            @PathVariable UUID id,
            @Valid @RequestBody CategoryUpdateReq req) {
        return categoryService.rename(id, req);
    }

    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void delete(@PathVariable UUID id) {
        categoryService.delete(id);
    }

    @DeleteMapping("/slug/{slug}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteSlug(@PathVariable String slug) {
        categoryService.deleteSlug(slug);
    }

}