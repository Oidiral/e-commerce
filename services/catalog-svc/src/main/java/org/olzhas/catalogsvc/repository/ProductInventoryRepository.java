package org.olzhas.catalogsvc.repository;

import org.olzhas.catalogsvc.model.ProductInventory;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface ProductInventoryRepository extends JpaRepository<ProductInventory, UUID> {
}