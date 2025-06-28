package org.olzhas.catalogsvc.service.impl;

import io.minio.*;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.ArgumentCaptor;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.olzhas.catalogsvc.config.MinioProperties;
import org.olzhas.catalogsvc.dto.ProductImageDto;
import org.olzhas.catalogsvc.exceptionHandler.NotFoundException;
import org.olzhas.catalogsvc.mapper.ProductImageMapper;
import org.olzhas.catalogsvc.model.Product;
import org.olzhas.catalogsvc.model.ProductImage;
import org.olzhas.catalogsvc.repository.ProductImageRepository;
import org.olzhas.catalogsvc.repository.ProductRepository;
import org.springframework.core.io.Resource;
import org.springframework.web.multipart.MultipartFile;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class ProductImageServiceImplTest {

    @Mock
    private ProductRepository productRepository;
    @Mock
    private ProductImageRepository productImageRepository;
    @Mock
    private ProductImageMapper productImageMapper;
    @Mock
    private MinioClient minioClient;
    @Mock
    private MultipartFile multipartFile;
    @Mock
    private MinioProperties properties;

    @InjectMocks
    private ProductImageServiceImpl productImageService;

    @Test
    void uploadStoresImageAndReturnsDto() throws Exception {
        UUID productId = UUID.randomUUID();
        Product product = new Product();
        product.setId(productId);

        when(productRepository.findById(productId)).thenReturn(Optional.of(product));
        when(multipartFile.getOriginalFilename()).thenReturn("image.jpg");
        when(multipartFile.getSize()).thenReturn(4L);
        when(multipartFile.getContentType()).thenReturn("image/jpeg");
        InputStream is = new ByteArrayInputStream("data".getBytes());
        when(multipartFile.getInputStream()).thenReturn(is);
        when(properties.getBucket()).thenReturn("bucket");
        when(properties.getUrl()).thenReturn("http://localhost");
        when(productImageRepository.save(any(ProductImage.class))).thenAnswer(inv -> inv.getArgument(0));
        ProductImageDto dto = new ProductImageDto(UUID.randomUUID(), "url", true);
        when(productImageMapper.toDto(any(ProductImage.class))).thenReturn(dto);

        ProductImageDto result = productImageService.upload(productId, multipartFile, true);

        assertSame(dto, result);
        verify(minioClient).putObject(any(PutObjectArgs.class));
        ArgumentCaptor<ProductImage> captor = ArgumentCaptor.forClass(ProductImage.class);
        verify(productImageRepository).save(captor.capture());
        ProductImage saved = captor.getValue();
        assertEquals(product, saved.getProduct());
        assertTrue(saved.getUrl().startsWith("http://localhost/bucket/"));
        assertTrue(saved.getIsPrimary());
    }

    @Test
    void uploadThrowsWhenProductMissing() {
        UUID productId = UUID.randomUUID();
        when(productRepository.findById(productId)).thenReturn(Optional.empty());

        assertThrows(NotFoundException.class,
                () -> productImageService.upload(productId, multipartFile, false));
    }

    @Test
    void downloadFetchesImage() throws Exception {
        UUID imageId = UUID.randomUUID();
        ProductImage image = new ProductImage();
        image.setId(imageId);
        image.setS3Key("key");
        when(productImageRepository.findById(imageId)).thenReturn(Optional.of(image));
        when(properties.getBucket()).thenReturn("bucket");
        GetObjectResponse response = mock(GetObjectResponse.class);
        when(response.readAllBytes()).thenReturn("data".getBytes());
        when(minioClient.getObject(any(GetObjectArgs.class))).thenReturn(response);

        Resource res = productImageService.download(imageId);

        assertNotNull(res);
        verify(minioClient).getObject(any(GetObjectArgs.class));
    }

    @Test
    void downloadThrowsWhenMissing() {
        when(productImageRepository.findById(any())).thenReturn(Optional.empty());

        assertThrows(NotFoundException.class, () -> productImageService.download(UUID.randomUUID()));
    }

    @Test
    void deleteRemovesImage() throws Exception {
        UUID imageId = UUID.randomUUID();
        ProductImage image = new ProductImage();
        image.setId(imageId);
        image.setS3Key("key");
        when(productImageRepository.findById(imageId)).thenReturn(Optional.of(image));
        when(properties.getBucket()).thenReturn("bucket");

        productImageService.delete(imageId);

        verify(minioClient).removeObject(any(RemoveObjectArgs.class));
        verify(productImageRepository).delete(image);
    }

    @Test
    void deleteThrowsWhenMissing() {
        when(productImageRepository.findById(any())).thenReturn(Optional.empty());

        assertThrows(NotFoundException.class, () -> productImageService.delete(UUID.randomUUID()));
    }
}