package org.olzhas.catalogsvc.controller.admin;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.dto.PriceCreateReq;
import org.olzhas.catalogsvc.dto.ProductPriceResponseDto;
import org.olzhas.catalogsvc.service.PriceService;
import org.springframework.http.HttpStatus;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@RequestMapping("/admin/prices")
@RequiredArgsConstructor
@PreAuthorize("hasRole('ADMIN')")
public class AdminPriceController {

    private final PriceService priceService;

    @PostMapping("/{productId}")
    @ResponseStatus(HttpStatus.CREATED)
    public ProductPriceResponseDto add(@Valid @RequestBody PriceCreateReq req, @PathVariable UUID productId) {
        return priceService.addPrice(productId, req);
    }
}