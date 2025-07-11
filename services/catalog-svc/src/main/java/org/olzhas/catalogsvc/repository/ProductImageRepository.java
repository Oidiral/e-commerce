package org.olzhas.catalogsvc.repository;

import org.olzhas.catalogsvc.model.ProductImage;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface ProductImageRepository extends JpaRepository<ProductImage, UUID> {
    @Query("SELECT pi FROM ProductImage pi WHERE pi.product.id = ?1")
    List<ProductImage> findByProductId(UUID productId);

}