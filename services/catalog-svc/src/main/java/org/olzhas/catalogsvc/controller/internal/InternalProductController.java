package org.olzhas.catalogsvc.controller.internal;

import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.dto.InternalProductDto;
import org.olzhas.catalogsvc.service.ProductService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

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


    @PostMapping("/{id}/reserve")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void reserveProduct(@PathVariable UUID id,
                               @RequestParam("qty") int quantity) {
        productService.reserve(id, quantity);
    }

    @PostMapping("/{id}/release")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void releaseProduct(@PathVariable UUID id,
                               @RequestParam("qty") int quantity) {
        productService.release(id, quantity);
    }
}
