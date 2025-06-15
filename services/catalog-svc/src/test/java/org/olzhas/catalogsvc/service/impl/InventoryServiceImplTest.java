package org.olzhas.catalogsvc.service.impl;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.olzhas.catalogsvc.model.ProductInventory;
import org.olzhas.catalogsvc.repository.ProductInventoryRepository;

import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class InventoryServiceImplTest {

    @Mock
    private ProductInventoryRepository productInventoryRepository;

    @InjectMocks
    private InventoryServiceImpl inventoryService;

    @Test
    void setQuantityUpdatesExistingInventory() {
        UUID productId = UUID.randomUUID();
        ProductInventory inventory = new ProductInventory();
        inventory.setId(productId);
        inventory.setQuantity(3);

        when(productInventoryRepository.findById(productId)).thenReturn(Optional.of(inventory));
        when(productInventoryRepository.save(inventory)).thenReturn(inventory);

        inventoryService.setQuantity(productId, 8);

        assertEquals(8, inventory.getQuantity());
        verify(productInventoryRepository).save(inventory);
    }

    @Test
    void setQuantityThrowsWhenMissing() {
        when(productInventoryRepository.findById(any())).thenReturn(Optional.empty());

        assertThrows(IllegalArgumentException.class,
                () -> inventoryService.setQuantity(UUID.randomUUID(), 2));
    }
}