package org.olzhas.catalogsvc.service;

import org.olzhas.catalogsvc.dto.ProductImageDto;
import org.springframework.core.io.Resource;
import org.springframework.web.multipart.MultipartFile;

import java.util.List;
import java.util.UUID;

public interface ImageService {
    List<ProductImageDto> upload(UUID productId, MultipartFile file, boolean primary);

    Resource download(UUID imageId);

    void delete(UUID imageId);
}
