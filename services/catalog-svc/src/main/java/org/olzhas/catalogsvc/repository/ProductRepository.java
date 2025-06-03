package org.olzhas.catalogsvc.repository;

import org.olzhas.catalogsvc.model.Product;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.Optional;
import java.util.UUID;

public interface ProductRepository extends JpaRepository<Product, UUID>,
        JpaSpecificationExecutor<Product> {

    @Query("SELECT p FROM ProductCategory pc JOIN pc.product p WHERE pc.category.id = :categoryId")
    Optional<Page<Product>> findByCategoryId(@Param("categoryId") UUID categoryId, Pageable pageable);
}