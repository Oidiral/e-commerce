package org.olzhas.catalogsvc.service.impl;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.olzhas.catalogsvc.exceptionHandler.BadRequestException;
import org.olzhas.catalogsvc.exceptionHandler.NotFoundException;
import org.olzhas.catalogsvc.mapper.ProductImageMapper;
import org.olzhas.catalogsvc.mapper.ProductMapper;
import org.olzhas.catalogsvc.model.ProductInventory;
import org.olzhas.catalogsvc.repository.ProductImageRepository;
import org.olzhas.catalogsvc.repository.ProductInventoryRepository;
import org.olzhas.catalogsvc.repository.ProductRepository;

import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class ProductServiceImplTest {

    @Mock
    private ProductRepository productRepository;
    @Mock
    private ProductMapper productMapper;
    @Mock
    private ProductImageRepository productImageRepository;
    @Mock
    private ProductImageMapper productImageMapper;
    @Mock
    private ProductInventoryRepository productInventoryRepository;

    @InjectMocks
    private ProductServiceImpl productService;

    @Test
    void reserveReducesQuantity() {
        UUID id = UUID.randomUUID();
        ProductInventory inv = new ProductInventory();
        inv.setId(id);
        inv.setQuantity(10);
        when(productInventoryRepository.findById(id)).thenReturn(Optional.of(inv));

        productService.reserve(id, 4);

        assertEquals(6, inv.getQuantity());
    }

    @Test
    void reserveThrowsBadRequestWhenNotEnoughStock() {
        UUID id = UUID.randomUUID();
        ProductInventory inv = new ProductInventory();
        inv.setId(id);
        inv.setQuantity(2);
        when(productInventoryRepository.findById(id)).thenReturn(Optional.of(inv));

        assertThrows(BadRequestException.class, () -> productService.reserve(id, 5));
    }

    @Test
    void reserveThrowsNotFoundWhenMissing() {
        UUID id = UUID.randomUUID();
        when(productInventoryRepository.findById(id)).thenReturn(Optional.empty());

        assertThrows(NotFoundException.class, () -> productService.reserve(id, 1));
    }

    @Test
    void releaseIncrementsQuantity() {
        UUID id = UUID.randomUUID();
        ProductInventory inv = new ProductInventory();
        inv.setId(id);
        inv.setQuantity(1);
        when(productInventoryRepository.findById(id)).thenReturn(Optional.of(inv));

        productService.release(id, 3);

        assertEquals(4, inv.getQuantity());
    }

    @Test
    void releaseThrowsNotFoundWhenMissing() {
        UUID id = UUID.randomUUID();
        when(productInventoryRepository.findById(id)).thenReturn(Optional.empty());

        assertThrows(NotFoundException.class, () -> productService.release(id, 1));
    }
}