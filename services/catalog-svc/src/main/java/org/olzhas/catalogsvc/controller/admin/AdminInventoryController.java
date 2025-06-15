package org.olzhas.catalogsvc.controller.admin;

import jakarta.validation.constraints.PositiveOrZero;
import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.service.InventoryService;
import org.springframework.http.HttpStatus;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@RequestMapping("/admin/inventory")
@RequiredArgsConstructor
@PreAuthorize("hasRole('ADMIN')")
public class AdminInventoryController {

    private final InventoryService inventoryService;

    @PutMapping("/{productId}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void setQuantity(
            @PathVariable UUID productId,
            @RequestParam("qty") @PositiveOrZero int qty) {

        inventoryService.setQuantity(productId, qty);
    }
}