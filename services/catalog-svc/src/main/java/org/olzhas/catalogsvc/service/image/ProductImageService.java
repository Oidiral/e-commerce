package org.olzhas.catalogsvc.service.image;

import org.olzhas.catalogsvc.dto.ProductImageDto;
import org.springframework.core.io.Resource;
import org.springframework.web.multipart.MultipartFile;

import java.util.UUID;

public interface ProductImageService {
    ProductImageDto upload(UUID productId, MultipartFile file, boolean primary);
    Resource download(UUID imageId);
}
