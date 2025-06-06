package org.olzhas.catalogsvc.service.impl;

import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.repository.ProductInventoryRepository;
import org.olzhas.catalogsvc.service.InventoryService;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
@RequiredArgsConstructor
public class InventoryServiceImpl implements InventoryService {

    private final ProductInventoryRepository productInventoryRepository;

    @Override
    public void setQuantity(UUID productId, int qty) {
        productInventoryRepository.findById(productId)
                .ifPresentOrElse(
                        inventory -> {
                            inventory.setQuantity(qty);
                            productInventoryRepository.save(inventory);
                        },
                        () -> {
                            throw new IllegalArgumentException("Product with ID " + productId + " not found");
                        }
                );
    }
}
