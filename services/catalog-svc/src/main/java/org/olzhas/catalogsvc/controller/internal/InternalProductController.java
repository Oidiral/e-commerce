package org.olzhas.catalogsvc.controller.internal;

import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.dto.InternalProductDto;
import org.olzhas.catalogsvc.service.ProductService;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.UUID;

@RestController
@RequiredArgsConstructor
@RequestMapping("/internal/products")
@PreAuthorize("hasAuthority('SERVICE_ORDER')")
public class InternalProductController {

    private final ProductService productService;

    @GetMapping("/{id}")
    public ResponseEntity<InternalProductDto> getProduct(@PathVariable UUID id){
        return ResponseEntity.ok(productService.getWithQuantity(id));
    }
}
