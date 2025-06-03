package org.olzhas.catalogsvc.service.impl;

import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.dto.*;
import org.olzhas.catalogsvc.exceptionHandler.NotFoundException;
import org.olzhas.catalogsvc.mapper.ProductImageMapper;
import org.olzhas.catalogsvc.mapper.ProductMapper;
import org.olzhas.catalogsvc.model.Product;
import org.olzhas.catalogsvc.repository.ProductImageRepository;
import org.olzhas.catalogsvc.repository.ProductRepository;
import org.olzhas.catalogsvc.repository.spec.ProductSpecificationBuilder;
import org.olzhas.catalogsvc.service.ProductService;
import org.olzhas.catalogsvc.utils.PageConverter;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.domain.Specification;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

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
        return productRepository.findInternalInfo(id)
                .map(info -> new InternalProductDto(
                        info.getProductId(),
                        info.getName(),
                        info.getQuantity(),
                        info.getLatestPrice()))
                .orElseThrow(() -> new NotFoundException("Product not found with id: " + id));
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


}
