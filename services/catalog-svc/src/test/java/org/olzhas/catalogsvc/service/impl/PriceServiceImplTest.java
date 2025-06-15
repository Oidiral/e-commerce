package org.olzhas.catalogsvc.service.impl;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.olzhas.catalogsvc.dto.PriceCreateReq;
import org.olzhas.catalogsvc.dto.ProductPriceResponseDto;
import org.olzhas.catalogsvc.model.Product;
import org.olzhas.catalogsvc.model.ProductPrice;
import org.olzhas.catalogsvc.repository.ProductPriceRepository;

import java.math.BigDecimal;
import java.time.Instant;
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class PriceServiceImplTest {

    @Mock
    private ProductPriceRepository productPriceRepository;

    @InjectMocks
    private PriceServiceImpl priceService;

    @Test
    void addPriceUpdatesExisting() {
        UUID productId = UUID.randomUUID();
        ProductPrice existing = new ProductPrice();
        existing.setId(UUID.randomUUID());
        existing.setProduct(Product.builder().id(productId).build());
        existing.setAmount(new BigDecimal("100.00"));
        existing.setCurrency("USD");
        existing.setCreatedAt(Instant.now());

        PriceCreateReq req = new PriceCreateReq();
        req.setAmount(new BigDecimal("120.00"));
        req.setCurrency("EUR");

        when(productPriceRepository.findByProductId(productId)).thenReturn(Optional.of(existing));
        when(productPriceRepository.save(existing)).thenReturn(existing);

        ProductPriceResponseDto dto = priceService.addPrice(productId, req);

        assertEquals(req.getAmount(), dto.getAmount());
        assertEquals(req.getCurrency(), dto.getCurrency());
        verify(productPriceRepository).save(existing);
    }

    @Test
    void addPriceCreatesNew() {
        UUID productId = UUID.randomUUID();
        PriceCreateReq req = new PriceCreateReq();
        req.setAmount(new BigDecimal("50.00"));
        req.setCurrency("USD");

        when(productPriceRepository.findByProductId(productId)).thenReturn(Optional.empty());
        when(productPriceRepository.save(any(ProductPrice.class))).thenAnswer(inv -> inv.getArgument(0));

        ProductPriceResponseDto dto = priceService.addPrice(productId, req);

        assertEquals(productId, dto.getProductId());
        assertEquals(req.getAmount(), dto.getAmount());
        assertEquals(req.getCurrency(), dto.getCurrency());
        verify(productPriceRepository).save(any(ProductPrice.class));
    }
}