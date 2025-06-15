package org.olzhas.catalogsvc.service.impl;

import io.minio.GetObjectArgs;
import io.minio.MinioClient;
import io.minio.PutObjectArgs;
import io.minio.RemoveObjectArgs;
import lombok.RequiredArgsConstructor;
import org.olzhas.catalogsvc.config.MinioProperties;
import org.olzhas.catalogsvc.dto.ProductImageDto;
import org.olzhas.catalogsvc.exceptionHandler.NotFoundException;
import org.olzhas.catalogsvc.mapper.ProductImageMapper;
import org.olzhas.catalogsvc.model.Product;
import org.olzhas.catalogsvc.model.ProductImage;
import org.olzhas.catalogsvc.repository.ProductImageRepository;
import org.olzhas.catalogsvc.repository.ProductRepository;
import org.olzhas.catalogsvc.service.ImageService;
import org.springframework.core.io.InputStreamResource;
import org.springframework.core.io.Resource;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.multipart.MultipartFile;

import java.io.InputStream;
import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class ProductImageServiceImpl implements ImageService {

    private final ProductRepository productRepository;
    private final ProductImageRepository productImageRepository;
    private final ProductImageMapper productImageMapper;
    private final MinioClient minioClient;
    private final MinioProperties properties;

    @Override
    @Transactional
    public List<ProductImageDto> upload(UUID productId, MultipartFile file, boolean primary) {
        Product product = productRepository.findById(productId)
                .orElseThrow(() -> new NotFoundException("Product not found with id: " + productId));
        String objectName = UUID.randomUUID() + "-" + file.getOriginalFilename();
        try (InputStream is = file.getInputStream()) {
            minioClient.putObject(
                    PutObjectArgs.builder()
                            .bucket(properties.getBucket())
                            .object(objectName)
                            .stream(is, file.getSize(), -1)
                            .contentType(file.getContentType())
                            .build());
        } catch (Exception e) {
            throw new RuntimeException("Failed to upload image", e);
        }
        ProductImage image = new ProductImage();
        image.setProduct(product);
        image.setS3Key(objectName);
        image.setUrl(properties.getUrl() + "/" + properties.getBucket() + "/" + objectName);
        image.setIsPrimary(primary);
        ProductImage saved = productImageRepository.save(image);
        return List.of(productImageMapper.toDto(saved));
    }

    @Override
    @Transactional(readOnly = true)
    public Resource download(UUID imageId) {
        ProductImage image = productImageRepository.findById(imageId)
                .orElseThrow(() -> new NotFoundException("Image not found with id: " + imageId));
        try {
            InputStream stream = minioClient.getObject(
                    GetObjectArgs.builder()
                            .bucket(properties.getBucket())
                            .object(image.getS3Key())
                            .build());
            return new InputStreamResource(stream);
        } catch (Exception e) {
            throw new RuntimeException("Failed to download image", e);
        }
    }

    @Override
    @Transactional
    public void delete(UUID imageId) {
        ProductImage image = productImageRepository.findById(imageId)
                .orElseThrow(() -> new NotFoundException("Image not found with id: " + imageId));
        try {
            minioClient.removeObject(
                    RemoveObjectArgs.builder()
                            .bucket(properties.getBucket())
                            .object(image.getS3Key())
                            .build()
            );
        } catch (Exception e) {
            throw new RuntimeException("Failed to delete file from S3", e);
        }
        productImageRepository.delete(image);
    }
}