package org.olzhas.catalogsvc.service.impl;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.dto.*;
import org.olzhas.catalogsvc.exceptionHandler.BadRequestException;
import org.olzhas.catalogsvc.exceptionHandler.NotFoundException;
import org.olzhas.catalogsvc.mapper.ProductImageMapper;
import org.olzhas.catalogsvc.mapper.ProductMapper;
import org.olzhas.catalogsvc.model.Product;
import org.olzhas.catalogsvc.model.ProductInventory;
import org.olzhas.catalogsvc.model.ProductPrice;
import org.olzhas.catalogsvc.repository.ProductImageRepository;
import org.olzhas.catalogsvc.repository.ProductInventoryRepository;
import org.olzhas.catalogsvc.repository.ProductPriceRepository;
import org.olzhas.catalogsvc.repository.ProductRepository;
import org.olzhas.catalogsvc.repository.spec.ProductSpecificationBuilder;
import org.olzhas.catalogsvc.service.ProductService;
import org.olzhas.catalogsvc.utils.PageConverter;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.domain.Specification;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;
import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class ProductServiceImpl implements ProductService {

    private final ProductRepository productRepository;
    private final ProductMapper productMapper;
    private final ProductImageRepository productImageRepository;
    private final ProductImageMapper productImageMapper;
    private final ProductInventoryRepository productInventoryRepository;
    private final ProductPriceRepository productPriceRepository;

    @Override
    public PageResponse<ProductDto> findAll(Pageable pageable) {
        Page<Product> page = productRepository.findAll(pageable);
        List<ProductDto> dtos = page.map(productMapper::toDto).getContent();
        return PageConverter.toPageResponse(dtos, page);
    }

    @Override
    public ProductDto getById(UUID id) {
        return productRepository.findById(id)
                .map(productMapper::toDto)
                .orElseThrow(() -> new NotFoundException("Product not found with id: " + id));
    }


    @Override
    public InternalProductDto getWithQuantity(UUID id) {
        return productRepository.findProductInfoBasic(id)
                .map(info -> new InternalProductDto(
                        info.getProductId(),
                        info.getName(),
                        info.getLatestPrice(),
                        info.getQuantity()))
                .orElseThrow(() -> new NotFoundException("Product not found with id: " + id));
    }

    @Override
    public ProductInventoryPriceDto getWithQuantityAndPrice(UUID id) {
        return productRepository.findProductInfoWithCurrency(id)
                .map(info -> new ProductInventoryPriceDto(
                                info.getProductId().toString(),
                                info.getQuantity(),
                                info.getCurrency(),
                                info.getLatestPrice()
                        )
                ).orElseThrow(() -> new NotFoundException("Product not found with id: " + id));
    }

    @Override
    public PageResponse<ProductDto> search(ProductFilter filter, Pageable pageable) {
        Specification<Product> spec = ProductSpecificationBuilder.build(filter);
        Page<Product> page = productRepository.findAll(spec, pageable);
        List<ProductDto> dtos = page.map(productMapper::toDto).getContent();
        return PageConverter.toPageResponse(dtos, page);
    }

    @Override
    public List<ProductImageDto> getImages(UUID productId) {
        return productImageRepository.findByProductId(productId)
                .stream()
                .map(productImageMapper::toDto)
                .toList();
    }

    @Override
    @Transactional
    public void reserve(UUID id, int qty) {
        productInventoryRepository.findById(id)
                .ifPresentOrElse(inv -> {
                    if (inv.getQuantity() < qty) {
                        throw new BadRequestException("Not enough stock");
                    }
                    inv.setQuantity(inv.getQuantity() - qty);
                }, () -> {
                    throw new NotFoundException("Product not found with id: " + id);
                });
    }

    @Override
    @Transactional
    public void release(UUID id, int qty) {
        productInventoryRepository.findById(id)
                .ifPresentOrElse(inv -> inv.setQuantity(inv.getQuantity() + qty),
                        () -> {
                            throw new NotFoundException("Product not found with id: " + id);
                        });
    }

    @Override
    @Transactional
    public ProductDto create(@Valid ProductCreateReq req) {
        Product product = productMapper.toEntity(req);
        if (req.getPrice() != null) {
            ProductPrice price = new ProductPrice();
            price.setProduct(product);
            price.setAmount(req.getPrice());
            productPriceRepository.save(price);
        }

        if (req.getQuantity() != null) {
            ProductInventory inventory = new ProductInventory();
            inventory.setId(product.getId());
            inventory.setProduct(product);
            inventory.setQuantity(req.getQuantity());
            productInventoryRepository.save(inventory);
        }
        productRepository.save(product);
        return productMapper.toDto(product);
    }

    @Override
    @Transactional
    public void delete(UUID id) {
        if (!productRepository.existsById(id)) {
            throw new NotFoundException("Product not found with id: " + id);
        }
        productRepository.deleteById(id);
    }

    @Override
    public ProductDto updateHard(UUID id, ProductUpdateReq req) {
        Product ex = productRepository.findById(id)
                .orElseThrow(() -> new NotFoundException("Product not found with id: " + id));

        Product newProduct = productMapper.toEntity(req);
        newProduct.setId(ex.getId());
        newProduct.setCreatedAt(ex.getCreatedAt());
        newProduct.setUpdatedAt(Instant.now());

        return productMapper.toDto(productRepository.save(newProduct));
    }

    @Override
    public ProductDto updateSoft(UUID id, ProductPatchReq request) {
        Product ex = productRepository.findById(id)
                .orElseThrow(() -> new NotFoundException("Product not found with id: " + id));

        Product updatedProduct = productMapper.partialUpdate(request, ex);
        updatedProduct.setUpdatedAt(Instant.now());

        return productMapper.toDto(productRepository.save(updatedProduct));
    }

}
