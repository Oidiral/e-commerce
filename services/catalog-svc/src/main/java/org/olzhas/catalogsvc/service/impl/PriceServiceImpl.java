package org.olzhas.catalogsvc.service.impl;

import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.dto.PriceCreateReq;
import org.olzhas.catalogsvc.dto.ProductPriceResponseDto;
import org.olzhas.catalogsvc.model.Product;
import org.olzhas.catalogsvc.model.ProductPrice;
import org.olzhas.catalogsvc.repository.ProductPriceRepository;
import org.olzhas.catalogsvc.service.PriceService;
import org.springframework.stereotype.Service;

import java.util.Optional;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class PriceServiceImpl implements PriceService {

    private final ProductPriceRepository productPriceRepository;

    @Override
    public ProductPriceResponseDto addPrice(UUID productId, PriceCreateReq req) {
        Optional<ProductPrice> productPrice = productPriceRepository.findByProductId(productId);
        if (productPrice.isPresent()) {
            ProductPrice existingPrice = productPrice.get();
            existingPrice.setAmount(req.getAmount());
            existingPrice.setCurrency(req.getCurrency());
            productPriceRepository.save(existingPrice);
            return new ProductPriceResponseDto(
                    existingPrice.getId(),
                    existingPrice.getProduct().getId(),
                    existingPrice.getAmount(),
                    existingPrice.getCurrency(),
                    existingPrice.getCreatedAt()
            );
        } else {
            ProductPrice newProductPrice = new ProductPrice();
            newProductPrice.setProduct(Product.builder()
                    .id(productId)
                    .build());
            newProductPrice.setAmount(req.getAmount());
            newProductPrice.setCurrency(req.getCurrency());
            productPriceRepository.save(newProductPrice);
            return new ProductPriceResponseDto(
                    newProductPrice.getId(),
                    newProductPrice.getProduct().getId(),
                    newProductPrice.getAmount(),
                    newProductPrice.getCurrency(),
                    newProductPrice.getCreatedAt()
            );
        }


    }
}
